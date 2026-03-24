<script lang="ts">
  import { browser } from '$app/environment';
  import { user } from '$lib/auth';
  import {
    saveCWSettings,
    getUserInfo,
    updateCallSign,
    updateEmail,
    sendVerificationEmail,
    verifyEmail,
    updatePassword,
    savePageSettings,
    type PageSettings
  } from '$lib/api';
  import { langPreference, setLangPreference, type Locale } from '$lib/i18n.svelte';
  import { locales, localizeHref } from '$lib/paraglide/runtime';
  import {
    getLocaleLongLabel,
    normalizeLocalePreference,
    type LocalePreference
  } from '$lib/locale';
  import { LESSONS } from '$lib/morse';
  import {
    applyClientPageSettings,
    normalizeLesson,
    readClientCwSettings,
    readClientPageSettings,
    restoreSettingsFromServer,
    saveClientCwSettings
  } from '$lib/cwSync';
  import { localizeApiError } from '$lib/errorLocalization';
  import { CW_STORAGE_KEYS } from '$lib/storageKeys';
  import { Settings } from 'lucide-svelte';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import * as m from '$lib/paraglide/messages';

  type Theme = 'auto' | 'dark' | 'light';

  // Account section
  let username = $state('');
  let callSign = $state('');
  let initialCallSign = $state('');
  let callSignSaving = $state(false);
  let callSignError = $state('');
  let callSignSaved = $state(false);
  let email = $state('');
  let initialEmail = $state('');
  let emailVerified = $state(false);
  let emailSaving = $state(false);
  let emailError = $state('');
  let emailSaved = $state(false);
  let verificationCode = $state('');
  let verificationSendLoading = $state(false);
  let verificationSendError = $state('');
  let verificationSent = $state(false);
  let verificationCheckLoading = $state(false);
  let verificationCheckError = $state('');
  let verificationSuccess = $state(false);

  // Password section
  let currentPassword = $state('');
  let newPassword = $state('');
  let confirmPassword = $state('');
  let passwordSaving = $state(false);
  let passwordError = $state('');
  let passwordSaved = $state(false);

  // Page settings section
  let pageTheme = $state<Theme>('auto');
  let initialPageTheme = $state<Theme>('auto');
  let pageLanguage = $state<LocalePreference>(langPreference.value);
  let initialPageLanguage = $state<LocalePreference>(langPreference.value);
  let pageLesson = $state(1);
  let initialPageLesson = $state(1);
  let pageSaving = $state(false);
  let pageError = $state('');
  let pageSaved = $state(false);

  // CW Settings section
  let charWpm = $state(20);
  let initialCharWpm = $state(20);
  let effWpm = $state(12);
  let initialEffWpm = $state(12);
  let freq = $state(600);
  let initialFreq = $state(600);
  let startDelay = $state(0);
  let initialStartDelay = $state(0);
  let cwSaving = $state(false);
  let cwError = $state('');
  let cwSaved = $state(false);

  let loading = $state(true);
  let loadError = $state('');
  let lastAuthLoaded = $state<boolean | null>(null);

  const passwordDirty = $derived(
    currentPassword.trim() !== '' || newPassword.trim() !== '' || confirmPassword.trim() !== ''
  );
  const callSignDirty = $derived(callSign.trim().toUpperCase() !== initialCallSign);
  const emailDirty = $derived(email.trim() !== initialEmail.trim());
  const pageDirty = $derived(
    pageTheme !== initialPageTheme ||
      pageLanguage !== initialPageLanguage ||
      normalizeLesson(pageLesson, LESSONS.length) !== initialPageLesson
  );
  const cwDirty = $derived(
    charWpm !== initialCharWpm ||
      effWpm !== initialEffWpm ||
      freq !== initialFreq ||
      startDelay !== initialStartDelay
  );

  $effect(() => {
    const isAuthenticated = Boolean($user);
    if (lastAuthLoaded === isAuthenticated) return;

    lastAuthLoaded = isAuthenticated;
    loadAll();
  });

  function getStoredLesson(): number {
    if (!browser) return 1;

    const rawLesson = localStorage.getItem(CW_STORAGE_KEYS.lesson);
    const parsedLesson = Number.parseInt(rawLesson ?? '1', 10);
    return normalizeLesson(parsedLesson, LESSONS.length);
  }

  function applyCwState(cw: { char_wpm: number; eff_wpm: number; freq: number; start_delay: number }) {
    charWpm = cw.char_wpm;
    initialCharWpm = cw.char_wpm;
    effWpm = cw.eff_wpm;
    initialEffWpm = cw.eff_wpm;
    freq = cw.freq;
    initialFreq = cw.freq;
    startDelay = cw.start_delay;
    initialStartDelay = cw.start_delay;
  }

  function applyPageState(page: PageSettings) {
    pageTheme = page.theme;
    initialPageTheme = page.theme;
    pageLanguage = normalizeLocalePreference(page.language);
    initialPageLanguage = pageLanguage;
    pageLesson = normalizeLesson(page.cur_lesson, LESSONS.length);
    initialPageLesson = pageLesson;
  }

  function resetEmailVerificationState() {
    verificationCode = '';
    verificationSent = false;
    verificationSuccess = false;
    verificationSendError = '';
    verificationCheckError = '';
  }

  async function loadAll() {
    loading = true;
    loadError = '';

    if (!$user) {
      const localCw = readClientCwSettings();
      const localPage = readClientPageSettings(getStoredLesson(), LESSONS.length, langPreference.value);

      applyCwState(localCw);
      applyPageState(localPage);

      loading = false;
      return;
    }

    try {
      const [info, settings] = await Promise.all([getUserInfo(), restoreSettingsFromServer()]);
      const { cw, page } = settings;
      username = info.username;
      callSign = (info.call_sign ?? '').trim().toUpperCase();
      initialCallSign = callSign;
      email = info.email;
      initialEmail = info.email;
      emailVerified = info.email_verified;
      resetEmailVerificationState();
      applyCwState(cw);
      applyPageState(page);
    } catch (e) {
      loadError = localizeApiError(e, () => m.settings_load_error());
    } finally {
      loading = false;
    }
  }

  async function saveCallSign(e: SubmitEvent) {
    e.preventDefault();
    if (!callSignDirty) return;
    callSignSaving = true;
    callSignError = '';
    callSignSaved = false;
    try {
      await updateCallSign(callSign.trim().toUpperCase());
      callSign = callSign.trim().toUpperCase();
      initialCallSign = callSign;
      callSignSaved = true;
      setTimeout(() => (callSignSaved = false), 3000);
    } catch (err) {
      callSignError = localizeApiError(err, () => m.settings_save_error());
    } finally {
      callSignSaving = false;
    }
  }

  function languageLabel(locale: Locale): string {
    return getLocaleLongLabel(locale);
  }

  async function saveEmail(e: SubmitEvent) {
    e.preventDefault();
    if (!emailDirty) return;
    emailSaving = true;
    emailError = '';
    emailSaved = false;
    try {
      await updateEmail(email);
      initialEmail = email;
      emailVerified = false;
      resetEmailVerificationState();
      emailSaved = true;
      setTimeout(() => (emailSaved = false), 3000);
    } catch (err) {
      emailError = localizeApiError(err, () => m.settings_save_error());
    } finally {
      emailSaving = false;
    }
  }

  async function requestEmailVerificationCode() {
    if (emailDirty) {
      verificationSendError = m.settings_email_verify_save_email_first();
      return;
    }
    if (emailVerified) {
      verificationSendError = '';
      return;
    }

    verificationSendLoading = true;
    verificationSendError = '';
    verificationSent = false;
    verificationSuccess = false;

    try {
      await sendVerificationEmail();
      verificationSent = true;
      setTimeout(() => (verificationSent = false), 8000);
    } catch (err) {
      verificationSendError = localizeApiError(err, () => m.settings_save_error());
    } finally {
      verificationSendLoading = false;
    }
  }

  async function submitEmailVerification(e: SubmitEvent) {
    e.preventDefault();

    if (emailDirty) {
      verificationCheckError = m.settings_email_verify_save_email_first();
      return;
    }

    const code = verificationCode.trim();
    if (code === '') {
      verificationCheckError = m.settings_email_verify_code_required();
      return;
    }

    verificationCheckLoading = true;
    verificationCheckError = '';
    verificationSuccess = false;

    try {
      await verifyEmail(code);
      emailVerified = true;
      verificationCode = '';
      verificationSent = false;
      verificationSuccess = true;
      setTimeout(() => (verificationSuccess = false), 5000);
    } catch (err) {
      verificationCheckError = localizeApiError(err, () => m.settings_save_error());
    } finally {
      verificationCheckLoading = false;
    }
  }

  async function saveCW(e: SubmitEvent) {
    e.preventDefault();
    if (!cwDirty) return;
    cwSaving = true;
    cwError = '';
    cwSaved = false;
    try {
      const normalized = saveClientCwSettings({
        char_wpm: charWpm,
        eff_wpm: effWpm,
        freq,
        start_delay: startDelay
      });

      charWpm = normalized.char_wpm;
      effWpm = normalized.eff_wpm;
      freq = normalized.freq;
      startDelay = normalized.start_delay;

      if ($user) {
        await saveCWSettings(normalized);
      }

      initialCharWpm = normalized.char_wpm;
      initialEffWpm = normalized.eff_wpm;
      initialFreq = normalized.freq;
      initialStartDelay = normalized.start_delay;
      cwSaved = true;
      setTimeout(() => (cwSaved = false), 3000);
    } catch (err) {
      cwError = localizeApiError(err, () => m.settings_save_error());
    } finally {
      cwSaving = false;
    }
  }

  async function savePassword(e: SubmitEvent) {
    e.preventDefault();
    if (!passwordDirty) return;
    passwordSaving = true;
    passwordError = '';
    passwordSaved = false;

    if (newPassword.length < 8) {
      passwordError = m.settings_password_min_length();
      passwordSaving = false;
      return;
    }

    if (newPassword !== confirmPassword) {
      passwordError = m.register_error_mismatch();
      passwordSaving = false;
      return;
    }

    try {
      await updatePassword(currentPassword, newPassword);
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
      passwordSaved = true;
      setTimeout(() => (passwordSaved = false), 3000);
    } catch (err) {
      passwordError = localizeApiError(err, () => m.settings_save_error());
    } finally {
      passwordSaving = false;
    }
  }

  async function savePage(e: SubmitEvent) {
    e.preventDefault();
    if (!pageDirty) return;
    pageSaving = true;
    pageError = '';
    pageSaved = false;
    pageLesson = normalizeLesson(pageLesson, LESSONS.length);

    try {
      const pagePayload: PageSettings = {
        theme: pageTheme,
        language: pageLanguage,
        cur_lesson: pageLesson
      };

      if ($user) {
        await savePageSettings(pagePayload);
      }

      applyClientPageSettings(pagePayload, LESSONS.length, setLangPreference);
      initialPageTheme = pageTheme;
      initialPageLanguage = pageLanguage;
      initialPageLesson = pageLesson;
      pageSaved = true;
      setTimeout(() => (pageSaved = false), 3000);
    } catch (err) {
      pageError = localizeApiError(err, () => m.settings_save_error());
    } finally {
      pageSaving = false;
    }
  }
</script>

<main class="page-narrow settings-page">
  <div class="settings-heading">
    <Settings class="settings-page-icon" aria-hidden="true" />
    <h1 class="page-title">{m.settings_title()}</h1>
  </div>

  {#if loading}
    <LoadingSpinner variant="spinner" padded />
  {:else if loadError}
    <ErrorAlert message={loadError} />
  {:else}
    {#if !$user}
      <section class="card settings-card">
        <p class="body-text">
          {m.trainer_guest_notice()}
          <a href={localizeHref('/login')} class="link">{m.nav_login()}</a>
          /
          <a href={localizeHref('/register')} class="link">{m.nav_register()}</a>
        </p>
      </section>
    {/if}

    {#if $user}
      <!-- Account -->
      <section class="card settings-card">
        <h2 class="card-label">{m.settings_account_section()}</h2>
        <form onsubmit={saveCallSign} class="settings-form">
          <label class="settings-field">
            <span class="label-text">{m.settings_username_label()}</span>
            <input type="text" value={username} class="input" disabled />
          </label>
          <label class="settings-field">
            <span class="label-text">{m.settings_call_sign_label()}</span>
            <div class="settings-input-action">
              <input
                type="text"
                bind:value={callSign}
                class="input"
                placeholder={m.settings_call_sign_placeholder()}
                maxlength="32"
              />
              {#if callSignDirty || callSignSaving || callSignSaved}
                <button type="submit" class="btn-primary settings-btn-compact" disabled={callSignSaving}>
                  {callSignSaved
                    ? m.settings_saved()
                    : callSignSaving
                      ? m.settings_saving()
                      : m.settings_update()}
                </button>
              {/if}
            </div>
          </label>
          {#if callSignError}
            <p class="settings-error">⚠ {callSignError}</p>
          {/if}
        </form>

        <hr class="settings-divider" />

        <form onsubmit={saveEmail} class="settings-form">
          <label class="settings-field">
            <span class="label-text">{m.settings_email_label()}</span>
            <div class="settings-input-action">
              <input type="email" bind:value={email} class="input" required />
              {#if emailDirty || emailSaving || emailSaved}
                <button type="submit" class="btn-primary settings-btn-compact" disabled={emailSaving}>
                  {emailSaved
                    ? m.settings_saved()
                    : emailSaving
                      ? m.settings_saving()
                      : m.settings_update()}
                </button>
              {/if}
            </div>
          </label>
          {#if emailError}
            <p class="settings-error">⚠ {emailError}</p>
          {/if}
        </form>

        <div class="settings-email-verification">
          <p class={`settings-email-status ${emailVerified ? 'is-verified' : 'is-unverified'}`}>
            {emailVerified
              ? m.settings_email_verify_status_verified()
              : m.settings_email_verify_status_unverified()}
          </p>

          {#if !emailVerified}
            <div class="settings-input-action">
              <button
                type="button"
                class="btn-primary settings-btn-compact"
                onclick={requestEmailVerificationCode}
                disabled={verificationSendLoading || emailDirty}
              >
                {verificationSendLoading
                  ? m.settings_saving()
                  : verificationSent
                    ? m.settings_email_verify_code_sent()
                    : m.settings_email_verify_send_code()}
              </button>
            </div>

            <form onsubmit={submitEmailVerification} class="settings-form settings-verification-form">
              <label class="settings-field">
                <span class="label-text">{m.settings_email_verify_code_label()}</span>
                <div class="settings-input-action">
                  <input
                    type="text"
                    bind:value={verificationCode}
                    class="input"
                    placeholder={m.settings_email_verify_code_placeholder()}
                    inputmode="numeric"
                    autocomplete="one-time-code"
                  />
                  <button
                    type="submit"
                    class="btn-primary settings-btn-compact"
                    disabled={verificationCheckLoading || emailDirty}
                  >
                    {verificationCheckLoading
                      ? m.settings_saving()
                      : verificationSuccess
                        ? m.settings_email_verify_verified()
                        : m.settings_email_verify_confirm()}
                  </button>
                </div>
              </label>
            </form>
          {/if}

          {#if verificationSendError}
            <p class="settings-error">⚠ {verificationSendError}</p>
          {/if}
          {#if verificationCheckError}
            <p class="settings-error">⚠ {verificationCheckError}</p>
          {/if}
        </div>

        <hr class="settings-divider" />

        <h3 class="settings-subtitle">{m.settings_password_section()}</h3>
        <form onsubmit={savePassword} class="settings-form">
          <label class="settings-field">
            <span class="label-text">{m.settings_current_password_label()}</span>
            <input
              type="password"
              bind:value={currentPassword}
              class="input"
              minlength="8"
              autocomplete="current-password"
              required
            />
          </label>
          <label class="settings-field">
            <span class="label-text">{m.settings_new_password_label()}</span>
            <input
              type="password"
              bind:value={newPassword}
              class="input"
              minlength="8"
              autocomplete="new-password"
              required
            />
          </label>
          <label class="settings-field">
            <span class="label-text">{m.settings_confirm_new_password_label()}</span>
            <input
              type="password"
              bind:value={confirmPassword}
              class="input"
              minlength="8"
              autocomplete="new-password"
              required
            />
          </label>
          {#if passwordError}
            <p class="settings-error">⚠ {passwordError}</p>
          {/if}
          {#if passwordDirty || passwordSaving || passwordSaved}
            <div class="settings-action-row">
              <button type="submit" class="btn-primary settings-btn-compact" disabled={passwordSaving}>
                {passwordSaved
                  ? m.settings_saved()
                  : passwordSaving
                    ? m.settings_saving()
                    : m.settings_update()}
              </button>
            </div>
          {/if}
        </form>
      </section>
    {/if}

    <!-- Page Settings -->
    <section class="card settings-card">
      <h2 class="card-label">{m.settings_page_section()}</h2>
      <form onsubmit={savePage} class="settings-form">
        <label class="settings-field">
          <span class="label-text">{m.settings_theme_label()}</span>
          <select bind:value={pageTheme} class="input">
            <option value="auto">{m.theme_auto()}</option>
            <option value="light">{m.theme_light()}</option>
            <option value="dark">{m.theme_dark()}</option>
          </select>
        </label>
        <label class="settings-field">
          <span class="label-text">{m.settings_language_label()}</span>
          <select bind:value={pageLanguage} class="input">
            <option value="auto">{m.theme_auto()}</option>
            {#each locales as locale (locale)}
              <option value={locale}>{languageLabel(locale)}</option>
            {/each}
          </select>
        </label>
        <label class="settings-field">
          <span class="label-text">{m.settings_current_lesson_label()}</span>
          <input type="number" bind:value={pageLesson} min="1" max={LESSONS.length} class="input" />
        </label>
        {#if pageError}
          <p class="settings-error">⚠ {pageError}</p>
        {/if}
        {#if pageDirty || pageSaving || pageSaved}
          <div class="settings-action-row">
            <button type="submit" class="btn-primary settings-btn-compact" disabled={pageSaving}>
              {pageSaved ? m.settings_saved() : pageSaving ? m.settings_saving() : m.settings_update()}
            </button>
          </div>
        {/if}
      </form>
    </section>

    <!-- CW Settings -->
    <section class="card settings-card">
      <h2 class="card-label">{m.settings_cw_section()}</h2>
      <form onsubmit={saveCW} class="settings-form">
        <label class="settings-field">
          <span class="label-text">{m.trainer_label_char_wpm()}</span>
          <input type="number" bind:value={charWpm} min="5" max="50" class="input" />
        </label>
        <label class="settings-field">
          <span class="label-text">{m.trainer_label_eff_wpm()}</span>
          <input type="number" bind:value={effWpm} min="5" max="50" class="input" />
        </label>
        <label class="settings-field">
          <span class="label-text">{m.trainer_label_freq()}</span>
          <input type="number" bind:value={freq} min="300" max="2000" class="input" />
        </label>
        <label class="settings-field">
          <span class="label-text">{m.trainer_label_start_delay()}</span>
          <input type="number" bind:value={startDelay} min="0" max="10" step="0.5" class="input" />
        </label>
        {#if cwError}
          <p class="settings-error">⚠ {cwError}</p>
        {/if}
        {#if cwDirty || cwSaving || cwSaved}
          <div class="settings-action-row">
            <button type="submit" class="btn-primary settings-btn-compact" disabled={cwSaving}>
              {cwSaved ? m.settings_saved() : cwSaving ? m.settings_saving() : m.settings_update()}
            </button>
          </div>
        {/if}
      </form>
    </section>
  {/if}
</main>

<style>
  .settings-page {
    padding-top: 2rem;
    padding-bottom: 2rem;
  }
  .settings-heading {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }
  :global(.settings-page-icon) {
    color: var(--accent);
    width: 2rem;
    height: 2rem;
    flex-shrink: 0;
  }
  .settings-card {
    margin-bottom: 1.25rem;
  }
  .settings-form {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 0.5rem;
  }
  .settings-error {
    color: var(--error, #ef4444);
    font-size: 0.875rem;
  }
  .settings-divider {
    border: none;
    border-top: 1px solid var(--border);
    margin: 1rem 0;
  }
  .settings-email-verification {
    margin-top: 0.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
  }
  .settings-email-status {
    margin: 0;
    font-size: 0.9rem;
    font-weight: 600;
  }
  .settings-email-status.is-verified {
    color: var(--success, #10b981);
  }
  .settings-email-status.is-unverified {
    color: var(--warning, #d97706);
  }
  .settings-verification-form {
    margin-top: 0;
  }
  .settings-subtitle {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
    color: var(--text-primary);
  }
  .settings-input-action {
    display: flex;
    gap: 0.6rem;
    align-items: stretch;
  }
  .settings-input-action .input {
    flex: 1;
    min-width: 0;
  }
  .settings-action-row {
    display: flex;
    justify-content: flex-end;
  }
  .settings-btn-compact {
    padding: var(--space-sm) 0.9rem;
    font-size: 0.875rem;
    line-height: 1.25rem;
    min-height: calc(1.25rem + (var(--space-sm) * 2) + 2px);
    width: auto;
    white-space: nowrap;
  }

  @media (max-width: 720px) {
    .settings-input-action {
      flex-direction: column;
      align-items: stretch;
    }
    .settings-btn-compact {
      width: 100%;
    }
  }
</style>

<script lang="ts">
  import { user } from '$lib/auth';
  import {
    saveCWSettings,
    getUserInfo,
    updateCallSign,
    updateEmail,
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
  import { applyClientPageSettings, normalizeLesson, restoreSettingsFromServer } from '$lib/cwSync';
  import { localizeApiError } from '$lib/errorLocalization';
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
  let emailSaving = $state(false);
  let emailError = $state('');
  let emailSaved = $state(false);

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
    if ($user) {
      loadAll();
    } else {
      loading = false;
    }
  });

  async function loadAll() {
    loading = true;
    loadError = '';
    try {
      const [info, settings] = await Promise.all([getUserInfo(), restoreSettingsFromServer()]);
      const { cw, page } = settings;
      username = info.username;
      callSign = (info.call_sign ?? '').trim().toUpperCase();
      initialCallSign = callSign;
      email = info.email;
      initialEmail = info.email;
      charWpm = cw.char_wpm;
      initialCharWpm = cw.char_wpm;
      effWpm = cw.eff_wpm;
      initialEffWpm = cw.eff_wpm;
      freq = cw.freq;
      initialFreq = cw.freq;
      startDelay = cw.start_delay;
      initialStartDelay = cw.start_delay;
      pageTheme = page.theme;
      initialPageTheme = page.theme;
      pageLanguage = normalizeLocalePreference(page.language);
      initialPageLanguage = pageLanguage;
      pageLesson = normalizeLesson(page.cur_lesson, LESSONS.length);
      initialPageLesson = pageLesson;
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
      emailSaved = true;
      setTimeout(() => (emailSaved = false), 3000);
    } catch (err) {
      emailError = localizeApiError(err, () => m.settings_save_error());
    } finally {
      emailSaving = false;
    }
  }

  async function saveCW(e: SubmitEvent) {
    e.preventDefault();
    if (!cwDirty) return;
    cwSaving = true;
    cwError = '';
    cwSaved = false;
    try {
      await saveCWSettings({ char_wpm: charWpm, eff_wpm: effWpm, freq, start_delay: startDelay });
      initialCharWpm = charWpm;
      initialEffWpm = effWpm;
      initialFreq = freq;
      initialStartDelay = startDelay;
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
      await savePageSettings(pagePayload);
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

  {#if !$user && !loading}
    <div class="card">
      <p class="body-text">
        {m.settings_not_logged_in()}
        <a href={localizeHref('/login')} class="link">{m.nav_login()}</a>
      </p>
    </div>
  {:else if loading}
    <LoadingSpinner variant="spinner" padded />
  {:else if loadError}
    <ErrorAlert message={loadError} />
  {:else}
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

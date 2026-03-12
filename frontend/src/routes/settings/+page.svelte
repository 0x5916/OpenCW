<script lang="ts">
  import { user } from '$lib/auth';
  import { getCWSettings, saveCWSettings, getUserInfo, updateEmail } from '$lib/api';
  import { Settings } from 'lucide-svelte';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import * as m from '$lib/paraglide/messages';

  // Account section
  let username = $state('');
  let email = $state('');
  let emailSaving = $state(false);
  let emailError = $state('');
  let emailSaved = $state(false);

  // CW Settings section
  let charWpm = $state(20);
  let effWpm = $state(12);
  let freq = $state(600);
  let startDelay = $state(0);
  let cwSaving = $state(false);
  let cwError = $state('');
  let cwSaved = $state(false);

  let loading = $state(true);
  let loadError = $state('');

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
      const [info, cw] = await Promise.all([getUserInfo(), getCWSettings()]);
      username = info.username;
      email = info.email;
      charWpm = cw.char_wpm;
      effWpm = cw.eff_wpm;
      freq = cw.freq;
      startDelay = cw.start_delay;
    } catch (e) {
      loadError = e instanceof Error ? e.message : m.settings_load_error();
    } finally {
      loading = false;
    }
  }

  async function saveEmail(e: SubmitEvent) {
    e.preventDefault();
    emailSaving = true;
    emailError = '';
    emailSaved = false;
    try {
      await updateEmail(email);
      emailSaved = true;
      setTimeout(() => (emailSaved = false), 3000);
    } catch (err) {
      emailError = err instanceof Error ? err.message : m.settings_save_error();
    } finally {
      emailSaving = false;
    }
  }

  async function saveCW(e: SubmitEvent) {
    e.preventDefault();
    cwSaving = true;
    cwError = '';
    cwSaved = false;
    try {
      await saveCWSettings({ char_wpm: charWpm, eff_wpm: effWpm, freq, start_delay: startDelay });
      cwSaved = true;
      setTimeout(() => (cwSaved = false), 3000);
    } catch (err) {
      cwError = err instanceof Error ? err.message : m.settings_save_error();
    } finally {
      cwSaving = false;
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
        {m.settings_not_logged_in()} <a href="/login" class="link">{m.nav_login()}</a>
      </p>
    </div>
  {:else if loading}
    <LoadingSpinner variant="dots" padded />
  {:else if loadError}
    <ErrorAlert message={loadError} />
  {:else}
    <!-- Account -->
    <section class="card settings-card">
      <h2 class="card-label">{m.settings_account_section()}</h2>
      <form onsubmit={saveEmail} class="settings-form">
        <label class="settings-field">
          <span class="label-text">{m.settings_username_label()}</span>
          <input type="text" value={username} class="input" disabled />
        </label>
        <label class="settings-field">
          <span class="label-text">{m.settings_email_label()}</span>
          <input type="email" bind:value={email} class="input" required />
        </label>
        {#if emailError}
          <p class="settings-error">⚠ {emailError}</p>
        {/if}
        <button type="submit" class="btn-primary" disabled={emailSaving}>
          {emailSaved
            ? m.settings_saved()
            : emailSaving
              ? m.settings_saving()
              : m.settings_email_save()}
        </button>
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
        <button type="submit" class="btn-primary" disabled={cwSaving}>
          {cwSaved ? m.settings_saved() : cwSaving ? m.settings_saving() : m.settings_save()}
        </button>
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
</style>

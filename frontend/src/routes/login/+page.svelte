<script lang="ts">
  import { login } from '$lib/auth';
  import { localizeApiError } from '$lib/errorLocalization';
  import { goto } from '$app/navigation';
  import { localizeHref } from '$lib/paraglide/runtime';
  import * as m from '$lib/paraglide/messages';

  let username = $state('');
  let password = $state('');
  let err = $state('');
  let loading = $state(false);

  async function handleLogin(e: SubmitEvent) {
    e.preventDefault();
    err = '';
    loading = true;
    try {
      await login(username, password);
      await goto(localizeHref('/'));
    } catch (error) {
      err = localizeApiError(error, () => m.login_error_unknown());
    } finally {
      loading = false;
    }
  }
</script>

<main class="page-narrow auth-page">
  <div class="card auth-card">
    <h1 class="page-title auth-title">{m.login_title()}</h1>
    <p class="body-text auth-subtitle">{m.login_subtitle()}</p>

    <form onsubmit={handleLogin} class="auth-form">
      <label class="settings-field">
        <span class="label-text">{m.login_username_label()}</span>
        <input
          type="text"
          bind:value={username}
          class="input"
          placeholder={m.login_username_placeholder()}
          autocomplete="username"
          required
        />
      </label>

      <label class="settings-field">
        <span class="label-text">{m.login_password_label()}</span>
        <input
          type="password"
          bind:value={password}
          class="input"
          placeholder={m.login_password_placeholder()}
          autocomplete="current-password"
          required
        />
      </label>

      {#if err}
        <div class="result-box result-bad">
          <p class="body-text error-text">⚠ {err}</p>
        </div>
      {/if}

      <button type="submit" class="btn-primary" disabled={loading}>
        {loading ? m.login_submitting() : m.login_submit()}
      </button>
    </form>

    <hr class="divider" />
    <p class="body-text auth-footer-text">
      {m.login_no_account()}
      <a href={localizeHref('/register')} class="link">{m.login_register_link()}</a>
    </p>
  </div>
</main>

<style>
  .auth-page {
    padding-top: 1.5rem;
    padding-bottom: 1.5rem;
  }

  .auth-card {
    max-width: 38rem;
    margin-left: auto;
    margin-right: auto;
    padding: 1.75rem;
  }

  .auth-title {
    font-size: 2rem;
    line-height: 2.25rem;
    margin-bottom: 0.35rem;
  }

  .auth-subtitle {
    margin-bottom: 1.25rem;
  }

  .auth-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .auth-footer-text {
    margin: 0;
    text-align: center;
  }

  .error-text {
    margin: 0;
  }

  @media (max-width: 639px) {
    .auth-card {
      padding: 1.25rem;
    }

    .auth-title {
      font-size: 1.75rem;
      line-height: 2rem;
    }
  }
</style>

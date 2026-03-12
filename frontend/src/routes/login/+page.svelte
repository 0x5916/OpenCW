<script lang="ts">
  import { login } from '$lib/auth';
  import { goto } from '$app/navigation';
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
      await goto('/');
    } catch (error) {
      err = error instanceof Error ? error.message : m.login_error_unknown();
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
      {m.login_no_account()} <a href="/register" class="link">{m.login_register_link()}</a>
    </p>
  </div>
</main>

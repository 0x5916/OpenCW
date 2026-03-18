<script lang="ts">
  import { register } from '$lib/auth';
  import { goto } from '$app/navigation';
  import * as m from '$lib/paraglide/messages';

  let username = $state('');
  let email = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let err = $state('');
  let loading = $state(false);

  async function handleRegister(e: SubmitEvent) {
    e.preventDefault();
    err = '';

    if (password !== confirmPassword) {
      err = m.register_error_mismatch();
      return;
    }

    loading = true;
    try {
      await register(username, email, password);
      await goto('/');
    } catch (error) {
      err = error instanceof Error ? error.message : m.register_error_unknown();
    } finally {
      loading = false;
    }
  }
</script>

<main class="page-narrow auth-page">
  <div class="card auth-card">
    <h1 class="page-title auth-title">{m.register_title()}</h1>
    <p class="body-text auth-subtitle">{m.register_subtitle()}</p>

    <form onsubmit={handleRegister} class="auth-form">
      <label class="settings-field">
        <span class="label-text">{m.register_username_label()}</span>
        <input
          type="text"
          bind:value={username}
          class="input"
          placeholder={m.register_username_placeholder()}
          autocomplete="username"
          minlength="3"
          maxlength="16"
          required
        />
      </label>

      <label class="settings-field">
        <span class="label-text">{m.register_email_label()}</span>
        <input
          type="email"
          bind:value={email}
          class="input"
          placeholder={m.register_email_placeholder()}
          autocomplete="email"
          required
        />
      </label>

      <label class="settings-field">
        <span class="label-text">{m.register_password_label()}</span>
        <input
          type="password"
          bind:value={password}
          class="input"
          placeholder={m.register_password_placeholder()}
          autocomplete="new-password"
          minlength="8"
          required
        />
      </label>

      <label class="settings-field">
        <span class="label-text">{m.register_confirm_label()}</span>
        <input
          type="password"
          bind:value={confirmPassword}
          class="input"
          placeholder={m.register_confirm_placeholder()}
          autocomplete="new-password"
          minlength="8"
          required
        />
      </label>

      {#if err}
        <div class="result-box result-bad">
          <p class="body-text error-text">⚠ {err}</p>
        </div>
      {/if}

      <button type="submit" class="btn-primary" disabled={loading}>
        {loading ? m.register_submitting() : m.register_submit()}
      </button>
    </form>

    <hr class="divider" />
    <p class="body-text auth-footer-text">
      {m.register_has_account()} <a href="/login" class="link">{m.register_login_link()}</a>
    </p>
  </div>
</main>

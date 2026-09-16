<script lang="ts">
  import { login } from '$lib/auth';
  import { localizeApiError } from '$lib/errorLocalization';
  import { goto } from '$app/navigation';
  import AuthCard from '$lib/components/AuthCard.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import { localizedHref } from '$lib/i18n.svelte';
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
      await goto(localizedHref('/'));
    } catch (error) {
      err = localizeApiError(error, () => m.login_error_unknown());
    } finally {
      loading = false;
    }
  }
</script>

<AuthCard title={m.login_title()} subtitle={m.login_subtitle()}>
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
      <ErrorAlert message={err} />
    {/if}

    <button type="submit" class="btn-primary" disabled={loading}>
      {loading ? m.login_submitting() : m.login_submit()}
    </button>
  </form>

  {#snippet footer()}
    {m.login_no_account()}
    <a href={localizedHref('/register')} class="link">{m.login_register_link()}</a>
  {/snippet}
</AuthCard>

<style>
  .auth-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
</style>

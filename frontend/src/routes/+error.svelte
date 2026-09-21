<script lang="ts">
  import { page } from '$app/state';
  import { Home, Radio } from '@lucide/svelte';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';

  // Unknown URLs land here with a 404 from SvelteKit; anything else is a real
  // failure (a load that threw, a 500 from the API) and gets generic wording.
  const isNotFound = $derived(page.status === 404);
</script>

<div class="error-page page-narrow">
  <div class="panel error-card">
    <p class="error-code">{page.status}</p>
    <h1 class="page-title error-title">
      {isNotFound ? m.error_title() : m.error_generic_title()}
    </h1>
    <p class="body-text error-body">
      {isNotFound ? m.error_body() : m.error_generic_body()}
    </p>
    <div class="error-actions">
      <a href={href('/')} class="btn-ghost error-btn">
        <Home size={16} aria-hidden="true" />
        {m.error_cta_home()}
      </a>
      <a href={href('/morse/learn')} class="btn-ghost error-btn error-btn-accent">
        <Radio size={16} aria-hidden="true" />
        {m.home_cta()}
      </a>
    </div>
  </div>
</div>

<style>
  /* `.page-narrow` owns the width and the shell owns the vertical gap, so this
     page needs no size rules of its own. */
  .error-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    padding: 2.5rem 1.5rem;
    text-align: center;
  }

  .error-code {
    margin: 0;
    font-size: 0.8125rem;
    font-weight: 700;
    letter-spacing: 0.16em;
    color: var(--text-muted);
  }

  .error-title {
    margin: 0;
    color: var(--accent);
  }

  .error-body {
    margin: 0;
    max-width: 30rem;
  }

  .error-actions {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 0.75rem;
    width: 100%;
    margin-top: 1rem;
  }

  .error-btn {
    flex: 0 1 auto;
    text-decoration: none;
  }

  .error-btn-accent {
    border-color: var(--accent);
    color: var(--accent);
  }
</style>

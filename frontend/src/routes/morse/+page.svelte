<script lang="ts">
  import { ArrowRight } from '@lucide/svelte';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import { SITE_NAME } from '$lib/seo';
  import * as m from '$lib/paraglide/messages';
</script>

<svelte:head>
  <title>{m.morse_title()} | {SITE_NAME}</title>
</svelte:head>

<div class="page-narrow">
  <header class="morse-header">
    <h1 class="page-title">{m.morse_title()}</h1>
    <p class="body-text">{m.morse_intro()}</p>
  </header>

  <ul class="row-list">
    <li>
      <a class="row-link tool-row" href={href('/morse/learn')}>
        <span class="tool-copy">
          <strong class="tool-title">{m.nav_train()}</strong>
          <span class="tool-note">{m.home_preview_caption()}</span>
        </span>
        <ArrowRight size={16} aria-hidden="true" />
      </a>
    </li>
  </ul>
</div>

<style>
  .morse-header .body-text {
    margin: var(--space-2) 0 0;
  }

  .tool-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-4);
  }

  .tool-copy {
    display: grid;
    gap: var(--space-1);
    min-width: 0;
  }

  .tool-title {
    font-size: var(--text-base);
    font-weight: 600;
    color: var(--text-primary);
  }

  .tool-note {
    font-size: var(--text-sm);
    color: var(--text-muted);
  }
</style>

<script lang="ts">
  import { ArrowRight, Radio } from '@lucide/svelte';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';
</script>

<svelte:head>
  <title>Morse tools - OpenCW</title>
</svelte:head>

<div class="page-narrow morse-index">
  <header class="morse-header">
    <p class="eyebrow">Morse tools</p>
    <h1 class="page-title">{m.home_hero_title()}</h1>
    <p class="body-text">{m.home_hero_subtitle()}</p>
  </header>

  <nav class="morse-tools" aria-label="Morse tools">
    <a class="panel morse-tool" href={href('/morse/learn')}>
      <span class="morse-tool-icon"><Radio size={22} aria-hidden="true" /></span>
      <span class="morse-tool-copy">
        <strong>{m.home_cta()}</strong>
        <span>{m.home_preview_caption()}</span>
      </span>
      <ArrowRight size={20} aria-hidden="true" />
    </a>
  </nav>
</div>

<style>
  .morse-index {
    display: grid;
    gap: 2rem;
  }

  .morse-header {
    display: grid;
    gap: 0.65rem;
  }

  .morse-header .page-title,
  .morse-header .body-text {
    margin: 0;
  }

  .eyebrow {
    margin: 0;
    color: var(--accent);
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  .morse-tools {
    display: grid;
    gap: 0.75rem;
  }

  .morse-tool {
    display: flex;
    align-items: center;
    gap: 0.9rem;
    padding: 1rem;
    color: inherit;
    text-decoration: none;
  }

  .morse-tool-icon {
    display: grid;
    flex: 0 0 2.75rem;
    width: 2.75rem;
    height: 2.75rem;
    place-items: center;
    color: var(--accent);
    background: var(--bg-inset);
    border-radius: 0.5rem;
  }

  .morse-tool-copy {
    display: grid;
    flex: 1;
    gap: 0.2rem;
  }

  .morse-tool-copy span {
    color: var(--text-muted);
    font-size: 0.875rem;
  }
</style>

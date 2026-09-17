<script lang="ts">
  import { Construction, MessagesSquare, Radio } from '@lucide/svelte';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';
</script>

<svelte:head>
  <title>{m.forum_title()} | OpenCW</title>
</svelte:head>

<div class="forum-page page-wide">
  <header class="forum-hero">
    <div class="eyebrow"><MessagesSquare size={16} /> {m.forum_eyebrow()}</div>
    <h1 class="page-title">{m.forum_title()}</h1>
    <p class="body-text">{m.forum_under_development_body()}</p>
  </header>

  <section class="card forum-state">
    <Construction size={34} aria-hidden="true" />
    <h2 class="state-title">{m.forum_under_development_title()}</h2>
    <p class="body-text state-note">{m.forum_under_development_note()}</p>
    <a class="btn-primary state-action" href={href('/learn')}>
      <Radio size={17} />
      {m.nav_learn()}
    </a>
  </section>
</div>

<style>
  .forum-page {
    padding-top: 2.25rem;
    padding-bottom: 3rem;
  }

  .forum-hero {
    margin-bottom: 2rem;
  }

  .eyebrow {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-bottom: 0.6rem;
    color: var(--accent);
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  .forum-hero .page-title {
    margin: 0;
  }

  .forum-hero .body-text {
    max-width: 40rem;
    margin: 0.45rem 0 0;
  }

  .forum-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    padding: 2.5rem 1.5rem;
    text-align: center;
  }

  .forum-state > :global(svg) {
    color: var(--accent);
  }

  .state-title {
    margin: 0;
    font-size: 1.35rem;
    font-weight: 700;
  }

  .state-note {
    max-width: 34rem;
    margin: 0;
  }

  .state-action {
    display: inline-flex;
    width: auto;
    align-items: center;
    gap: 0.5rem;
    margin-top: 0.5rem;
    text-decoration: none;
  }
</style>

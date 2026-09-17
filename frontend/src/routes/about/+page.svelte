<script lang="ts">
  import {
    ArrowRight,
    BookOpen,
    CheckCircle,
    ShieldCheck,
    UserRound,
    Wrench
  } from '@lucide/svelte';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';

  const GITHUB_URL = 'https://github.com/0x5916';

  // Static lists are built once per component instance; the layout remounts the
  // page when the language changes, so `m.*` never goes stale here.
  const toc = [
    { id: 'koch-method', label: m.about_koch_title() },
    { id: 'the-project', label: m.about_project_title() },
    { id: 'about-the-author', label: m.about_author_title() },
    { id: 'your-data', label: m.about_privacy_title() }
  ];

  const kochSteps = [
    { title: m.about_koch_step1_title(), body: m.about_koch_step1_body() },
    { title: m.about_koch_step2_title(), body: m.about_koch_step2_body() },
    { title: m.about_koch_step3_title(), body: m.about_koch_step3_body() },
    { title: m.about_koch_step4_title(), body: m.about_koch_step4_body() }
  ];
</script>

<div class="about-page page-narrow">
  <header class="about-header">
    <h1 class="page-title">{m.about_title()}</h1>
    <p class="body-text about-intro">{m.about_intro()}</p>
  </header>

  <nav class="about-toc" aria-label={m.about_toc_label()}>
    <span class="card-label about-toc-label">{m.about_toc_label()}</span>
    <ul class="about-toc-list">
      {#each toc as item (item.id)}
        <li><a class="link" href={`#${item.id}`}>{item.label}</a></li>
      {/each}
    </ul>
  </nav>

  <!-- Koch Method Card -->
  <section class="card about-card" id="koch-method">
    <div class="about-card-icon" aria-hidden="true">
      <BookOpen size={32} />
    </div>
    <h2 class="section-title">{m.about_koch_title()}</h2>
    <p class="body-text about-intro">{m.about_koch_intro()}</p>

    <hr class="divider" />

    <!-- How it works -->
    <h3 class="sub-title">{m.about_koch_how_title()}</h3>
    <ol class="step-list">
      {#each kochSteps as item, index (item.title)}
        <li class="step">
          <span class="step-number">{index + 1}</span>
          <div>
            <strong class="step-title">{item.title}</strong>
            <p class="body-text">{item.body}</p>
          </div>
        </li>
      {/each}
    </ol>

    <hr class="divider" />

    <!-- Why it works -->
    <h3 class="sub-title">{m.about_koch_why_title()}</h3>
    <p class="body-text">{m.about_koch_why_body()}</p>

    <hr class="divider" />

    <!-- Comparison -->
    <h3 class="sub-title">{m.about_koch_vs_title()}</h3>
    <div class="koch-vs-grid">
      <div class="koch-vs-card koch-vs-bad">
        <span class="card-label">{m.about_koch_vs_traditional_label()}</span>
        <p class="body-text">{m.about_koch_vs_traditional_body()}</p>
      </div>
      <div class="koch-vs-card koch-vs-good">
        <div class="vs-good-header">
          <CheckCircle size={16} aria-hidden="true" />
          <span class="card-label">{m.about_koch_vs_koch_label()}</span>
        </div>
        <p class="body-text">{m.about_koch_vs_koch_body()}</p>
      </div>
    </div>
  </section>

  <!-- The project -->
  <section class="card about-card" id="the-project">
    <div class="about-card-icon" aria-hidden="true"><Wrench size={32} /></div>
    <h2 class="section-title">{m.about_project_title()}</h2>
    <p class="body-text">{m.about_project_body()}</p>
    <p class="body-text about-links">
      <a href={GITHUB_URL} class="link" rel="noopener noreferrer" target="_blank"
        >{m.about_project_github()}</a
      >
    </p>
  </section>

  <!-- About the author -->
  <section class="card about-card" id="about-the-author">
    <div class="about-card-icon" aria-hidden="true"><UserRound size={32} /></div>
    <h2 class="section-title">{m.about_author_title()}</h2>
    <p class="body-text">{m.about_author_p1()}</p>
    <p class="body-text">
      {m.about_author_p2_pre()}
      <a href={GITHUB_URL} class="link" rel="noopener noreferrer" target="_blank"
        >{m.about_author_github()}</a
      >.
    </p>
  </section>

  <!-- Your data -->
  <section class="card about-card" id="your-data">
    <div class="about-card-icon" aria-hidden="true"><ShieldCheck size={32} /></div>
    <h2 class="section-title">{m.about_privacy_title()}</h2>
    <p class="body-text">{m.about_privacy_body()}</p>
  </section>

  <div class="about-cta">
    <a href={href('/learn')} class="btn-cta"
      >{m.home_cta()}<ArrowRight size={18} aria-hidden="true" /></a
    >
  </div>
</div>

<style>
  .about-page {
    padding-top: 2rem;
    padding-bottom: 3rem;
  }

  .about-header {
    margin-bottom: 0.25rem;
  }

  .about-intro {
    margin: 0.5rem 0 0;
  }

  .about-toc {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .about-toc-label {
    margin-bottom: 0;
  }

  .about-toc-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.35rem 1rem;
    margin: 0;
    padding: 0;
    list-style: none;
    font-size: 0.875rem;
  }

  .about-card {
    display: flex;
    flex-direction: column;
    gap: 0;
    /* Keep anchored sections clear of the sticky navbar. */
    scroll-margin-top: 4.5rem;
  }

  .about-card-icon {
    color: var(--accent);
    margin-bottom: 0.75rem;
  }

  .about-links {
    margin-top: 0.75rem;
  }

  .koch-vs-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    margin-top: 0.75rem;
  }

  .koch-vs-card {
    border-radius: var(--radius-md);
    padding: 1rem;
    border: 1px solid;
  }

  .koch-vs-bad {
    background-color: var(--result-bad-bg);
    border-color: var(--result-bad-bd);
  }

  .koch-vs-good {
    background-color: var(--result-good-bg);
    border-color: var(--result-good-bd);
  }

  .vs-good-header {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--accent);
    margin-bottom: 0.25rem;
  }

  .vs-good-header .card-label {
    margin-bottom: 0;
    color: var(--accent);
  }

  .about-cta {
    margin-top: 1.5rem;
    display: flex;
    justify-content: center;
  }

  @media (max-width: 480px) {
    .koch-vs-grid {
      grid-template-columns: 1fr;
    }
  }
</style>

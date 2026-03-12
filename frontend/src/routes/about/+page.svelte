<script lang="ts">
  import { BookOpen, CheckCircle } from 'lucide-svelte';
  import * as m from '$lib/paraglide/messages';
</script>

<main class="about-page">
  <header class="about-header">
    <h1 class="page-title">{m.about_title()}</h1>
  </header>

  <!-- Koch Method Card -->
  <section class="card about-card">
    <div class="about-card-icon" aria-hidden="true">
      <BookOpen size={32} />
    </div>
    <h2 class="section-title">{m.about_koch_title()}</h2>
    <p class="body-text about-intro">{m.about_koch_intro()}</p>

    <hr class="divider" />

    <!-- How it works -->
    <h3 class="sub-title">{m.about_koch_how_title()}</h3>
    <ol class="koch-steps">
      <li class="koch-step">
        <span class="step-number">1</span>
        <div>
          <strong class="step-title">{m.about_koch_step1_title()}</strong>
          <p class="body-text">{m.about_koch_step1_body()}</p>
        </div>
      </li>
      <li class="koch-step">
        <span class="step-number">2</span>
        <div>
          <strong class="step-title">{m.about_koch_step2_title()}</strong>
          <p class="body-text">{m.about_koch_step2_body()}</p>
        </div>
      </li>
      <li class="koch-step">
        <span class="step-number">3</span>
        <div>
          <strong class="step-title">{m.about_koch_step3_title()}</strong>
          <p class="body-text">{m.about_koch_step3_body()}</p>
        </div>
      </li>
      <li class="koch-step">
        <span class="step-number">4</span>
        <div>
          <strong class="step-title">{m.about_koch_step4_title()}</strong>
          <p class="body-text">{m.about_koch_step4_body()}</p>
        </div>
      </li>
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

    <div class="about-cta">
      <a href="/morse/learn" class="btn-cta">Start Training →</a>
    </div>
  </section>
</main>

<style>
  .about-page {
    max-width: var(--max-width-narrow);
    margin: 0 auto;
    padding: 2rem 1rem 3rem;
  }

  .about-header {
    margin-bottom: 1.5rem;
  }

  .about-card {
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .about-card-icon {
    color: var(--accent);
    margin-bottom: 0.75rem;
  }

  .about-intro {
    margin-top: 0.5rem;
  }

  .koch-steps {
    list-style: none;
    padding: 0;
    margin: 0.75rem 0 0;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .koch-step {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
  }

  .step-number {
    flex-shrink: 0;
    width: 1.75rem;
    height: 1.75rem;
    border-radius: 50%;
    background-color: var(--accent);
    color: var(--bg-base);
    font-weight: 700;
    font-size: 0.85rem;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 0.1rem;
  }

  .step-title {
    display: block;
    color: var(--text-primary);
    font-weight: 600;
    margin-bottom: 0.25rem;
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

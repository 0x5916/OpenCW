<script lang="ts">
  import {
    Activity,
    ArrowRight,
    Check,
    ClipboardCheck,
    MessageSquare,
    Smartphone
  } from '@lucide/svelte';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import { LESSONS } from '$lib/morse';
  import * as m from '$lib/paraglide/messages';

  // Static lists are built once per component instance; the layout remounts the
  // page when the language changes, so `m.*` never goes stale here.
  const heroChips = [
    m.home_chip_nosignup(),
    m.home_chip_free(),
    m.home_chip_languages(),
    m.home_chip_sync()
  ];

  const steps = [
    { title: m.home_step1_title(), body: m.home_step1_body() },
    { title: m.home_step2_title(), body: m.home_step2_body() },
    { title: m.home_step3_title(), body: m.home_step3_body() }
  ];

  // Decorative trainer mock-up: show the character set a learner has unlocked.
  const PREVIEW_LESSON = 7;
  const previewChars = LESSONS.slice(0, PREVIEW_LESSON).join('').split('');
</script>

<!-- Hero -->
<section class="hero">
  <h1 class="hero-title">{m.home_hero_title()}</h1>
  <p class="hero-sub">{m.home_hero_subtitle()}</p>
  <div class="hero-actions">
    <a href={href('/learn')} class="btn-cta"
      >{m.home_cta()}<ArrowRight size={18} aria-hidden="true" /></a
    >
    <a href={href('/about')} class="hero-link">{m.home_hero_cta_secondary()}</a>
  </div>
  <ul class="hero-chips">
    {#each heroChips as chip (chip)}
      <li class="hero-chip"><Check size={14} aria-hidden="true" />{chip}</li>
    {/each}
  </ul>
</section>

<!-- What a lesson looks like, next to how the method works -->
<section class="preview-grid">
  <div class="preview-col">
    <!-- Decorative: the interactive trainer lives on /learn -->
    <div class="preview-frame" aria-hidden="true">
      <span class="card-label">{m.trainer_label_lesson()} {PREVIEW_LESSON}</span>
      <p class="preview-chars">
        {#each previewChars as char (char)}<span class="preview-char">{char}</span>{/each}
      </p>
      <div class="preview-answer">
        <span class="preview-answer-text">{m.trainer_answer_placeholder()}</span>
        <span class="preview-check"
          ><ClipboardCheck size={15} aria-hidden="true" />{m.trainer_check()}</span
        >
      </div>
    </div>
    <p class="body-text preview-caption">{m.home_preview_caption()}</p>
  </div>

  <div>
    <h2 class="section-title">{m.home_steps_title()}</h2>
    <ol class="step-list">
      {#each steps as step, index (step.title)}
        <li class="step">
          <span class="step-number">{index + 1}</span>
          <div>
            <strong class="step-title">{step.title}</strong>
            <p class="body-text">{step.body}</p>
          </div>
        </li>
      {/each}
    </ol>
  </div>
</section>

<!-- Feature cards: browsable, independent, peer-grouped -> genuinely cards. -->
<section class="card-grid feature-grid">
  <div class="card">
    <div class="feature-icon" aria-hidden="true"><Activity /></div>
    <h2 class="feature-title">{m.home_feature_progress_title()}</h2>
    <p class="feature-body">{m.home_feature_progress_body()}</p>
  </div>
  <div class="card">
    <div class="feature-icon" aria-hidden="true"><Smartphone /></div>
    <h2 class="feature-title">{m.home_feature_anywhere_title()}</h2>
    <p class="feature-body">{m.home_feature_anywhere_body()}</p>
  </div>
  <a href={href('/forum')} class="card card--interactive">
    <div class="feature-icon" aria-hidden="true"><MessageSquare /></div>
    <h2 class="feature-title">{m.home_forum_link_label()}</h2>
    <p class="feature-body">{m.home_forum_link_note()}</p>
  </a>
</section>

<!-- Who builds it -->
<section class="author-band">
  <p class="body-text author-line">{m.home_author_line()}</p>
  <div class="author-links">
    <a href={href('/about')} class="link">{m.home_author_link()}</a>
    <a href="https://github.com/0x5916" class="link" rel="noopener noreferrer" target="_blank"
      >{m.home_author_github()}</a
    >
  </div>
</section>

<!-- Closing call to action -->
<section class="final-cta">
  <h2 class="section-title final-title">{m.home_final_title()}</h2>
  <p class="body-text final-body">{m.home_final_body()}</p>
  <a href={href('/learn')} class="btn-cta"
    >{m.home_cta()}<ArrowRight size={18} aria-hidden="true" /></a
  >
</section>

<style>
  /* One separation idiom for the whole page: every block below the hero carries
     the same top margin, and the redundant `<hr class="divider">` separators are
     gone — a transition used to be spaced three times over (hero bottom padding
     + rule margins + the next section's margin). */
  .hero {
    text-align: center;
  }

  .hero-title {
    color: var(--accent);
    font-size: var(--text-3xl);
    line-height: var(--leading-tight);
    font-weight: 800;
    letter-spacing: -0.03em;
    max-width: 34rem;
    margin: 0 auto 1rem;
  }

  @media (min-width: 640px) {
    .hero-title {
      font-size: var(--text-4xl);
    }
  }

  .hero-sub {
    color: var(--text-secondary);
    font-size: var(--text-lg);
    line-height: var(--leading-relaxed);
    max-width: 40rem;
    margin: 0 auto 1.75rem;
  }

  .hero-actions {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;
    gap: 0.75rem 1.25rem;
  }

  .hero-link {
    color: var(--accent);
    font-weight: 600;
    text-decoration: none;
  }

  .hero-link:hover {
    text-decoration: underline;
  }

  .hero-chips {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 0.5rem 0.75rem;
    max-width: 40rem;
    margin: 2rem auto 0;
    padding: 0;
    list-style: none;
  }

  .hero-chip {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.3rem 0.7rem;
    border: 1px solid var(--border);
    border-radius: 999px;
    background-color: var(--bg-surface);
    color: var(--text-muted);
    font-size: 0.8125rem;
    line-height: 1.2rem;
  }

  .hero-chip :global(svg) {
    flex-shrink: 0;
    color: var(--accent);
  }

  .preview-grid {
    display: grid;
    gap: var(--block-gap);
    margin-top: var(--section-gap);
  }

  @media (min-width: 800px) {
    .preview-grid {
      grid-template-columns: minmax(0, 20rem) minmax(0, 1fr);
      align-items: center;
    }
  }

  .preview-col {
    min-width: 0;
  }

  /* Decorative mock of the trainer: a frame, not a card. The dashed edge keeps
     it legible as an illustration — a solid edge at the card radius reads as the
     real trainer — and nothing inside it responds to hover. */
  .preview-frame {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
    padding: 1.25rem;
    border: 1px dashed var(--border-subtle);
    border-radius: var(--radius-lg);
    background-color: var(--bg-surface);
  }

  .preview-frame :global(.card-label) {
    margin-bottom: 0;
  }

  .preview-chars {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin: 0;
  }

  .preview-char {
    min-width: 2rem;
    padding: 0.35rem 0.5rem;
    text-align: center;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    font-size: 1.125rem;
    font-weight: 700;
    color: var(--accent);
    background-color: var(--bg-inset);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
  }

  .preview-answer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    background-color: var(--bg-inset);
  }

  .preview-answer-text {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: var(--text-muted);
    font-size: 0.8125rem;
  }

  /* Inside the mock the "check" affordance is a neutral chip, never a filled
     button: this whole frame is aria-hidden decoration, so anything that looks
     pressable here is a false affordance. */
  .preview-check {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    flex-shrink: 0;
    padding: 0.35rem 0.7rem;
    border-radius: var(--radius-sm);
    border: 1px solid var(--border);
    background-color: var(--bg-surface);
    color: var(--text-muted);
    font-size: 0.8125rem;
    font-weight: 600;
  }

  .preview-caption {
    margin: 0.85rem 0 0;
  }

  .author-band {
    margin-top: var(--section-gap);
    text-align: center;
  }

  .author-line {
    max-width: 40rem;
    margin: 0 auto 0.75rem;
  }

  .author-links {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 1.25rem;
  }

  .final-cta {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.4rem;
    margin-top: var(--section-gap);
    text-align: center;
  }

  .final-cta .final-title {
    margin: 0;
  }

  .final-body {
    max-width: 34rem;
    margin: 0 0 0.75rem;
  }

  .feature-grid {
    margin-top: var(--section-gap);
  }

  .feature-icon {
    color: var(--accent);
    font-size: 1.875rem;
    line-height: 2.25rem;
    margin-bottom: 0.75rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .feature-icon :global(svg) {
    width: 1.875rem;
    height: 1.875rem;
  }

  .feature-title {
    color: var(--text-primary);
    font-weight: 700;
    font-size: var(--text-lg);
    line-height: var(--leading-snug);
    margin-bottom: 0.25rem;
  }

  .feature-body {
    color: var(--text-secondary);
    font-size: var(--text-sm);
    line-height: var(--leading-normal);
  }
</style>

<script lang="ts">
  import { Trophy, TrendingUp, X, RefreshCw } from 'lucide-svelte';
  import type { DiffToken } from '$lib/score';
  import * as m from '$lib/paraglide/messages';
  import { onMount } from 'svelte';

  interface Props {
    result: number;
    diffTokens: DiffToken[];
    hasNextLesson: boolean;
    hasPrevLesson: boolean;
    shouldRegenerate: boolean;
    nextLessonNum: number;
    prevLessonNum: number;
    onClose: () => void;
    onNext: () => void;
    onPrev: () => void;
    onRegenerate: () => void;
  }

  let {
    result,
    diffTokens,
    hasNextLesson,
    hasPrevLesson,
    shouldRegenerate,
    nextLessonNum,
    prevLessonNum,
    onClose,
    onNext,
    onPrev,
    onRegenerate
  }: Props = $props();

  let ScoreIcon = $derived(result >= 0.9 ? Trophy : TrendingUp);
  let scoreText = $derived(
    result >= 0.9
      ? m.trainer_score_great()
      : result >= 0.7
        ? m.trainer_score_good()
        : m.trainer_score_bad()
  );
  let colorClass = $derived(
    result >= 0.9 ? 'overlay-good' : result >= 0.7 ? 'overlay-ok' : 'overlay-bad'
  );
  let pct = $derived(Math.round(result * 100) + '%');

  let panelRef: HTMLDivElement;

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault();
      onClose();
    }
  }

  function handleBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      onClose();
    }
  }

  onMount(() => {
    // Focus the close button when modal opens
    const closeBtn = panelRef?.querySelector('.overlay-close') as HTMLButtonElement;
    if (closeBtn) closeBtn.focus();
  });
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="overlay-backdrop" onclick={handleBackdropClick} role="presentation">
  <div class="overlay-panel" bind:this={panelRef} role="dialog" aria-modal="true">
    <!-- Close -->
    <button class="overlay-close" onclick={onClose} aria-label="Close"><X size={20} /></button>

    <!-- Score hero -->
    <header class="overlay-hero {colorClass}">
      <ScoreIcon size={40} class="overlay-score-icon" />
      <p class="overlay-pct">{pct}</p>
      <p class="overlay-score-text">{scoreText}</p>
    </header>

    <!-- Diff -->
    <section class="overlay-diff-section">
      <div class="overlay-diff-header">
        <h2 class="card-label">{m.trainer_diff_title()}</h2>
        <div class="diff-legend">
          <span class="diff-token diff-correct">{m.trainer_diff_legend_correct()}</span>
          <span class="diff-token diff-sub">{m.trainer_diff_legend_sub()}</span>
          <span class="diff-token diff-missing">{m.trainer_diff_legend_missing()}</span>
          <span class="diff-token diff-extra">{m.trainer_diff_legend_extra()}</span>
        </div>
      </div>
      <div class="diff-tokens overlay-diff-tokens">
        {#each diffTokens as tok, i (i)}
          {#if tok.type === 'correct'}
            <span class="diff-token diff-correct">{tok.ref}</span>
          {:else if tok.type === 'substitution'}
            <span class="diff-token diff-sub" title="expected: {tok.ref}"
              >{tok.inp}<span class="diff-expected"> ({tok.ref})</span></span
            >
          {:else if tok.type === 'missing'}
            <span class="diff-token diff-missing">{tok.ref}</span>
          {:else if tok.type === 'extra'}
            <span class="diff-token diff-extra">{tok.inp}</span>
          {/if}
        {/each}
      </div>
    </section>

    <!-- Actions -->
    <footer class="overlay-actions">
      {#if hasPrevLesson}
        <button class="btn-prev-lesson" onclick={onPrev}>
          {m.trainer_prev_lesson({ lesson: String(prevLessonNum) })}
        </button>
      {/if}
      {#if shouldRegenerate}
        <button class="btn-prev-lesson" onclick={onRegenerate}>
          <RefreshCw size={16} />{m.trainer_try_again()}
        </button>
      {/if}
      {#if hasNextLesson}
        <button class="btn-next-lesson" onclick={onNext}>
          {m.trainer_next_lesson({ lesson: String(nextLessonNum) })}
        </button>
      {/if}
      {#if !hasNextLesson && !hasPrevLesson && !shouldRegenerate}
        <button class="btn-regen" onclick={onClose}>
          <RefreshCw size={16} />{m.trainer_try_again()}
        </button>
      {/if}
    </footer>
  </div>
</div>

<style>
  .overlay-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(4px);
    z-index: 200;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
  }
  .overlay-panel {
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: 1rem;
    width: 100%;
    max-width: 680px;
    max-height: 90vh;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    position: relative;
    box-shadow: 0 24px 60px rgba(0, 0, 0, 0.4);
  }
  .overlay-close {
    position: absolute;
    top: 0.75rem;
    right: 0.75rem;
    background: none;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.35rem;
    border-radius: 0.4rem;
    display: flex;
    align-items: center;
    z-index: 1;
    transition:
      color 0.15s,
      background 0.15s;
  }
  .overlay-close:hover {
    color: var(--text-primary);
    background: var(--bg-inset);
  }

  /* Score hero */
  .overlay-hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.4rem;
    padding: 2.5rem 2rem 2rem;
    border-radius: 1rem 1rem 0 0;
  }
  .overlay-good {
    background: rgba(34, 197, 94, 0.12);
  }
  .overlay-ok {
    background: rgba(234, 179, 8, 0.12);
  }
  .overlay-bad {
    background: rgba(239, 68, 68, 0.12);
  }
  :global(.overlay-score-icon) {
    opacity: 0.85;
  }
  .overlay-good :global(.overlay-score-icon) {
    color: #22c55e;
  }
  .overlay-ok :global(.overlay-score-icon) {
    color: #ca8a04;
  }
  .overlay-bad :global(.overlay-score-icon) {
    color: #ef4444;
  }
  .overlay-pct {
    font-size: 3.5rem;
    font-weight: 900;
    line-height: 1;
    margin: 0;
  }
  .overlay-good .overlay-pct {
    color: #22c55e;
  }
  .overlay-ok .overlay-pct {
    color: #ca8a04;
  }
  .overlay-bad .overlay-pct {
    color: #ef4444;
  }
  .overlay-score-text {
    font-size: 1rem;
    color: var(--text-secondary);
    margin: 0;
  }

  /* Diff section */
  .overlay-diff-section {
    padding: 1.25rem 1.5rem;
    border-top: 1px solid var(--border);
  }
  .overlay-diff-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 0.875rem;
  }
  .overlay-diff-header :global(.card-label) {
    margin-bottom: 0;
  }
  .overlay-diff-tokens {
    max-height: 40vh;
    overflow-y: auto;
  }

  /* Actions */
  .overlay-actions {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    padding: 1rem 1.5rem 1.5rem;
    border-top: 1px solid var(--border);
  }
  .overlay-actions > * {
    flex: 1;
    min-width: 150px;
    justify-content: center;
    margin-top: 0;
  }
</style>

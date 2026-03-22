<script lang="ts">
  import { generateTimedLesson, LESSONS } from '$lib/morse';
  import MorsePlayer from '$lib/components/MorsePlayer.svelte';
  import ResultOverlay from '$lib/components/ResultOverlay.svelte';
  import { untrack } from 'svelte';
  import { ClipboardCheck } from 'lucide-svelte';
  import { score, diffWords } from '$lib/score';
  import type { DiffToken } from '$lib/score';
  import { user } from '$lib/auth';
  import { submitProgress } from '$lib/api';
  import { normalizeLesson } from '$lib/cwSync';
  import * as m from '$lib/paraglide/messages';

  let { data } = $props();

  let inputText = $state('');

  // svelte-ignore state_referenced_locally
  let chosenLesson = $state(data.lesson);

  let result = $state(-1);
  let showOverlay = $state(false);
  let diffTokens = $state<DiffToken[]>([]);
  let charWpm = $state(20);
  let effWpm = $state(10);
  let freq = $state(600);
  let volume = $state(1);
  let startDelay = $state(0.5);

  $effect(() => {
    const val = String(normalizeLesson(chosenLesson, LESSONS.length));
    localStorage.setItem('learn.lesson', val);
    document.cookie = `learn.lesson=${val}; path=/; max-age=31536000; SameSite=Lax`;
  });

  let lessonText = $derived(generateTimedLesson(chosenLesson, 60, charWpm, effWpm));
  let currentLessonWord = $derived(LESSONS.slice(0, chosenLesson).join(''));
  let currentLessonChars = $derived(currentLessonWord.split('').filter(Boolean));
  let selectedLessonChar = $state(LESSONS[0]?.[0] ?? '');
  let fullLessonPlayer = $state<{
    playNow: () => Promise<void>;
    stopNow: () => Promise<void>;
    isStarted: () => boolean;
  } | null>(null);
  function regenerate() {
    lessonText = generateTimedLesson(chosenLesson, 60, charWpm, effWpm);
    inputText = '';
    result = -1;
    diffTokens = [];
    showOverlay = false;
  }

  async function checkResult() {
    result = score(lessonText, inputText);
    diffTokens = diffWords(lessonText, inputText);
    showOverlay = true;
    if ($user && result > 0) {
      submitProgress(chosenLesson, charWpm, effWpm, result).catch(() => {});
    }
  }

  let hasNextLesson = $derived(result >= 0.9 && chosenLesson < LESSONS.length);
  let hasPrevLesson = $derived(result < 0.7 && chosenLesson > 1);
  let shouldRegenerate = $derived(result >= 0.7 && result < 0.9);

  function prevLesson() {
    chosenLesson -= 1;
    inputText = '';
    result = -1;
    diffTokens = [];
    showOverlay = false;
  }

  function nextLesson() {
    chosenLesson += 1;
    inputText = '';
    result = -1;
    diffTokens = [];
    showOverlay = false;
  }

  function onLessonSelectChange() {
    result = -1;
  }

  function onCwSettingInput() {}

  $effect(() => {
    const char = LESSONS[chosenLesson - 1] ?? '';
    const chars = currentLessonChars;
    untrack(() => {
      if (char) selectedLessonChar = char;
      if (chars.length === 0) {
        selectedLessonChar = '';
      } else if (!chars.includes(selectedLessonChar)) {
        selectedLessonChar = chars[0];
      }
    });
  });

  function onSelectedCharChange(event: Event) {
    selectedLessonChar = (event.currentTarget as HTMLSelectElement).value;
  }

  async function onAnswerInput(event: Event) {
    const value = (event.currentTarget as HTMLTextAreaElement).value;

    // Intercept a lone space to start the player
    if (value === ' ' && !fullLessonPlayer?.isStarted()) {
      inputText = '';
      (event.currentTarget as HTMLTextAreaElement).value = '';
      await fullLessonPlayer?.playNow();
      return;
    }

    inputText = value.toUpperCase();
  }


  async function onAnswerKeydown(event: KeyboardEvent) {
    if (event.code !== 'Space' || inputText.trim() !== '' || fullLessonPlayer?.isStarted()) return;
    event.preventDefault()
    await fullLessonPlayer?.playNow()
  }
</script>

<!-- Full-width heading -->
<header class="learn-heading">
  <h1 class="page-title">{m.trainer_title()}</h1>
  <p class="page-title-sub">
    <span class="accent-text">{m.trainer_subtitle_pre()}</span>
    {m.trainer_subtitle_post()}
  </p>
</header>

<main class="learn-page">
  <!-- Left column: lesson + settings + player -->
  <div class="learn-col-left">
    <div class="card-sm">
      <h2 class="card-label">{m.trainer_label_lesson()}</h2>
      <div class="lesson-row">
        <p class="lesson-current-label">{m.trainer_current_lesson()}</p>
        <select bind:value={chosenLesson} onchange={onLessonSelectChange} class="lesson-select">
          {#each LESSONS as lesson, index (index)}
            <option value={index + 1}>{index + 1} — {lesson.split('').join(', ')}</option>
          {/each}
        </select>
      </div>
    </div>

    <div class="card-sm">
      <h2 class="card-label">{m.trainer_label_current_chars()}</h2>
      <div class="lesson-char-row">
        <p class="lesson-char-preview">{m.trainer_choose_letter()}</p>
        <select
          bind:value={selectedLessonChar}
          onchange={onSelectedCharChange}
          class="lesson-select lesson-char-select"
        >
          {#each currentLessonChars as char (char)}
            <option value={char}>{char}</option>
          {/each}
        </select>
      </div>
      <MorsePlayer
        text={Array(5).fill(selectedLessonChar).join("")}
        {charWpm}
        {effWpm}
        {freq}
        {volume}
        compact
        showSettings
        mediaStyle
        onSettingsInput={onCwSettingInput}
        playLabel={currentLessonChars.length > 1
          ? `Play "${selectedLessonChar}"`
          : m.trainer_play_letter()}
      />
    </div>

  </div>

  <!-- Right column: answer + result -->
  <div class="learn-col-right">
    <div class="answer-card learn-answer-card">
      <MorsePlayer
        bind:this={fullLessonPlayer}
        text={lessonText}
        {charWpm}
        {effWpm}
        {freq}
        {volume}
        {startDelay}
        showSettings
        mediaStyle
        onSettingsInput={onCwSettingInput}
        label={m.player_label()}
      />
      <textarea
        placeholder={`${m.trainer_answer_placeholder()}\n${m.trainer_answer_shortcut_tip()}`}
        bind:value={inputText}
        oninput={onAnswerInput}
        onkeydown={onAnswerKeydown}
        autocapitalize="characters"
        autocomplete="off"
        autocorrect="off"
        spellcheck="false"
        class="textarea learn-answer-textarea"
      ></textarea>
      <button onclick={checkResult} class="btn-primary"
        ><ClipboardCheck size={16} />{m.trainer_check()}</button
      >
    </div>
  </div>
</main>

{#if showOverlay}
  <ResultOverlay
    {result}
    {diffTokens}
    {hasNextLesson}
    {hasPrevLesson}
    {shouldRegenerate}
    nextLessonNum={chosenLesson + 1}
    prevLessonNum={chosenLesson - 1}
    onClose={() => (showOverlay = false)}
    onNext={nextLesson}
    onPrev={prevLesson}
    onRegenerate={regenerate}
  />
{/if}

<style>
  .learn-page {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    min-width: 0;
  }

  .learn-col-left,
  .learn-col-right {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    min-width: 0;
  }

  .learn-heading {
    margin-bottom: 1rem;
  }

  .learn-page :global(.card-label) {
    color: var(--text-label);
  }

  .learn-answer-card :global(.card-label) {
    margin-bottom: 0;
  }

  .page-title-sub {
    color: var(--text-secondary);
    font-size: 0.875rem;
    line-height: 1.5;
    margin-top: 0.25rem;
  }

  .lesson-row {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.6rem;
  }

  .lesson-current-label {
    margin: 0;
    color: var(--accent);
    font-weight: 500;
    font-size: 0.875rem;
    line-height: 1.25rem;
    white-space: nowrap;
  }

  .lesson-select {
    flex: 1;
    min-width: 0;
    width: 100%;
    padding: 0.5rem 2rem 0.5rem 0.75rem;
    font-size: 0.875rem;
    line-height: 1.25rem;
    border-radius: var(--radius-md);
    border: 1px solid var(--border-subtle);
    background-color: var(--bg-inset);
    color: var(--text-primary);
    transition: background-color 0.2s, border-color 0.2s, color 0.2s;
  }

  .lesson-select:hover {
    border-color: color-mix(in srgb, var(--accent) 45%, var(--border-subtle));
  }

  .lesson-select:focus {
    outline: none;
    box-shadow: var(--focus-ring);
  }

  .lesson-char-row {
    display: grid;
    grid-template-columns: max-content minmax(0, 1fr);
    align-items: center;
    column-gap: 0.65rem;
    row-gap: 0.45rem;
    margin-bottom: 0.65rem;
  }

  .lesson-char-preview {
    color: var(--text-secondary);
    font-size: 0.8125rem;
    font-weight: 600;
    margin: 0;
    line-height: 2.15rem;
    white-space: nowrap;
  }

  .lesson-char-select {
    margin-left: 0;
    min-width: 0;
    width: 100%;
    max-width: none;
  }

  .answer-card {
    background-color: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: var(--radius-lg);
    padding: 1.25rem;
    min-width: 0;
    box-shadow: var(--shadow-soft);
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    transition: background-color var(--transition-base), border-color var(--transition-base), box-shadow var(--transition-base);
  }

  .answer-card:hover {
    box-shadow: var(--shadow-lift);
  }

  @media (max-width: 767px) {
    .learn-col-right .learn-answer-textarea {
      min-height: 11rem;
    }
  }

  @media (min-width: 768px) {
    .learn-page {
      display: grid;
      grid-template-columns: 5fr 7fr;
      column-gap: 1rem;
      row-gap: 1rem;
    }

    .learn-col-right > .learn-answer-card {
      flex: 1;
      display: flex;
      flex-direction: column;
    }

    .learn-col-right .learn-answer-textarea {
      flex: 1 1 0;
      min-height: 6rem;
      field-sizing: fixed;
    }
  }

  @media (max-width: 640px) {
    .lesson-char-row {
      grid-template-columns: max-content minmax(0, 1fr);
      align-items: center;
      column-gap: 0.45rem;
    }

    .lesson-char-select {
      max-width: none;
    }

    .lesson-char-preview {
      font-size: 0.78rem;
    }
  }

</style>

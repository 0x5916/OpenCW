<script lang="ts">
  import { browser } from '$app/environment';
  import { generateTimedLesson, LESSONS } from '$lib/morse';
  import MorsePlayer from '$lib/components/MorsePlayer.svelte';
  import ResultOverlay from '$lib/components/ResultOverlay.svelte';
  import GuestNotice from '$lib/components/GuestNotice.svelte';
  import { untrack, onDestroy } from 'svelte';
  import { ClipboardCheck } from '@lucide/svelte';
  import { CW_STORAGE_KEYS, UI_STORAGE_KEYS } from '$lib/storageKeys';
  import { writeCookie } from '$lib/cookies';
  import { langPreference } from '$lib/i18n.svelte';
  import { SCORE_GOOD, SCORE_OK, score, diffWords } from '$lib/score';
  import type { DiffToken } from '$lib/score';
  import { user } from '$lib/auth';
  import { saveProgressOfflineFirst } from '$lib/progressSync';
  import {
    normalizeLesson,
    readClientCwSettings,
    readClientPageSettings,
    restoreSettingsFromServer,
    saveClientCwSettings,
    syncSettingsToServer
  } from '$lib/cwSync';
  import * as m from '$lib/paraglide/messages';

  let { data } = $props();

  let inputText = $state('');

  // svelte-ignore state_referenced_locally
  let chosenLesson = $state(data.lesson);

  let result = $state(-1);
  let showOverlay = $state(false);
  let diffTokens = $state<DiffToken[]>([]);
  let showQuickStart = $state(false);
  let showQuickTips = $state(false);
  let charWpm = $state(20);
  let effWpm = $state(10);
  let freq = $state(600);
  let volume = $state(1);
  let startDelay = $state(0.5);
  let autoSyncTimeout: ReturnType<typeof setTimeout> | null = null;

  $effect(() => {
    if (!browser) return;
    showQuickStart = localStorage.getItem(UI_STORAGE_KEYS.quickstartDismissed) !== '1';
  });

  $effect(() => {
    if (!browser) return;
    const localCw = readClientCwSettings();
    charWpm = localCw.char_wpm;
    effWpm = localCw.eff_wpm;
    freq = localCw.freq;
    startDelay = localCw.start_delay;
  });

  function dismissQuickStart() {
    showQuickStart = false;
    showQuickTips = false;
    if (!browser) return;
    localStorage.setItem(UI_STORAGE_KEYS.quickstartDismissed, '1');
  }

  function openQuickTips() {
    showQuickTips = true;
  }

  function closeQuickTips() {
    showQuickTips = false;
  }

  $effect(() => {
    if (!$user) return;
    restoreSettingsFromServer()
      .then(({ cw }) => {
        charWpm = cw.char_wpm;
        effWpm = cw.eff_wpm;
        freq = cw.freq;
        startDelay = cw.start_delay;
        saveClientCwSettings(cw);
      })
      .catch(() => {
        // Keep local defaults if server restore fails.
      });
  });

  $effect(() => {
    const val = String(normalizeLesson(chosenLesson, LESSONS.length));
    localStorage.setItem(CW_STORAGE_KEYS.lesson, val);
    writeCookie(CW_STORAGE_KEYS.lesson, val);
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
    resetResultState({ clearInput: true });
  }

  async function checkResult() {
    result = score(lessonText, inputText);
    diffTokens = diffWords(lessonText, inputText);
    showOverlay = true;
    if (result > 0) {
      saveProgressOfflineFirst({
        lesson: chosenLesson,
        char_wpm: charWpm,
        eff_wpm: effWpm,
        accuracy: result
      }).catch(() => {});
    }
  }

  let hasNextLesson = $derived(result >= SCORE_GOOD && chosenLesson < LESSONS.length);
  let hasPrevLesson = $derived(result < SCORE_OK && chosenLesson > 1);
  let shouldRegenerate = $derived(result >= SCORE_OK && result < SCORE_GOOD);

  function prevLesson() {
    chosenLesson -= 1;
    resetResultState({ clearInput: true });
    scheduleApiSync();
  }

  function nextLesson() {
    chosenLesson += 1;
    resetResultState({ clearInput: true });
    scheduleApiSync();
  }

  function resetResultState(options: { clearInput?: boolean } = {}) {
    if (options.clearInput) {
      inputText = '';
    }

    result = -1;
    diffTokens = [];
    showOverlay = false;
  }

  function onLessonSelectChange() {
    result = -1;
    scheduleApiSync();
  }

  function onCwSettingInput() {
    saveClientCwSettings({ char_wpm: charWpm, eff_wpm: effWpm, freq, start_delay: startDelay });
    scheduleApiSync();
  }

  function scheduleApiSync() {
    if (!$user) return;
    if (autoSyncTimeout) clearTimeout(autoSyncTimeout);
    autoSyncTimeout = setTimeout(() => {
      void syncSettings();
    }, 1000);
  }

  async function syncSettings() {
    if (!$user) return;
    const cw = { char_wpm: charWpm, eff_wpm: effWpm, freq, start_delay: startDelay };
    const page = readClientPageSettings(chosenLesson, LESSONS.length, langPreference.value);
    try {
      await syncSettingsToServer(cw, page);
    } catch {
      // Keep training flow uninterrupted if sync fails.
    }
  }

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
    event.preventDefault();
    await fullLessonPlayer?.playNow();
  }

  onDestroy(() => {
    if (autoSyncTimeout) clearTimeout(autoSyncTimeout);
  });
</script>

<!-- Full-width heading -->
<header class="learn-heading">
  <h1 class="page-title">{m.trainer_title()}</h1>
</header>

{#if showQuickStart}
  <section class="card-sm quickstart-card" aria-label={m.trainer_quickstart_aria()}>
    <h2 class="quickstart-title">{m.trainer_quickstart_title()}</h2>
    <ol class="quickstart-steps">
      <li>{m.trainer_quickstart_step1()}</li>
      <li>{m.trainer_quickstart_step2()}</li>
      <li>{m.trainer_quickstart_step3()}</li>
    </ol>
    <div class="quickstart-actions">
      <button type="button" class="btn-success quickstart-btn" onclick={dismissQuickStart}
        >{m.trainer_quickstart_start()}</button
      >
      <button type="button" class="btn-ghost quickstart-btn" onclick={openQuickTips}
        >{m.trainer_quickstart_tips()}</button
      >
    </div>

    {#if showQuickTips}
      <div class="quickstart-modal-backdrop" role="presentation" onclick={closeQuickTips}>
        <div
          class="quickstart-modal card-sm"
          role="dialog"
          aria-modal="true"
          aria-label={m.trainer_quickstart_tips_title()}
          tabindex="-1"
          onclick={(event) => event.stopPropagation()}
          onkeydown={(event) => {
            if (event.key === 'Escape') closeQuickTips();
          }}
        >
          <h3 class="quickstart-modal-title">{m.trainer_quickstart_tips_title()}</h3>
          <ul class="quickstart-modal-list">
            <li>{m.trainer_quickstart_tip1()}</li>
            <li>{m.trainer_quickstart_tip2()}</li>
            <li>{m.trainer_quickstart_tip3()}</li>
            {#if !$user}
              <li><GuestNotice /></li>
            {/if}
          </ul>
          <div class="quickstart-actions">
            <button type="button" class="btn-ghost quickstart-btn" onclick={closeQuickTips}
              >{m.trainer_quickstart_close()}</button
            >
          </div>
        </div>
      </div>
    {/if}
  </section>
{/if}

<div class="learn-page">
  <!-- Left column: one tool panel — lesson select, then character select. -->
  <div class="learn-col-left">
    <section
      class="panel trainer-panel"
      class:lesson-char-highlight={showQuickStart && chosenLesson === 1}
    >
      <h2 class="card-title">{m.trainer_label_lesson()}</h2>
      <div class="lesson-row">
        <p class="lesson-current-label">{m.trainer_current_lesson()}</p>
        <select bind:value={chosenLesson} onchange={onLessonSelectChange} class="select">
          {#each LESSONS as lesson, index (index)}
            <option value={index + 1}>{index + 1} — {lesson.split('').join(', ')}</option>
          {/each}
        </select>
      </div>

      <hr class="trainer-divider" />

      <h2 class="card-title">{m.trainer_label_current_chars()}</h2>
      <div class="lesson-char-row">
        <p class="lesson-char-preview">{m.trainer_choose_letter()}</p>
        <select
          bind:value={selectedLessonChar}
          onchange={onSelectedCharChange}
          class="select lesson-char-select"
        >
          {#each currentLessonChars as char (char)}
            <option value={char}>{char}</option>
          {/each}
        </select>
      </div>
      <MorsePlayer
        text={Array(5).fill(selectedLessonChar).join('')}
        {charWpm}
        {effWpm}
        {freq}
        {volume}
        compact
        showSettings
        mediaStyle
        onSettingsInput={onCwSettingInput}
        playLabel={currentLessonChars.length > 1
          ? m.trainer_play_char({ char: selectedLessonChar })
          : m.trainer_play_letter()}
      />
    </section>
  </div>

  <!-- Right column: answer + result. Both columns are peer tool surfaces, so
       they share one primitive and therefore one padding. -->
  <div class="learn-col-right">
    <div class="panel learn-answer-card">
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
        class="textarea learn-answer-textarea"></textarea>
      <button onclick={checkResult} class="btn-primary"
        ><ClipboardCheck size={16} />{m.trainer_check()}</button
      >
    </div>
  </div>
</div>

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
    gap: var(--block-gap);
    min-width: 0;
  }

  .learn-col-left,
  .learn-col-right {
    display: flex;
    flex-direction: column;
    gap: var(--block-gap);
    min-width: 0;
  }

  .learn-heading {
    margin-bottom: var(--block-gap);
  }

  .trainer-panel {
    display: flex;
    flex-direction: column;
  }

  .trainer-divider {
    width: 100%;
    border: none;
    border-top: 1px solid var(--border);
    margin: var(--block-gap) 0;
  }

  .learn-answer-card {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .learn-answer-card :global(.card-title) {
    margin-bottom: 0;
  }

  /* Onboarding hint: flat like every other surface, with an accent edge to mark
     it as the transient one. Gradients are not part of the card system. */
  .quickstart-card {
    margin-bottom: var(--block-gap);
    border-color: color-mix(in srgb, var(--accent) 32%, var(--border-card));
  }

  .quickstart-title {
    margin: 0;
    color: var(--text-primary);
    font-size: var(--text-base);
    line-height: var(--leading-snug);
    font-weight: 600;
  }

  .quickstart-steps {
    margin: 0.75rem 0 0;
    padding-left: 1.25rem;
    color: var(--text-secondary);
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .quickstart-actions {
    margin-top: 0.9rem;
    display: flex;
    gap: 0.5rem;
  }

  .quickstart-btn {
    flex: 0 0 auto;
    width: auto;
    padding-inline: 1rem;
  }

  .quickstart-modal-backdrop {
    position: fixed;
    inset: 0;
    z-index: 60;
    display: grid;
    place-items: center;
    padding: 1rem;
    background: rgba(3, 7, 18, 0.62);
    backdrop-filter: blur(2px);
  }

  .quickstart-modal {
    width: min(26rem, calc(100vw - 2rem));
    border-color: color-mix(in srgb, var(--accent) 26%, var(--border-card));
    box-shadow: var(--shadow-overlay);
  }

  .quickstart-modal-title {
    margin: 0;
    color: var(--accent);
    font-size: 0.95rem;
  }

  .quickstart-modal-list {
    margin: 0.65rem 0 0;
    padding-left: 1.1rem;
    color: var(--text-secondary);
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .lesson-row {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.6rem;
  }

  .lesson-current-label {
    margin: 0;
    color: var(--text-secondary);
    font-weight: 600;
    font-size: var(--text-sm);
    line-height: var(--leading-snug);
    white-space: nowrap;
  }

  /* Both selects use the shared control vocabulary (`.select`); the row and the
     grid below supply the flex/grid sizing they need to share a line. */
  .lesson-row .select {
    flex: 1;
    min-width: 0;
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
    line-height: var(--leading-snug);
    white-space: nowrap;
  }

  .lesson-char-select {
    min-width: 0;
  }

  /* Transient onboarding emphasis: an accent edge, no shadow. Static surfaces
     are flat; elevation is reserved for things that stack above the page. */
  .lesson-char-highlight {
    border-color: color-mix(in srgb, var(--accent) 58%, var(--border-card));
  }

  @media (max-width: 767px) {
    .quickstart-actions {
      flex-wrap: wrap;
    }

    .quickstart-btn {
      flex: 1 1 auto;
      justify-content: center;
    }

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

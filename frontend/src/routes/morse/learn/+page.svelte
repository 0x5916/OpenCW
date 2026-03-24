<script lang="ts">
  import { browser } from '$app/environment';
  import { generateTimedLesson, LESSONS } from '$lib/morse';
  import MorsePlayer from '$lib/components/MorsePlayer.svelte';
  import ResultOverlay from '$lib/components/ResultOverlay.svelte';
  import { untrack, onDestroy } from 'svelte';
  import { ClipboardCheck } from 'lucide-svelte';
  import { CW_STORAGE_KEYS } from '$lib/storageKeys';
  import { langPreference } from '$lib/i18n.svelte';
  import { score, diffWords } from '$lib/score';
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
  import { localizeHref } from '$lib/paraglide/runtime';
  import * as m from '$lib/paraglide/messages';

  const QUICKSTART_DISMISSED_KEY = 'learn.quickstart.dismissed';

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
    showQuickStart = localStorage.getItem(QUICKSTART_DISMISSED_KEY) !== '1';
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
    localStorage.setItem(QUICKSTART_DISMISSED_KEY, '1');
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
    document.cookie = `${CW_STORAGE_KEYS.lesson}=${val}; path=/; max-age=31536000; SameSite=Lax`;
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

  let hasNextLesson = $derived(result >= 0.9 && chosenLesson < LESSONS.length);
  let hasPrevLesson = $derived(result < 0.7 && chosenLesson > 1);
  let shouldRegenerate = $derived(result >= 0.7 && result < 0.9);

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
  <p class="page-title-sub">
    <span class="accent-text">{m.trainer_subtitle_pre()}</span>
    {m.trainer_subtitle_post()}
  </p>
</header>

{#if showQuickStart}
  <section class="card-sm quickstart-card" aria-label="Quick start">
    <h2 class="quickstart-title">How to Train</h2>
    <ol class="quickstart-steps">
      <li>Pick a lesson or single letter</li>
      <li>Press ▶ to hear Morse code</li>
      <li>Type what you hear → Check Result</li>
    </ol>
    <div class="quickstart-actions">
      <button type="button" class="btn-success quickstart-btn" onclick={dismissQuickStart}
        >Start Training</button
      >
      <button type="button" class="btn-ghost quickstart-btn" onclick={openQuickTips}>Tips</button>
    </div>

    {#if showQuickTips}
      <div class="quickstart-modal-backdrop" role="presentation" onclick={closeQuickTips}>
        <div
          class="quickstart-modal card-sm"
          role="dialog"
          aria-modal="true"
          aria-label="Training tips"
          tabindex="-1"
          onclick={(event) => event.stopPropagation()}
          onkeydown={(event) => {
            if (event.key === 'Escape') closeQuickTips();
          }}
        >
          <h3 class="quickstart-modal-title">Tips</h3>
          <ul class="quickstart-modal-list">
            <li>Start with slower WPM and increase gradually.</li>
            <li>Use the single-letter player to isolate difficult characters.</li>
            <li>Check result often and focus on repeated mistakes.</li>
            {#if !$user}
              <li>
                {m.trainer_guest_notice()}
                <a href={localizeHref('/login')} class="link">{m.nav_login()}</a>
                /
                <a href={localizeHref('/register')} class="link">{m.nav_register()}</a>
              </li>
            {/if}
          </ul>
          <div class="quickstart-actions">
            <button type="button" class="btn-ghost quickstart-btn" onclick={closeQuickTips}
              >Close</button
            >
          </div>
        </div>
      </div>
    {/if}
  </section>
{/if}

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

    <div class="card-sm" class:lesson-char-highlight={showQuickStart && chosenLesson === 1}>
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

  .quickstart-card {
    margin-bottom: 1rem;
    border-color: color-mix(in srgb, var(--accent) 28%, var(--border));
    background:
      linear-gradient(170deg, color-mix(in srgb, var(--accent) 8%, transparent), transparent 45%),
      var(--bg-surface);
  }

  .quickstart-title {
    margin: 0;
    color: var(--accent);
    font-size: 1rem;
    line-height: 1.4;
    letter-spacing: 0.01em;
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
    border-color: color-mix(in srgb, var(--accent) 26%, var(--border));
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
    transition:
      background-color 0.2s,
      border-color 0.2s,
      color 0.2s;
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

  .lesson-char-highlight {
    border-color: color-mix(in srgb, var(--accent) 58%, var(--border));
    box-shadow:
      var(--shadow-soft),
      0 0 0 1px color-mix(in srgb, var(--accent) 28%, transparent),
      0 0 0 4px color-mix(in srgb, var(--accent) 10%, transparent);
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
    transition:
      background-color var(--transition-base),
      border-color var(--transition-base),
      box-shadow var(--transition-base);
  }

  .answer-card:hover {
    box-shadow: var(--shadow-lift);
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

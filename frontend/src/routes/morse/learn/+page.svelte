<script lang="ts">
  import { generateTimedLesson, LESSONS } from '$lib/morse';
  import MorsePlayer from '$lib/components/MorsePlayer.svelte';
  import ResultOverlay from '$lib/components/ResultOverlay.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import { untrack, onDestroy } from 'svelte';
  import { RefreshCw, Play, Square, ClipboardCheck, Upload, Download, Info } from 'lucide-svelte';
  import { score, diffWords } from '$lib/score';
  import type { DiffToken } from '$lib/score';
  import { user } from '$lib/auth';
  import { submitProgress, getCWSettings } from '$lib/api';
  import { performSync, performRestore, type CWSettings } from '$lib/cwSync';
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
  let startDelay = $state(0.5);

  let syncing = $state(false);
  let synced = $state(false);
  let syncError = $state('');
  let restoring = $state(false);
  let restoreError = $state('');

  let syncedTimeout: ReturnType<typeof setTimeout> | null = null;

  $effect(() => {
    if ($user) {
      getCWSettings()
        .then((cwSettings) => {
          charWpm = cwSettings.char_wpm;
          effWpm = cwSettings.eff_wpm;
          freq = cwSettings.freq;
          startDelay = cwSettings.start_delay;
        })
        .catch(() => {
          // Error silently handled, won't block UI
        });
    }
  });

  async function syncToAccount() {
    syncing = true;
    synced = false;
    syncError = '';
    if (syncedTimeout) {
      clearTimeout(syncedTimeout);
      syncedTimeout = null;
    }
    await performSync(
      { char_wpm: charWpm, eff_wpm: effWpm, freq, start_delay: startDelay },
      () => {
        synced = true;
        if (syncedTimeout) clearTimeout(syncedTimeout);
        syncedTimeout = setTimeout(() => {
          synced = false;
          syncedTimeout = null;
        }, 3000);
      },
      (msg: string) => {
        syncError = msg;
        setTimeout(() => (syncError = ''), 4000);
      }
    );
    syncing = false;
  }

  async function restoreFromAccount() {
    restoring = true;
    restoreError = '';
    await performRestore(
      (cwSettings: CWSettings) => {
        charWpm = cwSettings.char_wpm;
        effWpm = cwSettings.eff_wpm;
        freq = cwSettings.freq;
        startDelay = cwSettings.start_delay;
      },
      (msg: string) => {
        restoreError = msg;
        setTimeout(() => (restoreError = ''), 4000);
      }
    );
    restoring = false;
  }

  $effect(() => {
    const val = String(chosenLesson);
    localStorage.setItem('learn.lesson', val);
    document.cookie = `learn.lesson=${val}; path=/; max-age=31536000; SameSite=Lax`;
  });

  let lessonText = $derived(generateTimedLesson(chosenLesson, 60, charWpm, effWpm));
  let currentLessonWord = $derived(LESSONS.slice(0, chosenLesson).join(''));
  let currentLessonChars = $derived(currentLessonWord.split('').filter(Boolean));
  let selectedLessonChar = $state(LESSONS[0]?.[0] ?? '');
  let selectedCharPlayer = $state<{
    playNow: () => Promise<void>;
    stopNow: () => Promise<void>;
  } | null>(null);
  let fullLessonPlayer = $state<{
    playNow: () => Promise<void>;
    stopNow: () => Promise<void>;
    isStarted: () => boolean;
  } | null>(null);
  let charPlaying = $state(false);

  async function toggleSelectedLessonChar() {
    if (charPlaying) {
      await selectedCharPlayer?.stopNow();
      charPlaying = false;
    } else {
      charPlaying = true;
      await selectedCharPlayer?.playNow();
    }
  }

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
      submitProgress(String(chosenLesson), charWpm, effWpm, result).catch(() => {});
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

  onDestroy(() => {
    if (syncedTimeout) clearTimeout(syncedTimeout);
  });

  function onSelectedCharChange(event: Event) {
    selectedLessonChar = (event.currentTarget as HTMLSelectElement).value;
  }

  function onAnswerInput(event: Event) {
    inputText = (event.currentTarget as HTMLTextAreaElement).value.toUpperCase();
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
        <select bind:value={chosenLesson} onchange={() => (result = -1)} class="lesson-select">
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
      <div class="lesson-char-play-row">
        <button type="button" class="btn-regen" onclick={toggleSelectedLessonChar}>
          {#if charPlaying}
            <Square size={16} />{m.player_stop()}
          {:else}
            <Play size={16} />{m.trainer_play_letter()}
          {/if}
        </button>
      </div>
      <MorsePlayer
        bind:this={selectedCharPlayer}
        text={selectedLessonChar}
        {charWpm}
        {effWpm}
        {freq}
        compact
        hideControls
        playLabel={currentLessonChars.length > 1
          ? `Play "${selectedLessonChar}"`
          : m.trainer_play_letter()}
        repeat={3}
        onEnded={() => (charPlaying = false)}
      />
    </div>

    <div class="card-sm">
      <h2 class="card-label">{m.trainer_label_settings()}</h2>
      <div class="settings-grid settings-grid-2">
        <label class="settings-field">
          <span class="label-text">
            {m.trainer_label_char_wpm()}
            <span class="tooltip-icon" data-tooltip={m.trainer_tooltip_char_wpm()}><Info size={11} /></span>
          </span>
          <input type="number" bind:value={charWpm} min="5" max="50" class="input" />
        </label>
        <label class="settings-field">
          <span class="label-text">
            {m.trainer_label_eff_wpm()}
            <span class="tooltip-icon" data-tooltip={m.trainer_tooltip_eff_wpm()}><Info size={11} /></span>
          </span>
          <input type="number" bind:value={effWpm} min="5" max="50" class="input" />
        </label>
      </div>
      <details class="settings-adv">
        <summary class="settings-adv-toggle">Advanced</summary>
        <div class="settings-grid settings-grid-2">
          <label class="settings-field">
            <span class="label-text">{m.trainer_label_freq()}</span>
            <input type="number" bind:value={freq} min="100" max="2000" class="input" />
          </label>
          <label class="settings-field">
            <span class="label-text">{m.trainer_label_start_delay()}</span>
            <input
              type="number"
              bind:value={startDelay}
              min="0"
              max="10"
              step="0.5"
              class="input"
            />
          </label>
        </div>
      </details>
      {#if $user}
        <div class="lesson-char-play-row">
          <button
            type="button"
            class="btn-regen"
            class:btn-regen-success={synced}
            onclick={syncToAccount}
            disabled={syncing}
          >
            <Upload size={14} />{synced
              ? m.trainer_synced()
              : syncing
                ? m.trainer_syncing()
                : m.trainer_sync_settings()}
          </button>
          <button type="button" class="btn-regen" onclick={restoreFromAccount} disabled={restoring}>
            <Download size={14} />{restoring ? m.trainer_restoring() : m.trainer_restore_settings()}
          </button>
        </div>
        {#if syncError}
          <ErrorAlert message={syncError} />
        {:else if restoreError}
          <ErrorAlert message={restoreError} />
        {/if}
      {:else}
        <p class="body-text">
          {m.trainer_guest_notice()}
          <a href="/login" class="link">{m.nav_login()}</a>
          /
          <a href="/register" class="link">{m.nav_register()}</a>
        </p>
      {/if}
    </div>
  </div>

  <!-- Right column: answer + result -->
  <div class="learn-col-right">
    <div class="card-sm">
      <MorsePlayer
        bind:this={fullLessonPlayer}
        text={lessonText}
        {charWpm}
        {effWpm}
        {freq}
        {startDelay}
        label={m.player_label()}
      />
      <button onclick={regenerate} class="btn-regen"
        ><RefreshCw size={16} />{m.trainer_regenerate()}</button
      >
    </div>

    <div class="answer-card learn-answer-card">
      <h2 class="card-label">{m.trainer_answer_label()}</h2>
      <textarea
        placeholder={`${m.trainer_answer_placeholder()}\n${m.trainer_answer_shortcut_tip()}`}
        bind:value={inputText}
        oninput={onAnswerInput}
        onkeydown={onAnswerKeydown}
        autocapitalize="characters"
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

<script lang="ts">
  import { calculateDuration, getFarnsworthWpmSet, MORSE } from '$lib/morse';
  import { onDestroy } from 'svelte';
  import { Play, Square, Pause, SlidersHorizontal, SkipForward } from 'lucide-svelte';
  import * as m from '$lib/paraglide/messages';

  let {
    text = '',
    charWpm = $bindable(20),
    effWpm = $bindable(10),
    freq = $bindable(600),
    startDelay = 0,
    volume = $bindable(1.0),
    compact = false,
    playLabel = '',
    hideControls = false,
    label = '',
    showSettings = false,
    mediaStyle = false,
    onSettingsInput = () => {},
    onEnded = () => {}
  } = $props();

  let effectivePlayLabel = $derived(playLabel || m.player_play());

  let ctx = $state<AudioContext | null>(null);
  let started = $state(false);
  let paused = $state(false);
  let progress = $state(0);
  let timer = $state(0);
  let rafId = 0;
  let ctxStartTime = 0;
  let activeText = $state('');

  let duration = $derived(calculateDuration(activeText, charWpm, effWpm) + startDelay);
  let timings = $derived(getFarnsworthWpmSet(charWpm, effWpm));
  let fade = $derived(Math.min(timings.charDot * 0.1, 0.005));
  let elapsedTime = $derived(Math.max(timer, 0));
  let totalTime = $derived(Math.max(duration - startDelay, 0));

  function formatClock(seconds: number) {
    const safe = Math.max(seconds, 0);
    const min = Math.floor(safe / 60);
    const sec = Math.floor(safe % 60)
      .toString()
      .padStart(2, '0');
    return `${min}:${sec}`;
  }

  $effect(() => {
    if (!started) activeText = text;
  });

  function play() {
    if (ctx !== null) throw new Error('Audio is already playing');
    started = true;
    paused = false;

    ctx = new AudioContext();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.connect(gain);
    gain.connect(ctx.destination);
    osc.frequency.value = freq;
    gain.gain.setValueAtTime(0, ctx.currentTime);
    ctxStartTime = ctx.currentTime;
    let t = ctx.currentTime + startDelay;

    for (const char of activeText.toUpperCase()) {
      const morse = MORSE[char];
      if (!morse) {
        t += timings.wordSpace;
        continue;
      }

      for (let j = 0; j < morse.length; j++) {
        const dur = morse[j] === '.' ? timings.charDot : timings.dash;
        gain.gain.setValueAtTime(0, t);
        gain.gain.linearRampToValueAtTime(volume, t + fade);
        gain.gain.setValueAtTime(volume, t + dur - fade);
        gain.gain.linearRampToValueAtTime(0, t + dur);
        t += dur;
        if (j < morse.length - 1) t += timings.symbolSpace;
      }
      t += timings.letterSpace;
    }
    // remove last letter space
    t -= activeText.length > 0 ? timings.letterSpace : 0;

    osc.onended = () => stop();
    osc.start();
    osc.stop(t);

    function tick() {
      if (!ctx) return;
      const elapsed = ctx.currentTime - ctxStartTime;
      progress = Math.min(elapsed / duration, 1);
      timer = elapsed - startDelay;
      rafId = requestAnimationFrame(tick);
    }
    rafId = requestAnimationFrame(tick);
  }

  async function pause() {
    if (!ctx) return;
    await ctx.suspend();
    paused = true;
  }

  async function resume() {
    if (!ctx) return;
    await ctx.resume();
    paused = false;
  }

  async function stop(emitEnded = true) {
    cancelAnimationFrame(rafId);
    progress = 0;
    timer = 0;
    await ctx?.close();
    ctx = null;
    started = false;
    paused = false;
    if (emitEnded) onEnded();
  }

  async function togglePlayPause() {
    if (!started) {
      activeText = text;
      play();
      return;
    }
    if (paused) {
      await resume();
    } else {
      await pause();
    }
  }

  async function skipNext() {
    if (!started) return;

    const source = activeText.toUpperCase();
    if (!source) return;

    const ratio = totalTime > 0 ? elapsedTime / totalTime : 1;
    const nextIndex = Math.min(source.length, Math.floor(ratio * source.length) + 1);
    const remaining = source.slice(nextIndex);

    await stop(false);
    if (!remaining) {
      activeText = text;
      return;
    }

    activeText = remaining;
    play();
  }

  export async function playNow() {
    if (!text) return;
    if (ctx) await stop();
    play();
  }

  export async function stopNow() {
    await stop();
  }

  export function isStarted() {
    return started;
  }

  onDestroy(() => ctx?.close());
</script>

{#if !hideControls}
  <div class="player-wrapper">
    {#if label || !compact}
      <div class="player-header">
        {#if label}<p class="card-label player-label">{label}</p>{/if}
      </div>
    {/if}
    <div class="player-top">
      <div class="player-buttons">
        {#if mediaStyle}
          <div class="player-media-controls">
            <button
              type="button"
              class="player-icon-btn btn-ghost"
              onclick={() => stop()}
              disabled={!started}
              aria-label={m.player_stop()}
            >
              <Square size={16} />
            </button>
            <button
              type="button"
              class="player-icon-btn btn-success"
              onclick={togglePlayPause}
              aria-label={started && !paused ? m.player_pause() : effectivePlayLabel}
            >
              {#if started && !paused}
                <Pause size={16} />
              {:else}
                <Play size={16} />
              {/if}
            </button>
            <button
              type="button"
              class="player-icon-btn btn-ghost"
              onclick={skipNext}
              disabled={!started}
              aria-label="Skip to next character"
            >
              <SkipForward size={16} />
            </button>
            <span class="player-inline-timer" class:player-timer-delay={timer < 0 && started}>
              {formatClock(elapsedTime)} / {formatClock(totalTime)}
            </span>
          </div>
        {:else}
          <button onclick={play} class="btn-success" class:hidden={started}
            ><Play size={16} />{effectivePlayLabel}</button
          >
          <button onclick={() => stop()} class="btn-danger" class:hidden={!started}
            ><Square size={16} />{m.player_stop()}</button
          >
          {#if !compact}
            <button onclick={resume} class="btn-success" class:hidden={!paused}
              ><Play size={16} />{m.player_resume()}</button
            >
            <button onclick={pause} class="btn-ghost" class:hidden={paused} disabled={!started}
              ><Pause size={16} />{m.player_pause()}</button
            >
          {/if}
        {/if}
        {#if showSettings}
          <details class="player-settings-menu">
            <summary class="player-settings-trigger" aria-label={m.trainer_label_settings()}>
              <SlidersHorizontal size={14} />
            </summary>
            <div class="player-settings-popover card-sm">
              <label class="player-settings-field">
                <span class="label-text">{m.trainer_label_char_wpm()}</span>
                <input
                  type="number"
                  bind:value={charWpm}
                  min="5"
                  max="50"
                  class="input"
                  oninput={() => onSettingsInput()}
                />
              </label>
              <label class="player-settings-field">
                <span class="label-text">{m.trainer_label_eff_wpm()}</span>
                <input
                  type="number"
                  bind:value={effWpm}
                  min="5"
                  max="50"
                  class="input"
                  oninput={() => onSettingsInput()}
                />
              </label>
              <label class="player-settings-field">
                <span class="label-text">{m.trainer_label_freq()}</span>
                <input
                  type="number"
                  bind:value={freq}
                  min="100"
                  max="2000"
                  class="input"
                  oninput={() => onSettingsInput()}
                />
              </label>
              <label class="player-settings-field">
                <span class="label-text">Volume</span>
                <div class="player-volume-row">
                  <input
                    type="range"
                    bind:value={volume}
                    min="0"
                    max="1"
                    step="0.05"
                    class="player-volume-slider"
                    oninput={() => onSettingsInput()}
                  />
                  <span class="player-volume-value">{Math.round(volume * 100)}%</span>
                </div>
              </label>
              {#if mediaStyle}
                <label class="player-settings-field">
                  <span class="label-text">{m.trainer_label_start_delay()}</span>
                  <input
                    type="number"
                    bind:value={startDelay}
                    min="0"
                    max="10"
                    step="0.5"
                    class="input"
                    oninput={() => onSettingsInput()}
                  />
                </label>
              {/if}
              <p class="player-settings-hint">Visit Settings to save permanently</p>
            </div>
          </details>
        {/if}
      </div>
      {#if !mediaStyle && started}
        <div class="player-progress-row">
          <div
            class="player-progress"
            role="progressbar"
            aria-valuenow={Math.round(progress * 100)}
            aria-valuemin={0}
            aria-valuemax={100}
          >
            <div class="player-progress-fill" style="width: {progress * 100}%"></div>
          </div>
          <span class="player-timer" class:player-timer-delay={timer < 0}>
            {timer.toFixed(1)}s / {(duration - startDelay).toFixed(1)}s
          </span>
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .player-wrapper {
    --player-control-height: 2.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .player-header {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 0.5rem;
  }

  .player-label {
    margin-bottom: 0;
  }

  .player-top {
    display: flex;
    flex-direction: column;
    align-items: stretch;
    gap: 0.75rem;
  }

  .player-progress-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .player-progress {
    flex: 1;
    height: 4px;
    background: var(--border);
    border-radius: 2px;
    overflow: hidden;
  }

  .player-progress-fill {
    height: 100%;
    background: var(--accent);
    border-radius: 2px;
  }

  .player-timer {
    font-size: 0.7rem;
    font-variant-numeric: tabular-nums;
    color: var(--accent);
    white-space: nowrap;
    min-width: 9ch;
    text-align: right;
  }

  .player-timer-delay {
    color: var(--text-muted);
  }

  .player-buttons {
    display: flex;
    gap: 0.75rem;
    align-items: stretch;
    position: relative;
  }

  .player-buttons .hidden {
    display: none;
  }

  .player-media-controls {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex: 1;
  }

  .player-inline-timer {
    margin-left: auto;
    font-size: 0.78rem;
    font-variant-numeric: tabular-nums;
    color: var(--accent);
    white-space: nowrap;
    text-align: right;
  }

  .player-buttons > button {
    height: var(--player-control-height);
    line-height: 1;
    padding-top: 0;
    padding-bottom: 0;
  }

  .player-icon-btn {
    width: var(--player-control-height);
    height: var(--player-control-height);
    min-width: var(--player-control-height);
    min-height: var(--player-control-height);
    padding: 0;
    flex: 0 0 auto;
  }

  .player-settings-menu {
    margin-left: auto;
    position: relative;
    display: flex;
  }

  .player-settings-trigger {
    list-style: none;
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: center;
    width: var(--player-control-height);
    height: var(--player-control-height);
    line-height: 1;
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    background: var(--bg-inset);
    color: var(--text-secondary);
    cursor: pointer;
    transition: border-color var(--transition-fast), color var(--transition-fast);
  }

  .player-settings-trigger::-webkit-details-marker {
    display: none;
  }

  .player-settings-trigger:hover {
    border-color: color-mix(in srgb, var(--accent) 45%, var(--border-subtle));
    color: var(--text-primary);
  }

  .player-settings-menu[open] .player-settings-trigger {
    color: var(--accent);
    border-color: color-mix(in srgb, var(--accent) 55%, var(--border-subtle));
  }

  .player-settings-popover {
    position: absolute;
    right: 0;
    top: calc(100% + 0.45rem);
    z-index: 20;
    width: min(16rem, 75vw);
    padding: 0.8rem;
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
  }

  .player-settings-field {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .player-settings-hint {
    margin: 0.1rem 0 0;
    font-size: 0.72rem;
    color: var(--text-muted);
  }

  .player-volume-row {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    align-items: center;
    gap: 0.5rem;
  }

  .player-volume-slider {
    width: 100%;
    accent-color: var(--accent);
  }

  .player-volume-value {
    font-size: 0.75rem;
    color: var(--text-muted);
    min-width: 3ch;
    text-align: right;
  }

  @media (max-width: 900px) {
    .player-top {
      gap: 0.45rem;
    }
  }

  @media (max-width: 420px) {
    .player-buttons {
      flex-wrap: nowrap;
      overflow-x: auto;
      scrollbar-width: thin;
      padding-bottom: 0.2rem;
    }

    .player-buttons > * {
      flex: 0 0 auto;
    }

    .player-buttons > button {
      min-width: max-content;
    }

    .player-media-controls {
      min-width: max-content;
    }
  }
</style>

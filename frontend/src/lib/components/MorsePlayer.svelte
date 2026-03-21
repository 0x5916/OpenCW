<script lang="ts">
  import { calculateDuration, getFarnsworthWpmSet, MORSE } from '$lib/morse';
  import { onDestroy } from 'svelte';
  import { Play, Square, Pause } from 'lucide-svelte';
  import * as m from '$lib/paraglide/messages';

  let {
    text = '',
    charWpm = 20,
    effWpm = 10,
    freq = 600,
    startDelay = 0,
    volume = 1.0,
    compact = false,
    playLabel = '',
    hideControls = false,
    label = '',
    repeat = 1,
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

  let duration = $derived(calculateDuration(text, charWpm, effWpm) + startDelay);
  let timings = $derived(getFarnsworthWpmSet(charWpm, effWpm));
  let fade = $derived(Math.min(timings.charDot * 0.1, 0.005));

  function play() {
    if (ctx !== null) throw new Error('Audio is already playing');
    started = true;

    ctx = new AudioContext();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.connect(gain);
    gain.connect(ctx.destination);
    osc.frequency.value = freq;
    gain.gain.setValueAtTime(0, ctx.currentTime);
    ctxStartTime = ctx.currentTime;
    let t = ctx.currentTime + startDelay;

    for (let i = 0; i < repeat; i++) {
      for (const char of text.toUpperCase()) {
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
    }
    // remove last letter space
    t -= text.length > 0 ? timings.letterSpace : 0;

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

  async function stop() {
    cancelAnimationFrame(rafId);
    progress = 0;
    timer = 0;
    await ctx?.close();
    ctx = null;
    started = false;
    paused = false;
    onEnded();
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
        <button onclick={play} class="btn-success" class:hidden={started}
          ><Play size={16} />{effectivePlayLabel}</button
        >
        <button onclick={stop} class="btn-danger" class:hidden={!started}
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
      </div>
      {#if started}
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
  }

  .player-buttons .hidden {
    display: none;
  }

  @media (max-width: 900px) {
    .player-top {
      gap: 0.45rem;
    }
  }

  @media (max-width: 420px) {
    .player-buttons {
      flex-wrap: wrap;
    }

    .player-buttons > * {
      flex: 1 1 100%;
    }
  }
</style>

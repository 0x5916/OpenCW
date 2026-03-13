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
    startDelay = 0.5,
    volume = 1.0,
    compact = false,
    playLabel = '',
    hideControls = false,
    label = '',
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

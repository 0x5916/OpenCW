<script lang="ts">
  import { user } from '$lib/auth';
  import { getUserInfo, getCWSettings, getProgress } from '$lib/api';
  import type { ProgressRecord } from '$lib/api';
  import { LESSONS } from '$lib/morse';
  import {
    User,
    Radio,
    Calendar,
    Trophy,
    TrendingUp,
    Activity,
    Zap,
    BookOpen
  } from 'lucide-svelte';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import * as m from '$lib/paraglide/messages';

  let loading = $state(true);
  let loadError = $state('');

  let username = $state('');
  let email = $state('');
  let memberSince = $state('');
  let charWpm = $state(0);
  let effWpm = $state(0);
  let freq = $state(0);
  let records = $state<ProgressRecord[]>([]);

  // Derived stats
  let totalSessions = $derived(records.length);
  let bestAccuracy = $derived(records.length ? Math.max(...records.map((r) => r.accuracy)) : 0);
  let avgAccuracy = $derived(
    records.length ? records.reduce((s, r) => s + r.accuracy, 0) / records.length : 0
  );
  let uniqueLessons = $derived(new Set(records.map((r) => r.lesson)).size);
  let recentRecords = $derived([...records].reverse().slice(0, 20));

  $effect(() => {
    if ($user) loadAll();
    else loading = false;
  });

  async function loadAll() {
    loading = true;
    loadError = '';
    try {
      const [info, cw, prog] = await Promise.all([getUserInfo(), getCWSettings(), getProgress()]);
      username = info.username;
      email = info.email;
      memberSince = new Date(info.created_at).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
      charWpm = cw.char_wpm;
      effWpm = cw.eff_wpm;
      freq = cw.freq;
      records = prog;
    } catch (e) {
      loadError = e instanceof Error ? e.message : 'Failed to load profile';
    } finally {
      loading = false;
    }
  }

  function formatLesson(lessonStr: string): string {
    let cumulative = '';
    for (let i = 0; i < LESSONS.length; i++) {
      cumulative += LESSONS[i];
      if (cumulative === lessonStr.toUpperCase()) {
        return `${i + 1} — ${LESSONS[i].split('').join(', ')}`;
      }
    }
    // Fallback: truncate raw string
    return lessonStr.length > 12 ? lessonStr.slice(0, 12) + '…' : lessonStr;
  }

  function pct(v: number) {
    return Math.round(v * 100) + '%';
  }

  function accuracyClass(v: number) {
    return v >= 0.9 ? 'acc-good' : v >= 0.7 ? 'acc-ok' : 'acc-bad';
  }

  function formatDate(iso: string) {
    return new Date(iso).toLocaleDateString(undefined, {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  }
</script>

<main class="profile-page">
  {#if !$user && !loading}
    <div class="card">
      <p class="body-text">
        {m.profile_not_logged_in()} <a href="/login" class="link">{m.nav_login()}</a>
      </p>
    </div>
  {:else if loading}
    <LoadingSpinner variant="spinner" />
  {:else if loadError}
    <ErrorAlert message={loadError} />
  {:else}
    <!-- Header -->
    <header class="profile-header card">
      <div class="profile-avatar">
        <User size={48} />
      </div>
      <div class="profile-header-info">
        <h1 class="profile-username">{username}</h1>
        <p class="profile-email">{email}</p>
        <p class="profile-joined">
          <Calendar size={13} />
          {m.profile_member_since()}: {memberSince}
        </p>
      </div>
    </header>

    <!-- Stats row -->
    <section class="profile-stats-grid">
      <div class="profile-stat-card card">
        <Activity class="profile-stat-icon" aria-hidden="true" />
        <p class="profile-stat-value">{totalSessions}</p>
        <p class="profile-stat-label">{m.profile_stats_sessions()}</p>
      </div>
      <div class="profile-stat-card card">
        <Trophy class="profile-stat-icon profile-stat-icon--gold" aria-hidden="true" />
        <p class="profile-stat-value">{pct(bestAccuracy)}</p>
        <p class="profile-stat-label">{m.profile_stats_best()}</p>
      </div>
      <div class="profile-stat-card card">
        <TrendingUp class="profile-stat-icon profile-stat-icon--blue" aria-hidden="true" />
        <p class="profile-stat-value">{pct(avgAccuracy)}</p>
        <p class="profile-stat-label">{m.profile_stats_avg()}</p>
      </div>
      <div class="profile-stat-card card">
        <BookOpen class="profile-stat-icon profile-stat-icon--green" aria-hidden="true" />
        <p class="profile-stat-value">{uniqueLessons}</p>
        <p class="profile-stat-label">{m.profile_stats_lessons()}</p>
      </div>
    </section>

    <!-- CW Settings snapshot -->
    <section class="card profile-cw-card">
      <div class="profile-section-header">
        <Radio size={16} />
        <h2 class="card-label">{m.profile_cw_settings()}</h2>
      </div>
      <div class="profile-cw-grid">
        <div class="profile-cw-item">
          <Zap size={14} />
          <span class="profile-cw-val">{charWpm}</span>
          <span class="profile-cw-key">{m.profile_cw_char_wpm()}</span>
        </div>
        <div class="profile-cw-item">
          <Zap size={14} />
          <span class="profile-cw-val">{effWpm}</span>
          <span class="profile-cw-key">{m.profile_cw_eff_wpm()}</span>
        </div>
        <div class="profile-cw-item">
          <Zap size={14} />
          <span class="profile-cw-val">{freq} Hz</span>
          <span class="profile-cw-key">{m.profile_cw_freq()}</span>
        </div>
      </div>
    </section>

    <!-- Progress history -->
    <section class="card profile-history-card">
      <div class="profile-section-header">
        <Activity size={16} />
        <h2 class="card-label">{m.profile_history()}</h2>
      </div>
      {#if recentRecords.length === 0}
        <p class="body-text profile-empty">
          {m.profile_history_empty()} <a href="/morse/learn" class="link">{m.nav_learn()}</a>
        </p>
      {:else}
        <div class="profile-table-wrap">
          <table class="profile-table">
            <thead>
              <tr>
                <th>{m.profile_history_lesson()}</th>
                <th>{m.profile_history_accuracy()}</th>
                <th>{m.profile_history_wpm()}</th>
                <th>{m.profile_history_date()}</th>
              </tr>
            </thead>
            <tbody>
              {#each recentRecords as rec (rec.created_at)}
                <tr>
                  <td class="profile-lesson-cell">{formatLesson(rec.lesson)}</td>
                  <td
                    ><span class="profile-acc {accuracyClass(rec.accuracy)}"
                      >{pct(rec.accuracy)}</span
                    ></td
                  >
                  <td class="profile-wpm-cell">{rec.char_wpm} / {rec.eff_wpm}</td>
                  <td class="profile-date-cell">{formatDate(rec.created_at)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </section>
  {/if}
</main>

<style>
  .profile-page {
    max-width: 860px;
    margin: 0 auto;
    padding: 2rem 1rem;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  /* Header */
  .profile-header {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }
  .profile-avatar {
    width: 5rem;
    height: 5rem;
    border-radius: 50%;
    background-color: var(--bg-inset);
    border: 2px solid var(--border);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--accent);
    flex-shrink: 0;
  }
  .profile-header-info {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  .profile-username {
    font-size: 1.75rem;
    font-weight: 800;
    color: var(--text-primary);
    margin: 0;
  }
  .profile-email {
    font-size: 0.875rem;
    color: var(--text-muted);
    margin: 0;
  }
  .profile-joined {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.8rem;
    color: var(--text-muted);
    margin: 0;
  }

  /* Stats */
  .profile-stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
  }
  .profile-stat-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.25rem;
    padding: 1.25rem 0.75rem;
    text-align: center;
  }
  :global(.profile-stat-icon) {
    color: var(--text-muted);
    margin-bottom: 0.25rem;
  }
  :global(.profile-stat-icon--gold) {
    color: #f59e0b;
  }
  :global(.profile-stat-icon--blue) {
    color: #3b82f6;
  }
  :global(.profile-stat-icon--green) {
    color: #22c55e;
  }
  .profile-stat-value {
    font-size: 1.5rem;
    font-weight: 800;
    color: var(--text-primary);
    margin: 0;
  }
  .profile-stat-label {
    font-size: 0.7rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--text-muted);
    margin: 0;
  }

  /* CW Settings */
  .profile-section-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.75rem;
    color: var(--text-muted);
  }
  .profile-section-header :global(.card-label) {
    margin-bottom: 0;
  }
  .profile-cw-grid {
    display: flex;
    gap: 2rem;
  }
  .profile-cw-item {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--text-muted);
  }
  .profile-cw-val {
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--accent);
  }
  .profile-cw-key {
    font-size: 0.8rem;
    color: var(--text-muted);
  }

  /* History table */
  .profile-table-wrap {
    overflow-x: auto;
  }
  .profile-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.875rem;
  }
  .profile-table th {
    text-align: left;
    padding: 0.5rem 0.75rem;
    font-size: 0.7rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--text-muted);
    border-bottom: 1px solid var(--border);
  }
  .profile-table td {
    padding: 0.6rem 0.75rem;
    border-bottom: 1px solid var(--border-subtle);
    color: var(--text-secondary);
  }
  .profile-table tbody tr:last-child td {
    border-bottom: none;
  }
  .profile-table tbody tr:hover td {
    background-color: var(--bg-inset);
  }
  .profile-lesson-cell {
    font-family: monospace;
    font-size: 1rem;
    color: var(--text-primary) !important;
    letter-spacing: 0.15em;
  }
  .profile-acc {
    display: inline-block;
    padding: 0.15rem 0.5rem;
    border-radius: 999px;
    font-weight: 700;
    font-size: 0.8rem;
  }
  .acc-good {
    background: rgba(34, 197, 94, 0.15);
    color: #22c55e;
  }
  .acc-ok {
    background: rgba(234, 179, 8, 0.15);
    color: #ca8a04;
  }
  .acc-bad {
    background: rgba(239, 68, 68, 0.15);
    color: #ef4444;
  }
  .profile-wpm-cell {
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
  }
  .profile-date-cell {
    white-space: nowrap;
    font-size: 0.8rem;
  }
  .profile-empty {
    color: var(--text-muted);
    padding: 0.5rem 0;
  }

  /* Responsive */
  @media (max-width: 640px) {
    .profile-header {
      flex-direction: column;
      align-items: flex-start;
    }
    .profile-stats-grid {
      grid-template-columns: repeat(2, 1fr);
    }
    .profile-cw-grid {
      flex-direction: column;
      gap: 0.75rem;
    }
  }
</style>

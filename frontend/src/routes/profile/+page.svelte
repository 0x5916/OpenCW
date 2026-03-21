<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '$lib/auth';
  import { getUserInfo, getCWSettings, getProgress } from '$lib/api';
  import type { ProgressRecord } from '$lib/api';
  import { LESSONS } from '$lib/morse';
  import {
    User,
    Radio,
    Calendar,
    Activity,
    Zap,
    Check,
    X
  } from 'lucide-svelte';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import { localizeApiError } from '$lib/errorLocalization';
  import { localizeHref } from '$lib/paraglide/runtime';
  import * as m from '$lib/paraglide/messages';

  let loading = $state(true);
  let loadError = $state('');

  let username = $state('');
  let email = $state('');
  let callSign = $state<string | null>(null);
  let emailVerified = $state(false);
  let memberSince = $state('');
  let memberCreatedAtMs = $state<number | null>(null);
  let selectedYear = $state(new Date().getFullYear());
  let charWpm = $state(0);
  let effWpm = $state(0);
  let freq = $state(0);
  let records = $state<ProgressRecord[]>([]);
  let heatmapScrollEl = $state<HTMLDivElement | null>(null);

  const WEEKDAY_LABELS = ['Mon', '', 'Wed', '', 'Fri', '', ''] as const;
  const CURRENT_YEAR = new Date().getFullYear();
  const DAY_MS = 24 * 60 * 60 * 1000;
  const MONTH_FORMATTER = new Intl.DateTimeFormat(undefined, { month: 'short' });
  const DAY_OF_MONTH_UTC_FORMATTER = new Intl.DateTimeFormat('en-US', {
    day: 'numeric',
    timeZone: 'UTC'
  });
  const DAY_FORMATTER = new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });

  type HeatmapCell = {
    dateMs: number;
    dateKey: string;
    count: number;
    level: 0 | 1 | 2 | 3 | 4;
    outsideYear: boolean;
    beforeAccount: boolean;
    inFuture: boolean;
  };

  type HeatmapWeek = {
    startMs: number;
    monthLabel: string;
    cells: HeatmapCell[];
  };

  let availableYears = $derived(getAvailableYears(memberCreatedAtMs));
  let yearTotalSessions = $derived(countSessionsForYear(records, selectedYear));
  let heatmapWeeks = $derived(buildHeatmap(records, memberCreatedAtMs, selectedYear));
  let heatmapWeekCount = $derived(Math.max(1, heatmapWeeks.length));
  let todayDayMs = $derived(utcDayStart(Date.now()));
  let recentRecords = $derived([...records].reverse().slice(0, 20));

  $effect(() => {
    if ($user) loadAll();
    else loading = false;
  });

  function scrollHeatmapToEnd() {
    if (!heatmapScrollEl) return;
    heatmapScrollEl.scrollLeft = Math.max(0, heatmapScrollEl.scrollWidth - heatmapScrollEl.clientWidth);
  }

  function scrollHeatmapToCurrentWeek() {
    if (!heatmapScrollEl) return;

    const currentWeekEl = heatmapScrollEl.querySelector('.profile-heatmap-week.is-current-week') as
      | HTMLElement
      | null;
    const currentDayEl = heatmapScrollEl.querySelector('.profile-heatmap-cell.is-today') as
      | HTMLElement
      | null;
    const targetEl = currentWeekEl ?? currentDayEl;

    if (!targetEl) {
      scrollHeatmapToEnd();
      return;
    }

    const rightPadding = 8;
    const containerRect = heatmapScrollEl.getBoundingClientRect();
    const targetRect = targetEl.getBoundingClientRect();
    const target =
      heatmapScrollEl.scrollLeft + (targetRect.right - containerRect.right) + rightPadding;
    const maxScroll = Math.max(0, heatmapScrollEl.scrollWidth - heatmapScrollEl.clientWidth);
    heatmapScrollEl.scrollLeft = Math.max(0, Math.min(maxScroll, target));
  }

  onMount(() => {
    requestAnimationFrame(() => {
      if (selectedYear === CURRENT_YEAR) {
        scrollHeatmapToCurrentWeek();
      } else {
        if (!heatmapScrollEl) return;
        heatmapScrollEl.scrollLeft = 0;
      }
    });
  });

  $effect(() => {
    if (!availableYears.includes(selectedYear)) {
      selectedYear = availableYears[0] ?? new Date().getFullYear();
    }
  });

  $effect(() => {
    const weekCount = heatmapWeeks.length;
    if (weekCount === 0 || !heatmapScrollEl) return;

    requestAnimationFrame(() => {
      if (!heatmapScrollEl) return;
      if (selectedYear === CURRENT_YEAR) {
        scrollHeatmapToCurrentWeek();
      } else {
        heatmapScrollEl.scrollLeft = 0;
      }
    });
  });

  async function loadAll() {
    loading = true;
    loadError = '';
    try {
      const [info, cw, prog] = await Promise.all([getUserInfo(), getCWSettings(), getProgress()]);
      username = info.username;
      email = info.email;
      callSign = info.call_sign;
      emailVerified = info.email_verified;
      memberSince = new Date(info.created_at).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
      {
        const createdAtParsedMs = Date.parse(info.created_at);
        memberCreatedAtMs = Number.isNaN(createdAtParsedMs) ? null : utcDayStart(createdAtParsedMs);
      }
      selectedYear = new Date().getFullYear();
      charWpm = cw.char_wpm;
      effWpm = cw.eff_wpm;
      freq = cw.freq;
      records = prog;
    } catch (e) {
      loadError = localizeApiError(e, () => m.settings_load_error());
    } finally {
      loading = false;
    }
  }

  function formatLesson(lessonStr: string): string {
    const numericLesson = Number.parseInt(lessonStr, 10);
    if (Number.isInteger(numericLesson) && numericLesson >= 1 && numericLesson <= LESSONS.length) {
      return `${numericLesson} - ${LESSONS[numericLesson - 1].split('').join(', ')}`;
    }

    let cumulative = '';
    for (let i = 0; i < LESSONS.length; i++) {
      cumulative += LESSONS[i];
      if (cumulative === lessonStr.toUpperCase()) {
        return `${i + 1} - ${LESSONS[i].split('').join(', ')}`;
      }
    }
    // Fallback: truncate raw string
    return lessonStr.length > 12 ? lessonStr.slice(0, 12) + '…' : lessonStr;
  }

  function pct(v: number) {
    return Math.round(v * 100) + '%';
  }

  function utcDayStart(ms: number): number {
    return Math.floor(ms / DAY_MS) * DAY_MS;
  }

  function mondayIndexFromDayNumber(dayNumber: number): number {
    return ((dayNumber + 3) % 7 + 7) % 7;
  }

  function startOfWeekMondayMs(ms: number): number {
    const dayStart = utcDayStart(ms);
    const dayNumber = Math.floor(dayStart / DAY_MS);
    return dayStart - mondayIndexFromDayNumber(dayNumber) * DAY_MS;
  }

  function addDaysMs(ms: number, days: number): number {
    return ms + days * DAY_MS;
  }

  function toDateKey(ms: number): string {
    return String(ms);
  }

  function monthLabelForWeek(weekStartMs: number, todayMs: number): string {
    for (let offset = 0; offset < 7; offset += 1) {
      const dayMs = addDaysMs(weekStartMs, offset);
      if (dayMs > todayMs) break;
      if (DAY_OF_MONTH_UTC_FORMATTER.format(dayMs) === '1') {
        return MONTH_FORMATTER.format(dayMs);
      }
    }

    return '';
  }

  function contributionLevel(count: number, maxCount: number): 0 | 1 | 2 | 3 | 4 {
    if (count <= 0) return 0;
    if (maxCount <= 1) return 4;
    const scaled = Math.ceil((count / maxCount) * 4);
    return Math.min(4, Math.max(1, scaled)) as 1 | 2 | 3 | 4;
  }

  function getAvailableYears(createdAtMs: number | null): number[] {
    const currentYear = new Date().getFullYear();
    const startYear = createdAtMs === null ? currentYear : new Date(createdAtMs).getFullYear();

    const years: number[] = [];
    for (let year = currentYear; year >= startYear; year -= 1) {
      years.push(year);
    }

    return years;
  }

  function countSessionsForYear(items: ProgressRecord[], year: number): number {
    let count = 0;
    for (const rec of items) {
      const parsedMs = Date.parse(rec.created_at);
      if (Number.isNaN(parsedMs)) continue;
      if (new Date(parsedMs).getUTCFullYear() === year) {
        count += 1;
      }
    }
    return count;
  }

  function buildHeatmap(items: ProgressRecord[], createdAtMs: number | null, year: number): HeatmapWeek[] {
    const todayMs = utcDayStart(Date.now());
    const yearStartMs = utcDayStart(Date.UTC(year, 0, 1));
    const yearEndMs = utcDayStart(Date.UTC(year, 11, 31));
    const startWeekStartMs = startOfWeekMondayMs(yearStartMs);
    const endWeekStartMs = startOfWeekMondayMs(yearEndMs);
    const weekCount = Math.max(1, Math.floor((endWeekStartMs - startWeekStartMs) / (7 * DAY_MS)) + 1);

    const byDay: Record<string, number> = {};

    for (const rec of items) {
      const parsedMs = Date.parse(rec.created_at);
      if (Number.isNaN(parsedMs)) continue;

      const dayMs = utcDayStart(parsedMs);
      if (dayMs < yearStartMs || dayMs > yearEndMs) continue;
      if (createdAtMs !== null && dayMs < utcDayStart(createdAtMs)) continue;
      if (dayMs > todayMs) continue;

      const key = toDateKey(dayMs);
      byDay[key] = (byDay[key] ?? 0) + 1;
    }

    let maxCount = 0;
    for (const count of Object.values(byDay)) {
      if (count > maxCount) maxCount = count;
    }

    const weeks: HeatmapWeek[] = [];

    for (let weekIndex = 0; weekIndex < weekCount; weekIndex += 1) {
      const weekStartMs = addDaysMs(startWeekStartMs, weekIndex * 7);
      const cells: HeatmapCell[] = [];

      for (let dayIndex = 0; dayIndex < 7; dayIndex += 1) {
        const dateMs = addDaysMs(weekStartMs, dayIndex);
        const dateKey = toDateKey(dateMs);
        const outsideYear = dateMs < yearStartMs || dateMs > yearEndMs;
        const beforeAccount = createdAtMs !== null && dateMs < utcDayStart(createdAtMs);
        const inFuture = dateMs > todayMs;
        const count = outsideYear || inFuture || beforeAccount ? 0 : (byDay[dateKey] ?? 0);

        cells.push({
          dateMs,
          dateKey,
          count,
          level: outsideYear || inFuture || beforeAccount ? 0 : contributionLevel(count, maxCount),
          outsideYear,
          beforeAccount,
          inFuture
        });
      }

      weeks.push({
        startMs: weekStartMs,
        monthLabel: monthLabelForWeek(weekStartMs, yearEndMs),
        cells
      });
    }

    if (weeks[0] && weeks[0].monthLabel === '') {
      weeks[0].monthLabel = MONTH_FORMATTER.format(weeks[0].startMs);
    }

    return weeks;
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
        {m.profile_not_logged_in()} <a href={localizeHref('/login')} class="link">{m.nav_login()}</a>
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
        <div class="profile-meta-row">
          <p class="profile-meta-item">
            <span class="profile-meta-label">{m.profile_call_sign_label()}:</span>
            <span>{callSign ?? m.profile_call_sign_none()}</span>
          </p>
          <p class="profile-meta-item">
            {#if emailVerified}
              <Check size={14} class="profile-status-icon profile-status-icon-ok" aria-hidden="true" />
            {:else}
              <X size={14} class="profile-status-icon profile-status-icon-bad" aria-hidden="true" />
            {/if}
            <span
              >{emailVerified
                ? m.profile_email_status_verified()
                : m.profile_email_status_not_verified()}</span
            >
          </p>
        </div>
        <p class="profile-joined">
          <Calendar size={13} />
          {m.profile_member_since()}: {memberSince}
        </p>
      </div>
    </header>

    <section class="card profile-heatmap-card">
      <div class="profile-section-header profile-section-header--split">
        <div class="profile-section-title-wrap">
          <Activity size={16} />
          <h2 class="card-label">{yearTotalSessions} sessions in {selectedYear}</h2>
        </div>
        <div class="profile-heatmap-year-picker" role="tablist" aria-label="Heatmap year selector">
          {#each availableYears as year (year)}
            <button
              type="button"
              class={`profile-heatmap-year-btn ${year === selectedYear ? 'is-active' : ''}`}
              aria-pressed={year === selectedYear}
              onclick={() => {
                selectedYear = year;
              }}
            >
              {year}
            </button>
          {/each}
        </div>
      </div>

      <div class="profile-heatmap-scroll" bind:this={heatmapScrollEl}>
        <div
          class="profile-heatmap-shell"
          style={`--week-count:${heatmapWeekCount}`}
          role="img"
          aria-label={`Training activity heatmap for ${selectedYear}`}
        >
          <div class="profile-heatmap-months" aria-hidden="true">
            <div class="profile-heatmap-months-spacer"></div>
            <div class="profile-heatmap-months-grid">
              {#each heatmapWeeks as week, i (week.startMs)}
                <span class="profile-heatmap-month" style={`grid-column:${i + 1}`}>{week.monthLabel}</span>
              {/each}
            </div>
          </div>

          <div class="profile-heatmap-main">
            <div class="profile-heatmap-weekdays" aria-hidden="true">
              {#each WEEKDAY_LABELS as label, idx (`${idx}-${label}`)}
                <span>{label}</span>
              {/each}
            </div>

            <div class="profile-heatmap-grid">
              {#each heatmapWeeks as week (week.startMs)}
                <div
                  class={`profile-heatmap-week ${
                    selectedYear === CURRENT_YEAR &&
                    week.startMs <= todayDayMs &&
                    todayDayMs < week.startMs + 7 * DAY_MS
                      ? 'is-current-week'
                      : ''
                  }`}
                >
                  {#each week.cells as cell (cell.dateKey)}
                    <span
                      class={`profile-heatmap-cell level-${cell.level} ${cell.inFuture ? 'is-future' : ''} ${cell.beforeAccount ? 'is-before-account' : ''} ${cell.outsideYear ? 'is-outside-year' : ''} ${
                        selectedYear === CURRENT_YEAR && cell.dateMs === todayDayMs ? 'is-today' : ''
                      }`}
                      title={cell.outsideYear
                        ? `Outside ${selectedYear}`
                        : cell.beforeAccount
                        ? `Before account creation (${DAY_FORMATTER.format(cell.dateMs)})`
                        : `${cell.count} session${cell.count === 1 ? '' : 's'} on ${DAY_FORMATTER.format(cell.dateMs)}`}
                    ></span>
                  {/each}
                </div>
              {/each}
            </div>
          </div>
        </div>
      </div>

      <div class="profile-heatmap-legend" aria-hidden="true">
        <span>Less</span>
        <span class="profile-heatmap-cell level-0"></span>
        <span class="profile-heatmap-cell level-1"></span>
        <span class="profile-heatmap-cell level-2"></span>
        <span class="profile-heatmap-cell level-3"></span>
        <span class="profile-heatmap-cell level-4"></span>
        <span>More</span>
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
          {m.profile_history_empty()}
          <a href={localizeHref('/morse/learn')} class="link">{m.nav_learn()}</a>
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
  .profile-meta-row {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem 1rem;
    margin: 0;
  }
  .profile-meta-item {
    margin: 0;
    font-size: 0.8rem;
    color: var(--text-muted);
    display: inline-flex;
    gap: 0.25rem;
    align-items: center;
  }
  .profile-meta-label {
    color: var(--text-secondary);
    font-weight: 600;
  }
  :global(.profile-status-icon) {
    flex-shrink: 0;
  }
  :global(.profile-status-icon-ok) {
    color: #22c55e;
  }
  :global(.profile-status-icon-bad) {
    color: #ef4444;
  }

  /* Activity heatmap */
  .profile-heatmap-card {
    overflow: visible;
  }
  .profile-section-header--split {
    justify-content: space-between;
    gap: 0.75rem;
    flex-wrap: wrap;
  }
  .profile-section-title-wrap {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .profile-heatmap-year-picker {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
  }
  .profile-heatmap-year-btn {
    border: 1px solid var(--border);
    background: var(--bg-inset);
    color: var(--text-muted);
    border-radius: 0.5rem;
    padding: 0.2rem 0.5rem;
    font-size: 0.72rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.15s ease, color 0.15s ease, border-color 0.15s ease;
  }
  .profile-heatmap-year-btn:hover {
    color: var(--text-primary);
    border-color: color-mix(in srgb, var(--accent), var(--border) 45%);
  }
  .profile-heatmap-year-btn.is-active {
    background: color-mix(in srgb, var(--accent), transparent 88%);
    color: var(--accent);
    border-color: color-mix(in srgb, var(--accent), transparent 55%);
  }
  .profile-heatmap-scroll {
    overflow-x: auto;
    padding-right: 0.15rem;
    scrollbar-gutter: stable both-edges;
  }
  .profile-heatmap-shell {
    width: max-content;
  }
  .profile-heatmap-months {
    display: grid;
    grid-template-columns: 2.2rem max-content;
    align-items: center;
    margin-bottom: 0.4rem;
  }
  .profile-heatmap-months-spacer {
    width: 2.2rem;
  }
  .profile-heatmap-months-grid {
    display: grid;
    grid-template-columns: repeat(var(--week-count), 0.72rem);
    column-gap: 0.22rem;
    font-size: 0.74rem;
    color: var(--text-muted);
    line-height: 1;
  }
  .profile-heatmap-month {
    white-space: nowrap;
  }
  .profile-heatmap-main {
    display: grid;
    grid-template-columns: 2.2rem max-content;
    gap: 0.4rem;
  }
  .profile-heatmap-weekdays {
    display: grid;
    grid-template-rows: repeat(7, 0.72rem);
    row-gap: 0.22rem;
    font-size: 0.68rem;
    color: var(--text-muted);
    line-height: 1;
    align-items: center;
  }
  .profile-heatmap-grid {
    display: grid;
    grid-template-columns: repeat(var(--week-count), 0.72rem);
    column-gap: 0.22rem;
  }
  .profile-heatmap-week {
    display: grid;
    grid-template-rows: repeat(7, 0.72rem);
    row-gap: 0.22rem;
  }
  .profile-heatmap-cell {
    width: 0.72rem;
    height: 0.72rem;
    border-radius: 0.16rem;
    border: 1px solid color-mix(in srgb, var(--border), transparent 45%);
    background: color-mix(in srgb, var(--bg-inset), black 12%);
    display: inline-block;
  }
  .profile-heatmap-cell.level-1 {
    background: #0e4429;
    border-color: #0e4429;
  }
  .profile-heatmap-cell.level-2 {
    background: #006d32;
    border-color: #006d32;
  }
  .profile-heatmap-cell.level-3 {
    background: #26a641;
    border-color: #26a641;
  }
  .profile-heatmap-cell.level-4 {
    background: #39d353;
    border-color: #39d353;
  }
  .profile-heatmap-cell.is-future {
    opacity: 0.45;
  }
  .profile-heatmap-cell.is-outside-year {
    opacity: 0.2;
    border-style: dashed;
    border-color: color-mix(in srgb, var(--border), transparent 65%);
  }
  .profile-heatmap-cell.is-before-account {
    opacity: 0.32;
    border-color: color-mix(in srgb, var(--border), transparent 55%);
    background: color-mix(in srgb, var(--bg-inset), black 8%);
  }
  .profile-heatmap-legend {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 0.35rem;
    margin-top: 0.65rem;
    font-size: 0.76rem;
    color: var(--text-muted);
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
    .profile-cw-grid {
      flex-direction: column;
      gap: 0.75rem;
    }
    .profile-heatmap-year-picker {
      width: 100%;
      overflow-x: auto;
      padding-bottom: 0.2rem;
    }
  }
</style>

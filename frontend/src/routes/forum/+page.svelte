<script lang="ts">
  import { onMount } from 'svelte';
  import { afterNavigate, goto } from '$app/navigation';
  import { user } from '$lib/auth';
  import {
    FORUM_CATEGORIES,
    createForumThread,
    getForumThreads,
    isForumCategory,
    type ForumCategoryValue,
    type ForumThreadSummary
  } from '$lib/api';
  import { forumCategoryLabel, resolvePostingStatus, type PostingStatus } from '$lib/forum';
  import { localizeApiError } from '$lib/errorLocalization';
  import { authorLabel, formatDate } from '$lib/format';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import GuestNotice from '$lib/components/GuestNotice.svelte';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import { MessageSquare, MessagesSquare, Plus } from '@lucide/svelte';
  import * as m from '$lib/paraglide/messages';

  let threads = $state<ForumThreadSummary[]>([]);
  let nextCursor = $state<string | null>(null);
  let loading = $state(true);
  let loadingMore = $state(false);
  let loadFailed = $state('');
  let loadMoreError = $state('');
  let loadToken = 0;

  let composerOpen = $state(false);
  let postingStatus = $state<PostingStatus | null>(null);
  let composerCategory = $state<ForumCategoryValue>('general');
  let composerTitle = $state('');
  let composerBody = $state('');
  let composerError = $state('');
  let composerSubmitting = $state(false);

  // The active filter lives in the URL so the selection survives back/forward
  // navigation from a thread detail page and filtered lists can be shared. The
  // query string is read client-side only: `url.searchParams` is inaccessible
  // while this page is being prerendered.
  let filter = $state<ForumCategoryValue | null>(null);
  let filterInitialized = $state(false);

  function categoryFromUrl(value: string | null): ForumCategoryValue | null {
    return isForumCategory(value) ? value : null;
  }

  function syncFilterFromUrl(): void {
    if (typeof window === 'undefined') return;
    const next = categoryFromUrl(new URLSearchParams(window.location.search).get('category'));
    filterInitialized = true;
    if (next !== filter) filter = next;
  }

  onMount(syncFilterFromUrl);
  afterNavigate(syncFilterFromUrl);

  $effect(() => {
    if (!filterInitialized) return;
    const category = filter;
    void loadFirstPage(category);
  });

  $effect(() => {
    const username = $user?.username ?? null;
    let cancelled = false;
    postingStatus = null;
    void resolvePostingStatus(username).then((status) => {
      if (!cancelled) postingStatus = status;
    });
    return () => {
      cancelled = true;
    };
  });

  async function loadFirstPage(category: ForumCategoryValue | null): Promise<void> {
    const token = ++loadToken;
    loading = true;
    loadFailed = '';
    loadMoreError = '';
    nextCursor = null;

    try {
      const result = await getForumThreads({ category });
      if (token !== loadToken) return;
      threads = result.data;
      nextCursor = result.next_cursor;
    } catch (error) {
      if (token !== loadToken) return;
      threads = [];
      loadFailed = localizeApiError(error, () => m.api_error_forum_query_failed());
    } finally {
      if (token === loadToken) loading = false;
    }
  }

  async function loadMore(): Promise<void> {
    const cursor = nextCursor;
    if (!cursor || loadingMore) return;

    const token = loadToken;
    loadingMore = true;
    loadMoreError = '';

    try {
      const result = await getForumThreads({ category: filter, cursor });
      if (token !== loadToken) return;
      threads = [...threads, ...result.data];
      nextCursor = result.next_cursor;
    } catch (error) {
      if (token === loadToken) {
        loadMoreError = localizeApiError(error, () => m.api_error_forum_query_failed());
      }
    } finally {
      loadingMore = false;
    }
  }

  function selectCategory(category: ForumCategoryValue | null): void {
    if (category === filter) return;
    const path = category ? `${href('/forum')}?category=${category}` : href('/forum');
    void goto(path, { keepFocus: true, noScroll: true });
  }

  function closeComposer(): void {
    composerOpen = false;
    composerError = '';
  }

  function resetComposer(): void {
    composerCategory = 'general';
    composerTitle = '';
    composerBody = '';
    composerError = '';
  }

  async function submitThread(event: SubmitEvent): Promise<void> {
    event.preventDefault();
    if (composerSubmitting) return;

    const title = composerTitle.trim();
    const body = composerBody.trim();

    if (title.length < 3 || title.length > 200) {
      composerError = m.forum_thread_title_error();
      return;
    }

    if (!body) {
      composerError = m.forum_thread_body_error();
      return;
    }

    composerSubmitting = true;
    composerError = '';

    try {
      const createdCategory = composerCategory;
      await createForumThread(createdCategory, title, body);
      composerOpen = false;
      resetComposer();
      // Make the new thread visible: jump to its category when the active
      // filter would hide it, otherwise refresh the list the visitor sees.
      if (filter !== null && filter !== createdCategory) {
        selectCategory(createdCategory);
      } else {
        void loadFirstPage(filter);
      }
    } catch (error) {
      composerError = localizeApiError(error, () => m.api_error_forum_create_failed());
    } finally {
      composerSubmitting = false;
    }
  }
</script>

<svelte:head>
  <title>{m.forum_title()} | OpenCW</title>
</svelte:head>

<div class="forum-page page-wide">
  <header class="forum-hero">
    <div class="eyebrow"><MessagesSquare size={16} /> {m.forum_eyebrow()}</div>
    <h1 class="page-title">{m.forum_title()}</h1>
    <p class="body-text">{m.forum_intro()}</p>
  </header>

  <div class="forum-toolbar">
    <nav class="category-tabs">
      <button
        type="button"
        class="tab"
        class:is-active={filter === null}
        aria-pressed={filter === null}
        onclick={() => selectCategory(null)}
      >
        {m.forum_category_all()}
      </button>
      {#each FORUM_CATEGORIES as category (category)}
        <button
          type="button"
          class="tab"
          class:is-active={filter === category}
          aria-pressed={filter === category}
          onclick={() => selectCategory(category)}
        >
          {forumCategoryLabel(category)}
        </button>
      {/each}
    </nav>
    <button
      type="button"
      class="btn-primary new-thread-btn"
      aria-expanded={composerOpen}
      onclick={() => (composerOpen ? closeComposer() : (composerOpen = true))}
    >
      <Plus size={16} aria-hidden="true" />
      {m.forum_new_thread()}
    </button>
  </div>

  {#if composerOpen}
    <section class="panel composer">
      <h2 class="composer-title">{m.forum_new_thread_heading()}</h2>

      {#if postingStatus === 'guest'}
        <GuestNotice class="body-text" message={m.forum_guest_notice()} />
        <div class="composer-actions">
          <button type="button" class="btn-ghost" onclick={closeComposer}>
            {m.forum_cancel()}
          </button>
        </div>
      {:else if postingStatus === 'unverified'}
        <div class="notice gate-notice">
          <p class="gate-note">{m.forum_email_unverified_notice()}</p>
          <a class="link" href={href('/settings')}>{m.nav_settings()}</a>
        </div>
        <div class="composer-actions">
          <button type="button" class="btn-ghost" onclick={closeComposer}>
            {m.forum_cancel()}
          </button>
        </div>
      {:else if postingStatus === 'ready'}
        <form class="composer-form" onsubmit={submitThread}>
          <label class="composer-field">
            <span class="label-text">{m.forum_thread_category_label()}</span>
            <select class="select" bind:value={composerCategory}>
              {#each FORUM_CATEGORIES as category (category)}
                <option value={category}>{forumCategoryLabel(category)}</option>
              {/each}
            </select>
          </label>

          <label class="composer-field">
            <span class="label-text">{m.forum_thread_title_label()}</span>
            <input
              class="input"
              type="text"
              maxlength="200"
              placeholder={m.forum_thread_title_placeholder()}
              bind:value={composerTitle}
            />
          </label>

          <label class="composer-field">
            <span class="label-text">{m.forum_thread_body_label()}</span>
            <textarea
              class="input composer-textarea"
              maxlength="10000"
              placeholder={m.forum_thread_body_placeholder()}
              bind:value={composerBody}></textarea>
          </label>

          {#if composerError}
            <ErrorAlert message={composerError} />
          {/if}

          <div class="composer-actions">
            <button type="submit" class="btn-primary" disabled={composerSubmitting}>
              {composerSubmitting ? m.settings_saving() : m.forum_thread_submit()}
            </button>
            <button type="button" class="btn-ghost" onclick={closeComposer}>
              {m.forum_cancel()}
            </button>
          </div>
        </form>
      {:else}
        <LoadingSpinner />
      {/if}
    </section>
  {/if}

  {#if loading}
    <LoadingSpinner />
  {:else if loadFailed}
    <div class="forum-error">
      <ErrorAlert message={loadFailed} />
      <button type="button" class="btn-ghost retry-btn" onclick={() => loadFirstPage(filter)}>
        {m.forum_retry()}
      </button>
    </div>
  {:else if threads.length === 0}
    <section class="panel forum-empty">
      <MessagesSquare size={34} aria-hidden="true" />
      <h2 class="empty-title">{m.forum_threads_empty_title()}</h2>
      <p class="body-text empty-note">{m.forum_threads_empty_body()}</p>
    </section>
  {:else}
    <div class="thread-list">
      {#each threads as thread (thread.id)}
        <a class="card card--interactive thread-card" href={href(`/forum/${thread.id}`)}>
          <div class="thread-card-top">
            <span class="category-badge">{forumCategoryLabel(thread.category)}</span>
            <h3 class="thread-title">{thread.title}</h3>
          </div>
          <div class="thread-meta">
            <span>{authorLabel(thread.author)}</span>
            {#if thread.author.call_sign}
              <span class="callsign">{thread.author.call_sign}</span>
            {/if}
            <span>{formatDate(thread.created_at)}</span>
            <span class="thread-replies">
              <MessageSquare size={14} aria-hidden="true" />
              {thread.reply_count}
              <span class="sr-only">{m.forum_replies_label()}</span>
            </span>
          </div>
        </a>
      {/each}
    </div>

    {#if loadMoreError}
      <ErrorAlert message={loadMoreError} />
    {/if}

    {#if nextCursor}
      <div class="load-more-row">
        <button
          type="button"
          class="btn-ghost load-more-btn"
          disabled={loadingMore}
          onclick={loadMore}
        >
          {m.forum_threads_load_more()}
        </button>
      </div>
    {/if}
  {/if}
</div>

<style>
  .forum-hero {
    margin-bottom: var(--section-gap);
  }

  .eyebrow {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-bottom: 0.6rem;
    color: var(--accent);
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  .forum-hero .page-title {
    margin: 0;
  }

  .forum-hero .body-text {
    margin: 0.45rem 0 0;
    max-width: var(--max-width-narrow);
  }

  .forum-toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    flex-wrap: wrap;
    margin-bottom: var(--block-gap);
  }

  .category-tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    min-width: 0;
  }

  .tab {
    padding: 0.4rem 0.85rem;
    font-size: 0.875rem;
    line-height: 1.25rem;
    border-radius: 999px;
    border: 1px solid var(--border-subtle);
    background: transparent;
    color: var(--text-secondary);
    cursor: pointer;
    transition:
      border-color var(--transition-fast),
      color var(--transition-fast),
      background-color var(--transition-fast);
  }

  .tab:hover {
    border-color: var(--accent);
    color: var(--accent);
  }

  .tab.is-active {
    border-color: var(--accent);
    background-color: color-mix(in srgb, var(--accent) 12%, transparent);
    color: var(--accent);
    font-weight: 600;
  }

  .new-thread-btn {
    width: auto;
    padding: 0.6rem 1rem;
    font-size: 0.9rem;
  }

  .composer {
    display: flex;
    flex-direction: column;
    gap: var(--block-gap);
    margin-bottom: var(--block-gap);
  }

  .composer-title {
    margin: 0;
    font-size: 1.15rem;
    font-weight: 700;
  }

  .composer-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .composer-field {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .composer-textarea {
    min-height: 8rem;
    font-family: inherit;
    resize: vertical;
  }

  .composer-actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .composer-actions .btn-primary,
  .composer-actions .btn-ghost {
    flex: 1 1 0;
  }

  .gate-notice {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .gate-note {
    margin: 0;
  }

  .forum-error {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    align-items: flex-start;
  }

  .retry-btn {
    width: auto;
  }

  .forum-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    padding: 2.5rem 1.5rem;
    text-align: center;
  }

  .forum-empty > :global(svg) {
    color: var(--accent);
  }

  .empty-title {
    margin: 0;
    font-size: 1.35rem;
    font-weight: 700;
  }

  .empty-note {
    max-width: 34rem;
    margin: 0;
  }

  .thread-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .thread-card {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    padding: 1rem 1.25rem;
  }

  .thread-card-top {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-wrap: wrap;
  }

  .category-badge {
    padding: 0.15rem 0.6rem;
    font-size: 0.72rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    border-radius: var(--radius-sm);
    border: 1px solid color-mix(in srgb, var(--accent) 35%, transparent);
    background-color: color-mix(in srgb, var(--accent) 10%, transparent);
    color: var(--accent);
  }

  .thread-title {
    margin: 0;
    font-size: 1.05rem;
    font-weight: 600;
    line-height: 1.35;
    color: var(--text-primary);
  }

  .thread-meta {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
    color: var(--text-muted);
    font-size: 0.82rem;
  }

  .callsign {
    padding: 0.1rem 0.4rem;
    font-family:
      ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
      monospace;
    font-size: 0.75rem;
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    color: var(--accent);
  }

  .thread-replies {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
  }

  .load-more-row {
    display: flex;
    justify-content: center;
    margin-top: var(--block-gap);
  }

  .load-more-btn {
    width: auto;
    padding: 0.6rem 1.25rem;
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
  }
</style>

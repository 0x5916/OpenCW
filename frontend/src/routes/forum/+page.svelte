<script lang="ts">
  import { MessageSquare, Plus, Pin, Lock, MessagesSquare, ListFilter } from '@lucide/svelte';
  import { user } from '$lib/auth';
  import { getForumCategories, getForumCategoryThreads, type ForumThread } from '$lib/api';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import { localizeApiError } from '$lib/errorLocalization';
  import { lang } from '$lib/i18n.svelte';
  import { localizeHref } from '$lib/paraglide/runtime';
  import * as m from '$lib/paraglide/messages';
  type ThreadPreview = ForumThread & {
    categoryName: string;
  };

  type SortMode = 'activity' | 'created' | 'category';

  let loading = $state(true);
  let loadError = $state('');
  let threads = $state<ThreadPreview[]>([]);
  let sortMode = $state<SortMode>('activity');

  let sortedThreads = $derived(sortThreads(threads, sortMode));

  $effect(() => {
    void loadForum();
  });
  function href(path: string): string {
    return localizeHref(path, { locale: lang.value });
  }
  async function loadForum() {
    loading = true;
    loadError = '';
    try {
      const categories = await getForumCategories();
      const categoryThreads = await Promise.all(
        categories.map(async (category) => {
          const result = await getForumCategoryThreads(category.id, 1, 100);
          return result.data.map((thread) => ({ ...thread, categoryName: category.name }));
        })
      );
      threads = categoryThreads.flat();
    } catch (error) {
      loadError = localizeApiError(error, () => m.forum_load_error());
    } finally {
      loading = false;
    }
  }

  function sortThreads(items: ThreadPreview[], mode: SortMode): ThreadPreview[] {
    return [...items].sort((left, right) => {
      if (mode === 'category') {
        const categoryOrder = left.categoryName.localeCompare(right.categoryName);
        if (categoryOrder !== 0) return categoryOrder;
      }

      const leftTime = Date.parse(
        mode === 'created'
          ? left.created_at
          : (left.latest_activity ?? left.updated_at ?? left.created_at)
      );
      const rightTime = Date.parse(
        mode === 'created'
          ? right.created_at
          : (right.latest_activity ?? right.updated_at ?? right.created_at)
      );
      return (Number.isNaN(rightTime) ? 0 : rightTime) - (Number.isNaN(leftTime) ? 0 : leftTime);
    });
  }
  function authorLabel(thread: ForumThread): string {
    return thread.username ?? thread.author ?? m.forum_community_member();
  }
  function formatDate(value?: string): string {
    if (!value) return '';
    const date = new Date(value);
    return Number.isNaN(date.getTime())
      ? ''
      : date.toLocaleDateString(undefined, {
          year: 'numeric',
          month: 'short',
          day: 'numeric'
        });
  }
</script>

<svelte:head>
  <title>{m.forum_title()} | OpenCW</title>
</svelte:head>
<main class="forum-page page-wide">
  <header class="forum-hero">
    <div>
      <div class="eyebrow"><MessageSquare size={16} /> {m.forum_eyebrow()}</div>
      <h1 class="page-title">{m.forum_title()}</h1>
      <p class="body-text">{m.forum_intro()}</p>
    </div>
    {#if $user}
      <a class="btn-primary forum-action" href={href('/forum/new')}>
        <Plus size={17} />
        {m.forum_new_thread()}
      </a>
    {:else}
      <a class="btn-ghost forum-action" href={href('/login')}>
        <Plus size={17} />
        {m.forum_login_to_post()}
      </a>
    {/if}
  </header>

  {#if loading}
    <div class="forum-state card" aria-live="polite"><LoadingSpinner /></div>
  {:else if loadError}
    <ErrorAlert message={loadError} />
  {:else if threads.length === 0}
    <div class="forum-state card">
      <MessagesSquare size={30} />
      <p>{m.forum_no_categories()}</p>
    </div>
  {:else}
    <div class="feed-toolbar">
      <div class="feed-label"><ListFilter size={16} /> {m.forum_sort_label()}</div>
      <label class="sort-control" for="forum-sort">
        <span class="sr-only">{m.forum_sort_label()}</span>
        <select id="forum-sort" bind:value={sortMode}>
          <option value="activity">{m.forum_sort_activity()}</option>
          <option value="created">{m.forum_sort_created()}</option>
          <option value="category">{m.forum_sort_category()}</option>
        </select>
      </label>
    </div>

    <div class="thread-feed">
      {#each sortedThreads as thread (thread.id)}
        <article class="thread-card card">
          <div class="thread-card-main">
            <a class="category-label" href={href(`/forum/category/${thread.category_id}`)}>
              {thread.categoryName}
            </a>
            <a class="thread-link" href={href(`/forum/thread/${thread.id}`)}>
              <span class="thread-title">
                {#if thread.is_pinned}<Pin size={14} aria-label={m.forum_pinned()} />{/if}
                {#if thread.is_locked}<Lock size={14} aria-label={m.forum_locked()} />{/if}
                <span>{thread.title}</span>
              </span>
              <span class="thread-meta">
                {authorLabel(thread)} · {formatDate(
                  thread.latest_activity ?? thread.updated_at ?? thread.created_at
                )}
              </span>
            </a>
          </div>
          {#if thread.post_count !== undefined}
            <span class="post-count">{m.forum_post_count({ count: thread.post_count })}</span>
          {/if}
        </article>
      {/each}
    </div>
  {/if}
</main>

<style>
  .forum-page {
    padding-top: 2.25rem;
    padding-bottom: 3rem;
  }

  .forum-hero {
    display: flex;
    align-items: end;
    justify-content: space-between;
    gap: 1.5rem;
    margin-bottom: 2rem;
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
    max-width: 38rem;
    margin: 0.45rem 0 0;
  }
  .forum-action {
    display: inline-flex;
    width: auto;
    flex: 0 0 auto;
    flex-shrink: 0;
    text-decoration: none;
  }

  .feed-toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 0.75rem;
  }

  .feed-label {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--text-muted);
    font-size: 0.8rem;
    font-weight: 650;
  }

  .sort-control select {
    min-width: 11rem;
    padding: 0.5rem 0.7rem;
    color: var(--text-secondary);
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    font: inherit;
    font-size: 0.8rem;
  }

  .sort-control select:focus {
    outline: none;
    border-color: var(--accent);
    box-shadow: var(--focus-ring);
  }

  .thread-feed {
    display: grid;
    gap: 0.75rem;
  }

  .thread-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 1.1rem 1.25rem;
  }

  .thread-card-main {
    min-width: 0;
  }

  .category-label {
    display: inline-block;
    margin-bottom: 0.35rem;
    color: var(--accent);
    font-size: 0.73rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-decoration: none;
    text-transform: uppercase;
  }

  .category-label:hover {
    color: var(--accent-hover);
    text-decoration: underline;
  }

  .thread-link {
    display: grid;
    gap: 0.15rem;
    color: var(--text-primary);
    text-decoration: none;
  }

  .thread-link:hover .thread-title,
  .thread-link:focus-visible .thread-title {
    color: var(--accent);
  }
  .thread-title {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    min-width: 0;
    font-size: 0.92rem;
    font-weight: 650;
  }

  .thread-title span:last-child {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .thread-title :global(svg) {
    flex-shrink: 0;
    color: var(--accent);
  }
  .thread-meta,
  .post-count {
    color: var(--text-muted);
    font-size: 0.78rem;
  }

  .post-count {
    flex-shrink: 0;
    min-width: 2rem;
    text-align: right;
  }

  .forum-state {
    display: grid;
    place-items: center;
    gap: 0.75rem;
    min-height: 12rem;
    color: var(--text-muted);
    text-align: center;
  }

  .forum-state p {
    margin: 0;
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

  @media (max-width: 700px) {
    .forum-hero {
      align-items: stretch;
      flex-direction: column;
    }

    .forum-action {
      align-self: start;
    }

    .feed-toolbar {
      align-items: stretch;
      flex-direction: column;
    }

    .sort-control select {
      width: 100%;
    }

    .thread-card {
      align-items: flex-start;
    }
  }
</style>

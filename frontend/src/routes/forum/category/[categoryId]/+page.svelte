<script lang="ts">
  import { ArrowLeft, Lock, MessageSquare, Pin } from '@lucide/svelte';
  import { page } from '$app/state';
  import {
    getForumCategories,
    getForumCategoryThreads,
    type ForumCategory,
    type ForumThread
  } from '$lib/api';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import Pagination from '$lib/components/Pagination.svelte';
  import { localizeApiError } from '$lib/errorLocalization';
  import { authorLabel, formatDate } from '$lib/format';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';

  const DEFAULT_LIMIT = 20;
  const MAX_LIMIT = 100;
  let category = $state<ForumCategory | null>(null);
  let threads = $state<ForumThread[]>([]);
  let currentPage = $state(1);
  let totalPages = $state(1);
  let loading = $state(true);
  let error = $state('');

  $effect(() => {
    const categoryId = page.params.categoryId;
    if (!categoryId) {
      error = m.api_error_forum_categories_fetch_failed();
      loading = false;
      return;
    }
    const requestedPage = Number.parseInt(page.url.searchParams.get('page') ?? '1', 10);
    currentPage = Number.isInteger(requestedPage) && requestedPage > 0 ? requestedPage : 1;
    void loadCategory(categoryId, currentPage);
  });

  async function loadCategory(categoryId: string, requestedPage: number) {
    loading = true;
    error = '';
    try {
      const categories = await getForumCategories();
      category = categories.find((item) => item.id === categoryId) ?? null;
      const result = await getForumCategoryThreads(categoryId, requestedPage, DEFAULT_LIMIT);
      threads = result.data;
      totalPages =
        result.total_pages ??
        (result.total ? Math.max(1, Math.ceil(result.total / result.limit)) : 1);
    } catch (cause) {
      error = localizeApiError(cause, () => m.forum_threads_load_error());
    } finally {
      loading = false;
    }
  }

  function pageHref(nextPage: number): string {
    return href(
      `/forum/category/${page.params.categoryId}?page=${Math.min(MAX_LIMIT, Math.max(1, nextPage))}`
    );
  }
</script>

<svelte:head><title>{category?.name ?? m.forum_title()} | OpenCW</title></svelte:head>

<main class="forum-page page-wide">
  <a class="back-link" href={href('/forum')}><ArrowLeft size={16} /> {m.forum_title()}</a>
  {#if category}
    <header class="route-header">
      <div>
        <div class="eyebrow"><MessageSquare size={16} /> {m.forum_eyebrow()}</div>
        <h1 class="page-title">{category.name}</h1>
        {#if category.description}<p class="body-text">{category.description}</p>{/if}
      </div>
      <a class="btn-primary" href={href(`/forum/new?category=${category.id}`)}
        >+ {m.forum_new_thread()}</a
      >
    </header>
  {/if}

  {#if loading}
    <div class="forum-state card"><LoadingSpinner /></div>
  {:else if error}
    <ErrorAlert message={error} />
  {:else if !category}
    <div class="forum-state card"><p>{m.api_error_forum_categories_fetch_failed()}</p></div>
  {:else if threads.length === 0}
    <div class="forum-state card">
      <MessageSquare size={30} />
      <p>{m.forum_no_threads()}</p>
    </div>
  {:else}
    <section class="thread-table card">
      {#each threads as thread (thread.id)}
        <a class="thread-row" href={href(`/forum/thread/${thread.id}`)}>
          <span class="thread-main"
            ><span class="thread-title">
              {#if thread.is_pinned}<Pin size={15} />{/if}{#if thread.is_locked}<Lock
                  size={15}
                />{/if}<span>{thread.title}</span>
            </span><span class="thread-meta"
              >{authorLabel(thread)} · {formatDate(
                thread.latest_activity ?? thread.updated_at ?? thread.created_at
              )}</span
            ></span
          >
          <span class="thread-count">{m.forum_post_count({ count: thread.post_count ?? 0 })}</span>
        </a>
      {/each}
    </section>
    <Pagination
      {currentPage}
      {totalPages}
      prevHref={pageHref(currentPage - 1)}
      nextHref={pageHref(currentPage + 1)}
    />
  {/if}
</main>

<style>
  .forum-page {
    padding-top: 2.25rem;
    padding-bottom: 3rem;
  }
  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--text-muted);
    text-decoration: none;
    font-size: 0.85rem;
  }
  .back-link:hover {
    color: var(--accent);
  }
  .route-header {
    display: flex;
    align-items: end;
    justify-content: space-between;
    gap: 1rem;
    margin: 1.5rem 0 2rem;
  }
  .route-header .page-title {
    margin: 0;
  }
  .route-header .body-text {
    margin: 0.4rem 0 0;
  }
  .eyebrow {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--accent);
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    margin-bottom: 0.6rem;
  }
  .thread-table {
    padding: 0 1.25rem;
  }
  .thread-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 1rem 0;
    color: var(--text-primary);
    text-decoration: none;
    border-bottom: 1px solid var(--border);
  }
  .thread-row:last-child {
    border-bottom: 0;
  }
  .thread-row:hover .thread-title {
    color: var(--accent);
  }
  .thread-main {
    min-width: 0;
  }
  .thread-title {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-weight: 650;
  }
  .thread-title span:last-child {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .thread-title :global(svg) {
    flex-shrink: 0;
    color: var(--accent);
  }
  .thread-meta {
    display: block;
    color: var(--text-muted);
    font-size: 0.8rem;
    margin-top: 0.15rem;
  }
  .thread-count {
    color: var(--text-muted);
    font-size: 0.8rem;
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
  @media (max-width: 600px) {
    .route-header {
      align-items: stretch;
      flex-direction: column;
    }
    .route-header .btn-primary {
      align-self: start;
    }
    .thread-table {
      padding: 0 1rem;
    }
  }
</style>

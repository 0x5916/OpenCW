<script lang="ts">
  import { ArrowLeft, Lock, MessageSquare, Send } from '@lucide/svelte';
  import { page } from '$app/state';
  import { user } from '$lib/auth';
  import {
    createForumPost,
    getForumThread,
    getForumThreadPosts,
    type ForumPost,
    type ForumThread
  } from '$lib/api';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import Pagination from '$lib/components/Pagination.svelte';
  import { localizeApiError } from '$lib/errorLocalization';
  import { authorLabel, formatDateTime } from '$lib/format';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';

  let thread = $state<ForumThread | null>(null);
  let posts = $state<ForumPost[]>([]);
  let currentPage = $state(1);
  let totalPages = $state(1);
  let loading = $state(true);
  let error = $state('');
  let replyBody = $state('');
  let replyError = $state('');
  let submitting = $state(false);

  $effect(() => {
    const threadId = page.params.threadId;
    if (!threadId) {
      error = m.api_error_forum_thread_not_found();
      loading = false;
      return;
    }
    const requestedPage = Number.parseInt(page.url.searchParams.get('page') ?? '1', 10);
    currentPage = Number.isInteger(requestedPage) && requestedPage > 0 ? requestedPage : 1;
    void loadThread(threadId, currentPage);
  });

  async function loadThread(threadId: string, requestedPage: number) {
    loading = true;
    error = '';
    try {
      const [threadResult, postsResult] = await Promise.all([
        getForumThread(threadId),
        getForumThreadPosts(threadId, requestedPage, 20)
      ]);
      thread = threadResult;
      posts = postsResult.data;
      totalPages =
        postsResult.total_pages ??
        (postsResult.total ? Math.max(1, Math.ceil(postsResult.total / postsResult.limit)) : 1);
    } catch (cause) {
      error = localizeApiError(cause, () => m.forum_load_error());
    } finally {
      loading = false;
    }
  }

  async function submitReply(event: SubmitEvent) {
    event.preventDefault();
    if (!thread || !replyBody.trim()) {
      replyError = m.forum_reply_required();
      return;
    }
    submitting = true;
    replyError = '';
    try {
      await createForumPost(thread.id, replyBody.trim());
      replyBody = '';
      await loadThread(thread.id, 1);
    } catch (cause) {
      replyError = localizeApiError(cause, () => m.api_error_forum_post_create_failed());
    } finally {
      submitting = false;
    }
  }

  function pageHref(nextPage: number): string {
    return href(`/forum/thread/${page.params.threadId}?page=${Math.max(1, nextPage)}`);
  }
</script>

<svelte:head><title>{thread?.title ?? m.forum_title()} | OpenCW</title></svelte:head>

<main class="forum-page page-wide">
  <a class="back-link" href={href('/forum')}><ArrowLeft size={16} /> {m.forum_title()}</a>
  {#if loading}<div class="forum-state card"><LoadingSpinner /></div>
  {:else if error}<ErrorAlert message={error} />
  {:else if thread}
    <header class="thread-header">
      <div class="eyebrow"><MessageSquare size={16} /> {m.forum_thread()}</div>
      <h1 class="page-title">{thread.title}</h1>
      <p class="thread-meta">
        {authorLabel(thread)} · {formatDateTime(thread.created_at)}
      </p>
      {#if thread.is_locked}<div class="locked-note">
          <Lock size={16} />
          {m.forum_thread_locked()}
        </div>{/if}
    </header>

    <section class="posts" aria-label={m.forum_posts()}>
      {#each posts as post, index (post.id)}
        <article class="post card">
          <header class="post-header">
            <strong>{authorLabel(post)}</strong><time datetime={post.created_at}
              >{formatDateTime(post.created_at)}</time
            >
          </header>
          {#if post.parent_id}<div class="parent-note">{m.forum_reply_to()}</div>{/if}
          <p>{post.body}</p>
          <span class="post-number">#{(currentPage - 1) * 20 + index + 1}</span>
        </article>
      {/each}
    </section>

    <Pagination
      {currentPage}
      {totalPages}
      prevHref={pageHref(currentPage - 1)}
      nextHref={pageHref(currentPage + 1)}
    />

    {#if $user && !thread.is_locked}
      <form class="reply-form card" onsubmit={submitReply}>
        <label for="reply-body">{m.forum_reply_label()}</label>
        <textarea
          id="reply-body"
          bind:value={replyBody}
          rows="5"
          maxlength="10000"
          placeholder={m.forum_reply_placeholder()}></textarea>
        {#if replyError}<p class="form-error">{replyError}</p>{/if}
        <button class="btn-primary" type="submit" disabled={submitting}
          ><Send size={16} />{submitting ? m.forum_replying() : m.forum_reply_submit()}</button
        >
      </form>
    {:else if thread.is_locked}
      <div class="locked-note reply-disabled"><Lock size={16} /> {m.forum_thread_locked()}</div>
    {:else}
      <p class="login-prompt"><a href={href('/login')}>{m.forum_login_to_post()}</a></p>
    {/if}
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
  .back-link:hover,
  .login-prompt a:hover {
    color: var(--accent);
  }
  .thread-header {
    margin: 1.5rem 0 2rem;
  }
  .thread-header .page-title {
    margin: 0;
    max-width: 52rem;
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
  .thread-meta,
  .post-header time {
    color: var(--text-muted);
    font-size: 0.82rem;
  }
  .thread-meta {
    margin: 0.5rem 0 0;
  }
  .posts {
    display: grid;
    gap: 1rem;
  }
  .post {
    position: relative;
    padding: 1.25rem;
  }
  .post-header {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--border);
  }
  .post p {
    white-space: pre-wrap;
    overflow-wrap: anywhere;
    margin: 1rem 0 0.5rem;
    color: var(--text-secondary);
  }
  .post-number {
    color: var(--text-muted);
    font-size: 0.75rem;
  }
  .parent-note {
    color: var(--accent);
    font-size: 0.78rem;
    margin-top: 0.75rem;
  }
  .locked-note {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    margin-top: 0.8rem;
    color: var(--text-muted);
    font-size: 0.85rem;
  }
  .forum-state {
    display: grid;
    place-items: center;
    min-height: 12rem;
  }
  .reply-form {
    display: grid;
    gap: 0.75rem;
    margin-top: 2rem;
  }
  .reply-form label {
    color: var(--text-label);
    font-weight: 650;
  }
  textarea {
    width: 100%;
    resize: vertical;
    padding: 0.75rem;
    color: var(--text-primary);
    background: var(--bg-inset);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    font: inherit;
  }
  textarea:focus {
    outline: none;
    box-shadow: var(--focus-ring);
    border-color: var(--accent);
  }
  .form-error {
    margin: 0;
    color: var(--danger);
    font-size: 0.85rem;
  }
  .reply-disabled {
    padding: 1rem;
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
  }
  .login-prompt {
    margin-top: 2rem;
    text-align: center;
  }
  .login-prompt a {
    color: var(--accent);
  }
</style>

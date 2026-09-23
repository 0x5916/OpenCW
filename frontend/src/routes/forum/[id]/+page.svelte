<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { user } from '$lib/auth';
  import {
    ApiError,
    createForumReply,
    deleteForumReply,
    deleteForumThread,
    getForumThread,
    getForumThreadReplies,
    type ForumReply,
    type ForumThreadDetail
  } from '$lib/api';
  import { forumCategoryLabel, resolvePostingStatus, type PostingStatus } from '$lib/forum';
  import { localizeApiError } from '$lib/errorLocalization';
  import { extractErrorCode } from '$lib/errorCode';
  import { authorLabel, formatDate } from '$lib/format';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import GuestNotice from '$lib/components/GuestNotice.svelte';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ForumReplyItem from './ForumReplyItem.svelte';
  import { ArrowLeft } from '@lucide/svelte';
  import * as m from '$lib/paraglide/messages';

  const threadId = $derived(page.params.id);

  let thread = $state<ForumThreadDetail | null>(null);
  let replies = $state<ForumReply[]>([]);
  let loading = $state(true);
  let loadFailed = $state('');
  let notFound = $state(false);
  let loadToken = 0;

  let postingStatus = $state<PostingStatus | null>(null);

  let replyBody = $state('');
  let replyTarget = $state<{ id: string; username: string } | null>(null);
  let replyError = $state('');
  let replySubmitting = $state(false);
  let repliesError = $state('');

  let threadDeleteArmed = $state(false);
  let deletingThread = $state(false);
  let deleteError = $state('');

  const isThreadAuthor = $derived(
    thread !== null && $user !== null && thread.author.username === $user.username
  );

  // Replies arrive as a nested tree; the visible count is the number of nodes.
  function countReplies(list: ForumReply[]): number {
    return list.reduce((sum, reply) => sum + 1 + countReplies(reply.children ?? []), 0);
  }

  let replyCount = $derived(countReplies(replies));
  let replyCountLabel = $derived(
    replyCount === 1 ? m.forum_reply_one() : m.forum_reply_many({ count: String(replyCount) })
  );

  $effect(() => {
    const id = threadId;
    if (!id) return;
    void loadThread(id);
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

  function isThreadNotFound(error: unknown): boolean {
    if (error instanceof ApiError && error.status === 404) return true;
    return extractErrorCode(error) === 'THREAD_NOT_FOUND';
  }

  async function loadThread(id: string): Promise<void> {
    const token = ++loadToken;
    loading = true;
    loadFailed = '';
    notFound = false;
    repliesError = '';
    deleteError = '';
    thread = null;
    replies = [];
    replyTarget = null;
    threadDeleteArmed = false;

    try {
      const [threadData, replyData] = await Promise.all([
        getForumThread(id),
        getForumThreadReplies(id)
      ]);
      if (token !== loadToken) return;
      thread = threadData;
      replies = replyData.data;
    } catch (error) {
      if (token !== loadToken) return;
      if (isThreadNotFound(error)) {
        notFound = true;
      } else {
        loadFailed = localizeApiError(error, () => m.api_error_forum_query_failed());
      }
    } finally {
      if (token === loadToken) loading = false;
    }
  }

  async function refreshReplies(): Promise<void> {
    const id = threadId;
    if (!id) return;

    try {
      const data = await getForumThreadReplies(id);
      replies = data.data;
      repliesError = '';
    } catch (error) {
      repliesError = localizeApiError(error, () => m.api_error_forum_query_failed());
    }
  }

  function retryLoad(): void {
    if (threadId) void loadThread(threadId);
  }

  function setReplyTarget(reply: ForumReply): void {
    replyTarget = { id: reply.id, username: authorLabel(reply.author) };
    replyError = '';
  }

  async function submitReply(event: SubmitEvent): Promise<void> {
    event.preventDefault();
    const id = threadId;
    if (!id || replySubmitting) return;

    const body = replyBody.trim();
    if (!body) {
      replyError = m.forum_reply_body_error();
      return;
    }

    replySubmitting = true;
    replyError = '';

    try {
      await createForumReply(id, body, replyTarget?.id);
      replyBody = '';
      replyTarget = null;
      await refreshReplies();
    } catch (error) {
      replyError = localizeApiError(error, () => m.api_error_forum_create_failed());
    } finally {
      replySubmitting = false;
    }
  }

  async function onDeleteReply(replyId: string): Promise<void> {
    try {
      await deleteForumReply(replyId);
      await refreshReplies();
    } catch (error) {
      repliesError = localizeApiError(error, () => m.api_error_forum_delete_failed());
    }
  }

  async function deleteThread(): Promise<void> {
    const id = threadId;
    if (!id || deletingThread) return;

    if (!threadDeleteArmed) {
      threadDeleteArmed = true;
      return;
    }

    deletingThread = true;
    deleteError = '';

    try {
      await deleteForumThread(id);
      await goto(href('/forum'));
    } catch (error) {
      deleteError = localizeApiError(error, () => m.api_error_forum_delete_failed());
    } finally {
      deletingThread = false;
      threadDeleteArmed = false;
    }
  }
</script>

<svelte:head>
  <title>{thread ? `${thread.title} | OpenCW` : `${m.forum_title()} | OpenCW`}</title>
</svelte:head>

<div class="thread-page page-narrow">
  <nav class="thread-crumbs" aria-label={m.forum_title()}>
    <a class="crumb" href={href('/forum')}>
      <ArrowLeft size={14} aria-hidden="true" />
      {m.forum_back_to_threads()}
    </a>
  </nav>

  {#if loading}
    <div class="thread-skeleton">
      <span class="skeleton skeleton-title"></span>
      <span class="skeleton skeleton-line"></span>
      <span class="skeleton skeleton-line short"></span>
    </div>
    <p class="sr-only" role="status">{m.common_loading()}</p>
  {:else if notFound}
    <section class="state-block">
      <h1 class="state-title">{m.forum_thread_not_found_title()}</h1>
      <p class="body-text state-note">{m.forum_thread_not_found_body()}</p>
      <a class="btn-primary" href={href('/forum')}>{m.forum_back_to_threads()}</a>
    </section>
  {:else if loadFailed}
    <div class="thread-error">
      <ErrorAlert message={loadFailed} />
      <button type="button" class="btn-ghost" onclick={retryLoad}>{m.forum_retry()}</button>
    </div>
  {:else if thread}
    <article class="thread-main">
      <header class="thread-head">
        <p class="thread-cat">
          <span class="chip-dot"></span>{forumCategoryLabel(thread.category)}
        </p>
        <h1 class="thread-title">{thread.title}</h1>
        <p class="thread-meta">
          <span class="thread-author">{authorLabel(thread.author)}</span>
          {#if thread.author.call_sign}
            <span class="callsign">{thread.author.call_sign}</span>
          {/if}
          <span>{formatDate(thread.created_at)}</span>
        </p>
      </header>
      <p class="thread-body">{thread.body}</p>

      {#if isThreadAuthor}
        <div class="thread-actions">
          {#if deleteError}
            <ErrorAlert message={deleteError} />
          {/if}
          <div class="delete-row">
            <button
              type="button"
              class="btn-danger"
              disabled={deletingThread}
              onclick={deleteThread}
            >
              {threadDeleteArmed ? m.forum_delete_confirm() : m.forum_delete_thread()}
            </button>
            {#if threadDeleteArmed}
              <button type="button" class="btn-ghost" onclick={() => (threadDeleteArmed = false)}>
                {m.forum_cancel()}
              </button>
            {/if}
          </div>
        </div>
      {/if}
    </article>

    <section class="replies-section" aria-labelledby="replies-title">
      <h2 id="replies-title" class="replies-title">{replyCountLabel}</h2>

      {#if repliesError}
        <ErrorAlert message={repliesError} />
      {/if}

      {#if replies.length === 0}
        <p class="body-text replies-empty">{m.forum_replies_empty()}</p>
      {:else}
        <div class="reply-tree">
          {#each replies as reply (reply.id)}
            <ForumReplyItem
              {reply}
              currentUsername={$user?.username ?? null}
              threadAuthor={thread.author?.username ?? null}
              onReply={setReplyTarget}
              onDelete={onDeleteReply}
            />
          {/each}
        </div>
      {/if}
    </section>

    <section class="reply-composer">
      {#if replyTarget}
        <div class="reply-target">
          <span>{m.forum_replying_to({ username: replyTarget.username })}</span>
          <button type="button" class="quiet-btn" onclick={() => (replyTarget = null)}>
            {m.forum_cancel()}
          </button>
        </div>
      {/if}

      {#if postingStatus === 'guest'}
        <GuestNotice class="body-text" message={m.forum_guest_notice()} />
      {:else if postingStatus === 'unverified'}
        <div class="notice gate-notice">
          <p class="gate-note">{m.forum_email_unverified_notice()}</p>
          <a class="link" href={href('/settings')}>{m.nav_settings()}</a>
        </div>
      {:else if postingStatus === 'ready'}
        <form class="composer-form" onsubmit={submitReply}>
          <label class="field">
            <span class="label-text">{m.forum_reply_label()}</span>
            <textarea
              class="textarea composer-textarea"
              maxlength="10000"
              placeholder={m.forum_reply_placeholder()}
              bind:value={replyBody}></textarea>
          </label>

          {#if replyError}
            <ErrorAlert message={replyError} />
          {/if}

          <div class="delete-row">
            <button type="submit" class="btn-primary" disabled={replySubmitting}>
              {replySubmitting ? m.settings_saving() : m.forum_reply_submit()}
            </button>
          </div>
        </form>
      {:else}
        <LoadingSpinner />
      {/if}
    </section>
  {/if}
</div>

<style>
  .thread-crumbs {
    margin-bottom: var(--block-gap);
  }

  .crumb {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--text-secondary);
    font-size: var(--text-sm);
    text-decoration: none;
    transition: color var(--transition-fast);
  }

  .crumb:hover {
    color: var(--accent);
  }

  .state-block {
    padding: var(--space-8) 0;
  }

  .state-title {
    margin: 0 0 var(--space-2);
    font-size: var(--text-lg);
    font-weight: 600;
  }

  .state-note {
    margin: 0 0 var(--space-4);
  }

  .thread-error {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    align-items: flex-start;
  }

  .thread-skeleton {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }

  .skeleton-title {
    height: 1.25rem;
    width: min(26rem, 85%);
  }

  .skeleton-line {
    height: 0.8rem;
    width: 100%;
  }

  .skeleton-line.short {
    width: 60%;
  }

  .thread-main {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  .thread-head {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .thread-cat {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    margin: 0;
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-secondary);
  }

  .thread-cat .chip-dot {
    background: var(--accent);
  }

  .thread-title {
    margin: 0;
    font-size: var(--text-2xl);
    line-height: var(--leading-tight);
    font-weight: 600;
    letter-spacing: -0.015em;
  }

  .thread-meta {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.4rem 0.9rem;
    margin: 0;
    font-size: var(--text-xs);
    color: var(--text-muted);
  }

  .thread-author {
    color: var(--text-primary);
    font-weight: 600;
  }

  .callsign {
    font-family: var(--font-mono);
    letter-spacing: 0.04em;
    color: var(--text-secondary);
  }

  /* The post body is the only long-form text on the page: 16px on a relaxed
     leading, kept on the same measure as the replies below it. */
  .thread-body {
    margin: 0;
    padding-top: var(--space-4);
    border-top: 1px solid var(--border);
    font-size: var(--text-base);
    line-height: var(--leading-relaxed);
    color: var(--text-primary);
    white-space: pre-wrap;
    overflow-wrap: anywhere;
  }

  .thread-actions {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    padding-top: var(--space-4);
    border-top: 1px solid var(--border);
  }

  .delete-row {
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }

  .replies-section {
    margin-top: var(--section-gap);
  }

  .replies-title {
    margin: 0 0 var(--space-4);
    padding-bottom: var(--space-3);
    border-bottom: 1px solid var(--border);
    font-size: var(--text-sm);
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-muted);
  }

  .replies-empty {
    margin: 0;
  }

  .reply-tree {
    display: flex;
    flex-direction: column;
  }

  .reply-composer {
    margin-top: var(--section-gap);
    padding-top: var(--space-4);
    border-top: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  .reply-target {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    background-color: var(--bg-inset);
    font-size: var(--text-sm);
    color: var(--text-secondary);
  }

  .quiet-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.15rem 0.3rem;
    border: none;
    background: transparent;
    color: var(--text-muted);
    font-size: var(--text-xs);
    font-weight: 500;
    cursor: pointer;
    transition: color var(--transition-fast);
  }

  .quiet-btn:hover {
    color: var(--text-primary);
  }

  .gate-notice {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-2);
  }

  .gate-note {
    margin: 0;
  }

  .composer-form {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  .composer-textarea {
    min-height: 8rem;
    font-family: var(--font-ui);
    resize: vertical;
  }
</style>

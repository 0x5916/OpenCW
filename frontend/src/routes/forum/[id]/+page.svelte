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
  import { ArrowLeft, MessageSquare } from '@lucide/svelte';
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

<div class="thread-page page-wide">
  <a class="back-link" href={href('/forum')}>
    <ArrowLeft size={16} aria-hidden="true" />
    {m.forum_back_to_threads()}
  </a>

  {#if loading}
    <LoadingSpinner />
  {:else if notFound}
    <section class="panel state-panel">
      <h1 class="state-title">{m.forum_thread_not_found_title()}</h1>
      <p class="body-text state-note">{m.forum_thread_not_found_body()}</p>
      <a class="btn-primary state-action" href={href('/forum')}>{m.forum_back_to_threads()}</a>
    </section>
  {:else if loadFailed}
    <div class="thread-error">
      <ErrorAlert message={loadFailed} />
      <button type="button" class="btn-ghost retry-btn" onclick={retryLoad}>
        {m.forum_retry()}
      </button>
    </div>
  {:else if thread}
    <article class="panel thread-main">
      <div class="thread-head">
        <span class="category-badge">{forumCategoryLabel(thread.category)}</span>
        <h1 class="thread-title">{thread.title}</h1>
        <div class="thread-meta">
          <span>{authorLabel(thread.author)}</span>
          {#if thread.author.call_sign}
            <span class="callsign">{thread.author.call_sign}</span>
          {/if}
          <span>{formatDate(thread.created_at)}</span>
        </div>
      </div>
      <p class="thread-body">{thread.body}</p>

      {#if isThreadAuthor}
        <div class="thread-actions">
          {#if deleteError}
            <ErrorAlert message={deleteError} />
          {/if}
          <div class="delete-row">
            <button
              type="button"
              class="btn-danger delete-btn"
              disabled={deletingThread}
              onclick={deleteThread}
            >
              {threadDeleteArmed ? m.forum_delete_confirm() : m.forum_delete_thread()}
            </button>
            {#if threadDeleteArmed}
              <button
                type="button"
                class="btn-ghost cancel-delete-btn"
                onclick={() => (threadDeleteArmed = false)}
              >
                {m.forum_cancel()}
              </button>
            {/if}
          </div>
        </div>
      {/if}
    </article>

    <section class="replies-section">
      <h2 class="replies-title">
        <MessageSquare size={18} aria-hidden="true" />
        {m.forum_replies_label()}
      </h2>

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
              onReply={setReplyTarget}
              onDelete={onDeleteReply}
            />
          {/each}
        </div>
      {/if}
    </section>

    <section class="panel reply-composer">
      {#if replyTarget}
        <div class="reply-target">
          <span>{m.forum_replying_to({ username: replyTarget.username })}</span>
          <button type="button" class="action-btn" onclick={() => (replyTarget = null)}>
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
          <textarea
            class="input composer-textarea"
            maxlength="10000"
            placeholder={m.forum_reply_placeholder()}
            bind:value={replyBody}></textarea>

          {#if replyError}
            <ErrorAlert message={replyError} />
          {/if}

          <button type="submit" class="btn-primary" disabled={replySubmitting}>
            {replySubmitting ? m.settings_saving() : m.forum_reply_submit()}
          </button>
        </form>
      {:else}
        <LoadingSpinner />
      {/if}
    </section>
  {/if}
</div>

<style>
  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    margin-bottom: var(--block-gap);
    color: var(--text-secondary);
    font-size: 0.875rem;
    text-decoration: none;
    transition: color var(--transition-fast);
  }

  .back-link:hover {
    color: var(--accent);
  }

  .state-panel {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    padding: 2.5rem 1.5rem;
    text-align: center;
  }

  .state-title {
    margin: 0;
    font-size: 1.35rem;
    font-weight: 700;
  }

  .state-note {
    max-width: 34rem;
    margin: 0;
  }

  .state-action {
    display: inline-flex;
    width: auto;
    align-items: center;
    gap: 0.5rem;
    margin-top: 0.5rem;
    text-decoration: none;
  }

  .thread-error {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    align-items: flex-start;
  }

  .retry-btn {
    width: auto;
  }

  .thread-main {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .thread-head {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .category-badge {
    align-self: flex-start;
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
    font-size: 1.5rem;
    line-height: 1.3;
    font-weight: 700;
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
    font-size: 0.72rem;
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    color: var(--accent);
  }

  .thread-body {
    margin: 0;
    max-width: var(--max-width-narrow);
    color: var(--text-secondary);
    white-space: pre-wrap;
    overflow-wrap: anywhere;
  }

  .thread-actions {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    border-top: 1px solid var(--border);
    padding-top: 1rem;
  }

  .delete-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .delete-btn,
  .cancel-delete-btn {
    flex: 0 0 auto;
    width: auto;
  }

  .replies-section {
    margin-top: var(--block-gap);
  }

  .replies-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 0 0 1rem;
    font-size: 1.1rem;
    font-weight: 700;
  }

  .replies-empty {
    margin: 0;
  }

  .reply-tree {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .reply-composer {
    margin-top: var(--block-gap);
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .reply-target {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 0.5rem 0.75rem;
    border-radius: var(--radius-md);
    background-color: var(--bg-inset);
    border: 1px solid var(--border-card);
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .action-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.2rem 0.4rem;
    border: none;
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--text-muted);
    font-size: 0.82rem;
    cursor: pointer;
    transition: color var(--transition-fast);
  }

  .action-btn:hover {
    color: var(--accent);
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

  .composer-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .composer-textarea {
    min-height: 6rem;
    font-family: inherit;
    resize: vertical;
  }
</style>

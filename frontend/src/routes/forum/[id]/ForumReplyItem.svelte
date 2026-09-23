<script lang="ts">
  import { CornerDownRight } from '@lucide/svelte';
  import type { ForumReply } from '$lib/api';
  import { authorLabel, formatDate } from '$lib/format';
  import * as m from '$lib/paraglide/messages';
  import Self from './ForumReplyItem.svelte';

  interface Props {
    reply: ForumReply;
    currentUsername: string | null;
    /** Username of the thread author, used to mark their replies. */
    threadAuthor?: string | null;
    /** Display name of the reply this one answers; shown when nesting is flat. */
    parentLabel?: string | null;
    /** Nesting level used to cap the visual indentation. */
    depth?: number;
    onReply: (reply: ForumReply) => void;
    onDelete: (replyId: string) => Promise<void>;
  }

  let {
    reply,
    currentUsername,
    threadAuthor = null,
    parentLabel = null,
    depth = 0,
    onReply,
    onDelete
  }: Props = $props();

  let confirmingDelete = $state(false);
  let deleting = $state(false);

  const isAuthor = $derived(
    !reply.is_deleted && currentUsername !== null && reply.author?.username === currentUsername
  );

  const isOp = $derived(
    !reply.is_deleted && threadAuthor !== null && reply.author?.username === threadAuthor
  );

  async function handleDelete(): Promise<void> {
    if (!confirmingDelete) {
      confirmingDelete = true;
      return;
    }

    if (deleting) return;
    deleting = true;
    try {
      // The page owns error reporting; the item is removed on success because
      // the reply tree is refetched.
      await onDelete(reply.id);
    } finally {
      deleting = false;
      confirmingDelete = false;
    }
  }
</script>

<div class="reply">
  {#if reply.is_deleted}
    <div class="reply-main is-deleted">
      <p class="reply-text deleted-text">{m.forum_reply_deleted()}</p>
    </div>
  {:else}
    <div class="reply-main">
      <div class="reply-meta">
        <span class="reply-author">{authorLabel(reply.author)}</span>
        {#if isOp}
          <span class="chip">{m.forum_op_tag()}</span>
        {/if}
        {#if reply.author?.call_sign}
          <span class="callsign">{reply.author.call_sign}</span>
        {/if}
        <span class="reply-date">{formatDate(reply.created_at)}</span>
      </div>
      {#if depth >= 3 && parentLabel}
        <p class="reply-context">{m.forum_replying_to({ username: parentLabel })}</p>
      {/if}
      <p class="reply-text">{reply.body}</p>
      <div class="reply-actions">
        <button type="button" class="quiet-btn" onclick={() => onReply(reply)}>
          <CornerDownRight size={14} aria-hidden="true" />
          {m.forum_reply_button()}
        </button>
        {#if isAuthor}
          {#if confirmingDelete}
            <button
              type="button"
              class="quiet-btn is-danger"
              disabled={deleting}
              onclick={handleDelete}
            >
              {m.forum_delete_confirm()}
            </button>
            <button type="button" class="quiet-btn" onclick={() => (confirmingDelete = false)}>
              {m.forum_cancel()}
            </button>
          {:else}
            <button type="button" class="quiet-btn is-danger" onclick={handleDelete}>
              {m.forum_delete_reply()}
            </button>
          {/if}
        {/if}
      </div>
    </div>
  {/if}

  {#if reply.children.length > 0}
    <div class="reply-children" class:flattened={depth >= 3}>
      {#each reply.children as child (child.id)}
        <Self
          reply={child}
          {currentUsername}
          {threadAuthor}
          depth={depth + 1}
          parentLabel={reply.is_deleted ? null : authorLabel(reply.author)}
          {onReply}
          {onDelete}
        />
      {/each}
    </div>
  {/if}
</div>

<style>
  /* Replies are a thread, not a grid of cards: hairline separators between
     replies, and one left rule per nesting level. */
  .reply {
    padding: var(--space-4) 0;
    border-bottom: 1px solid var(--border);
  }

  .reply-main {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .reply-main.is-deleted {
    opacity: 0.7;
  }

  .deleted-text {
    margin: 0;
    color: var(--text-muted);
    font-style: italic;
    font-size: var(--text-sm);
  }

  .reply-meta {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.4rem 0.7rem;
    font-size: var(--text-xs);
    color: var(--text-muted);
  }

  .reply-author {
    color: var(--text-primary);
    font-weight: 600;
  }

  .reply-context {
    margin: 0;
    font-size: var(--text-xs);
    color: var(--text-muted);
  }

  .callsign {
    font-family: var(--font-mono);
    letter-spacing: 0.04em;
    color: var(--text-secondary);
  }

  .reply-text {
    margin: 0;
    font-size: var(--text-sm);
    line-height: var(--leading-normal);
    color: var(--text-primary);
    white-space: pre-wrap;
    overflow-wrap: anywhere;
  }

  .reply-actions {
    display: flex;
    align-items: center;
    gap: var(--space-3);
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

  .quiet-btn.is-danger:hover {
    color: var(--danger);
  }

  /* Nesting: one hairline rule per level, capped at three; deeper replies are
     flattened and name the reply they answer instead. */
  .reply-children {
    display: flex;
    flex-direction: column;
    margin-top: var(--space-1);
    margin-left: var(--space-4);
    padding-left: var(--space-4);
    border-left: 1px solid var(--border);
  }

  .reply-children.flattened {
    margin-left: 0;
    padding-left: 0;
    border-left: none;
  }
</style>

<script lang="ts">
  import { CornerDownRight } from '@lucide/svelte';
  import type { ForumReply } from '$lib/api';
  import { authorLabel, formatDate } from '$lib/format';
  import * as m from '$lib/paraglide/messages';
  import Self from './ForumReplyItem.svelte';

  interface Props {
    reply: ForumReply;
    currentUsername: string | null;
    /** Nesting level used to cap the visual indentation. */
    depth?: number;
    onReply: (reply: ForumReply) => void;
    onDelete: (replyId: string) => Promise<void>;
  }

  let { reply, currentUsername, depth = 0, onReply, onDelete }: Props = $props();

  let confirmingDelete = $state(false);
  let deleting = $state(false);

  const isAuthor = $derived(
    !reply.is_deleted && currentUsername !== null && reply.author?.username === currentUsername
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
        {#if reply.author?.call_sign}
          <span class="callsign">{reply.author.call_sign}</span>
        {/if}
        <span class="reply-date">{formatDate(reply.created_at)}</span>
      </div>
      <p class="reply-text">{reply.body}</p>
      <div class="reply-actions">
        <button type="button" class="action-btn" onclick={() => onReply(reply)}>
          <CornerDownRight size={14} aria-hidden="true" />
          {m.forum_reply_button()}
        </button>
        {#if isAuthor}
          {#if confirmingDelete}
            <button
              type="button"
              class="action-btn is-danger"
              disabled={deleting}
              onclick={handleDelete}
            >
              {m.forum_delete_confirm()}
            </button>
            <button type="button" class="action-btn" onclick={() => (confirmingDelete = false)}>
              {m.forum_cancel()}
            </button>
          {:else}
            <button type="button" class="action-btn is-danger" onclick={handleDelete}>
              {m.forum_delete_reply()}
            </button>
          {/if}
        {/if}
      </div>
    </div>
  {/if}

  {#if reply.children.length > 0}
    <div class="reply-children" class:no-indent={depth >= 3}>
      {#each reply.children as child (child.id)}
        <Self reply={child} {currentUsername} depth={depth + 1} {onReply} {onDelete} />
      {/each}
    </div>
  {/if}
</div>

<style>
  .reply-main {
    padding: 0.85rem 1rem;
    border: 1px solid var(--border-card);
    border-radius: var(--radius-md);
    background-color: var(--bg-surface);
  }

  .reply-main.is-deleted {
    background-color: var(--bg-inset);
  }

  .deleted-text {
    margin: 0;
    color: var(--text-muted);
    font-style: italic;
  }

  .reply-meta {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-wrap: wrap;
    margin-bottom: 0.35rem;
    color: var(--text-muted);
    font-size: 0.82rem;
  }

  .reply-author {
    color: var(--text-primary);
    font-weight: 600;
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

  .reply-text {
    margin: 0;
    color: var(--text-secondary);
    white-space: pre-wrap;
    overflow-wrap: anywhere;
  }

  .reply-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-top: 0.5rem;
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

  .action-btn.is-danger:hover {
    color: var(--danger);
  }

  .reply-children {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 0.75rem;
    margin-left: 1.25rem;
    padding-left: 0.75rem;
    border-left: 1px solid var(--border-subtle);
  }

  .reply-children.no-indent {
    margin-left: 0;
    padding-left: 0;
    border-left-color: transparent;
  }
</style>

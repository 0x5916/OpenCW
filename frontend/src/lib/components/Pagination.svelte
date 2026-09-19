<script lang="ts">
  // Not mounted anywhere yet: the forum is still a placeholder, so this ships
  // ahead of the thread list it will paginate.
  import { ChevronLeft, ChevronRight } from '@lucide/svelte';
  import * as m from '$lib/paraglide/messages';

  interface Props {
    currentPage: number;
    totalPages: number;
    prevHref: string;
    nextHref: string;
  }

  let { currentPage, totalPages, prevHref, nextHref }: Props = $props();
</script>

{#if totalPages > 1}
  <nav class="pagination" aria-label={m.forum_pagination()}>
    {#if currentPage > 1}
      <a
        class="icon-link"
        href={prevHref}
        aria-label={m.forum_previous_page()}
        title={m.forum_previous_page()}><ChevronLeft size={18} /></a
      >
    {/if}
    <span>{m.forum_page_of({ page: currentPage, total: totalPages })}</span>
    {#if currentPage < totalPages}
      <a
        class="icon-link"
        href={nextHref}
        aria-label={m.forum_next_page()}
        title={m.forum_next_page()}><ChevronRight size={18} /></a
      >
    {/if}
  </nav>
{/if}

<style>
  .pagination {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    margin: 1.25rem 0;
    color: var(--text-secondary);
    font-size: 0.85rem;
  }

  .icon-link {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 2.25rem;
    height: 2.25rem;
    color: var(--accent);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    text-decoration: none;
  }

  .icon-link:hover {
    background: var(--bg-inset);
  }
</style>

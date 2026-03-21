<script lang="ts">
  import { MessageSquare } from 'lucide-svelte';
  import * as m from '$lib/paraglide/messages';
</script>

<main class="placeholder-page">
  <header class="placeholder-header">
    <div class="placeholder-icon" aria-hidden="true">
      <MessageSquare size={48} />
    </div>
    <h1 class="page-title">{m.forum_title()}</h1>
  </header>
  <p class="body-text placeholder-message">{m.forum_wip()}</p>
</main>

<style>
  .placeholder-page {
    padding: 3rem 1rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: calc(100vh - 6rem);
  }

  .placeholder-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1.5rem;
    text-align: center;
  }

  .placeholder-header :global(.page-title) {
    margin: 0;
  }

  .placeholder-icon {
    color: var(--accent);
    font-size: 3rem;
  }

  .placeholder-message {
    margin: 0;
    color: var(--text-secondary);
    max-width: 400px;
    text-align: center;
  }
</style>

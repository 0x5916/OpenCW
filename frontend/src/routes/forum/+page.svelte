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

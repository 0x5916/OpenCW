<script lang="ts">
  import { Wrench } from 'lucide-svelte';
  import * as m from '$lib/paraglide/messages';
</script>

<main class="placeholder-page">
  <header class="placeholder-header">
    <div class="placeholder-icon" aria-hidden="true">
      <Wrench size={48} />
    </div>
    <h1 class="page-title">{m.about_title()}</h1>
  </header>
  <p class="body-text placeholder-message">{m.about_wip()}</p>
</main>

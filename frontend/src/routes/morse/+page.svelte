<script lang="ts">
  import { localizedHref as href } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';
</script>

<svelte:head>
  <meta name="robots" content="noindex,follow" />
</svelte:head>

<!-- The load function redirects before this renders; the link is a fallback
     for the rare case where neither the meta refresh nor JavaScript fires. -->
<div class="page-narrow">
  <p class="body-text">
    <a href={href('/learn')}>{m.home_cta()}</a>
  </p>
</div>

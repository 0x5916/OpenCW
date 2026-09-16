<script lang="ts">
  import { localizedHref } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';

  interface Props {
    /** Optional classes for the wrapping <p> (e.g. `body-text`). */
    class?: string;
  }

  let { class: className }: Props = $props();
</script>

<p class={className}>
  {m.trainer_guest_notice()}
  <a href={localizedHref('/login')} class="link">{m.nav_login()}</a>
  /
  <a href={localizedHref('/register')} class="link">{m.nav_register()}</a>
</p>

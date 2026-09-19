<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    title: string;
    subtitle: string;
    footer: Snippet;
    children: Snippet;
  }

  let { title, subtitle, footer, children }: Props = $props();
</script>

<div class="page-narrow auth-page">
  <div class="card auth-card">
    <h1 class="page-title auth-title">{title}</h1>
    <p class="body-text auth-subtitle">{subtitle}</p>

    {@render children()}

    <hr class="divider" />
    <p class="body-text auth-footer-text">{@render footer()}</p>
  </div>
</div>

<style>
  .auth-card {
    max-width: 38rem;
    margin-left: auto;
    margin-right: auto;
    padding: 1.75rem;
  }

  /* `.page-title` already carries the h1 scale; the card only owns the gap. */
  .auth-title {
    margin-bottom: 0.35rem;
  }

  .auth-subtitle {
    margin-bottom: 1.25rem;
  }

  .auth-footer-text {
    margin: 0;
    text-align: center;
  }

  @media (max-width: 639px) {
    .auth-card {
      padding: 1.25rem;
    }
  }
</style>

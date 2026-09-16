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

<main class="page-narrow auth-page">
  <div class="card auth-card">
    <h1 class="page-title auth-title">{title}</h1>
    <p class="body-text auth-subtitle">{subtitle}</p>

    {@render children()}

    <hr class="divider" />
    <p class="body-text auth-footer-text">{@render footer()}</p>
  </div>
</main>

<style>
  .auth-page {
    padding-top: 1.5rem;
    padding-bottom: 1.5rem;
  }

  .auth-card {
    max-width: 38rem;
    margin-left: auto;
    margin-right: auto;
    padding: 1.75rem;
  }

  .auth-title {
    font-size: 2rem;
    line-height: 2.25rem;
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

    .auth-title {
      font-size: 1.75rem;
      line-height: 2rem;
    }
  }
</style>

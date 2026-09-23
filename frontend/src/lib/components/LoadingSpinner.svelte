<script lang="ts">
  import * as m from '$lib/paraglide/messages';

  interface Props {
    variant?: 'spinner' | 'dots';
    padded?: boolean;
    /** Accessible status text; override where a more specific state exists. */
    label?: string;
  }

  let { variant = 'spinner', padded = true, label = '' }: Props = $props();

  let statusLabel = $derived(label || m.common_loading());
</script>

{#if padded}
  <div class="loading-container" role="status" aria-live="polite">
    <span class="loading-{variant}" aria-hidden="true"></span>
    <span class="sr-only">{statusLabel}</span>
  </div>
{:else}
  <span class="loading-inline" role="status" aria-live="polite">
    <span class="loading-{variant}" aria-hidden="true"></span>
    <span class="sr-only">{statusLabel}</span>
  </span>
{/if}

<style>
  .loading-container {
    display: flex;
    justify-content: center;
    padding: 2rem;
  }

  .loading-inline {
    display: inline-flex;
    vertical-align: middle;
  }

  .loading-spinner {
    width: 1.5rem;
    height: 1.5rem;
    border: 2px solid var(--border);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .loading-dots::after {
    content: '.';
    animation: dots 1.2s steps(4, end) infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  @keyframes dots {
    0%,
    100% {
      content: '.';
    }
    33% {
      content: '..';
    }
    66% {
      content: '...';
    }
  }
</style>

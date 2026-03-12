<script lang="ts">
  let { variant = 'spinner', padded = true } = $props();
</script>

{#if padded}
  <div class="loading-container">
    <span class="loading-{variant}"></span>
  </div>
{:else}
  <span class="loading-{variant}"></span>
{/if}

<style>
  .loading-container {
    display: flex;
    justify-content: center;
    padding: 2rem;
  }

  .loading-spinner {
    width: 2rem;
    height: 2rem;
    border: 3px solid var(--border);
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

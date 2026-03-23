<svelte:head>
  <title>Offline | OpenCW</title>
</svelte:head>

<main class="offline-page card-sm" aria-live="polite">
  <h1 class="offline-title">You're offline</h1>
  <p class="offline-copy">
    OpenCW will keep working for cached pages and queue your training progress for upload when your
    connection is restored.
  </p>
  <p class="offline-copy offline-tip">Reconnect to sync profile data and cloud settings.</p>
</main>

<style>
  .offline-page {
    max-width: 40rem;
    margin: 3rem auto;
    text-align: center;
  }

  .offline-title {
    margin: 0 0 0.75rem;
    font-size: 1.75rem;
    line-height: 1.2;
    color: var(--accent);
  }

  .offline-copy {
    margin: 0.5rem 0;
    color: var(--text-secondary);
  }

  .offline-tip {
    color: var(--text-muted);
    font-size: 0.95rem;
  }
</style>

<script lang="ts">
  interface Props {
    message: string;
    variant?: 'error' | 'warning' | 'info';
  }

  let { message, variant = 'error' }: Props = $props();
</script>

<div class="error-alert error-alert-{variant}">
  <p class="error-message">⚠ {message}</p>
</div>

<style>
  .error-alert {
    padding: 1rem;
    border-radius: var(--radius-md);
    border: 1px solid;
  }

  .error-alert-error {
    background: var(--status-bad-tint);
    border-color: color-mix(in srgb, var(--status-bad) 30%, transparent);
    color: var(--status-bad);
  }

  .error-alert-warning {
    background: var(--status-ok-tint);
    border-color: color-mix(in srgb, var(--status-ok) 30%, transparent);
    color: var(--status-ok);
  }

  /* No dedicated info role exists, so the informational banner borrows the
     accent (which is already tuned per theme). */
  .error-alert-info {
    background: color-mix(in srgb, var(--accent) 12%, transparent);
    border-color: color-mix(in srgb, var(--accent) 30%, transparent);
    color: var(--accent);
  }

  .error-message {
    margin: 0;
    font-size: 0.9375rem;
  }
</style>

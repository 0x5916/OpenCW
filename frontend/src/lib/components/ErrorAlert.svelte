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
    border-radius: 0.75rem;
    border: 1px solid;
  }

  .error-alert-error {
    background: rgba(239, 68, 68, 0.1);
    border-color: rgba(239, 68, 68, 0.3);
    color: #ef4444;
  }

  .error-alert-warning {
    background: rgba(234, 179, 8, 0.1);
    border-color: rgba(234, 179, 8, 0.3);
    color: #ca8a04;
  }

  .error-alert-info {
    background: rgba(59, 130, 246, 0.1);
    border-color: rgba(59, 130, 246, 0.3);
    color: #3b82f6;
  }

  .error-message {
    margin: 0;
    font-size: 0.9375rem;
  }
</style>

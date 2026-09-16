<script lang="ts">
  import * as m from '$lib/paraglide/messages';

  interface Props {
    saving: boolean;
    saved: boolean;
  }

  let { saving, saved }: Props = $props();

  let label = $derived(
    saved ? m.settings_saved() : saving ? m.settings_saving() : m.settings_update()
  );
</script>

<button type="submit" class="btn-primary settings-btn-compact" disabled={saving}>
  {label}
</button>

<style>
  .settings-btn-compact {
    padding: var(--space-sm) 0.9rem;
    font-size: 0.875rem;
    line-height: 1.25rem;
    min-height: calc(1.25rem + (var(--space-sm) * 2) + 2px);
    width: auto;
    white-space: nowrap;
  }

  @media (max-width: 720px) {
    .settings-btn-compact {
      width: 100%;
    }
  }
</style>

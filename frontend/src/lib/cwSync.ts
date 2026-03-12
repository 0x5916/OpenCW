import { saveCWSettings, getCWSettings, type CWSettings } from '$lib/api';
export type { CWSettings };

/**
 * Consolidates common sync/restore patterns for CW settings
 * These are utility functions used in learn and settings pages for consistency
 */

/**
 * Sync CW settings to account with error handling
 * Callbacks manage the component state (syncing, synced, syncError)
 */
export async function performSync(
  settings: CWSettings,
  onSuccess: () => void,
  onError: (msg: string) => void
) {
  try {
    await saveCWSettings(settings);
    onSuccess();
  } catch (error) {
    const msg = error instanceof Error ? error.message : 'Failed to sync settings';
    onError(msg);
  }
}

/**
 * Restore CW settings from account with error handling
 * Callbacks manage the component state (restoring, restoreError)
 */
export async function performRestore(
  onSuccess: (cw: CWSettings) => void,
  onError: (msg: string) => void
) {
  try {
    const cw = await getCWSettings();
    onSuccess(cw);
  } catch (error) {
    const msg = error instanceof Error ? error.message : 'Failed to load settings';
    onError(msg);
  }
}

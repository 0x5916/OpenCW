import type { PageServerLoad } from './$types';
import { CW_STORAGE_KEYS } from '$lib/storageKeys';

export const load: PageServerLoad = ({ cookies }) => {
  const stored = cookies.get(CW_STORAGE_KEYS.lesson);
  const lesson = stored ? Number(stored) || 1 : 1;
  return { lesson };
};

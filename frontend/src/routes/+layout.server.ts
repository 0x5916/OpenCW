import type { LayoutServerLoad } from './$types';
import { isLocale, baseLocale } from '$lib/paraglide/runtime';

export const load: LayoutServerLoad = ({ cookies }) => {
  const stored = cookies.get('PARAGLIDE_LOCALE');
  const locale = stored && isLocale(stored) ? stored : baseLocale;
  return { locale };
};

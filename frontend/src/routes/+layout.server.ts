import type { LayoutServerLoad } from './$types';
import { LOCALE_COOKIE, normalizeLocalePreference, type LocalePreference } from '$lib/locale';
import { getLocale } from '$lib/paraglide/runtime';

export const load: LayoutServerLoad = ({ cookies }) => {
  const storedPreference = cookies.get(LOCALE_COOKIE);
  const localePreference = normalizeLocalePreference(storedPreference);
  const locale = getLocale();

  return { locale, localePreference: localePreference as LocalePreference };
};

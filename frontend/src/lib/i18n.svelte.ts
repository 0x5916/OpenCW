import { baseLocale, isLocale, locales, overwriteGetLocale } from '$lib/paraglide/runtime';

export type Locale = (typeof locales)[number];

const COOKIE_NAME = 'PARAGLIDE_LOCALE';

export const lang = $state<{ value: Locale }>({ value: baseLocale as Locale });

export function initLang(initialLocale: string) {
  lang.value = isLocale(initialLocale)
    ? (initialLocale as Locale)
    : (baseLocale as Locale);
  overwriteGetLocale(() => lang.value);
}

export function setLang(l: Locale) {
  lang.value = l;
  document.cookie = `${COOKIE_NAME}=${l}; path=/; max-age=31536000; SameSite=Lax`;
}

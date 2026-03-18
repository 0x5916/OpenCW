import { baseLocale, localizeHref, overwriteGetLocale } from '$lib/paraglide/runtime';
import {
  LOCALE_COOKIE,
  LOCALE_PREFERENCE_STORAGE_KEY,
  matchLocaleCandidate,
  normalizeLocalePreference,
  detectLocaleFromAcceptedLanguages,
  type Locale,
  type LocalePreference
} from '$lib/locale';

const ONE_YEAR_SECONDS = 31536000;

export const lang = $state<{ value: Locale }>({ value: baseLocale as Locale });
export const langPreference = $state<{ value: LocalePreference }>({ value: 'auto' });

export type { Locale, LocalePreference };

function applyDocumentLocale(locale: Locale): void {
  if (typeof document === 'undefined') return;
  document.documentElement.lang = locale;
}

function browserLanguages(): readonly string[] {
  if (typeof navigator === 'undefined') return [];
  if (Array.isArray(navigator.languages) && navigator.languages.length > 0) {
    return navigator.languages;
  }
  return navigator.language ? [navigator.language] : [];
}

function applyPreference(preference: LocalePreference): void {
  langPreference.value = preference;
  lang.value =
    preference === 'auto'
      ? detectLocaleFromAcceptedLanguages(browserLanguages())
      : preference;
  overwriteGetLocale(() => lang.value);
  applyDocumentLocale(lang.value);
}

function navigateToLocalizedPathIfNeeded(): void {
  if (typeof window === 'undefined') return;

  const current = `${window.location.pathname}${window.location.search}${window.location.hash}`;
  const next = localizeHref(current, { locale: lang.value });

  const currentUrl = new URL(window.location.href);
  const nextUrl = new URL(next, window.location.origin);

  if (currentUrl.href !== nextUrl.href) {
    window.location.assign(nextUrl.href);
  }
}

export function initLang(initialLocale: string, initialPreference: string = 'auto') {
  const preference = normalizeLocalePreference(initialPreference);
  const initial = matchLocaleCandidate(initialLocale) ?? (baseLocale as Locale);

  if (preference === 'auto') {
    if (typeof window === 'undefined') {
      // Keep SSR locale for hydration consistency; browser detection is client-only.
      langPreference.value = 'auto';
      lang.value = initial;
      overwriteGetLocale(() => lang.value);
      applyDocumentLocale(lang.value);
      return;
    }

    applyPreference('auto');
    navigateToLocalizedPathIfNeeded();
    return;
  }

  applyPreference(preference);
}

export function setLangPreference(preference: LocalePreference, options: { navigate?: boolean } = {}) {
  applyPreference(preference);
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(LOCALE_PREFERENCE_STORAGE_KEY, preference);
  }
  if (typeof document !== 'undefined') {
    document.cookie = `${LOCALE_COOKIE}=${preference}; path=/; max-age=${ONE_YEAR_SECONDS}; SameSite=Lax`;
  }

  if (options.navigate !== false) {
    navigateToLocalizedPathIfNeeded();
  }
}

export function setLang(locale: Locale) {
  setLangPreference(locale);
}

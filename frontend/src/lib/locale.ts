import { baseLocale, isLocale, locales } from '$lib/paraglide/runtime';

export type Locale = (typeof locales)[number];
export type LocalePreference = Locale | 'auto';

export const LOCALE_COOKIE = 'PARAGLIDE_LOCALE';
export const LOCALE_PREFERENCE_STORAGE_KEY = 'PARAGLIDE_LOCALE_PREF';

export const LOCALE_DISPLAY: Record<Locale, { short: string; native: string; english: string }> = {
  en: { short: 'EN', native: 'English', english: 'English' },
  de: { short: 'DE', native: 'Deutsch', english: 'German' },
  ja: { short: 'JA', native: '日本語', english: 'Japanese' },
  'zh-Hans': { short: '简', native: '简体中文', english: 'Chinese (Simplified)' },
  'zh-Hant': { short: '繁', native: '繁體中文', english: 'Chinese (Traditional)' }
};

type WeightedLanguage = { value: string; q: number };

function parseAcceptLanguage(header: string): WeightedLanguage[] {
  return header
    .split(',')
    .map((part) => {
      const [rawValue, ...params] = part.trim().split(';');
      const value = rawValue?.trim();
      if (!value) return null;

      const qParam = params.map((p) => p.trim()).find((p) => p.startsWith('q='));
      const q = qParam ? Number.parseFloat(qParam.slice(2)) : 1;

      return {
        value,
        q: Number.isFinite(q) ? q : 1
      };
    })
    .filter((item): item is WeightedLanguage => item !== null)
    .sort((a, b) => b.q - a.q);
}

export function normalizeLocalePreference(value: string | null | undefined): LocalePreference {
  if (!value || value === 'auto') return 'auto';
  if (isLocale(value)) return value as Locale;

  const matched = matchLocaleCandidate(value);
  return matched ?? 'auto';
}

export function matchLocaleCandidate(input: string | null | undefined): Locale | null {
  if (!input) return null;

  const cleaned = input.trim();
  if (!cleaned) return null;

  if (isLocale(cleaned)) {
    return cleaned as Locale;
  }

  const lower = cleaned.toLowerCase();
  const exact = locales.find((locale) => locale.toLowerCase() === lower);
  if (exact) return exact as Locale;

  const base = lower.split('-')[0];
  const primaryMatch = locales.find((locale) => locale.toLowerCase().split('-')[0] === base);
  return primaryMatch ? (primaryMatch as Locale) : null;
}

export function detectLocaleFromAcceptedLanguages(candidates: readonly string[]): Locale {
  for (const candidate of candidates) {
    const matched = matchLocaleCandidate(candidate);
    if (matched) return matched;
  }

  return baseLocale as Locale;
}

export function detectLocaleFromAcceptLanguageHeader(header: string | null | undefined): Locale {
  if (!header) return baseLocale as Locale;
  const accepted = parseAcceptLanguage(header).map((entry) => entry.value);
  return detectLocaleFromAcceptedLanguages(accepted);
}

export function resolveLocalePreference(
  preference: LocalePreference,
  acceptedLanguages: readonly string[] = []
): Locale {
  if (preference !== 'auto') return preference;
  return detectLocaleFromAcceptedLanguages(acceptedLanguages);
}

export function getLocaleShortLabel(locale: Locale): string {
  return LOCALE_DISPLAY[locale]?.short ?? locale.toUpperCase();
}

export function getLocaleLongLabel(locale: Locale): string {
  const label = LOCALE_DISPLAY[locale];
  if (!label) return locale;
  if (label.native === label.english) return label.native;
  return `${label.native} (${label.english})`;
}

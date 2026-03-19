import {
  getSettings,
  saveCWSettings,
  savePageSettings,
  type CWSettings,
  type PageSettings
} from '$lib/api';
import { isLocale } from '$lib/paraglide/runtime';
import {
  LOCALE_COOKIE,
  LOCALE_PREFERENCE_STORAGE_KEY,
  normalizeLocalePreference,
  type LocalePreference
} from '$lib/locale';

export type { CWSettings, PageSettings };

const LESSON_COOKIE = 'learn.lesson';
const ONE_YEAR_SECONDS = 31536000;

export function normalizeLesson(lesson: number, maxLesson: number): number {
  if (!Number.isFinite(lesson)) return 1;
  return Math.min(Math.max(1, Math.trunc(lesson)), Math.max(1, maxLesson));
}

function readCookie(name: string): string | null {
  if (typeof document === 'undefined') return null;
  const target = `${name}=`;
  const item = document.cookie.split(';').map((part) => part.trim()).find((part) => part.startsWith(target));
  return item ? decodeURIComponent(item.slice(target.length)) : null;
}

export function readClientPageSettings(
  currentLesson: number,
  maxLesson: number,
  fallbackLanguagePreference: LocalePreference
): PageSettings {
  const themeRaw = typeof localStorage === 'undefined' ? null : localStorage.getItem('theme');
  const theme: PageSettings['theme'] =
    themeRaw === 'dark' || themeRaw === 'light' || themeRaw === 'auto' ? themeRaw : 'auto';

  const localLang =
    typeof localStorage === 'undefined' ? null : localStorage.getItem(LOCALE_PREFERENCE_STORAGE_KEY);
  const cookieLang = readCookie(LOCALE_COOKIE);
  const language = normalizeLocalePreference(localLang ?? cookieLang ?? fallbackLanguagePreference);

  return {
    theme,
    language,
    cur_lesson: normalizeLesson(currentLesson, maxLesson)
  };
}

export function applyClientPageSettings(
  page: PageSettings,
  maxLesson: number,
  onLocale: (preference: LocalePreference, options?: { navigate?: boolean }) => void,
  options: { applyLanguage?: boolean; navigate?: boolean } = {}
): number {
  const lesson = normalizeLesson(page.cur_lesson, maxLesson);
  const language = normalizeLocalePreference(page.language);
  const applyLanguage = options.applyLanguage ?? true;

  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('theme', page.theme);
    localStorage.setItem('learn.lesson', String(lesson));
    if (applyLanguage) {
      localStorage.setItem(LOCALE_PREFERENCE_STORAGE_KEY, language);
    }
  }

  if (typeof document !== 'undefined') {
    if (page.theme === 'auto') document.documentElement.removeAttribute('data-theme');
    else document.documentElement.setAttribute('data-theme', page.theme);

    document.cookie = `${LESSON_COOKIE}=${lesson}; path=/; max-age=${ONE_YEAR_SECONDS}; SameSite=Lax`;

    if (applyLanguage) {
      document.cookie = `${LOCALE_COOKIE}=${language}; path=/; max-age=${ONE_YEAR_SECONDS}; SameSite=Lax`;
    }
  }

  if (applyLanguage && (language === 'auto' || isLocale(language))) {
    onLocale(language, { navigate: options.navigate });
  }

  return lesson;
}

export async function syncSettingsToServer(cw: CWSettings, page: PageSettings): Promise<void> {
  await Promise.all([saveCWSettings(cw), savePageSettings(page)]);
}

export async function restoreSettingsFromServer(): Promise<{ cw: CWSettings; page: PageSettings }> {
  const settings = await getSettings();
  return { cw: settings.cw_settings, page: settings.page_settings };
}

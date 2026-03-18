import {
  getCWSettings,
  getPageSettings,
  saveCWSettings,
  savePageSettings,
  type CWSettings,
  type PageSettings
} from '$lib/api';
import { isLocale } from '$lib/paraglide/runtime';
import type { Locale } from '$lib/i18n.svelte';

export type { CWSettings, PageSettings };

const LESSON_COOKIE = 'learn.lesson';
const LANG_COOKIE = 'PARAGLIDE_LOCALE';
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

export function readClientPageSettings(currentLesson: number, maxLesson: number, fallbackLanguage: Locale): PageSettings {
  const themeRaw = typeof localStorage === 'undefined' ? null : localStorage.getItem('theme');
  const theme: PageSettings['theme'] =
    themeRaw === 'dark' || themeRaw === 'light' || themeRaw === 'auto' ? themeRaw : 'auto';

  const localLang = typeof localStorage === 'undefined' ? null : localStorage.getItem(LANG_COOKIE);
  const cookieLang = readCookie(LANG_COOKIE);
  const candidateLang = localLang ?? cookieLang ?? fallbackLanguage;
  const language = candidateLang === 'auto' || isLocale(candidateLang) ? candidateLang : fallbackLanguage;

  return {
    theme,
    language,
    cur_lesson: normalizeLesson(currentLesson, maxLesson)
  };
}

export function applyClientPageSettings(
  page: PageSettings,
  maxLesson: number,
  onLocale: (locale: Locale) => void
): number {
  const lesson = normalizeLesson(page.cur_lesson, maxLesson);

  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('theme', page.theme);
    localStorage.setItem('learn.lesson', String(lesson));
  }

  if (typeof document !== 'undefined') {
    if (page.theme === 'auto') document.documentElement.removeAttribute('data-theme');
    else document.documentElement.setAttribute('data-theme', page.theme);

    document.cookie = `${LESSON_COOKIE}=${lesson}; path=/; max-age=${ONE_YEAR_SECONDS}; SameSite=Lax`;

    if (page.language !== 'auto') {
      document.cookie = `${LANG_COOKIE}=${page.language}; path=/; max-age=${ONE_YEAR_SECONDS}; SameSite=Lax`;
    }
  }

  if (page.language !== 'auto' && isLocale(page.language)) {
    onLocale(page.language as Locale);
  }

  return lesson;
}

export async function syncSettingsToServer(cw: CWSettings, page: PageSettings): Promise<void> {
  await Promise.all([saveCWSettings(cw), savePageSettings(page)]);
}

export async function restoreSettingsFromServer(): Promise<{ cw: CWSettings; page: PageSettings }> {
  const [cw, page] = await Promise.all([getCWSettings(), getPageSettings()]);
  return { cw, page };
}

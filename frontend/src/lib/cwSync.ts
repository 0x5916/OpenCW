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
const CW_SETTINGS_STORAGE_KEY = 'cw.settings.v1';
const CW_SETTINGS_UPDATED_AT_STORAGE_KEY = 'cw.settings.updated_at.v1';
const PAGE_SETTINGS_UPDATED_AT_STORAGE_KEY = 'cw.page_settings.updated_at.v1';

const DEFAULT_CW_SETTINGS: CWSettings = {
  char_wpm: 20,
  eff_wpm: 12,
  freq: 600,
  start_delay: 0.5
};

export function normalizeLesson(lesson: number, maxLesson: number): number {
  if (!Number.isFinite(lesson)) return 1;
  return Math.min(Math.max(1, Math.trunc(lesson)), Math.max(1, maxLesson));
}

function clampNumber(value: number, min: number, max: number): number {
  if (!Number.isFinite(value)) return min;
  return Math.min(max, Math.max(min, value));
}

function nowIso(): string {
  return new Date().toISOString();
}

function parseUpdatedAt(value: string | null | undefined): number {
  if (!value || value.trim() === '') return 0;

  const normalized = value.trim();
  if (normalized === '0001-01-01T00:00:00Z') return 0;

  const parsed = Date.parse(normalized);
  return Number.isNaN(parsed) ? 0 : parsed;
}

function readLocalUpdatedAt(storageKey: string): number {
  if (typeof localStorage === 'undefined') return 0;
  return parseUpdatedAt(localStorage.getItem(storageKey));
}

function writeLocalUpdatedAt(storageKey: string, isoTimestamp: string): void {
  if (typeof localStorage === 'undefined') return;
  localStorage.setItem(storageKey, isoTimestamp);
}

export function touchLocalPageSettingsUpdatedAt(isoTimestamp: string = nowIso()): void {
  writeLocalUpdatedAt(PAGE_SETTINGS_UPDATED_AT_STORAGE_KEY, isoTimestamp);
}

export function touchLocalCwSettingsUpdatedAt(isoTimestamp: string = nowIso()): void {
  writeLocalUpdatedAt(CW_SETTINGS_UPDATED_AT_STORAGE_KEY, isoTimestamp);
}

function readStoredLesson(maxLesson: number): number {
  if (typeof localStorage === 'undefined') return 1;
  const rawLesson = localStorage.getItem('learn.lesson');
  const parsedLesson = Number.parseInt(rawLesson ?? '1', 10);
  return normalizeLesson(parsedLesson, maxLesson);
}

export function normalizeClientCwSettings(raw: Partial<CWSettings> | null | undefined): CWSettings {
  return {
    char_wpm: clampNumber(raw?.char_wpm ?? DEFAULT_CW_SETTINGS.char_wpm, 5, 50),
    eff_wpm: clampNumber(raw?.eff_wpm ?? DEFAULT_CW_SETTINGS.eff_wpm, 5, 50),
    freq: clampNumber(raw?.freq ?? DEFAULT_CW_SETTINGS.freq, 300, 2000),
    start_delay: clampNumber(raw?.start_delay ?? DEFAULT_CW_SETTINGS.start_delay, 0, 10)
  };
}

export function readClientCwSettings(): CWSettings {
  if (typeof localStorage === 'undefined') {
    return { ...DEFAULT_CW_SETTINGS };
  }

  const raw = localStorage.getItem(CW_SETTINGS_STORAGE_KEY);
  if (!raw) {
    return { ...DEFAULT_CW_SETTINGS };
  }

  try {
    const parsed = JSON.parse(raw) as Partial<CWSettings>;
    return normalizeClientCwSettings(parsed);
  } catch {
    return { ...DEFAULT_CW_SETTINGS };
  }
}

export function saveClientCwSettings(settings: CWSettings): CWSettings {
  const normalized = normalizeClientCwSettings(settings);

  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(CW_SETTINGS_STORAGE_KEY, JSON.stringify(normalized));
  }

  touchLocalCwSettingsUpdatedAt();

  return normalized;
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
  options: { applyTheme?: boolean; applyLanguage?: boolean; navigate?: boolean } = {}
): number {
  const lesson = normalizeLesson(page.cur_lesson, maxLesson);
  const language = normalizeLocalePreference(page.language);
  const applyTheme = options.applyTheme ?? true;
  const applyLanguage = options.applyLanguage ?? true;

  if (typeof localStorage !== 'undefined') {
    if (applyTheme) {
      localStorage.setItem('theme', page.theme);
    }
    localStorage.setItem('learn.lesson', String(lesson));
    if (applyLanguage) {
      localStorage.setItem(LOCALE_PREFERENCE_STORAGE_KEY, language);
    }
  }

  if (typeof document !== 'undefined') {
    if (applyTheme) {
      if (page.theme === 'auto') document.documentElement.removeAttribute('data-theme');
      else document.documentElement.setAttribute('data-theme', page.theme);
    }

    document.cookie = `${LESSON_COOKIE}=${lesson}; path=/; max-age=${ONE_YEAR_SECONDS}; SameSite=Lax`;

    if (applyLanguage) {
      document.cookie = `${LOCALE_COOKIE}=${language}; path=/; max-age=${ONE_YEAR_SECONDS}; SameSite=Lax`;
    }
  }

  if (applyLanguage && (language === 'auto' || isLocale(language))) {
    onLocale(language, { navigate: options.navigate });
  }

  touchLocalPageSettingsUpdatedAt();

  return lesson;
}

export async function syncSettingsToServer(cw: CWSettings, page: PageSettings): Promise<void> {
  await Promise.all([saveCWSettings(cw), savePageSettings(page)]);
}

export async function restoreSettingsFromServer(): Promise<{ cw: CWSettings; page: PageSettings }> {
  const settings = await getSettings();
  return { cw: settings.cw_settings, page: settings.page_settings };
}

export async function reconcileSettingsWithServer(args: {
  maxLesson: number;
  fallbackLanguagePreference: LocalePreference;
  onLocale: (preference: LocalePreference, options?: { navigate?: boolean }) => void;
}): Promise<void> {
  if (typeof localStorage === 'undefined') return;

  const settings = await getSettings();

  const serverCw = normalizeClientCwSettings(settings.cw_settings);
  const serverPage: PageSettings = {
    theme: settings.page_settings.theme,
    language: normalizeLocalePreference(settings.page_settings.language),
    cur_lesson: normalizeLesson(settings.page_settings.cur_lesson, args.maxLesson)
  };

  const localCw = readClientCwSettings();
  const localPage = readClientPageSettings(
    readStoredLesson(args.maxLesson),
    args.maxLesson,
    args.fallbackLanguagePreference
  );

  const serverCwUpdatedAt = parseUpdatedAt(settings.cw_settings.updated_at);
  const localCwUpdatedAt = readLocalUpdatedAt(CW_SETTINGS_UPDATED_AT_STORAGE_KEY);
  const serverPageUpdatedAt = parseUpdatedAt(settings.page_settings.updated_at);
  const localPageUpdatedAt = readLocalUpdatedAt(PAGE_SETTINGS_UPDATED_AT_STORAGE_KEY);

  if (localCwUpdatedAt > serverCwUpdatedAt) {
    await saveCWSettings(localCw);
  } else if (serverCwUpdatedAt > localCwUpdatedAt) {
    saveClientCwSettings(serverCw);
    if (settings.cw_settings.updated_at && serverCwUpdatedAt > 0) {
      touchLocalCwSettingsUpdatedAt(settings.cw_settings.updated_at);
    }
  }

  if (localPageUpdatedAt > serverPageUpdatedAt) {
    await savePageSettings(localPage);
  } else if (serverPageUpdatedAt > localPageUpdatedAt) {
    applyClientPageSettings(serverPage, args.maxLesson, args.onLocale, { navigate: false });
    if (settings.page_settings.updated_at && serverPageUpdatedAt > 0) {
      touchLocalPageSettingsUpdatedAt(settings.page_settings.updated_at);
    }
  }
}

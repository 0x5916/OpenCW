/**
 * Client-side cookie helpers.
 *
 * The app persists locale + lesson preferences in cookies so the server can read
 * them on the next request (see `+layout.server.ts` and `morse/learn/+page.server.ts`).
 */

/** Lifetime used for preference cookies: one year. */
export const ONE_YEAR_SECONDS = 60 * 60 * 24 * 365;

/** Read a cookie value by name. Returns `null` on the server or when absent. */
export function readCookie(name: string): string | null {
  if (typeof document === 'undefined') return null;

  const target = `${name}=`;
  const item = document.cookie
    .split(';')
    .map((part) => part.trim())
    .find((part) => part.startsWith(target));

  return item ? decodeURIComponent(item.slice(target.length)) : null;
}

/** Persist a site-wide cookie. No-op on the server. */
export function writeCookie(name: string, value: string, maxAgeSeconds = ONE_YEAR_SECONDS): void {
  if (typeof document === 'undefined') return;

  document.cookie = `${name}=${value}; path=/; max-age=${maxAgeSeconds}; SameSite=Lax`;
}

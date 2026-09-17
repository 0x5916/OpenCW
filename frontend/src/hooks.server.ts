import type { Handle } from '@sveltejs/kit';
import { paraglideMiddleware } from '$lib/paraglide/server';
import {
  deLocalizeHref,
  getLocale,
  getStrategyForUrl,
  locales,
  localizeHref,
  shouldRedirect
} from '$lib/paraglide/runtime';

function isCrawlerUserAgent(userAgent: string | null): boolean {
  if (!userAgent) return false;
  return /bot|crawler|spider|slurp|bingpreview|facebookexternalhit|linkedinbot|duckduckbot/i.test(
    userAgent
  );
}

function appendVary(headerValue: string | null, token: string): string {
  if (!headerValue) return token;
  const values = headerValue
    .split(',')
    .map((value) => value.trim())
    .filter(Boolean);

  if (!values.includes(token)) values.push(token);
  return values.join(', ');
}

function isLikelyPageRequest(request: Request): boolean {
  if (request.method !== 'GET' && request.method !== 'HEAD') return false;

  const url = new URL(request.url);
  const pathname = url.pathname;
  const lastSegment = pathname.split('/').filter(Boolean).at(-1) ?? '';

  // Skip static assets and extension-based files.
  if (lastSegment.includes('.')) return false;

  const accept = request.headers.get('accept') ?? '';
  const fetchDest = request.headers.get('sec-fetch-dest');

  return (
    fetchDest === 'document' ||
    accept.includes('text/html') ||
    accept.includes('application/xhtml+xml') ||
    accept.includes('*/*')
  );
}

/** Locale prefix carried by a path, if it has one. */
function localePrefixFor(pathname: string): string | null {
  return (
    locales.find((locale) => pathname === `/${locale}` || pathname.startsWith(`/${locale}/`)) ??
    null
  );
}

/** Paths that moved; the old URLs keep working with a single permanent redirect. */
const LEGACY_ROUTE_REDIRECTS: Record<string, string> = {
  '/morse': '/learn',
  '/morse/learn': '/learn'
};

/**
 * The trainer moved from `/morse/learn` to `/learn`. Resolving the old paths
 * here rather than in a route keeps it to one hop, because any bare path is
 * already redirected to its localized form before routing happens.
 */
function legacyRedirectFor(request: Request): Response | null {
  if (!isLikelyPageRequest(request)) return null;

  const url = new URL(request.url);
  const targetPath = LEGACY_ROUTE_REDIRECTS[deLocalizeHref(url.pathname)];
  if (!targetPath) return null;

  // Reuse the prefix the request already carried, so the redirect stays in the
  // visitor's language even where locale detection has to fall back.
  const prefix = localePrefixFor(url.pathname);
  const target = new URL(
    prefix ? `${prefix}${targetPath}` : localizeHref(targetPath, { locale: getLocale() }),
    url.origin
  );
  target.search = url.search;

  return new Response(null, { status: 308, headers: { Location: target.href } });
}

export const handle: Handle = async ({ event, resolve }) => {
  const isCrawler = isCrawlerUserAgent(event.request.headers.get('user-agent'));

  if (isLikelyPageRequest(event.request)) {
    const decision = await shouldRedirect({ request: event.request });
    if (decision.shouldRedirect && decision.redirectUrl) {
      const headers = new Headers({ Location: decision.redirectUrl.href });
      if (getStrategyForUrl(event.request.url).includes('preferredLanguage')) {
        headers.set('Vary', 'Accept-Language');
      }
      // The destination can depend on the locale cookie, so a shared cache must
      // not hand this redirect to a visitor with a different preference.
      headers.set('Vary', appendVary(headers.get('Vary'), 'Cookie'));

      return new Response(null, { status: 307, headers });
    }
  }

  return paraglideMiddleware(event.request, async () => {
    const legacyRedirect = legacyRedirectFor(event.request);
    if (legacyRedirect) return legacyRedirect;

    const response = await resolve(event, {
      // Resolve `%paraglide.lang%` in app.html so the served markup carries the
      // request locale instead of a hard-coded `lang="en"`.
      transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', getLocale())
    });

    if (isCrawler && isLikelyPageRequest(event.request)) {
      response.headers.set('Cache-Control', 'no-cache, max-age=0, must-revalidate');
      response.headers.set('Vary', appendVary(response.headers.get('Vary'), 'User-Agent'));
    }

    return response;
  });
};

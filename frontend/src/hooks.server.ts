import type { Handle } from '@sveltejs/kit';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { getStrategyForUrl, shouldRedirect } from '$lib/paraglide/runtime';

const NO_CACHE_HEADER = 'no-cache, no-store, must-revalidate';

function shouldDisableCachingForPwaAsset(pathname: string): boolean {
  return (
    pathname === '/service-worker.js' ||
    pathname === '/sw.js' ||
    pathname === '/site.webmanifest' ||
    pathname.startsWith('/workbox-') ||
    pathname.endsWith('/_app/version.json')
  );
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

export const handle: Handle = async ({ event, resolve }) => {
  if (isLikelyPageRequest(event.request)) {
    const decision = await shouldRedirect({ request: event.request });
    if (decision.shouldRedirect && decision.redirectUrl) {
      const headers = new Headers({ Location: decision.redirectUrl.href });
      if (getStrategyForUrl(event.request.url).includes('preferredLanguage')) {
        headers.set('Vary', 'Accept-Language');
      }

      return new Response(null, { status: 307, headers });
    }
  }

  return paraglideMiddleware(event.request, async () => {
    const response = await resolve(event);

    if (shouldDisableCachingForPwaAsset(event.url.pathname)) {
      response.headers.set('Cache-Control', NO_CACHE_HEADER);
      response.headers.set('Pragma', 'no-cache');
      response.headers.set('Expires', '0');
    }

    return response;
  });
};

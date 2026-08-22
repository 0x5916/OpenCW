/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
/// <reference types="@sveltejs/kit" />

import { build, files, prerendered, version } from '$service-worker';

// Gives `self` the correct ServiceWorkerGlobalScope types.
const self = globalThis.self as unknown as ServiceWorkerGlobalScope;

const CACHE = `opencw-${version}`;
const ASSETS = [...build, ...files];
const OFFLINE = '/offline.html';

function normalizePathname(pathname: string): string {
  return pathname.length > 1 && pathname.endsWith('/') ? pathname.slice(0, -1) : pathname;
}

function isCacheable(response: Response): boolean {
  if (response.status !== 200 || response.type !== 'basic') return false;
  const cacheControl = response.headers.get('cache-control') ?? '';
  return !/(no-store|no-cache|private)/i.test(cacheControl);
}

self.addEventListener('install', (event) => {
  event.waitUntil(
    (async () => {
      const cache = await caches.open(CACHE);
      // Precache the app shell, static files, and every prerendered page so the
      // whole app works offline. A single failed asset must not abort the whole
      // install, so failures are swallowed per-file.
      await Promise.all(
        [...ASSETS, ...prerendered].map((asset) => cache.add(asset).catch(() => undefined))
      );
      await self.skipWaiting();
    })()
  );
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    (async () => {
      const keys = await caches.keys();
      await Promise.all(keys.filter((key) => key !== CACHE).map((key) => caches.delete(key)));
      await self.clients.claim();
    })()
  );
});

self.addEventListener('fetch', (event) => {
  const request = event.request;
  if (request.method !== 'GET') return;

  const url = new URL(request.url);
  // Only handle same-origin requests. The API (`PUBLIC_API_BASE`) and any other
  // cross-origin traffic pass straight through to the network.
  if (url.origin !== self.location.origin) return;

  // App shell (JS/CSS) and static assets: cache-first.
  if (ASSETS.includes(url.pathname)) {
    event.respondWith(
      (async () => {
        const cache = await caches.open(CACHE);
        const cached = await cache.match(request);
        if (cached) return cached;
        try {
          const response = await fetch(request);
          if (isCacheable(response)) await cache.put(request, response.clone());
          return response;
        } catch {
          return new Response('', { status: 504, statusText: 'Offline' });
        }
      })()
    );
    return;
  }

  // Navigations: network-first, falling back to the cached (prerendered) page,
  // then the standalone offline page.
  if (request.mode === 'navigate') {
    event.respondWith(
      (async () => {
        const cache = await caches.open(CACHE);
        try {
          const response = await fetch(request);
          if (isCacheable(response)) await cache.put(request, response.clone());
          return response;
        } catch {
          const pathname = normalizePathname(url.pathname);
          const cached =
            (await cache.match(request)) ??
            (await cache.match(pathname)) ??
            // The bare root is the PWA start_url; only localized paths get
            // cached, so fall back to the base-locale home page if available.
            (pathname === '/' ? await cache.match('/en') : undefined) ??
            (await cache.match(OFFLINE));
          if (cached) return cached;
          return new Response('Offline', {
            status: 504,
            statusText: 'Offline',
            headers: { 'Content-Type': 'text/plain; charset=utf-8' }
          });
        }
      })()
    );
    return;
  }

  // Any other same-origin GET: network-first with cache fallback.
  event.respondWith(
    (async () => {
      const cache = await caches.open(CACHE);
      try {
        const response = await fetch(request);
        if (isCacheable(response)) await cache.put(request, response.clone());
        return response;
      } catch (error) {
        const cached = await cache.match(request);
        if (cached) return cached;
        throw error;
      }
    })()
  );
});

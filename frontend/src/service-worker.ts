/// <reference no-default-lib="true" />
/// <reference lib="esnext" />
/// <reference lib="webworker" />
/// <reference types="@sveltejs/kit" />

import { build, files, prerendered, version } from '$service-worker';

const self = globalThis.self as unknown as ServiceWorkerGlobalScope;

const CACHE = `opencw-${version}`;
const OFFLINE_URL = '/offline.html';
const API_PREFIXES = ['/auth', '/user', '/settings', '/cw', '/v1'];
const NAVIGATION_FALLBACKS = ['/en/morse/learn', '/morse/learn', '/en', '/'];

const ASSETS = new Set([...build, ...files, ...prerendered]);
ASSETS.add(OFFLINE_URL);

self.addEventListener('install', (event) => {
  event.waitUntil(
    (async () => {
      const cache = await caches.open(CACHE);
      await cache.addAll([...ASSETS]);
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

function shouldBypassRuntimeCache(pathname: string): boolean {
  return API_PREFIXES.some((prefix) => pathname.startsWith(prefix));
}

async function findCachedNavigationFallback(cache: Cache): Promise<Response | undefined> {
  for (const path of NAVIGATION_FALLBACKS) {
    const candidate = await cache.match(path);
    if (candidate) return candidate;
  }

  return undefined;
}

self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return;

  const url = new URL(event.request.url);
  if (url.origin !== self.location.origin) return;

  const pathname = url.pathname;

  event.respondWith(
    (async () => {
      const cache = await caches.open(CACHE);

      if (ASSETS.has(pathname)) {
        const cachedAsset = await cache.match(pathname);
        if (cachedAsset) return cachedAsset;
      }

      if (event.request.mode === 'navigate') {
        try {
          const response = await fetch(event.request);
          if (response.ok) {
            await cache.put(event.request, response.clone());
          }
          return response;
        } catch {
          const cachedPage = await cache.match(event.request);
          if (cachedPage) return cachedPage;

          const cachedFallbackPage = await findCachedNavigationFallback(cache);
          if (cachedFallbackPage) return cachedFallbackPage;

          const offline = await cache.match(OFFLINE_URL);
          if (offline) return offline;

          return new Response('Offline', { status: 503, statusText: 'Offline' });
        }
      }

      if (shouldBypassRuntimeCache(pathname)) {
        return fetch(event.request);
      }

      try {
        const response = await fetch(event.request);
        if (response.ok && response.type === 'basic') {
          await cache.put(event.request, response.clone());
        }
        return response;
      } catch {
        const cached = await cache.match(event.request);
        if (cached) return cached;

        return new Response('Offline', { status: 503, statusText: 'Offline' });
      }
    })()
  );
});

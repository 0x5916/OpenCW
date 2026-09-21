import adapter from '@sveltejs/adapter-static';
import { readFileSync } from 'node:fs';

// Locales are read from the committed inlang project (the generated paraglide
// runtime is gitignored and does not exist before `vite build` starts).
const { locales } = JSON.parse(
  readFileSync(new URL('./project.inlang/settings.json', import.meta.url), 'utf8')
);

/** Every static route, without a locale prefix. */
const ROUTE_PATHS = [
  '/',
  '/about',
  '/forum',
  '/morse',
  '/morse/learn',
  '/login',
  '/more',
  '/profile',
  '/register',
  '/settings'
];

/** Mirrors `buildLocalizedPath` in src/lib/seo.ts. */
function localizedPath(routePath, locale) {
  return routePath === '/' ? `/${locale}` : `/${locale}${routePath}`;
}

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      // Unknown URLs fall back to a SvelteKit-rendered 404 page. The adapter's
      // strict check is disabled with a fallback, so the build output is
      // verified explicitly (see scripts/validate-seo.ts and the build docs).
      fallback: '404.html'
    }),
    prerender: {
      // The language switcher is JavaScript-driven, so the crawler cannot
      // discover the other locale variants from `<a>` links — list them.
      entries: [
        '*',
        ...locales.flatMap((locale) =>
          ROUTE_PATHS.map((routePath) => localizedPath(routePath, locale))
        )
      ],
      // Absolute URLs (canonical, hreflang, Open Graph, sitemap.xml) are baked
      // into the prerendered files, so the origin is a build-time value.
      origin: process.env.PRERENDER_ORIGIN ?? 'https://opencw.net'
    }
  }
};

export default config;

import { readFileSync } from 'node:fs';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';
import { paraglideVitePlugin } from '@inlang/paraglide-js';

// The committed inlang project owns the locale list: deriving the URL patterns
// from it keeps a new locale from silently missing its own routes.
const { locales } = JSON.parse(
  readFileSync(new URL('./project.inlang/settings.json', import.meta.url), 'utf8')
) as { locales: string[] };

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit(),
    paraglideVitePlugin({
      project: './project.inlang',
      outdir: './src/lib/paraglide',
      strategy: ['url', 'preferredLanguage', 'globalVariable', 'baseLocale'],
      urlPatterns: [
        {
          pattern: ':protocol://:domain(.*)::port?/:path(.*)?',
          localized: locales.map((locale): [string, string] => [
            locale,
            `:protocol://:domain(.*)::port?/${locale}/:path(.*)?`
          ])
        }
      ]
    })
  ]
});

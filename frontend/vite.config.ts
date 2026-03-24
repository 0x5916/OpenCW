import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';
import { paraglideVitePlugin } from '@inlang/paraglide-js';

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit(),
    paraglideVitePlugin({
      project: './project.inlang',
      outdir: './src/lib/paraglide',
      strategy: ['url', 'cookie', 'preferredLanguage', 'globalVariable', 'baseLocale'],
      urlPatterns: [
        {
          pattern: ':protocol://:domain(.*)::port?/:path(.*)?',
          localized: [
            ['en', ':protocol://:domain(.*)::port?/en/:path(.*)?'],
            ['zh-Hant', ':protocol://:domain(.*)::port?/zh-Hant/:path(.*)?'],
            ['zh-Hans', ':protocol://:domain(.*)::port?/zh-Hans/:path(.*)?'],
            ['ja', ':protocol://:domain(.*)::port?/ja/:path(.*)?'],
            ['de', ':protocol://:domain(.*)::port?/de/:path(.*)?']
          ]
        }
      ]
    })
  ]
});

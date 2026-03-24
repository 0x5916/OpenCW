import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';
import { paraglideVitePlugin } from '@inlang/paraglide-js';
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
		VitePWA({
			injectRegister: false,
			registerType: 'autoUpdate',
			includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'favicon-96x96.png'],
			manifest: {
				name: 'OpenCW',
				short_name: 'OpenCW',
				description:
					'Practice Morse code (CW) with OpenCW using the Koch method. Train speed, accuracy, and consistency online or offline.',
				theme_color: '#030712',
				background_color: '#030712',
				display: 'standalone',
				start_url: '/',
				scope: '/',
				icons: [
					{
						src: '/web-app-manifest-192x192.png',
						sizes: '192x192',
						type: 'image/png'
					},
					{
						src: '/web-app-manifest-512x512.png',
						sizes: '512x512',
						type: 'image/png'
					},
					{
						src: '/web-app-manifest-512x512.png',
						sizes: '512x512',
						type: 'image/png',
						purpose: 'maskable'
					}
				]
			},
			workbox: {
				navigateFallback: '/offline.html',
				navigateFallbackDenylist: [/^\/v1\//],
				globPatterns: ['**/*.{js,css,html,ico,png,svg,webmanifest,json}'],
				additionalManifestEntries: [
					{ url: '/', revision: null },
					{ url: '/en', revision: null },
					{ url: '/zh-Hant', revision: null },
					{ url: '/zh-Hans', revision: null },
					{ url: '/ja', revision: null },
					{ url: '/de', revision: null },
					{ url: '/en/morse/learn', revision: null },
					{ url: '/zh-Hant/morse/learn', revision: null },
					{ url: '/zh-Hans/morse/learn', revision: null },
					{ url: '/ja/morse/learn', revision: null },
					{ url: '/de/morse/learn', revision: null },
					{ url: '/offline.html', revision: null }
				],
				runtimeCaching: [
					{
						urlPattern: ({ url }) => url.pathname.startsWith('/v1/'),
						handler: 'NetworkOnly'
					},
					{
						urlPattern: ({ request }) => request.mode === 'navigate',
						handler: 'NetworkFirst',
						options: {
							cacheName: 'pages-cache',
							networkTimeoutSeconds: 3,
							expiration: {
								maxEntries: 40,
								maxAgeSeconds: 60 * 60 * 24 * 7
							}
						}
					},
					{
						urlPattern: ({ request }) =>
							request.destination === 'style' ||
							request.destination === 'script' ||
							request.destination === 'worker',
						handler: 'StaleWhileRevalidate',
						options: {
							cacheName: 'assets-cache',
							expiration: {
								maxEntries: 80,
								maxAgeSeconds: 60 * 60 * 24 * 14
							}
						}
					},
					{
						urlPattern: ({ request }) => request.destination === 'image',
						handler: 'CacheFirst',
						options: {
							cacheName: 'images-cache',
							expiration: {
								maxEntries: 120,
								maxAgeSeconds: 60 * 60 * 24 * 30
							}
						}
					}
				]
			}
		}),
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

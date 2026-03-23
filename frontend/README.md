# sv

Everything you need to build a Svelte project, powered by [`sv`](https://github.com/sveltejs/cli).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```sh
# create a new project
npx sv create my-app
```

To recreate this project with the same configuration:

```sh
# recreate this project
npx sv create --template minimal --types ts --add prettier eslint --install npm cw-frontend
```

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```sh
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```sh
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.

## PWA and offline access

This project includes a SvelteKit service worker at `src/service-worker.ts` and an offline fallback page at `static/offline.html`.

Behavior summary:

- Pre-caches build output and static assets.
- Serves cached assets while offline.
- Uses network-first navigation with fallback to the cached page or `offline.html`.
- Skips runtime cache for dynamic API paths (`/auth`, `/user`, `/settings`, `/cw`, `/v1`).

The web app manifest is `static/site.webmanifest`.

### Cloudflare cache settings for opencw.net and dev.opencw.net

To avoid stale service worker updates, disable edge/browser caching for these paths:

- `/service-worker.js`
- `/sw.js`
- `/workbox-*`
- `/_app/version.json`
- `/site.webmanifest`

Recommended cache policy for the paths above:

- `Cache-Control: no-cache, no-store, must-revalidate`

All other fingerprinted static assets can remain long-lived and immutable.

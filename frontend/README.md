# OpenCW Frontend

Web frontend for OpenCW — a Morse code (CW) training site with Koch-method lessons,
progress tracking, and a community forum.

Built with SvelteKit (Svelte 5 runes), TypeScript, Paraglide i18n (en, de, ja, zh-Hans,
zh-Hant) and Tailwind CSS (preflight; most UI is plain CSS/utility classes in `src/app.css`).
The site is fully static: every page is prerendered for all five locales at build time
and served as plain files — there is no server runtime.

## Requirements

- Node.js 24 (see the `Dockerfile`; `npm` 11+ recommended)
- A running OpenCW API (defaults to `http://localhost:8080/v1` via `PUBLIC_API_BASE`)

## Setup

```sh
npm ci
cp example.env .env   # adjust PUBLIC_API_BASE if needed
npm run dev
```

## Scripts

| Script                 | Purpose                                                        |
| ---------------------- | -------------------------------------------------------------- |
| `npm run dev`          | Vite dev server with HMR                                       |
| `npm run build`        | Static production build (prerendered site in `build/`)         |
| `npm run start`        | Serve the production build locally (`vite preview`)            |
| `npm run preview`      | Preview the production build                                   |
| `npm run check`        | `svelte-kit sync` + `svelte-check` (TypeScript/Svelte errors)  |
| `npm run lint`         | Prettier check + ESLint                                        |
| `npm run format`       | Prettier write                                                 |
| `npm run seo:validate` | SEO gate: sitemap coverage, metadata lengths, duplicate titles |
| `npm run icons`        | Regenerate favicons/social images from the source SVG          |

## Project layout

```
messages/          Paraglide translation catalogs (one JSON per locale)
scripts/           Icon generation + SEO validation
src/
  app.css          Design tokens, base styles, shared UI classes
  hooks.server.ts  Paraglide middleware (dev server + prerendering only)
  lib/
    api.ts         Typed API client (settings, user, forum, progress)
    auth.ts        Token handling (login/register/refresh/logout) + apiFetch
    cookies.ts     One-time migration from legacy preference cookies
    cwSync.ts      Client/server CW + page settings reconciliation
    errorCode.ts   API error-code extraction
    errorLocalization.ts  Error code → localized message
    format.ts      Date / lesson / percentage formatting helpers
    i18n.svelte.ts Locale state + locale-aware hrefs
    locale.ts      Locale matching + display labels
    morse.ts       Koch lessons, Morse table, Farnsworth timing
    progressSync.ts Offline-first progress queue
    score.ts       Accuracy scoring, word-level diff, grade thresholds
    seo.ts         Route metadata, sitemap URL builder (used by scripts)
    storageKeys.ts localStorage key registry
    theme.ts       Theme normalization + apply helpers
    components/    Shared UI (auth card, dropdown, Morse player, alerts, …)
  routes/          Pages (home, about, learn, login, register, profile,
                   settings, a placeholder /forum, and legacy /morse redirects)
```

`/forum` currently renders an "under development" card. The forum parts of
`api.ts` and `components/Pagination.svelte` are intentionally kept (and still
covered by localized `api_error_forum_*` messages) for when the forum returns.

## Internationalization

- Message keys live in `messages/*.json`; all locales must define the same keys.
- Paraglide output is generated into `src/lib/paraglide/` (gitignored) by the Vite
  plugin. After changing keys in `messages/*.json`, run `npm run build` (or
  `npm run check`) before `svelte-check` can resolve new `m.*` imports.
- URLs are localized (`/de/...`, `/ja/...`, …). Every locale variant is
  prerendered, so the active locale always comes from the URL.
- Bare paths (`/about`) and legacy `/morse` URLs exist only in the base locale in
  the build; once the app hydrates they redirect to the visitor's preferred
  locale (stored in `localStorage`). Existing visitors' legacy preference
  cookies are migrated once and then removed.

## SEO validation gate

Run before release:

```sh
npm run seo:validate
```

It checks sitemap URL coverage for all locales and indexable public routes, excludes
noindex routes, and sanity-checks title/description lengths and duplicate titles per
locale.

## Docker

```sh
docker build \
  --build-arg PUBLIC_API_BASE=https://api.example.com/v1 \
  --build-arg PRERENDER_ORIGIN=https://opencw.net \
  -t opencw-frontend .
docker run -p 3000:80 opencw-frontend
```

The image builds the static site and serves it with nginx. Both arguments are
build-time values: `PUBLIC_API_BASE` is baked into the client bundle through
`$env/static/public`, and `PRERENDER_ORIGIN` (default `https://opencw.net`) into
the canonical/hreflang/Open Graph URLs and `sitemap.xml`.

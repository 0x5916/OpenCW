# OpenCW Frontend

Web frontend for OpenCW — a Morse code (CW) training site with Koch-method lessons,
progress tracking, and a community forum.

Built with SvelteKit (Svelte 5 runes), TypeScript, Paraglide i18n (en, de, ja, zh-Hans,
zh-Hant) and Tailwind CSS (preflight; most UI is plain CSS/utility classes in `src/app.css`).

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
| `npm run build`        | Production build (adapter-node output in `build/`)             |
| `npm run start`        | Run the production server (`node build`)                       |
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
  hooks.server.ts  Paraglide middleware + crawler/redirect handling
  lib/
    api.ts         Typed API client (settings, user, forum, progress)
    auth.ts        Token handling (login/register/refresh/logout) + apiFetch
    cwSync.ts      Client/server CW + page settings reconciliation
    progressSync.ts Offline-first progress queue
    morse.ts       Koch lessons, Morse table, Farnsworth timing
    score.ts       Accuracy scoring + word-level diff
    seo.ts         Route metadata, sitemap URL builder (used by scripts)
  routes/          Pages (home, about, forum, morse/learn, login, register, …)
```

## Internationalization

- Message keys live in `messages/*.json`; all locales must define the same keys.
- Paraglide output is generated into `src/lib/paraglide/` (gitignored) by the Vite
  plugin. After changing keys in `messages/*.json`, run `npm run build` (or
  `npm run check`) before `svelte-check` can resolve new `m.*` imports.
- URLs are localized (`/de/...`, `/ja/...`, …); locale resolution uses the
  `PARAGLIDE_LOCALE` cookie, then `preferredLanguage`, then the base locale.

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
docker build --build-arg PUBLIC_API_BASE=https://api.example.com/v1 -t opencw-frontend .
docker run -p 3000:3000 opencw-frontend
```

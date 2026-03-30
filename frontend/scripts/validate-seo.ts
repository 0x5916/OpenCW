import { locales } from '../src/lib/paraglide/runtime.js';
import {
  buildSitemapUrlSet,
  buildLocalizedPath,
  getIndexablePublicRoutePaths,
  isRouteIndexable,
  resolveSeoMetadata
} from '../src/lib/seo.ts';

type CheckResult = { ok: true; message: string } | { ok: false; message: string };

const INDEXABLE_ROUTES = getIndexablePublicRoutePaths();
const NOINDEX_ROUTES = ['/login', '/register', '/profile', '/settings', '/offline'];

function absolute(origin: string, path: string): string {
  return new URL(path, origin).toString();
}

function expectedLocalizedPath(locale: string, routePath: string): string {
  return buildLocalizedPath(routePath, locale);
}

function validateMetadata(): CheckResult[] {
  const checks: CheckResult[] = [];

  for (const locale of locales) {
    const seenTitles = new Map<string, string>();

    for (const routePath of INDEXABLE_ROUTES) {
      const metadata = resolveSeoMetadata(routePath, locale);
      const titleLen = metadata.title.trim().length;
      const descriptionLen = metadata.description.trim().length;

      if (!metadata.title.trim()) {
        checks.push({ ok: false, message: `[${locale}] ${routePath}: missing title` });
      }

      if (!metadata.description.trim()) {
        checks.push({ ok: false, message: `[${locale}] ${routePath}: missing description` });
      }

      if (titleLen < 20 || titleLen > 70) {
        checks.push({
          ok: false,
          message: `[${locale}] ${routePath}: title length ${titleLen} outside recommended range (20-70)`
        });
      }

      if (descriptionLen < 80 || descriptionLen > 180) {
        checks.push({
          ok: false,
          message: `[${locale}] ${routePath}: description length ${descriptionLen} outside recommended range (80-180)`
        });
      }

      const duplicateRoute = seenTitles.get(metadata.title);
      if (duplicateRoute) {
        checks.push({
          ok: false,
          message: `[${locale}] duplicate title shared by ${duplicateRoute} and ${routePath}`
        });
      } else {
        seenTitles.set(metadata.title, routePath);
      }
    }
  }

  return checks;
}

async function validateSitemap(): Promise<CheckResult[]> {
  const checks: CheckResult[] = [];
  const origin = 'https://opencw.net';
  const generatedUrls = buildSitemapUrlSet(origin, locales);
  const urlSet = new Set(generatedUrls);

  if (generatedUrls.length !== urlSet.size) {
    checks.push({ ok: false, message: 'generated sitemap URL set contains duplicates' });
  }

  const expectedUrls = new Set<string>([absolute(origin, '/')]);
  for (const locale of locales) {
    for (const routePath of INDEXABLE_ROUTES) {
      expectedUrls.add(absolute(origin, expectedLocalizedPath(locale, routePath)));
    }
  }

  for (const expected of expectedUrls) {
    if (!urlSet.has(expected)) {
      checks.push({ ok: false, message: `sitemap missing expected URL: ${expected}` });
    }
  }

  for (const routePath of NOINDEX_ROUTES) {
    if (isRouteIndexable(routePath)) {
      checks.push({ ok: false, message: `route marked noindex expected but currently indexable: ${routePath}` });
    }

    for (const locale of locales) {
      const localized = absolute(origin, expectedLocalizedPath(locale, routePath));
      if (urlSet.has(localized)) {
        checks.push({ ok: false, message: `sitemap should not include noindex URL: ${localized}` });
      }
    }
  }

  if (checks.length === 0) {
    checks.push({ ok: true, message: `sitemap URL set validated (${urlSet.size} URLs)` });
  }

  return checks;
}

async function main() {
  const results: CheckResult[] = [];

  results.push(...validateMetadata());
  results.push(...(await validateSitemap()));

  const failures = results.filter((result) => !result.ok);
  const successes = results.filter((result) => result.ok);

  for (const success of successes) {
    console.log(`OK: ${success.message}`);
  }

  for (const failure of failures) {
    console.error(`FAIL: ${failure.message}`);
  }

  if (failures.length > 0) {
    process.exitCode = 1;
    return;
  }

  console.log('SEO validation passed.');
}

await main();

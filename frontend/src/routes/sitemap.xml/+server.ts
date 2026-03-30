import type { RequestHandler } from './$types';
import { locales } from '$lib/paraglide/runtime';
import { buildAbsoluteUrl, normalizePathname, PUBLIC_ROUTE_PATHS } from '$lib/seo';

function escapeXml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&apos;');
}

function buildSitemapUrls(origin: string): string[] {
  const urls = new Set<string>();

  urls.add(buildAbsoluteUrl(origin, '/'));

  for (const locale of locales) {
    for (const routePath of PUBLIC_ROUTE_PATHS) {
      const localizedPath =
        routePath === '/' ? normalizePathname(`/${locale}`) : normalizePathname(`/${locale}${routePath}`);
      urls.add(buildAbsoluteUrl(origin, localizedPath));
    }
  }

  return [...urls].sort();
}

function buildXml(urls: string[]): string {
  const now = new Date().toISOString().slice(0, 10);
  const body = urls
    .map((url) => {
      const priority = url.endsWith('/en') || url.endsWith('/de') || url.endsWith('/ja') || url.endsWith('/zh-Hans') || url.endsWith('/zh-Hant') || url.endsWith('/')
        ? '1.0'
        : url.includes('/morse/learn')
          ? '0.9'
          : '0.8';
      const changefreq = url.includes('/forum') ? 'daily' : url.includes('/morse/learn') ? 'monthly' : 'weekly';

      return `  <url>\n    <loc>${escapeXml(url)}</loc>\n    <lastmod>${now}</lastmod>\n    <changefreq>${changefreq}</changefreq>\n    <priority>${priority}</priority>\n  </url>`;
    })
    .join('\n');

  return `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${body}\n</urlset>\n`;
}

export const GET: RequestHandler = ({ url }) => {
  const origin = `${url.protocol}//${url.host}`;
  const xml = buildXml(buildSitemapUrls(origin));

  return new Response(xml, {
    headers: {
      'Content-Type': 'application/xml; charset=utf-8',
      'Cache-Control': 'public, max-age=3600, stale-while-revalidate=86400'
    }
  });
};

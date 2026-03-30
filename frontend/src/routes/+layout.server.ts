import type { LayoutServerLoad } from './$types';
import { LOCALE_COOKIE, normalizeLocalePreference, type LocalePreference } from '$lib/locale';
import { deLocalizeHref, getLocale, locales, localizeHref } from '$lib/paraglide/runtime';
import {
  buildAbsoluteUrl,
  DEFAULT_OG_IMAGE_PATH,
  getOpenGraphLocale,
  normalizePathname,
  resolveSeoMetadata,
  SITE_NAME
} from '$lib/seo';

export const load: LayoutServerLoad = ({ cookies, url, route }) => {
  const storedPreference = cookies.get(LOCALE_COOKIE);
  const localePreference = normalizeLocalePreference(storedPreference);
  const locale = getLocale();
  const normalizedPath = normalizePathname(url.pathname);
  const basePath = normalizePathname(deLocalizeHref(normalizedPath));
  const metadata = resolveSeoMetadata(route.id, locale);
  const canonicalUrl = buildAbsoluteUrl(url.origin, normalizedPath);

  const alternates = locales.map((alternateLocale) => {
    const localizedPath = localizeHref(basePath, { locale: alternateLocale });

    return {
      locale: alternateLocale,
      href: buildAbsoluteUrl(url.origin, localizedPath)
    };
  });

  const xDefaultHref = alternates.find((alternate) => alternate.locale === 'en')?.href ?? canonicalUrl;
  const isIndexable = !metadata.robots.toLowerCase().includes('noindex');
  const structuredData = isIndexable
    ? [
        {
          '@context': 'https://schema.org',
          '@type': 'Organization',
          name: SITE_NAME,
          url: buildAbsoluteUrl(url.origin, '/'),
          logo: buildAbsoluteUrl(url.origin, '/apple-touch-icon.png'),
          sameAs: ['https://github.com/0x5916']
        },
        {
          '@context': 'https://schema.org',
          '@type': 'WebSite',
          name: SITE_NAME,
          url: buildAbsoluteUrl(url.origin, '/'),
          inLanguage: locale,
          description: metadata.description
        },
        ...(route.id === '/morse/learn'
          ? [
              {
                '@context': 'https://schema.org',
                '@type': 'Course',
                name: metadata.title,
                description: metadata.description,
                inLanguage: locale,
                provider: {
                  '@type': 'Organization',
                  name: SITE_NAME,
                  url: buildAbsoluteUrl(url.origin, '/')
                },
                educationalLevel: 'Beginner to Intermediate',
                url: canonicalUrl
              }
            ]
          : [])
      ]
    : [];

  return {
    locale,
    localePreference: localePreference as LocalePreference,
    seo: {
      ...metadata,
      siteName: SITE_NAME,
      canonicalUrl,
      alternates,
      xDefaultHref,
      openGraphLocale: getOpenGraphLocale(locale),
      openGraphImage: buildAbsoluteUrl(url.origin, metadata.ogImagePath ?? DEFAULT_OG_IMAGE_PATH),
      structuredData
    }
  };
};

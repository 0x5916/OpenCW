import type { Locale } from '$lib/locale';

export type SeoMetadata = {
  title: string;
  description: string;
  robots: string;
  ogType: 'website' | 'article';
  ogImagePath: string;
};

export const SITE_NAME = 'OpenCW';
export const DEFAULT_OG_IMAGE_PATH = '/web-app-manifest-512x512.png';
export const PUBLIC_ROUTE_PATHS = ['/', '/about', '/forum', '/morse/learn'] as const;

type LocalizedSeoText = {
  title: string;
  description: string;
};

type SeoRouteOverride = {
  robots?: string;
  ogType?: 'website' | 'article';
  ogImagePath?: string;
  localized: Partial<Record<Locale, LocalizedSeoText>>;
};

const DEFAULT_ROBOTS = 'index,follow,max-snippet:-1,max-image-preview:large,max-video-preview:-1';

const DEFAULT_LOCALIZED_TEXT: Record<Locale, LocalizedSeoText> = {
  en: {
    title: 'OpenCW - Morse Code (CW) Training & Practice',
    description:
      "Practice Morse code (CW) with OpenCW's free Koch method trainer. Track WPM progress and join the amateur radio community."
  },
  de: {
    title: 'OpenCW - Morsecode (CW) Training und Uebung',
    description:
      'Lerne Morsecode (CW) mit dem kostenlosen Koch-Trainer von OpenCW. Verfolge deinen WPM-Fortschritt und tausche dich mit der Funk-Community aus.'
  },
  ja: {
    title: 'OpenCW - Morse (CW) no Renshuu',
    description:
      'OpenCW no muryo Koch methodo toreena de Morse (CW) o renshuu. WPM no seichou o kiroku shi, amateur radio community ni sanka dekimasu.'
  },
  'zh-Hans': {
    title: 'OpenCW - Moersi (CW) Xunlian yu Lianxi',
    description:
      'Shi yong OpenCW mianfei Koch xunlianqi lianxi Moersi ma(CW), genzong WPM jinbu, bing jiaru wuxian dian shequ.'
  },
  'zh-Hant': {
    title: 'OpenCW - Moshi (CW) Xunlian yu Lianxi',
    description:
      'Shi yong OpenCW mianfei Koch xunlianqi lianxi Moshi ma(CW), zhuizong WPM jinbu, bing jiaru yuyu diantai shequ.'
  }
};

const ROUTE_SEO: Record<string, SeoRouteOverride> = {
  '/': {
    localized: {
      en: {
        title: 'OpenCW - Morse Code (CW) Training & Practice',
        description:
          "Practice Morse code (CW) with OpenCW's free Koch method trainer. Track WPM progress and join the amateur radio community."
      }
    }
  },
  '/about': {
    localized: {
      en: {
        title: 'About OpenCW - Koch Method Morse Code Training',
        description:
          'Learn how OpenCW uses the Koch method to teach Morse code effectively for amateur radio and CW operators.'
      },
      de: {
        title: 'Ueber OpenCW - Koch-Methode fuer Morsecode',
        description:
          'Erfahre, wie OpenCW mit der Koch-Methode Morsecode fuer Funkamateure und CW-Anwender effektiv vermittelt.'
      },
      ja: {
        title: 'OpenCW ni tsuite - Koch methodo de Manabu Morse',
        description:
          'OpenCW ga Koch methodo de douyatte koukateki ni Morse o oshieru ka o setsumei shimasu.'
      },
      'zh-Hans': {
        title: 'Guanyu OpenCW - Koch Fangfa Mosi Xunlian',
        description: 'Liaojie OpenCW ruhe tongguo Koch fangfa gaoxiao xunlian Mosi ma, mianxiang yuhuo aihaozhe.'
      },
      'zh-Hant': {
        title: 'Guanyu OpenCW - Koch Fangfa Mosi Xunlian',
        description: 'Liaojie OpenCW ruhe tongguo Koch fangfa gao xiao xunlian Mosi ma, mianxiang yuyu dian aihaozhe.'
      }
    }
  },
  '/morse/learn': {
    ogImagePath: '/web-app-manifest-512x512.png',
    localized: {
      en: {
        title: 'Learn Morse Code - OpenCW Koch Trainer',
        description:
          'Train Morse code with adaptive Koch lessons, listening drills, and progress tracking built for practical CW improvement.'
      },
      de: {
        title: 'Morse lernen - OpenCW Koch-Trainer',
        description:
          'Trainiere Morsecode mit adaptiven Koch-Lektionen, Hoeruebungen und Fortschrittsverfolgung fuer echte CW-Verbesserung.'
      },
      ja: {
        title: 'Morse o Manabu - OpenCW Koch Toreena',
        description:
          'Koch lesson, listening drill, shinchoku tsuiseki de jissen-teki ni Morse (CW) no nouryoku o nobasu kunren ga dekimasu.'
      },
      'zh-Hans': {
        title: 'Xuexi Mosi Ma - OpenCW Koch Xunlianqi',
        description: 'Tongguo zishiying Koch kecheng, tingli lianxi he jindu zhuizong, wending tisheng CW shizhan nengli.'
      },
      'zh-Hant': {
        title: 'Xuexi Mosi Ma - OpenCW Koch Xunlianqi',
        description: 'Tongguo zishiying Koch kecheng, tingli lianxi he jindu zhuizong, wending tisheng CW shizhan nengli.'
      }
    }
  },
  '/forum': {
    localized: {
      en: {
        title: 'OpenCW Forum - Morse Code Community',
        description:
          'Join OpenCW community discussions for Morse code learning, amateur radio tips, and CW training support.'
      },
      de: {
        title: 'OpenCW Forum - Morsecode Community',
        description: 'Tausche dich im OpenCW Forum ueber Morse-Lernen, Amateurfunk-Tipps und CW-Training aus.'
      },
      ja: {
        title: 'OpenCW Forum - Morse Community',
        description: 'Morse gakushu, amateur radio no chie, CW no kunren ni tsuite komyuniti de jiyuu ni hanashimashou.'
      },
      'zh-Hans': {
        title: 'OpenCW Luntan - Mosi Ma Shequ',
        description: 'Canyu OpenCW shequ taolun, jiaoliu Mosi xuexi, yuhuo jiqiao he CW xunlian jingyan.'
      },
      'zh-Hant': {
        title: 'OpenCW Luntan - Mosi Ma Shequ',
        description: 'Canyu OpenCW shequ taolun, jiaoliu Mosi xuexi, yuyu jiqiao he CW xunlian jingyan.'
      }
    }
  },
  '/login': {
    robots: 'noindex,nofollow',
    localized: {
      en: {
        title: 'Login - OpenCW',
        description: 'Sign in to OpenCW to continue your Morse code training and sync your learning progress.'
      }
    }
  },
  '/register': {
    robots: 'noindex,nofollow',
    localized: {
      en: {
        title: 'Register - OpenCW',
        description: 'Create an OpenCW account to save your Morse code training progress and settings.'
      }
    }
  },
  '/profile': {
    robots: 'noindex,nofollow',
    localized: {
      en: {
        title: 'Profile - OpenCW',
        description: 'Review your OpenCW profile and training milestones.'
      }
    }
  },
  '/settings': {
    robots: 'noindex,nofollow',
    localized: {
      en: {
        title: 'Settings - OpenCW',
        description: 'Manage your OpenCW preferences for language, practice behavior, and account options.'
      }
    }
  },
  '/offline': {
    robots: 'noindex,nofollow',
    localized: {
      en: {
        title: 'Offline - OpenCW',
        description: 'You are currently offline. Reconnect to continue synchronized Morse code training.'
      }
    }
  }
};

const OG_LOCALE_BY_LOCALE: Record<Locale, string> = {
  en: 'en_US',
  de: 'de_DE',
  ja: 'ja_JP',
  'zh-Hans': 'zh_CN',
  'zh-Hant': 'zh_TW'
};

export function normalizePathname(pathname: string): string {
  if (!pathname) return '/';
  if (pathname === '/') return '/';
  return pathname.endsWith('/') ? pathname.slice(0, -1) : pathname;
}

export function buildAbsoluteUrl(origin: string, pathname: string): string {
  const normalizedPath = normalizePathname(pathname);
  return new URL(normalizedPath, origin).toString();
}

export function resolveSeoMetadata(routeId: string | null | undefined, locale: Locale): SeoMetadata {
  const routeSeo = routeId ? ROUTE_SEO[routeId] : undefined;
  const localizedDefaults = DEFAULT_LOCALIZED_TEXT[locale] ?? DEFAULT_LOCALIZED_TEXT.en;
  const localizedRouteContent = routeSeo?.localized[locale] ?? routeSeo?.localized.en;

  return {
    title: localizedRouteContent?.title ?? localizedDefaults.title,
    description: localizedRouteContent?.description ?? localizedDefaults.description,
    robots: routeSeo?.robots ?? DEFAULT_ROBOTS,
    ogType: routeSeo?.ogType ?? 'website',
    ogImagePath: routeSeo?.ogImagePath ?? DEFAULT_OG_IMAGE_PATH
  };
}

export function getOpenGraphLocale(locale: Locale): string {
  return OG_LOCALE_BY_LOCALE[locale] ?? 'en_US';
}

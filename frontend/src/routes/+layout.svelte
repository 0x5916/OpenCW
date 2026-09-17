<script lang="ts">
  import { onMount } from 'svelte';
  import '../app.css';
  import { user, initAuth, logout } from '$lib/auth';
  import { flushQueuedProgress, initializeProgressSync } from '$lib/progressSync';
  import { LESSONS } from '$lib/morse';
  import { reconcileSettingsWithServer, touchLocalPageSettingsUpdatedAt } from '$lib/cwSync';
  import { goto, afterNavigate } from '$app/navigation';
  import { page } from '$app/state';
  import {
    Ellipsis,
    Info,
    Languages,
    LayoutDashboard,
    LogIn,
    LogOut,
    MessageSquare,
    Monitor,
    Moon,
    Radio,
    Settings,
    Sun,
    User,
    UserPlus
  } from '@lucide/svelte';
  import {
    lang,
    setLang,
    setLangPreference,
    initLang,
    localizedHref as href
  } from '$lib/i18n.svelte';
  import { locales } from '$lib/paraglide/runtime';
  import Dropdown from '$lib/components/Dropdown.svelte';
  import { getLocaleLongLabel, getLocaleShortLabel } from '$lib/locale';
  import { THEME_CYCLE, applyTheme, readStoredTheme, setTheme } from '$lib/theme';
  import * as m from '$lib/paraglide/messages';
  import type { Locale } from '$lib/i18n.svelte';
  import type { Theme } from '$lib/theme';

  let { children, data } = $props();

  /** Desktop link cluster. Home is the brand logo, so it is not repeated here. */
  const PRIMARY_NAV = [
    { path: '/learn', label: m.nav_learn, icon: Radio },
    { path: '/forum', label: m.nav_forum, icon: MessageSquare },
    { path: '/about', label: m.nav_about, icon: Info }
  ];

  /** Pinned phone tab bar. All four slots are real destinations. */
  const TAB_NAV = [
    { path: '/learn', label: m.nav_learn, icon: Radio },
    { path: '/profile', label: m.nav_profile, icon: LayoutDashboard },
    { path: '/forum', label: m.nav_forum, icon: MessageSquare },
    { path: '/more', label: m.nav_more, icon: Ellipsis }
  ];

  const GITHUB_URL = 'https://github.com/0x5916';

  let theme = $state<Theme>('auto');
  let ThemeIcon = $derived(themeIconFor(theme));
  let reconciledSettingsForUser = $state<string | null>(null);

  const structuredDataScripts = $derived(
    data.seo.structuredData.map((schema) => {
      // eslint-disable-next-line no-useless-escape -- Svelte ends script blocks at a literal closing script tag
      return `<script type="application/ld+json">${JSON.stringify(schema)}<\/script>`;
    })
  );

  // initLang receives the locale the server read from the cookie —
  // so SSR renders the correct language from the very first request.
  // svelte-ignore state_referenced_locally
  initLang(data.locale, data.localePreference);

  $effect(() => {
    theme = readStoredTheme();
    applyTheme(theme);
    initAuth();
  });

  $effect(() => {
    if (!$user) {
      reconciledSettingsForUser = null;
      return;
    }

    void flushQueuedProgress();

    if (reconciledSettingsForUser === $user.username) return;
    reconciledSettingsForUser = $user.username;

    void reconcileSettingsWithServer({
      maxLesson: LESSONS.length,
      fallbackLanguagePreference: data.localePreference,
      onLocale: setLangPreference
    }).catch(() => {
      // Keep app startup/login resilient if reconciliation fails.
    });
  });

  onMount(() => {
    initializeProgressSync();
  });

  function langLabel(locale: Locale): string {
    return getLocaleShortLabel(locale);
  }

  function languageLabel(locale: Locale): string {
    return getLocaleLongLabel(locale);
  }

  function setLanguage(locale: Locale): void {
    setLang(locale);
    touchLocalPageSettingsUpdatedAt();
  }

  // Theme is a device-local preference: changing it must not mark the synced page
  // settings (language, lesson) as locally newer than the server's copy.
  function changeTheme(nextTheme: Theme) {
    theme = setTheme(nextTheme);
  }

  function cycleTheme() {
    changeTheme(THEME_CYCLE[theme]);
  }

  async function handleLogout() {
    await logout();
    await goto(href('/'));
  }

  function stripTrailingSlash(path: string): string {
    return path.replace(/\/+$/, '') || '/';
  }

  function isActive(path: string): boolean {
    return stripTrailingSlash(page.url.pathname) === stripTrailingSlash(href(path));
  }

  function themeIconFor(currentTheme: Theme) {
    if (currentTheme === 'light') return Sun;
    if (currentTheme === 'dark') return Moon;
    return Monitor;
  }

  // Toggle a body class when an input/textarea gains or loses focus so the
  // mobile bottom nav can hide while the soft keyboard is open. Selects are
  // excluded on purpose: they open a picker, not a keyboard, and leaving one
  // focused would keep the tab bar hidden.
  function onDocumentFocusIn(event: FocusEvent) {
    const target = event.target;
    if (target instanceof HTMLElement && target.matches('input, textarea')) {
      document.body.classList.add('keyboard-open');
    }
  }

  function onDocumentFocusOut(event: FocusEvent) {
    const target = event.target;
    if (target instanceof HTMLElement && target.matches('input, textarea')) {
      document.body.classList.remove('keyboard-open');
    }
  }

  // The desktop theme button keeps a stored preference that other pages can
  // change (Settings, the More page), so refresh it on every navigation.
  afterNavigate(() => {
    theme = readStoredTheme();
  });

  $effect(() => {
    if (typeof document === 'undefined') return;

    document.addEventListener('focusin', onDocumentFocusIn);
    document.addEventListener('focusout', onDocumentFocusOut);

    return () => {
      document.removeEventListener('focusin', onDocumentFocusIn);
      document.removeEventListener('focusout', onDocumentFocusOut);
    };
  });
</script>

<svelte:head>
  <title>{data.seo.title}</title>
  <meta name="description" content={data.seo.description} />
  <meta name="robots" content={data.seo.robots} />
  <link rel="canonical" href={data.seo.canonicalUrl} />
  {#each data.seo.alternates as alternate (alternate.locale)}
    <link rel="alternate" hreflang={alternate.locale} href={alternate.href} />
  {/each}
  <link rel="alternate" hreflang="x-default" href={data.seo.xDefaultHref} />
  <meta property="og:site_name" content={data.seo.siteName} />
  <meta property="og:type" content={data.seo.ogType} />
  <meta property="og:title" content={data.seo.title} />
  <meta property="og:description" content={data.seo.description} />
  <meta property="og:url" content={data.seo.canonicalUrl} />
  <meta property="og:locale" content={data.seo.openGraphLocale} />
  <meta property="og:image" content={data.seo.openGraphImage} />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content={data.seo.title} />
  <meta name="twitter:description" content={data.seo.description} />
  <meta name="twitter:image" content={data.seo.openGraphImage} />
  {#each structuredDataScripts as scriptTag (scriptTag)}
    <!-- eslint-disable-next-line svelte/no-at-html-tags -- server-built JSON-LD, no user input -->
    {@html scriptTag}
  {/each}
</svelte:head>

<div class="page-wrapper">
  <nav class="navbar">
    <div class="navbar-inner">
      <!-- Brand -->
      <!-- Brand is the home affordance, so "Home" is not duplicated in the links -->
      <a
        href={href('/')}
        class="navbar-brand"
        title={m.nav_home()}
        aria-current={isActive('/') ? 'page' : undefined}
      >
        <!-- The wordmark next to it already names the link. -->
        <img src="/favicon.svg" alt="" />
        OpenCW
      </a>

      <!-- Desktop: all links + user menu on the right -->
      <div class="navbar-right navbar-desktop">
        {#each PRIMARY_NAV as item (item.path)}
          <a
            href={href(item.path)}
            class="navbar-link"
            class:active={isActive(item.path)}
            aria-current={isActive(item.path) ? 'page' : undefined}>{item.label()}</a
          >
        {/each}
        <div class="navbar-divider"></div>
        {#if $user}
          <Dropdown id="user-menu">
            {#snippet trigger()}
              <User class="nav-icon" aria-hidden="true" />
              {$user.username}
            {/snippet}
            {#snippet menu()}
              <a href={href('/profile')} class="user-dropdown-item" role="menuitem"
                ><LayoutDashboard size={14} style="pointer-events:none" /> {m.nav_profile()}</a
              >
              <a href={href('/settings')} class="user-dropdown-item" role="menuitem"
                ><Settings size={14} style="pointer-events:none" /> {m.nav_settings()}</a
              >
              <button
                type="button"
                onclick={() => void handleLogout()}
                class="user-dropdown-item"
                role="menuitem"
                ><LogOut size={14} style="pointer-events:none" /> {m.nav_logout()}</button
              >
            {/snippet}
          </Dropdown>
        {:else}
          <Dropdown id="guest-menu">
            {#snippet trigger()}
              <User class="nav-icon" aria-hidden="true" />
              {m.nav_guest()}
            {/snippet}
            {#snippet menu()}
              <a href={href('/profile')} class="user-dropdown-item" role="menuitem"
                ><LayoutDashboard size={14} style="pointer-events:none" /> {m.nav_profile()}</a
              >
              <a href={href('/settings')} class="user-dropdown-item" role="menuitem"
                ><Settings size={14} style="pointer-events:none" /> {m.nav_settings()}</a
              >
              <a href={href('/login')} class="user-dropdown-item" role="menuitem"
                ><LogIn size={14} style="pointer-events:none" /> {m.nav_login()}</a
              >
              <a href={href('/register')} class="user-dropdown-item" role="menuitem"
                ><UserPlus size={14} style="pointer-events:none" /> {m.nav_register()}</a
              >
            {/snippet}
          </Dropdown>
        {/if}
        <button
          type="button"
          onclick={cycleTheme}
          class="theme-icon-btn"
          title={m.nav_theme_cycle()}
          aria-label={m.nav_theme_cycle()}
        >
          <span class="nav-label-icon">
            <ThemeIcon class="nav-icon" aria-hidden="true" />
            {theme === 'auto'
              ? m.theme_auto()
              : theme === 'light'
                ? m.theme_light()
                : m.theme_dark()}
          </span>
        </button>
        <Dropdown id="lang-menu" label={m.settings_language_label()}>
          {#snippet trigger()}
            <Languages class="nav-icon" aria-hidden="true" />
            {langLabel(lang.value)}
          {/snippet}
          {#snippet menu()}
            {#each locales as locale (locale)}
              <button
                type="button"
                class="user-dropdown-item"
                role="menuitem"
                onclick={() => setLanguage(locale as Locale)}
              >
                {languageLabel(locale as Locale)}
              </button>
            {/each}
          {/snippet}
        </Dropdown>
      </div>
    </div>
  </nav>

  <main class="page-content">
    {#key lang.value}
      {@render children()}
    {/key}
  </main>

  <footer class="footer">
    <span>{m.footer_text()}</span>
    <span class="footer-links">
      <a href={href('/about')} class="footer-link">{m.nav_about()}</a>
      <a href={GITHUB_URL} class="footer-link" rel="noopener noreferrer" target="_blank"
        >{m.footer_link_github()}</a
      >
    </span>
  </footer>

  <!-- Mobile: bottom navigation bar (hidden on desktop) -->
  <nav class="bottom-nav" aria-label={m.nav_primary()}>
    {#each TAB_NAV as item (item.path)}
      <a
        href={href(item.path)}
        class="bottom-nav-item"
        class:active={isActive(item.path)}
        aria-current={isActive(item.path) ? 'page' : undefined}
        ><item.icon size={20} class="bottom-nav-icon" aria-hidden="true" /><span
          class="bottom-nav-label">{item.label()}</span
        ></a
      >
    {/each}
  </nav>
</div>

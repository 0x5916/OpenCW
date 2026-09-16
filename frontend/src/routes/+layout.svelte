<script lang="ts">
  import { onMount } from 'svelte';
  import favicon from '$lib/assets/favicon.svg';
  import '../app.css';
  import { user, initAuth, logout } from '$lib/auth';
  import { flushQueuedProgress, initializeProgressSync } from '$lib/progressSync';
  import { LESSONS } from '$lib/morse';
  import { reconcileSettingsWithServer, touchLocalPageSettingsUpdatedAt } from '$lib/cwSync';
  import { goto, afterNavigate } from '$app/navigation';
  import { page } from '$app/state';
  import {
    Menu,
    X,
    Monitor,
    Sun,
    Moon,
    Home,
    Radio,
    MessageSquare,
    Info,
    LogIn,
    UserPlus,
    LogOut,
    Languages,
    User,
    Settings,
    LayoutDashboard
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
  import { UI_STORAGE_KEYS } from '$lib/storageKeys';
  import * as m from '$lib/paraglide/messages';
  import type { Locale } from '$lib/i18n.svelte';

  let { children, data } = $props();

  type Theme = 'auto' | 'light' | 'dark';

  const CYCLE: Record<Theme, Theme> = { auto: 'light', light: 'dark', dark: 'auto' };

  function load(): Theme {
    if (typeof localStorage === 'undefined') return 'auto';
    return (localStorage.getItem(UI_STORAGE_KEYS.theme) as Theme) ?? 'auto';
  }

  function apply(t: Theme) {
    if (t === 'auto') document.documentElement.removeAttribute('data-theme');
    else document.documentElement.setAttribute('data-theme', t);
  }

  let theme = $state<Theme>('auto');
  let menuOpen = $state(false);
  let navEl = $state<HTMLElement | null>(null);
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
    theme = load();
    apply(theme);
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

  function setTheme(nextTheme: Theme) {
    theme = nextTheme;
    localStorage.setItem(UI_STORAGE_KEYS.theme, theme);
    apply(theme);
    touchLocalPageSettingsUpdatedAt();
  }

  function cycleTheme() {
    setTheme(CYCLE[theme]);
  }

  async function handleLogout() {
    await logout();
    await goto(href('/'));
  }

  function isActive(path: string): boolean {
    const strip = (p: string) => p.replace(/\/+$/, '') || '/';
    return strip(page.url.pathname) === strip(href(path));
  }

  function themeIconFor(currentTheme: Theme) {
    if (currentTheme === 'light') return Sun;
    if (currentTheme === 'dark') return Moon;
    return Monitor;
  }

  function closeMobileMenu() {
    menuOpen = false;
  }

  function onDocumentClick(event: MouseEvent) {
    const target = event.target;
    if (!(target instanceof Node)) return;

    if (menuOpen && navEl && !navEl.contains(target)) {
      menuOpen = false;
    }
  }

  function onDocumentKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      closeMobileMenu();
    }
  }

  // Toggle a body class when an input/textarea/select gains or loses focus so
  // the mobile bottom nav can hide while the soft keyboard is open.
  function onDocumentFocusIn(event: FocusEvent) {
    const target = event.target;
    if (target instanceof HTMLElement && target.matches('input, textarea, select')) {
      document.body.classList.add('keyboard-open');
    }
  }

  function onDocumentFocusOut(event: FocusEvent) {
    const target = event.target;
    if (target instanceof HTMLElement && target.matches('input, textarea, select')) {
      document.body.classList.remove('keyboard-open');
    }
  }

  afterNavigate(() => {
    closeMobileMenu();
  });

  $effect(() => {
    if (typeof document === 'undefined') return;

    document.addEventListener('click', onDocumentClick);
    document.addEventListener('keydown', onDocumentKeydown);
    document.addEventListener('focusin', onDocumentFocusIn);
    document.addEventListener('focusout', onDocumentFocusOut);

    return () => {
      document.removeEventListener('click', onDocumentClick);
      document.removeEventListener('keydown', onDocumentKeydown);
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
  <link rel="icon" href={favicon} />
</svelte:head>

<div class="page-wrapper">
  <nav class="navbar" bind:this={navEl}>
    <div class="navbar-inner">
      <!-- Brand -->
      <a href={href('/')} class="navbar-brand">
        <img src={favicon} alt="OpenCW" />
        OpenCW
      </a>

      <!-- Desktop: all links + user menu on the right -->
      <div class="navbar-right navbar-desktop">
        <a href={href('/')} class="navbar-link">{m.nav_home()}</a>
        <a href={href('/morse/learn')} class="navbar-link">{m.nav_learn()}</a>
        <a href={href('/forum')} class="navbar-link">{m.nav_forum()}</a>
        <a href={href('/about')} class="navbar-link">{m.nav_about()}</a>
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
              Guest
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
          title="Cycle theme"
          aria-label="Cycle theme"
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

      <!-- Mobile: hamburger only -->
      <div class="navbar-mobile-controls">
        <button
          type="button"
          onclick={() => (menuOpen = !menuOpen)}
          class="hamburger"
          aria-label="Menu"
          aria-expanded={menuOpen}
          aria-controls="mobile-nav-menu"
        >
          {#if menuOpen}
            <X class="nav-icon" aria-hidden="true" />
          {:else}
            <Menu class="nav-icon" aria-hidden="true" />
          {/if}
        </button>
      </div>
    </div>

    <!-- Mobile dropdown menu -->
    {#if menuOpen}
      <div class="mobile-menu" id="mobile-nav-menu">
        <a href={href('/')} class="mobile-link" onclick={() => (menuOpen = false)}
          ><Home size={16} />{m.nav_home()}</a
        >
        <a href={href('/morse/learn')} class="mobile-link" onclick={() => (menuOpen = false)}
          ><Radio size={16} />{m.nav_learn()}</a
        >
        <a href={href('/forum')} class="mobile-link" onclick={() => (menuOpen = false)}
          ><MessageSquare size={16} />{m.nav_forum()}</a
        >
        <a href={href('/about')} class="mobile-link" onclick={() => (menuOpen = false)}
          ><Info size={16} />{m.nav_about()}</a
        >
        <div class="mobile-divider"></div>
        {#if $user}
          <a href={href('/profile')} class="mobile-link" onclick={() => (menuOpen = false)}
            ><LayoutDashboard size={16} />{m.nav_profile()}</a
          >
          <a href={href('/settings')} class="mobile-link" onclick={() => (menuOpen = false)}
            ><Settings size={16} />{m.nav_settings()}</a
          >
          <button
            type="button"
            onclick={() => {
              handleLogout();
              menuOpen = false;
            }}
            class="mobile-link mobile-link-btn"
            ><LogOut size={16} />{m.nav_logout()} ({$user.username})</button
          >
        {:else}
          <a href={href('/profile')} class="mobile-link" onclick={() => (menuOpen = false)}
            ><LayoutDashboard size={16} />{m.nav_profile()}</a
          >
          <a href={href('/settings')} class="mobile-link" onclick={() => (menuOpen = false)}
            ><Settings size={16} />{m.nav_settings()}</a
          >
          <a href={href('/login')} class="mobile-link" onclick={() => (menuOpen = false)}
            ><LogIn size={16} />{m.nav_login()}</a
          >
          <a href={href('/register')} class="mobile-link" onclick={() => (menuOpen = false)}
            ><UserPlus size={16} />{m.nav_register()}</a
          >
        {/if}
        <div class="mobile-divider"></div>
        <button type="button" onclick={cycleTheme} class="mobile-link mobile-link-btn"
          ><ThemeIcon size={16} />{theme === 'auto'
            ? m.theme_auto()
            : theme === 'light'
              ? m.theme_light()
              : m.theme_dark()}</button
        >
        <div class="mobile-divider"></div>
        {#each locales as locale (locale)}
          <button
            type="button"
            class="mobile-link mobile-link-btn"
            onclick={() => setLanguage(locale)}
            ><Languages size={16} />{languageLabel(locale)}</button
          >
        {/each}
      </div>
    {/if}
  </nav>

  <main class="page-content">
    {#key lang.value}
      {@render children()}
    {/key}
  </main>

  <footer class="footer">{m.footer_text()}</footer>

  <!-- Mobile: bottom navigation bar (hidden on desktop) -->
  <nav class="bottom-nav" aria-label={m.nav_primary()}>
    <a
      href={href('/')}
      class="bottom-nav-item"
      class:active={isActive('/')}
      aria-current={isActive('/') ? 'page' : undefined}
      ><Home size={20} class="bottom-nav-icon" aria-hidden="true" /><span class="bottom-nav-label"
        >{m.nav_home()}</span
      ></a
    >
    <a
      href={href('/morse/learn')}
      class="bottom-nav-item"
      class:active={isActive('/morse/learn')}
      aria-current={isActive('/morse/learn') ? 'page' : undefined}
      ><Radio size={20} class="bottom-nav-icon" aria-hidden="true" /><span class="bottom-nav-label"
        >{m.nav_learn()}</span
      ></a
    >
    <a
      href={href('/forum')}
      class="bottom-nav-item"
      class:active={isActive('/forum')}
      aria-current={isActive('/forum') ? 'page' : undefined}
      ><MessageSquare size={20} class="bottom-nav-icon" aria-hidden="true" /><span
        class="bottom-nav-label">{m.nav_forum()}</span
      ></a
    >
    <a
      href={href('/about')}
      class="bottom-nav-item"
      class:active={isActive('/about')}
      aria-current={isActive('/about') ? 'page' : undefined}
      ><Info size={20} class="bottom-nav-icon" aria-hidden="true" /><span class="bottom-nav-label"
        >{m.nav_about()}</span
      ></a
    >
  </nav>
</div>

<script lang="ts">
  import { onMount } from 'svelte';
  import favicon from '$lib/assets/favicon.svg';
  import '../app.css';
  import { user, initAuth, logout } from '$lib/auth';
  import { flushQueuedProgress, initializeProgressSync } from '$lib/progressSync';
  import { LESSONS } from '$lib/morse';
  import { reconcileSettingsWithServer, touchLocalPageSettingsUpdatedAt } from '$lib/cwSync';
  import { goto, afterNavigate } from '$app/navigation';
  import { registerSW } from 'virtual:pwa-register';
  import {
    ChevronDown,
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
  } from 'lucide-svelte';
  import { lang, setLang, setLangPreference, initLang } from '$lib/i18n.svelte';
  import { locales, localizeHref } from '$lib/paraglide/runtime';
  import { getLocaleLongLabel, getLocaleShortLabel } from '$lib/locale';
  import * as m from '$lib/paraglide/messages';
  import type { Locale } from '$lib/i18n.svelte';

  let { children, data } = $props();

  type Theme = 'auto' | 'light' | 'dark';

  const CYCLE: Record<Theme, Theme> = { auto: 'light', light: 'dark', dark: 'auto' };

  function load(): Theme {
    if (typeof localStorage === 'undefined') return 'auto';
    return (localStorage.getItem('theme') as Theme) ?? 'auto';
  }

  function apply(t: Theme) {
    if (t === 'auto') document.documentElement.removeAttribute('data-theme');
    else document.documentElement.setAttribute('data-theme', t);
  }

  let theme = $state<Theme>('auto');
  let menuOpen = $state(false);
  let userMenuOpen = $state(false);
  let guestMenuOpen = $state(false);
  let userMenuLeaveTimer = 0;
  let guestMenuLeaveTimer = 0;
  let langMenuLeaveTimer = 0;
  let navEl = $state<HTMLElement | null>(null);
  let userMenuEl = $state<HTMLElement | null>(null);
  let guestMenuEl = $state<HTMLElement | null>(null);
  let langMenuEl = $state<HTMLElement | null>(null);
  let langMenuOpen = $state(false);
  let ThemeIcon = $derived(themeIconFor(theme));
  let reconciledSettingsForUser = $state<string | null>(null);

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

    const updateServiceWorker = registerSW({
      immediate: true
    });

    if ('serviceWorker' in navigator) {
      const localePrefixes = ['/en', '/zh-Hant', '/zh-Hans', '/ja', '/de'];

      void (async () => {
        try {
          const registrations = await navigator.serviceWorker.getRegistrations();

          for (const registration of registrations) {
            const scriptUrl = registration.active?.scriptURL ?? registration.installing?.scriptURL ?? registration.waiting?.scriptURL;
            const scriptPath = scriptUrl ? new URL(scriptUrl).pathname : '';
            const scopePath = new URL(registration.scope).pathname.replace(/\/$/, '') || '/';
            const scopeIsLocale = localePrefixes.some((prefix) => scopePath === prefix || scopePath.startsWith(`${prefix}/`));
            const scriptIsLegacy = scriptPath === '/service-worker.js' || scriptPath.endsWith('/sw.js') && scriptPath !== '/sw.js';

            if (scopeIsLocale || scriptIsLegacy) {
              await registration.unregister();
            }
          }
        } catch {
          // Ignore cleanup failures and continue with root registration.
        }

        updateServiceWorker();
      })().catch(() => {
        // Avoid breaking page initialization if SW registration fails.
      });

      return;
    }

    updateServiceWorker();
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
    langMenuOpen = false;
  }

  function setTheme(nextTheme: Theme) {
    theme = nextTheme;
    localStorage.setItem('theme', theme);
    apply(theme);
    touchLocalPageSettingsUpdatedAt();
  }

  function cycleTheme() {
    setTheme(CYCLE[theme]);
  }

  async function handleLogout() {
    await logout();
    await goto(localizeHref('/', { locale: lang.value }));
  }

  function href(path: string) {
    return localizeHref(path, { locale: lang.value });
  }

  function themeIconFor(currentTheme: Theme) {
    if (currentTheme === 'light') return Sun;
    if (currentTheme === 'dark') return Moon;
    return Monitor;
  }

  function closeMenus() {
    menuOpen = false;
    userMenuOpen = false;
    guestMenuOpen = false;
    langMenuOpen = false;
  }

  function onDocumentClick(event: MouseEvent) {
    const target = event.target;
    if (!(target instanceof Node)) return;

    if (menuOpen && navEl && !navEl.contains(target)) {
      menuOpen = false;
    }

    if (userMenuOpen && userMenuEl && !userMenuEl.contains(target)) {
      userMenuOpen = false;
    }

    if (guestMenuOpen && guestMenuEl && !guestMenuEl.contains(target)) {
      guestMenuOpen = false;
    }

    if (langMenuOpen && langMenuEl && !langMenuEl.contains(target)) {
      langMenuOpen = false;
    }
  }

  function onDocumentKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      closeMenus();
    }
  }

  afterNavigate(() => {
    closeMenus();
  });

  $effect(() => {
    if (typeof document === 'undefined') return;

    document.addEventListener('click', onDocumentClick);
    document.addEventListener('keydown', onDocumentKeydown);

    return () => {
      document.removeEventListener('click', onDocumentClick);
      document.removeEventListener('keydown', onDocumentKeydown);
    };
  });
</script>

<svelte:head>
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
          <div
            class="user-menu-wrapper"
            role="group"
            bind:this={userMenuEl}
            onmouseenter={() => {
              clearTimeout(userMenuLeaveTimer);
              userMenuOpen = true;
            }}
            onmouseleave={() => {
              userMenuLeaveTimer = window.setTimeout(() => (userMenuOpen = false), 150);
            }}
          >
            <button
              type="button"
              onclick={() => (userMenuOpen = !userMenuOpen)}
              class="navbar-user-btn"
              aria-expanded={userMenuOpen}
              aria-haspopup="menu"
              aria-controls="user-menu"
            >
              <span class="nav-label-icon">
                <User class="nav-icon" aria-hidden="true" />
                {$user.username}
                <ChevronDown class="nav-icon" aria-hidden="true" />
              </span>
            </button>
            {#if userMenuOpen}
              <div class="user-dropdown" id="user-menu" role="menu">
                <a
                  href={href('/profile')}
                  onclick={() => (userMenuOpen = false)}
                  class="user-dropdown-item"
                  role="menuitem"
                  ><LayoutDashboard size={14} style="pointer-events:none" /> {m.nav_profile()}</a
                >
                <a
                  href={href('/settings')}
                  onclick={() => (userMenuOpen = false)}
                  class="user-dropdown-item"
                  role="menuitem"
                  ><Settings size={14} style="pointer-events:none" /> {m.nav_settings()}</a
                >
                <button
                  type="button"
                  onclick={() => {
                    void handleLogout();
                    userMenuOpen = false;
                  }}
                  class="user-dropdown-item"
                  role="menuitem"
                  ><LogOut size={14} style="pointer-events:none" /> {m.nav_logout()}</button
                >
              </div>
            {/if}
          </div>
        {:else}
          <div
            class="user-menu-wrapper"
            role="group"
            bind:this={guestMenuEl}
            onmouseenter={() => {
              clearTimeout(guestMenuLeaveTimer);
              guestMenuOpen = true;
            }}
            onmouseleave={() => {
              guestMenuLeaveTimer = window.setTimeout(() => (guestMenuOpen = false), 150);
            }}
          >
            <button
              type="button"
              onclick={() => (guestMenuOpen = !guestMenuOpen)}
              class="navbar-user-btn"
              aria-expanded={guestMenuOpen}
              aria-haspopup="menu"
              aria-controls="guest-menu"
            >
              <span class="nav-label-icon">
                <User class="nav-icon" aria-hidden="true" />
                Guest
                <ChevronDown class="nav-icon" aria-hidden="true" />
              </span>
            </button>
            {#if guestMenuOpen}
              <div class="user-dropdown" id="guest-menu" role="menu">
                <a
                  href={href('/profile')}
                  class="user-dropdown-item"
                  role="menuitem"
                  onclick={() => (guestMenuOpen = false)}
                  ><LayoutDashboard size={14} style="pointer-events:none" /> {m.nav_profile()}</a
                >
                <a
                  href={href('/settings')}
                  class="user-dropdown-item"
                  role="menuitem"
                  onclick={() => (guestMenuOpen = false)}
                  ><Settings size={14} style="pointer-events:none" /> {m.nav_settings()}</a
                >
                <a
                  href={href('/login')}
                  class="user-dropdown-item"
                  role="menuitem"
                  onclick={() => (guestMenuOpen = false)}
                  ><LogIn size={14} style="pointer-events:none" /> {m.nav_login()}</a
                >
                <a
                  href={href('/register')}
                  class="user-dropdown-item"
                  role="menuitem"
                  onclick={() => (guestMenuOpen = false)}
                  ><UserPlus size={14} style="pointer-events:none" /> {m.nav_register()}</a
                >
              </div>
            {/if}
          </div>
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
        <div
          class="user-menu-wrapper"
          role="group"
          bind:this={langMenuEl}
          onmouseenter={() => {
            clearTimeout(langMenuLeaveTimer);
            langMenuOpen = true;
          }}
          onmouseleave={() => {
            langMenuLeaveTimer = window.setTimeout(() => (langMenuOpen = false), 150);
          }}
        >
          <button
            type="button"
            onclick={() => (langMenuOpen = !langMenuOpen)}
            class="navbar-user-btn"
            aria-expanded={langMenuOpen}
            aria-haspopup="menu"
            aria-controls="lang-menu"
            title={m.settings_language_label()}
            aria-label={m.settings_language_label()}
          >
            <span class="nav-label-icon">
              <Languages class="nav-icon" aria-hidden="true" />
              {langLabel(lang.value)}
              <ChevronDown class="nav-icon" aria-hidden="true" />
            </span>
          </button>
          {#if langMenuOpen}
            <div class="user-dropdown" id="lang-menu" role="menu">
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
            </div>
          {/if}
        </div>
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
          <button type="button" class="mobile-link mobile-link-btn" onclick={() => setLanguage(locale)}
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
</div>

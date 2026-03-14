<script lang="ts">
  import favicon from '$lib/assets/favicon.svg';
  import '../app.css';
  import { user, initAuth, logout } from '$lib/auth';
  import { goto, afterNavigate } from '$app/navigation';
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
  import { lang, setLang, initLang } from '$lib/i18n.svelte';
  import { locales } from '$lib/paraglide/runtime';
  import * as m from '$lib/paraglide/messages';
  import type { Locale } from '$lib/i18n.svelte';

  let { children, data } = $props();

  type Theme = 'auto' | 'light' | 'dark';

  const LANG_LABELS: Partial<Record<Locale, string>> = {
    en: 'EN',
    'zh-Hant': 'ZH-T',
    'zh-Hans': 'ZH-S'
  };
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
  let navEl = $state<HTMLElement | null>(null);
  let userMenuEl = $state<HTMLElement | null>(null);
  let guestMenuEl = $state<HTMLElement | null>(null);
  let ThemeIcon = $derived(themeIconFor(theme));

  // initLang receives the locale the server read from the cookie —
  // so SSR renders the correct language from the very first request.
  // svelte-ignore state_referenced_locally
  initLang(data.locale);

  $effect(() => {
    theme = load();
    apply(theme);
    initAuth();
  });

  function toggleLang() {
    const current = locales.indexOf(lang.value);
    const next = locales[(current + 1) % locales.length] as Locale;
    setLang(next);
  }

  function langLabel(locale: Locale) {
    return LANG_LABELS[locale] ?? locale.toUpperCase();
  }

  function setTheme(nextTheme: Theme) {
    theme = nextTheme;
    localStorage.setItem('theme', theme);
    apply(theme);
  }

  function cycleTheme() {
    setTheme(CYCLE[theme]);
  }

  function handleLogout() {
    logout();
    goto('/');
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
      <a href="/" class="navbar-brand">
        <img src={favicon} alt="OpenCW" />
        OpenCW
      </a>

      <!-- Desktop: all links + user menu on the right -->
      <div class="navbar-right navbar-desktop">
        <a href="/" class="navbar-link">{m.nav_home()}</a>
        <a href="/morse/learn" class="navbar-link">{m.nav_learn()}</a>
        <a href="/forum" class="navbar-link">{m.nav_forum()}</a>
        <a href="/about" class="navbar-link">{m.nav_about()}</a>
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
                  href="/profile"
                  onclick={() => (userMenuOpen = false)}
                  class="user-dropdown-item"
                  role="menuitem"
                  ><LayoutDashboard size={14} style="pointer-events:none" /> {m.nav_profile()}</a
                >
                <a
                  href="/settings"
                  onclick={() => (userMenuOpen = false)}
                  class="user-dropdown-item"
                  role="menuitem"
                  ><Settings size={14} style="pointer-events:none" /> {m.nav_settings()}</a
                >
                <button
                  type="button"
                  onclick={() => {
                    handleLogout();
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
                {m.nav_login()}
                <ChevronDown class="nav-icon" aria-hidden="true" />
              </span>
            </button>
            {#if guestMenuOpen}
              <div class="user-dropdown" id="guest-menu" role="menu">
                <a
                  href="/login"
                  class="user-dropdown-item"
                  role="menuitem"
                  onclick={() => (guestMenuOpen = false)}>{m.nav_login()}</a
                >
                <a
                  href="/register"
                  class="user-dropdown-item"
                  role="menuitem"
                  onclick={() => (guestMenuOpen = false)}>{m.nav_register()}</a
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
        <button
          type="button"
          onclick={toggleLang}
          class="theme-icon-btn"
          title="Switch language"
          aria-label="Switch language"
        >
          <span class="nav-label-icon"
            ><Languages class="nav-icon" aria-hidden="true" />{langLabel(lang.value)}</span
          >
        </button>
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
        <a href="/" class="mobile-link" onclick={() => (menuOpen = false)}
          ><Home size={16} />{m.nav_home()}</a
        >
        <a href="/morse/learn" class="mobile-link" onclick={() => (menuOpen = false)}
          ><Radio size={16} />{m.nav_learn()}</a
        >
        <a href="/forum" class="mobile-link" onclick={() => (menuOpen = false)}
          ><MessageSquare size={16} />{m.nav_forum()}</a
        >
        <a href="/about" class="mobile-link" onclick={() => (menuOpen = false)}
          ><Info size={16} />{m.nav_about()}</a
        >
        <div class="mobile-divider"></div>
        {#if $user}
          <a href="/profile" class="mobile-link" onclick={() => (menuOpen = false)}
            ><LayoutDashboard size={16} />{m.nav_profile()}</a
          >
          <a href="/settings" class="mobile-link" onclick={() => (menuOpen = false)}
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
          <a href="/login" class="mobile-link" onclick={() => (menuOpen = false)}
            ><LogIn size={16} />{m.nav_login()}</a
          >
          <a href="/register" class="mobile-link" onclick={() => (menuOpen = false)}
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
        <button type="button" onclick={toggleLang} class="mobile-link mobile-link-btn"
          ><Languages size={16} />{langLabel(lang.value)}</button
        >
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

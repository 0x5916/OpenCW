<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { House, Info, LogIn, LogOut, Settings, UserPlus } from '@lucide/svelte';
  import { logout, user } from '$lib/auth';
  import { touchLocalPageSettingsUpdatedAt } from '$lib/cwSync';
  import { localizedHref as href, setLangPreference, langPreference } from '$lib/i18n.svelte';
  import { getLocaleLongLabel } from '$lib/locale';
  import { locales } from '$lib/paraglide/runtime';
  import { readStoredTheme, setTheme, type Theme } from '$lib/theme';
  import * as m from '$lib/paraglide/messages';
  import type { Locale, LocalePreference } from '$lib/i18n.svelte';

  // Mirrors the `max-width: 639px` block in layout.css, where the top bar and the
  // tab bar hand over to each other.
  const DESKTOP_QUERY = '(min-width: 640px)';

  let theme = $state<Theme>('auto');
  let language = $state<LocalePreference>('auto');

  // Read the stored preferences on the client only, so the server-rendered
  // defaults hydrate cleanly.
  onMount(() => {
    theme = readStoredTheme();
    language = langPreference.value;
  });

  // This is a phone-only surface: on desktop every entry below already lives in
  // the top bar, so send the visitor to the settings page instead.
  $effect(() => {
    if (typeof window === 'undefined') return;
    if (window.matchMedia(DESKTOP_QUERY).matches) {
      void goto(href('/settings'), { replaceState: true });
    }
  });

  // Theme is device-local, so it does not touch the synced page-settings timestamp.
  function changeTheme(next: Theme) {
    theme = setTheme(next);
  }

  function changeLanguage(next: LocalePreference) {
    language = next;
    setLangPreference(next);
    touchLocalPageSettingsUpdatedAt();
  }

  function languageLabel(locale: Locale): string {
    return getLocaleLongLabel(locale);
  }

  async function handleLogout() {
    await logout();
    await goto(href('/'));
  }
</script>

<div class="page-narrow">
  <!-- Same heading shape as the Settings page: mark + page title. -->
  <div class="more-heading">
    <img src="/favicon.svg" alt="" class="more-logo" />
    <h1 class="page-title">{m.nav_more()}</h1>
  </div>

  <section class="panel">
    <h2 class="card-title">{m.more_app_section()}</h2>
    <div class="more-list">
      <a href={href('/')} class="more-link"
        ><House size={16} aria-hidden="true" /><span class="more-link-text">{m.nav_home()}</span></a
      >
      <a href={href('/about')} class="more-link"
        ><Info size={16} aria-hidden="true" /><span class="more-link-text">{m.nav_about()}</span></a
      >
      <a href={href('/settings')} class="more-link"
        ><Settings size={16} aria-hidden="true" /><span class="more-link-text"
          >{m.nav_settings()}</span
        ></a
      >
    </div>
  </section>

  <section class="panel">
    <h2 class="card-title">{m.settings_account_section()}</h2>
    <div class="more-list">
      {#if $user}
        <button type="button" class="more-link" onclick={() => void handleLogout()}
          ><LogOut size={16} aria-hidden="true" /><span class="more-link-text"
            >{m.nav_logout()}</span
          ></button
        >
      {:else}
        <a href={href('/login')} class="more-link"
          ><LogIn size={16} aria-hidden="true" /><span class="more-link-text">{m.nav_login()}</span
          ></a
        >
        <a href={href('/register')} class="more-link"
          ><UserPlus size={16} aria-hidden="true" /><span class="more-link-text"
            >{m.nav_register()}</span
          ></a
        >
      {/if}
    </div>
  </section>

  <section class="panel">
    <h2 class="card-title">{m.settings_page_section()}</h2>
    <div class="more-form">
      <label class="field">
        <span class="label-text">{m.settings_theme_label()}</span>
        <select
          class="input"
          value={theme}
          onchange={(e) => changeTheme(e.currentTarget.value as Theme)}
        >
          <option value="auto">{m.theme_auto()}</option>
          <option value="light">{m.theme_light()}</option>
          <option value="dark">{m.theme_dark()}</option>
        </select>
      </label>
      <label class="field">
        <span class="label-text">{m.settings_language_label()}</span>
        <select
          class="input"
          value={language}
          onchange={(e) => changeLanguage(e.currentTarget.value as LocalePreference)}
        >
          <option value="auto">{m.theme_auto()}</option>
          {#each locales as locale (locale)}
            <option value={locale}>{languageLabel(locale)}</option>
          {/each}
        </select>
      </label>
    </div>
  </section>
</div>

<style>
  /* The app mark opens this phone-only surface; its desktop counterpart
     (Settings) opens with the shared eyebrow above the title instead. */
  .more-heading {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }

  .more-logo {
    width: 2rem;
    height: 2rem;
    flex-shrink: 0;
  }

  .more-list {
    display: flex;
    flex-direction: column;
    margin-top: 0.5rem;
  }

  /* The panel already draws the group's boundary, so a row carries no box of
     its own — just a hairline from the next one. Stacking a bordered, inset
     control-shaped box per row inside a bordered panel was a box of boxes.
     (Real controls, like the selects below, still get `.input` chrome: their
     edge is an affordance, not decoration.) */
  .more-link {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    width: 100%;
    padding: 0.7rem 0;
    font-family: inherit;
    font-size: var(--text-sm);
    line-height: var(--leading-snug);
    text-align: left;
    text-decoration: none;
    border: none;
    border-bottom: 1px solid var(--border);
    background: none;
    color: var(--text-primary);
    cursor: pointer;
    transition: color 0.15s;
  }

  .more-link:last-child {
    border-bottom: none;
  }

  .more-link:hover {
    color: var(--accent);
  }

  .more-link-text {
    flex: 1 1 auto;
  }

  /* Same rhythm as `.settings-form`. */
  .more-form {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 0.5rem;
  }
</style>

<script lang="ts">
  import { afterNavigate } from '$app/navigation';
  import { ChevronDown } from '@lucide/svelte';
  import type { Snippet } from 'svelte';

  interface Props {
    /** id of the menu element, referenced by the trigger's aria-controls */
    id: string;
    /** optional accessible label for the trigger button */
    label?: string;
    /** trigger button content (icon + text; the chevron is added automatically) */
    trigger: Snippet;
    /** menu items (role="menuitem" is expected on each) */
    menu: Snippet;
  }

  let { id, label = '', trigger, menu }: Props = $props();

  let open = $state(false);
  let wrapperEl = $state<HTMLElement | null>(null);
  let triggerEl = $state<HTMLButtonElement | null>(null);
  let leaveTimer = 0;
  let openedByHover = false;

  function canHover(): boolean {
    return (
      typeof window !== 'undefined' &&
      window.matchMedia('(hover: hover) and (pointer: fine)').matches
    );
  }

  function close() {
    open = false;
    openedByHover = false;
  }

  function toggle() {
    // A click directly after a hover-open should confirm the menu rather than
    // dismiss it (touch devices can fire mouseenter before click).
    if (open && openedByHover) {
      openedByHover = false;
      return;
    }
    open = !open;
    openedByHover = false;
  }

  function onMouseEnter() {
    if (!canHover()) return;
    clearTimeout(leaveTimer);
    if (!open) openedByHover = true;
    open = true;
  }

  function onMouseLeave() {
    if (!canHover()) return;
    clearTimeout(leaveTimer);
    leaveTimer = window.setTimeout(close, 150);
  }

  function onDocumentClick(event: MouseEvent) {
    const target = event.target;
    if (!open || !wrapperEl || !(target instanceof Node)) return;

    // Close on outside clicks and on menu-item clicks; the trigger itself is
    // excluded because its own click handler toggles the menu.
    if (!wrapperEl.contains(target) || !triggerEl?.contains(target)) close();
  }

  function onDocumentKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') close();
  }

  // Close on back/forward and client-side navigation.
  afterNavigate(close);

  $effect(() => {
    document.addEventListener('click', onDocumentClick);
    document.addEventListener('keydown', onDocumentKeydown);

    return () => {
      document.removeEventListener('click', onDocumentClick);
      document.removeEventListener('keydown', onDocumentKeydown);
    };
  });
</script>

<div
  class="user-menu-wrapper"
  role="group"
  bind:this={wrapperEl}
  onmouseenter={onMouseEnter}
  onmouseleave={onMouseLeave}
>
  <button
    type="button"
    bind:this={triggerEl}
    onclick={toggle}
    class="navbar-user-btn"
    aria-expanded={open}
    aria-haspopup="menu"
    aria-controls={id}
    title={label || undefined}
    aria-label={label || undefined}
  >
    <span class="nav-label-icon">
      {@render trigger()}
      <ChevronDown class="nav-icon" aria-hidden="true" />
    </span>
  </button>
  {#if open}
    <div class="user-dropdown" {id} role="menu">
      {@render menu()}
    </div>
  {/if}
</div>

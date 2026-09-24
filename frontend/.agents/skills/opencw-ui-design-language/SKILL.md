---
name: opencw-ui-design-language
description: "Use when designing, changing, reviewing, or verifying OpenCW's site-wide UI: visual design language, shared CSS tokens/classes, layout shell, navigation, responsive behavior, cards versus row/list layouts, forms, buttons, typography, colors, contrast, spacing, borders, or visual consistency across pages such as forum, settings, profile, auth, about, and home. Do not use for backend/API-only work, data-model changes, non-visual bug fixes, or narrowly trainer-specific UI behavior covered by opencw-trainer-ui unless the request explicitly asks for cross-site visual consistency."
argument-hint: "Describe the OpenCW UI design, review, or visual consistency task"
---

# OpenCW UI Design Language

Apply OpenCW's existing visual language when making or reviewing site-wide UI. Preserve page-specific product differences: consistency means shared principles and primitives, not making the trainer, forum, settings, auth, and marketing pages look identical.

## Inspect before changing

Read the relevant route/component first, then cross-check these sources:

- `src/app.css` — design tokens and shared primitives: colors, type, spacing, cards/panels, buttons, inputs, chips, metrics, row lists, focus, contrast modes.
- `src/lib/styles/layout.css` — page shell, width tiers, desktop navbar, mobile bottom navigation, footer, safe-area and bottom-nav clearance.
- `src/routes/+layout.svelte` — global navigation structure and mobile tab destinations.
- `src/lib/components/` — small shared components: `AuthCard`, `Dropdown`, `ErrorAlert`, `GuestNotice`, `LoadingSpinner`, `SaveButton`; there is no generic component-library layer for every button/card/input.
- Representative pages before broad visual edits:
  - `src/routes/+page.svelte` for marketing rhythm and masthead/CTA treatment.
  - `src/routes/forum/+page.svelte` and `src/routes/forum/[id]/+page.svelte` for discussion lists, threads, reply trees, quiet metadata, and inline composers.
  - `src/routes/settings/+page.svelte` for dense forms and ledger-like sections.
  - `src/routes/more/+page.svelte` for phone-only navigation surfaces.
  - `src/routes/morse/learn/+page.svelte`, `MorsePlayer.svelte`, and `ResultOverlay.svelte` only for cross-site consistency; load `opencw-trainer-ui` for trainer-specific behavior or layout.

## Core language

- Use the warm Morse-workshop palette from `src/app.css`: dark/light theme variables, amber as the single product accent, and semantic status colors for correctness/warnings/errors. Do not hard-code new hex colors unless adding a token with a clear role.
- Maintain contrast roles: `--text-primary` for reading content, `--text-secondary` for supporting metadata, `--text-muted` for quiet labels. Use `--border-control` for form edges that need non-text contrast; use `--border`/`--border-subtle` for decorative hairlines.
- Keep typography simple: UI system font for interface/prose, monospace only for Morse text, timers, WPM/Hz, lesson numbers, call signs, tabular stats, and code-like values.
- Follow the 4px spacing scale (`--space-*`) and page rhythm tokens (`--section-gap`, `--block-gap`). Let `.page-content` own page top/bottom padding; do not add route-level padding to compensate for the navbar or phone tab bar.
- Prefer flat, mostly rectangular surfaces. Radii are small; shadows belong to menus, popovers, dialogs, and overlays, not static cards.

## Layout and hierarchy

- Choose width by content: default container for normal pages, `.page-narrow` for focused forms/details/settings/thread reading. Add `.page-stack` when a page needs the block rhythm on its children without the narrow measure (profile).
- Let text and primary task content lead. Use headings, whitespace, alignment, and quiet metadata before adding boxes, borders, badges, or extra copy.
- Use cards/panels when content is an independent unit or dense bounded workspace. Use simpler vertical stacks, `.row-list`, or ledger-like sections for lists, settings, histories, and discussion flows.
- Avoid “box of boxes”: do not wrap every row in a card, stack borders directly against section gaps, or add separators when spacing already groups content.
- Keep page differences useful:
  - Trainer can feel like a focused console/workspace.
  - Forum should read like conversations and rows, not dashboard cards.
  - Settings should read like a ledger of form sections.
  - Auth can remain a centered card.
  - Marketing/about can use larger rhythm and CTA emphasis.

## Controls and interactions

- Use shared classes before inventing new controls: `.btn-primary`, `.btn-ghost`, `.btn-cta`, `.btn-icon` (one shared base, plus a shared size tier for the two text-sized buttons), `.input`, `.select`, `.textarea`, `.field`, `.notice`, `.chip`, `.callsign`, `.thread-cat`, `.metric`, `.panel-label`, `.composer-form` / `.composer-textarea` / `.composer-actions` (with `.is-inline`), `.gate-notice` / `.gate-note`, `.state-block` / `.state-title` / `.state-note` / `.state-error`, `.skeleton-rows` with `.skeleton-row` / `.skeleton-title` / `.skeleton-line` / `.skeleton-meta`, `.row-list`.
- One row language: every whole-row link list — home destinations, the Morse hub, and the forum thread list — is a `.row-list` of `.row-link`s: hairline separators, the shared row padding, and the floor + 2px amber edge on hover and keyboard focus. Do not give one list its own padding, radius or separator strategy. The quieter treatments are deliberate and local, not drift: dropdown items, the profile history table, the desktop navbar, the bottom tab bar, and `/more` rows (colour-only hover).
- Labels share one recipe: `.card-label` (which keeps its bottom margin), `.panel-label` and `.metric-label` (margin-free) are the same caps label. Do not author a fourth.
- Whole-row links (`.row-link`) read as one target: hover and keyboard focus floor the row with `--bg-inset` and run the shared 2px amber edge. Destructive row actions stay quiet (unbordered, danger on hover) rather than filled buttons.
- Top-bar controls share one treatment: a secondary-ink label over a 2px amber rule that grows on hover and `:focus-visible`, and stays put while that control's dropdown is open. The rule is an absolutely positioned child, so the shared block in `layout.css` declares `position: relative` for the controls it applies to — a new top-bar control must join that selector list, or the sticky `.navbar` becomes the containing block and the rule paints across the whole bar. The rule stays local to its control; never draw a full-width bar indicator.
- Keep hover/focus restrained: color, background, border, or the established 2px amber rule; no lift/glow for static content. Always preserve visible keyboard focus.
- Reserve filled accent buttons for the main action in the current context. Use ghost/quiet text-like actions for secondary actions and destructive actions that should not dominate until hovered, focused, or confirmed.
- Preserve semantic status colors. Do not use green as brand success decoration or amber as an error substitute.
- Keep form labels visible above controls; use hints and notices sparingly. Group related controls with spacing before adding panels.

## Responsive behavior

- Check desktop and narrow mobile widths, especially 320–430 px. Measure `scrollWidth > innerWidth`, wrapping, clipping, focus visibility, and bottom-nav overlap.
- On phones, assume the desktop navbar is absent and the fixed bottom navigation owns primary destinations. Leave bottom clearance to `.page-content` and `--bottom-nav-clearance` unless the shell itself is the problem.
- Let action rows wrap or stack before they squeeze labels. Full-width buttons are appropriate on mobile when the action row would otherwise look cramped.
- Avoid fixed heights for content panels unless the interaction genuinely needs a bounded scroll area. Prefer natural document flow.

## Review checklist

- Confirm the change uses existing tokens/classes or adds a clearly named reusable token/class.
- Confirm color choices work in dark theme, light theme, `prefers-contrast: more`, and keyboard focus states.
- Confirm cards, borders, dividers, and empty-state copy are necessary and not decorative noise.
- Confirm the page keeps its intended character instead of copying another page's layout wholesale.
- Confirm rows follow the one row language, and that any deviation is one of the documented quieter treatments (menus, tables, chrome, `/more`).
- Confirm unrelated behavior, routes, data loading, posting, deletion, storage, and accessibility semantics are unchanged.
- For trainer-specific UI work, switch to `opencw-trainer-ui`; keep this skill focused on site-wide language and cross-page consistency.

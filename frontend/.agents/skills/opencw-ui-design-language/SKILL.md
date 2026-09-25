---
name: opencw-ui-design-language
description: "Use for OpenCW's shared visual language: design tokens, CSS primitives, layout shell, navigation treatment, typography, colour, spacing, borders, contrast, responsive rules, and visual consistency across the site's pages (home, forum, settings, profile, auth, about, more). It also applies to the trainer's visual work, where it supplies the shared principles while opencw-trainer-ui owns behaviour, states and trainer layout constraints — load both for a trainer visual change. Skip it for backend/API work, data-model changes, non-visual fixes, or work that is only trainer behaviour (opencw-trainer-ui alone)."
argument-hint: "Describe the OpenCW UI design, review, or visual consistency task"
---

# OpenCW UI Design Language

Apply OpenCW's existing visual language when making or reviewing UI anywhere in the app. Preserve page-specific product differences: consistency means shared principles and primitives, not making the trainer, forum, settings, auth, and marketing pages look identical.

On the trainer this skill and `opencw-trainer-ui` both apply. This skill owns the shared visual language — tokens, primitives, shell, cross-page consistency; the trainer skill owns what is specific to the trainer (modes, interaction contracts, keyboard behaviour, panel layout constraints, its verification states). A trainer visual change takes the shared principles from here and the trainer requirements from there; neither skill overrides the other outside its scope.

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
  - Trainer surfaces, only for cross-site consistency: `opencw-trainer-ui` names the files and owns their behaviour, states and layout constraints.

## Core language

- Use the warm Morse-workshop palette from `src/app.css`: dark/light theme variables, amber as the single product accent, and semantic status colors for correctness/warnings/errors. Do not hard-code new hex colors unless adding a token with a clear role.
- Maintain contrast roles: `--text-primary` for reading content, `--text-secondary` for supporting metadata, `--text-muted` for quiet labels. Use `--border-control` for form edges that need non-text contrast; use `--border`/`--border-subtle` for decorative hairlines.
- Keep typography simple: UI system font for interface/prose, monospace only for Morse text, timers, WPM/Hz, lesson numbers, call signs, tabular stats, and code-like values.
- Follow the 4px spacing scale (`--space-*`) and page rhythm tokens (`--section-gap`, `--block-gap`). Let `.page-content` own page top/bottom padding; do not add route-level padding to compensate for the navbar or phone tab bar.
- Prefer flat, mostly rectangular surfaces. Radii are small; shadows belong to menus, popovers, dialogs, and overlays, not static cards.

## Layout and hierarchy

- Choose width by content: default container for normal pages, `.page-narrow` for focused reading, forms, details, settings and thread pages (about and the Morse hub use it too). Add `.page-stack` when a page needs the block rhythm on its children without the narrow measure (profile).
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
- One row language: every whole-row link list — home destinations, the Morse hub, and the forum thread list — is a `.row-list` of `.row-link`s, each reading as one target: hairline separators, the shared row padding, and the floor + 2px amber edge on hover and keyboard focus. Do not give one list its own padding, radius or separator strategy, and keep destructive row actions quiet (unbordered, danger on hover) rather than filled buttons. The quieter treatments are deliberate and local, not drift: dropdown items, the profile history table, the desktop navbar, the bottom tab bar, and `/more` rows (colour-only hover).
- Labels share one recipe: `.card-label` (which keeps its bottom margin), `.panel-label` and `.metric-label` (margin-free) are the same caps label. Do not author a fourth.
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
- For trainer work, load `opencw-trainer-ui` as well: apply the checklist above to the shared visual side, and leave trainer behaviour, states and its own layout constraints to that skill.

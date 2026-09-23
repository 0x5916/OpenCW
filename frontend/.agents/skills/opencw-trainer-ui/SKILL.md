---
name: opencw-trainer-ui
description: 'Use when changing, fixing, reviewing, or verifying OpenCW''s Morse practice trainer (src/routes/morse/learn/+page.svelte): trainer UI, Character drill, Copy passage, keyboard shortcuts, checked-copy review, responsive or mobile layout, panel height, and controls that collide with the fixed bottom navigation. Also use when asked why the drill or passage panel behaves, wraps, or measures differently than expected.'
argument-hint: 'Describe the trainer change to make, or the states to verify'
---

# OpenCW Trainer UI

One route holds both practices: **Character drill** (listening to a selected character)
and **Copy passage** (transcribing a timed transmission). They share a lesson model,
playback component and mode switch, and they are deliberately different tools — most
mistakes here come from blurring them together.

## Inspect before editing (reflect the code, not assumptions)

- `src/routes/morse/learn/+page.svelte` — both modes plus their styles.
- `src/app.css` — design tokens and shared classes (buttons, fields, panels).
- `src/lib/styles/layout.css` — page shell, footer, fixed mobile bottom navigation and the
  clearance that keeps content above it.
- Direct components: `MorsePlayer` (playback + settings), `ResultOverlay` (review).
- `src/lib/storageKeys.ts` — `CW_STORAGE_KEYS.lesson`, `UI_STORAGE_KEYS.trainerMode`,
  `UI_STORAGE_KEYS.trainerIntroLesson`. Never rename these.
- `messages/*.json` — message keys; English is the source and the fallback for other
  locales. Paraglide only regenerates while the dev server runs.

Edit surface: the route and those direct components. Touch the shell only when the issue
genuinely belongs to it (footer or fixed-nav overlap), and never stretch the trainer card
to fix a page-level problem.

Contracts to keep unless the request explicitly changes them: lesson progression and
storage, the `MorsePlayer` handle (`playNow` / `stopNow` / `isStarted`) and its settings,
`saveProgressOfflineFirst` on checked copies, `generateTimedLesson` / `score` / `diffWords`,
keyboard shortcuts, and the established visual style.

## Mode rules

- **Character drill** is selected-character listening: character, Morse pattern, selector
  limited to the lesson's characters (defaulting to the newest one), Play/Replay, and one
  hand-over action into passage practice. It has no answer field, no checking and produces
  no results.
- **Copy passage** owns transcription: playback, answer field, state line, state-driven
  hints, Check transcription / Generate another passage, and the source passage rendered
  **only after the copy has been checked**, with review in the overlay.
- The mode switch is always available; the drill leads a lesson the learner has not opened
  yet, and the chosen mode is remembered.

## Procedure

1. **Inspect both modes** in the states that exist — drill: ready, playing, selected
   character changed; passage: ready, playing, typed, checked, review open — plus the
   lesson change path that resets a mode.
2. **Find the smallest change** that solves the reported issue. Prefer one rule or one
   condition over a restructure; leave unrelated layout, copy and behaviour alone.
3. **Keep the modes separate.** Do not give the drill answer checking, and do not weaken
   the passage's source-hidden-until-review rule.
4. **Layout discipline:** no fixed card heights or filler space; the panel ends after its
   final action row with the panel's normal padding. Watch for wrapping action rows on
   narrow screens (stack them full width when they wrap) and for anything hidden behind the
   fixed bottom navigation — that clearance belongs to the page shell, not the card.
5. **Delete, don't hide.** When behaviour changes, remove the obsolete markup, state,
   handlers, styles and message keys together so nothing dead is left behind.
6. **Verify, then report.** Run `npm run check`, `npm run lint` and `npm run build`. Then
   check the affected states in the browser (with `npm run dev` running — Paraglide only
   regenerates messages while it does) on desktop and at narrow mobile widths (about
   320–430 px), including the review state. Trust measured geometry (bounding boxes,
   `scrollWidth`/`clientWidth`, computed styles, document height) over screenshots, and
   confirm no clipping, no unwanted internal scroll area, and no overlap with the bottom
   navigation. Report what changed, what you verified, and what you could not verify.

## Checklist before finishing

- [ ] Both modes still work in their ready, active and review states.
- [ ] Audio plays, replays and stops as before; playback settings still apply.
- [ ] Lesson selection, default lesson and remembered mode are unchanged.
- [ ] Keyboard: Space starts/replays and Esc stops from outside controls, Enter checks a
      copy, and spaces typed in the transcription field are never intercepted.
- [ ] Drill shows the selected character and its pattern, and never asks for an answer.
- [ ] Passage keeps the source hidden until the copy is checked.
- [ ] The current next-step action is the primary one — Start session while ready, then
      Check transcription once an answer can be submitted (and it stays disabled until
      an answer exists); the drill's hand-over action stays secondary-looking until used.
- [ ] No fixed panel height, no filler content, no accidental mobile button wrapping, no
      control under the bottom navigation.
- [ ] No dead state, selectors or message keys; checks and build pass.

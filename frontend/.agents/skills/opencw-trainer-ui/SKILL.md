---
name: opencw-trainer-ui
description: 'Use when changing, fixing, reviewing, or verifying OpenCW''s Morse practice trainer (src/routes/morse/learn/+page.svelte): Character drill, Copy passage, playback, keyboard shortcuts, checked-copy review, lesson/mode state, panel layout constraints, and controls that collide with the fixed bottom navigation. Visual changes here also follow opencw-ui-design-language for shared tokens, primitives and cross-page consistency — load both when the work is visual, and this skill alone when it is trainer behaviour. Also use when asked why the drill or passage panel behaves, wraps, or measures differently than expected.'
argument-hint: 'Describe the trainer change to make, or the states to verify'
---

# OpenCW Trainer UI

One route holds both practices: **Character drill** (listening to a selected character)
and **Copy passage** (transcribing a timed transmission). They share a lesson model,
playback component and mode switch, and they are deliberately different tools — most
mistakes here come from blurring them together.

This skill and `opencw-ui-design-language` both apply to trainer visual work. The
design-language skill owns the shared visual language (tokens, primitives, shell,
cross-page consistency); this one owns what is trainer-specific — modes, interaction
contracts, keyboard behaviour, panel layout constraints, and the verification states
below. Neither overrides the other outside its scope.

## Inspect before editing (reflect the code, not assumptions)

- `src/routes/morse/learn/+page.svelte` — both modes plus their styles.
- `src/app.css` — design tokens and shared classes (buttons, fields, panels).
- `src/lib/styles/layout.css` — page shell, footer and the fixed mobile bottom navigation.
- Direct components, both in `src/lib/components/`: `MorsePlayer.svelte` (playback +
  settings) and `ResultOverlay.svelte` (review).
- `src/lib/storageKeys.ts` — `CW_STORAGE_KEYS.lesson`, `UI_STORAGE_KEYS.trainerMode`,
  `UI_STORAGE_KEYS.trainerIntroLesson`, `UI_STORAGE_KEYS.quickstartDismissed` (which decides
  whether the quickstart strip shows). Never rename these.
- `messages/*.json` — message keys; English is the source and the fallback for other
  locales. The generated Paraglide modules are gitignored and only exist after a dev run or
  a build, so regenerate (dev server or `npm run build`) before `npm run check`.

Edit surface: the route and those direct components. Touch the shell only when the issue
genuinely belongs to it (footer or fixed-nav overlap), and never stretch the trainer card
to fix a page-level problem.

Contracts to keep unless the request explicitly changes them: lesson progression and
storage, the `MorsePlayer` handle (`playNow` / `stopNow` / `isStarted`, typed once as
`PlayerHandle` in the route) and its settings, `saveProgressOfflineFirst` on checked
copies, `generateTimedLesson` / `score` / `diffWords`, `getLessonChars()` from `morse.ts`,
the shared `formatClock` / `percentage` helpers in `format.ts`, keyboard shortcuts, and
the shared visual style (`opencw-ui-design-language`).

## Mode rules

- **Character drill** is selected-character listening: character, Morse pattern, selector
  limited to the lesson's characters (defaulting to the newest one), Play/Replay, and one
  hand-over action into passage practice. It has no answer field, no checking and produces
  no results.
- **Copy passage** owns transcription: playback, answer field, state line, state-driven
  hints, Check transcription / Generate another passage, and the source passage rendered
  **only after the copy has been checked**, with review in the overlay. Starting playback —
  the player's Play button, Start session, a replay, a skip — puts the cursor in the answer
  field, and Escape there stops the transmission through the same handler the page-level
  shortcut uses.
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
   narrow screens (stack them full width when they wrap) and for content hidden behind the
   fixed bottom navigation — bottom-nav clearance is a shell rule, covered by
   `opencw-ui-design-language`.
   Small `min-height`s are not fixed heights: the live status line and the answer box keep
   theirs on purpose so the panel does not jump between states.
5. **Delete, don't hide.** When behaviour changes, remove the obsolete markup, state,
   handlers, styles and message keys together so nothing dead is left behind. Run
   `npm run messages:validate` after touching a message key: it fails on keys no source
   file references and on any locale that lags the English source.
6. **Verify, then report.** Run `npm run check`, `npm run lint` and `npm run build`. Then
   exercise the affected states in the browser — including the review state — at the widths,
   themes and focus states the design-language skill's responsive checks call for.
   Trainer-specific geometry to confirm: no clipping, no unwanted internal scroll area, and
   no control under the bottom navigation; trust measured geometry (bounding boxes,
   computed styles, document height) over screenshots. Report what changed, what you
   verified, and what you could not verify.

## Checklist before finishing

- [ ] Both modes still work in their ready, active and review states.
- [ ] Audio plays, replays and stops as before; playback settings still apply.
- [ ] Lesson selection, default lesson and remembered mode are unchanged.
- [ ] Keyboard: Space starts/replays from outside controls and Esc stops the transmission,
      whether focus is outside a control or in the answer field; Enter checks a copy, and
      keys typed in the transcription field (including Space) are never intercepted.
- [ ] Drill shows the selected character and its pattern, and never asks for an answer.
- [ ] Passage keeps the source hidden until the copy is checked.
- [ ] The current next-step action is the primary one — Start session while ready, then
      Check transcription once an answer can be submitted (and it stays disabled until
      an answer exists); the drill's hand-over action stays secondary-looking until used.
- [ ] No fixed panel height, no filler content, no accidental mobile button wrapping, no
      control under the bottom navigation.
- [ ] Labels use the shared caps recipe (`.panel-label`, the same class the review overlay
      uses), and the quickstart — when it is shown — keeps its numbered steps and bulleted
      tips (Tailwind preflight removes list markers, so the rules have to restore them).
- [ ] The shared visual checks from `opencw-ui-design-language` hold for the trainer
      surface too (themes, contrast, focus visibility, widths) — its checklist applies here.
- [ ] No dead state, selectors or message keys; checks, `messages:validate` and build pass.

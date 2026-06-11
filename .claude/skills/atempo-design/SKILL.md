---
name: atempo-design
description: Use whenever building any UI, mock, prototype, slide, or visual asset for Atempo — the warm espresso-and-teal flow-state studio for a solo builder running many AI-assisted projects. Covers production UI and throwaway prototypes alike. Provides the invariants, color/type/spacing tokens, fonts, the tempo-meter brand mark, nine reusable components, and the Cockpit UI kit. Invoke before writing any Atempo-facing markup or CSS so the output is on-brand.
user-invocable: true
---

# Designing for Atempo

Everything lives under `design/` (paths below are from the repo root). The system is
**plain HTML/CSS + a vendored React runtime, no build step** — that's the feature: a human
can render it in a browser and you can read it as source.

## Read selectively — spend tokens proportionally

Don't read the whole system every time. Pull in only what the task needs:

- **Any task →** `design/PRINCIPLES.md` (the six invariants + voice rules, ~60 lines). This
  is the floor; read it every time. For the full spec (visual foundations, iconography,
  manifest) read `design/README.md` — but only when you need that depth.
- **Building UI / components →** also `design/styles.css` (link it to inherit every token +
  the webfonts) and the relevant `design/components/**/<Name>.prompt.md` (each is a
  standalone usage card: snippet, props, the one rule that's easy to violate). The `.d.ts`
  next to it is the typed contract; the `.jsx` is the live implementation.
- **A new surface / screen / mock →** also `design/ui_kits/cockpit/index.html` as the
  reference composition (how the primitives assemble into real screens).
- **A brand / marketing asset →** also `design/guidelines/brand-tempo-meter.html` (the logo
  + the tempo meter, the one motif).

## The tokens & components (what you're composing with)

- `design/tokens/colors.css` — espresso ground, cream foreground ramp, the one teal accent,
  reserved amber, hairlines. Reach for the semantic aliases (`--canvas`, `--accent`,
  `--text-muted`, `--line-08`, `--tempo-cold`…).
- `design/tokens/typography.css` — the mono/human split: `--font-display` (Bricolage
  Grotesque), `--font-ui` (Hanken Grotesk), `--font-mono` (JetBrains Mono). All three have
  system-font fallbacks so output still reads correctly offline.
- `design/tokens/spacing.css`, `design/tokens/fonts.css` (Google Fonts CDN).
- `design/components/` — nine primitives: `core/` (Button, Badge, Tag, KeyCap), `status/`
  (StatusDot, TempoMeter, RuntimePill), `surfaces/` (Card, EyebrowLabel). One source of
  truth — `.jsx` + `.d.ts` + `.prompt.md` each, no duplicates.
- `design/vendor/` — pinned React + ReactDOM (production) + Babel standalone. To make a
  React mock render with no build step, link `design/styles.css`, load these three, then
  load the component `.jsx` files via `<script type="text/babel" src="...">`. They attach to
  `window.AtempoDS`. (Needs HTTP, not `file://`.) For a pure HTML/CSS mock you only need
  `styles.css` and the tokens.

## Non-negotiables (the fast self-check)

- **Warm, not server-black.** Espresso ground, one teal accent, amber only for
  "glance-at-this" signals.
- **The mono/human split is sacred** — JetBrains Mono for anything the machine touches
  (paths, ports, branches, commits, commands); a human face for everything warm.
- **Tempo is pace, not health** — in tempo / steady / cold. Cold decays on its own; never
  nudge or shame a cold project.
- **The recap is the moment** — where you left off + one concrete next beat + open threads,
  assembled from real signals.
- **Orchestration is shown, never operated** — no gauges, logs, uptime, or per-service
  controls. This is a re-entry surface.
- **No emoji.** Functional Unicode marks only (`▶ → ✓ ? · ●`). The brand mark is the 4-bar
  tempo meter.

## Verification loop (do this before declaring an artifact done)

1. Serve the folder: `python3 -m http.server 4400 -d design/` (or `design/serve.sh`). For a
   standalone mock, serve its directory the same way — the `text/babel` scripts need HTTP.
2. Render it and take a screenshot (headless Chrome works:
   `--headless=new --screenshot=out.png --virtual-time-budget=8000 <url>`).
3. Read the screenshot back and compare against the self-check above: no emoji; mono/human
   split intact; espresso + teal only (amber sparing); tempo reads as pace; no gauges. Fix
   and re-render until it passes.

## When invoked with no other guidance

Ask what they want to build and a couple of scoping questions, then act as an expert Atempo
designer. Copy assets out and produce static HTML the user can open for mocks/slides; for
production code, link the tokens and follow the rules here.

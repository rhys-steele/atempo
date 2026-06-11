# Atempo Design System

> Atempo is a flow‑state studio for a solo builder running many AI‑assisted projects at once. It abolishes the cold start. When you build a lot of things in parallel, the work survives but the *rhythm* doesn't — you switch away, come back weeks later, and the thread is gone, both in your head (what you were doing, the next move) and in the machine (what was running, which port). Atempo lets you **re‑enter** instead of restart. The name is the thesis: *a tempo* is the instruction to rejoin the music in rhythm, exactly where it was going.

This design system encodes the **Cockpit** direction — Atempo as a cockpit for one. The whole fleet's state is readable in a single glance, the launch ramp is under your thumb, and runtime is surfaced confidently — but it stays a **re‑entry surface, never a monitoring console.** No gauges, CPU, uptime, live logs, or per‑service controls. That's the foil (k9s, lazydocker, mission control), not this.

## Sources
- Built from the in‑product direction studies in this project: `The Cockpit.dc.html` (chosen), plus `The Sequencer.dc.html` and `The Studio.dc.html` (alternates), and `Atempo Directions.dc.html` (index).
- Original product screenshots provided by the builder (dark espresso desktop app, teal accent, "Welcome back" overview + `pricing-page` project detail).
- No external codebase or Figma was attached; the system is the source of truth for these visuals. If a real codebase exists, re‑attach it and this system should be reconciled against it.

---

## The invariants (true in every Atempo surface)
1. **Tempo is status as a feeling of pace** — *in tempo / steady / cold.* It decays on its own, so "cold" is physics, never an accusation. The studio sorts by it. Never nudge or shame a cold project; returning should feel free.
2. **The recap is the signature moment** — a "where you left off" + one concrete "next beat" + open threads (unanswered decisions). Functional, assembled from real signals (branch, last commit, changed files, running services) — not lyrical prose.
3. **Orchestration is done‑for‑you‑and‑shown, never operated.** One Run brings things up, confirms they're live, and gets out of the way. The launch ramp shows the ordered entrances and deconflicted ports as a *result*, not a control panel.
4. **A project is the thing you pick up.** Larger projects have multiple **parts**, each its own repo/subdir with its own runtime or none. Composition, not tags. A single‑part project stays visually simple — hierarchy appears only when there are real parts.
5. **The mono/human split is sacred** — monospace for anything the machine touches (paths, ports, branches, commits, commands), a human face for everything warm.
6. **Warm, not server‑black.** Espresso ground + a single teal accent.

---

## CONTENT FUNDAMENTALS

**Voice.** Calm, plain, second‑person. Atempo speaks *to* the builder ("You were wiring caption search"; "where you left off") and never about itself in marketing language. It states facts and hands over the next move. It never congratulates, never scolds, never hypes.

**Casing.** Sentence case everywhere in prose and buttons ("Resume with Claude", "Open app", "Pick up in code"). UPPERCASE is reserved for tracked‑out mono **eyebrow labels** (`WHERE YOU LEFT OFF`, `NEXT BEAT`, `OPEN THREADS`, `LAUNCH RAMP`, `PARTS · 3`). Project names are lowercase mono‑adjacent tokens (`atlas`, `pricing-page`, `recipe-box`) — they read like directory names because they are.

**I vs you.** Always "you". Never "I" or "we". Atempo is an instrument, not a character.

**Tempo language.** Status is written as pace, not health: *in tempo*, *steady*, *cold*. Time is concrete and relative: `2h`, `2d`, `6w`, `2mo`, `3mo`. Cold copy is matter‑of‑fact and inviting, never guilt‑laden: "Figma mocks only — a clean slate when ready", "Engine plays cleanly — pick up at save/export". The cold divider says it outright: *"COLD — NO RUSH, EVERYTHING KEPT WARM FOR RE‑ENTRY"* and the footer reminds: *"tempo decays on its own — cold is physics, not a verdict."*

**Recap copy.** Two registers in one sentence: a human clause for what you were doing, with **machine tokens inline** for the specifics. e.g. *"Wire `SearchList.swift` to `/v2/captions/search`, drop the mock."* The "next beat" is always one concrete, do‑able action — a verb, a file, a destination.

**Machine copy.** Terse, lowercase, exact: `feat/captions`, `4 changed · 2h`, `:4100 :1420`, `$ atempo run atlas`, `postgres:16 ✓ up`. Commands are shown as the literal thing Atempo runs for you.

**Emoji.** None. The only glyphs used as marks are functional: the play triangle `▶ / ▸`, the right arrow `→`, the open‑question mark `?`, the check `✓`, the bullet `·`, and status dots `●`. (See ICONOGRAPHY.)

---

## VISUAL FOUNDATIONS

**Overall vibe.** A warm, low‑lit desktop app rendered as a floating macOS window on a near‑black espresso canvas. Dense but unhurried. It reads like a well‑set terminal that learned typography.

**Color.** Espresso ground (`--espresso-1000 #16110d` canvas, `--espresso-900 #1a1410` window) with a full cream foreground ramp (`#f2e9dd → #574c3f`, eight steps) doing nearly all the hierarchy work. A **single teal accent** (`--teal-500 #56b3a3`) carries every active meaning: in‑tempo, live, primary action, the playhead. A reserved **amber** (`--amber-400 #e0a458`) appears only for "give this a glance" signals — uncommitted changes, self‑notes, open questions. Saturation is low and warm throughout; whites are creamed, blacks are browned. No blues, no purples, no gradients-as-decoration.

**Type.** Three families, strict split. **Bricolage Grotesque** (display) for headings, project names, and big numbers. **Hanken Grotesk** (UI) for body, descriptions, prose. **JetBrains Mono** for everything the machine touches. Display runs 600–700 weight, tight leading (1.05). Mono labels are 11px, uppercase, `letter-spacing: 0.16–0.20em`.

**Backgrounds.** Flat warm fills only. No imagery, no full‑bleed photos, no patterns, no gradients except one barely‑there diagonal teal wash on the Studio hero (`linear-gradient(135deg, rgba(86,179,163,0.05), rgba(86,179,163,0.01))`). Depth comes from layered espresso surfaces + hairlines, not shadow.

**Cards & surfaces.** A card is a flat espresso surface with a **hairline border** (`rgba(242,233,221,0.04–0.12)` depending on emphasis) and a radius of 12–14px (20px for heroes). In‑tempo / active surfaces gain a teal‑tinted fill (`--teal-tint-04`) and a teal hairline (`--teal-line-20`). Cold rows fade their border toward `--line-04` and drop text down the cream ramp. **No drop shadows on in‑canvas cards** — shadow is reserved for the floating window itself (`0 32px 80px rgba(0,0,0,0.55)`).

**Borders & hairlines.** The system separates with 1px hairlines on the cream ramp at low alpha, not with fills. Dashed hairlines mean "intentionally inert" — a bypassed part, a skipped ramp step (`1px dashed rgba(242,233,221,0.12)`).

**Status dots & the tempo meter.** A 5–6px round dot, teal for live, animated with the single `pulseDot` keyframe (a slow 2.2s breath, opacity 1→0.3→1). Tempo itself renders as a **4‑bar meter** (like an equalizer): tall+teal = in tempo, mid+tan = steady, short+faint‑brown = cold. The bars descend left‑to‑right.

**Hover / press.** Hover lifts a surface one step (teal‑tint‑04 → teal‑tint‑07, or transparent → `rgba(242,233,221,0.02)`); primary buttons lighten teal toward `--teal-400`. Press has no scale/bounce — at most a slight darken. Motion is restrained: fades and the live pulse only. Nothing bounces; cold is reached by decay.

**Layout.** Content sits in fixed‑width window frames (1440px for the app) centered on the canvas with `min-width: fit-content` so narrow viewports scroll rather than clip. Inside, dense screens use CSS‑grid rows with explicit column templates and generous `gap`; restful screens use wide padding (up to 120px) and large type. Window chrome is always present: three traffic lights, a mono breadcrumb (`atempo / fleet / atlas`), and a right‑aligned live/tempo readout.

**Transparency & blur.** Used sparingly — teal and cream tints over the espresso ground for emphasis. No backdrop blur; the app is opaque and grounded.

---

## ICONOGRAPHY

Atempo is **near‑iconless by design** — hierarchy is carried by type, color, and the tempo meter, not an icon set. The few marks in use are **functional Unicode glyphs**, not an icon font or SVG library:

- `▶` / `▸` — run / play (primary orchestration action)
- `→` — re‑enter / open / proceed
- `←` — back (breadcrumb)
- `?` — an open thread / unanswered decision (rendered in amber)
- `✓` — confirmed live / done
- `●` — status dot (teal = live, amber = attention, cream = neutral)
- `·` — the separator between machine facts
- `⌕` — search affordance in the command bar
- `+` — new project
- `⌘K` — the command palette shortcut, shown in a hairline key‑cap

The **logo mark** is a 4‑bar equalizer (the same motif as the tempo meter) — bars of descending height in teal, doubling as the "in tempo" glyph. No emoji anywhere. The macOS traffic lights (`#ff5f57 / #febc2e / #28c840`) appear only in window chrome and are decorative, not interactive.

If you need a broader icon set for a future surface, substitute **Lucide** (matching the thin, geometric, low‑ornament feel) and flag the addition — but prefer adding none.

---

## INDEX / MANIFEST

- **`styles.css`** — root entry; `@import`s all tokens. Link this.
- **`tokens/`** — `colors.css`, `typography.css`, `spacing.css`, `fonts.css`.
- **`guidelines/`** — foundation specimen cards (Type, Colors, Spacing, Brand) shown in the Design System tab.
- **`components/`** — reusable primitives: `core/` (Button, Badge, Tag, KeyCap), `status/` (StatusDot, TempoMeter, RuntimePill), `surfaces/` (Card, EyebrowLabel).
- **`ui_kits/cockpit/`** — the Cockpit product kit: Fleet overview + Atlas project briefing, click‑through.
- **`SKILL.md`** — portable skill manifest for use in Claude Code.

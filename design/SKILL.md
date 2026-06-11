---
name: atempo-design
description: Use this skill to generate well-branded interfaces and assets for Atempo — a warm, espresso-and-teal flow-state studio for a solo builder running many AI-assisted projects. Use for production UI or throwaway prototypes/mocks. Contains design guidelines, color/type/spacing tokens, fonts, the tempo-meter brand mark, reusable components, and the Cockpit UI kit.
user-invocable: true
---

Read `readme.md` in this skill first — it carries the product context, the six invariants, content fundamentals, visual foundations, and iconography. Then explore:

- `styles.css` + `tokens/` — link `styles.css` to inherit every color/type/spacing token and the webfonts (Bricolage Grotesque, Hanken Grotesk, JetBrains Mono).
- `components/` — reusable primitives (Button, Badge, Tag, KeyCap, TempoMeter, StatusDot, RuntimePill, Card, EyebrowLabel). Each has a `.prompt.md` with usage.
- `ui_kits/cockpit/` — the full product surface (fleet overview + project briefing).
- `guidelines/` — foundation specimen cards.

Non-negotiables when designing for Atempo:
- **Warm, not server-black.** Espresso ground, one teal accent, a reserved amber for "glance at this" signals only.
- **The mono/human split is sacred** — JetBrains Mono for anything the machine touches (paths, ports, branches, commits, commands); a human face for everything warm.
- **Tempo is pace, not health** — in tempo / steady / cold. Cold decays on its own; never nudge or shame a cold project.
- **The recap is the moment** — where you left off + one concrete next beat + open threads, assembled from real signals.
- **Orchestration is shown, never operated** — no gauges, logs, uptime, or per-service controls. This is a re-entry surface.
- **No emoji.** Functional Unicode marks only (▶ → ✓ ? · ●). The logo/brand mark is the 4-bar tempo meter.

If creating visual artifacts (slides, mocks, throwaway prototypes), copy assets out and produce static HTML for the user to view. If working in production code, copy the tokens and read the rules here. If invoked with no other guidance, ask what they want to build, ask a few questions, and act as an expert Atempo designer outputting HTML or production code.

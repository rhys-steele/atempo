# Atempo — design principles (the rules, extracted)

The cheap-to-read core of the system: the invariants + the voice. Read this on **any**
Atempo design task. For the full spec (visual foundations, iconography, manifest) read
`README.md`. These two together are the system.

> Atempo is a flow-state studio for a solo builder running many AI-assisted projects at
> once. It abolishes the cold start. The whole value is **re-entry**: come back after weeks
> and immediately know a project's state, history, and next move. The name is the thesis —
> *a tempo* is the instruction to rejoin the music in rhythm, exactly where it was going.
> The chosen direction is the **Cockpit**: a cockpit for one — a re-entry surface, never a
> monitoring console.

## The six invariants (true in every Atempo surface)

1. **Tempo is status as a feeling of pace** — *in tempo / steady / cold.* It decays on its
   own, so "cold" is physics, never an accusation. The studio sorts by it. Never nudge or
   shame a cold project; returning should feel free.
2. **The recap is the signature moment** — a "where you left off" + one concrete "next beat"
   + open threads (unanswered decisions). Functional, assembled from real signals (branch,
   last commit, changed files, running services) — not lyrical prose.
3. **Orchestration is done-for-you-and-shown, never operated.** One Run brings things up,
   confirms they're live, and gets out of the way. The launch ramp shows the ordered
   entrances and deconflicted ports as a *result*, not a control panel.
4. **A project is the thing you pick up.** Larger projects have multiple **parts**, each its
   own repo/subdir with its own runtime or none. Composition, not tags. A single-part
   project stays visually simple — hierarchy appears only when there are real parts.
5. **The mono/human split is sacred** — monospace for anything the machine touches (paths,
   ports, branches, commits, commands), a human face for everything warm.
6. **Warm, not server-black.** Espresso ground + a single teal accent. (Amber is reserved
   for "give this a glance" signals only — uncommitted changes, self-notes, open questions.)

## Voice & content rules

- **Voice.** Calm, plain, second-person. Atempo speaks *to* the builder ("You were wiring
  caption search"). It states facts and hands over the next move. Never congratulates, never
  scolds, never hypes. It's an instrument, not a character.
- **I vs you.** Always "you". Never "I" or "we".
- **Casing.** Sentence case everywhere in prose and buttons ("Resume with Claude", "Open
  app"). UPPERCASE is reserved for tracked-out mono **eyebrow labels** (`WHERE YOU LEFT OFF`,
  `NEXT BEAT`, `LAUNCH RAMP`, `PARTS · 3`). Project names are lowercase mono-adjacent tokens
  (`atlas`, `pricing-page`) — they read like directory names because they are.
- **Tempo language.** Status is written as pace, not health: *in tempo*, *steady*, *cold*.
  Time is concrete and relative: `2h`, `2d`, `6w`, `2mo`. Cold copy is matter-of-fact and
  inviting, never guilt-laden.
- **Recap copy.** Two registers in one sentence: a human clause for what you were doing, with
  **machine tokens inline** for the specifics — e.g. *"Wire `SearchList.swift` to
  `/v2/captions/search`, drop the mock."* The next beat is always one concrete, do-able
  action: a verb, a file, a destination.
- **Machine copy.** Terse, lowercase, exact: `feat/captions`, `4 changed · 2h`, `:4100 :1420`,
  `$ atempo run atlas`, `postgres:16 ✓ up`.
- **No emoji, ever.** The only marks are functional Unicode glyphs: `▶ / ▸` (run), `→` (open /
  proceed), `←` (back), `?` (open thread, amber), `✓` (confirmed), `●` (status dot), `·`
  (separator), `⌕` (search), `+` (new), `⌘K` (palette, in a hairline KeyCap). The logo/brand
  mark is the 4-bar tempo meter.

## The fast self-check (does it pass?)

- No emoji; only functional Unicode marks.
- Mono for everything the machine touches; a human face for everything warm.
- Espresso ground + one teal accent; amber only for glance-signals. No blues/purples, no
  decorative gradients, no drop shadows on in-canvas cards.
- Tempo reads as pace (in tempo / steady / cold), never as health/uptime.
- No gauges, CPU, live logs, or per-service controls — this is re-entry, not monitoring.

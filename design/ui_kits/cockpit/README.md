# Cockpit UI Kit

The chosen Atempo direction — **a cockpit for one.** The whole fleet's state readable in a single glance, the launch ramp under your thumb, runtime surfaced confidently — yet still a **re-entry surface, never a monitoring console.**

## Screens (`index.html`, click-through)
1. **Fleet overview** — every project as a dense row sorted by tempo (in tempo → steady → cold). Each row carries a `TempoMeter`, the project name, its one **next beat**, git state, read-only runtime, and a single go-action. Click **Open →** on `atlas` to enter the detail.
2. **Atlas briefing** — a multi-part project (api + desktop + site). The cross-part **recap trio** (where you left off · next beat · open threads), the **launch ramp** (ordered entrances, deconflicted ports, shown as a result), and the **parts** strip with per-part git + runtime. **← Fleet** returns.

## Composition
Built entirely from the system's primitives — `Button`, `Badge`, `Tag`, `KeyCap`, `TempoMeter`, `StatusDot`, `RuntimePill`, `Card`, `EyebrowLabel` — resolved from the generated DS bundle. Screens are defined inline in `index.html` for standalone reliability rather than split into separate bundled files (a screen `.jsx` with no `.d.ts` would otherwise be swept into the component namespace).

## What it deliberately omits
No CPU/memory gauges, uptime, live log tails, or per-service start/stop. Runtime is reported, not operated. That restraint is the product.


# Atempo

A flow-state studio for solo builders running many AI-assisted projects at once.

Ship a project in a weekend with Claude. Come back three weeks later and pick up *in tempo* — not from a cold start. Atempo holds everything Claude can't see above the project: what you have, what's running, where you left off, and how to get back into it in one move.

Available as a **CLI** (drive it yourself) and a **desktop app** (it drives itself and shows you it's handled).

```bash
atempo create laravel:12 my-app
atempo start my-app
```

---

## The idea

When you build a lot of things at once, the work survives but the rhythm doesn't. You switch away, come back weeks later, and the thread is gone — in two places at once:

* **In your head** — what you were mid-thought on, why that file was open, what the next move was.
* **In the machine** — what was running, which env, which port, what the run command even was.

Both of them stall you on the doorstep. The expensive part of coming back isn't the code; it's rebuilding that context. Atempo's whole job is to make re-entry cost almost nothing — and to protect flow once you're back in.

The name is the thesis. *A tempo* is the musical instruction to return to the original tempo after you've slowed or drifted: you don't restart the piece, you rejoin it, in rhythm, exactly where it was going.

**A principle worth stating plainly:** orchestration is a cost, not a feature. Spinning up a stack, deconflicting ports across a fleet, remembering what runs where — people do this because they're forced to, not because they enjoy it. Atempo does that heavy lifting underneath so the studio stays calm. The engine eats orchestration so flow survives.

---

## How Claude fits

Claude is excellent *inside* a project — it reads your code, understands your patterns, helps you move fast. But it has no awareness *above* the project: it doesn't know what else you're building, what was running last week, or where you left off when you switched away. And inside a session, it can't see your *running* environment — only your code.

Atempo closes both gaps. It keeps the cross-project context Claude lacks, and via MCP it feeds Claude live runtime state — which services are up, the actual port assignments, recent errors, environment state. Claude goes from reading your code to understanding your running environment. **This runtime bridge is the thing nothing else does.**

---

## One engine, three faces

All orchestration logic lives in a single Go core. Everything else is a thin face over it:

* **The CLI** is the cockpit — full manual control, for people who want to live in the terminal.
* **The desktop app** is the calm studio — fleet-at-a-glance and visual re-entry. Here, orchestration is  *done-for-you-and-shown* , never  *operated* .
* **The MCP server** is the bridge — it hands Claude the live runtime context it otherwise can't see.

The split isn't cosmetic. It's the product: the CLI is where you drive, the desktop is the doorway, and the same engine powers both — plus Claude.

---

## What's in this repo

Atempo is a monorepo. Each top-level directory is one face of the product:

```
cli/        # The Go engine + CLI. All orchestration logic lives here:
            # scaffolding, the project registry, docker-compose generation,
            # port assignment, and the embedded framework templates.

desktop/    # The Tauri desktop app — Rust backend + React/Vite/Tailwind
            # frontend. It never reimplements CLI logic; it runs the Go
            # binary as a sidecar and talks to it over a structured JSON
            # protocol (--json for one-shot calls, --stream-json for
            # streaming progress and logs).

design/     # The design system — tokens, components, guidelines, and a
            # full UI kit as plain HTML/CSS. The visual source of truth
            # for every Atempo surface, readable by humans and by Claude.
            # See "Design system" below.
```

Planned siblings: a marketing/docs site and, later, a web app — both will draw on `design/` directly.

---

## What it does

Atempo manages three moments of flow.

### Getting back in

Every project carries a **recap** — assembled from real signals, not prose you have to write: current branch, last commit and when, changed files ("you were in…"), what's running, and any **open threads** you jotted (the unanswered decisions — *"session store: redis or a signed cookie?"* — that are the real reason re-entry hurts). Where available, it surfaces a  **next beat** : the smallest concrete next action.

> The recap is functional first. A one-line narrative summary ("you were wiring the checkout flow") is a later premium layer that sits *above* the signals — the card is built to receive it, but it isn't required for the recap to be useful today.

Each project also ships an `ai/context.md`, pre-populated with stack, structure, and patterns — something to hand straight to Claude instead of starting cold.

### Staying in

Orchestration never interrupts you. One move brings the stack up; Atempo confirms it's live and gets out of the way. And because Claude has your runtime over MCP, you never stop mid-flow to narrate your environment — it already knows.

### Switching across the fleet

The desktop studio is a live view of everything you're building — what's running, what's dormant, last touched — so moving between projects doesn't mean losing your place in either.

---

## Tempo — status as a feeling of pace

Projects don't have a server "status," they have a  **tempo** :

* **in tempo** — actively in flow
* **steady** — shipped and stable
* **cold** — drifted

Tempo decays on its own from last activity, so *cold* is physics, not an accusation — nothing failed, the project just cooled. The studio sorts by it so the thing most worth resuming is the first thing you see. And **letting a project rest** is a first-class action, distinct from archiving: permission to pause without guilt, so returns happen from desire rather than obligation.

> Tempo decay, the recap card, open threads, and "let it rest" are the product's organizing model. Some are directional relative to the current CLI — treat them as the target the project is moving toward.

---

## What Atempo is *not*

Not a monitoring console. The foil is k9s / lazydocker / mission-control — dashboards you live inside, watching CPU and uptime. Atempo doesn't make you *operate* your infrastructure; it runs it and shows you it's handled, then steps back.

That line splits cleanly across the two surfaces:

* **In the desktop app:** Run, Open site, Open in editor, and a confirmation that the stack came up (`app up on :4321 ✓`). It shows the commands it ran underneath as transparency and an escape hatch — *here's what I ran for you* — not as a control panel.
* **In the CLI:** the full power — `logs`, `status`, per-service control, everything. Monitoring lives here, where it belongs, for people who want it.

---

## A project is a registered thing, not just a stack

A project is something you've registered with Atempo: a path, a tempo, a recap — and an  **optional runtime adapter** . Docker Compose is the primary adapter (and the only one today), but the model is deliberately looser than "a project is a compose stack," so that ideas, specs, and writing-only projects are first-class members of the portfolio. A markdown-only project has no launch ramp — it gets "open the doc" instead — but it still has a tempo, a recap, and a place in the studio.

> Current reality: the CLI scaffolds and runs Dockerized projects. First-class markdown/non-container projects are part of the intended model and the data layer should be built to accommodate them — a project should not be required to have a runtime.

---

## Installation

### CLI

Requires Go 1.23+:

```bash
go install github.com/rhys-atempo/atempo@latest
```

Or build from source:

```bash
git clone https://github.com/rhys-atempo/atempo
cd atempo/cli
go build -o atempo ./cmd/atempo
mv atempo /usr/local/bin/
```

### Desktop app

Build the CLI sidecar, then launch the Tauri app:

```bash
cd cli
go build -o ../desktop/backend/binaries/atempo-aarch64-apple-darwin ./cmd/atempo
cd ../desktop/backend
cargo tauri dev
```

**Prerequisites:** Docker with the Compose v2 plugin, Git.

Run `atempo doctor` to verify your environment.

---

## Usage

### Create a project

```bash
atempo create laravel:12 my-app
atempo create fastapi my-api
atempo create bun my-site
atempo create python my-worker
```

This produces:

```
my-app/
├── src/                    # Framework source code
├── ai/                     # AI context files for Claude
│   └── context.md
├── infra/
│   └── docker/
│       ├── docker-compose.yml
│       └── Dockerfile
└── atempo.json             # Project metadata + port assignments
```

Ports are assigned automatically — no conflicts across your fleet.

### Manage projects

```bash
atempo start [project]         # docker compose up -d
atempo stop [project]          # docker compose down
atempo logs [project]          # stream container logs
atempo status [project]        # show service health
atempo projects                # list all registered projects
atempo describe [project]      # show project details
atempo add-service <service>   # add minio, elasticsearch, rabbitmq, mongodb
atempo reconfigure             # regenerate docker-compose.yml
atempo remove <project>        # remove from registry
atempo doctor                  # check system dependencies
```

### Interactive shell

Running `atempo` with no arguments starts an interactive shell with a git-aware prompt, command history (↑/↓), and tab completion for both atempo commands and filesystem paths.

---

## Supported frameworks

| Framework   | Version syntax             | Stack                                          |
| ----------- | -------------------------- | ---------------------------------------------- |
| `laravel` | `laravel`,`laravel:12` | PHP-FPM, Nginx, MySQL, Redis, Mailhog          |
| `fastapi` | `fastapi`                | Python 3.12, Uvicorn, PostgreSQL, Redis        |
| `bun`     | `bun`                    | Bun runtime, TypeScript,`Bun.serve()`        |
| `python`  | `python`                 | Python 3.12, plain container (scripts/workers) |

Heavy stacks aren't meant to all run at once — and that's the point. Cold projects are stopped; only what's in tempo runs. The architecture and the tempo model agree.

---

## Predefined services

Add to any existing project with `atempo add-service <name>`:

| Service           | Image                 | Ports       |
| ----------------- | --------------------- | ----------- |
| `minio`         | minio/minio           | 9000, 9001  |
| `elasticsearch` | elasticsearch:8.8.0   | 9200        |
| `rabbitmq`      | rabbitmq:3-management | 5672, 15672 |
| `mongodb`       | mongo:6               | 27017       |

---

## Design system

The visual identity lives in `design/` as a self-contained design system — plain HTML and CSS, no build step. It's the source of truth for every Atempo surface: the desktop app today, the website and web app later. Because it's static files, it serves two audiences at once: a human can open any file in a browser and see rendered specimens, and Claude can read the same files as code and design new surfaces that stay on-brand.

```
design/
├── index.html            # The system viewer/gallery — serve, then open this
├── PRINCIPLES.md         # The cheap-to-read rules: invariants + voice
├── README.md             # The full system spec: content + visual foundations,
│                         # iconography, manifest
├── styles.css            # Root entry — @imports all tokens. Link this and
│                         # you inherit the whole system.
├── tokens/               # colors.css, typography.css, spacing.css, fonts.css
├── guidelines/           # Foundation specimen cards: type, color ramps,
│                         # spacing scale, radii, the tempo-meter brand mark
├── components/           # 9 primitives, each as .jsx + .d.ts + a
│   ├── core/             #   .prompt.md usage guide: Button, Badge, Tag, KeyCap
│   ├── status/           #   StatusDot, TempoMeter, RuntimePill
│   └── surfaces/         #   Card, EyebrowLabel
├── ui_kits/cockpit/      # The full product kit — fleet overview + project
│                         # briefing as a click-through, built from the primitives
├── vendor/               # Pinned React + ReactDOM + Babel — renders offline,
│                         # no build step
├── archive/              # Direction studies (*.dc.html), the generated
│                         # support.js, and the source screenshots — history
└── SKILL.md              # Pointer to the .claude/skills/atempo-design/ skill
```

### The design language in brief

* **Warm, not server-black.** An espresso ground (`#16110d`) with a cream foreground ramp doing the hierarchy work. One teal accent (`#56b3a3`) carries every active meaning — in tempo, live, primary action. A reserved amber appears only for "give this a glance" signals: uncommitted changes, open questions.
* **The mono/human split is sacred.** **Bricolage Grotesque** for display, **Hanken Grotesk** for UI prose, and **JetBrains Mono** strictly for anything the machine touches — paths, ports, branches, commits, commands, the tempo mark. The human parts feel human; the machine parts look like machine; the line between them is intentional.
* **Tempo is pace, not health.** Status renders as a **4-bar meter** — tall and teal for in tempo, mid tan for steady, short and faint for cold. The same 4-bar mark is the logo. Cold copy is matter-of-fact and inviting, never guilt-laden.
* **Near-iconless.** No icon font, no emoji — a handful of functional Unicode glyphs (`▶ → ✓ ? · ●`) and type doing the work.
* **Copy is plain, sentence case, warm.** Always "you", never "we". Empty and idle states are invitations, not apologies. Load-bearing phrases: *where you left off*, *next beat*, *pick up*, *return to tempo*.

### Using it

To view it, serve the folder — `python3 -m http.server 4400 -d design/` (or `design/serve.sh`) — and open `localhost:4400` for the gallery; `ui_kits/cockpit/index.html` is the hero. Serving (not opening over `file://`) is required because the component cards and the kit render the real components through Babel's `text/babel` scripts. To build with it, link `styles.css` for the full token set and webfonts, and read `design/PRINCIPLES.md` for the invariants and voice (or `design/README.md` for the full spec). When designing with Claude, the `atempo-design` skill (`.claude/skills/atempo-design/`) wires the whole system in — invoke `/atempo-design` and Claude designs as an Atempo designer rather than from scratch.

---

## Architecture notes

* **One Go engine** holds all orchestration; the desktop runs it as a sidecar and the MCP server exposes it to Claude. Keep the desktop↔CLI contract **structured** (JSON / a local RPC), not stdout-scraped — it's the load-bearing seam.
* **Three places hold state** — `atempo.json`, the global registry, and actual Docker reality. Decide which is authoritative and  **reconcile against Docker as ground truth** ; don't trust the ledger for "is it running" or port allocation.
* **Generated config needs an owner-editable escape hatch** (e.g. a `docker-compose.override.yml`) so `reconfigure` never clobbers a user's hand edits.
* **The MCP runtime bridge** is the moat. Ports / services / env are solid and shippable; "recent errors" is honestly crude for now (tail of stderr), not structured error capture — keep it labelled as such.

---

## License

[Atempo License v1.0](LICENCE) — open source with attribution.

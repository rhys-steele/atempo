# Atempo — Claude Code Context

## Product vision

Atempo is a portfolio manager for AI-assisted projects. The core problem: Claude is excellent inside a project but has no awareness above the project level. When you're shipping multiple projects, the between-session context — what's running, where you left off, what you were trying to solve — lives nowhere and evaporates.

The core value is **re-entry**: come back to something after three weeks and immediately know its state, history, and what to pick up next. Every feature should serve this.

- The dashboard is a portfolio, not a container manager
- `ai/context.md` is a breadcrumb trail, not documentation
- Start/stop is "wake this project up", not DevOps
- MCP context injection means Claude knows your running environment, not just your code

Atempo has two components that work together:

- **CLI** (`cli/`) — Go binary that scaffolds and manages Docker-based projects
- **Desktop** (`desktop/`) — Tauri app (Rust backend + React frontend) that drives the CLI as a sidecar

## Repo layout

```
cli/                           # Go CLI
  cmd/atempo/
    main.go                    # entry point + embed directive
    templates/                 # embedded at build time (all:templates)
      frameworks/
        laravel/               # PHP-FPM + Nginx + MySQL + Redis + Mailhog
        fastapi/               # Python 3.12 + Uvicorn + PostgreSQL + Redis
        bun/                   # Bun runtime + TypeScript
        python/                # Python 3.12 plain container (scripts/workers)
  internal/
    app/commands/              # one file per CLI command
      command.go               # Command interface, BaseCommand, CommandContext, ANSI consts
      shell.go                 # interactive shell (readline, tab completion, history)
      create.go                # atempo create <framework[:version]> <name>
      start/stop/logs/status/  # thin wrappers around docker compose
      add_service.go           # append predefined service + regenerate compose
      reconfigure.go           # regenerate docker-compose.yml from atempo.json
      describe/projects/remove/doctor/registry.go
    compose/
      generator.go             # LoadConfig, SaveConfig, GenerateDockerCompose, AddPredefinedService
      ports.go                 # port scanning + conflict-free assignment
    scaffold/
      scaffold.go              # Run() — installer, file copy, token expansion, registration
    registry/
      types.go                 # Project struct
    repositories/
      registry.go              # FileRegistryRepository (JSON file, ~/.config/atempo/registry.json)
      cache.go                 # MemoryCacheRepository (sync.Mutex, TTL-based)
    interfaces/
      registry.go              # RegistryRepository interface
      cache.go                 # CacheRepository interface
    docker/
      docker.go                # exec helpers: RunCompose, StreamCompose, RunDockerCommand
    dependency/
      manager.go               # doctor checks (docker, git, compose v2)
    errors/
      types.go                 # AtempoError with Code, Hint, Details
      utils.go                 # FormatUserMessage, GetRecoveryAction, GetDetails
    logger/
      logger.go                # structured logger

desktop/
  backend/                     # Tauri (Rust)
    src/
      lib.rs                   # plugin setup + invoke_handler registration
      sidecar.rs               # run_json() and stream_json() helpers
      commands/
        projects.rs            # list_projects, describe_project, remove_project
        compose.rs             # start_project, stop_project, get_status, reconfigure_project
        create.rs              # create_project (streaming)
        logs.rs                # stream_logs
        doctor.rs              # run_doctor
    binaries/                  # sidecar binary lives here (gitignored, built manually)
      atempo-aarch64-apple-darwin
    capabilities/
      default.json             # Tauri v2 permission declarations
    tauri.conf.json
  frontend/                    # React + Vite + Tailwind
    src/
      lib/tauri.ts             # invoke() wrappers + shared TypeScript types
      hooks/useProjects.ts     # React Query hooks
      pages/                   # Dashboard, CreateProject, ProjectDetail, Doctor
      components/              # Layout, Sidebar, ProjectCard, ServiceBadge, LogStream

design/                        # Design system — plain HTML/CSS, no build step
  index.html                   # the system viewer/gallery (serve, then open this)
  PRINCIPLES.md                # the invariants + voice (cheap-to-read rules)
  README.md                    # the full system spec + manifest
  styles.css                   # links all tokens/
  tokens/                      # colors, typography, spacing, fonts (Google Fonts CDN)
  guidelines/                  # foundation specimen cards (type, color, spacing, brand)
  components/                  # 9 primitives — one source of truth (.jsx + .d.ts + .prompt.md)
  ui_kits/cockpit/             # the chosen product kit (fleet → atlas, click-through)
  vendor/                      # pinned React + ReactDOM + Babel (offline, no build)
  archive/                     # original direction studies + support.js + screenshots
```

The design system is consumed by the `atempo-design` Claude Code skill
(`.claude/skills/atempo-design/`) — invoke it before building any Atempo UI or visual asset.
Serve it with `python3 -m http.server 4400 -d design/` (or `design/serve.sh`); it needs HTTP,
not `file://`, because the component cards and kit use Babel's `text/babel` scripts.

---

## Build & run

### CLI

```bash
cd cli
go build -o atempo ./cmd/atempo   # produces ./atempo binary
go test ./...
go test -race ./...               # used in CI
go vet ./...
```

### Desktop (dev)

1. Build the CLI sidecar first:
   ```bash
   cd cli
   go build -o ../desktop/backend/binaries/atempo-aarch64-apple-darwin ./cmd/atempo
   ```

2. Start the Tauri dev server (launches Vite automatically via `beforeDevCommand`):
   ```bash
   cd desktop/backend
   cargo tauri dev
   ```

3. After Go changes, rebuild the sidecar (step 1) — no restart needed.

CI runs on push to `main`/`develop` and PRs to `main` (`.github/workflows/build.yml`).

---

## CLI entry point

`cli/cmd/atempo/main.go` — embeds `all:templates` via `embed.FS`, wires the repository layer, builds the command registry, and dispatches. No arguments → interactive shell.

---

## CLI key conventions

**Templates are embedded, never copied as loose files.** All framework templates live under `cli/cmd/atempo/templates/` and are embedded via `//go:embed all:templates`. The `all:` prefix includes dotfiles (e.g. `.env.example`).

**Path prefix stripping.** The `embed.FS` is rooted at `cli/cmd/atempo/` so template paths start with `templates/`. Before passing to `scaffold.Run`, `create.go` calls `fs.Sub(templatesFS, "templates")` so scaffold sees `frameworks/laravel/...` directly.

**atempo.json is the source of truth.** Every project has an `atempo.json` that drives both docker-compose generation and the registry. `compose.GenerateDockerCompose(projectPath)` reads it, assigns free ports, rewrites it with assigned ports, and writes `infra/docker/docker-compose.yml`.

**Port assignment is automatic.** `compose/ports.go` scans all registered projects' atempo.json files and finds the next free port above each service's preferred base port. No manual port configuration needed.

**Installer types.** Framework `atempo.json` templates declare one of:
- `"type": "docker"` — runs the installer command inside Docker (e.g. `composer create-project`)
- `"type": "template"` — skip the installer; source files are included in the embedded template

**Token expansion.** Template files use `{{project}}`, `{{version}}`, `{{majorVersion}}` placeholders. `scaffold.expandTokens` replaces them at copy time. `{{majorVersion}}` of `""` → `"latest"`.

**Version constraints are applied in code, not templates.** `laravel/laravel` in the installer command becomes `laravel/laravel:^12.0` when `--version 12` is passed. Handled by `applyVersionConstraints` in scaffold.go.

**No globals, no shared mutable state.** All dependencies flow through `CommandContext` (injected at construction) or function parameters.

**No `os.Chdir` in commands.** The interactive shell (`shell.go`) calls `os.Chdir` only for the `cd` builtin. All other commands work with absolute paths resolved at call time.

**Docker Compose v2 only.** Commands use `docker compose` (space, not hyphen). `dependency/manager.go` verifies this in `atempo doctor`.

---

## Desktop ↔ CLI protocol

The desktop never reimplements CLI logic — it calls the sidecar binary and parses its output.

- `sidecar::run_json(app, args)` — appends `--json`, runs to completion, parses stdout as JSON
- `sidecar::stream_json(app, args, event_name)` — appends `--stream-json`, spawns the process, emits each stdout line as a Tauri event

On the CLI side, commands opt into these modes via interfaces in `command.go`:
- `JSONCommand` — implements `ExecuteJSON` for one-shot structured output
- `StreamCommand` — implements `ExecuteStream` for line-by-line progress events (create, logs)

The registry in `registry.go` detects `--json` / `--stream-json` flags and routes to the appropriate interface automatically.

---

## CLI commands

| Command | File | Notes |
|---|---|---|
| `create <fw[:ver]> <name>` | create.go | scaffold + register |
| `start [project]` | start.go | `docker compose up -d` |
| `stop [project]` | stop.go | `docker compose down` |
| `logs [project]` | logs.go | streaming compose logs |
| `status [project]` | status.go | service health |
| `projects` | projects.go | list registry |
| `describe [project]` | describe.go | project details |
| `add-service <type>` | add_service.go | append + regenerate |
| `reconfigure` | reconfigure.go | regenerate compose |
| `remove <project>` | remove.go | deregister |
| `doctor` | doctor.go | dependency checks |
| `shell` (no args) | shell.go | interactive REPL |

## Frameworks

| Name | Installer | Key services |
|---|---|---|
| `laravel[:12]` | docker (composer) | PHP-FPM, Nginx, MySQL 8, Redis, Mailhog |
| `fastapi` | template | Python 3.12, Uvicorn, PostgreSQL 15, Redis |
| `bun` | template | Bun runtime, TypeScript |
| `python` | template | Python 3.12, no ports, restart: "no" |

## Predefined services (add-service)

`minio`, `elasticsearch`, `rabbitmq`, `mongodb` — defined in `compose/generator.go`.

## Interactive shell

`shell.go` uses `github.com/chzyer/readline` v1.5.1 for history (↑/↓) and tab completion. The `AutoComplete` tree covers every atempo command plus common shell commands (`cd`, `ls`, `cat`, `grep`, `git`, `docker`), each with a dynamic file/directory completer. The `PcItemDynamic` callback receives the full input line — the last word is extracted manually with `strings.Fields` before passing to `fsCompletions`. The `~/` prefix is preserved in candidates so readline's suffix-match works correctly. Unknown commands are passed to `$SHELL -c`.

## Project file layout (created by `atempo create`)

```
<project>/
  src/                       # framework source code
  ai/
    context.md               # AI context pre-populated with stack details
  infra/docker/
    docker-compose.yml       # generated from atempo.json
    Dockerfile               # (for build-type services)
  atempo.json                # project config + assigned ports
```

## Registry storage

`~/.config/atempo/registry.json` — JSON array of `registry.Project` structs. Read on demand, cached in `MemoryCacheRepository` with a 5-minute TTL. The cache uses `sync.Mutex` (not `sync.RWMutex`) because expired-key deletion during `Get` requires a write lock.

## V1 scope boundaries

- No remote registries, no cloud deployments, no Kubernetes
- No project templates beyond the 4 supported frameworks
- No plugin system
- Port conflict resolution is best-effort scan; no reservation lock file

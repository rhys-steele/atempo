# Atempo — Rebuild Plan
**Date**: 2026-05-29  
**Repo**: `~/Documents/projects/atempo`  
**Audit**: See `rebuild-audit-2026-05-29.md`

---

## What We're Building

A CLI tool that bootstraps Docker-based web development environments. You run one command, get a fully wired project — framework source code, Docker Compose stack, AI context files, and project registration. Everything else (DNS, CI, AI generation, MCP discovery) comes later once the core is solid.

**The pitch**: smarter `create-react-app`, framework-agnostic, AI-aware from day one.

---

## V1 Scope — Ship This, Nothing Else

### Commands to build

| Command | Description | Status in old codebase |
|---|---|---|
| `atempo create <framework>[:<version>] [name]` | Scaffold + compose + register | Works, needs cleanup |
| `atempo start [project]` | `docker compose up -d` | Works |
| `atempo stop [project]` | `docker compose down` | Works |
| `atempo logs [project]` | Stream container logs | Works |
| `atempo status [project]` | Show service health | Works |
| `atempo doctor` | Check Docker, Git dependencies | Works |
| `atempo projects` | List registered projects | Works |
| `atempo describe [project]` | Show project details | Works |
| `atempo add-service <service>` | Add Redis/MySQL/etc to compose | Works |
| `atempo reconfigure` | Regenerate docker-compose.yml | Works |
| `atempo remove <project>` | Remove from registry | Works |
| `atempo shell` | Interactive shell (default when no args) | Works, needs readline |

### Frameworks to support in V1

| Framework | Templates | Docker stack | Installer |
|---|---|---|---|
| `laravel` | ✅ complete | PHP-FPM + Nginx + MySQL + Redis + Mailhog | `composer create-project` in Docker |
| `django` | ✅ complete | Python + Postgres + Redis + Celery + Mailhog | `django-admin startproject` in Docker |
| `express` | ✅ complete | Node.js | `npm init` + server.js template |

**Arduino, MicroPython, lambda-node** — templates exist but are wrong domain or incomplete. Cut from V1.

### Project output structure (what `create` produces)

```
my-app/
├── src/                    # Framework source code (installed by official installer)
├── ai/                     # Static AI context files (copied from templates)
│   ├── context.md
│   ├── project-overview.md
│   ├── codebase-map.md
│   ├── development-workflows.md
│   ├── patterns-and-conventions.md
│   └── ui-ux-guidelines.md
├── infra/
│   └── docker/
│       ├── docker-compose.yml  (generated from atempo.json)
│       ├── Dockerfile
│       └── <framework config files>
├── atempo.json             # Project metadata + service config
└── README.md
```

### What's explicitly NOT in V1

- DNS/SSL management (`atempo dns`, `atempo ssl`) — Docker-based dnsmasq system is macOS-only and adds hard setup requirements. Cut entirely. Use localhost:PORT.
- AI generation via API (`atempo ai`) — all stubs, no working implementation. Copy static templates instead.
- CI/CD generation (`atempo ci`) — internals work but command layer is stubs. Defer.
- MCP server management — discovery is a stub, community servers endpoint doesn't exist.
- `atempo repair` — regenerate_compose case returns "not yet implemented".
- `atempo setup` — removed entirely (had curl-pipe-sudo security issue).
- `atempo auth` — deprecated in old codebase, cut.
- `atempo test` — shell injection bug, cut.
- `atempo reset` — dangerous, needs redesign.
- `atempo audit` — invokes `claude /audit` shell command, not a real feature.
- `atempo report` — error reporting system, low value for V1.
- `atempo device` — hardware package, wrong domain entirely.

---

## Architecture Decisions

### Package structure (new repo)

```
atempo/
├── cmd/atempo/main.go          # Entry point, no logic
├── internal/
│   ├── errors/                 # Error types + user message formatting
│   ├── registry/               # Types only (Project, Port, Service, Registry)
│   ├── interfaces/             # CacheRepository, RegistryRepository
│   ├── repositories/           # FileRegistryRepository, MemoryCacheRepository (fixed)
│   ├── logger/                 # Project operation logger
│   ├── dependency/             # Docker/Git dependency checks
│   ├── scaffold/               # Framework scaffolding runner
│   ├── compose/                # atempo.json → docker-compose.yml generator
│   ├── docker/                 # Docker lifecycle (start/stop/logs/status)
│   └── app/
│       └── commands/           # One file per command
└── templates/                  # Embedded at build time
    ├── frameworks/
    │   ├── laravel/
    │   ├── django/
    │   └── express/
    └── shared/
```

### Key architecture rules for the rebuild

1. **embed.FS wired through properly** — `main.go` embeds `templates/**` and passes the FS into every package that needs templates. No filesystem fallback except in development (`ATEMPO_DEV=1`).

2. **One registry access path** — `internal/repositories/FileRegistryRepository` implementing `internal/interfaces/RegistryRepository`. No direct file I/O from commands. The old `internal/registry/registry.go` file methods are gone.

3. **No `os.Chdir()`** — every operation receives an explicit `projectPath string`. Commands resolve the path once at the top and pass it down.

4. **No global mutable state** — no `GlobalStatus`, `GlobalCommandFactory`, etc. Everything injected from `main.go` via the `CommandRegistry`.

5. **Docker Compose v2** — `docker compose` (space, not hyphen). Add v2 availability check to `doctor`.

6. **Ports are auto-assigned — user never thinks about them** — during `create`, Atempo scans for available ports starting from each service's preferred default (e.g. 8000 for the app, 3306 for MySQL) and picks the first free one. Assigned ports are written into `atempo.json` and used verbatim in the generated `docker-compose.yml`. The `atempo.json` is the single source of truth. Ports are shown in the output after create so the user knows where to go, but they never have to choose.

   Port detection: attempt to `net.Listen("tcp", ":<port>")` — if it binds, the port is free. Increment and retry if not.

   Each framework template defines preferred defaults per service:
   - Laravel: app=8000, mysql=3306, redis=6379, mailhog=8025
   - Django: app=8000, postgres=5432, redis=6379, mailhog=8025
   - Express: app=3000

7. **Shell injection prevention** — any user input that touches `exec.Command` is passed as a separate arg, never interpolated into a string.

8. **No API keys at rest** — V1 doesn't store any API keys. AI context is static template copy only.

8. **`strings.Title` gone** — use `golang.org/x/text/cases` everywhere.

### Fixed bugs from the audit (applied during rebuild)

- Cache `Get` data race: use `Lock()` not `RLock()` when deleting expired entries
- Auth expiry logic: compare against `time.Now().Unix()`
- No unguarded `Choices[0]` indexing (no AI API calls in V1 anyway)
- No unsafe type assertions on Docker JSON output — use typed structs
- `FindProject` side effect (writing on read) — drop the `LastAccessed` touch on read, or make it explicit

---

## Dependencies (go.mod)

**Keep**:
- `golang.org/x/term` — terminal state management for shell
- `golang.org/x/text` — replace `strings.Title`  
- `gopkg.in/yaml.v3` — docker-compose.yml generation

**Add**:
- `github.com/chzyer/readline` — shell readline (already in old codebase, keep it)

**Drop**:
- `github.com/chzyer/readline` may be replaceable with stdlib `bufio` for simple shell — evaluate during shell implementation

**Never add** (was in old codebase via hardware package dependencies):
- Any USB/serial port libraries
- Any hardware detection libraries

---

## Build Phases & Commit Points

### Phase 1 — Foundation (Commit 1)
Compiles, no commands yet. Just the shared infrastructure.

**Packages**:
- `internal/errors/` — types + utils (clean port from old, no deprecation)
- `internal/registry/types.go` — Project, Port, Service, Registry structs only
- `internal/interfaces/` — CacheRepository, RegistryRepository
- `internal/repositories/cache.go` — MemoryCacheRepository (data race fixed)
- `internal/repositories/registry.go` — FileRegistryRepository (fire-and-forget removed)
- `internal/logger/logger.go` — clean port

**Validation**: `go build ./...` passes, `go vet ./...` clean.

---

### Phase 2 — Entry Point + Command Skeleton (Commit 2)
Binary builds and runs. Prints usage for unknown commands.

**Packages**:
- `internal/app/commands/command.go` — Command interface, BaseCommand, CommandRegistry (no `strings.Title`)
- `internal/app/commands/registry.go` — CommandRegistry.Execute, IsHelpCommand, ShowUsage
- `cmd/atempo/main.go` — arg parsing, error display, `embed.FS` wired

**Validation**: `go build -o /tmp/atempo ./cmd/atempo && /tmp/atempo --help` shows usage.

---

### Phase 3 — Doctor + Dependency Check (Commit 3)
First real command. Validates the environment.

**Packages**:
- `internal/dependency/manager.go` — Docker (v2), Git checks (use `runtime.GOOS` not `exec.LookPath`)
- `internal/app/commands/doctor.go`

**Validation**: `atempo doctor` runs and correctly identifies installed/missing deps.

---

### Phase 4 — Project Registry Commands (Commit 4)
Projects, Describe, Remove.

**Packages**:
- `internal/app/commands/projects.go`
- `internal/app/commands/describe.go`
- `internal/app/commands/remove.go`

**Validation**: Create a dummy `~/.atempo/registry.json`, run `atempo projects` and `atempo describe`.

---

### Phase 5 — Compose Generator (Commit 5)
The `atempo.json` → `docker-compose.yml` translation layer.

**Packages**:
- `internal/compose/generator.go` — clean port, remove DNS integration (deferred), ports read directly from `atempo.json` service definitions
- `internal/app/commands/reconfigure.go`

**Validation**: Feed a laravel `atempo.json` with user-specified ports in, get a valid `docker-compose.yml` out. Write a test.

---

### Phase 6 — Docker Lifecycle (Commit 6)
Start, stop, logs, status.

**Packages**:
- `internal/docker/docker.go` — upgrade to `docker compose` v2, strip `docker-compose` v1 references, remove `shouldShowStatusLine` dead code
- `internal/app/commands/start.go` (renamed from implicit create flow)
- `internal/app/commands/stop.go`
- `internal/app/commands/logs.go`
- `internal/app/commands/status.go`

**Validation**: Against a real docker-compose project, all four commands work end-to-end.

---

### Phase 7 — Scaffold + Create (Commit 7)
The main event. Framework installation + template copy + registration.

**Packages**:
- `internal/scaffold/scaffold.go` — clean port, use embedded FS only (no filesystem fallback), pass explicit `projectPath` not `os.Chdir`
- `templates/` — copy laravel, django, express templates verbatim from old repo (they're already good)
- `internal/app/commands/create.go` — clean port, no AI planning stub, no `os.Chdir`

**Validation**: `atempo create laravel testing/my-app` produces correct directory structure with all files.

---

### Phase 8 — Add-Service + Shell (Commit 8)
Round out the command surface.

**Packages**:
- `internal/app/commands/add_service.go`
- `internal/app/commands/shell.go` — interactive shell with readline, git-aware prompt

**Validation**: Shell starts, basic commands work, `exit` exits cleanly.

---

### Phase 9 — CI + Tests (Commit 9)
CI pipeline and test coverage for critical paths.

**Files**:
- `.github/workflows/build.yml` — Go 1.23, `go vet`, `golangci-lint`, build + test. Use `actions/upload-artifact@v4`.
- Tests: compose generator, registry repository, error types, dependency checks
- `internal/repositories/cache_test.go` — specifically test the fixed data race path
- `internal/repositories/registry_test.go`

**Validation**: CI passes on push. `go test ./...` all green.

---

## V2 Backlog (After V1 Ships)

These are real features worth building, just not blocking V1:

| Feature | Notes |
|---|---|
| **DNS / `.local` domains** | Fully Dockerized dnsmasq+nginx system is well designed. macOS only. Needs `atempo dns setup/status/teardown`. |
| **`atempo ci`** | GitHub Actions + GitLab CI generation. Internals exist in `internal/ci/`, command layer needs rebuilding. |
| **Dynamic AI context** | Replace static template copy with API-driven generation. Use `internal/ai/client.go` as base (single HTTP client, remove the duplicate). |
| **Port manager** | Dynamic port allocation (8000–8999 range) to prevent conflicts between multiple projects. |
| **`atempo cd <project>`** | Output a `cd /path/to/project` for shell eval. Simple but high value. |
| **Framework passthrough** | `atempo artisan migrate`, `atempo manage createsuperuser` — exec into the right container. |
| **`atempo env`** | Switch `.env` profiles between dev/staging. |
| **MCP server generation** | Static MCP config generation (no community registry needed). |
| **SSL certificates** | Self-signed wildcard certs for `.local`. Needs DNS first. |

---

## Template Audit (What to Copy)

Copy these as-is — they're complete and correct:

```
templates/frameworks/laravel/   ← copy all
templates/frameworks/django/    ← copy all  
templates/frameworks/express/   ← copy all (includes .env.example, server.js)
```

Leave behind:
```
templates/frameworks/arduino/      ← wrong domain
templates/frameworks/micropython/  ← wrong domain
templates/frameworks/lambda-node/  ← only has CI templates, nothing else
templates/ai/                      ← these are generation prompts for the AI API feature (V2)
```

The `templates/shared/ci/` files (github.yml, gitlab.yml) — copy to V2 when CI command is implemented.

---

## Definition of Done for V1

- [ ] `atempo create laravel my-app` produces a working project that `docker compose up` starts successfully
- [ ] `atempo create django my-api` same
- [ ] `atempo create express my-site` same
- [ ] `atempo doctor` correctly identifies Docker v2, Git
- [ ] `atempo projects` lists registered projects
- [ ] `atempo start/stop/logs/status` all work against a real running project
- [ ] `atempo shell` launches interactive shell with git-aware prompt
- [ ] `go build ./...` with zero warnings
- [ ] `go vet ./...` clean
- [ ] `go test ./...` all pass
- [ ] CI passes on push to main
- [ ] Binary is a single self-contained executable (all templates embedded)

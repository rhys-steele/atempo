# Atempo — Claude Code Context

A Go CLI tool that bootstraps Docker-based web development environments. One command produces a fully wired project: framework source, Docker Compose stack, port assignments, and registry entry.

## Build & test

```bash
go build -o atempo ./cmd/atempo   # produces ./atempo binary
go test ./...                      # all packages
go test -race ./...               # with race detector (used in CI)
go vet ./...                      # static analysis
```

CI runs on push to `main`/`develop` and PRs to `main` (`.github/workflows/build.yml`).

## Entry point

`cmd/atempo/main.go` — embeds `all:templates` via `embed.FS`, wires the repository layer, builds the command registry, and dispatches. No arguments → interactive shell.

## Package structure

```
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
```

## Key conventions

**Templates are embedded, never copied as loose files.** All framework templates live under `cmd/atempo/templates/` and are embedded via `//go:embed all:templates`. The `all:` prefix includes dotfiles (e.g. `.env.example`).

**Path prefix stripping.** The `embed.FS` is rooted at `cmd/atempo/` so template paths start with `templates/`. Before passing to `scaffold.Run`, `create.go` calls `fs.Sub(templatesFS, "templates")` so scaffold sees `frameworks/laravel/...` directly.

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

## Commands

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

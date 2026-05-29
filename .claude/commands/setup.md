# Project Setup

Read the following to understand this codebase before doing any work:

1. `README.md` — what Atempo does and the command surface
2. `.me/rebuild-plan-2026-05-29.md` — architecture decisions, V1 scope, build phases, and commit points
3. `.me/rebuild-audit-2026-05-29.md` — known issues from the old codebase and what was fixed in the rebuild

## Key rules for this codebase

- **No `os.Chdir()`** — pass explicit `projectPath string` to all operations
- **No global mutable state** — everything injected from `main.go` via `CommandRegistry`
- **Docker Compose v2** — `docker compose` (space), never `docker-compose` (hyphen)
- **No shell interpolation** — user input always passed as separate args to `exec.Command`
- **embed.FS is the only template source** — no filesystem fallback in production
- **Ports are auto-assigned** — `net.Listen` probe from preferred default, written into `atempo.json`
- **`golang.org/x/text/cases`** — never `strings.Title`

## Package layout

```
internal/errors/          Error types + user message formatting
internal/registry/        Data types only (Project, Service, Port, Registry)
internal/interfaces/      CacheRepository, RegistryRepository interfaces
internal/repositories/    FileRegistryRepository, MemoryCacheRepository
internal/logger/          Per-operation file logger
internal/dependency/      Docker v2 + Git availability checks
internal/scaffold/        Framework installer runner
internal/compose/         atempo.json → docker-compose.yml
internal/docker/          Docker lifecycle (start/stop/logs/status)
internal/app/commands/    One file per CLI command
```

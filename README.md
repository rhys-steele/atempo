# Atempo

A CLI tool for bootstrapping Docker-based web development environments. One command gives you a wired project — framework source code, Docker Compose stack, AI context files, and project registration.

```bash
atempo create laravel:12 my-app
atempo start my-app
```

---

## Installation

Requires Go 1.23+:

```bash
go install github.com/rhys-atempo/atempo@latest
```

Or build from source:

```bash
git clone https://github.com/rhys-atempo/atempo
cd atempo
go build -o atempo ./cmd/atempo
mv atempo /usr/local/bin/
```

**Prerequisites**: Docker with Compose v2 plugin, Git.

Run `atempo doctor` to verify your environment.

---

## Usage

### Create a project

```bash
atempo create laravel:12 my-app
atempo create django my-api
atempo create express my-site
```

This produces:

```
my-app/
├── src/                    # Framework source code
├── ai/                     # AI context files for Claude
├── infra/
│   └── docker/
│       ├── docker-compose.yml
│       └── Dockerfile
└── atempo.json             # Project metadata
```

Ports are assigned automatically — no conflicts, no configuration needed.

### Manage projects

```bash
atempo start [project]         # docker compose up -d
atempo stop [project]          # docker compose down
atempo logs [project]          # stream container logs
atempo status [project]        # show service health
atempo projects                # list all registered projects
atempo describe [project]      # show project details
atempo add-service <service>   # add Redis, MySQL, etc.
atempo reconfigure             # regenerate docker-compose.yml
atempo remove <project>        # remove from registry
atempo doctor                  # check system dependencies
```

### Interactive shell

Running `atempo` with no arguments starts an interactive shell with a git-aware prompt and tab completion for all commands.

---

## Supported Frameworks

| Framework | Version syntax | Stack |
|---|---|---|
| Laravel | `laravel`, `laravel:12` | PHP-FPM, Nginx, MySQL, Redis, Mailhog |
| Django | `django`, `django:5` | Python, PostgreSQL, Redis, Celery, Mailhog |
| Express | `express` | Node.js |

---

## Project structure convention

Test and experimental projects should live in a `testing/` subdirectory:

```bash
atempo create laravel testing/experiment    # ✅
atempo create laravel my-production-app    # ✅ for real projects
```

---

## AI context

Every project gets a set of AI context files in `ai/` ready for use with Claude Code or any MCP-compatible client:

- `context.md` — master index
- `project-overview.md` — architecture and features
- `codebase-map.md` — file structure and relationships
- `development-workflows.md` — commands and processes
- `patterns-and-conventions.md` — coding standards
- `ui-ux-guidelines.md` — design standards

---

## License

[Atempo License v1.0](./LICENCE) — open source with attribution.

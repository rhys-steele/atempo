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

Ports are assigned automatically — no conflicts, no configuration needed.

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

## Supported Frameworks

| Framework | Version syntax | Stack |
|---|---|---|
| `laravel` | `laravel`, `laravel:12` | PHP-FPM, Nginx, MySQL, Redis, Mailhog |
| `fastapi` | `fastapi` | Python 3.12, Uvicorn, PostgreSQL, Redis |
| `bun` | `bun` | Bun runtime, TypeScript, `Bun.serve()` |
| `python` | `python` | Python 3.12, plain container (scripts/workers) |

---

## Predefined services

Add any of these to an existing project with `atempo add-service <name>`:

| Service | Image | Ports |
|---|---|---|
| `minio` | minio/minio | 9000, 9001 |
| `elasticsearch` | elasticsearch:8.8.0 | 9200 |
| `rabbitmq` | rabbitmq:3-management | 5672, 15672 |
| `mongodb` | mongo:6 | 27017 |

---

## AI context

Every project gets an `ai/context.md` file pre-populated with the stack, structure, key patterns, and common Docker commands — ready to load into Claude Code or any MCP-compatible client.

---

## License

[Atempo License v1.0](./LICENCE) — open source with attribution.

# {{project}}

A FastAPI application scaffolded with [Atempo](https://github.com/atempo-ai/atempo).

## Stack

- **FastAPI** — Python HTTP API framework
- **PostgreSQL 16** — Database
- **Redis** — Cache / task queue
- **Docker Compose** — Local development environment

## Getting started

```bash
# Start all services
atempo start

# Stream logs
atempo logs

# Check service health
atempo status
```

## Development

The app runs with `--reload` in development, so changes to `src/` restart automatically.

```
src/
├── main.py          # Application entry point
├── requirements.txt # Python dependencies
└── .env.example     # Environment variable template
```

Copy `.env.example` to `.env` and adjust values as needed.

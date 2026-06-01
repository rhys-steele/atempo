# {{project}} — AI Context

## Project

**Framework**: FastAPI  
**Language**: Python  
**Created**: {{date}}

## Stack

| Service | Image | Port |
|---------|-------|------|
| app | python:3.12-slim (build) | 8000 |
| postgres | postgres:16-alpine | 5432 |
| redis | redis:alpine | 6379 |

## Project structure

```
src/
├── main.py          # FastAPI application
├── requirements.txt # Python dependencies
└── .env.example     # Environment variables
infra/docker/
├── Dockerfile
└── requirements.txt
```

## Key patterns

- Routes defined in `src/main.py` using FastAPI decorators
- Database connection string via `DATABASE_URL` environment variable
- Redis connection string via `REDIS_URL` environment variable
- Health check endpoint at `GET /health`

## Running commands

```bash
# Enter app container
docker compose exec app bash

# Install new dependency
docker compose exec app pip install <package>
# Then add to src/requirements.txt

# Run database migrations (if using Alembic)
docker compose exec app alembic upgrade head
```

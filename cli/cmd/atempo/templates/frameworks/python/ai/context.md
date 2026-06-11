# {{project}} — AI Context

## Project

**Framework**: Python (plain)  
**Language**: Python  
**Created**: {{date}}

## Stack

| Service | Image |
|---------|-------|
| app | python:3.12-slim (build) |

No exposed ports — this is a script/worker, not an HTTP server.

## Project structure

```
src/
├── main.py          # Entry point — start here
└── requirements.txt # Python dependencies
infra/docker/
└── Dockerfile
```

## Key patterns

- `src/main.py` is the entry point: `python main.py`
- Add dependencies to `src/requirements.txt`, then rebuild: `docker compose build`
- Container exits after `main()` returns (use `restart: unless-stopped` for workers)

## Running commands

```bash
# Run the script
docker compose run --rm app

# Enter interactive shell
docker compose run --rm app python

# Install new dependency
docker compose exec app pip install <package>
# Then add to src/requirements.txt and rebuild
```

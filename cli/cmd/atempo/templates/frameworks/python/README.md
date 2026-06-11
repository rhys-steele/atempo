# {{project}}

A containerised Python script scaffolded with [Atempo](https://github.com/atempo-ai/atempo).

## Stack

- **Python 3.12** — Runtime
- **Docker Compose** — Local development environment

## Getting started

```bash
# Run once
atempo start

# Stream logs
atempo logs

# Check service health
atempo status
```

## Development

Add your code to `src/main.py`. Add dependencies to `src/requirements.txt`.

```
src/
├── main.py          # Entry point
└── requirements.txt # Python dependencies
infra/docker/
└── Dockerfile
```

The container runs `python main.py` and exits. For a long-running process (worker/daemon), change `CMD` in the Dockerfile and set `restart: unless-stopped` in `atempo.json`.

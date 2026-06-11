# {{project}}

A minimal Bun HTTP server scaffolded with [Atempo](https://github.com/atempo-ai/atempo).

## Stack

- **Bun** — JavaScript/TypeScript runtime (`oven/bun`)
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

The server is at `src/index.ts`. Add routes inside the `fetch` handler.

```
src/
├── index.ts      # Bun.serve entry point
└── package.json
```

Hot reload: `bun --watch index.ts` (enabled automatically in the dev container).

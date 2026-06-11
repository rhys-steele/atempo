# {{project}} — AI Context

## Project

**Framework**: Bun  
**Language**: TypeScript  
**Created**: {{date}}

## Stack

| Service | Image | Port |
|---------|-------|------|
| app | oven/bun:1 (build) | 3000 |

## Project structure

```
src/
├── index.ts      # Bun.serve entry point — all routes here
└── package.json
infra/docker/
└── Dockerfile
```

## Key patterns

- HTTP server via `Bun.serve({ fetch(req) { ... } })`
- Route dispatch by `url.pathname` inside `fetch`
- No framework — raw Bun.serve only
- JSON responses via `Response.json()`

## Running commands

```bash
# Enter app container
docker compose exec app sh

# Add a dependency
docker compose exec app bun add <package>
# Then commit the updated package.json and bun.lockb
```

# DNS System Redesign: Simple & User-Friendly

## Goals
1. ✅ Custom domain names: `new-project.local`
2. ✅ Subdomain routing: `mailhog.new-project.local`  
3. ✅ Simple, friendly UX for DNS setup
4. ✅ Auto-registration on `create` (no sudo per project)
5. ✅ No local installations required

## Ultra-Simple Architecture

### Single Container Solution
```
┌─────────────────────────────────────────┐
│ atempo-dns (single container)          │
│ ┌─────────────┐ ┌─────────────────────┐ │
│ │   dnsmasq   │ │    nginx (simple)   │ │
│ │  port: 5353 │ │    port: 80         │ │
│ └─────────────┘ └─────────────────────┘ │
└─────────────────────────────────────────┘
```

### Ultra-Simple Flow
1. **One-time setup**: `atempo dns setup` → creates resolver, starts container
2. **Project creation**: `atempo create laravel my-app` → auto-adds DNS entries
3. **Access**: Browser `my-app.local` → nginx routes to localhost:8080

## Implementation: Maximum Simplicity

### Phase 1: Single File DNS Manager (15 min)
- One simple struct with 4 methods: Setup(), Start(), AddProject(), Status()
- No complex state management, no mutexes, no orchestration

### Phase 2: Single Container with Both Services (10 min)
- nginx:alpine base image
- Install dnsmasq via apk add
- Simple startup script runs both services

### Phase 3: File-Based Config (5 min)
- Write DNS entries to simple files
- Restart container to reload (simple and reliable)

## Benefits of This Approach
- **95% less code**: Remove all complex managers, orchestration, networking
- **Zero silent failures**: Simple error handling, obvious what went wrong
- **Debuggable**: One container, clear logs, easy to inspect
- **User-friendly**: One command setup, automatic project registration
- **Maintainable**: Simple, direct, no clever tricks

## Why This is Better
Current system tries to be too clever with:
- Complex nginx proxy orchestration
- Multi-container coordination  
- Mutex locks and state management
- Silent fallback logic

New system is dumb and simple:
- One container does everything
- File-based config (restart to reload)
- Direct error reporting
- Zero clever abstractions
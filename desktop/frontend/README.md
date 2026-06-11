# Atempo Desktop — Frontend

React + Vite + Tailwind frontend for the Atempo desktop app. Runs inside a Tauri webview — do not open standalone in a browser.

## Dev

Run from `desktop/backend` (Tauri manages starting this):

```bash
cargo tauri dev
```

Or run the Vite server alone (no Tauri invoke calls will work):

```bash
npm run dev
```

## Structure

```
src/
  lib/tauri.ts          # invoke() wrappers + shared TypeScript types
  hooks/useProjects.ts  # React Query hooks
  pages/
    Dashboard.tsx        # project list
    CreateProject.tsx    # multi-step create wizard (streaming logs)
    ProjectDetail.tsx    # start/stop/logs/reconfigure/remove
    Doctor.tsx           # dependency checks
  components/
    Layout.tsx           # sidebar + outlet
    Sidebar.tsx          # nav links
    ProjectCard.tsx      # project summary card
    ServiceBadge.tsx     # status pill
    LogStream.tsx        # scrolling log output
```

## Notes

- `verbatimModuleSyntax` is enabled — use `import type` for type-only imports
- All backend communication goes through `@tauri-apps/api/core` `invoke()` — no HTTP
- Streaming operations (create, logs) use Tauri events (`listen()`) emitted by the Rust sidecar layer

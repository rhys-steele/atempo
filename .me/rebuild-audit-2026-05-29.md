# Steele / Atempo — Rebuild Audit
**Date**: 2026-05-29  
**Codebase**: ~25,000+ lines across 100+ Go files  
**Module**: `atempo` (Go 1.23.0)  
**Purpose**: CLI tool for bootstrapping Docker-based web development environments (Laravel, Django, Express)

---

## Executive Summary

**Overall Health: 4/10 — Not shippable in current state.**

The codebase grew fast without consolidation checkpoints. Three generations of the same abstractions coexist, ~5,000 lines belong to a completely different product (embedded systems/Arduino), several critical paths are stub-only, there are compile-breaking duplicate function declarations, a data race, multiple runtime panic risks, and real security issues. The foundation (Docker Compose generation, project registry, scaffold templates) is solid — the surrounding architecture needs a hard cut.

### What works and is worth keeping
- `internal/compose/generator.go` — Docker Compose file generation
- `internal/docker/` — Docker lifecycle management (needs v2 upgrade)
- `internal/scaffold/` — Framework scaffolding templates
- `internal/repositories/` — Repository pattern (cleanest abstraction layer)
- `internal/errors/` — Error types and recovery suggestions
- `internal/logger/` — Project logger
- `internal/ci/` — CI/CD config generation

### What needs to be cut entirely
- `internal/hardware/` — 5,004 lines of Arduino/ESP32 embedded systems tooling. Zero connection to web scaffolding. Delete.
- `internal/app/commands/auth.go` — Fully deprecated, shows deprecation notice on execute.
- `internal/ai/auth.go` — Duplicate HTTP client alongside `client.go`. Delete.
- `internal/app/commands/ai_manifest_clean.go` — `generateTemplatedFile` only prints to stdout, never writes files. Pure demo code.
- `internal/app/commands/template_loader.go` — Duplicate type definitions from `internal/ai/context.go`. Delete.
- `internal/registry/registry.go` — Deprecated registry access layer. Superseded by `internal/repositories/`.
- `internal/services/registry.go` + `registry_service.go` — Duplicate `NewRegistryService` declarations (compile error).

---

## Critical Bugs (Fix Before Anything Else)

### 1. Compile Error — Duplicate `NewRegistryService`
**Files**: `internal/services/registry.go:27`, `internal/services/registry_service.go:51`  
Both define `NewRegistryService` in the same `services` package with different signatures. The package cannot build.

**Fix**: Delete one entirely. The context-based `registry_service.go` variant (lowercase `registryService`) is the cleaner design — keep it, delete `registry.go`.

### 2. Runtime Panic — Unsafe Type Assertions
**File**: `internal/registry/registry.go`  
```go
serviceName := serviceData["Service"].(string)  // panics if key missing
state := serviceData["State"].(string)           // panics if key missing
```
If Docker outputs a service without the expected JSON keys, this crashes the entire process.

**Fix**: Use comma-ok assertions or `json.Unmarshal` into a typed struct.

### 3. Runtime Panic — Unbounded Slice Index
**File**: `internal/services/ai.go`  
`response.Choices[0]` is indexed without checking `len(response.Choices) > 0` in both `GenerateManifest` and `AnalyzeProject`. Any AI provider returning an empty choices array crashes the process.

**Fix**: Guard with `if len(response.Choices) == 0 { return "", fmt.Errorf("empty response") }`.

### 4. Data Race — Write Under Read Lock
**File**: `internal/repositories/cache.go`  
```go
func (c *MemoryCacheRepository) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    // ...
    delete(c.items, key)  // WRITE OPERATION under RLock — data race
```
**Fix**: Upgrade to `c.mu.Lock()` / `c.mu.Unlock()` for the `Get` method, or use a separate cleanup pass.

### 5. Broken Token Expiry Logic
**File**: `internal/interfaces/auth.go`  
```go
func (c *Credential) IsValid() bool   { return c.ExpiresAt > 0 }   // always true for any real timestamp
func (c *Credential) IsExpired() bool { return c.ExpiresAt < 0 }   // always false
```
Neither function compares against `time.Now().Unix()`.

**Fix**:
```go
func (c *Credential) IsExpired() bool { return time.Now().Unix() > c.ExpiresAt }
func (c *Credential) IsValid() bool   { return c.ExpiresAt > 0 && !c.IsExpired() }
```

---

## Security Issues

### S1 — Curl-Pipe-Sudo Install Script [CRITICAL]
**File**: `internal/app/commands/setup.go`  
Downloads `https://get.docker.com` and runs it via `sudo sh get-docker.sh`. Classic supply-chain attack vector — network interception or CDN compromise executes arbitrary code as root.

**Fix**: Remove this entirely from the CLI. Link to official Docker docs instead. If you must automate, verify a checksum against a pinned hash before executing.

### S2 — API Keys in Plaintext Files [HIGH]
**File**: `internal/ai/config.go`  
API keys written to `~/.atempo/auth/*.key` with `0600`. The code already has a comment acknowledging this should use OS keychain.

**Fix on rebuild**: Use `github.com/zalando/go-keyring` for OS keychain storage on all three platforms. It's a small dependency and a one-line swap.

### S3 — Token Auth Is File-Existence Check Only [HIGH]
**File**: `internal/app/commands/auth_checker.go`  
```go
if _, err := os.Stat(tokenPath); os.IsNotExist(err) { return false }
return true  // any file at ~/.atempo/auth.token grants full auth
```
Creating an empty file grants authenticated access. Token content is never read, let alone validated.

**Fix**: Read the token, validate its format, and if the token is meant to be validated against a server, do that call.

### S4 — Shell Command Injection [HIGH]
**File**: `internal/app/commands/test.go`  
```go
fmt.Sprintf("python manage.py test %s", testSuite)
```
`testSuite` is raw user input concatenated into a shell command. A value like `; rm -rf ~` executes as-is.

**Fix**: Use `exec.Command("python", "manage.py", "test", testSuite)` directly — never shell-interpolate user input.

### S5 — Error Reports Capture Environment Variables [MEDIUM]
**File**: `internal/reporting/error_reporter.go`  
Error reports stored in `~/.atempo/reports/` capture the full environment, which may include `AWS_SECRET_ACCESS_KEY`, `GITHUB_TOKEN`, database passwords, etc.

**Fix**: Allowlist specific safe env vars (`HOME`, `PATH`, `SHELL`, `TERM`, `ATEMPO_HOME`) rather than capturing all of them.

---

## Architectural Problems

### A1 — Three Generations of Registry Access
Commands import from three different layers with no consistent path:

| Layer | File | Status |
|---|---|---|
| Direct file I/O | `internal/registry/registry.go` | Has deprecated marker |
| Service wrapper | `internal/services/project.go` | Calls deprecated `LoadRegistry()` |
| Repository pattern | `internal/repositories/registry.go` | Cleanest, used by newer code |

**Fix on rebuild**: One path only — repository pattern through an interface. Delete the other two.

### A2 — `os.Chdir()` as Side Effect
**Files**: `internal/app/commands/create.go`, `internal/app/commands/registry_projects.go`  
Calling `os.Chdir()` changes the working directory for the entire process. Not goroutine-safe. If anything runs concurrently, it silently uses the wrong directory.

**Fix**: Pass explicit paths to all operations. Never change the process CWD.

### A3 — `internal/ai/context.go` Is a God Object
1,799-line file with functions exceeding 200 lines each. Generates all AI context `.md` files in one monolithic function chain.

**Fix on rebuild**: One struct per context document, each with a `Generate() string` method. Small, testable, composable.

### A4 — Docker Compose v1 Throughout
All Docker management uses `docker-compose` (v1 binary), which is deprecated and removed from modern Docker installs. v2 is `docker compose` (space, not hyphen) and ships as a plugin.

**Fix on rebuild**: Replace all `docker-compose` invocations with `docker compose`. Add a preflight check in `doctor.go` that validates the v2 plugin is available.

### A5 — Scaffold Only Works From Repo Root
**File**: `internal/services/scaffold.go`  
Validates framework by checking relative path `templates/frameworks/`. The binary must be run from the repo root, or it silently fails to find templates.

**Fix**: Use `embed.FS` properly (see below) or resolve paths relative to the binary's install location.

### A6 — Dead `embed.FS` Declaration
**File**: `cmd/atempo/main.go`  
```go
//go:embed templates
var templateFS embed.FS  // zero value — templates loaded from filesystem
```
The embed directive is present but the `FS` is never passed to any package that uses it. Templates fall through to filesystem loading. This means shipping the binary without the `templates/` directory alongside it breaks scaffolding.

**Fix on rebuild**: Actually thread `templateFS` through to `internal/scaffold/` and `internal/ai/context.go`. Then remove the filesystem fallback.

### A7 — Duplicate Executor Implementations
`internal/docker/enhanced_command_executor.go` (440 lines) and `internal/docker/enhanced_executor.go` (345 lines) both implement enhanced Docker command execution with retry/timeout. They coexist without cross-referencing.

**Fix**: Pick one, delete the other.

### A8 — Global Mutable State
Multiple packages declare global mutable singletons:
- `commands.GlobalStatus`
- `utils.GlobalCommandFactory`
- `utils.GlobalFrameworkDetector`
- `utils.GlobalStatusDisplay`
- `utils.GlobalTemplateProcessor`

These make testing impossible without side effects and create implicit coupling.

**Fix on rebuild**: Dependency injection from `main.go` downward. No package-level `var Global*`.

### A9 — Platform Detection by Path Lookup
**File**: `internal/recovery/manager.go`  
```go
if _, err := exec.LookPath("open"); err == nil { return "macos" }
```
`open` exists on some Linux systems too. Use `runtime.GOOS` instead.

---

## Dead Code / Stubs to Cut

| File | Dead Code |
|---|---|
| `commands/ai.go` | `handleInit`, `handleAnalyze`, `handleSuggest`, `handleValidate` — all return "coming soon" |
| `commands/ci.go` | `executeRun`, `executeValidate`, `executeRemove` — all unimplemented |
| `commands/docker_recovery.go` | `confirmRollback()` always returns `true`, never prompts |
| `commands/repair.go` | `"regenerate_compose"` case returns "not yet implemented" |
| `mcp/discovery.go` | `fetchCommunityServers` returns `nil`; `getServerTemplatePath` returns `""` |
| `hardware/device_detector.go` | `detectMacOS()` and `detectWindows()` discard all output with `_ = output // TODO` |
| `auth/providers.go` | `AtempoProvider.Authenticate` returns "not yet implemented" |

Cut these or implement them. Stub commands that return "coming soon" erode user trust more than just omitting the subcommand.

---

## `strings.Title` Deprecation
Used in 13 locations across 7 files. Replace with `golang.org/x/text/cases`:
```go
// Before
strings.Title(s)

// After
cases.Title(language.English).String(s)
```
Files: `command.go`, `context.go`, `template_loader.go`, `utils/command.go`, `utils/template.go`, `services/status.go`, `ci/interactive_prompter.go`.

---

## CI/CD Issues

**File**: `.github/workflows/build.yml`

| Issue | Detail |
|---|---|
| Go version mismatch | CI uses Go 1.21, `go.mod` specifies 1.23.0 with toolchain 1.24.4 |
| Deprecated actions | `actions/upload-artifact@v3` and `actions/download-artifact@v3` (v4 is current) |
| No linting step | No `golangci-lint` or `go vet` in CI |
| No security scanning | No `govulncheck` or SAST tooling |

---

## Test Coverage Gaps

Current tests are concentrated in `internal/utils/` and `internal/repositories/`. Major gaps:

| Area | Gap |
|---|---|
| `internal/compose/generator.go` | No tests for compose file generation or port allocation |
| `internal/docker/` | Minimal — only basic docker_test.go |
| `internal/ai/` | Only `client_test.go` |
| `internal/scaffold/` | No tests |
| `internal/hardware/` | No tests (but this should be deleted) |
| Integration | No end-to-end test for the create → scaffold → docker-up workflow |

---

## README Issue
`README.md` line 143 contains a bare `=======` — a git merge conflict marker that was never resolved.

---

## Rebuild Action Plan

### Phase 1 — Foundation (Before Writing New Code)

These are preconditions for a stable rebuild. Do these first.

1. **Delete `internal/hardware/`** — 5,004 lines of wrong-domain code. If this is a separate product, start a new repo for it.
2. **Fix compile error** — Delete `internal/services/registry.go`, keep `registry_service.go`.
3. **Decide on one registry access pattern** — Repository pattern only. Delete `internal/registry/registry.go`.
4. **Fix broken auth expiry** in `internal/interfaces/auth.go`.
5. **Fix data race** in `internal/repositories/cache.go`.
6. **Fix panic risks** — type assertions in `registry.go`, unbounded index in `services/ai.go`.
7. **Remove `setup.go` curl-pipe-sudo** — Replace with a link to Docker install docs.
8. **Fix shell injection** in `test.go` — Use `exec.Command` with separate args.

### Phase 2 — Architecture Cleanup

1. **Thread `embed.FS` properly** from `main.go` through scaffold and AI context packages. Remove filesystem fallback.
2. **Replace all `os.Chdir()` calls** with explicit path parameters.
3. **Migrate to Docker Compose v2** (`docker compose` not `docker-compose`).
4. **Delete all duplicate implementations**: `ai/auth.go`, `docker/enhanced_executor.go` (one of the two), `template_loader.go`, `ai/context.go` duplicate types.
5. **Break up `ai/context.go`** into per-document structs.
6. **Replace all global mutable state** with constructor-injected dependencies.
7. **Move API key storage to OS keychain** using `go-keyring`.

### Phase 3 — Feature Completeness

Only ship commands that work. For the rebuild, decide which surface area is V1:

**Recommended V1 scope** (everything that already has working internals):
- `create` — scaffold + compose generation ✓
- `start` / `stop` / `logs` / `status` — Docker lifecycle ✓
- `doctor` — dependency check ✓
- `describe` — project info ✓
- `add-service` — compose service addition ✓

**Defer to V2** (internals are stubs or broken):
- `ai` subcommands
- `ci` subcommands  
- `repair`
- `audit`
- MCP discovery/community servers

Remove the deprecated `auth` command entirely.

### Phase 4 — Quality Gates

1. Add `golangci-lint` to CI with at minimum: `errcheck`, `govet`, `staticcheck`, `gosec`.
2. Add `govulncheck` to CI.
3. Fix Go version mismatch in CI (use 1.23 or 1.24, match `go.mod`).
4. Upgrade to `actions/upload-artifact@v4`.
5. Add integration test for the full create workflow.
6. Fix `strings.Title` deprecation across 7 files.
7. Fix README merge conflict marker.

---

## Dependency Assessment

**`go.mod` dependencies to review on rebuild**:
- `github.com/spf13/cobra` — keep, solid CLI framework
- `github.com/charmbracelet/bubbletea` + `lipgloss` + `bubbles` — keep for TUI
- `gopkg.in/yaml.v3` — keep
- `github.com/go-playground/validator/v10` — keep for input validation
- Any hardware/serial/USB dependencies — delete with the hardware package

**Add on rebuild**:
- `github.com/zalando/go-keyring` — OS keychain for API keys
- `golang.org/x/text` — replace deprecated `strings.Title`
- `github.com/google/shlex` — safe shell argument splitting if you need to construct commands from user input

---

## Summary Table

| Category | Issue Count | Most Critical |
|---|---|---|
| Compile errors | 1 | Duplicate `NewRegistryService` |
| Runtime panics | 2 | Type assertions in registry, slice index in ai service |
| Data races | 1 | Cache RLock/write |
| Security | 5 | Curl-pipe-sudo, plaintext keys, shell injection |
| Wrong domain code | 1 | Entire `hardware/` package |
| Dead stubs | 7 | AI, CI, MCP subcommands |
| Architectural debt | 9 | Three registry layers, duplicate executors, global state |
| Deprecations | 13 | `strings.Title` across 7 files |
| CI issues | 4 | Go version mismatch, old actions |

**Bottom line**: Cut the hardware package, fix the 4 correctness bugs, pick one architecture layer for registry access, and implement only the features that actually work. The scaffolding and Docker lifecycle core is solid. Everything else needs either completion or deletion — there is no middle ground on stub commands.

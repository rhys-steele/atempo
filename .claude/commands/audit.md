# Codebase Audit

Run a thorough audit of this Go CLI codebase. Focus on correctness, security, and maintainability — not hypothetical improvements.

## Audit Areas

### Correctness
- Data races (check RLock/Lock usage, shared state across goroutines)
- Panic risks (unguarded type assertions, unguarded slice indexing)
- Error handling (swallowed errors, missing error checks)
- Logic bugs (off-by-one, incorrect comparisons, wrong return values)

### Security
- Shell injection (user input interpolated into command strings — use `exec.Command` with separate args)
- Sensitive data at rest (API keys, tokens written to files)
- File permissions (world-readable credential files)
- Input validation at system boundaries

### Architecture
- Dead code and unimplemented stubs
- Duplicate implementations of the same thing
- Global mutable state
- `os.Chdir()` as a side effect
- Packages that import each other creating cycles

### Go-specific
- Deprecated API usage (`strings.Title` → `golang.org/x/text/cases`)
- Context propagation (HTTP clients without context, no timeout)
- Resource cleanup (deferred Close, goroutine leaks)
- Interface satisfaction (check implementations match their interfaces)

### Tests
- Missing tests for critical paths (registry, compose generation, port assignment)
- Tests that test the wrong thing
- Race conditions in tests

## Deliverables

Write findings to `.me/codebase-audit-YYYY-MM-DD.md` with:
- Executive summary and health score
- Each issue with file path, line number, severity (Critical/High/Medium/Low)
- Concrete fix for each issue (not just "improve this")
- Prioritised action plan

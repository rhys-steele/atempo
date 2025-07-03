# Patterns & Conventions

## Code Architecture Patterns

### Clean Architecture Implementation
The project follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────────┐
│    CLI Interface    │  ← cmd/atempo/main.go
├─────────────────────┤
│   Command Layer     │  ← internal/app/commands/
├─────────────────────┤
│  Business Logic     │  ← internal/scaffold, internal/registry
├─────────────────────┤
│  Infrastructure     │  ← internal/docker, internal/logger
└─────────────────────┘
```

### Command Pattern Implementation
All CLI commands implement a consistent interface:

```go
type Command interface {
    Execute(ctx *CommandContext) error
    GetDescription() string
}

type BaseCommand struct {
    Name        string
    Description string
    Usage       string
}
```

### Registry Pattern
Central command discovery and routing:

```go
type CommandRegistry struct {
    commands map[string]Command
    projects *registry.Registry
}

func (r *CommandRegistry) RegisterCommand(name string, cmd Command) {
    r.commands[name] = cmd
}
```

## Naming Conventions

### File Naming
- **Snake case for files**: `ai_manifest.go`, `interactive_prompt.go`
- **Descriptive names**: Files clearly indicate their purpose
- **Grouped by responsibility**: Commands in `commands/`, business logic in separate packages

### Go Naming Conventions
- **Public functions**: PascalCase (`GenerateManifest`)
- **Private functions**: camelCase (`generateProjectIntent`)
- **Constants**: UPPER_SNAKE_CASE (`DEFAULT_TIMEOUT`)
- **Interfaces**: Often end with -er (`Command`, `Generator`)

### Package Organization
```go
// Package names are lowercase, single words
package commands
package scaffold
package registry

// Import aliases for clarity
import (
    "github.com/chzyer/readline"
    "gopkg.in/yaml.v3"
)
```

## Error Handling Patterns

### Consistent Error Messages
```go
// Descriptive error messages with context
return fmt.Errorf("failed to generate docker-compose.yml: %w", err)

// Include relevant details
return fmt.Errorf("framework %s version %s not found", framework, version)

// Use error wrapping for context
if err := scaffold.Run(config); err != nil {
    return fmt.Errorf("scaffolding failed: %w", err)
}
```

### Error Types
```go
// Custom error types for specific cases
type TemplateNotFoundError struct {
    Framework string
    Version   string
}

func (e *TemplateNotFoundError) Error() string {
    return fmt.Sprintf("template not found: %s:%s", e.Framework, e.Version)
}
```

## Configuration Patterns

### Template Configuration Structure
```json
{
  "name": "{{project}}",
  "framework": "laravel",
  "language": "php",
  "installer": {
    "type": "docker",
    "command": ["docker", "run", "--rm", "-v", "{{cwd}}:/workspace", "..."]
  },
  "services": {
    "app": {
      "type": "build",
      "dockerfile": "infra/docker/Dockerfile"
    }
  }
}
```

### Variable Substitution Pattern
```go
// Template variables
var templateVars = map[string]string{
    "{{project}}": projectName,
    "{{cwd}}":     currentDir,
    "{{version}}": version,
}

// Safe substitution
func substituteVariables(content string, vars map[string]string) string {
    result := content
    for key, value := range vars {
        result = strings.ReplaceAll(result, key, value)
    }
    return result
}
```

## Logging Patterns

### Structured Logging
```go
// Logger with context
logger := logger.New(projectName, projectPath)

// Step-based logging
step := logger.StartStep("Installing dependencies")
step.SetStatus(logger.StatusRunning)
// ... operation
step.SetStatus(logger.StatusCompleted)
```

### Progress Tracking
```go
// Progress with stages
tracker := progress.NewTracker(4)
tracker.StartStep("AI Planning")
// ... AI work
tracker.CompleteStep()

tracker.StartStep("Template Loading")
// ... template work
tracker.CompleteStep()
```

## Testing Patterns

### Test Organization
```go
// Test files follow _test.go convention
// scaffold_test.go
// registry_test.go
// docker_test.go

// Test function naming
func TestScaffoldRun(t *testing.T) { /* ... */ }
func TestRegistryLoadProjects(t *testing.T) { /* ... */ }
```

### Mock Patterns
```go
// Interface-based mocking
type MockDocker struct {
    commands []string
}

func (m *MockDocker) ExecuteCommand(cmd string) error {
    m.commands = append(m.commands, cmd)
    return nil
}
```

## Documentation Patterns

### Code Documentation
```go
// Package documentation
// Package commands provides the CLI command interface for the Atempo tool.
// It implements a registry pattern for command discovery and routing.
package commands

// Function documentation
// GenerateManifest creates AI-powered project manifest files based on user input
// and framework analysis. It returns the generated manifest path or an error.
func GenerateManifest(config *Config) (string, error) {
    // Implementation
}
```

### README Structure
```markdown
# Component Name

## Purpose
Brief description of the component's role

## Key Features
- Feature 1
- Feature 2

## Usage
```go
example code
```

## Architecture
Description of how it fits into the overall system
```

## Data Structure Patterns

### Registry Data Structure
```go
type Project struct {
    Name        string                 `json:"name"`
    Path        string                 `json:"path"`
    Framework   string                 `json:"framework"`
    Language    string                 `json:"language"`
    Services    map[string]Service     `json:"services"`
    Status      string                 `json:"status"`
    CreatedAt   time.Time             `json:"created_at"`
    UpdatedAt   time.Time             `json:"updated_at"`
}
```

### Configuration Validation
```go
func (c *AtempoConfig) Validate() error {
    if c.Name == "" {
        return errors.New("project name is required")
    }
    if c.Framework == "" {
        return errors.New("framework is required")
    }
    // ... additional validation
    return nil
}
```

## Concurrency Patterns

### Safe Concurrent Access
```go
type Registry struct {
    mu       sync.RWMutex
    projects map[string]*Project
}

func (r *Registry) GetProject(name string) (*Project, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    project, exists := r.projects[name]
    return project, exists
}
```

### Timeout Handling
```go
func ExecuteWithTimeout(cmd string, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    return exec.CommandContext(ctx, "docker-compose", cmd).Run()
}
```

## Security Patterns

### Credential Management
```go
// Secure credential storage
type AuthManager struct {
    credentialsPath string
}

func (a *AuthManager) StoreCredentials(apiKey string) error {
    // Encrypt before storing
    encrypted, err := encrypt(apiKey)
    if err != nil {
        return err
    }
    return os.WriteFile(a.credentialsPath, encrypted, 0600)
}
```

### Input Validation
```go
func validateFrameworkVersion(framework, version string) error {
    // Validate framework name
    if !isValidFramework(framework) {
        return fmt.Errorf("invalid framework: %s", framework)
    }
    
    // Validate version format
    if !isValidVersion(version) {
        return fmt.Errorf("invalid version: %s", version)
    }
    
    return nil
}
```

## Performance Patterns

### Lazy Loading
```go
type Registry struct {
    projects map[string]*Project
    loaded   bool
}

func (r *Registry) LoadProjects() error {
    if r.loaded {
        return nil
    }
    // Load projects from disk
    r.loaded = true
    return nil
}
```

### Caching
```go
type TemplateCache struct {
    cache map[string]*Template
    mu    sync.RWMutex
}

func (c *TemplateCache) Get(key string) (*Template, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    template, exists := c.cache[key]
    return template, exists
}
```

## Integration Patterns

### Docker Integration
```go
// Standard Docker command execution
func (d *DockerClient) ExecuteCommand(args ...string) error {
    cmd := exec.Command("docker-compose", args...)
    cmd.Dir = d.projectPath
    return cmd.Run()
}

// Health checking
func (d *DockerClient) CheckHealth() ([]ServiceStatus, error) {
    // Implementation
}
```

### Git Integration
```go
type GitInfo struct {
    Branch   string `json:"branch"`
    Status   string `json:"status"`
    HasChanges bool `json:"has_changes"`
}

func GetGitInfo(projectPath string) (*GitInfo, error) {
    // Implementation
}
```

## Extension Patterns

### Plugin Architecture (Future)
```go
type Plugin interface {
    Name() string
    Execute(ctx *Context) error
    Register(registry *CommandRegistry)
}

type PluginManager struct {
    plugins []Plugin
}

func (pm *PluginManager) LoadPlugins(dir string) error {
    // Implementation
}
```

### Framework Support
```go
type Framework interface {
    Name() string
    SupportedVersions() []string
    Install(config *Config) error
    Validate(config *Config) error
}

type FrameworkRegistry struct {
    frameworks map[string]Framework
}
```

## Common Utilities

### File Operations
```go
func CopyDir(src, dst string) error {
    return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        // Implementation
    })
}
```

### Version Comparison
```go
func CompareVersions(v1, v2 string) (int, error) {
    // Semantic version comparison
    // Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
}
```

## Best Practices Summary

### Code Organization
- Use Clear Architecture layers
- Implement consistent interfaces
- Follow Go naming conventions
- Keep functions focused and small

### Error Handling
- Provide descriptive error messages
- Use error wrapping for context
- Handle edge cases gracefully
- Log errors appropriately

### Performance
- Use lazy loading where appropriate
- Implement caching for expensive operations
- Handle timeouts for external calls
- Use goroutines for concurrent operations

### Security
- Validate all inputs
- Store credentials securely
- Use appropriate file permissions
- Sanitize user-provided data

### Testing
- Write unit tests for business logic
- Use mocks for external dependencies
- Test error conditions
- Provide integration tests

### Documentation
- Document public APIs
- Provide usage examples
- Keep documentation up to date
- Use clear, concise language

These patterns and conventions ensure consistency across the Atempo codebase and provide guidance for future development and contributions.
---
last_updated: 2024-07-19
version: 1.0
dependencies: [development-workflows.md, codebase-map.md]
priority: high
---

# Atempo CLI API Reference

## Command Interface Specification

All commands implement the base `Command` interface:

```go
type Command interface {
    Execute(ctx *CommandContext) error
    GetDescription() string
}

type CommandContext struct {
    Args     []string
    Registry *registry.Registry
    Logger   *logger.Logger
    Config   *config.Config
}
```

## Core Commands

### Project Creation Commands

#### `create <framework>:<version> [project-name]`

**Purpose**: Scaffold new framework projects with Docker and AI integration

**Parameters**:
- `framework`: Framework name (laravel, django)
- `version`: Framework version (e.g., 11, 5.0)
- `project-name`: Optional project directory name

**Flags**:
- `--ai-manifest`: Generate AI project manifest
- `--interactive`: Use interactive setup prompts
- `--debug`: Enable debug output
- `--timeout <duration>`: Set operation timeout

**Examples**:
```bash
atempo create laravel:11                     # Create in current directory
atempo create laravel:11 my-app             # Create in ./my-app/
atempo create django:5.0 testing/api-test   # Create test project
atempo create laravel:11 --ai-manifest      # With AI features
```

**Return Codes**:
- `0`: Success
- `1`: General error
- `2`: Template not found
- `3`: Docker error
- `4`: Permission error

**Error Conditions**:
- Framework not supported
- Version not available
- Directory already exists
- Docker not available
- Insufficient permissions

---

### Project Management Commands

#### `projects`

**Purpose**: List all registered projects with status

**Parameters**: None

**Output Format**:
```
Project Name          Framework    Status    Path
────────────────────────────────────────────────
my-app               laravel      running   /path/to/my-app
api-project          django       stopped   /path/to/api-project
```

**Return Codes**:
- `0`: Success
- `1`: Registry read error

---

#### `status [project-name]`

**Purpose**: Show detailed project status and health

**Parameters**:
- `project-name`: Optional specific project name

**Output Sections**:
- Project Information
- Docker Services Status
- Port Mappings
- Git Information
- Health Checks

**Example Output**:
```
Project Status: my-app
──────────────────────
Framework: Laravel 11
Status: running
Path: /Users/user/projects/my-app

Services:
 my-app-redis           running
 my-app-mysql           running  
 my-app-app             running
 my-app-webserver       running

URLs:
 Application: http://localhost:8000
 Mailhog: http://localhost:8025

Git:
 Branch: main
 Status: clean
```

---

#### `describe [project-name]`

**Purpose**: Show comprehensive project details and configuration

**Parameters**:
- `project-name`: Optional project name (uses current directory if omitted)

**Output Includes**:
- Project metadata
- Service configurations
- Environment details
- File structure overview

---

#### `remove [project-name]`

**Purpose**: Remove project from registry (does not delete files)

**Parameters**:
- `project-name`: Project to remove from registry

**Confirmation**: Requires user confirmation

**Return Codes**:
- `0`: Success
- `1`: Project not found
- `2`: User cancelled

---

### Docker Operations Commands

#### `docker <command> [args...]`

**Purpose**: Execute Docker Compose commands with project context

**Common Commands**:
- `up [-d]`: Start services
- `down`: Stop services
- `logs [service]`: View logs
- `ps`: List containers
- `exec <service> <command>`: Execute in container
- `build [service]`: Build services
- `pull [service]`: Pull images

**Global Flags**:
- `--timeout <duration>`: Command timeout (default: 3m)
- `--project-path <path>`: Explicit project path

**Examples**:
```bash
atempo docker up -d                    # Start all services
atempo docker logs app                 # View app logs
atempo docker exec app bash            # Shell into app container
atempo docker up --timeout 300        # Custom timeout (5 minutes)
```

**Timeout Handling**:
- Default timeout: 3 minutes
- Configurable per command
- Graceful cancellation on timeout

---

### Interactive Shell Commands

#### `shell` or `atempo` (no args)

**Purpose**: Enter interactive shell mode

**Features**:
- Command history and completion
- Git-aware prompt
- Project context awareness
- Bash command passthrough

**Prompt Format**:
```
atempo (project-name) [git-branch] $
```

**Special Commands**:
- `exit`: Leave shell
- `help`: Show available commands
- `clear`: Clear screen

---

### Configuration Commands

#### `reconfigure [project-name]`

**Purpose**: Update existing project configuration

**Parameters**:
- `project-name`: Optional project name

**Interactive Process**:
1. Load current configuration
2. Prompt for changes
3. Update docker-compose.yml
4. Restart services if needed

---

#### `add-service <service-name> [project-name]`

**Purpose**: Add new service to existing project

**Parameters**:
- `service-name`: Name of service to add
- `project-name`: Optional project name

**Supported Services**:
- `redis`: Redis cache
- `elasticsearch`: Search engine
- `minio`: S3-compatible storage
- `mailhog`: Email testing
- `postgres`: PostgreSQL database

---

## AI Integration Commands

### `ai-manifest [project-name]`

**Purpose**: Generate AI project manifest and context files

**Parameters**:
- `project-name`: Optional project name

**Generated Files**:
- `.ai/project-intent.md`: Project goals and features
- `.ai/architecture-hints.md`: Technical recommendations
- `.ai/context.json`: Structured project data

**Process**:
1. Analyze project structure
2. Extract framework features
3. Generate AI-ready documentation
4. Create Claude-compatible context

---

### `interactive-prompt`

**Purpose**: Guided project setup with AI assistance

**Flow**:
1. Framework selection
2. Feature requirements gathering
3. Complexity assessment
4. AI manifest generation
5. Project scaffolding

---

## Project-Specific Command Routing

When inside a project directory, commands can be prefixed with project name:

```bash
# Inside /path/to/my-app/
my-app docker up        # Equivalent to: atempo docker up
my-app status          # Equivalent to: atempo status
my-app logs            # Equivalent to: atempo docker logs
```

**Resolution Order**:
1. Check for project in current directory
2. Check registry for project by name
3. Prompt user if ambiguous

---

## Error Handling Patterns

### Standard Error Format

```go
type CommandError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Context string `json:"context,omitempty"`
    Cause   error  `json:"cause,omitempty"`
}
```

### Common Error Codes

- `1`: General command error
- `2`: Template/resource not found
- `3`: Docker operation failed
- `4`: Permission denied
- `5`: Invalid configuration
- `6`: Network/connectivity error
- `7`: Timeout exceeded
- `8`: User cancellation

### Error Recovery

Commands provide recovery suggestions:

```
✗ Error: Docker service failed to start
  Cause: Port 3306 already in use
  
  Suggestions:
  - Stop conflicting service: sudo lsof -ti:3306 | xargs kill
  - Use different port in docker-compose.yml
  - Run: atempo docker down && atempo docker up
```

---

## Configuration System

### Global Configuration

Location: `~/.atempo/config.json`

```json
{
  "default_timeout": "3m",
  "auto_start_services": true,
  "ui_theme": "default",
  "ai_features_enabled": true,
  "docker_compose_version": "v2"
}
```

### Project Configuration

Location: `<project>/atempo.json`

```json
{
  "name": "my-app",
  "framework": "laravel",
  "version": "11",
  "language": "php",
  "installer": {
    "type": "docker",
    "command": ["docker", "run", "--rm", "..."]
  },
  "services": {
    "app": {
      "type": "build",
      "dockerfile": "infra/docker/Dockerfile"
    }
  }
}
```

---

## Registry System

### Registry Location

- **File**: `~/.atempo/registry.json`
- **Format**: JSON with project metadata
- **Backup**: Automatic backup on modifications

### Project Registration

```go
type Project struct {
    Name        string                 `json:"name"`
    Path        string                 `json:"path"`
    Framework   string                 `json:"framework"`
    Version     string                 `json:"version"`
    Language    string                 `json:"language"`
    Services    map[string]Service     `json:"services"`
    Status      string                 `json:"status"`
    URLs        map[string]string      `json:"urls"`
    CreatedAt   time.Time             `json:"created_at"`
    UpdatedAt   time.Time             `json:"updated_at"`
    GitBranch   string                `json:"git_branch,omitempty"`
    GitStatus   string                `json:"git_status,omitempty"`
}
```

### Service Status

```go
type Service struct {
    Name       string            `json:"name"`
    Status     string            `json:"status"`
    Image      string            `json:"image"`
    Ports      map[string]string `json:"ports"`
    Health     string            `json:"health"`
    UpdatedAt  time.Time         `json:"updated_at"`
}
```

---

## Logging System

### Log Locations

- **Global Logs**: `~/.atempo/logs/atempo.log`
- **Project Logs**: `~/.atempo/logs/<project-name>_<timestamp>.log`
- **Command Logs**: `~/.atempo/logs/commands/`

### Log Levels

- `DEBUG`: Detailed execution information
- `INFO`: General operational messages
- `WARN`: Warning conditions
- `ERROR`: Error conditions
- `FATAL`: Critical failures

### Log Format

```
2024-07-19T10:30:45Z [INFO] [project:my-app] [command:docker] Starting services
2024-07-19T10:30:46Z [DEBUG] [project:my-app] [docker] Executing: docker-compose up -d
2024-07-19T10:30:50Z [INFO] [project:my-app] [docker] Services started successfully
```

---

## Extension Points

### Custom Commands

```go
type CustomCommand struct {
    BaseCommand
    handler func(*CommandContext) error
}

func (c *CustomCommand) Execute(ctx *CommandContext) error {
    return c.handler(ctx)
}
```

### Plugin Interface (Future)

```go
type Plugin interface {
    Name() string
    Version() string
    Commands() []Command
    Initialize(config *Config) error
    Cleanup() error
}
```

### Framework Extensions

```go
type FrameworkExtension interface {
    Name() string
    SupportedVersions() []string
    TemplateURL() string
    Validate(config *AtempoConfig) error
    PostInstall(project *Project) error
}
```

This API reference provides comprehensive documentation for all command interfaces, parameters, and behaviors in the Atempo CLI system.
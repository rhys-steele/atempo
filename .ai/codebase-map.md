# Atempo Codebase Map

## Directory Structure Overview

```
atempo/
├── cmd/atempo/main.go              # CLI entry point (66 lines)
├── internal/                       # Core application logic
│   ├── app/commands/              # Command layer
│   │   ├── command.go             # Base command interface (59 lines)
│   │   ├── registry.go            # Command registry & routing (506 lines)
│   │   ├── create.go              # Project scaffolding (224 lines)
│   │   ├── docker.go              # Docker operations (196 lines)
│   │   ├── shell.go               # Interactive shell (474 lines)
│   │   ├── projects.go            # Project listing (53 lines)
│   │   ├── status.go              # Project dashboard (149 lines)
│   │   ├── other.go               # Additional commands (441 lines)
│   │   ├── ai_manifest.go         # AI manifest generation (435 lines)
│   │   ├── interactive_prompt.go  # User input gathering (390 lines)
│   │   ├── progress_tracker.go    # Progress indicators
│   │   └── status_indicator.go    # Visual status indicators
│   ├── scaffold/                  # Scaffolding engine
│   │   └── scaffold.go            # Core scaffolding logic (847 lines)
│   ├── registry/                  # Project registry
│   │   └── registry.go            # Project tracking (484 lines)
│   ├── docker/                    # Docker integration
│   │   └── docker.go              # Docker Compose operations (367 lines)
│   ├── compose/                   # Compose generation
│   │   └── generator.go           # YAML generation (378 lines)
│   ├── logger/                    # Logging system
│   │   └── logger.go              # Project logging (335 lines)
│   └── utils/                     # Utilities
│       └── file.go                # File operations (150 lines)
├── templates/                     # Template system
│   ├── frameworks/                # Framework templates
│   │   ├── laravel/atempo.json    # Laravel config (83 lines)
│   │   └── django/atempo.json     # Django config (102 lines)
│   └── ai/                        # AI templates
│       └── manifests/             # AI manifest templates
├── go.mod                         # Go module (14 lines)
├── go.sum                         # Dependency checksums
└── CLAUDE.md                      # Project documentation
```

## File Descriptions & Responsibilities

### Entry Point
- **`cmd/atempo/main.go`**: Minimal CLI entry point that delegates to command registry
  - Handles no-args case (shell mode), help commands, and error handling
  - Single responsibility: route commands to registry

### Command Layer Architecture

#### Core Command System
- **`internal/app/commands/command.go`**: Base interfaces and structures
  - Defines `Command` interface with `Execute()` method
  - Provides `BaseCommand` struct for common functionality
  - Establishes command pattern for all CLI operations

- **`internal/app/commands/registry.go`**: Command registry and routing system
  - Central command discovery and execution
  - Project-specific command routing (`my-app up`)
  - VS Code integration and browser launching
  - Help system and usage display

#### Primary Commands
- **`internal/app/commands/create.go`**: Project scaffolding with AI enhancement
  - 4-step process: AI Planning, Template Loading, Framework Installation, AI Context
  - Progress tracking with real-time feedback
  - Integration with AI manifest generation

- **`internal/app/commands/docker.go`**: Docker operations with timeout support
  - Custom timeout configurations
  - Project path resolution
  - Exec command handling for container interaction

- **`internal/app/commands/shell.go`**: Interactive shell interface
  - readline support with command history
  - Tab completion and auto-suggestions
  - Git-aware prompt with branch display
  - Bash command passthrough

- **`internal/app/commands/status.go`**: Project dashboard with health monitoring
  - Docker service status checking
  - Port mappings and URL detection
  - Git branch and status information
  - Service health indicators

#### AI Integration Commands
- **`internal/app/commands/ai_manifest.go`**: AI-powered project intent generation
  - Framework analysis and feature extraction
  - Architecture hints and recommendations
  - Project intent file generation

- **`internal/app/commands/interactive_prompt.go`**: User input gathering
  - Multi-step questionnaire system
  - Complexity analysis and feature selection
  - Integration with AI manifest generation

#### Utility Commands
- **`internal/app/commands/other.go`**: Additional commands
  - `ReconfigureCommand`: Update existing projects
  - `AddServiceCommand`: Add new services to projects
  - `LogsCommand`: View project logs
  - `DescribeCommand`: Show project details
  - `RemoveCommand`: Remove projects from registry

### Core Business Logic

#### Scaffolding Engine
- **`internal/scaffold/scaffold.go`**: Core scaffolding business logic
  - Version validation and compatibility checking
  - Template processing with variable substitution
  - Post-installation setup and configuration
  - Integration with registry and logging systems

#### Project Registry
- **`internal/registry/registry.go`**: Project registry management
  - JSON-based project storage in `~/.atempo/registry.json`
  - Docker status monitoring and health checking
  - Git integration for branch tracking
  - URL detection and service mapping

#### Docker Integration
- **`internal/docker/docker.go`**: Docker Compose operations
  - Command execution with timeout support
  - Bake support detection
  - Framework detection from compose files
  - Service listing and status checking

#### Compose Generation
- **`internal/compose/generator.go`**: Generate docker-compose.yml
  - Convert `atempo.json` to Docker Compose format
  - Service definition and configuration
  - Predefined service templates (minio, elasticsearch, etc.)
  - Network and volume management

#### Logging System
- **`internal/logger/logger.go`**: Project-specific logging
  - Command output capture and timestamping
  - Step-by-step progress tracking
  - Log file management in `~/.atempo/logs/`
  - Quiet mode and verbose output options

### Template System

#### Framework Templates
- **`templates/frameworks/laravel/atempo.json`**: Laravel project configuration
  - Docker-based Composer installation
  - Full Laravel stack: app, nginx, mysql, redis, mailhog
  - Service dependencies and networking

- **`templates/frameworks/django/atempo.json`**: Django project configuration
  - Docker-based Django installation
  - Async stack: web, postgres, redis, celery worker/beat
  - Mail testing with mailhog

#### AI Templates
- **`templates/ai/manifests/`**: AI manifest templates
  - Project intent templates for AI tools
  - Structured context for Claude integration
  - MCP server configuration templates

### Utility Functions
- **`internal/utils/file.go`**: File operations and utilities
  - Recursive directory copying
  - File existence checking
  - Semantic version comparison
  - Path manipulation utilities

## Key Relationships & Data Flow

### Command Execution Flow
1. **CLI Entry**: `main.go` → `registry.go`
2. **Command Routing**: Registry identifies command and project context
3. **Command Execution**: Specific command implements business logic
4. **Business Logic**: Commands delegate to appropriate service layer

### Project Scaffolding Flow
1. **User Input**: Interactive prompts gather project requirements
2. **AI Planning**: Generate project intent and manifest files
3. **Template Selection**: Choose appropriate framework template
4. **Scaffolding**: Copy template files and run installers
5. **Registration**: Add project to registry with metadata
6. **Docker Setup**: Generate and start Docker Compose services

### Docker Management Flow
1. **Project Detection**: Registry provides project metadata
2. **Compose Generation**: Generate docker-compose.yml from atempo.json
3. **Docker Operations**: Execute Docker Compose commands
4. **Status Monitoring**: Check service health and update registry

### AI Integration Flow
1. **Intent Gathering**: Interactive prompts collect requirements
2. **Manifest Generation**: AI processes intent into structured data
3. **Context Creation**: Generate AI-ready project context
4. **File Generation**: Create manifest and configuration files

## Architecture Patterns

### Design Patterns Used
- **Command Pattern**: All CLI commands implement unified interface
- **Registry Pattern**: Central project registry with metadata
- **Template Pattern**: Framework-agnostic scaffolding system
- **Factory Pattern**: Command creation and registration
- **Observer Pattern**: Progress tracking and status updates

### Architectural Principles
- **Clean Architecture**: Clear separation of concerns
- **Dependency Inversion**: Business logic doesn't depend on infrastructure
- **Single Responsibility**: Each file has focused responsibility
- **Open/Closed**: Easy to extend with new commands and frameworks

## Key Interfaces & Types

### Command Interface
```go
type Command interface {
    Execute(ctx *CommandContext) error
    GetDescription() string
}
```

### Project Registry Types
```go
type Project struct {
    Name        string
    Path        string
    Framework   string
    Services    map[string]Service
    Status      string
    // ...
}
```

### Template Configuration
```go
type AtempoConfig struct {
    Name      string
    Framework string
    Language  string
    Installer InstallerConfig
    Services  map[string]ServiceConfig
    // ...
}
```

## External Dependencies

### Runtime Dependencies
- **Docker**: Required for project scaffolding and management
- **Docker Compose**: Service orchestration and management
- **Git**: Version control integration and project tracking

### Go Dependencies
- **`github.com/chzyer/readline`**: Interactive shell functionality
- **`gopkg.in/yaml.v3`**: YAML processing for Docker Compose
- **Standard Library**: Primarily uses Go standard library

## Development & Extension Points

### Adding New Frameworks
1. Create template directory in `templates/frameworks/`
2. Define `atempo.json` with framework configuration
3. Include Docker setup files and framework-specific templates
4. Test scaffolding process

### Adding New Commands
1. Create command file in `internal/app/commands/`
2. Implement `Command` interface
3. Register command in `registry.go`
4. Add command-specific business logic

### Extending AI Features
1. Enhance prompt templates in `templates/ai/`
2. Extend manifest generation in `ai_manifest.go`
3. Add new interactive prompts in `interactive_prompt.go`
4. Integrate with MCP server discovery

This codebase demonstrates sophisticated software architecture with clear separation of concerns, making it maintainable and extensible for future development.
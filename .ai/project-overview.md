# Atempo - AI-Enhanced CLI Project Overview

## Project Mission
Atempo is a developer-first, AI-enhanced CLI tool for scaffolding framework-agnostic projects with built-in support for Claude, MCP (Model-Context-Protocol), Docker, and best-practice architecture. It operates external to projects while bootstrapping them with AI-ready context systems.

## Core Value Proposition
- **AI-First Development**: Built-in Claude and MCP integration for enhanced developer workflow
- **Framework Agnostic**: Supports multiple frameworks (Laravel, Django, and expanding)
- **Docker-Native**: All scaffolding uses Docker containers for consistency
- **Template-Driven**: Flexible template system with variable substitution
- **Developer Experience**: Interactive shell, progress tracking, and intuitive commands

## Architecture Overview

### Clean Architecture Pattern
The project follows Clean Architecture principles with clear separation of concerns:

```
cmd/atempo/main.go              # Entry point (65 lines)
├── internal/app/commands/      # Command Layer
│   ├── registry.go            # Command registry and routing
│   ├── create.go              # Project scaffolding
│   ├── docker.go              # Docker operations
│   ├── shell.go               # Interactive shell
│   └── [other commands]
├── internal/scaffold/          # Business Logic Layer
├── internal/registry/          # Project registry management
├── internal/compose/           # Docker Compose generation
├── internal/mcp/               # MCP server integration
├── internal/auth/              # Authentication providers
└── templates/                  # Template System
    ├── frameworks/            # Framework-specific templates
    └── ai/                    # AI-specific templates
```

### Key Components

#### 1. Command System
- **Modular Architecture**: Each command is a separate, testable unit
- **Registry Pattern**: Central command routing and discovery
- **Interactive Shell**: Real-time command execution with progress tracking
- **Project-Aware**: Commands adapt based on current project context

#### 2. Template System
- **Metadata-Driven**: Each template has `atempo.json` configuration
- **Variable Substitution**: `{{project}}`, `{{cwd}}`, `{{version}}`
- **Docker-First**: All installers use Docker containers
- **Service Definition**: Complex multi-service project support

#### 3. MCP Integration
- **Server Discovery**: Automatic MCP server detection and installation
- **Framework-Specific**: Custom MCP servers for each framework
- **Claude-Ready**: Built-in Claude integration capabilities
- **Context Management**: Automatic project context generation

#### 4. Docker Integration
- **Compose Generation**: Automatic `docker-compose.yml` from templates
- **Service Management**: Start/stop/status of multi-service projects
- **Health Monitoring**: Service health checking and status reporting
- **Development Workflow**: Streamlined Docker-based development

#### 5. DNS Management System
- **Local Domain Resolution**: DNSmasq-powered custom domains (`.local` TLD)
- **Service Subdomains**: Automatic subdomain generation for multi-service projects
- **Nginx Reverse Proxy**: Clean URL routing instead of port numbers
- **Port Management**: Dynamic port allocation with conflict resolution
- **macOS Integration**: System-level DNS resolver configuration

## Current State (Feature/AI-Prompting-Manifests Branch)

### ✅ Completed Features
- Core CLI architecture with command registry
- Template system with Laravel and Django support
- Docker Compose generation and management
- Project registry with status tracking
- Interactive shell mode with progress tracking
- MCP server discovery framework
- Basic scaffolding workflow
- DNS management system with DNSmasq integration
- Nginx reverse proxy for clean URL routing
- Dynamic port allocation and conflict resolution

### 🔄 In Development (Current Branch)
- AI manifest generation system
- Interactive prompting for project setup
- Authentication system for AI features
- Advanced MCP server templates
- Claude integration workflows

### 📋 Planned Features
- Additional framework support (Node.js, React, Astro)
- Framework command passthrough (`atempo artisan migrate`)
- AI context editing (`atempo context edit`)
- Direct Claude integration (`atempo claude "prompt"`)
- Plugin architecture for custom frameworks

## Technical Highlights

### Go Implementation
- **Minimal Dependencies**: Uses standard library with yaml.v3 and readline only
- **Clean Code**: Follows Go best practices with comprehensive error handling
- **Modular Design**: Easy to extend with new commands and frameworks
- **Testing-Ready**: Architecture supports unit testing of individual components

### AI Features
- **MCP Protocol**: Native Model-Context-Protocol support
- **Claude Integration**: Built-in Claude API integration
- **Context Generation**: Automatic project context for AI systems
- **Interactive Prompts**: AI-powered project setup workflows

### Developer Experience
- **Progress Tracking**: Real-time visual progress indicators
- **Interactive Shell**: Persistent shell session for project management
- **Status Dashboard**: Comprehensive project health monitoring
- **Error Handling**: Descriptive error messages and recovery suggestions

## Key Files for LLM Understanding

### Entry Points
- `cmd/atempo/main.go`: CLI entry point
- `internal/app/commands/registry.go`: Command system core

### Core Business Logic
- `internal/scaffold/scaffold.go`: Project scaffolding engine
- `internal/registry/registry.go`: Project registry management
- `internal/compose/generator.go`: Docker Compose generation
- `internal/docker/dns_manager.go`: DNS management and DNSmasq integration
- `internal/docker/nginx_proxy.go`: Nginx reverse proxy management
- `internal/docker/port_manager.go`: Dynamic port allocation system

### Template System
- `templates/frameworks/*/atempo.json`: Framework metadata
- `templates/ai/manifests/`: AI manifest templates

### Current Development Focus
- `internal/app/commands/ai_manifest.go`: AI manifest generation
- `internal/app/commands/interactive_prompt.go`: Interactive AI prompts
- `internal/app/commands/auth_checker.go`: Authentication for AI features

## Development Workflows

### Building and Testing
```bash
go build -o atempo cmd/atempo/main.go
./atempo create laravel:11
./atempo docker up
./atempo shell
```

### Module Management
```bash
go mod tidy
go mod download
```

### Adding New Frameworks
1. Create `templates/frameworks/[framework]/` directory
2. Add `atempo.json` with framework metadata
3. Include Docker setup files
4. Test with `atempo create [framework]:[version]`

### Adding New Commands
1. Create command file in `internal/app/commands/`
2. Implement `Command` interface
3. Register in `registry.go`
4. Add command-specific logic

This overview provides the foundational understanding needed to work effectively with the Atempo codebase, its architecture, and current development direction.
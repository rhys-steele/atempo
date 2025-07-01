# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Atempo is a developer-first, AI-enhanced CLI tool for scaffolding framework-agnostic projects with built-in support for Claude, MCP (Model-Context-Protocol), Docker, and best-practice architecture. It's designed as a reusable CLI that operates external to projects while bootstrapping them with AI-ready context systems.

## Build and Development Commands

### Building the CLI
```bash
go build -o atempo cmd/atempo/main.go
```

### Installation from source
```bash
go build -o atempo .
mv atempo /usr/local/bin/  # Optional: add to PATH
```

### Module management
```bash
go mod tidy
```

## Architecture

### Core Components

- **`cmd/atempo/main.go`**: CLI entry point that handles command parsing and delegates to scaffold package
- **`internal/scaffold/scaffold.go`**: Core scaffolding logic that reads template metadata and executes installation commands
- **`templates/`**: Framework-specific templates with metadata files (`atempo.json`)
- **`internal/context/`**: AI context management (currently empty, placeholder for future MCP integration)
- **`internal/utils/`**: Shared utilities (currently empty)

### Template System

Templates are stored in `templates/<framework>/` directories, each containing:
- `atempo.json`: Metadata file defining framework, language, installer commands, and working directory
- Framework-specific files (Dockerfile, docker-compose.yml, context files, etc.)

The `atempo.json` structure:
```json
{
  "framework": "laravel",
  "language": "php", 
  "installer": {
    "type": "composer",
    "command": ["composer", "create-project", "..."],
    "work-dir": "{{cwd}}"
  },
  "working-dir": "/var/www",
  "min-version": "10.0"
}
```

Template variables supported:
- `{{name}}`: Target directory name (defaults to "src")
- `{{cwd}}`: Current working directory 
- `{{project}}`: Project name (basename of current directory)

### Current CLI Usage
```bash
atempo start <framework>:<version>
```

Example: `atempo start laravel:12`

### Planned Features (from README)
- Docker integration (`atempo docker up`, `atempo docker bash`)
- Framework command passthrough (`atempo artisan migrate:fresh`)
- AI context editing (`atempo context edit`)
- Direct Claude integration (`atempo claude "prompt"`)
- Additional framework support (Node.js, Django, React, Astro)

## Development Notes

- This is a Go CLI application using standard library only (no external dependencies in go.mod)
- Templates use JSON for metadata configuration
- The scaffold system uses `os/exec` to run installer commands (composer, docker, etc.)
- Error handling returns descriptive messages to help users understand scaffold failures
- The codebase follows Go conventions with clear package separation
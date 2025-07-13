# Atempo Project AI Context System

## Overview
The `.ai` directory contains structured context files designed to provide comprehensive project understanding for AI assistants working on the Atempo CLI tool. This system ensures consistent, informed development aligned with project standards and goals.

## Directory Structure & File Index

### `.ai/context.md` (This File)
**Purpose**: Master context file serving as an index to all AI context files
**Format**: Markdown with structured sections
**Contents**:
- Complete directory structure and file purpose explanations
- Development standards and architectural guidelines
- Current priorities and security requirements
- Cross-references to all other context files

### `.ai/project-overview.md`
**Purpose**: High-level project mission, architecture, and feature overview
**Format**: Markdown with architecture diagrams and feature lists
**Contents**:
- Project mission and value proposition
- Clean Architecture implementation details
- Current development state and feature status
- Key components (Command System, Template System, MCP Integration, Docker Integration, DNS Management)
- Technical highlights and development workflows

### `.ai/codebase-map.md`
**Purpose**: Detailed technical documentation of the entire codebase structure
**Format**: Markdown with code examples and file relationships
**Contents**:
- Complete directory structure with line counts
- File descriptions and responsibilities for each component
- Architecture patterns and design principles
- Key interfaces, types, and data flow diagrams
- External dependencies and development/extension points

### `.ai/development-workflows.md`
**Purpose**: Comprehensive guide to development processes and CLI commands
**Format**: Markdown with code examples and workflow instructions
**Contents**:
- Quick start and build instructions
- Core development workflows (testing, framework development, command development)
- Complete CLI command reference with examples
- Testing, debugging, and troubleshooting procedures
- Git workflow, deployment, and performance optimization

### `.ai/patterns-and-conventions.md`
**Purpose**: Detailed coding standards, patterns, and architectural conventions
**Format**: Markdown with Go code examples and pattern explanations
**Contents**:
- Code architecture patterns (Clean Architecture, Command Pattern, Registry Pattern)
- Naming conventions and Go best practices
- Error handling patterns and configuration structures
- Logging, testing, and documentation patterns
- Security, performance, and extension patterns

### `.ai/ui-ux-guidelines.md`
**Purpose**: Strict UI/UX standards for CLI interface design
**Format**: Markdown with examples and code snippets
**Contents**:
- Core design principles for modern, elegant CLI interface
- Emoji usage policy (minimal, professional only)
- Status indicators and output formatting standards
- Text formatting conventions and professional language guidelines
- Good/bad examples with visual comparisons

## AI Context Provision Strategy

This structure is optimized for AI assistants by providing:

1. **Hierarchical Information**: Master context file references specialized guidelines
2. **Specific Standards**: Clear, actionable rules rather than vague suggestions
3. **Examples**: Concrete good/bad examples in UI guidelines
4. **Cross-References**: Links between related context files
5. **Structured Format**: Consistent markdown formatting for easy parsing

## Development Standards

### Code Quality
- Follow Go best practices and conventions
- Clean Architecture with clear separation of concerns
- Comprehensive error handling with descriptive messages
- No external dependencies beyond Go standard library

### User Experience Requirements
- **CRITICAL**: Strictly follow UI/UX guidelines in `ui-ux-guidelines.md`
- Minimal emoji usage - professional CLI interface only
- Clean, modern output with consistent formatting
- Focus on functionality over decoration

### Testing Standards
- Test all changes thoroughly before completion
- Verify DNS resolution actually works end-to-end
- Test error cases and edge conditions
- Ensure backwards compatibility

### Architecture Guidelines
- Commands are modular and isolated in `internal/app/commands/`
- Shared utilities in `internal/utils/` prevent code duplication
- DNS system supports both Docker and local modes
- Dynamic port allocation prevents conflicts
- Graceful fallbacks for missing dependencies

## Current Development Priorities
1. Professional CLI interface (strictly enforce ui-ux-guidelines.md)
2. Reliable DNS resolution and comprehensive testing
3. Clear error messages and actionable user guidance
4. Modern, elegant user experience without excessive decoration
5. **CRITICAL: Always use `/testing` directory for test projects**

## **CRITICAL: Testing Directory Convention**
**ALL testing and experimental projects MUST be created in the `/testing` directory!**

### Requirements:
- **ALWAYS use**: `atempo create <framework> testing/<project-name>`
- **Examples**: `atempo create laravel testing/my-test-app`, `atempo create django testing/api-test`
- **Never create test projects in the root directory**
- **This applies to ALL testing, debugging, and experimental work**

## Security Requirements
- Never log or expose sensitive information
- Validate all user inputs thoroughly
- Use secure defaults for all configurations
- Require explicit confirmation for destructive actions
- Follow defensive programming practices

## Project Context
**Atempo** is a developer-first, AI-enhanced CLI tool for scaffolding framework-agnostic projects with built-in support for Claude, MCP (Model-Context-Protocol), Docker, and best-practice architecture. It operates external to projects while bootstrapping them with AI-ready context systems.

**Key Technologies**: Go (standard library only), Docker, Docker Compose, JSON templates
**Architecture**: Clean Architecture with command-based modular design
**Target Users**: Developers seeking rapid, AI-enhanced project scaffolding
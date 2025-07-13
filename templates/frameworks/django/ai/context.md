# {{project}} Project AI Context System

## Overview
The `.ai` directory contains structured context files designed to provide comprehensive project understanding for AI assistants working on the {{project}} Django application. This system ensures consistent, informed development aligned with project standards and goals.

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
- Django application architecture details
- Current development state and feature status
- Key components and technical highlights
- Development workflows and deployment processes

### `.ai/codebase-map.md`
**Purpose**: Detailed technical documentation of the entire codebase structure
**Format**: Markdown with code examples and file relationships
**Contents**:
- Complete directory structure with descriptions
- File descriptions and responsibilities for each component
- Architecture patterns and design principles
- Key interfaces, models, and data flow diagrams
- External dependencies and integration points

### `.ai/development-workflows.md`
**Purpose**: Comprehensive guide to development processes and commands
**Format**: Markdown with code examples and workflow instructions
**Contents**:
- Quick start and setup instructions
- Core development workflows (testing, database, deployment)
- Command reference for Django and Docker operations
- Testing, debugging, and troubleshooting procedures
- Git workflow and deployment procedures

### `.ai/patterns-and-conventions.md`
**Purpose**: Detailed coding standards, patterns, and architectural conventions
**Format**: Markdown with Python code examples and pattern explanations
**Contents**:
- Code architecture patterns and best practices
- Naming conventions and Python coding standards
- Error handling patterns and validation approaches
- Testing patterns and documentation standards
- Security, performance, and maintainability guidelines

### `.ai/ui-ux-guidelines.md`
**Purpose**: UI/UX standards and design guidelines for the application
**Format**: Markdown with examples and design principles
**Contents**:
- Design system and component guidelines
- User experience principles and accessibility standards
- Frontend architecture and styling conventions
- API design and response formatting standards
- User interface patterns and interaction guidelines

## AI Context Provision Strategy

This structure is optimized for AI assistants by providing:

1. **Hierarchical Information**: Master context file references specialized guidelines
2. **Specific Standards**: Clear, actionable rules rather than vague suggestions
3. **Examples**: Concrete code examples and patterns
4. **Cross-References**: Links between related context files
5. **Structured Format**: Consistent markdown formatting for easy parsing

## Development Standards

### Code Quality
- Follow Python and Django best practices and conventions
- Implement proper error handling and validation
- Write comprehensive tests for all features
- Use Django's built-in patterns and avoid reinventing wheels

### Security Requirements
- Never log or expose sensitive information
- Validate all user inputs thoroughly
- Use secure defaults for all configurations
- Implement proper authentication and authorization
- Follow Django security guidelines and OWASP recommendations

### Testing Standards
- Write unit tests for all models and business logic
- Implement integration tests for views and API endpoints
- Test database migrations and data integrity
- Verify security and validation rules
- Test error cases and edge conditions

### Architecture Guidelines
- Follow Django conventions and best practices
- Implement clean separation of concerns with Django apps
- Use Django's service layer patterns
- Implement proper caching strategies with Redis
- Design for scalability and maintainability

## Current Development Priorities
1. Core application functionality and features
2. Comprehensive testing and quality assurance
3. Performance optimization and scalability
4. Security hardening and vulnerability assessment
5. Documentation and developer experience

## Project Context
**{{project}}** is a Django application scaffolded with Atempo, featuring built-in Docker development environment, AI-ready context system, and best-practice architecture.

**Key Technologies**: Django, Python, Docker, PostgreSQL, Redis, Celery
**Architecture**: Clean Django architecture with app-based modular design
**Development Environment**: Docker-based with hot reload and debugging support
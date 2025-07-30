# {{project}} Project AI Context System

## Overview
The `.ai` directory contains structured context files designed to provide comprehensive project understanding for AI assistants working on the {{project}} Express.js application. This system ensures consistent, informed development aligned with project standards and goals.

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
- Express.js application architecture details
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
- Command reference for Express.js and Docker operations
- Testing, debugging, and troubleshooting procedures
- Git workflow and deployment procedures

### `.ai/patterns-and-conventions.md`
**Purpose**: Detailed coding standards, patterns, and architectural conventions
**Format**: Markdown with JavaScript code examples and pattern explanations
**Contents**:
- Code architecture patterns (MVC, Service Layer, Repository)
- Naming conventions and JavaScript best practices
- Error handling patterns and middleware structures
- Testing, logging, and documentation patterns
- Security, performance, and API design patterns

### `.ai/ui-ux-guidelines.md`
**Purpose**: API design and response formatting standards
**Format**: Markdown with examples and code snippets
**Contents**:
- REST API design principles
- JSON response formatting standards
- Error response conventions
- API documentation guidelines
- Client integration examples

## Development Standards

### Code Quality
- Follow JavaScript/Node.js best practices and conventions
- Use ESLint and Prettier for code formatting
- Implement comprehensive error handling
- Use async/await for asynchronous operations

### Security Requirements
- **CRITICAL**: Implement helmet for security headers
- Use environment variables for sensitive configuration
- Implement rate limiting for API endpoints
- Validate all input data thoroughly
- Use HTTPS in production environments

### Testing Standards
- Write unit tests for business logic
- Write integration tests for API endpoints
- Use Jest as the testing framework
- Aim for good test coverage
- Test error scenarios and edge cases

### Architecture Guidelines
- Follow Express.js middleware pattern
- Use modular route organization with Express Router
- Implement service layer for business logic
- Use middleware for cross-cutting concerns
- Keep controllers thin and focused

## Current Development Priorities
1. RESTful API design and implementation
2. Proper error handling and logging
3. Database integration (MongoDB/PostgreSQL)
4. Authentication and authorization
5. API documentation and testing

## Security Requirements
- Never expose sensitive information in responses
- Validate all user inputs thoroughly
- Use secure defaults for all configurations
- Implement proper authentication and authorization
- Follow OWASP security guidelines

## Project Context
**{{project}}** is an Express.js application created with Atempo, designed for building robust REST APIs and web services. It follows modern Node.js development practices with Docker containerization and comprehensive tooling.

**Key Technologies**: Node.js, Express.js, MongoDB/Redis, Docker, Jest
**Architecture**: RESTful API with middleware-based architecture
**Target Use Cases**: REST APIs, microservices, web backends, GraphQL APIs
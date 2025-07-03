# Development Workflows & Commands

## Quick Start Guide

### Building the CLI
```bash
# Build the binary
go build -o atempo cmd/atempo/main.go

# Alternative: Build with go build (uses go.mod)
go build -o atempo .

# Install to system PATH (optional)
mv atempo /usr/local/bin/
```

### Module Management
```bash
# Keep dependencies clean
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

## Core Development Workflows

### 1. Testing the CLI Locally
```bash
# Build and test basic functionality
go build -o atempo cmd/atempo/main.go
./atempo --help
./atempo

# Test project creation
./atempo create laravel:11
./atempo create django:5.0

# Test project management
./atempo projects
./atempo status
./atempo docker up
```

### 2. Framework Template Development
```bash
# Create new framework template
mkdir -p templates/frameworks/new-framework
cd templates/frameworks/new-framework

# Create template configuration
cat > atempo.json << 'EOF'
{
  "name": "{{project}}",
  "framework": "new-framework",
  "language": "javascript",
  "installer": {
    "type": "docker",
    "command": ["docker", "run", "--rm", "-v", "{{cwd}}:/workspace", "node:18", "npm", "init", "-y"]
  },
  "services": {
    "app": {
      "type": "build",
      "dockerfile": "Dockerfile"
    }
  }
}
EOF

# Test the new template
cd ../../../
./atempo create new-framework:1.0
```

### 3. Command Development Workflow
```bash
# Create new command file
touch internal/app/commands/new_command.go

# Implement command interface
# Register in registry.go
# Test the command
./atempo new-command --help
```

### 4. AI Feature Development
```bash
# Test AI manifest generation
./atempo create laravel:11 --ai-manifest

# Test interactive prompts
./atempo create --interactive

# Debug AI features
./atempo create laravel:11 --debug
```

## Available CLI Commands

### Core Commands
```bash
# Project scaffolding
atempo create <framework>:<version>      # Create new project
atempo create laravel:11                 # Create Laravel 11 project
atempo create django:5.0                 # Create Django 5.0 project

# Project management
atempo projects                          # List all projects
atempo status                            # Show current project status
atempo describe                          # Show project details
atempo remove                            # Remove project from registry

# Docker operations
atempo docker <command>                  # Run Docker Compose commands
atempo docker up                         # Start services
atempo docker down                       # Stop services
atempo docker logs                       # View logs
atempo docker ps                         # List containers
atempo docker exec <service> <command>   # Execute in container

# Interactive shell
atempo shell                             # Enter interactive shell
atempo                                   # Enter interactive shell (no args)
```

### Project-Specific Commands
```bash
# When inside a project directory
my-app docker up                         # Start this project's services
my-app status                            # Show this project's status
my-app logs                              # View this project's logs
my-app docker exec app bash              # Execute in app container
```

### Configuration Commands
```bash
# Project configuration
atempo reconfigure                       # Reconfigure existing project
atempo add-service <service>             # Add new service to project
```

### AI Commands (In Development)
```bash
# AI-powered features
atempo create <framework> --ai-manifest  # Generate AI manifest
atempo create --interactive              # Interactive project setup
atempo context edit                      # Edit AI context (planned)
atempo claude "prompt"                   # Direct Claude integration (planned)
```

## Development Environment Setup

### Prerequisites
```bash
# Required tools
- Go 1.23+ (with toolchain go1.24.4)
- Docker & Docker Compose
- Git

# Optional tools for development
- VS Code with Go extension
- golangci-lint for code quality
- air for live reloading during development
```

### Setting Up Development Environment
```bash
# Clone the repository
git clone <repository-url>
cd atempo

# Install dependencies
go mod download

# Build and test
go build -o atempo cmd/atempo/main.go
./atempo --help

# Run in development mode
go run cmd/atempo/main.go --help
```

## Testing Workflows

### Manual Testing
```bash
# Test complete workflow
./atempo create laravel:11               # Create project
cd laravel-project
../atempo docker up                      # Start services
../atempo status                         # Check status
../atempo docker logs                    # View logs
../atempo docker down                    # Stop services
```

### Testing New Features
```bash
# Test AI features
./atempo create laravel:11 --debug       # Debug mode
./atempo create --interactive            # Interactive mode

# Test Docker integration
./atempo docker up --timeout 300         # Custom timeout
./atempo docker exec app bash           # Container access
```

### Template Testing
```bash
# Test framework templates
./atempo create laravel:11
./atempo create django:5.0

# Verify generated files
ls -la <project-directory>/
docker-compose -f <project-directory>/docker-compose.yml config
```

## Debugging & Troubleshooting

### Debug Mode
```bash
# Enable debug output
./atempo create laravel:11 --debug

# Check logs
ls ~/.atempo/logs/
cat ~/.atempo/logs/<project-name>_*.log
```

### Common Issues
```bash
# Docker issues
docker --version                         # Check Docker
docker-compose --version                # Check Compose

# Permission issues
chmod +x atempo                         # Make executable
sudo mv atempo /usr/local/bin/          # Install globally

# Registry issues
cat ~/.atempo/registry.json             # Check registry
rm ~/.atempo/registry.json              # Reset registry
```

### Project Registry Management
```bash
# Registry location
~/.atempo/registry.json

# Reset registry
rm ~/.atempo/registry.json

# Manually edit registry
code ~/.atempo/registry.json
```

## Code Quality & Standards

### Go Best Practices
```bash
# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run

# Vet code
go vet ./...

# Run tests
go test ./...
```

### Code Structure Guidelines
- Follow Clean Architecture principles
- Keep commands focused and testable
- Use interfaces for testability
- Handle errors gracefully with descriptive messages
- Log important operations for debugging

## Git Workflow

### Branch Strategy
```bash
# Current development branch
git checkout feature/ai-prompting-manifests

# Create feature branch
git checkout -b feature/new-feature

# Commit changes
git add .
git commit -m "feat: implement new feature"

# Push and create PR
git push origin feature/new-feature
```

### Commit Message Format
```bash
# Use conventional commits
feat: add new command
fix: resolve Docker timeout issue
docs: update README
refactor: improve error handling
test: add unit tests for scaffold
```

## Deployment & Distribution

### Building for Distribution
```bash
# Build for current platform
go build -o atempo cmd/atempo/main.go

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o atempo-linux-amd64 cmd/atempo/main.go
GOOS=darwin GOARCH=amd64 go build -o atempo-darwin-amd64 cmd/atempo/main.go
GOOS=windows GOARCH=amd64 go build -o atempo-windows-amd64.exe cmd/atempo/main.go
```

### Installation Methods
```bash
# From source
go install ./cmd/atempo

# Manual installation
go build -o atempo cmd/atempo/main.go
mv atempo /usr/local/bin/

# Development installation
go build -o atempo cmd/atempo/main.go
export PATH=$PATH:$(pwd)
```

## Performance & Optimization

### Monitoring Performance
```bash
# Profile the application
go build -o atempo cmd/atempo/main.go
./atempo create laravel:11 --profile

# Check memory usage
go build -race -o atempo cmd/atempo/main.go
./atempo create laravel:11
```

### Optimization Tips
- Use minimal dependencies (currently only readline and yaml)
- Implement timeout handling for Docker operations
- Use buffered I/O for large file operations
- Leverage Go's concurrency for parallel operations

## Integration Testing

### End-to-End Testing
```bash
# Test complete workflow
./test-e2e.sh laravel:11
./test-e2e.sh django:5.0

# Test Docker integration
./test-docker.sh
```

### Framework Integration Testing
```bash
# Test Laravel integration
./atempo create laravel:11
cd laravel-project
docker-compose up -d
curl http://localhost:8000
docker-compose down

# Test Django integration
./atempo create django:5.0
cd django-project
docker-compose up -d
curl http://localhost:8000
docker-compose down
```

This comprehensive guide covers all aspects of development workflows for the Atempo project, from basic setup to advanced debugging and deployment strategies.
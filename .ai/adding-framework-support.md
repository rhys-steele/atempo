# Adding Framework Support to Atempo

This document outlines the complete process for adding a new framework to Atempo, based on the experience of implementing Express.js support.

## Overview

Adding a new framework to Atempo requires creating template files, configuration, and Go code changes. The process involves several interconnected components that work together to provide a complete scaffolding experience.

## Required Components

### 1. Template Directory Structure

Create the framework template directory:

```
templates/frameworks/{framework}/
├── README.md                    # Framework-specific documentation
├── atempo.json                 # Core configuration file
├── {framework-files}           # Framework starter files (e.g., server.js)
├── ai/                         # AI context files
│   ├── ai-config.json
│   ├── context.md
│   ├── project-overview.md
│   ├── codebase-map.md
│   ├── development-workflows.md
│   ├── patterns-and-conventions.md
│   └── ui-ux-guidelines.md
├── claude/                     # Claude integration
│   └── commands/
│       └── setup.md
└── infra/                      # Infrastructure files
    └── docker/
        ├── Dockerfile
        ├── docker-compose.yml
        └── {config-files}      # Framework-specific configs
```

### 2. Core Configuration (atempo.json)

The `atempo.json` file defines how the framework is installed and configured:

```json
{
  "name": "{{project}}",
  "framework": "express",
  "language": "javascript",
  "installer": {
    "type": "docker",
    "command": [
      "docker", "run", "--rm", "-v", "{{cwd}}:/workspace",
      "-w", "/workspace/src", "node:18-alpine", "sh", "-c",
      "mkdir -p /workspace/src && cd /workspace/src && npm init -y && npm install {packages}"
    ],
    "work-dir": "{{cwd}}"
  },
  "working-dir": "/app",
  "min-version": "4.18.0",
  "services": {
    // Docker services configuration
  },
  "volumes": {
    // Docker volumes
  },
  "networks": {
    // Docker networks
  }
}
```

**Key considerations:**
- The installer creates the framework project structure
- Use `{{project}}`, `{{cwd}}`, `{{name}}`, `{{version}}` for templating
- Services define the Docker containerization
- Dynamic port allocation is handled automatically by Atempo

### 3. Docker Configuration

#### Dockerfile
- Must be flexible to work with or without package.json initially
- Should handle the case where dependencies aren't installed yet
- Use multi-stage approach if needed for production builds

#### docker-compose.yml Template
- This is a reference template - the actual docker-compose.yml is generated
- Port mappings use internal container ports
- Atempo handles external port allocation dynamically
- Include necessary services (database, cache, etc.)

### 4. AI Context Files

Each framework needs comprehensive AI context files:

#### ai-config.json
```json
{
  "framework": "express",
  "language": "JavaScript",
  "latest_version": "4.18.0",
  "ai_features": {
    "default_project_types": ["REST API", "Web Service"],
    "core_features": ["Express Router", "Middleware System"],
    "architecture_patterns": {
      "mvc_pattern": "Description of MVC in this framework"
    }
  },
  "development_context": {
    "package_manager": "npm",
    "structure": {
      "source_root": "src/",
      "routes_dir": "src/routes/"
    },
    "commands": {
      "install_dependencies": "npm install",
      "start_dev": "npm run dev"
    }
  }
}
```

#### Other AI Files
- `context.md`: Master index file
- `project-overview.md`: High-level architecture and mission
- `codebase-map.md`: Detailed technical documentation
- `development-workflows.md`: Development processes and commands
- `patterns-and-conventions.md`: Coding standards and patterns
- `ui-ux-guidelines.md`: Interface/API design guidelines

### 5. Go Code Changes

#### Add Framework Validation

In `internal/scaffold/scaffold.go`, add version validation:

```go
// In validateVersion function switch statement
case "express":
    return validateExpressVersion(requestedVersion)

// Add validation function
func validateExpressVersion(version string) error {
    majorVersion := utils.ParseVersionPart(strings.Split(version, ".")[0])
    
    if majorVersion < 4 {
        return fmt.Errorf("Express.js version %s is too old (minimum supported: 4.0)", version)
    }
    
    if majorVersion > 5 {
        return fmt.Errorf("Express.js version %s is not yet supported (maximum: 5.x)", version)
    }
    
    return nil
}
```

#### Add Post-Install Setup

In `internal/scaffold/scaffold.go`:

```go
// In runPostInstall function
if meta.Framework == "express" {
    return setupExpress(log, step, projectDir)
}

// Add setup function
func setupExpress(log *logger.Logger, step *logger.Step, projectDir string) error {
    srcDir := filepath.Join(projectDir, "src")
    
    // Copy template files from project root to src/
    templateFiles := map[string]string{
        "server.js":     filepath.Join(srcDir, "server.js"),
        ".env.example":  filepath.Join(srcDir, ".env.example"),
    }
    
    projectName := filepath.Base(projectDir)
    
    for templateFile, dstPath := range templateFiles {
        srcPath := filepath.Join(projectDir, templateFile)
        if utils.FileExists(srcPath) {
            // Process template variables
            content, err := os.ReadFile(srcPath)
            if err != nil {
                log.WarningStep(step, fmt.Sprintf("Failed to read %s: %v", templateFile, err))
                continue
            }
            
            processedContent := processTemplateContent(string(content), projectName, projectDir, "")
            
            if err := os.WriteFile(dstPath, []byte(processedContent), 0644); err != nil {
                log.WarningStep(step, fmt.Sprintf("Failed to copy %s: %v", templateFile, err))
                continue
            }
            
            // Remove template file from project root
            os.Remove(srcPath)
        }
    }
    
    // Additional setup (env files, package.json updates, etc.)
    
    return nil
}
```

#### Add Template File Copying

In the `copyTemplateFiles` function:

```go
// Copy framework-specific template files
if framework == "express" {
    expressFiles := map[string]string{
        "server.js":     "server.js",
        ".env.example":  ".env.example",
    }
    
    for templateFile, dstFile := range expressFiles {
        dstPath := filepath.Join(projectDir, dstFile)
        
        // Try embedded first, fallback to filesystem
        embeddedPath := fmt.Sprintf("templates/frameworks/%s/%s", framework, templateFile)
        if err := copyEmbeddedFileWithContext(templatesFS, embeddedPath, dstPath, projectName, projectDir, version); err != nil {
            // Fallback to filesystem
            srcPath, pathErr := getFilesystemTemplatePath(framework, templateFile)
            if pathErr == nil {
                if err := copyFilesystemFileWithContext(srcPath, dstPath, projectName, projectDir, version); err != nil {
                    log.WarningStep(step, fmt.Sprintf("Failed to copy %s: %v", templateFile, err))
                }
            }
        }
    }
}
```

## Key Insights and Gotchas

### 1. Template Variable Processing
- Use `{{project}}`, `{{cwd}}`, `{{name}}`, `{{version}}` in template files
- These are processed during file copying and post-install setup
- Template processing happens in multiple places - ensure consistency

### 2. Docker Integration
- **Don't hardcode ports** in atempo.json - use internal container ports
- Atempo's port manager handles dynamic external port allocation
- Volume mapping should match the installer's directory structure
- The installer creates files where Docker expects them

### 3. File Structure Flow
```
1. Installer runs → creates package.json, dependencies in src/
2. Template files get copied → server.js, .env.example to project root
3. Post-install runs → moves template files from root to src/, processes variables
4. Docker Compose generated → uses dynamic ports, proper volume mapping
```

### 4. Installer Command Structure
- Use Docker containers for consistent environments
- The installer should create the basic project structure
- Don't try to copy template files in the installer - let the scaffolding system handle it
- Work directory should match the atempo.json volume mapping

### 5. Testing Strategy
- Test with `--no-ai` flag first to avoid authentication requirements
- Check file placement after each step: installer, template copying, post-install
- Verify Docker services start correctly
- Test API endpoints if applicable

### 6. AI Context Importance
- AI context files are crucial for developer experience
- They guide AI assistants on how to work with the framework
- Include comprehensive examples and patterns
- Keep them consistent with other frameworks

## Common Pitfalls

1. **Volume Mapping Misalignment**: Ensure Docker volume mapping matches where the installer creates files
2. **Template Variable Scope**: Remember that template processing happens in multiple phases
3. **Port Configuration**: Never hardcode external ports - always use container ports
4. **File Permissions**: Ensure Docker containers can read/write files properly
5. **Dependency Installation**: The installer must complete successfully before post-install runs

## Testing Checklist

- [ ] Template directory structure is complete
- [ ] atempo.json configuration is valid
- [ ] Installer creates proper file structure
- [ ] Template files are copied to correct locations
- [ ] Post-install processing works correctly
- [ ] Docker services start successfully
- [ ] Application is accessible via dynamic ports
- [ ] AI context files are comprehensive and accurate
- [ ] Version validation works correctly

## Final Notes

Adding framework support to Atempo is a multi-step process that requires understanding the interaction between:
- Template system
- Go scaffolding code
- Docker containerization
- Dynamic port allocation
- AI context generation

The key is ensuring all these components work together cohesively to provide a seamless developer experience.
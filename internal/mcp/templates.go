package mcp

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var templateFiles embed.FS

// ServerTemplate represents an MCP server template
type ServerTemplate struct {
	PackageJSON string
	IndexJS     string
	README      string
}

// ProjectInfo contains information about the project for template generation
type ProjectInfo struct {
	Name      string
	Framework string
	Version   string
	Path      string
}

// GenerateServerFromTemplate creates an MCP server from templates
func GenerateServerFromTemplate(server MCPServer, projectInfo ProjectInfo, mcpDir string) error {
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		return err
	}

	template := getServerTemplate(server.Framework)

	// Generate package.json
	if err := generateFile(mcpDir, "package.json", template.PackageJSON, server, projectInfo); err != nil {
		return fmt.Errorf("failed to generate package.json: %w", err)
	}

	// Generate index.js
	if err := generateFile(mcpDir, "index.js", template.IndexJS, server, projectInfo); err != nil {
		return fmt.Errorf("failed to generate index.js: %w", err)
	}

	// Generate README.md
	if err := generateFile(mcpDir, "README.md", template.README, server, projectInfo); err != nil {
		return fmt.Errorf("failed to generate README.md: %w", err)
	}

	// Make index.js executable
	indexPath := filepath.Join(mcpDir, "index.js")
	return os.Chmod(indexPath, 0755)
}

// generateFile creates a file from a template
func generateFile(dir, filename, templateContent string, server MCPServer, projectInfo ProjectInfo) error {
	tmpl, err := template.New(filename).Parse(templateContent)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		Server  MCPServer
		Project ProjectInfo
	}{
		Server:  server,
		Project: projectInfo,
	}

	return tmpl.Execute(file, data)
}

// getServerTemplate returns the template for a specific framework
func getServerTemplate(framework string) ServerTemplate {
	switch framework {
	case "laravel":
		return getLaravelTemplate()
	case "django":
		return getDjangoTemplate()
	default:
		return getUniversalTemplate()
	}
}

// getLaravelTemplate returns Laravel-specific MCP server template
func getLaravelTemplate() ServerTemplate {
	return ServerTemplate{
		PackageJSON: readTemplateFile("templates/laravel/package.json"),
		IndexJS:     readTemplateFile("templates/laravel/index.js"),
		README:      readTemplateFile("templates/laravel/README.md"),
	}
}

// getDjangoTemplate returns Django-specific MCP server template
func getDjangoTemplate() ServerTemplate {
	return ServerTemplate{
		PackageJSON: readTemplateFile("templates/django/package.json"),
		IndexJS:     readTemplateFile("templates/django/index.js"),
		README:      readTemplateFile("templates/django/README.md"),
	}
}

// getUniversalTemplate returns a universal MCP server template
func getUniversalTemplate() ServerTemplate {
	return ServerTemplate{
		PackageJSON: readTemplateFile("templates/universal/package.json"),
		IndexJS:     readTemplateFile("templates/universal/index.js"),
		README:      readTemplateFile("templates/universal/README.md"),
	}
}

// readTemplateFile reads a template file from the embedded filesystem
func readTemplateFile(path string) string {
	content, err := templateFiles.ReadFile(path)
	if err != nil {
		// Fallback to basic template if file not found
		return getBasicTemplate(filepath.Base(path))
	}
	return string(content)
}

// getBasicTemplate provides fallback templates
func getBasicTemplate(filename string) string {
	switch filename {
	case "package.json":
		return `{
  "name": "{{.Server.Name}}",
  "version": "1.0.0",
  "description": "{{.Server.Description}}",
  "main": "index.js",
  "dependencies": {
    "@modelcontextprotocol/sdk": "^1.0.0"
  }
}`
	case "index.js":
		return `#!/usr/bin/env node
console.log("Basic MCP Server for {{.Project.Name}}");`
	case "README.md":
		return `# {{.Server.Name}}

Basic MCP server for {{.Project.Name}}.`
	default:
		return ""
	}
}

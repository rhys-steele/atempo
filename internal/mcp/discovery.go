package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // "official", "community", "generated"
	Framework   string            `json:"framework"`
	Repository  string            `json:"repository,omitempty"`
	NPMPackage  string            `json:"npm_package,omitempty"`
	Command     []string          `json:"command"`
	Args        []string          `json:"args,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Description string            `json:"description"`
	Version     string            `json:"version,omitempty"`
}

// MCPRegistry holds information about available MCP servers
type MCPRegistry struct {
	LastUpdated time.Time             `json:"last_updated"`
	Servers     map[string]MCPServer  `json:"servers"`
}

// DiscoveryResult contains the result of MCP server discovery
type DiscoveryResult struct {
	Official  []MCPServer `json:"official"`
	Community []MCPServer `json:"community"`
	Generated *MCPServer  `json:"generated,omitempty"`
}

// registryURLs contains known MCP server registries
var registryURLs = []string{
	"https://raw.githubusercontent.com/modelcontextprotocol/servers/main/registry.json",
	"https://api.github.com/search/repositories?q=mcp-server+topic:model-context-protocol&sort=stars",
}

// frameworkKeywords maps frameworks to search keywords
var frameworkKeywords = map[string][]string{
	"laravel": {"laravel", "php", "artisan", "eloquent", "blade"},
	"django":  {"django", "python", "manage.py", "orm", "templates"},
	"rails":   {"rails", "ruby", "activerecord", "erb"},
	"express": {"express", "nodejs", "npm", "javascript"},
}

// DiscoverMCPServers discovers available MCP servers for a framework
func DiscoverMCPServers(framework string) (*DiscoveryResult, error) {
	result := &DiscoveryResult{
		Official:  []MCPServer{},
		Community: []MCPServer{},
	}

	// Try to load cached registry
	registry, err := loadCachedRegistry()
	if err != nil || time.Since(registry.LastUpdated) > 24*time.Hour {
		// Refresh registry if old or missing
		registry, err = refreshRegistry()
		if err != nil {
			fmt.Printf("⚠️  Failed to refresh MCP registry: %v\n", err)
			// Continue with empty registry
			registry = &MCPRegistry{
				LastUpdated: time.Now(),
				Servers:     make(map[string]MCPServer),
			}
		}
	}

	// Search for framework-specific servers
	keywords := frameworkKeywords[framework]
	if keywords == nil {
		keywords = []string{framework}
	}

	for _, server := range registry.Servers {
		if matchesFramework(server, framework, keywords) {
			if server.Type == "official" {
				result.Official = append(result.Official, server)
			} else {
				result.Community = append(result.Community, server)
			}
		}
	}

	// If no servers found, generate a custom one
	if len(result.Official) == 0 && len(result.Community) == 0 {
		generated, err := generateCustomServer(framework)
		if err != nil {
			return nil, fmt.Errorf("failed to generate custom MCP server: %w", err)
		}
		result.Generated = generated
	}

	return result, nil
}

// loadCachedRegistry loads the cached MCP registry
func loadCachedRegistry() (*MCPRegistry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cacheFile := filepath.Join(homeDir, ".steele", "mcp-registry.json")
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var registry MCPRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, err
	}

	return &registry, nil
}

// refreshRegistry fetches the latest MCP server registry
func refreshRegistry() (*MCPRegistry, error) {
	registry := &MCPRegistry{
		LastUpdated: time.Now(),
		Servers:     make(map[string]MCPServer),
	}

	// Add known official servers
	addKnownServers(registry)

	// Try to fetch from GitHub API for community servers
	if err := fetchCommunityServers(registry); err != nil {
		fmt.Printf("⚠️  Failed to fetch community servers: %v\n", err)
	}

	// Cache the registry
	if err := cacheRegistry(registry); err != nil {
		fmt.Printf("⚠️  Failed to cache registry: %v\n", err)
	}

	return registry, nil
}

// addKnownServers adds known official and high-quality community servers
func addKnownServers(registry *MCPRegistry) {
	knownServers := []MCPServer{
		{
			Name:        "filesystem",
			Type:        "official",
			Framework:   "universal",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			NPMPackage:  "@modelcontextprotocol/server-filesystem",
			Command:     []string{"npx"},
			Args:        []string{"@modelcontextprotocol/server-filesystem"},
			Description: "File system operations and file reading",
		},
		{
			Name:        "database",
			Type:        "official", 
			Framework:   "universal",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			NPMPackage:  "@modelcontextprotocol/server-database",
			Command:     []string{"npx"},
			Args:        []string{"@modelcontextprotocol/server-database"},
			Description: "Database operations and schema inspection",
		},
		{
			Name:        "git",
			Type:        "official",
			Framework:   "universal",
			Repository:  "https://github.com/modelcontextprotocol/servers",
			NPMPackage:  "@modelcontextprotocol/server-git",
			Command:     []string{"npx"},
			Args:        []string{"@modelcontextprotocol/server-git"},
			Description: "Git repository operations and history",
		},
	}

	for _, server := range knownServers {
		registry.Servers[server.Name] = server
	}
}

// fetchCommunityServers searches GitHub for community MCP servers
func fetchCommunityServers(registry *MCPRegistry) error {
	// This would implement GitHub API search for MCP servers
	// For now, we'll skip this to avoid API rate limits
	return nil
}

// cacheRegistry saves the registry to cache
func cacheRegistry(registry *MCPRegistry) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	steeleDir := filepath.Join(homeDir, ".steele")
	if err := os.MkdirAll(steeleDir, 0755); err != nil {
		return err
	}

	cacheFile := filepath.Join(steeleDir, "mcp-registry.json")
	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0644)
}

// matchesFramework checks if a server matches the given framework
func matchesFramework(server MCPServer, framework string, keywords []string) bool {
	// Exact framework match
	if server.Framework == framework || server.Framework == "universal" {
		return true
	}

	// Check if any keywords match the server name or description
	searchText := strings.ToLower(server.Name + " " + server.Description)
	for _, keyword := range keywords {
		if strings.Contains(searchText, strings.ToLower(keyword)) {
			return true
		}
	}

	return false
}

// generateCustomServer creates a custom MCP server for the framework
func generateCustomServer(framework string) (*MCPServer, error) {
	switch framework {
	case "laravel":
		return generateLaravelServer(), nil
	case "django":
		return generateDjangoServer(), nil
	default:
		return generateUniversalServer(framework), nil
	}
}

// generateLaravelServer creates a Laravel-specific MCP server
func generateLaravelServer() *MCPServer {
	return &MCPServer{
		Name:        "steele-laravel",
		Type:        "generated",
		Framework:   "laravel",
		Command:     []string{"node"},
		Args:        []string{"ai/mcp-server/index.js"},
		Description: "Generated Laravel MCP server with Artisan, Eloquent, and Blade support",
		Environment: map[string]string{
			"LARAVEL_PROJECT": "true",
			"FRAMEWORK":       "laravel",
		},
	}
}

// generateDjangoServer creates a Django-specific MCP server
func generateDjangoServer() *MCPServer {
	return &MCPServer{
		Name:        "steele-django",
		Type:        "generated",
		Framework:   "django",
		Command:     []string{"node"},
		Args:        []string{"ai/mcp-server/index.js"},
		Description: "Generated Django MCP server with manage.py, ORM, and templates support",
		Environment: map[string]string{
			"DJANGO_PROJECT": "true",
			"FRAMEWORK":      "django",
		},
	}
}

// generateUniversalServer creates a universal MCP server
func generateUniversalServer(framework string) *MCPServer {
	return &MCPServer{
		Name:        fmt.Sprintf("steele-%s", framework),
		Type:        "generated",
		Framework:   framework,
		Command:     []string{"node"},
		Args:        []string{"ai/mcp-server/index.js"},
		Description: fmt.Sprintf("Generated %s MCP server with basic project support", framework),
		Environment: map[string]string{
			"FRAMEWORK": framework,
		},
	}
}

// InstallMCPServer installs the selected MCP server
func InstallMCPServer(server MCPServer, projectDir string) error {
	mcpDir := filepath.Join(projectDir, "ai", "mcp-server")
	
	switch server.Type {
	case "official", "community":
		return installOfficialServer(server, mcpDir)
	case "generated":
		return installGeneratedServer(server, mcpDir)
	default:
		return fmt.Errorf("unknown server type: %s", server.Type)
	}
}

// installOfficialServer installs an official or community MCP server
func installOfficialServer(server MCPServer, mcpDir string) error {
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		return err
	}

	// Create package.json for the server
	packageJSON := map[string]interface{}{
		"name":         server.Name,
		"version":      "1.0.0",
		"description":  server.Description,
		"main":         "index.js",
		"dependencies": map[string]string{},
	}

	if server.NPMPackage != "" {
		packageJSON["dependencies"].(map[string]string)[server.NPMPackage] = "latest"
	}

	packageData, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return err
	}

	packageFile := filepath.Join(mcpDir, "package.json")
	if err := os.WriteFile(packageFile, packageData, 0644); err != nil {
		return err
	}

	// Create a simple index.js that imports the official server
	indexJS := fmt.Sprintf(`#!/usr/bin/env node
// Auto-generated MCP server wrapper for %s
const { spawn } = require('child_process');

// Forward to the official server
const args = process.argv.slice(2);
const cmd = spawn('%s', %s.concat(args), { 
  stdio: 'inherit',
  env: { ...process.env, %s }
});

cmd.on('close', (code) => {
  process.exit(code);
});
`, server.Name, server.Command[0], formatArgs(server.Args), formatEnv(server.Environment))

	indexFile := filepath.Join(mcpDir, "index.js")
	return os.WriteFile(indexFile, []byte(indexJS), 0755)
}

// installGeneratedServer installs a generated MCP server
func installGeneratedServer(server MCPServer, mcpDir string) error {
	// Generated servers are handled by the GenerateServerFromTemplate function
	// This function should not be called for generated servers
	return fmt.Errorf("generated servers should use GenerateServerFromTemplate")
}

// Helper functions
func formatArgs(args []string) string {
	if len(args) == 0 {
		return "[]"
	}
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = fmt.Sprintf("'%s'", arg)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func formatEnv(env map[string]string) string {
	if len(env) == 0 {
		return ""
	}
	pairs := make([]string, 0, len(env))
	for k, v := range env {
		pairs = append(pairs, fmt.Sprintf("%s: '%s'", k, v))
	}
	return strings.Join(pairs, ", ")
}

func getServerTemplatePath(framework string) string {
	// This would return the path to server templates
	// For now, return empty to indicate fallback needed
	return ""
}
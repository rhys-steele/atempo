package providers

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"atempo/internal/logger"
	"gopkg.in/yaml.v3"
)

// CIProvider represents the CI/CD provider (duplicated to avoid import cycle)
type CIProvider string

const (
	ProviderGitHub CIProvider = "github"
	ProviderGitLab CIProvider = "gitlab"
)

// CIConfig represents CI configuration (duplicated to avoid import cycle)
type CIConfig struct {
	Provider      CIProvider              `json:"provider"`
	Framework     string                  `json:"framework"`
	ProjectName   string                  `json:"project_name"`
	ProjectPath   string                  `json:"project_path"`
	RepoURL       string                  `json:"repo_url"`
	Settings      map[string]interface{}  `json:"settings"`
}

// GitHubProvider implements the Provider interface for GitHub Actions
type GitHubProvider struct {
	logger *logger.Logger
}

// NewGitHubProvider creates a new GitHub Actions provider
func NewGitHubProvider(logger *logger.Logger) *GitHubProvider {
	return &GitHubProvider{
		logger: logger,
	}
}

// Name returns the provider name
func (g *GitHubProvider) Name() string {
	return "github"
}

// GetConfigFileName returns the GitHub Actions workflow filename
func (g *GitHubProvider) GetConfigFileName() string {
	return "ci.yml"
}

// GetConfigPath returns the full path where the GitHub Actions config should be written
func (g *GitHubProvider) GetConfigPath(projectPath string) string {
	return filepath.Join(projectPath, ".github", "workflows", "ci.yml")
}

// SupportedFrameworks returns the frameworks supported by GitHub Actions
func (g *GitHubProvider) SupportedFrameworks() []string {
	return []string{"laravel", "django", "express", "lambda-node"}
}

// ValidateConfig validates a CI configuration for GitHub Actions
func (g *GitHubProvider) ValidateConfig(config *CIConfig) error {
	if config == nil {
		return fmt.Errorf("configuration is nil")
	}

	if config.Provider != ProviderGitHub {
		return fmt.Errorf("invalid provider '%s' for GitHub Actions", config.Provider)
	}

	// Validate framework support
	supported := false
	for _, framework := range g.SupportedFrameworks() {
		if config.Framework == framework {
			supported = true
			break
		}
	}
	if !supported {
		return fmt.Errorf("framework '%s' is not supported by GitHub Actions. Supported: %v", 
			config.Framework, g.SupportedFrameworks())
	}

	// Validate settings if present
	if config.Settings != nil {
		return g.validateSettings(config.Settings, config.Framework)
	}

	return nil
}

// GenerateConfig generates a GitHub Actions workflow configuration
func (g *GitHubProvider) GenerateConfig(config *CIConfig, templateFS embed.FS) ([]byte, error) {
	// Template path based on framework (new structure)
	templatePath := fmt.Sprintf("templates/frameworks/%s/ci/github.yml", config.Framework)
	
	// Try to read the framework-specific template
	templateContent, err := templateFS.ReadFile(templatePath)
	if err != nil {
		// Fallback to filesystem for development (when embedding is disabled)
		possiblePaths := []string{
			templatePath,
			fmt.Sprintf("../%s", templatePath),
			fmt.Sprintf("../../%s", templatePath),
			fmt.Sprintf("../../../%s", templatePath),
		}

		var fsErr error
		found := false
		for _, fallbackPath := range possiblePaths {
			if templateContent, fsErr = os.ReadFile(fallbackPath); fsErr == nil {
				found = true
				break
			}
		}

		if !found {
			// Try shared template
			sharedPath := "templates/shared/ci/github.yml"
			templateContent, err = templateFS.ReadFile(sharedPath)
			if err != nil {
				// Fallback to filesystem for shared template
				sharedPossiblePaths := []string{
					sharedPath,
					fmt.Sprintf("../%s", sharedPath),
					fmt.Sprintf("../../%s", sharedPath),
					fmt.Sprintf("../../../%s", sharedPath),
				}
				
				for _, fallbackPath := range sharedPossiblePaths {
					if templateContent, fsErr = os.ReadFile(fallbackPath); fsErr == nil {
						found = true
						break
					}
				}
				
				if !found {
					return nil, fmt.Errorf("template not found for framework '%s' and basic template unavailable (tried embedded and filesystem): %w", config.Framework, err)
				}
			}
		}
	}

	// Parse settings into structured format
	settings, err := g.parseSettings(config.Settings, config.Framework)
	if err != nil {
		return nil, fmt.Errorf("invalid settings: %w", err)
	}

	// Execute template with settings
	return g.executeTemplate(templateContent, settings)
}

// GetDefaultSettings returns default settings for a framework
func (g *GitHubProvider) GetDefaultSettings(framework string) map[string]interface{} {
	defaults := map[string]interface{}{
		"workflow_name": "CI",
		"triggers":      []string{"push", "pull_request"},
		"branches":      []string{"main", "develop"},
		"timeout":       30, // minutes
	}

	// Framework-specific defaults
	switch framework {
	case "laravel":
		defaults["php_version"] = "8.2"
		defaults["services"] = []string{"mysql", "redis"}
		defaults["cache_paths"] = []string{"vendor", "node_modules"}
		defaults["environment"] = map[string]string{
			"APP_ENV": "testing",
			"DB_CONNECTION": "mysql",
			"DB_DATABASE": "testing",
		}
	case "django":
		defaults["python_version"] = "3.11"
		defaults["services"] = []string{"postgres", "redis"}
		defaults["cache_paths"] = []string{"~/.cache/pip"}
		defaults["environment"] = map[string]string{
			"DJANGO_SETTINGS_MODULE": "settings.test",
			"DATABASE_URL": "postgres://test:test@localhost:5432/test",
		}
	case "express":
		defaults["node_version"] = "18"
		defaults["services"] = []string{"redis"}
		defaults["cache_paths"] = []string{"node_modules"}
		defaults["environment"] = map[string]string{
			"NODE_ENV": "test",
		}
	case "lambda-node":
		defaults["node_version"] = "18"
		defaults["services"] = []string{}
		defaults["cache_paths"] = []string{"node_modules"}
		defaults["environment"] = map[string]string{
			"NODE_ENV": "test",
			"AWS_DEFAULT_REGION": "us-east-1",
		}
	}

	return defaults
}

// PromptForSettings prompts the user for provider-specific settings
func (g *GitHubProvider) PromptForSettings(framework string) (map[string]interface{}, error) {
	// For Phase 2, return default settings
	// In Phase 3, this will be implemented with interactive prompts
	defaults := g.GetDefaultSettings(framework)
	
	// Note: Logger integration will be added in Phase 3
	
	return defaults, nil
}

// validateSettings validates GitHub Actions specific settings
func (g *GitHubProvider) validateSettings(settings map[string]interface{}, framework string) error {
	// Validate workflow name
	if workflowName, exists := settings["workflow_name"]; exists {
		if name, ok := workflowName.(string); ok {
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf("workflow_name cannot be empty")
			}
		} else {
			return fmt.Errorf("workflow_name must be a string")
		}
	}

	// Validate triggers
	if triggers, exists := settings["triggers"]; exists {
		if triggerList, ok := triggers.([]interface{}); ok {
			validTriggers := []string{"push", "pull_request", "schedule", "workflow_dispatch"}
			for _, trigger := range triggerList {
				if triggerStr, ok := trigger.(string); ok {
					isValid := false
					for _, validTrigger := range validTriggers {
						if triggerStr == validTrigger {
							isValid = true
							break
						}
					}
					if !isValid {
						return fmt.Errorf("invalid trigger '%s', must be one of: %v", triggerStr, validTriggers)
					}
				} else {
					return fmt.Errorf("trigger must be a string")
				}
			}
		} else {
			return fmt.Errorf("triggers must be an array")
		}
	}

	// Validate timeout
	if timeout, exists := settings["timeout"]; exists {
		if timeoutNum, ok := timeout.(float64); ok {
			if timeoutNum <= 0 || timeoutNum > 360 {
				return fmt.Errorf("timeout must be between 1 and 360 minutes")
			}
		} else {
			return fmt.Errorf("timeout must be a number")
		}
	}

	return nil
}

// parseSettings converts raw settings into a structured format for template execution
func (g *GitHubProvider) parseSettings(settings map[string]interface{}, framework string) (map[string]interface{}, error) {
	if settings == nil {
		settings = g.GetDefaultSettings(framework)
	}

	// Ensure all required fields have defaults
	defaults := g.GetDefaultSettings(framework)
	for key, value := range defaults {
		if _, exists := settings[key]; !exists {
			settings[key] = value
		}
	}

	// Convert to template-friendly format
	templateData := make(map[string]interface{})

	// Basic workflow settings
	templateData["WorkflowName"] = settings["workflow_name"]
	templateData["Triggers"] = settings["triggers"]
	templateData["Branches"] = settings["branches"]
	templateData["Timeout"] = settings["timeout"]

	// Framework-specific settings
	switch framework {
	case "laravel":
		templateData["PHPVersion"] = settings["php_version"]
		templateData["Services"] = settings["services"]
		templateData["CachePaths"] = settings["cache_paths"]
		templateData["Environment"] = settings["environment"]
	case "django":
		templateData["PythonVersion"] = settings["python_version"]
		templateData["Services"] = settings["services"]
		templateData["CachePaths"] = settings["cache_paths"]
		templateData["Environment"] = settings["environment"]
	case "express", "lambda-node":
		templateData["NodeVersion"] = settings["node_version"]
		templateData["Services"] = settings["services"]
		templateData["CachePaths"] = settings["cache_paths"]
		templateData["Environment"] = settings["environment"]
	}

	return templateData, nil
}

// executeTemplate executes a template with the given settings and validates the output
func (g *GitHubProvider) executeTemplate(templateContent []byte, settings interface{}) ([]byte, error) {
	tmpl, err := template.New("github-workflow").Parse(string(templateContent))
	if err != nil {
		return nil, fmt.Errorf("template parsing error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, settings); err != nil {
		return nil, fmt.Errorf("template execution error: %w", err)
	}

	// Validate generated YAML
	var yamlCheck interface{}
	if err := yaml.Unmarshal(buf.Bytes(), &yamlCheck); err != nil {
		return nil, fmt.Errorf("generated invalid YAML: %w", err)
	}

	return buf.Bytes(), nil
}
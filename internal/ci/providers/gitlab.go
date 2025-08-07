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

// GitLabProvider implements the Provider interface for GitLab CI
type GitLabProvider struct {
	logger *logger.Logger
}

// NewGitLabProvider creates a new GitLab CI provider
func NewGitLabProvider(logger *logger.Logger) *GitLabProvider {
	return &GitLabProvider{
		logger: logger,
	}
}

// Name returns the provider name
func (g *GitLabProvider) Name() string {
	return "gitlab"
}

// GetConfigFileName returns the GitLab CI configuration filename
func (g *GitLabProvider) GetConfigFileName() string {
	return ".gitlab-ci.yml"
}

// GetConfigPath returns the full path where the GitLab CI config should be written
func (g *GitLabProvider) GetConfigPath(projectPath string) string {
	return filepath.Join(projectPath, ".gitlab-ci.yml")
}

// SupportedFrameworks returns the frameworks supported by GitLab CI
func (g *GitLabProvider) SupportedFrameworks() []string {
	return []string{"laravel", "django", "express", "lambda-node"}
}

// ValidateConfig validates a CI configuration for GitLab CI
func (g *GitLabProvider) ValidateConfig(config *CIConfig) error {
	if config == nil {
		return fmt.Errorf("configuration is nil")
	}

	if config.Provider != ProviderGitLab {
		return fmt.Errorf("invalid provider '%s' for GitLab CI", config.Provider)
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
		return fmt.Errorf("framework '%s' is not supported by GitLab CI. Supported: %v", 
			config.Framework, g.SupportedFrameworks())
	}

	// Validate settings if present
	if config.Settings != nil {
		return g.validateSettings(config.Settings, config.Framework)
	}

	return nil
}

// GenerateConfig generates a GitLab CI configuration
func (g *GitLabProvider) GenerateConfig(config *CIConfig, templateFS embed.FS) ([]byte, error) {
	// Template path based on framework (new structure)  
	templatePath := fmt.Sprintf("templates/frameworks/%s/ci/gitlab.yml", config.Framework)
	
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
			sharedPath := "templates/shared/ci/gitlab.yml"
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
func (g *GitLabProvider) GetDefaultSettings(framework string) map[string]interface{} {
	defaults := map[string]interface{}{
		"stages":        []string{"build", "test"},
		"before_script": []string{},
		"variables":     map[string]string{},
		"cache_paths":   []string{},
		"services":      []string{},
		"artifact_paths": []string{},
	}

	// Framework-specific defaults
	switch framework {
	case "laravel":
		defaults["image"] = "php:8.2"
		defaults["services"] = []string{"mysql:8.0", "redis:7"}
		defaults["cache_paths"] = []string{"vendor/", "node_modules/"}
		defaults["before_script"] = []string{
			"apt-get update -qq && apt-get install -y -qq git curl libmcrypt-dev libjpeg-dev libpng-dev libfreetype6-dev libbz2-dev",
			"docker-php-ext-install pdo_mysql zip",
			"curl -sS https://getcomposer.org/installer | php",
			"php composer.phar install --no-dev --no-scripts",
		}
		defaults["variables"] = map[string]string{
			"MYSQL_DATABASE": "testing",
			"MYSQL_ROOT_PASSWORD": "secret",
			"APP_ENV": "testing",
		}
	case "django":
		defaults["image"] = "python:3.11"
		defaults["services"] = []string{"postgres:14", "redis:7"}
		defaults["cache_paths"] = []string{".cache/pip/"}
		defaults["before_script"] = []string{
			"python -V",
			"pip install -r requirements.txt",
		}
		defaults["variables"] = map[string]string{
			"POSTGRES_DB": "test",
			"POSTGRES_USER": "test", 
			"POSTGRES_PASSWORD": "test",
			"DJANGO_SETTINGS_MODULE": "settings.test",
		}
	case "express":
		defaults["image"] = "node:18"
		defaults["services"] = []string{"redis:7"}
		defaults["cache_paths"] = []string{"node_modules/"}
		defaults["before_script"] = []string{
			"node --version",
			"npm install",
		}
		defaults["variables"] = map[string]string{
			"NODE_ENV": "test",
		}
	case "lambda-node":
		defaults["image"] = "node:18"
		defaults["services"] = []string{}
		defaults["cache_paths"] = []string{"node_modules/"}
		defaults["before_script"] = []string{
			"node --version",
			"npm install",
		}
		defaults["variables"] = map[string]string{
			"NODE_ENV": "test",
			"AWS_DEFAULT_REGION": "us-east-1",
		}
	}

	return defaults
}

// PromptForSettings prompts the user for provider-specific settings
func (g *GitLabProvider) PromptForSettings(framework string) (map[string]interface{}, error) {
	// For Phase 2, return default settings
	// In Phase 3, this will be implemented with interactive prompts
	defaults := g.GetDefaultSettings(framework)
	
	// Note: Logger integration will be added in Phase 3
	
	return defaults, nil
}

// validateSettings validates GitLab CI specific settings
func (g *GitLabProvider) validateSettings(settings map[string]interface{}, framework string) error {
	// Validate image
	if image, exists := settings["image"]; exists {
		if imageStr, ok := image.(string); ok {
			if strings.TrimSpace(imageStr) == "" {
				return fmt.Errorf("image cannot be empty")
			}
		} else {
			return fmt.Errorf("image must be a string")
		}
	}

	// Validate stages
	if stages, exists := settings["stages"]; exists {
		if stageList, ok := stages.([]interface{}); ok {
			if len(stageList) == 0 {
				return fmt.Errorf("stages cannot be empty")
			}
			for _, stage := range stageList {
				if stageStr, ok := stage.(string); ok {
					if strings.TrimSpace(stageStr) == "" {
						return fmt.Errorf("stage name cannot be empty")
					}
				} else {
					return fmt.Errorf("stage must be a string")
				}
			}
		} else {
			return fmt.Errorf("stages must be an array")
		}
	}

	// Validate variables
	if variables, exists := settings["variables"]; exists {
		if variableMap, ok := variables.(map[string]interface{}); ok {
			for key, value := range variableMap {
				if strings.TrimSpace(key) == "" {
					return fmt.Errorf("variable key cannot be empty")
				}
				if _, ok := value.(string); !ok {
					return fmt.Errorf("variable value must be a string")
				}
			}
		} else {
			return fmt.Errorf("variables must be a map")
		}
	}

	return nil
}

// parseSettings converts raw settings into a structured format for template execution
func (g *GitLabProvider) parseSettings(settings map[string]interface{}, framework string) (map[string]interface{}, error) {
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

	// Basic GitLab CI settings
	templateData["Image"] = settings["image"]
	templateData["Stages"] = settings["stages"]
	templateData["BeforeScript"] = settings["before_script"]
	templateData["Variables"] = settings["variables"]
	templateData["CachePaths"] = settings["cache_paths"]
	templateData["Services"] = settings["services"]
	templateData["ArtifactPaths"] = settings["artifact_paths"]

	// Framework-specific settings
	switch framework {
	case "laravel":
		templateData["TestCommand"] = "vendor/bin/phpunit --coverage-text --colors=never"
		templateData["BuildCommand"] = "composer install --no-progress --prefer-dist --optimize-autoloader"
	case "django":
		templateData["TestCommand"] = "python manage.py test"
		templateData["BuildCommand"] = "pip install -r requirements.txt"
	case "express":
		templateData["TestCommand"] = "npm test"
		templateData["BuildCommand"] = "npm install"
	case "lambda-node":
		templateData["TestCommand"] = "npm test"
		templateData["BuildCommand"] = "npm install"
	}

	return templateData, nil
}

// executeTemplate executes a template with the given settings and validates the output
func (g *GitLabProvider) executeTemplate(templateContent []byte, settings interface{}) ([]byte, error) {
	tmpl, err := template.New("gitlab-ci").Parse(string(templateContent))
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
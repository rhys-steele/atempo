package commands

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateLoader handles loading and parsing of AI manifest templates
type TemplateLoader struct {
	templatesFS fs.FS
}

// NewTemplateLoader creates a new template loader
func NewTemplateLoader(templatesFS fs.FS) *TemplateLoader {
	return &TemplateLoader{
		templatesFS: templatesFS,
	}
}

// ManifestTemplate represents a loaded manifest template
type ManifestTemplate struct {
	TemplateVersion  string `json:"template_version"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	MarkdownTemplate string `json:"markdown_template"`
}

// FrameworkConfig represents framework-specific AI configuration
type FrameworkConfig struct {
	Framework          string             `json:"framework"`
	Language           string             `json:"language"`
	LatestVersion      string             `json:"latest_version"`
	AIFeatures         AIFeatures         `json:"ai_features"`
	DevelopmentContext DevelopmentContext `json:"development_context"`
	MCPConfig          MCPConfig          `json:"mcp_config"`
}

// AIFeatures contains AI-specific framework configuration
type AIFeatures struct {
	DefaultProjectTypes       []string          `json:"default_project_types"`
	CoreFeatures              []string          `json:"core_features"`
	ArchitecturePatterns      map[string]string `json:"architecture_patterns"`
	FrameworkPatternsTemplate string            `json:"framework_patterns_template"`
	TechnicalStack            []string          `json:"technical_stack"`
	ProjectAnalysisKeywords   map[string]string `json:"project_analysis_keywords"`
}

// DevelopmentContext contains development environment configuration
type DevelopmentContext struct {
	PackageManager  string              `json:"package_manager"`
	Structure       map[string]string   `json:"structure"`
	Commands        map[string]string   `json:"commands"`
	Docker          DockerConfig        `json:"docker"`
	Patterns        map[string][]string `json:"patterns"`
	BestPractices   []string            `json:"best_practices"`
	Environment     EnvironmentConfig   `json:"environment"`
	Troubleshooting map[string]string   `json:"troubleshooting"`
	CodeTemplates   map[string]string   `json:"code_templates"`
	AppsStructure   AppsStructureConfig `json:"apps_structure,omitempty"`
}

type DockerConfig struct {
	AppContainer      string `json:"app_container"`
	DatabaseContainer string `json:"database_container"`
	RedisContainer    string `json:"redis_container"`
	WorkingDirectory  string `json:"working_directory"`
}

type EnvironmentConfig struct {
	RequiredEnvVars  []string `json:"required_env_vars"`
	DevelopmentTools []string `json:"development_tools"`
}

type AppsStructureConfig struct {
	RecommendedApps []string `json:"recommended_apps"`
	AppStructure    string   `json:"app_structure"`
}

// MCPConfig contains MCP server configuration
type MCPConfig struct {
	Servers map[string]MCPServer `json:"servers"`
}

type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Cwd     string            `json:"cwd"`
	Env     map[string]string `json:"env"`
}

// InteractivePrompts represents the interactive prompting configuration
type InteractivePrompts struct {
	Prompts    PromptsConfig `json:"prompts"`
	UIElements UIElements    `json:"ui_elements"`
}

type PromptsConfig struct {
	ProjectDescription ProjectDescriptionPrompt `json:"project_description"`
	AdditionalFeatures AdditionalFeaturesPrompt `json:"additional_features"`
	Complexity         ComplexityPrompt         `json:"complexity"`
}

type ProjectDescriptionPrompt struct {
	Question string   `json:"question"`
	Subtitle string   `json:"subtitle"`
	Examples []string `json:"examples"`
}

type AdditionalFeaturesPrompt struct {
	Question string   `json:"question"`
	Subtitle string   `json:"subtitle"`
	Options  []string `json:"options"`
}

type ComplexityPrompt struct {
	Question string             `json:"question"`
	Options  []ComplexityOption `json:"options"`
}

type ComplexityOption struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type UIElements struct {
	Header       string              `json:"header"`
	Separator    string              `json:"separator"`
	Subtitle     string              `json:"subtitle"`
	AuthRequired AuthRequiredElement `json:"auth_required"`
}

type AuthRequiredElement struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Action   string `json:"action"`
	Fallback string `json:"fallback"`
}

// LoadManifestTemplate loads a manifest template by name
func (tl *TemplateLoader) LoadManifestTemplate(templateName string) (*ManifestTemplate, error) {
	templatePath := filepath.Join("ai", "manifests", templateName+".json")

	data, err := fs.ReadFile(tl.templatesFS, templatePath)
	if err != nil {
		// Fallback to filesystem for development (when embedding is disabled)
		// Try multiple relative paths since the working directory varies
		possiblePaths := []string{
			filepath.Join("templates", templatePath),
			filepath.Join("../templates", templatePath),
			filepath.Join("../../templates", templatePath),
		}

		for _, fallbackPath := range possiblePaths {
			if data, err = os.ReadFile(fallbackPath); err == nil {
				break
			}
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read template %s (tried embedded and filesystem): %w", templateName, err)
		}
	}

	var manifestTemplate ManifestTemplate
	if err := json.Unmarshal(data, &manifestTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	return &manifestTemplate, nil
}

// LoadFrameworkConfig loads framework-specific AI configuration
func (tl *TemplateLoader) LoadFrameworkConfig(framework string) (*FrameworkConfig, error) {
	configPath := filepath.Join("frameworks", framework, "ai", "ai-config.json")

	data, err := fs.ReadFile(tl.templatesFS, configPath)
	if err != nil {
		// Fallback to filesystem for development
		possiblePaths := []string{
			filepath.Join("templates", configPath),
			filepath.Join("../templates", configPath),
			filepath.Join("../../templates", configPath),
		}

		for _, fallbackPath := range possiblePaths {
			if data, err = os.ReadFile(fallbackPath); err == nil {
				break
			}
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read framework config %s (tried embedded and filesystem): %w", framework, err)
		}
	}

	var config FrameworkConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse framework config %s: %w", framework, err)
	}

	return &config, nil
}

// LoadInteractivePrompts loads the interactive prompting configuration
func (tl *TemplateLoader) LoadInteractivePrompts() (*InteractivePrompts, error) {
	promptsPath := filepath.Join("ai", "prompts", "interactive-prompts.json")

	data, err := fs.ReadFile(tl.templatesFS, promptsPath)
	if err != nil {
		// Fallback to filesystem for development
		possiblePaths := []string{
			filepath.Join("templates", promptsPath),
			filepath.Join("../templates", promptsPath),
			filepath.Join("../../templates", promptsPath),
		}

		for _, fallbackPath := range possiblePaths {
			if data, err = os.ReadFile(fallbackPath); err == nil {
				break
			}
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read interactive prompts (tried embedded and filesystem): %w", err)
		}
	}

	var prompts InteractivePrompts
	if err := json.Unmarshal(data, &prompts); err != nil {
		return nil, fmt.Errorf("failed to parse interactive prompts: %w", err)
	}

	return &prompts, nil
}

// GenerateFromTemplate generates content from a template with the given data
func (tl *TemplateLoader) GenerateFromTemplate(templateContent string, data interface{}) (string, error) {
	// Create template with helper functions
	tmpl := template.New("manifest").Funcs(template.FuncMap{
		"title": strings.Title,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	})

	tmpl, err := tmpl.Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"atempo/internal/utils"
)

// ContextConfig defines configuration for AI context generation
type ContextConfig struct {
	Framework           string                 `json:"framework"`
	Language            string                 `json:"language"`
	LatestVersion       string                 `json:"latest_version"`
	AIFeatures          AIFeatures             `json:"ai_features"`
	DevelopmentContext  DevelopmentContext     `json:"development_context"`
	MCPConfig           MCPConfig              `json:"mcp_config"`
}

// AIFeatures defines AI-specific features and patterns
type AIFeatures struct {
	DefaultProjectTypes      []string                   `json:"default_project_types"`
	CoreFeatures             []string                   `json:"core_features"`
	ArchitecturePatterns     map[string]string          `json:"architecture_patterns"`
	FrameworkPatternsTemplate string                    `json:"framework_patterns_template"`
	TechnicalStack           []string                   `json:"technical_stack"`
	ProjectAnalysisKeywords  map[string]string          `json:"project_analysis_keywords"`
}

// DevelopmentContext defines development-specific context
type DevelopmentContext struct {
	PackageManager  string                 `json:"package_manager"`
	Structure       map[string]string      `json:"structure"`
	Commands        map[string]string      `json:"commands"`
	Docker          DockerConfig           `json:"docker"`
	Patterns        map[string][]string    `json:"patterns"`
	BestPractices   []string              `json:"best_practices"`
	Environment     EnvironmentConfig      `json:"environment"`
	Troubleshooting map[string]string      `json:"troubleshooting"`
	CodeTemplates   map[string]string      `json:"code_templates"`
}

// DockerConfig defines Docker-specific configuration
type DockerConfig struct {
	AppContainer      string `json:"app_container"`
	DatabaseContainer string `json:"database_container"`
	RedisContainer    string `json:"redis_container"`
	WorkingDirectory  string `json:"working_directory"`
}

// EnvironmentConfig defines environment-specific configuration
type EnvironmentConfig struct {
	RequiredEnvVars    []string `json:"required_env_vars"`
	DevelopmentTools   []string `json:"development_tools"`
}

// MCPConfig defines MCP server configuration
type MCPConfig struct {
	Servers map[string]MCPServer `json:"servers"`
}

// MCPServer defines individual MCP server configuration
type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	CWD     string            `json:"cwd"`
	Env     map[string]string `json:"env"`
}

// GenerateAIContext generates AI context files for a project
func GenerateAIContext(projectDir, projectName, framework, language, version string, aiEnabled bool) error {
	if !aiEnabled {
		// Use static templates - this is already handled by scaffold.go
		return nil
	}

	// For AI-enabled mode, generate dynamic context
	return generateDynamicContext(projectDir, projectName, framework, language, version)
}

// generateDynamicContext generates AI context dynamically based on project analysis
func generateDynamicContext(projectDir, projectName, framework, language, version string) error {
	// Load framework-specific configuration
	config, err := loadFrameworkConfig(framework)
	if err != nil {
		return fmt.Errorf("failed to load framework config: %w", err)
	}

	// Analyze project structure
	projectAnalysis, err := analyzeProjectStructure(projectDir, framework)
	if err != nil {
		return fmt.Errorf("failed to analyze project: %w", err)
	}

	// Generate context files
	aiDir := filepath.Join(projectDir, ".ai")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		return fmt.Errorf("failed to create .ai directory: %w", err)
	}

	// Generate master context file
	if err := generateContextMD(aiDir, projectName, framework, language, version, config, projectAnalysis); err != nil {
		return fmt.Errorf("failed to generate context.md: %w", err)
	}

	// Generate project overview
	if err := generateProjectOverview(aiDir, projectName, framework, language, version, config, projectAnalysis); err != nil {
		return fmt.Errorf("failed to generate project-overview.md: %w", err)
	}

	// Generate codebase map
	if err := generateCodebaseMap(aiDir, projectName, framework, language, config, projectAnalysis); err != nil {
		return fmt.Errorf("failed to generate codebase-map.md: %w", err)
	}

	// Generate development workflows
	if err := generateDevelopmentWorkflows(aiDir, projectName, framework, language, config); err != nil {
		return fmt.Errorf("failed to generate development-workflows.md: %w", err)
	}

	// Generate patterns and conventions
	if err := generatePatternsAndConventions(aiDir, projectName, framework, language, config); err != nil {
		return fmt.Errorf("failed to generate patterns-and-conventions.md: %w", err)
	}

	// Generate UI/UX guidelines
	if err := generateUIUXGuidelines(aiDir, projectName, framework, language, config); err != nil {
		return fmt.Errorf("failed to generate ui-ux-guidelines.md: %w", err)
	}

	return nil
}

// loadFrameworkConfig loads framework-specific configuration
func loadFrameworkConfig(framework string) (*ContextConfig, error) {
	// Try to load from templates directory
	configPath := filepath.Join("templates", "frameworks", framework, "ai", "ai-config.json")
	
	// Check if file exists
	if !utils.FileExists(configPath) {
		return nil, fmt.Errorf("ai-config.json not found for framework %s", framework)
	}

	// Read and parse configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ai-config.json: %w", err)
	}

	var config ContextConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse ai-config.json: %w", err)
	}

	return &config, nil
}

// ProjectAnalysis contains analysis results of the project
type ProjectAnalysis struct {
	ProjectType     string
	Features        []string
	Architecture    string
	Dependencies    []string
	Structure       map[string][]string
	DatabaseTables  []string
	APIEndpoints    []string
	HasTests        bool
	HasDocker       bool
	HasCI           bool
}

// analyzeProjectStructure analyzes the project structure and infers characteristics
func analyzeProjectStructure(projectDir, framework string) (*ProjectAnalysis, error) {
	analysis := &ProjectAnalysis{
		ProjectType:  "Web Application",
		Features:     []string{},
		Architecture: "MVC",
		Dependencies: []string{},
		Structure:    make(map[string][]string),
	}

	// Analyze based on framework
	switch framework {
	case "laravel":
		return analyzeLaravelProject(projectDir, analysis)
	case "django":
		return analyzeDjangoProject(projectDir, analysis)
	}

	return analysis, nil
}

// analyzeLaravelProject analyzes a Laravel project
func analyzeLaravelProject(projectDir string, analysis *ProjectAnalysis) (*ProjectAnalysis, error) {
	srcDir := filepath.Join(projectDir, "src")
	
	// Check if Laravel project exists
	if !utils.FileExists(filepath.Join(srcDir, "artisan")) {
		return analysis, nil // Project not yet created
	}

	// Analyze controllers
	controllerDir := filepath.Join(srcDir, "app", "Http", "Controllers")
	if utils.FileExists(controllerDir) {
		controllers, _ := filepath.Glob(filepath.Join(controllerDir, "*.php"))
		analysis.Structure["controllers"] = controllers
		
		// Infer features from controllers
		for _, controller := range controllers {
			name := strings.TrimSuffix(filepath.Base(controller), ".php")
			if strings.Contains(strings.ToLower(name), "api") {
				analysis.Features = append(analysis.Features, "REST API")
			}
			if strings.Contains(strings.ToLower(name), "admin") {
				analysis.Features = append(analysis.Features, "Admin Dashboard")
			}
		}
	}

	// Analyze models
	modelDir := filepath.Join(srcDir, "app", "Models")
	if utils.FileExists(modelDir) {
		models, _ := filepath.Glob(filepath.Join(modelDir, "*.php"))
		analysis.Structure["models"] = models
	}

	// Analyze migrations to infer database tables
	migrationDir := filepath.Join(srcDir, "database", "migrations")
	if utils.FileExists(migrationDir) {
		migrations, _ := filepath.Glob(filepath.Join(migrationDir, "*.php"))
		for _, migration := range migrations {
			name := filepath.Base(migration)
			if strings.Contains(name, "create_") {
				tableName := extractTableName(name)
				if tableName != "" {
					analysis.DatabaseTables = append(analysis.DatabaseTables, tableName)
				}
			}
		}
	}

	// Check for API routes
	apiRoutesFile := filepath.Join(srcDir, "routes", "api.php")
	if utils.FileExists(apiRoutesFile) {
		analysis.Features = append(analysis.Features, "REST API")
	}

	// Check for tests
	testDir := filepath.Join(srcDir, "tests")
	if utils.FileExists(testDir) {
		analysis.HasTests = true
	}

	// Check for Docker
	dockerFile := filepath.Join(projectDir, "docker-compose.yml")
	if utils.FileExists(dockerFile) {
		analysis.HasDocker = true
	}

	// Analyze composer.json for dependencies
	composerFile := filepath.Join(srcDir, "composer.json")
	if utils.FileExists(composerFile) {
		deps, _ := analyzeComposerDependencies(composerFile)
		analysis.Dependencies = deps
	}

	return analysis, nil
}

// analyzeDjangoProject analyzes a Django project
func analyzeDjangoProject(projectDir string, analysis *ProjectAnalysis) (*ProjectAnalysis, error) {
	srcDir := filepath.Join(projectDir, "src")
	
	// Check if Django project exists
	if !utils.FileExists(filepath.Join(srcDir, "manage.py")) {
		return analysis, nil // Project not yet created
	}

	// Analyze Django apps
	entries, err := os.ReadDir(srcDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				appPath := filepath.Join(srcDir, entry.Name())
				if utils.FileExists(filepath.Join(appPath, "models.py")) {
					analysis.Structure["apps"] = append(analysis.Structure["apps"], entry.Name())
				}
			}
		}
	}

	// Check for API framework
	requirementsFile := filepath.Join(srcDir, "requirements.txt")
	if utils.FileExists(requirementsFile) {
		content, err := os.ReadFile(requirementsFile)
		if err == nil {
			if strings.Contains(string(content), "djangorestframework") {
				analysis.Features = append(analysis.Features, "REST API")
			}
		}
	}

	// Check for tests
	if utils.FileExists(filepath.Join(srcDir, "tests")) {
		analysis.HasTests = true
	}

	// Check for Docker
	dockerFile := filepath.Join(projectDir, "docker-compose.yml")
	if utils.FileExists(dockerFile) {
		analysis.HasDocker = true
	}

	return analysis, nil
}

// extractTableName extracts table name from Laravel migration filename
func extractTableName(filename string) string {
	// Laravel migration naming: 2023_01_01_000000_create_users_table.php
	parts := strings.Split(filename, "_")
	if len(parts) >= 5 && parts[4] == "create" {
		// Extract table name (remove "table.php" suffix)
		tableName := strings.Join(parts[5:], "_")
		tableName = strings.TrimSuffix(tableName, "_table.php")
		return tableName
	}
	return ""
}

// analyzeComposerDependencies analyzes composer.json for dependencies
func analyzeComposerDependencies(composerFile string) ([]string, error) {
	content, err := os.ReadFile(composerFile)
	if err != nil {
		return nil, err
	}

	var composer map[string]interface{}
	if err := json.Unmarshal(content, &composer); err != nil {
		return nil, err
	}

	var deps []string
	if require, ok := composer["require"].(map[string]interface{}); ok {
		for pkg := range require {
			if !strings.HasPrefix(pkg, "php") {
				deps = append(deps, pkg)
			}
		}
	}

	return deps, nil
}

// generateContextMD generates the master context.md file
func generateContextMD(aiDir, projectName, framework, language, version string, config *ContextConfig, analysis *ProjectAnalysis) error {
	content := fmt.Sprintf("# %s Project AI Context System\n\n## Overview\nThe .ai directory contains structured context files designed to provide comprehensive project understanding for AI assistants working on the %s %s application. This system ensures consistent, informed development aligned with project standards and goals.\n\n## Directory Structure & File Index\n\n### .ai/context.md (This File)\n**Purpose**: Master context file serving as an index to all AI context files\n**Format**: Markdown with structured sections\n**Contents**:\n- Complete directory structure and file purpose explanations\n- Development standards and architectural guidelines\n- Current priorities and security requirements\n- Cross-references to all other context files\n\n### .ai/project-overview.md\n**Purpose**: High-level project mission, architecture, and feature overview\n**Format**: Markdown with architecture diagrams and feature lists\n**Contents**:\n- Project mission and value proposition\n- %s application architecture details\n- Current development state and feature status\n- Key components and technical highlights\n- Development workflows and deployment processes\n\n### .ai/codebase-map.md\n**Purpose**: Detailed technical documentation of the entire codebase structure\n**Format**: Markdown with code examples and file relationships\n**Contents**:\n- Complete directory structure with descriptions\n- File descriptions and responsibilities for each component\n- Architecture patterns and design principles\n- Key interfaces, models, and data flow diagrams\n- External dependencies and integration points\n\n### .ai/development-workflows.md\n**Purpose**: Comprehensive guide to development processes and commands\n**Format**: Markdown with code examples and workflow instructions\n**Contents**:\n- Quick start and setup instructions\n- Core development workflows (testing, database, deployment)\n- Command reference for %s and Docker operations\n- Testing, debugging, and troubleshooting procedures\n- Git workflow and deployment procedures\n\n### .ai/patterns-and-conventions.md\n**Purpose**: Detailed coding standards, patterns, and architectural conventions\n**Format**: Markdown with %s code examples and pattern explanations\n**Contents**:\n- Code architecture patterns and best practices\n- Naming conventions and %s coding standards\n- Error handling patterns and validation approaches\n- Testing patterns and documentation standards\n- Security, performance, and maintainability guidelines\n\n### .ai/ui-ux-guidelines.md\n**Purpose**: UI/UX standards and design guidelines for the application\n**Format**: Markdown with examples and design principles\n**Contents**:\n- Design system and component guidelines\n- User experience principles and accessibility standards\n- Frontend architecture and styling conventions\n- API design and response formatting standards\n- User interface patterns and interaction guidelines\n\n## AI Context Provision Strategy\n\nThis structure is optimized for AI assistants by providing:\n\n1. **Hierarchical Information**: Master context file references specialized guidelines\n2. **Specific Standards**: Clear, actionable rules rather than vague suggestions\n3. **Examples**: Concrete code examples and patterns\n4. **Cross-References**: Links between related context files\n5. **Structured Format**: Consistent markdown formatting for easy parsing\n\n## Development Standards\n\n### Code Quality\n- Follow %s best practices and conventions\n- Implement proper error handling and validation\n- Write comprehensive tests for all features\n- Use dependency injection and clean architecture patterns\n\n### Security Requirements\n- Never log or expose sensitive information\n- Validate all user inputs thoroughly\n- Use secure defaults for all configurations\n- Implement proper authentication and authorization\n- Follow security guidelines for %s applications\n\n### Testing Standards\n- Write unit tests for all business logic\n- Implement integration tests for API endpoints\n- Test database migrations and seeders\n- Verify security and validation rules\n- Test error cases and edge conditions\n\n### Architecture Guidelines\n- Follow %s conventions and best practices\n- Implement clean separation of concerns\n- Use service layer for business logic\n- Implement proper caching strategies\n- Design for scalability and maintainability\n\n## Project Analysis\n\n### Project Type\n**Detected Type**: %s\n\n### Core Features\n%s\n\n### Architecture Pattern\n**Pattern**: %s\n\n### Database Schema\n%s\n\n### Technical Stack\n%s\n\n## Current Development Priorities\n1. Core application functionality and features\n2. Comprehensive testing and quality assurance\n3. Performance optimization and scalability\n4. Security hardening and vulnerability assessment\n5. Documentation and developer experience\n\n## Project Context\n**%s** is a %s application scaffolded with Atempo, featuring built-in Docker development environment, AI-ready context system, and best-practice architecture.\n\n**Key Technologies**: %s, %s, Docker, %s\n**Architecture**: Clean %s architecture with service layer pattern\n**Development Environment**: Docker-based with hot reload and debugging support\n**Last Updated**: %s",
		projectName, projectName, framework, framework, framework, language, language, 
		language, framework, framework, analysis.ProjectType, 
		formatFeatures(analysis.Features), analysis.Architecture, 
		formatDatabaseTables(analysis.DatabaseTables), formatTechnicalStack(config.AIFeatures.TechnicalStack),
		projectName, framework, framework, language, getDatabaseType(framework), framework,
		time.Now().Format("2006-01-02"))

	return os.WriteFile(filepath.Join(aiDir, "context.md"), []byte(content), 0644)
}

// Helper functions for formatting
func formatFeatures(features []string) string {
	if len(features) == 0 {
		return "- Core application functionality"
	}
	
	result := ""
	for _, feature := range features {
		result += fmt.Sprintf("- %s\n", feature)
	}
	return result
}

func formatDatabaseTables(tables []string) string {
	if len(tables) == 0 {
		return "- Standard application tables"
	}
	
	result := ""
	for _, table := range tables {
		result += fmt.Sprintf("- %s\n", table)
	}
	return result
}

func formatTechnicalStack(stack []string) string {
	result := ""
	for _, item := range stack {
		result += fmt.Sprintf("- %s\n", item)
	}
	return result
}

func getDatabaseType(framework string) string {
	switch framework {
	case "laravel":
		return "MySQL"
	case "django":
		return "PostgreSQL"
	default:
		return "Database"
	}
}

// generateProjectOverview generates the project-overview.md file
func generateProjectOverview(aiDir, projectName, framework, language, version string, config *ContextConfig, analysis *ProjectAnalysis) error {
	content := fmt.Sprintf("# %s - Project Overview\n\n## Mission Statement\n%s is a %s application designed to provide [describe your project's purpose and goals]. Built with modern development practices and a focus on maintainability, scalability, and developer experience.\n\n## Technology Stack\n- **Backend Framework**: %s %s (%s)\n- **Database**: %s\n- **Development Environment**: Docker & Docker Compose\n- **Package Manager**: %s\n\n## Project Analysis\n\n### Detected Project Type\n**%s**\n\n### Core Features\n%s\n\n### Architecture Pattern\n**%s Architecture**\n\n%s\n\n## Development Workflow\n\n### Local Development\n%s\n\n### Key Development Commands\n%s\n\n## Current Project Status\n**Framework**: %s %s\n**Language**: %s\n**Development Status**: Active Development\n**Last Updated**: %s\n\n## Next Steps\n1. Complete core feature implementation\n2. Add comprehensive test coverage\n3. Implement proper error handling\n4. Optimize performance and scalability\n5. Add documentation and deployment guides",
		projectName, projectName, framework, framework, version, language, 
		getDatabaseType(framework), config.DevelopmentContext.PackageManager,
		analysis.ProjectType, formatFeatures(analysis.Features), analysis.Architecture,
		config.AIFeatures.FrameworkPatternsTemplate, generateDevelopmentInstructions(framework),
		generateKeyCommands(config.DevelopmentContext.Commands), framework, version, language,
		time.Now().Format("2006-01-02"))

	return os.WriteFile(filepath.Join(aiDir, "project-overview.md"), []byte(content), 0644)
}

// generateDevelopmentInstructions generates development setup instructions
func generateDevelopmentInstructions(framework string) string {
	switch framework {
	case "laravel":
		return `1. **Setup**: docker-compose up -d
2. **Dependencies**: docker-compose exec app composer install
3. **Database**: docker-compose exec app php artisan migrate --seed
4. **Testing**: docker-compose exec app php artisan test`
	case "django":
		return `1. **Setup**: docker-compose up -d
2. **Dependencies**: docker-compose exec web pip install -r requirements.txt
3. **Database**: docker-compose exec web python manage.py migrate
4. **Testing**: docker-compose exec web python manage.py test`
	default:
		return `1. **Setup**: docker-compose up -d
2. **Dependencies**: Install project dependencies
3. **Database**: Run database migrations
4. **Testing**: Run test suite`
	}
}

// generateKeyCommands generates key development commands
func generateKeyCommands(commands map[string]string) string {
	result := "```bash\n"
	for cmd, description := range commands {
		result += fmt.Sprintf("# %s\n%s\n\n", cmd, description)
	}
	result += "```"
	return result
}

// Stub functions for other generators (these would be implemented similarly)
func generateCodebaseMap(aiDir, projectName, framework, language string, config *ContextConfig, analysis *ProjectAnalysis) error {
	// Implementation would be similar to generateProjectOverview but focused on codebase structure
	return nil
}

func generateDevelopmentWorkflows(aiDir, projectName, framework, language string, config *ContextConfig) error {
	// Implementation would generate development workflow documentation
	return nil
}

func generatePatternsAndConventions(aiDir, projectName, framework, language string, config *ContextConfig) error {
	// Implementation would generate coding patterns and conventions
	return nil
}

func generateUIUXGuidelines(aiDir, projectName, framework, language string, config *ContextConfig) error {
	// Implementation would generate UI/UX guidelines
	return nil
}
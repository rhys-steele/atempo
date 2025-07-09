package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atempo/internal/scaffold"
	"atempo/internal/utils"
)

// scaffoldService implements ScaffoldService
type scaffoldService struct {
	templateService TemplateService
}

// NewScaffoldService creates a new ScaffoldService implementation
func NewScaffoldService(templateService TemplateService) ScaffoldService {
	return &scaffoldService{
		templateService: templateService,
	}
}

// CreateProject scaffolds a new project using the specified framework and version
func (s *scaffoldService) CreateProject(ctx context.Context, req CreateProjectRequest) error {
	// Validate framework and version first
	if err := s.ValidateFramework(ctx, req.Framework, req.Version); err != nil {
		return fmt.Errorf("framework validation failed: %w", err)
	}

	// Create scaffold configuration
	config := scaffold.Config{
		Framework:   req.Framework,
		Version:     req.Version,
		ProjectPath: req.ProjectPath,
		ProjectName: req.ProjectName,
		EnableAI:    req.EnableAI,
		AIManifest:  req.AIManifest,
		Interactive: req.Interactive,
	}

	// Run the scaffolding process
	if err := scaffold.Run(&config); err != nil {
		return fmt.Errorf("scaffolding failed: %w", err)
	}

	return nil
}

// ValidateFramework checks if a framework and version combination is supported
func (s *scaffoldService) ValidateFramework(ctx context.Context, framework, version string) error {
	// Check if template path exists
	templatePath := filepath.Join("templates", "frameworks", framework)
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("framework '%s' is not supported", framework)
	}

	// Load template configuration to check supported versions
	templateConfigPath := filepath.Join(templatePath, "atempo.json")
	templateConfig, err := s.templateService.LoadTemplate(ctx, templateConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load template configuration: %w", err)
	}

	// Validate framework-specific version constraints
	switch framework {
	case "laravel":
		return s.validateLaravelVersion(version)
	case "django":
		return s.validateDjangoVersion(version)
	default:
		// For other frameworks, just check if the template exists
		return s.templateService.ValidateTemplate(ctx, templateConfig)
	}
}

// validateLaravelVersion validates Laravel version constraints
func (s *scaffoldService) validateLaravelVersion(version string) error {
	// Laravel-specific version validation
	supportedVersions := []string{"8", "9", "10", "11"}
	majorVersion := strings.Split(version, ".")[0]
	
	for _, supported := range supportedVersions {
		if majorVersion == supported {
			return nil
		}
	}
	
	return fmt.Errorf("Laravel version %s is not supported. Supported versions: %s", 
		version, strings.Join(supportedVersions, ", "))
}

// validateDjangoVersion validates Django version constraints  
func (s *scaffoldService) validateDjangoVersion(version string) error {
	// Django-specific version validation
	versionFloat, err := utils.ParseSemanticVersion(version)
	if err != nil {
		return fmt.Errorf("invalid Django version format: %w", err)
	}
	
	if versionFloat < 3.0 {
		return fmt.Errorf("Django version %s is not supported. Minimum version: 3.0", version)
	}
	
	if versionFloat > 5.0 {
		return fmt.Errorf("Django version %s is not yet supported. Maximum version: 5.0", version)
	}
	
	return nil
}

// GetSupportedFrameworks returns a list of supported frameworks
func (s *scaffoldService) GetSupportedFrameworks(ctx context.Context) ([]FrameworkInfo, error) {
	frameworksPath := filepath.Join("templates", "frameworks")
	
	// Check if frameworks directory exists
	if _, err := os.Stat(frameworksPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("frameworks directory not found: %s", frameworksPath)
	}

	// Read frameworks directory
	entries, err := os.ReadDir(frameworksPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read frameworks directory: %w", err)
	}

	var frameworks []FrameworkInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		frameworkName := entry.Name()
		
		// Load template configuration
		templateConfigPath := filepath.Join(frameworksPath, frameworkName, "atempo.json")
		templateConfig, err := s.templateService.LoadTemplate(ctx, templateConfigPath)
		if err != nil {
			// Skip frameworks with invalid configuration
			continue
		}

		// Create framework info
		frameworkInfo := FrameworkInfo{
			Name:        frameworkName,
			Language:    templateConfig.Language,
			Description: s.getFrameworkDescription(frameworkName),
		}

		// Get supported versions for this framework
		frameworkInfo.SupportedVersions = s.getSupportedVersions(frameworkName)

		frameworks = append(frameworks, frameworkInfo)
	}

	return frameworks, nil
}

// getFrameworkDescription returns a description for the framework
func (s *scaffoldService) getFrameworkDescription(framework string) string {
	descriptions := map[string]string{
		"laravel": "The PHP Framework for Web Artisans",
		"django":  "The Web framework for perfectionists with deadlines",
		"react":   "A JavaScript library for building user interfaces",
		"vue":     "The Progressive JavaScript Framework",
		"angular": "Platform for building mobile and desktop web applications",
		"express": "Fast, unopinionated, minimalist web framework for Node.js",
	}

	if desc, exists := descriptions[framework]; exists {
		return desc
	}

	return fmt.Sprintf("Framework: %s", framework)
}

// getSupportedVersions returns supported versions for a framework
func (s *scaffoldService) getSupportedVersions(framework string) []string {
	versions := map[string][]string{
		"laravel": {"8", "9", "10", "11"},
		"django":  {"3.0", "3.1", "3.2", "4.0", "4.1", "4.2", "5.0"},
		"react":   {"17", "18"},
		"vue":     {"2", "3"},
		"angular": {"13", "14", "15", "16", "17"},
		"express": {"4"},
	}

	if supportedVersions, exists := versions[framework]; exists {
		return supportedVersions
	}

	// Default version if not specified
	return []string{"latest"}
}
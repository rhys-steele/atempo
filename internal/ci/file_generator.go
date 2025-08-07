package ci

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"atempo/internal/ci/providers"
	"atempo/internal/registry"
	"atempo/internal/utils"
)

// FileGenerator handles CI file generation and validation
type FileGenerator struct {
	providerRegistry *providers.ProviderRegistry
	templatesFS      embed.FS
}

// NewFileGenerator creates a new file generator
func NewFileGenerator(providerRegistry *providers.ProviderRegistry, templatesFS embed.FS) *FileGenerator {
	return &FileGenerator{
		providerRegistry: providerRegistry,
		templatesFS:      templatesFS,
	}
}

// GenerationResult represents the result of CI file generation
type GenerationResult struct {
	ConfigFile   string   `json:"config_file"`
	GeneratedAt  string   `json:"generated_at"`
	Provider     string   `json:"provider"`
	Framework    string   `json:"framework"`
	FilesCreated []string `json:"files_created"`
	Success      bool     `json:"success"`
	Error        string   `json:"error,omitempty"`
}

// GenerateFiles generates CI configuration files based on the provided config
func (fg *FileGenerator) GenerateFiles(ciConfig *CIConfig, projectPath string) (*GenerationResult, error) {
	result := &GenerationResult{
		Provider:     string(ciConfig.Provider),
		Framework:    ciConfig.Framework,
		FilesCreated: []string{},
		Success:      false,
	}

	// Get provider implementation  
	provider, err := fg.providerRegistry.Get(providers.CIProvider(ciConfig.Provider))
	if err != nil {
		result.Error = fmt.Sprintf("Provider not found: %v", err)
		return result, err
	}

	// Convert ciConfig to providers.CIConfig format
	providerConfig := &providers.CIConfig{
		Provider:    providers.CIProvider(ciConfig.Provider),
		Framework:   ciConfig.Framework,
		ProjectName: ciConfig.ProjectName,
		ProjectPath: ciConfig.ProjectPath,
		RepoURL:     ciConfig.RepoURL,
		Settings:    ciConfig.Settings,
	}

	// Generate CI configuration content
	configContent, err := provider.GenerateConfig(providerConfig, fg.templatesFS)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to generate configuration: %v", err)
		return result, err
	}

	// Determine output file path
	configFilePath := provider.GetConfigPath(projectPath)
	result.ConfigFile = configFilePath

	// Ensure directory exists
	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		result.Error = fmt.Sprintf("Failed to create directory %s: %v", configDir, err)
		return result, err
	}

	// Write configuration file
	if err := os.WriteFile(configFilePath, configContent, 0644); err != nil {
		result.Error = fmt.Sprintf("Failed to write configuration file: %v", err)
		return result, err
	}

	result.FilesCreated = append(result.FilesCreated, configFilePath)

	// Save CI configuration to project registry
	if err := fg.saveCIConfig(ciConfig, projectPath); err != nil {
		result.Error = fmt.Sprintf("Failed to save CI configuration: %v", err)
		return result, err
	}

	result.Success = true
	return result, nil
}

// ValidateFiles validates existing CI configuration files
func (fg *FileGenerator) ValidateFiles(projectPath string, provider CIProvider) (*ValidationResult, error) {
	providerImpl, err := fg.providerRegistry.Get(providers.CIProvider(provider))
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	configPath := providerImpl.GetConfigPath(projectPath)
	
	result := &ValidationResult{
		Provider:   provider,
		ConfigPath: configPath,
		Valid:      false,
	}

	// Check if file exists
	if !utils.FileExists(configPath) {
		result.Errors = append(result.Errors, fmt.Sprintf("Configuration file not found: %s", configPath))
		return result, nil
	}

	// Basic YAML/file validation would go here
	// For now, we'll do a simple existence check
	result.Valid = true
	return result, nil
}

// saveCIConfig saves the CI configuration to the project's atempo config
func (fg *FileGenerator) saveCIConfig(ciConfig *CIConfig, projectPath string) error {
	// Load registry
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}

	// Find project by path and save CI config
	for i, project := range reg.Projects {
		if project.Path == projectPath {
			// Convert to registry CIConfig format
			regConfig := &registry.CIConfig{
				Provider:      string(ciConfig.Provider),
				Framework:     ciConfig.Framework,
				ProjectName:   ciConfig.ProjectName,
				ProjectPath:   ciConfig.ProjectPath,
				RepoURL:       ciConfig.RepoURL,
				Settings:      ciConfig.Settings,
				CreatedAt:     ciConfig.CreatedAt,
				UpdatedAt:     ciConfig.UpdatedAt,
				LastRunAt:     ciConfig.LastRunAt,
				LastRunStatus: ciConfig.LastRunStatus,
			}

			reg.Projects[i].CIConfig = regConfig
			reg.Projects[i].CIStatus = "enabled"
			return reg.SaveRegistry()
		}
	}

	return fmt.Errorf("project not found in registry for path: %s", projectPath)
}

// GetProjectByPath finds a project in the registry by path
func (fg *FileGenerator) GetProjectByPath(projectPath string) (*registry.Project, error) {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to load project registry: %w", err)
	}

	for _, project := range reg.Projects {
		if project.Path == projectPath {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("project not found for path: %s", projectPath)
}

// GetProjectName attempts to determine project name from various sources
func (fg *FileGenerator) GetProjectName(projectPath string) string {
	// Try to get from registry first
	if project, err := fg.GetProjectByPath(projectPath); err == nil {
		return project.Name
	}

	// Fallback to directory name
	return filepath.Base(projectPath)
}

// CreateAtempoDirectory ensures the .atempo directory exists in the project
func (fg *FileGenerator) CreateAtempoDirectory(projectPath string) error {
	atempoDir := filepath.Join(projectPath, ".atempo")
	return os.MkdirAll(atempoDir, 0755)
}

// BackupExistingConfig creates a backup of existing CI configuration if it exists
func (fg *FileGenerator) BackupExistingConfig(projectPath string, provider CIProvider) error {
	providerImpl, err := fg.providerRegistry.Get(providers.CIProvider(provider))
	if err != nil {
		return err
	}

	configPath := providerImpl.GetConfigPath(projectPath)
	if !utils.FileExists(configPath) {
		return nil // No existing config to backup
	}

	// Create backup with timestamp
	backupPath := configPath + ".backup"
	
	// Read existing content
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read existing config: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}

// ValidateProjectStructure validates that the project has the expected structure for CI
func (fg *FileGenerator) ValidateProjectStructure(projectPath, framework string) error {
	switch framework {
	case "laravel":
		requiredFiles := []string{"composer.json", "artisan"}
		for _, file := range requiredFiles {
			if !utils.FileExists(filepath.Join(projectPath, file)) {
				return fmt.Errorf("required Laravel file not found: %s", file)
			}
		}
	case "django":
		requiredFiles := []string{"manage.py"}
		for _, file := range requiredFiles {
			if !utils.FileExists(filepath.Join(projectPath, file)) {
				return fmt.Errorf("required Django file not found: %s", file)
			}
		}
	case "express", "lambda-node":
		requiredFiles := []string{"package.json"}
		for _, file := range requiredFiles {
			if !utils.FileExists(filepath.Join(projectPath, file)) {
				return fmt.Errorf("required Node.js file not found: %s", file)
			}
		}
	}

	return nil
}

// ShowGenerationSummary displays a summary of the generation results
func (fg *FileGenerator) ShowGenerationSummary(result *GenerationResult) {
	if result.Success {
		fmt.Printf("\n%s✅ CI Configuration Generated Successfully!%s\n", ColorGreen, ColorReset)
		fmt.Printf("═══════════════════════════════════════════════════════════\n")
		fmt.Printf("%sProvider:%s %s\n", ColorCyan, ColorReset, result.Provider)
		fmt.Printf("%sFramework:%s %s\n", ColorCyan, ColorReset, result.Framework)
		fmt.Printf("%sConfig File:%s %s\n", ColorCyan, ColorReset, result.ConfigFile)
		
		fmt.Printf("\n%sFiles Created:%s\n", ColorCyan, ColorReset)
		for _, file := range result.FilesCreated {
			fmt.Printf("  • %s\n", file)
		}

		fmt.Printf("\n%sNext Steps:%s\n", ColorYellow, ColorReset)
		fmt.Printf("  1. Review the generated configuration file\n")
		fmt.Printf("  2. Commit and push to trigger CI pipeline\n")
		fmt.Printf("  3. Test locally with: %sci run%s\n", ColorCyan, ColorReset)
		fmt.Printf("  4. Validate configuration: %sci validate%s\n", ColorCyan, ColorReset)
	} else {
		fmt.Printf("\n%s❌ CI Configuration Generation Failed%s\n", ColorRed, ColorReset)
		fmt.Printf("═══════════════════════════════════════════════════════════\n")
		fmt.Printf("%sError:%s %s\n", ColorRed, ColorReset, result.Error)
	}
}
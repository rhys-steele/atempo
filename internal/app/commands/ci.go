package commands

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"atempo/internal/ci"
	"atempo/internal/ci/providers"
	"atempo/internal/registry"
	"atempo/internal/utils"
)

// Type aliases to work with both ci and providers types
type CIProvider = providers.CIProvider

// CICommand handles CI/CD operations for projects
type CICommand struct {
	*BaseCommand
	registry         *registry.Registry
	providerRegistry *providers.ProviderRegistry
	templatesFS      embed.FS
}

// NewCICommand creates a new CI command
func NewCICommand(ctx *CommandContext, templatesFS embed.FS) *CICommand {
	return &CICommand{
		BaseCommand: NewBaseCommand(
			"ci",
			"Manage CI/CD configurations for projects",
			"ci <subcommand> [options]",
			ctx,
		),
		providerRegistry: providers.NewProviderRegistry(nil), // No logger for now
		templatesFS:      templatesFS,
	}
}

// Execute runs the CI command
func (c *CICommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return c.showUsage()
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "init":
		return c.executeInit(ctx, subArgs)
	case "run":
		return c.executeRun(ctx, subArgs)
	case "validate":
		return c.executeValidate(ctx, subArgs)
	case "status":
		return c.executeStatus(ctx, subArgs)
	case "remove":
		return c.executeRemove(ctx, subArgs)
	case "providers":
		return c.executeProviders(ctx, subArgs)
	default:
		return fmt.Errorf("unknown CI subcommand: %s", subcommand)
	}
}

// showUsage displays the CI command usage
func (c *CICommand) showUsage() error {
	fmt.Printf(`CI Commands:
  ci init      Initialize CI configuration for the project
  ci run       Run CI pipeline locally using Docker
  ci validate  Validate CI configuration files
  ci status    Show CI configuration status
  ci remove    Remove CI configuration
  ci providers List available CI providers

Usage:
  atempo ci <command> [options]
  {project} ci <command> [options]

Examples:
  atempo ci init                    # Initialize CI for current project
  my-app ci init                    # Initialize CI for specific project
  my-app ci run                     # Run CI pipeline locally
  my-app ci validate                # Validate CI configuration
  my-app ci status                  # Show CI status
  atempo ci providers               # List available providers
`)
	return nil
}

// executeInit initializes CI configuration for a project
func (c *CICommand) executeInit(ctx context.Context, args []string) error {
	// Step 1: Determine project context
	projectPath, err := c.resolveProjectPath(args)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Create interactive prompter
	prompter := ci.NewInteractivePrompter()
	
	// Show header
	prompter.ShowHeader()

	// Step 2: Detect framework using existing utilities
	framework, err := utils.DetectFramework(projectPath)
	if err != nil {
		ShowWarning(fmt.Sprintf("Could not auto-detect framework: %v", err))
		framework = "unknown"
	}

	// Step 3: Get available providers
	providers := c.providerRegistry.List()
	if len(providers) == 0 {
		return fmt.Errorf("no CI providers available")
	}

	// Step 4: Interactive provider selection
	// Convert providers slice to ci.CIProvider slice
	ciProviders := make([]ci.CIProvider, len(providers))
	for i, p := range providers {
		ciProviders[i] = ci.CIProvider(p)
	}
	selectedProvider, err := prompter.PromptProviderSelection(ciProviders)
	if err != nil {
		return fmt.Errorf("provider selection failed: %w", err)
	}

	// Get provider implementation - convert ci.CIProvider to providers.CIProvider
	// Use string values directly to avoid import issues
	providerImpl, err := c.providerRegistry.Get(CIProvider(string(selectedProvider)))
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	// Step 5: Interactive framework selection/confirmation
	supportedFrameworks := providerImpl.SupportedFrameworks()
	selectedFramework, err := prompter.PromptFrameworkSelection(
		framework,
		supportedFrameworks,
	)
	if err != nil {
		return fmt.Errorf("framework selection failed: %w", err)
	}

	// Validate framework compatibility
	if err := c.providerRegistry.ValidateProviderFrameworkCombination(CIProvider(string(selectedProvider)), selectedFramework); err != nil {
		return fmt.Errorf("provider/framework compatibility error: %w", err)
	}

	// Step 6: Get default settings and prompt for customization
	defaultSettings := providerImpl.GetDefaultSettings(selectedFramework)
	settings, err := prompter.PromptSettings(selectedProvider, selectedFramework, defaultSettings)
	if err != nil {
		return fmt.Errorf("settings configuration failed: %w", err)
	}

	// Step 7: Create CI configuration
	projectName := c.getProjectName(projectPath)
	ciConfig := &ci.CIConfig{
		Provider:    selectedProvider,
		Framework:   selectedFramework,
		ProjectName: projectName,
		ProjectPath: projectPath,
		Settings:    settings,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Show configuration summary
	prompter.ShowConfigurationSummary(ciConfig)

	// Step 8: Confirm generation
	if !prompter.ConfirmGeneration() {
		fmt.Printf("\n%s⏺%s CI initialization cancelled.\n", ColorYellow, ColorReset)
		return nil
	}

	// Step 9: Generate CI files
	fileGenerator := ci.NewFileGenerator(c.providerRegistry, c.templatesFS)
	
	// Backup existing config if it exists
	if err := fileGenerator.BackupExistingConfig(projectPath, selectedProvider); err != nil {
		ShowWarning(fmt.Sprintf("Failed to backup existing config: %v", err))
	}

	// Generate files
	ShowWorking("Generating CI configuration files...")
	result, err := fileGenerator.GenerateFiles(ciConfig, projectPath)
	if err != nil {
		ShowError("CI file generation failed", err.Error())
		return err
	}

	// Show results
	fileGenerator.ShowGenerationSummary(result)

	return nil
}

// executeRun runs the CI pipeline locally
func (c *CICommand) executeRun(ctx context.Context, args []string) error {
	// For now, return a placeholder implementation
	fmt.Println("CI local execution is not yet implemented")
	fmt.Println("This feature will be available in a future release")
	return nil
}

// executeValidate validates the CI configuration
func (c *CICommand) executeValidate(ctx context.Context, args []string) error {
	// For now, return a placeholder implementation
	fmt.Println("CI validation is not yet implemented")
	fmt.Println("This feature will be available in a future release")
	return nil
}

// executeStatus shows CI configuration status
func (c *CICommand) executeStatus(ctx context.Context, args []string) error {
	// Load registry
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}

	// If no arguments, show status for current directory
	var projectName string
	if len(args) == 0 {
		// Try to determine project from current directory
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		
		// Find project by path
		for _, project := range reg.Projects {
			if strings.HasPrefix(currentDir, project.Path) {
				projectName = project.Name
				break
			}
		}
		
		if projectName == "" {
			return fmt.Errorf("no project found for current directory")
		}
	} else {
		projectName = args[0]
	}

	// Get project CI configuration
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project '%s' not found", projectName)
	}

	fmt.Printf("CI Status for Project: %s\n", project.Name)
	fmt.Printf("─────────────────────────────\n\n")

	if project.CIConfig == nil {
		fmt.Printf("✗ CI not configured\n")
		fmt.Printf("  Run '%s ci init' to set up CI/CD\n", project.Name)
		return nil
	}

	fmt.Printf("✓ CI configured\n")
	fmt.Printf("  Provider: %s\n", project.CIConfig.Provider)
	fmt.Printf("  Framework: %s\n", project.CIConfig.Framework)
	fmt.Printf("  Status: %s\n", project.CIStatus)
	fmt.Printf("  Created: %s\n", project.CIConfig.CreatedAt.Format("2006-01-02 15:04"))
	
	if project.CIConfig.LastRunAt != nil {
		fmt.Printf("  Last Run: %s (%s)\n", 
			project.CIConfig.LastRunAt.Format("2006-01-02 15:04"),
			project.CIConfig.LastRunStatus)
	}

	// Check if CI files exist
	configPath := filepath.Join(project.Path, ".github", "workflows", "ci.yml")
	if project.CIConfig.Provider == "gitlab" {
		configPath = filepath.Join(project.Path, ".gitlab-ci.yml")
	}

	if utils.FileExists(configPath) {
		fmt.Printf("  Config File: %s ✓\n", configPath)
	} else {
		fmt.Printf("  Config File: %s ✗\n", configPath)
	}

	return nil
}

// executeRemove removes CI configuration
func (c *CICommand) executeRemove(ctx context.Context, args []string) error {
	// For now, return a placeholder implementation
	fmt.Println("CI removal is not yet implemented")
	fmt.Println("This feature will be available in a future release")
	return nil
}

// executeProviders lists available CI providers and their supported frameworks
func (c *CICommand) executeProviders(ctx context.Context, args []string) error {
	fmt.Printf("Available CI Providers\n")
	fmt.Printf("─────────────────────\n\n")

	providers := c.providerRegistry.List()
	for _, providerName := range providers {
		provider, err := c.providerRegistry.Get(providerName)
		if err != nil {
			continue
		}

		fmt.Printf("✓ %s\n", string(providerName))
		
		frameworks := provider.SupportedFrameworks()
		fmt.Printf("  Supported frameworks: %s\n", strings.Join(frameworks, ", "))
		
		if len(args) > 0 && args[0] == "detailed" {
			// Show detailed info for each framework
			for _, framework := range frameworks {
				defaults := provider.GetDefaultSettings(framework)
				fmt.Printf("    %s defaults:\n", framework)
				for key, value := range defaults {
					fmt.Printf("      %s: %v\n", key, value)
				}
			}
		}
		fmt.Println()
	}

	if len(args) == 0 {
		fmt.Printf("Use 'atempo ci providers detailed' for detailed framework settings\n")
	}

	return nil
}

// Helper methods

// resolveProjectPath resolves the project path from arguments or current directory
func (c *CICommand) resolveProjectPath(args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		// Use provided path/project name
		return registry.ResolveProjectPath(args[0])
	}
	
	// Use current directory
	return os.Getwd()
}

// getProjectName determines the project name from path or registry
func (c *CICommand) getProjectName(projectPath string) string {
	// Try to find in registry first
	reg, err := registry.LoadRegistry()
	if err == nil {
		for _, project := range reg.Projects {
			if project.Path == projectPath {
				return project.Name
			}
		}
	}
	
	// Fallback to directory basename
	return filepath.Base(projectPath)
}
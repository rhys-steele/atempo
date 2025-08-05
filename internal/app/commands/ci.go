package commands

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atempo/internal/registry"
	"atempo/internal/utils"
)

// CICommand handles CI/CD operations for projects
type CICommand struct {
	*BaseCommand
	registry    *registry.Registry
	templatesFS embed.FS
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
		templatesFS: templatesFS,
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

Usage:
  atempo ci <command> [options]
  {project} ci <command> [options]

Examples:
  atempo ci init                    # Initialize CI for current project
  my-app ci init                    # Initialize CI for specific project
  my-app ci run                     # Run CI pipeline locally
  my-app ci validate                # Validate CI configuration
  my-app ci status                  # Show CI status
`)
	return nil
}

// executeInit initializes CI configuration for a project
func (c *CICommand) executeInit(ctx context.Context, args []string) error {
	// For now, return a placeholder implementation
	fmt.Println("CI initialization is not yet implemented")
	fmt.Println("This feature will be available in a future release")
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
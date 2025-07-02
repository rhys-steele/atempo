package commands

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atempo/internal/scaffold"
)

// StartCommand handles the 'start' command for scaffolding new projects
type StartCommand struct {
	*BaseCommand
	templatesFS   embed.FS
	mcpServersFS  embed.FS
}

// NewStartCommand creates a new start command
func NewStartCommand(ctx *CommandContext, templatesFS, mcpServersFS embed.FS) *StartCommand {
	return &StartCommand{
		BaseCommand: NewBaseCommand(
			"start",
			"Create a new project",
			"atempo start <framework>[:<version>] [project_name]",
			ctx,
		),
		templatesFS:  templatesFS,
		mcpServersFS: mcpServersFS,
	}
}

// Execute runs the start command
func (c *StartCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s\nExamples:\n  atempo start laravel my-app     # Laravel latest in ./my-app/\n  atempo start laravel:11 my-app  # Laravel 11 in ./my-app/\n  atempo start laravel            # Laravel latest in current directory", c.Usage())
	}

	// Parse framework and optional version
	frameworkArg := args[0]
	var framework, version string
	
	if strings.Contains(frameworkArg, ":") {
		parts := strings.Split(frameworkArg, ":")
		if len(parts) != 2 {
			return fmt.Errorf("error: expected format <framework>[:<version>]")
		}
		framework = parts[0]
		version = parts[1]
	} else {
		framework = frameworkArg
		version = c.getLatestVersion(framework)
	}

	// Parse optional project name
	var projectDir string
	var projectName string
	
	if len(args) >= 2 {
		// Project name specified - create directory
		projectName = args[1]
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectDir = filepath.Join(cwd, projectName)
		
		// Create project directory
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
		
		// Change to project directory
		if err := os.Chdir(projectDir); err != nil {
			return fmt.Errorf("failed to change to project directory: %w", err)
		}
	} else {
		// Use current directory
		var err error
		projectDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectName = filepath.Base(projectDir)
	}

	// Trigger the scaffold process
	err := scaffold.Run(framework, version, c.templatesFS, c.mcpServersFS)
	if err != nil {
		return fmt.Errorf("scaffold error: %w", err)
	}

	fmt.Println("✅ Project scaffolding complete.")
	return nil
}

// getLatestVersion returns the latest supported version for a framework
func (c *StartCommand) getLatestVersion(framework string) string {
	switch framework {
	case "laravel":
		return "11" // Laravel 11 is the latest LTS
	case "django":
		return "5"  // Django 5 is the latest major version
	default:
		return "latest"
	}
}
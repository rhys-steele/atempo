package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atempo/internal/registry"
	"atempo/internal/utils"
)

// DescribeCommand provides detailed project description using context
type DescribeCommand struct {
	*BaseCommand
}

// NewDescribeCommand creates a new describe command
func NewDescribeCommand(ctx *CommandContext) *DescribeCommand {
	return &DescribeCommand{
		BaseCommand: NewBaseCommand(
			"describe",
			utils.GetStandardDescription("describe"),
			utils.CreateStandardUsage("describe", utils.PatternWithProjectContext),
			ctx,
		),
	}
}

// Execute runs the describe command
func (c *DescribeCommand) Execute(ctx context.Context, args []string) error {
	// Resolve project path from arguments
	resolution, err := utils.ResolveProjectPathFromArgs(args)
	if err != nil {
		return err
	}
	projectPath := resolution.Path
	projectName := resolution.Name

	// Load project from registry if available
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	var project *registry.Project
	if len(args) >= 1 {
		// Try to find project by name first
		project, _ = reg.FindProject(projectName)
	}

	// If not found in registry, scan directory for atempo.json
	if project == nil {
		atempoJsonPath := filepath.Join(projectPath, "atempo.json")
		if !utils.FileExists(atempoJsonPath) {
			return fmt.Errorf("no atempo.json found in %s\nThis doesn't appear to be an Atempo project.\nRun 'atempo create <framework>' to create a new project", projectPath)
		}

		// Read atempo.json to get basic info
		project = &registry.Project{
			Name: filepath.Base(projectPath),
			Path: projectPath,
		}

		if content, err := os.ReadFile(atempoJsonPath); err == nil {
			var config struct {
				Framework string `json:"framework"`
				Version   string `json:"version"`
				Name      string `json:"name"`
			}
			if json.Unmarshal(content, &config) == nil {
				project.Framework = config.Framework
				project.Version = config.Version
				if config.Name != "" && !strings.Contains(config.Name, "{{") {
					project.Name = config.Name
				}
			}
		}
	}

	// Update project status if it's in the registry
	if len(args) >= 1 {
		reg.UpdateProjectStatus(project.Name)
		// Reload to get updated status
		if updatedProject, err := reg.FindProject(project.Name); err == nil {
			project = updatedProject
		}
	}

	c.displayProjectInfo(project)
	return nil
}

// displayProjectInfo displays comprehensive project information
func (c *DescribeCommand) displayProjectInfo(project *registry.Project) {
	fmt.Printf("Project Description: %s\n", project.Name)
	fmt.Println(strings.Repeat("=", 50))

	// Basic project information
	fmt.Printf("Name: %s\n", project.Name)
	fmt.Printf("Path: %s\n", project.Path)

	if project.Framework != "" {
		fmt.Printf("Framework: %s", project.Framework)
		if project.Version != "" {
			fmt.Printf(" %s", project.Version)
		}
		fmt.Println()
	}

	// Project status if available
	if project.Status != "" {
		statusIcon := utils.GetStatusIcon(project.Status)
		fmt.Printf("Status: %s %s\n", statusIcon, project.Status)
	}

	// Git information
	if project.GitBranch != "" {
		fmt.Printf("Git: %s", project.GitBranch)
		if project.GitStatus != "" && project.GitStatus != "clean" {
			fmt.Printf(" (%s)", project.GitStatus)
		}
		fmt.Println()
	}

	// URLs if available
	if len(project.URLs) > 0 {
		fmt.Printf("URLs: %s\n", strings.Join(project.URLs, ", "))
	}

	fmt.Println()

	// Services information
	if len(project.Services) > 0 {
		fmt.Println("Docker Services")
		fmt.Println(strings.Repeat("-", 30))
		for _, service := range project.Services {
			var serviceIcon string
			switch service.Status {
			case "running":
				serviceIcon = "✓"
			case "stopped":
				serviceIcon = "✗"
			default:
				serviceIcon = "⚠"
			}
			fmt.Printf("  %s %s", serviceIcon, service.Name)
			if service.URL != "" {
				fmt.Printf(" → %s", service.URL)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Port mappings
	if len(project.Ports) > 0 {
		fmt.Println("Port Mappings")
		fmt.Println(strings.Repeat("-", 30))
		for _, port := range project.Ports {
			fmt.Printf("  %s: localhost:%d → container:%d\n", port.Service, port.External, port.Internal)
		}
		fmt.Println()
	}

	// Quick actions
	fmt.Println("Quick Actions")
	fmt.Println(strings.Repeat("-", 30))
	if project.Status == "stopped" || project.Status == "no-docker" {
		fmt.Printf("  atempo docker up %s      # Start services\n", project.Name)
	} else if project.Status == "running" {
		fmt.Printf("  atempo docker down %s    # Stop services\n", project.Name)
		fmt.Printf("  atempo docker logs %s    # View logs\n", project.Name)
	}
	fmt.Printf("  atempo logs %s           # View setup logs\n", project.Name)
	fmt.Printf("  cd %s           # Navigate to project\n", project.Path)
}
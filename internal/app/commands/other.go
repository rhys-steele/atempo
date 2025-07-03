package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atempo/internal/compose"
	"atempo/internal/logger"
	"atempo/internal/registry"
	"atempo/internal/utils"
)

// ReconfigureCommand regenerates docker-compose.yml from atempo.json
type ReconfigureCommand struct {
	*BaseCommand
}

// NewReconfigureCommand creates a new reconfigure command
func NewReconfigureCommand(ctx *CommandContext) *ReconfigureCommand {
	return &ReconfigureCommand{
		BaseCommand: NewBaseCommand(
			"reconfigure",
			"Regenerate docker-compose.yml from atempo.json",
			"atempo reconfigure [project]",
			ctx,
		),
	}
}

// Execute runs the reconfigure command
func (c *ReconfigureCommand) Execute(ctx context.Context, args []string) error {
	var projectPath string
	
	if len(args) > 0 {
		resolvedPath, err := registry.ResolveProjectPath(args[0])
		if err != nil {
			return fmt.Errorf("failed to resolve project: %w", err)
		}
		projectPath = resolvedPath
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectPath = cwd
	}

	fmt.Printf("→ Regenerating docker-compose.yml from atempo.json in %s...\n", projectPath)
	
	if err := compose.GenerateDockerCompose(projectPath); err != nil {
		return fmt.Errorf("failed to regenerate docker-compose.yml: %w", err)
	}

	fmt.Println("✅ docker-compose.yml regenerated successfully!")
	return nil
}

// AddServiceCommand adds a predefined service to a project
type AddServiceCommand struct {
	*BaseCommand
}

// NewAddServiceCommand creates a new add-service command
func NewAddServiceCommand(ctx *CommandContext) *AddServiceCommand {
	return &AddServiceCommand{
		BaseCommand: NewBaseCommand(
			"add-service",
			"Add predefined services (minio, elasticsearch, etc.)",
			"atempo add-service <service_type> [project]",
			ctx,
		),
	}
}

// Execute runs the add-service command
func (c *AddServiceCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: atempo add-service <service_type> [project]")
		fmt.Println("\nAvailable services:")
		for _, service := range compose.ListPredefinedServices() {
			fmt.Printf("  %s\n", service)
		}
		return fmt.Errorf("service type required")
	}

	serviceType := args[0]
	var projectPath string

	if len(args) > 1 {
		resolvedPath, err := registry.ResolveProjectPath(args[1])
		if err != nil {
			return fmt.Errorf("failed to resolve project: %w", err)
		}
		projectPath = resolvedPath
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectPath = cwd
	}

	fmt.Printf("→ Adding %s service to project...\n", serviceType)
	
	if err := compose.AddPredefinedService(projectPath, serviceType); err != nil {
		return fmt.Errorf("failed to add service: %w", err)
	}

	fmt.Printf("✅ %s service added to atempo.json\n", serviceType)
	fmt.Println("Run 'atempo reconfigure' to update docker-compose.yml")
	return nil
}

// LogsCommand displays setup logs for a project
type LogsCommand struct {
	*BaseCommand
}

// NewLogsCommand creates a new logs command
func NewLogsCommand(ctx *CommandContext) *LogsCommand {
	return &LogsCommand{
		BaseCommand: NewBaseCommand(
			"logs",
			"View setup logs for a project",
			"atempo logs <project_name>",
			ctx,
		),
	}
}

// Execute runs the logs command
func (c *LogsCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: atempo logs <project_name>")
		fmt.Println("\nExample: atempo logs my-laravel-app")
		return fmt.Errorf("project name required")
	}

	projectName := args[0]

	// Get the latest log file for the project
	logFile, err := logger.GetLatestLogFile(projectName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nTip: Project logs are created during 'atempo create'. Available projects:")
		
		// Show available projects
		reg, regErr := registry.LoadRegistry()
		if regErr == nil {
			projects := reg.ListProjects()
			for _, project := range projects {
				fmt.Printf("  - %s\n", project.Name)
			}
		}
		return err
	}

	// Display the log file
	fmt.Printf("📄 Setup logs for project: %s\n", projectName)
	fmt.Printf("🔗 Log file: %s\n\n", logFile)

	// Read and display the file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	fmt.Print(string(content))

	// Show all available log files if there are multiple
	allLogs, err := logger.GetAllLogFiles(projectName)
	if err == nil && len(allLogs) > 1 {
		fmt.Printf("\nOther available logs for %s:\n", projectName)
		for i, logPath := range allLogs {
			if logPath != logFile {
				fmt.Printf("  %d. %s\n", i+1, logPath)
			}
		}
	}

	return nil
}

// DescribeCommand provides detailed project description using context
type DescribeCommand struct {
	*BaseCommand
}

// NewDescribeCommand creates a new describe command
func NewDescribeCommand(ctx *CommandContext) *DescribeCommand {
	return &DescribeCommand{
		BaseCommand: NewBaseCommand(
			"describe",
			"Show detailed project description and context",
			"atempo describe [project]",
			ctx,
		),
	}
}

// Execute runs the describe command
func (c *DescribeCommand) Execute(ctx context.Context, args []string) error {
	var projectPath string
	var projectName string
	
	// Parse optional project argument
	if len(args) >= 1 {
		// Try to resolve project by name or path
		resolvedPath, err := registry.ResolveProjectPath(args[0])
		if err != nil {
			return fmt.Errorf("failed to resolve project: %w", err)
		}
		projectPath = resolvedPath
		projectName = args[0]
	} else {
		// Use current directory
		var err error
		projectPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectName = filepath.Base(projectPath)
	}

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
	fmt.Printf("📋 Project Description: %s\n", project.Name)
	fmt.Println(strings.Repeat("=", 50))
	
	// Basic project information
	fmt.Printf("🏷️  Name: %s\n", project.Name)
	fmt.Printf("📁 Path: %s\n", project.Path)
	
	if project.Framework != "" {
		fmt.Printf("🛠️  Framework: %s", project.Framework)
		if project.Version != "" {
			fmt.Printf(" %s", project.Version)
		}
		fmt.Println()
	}

	// Project status if available
	if project.Status != "" {
		var statusIcon string
		switch project.Status {
		case "running":
			statusIcon = "🟢"
		case "partial":
			statusIcon = "🟡"
		case "stopped", "no-docker", "no-services":
			statusIcon = "🔴"
		case "docker-error":
			statusIcon = "❌"
		default:
			statusIcon = "❓"
		}
		fmt.Printf("⚡ Status: %s %s\n", statusIcon, project.Status)
	}

	// Git information
	if project.GitBranch != "" {
		fmt.Printf("🌿 Git: %s", project.GitBranch)
		if project.GitStatus != "" && project.GitStatus != "clean" {
			fmt.Printf(" (%s)", project.GitStatus)
		}
		fmt.Println()
	}

	// URLs if available
	if len(project.URLs) > 0 {
		fmt.Printf("🌐 URLs: %s\n", strings.Join(project.URLs, ", "))
	}

	fmt.Println()

	// Services information
	if len(project.Services) > 0 {
		fmt.Println("🐳 Docker Services")
		fmt.Println(strings.Repeat("-", 30))
		for _, service := range project.Services {
			var serviceIcon string
			switch service.Status {
			case "running":
				serviceIcon = "🟢"
			case "stopped":
				serviceIcon = "🔴"
			default:
				serviceIcon = "🟡"
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
		fmt.Println("🔌 Port Mappings")
		fmt.Println(strings.Repeat("-", 30))
		for _, port := range project.Ports {
			fmt.Printf("  %s: localhost:%d → container:%d\n", port.Service, port.External, port.Internal)
		}
		fmt.Println()
	}

	// Quick actions
	fmt.Println("💡 Quick Actions")
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

// RemoveCommand removes a project from the registry
type RemoveCommand struct {
	*BaseCommand
}

// NewRemoveCommand creates a new remove command
func NewRemoveCommand(ctx *CommandContext) *RemoveCommand {
	return &RemoveCommand{
		BaseCommand: NewBaseCommand(
			"remove",
			"Remove a project from the registry",
			"atempo remove <project_name>",
			ctx,
		),
	}
}

// Execute runs the remove command
func (c *RemoveCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s\nExample: atempo remove my-app", c.Usage())
	}

	projectName := args[0]

	// Load registry
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Check if project exists
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project '%s' not found in registry", projectName)
	}

	// Confirm with user
	fmt.Printf("Are you sure you want to remove project '%s'?\n", projectName)
	fmt.Printf("Path: %s\n", project.Path)
	fmt.Printf("Framework: %s %s\n", project.Framework, project.Version)
	fmt.Print("This will only remove it from the registry, not delete the files. [y/N]: ")

	var response string
	fmt.Scanln(&response)
	
	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("Cancelled.")
		return nil
	}

	// Remove project from registry
	err = reg.RemoveProject(projectName)
	if err != nil {
		return fmt.Errorf("failed to remove project: %w", err)
	}

	fmt.Printf("✅ Project '%s' removed from registry successfully!\n", projectName)
	fmt.Printf("💡 Project files at %s are still intact.\n", project.Path)
	
	return nil
}
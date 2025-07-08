package commands

import (
	"context"
	"fmt"
	"strings"

	"atempo/internal/registry"
)

// StatusCommand displays project status dashboard
type StatusCommand struct {
	*BaseCommand
}

// NewStatusCommand creates a new status command
func NewStatusCommand(ctx *CommandContext) *StatusCommand {
	return &StatusCommand{
		BaseCommand: NewBaseCommand(
			"status",
			"Show project dashboard with health status",
			"atempo status",
			ctx,
		),
	}
}

// Execute runs the status command
func (c *StatusCommand) Execute(ctx context.Context, args []string) error {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Check if specific project is requested
	if len(args) > 0 {
		return c.showProjectStatus(reg, args[0])
	}

	projects := reg.ListProjects()
	if len(projects) == 0 {
		fmt.Println("No projects registered yet.")
		fmt.Println("Projects are automatically registered when you run 'atempo create'")
		return nil
	}

	// Update all project statuses
	fmt.Print("Checking project status...")
	if err := reg.UpdateAllProjectsStatus(); err != nil {
		fmt.Printf(" failed: %v\n", err)
	} else {
		fmt.Println(" done")
	}

	// Reload registry to get updated statuses
	reg, err = registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to reload registry: %w", err)
	}

	projects = reg.ListProjects()

	fmt.Println("\nProject Status")
	fmt.Println(strings.Repeat("─", 80))

	runningCount := 0
	stoppedCount := 0
	errorCount := 0

	for i, project := range projects {
		if i > 0 {
			fmt.Println()
		}

		var status string
		var statusColor string
		switch project.Status {
		case "running":
			status = "running"
			statusColor = "\033[32m" // green
			runningCount++
		case "partial":
			status = "partial"
			statusColor = "\033[33m" // yellow
			runningCount++
		case "stopped", "no-docker", "no-services":
			status = "stopped"
			statusColor = "\033[31m" // red
			stoppedCount++
		case "docker-error":
			status = "error"
			statusColor = "\033[31m" // red
			errorCount++
		default:
			status = "unknown"
			statusColor = "\033[37m" // gray
			errorCount++
		}

		// Project header with colored status
		fmt.Printf("%s %s%s\033[0m\n", project.Name, statusColor, status)

		// Framework on separate line
		if project.Framework != "" {
			fmt.Printf("  Framework: %s", project.Framework)
			if project.Version != "" {
				fmt.Printf(" %s", project.Version)
			}
			fmt.Println()
		}

		// Git information (compact)
		if project.GitBranch != "" {
			fmt.Printf("  Branch: %s", project.GitBranch)
			if project.GitStatus != "" && project.GitStatus != "clean" {
				fmt.Printf(" (%s)", project.GitStatus)
			}
			fmt.Println()
		}

		// Services (running only, compact)
		if len(project.Services) > 0 {
			runningServices := []string{}
			for _, service := range project.Services {
				if service.Status == "running" {
					runningServices = append(runningServices, service.Name)
				}
			}
			if len(runningServices) > 0 {
				fmt.Printf("  Services: %s\n", strings.Join(runningServices, ", "))
			}
		}

		// URLs (clean format, filter valid URLs)
		if len(project.URLs) > 0 {
			validURLs := []string{}
			for _, url := range project.URLs {
				// Filter out invalid URLs (port 0, duplicates)
				if !strings.Contains(url, ":0") && !contains(validURLs, url) {
					validURLs = append(validURLs, url)
				}
			}
			if len(validURLs) > 0 {
				fmt.Printf("  URLs: %s\n", strings.Join(validURLs, ", "))
			}
		}

		// Path (compact)
		fmt.Printf("  Path: %s\n", project.Path)
	}

	// Summary footer
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("%d projects", len(projects))
	if runningCount > 0 {
		fmt.Printf(" • %d running", runningCount)
	}
	if stoppedCount > 0 {
		fmt.Printf(" • %d stopped", stoppedCount)
	}
	if errorCount > 0 {
		fmt.Printf(" • %d errors", errorCount)
	}
	fmt.Println()

	return nil
}

// showProjectStatus displays status for a specific project
func (c *StatusCommand) showProjectStatus(reg *registry.Registry, projectName string) error {
	// Update status for the specific project
	if err := reg.UpdateProjectStatus(projectName); err != nil {
		return fmt.Errorf("failed to update project status: %w", err)
	}

	// Find the project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project '%s' not found", projectName)
	}

	var status string
	var statusColor string
	switch project.Status {
	case "running":
		status = "running"
		statusColor = "\033[32m" // green
	case "partial":
		status = "partial"
		statusColor = "\033[33m" // yellow
	case "stopped", "no-docker", "no-services":
		status = "stopped"
		statusColor = "\033[31m" // red
	case "docker-error":
		status = "error"
		statusColor = "\033[31m" // red
	default:
		status = "unknown"
		statusColor = "\033[37m" // gray
	}

	fmt.Printf("Project: %s %s%s\033[0m\n", project.Name, statusColor, status)

	if project.Framework != "" {
		fmt.Printf("Framework: %s", project.Framework)
		if project.Version != "" {
			fmt.Printf(" %s", project.Version)
		}
		fmt.Println()
	}

	if project.GitBranch != "" {
		fmt.Printf("Branch: %s", project.GitBranch)
		if project.GitStatus != "" && project.GitStatus != "clean" {
			fmt.Printf(" (%s)", project.GitStatus)
		}
		fmt.Println()
	}

	if len(project.Services) > 0 {
		fmt.Println("\nServices:")
		for _, service := range project.Services {
			var serviceStatus string
			var serviceColor string
			switch service.Status {
			case "running":
				serviceStatus = "running"
				serviceColor = "\033[32m" // green
			case "stopped":
				serviceStatus = "stopped"
				serviceColor = "\033[31m" // red
			default:
				serviceStatus = "unknown"
				serviceColor = "\033[37m" // gray
			}
			fmt.Printf("  %-15s %s%s\033[0m", service.Name, serviceColor, serviceStatus)
			if service.URL != "" && !strings.Contains(service.URL, ":0") {
				fmt.Printf(" → %s", service.URL)
			}
			fmt.Println()
		}
	}

	if len(project.Ports) > 0 {
		fmt.Println("\nPorts:")
		seenPorts := make(map[string]bool)
		for _, port := range project.Ports {
			portKey := fmt.Sprintf("%s:%d→%d", port.Service, port.External, port.Internal)
			if !seenPorts[portKey] && port.External != 0 {
				fmt.Printf("  %-15s localhost:%d → container:%d\n",
					port.Service, port.External, port.Internal)
				seenPorts[portKey] = true
			}
		}
	}

	fmt.Printf("\nPath: %s\n", project.Path)

	return nil
}

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

	projects := reg.ListProjects()
	if len(projects) == 0 {
		fmt.Println("No projects registered yet.")
		fmt.Println("Projects are automatically registered when you run 'atempo start'")
		return nil
	}

	// Update all project statuses
	fmt.Print("🔄 Checking project status...")
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
	
	fmt.Println("\n🚀 Atempo Project Dashboard")
	fmt.Println(strings.Repeat("=", 50))

	runningCount := 0
	stoppedCount := 0
	errorCount := 0

	for _, project := range projects {
		var statusIcon string
		
		switch project.Status {
		case "running":
			statusIcon = "🟢"
			runningCount++
		case "partial":
			statusIcon = "🟡"
			runningCount++
		case "stopped", "no-docker", "no-services":
			statusIcon = "🔴"
			stoppedCount++
		case "docker-error":
			statusIcon = "❌"
			errorCount++
		default:
			statusIcon = "❓"
			errorCount++
		}

		fmt.Printf("\n%s %s (%s %s)\n", statusIcon, project.Name, project.Framework, project.Version)
		fmt.Printf("   📁 %s\n", project.Path)
		
		if project.GitBranch != "" {
			fmt.Printf("   🌿 %s", project.GitBranch)
			if project.GitStatus != "" && project.GitStatus != "clean" {
				fmt.Printf(" • %s", project.GitStatus)
			}
			fmt.Printf("\n")
		}

		if len(project.URLs) > 0 {
			fmt.Printf("   🌐 URLs: %s\n", strings.Join(project.URLs, ", "))
		}

		if len(project.Services) > 0 {
			fmt.Printf("   🐳 Services: ")
			serviceStrs := make([]string, len(project.Services))
			for i, service := range project.Services {
				var serviceIcon string
				switch service.Status {
				case "running":
					serviceIcon = "🟢"
				case "stopped":
					serviceIcon = "🔴"
				default:
					serviceIcon = "🟡"
				}
				serviceStrs[i] = fmt.Sprintf("%s %s", serviceIcon, service.Name)
			}
			fmt.Println(strings.Join(serviceStrs, ", "))
		}

		if len(project.Ports) > 0 {
			fmt.Printf("   🔌 Ports: ")
			portStrs := make([]string, len(project.Ports))
			for i, port := range project.Ports {
				portStrs[i] = fmt.Sprintf("%s:%d→%d", port.Service, port.External, port.Internal)
			}
			fmt.Println(strings.Join(portStrs, ", "))
		}
	}

	// Summary
	fmt.Printf("\n%s Summary\n", strings.Repeat("=", 50))
	fmt.Printf("Total Projects: %d\n", len(projects))
	if runningCount > 0 {
		fmt.Printf("🟢 Running: %d\n", runningCount)
	}
	if stoppedCount > 0 {
		fmt.Printf("🔴 Stopped: %d\n", stoppedCount)
	}
	if errorCount > 0 {
		fmt.Printf("❌ Errors: %d\n", errorCount)
	}

	fmt.Println("\n💡 Quick Actions:")
	fmt.Println("  atempo docker up [project]     # Start project services")
	fmt.Println("  atempo docker down [project]   # Stop project services")
	fmt.Println("  atempo docker logs [project]   # View service logs")
	fmt.Println("  atempo logs [project]          # View setup logs")

	return nil
}
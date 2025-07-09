package commands

import (
	"context"
	"fmt"

	"atempo/internal/docker"
	"atempo/internal/registry"
	"atempo/internal/utils"
)

// StopCommand stops all running projects
type StopCommand struct {
	*BaseCommand
}

// NewStopCommand creates a new stop command
func NewStopCommand(ctx *CommandContext) *StopCommand {
	return &StopCommand{
		BaseCommand: NewBaseCommand(
			"stop",
			utils.GetStandardDescription("stop"),
			utils.CreateStandardUsage("stop", utils.PatternSimple),
			ctx,
		),
	}
}

// Execute runs the stop command
func (c *StopCommand) Execute(ctx context.Context, args []string) error {
	// Load registry to get all projects
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	projects := reg.ListProjects()
	if len(projects) == 0 {
		fmt.Println("No projects found in registry.")
		return nil
	}

	// Filter for running projects
	runningProjects := []registry.Project{}
	for _, project := range projects {
		if project.Status == "running" || project.Status == "partial" {
			runningProjects = append(runningProjects, project)
		}
	}

	if len(runningProjects) == 0 {
		fmt.Println("No running projects found.")
		return nil
	}

	fmt.Printf("Found %d running project(s):\n", len(runningProjects))
	for _, project := range runningProjects {
		fmt.Printf("  • %s (%s)\n", project.Name, project.Status)
	}
	fmt.Println()

	// Stop each running project
	stoppedCount := 0
	for _, project := range runningProjects {
		fmt.Printf("→ Stopping %s...\n", project.Name)

		// Use docker-compose down to stop all services
		if err := docker.ExecuteCommand("down", project.Path, []string{}); err != nil {
			fmt.Printf("✗ Failed to stop %s: %v\n", project.Name, err)
			continue
		}

		fmt.Printf("✓ Stopped %s\n", project.Name)
		stoppedCount++
	}

	fmt.Printf("\nSuccessfully stopped %d of %d running projects.\n", stoppedCount, len(runningProjects))
	return nil
}
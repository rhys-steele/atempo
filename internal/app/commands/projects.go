package commands

import (
	"context"
	"fmt"

	"atempo/internal/registry"
)

// ProjectsCommand handles listing all registered projects
type ProjectsCommand struct {
	*BaseCommand
}

// NewProjectsCommand creates a new projects command
func NewProjectsCommand(ctx *CommandContext) *ProjectsCommand {
	return &ProjectsCommand{
		BaseCommand: NewBaseCommand(
			"projects",
			"List all registered projects",
			"atempo projects",
			ctx,
		),
	}
}

// Execute runs the projects command
func (c *ProjectsCommand) Execute(ctx context.Context, args []string) error {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	projects := reg.ListProjects()
	if len(projects) == 0 {
		fmt.Println("No projects registered yet.")
		fmt.Println("Projects are automatically registered when you run 'atempo create'")
		return nil
	}

	fmt.Println("Registered Atempo Projects:")
	fmt.Println()

	for _, project := range projects {
		fmt.Printf("  %s\n", project.Name)
		fmt.Printf("    Framework: %s %s\n", project.Framework, project.Version)
		fmt.Printf("    Path: %s\n", project.Path)
		fmt.Printf("    Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Println()
	}

	return nil
}

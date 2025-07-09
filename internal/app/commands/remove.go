package commands

import (
	"context"
	"fmt"
	"strings"

	"atempo/internal/registry"
	"atempo/internal/utils"
)

// RemoveCommand removes a project from the registry
type RemoveCommand struct {
	*BaseCommand
}

// NewRemoveCommand creates a new remove command
func NewRemoveCommand(ctx *CommandContext) *RemoveCommand {
	return &RemoveCommand{
		BaseCommand: NewBaseCommand(
			"remove",
			utils.GetStandardDescription("remove"),
			utils.CreateStandardUsage("remove", utils.PatternWithRequiredArgs, "<project_name>"),
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

	fmt.Printf("✓ Project '%s' removed from registry successfully!\n", projectName)
	fmt.Printf("Project files at %s are still intact.\n", project.Path)

	return nil
}
package commands

import (
	"context"
	"fmt"

	"atempo/internal/registry"
	"atempo/internal/utils"
)

// ResetCommand implements the reset command
type ResetCommand struct {
	*BaseCommand
}

// NewResetCommand creates a new reset command
func NewResetCommand(ctx *CommandContext) *ResetCommand {
	return &ResetCommand{
		BaseCommand: NewBaseCommand(
			"reset",
			"Delete all projects from filesystem and registry",
			"atempo reset [--confirm]",
			ctx,
		),
	}
}

// Execute runs the reset command
func (c *ResetCommand) Execute(ctx context.Context, args []string) error {
	// Check for --confirm flag
	skipConfirmation := false
	for _, arg := range args {
		if arg == "--confirm" {
			skipConfirmation = true
			break
		}
	}

	// Load registry to get all projects
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}

	// Get all projects
	projects := reg.ListProjects()
	if len(projects) == 0 {
		fmt.Println("✓ No projects found - registry is already empty")
		return nil
	}

	// Show warning and confirmation
	if !skipConfirmation {
		fmt.Printf("⚠️  DANGER: This will permanently delete ALL %d projects!\n\n", len(projects))
		fmt.Println("Projects to be deleted:")
		for _, project := range projects {
			fmt.Printf("  • %s (%s %s) - %s\n",
				project.Name,
				project.Framework,
				project.Version,
				project.Path)
		}
		fmt.Printf("\nThis will:\n")
		fmt.Printf("  • Move all project directories to Trash\n")
		fmt.Printf("  • Clear the entire atempo registry\n")
		fmt.Printf("  • This action CANNOT be undone!\n\n")
		fmt.Print("Type 'RESET' in ALL CAPS to confirm, or anything else to cancel: ")

		var response string
		fmt.Scanln(&response)

		if response != "RESET" {
			fmt.Println("✗ Cancelled - no projects were deleted.")
			return nil
		}
	}

	// Delete all projects
	failedDeletions := []string{}
	successCount := 0

	for _, project := range projects {
		fmt.Printf("Deleting %s...", project.Name)

		// Move project directory to trash
		if err := utils.MoveToTrash(project.Path); err != nil {
			fmt.Printf(" ✗ Failed to delete files: %v\n", err)
			failedDeletions = append(failedDeletions, fmt.Sprintf("%s (filesystem)", project.Name))
		} else {
			fmt.Printf(" ✓ Files moved to trash\n")
			successCount++
		}
	}

	// Clear the entire registry
	emptyRegistry := &registry.Registry{
		Projects: []registry.Project{},
		Version:  "1.0",
	}

	if err := emptyRegistry.SaveRegistry(); err != nil {
		failedDeletions = append(failedDeletions, "registry clear")
		fmt.Printf("✗ Failed to clear registry: %v\n", err)
	} else {
		fmt.Printf("✓ Registry cleared\n")
	}

	// Show results
	fmt.Printf("\n")
	if len(failedDeletions) == 0 {
		fmt.Printf("🎉 All %d projects successfully deleted!\n", successCount)
		fmt.Printf("  ⎿ Project files moved to Trash\n")
		fmt.Printf("  ⎿ Registry cleared\n")
	} else {
		fmt.Printf("⚠️  Reset completed with %d errors:\n", len(failedDeletions))
		for _, failed := range failedDeletions {
			fmt.Printf("  • %s\n", failed)
		}
		if successCount > 0 {
			fmt.Printf("\n✓ Successfully deleted %d projects\n", successCount)
		}
	}

	return nil
}

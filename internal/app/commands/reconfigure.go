package commands

import (
	"context"
	"fmt"

	"atempo/internal/compose"
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
			utils.GetStandardDescription("reconfigure"),
			utils.CreateStandardUsage("reconfigure", utils.PatternWithProjectContext),
			ctx,
		),
	}
}

// Execute runs the reconfigure command
func (c *ReconfigureCommand) Execute(ctx context.Context, args []string) error {
	// Resolve project path from arguments
	resolution, err := utils.ResolveProjectPathFromArgs(args)
	if err != nil {
		return err
	}
	projectPath := resolution.Path

	fmt.Printf("→ Regenerating docker-compose.yml from atempo.json in %s...\n", projectPath)

	if err := compose.GenerateDockerCompose(projectPath); err != nil {
		return fmt.Errorf("failed to regenerate docker-compose.yml: %w", err)
	}

	fmt.Println("✓ docker-compose.yml regenerated successfully!")
	return nil
}
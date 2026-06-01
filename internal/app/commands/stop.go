package commands

import (
	"context"
	"fmt"

	"atempo/internal/docker"
)

type StopCommand struct{ *BaseCommand }

func NewStopCommand(ctx *CommandContext) *StopCommand {
	cmd := &StopCommand{newBase(
		"stop",
		"Stop project containers",
		"atempo stop [project]",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo stop [project]

Runs 'docker compose down' for the project, stopping and removing containers.
Volumes and images are preserved.

Arguments:
  project   Registered project name or path (uses current directory if omitted)

Examples:
  atempo stop             Stop project in current directory
  atempo stop my-app      Stop registered project 'my-app'
`)
	return cmd
}

func (c *StopCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	projectPath, projectName, err := resolveProject(args, c.ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Stopping %s%s%s...\n", ColorBold, projectName, ColorReset)

	if err := docker.Down(projectPath); err != nil {
		return fmt.Errorf("stop %s: %w", projectName, err)
	}

	fmt.Printf("\n%s✓%s %s stopped\n\n", ColorGreen, ColorReset, projectName)
	return nil
}

package commands

import (
	"context"
	"fmt"

	"atempo/internal/docker"
)

type StartCommand struct{ *BaseCommand }

func NewStartCommand(ctx *CommandContext) *StartCommand {
	cmd := &StartCommand{newBase(
		"start",
		"Start project containers",
		"atempo start [project]",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo start [project]

Runs 'docker compose up --build -d' for the project. Images are (re)built if
needed, then all services are started in detached mode.

Arguments:
  project   Registered project name or path (uses current directory if omitted)

Examples:
  atempo start             Start project in current directory
  atempo start my-app      Start registered project 'my-app'
`)
	return cmd
}

func (c *StartCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	projectPath, projectName, err := resolveProject(args, c.ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Starting %s%s%s...\n", ColorBold, projectName, ColorReset)

	if err := docker.Up(projectPath); err != nil {
		return fmt.Errorf("start %s: %w", projectName, err)
	}

	fmt.Printf("\n%s✓%s %s is running\n", ColorGreen, ColorReset, projectName)
	fmt.Printf("  Run %satempo status %s%s to see service health\n\n",
		ColorCyan, projectName, ColorReset)
	return nil
}

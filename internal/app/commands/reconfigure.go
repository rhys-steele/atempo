package commands

import (
	"context"
	"fmt"
	"os"

	"atempo/internal/compose"
)

type ReconfigureCommand struct{ *BaseCommand }

func NewReconfigureCommand(ctx *CommandContext) *ReconfigureCommand {
	cmd := &ReconfigureCommand{newBase(
		"reconfigure",
		"Regenerate docker-compose.yml from atempo.json",
		"atempo reconfigure [project|path]",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo reconfigure [project|path]

Regenerates docker-compose.yml from the project's atempo.json. Ports are
re-probed and any conflicts are resolved automatically.

Arguments:
  project   Registered project name
  path      Explicit path to a project directory
  (none)    Uses the current directory

Example:
  atempo reconfigure my-app
`)
	return cmd
}

func (c *ReconfigureCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	projectPath, err := c.resolveProjectPath(args)
	if err != nil {
		return err
	}

	fmt.Printf("Regenerating docker-compose.yml for %s...\n", projectPath)

	assignments, err := compose.GenerateDockerCompose(projectPath)
	if err != nil {
		return fmt.Errorf("reconfigure: %w", err)
	}

	fmt.Printf("%s✓%s docker-compose.yml updated\n\n", ColorGreen, ColorReset)
	printPortAssignments(assignments)

	return nil
}

func (c *ReconfigureCommand) resolveProjectPath(args []string) (string, error) {
	if len(args) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get working directory: %w", err)
		}
		if !hasAtempoJSON(cwd) {
			return "", fmt.Errorf("no atempo.json in current directory — pass a project name or path")
		}
		return cwd, nil
	}

	// Try registry first.
	if p, err := c.ctx.Registry.FindProject(args[0]); err == nil {
		return p.Path, nil
	}

	// Fall back to treating the arg as a path.
	if hasAtempoJSON(args[0]) {
		return args[0], nil
	}

	return "", fmt.Errorf("project %q not found — run 'atempo projects' to see registered projects", args[0])
}

func printPortAssignments(assignments []compose.PortAssignment) {
	if len(assignments) == 0 {
		return
	}
	fmt.Printf("%sPorts%s\n", ColorBold, ColorReset)
	for _, a := range assignments {
		fmt.Printf("  %-20s localhost:%d → %d\n", a.Service, a.External, a.Internal)
	}
	fmt.Println()
}

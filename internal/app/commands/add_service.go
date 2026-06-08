package commands

import (
	"context"
	"fmt"
	"strings"

	"atempo/internal/compose"
)

type AddServiceCommand struct{ *BaseCommand }

func NewAddServiceCommand(ctx *CommandContext) *AddServiceCommand {
	cmd := &AddServiceCommand{newBase(
		"add-service",
		"Add a predefined service to the project",
		"atempo add-service <service> [project]",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo add-service <service> [project]

Adds a predefined service to the project's atempo.json and regenerates
docker-compose.yml with a free port assigned.

Available services:
  minio           Object storage (ports 9000, 9001)
  elasticsearch   Search engine (port 9200)
  rabbitmq        Message broker (ports 5672, 15672)
  mongodb         Document database (port 27017)

Arguments:
  service   Name of the service to add
  project   Registered project name or path (uses current directory if omitted)

Examples:
  atempo add-service redis              Add redis to project in current directory
  atempo add-service mongodb my-app    Add mongodb to registered project 'my-app'
`)
	return cmd
}

func (c *AddServiceCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}
	if len(args) == 0 {
		available := strings.Join(compose.ListPredefinedServices(), ", ")
		return fmt.Errorf("usage: %s\nAvailable services: %s", c.usage, available)
	}

	serviceType := args[0]

	// Project is optional second arg; default to CWD.
	projectArgs := args[1:]
	projectPath, projectName, err := resolveProject(projectArgs, c.ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Adding %s%s%s to %s%s%s...\n",
		ColorBold, serviceType, ColorReset,
		ColorBold, projectName, ColorReset)

	if err := compose.AddPredefinedService(projectPath, serviceType); err != nil {
		return fmt.Errorf("add-service: %w", err)
	}

	// Regenerate compose with updated service and reassigned ports.
	assignments, err := compose.GenerateDockerCompose(projectPath)
	if err != nil {
		return fmt.Errorf("regenerate docker-compose.yml: %w", err)
	}

	fmt.Printf("%s✓%s %s added to %s\n\n", ColorGreen, ColorReset, serviceType, projectName)
	printPortAssignments(assignments)
	fmt.Printf("  Run %satempo start %s%s to apply the changes\n\n",
		ColorCyan, projectName, ColorReset)

	return nil
}

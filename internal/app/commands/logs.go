package commands

import (
	"context"
	"fmt"

	"atempo/internal/docker"
)

type LogsCommand struct{ *BaseCommand }

func NewLogsCommand(ctx *CommandContext) *LogsCommand {
	cmd := &LogsCommand{newBase(
		"logs",
		"Stream container logs",
		"atempo logs [project] [service...]",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo logs [project] [service...]

Streams live container logs via 'docker compose logs -f'. Press Ctrl-C to stop.

Arguments:
  project    Registered project name or path (uses current directory if omitted)
  service    One or more service names to filter (optional)

Examples:
  atempo logs                        Stream all logs for current directory project
  atempo logs my-app                 Stream all logs for 'my-app'
  atempo logs my-app webserver       Stream only the webserver service logs
  atempo logs my-app app mysql       Stream app and mysql logs
`)
	return cmd
}

func (c *LogsCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	// First arg may be a project name or a service name — try project first.
	var projectPath, projectName string
	var services []string
	var err error

	if len(args) > 0 {
		projectPath, projectName, err = resolveProject(args[:1], c.ctx)
		if err != nil {
			// Not a project — treat everything as service filters in CWD.
			projectPath, projectName, err = resolveProject(nil, c.ctx)
			if err != nil {
				return err
			}
			services = args
		} else {
			services = args[1:]
		}
	} else {
		projectPath, projectName, err = resolveProject(nil, c.ctx)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Streaming logs for %s%s%s (Ctrl-C to stop)...\n\n",
		ColorBold, projectName, ColorReset)

	return docker.Logs(projectPath, services)
}

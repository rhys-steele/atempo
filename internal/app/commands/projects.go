package commands

import (
	"context"
	"fmt"
	"strings"
)

type ProjectsCommand struct{ *BaseCommand }

func NewProjectsCommand(ctx *CommandContext) *ProjectsCommand {
	cmd := &ProjectsCommand{newBase(
		"projects",
		"List registered projects",
		"atempo projects",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo projects

Lists all projects registered with Atempo, showing framework, status, and path.
Projects are registered automatically when created with 'atempo create'.

Registry location: ~/.atempo/registry.json
`)
	return cmd
}

func (c *ProjectsCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	projects, err := c.ctx.Registry.ListProjects()
	if err != nil {
		return fmt.Errorf("load registry: %w", err)
	}

	if len(projects) == 0 {
		fmt.Printf("\nNo projects registered yet.\n")
		fmt.Printf("Run %satempo create laravel my-app%s to get started.\n\n", ColorCyan, ColorReset)
		return nil
	}

	fmt.Printf("\n%sProjects%s\n", ColorBold, ColorReset)
	fmt.Println(strings.Repeat("─", 52))

	for _, p := range projects {
		status, color := statusLabel(p.Status)
		fw := p.Framework
		if p.Version != "" {
			fw += " " + p.Version
		}

		fmt.Printf("%s%-24s%s %s%-12s%s %s%s%s\n",
			ColorBold, p.Name, ColorReset,
			color, status, ColorReset,
			ColorGray, fw, ColorReset,
		)

		if len(p.Ports) > 0 {
			for _, port := range p.Ports {
				if port.External != 0 {
					fmt.Printf("  %s%s%s → localhost:%d\n",
						ColorGray, port.Service, ColorReset, port.External)
				}
			}
		}

		fmt.Printf("  %s%s%s\n", ColorGray, p.Path, ColorReset)
		fmt.Println()
	}

	fmt.Printf("%d project(s) registered\n\n", len(projects))
	return nil
}

// statusLabel returns a display label and colour for a project status string.
func statusLabel(status string) (string, string) {
	switch status {
	case "running":
		return "running", ColorGreen
	case "stopped", "":
		return "stopped", ColorGray
	case "healthy":
		return "healthy", ColorGreen
	case "unhealthy":
		return "unhealthy", ColorRed
	default:
		return status, ColorYellow
	}
}

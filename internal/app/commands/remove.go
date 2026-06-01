package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type RemoveCommand struct{ *BaseCommand }

func NewRemoveCommand(ctx *CommandContext) *RemoveCommand {
	cmd := &RemoveCommand{newBase(
		"remove",
		"Remove a project from the registry",
		"atempo remove <project>",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo remove <project>

Removes a project from the Atempo registry. Project files on disk are not
deleted — only the registry entry is removed.

To stop containers before removing, run 'atempo stop <project>' first.

Example:
  atempo remove my-app
`)
	return cmd
}

func (c *RemoveCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}
	if len(args) < 1 {
		return fmt.Errorf("project name required\nUsage: atempo remove <project>")
	}

	name := args[0]

	project, err := c.ctx.Registry.FindProject(name)
	if err != nil {
		return fmt.Errorf("project %q not found — run 'atempo projects' to see registered projects", name)
	}

	fmt.Printf("\nRemove %s%s%s from registry?\n\n", ColorBold, name, ColorReset)
	fmt.Printf("  Framework  %s %s\n", project.Framework, project.Version)
	fmt.Printf("  Path       %s\n", project.Path)
	fmt.Printf("\n  Files on disk will %snot%s be deleted.\n\n", ColorBold, ColorReset)
	fmt.Printf("Confirm [y/N]: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.TrimSpace(scanner.Text())

	if !strings.EqualFold(response, "y") && !strings.EqualFold(response, "yes") {
		fmt.Println("Cancelled.")
		return nil
	}

	if err := c.ctx.Registry.RemoveProject(name); err != nil {
		return fmt.Errorf("remove project: %w", err)
	}

	fmt.Printf("\n%s✓%s %s removed from registry.\n", ColorGreen, ColorReset, name)
	fmt.Printf("Files at %s are still intact.\n\n", project.Path)
	return nil
}

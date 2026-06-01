package commands

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"atempo/internal/scaffold"
)

type CreateCommand struct {
	*BaseCommand
	templatesFS embed.FS
}

func NewCreateCommand(ctx *CommandContext, templatesFS embed.FS) *CreateCommand {
	cmd := &CreateCommand{
		BaseCommand: newBase(
			"create",
			"Scaffold a new project",
			"atempo create <framework>[:<version>] [name]",
			ctx,
		),
		templatesFS: templatesFS,
	}
	cmd.setHelp(`Usage: atempo create <framework>[:<version>] [name]

Scaffolds a new project: copies template files, runs the framework installer
(if applicable), assigns free ports, generates docker-compose.yml, and
registers the project.

Arguments:
  framework   laravel, fastapi, bun, python
  version     Optional version hint (default: latest)
  name        Project directory name (default: current directory name)

Examples:
  atempo create laravel my-app          Laravel (latest) in ./my-app/
  atempo create laravel:12 my-app       Laravel 12 in ./my-app/
  atempo create fastapi my-api          FastAPI in ./my-api/
  atempo create bun my-site             Bun in ./my-site/
  atempo create python my-worker        Plain Python in ./my-worker/
`)
	return cmd
}

func (c *CreateCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}
	if len(args) == 0 {
		return fmt.Errorf("usage: %s", c.usage)
	}

	framework, version := parseFrameworkArg(args[0])

	// Resolve project directory and name.
	var projectPath, projectName string
	if len(args) >= 2 {
		rawName := args[1]
		if filepath.IsAbs(rawName) {
			projectPath = filepath.Clean(rawName)
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get working directory: %w", err)
			}
			projectPath = filepath.Join(cwd, rawName)
		}
		projectName = filepath.Base(projectPath)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		projectPath = cwd
		projectName = filepath.Base(cwd)
	}

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("create project directory: %w", err)
	}

	// Strip the "templates/" prefix from the embedded FS so scaffold sees
	// paths like "frameworks/laravel/..." directly.
	templatesSubFS, err := fs.Sub(c.templatesFS, "templates")
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	fmt.Printf("\nCreating %s%s%s", ColorBold, projectName, ColorReset)
	if version != "" {
		fmt.Printf(" (%s %s)", framework, version)
	} else {
		fmt.Printf(" (%s)", framework)
	}
	fmt.Printf("\n%s\n\n", strings.Repeat("─", 50))

	fmt.Printf("  %s→%s Copying template files...\n", ColorCyan, ColorReset)

	opts := scaffold.Options{
		Framework:   framework,
		Version:     version,
		ProjectPath: projectPath,
		ProjectName: projectName,
		TemplatesFS: templatesSubFS,
		Registry:    c.ctx.Registry,
	}

	assignments, err := scaffold.Run(opts)
	if err != nil {
		return err
	}

	fmt.Printf("\n%s✓%s %s created at %s\n", ColorGreen, ColorReset, projectName, projectPath)
	printPortAssignments(assignments)
	fmt.Printf("\n  Run %satempo start %s%s to start the project\n\n",
		ColorCyan, projectName, ColorReset)

	return nil
}

// parseFrameworkArg splits "framework:version" into its parts.
// "laravel" → ("laravel", ""), "laravel:12" → ("laravel", "12")
func parseFrameworkArg(arg string) (framework, version string) {
	if i := strings.IndexByte(arg, ':'); i >= 0 {
		return arg[:i], arg[i+1:]
	}
	return arg, ""
}


package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

// resolveProject resolves a project path and display name from CLI args.
//
// Resolution order:
//  1. args[0] is a registered project name → use registry path
//  2. args[0] is an existing directory with atempo.json → use it as path
//  3. args is empty → use current working directory
func resolveProject(args []string, ctx *CommandContext) (projectPath, projectName string, err error) {
	if len(args) == 0 {
		return resolveFromCWD()
	}

	target := args[0]

	// Check registry first.
	if ctx.Registry != nil {
		if p, err := ctx.Registry.FindProject(target); err == nil {
			return p.Path, p.Name, nil
		}
	}

	// Treat as a filesystem path.
	abs, err := filepath.Abs(target)
	if err != nil {
		return "", "", fmt.Errorf("resolve path %q: %w", target, err)
	}
	if hasAtempoJSON(abs) {
		return abs, filepath.Base(abs), nil
	}

	return "", "", fmt.Errorf(
		"project %q not found — run 'atempo projects' to see registered projects",
		target,
	)
}

func resolveFromCWD() (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("get working directory: %w", err)
	}
	if !hasAtempoJSON(cwd) {
		return "", "", fmt.Errorf(
			"no atempo.json in current directory — pass a project name or run from a project directory",
		)
	}
	return cwd, filepath.Base(cwd), nil
}

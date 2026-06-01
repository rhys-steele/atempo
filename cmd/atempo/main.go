package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"

	"atempo/internal/app/commands"
	"atempo/internal/errors"
	"atempo/internal/interfaces"
	"atempo/internal/repositories"
)

//go:embed all:templates
var templatesFS embed.FS

func main() {
	repo := buildRegistry()
	registry := commands.NewCommandRegistry(repo, templatesFS)

	// No arguments — launch interactive shell
	if len(os.Args) < 2 {
		ctx := context.Background()
		if err := registry.Execute(ctx, "shell", []string{}); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting shell: %v\n", err)
			os.Exit(1)
		}
		return
	}

	name := os.Args[1]
	args := os.Args[2:]

	if commands.IsHelpCommand(name) {
		registry.ShowUsage()
		return
	}

	ctx := context.Background()
	if err := registry.Execute(ctx, name, args); err != nil {
		printError(err, name, registry)
		os.Exit(1)
	}
}

// buildRegistry wires the repository layer used by all commands.
func buildRegistry() interfaces.RegistryRepository {
	cache := repositories.NewMemoryCacheRepository()
	return repositories.NewFileRegistryRepository(cache)
}

// printError formats and prints an error to stderr with recovery hints.
func printError(err error, commandName string, registry *commands.CommandRegistry) {
	fmt.Fprintln(os.Stderr)

	if ae, ok := err.(*errors.AtempoError); ok {
		fmt.Fprintf(os.Stderr, "%s\n", errors.FormatUserMessage(ae))
		if hint := errors.GetRecoveryAction(ae); hint != "" {
			fmt.Fprintf(os.Stderr, "\n• %s\n", hint)
		}
		if details := errors.GetDetails(ae); len(details) > 0 {
			if projects, ok := details["available_projects"].([]string); ok && len(projects) > 0 {
				fmt.Fprintf(os.Stderr, "\nAvailable projects: %s\n", formatList(projects))
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	fmt.Fprintln(os.Stderr)
}

func formatList(items []string) string {
	if len(items) <= 3 {
		return strings.Join(items, ", ")
	}
	return fmt.Sprintf("%s ... (%d total)", strings.Join(items[:3], ", "), len(items))
}

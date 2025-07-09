package commands

import (
	"context"
	"fmt"
	"os"

	"atempo/internal/logger"
	"atempo/internal/registry"
)

// LogsCommand displays setup logs for a project
type LogsCommand struct {
	*BaseCommand
}

// NewLogsCommand creates a new logs command
func NewLogsCommand(ctx *CommandContext) *LogsCommand {
	return &LogsCommand{
		BaseCommand: NewBaseCommand(
			"logs",
			"View setup logs for a project",
			"atempo logs <project_name>",
			ctx,
		),
	}
}

// Execute runs the logs command
func (c *LogsCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: atempo logs <project_name>")
		fmt.Println("\nExample: atempo logs my-laravel-app")
		return fmt.Errorf("project name required")
	}

	projectName := args[0]

	// Get the latest log file for the project
	logFile, err := logger.GetLatestLogFile(projectName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nTip: Project logs are created during 'atempo create'. Available projects:")

		// Show available projects
		reg, regErr := registry.LoadRegistry()
		if regErr == nil {
			projects := reg.ListProjects()
			for _, project := range projects {
				fmt.Printf("  - %s\n", project.Name)
			}
		}
		return err
	}

	// Display the log file
	fmt.Printf("Setup logs for project: %s\n", projectName)
	fmt.Printf("Log file: %s\n\n", logFile)

	// Read and display the file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	fmt.Print(string(content))

	// Show all available log files if there are multiple
	allLogs, err := logger.GetAllLogFiles(projectName)
	if err == nil && len(allLogs) > 1 {
		fmt.Printf("\nOther available logs for %s:\n", projectName)
		for i, logPath := range allLogs {
			if logPath != logFile {
				fmt.Printf("  %d. %s\n", i+1, logPath)
			}
		}
	}

	return nil
}
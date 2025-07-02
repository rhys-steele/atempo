package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// ShellCommand provides an interactive shell interface
type ShellCommand struct {
	ctx      *CommandContext
	registry *CommandRegistry
}

// NewShellCommand creates a new shell command
func NewShellCommand(ctx *CommandContext, registry *CommandRegistry) *ShellCommand {
	return &ShellCommand{
		ctx:      ctx,
		registry: registry,
	}
}

// Name returns the command name
func (c *ShellCommand) Name() string {
	return "shell"
}

// Description returns the command description
func (c *ShellCommand) Description() string {
	return "Start interactive shell mode"
}

// Usage returns the command usage
func (c *ShellCommand) Usage() string {
	return "shell"
}

// Execute runs the shell command
func (c *ShellCommand) Execute(ctx context.Context, args []string) error {
	c.showWelcome()
	return c.runInteractiveLoop(ctx)
}

// showWelcome displays the welcome screen with ASCII art and tips
func (c *ShellCommand) showWelcome() {

	// Display ASCII art
	fmt.Println()
	fmt.Println("")
	fmt.Println()

	fmt.Println("╭──────────────────────────────────────────────────────────╮")
	fmt.Println("│                                                          │")
	fmt.Println("│    █████╗ ████████╗███████╗███╗   ███╗██████╗  ██████╗   |")
	fmt.Println("│   ██╔══██╗╚══██╔══╝██╔════╝████╗ ████║██╔══██╗██╔═══██╗  │")
	fmt.Println("│   ███████║   ██║   █████╗  ██╔████╔██║██████╔╝██║   ██║  │")
	fmt.Println("│   ██╔══██║   ██║   ██╔══╝  ██║╚██╔╝██║██╔═══╝ ██║   ██║  │")
	fmt.Println("│   ██║  ██║   ██║   ███████╗██║ ╚═╝ ██║██║     ╚██████╔╝  │")
	fmt.Println("│   ╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝     ╚═╝╚═╝      ╚═════╝   │")
	fmt.Println("│                                                          │")
	fmt.Println("│   Type 'help' for help, 'exit' to quit                   │")
	fmt.Println("│                                                          │")
	fmt.Println("╰──────────────────────────────────────────────────────────╯")
	fmt.Println()

	c.showTips()
}

// showTips displays helpful tips for getting started
func (c *ShellCommand) showTips() {
	fmt.Println(" Tips for getting started:")
	fmt.Println()
	fmt.Println(" 1. Create a new project: create laravel my-app")
	fmt.Println(" 2. Check project status: status")
	fmt.Println(" 3. Start Docker services: docker up")
	fmt.Println(" 4. View all projects: projects")
	fmt.Println(" 5. Get project details: describe [project-name]")
	fmt.Println()
	fmt.Println(" ※ Tip: All regular atempo commands work here without 'atempo' prefix")
	fmt.Println()
}

// runInteractiveLoop handles the main interactive loop
func (c *ShellCommand) runInteractiveLoop(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Simple prompt
		fmt.Print("atempo > ")
		
		if !scanner.Scan() {
			// EOF or error
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		// Handle built-in shell commands
		if c.handleBuiltinCommand(input) {
			continue
		}

		// Parse command and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		commandName := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		// Execute the command with status indicators
		c.executeCommandWithStatus(ctx, commandName, args)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

// handleBuiltinCommand handles shell-specific commands
func (c *ShellCommand) handleBuiltinCommand(input string) bool {
	switch strings.ToLower(input) {
	case "exit", "quit", "q":
		ShowInfo("Shutting down Atempo shell...")
		fmt.Println("Goodbye!")
		os.Exit(0)
		return true
	case "clear", "cls":
		ShowWorking("Clearing screen...")
		fmt.Print("\033[2J\033[H") // ANSI clear screen
		c.showWelcome()
		ShowSuccess("Screen cleared", "Welcome screen refreshed")
		return true
	case "help":
		ShowThinking("Loading command help...")
		c.registry.ShowUsage()
		ShowSuccess("Help displayed", "All available commands shown")
		return true
	case "tips":
		ShowInfo("Showing helpful tips...")
		c.showTips()
		return true
	}
	return false
}


// executeCommandWithStatus executes a command with Claude Code-style status indicators
func (c *ShellCommand) executeCommandWithStatus(ctx context.Context, commandName string, args []string) {
	// Check if command exists
	if !c.registry.HasCommand(commandName) {
		ShowError(fmt.Sprintf("Unknown command: %s", commandName), "Type 'help' to see available commands")
		return
	}
	
	// Show thinking indicator
	ShowThinking(fmt.Sprintf("Executing %s...", commandName))
	
	// Execute the command
	err := c.registry.Execute(ctx, commandName, args)
	
	if err != nil {
		ShowError(fmt.Sprintf("Command failed: %s", commandName), err.Error())
	} else {
		ShowSuccess(fmt.Sprintf("Completed %s", commandName), fmt.Sprintf("Successfully executed %s %s", commandName, strings.Join(args, " ")))
	}
}


package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

// ShellCommand provides an interactive shell interface
type ShellCommand struct {
	ctx           *CommandContext
	registry      *CommandRegistry
	cachedPrompt  string
	lastDir       string
	lastGitBranch string
	lastGitStatus string
}

// NewShellCommand creates a new shell command
func NewShellCommand(ctx *CommandContext, registry *CommandRegistry) *ShellCommand {
	cmd := &ShellCommand{
		ctx:      ctx,
		registry: registry,
	}
	cmd.cachedPrompt = fmt.Sprintf("%s>%s ", ColorPurple, ColorReset) // Simple arrow prompt
	return cmd
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

	// c.showTips()
}

// showTips displays helpful tips for getting started
func (c *ShellCommand) showTips() {
	fmt.Printf("%sTips for using atempo shell:%s\n\n", ColorCyan, ColorReset)
	fmt.Println(" Shell Features:")
	fmt.Println(" • Use ↑/↓ arrow keys to navigate command history")
	fmt.Println(" • Use Tab key for command auto-completion")
	fmt.Println(" • Command history is saved to ~/.atempo_history")
	fmt.Println()
	fmt.Println(" Getting Started:")
	fmt.Println(" 1. Create a new project: create laravel my-app")
	fmt.Println(" 2. View all projects: projects")
	fmt.Println()
	fmt.Println(" Project Commands (use: {project} {command}):")
	fmt.Println(" • my-app up          Start services")
	fmt.Println(" • my-app down        Stop services")
	fmt.Println(" • my-app status      Check status")
	fmt.Println(" • my-app logs        View logs")
	fmt.Println(" • my-app shell       Enter container")
	fmt.Println(" • my-app code        Open in VS Code")
	fmt.Println()
}

// runInteractiveLoop handles the main interactive loop
func (c *ShellCommand) runInteractiveLoop(ctx context.Context) error {
	// Print initial directory info
	c.printDirectoryInfo()
	
	// Configure readline with history and completion
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          c.cachedPrompt, // Use cached prompt
		HistoryFile:     os.ExpandEnv("$HOME/.atempo_history"),
		AutoComplete:    c.createAutoCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			// Handle EOF (Ctrl+D) or interrupt
			if err == readline.ErrInterrupt {
				continue
			}
			break
		}

		input := strings.TrimSpace(line)

		if input == "" {
			// Don't update prompt on empty input - just continue
			continue
		}

		// Handle built-in shell commands
		if c.handleBuiltinCommand(input) {
			// Print directory info after built-in commands
			c.printDirectoryInfo()
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

		// Execute the command with status indicators (includes bash passthrough)
		c.executeCommandWithStatus(ctx, commandName, args)
		
		// Print directory info after command execution
		c.printDirectoryInfo()
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

// executeCommandWithStatus executes a command with bash passthrough support
func (c *ShellCommand) executeCommandWithStatus(ctx context.Context, commandName string, args []string) {
	// First, try atempo commands (global or project commands)
	if c.registry.HasCommand(commandName) || c.registry.IsProjectName(commandName) {
		// Show thinking indicator for atempo commands
		if c.registry.IsProjectName(commandName) && len(args) > 0 {
			ShowThinking(fmt.Sprintf("Executing %s %s...", commandName, args[0]))
		} else {
			ShowThinking(fmt.Sprintf("Executing %s...", commandName))
		}

		// Execute the atempo command
		err := c.registry.Execute(ctx, commandName, args)

		if err != nil {
			if c.registry.IsProjectName(commandName) && len(args) > 0 {
				ShowError(fmt.Sprintf("Command failed: %s %s", commandName, args[0]), err.Error())
			} else {
				ShowError(fmt.Sprintf("Command failed: %s", commandName), err.Error())
			}
		} else {
			if c.registry.IsProjectName(commandName) && len(args) > 0 {
				ShowSuccess(fmt.Sprintf("Completed %s %s", commandName, args[0]), "")
			} else {
				ShowSuccess(fmt.Sprintf("Completed %s", commandName), "")
			}
		}
		return
	}

	// If not an atempo command, try bash passthrough
	c.executeBashCommand(commandName, args)
}

// executeBashCommand executes a bash command with proper output handling
func (c *ShellCommand) executeBashCommand(commandName string, args []string) {
	// Handle special cd command (change directory)
	if commandName == "cd" {
		c.handleCdCommand(args)
		return
	}

	// Check if command exists in PATH
	_, err := exec.LookPath(commandName)
	if err != nil {
		ShowError(fmt.Sprintf("Command not found: %s", commandName), "Not an atempo command or system command")
		return
	}

	// Create the command
	cmd := exec.Command(commandName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Execute the command
	err = cmd.Run()
	
	// Only show status for commands that might change state
	// Don't show status for simple commands like ls, pwd, etc.
	stateChangingCommands := map[string]bool{
		"mkdir": true, "rmdir": true, "rm": true, "cp": true, "mv": true,
		"touch": true, "chmod": true, "chown": true, "ln": true,
		"git": true, "npm": true, "yarn": true, "composer": true,
		"docker": true, "docker-compose": true,
	}
	
	if stateChangingCommands[commandName] {
		if err != nil {
			ShowError(fmt.Sprintf("Command failed: %s", commandName), err.Error())
		} else {
			ShowSuccess(fmt.Sprintf("Completed: %s", commandName), "")
		}
	}
}

// handleCdCommand handles directory changes within the shell
func (c *ShellCommand) handleCdCommand(args []string) {
	var targetDir string
	
	if len(args) == 0 {
		// cd with no arguments goes to home directory
		home, err := os.UserHomeDir()
		if err != nil {
			ShowError("Failed to get home directory", err.Error())
			return
		}
		targetDir = home
	} else {
		targetDir = args[0]
		
		// Handle ~ expansion
		if strings.HasPrefix(targetDir, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				ShowError("Failed to get home directory", err.Error())
				return
			}
			targetDir = filepath.Join(home, targetDir[2:])
		}
	}
	
	// Change directory
	err := os.Chdir(targetDir)
	if err != nil {
		ShowError(fmt.Sprintf("cd: %s", targetDir), err.Error())
		return
	}
	
	// Show success for cd command
	pwd, _ := os.Getwd()
	ShowSuccess("Changed directory", pwd)
}

// createAutoCompleter creates an auto-completer for the shell
func (c *ShellCommand) createAutoCompleter() readline.AutoCompleter {
	// Get all available command names
	commands := c.registry.GetCommandNames()

	// Add built-in shell commands
	builtins := []string{"exit", "quit", "q", "clear", "cls", "help", "tips"}
	commands = append(commands, builtins...)
	
	// Add common bash commands for auto-completion
	bashCommands := []string{
		"ls", "cd", "pwd", "mkdir", "rmdir", "rm", "cp", "mv", "touch",
		"cat", "less", "head", "tail", "grep", "find", "which", "whereis",
		"git", "npm", "yarn", "composer", "docker", "docker-compose",
		"chmod", "chown", "ln", "ps", "kill", "top", "df", "du", "free",
	}
	commands = append(commands, bashCommands...)
	
	// Get project names for project-based completion
	projects := c.registry.GetProjectNames()
	
	// Create completion items
	var items []readline.PrefixCompleterInterface
	
	// Add global commands
	for _, cmd := range commands {
		items = append(items, readline.PcItem(cmd))
	}
	
	// Add project commands with sub-commands
	projectCommands := []string{"up", "down", "status", "logs", "describe", "shell", "reconfigure", "code", "cd", "open", "delete"}
	for _, project := range projects {
		// Create sub-completers for each project
		subItems := make([]readline.PrefixCompleterInterface, len(projectCommands))
		for i, subCmd := range projectCommands {
			subItems[i] = readline.PcItem(subCmd)
		}
		items = append(items, readline.PcItem(project, subItems...))
	}

	return readline.NewPrefixCompleter(items...)
}

// printDirectoryInfo prints directory and git info as a separate line
func (c *ShellCommand) printDirectoryInfo() {
	dir := getCurrentDirShort()
	branch := getGitBranch()
	status := getGitStatus()
	
	// Build the info line
	var info strings.Builder
	
	// Directory (cyan)
	info.WriteString(fmt.Sprintf("%s%s%s", ColorCyan, dir, ColorReset))
	
	if branch != "" {
		// Git branch (green for git:)
		info.WriteString(fmt.Sprintf(" %sgit:(%s%s%s)%s", ColorGreen, ColorYellow, branch, ColorGreen, ColorReset))
		
		// Git status
		if status == "✓" {
			info.WriteString(fmt.Sprintf(" %s%s%s", ColorGreen, status, ColorReset))
		} else if status == "✗" {
			info.WriteString(fmt.Sprintf(" %s%s%s", ColorRed, status, ColorReset))
		}
	}
	
	fmt.Printf("\n%s\n", info.String())
}

// updatePromptIfNeeded updates the cached prompt only when needed (like zsh)
func (c *ShellCommand) updatePromptIfNeeded() {
	currentDir := getCurrentDirShort()
	
	// Always check git info fresh after each command (like zsh)
	currentBranch := getGitBranch()
	currentStatus := getGitStatus()
	
	// Only regenerate prompt if something actually changed
	if currentDir != c.lastDir || currentBranch != c.lastGitBranch || currentStatus != c.lastGitStatus {
		c.lastDir = currentDir
		c.lastGitBranch = currentBranch
		c.lastGitStatus = currentStatus
		c.cachedPrompt = c.generatePrompt(currentDir, currentBranch, currentStatus)
	}
}

// generateSimplePrompt creates a simple prompt without git info for debugging
func (c *ShellCommand) generateSimplePrompt(dir string) string {
	return fmt.Sprintf("%s%s%s\n%s❯%s ", ColorCyan, dir, ColorReset, ColorPurple, ColorReset)
}

// generatePrompt creates a zsh-style prompt with provided directory and git info
func (c *ShellCommand) generatePrompt(dir, branch, status string) string {
	var prompt strings.Builder
	
	// Directory (cyan)
	prompt.WriteString(fmt.Sprintf("%s%s%s", ColorCyan, dir, ColorReset))
	
	if branch != "" {
		// Git branch (green for git:)
		prompt.WriteString(fmt.Sprintf(" %sgit:(%s%s%s)%s", ColorGreen, ColorYellow, branch, ColorGreen, ColorReset))
		
		// Git status
		if status == "✓" {
			prompt.WriteString(fmt.Sprintf(" %s%s%s", ColorGreen, status, ColorReset))
		} else if status == "✗" {
			prompt.WriteString(fmt.Sprintf(" %s%s%s", ColorRed, status, ColorReset))
		}
	}
	
	// Prompt character (purple)
	prompt.WriteString(fmt.Sprintf("\n%s❯%s ", ColorPurple, ColorReset))
	
	return prompt.String()
}

// Git detection and prompt generation functions

// getCurrentDirShort returns the current directory in a shortened format
func getCurrentDirShort() string {
	pwd, err := os.Getwd()
	if err != nil {
		return "~"
	}
	
	// Replace home directory with ~
	if home, err := os.UserHomeDir(); err == nil {
		if strings.HasPrefix(pwd, home) {
			pwd = "~" + pwd[len(home):]
		}
	}
	
	// Shorten long paths by showing only last 2 components
	parts := strings.Split(pwd, string(filepath.Separator))
	if len(parts) > 3 && !strings.HasPrefix(pwd, "~") {
		return ".../" + filepath.Join(parts[len(parts)-2:]...)
	}
	
	return pwd
}

// getGitBranch returns the current git branch if in a git repository
func getGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stderr = nil // Suppress errors
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getGitStatus returns git status indicators (clean/dirty)
func getGitStatus() string {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Stderr = nil // Suppress errors
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	status := strings.TrimSpace(string(output))
	if status == "" {
		return "✓" // Clean
	}
	return "✗" // Dirty
}


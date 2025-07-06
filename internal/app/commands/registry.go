package commands

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"atempo/internal/registry"
)

// CommandRegistry manages all available commands
type CommandRegistry struct {
	commands map[string]Command
	ctx      *CommandContext
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry(templatesFS, mcpServersFS embed.FS) *CommandRegistry {
	ctx := &CommandContext{}
	
	registry := &CommandRegistry{
		commands: make(map[string]Command),
		ctx:      ctx,
	}
	
	// Register all commands
	registry.register(NewCreateCommand(ctx, templatesFS, mcpServersFS))
	registry.register(NewAuthCommand(ctx))
	registry.register(NewDockerCommand(ctx))
	registry.register(NewProjectsCommand(ctx))
	registry.register(NewStatusCommand(ctx))
	registry.register(NewReconfigureCommand(ctx))
	registry.register(NewAddServiceCommand(ctx))
	registry.register(NewLogsCommand(ctx))
	registry.register(NewDescribeCommand(ctx))
	registry.register(NewRemoveCommand(ctx))
	registry.register(NewStopCommand(ctx))
	registry.register(NewTestCommand(ctx))
	registry.register(NewResetCommand(ctx))
	registry.register(NewShellCommand(ctx, registry))
	
	return registry
}

// register adds a command to the registry
func (r *CommandRegistry) register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Execute runs a command by name or routes project commands
func (r *CommandRegistry) Execute(ctx context.Context, commandName string, args []string) error {
	// First check if it's a registered global command
	if cmd, exists := r.commands[commandName]; exists {
		return cmd.Execute(ctx, args)
	}
	
	// Check if commandName is a project name
	if r.IsProjectName(commandName) {
		if len(args) == 0 {
			return fmt.Errorf("project command required. Usage: %s <command>", commandName)
		}
		
		// Route to project command handler
		projectCommand := args[0]
		projectArgs := args[1:]
		return r.executeProjectCommand(ctx, commandName, projectCommand, projectArgs)
	}
	
	return fmt.Errorf("unknown command: %s", commandName)
}

// GetCommand returns a command by name
func (r *CommandRegistry) GetCommand(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// ListCommands returns all available commands
func (r *CommandRegistry) ListCommands() []Command {
	commands := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// ShowUsage displays the main help message
func (r *CommandRegistry) ShowUsage() {
	fmt.Println(`
     █████╗ ████████╗███████╗███╗   ███╗██████╗  ██████╗ 
    ██╔══██╗╚══██╔══╝██╔════╝████╗ ████║██╔══██╗██╔═══██╗
    ███████║   ██║   █████╗  ██╔████╔██║██████╔╝██║   ██║
    ██╔══██║   ██║   ██╔══╝  ██║╚██╔╝██║██╔═══╝ ██║   ██║
    ██║  ██║   ██║   ███████╗██║ ╚═╝ ██║██║     ╚██████╔╝
    ╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝     ╚═╝╚═╝      ╚═════╝ 

Usage:
  atempo <command> [arguments]

Commands:`)

	// Display commands in a logical order
	commandOrder := []string{
		"create", "auth", "status", "describe", "docker", 
		"reconfigure", "add-service", "projects", "remove", "logs", "stop", "test", "reset",
	}
	
	for _, cmdName := range commandOrder {
		if cmd, exists := r.commands[cmdName]; exists {
			fmt.Printf("  %-20s %s\n", cmdName, cmd.Description())
		}
	}

	fmt.Println(`
Examples:
  atempo create laravel my-app          Create Laravel (latest) in ./my-app/
  atempo create laravel:11 my-app       Create Laravel 11 in ./my-app/
  atempo create django                  Create Django (latest) in current directory
  atempo create django:5                Create Django 5 in current directory
  atempo status                         Show dashboard with all project statuses
  atempo describe my-app                Show detailed description of 'my-app' project
  atempo describe                       Describe project in current directory
  atempo docker up                      Start services in current directory
  atempo docker up my-app               Start services for registered project 'my-app'
  atempo reconfigure                    Regenerate docker-compose.yml from atempo.json
  atempo add-service minio              Add MinIO object storage service
  atempo projects                       List all registered projects
  atempo logs my-app                    View setup logs for 'my-app' project
  atempo stop                           Stop all running projects
  atempo test                           Run all tests in current project
  atempo test my-app                    Run all tests for 'my-app' project
  atempo test my-app UserTest           Run specific test suite for Laravel
  atempo test accounts                  Run tests for Django 'accounts' app

Project Management:
  - Projects are automatically registered when created with 'atempo create'
  - Use project names instead of paths: 'atempo docker up my-laravel-app'
  - Services defined in atempo.json generate docker-compose.yml automatically

For more information about specific commands:
  atempo <command> --help`)
}

// HasCommand checks if a command exists
func (r *CommandRegistry) HasCommand(name string) bool {
	_, exists := r.commands[name]
	return exists
}

// IsHelpCommand checks if the argument is a help request
func IsHelpCommand(arg string) bool {
	helpCommands := []string{"help", "--help", "-h"}
	for _, helpCmd := range helpCommands {
		if strings.EqualFold(arg, helpCmd) {
			return true
		}
	}
	return false
}

// GetCommandNames returns a slice of all registered command names
func (r *CommandRegistry) GetCommandNames() []string {
	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	return names
}

// GetProjectNames returns a slice of all registered project names
func (r *CommandRegistry) GetProjectNames() []string {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return []string{}
	}
	
	projects := reg.ListProjects()
	names := make([]string, len(projects))
	for i, project := range projects {
		names[i] = project.Name
	}
	return names
}

// IsProjectName checks if a name matches a registered project
func (r *CommandRegistry) IsProjectName(name string) bool {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return false
	}
	
	_, err = reg.FindProject(name)
	return err == nil
}

// executeProjectCommand handles project-specific commands
func (r *CommandRegistry) executeProjectCommand(ctx context.Context, projectName, command string, args []string) error {
	// Map project commands to existing global commands
	switch command {
	case "up", "start":
		// Execute docker up for this project
		dockerCmd := r.commands["docker"]
		return dockerCmd.Execute(ctx, append([]string{"up", projectName}, args...))
	
	case "down", "stop":
		// Execute docker down for this project  
		dockerCmd := r.commands["docker"]
		return dockerCmd.Execute(ctx, append([]string{"down", projectName}, args...))
	
	case "status":
		// Execute status for this project
		statusCmd := r.commands["status"]
		return statusCmd.Execute(ctx, append([]string{projectName}, args...))
	
	case "logs":
		// Execute logs for this project
		logsCmd := r.commands["logs"]
		return logsCmd.Execute(ctx, append([]string{projectName}, args...))
	
	case "describe", "info":
		// Execute describe for this project
		describeCmd := r.commands["describe"]
		return describeCmd.Execute(ctx, append([]string{projectName}, args...))
	
	case "shell", "bash", "exec":
		// Execute shell access for this project
		dockerCmd := r.commands["docker"]
		return dockerCmd.Execute(ctx, append([]string{"bash", projectName}, args...))
	
	case "reconfigure", "reconfig":
		// Execute reconfigure for this project
		reconfigCmd := r.commands["reconfigure"]
		return reconfigCmd.Execute(ctx, append([]string{projectName}, args...))
	
	case "code":
		// Open project in VS Code
		return r.openProjectInVSCode(projectName)
	
	case "cd":
		// Change directory to project (note: this only works in shell session)
		return r.changeToProjectDirectory(projectName)
	
	case "delete", "remove":
		// Delete project files and remove from registry
		return r.deleteProject(projectName)
	
	case "open":
		// Open project or specific service in browser
		return r.openProjectInBrowser(projectName, args)
	
	default:
		return fmt.Errorf("unknown project command: %s. Available: up, down, status, logs, describe, shell, reconfigure, code, cd, delete, open", command)
	}
}

// openProjectInVSCode opens the specified project in VS Code
func (r *CommandRegistry) openProjectInVSCode(projectName string) error {
	// Load registry to get project details
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}
	
	// Find the project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project not found: %s", projectName)
	}
	
	// Check if VS Code is installed
	codePath, err := exec.LookPath("code")
	if err != nil {
		return fmt.Errorf("VS Code (code command) not found. Please install VS Code and ensure 'code' is in your PATH")
	}
	
	// Open the project directory in VS Code
	cmd := exec.Command(codePath, project.Path)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open project in VS Code: %w", err)
	}
	
	// Don't wait for VS Code to close - let it run in background
	go func() {
		cmd.Wait()
	}()
	
	return nil
}

// changeToProjectDirectory changes the current working directory to the project path
func (r *CommandRegistry) changeToProjectDirectory(projectName string) error {
	// Load registry to get project details
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}
	
	// Find the specified project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project '%s' not found in registry", projectName)
	}
	
	// Change to the project directory
	if err := os.Chdir(project.Path); err != nil {
		return fmt.Errorf("failed to change to project directory %s: %w", project.Path, err)
	}
	
	fmt.Printf("  ⎿ Changed to project directory: %s\n", project.Path)
	return nil
}

// deleteProject deletes project files and removes from registry
func (r *CommandRegistry) deleteProject(projectName string) error {
	// Load registry to get project details
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}
	
	// Find the specified project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project '%s' not found in registry", projectName)
	}
	
	// Show confirmation prompt with detailed info
	fmt.Printf("! Are you sure you want to delete project '%s'?\n\n", projectName)
	fmt.Printf("  Path: %s\n", project.Path)
	fmt.Printf("  Framework: %s %s\n", project.Framework, project.Version)
	fmt.Printf("  Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("\n  This will:\n")
	fmt.Printf("  • Move project files to Trash\n")
	fmt.Printf("  • Remove project from atempo registry\n")
	fmt.Printf("  • This action cannot be undone!\n\n")
	fmt.Print("Type 'delete' to confirm, or anything else to cancel: ")

	var response string
	fmt.Scanln(&response)
	
	if strings.ToLower(response) != "delete" {
		fmt.Println("✗ Cancelled - project not deleted.")
		return nil
	}

	// Move to trash using macOS 'trash' command or fallback
	err = moveToTrash(project.Path)
	if err != nil {
		return fmt.Errorf("failed to move project to trash: %w", err)
	}

	// Remove project from registry
	err = reg.RemoveProject(projectName)
	if err != nil {
		return fmt.Errorf("failed to remove project from registry: %w", err)
	}

	fmt.Printf("✓ Project '%s' deleted successfully!\n", projectName)
	fmt.Printf("  ⎿ Files moved to Trash\n")
	fmt.Printf("  ⎿ Removed from registry\n")
	
	return nil
}

// moveToTrash moves a directory to the trash
func moveToTrash(path string) error {
	// Try using 'trash' command if available (brew install trash)
	if _, err := exec.LookPath("trash"); err == nil {
		cmd := exec.Command("trash", path)
		return cmd.Run()
	}
	
	// Fallback: use 'mv' to move to ~/.Trash (macOS)
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	trashDir := filepath.Join(home, ".Trash")
	
	// Ensure .Trash directory exists
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("failed to create .Trash directory: %w", err)
	}
	
	basename := filepath.Base(path)
	
	// Create unique name if file already exists in trash
	trashPath := filepath.Join(trashDir, basename)
	counter := 1
	for {
		if _, err := os.Stat(trashPath); os.IsNotExist(err) {
			break
		}
		trashPath = filepath.Join(trashDir, fmt.Sprintf("%s (%d)", basename, counter))
		counter++
	}
	
	return os.Rename(path, trashPath)
}

// openProjectInBrowser opens the project or specific service in the default browser
func (r *CommandRegistry) openProjectInBrowser(projectName string, args []string) error {
	// Load registry to get project details
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}
	
	// Find the project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project not found: %s", projectName)
	}
	
	// Update project status to get current URLs and services
	err = reg.UpdateProjectStatus(projectName)
	if err != nil {
		return fmt.Errorf("failed to update project status: %w", err)
	}
	
	// Reload to get updated project info
	project, err = reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("failed to reload project: %w", err)
	}
	
	// Check if project has running services
	if project.Status == "stopped" || project.Status == "no-docker" || project.Status == "no-services" {
		return fmt.Errorf("project '%s' is not running. Start it with: %s up", projectName, projectName)
	}
	
	var targetURL string
	
	if len(args) == 0 {
		// Open main application (first available URL)
		if len(project.URLs) == 0 {
			return fmt.Errorf("no web URLs found for project '%s'. Make sure services are running and have exposed web ports", projectName)
		}
		targetURL = project.URLs[0]
		ShowInfo(fmt.Sprintf("Opening main application: %s", targetURL))
	} else {
		// Open specific service
		serviceName := args[0]
		found := false
		
		for _, service := range project.Services {
			if service.Name == serviceName {
				if service.URL == "" {
					return fmt.Errorf("service '%s' doesn't have a web URL (no exposed web ports)", serviceName)
				}
				targetURL = service.URL
				found = true
				break
			}
		}
		
		if !found {
			// List available services for user
			var availableServices []string
			for _, service := range project.Services {
				if service.URL != "" {
					availableServices = append(availableServices, service.Name)
				}
			}
			
			if len(availableServices) == 0 {
				return fmt.Errorf("no services with web URLs found for project '%s'", projectName)
			}
			
			return fmt.Errorf("service '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", "))
		}
		
		ShowInfo(fmt.Sprintf("Opening service '%s': %s", serviceName, targetURL))
	}
	
	// Open URL in default browser
	return openURL(targetURL)
}

// openURL opens a URL in the default browser (cross-platform)
func openURL(url string) error {
	var cmd *exec.Cmd
	
	// Try platform-specific commands
	if err := exec.Command("which", "open").Run(); err == nil {
		// macOS
		cmd = exec.Command("open", url)
	} else if err := exec.Command("which", "xdg-open").Run(); err == nil {
		// Linux
		cmd = exec.Command("xdg-open", url)
	} else if err := exec.Command("where", "cmd").Run(); err == nil {
		// Windows
		cmd = exec.Command("cmd", "/c", "start", url)
	} else {
		return fmt.Errorf("unable to find browser command (tried: open, xdg-open, cmd)")
	}
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open URL in browser: %w", err)
	}
	
	// Don't wait for browser to close
	go func() {
		cmd.Wait()
	}()
	
	ShowSuccess("Browser opened", url)
	return nil
}
package commands

import (
	"context"
	"embed"
	"fmt"
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
	registry.register(NewAuthCommand(ctx)) // Deprecated - shows migration notice
	registry.register(NewAICommand())
	registry.register(NewDockerCommand(ctx))
	registry.register(NewProjectsCommand(ctx))
	registry.register(NewReconfigureCommand(ctx))
	registry.register(NewAddServiceCommand(ctx))
	registry.register(NewLogsCommand(ctx))
	registry.register(NewDescribeCommand(ctx))
	registry.register(NewRemoveCommand(ctx))
	registry.register(NewStopCommand(ctx))
	registry.register(NewTestCommand(ctx))
	registry.register(NewResetCommand(ctx))
	registry.register(NewDNSCommand(ctx))
	registry.register(NewSSLCommand(ctx))
	registry.register(NewAuditCommand(ctx))
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
		"create", "ai", "projects", "describe", "docker", "dns", "ssl",
		"reconfigure", "add-service", "remove", "logs", "stop", "test", "audit", "reset",
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
  atempo projects                       Show all projects with their status
  atempo projects my-app                Show detailed status for 'my-app' project
  atempo describe my-app                Show detailed description of 'my-app' project
  atempo describe                       Describe project in current directory
  atempo docker up                      Start services in current directory
  atempo docker up my-app               Start services for registered project 'my-app'
  atempo reconfigure                    Regenerate docker-compose.yml from atempo.json
  atempo add-service minio              Add MinIO object storage service
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

package commands

import (
	"context"
	"embed"
	"fmt"
	"strings"
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
	
	return registry
}

// register adds a command to the registry
func (r *CommandRegistry) register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Execute runs a command by name
func (r *CommandRegistry) Execute(ctx context.Context, commandName string, args []string) error {
	cmd, exists := r.commands[commandName]
	if !exists {
		return fmt.Errorf("unknown command: %s", commandName)
	}
	
	return cmd.Execute(ctx, args)
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
		"reconfigure", "add-service", "projects", "logs",
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
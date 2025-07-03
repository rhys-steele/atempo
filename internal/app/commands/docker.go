package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"atempo/internal/docker"
	"atempo/internal/registry"
)

// DockerCommand handles Docker-related subcommands
type DockerCommand struct {
	*BaseCommand
}

// NewDockerCommand creates a new docker command
func NewDockerCommand(ctx *CommandContext) *DockerCommand {
	return &DockerCommand{
		BaseCommand: NewBaseCommand(
			"docker",
			"Run Docker operations on projects",
			"atempo docker <command> [project] [options]",
			ctx,
		),
	}
}

// Execute runs the docker command
func (c *DockerCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s\n\n%s", c.Usage(), c.getDockerUsage())
	}

	// Validate Docker installation
	if err := docker.ValidateDockerCompose(); err != nil {
		return fmt.Errorf("docker validation failed: %w", err)
	}

	dockerCmd := args[0]
	var projectPath string
	var additionalArgs []string

	// Parse arguments: atempo docker <command> [project_name_or_path] [additional_args...]
	if len(args) > 1 {
		// Check if the second argument is a path/name or additional docker args
		potentialIdentifier := args[1]
		if strings.HasPrefix(potentialIdentifier, "-") || c.isDockerArg(potentialIdentifier) {
			// It's a docker argument, use current directory
			projectPath = ""
			additionalArgs = args[1:]
		} else {
			// It's a project identifier (name or path)
			resolvedPath, err := registry.ResolveProjectPath(potentialIdentifier)
			if err != nil {
				return fmt.Errorf("failed to resolve project: %w", err)
			}
			projectPath = resolvedPath
			if len(args) > 2 {
				additionalArgs = args[2:]
			}
		}
	}

	// Check for timeout flag in additional args
	timeout, filteredArgs := c.parseTimeoutFlag(additionalArgs)
	
	// Handle special commands
	switch dockerCmd {
	case "exec":
		return c.handleDockerExec(projectPath, filteredArgs)
	case "services":
		return c.handleDockerServices(projectPath)
	default:
		// Standard docker-compose command with optional custom timeout
		if timeout > 0 {
			return docker.ExecuteWithCustomTimeout(dockerCmd, projectPath, filteredArgs, timeout)
		}
		return docker.ExecuteCommand(dockerCmd, projectPath, filteredArgs)
	}
}

// handleDockerExec processes docker exec commands
func (c *DockerCommand) handleDockerExec(projectPath string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: atempo docker exec <service> [command...]\nExample: atempo docker exec app bash")
	}

	service := args[0]
	cmdArgs := []string{"bash"} // default to bash
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	return docker.ExecuteExecCommand(service, projectPath, cmdArgs)
}

// handleDockerServices lists available services
func (c *DockerCommand) handleDockerServices(projectPath string) error {
	return docker.ListServices(projectPath)
}

// isDockerArg checks if a string looks like a Docker argument
func (c *DockerCommand) isDockerArg(arg string) bool {
	dockerArgs := []string{"--force-recreate", "--build", "--no-deps", "--remove-orphans", "-V", "--volumes"}
	for _, dockerArg := range dockerArgs {
		if arg == dockerArg {
			return true
		}
	}
	return false
}

// getDockerUsage returns detailed Docker usage information
func (c *DockerCommand) getDockerUsage() string {
	return `Atempo Docker Commands

Usage:
  atempo docker <command> [project_name_or_path] [options]

Common Commands:
  up [project]           Start services in detached mode
  down [project]         Stop and remove containers  
  build [project]        Build or rebuild services
  logs [project] [svc]   View output from containers
  ps [project]           List containers
  restart [project]      Restart services
  stop [project]         Stop running containers
  exec <service> [cmd]   Execute command in container
  services [project]     List available services

Examples:
  atempo docker up                    # Start services in current directory
  atempo docker up my-laravel-app    # Start services for registered project
  atempo docker up ../myproject      # Start services in relative path
  atempo docker logs app             # View app container logs
  atempo docker exec app bash        # Open bash in app container
  atempo docker exec web python manage.py shell  # Django shell
  atempo docker down --volumes       # Stop and remove volumes

Project Resolution:
  - Project name (from registry): 'my-laravel-app'
  - Relative path: '../myproject'  
  - Absolute path: '/full/path/to/project'
  - Current directory if no argument provided`
}

// parseTimeoutFlag extracts timeout flag from arguments and returns filtered args
func (c *DockerCommand) parseTimeoutFlag(args []string) (time.Duration, []string) {
	var filteredArgs []string
	var timeout time.Duration
	
	for i := 0; i < len(args); i++ {
		arg := args[i]
		
		// Check for --timeout flag
		if arg == "--timeout" && i+1 < len(args) {
			if duration, err := c.parseTimeoutValue(args[i+1]); err == nil {
				timeout = duration
				i++ // Skip the next argument (timeout value)
				continue
			}
		}
		
		// Check for --timeout=value format
		if strings.HasPrefix(arg, "--timeout=") {
			value := strings.TrimPrefix(arg, "--timeout=")
			if duration, err := c.parseTimeoutValue(value); err == nil {
				timeout = duration
				continue
			}
		}
		
		// Keep all other arguments
		filteredArgs = append(filteredArgs, arg)
	}
	
	return timeout, filteredArgs
}

// parseTimeoutValue parses timeout string into duration (supports suffixes like 5m, 30s, etc.)
func (c *DockerCommand) parseTimeoutValue(value string) (time.Duration, error) {
	// Try parsing as duration first (5m, 30s, etc.)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration, nil
	}
	
	// Try parsing as plain number (assume minutes)
	if minutes, err := strconv.Atoi(value); err == nil {
		return time.Duration(minutes) * time.Minute, nil
	}
	
	return 0, fmt.Errorf("invalid timeout format: %s", value)
}
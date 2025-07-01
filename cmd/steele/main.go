package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"steele/internal/compose"
	"steele/internal/docker"
	"steele/internal/registry"
	"steele/internal/scaffold"
)

//go:embed templates
var templatesFS embed.FS

//go:embed mcp-servers/*
var mcpServersFS embed.FS

// main is the entry point for the Steele CLI.
// Steele is a command-line tool for scaffolding and managing developer projects
// using AI-first principles and an MCP-ready context architecture.
func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "start":
		handleStartCommand()
	case "docker":
		handleDockerCommand()
	case "reconfigure":
		handleReconfigureCommand()
	case "projects":
		handleProjectsCommand()
	case "add-service":
		handleAddServiceCommand()
	case "help", "--help", "-h":
		showUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		showUsage()
		os.Exit(1)
	}
}

// handleStartCommand processes the start subcommand for scaffolding
func handleStartCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: steele start <framework>:<version>")
		os.Exit(1)
	}

	// Extract framework and version
	arg := os.Args[2]
	parts := strings.Split(arg, ":")
	if len(parts) != 2 {
		fmt.Println("Error: expected format <framework>:<version>")
		os.Exit(1)
	}

	framework := parts[0]
	version := parts[1]

	// Trigger the scaffold process
	err := scaffold.Run(framework, version, templatesFS, mcpServersFS)
	if err != nil {
		fmt.Printf("Scaffold error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Project scaffolding complete.")
}

// handleDockerCommand processes Docker-related subcommands
func handleDockerCommand() {
	if len(os.Args) < 3 {
		showDockerUsage()
		os.Exit(1)
	}

	// Validate Docker installation
	if err := docker.ValidateDockerCompose(); err != nil {
		fmt.Printf("Docker validation failed: %v\n", err)
		os.Exit(1)
	}

	dockerCmd := os.Args[2]
	var projectPath string
	var additionalArgs []string

	// Parse arguments: steele docker <command> [project_name_or_path] [additional_args...]
	if len(os.Args) > 3 {
		// Check if the third argument is a path/name or additional docker args
		potentialIdentifier := os.Args[3]
		if strings.HasPrefix(potentialIdentifier, "-") || isDockerArg(potentialIdentifier) {
			// It's a docker argument, use current directory
			projectPath = ""
			additionalArgs = os.Args[3:]
		} else {
			// It's a project identifier (name or path)
			resolvedPath, err := registry.ResolveProjectPath(potentialIdentifier)
			if err != nil {
				fmt.Printf("Failed to resolve project: %v\n", err)
				os.Exit(1)
			}
			projectPath = resolvedPath
			if len(os.Args) > 4 {
				additionalArgs = os.Args[4:]
			}
		}
	}

	// Handle special commands
	switch dockerCmd {
	case "exec":
		handleDockerExec(projectPath, additionalArgs)
	case "services":
		handleDockerServices(projectPath)
	default:
		// Standard docker-compose command
		err := docker.ExecuteCommand(dockerCmd, projectPath, additionalArgs)
		if err != nil {
			fmt.Printf("Docker command failed: %v\n", err)
			os.Exit(1)
		}
	}
}

// handleDockerExec processes docker exec commands
func handleDockerExec(projectPath string, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: steele docker exec <service> [command...]")
		fmt.Println("Example: steele docker exec app bash")
		os.Exit(1)
	}

	service := args[0]
	cmdArgs := []string{"bash"} // default to bash
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	err := docker.ExecuteExecCommand(service, projectPath, cmdArgs)
	if err != nil {
		fmt.Printf("Docker exec failed: %v\n", err)
		os.Exit(1)
	}
}

// handleDockerServices lists available services
func handleDockerServices(projectPath string) {
	err := docker.ListServices(projectPath)
	if err != nil {
		fmt.Printf("Failed to list services: %v\n", err)
		os.Exit(1)
	}
}

// isDockerArg checks if a string looks like a Docker argument
func isDockerArg(arg string) bool {
	dockerArgs := []string{"--force-recreate", "--build", "--no-deps", "--remove-orphans", "-V", "--volumes"}
	for _, dockerArg := range dockerArgs {
		if arg == dockerArg {
			return true
		}
	}
	return false
}

// showUsage displays the main help message
func showUsage() {
	fmt.Println(`Steele - AI-first project scaffolding and management

Usage:
  steele <command> [arguments]

Commands:
  start <framework>:<version>   Create a new project with the specified framework
  docker <command> [project]    Run Docker operations on projects
  reconfigure [project]         Regenerate docker-compose.yml from steele.json
  add-service <type> [project]  Add predefined services (minio, elasticsearch, etc.)
  projects                      List all registered projects
  help                          Show this help message

Examples:
  steele start laravel:11       Create a new Laravel 11 project
  steele start django:5         Create a new Django 5 project
  steele docker up              Start services in current directory
  steele docker up my-app       Start services for registered project 'my-app'
  steele reconfigure            Regenerate docker-compose.yml from steele.json
  steele add-service minio      Add MinIO object storage service
  steele projects               List all registered projects

Project Management:
  - Projects are automatically registered when created with 'steele start'
  - Use project names instead of paths: 'steele docker up my-laravel-app'
  - Services defined in steele.json generate docker-compose.yml automatically

For more information about specific commands:
  steele docker --help`)
}

// handleReconfigureCommand regenerates docker-compose.yml from steele.json
func handleReconfigureCommand() {
	var projectPath string
	
	if len(os.Args) > 2 {
		resolvedPath, err := registry.ResolveProjectPath(os.Args[2])
		if err != nil {
			fmt.Printf("Failed to resolve project: %v\n", err)
			os.Exit(1)
		}
		projectPath = resolvedPath
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get current directory: %v\n", err)
			os.Exit(1)
		}
		projectPath = cwd
	}

	fmt.Printf("→ Regenerating docker-compose.yml from steele.json in %s...\n", projectPath)
	
	if err := compose.GenerateDockerCompose(projectPath); err != nil {
		fmt.Printf("Failed to regenerate docker-compose.yml: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ docker-compose.yml regenerated successfully!")
}

// handleProjectsCommand lists all registered projects
func handleProjectsCommand() {
	reg, err := registry.LoadRegistry()
	if err != nil {
		fmt.Printf("Failed to load registry: %v\n", err)
		os.Exit(1)
	}

	projects := reg.ListProjects()
	if len(projects) == 0 {
		fmt.Println("No projects registered yet.")
		fmt.Println("Projects are automatically registered when you run 'steele start'")
		return
	}

	fmt.Println("Registered Steele Projects:")
	fmt.Println()
	
	for _, project := range projects {
		fmt.Printf("  %s\n", project.Name)
		fmt.Printf("    Framework: %s %s\n", project.Framework, project.Version)
		fmt.Printf("    Path: %s\n", project.Path)
		fmt.Printf("    Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Println()
	}
}

// handleAddServiceCommand adds a predefined service to a project
func handleAddServiceCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: steele add-service <service_type> [project]")
		fmt.Println("\nAvailable services:")
		for _, service := range compose.ListPredefinedServices() {
			fmt.Printf("  %s\n", service)
		}
		os.Exit(1)
	}

	serviceType := os.Args[2]
	var projectPath string

	if len(os.Args) > 3 {
		resolvedPath, err := registry.ResolveProjectPath(os.Args[3])
		if err != nil {
			fmt.Printf("Failed to resolve project: %v\n", err)
			os.Exit(1)
		}
		projectPath = resolvedPath
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get current directory: %v\n", err)
			os.Exit(1)
		}
		projectPath = cwd
	}

	fmt.Printf("→ Adding %s service to project...\n", serviceType)
	
	if err := compose.AddPredefinedService(projectPath, serviceType); err != nil {
		fmt.Printf("Failed to add service: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ %s service added to steele.json\n", serviceType)
	fmt.Println("Run 'steele reconfigure' to update docker-compose.yml")
}

// showDockerUsage displays Docker-specific help
func showDockerUsage() {
	fmt.Println(`Steele Docker Commands

Usage:
  steele docker <command> [project_name_or_path] [options]

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
  steele docker up                    # Start services in current directory
  steele docker up my-laravel-app    # Start services for registered project
  steele docker up ../myproject      # Start services in relative path
  steele docker logs app             # View app container logs
  steele docker exec app bash        # Open bash in app container
  steele docker exec web python manage.py shell  # Django shell
  steele docker down --volumes       # Stop and remove volumes

Project Resolution:
  - Project name (from registry): 'my-laravel-app'
  - Relative path: '../myproject'  
  - Absolute path: '/full/path/to/project'
  - Current directory if no argument provided`)
}

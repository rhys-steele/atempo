package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"steele/internal/utils"
)

// DockerCommand represents available Docker operations
type DockerCommand struct {
	Name        string
	Description string
	Args        []string
}

// Common Docker commands for Steele projects
var SupportedCommands = map[string]DockerCommand{
	"up": {
		Name:        "up",
		Description: "Start services in detached mode",
		Args:        []string{"up", "-d"},
	},
	"down": {
		Name:        "down",
		Description: "Stop and remove containers",
		Args:        []string{"down"},
	},
	"build": {
		Name:        "build",
		Description: "Build or rebuild services",
		Args:        []string{"build"},
	},
	"logs": {
		Name:        "logs",
		Description: "View output from containers",
		Args:        []string{"logs", "-f"},
	},
	"ps": {
		Name:        "ps",
		Description: "List containers",
		Args:        []string{"ps"},
	},
	"restart": {
		Name:        "restart",
		Description: "Restart services",
		Args:        []string{"restart"},
	},
	"stop": {
		Name:        "stop",
		Description: "Stop running containers",
		Args:        []string{"stop"},
	},
	"pull": {
		Name:        "pull",
		Description: "Pull service images",
		Args:        []string{"pull"},
	},
}

// ExecuteCommand runs a Docker Compose command in the specified project directory
func ExecuteCommand(command string, projectPath string, additionalArgs []string) error {
	// Resolve project path
	resolvedPath, err := resolveProjectPath(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Validate that docker-compose.yml exists
	composePath := filepath.Join(resolvedPath, "docker-compose.yml")
	if !utils.FileExists(composePath) {
		return fmt.Errorf("docker-compose.yml not found in %s", resolvedPath)
	}

	// Get the Docker command configuration
	dockerCmd, exists := SupportedCommands[command]
	if !exists {
		return fmt.Errorf("unsupported Docker command: %s", command)
	}

	// Build the full command
	args := append(dockerCmd.Args, additionalArgs...)
	fullCommand := append([]string{"docker-compose"}, args...)

	fmt.Printf("→ Running: %s (in %s)\n", strings.Join(fullCommand, " "), resolvedPath)

	// Execute the command
	cmd := exec.Command(fullCommand[0], fullCommand[1:]...)
	cmd.Dir = resolvedPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ExecuteExecCommand runs a command inside a container (docker-compose exec)
func ExecuteExecCommand(service string, projectPath string, cmdArgs []string) error {
	// Resolve project path
	resolvedPath, err := resolveProjectPath(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Validate that docker-compose.yml exists
	composePath := filepath.Join(resolvedPath, "docker-compose.yml")
	if !utils.FileExists(composePath) {
		return fmt.Errorf("docker-compose.yml not found in %s", resolvedPath)
	}

	// Build the exec command
	args := append([]string{"docker-compose", "exec", service}, cmdArgs...)

	fmt.Printf("→ Running: %s (in %s)\n", strings.Join(args, " "), resolvedPath)

	// Execute the command
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = resolvedPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ListServices shows available services in the docker-compose.yml
func ListServices(projectPath string) error {
	// Resolve project path
	resolvedPath, err := resolveProjectPath(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Validate that docker-compose.yml exists
	composePath := filepath.Join(resolvedPath, "docker-compose.yml")
	if !utils.FileExists(composePath) {
		return fmt.Errorf("docker-compose.yml not found in %s", resolvedPath)
	}

	fmt.Printf("→ Services in %s:\n", resolvedPath)

	// Run docker-compose config --services
	cmd := exec.Command("docker-compose", "config", "--services")
	cmd.Dir = resolvedPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// resolveProjectPath determines the correct project directory
func resolveProjectPath(projectPath string) (string, error) {
	var targetPath string

	if projectPath == "" {
		// Use current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		targetPath = cwd
	} else {
		// Use provided path
		if filepath.IsAbs(projectPath) {
			targetPath = projectPath
		} else {
			// Convert relative path to absolute
			cwd, err := os.Getwd()
			if err != nil {
				return "", fmt.Errorf("failed to get current directory: %w", err)
			}
			targetPath = filepath.Join(cwd, projectPath)
		}
	}

	// Verify the directory exists
	if !utils.FileExists(targetPath) {
		return "", fmt.Errorf("project directory does not exist: %s", targetPath)
	}

	return targetPath, nil
}

// ValidateDockerCompose checks if Docker and Docker Compose are available
func ValidateDockerCompose() error {
	// Check if docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker command not found. Please install Docker")
	}

	// Check if docker-compose is available
	if _, err := exec.LookPath("docker-compose"); err != nil {
		return fmt.Errorf("docker-compose command not found. Please install Docker Compose")
	}

	return nil
}

// DetectFramework attempts to detect the framework based on project files
func DetectFramework(projectPath string) (string, error) {
	resolvedPath, err := resolveProjectPath(projectPath)
	if err != nil {
		return "", err
	}

	// Check for Laravel indicators
	if utils.FileExists(filepath.Join(resolvedPath, "src", "artisan")) ||
		utils.FileExists(filepath.Join(resolvedPath, "src", "composer.json")) {
		return "laravel", nil
	}

	// Check for Django indicators
	if utils.FileExists(filepath.Join(resolvedPath, "src", "manage.py")) ||
		utils.FileExists(filepath.Join(resolvedPath, "src", "requirements.txt")) {
		return "django", nil
	}

	return "unknown", nil
}

// GetFrameworkServices returns common services for different frameworks
func GetFrameworkServices(framework string) []string {
	switch framework {
	case "laravel":
		return []string{"app", "webserver", "mysql", "redis"}
	case "django":
		return []string{"web", "postgres", "redis", "worker", "beat"}
	default:
		return []string{}
	}
}
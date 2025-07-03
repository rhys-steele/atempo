package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"atempo/internal/utils"
)

// Bake detection cache
var (
	bakeSupported   *bool
	bakeMutex       sync.Mutex
)

// DockerCommand represents available Docker operations
type DockerCommand struct {
	Name        string
	Description string
	Args        []string
	Timeout     time.Duration // Default timeout for this command
}

// Common Docker commands for Atempo projects
var SupportedCommands = map[string]DockerCommand{
	"up": {
		Name:        "up",
		Description: "Build and start services in detached mode",
		Args:        []string{"up", "-d", "--build"},
		Timeout:     10 * time.Minute, // Long timeout for building + pulling images
	},
	"down": {
		Name:        "down",
		Description: "Stop and remove containers",
		Args:        []string{"down"},
		Timeout:     2 * time.Minute, // Shorter timeout for stopping
	},
	"build": {
		Name:        "build",
		Description: "Build or rebuild services",
		Timeout:     8 * time.Minute, // Long timeout for building
		Args:        []string{"build"},
	},
	"logs": {
		Name:        "logs",
		Description: "View output from containers",
		Args:        []string{"logs", "-f"},
		Timeout:     0, // No timeout for logs (streaming)
	},
	"ps": {
		Name:        "ps",
		Description: "List containers",
		Args:        []string{"ps"},
		Timeout:     30 * time.Second, // Quick command
	},
	"restart": {
		Name:        "restart",
		Description: "Restart services",
		Args:        []string{"restart"},
		Timeout:     3 * time.Minute, // Medium timeout
	},
	"stop": {
		Name:        "stop",
		Timeout:     2 * time.Minute, // Medium timeout
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
	// Get the Docker command configuration
	dockerCmd, exists := SupportedCommands[command]
	if !exists {
		return fmt.Errorf("unsupported Docker command: %s", command)
	}
	
	// Use the shared execution logic
	return executeWithCommand(dockerCmd, projectPath, additionalArgs)
}

// ExecuteWithCustomTimeout allows overriding the default timeout for a command
func ExecuteWithCustomTimeout(command string, projectPath string, additionalArgs []string, customTimeout time.Duration) error {
	// Get the command configuration
	dockerCmd, exists := SupportedCommands[command]
	if !exists {
		return fmt.Errorf("unsupported Docker command: %s", command)
	}
	
	// Override the timeout
	dockerCmd.Timeout = customTimeout
	
	// Use the same logic but with custom timeout
	return executeWithCommand(dockerCmd, projectPath, additionalArgs)
}

// executeWithCommand is the core execution logic extracted for reuse
func executeWithCommand(dockerCmd DockerCommand, projectPath string, additionalArgs []string) error {
	// Resolve project path
	resolvedPath, err := resolveProjectPath(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Look for docker-compose.yml in project root first (new architecture)
	rootComposePath := filepath.Join(resolvedPath, "docker-compose.yml")
	dockerDir := resolvedPath
	var composeFile string
	
	if utils.FileExists(rootComposePath) {
		// Use compose file in project root
		composeFile = "docker-compose.yml"
	} else {
		// Fallback: check infra/docker subdirectory for legacy projects
		legacyDockerDir := filepath.Join(resolvedPath, "infra", "docker")
		legacyComposePath := filepath.Join(legacyDockerDir, "docker-compose.yml")
		if !utils.FileExists(legacyComposePath) {
			return fmt.Errorf("docker-compose.yml not found in %s or %s", resolvedPath, legacyDockerDir)
		}
		// Use compose file in subdirectory with -f flag, but run from project root
		composeFile = "infra/docker/docker-compose.yml"
	}

	// Build the full command with -f flag for compose file location
	baseArgs := []string{"-f", composeFile}
	args := append(baseArgs, dockerCmd.Args...)
	args = append(args, additionalArgs...)
	fullCommand := append([]string{"docker-compose"}, args...)

	// Create context with timeout
	var ctx context.Context
	var cancel context.CancelFunc
	
	if dockerCmd.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), dockerCmd.Timeout)
		defer cancel()
		fmt.Printf("→ Running: %s (in %s, timeout: %v)\n", strings.Join(fullCommand, " "), dockerDir, dockerCmd.Timeout)
	} else {
		ctx = context.Background()
		fmt.Printf("→ Running: %s (in %s, no timeout)\n", strings.Join(fullCommand, " "), dockerDir)
	}

	// Execute the command with timeout
	cmd := exec.CommandContext(ctx, fullCommand[0], fullCommand[1:]...)
	cmd.Dir = dockerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	// Setup Bake environment for build commands
	if dockerCmd.Name == "up" || dockerCmd.Name == "build" {
		setupBakeEnvironment(cmd)
	}

	err = cmd.Run()
	
	// Check if the command was cancelled due to timeout
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timed out after %v", dockerCmd.Timeout)
	}
	
	return err
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

// detectBakeSupport checks if Docker Bake is available on the system
func detectBakeSupport() bool {
	bakeMutex.Lock()
	defer bakeMutex.Unlock()
	
	// Return cached result if we've already checked
	if bakeSupported != nil {
		return *bakeSupported
	}
	
	supported := false
	
	// Check if docker buildx is available
	if _, err := exec.LookPath("docker"); err == nil {
		// Test if buildx bake command exists
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		cmd := exec.CommandContext(ctx, "docker", "buildx", "bake", "--help")
		cmd.Stdout = nil // Suppress output
		cmd.Stderr = nil // Suppress errors
		
		if err := cmd.Run(); err == nil {
			supported = true
		}
	}
	
	// Cache the result
	bakeSupported = &supported
	return supported
}

// setupBakeEnvironment sets up environment variables for optimal build performance
func setupBakeEnvironment(cmd *exec.Cmd) {
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	
	// Check if COMPOSE_BAKE is already set by user
	bakeAlreadySet := false
	for _, env := range cmd.Env {
		if strings.HasPrefix(env, "COMPOSE_BAKE=") {
			bakeAlreadySet = true
			break
		}
	}
	
	// Only set if not already configured by user
	if !bakeAlreadySet {
		if detectBakeSupport() {
			cmd.Env = append(cmd.Env, "COMPOSE_BAKE=true")
		} else {
			cmd.Env = append(cmd.Env, "COMPOSE_BAKE=false")
		}
	}
}
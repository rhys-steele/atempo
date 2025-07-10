package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
	bakeSupported *bool
	bakeMutex     sync.Mutex
)

// DockerCommand represents available Docker operations
type DockerCommand struct {
	Name        string
	Description string
	Args        []string
	Timeout     time.Duration // Default timeout for this command
}

// validateDockerAvailability checks if Docker and Docker Compose are available
func validateDockerAvailability() error {
	// Check if Docker is running
	if err := exec.Command("docker", "info").Run(); err != nil {
		return fmt.Errorf("Docker is not running or not available: %w", err)
	}

	// Check if Docker Compose is available
	if err := exec.Command("docker-compose", "--version").Run(); err != nil {
		return fmt.Errorf("Docker Compose is not available: %w", err)
	}

	return nil
}

// Common Docker commands for Atempo projects
var SupportedCommands = map[string]DockerCommand{
	"up": {
		Name:        "up",
		Description: "Start services in detached mode (use --build to force rebuild)",
		Args:        []string{"up", "-d"},
		Timeout:     3 * time.Minute, // Shorter timeout since no building by default
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
	// Validate Docker availability before proceeding
	if err := validateDockerAvailability(); err != nil {
		return fmt.Errorf("Docker validation failed: %w", err)
	}

	// Check if --build flag is present for up command and adjust timeout
	if dockerCmd.Name == "up" && containsBuildFlag(additionalArgs) {
		dockerCmd.Timeout = 10 * time.Minute // Increase timeout when building
	}
	// Resolve project path
	resolvedPath, err := utils.ResolveProjectPath(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Find docker-compose.yml file using shared utility
	composeFile := utils.FindDockerComposeFile(resolvedPath)
	if composeFile == "" {
		return fmt.Errorf("docker-compose.yml not found in %s or %s/infra/docker", resolvedPath, resolvedPath)
	}
	dockerDir := resolvedPath

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
		fmt.Printf("⎿ Running: %s\n", strings.Join(fullCommand, " "))
	} else {
		ctx = context.Background()
		fmt.Printf("⎿ Running: %s\n", strings.Join(fullCommand, " "))
	}

	// Execute the command with timeout
	cmd := exec.CommandContext(ctx, fullCommand[0], fullCommand[1:]...)
	cmd.Dir = dockerDir
	cmd.Stdin = os.Stdin

	// Setup Bake environment for build commands
	if dockerCmd.Name == "up" || dockerCmd.Name == "build" {
		setupBakeEnvironment(cmd)
	}

	// For 'up' and 'down' commands, filter output to reduce verbosity
	if dockerCmd.Name == "up" || dockerCmd.Name == "down" {
		err = executeWithFilteredOutput(cmd)
	} else {
		// For other commands, use normal output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}

	// Check if the command was cancelled due to timeout
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timed out after %v", dockerCmd.Timeout)
	}

	return err
}

// executeWithFilteredOutput runs a command and filters Docker output for cleaner display
func executeWithFilteredOutput(cmd *exec.Cmd) error {
	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Channel to collect errors from goroutines
	errChan := make(chan error, 2)

	// Process stdout
	go func() {
		errChan <- filterDockerOutput(stdout, false)
	}()

	// Process stderr
	go func() {
		errChan <- filterDockerOutput(stderr, true)
	}()

	// Wait for both goroutines to complete
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			// Log error but don't fail the command
			fmt.Printf("Output processing error: %v\n", err)
		}
	}

	// Wait for the command to complete
	return cmd.Wait()
}

// filterDockerOutput filters Docker output to show only relevant information
func filterDockerOutput(reader io.Reader, isStderr bool) error {
	scanner := bufio.NewScanner(reader)
	var finalContainerStatus []string

	for scanner.Scan() {
		line := scanner.Text()

		// Skip verbose build output
		if shouldSkipBuildLine(line) {
			continue
		}

		// Skip intermediate status messages - only collect final status
		if strings.Contains(line, "Container") && (strings.Contains(line, "Running") || strings.Contains(line, "Started")) {
			finalContainerStatus = append(finalContainerStatus, line)
			continue
		}

		// For down command, show final container removal status
		if strings.Contains(line, "Container") && strings.Contains(line, "Removed") {
			finalContainerStatus = append(finalContainerStatus, line)
			continue
		}

		// Skip all creating/starting/created intermediate messages
		if shouldSkipIntermediateStatus(line) {
			continue
		}

		// Skip warnings we don't want to show
		if shouldSkipWarning(line) {
			continue
		}

		// Show critical errors only
		if isStderr && strings.Contains(line, "ERROR[") {
			fmt.Fprintln(os.Stderr, line)
			continue
		}

		// Show build completion
		if strings.Contains(line, "Built") && !strings.Contains(line, "[+] Building") {
			if parts := strings.Fields(line); len(parts) >= 2 {
				serviceName := parts[0]
				statusText := parts[len(parts)-1]
				paddedName := fmt.Sprintf("%-25s", serviceName)
				fmt.Printf(" %s \033[32m%s\033[0m\n", paddedName, strings.ToLower(statusText))
			}
		}
	}

	// Show final container status in a clean format
	for _, status := range finalContainerStatus {
		formattedLine := formatContainerStatusLine(status)
		// Write directly to os.Stdout and force colors
		if parts := strings.Fields(status); len(parts) >= 3 {
			containerName := parts[1]
			statusText := parts[len(parts)-1]
			paddedName := fmt.Sprintf("%-25s", containerName)

			// Apply colors directly when writing to terminal
			switch strings.ToLower(statusText) {
			case "running", "started", "created", "built":
				fmt.Printf(" %s \033[32m%s\033[0m\n", paddedName, strings.ToLower(statusText))
			case "stopped", "removed":
				fmt.Printf(" %s \033[31m%s\033[0m\n", paddedName, strings.ToLower(statusText))
			case "creating", "starting", "stopping", "removing":
				fmt.Printf(" %s \033[33m%s\033[0m\n", paddedName, strings.ToLower(statusText))
			default:
				fmt.Printf(" %s %s\n", paddedName, strings.ToLower(statusText))
			}
		} else {
			fmt.Fprintln(os.Stdout, formattedLine)
		}
	}

	return scanner.Err()
}

// shouldSkipBuildLine determines if a build output line should be skipped
func shouldSkipBuildLine(line string) bool {
	// Skip Docker build steps
	if strings.Contains(line, "=> [internal] load") ||
		strings.Contains(line, "=> => reading") ||
		strings.Contains(line, "=> => transferring") ||
		strings.Contains(line, "=> FROM docker.io") ||
		strings.Contains(line, "=> [stage-0") ||
		strings.Contains(line, "=> CACHED") ||
		strings.Contains(line, "=> exporting") ||
		strings.Contains(line, "=> => exporting layers") ||
		strings.Contains(line, "=> => writing image") ||
		strings.Contains(line, "=> => naming to") ||
		strings.Contains(line, "=> resolving provenance") {
		return true
	}

	// Skip build progress indicators
	if strings.Contains(line, "[+] Building") && strings.Contains(line, "FINISHED") {
		return true
	}

	return false
}

// shouldShowStatusLine determines if a status line should be shown
func shouldShowStatusLine(line string) bool {
	// Show all container status lines
	if strings.Contains(line, "Container") {
		return true
	}

	// Show service build status
	if strings.Contains(line, "Built") {
		return true
	}

	// Show network and volume status
	if strings.Contains(line, "Network") || strings.Contains(line, "Volume") {
		return true
	}

	// Show final summary
	if strings.Contains(line, "[+] Running") && strings.Contains(line, "/") {
		return true
	}

	return false
}

// shouldSkipIntermediateStatus determines if an intermediate status line should be skipped
func shouldSkipIntermediateStatus(line string) bool {
	// Skip all intermediate container status messages
	if strings.Contains(line, "Container") && (strings.Contains(line, "Creating") ||
		strings.Contains(line, "Created") ||
		strings.Contains(line, "Starting") ||
		strings.Contains(line, "Stopping") ||
		strings.Contains(line, "Stopped") ||
		strings.Contains(line, "Removing") ||
		strings.Contains(line, "Removed")) {
		return true
	}

	// Skip network and volume creation messages
	if strings.Contains(line, "Network") || strings.Contains(line, "Volume") {
		return true
	}

	// Skip running progress indicators
	if strings.Contains(line, "[+] Running") {
		return true
	}

	return false
}

// shouldSkipWarning determines if a warning should be skipped
func shouldSkipWarning(line string) bool {
	// Skip version warnings
	if strings.Contains(line, "version") && strings.Contains(line, "obsolete") {
		return true
	}

	// Skip platform warnings
	if strings.Contains(line, "platform") && strings.Contains(line, "does not match") {
		return true
	}

	// Skip other docker-compose warnings
	if strings.Contains(line, "WARN[") {
		return true
	}

	return false
}

// formatContainerStatusLine formats container status lines with proper alignment and coloring
func formatContainerStatusLine(line string) string {
	// Handle container status lines (e.g., " Container my-app-2-redis  Running")
	if strings.Contains(line, "Container") {
		// Extract container name and status
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			containerName := parts[1]
			status := parts[len(parts)-1]

			// Format with proper alignment (pad to 25 characters)
			paddedName := fmt.Sprintf("%-25s", containerName)

			// Apply color based on status
			coloredStatus := formatStatusWithColor(status)

			return fmt.Sprintf(" %s %s", paddedName, coloredStatus)
		}
	}

	// Handle network and volume lines
	if strings.Contains(line, "Network") || strings.Contains(line, "Volume") {
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			resourceType := parts[0] // "Network" or "Volume"
			resourceName := parts[1]
			status := parts[len(parts)-1]

			// Format with proper alignment
			paddedName := fmt.Sprintf("%-25s", resourceName)
			coloredStatus := formatStatusWithColor(status)

			return fmt.Sprintf(" %s %s %s", resourceType, paddedName, coloredStatus)
		}
	}

	// Handle built status lines
	if strings.Contains(line, "Built") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			serviceName := parts[0]
			status := parts[len(parts)-1]

			paddedName := fmt.Sprintf("%-25s", serviceName)
			coloredStatus := formatStatusWithColor(status)

			return fmt.Sprintf(" %s %s", paddedName, coloredStatus)
		}
	}

	// Return original line if no specific formatting needed
	return line
}

// formatStatusWithColor applies color to status text
func formatStatusWithColor(status string) string {
	statusLower := strings.ToLower(status)

	// Force colors for interactive CLI usage
	const (
		Green  = "\033[32m"
		Red    = "\033[31m"
		Yellow = "\033[33m"
		Reset  = "\033[0m"
	)

	switch statusLower {
	case "running", "started", "created", "built":
		return Green + statusLower + Reset
	case "stopped", "removed":
		return Red + statusLower + Reset
	case "creating", "starting", "stopping", "removing":
		return Yellow + statusLower + Reset
	default:
		return statusLower
	}
}

// ExecuteExecCommand runs a command inside a container (docker-compose exec)
func ExecuteExecCommand(service string, projectPath string, cmdArgs []string) error {
	// Validate Docker availability before proceeding
	if err := validateDockerAvailability(); err != nil {
		return fmt.Errorf("Docker validation failed: %w", err)
	}

	// Resolve project path
	resolvedPath, err := utils.ResolveProjectPath(projectPath)
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
	resolvedPath, err := utils.ResolveProjectPath(projectPath)
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
	resolvedPath, err := utils.ResolveProjectPath(projectPath)
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

// containsBuildFlag checks if --build flag is present in additional arguments
func containsBuildFlag(args []string) bool {
	return utils.ContainsFlag(args, "--build")
}

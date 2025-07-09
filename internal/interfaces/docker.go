package interfaces

import "atempo/internal/registry"

// DockerClient defines the interface for Docker operations
type DockerClient interface {
	// ExecuteCommand executes a docker-compose command
	ExecuteCommand(cmd string, projectPath string, args []string) error
	
	// ExecuteExecCommand executes a command inside a running container
	ExecuteExecCommand(serviceName, projectPath string, args []string) error
	
	// ValidateAvailability checks if Docker and Docker Compose are available
	ValidateAvailability() error
	
	// GetServices returns the list of services for a project
	GetServices(projectPath string) ([]registry.Service, error)
	
	// GetProjectStatus returns the overall status of a project's Docker services
	GetProjectStatus(projectPath string) (string, error)
	
	// SupportsBake checks if Docker Bake is supported
	SupportsBake() bool
	
	// DetectFrameworkFromCompose detects the framework from docker-compose.yml
	DetectFrameworkFromCompose(projectPath string) (string, error)
}
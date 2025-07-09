package services

import (
	"context"
	"fmt"

	"atempo/internal/docker"
	"atempo/internal/registry"
)

// dockerService implements DockerService
type dockerService struct{}

// NewDockerService creates a new DockerService implementation
func NewDockerService() DockerService {
	return &dockerService{}
}

// ExecuteCommand executes a docker-compose command in the project directory
func (s *dockerService) ExecuteCommand(ctx context.Context, cmd string, projectPath string, args []string) error {
	return docker.ExecuteCommand(cmd, projectPath, args)
}

// ValidateAvailability checks if Docker and Docker Compose are available
func (s *dockerService) ValidateAvailability(ctx context.Context) error {
	return docker.ValidateAvailability()
}

// GetServices returns the list of services for a project
func (s *dockerService) GetServices(ctx context.Context, projectPath string) ([]registry.Service, error) {
	services, err := docker.GetServices(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker services: %w", err)
	}
	
	// Convert docker.Service to registry.Service
	registryServices := make([]registry.Service, len(services))
	for i, service := range services {
		registryServices[i] = registry.Service{
			Name:   service.Name,
			Status: service.Status,
			URL:    service.URL,
		}
	}
	
	return registryServices, nil
}

// GetProjectStatus returns the overall status of a project's Docker services
func (s *dockerService) GetProjectStatus(ctx context.Context, projectPath string) (string, error) {
	services, err := s.GetServices(ctx, projectPath)
	if err != nil {
		return "docker-error", fmt.Errorf("failed to get services: %w", err)
	}
	
	if len(services) == 0 {
		return "no-services", nil
	}
	
	runningCount := 0
	for _, service := range services {
		if service.Status == "running" {
			runningCount++
		}
	}
	
	if runningCount == 0 {
		return "stopped", nil
	} else if runningCount == len(services) {
		return "running", nil
	} else {
		return "partial", nil
	}
}

// ExecuteInContainer executes a command inside a running container
func (s *dockerService) ExecuteInContainer(ctx context.Context, service, projectPath string, args []string) error {
	return docker.ExecuteExecCommand(service, projectPath, args)
}
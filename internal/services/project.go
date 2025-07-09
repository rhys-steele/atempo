package services

import (
	"context"
	"fmt"
	"path/filepath"

	"atempo/internal/registry"
	"atempo/internal/scaffold"
)

// projectService implements ProjectService
type projectService struct {
	dockerService DockerService
}

// NewProjectService creates a new ProjectService implementation
func NewProjectService(dockerService DockerService) ProjectService {
	return &projectService{
		dockerService: dockerService,
	}
}

// Create creates a new project with the given configuration
func (s *projectService) Create(ctx context.Context, req CreateProjectRequest) (*registry.Project, error) {
	// Create scaffold configuration
	config := scaffold.Config{
		Framework:   req.Framework,
		Version:     req.Version,
		ProjectPath: req.ProjectPath,
		ProjectName: req.ProjectName,
		EnableAI:    req.EnableAI,
		AIManifest:  req.AIManifest,
		Interactive: req.Interactive,
	}
	
	// Run the scaffolding process
	if err := scaffold.Run(&config); err != nil {
		return nil, fmt.Errorf("failed to scaffold project: %w", err)
	}
	
	// Load the registry to add the new project
	reg, err := registry.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}
	
	// Add project to registry
	project := &registry.Project{
		Name:      req.ProjectName,
		Path:      req.ProjectPath,
		Framework: req.Framework,
		Version:   req.Version,
		Status:    "created",
	}
	
	if err := reg.AddProject(project); err != nil {
		return nil, fmt.Errorf("failed to add project to registry: %w", err)
	}
	
	return project, nil
}

// GetByName retrieves a project by name from the registry
func (s *projectService) GetByName(ctx context.Context, name string) (*registry.Project, error) {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}
	
	project, err := reg.FindProject(name)
	if err != nil {
		return nil, fmt.Errorf("project '%s' not found: %w", name, err)
	}
	
	return project, nil
}

// List returns all projects in the registry
func (s *projectService) List(ctx context.Context) ([]registry.Project, error) {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}
	
	return reg.ListProjects(), nil
}

// Delete removes a project from the registry
func (s *projectService) Delete(ctx context.Context, name string) error {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}
	
	// Check if project exists
	_, err = reg.FindProject(name)
	if err != nil {
		return fmt.Errorf("project '%s' not found: %w", name, err)
	}
	
	// Remove project from registry
	if err := reg.RemoveProject(name); err != nil {
		return fmt.Errorf("failed to remove project from registry: %w", err)
	}
	
	return nil
}

// UpdateStatus updates the status of a project
func (s *projectService) UpdateStatus(ctx context.Context, name string) error {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}
	
	// Update project status using the registry's existing method
	if err := reg.UpdateProjectStatus(name); err != nil {
		return fmt.Errorf("failed to update project status: %w", err)
	}
	
	return nil
}

// FindByPath attempts to find a project by path
func (s *projectService) FindByPath(ctx context.Context, path string) (*registry.Project, error) {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}
	
	projects := reg.ListProjects()
	for _, project := range projects {
		// Compare absolute paths
		if absPath, err := filepath.Abs(project.Path); err == nil {
			if targetPath, err := filepath.Abs(path); err == nil {
				if absPath == targetPath {
					return &project, nil
				}
			}
		}
	}
	
	return nil, fmt.Errorf("no project found at path: %s", path)
}
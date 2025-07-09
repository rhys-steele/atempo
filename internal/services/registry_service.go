package services

import (
	"fmt"
	"time"

	"atempo/internal/interfaces"
	"atempo/internal/registry"
	"atempo/internal/repositories"
)

// RegistryService implements the registry business logic using the repository pattern
type RegistryService struct {
	repo interfaces.RegistryRepository
}

// NewRegistryService creates a new registry service with the given repository
func NewRegistryService(repo interfaces.RegistryRepository) *RegistryService {
	return &RegistryService{
		repo: repo,
	}
}

// NewDefaultRegistryService creates a new registry service with default file repository and cache
func NewDefaultRegistryService() *RegistryService {
	cache := repositories.NewMemoryCacheRepository()
	repo := repositories.NewFileRegistryRepository(cache)
	return NewRegistryService(repo)
}

// LoadRegistry loads the project registry
func (s *RegistryService) LoadRegistry() (*registry.Registry, error) {
	return s.repo.LoadRegistry()
}

// SaveRegistry saves the project registry
func (s *RegistryService) SaveRegistry(reg *registry.Registry) error {
	return s.repo.SaveRegistry(reg)
}

// AddProject adds a new project to the registry
func (s *RegistryService) AddProject(project *registry.Project) error {
	return s.repo.AddProject(project)
}

// FindProject finds a project by name
func (s *RegistryService) FindProject(name string) (*registry.Project, error) {
	return s.repo.FindProject(name)
}

// RemoveProject removes a project from the registry
func (s *RegistryService) RemoveProject(name string) error {
	return s.repo.RemoveProject(name)
}

// ListProjects returns all projects in the registry
func (s *RegistryService) ListProjects() ([]registry.Project, error) {
	return s.repo.ListProjects()
}

// UpdateProjectStatus updates the status of a project
func (s *RegistryService) UpdateProjectStatus(name string) error {
	return s.repo.UpdateProjectStatus(name)
}

// ScanForProjects scans a directory for Atempo projects
func (s *RegistryService) ScanForProjects(scanPath string) error {
	return s.repo.ScanForProjects(scanPath)
}

// GetProjectByPath finds a project by its path
func (s *RegistryService) GetProjectByPath(projectPath string) (*registry.Project, error) {
	projects, err := s.ListProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to load projects: %w", err)
	}

	for _, project := range projects {
		if project.Path == projectPath {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("no project found at path: %s", projectPath)
}

// UpdateProject updates an existing project in the registry
func (s *RegistryService) UpdateProject(project *registry.Project) error {
	reg, err := s.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	for i, p := range reg.Projects {
		if p.Name == project.Name {
			project.LastAccessed = time.Now()
			reg.Projects[i] = *project
			return s.SaveRegistry(reg)
		}
	}

	return fmt.Errorf("project '%s' not found in registry", project.Name)
}

// GetProjectCount returns the total number of registered projects
func (s *RegistryService) GetProjectCount() (int, error) {
	projects, err := s.ListProjects()
	if err != nil {
		return 0, fmt.Errorf("failed to load projects: %w", err)
	}

	return len(projects), nil
}

// GetActiveProjects returns projects that are currently running
func (s *RegistryService) GetActiveProjects() ([]registry.Project, error) {
	projects, err := s.ListProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to load projects: %w", err)
	}

	var activeProjects []registry.Project
	for _, project := range projects {
		if project.Status == "running" || project.Status == "healthy" {
			activeProjects = append(activeProjects, project)
		}
	}

	return activeProjects, nil
}

// GetRecentProjects returns the most recently accessed projects
func (s *RegistryService) GetRecentProjects(limit int) ([]registry.Project, error) {
	projects, err := s.ListProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to load projects: %w", err)
	}

	// Sort projects by last accessed time (most recent first)
	for i := 0; i < len(projects)-1; i++ {
		for j := i + 1; j < len(projects); j++ {
			if projects[i].LastAccessed.Before(projects[j].LastAccessed) {
				projects[i], projects[j] = projects[j], projects[i]
			}
		}
	}

	// Return up to the limit
	if limit > 0 && limit < len(projects) {
		return projects[:limit], nil
	}

	return projects, nil
}

// ProjectExists checks if a project with the given name exists
func (s *RegistryService) ProjectExists(name string) (bool, error) {
	_, err := s.FindProject(name)
	if err != nil {
		if err.Error() == fmt.Sprintf("project '%s' not found in registry", name) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetProjectStats returns statistics about the registry
func (s *RegistryService) GetProjectStats() (*ProjectStats, error) {
	projects, err := s.ListProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to load projects: %w", err)
	}

	stats := &ProjectStats{
		TotalProjects: len(projects),
		Frameworks:    make(map[string]int),
		Statuses:      make(map[string]int),
	}

	for _, project := range projects {
		stats.Frameworks[project.Framework]++
		stats.Statuses[project.Status]++

		if project.Status == "running" || project.Status == "healthy" {
			stats.ActiveProjects++
		}
	}

	return stats, nil
}

// ProjectStats represents statistics about the project registry
type ProjectStats struct {
	TotalProjects  int            `json:"total_projects"`
	ActiveProjects int            `json:"active_projects"`
	Frameworks     map[string]int `json:"frameworks"`
	Statuses       map[string]int `json:"statuses"`
}
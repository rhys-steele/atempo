package interfaces

import "atempo/internal/registry"

// RegistryRepository defines the interface for registry data access operations
type RegistryRepository interface {
	// LoadRegistry loads the project registry from storage
	LoadRegistry() (*registry.Registry, error)
	
	// SaveRegistry saves the project registry to storage
	SaveRegistry(reg *registry.Registry) error
	
	// AddProject adds a new project to the registry
	AddProject(project *registry.Project) error
	
	// FindProject finds a project by name
	FindProject(name string) (*registry.Project, error)
	
	// RemoveProject removes a project from the registry
	RemoveProject(name string) error
	
	// ListProjects returns all projects in the registry
	ListProjects() ([]registry.Project, error)
	
	// UpdateProjectStatus updates the status of a project
	UpdateProjectStatus(name string) error
	
	// ScanForProjects scans a directory for Atempo projects
	ScanForProjects(scanPath string) error
}
package interfaces

import "atempo/internal/registry"

// RegistryRepository defines all registry data access operations.
type RegistryRepository interface {
	LoadRegistry() (*registry.Registry, error)
	SaveRegistry(reg *registry.Registry) error
	AddProject(project *registry.Project) error
	FindProject(name string) (*registry.Project, error)
	RemoveProject(name string) error
	ListProjects() ([]registry.Project, error)
}

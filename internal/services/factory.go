package services

// ServiceFactory provides a centralized way to create and manage service instances
type ServiceFactory struct {
	pathResolver PathResolverService
	docker       DockerService
	project      ProjectService
}

// NewServiceFactory creates a new ServiceFactory with all service dependencies
func NewServiceFactory() *ServiceFactory {
	// Create service instances with proper dependencies
	pathResolver := NewPathResolverService()
	docker := NewDockerService()
	project := NewProjectService(docker)
	
	return &ServiceFactory{
		pathResolver: pathResolver,
		docker:       docker,
		project:      project,
	}
}

// PathResolver returns the PathResolverService instance
func (f *ServiceFactory) PathResolver() PathResolverService {
	return f.pathResolver
}

// Docker returns the DockerService instance
func (f *ServiceFactory) Docker() DockerService {
	return f.docker
}

// Project returns the ProjectService instance
func (f *ServiceFactory) Project() ProjectService {
	return f.project
}

// Services returns all services in a convenient struct
func (f *ServiceFactory) Services() Services {
	return Services{
		PathResolver: f.pathResolver,
		Docker:       f.docker,
		Project:      f.project,
	}
}

// Services groups all service interfaces together for easy access
type Services struct {
	PathResolver PathResolverService
	Docker       DockerService
	Project      ProjectService
}
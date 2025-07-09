package services

import (
	"context"

	"atempo/internal/registry"
)

// ProjectService provides business operations for project management
type ProjectService interface {
	// Create creates a new project with the given configuration
	Create(ctx context.Context, req CreateProjectRequest) (*registry.Project, error)
	
	// GetByName retrieves a project by name from the registry
	GetByName(ctx context.Context, name string) (*registry.Project, error)
	
	// List returns all projects in the registry
	List(ctx context.Context) ([]registry.Project, error)
	
	// Delete removes a project from the registry
	Delete(ctx context.Context, name string) error
	
	// UpdateStatus updates the status of a project
	UpdateStatus(ctx context.Context, name string) error
	
	// FindByPath attempts to find a project by path
	FindByPath(ctx context.Context, path string) (*registry.Project, error)
}

// DockerService provides business operations for Docker management
type DockerService interface {
	// ExecuteCommand executes a docker-compose command in the project directory
	ExecuteCommand(ctx context.Context, cmd string, projectPath string, args []string) error
	
	// ValidateAvailability checks if Docker and Docker Compose are available
	ValidateAvailability(ctx context.Context) error
	
	// GetServices returns the list of services for a project
	GetServices(ctx context.Context, projectPath string) ([]registry.Service, error)
	
	// GetProjectStatus returns the overall status of a project's Docker services
	GetProjectStatus(ctx context.Context, projectPath string) (string, error)
	
	// ExecuteInContainer executes a command inside a running container
	ExecuteInContainer(ctx context.Context, service, projectPath string, args []string) error
}

// ScaffoldService provides business operations for project scaffolding
type ScaffoldService interface {
	// CreateProject scaffolds a new project using the specified framework and version
	CreateProject(ctx context.Context, req CreateProjectRequest) error
	
	// ValidateFramework checks if a framework and version combination is supported
	ValidateFramework(ctx context.Context, framework, version string) error
	
	// GetSupportedFrameworks returns a list of supported frameworks
	GetSupportedFrameworks(ctx context.Context) ([]FrameworkInfo, error)
}

// TemplateService provides business operations for template management
type TemplateService interface {
	// ProcessTemplate processes a template with variable substitution
	ProcessTemplate(ctx context.Context, req ProcessTemplateRequest) error
	
	// LoadTemplate loads template configuration from the specified path
	LoadTemplate(ctx context.Context, templatePath string) (*TemplateConfig, error)
	
	// ValidateTemplate validates that a template configuration is correct
	ValidateTemplate(ctx context.Context, config *TemplateConfig) error
}

// PathResolverService provides business operations for path resolution
type PathResolverService interface {
	// ResolveProjectFromArgs resolves a project identifier from command arguments
	ResolveProjectFromArgs(ctx context.Context, args []string) (*ProjectResolution, error)
	
	// ResolveProjectPath resolves a project identifier to an absolute path
	ResolveProjectPath(ctx context.Context, identifier string) (string, error)
	
	// GetCurrentProjectPath returns the current working directory as a project path
	GetCurrentProjectPath(ctx context.Context) (*ProjectResolution, error)
}

// Request/Response types for service operations

// CreateProjectRequest represents a request to create a new project
type CreateProjectRequest struct {
	Framework    string
	Version      string
	ProjectPath  string
	ProjectName  string
	EnableAI     bool
	AIManifest   string
	Interactive  bool
}

// ProcessTemplateRequest represents a request to process a template
type ProcessTemplateRequest struct {
	SourcePath      string
	DestinationPath string
	Variables       map[string]string
	SkipFiles       []string
}

// ProjectResolution represents the result of resolving a project identifier
type ProjectResolution struct {
	Path string // The resolved absolute path to the project
	Name string // The project name (either from args or directory basename)
}

// TemplateConfig represents the configuration of a project template
type TemplateConfig struct {
	Name      string                    `json:"name"`
	Framework string                    `json:"framework"`
	Language  string                    `json:"language"`
	Installer InstallerConfig           `json:"installer"`
	Services  map[string]ServiceConfig  `json:"services"`
	Variables map[string]string         `json:"variables,omitempty"`
}

// InstallerConfig represents the installer configuration for a template
type InstallerConfig struct {
	Type     string   `json:"type"`
	Command  []string `json:"command"`
	WorkDir  string   `json:"work-dir,omitempty"`
	Timeout  int      `json:"timeout,omitempty"`
}

// ServiceConfig represents the configuration for a Docker service
type ServiceConfig struct {
	Type       string            `json:"type"`
	Image      string            `json:"image,omitempty"`
	Dockerfile string            `json:"dockerfile,omitempty"`
	Ports      []string          `json:"ports,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Volumes    []string          `json:"volumes,omitempty"`
	DependsOn  []string          `json:"depends_on,omitempty"`
}

// FrameworkInfo represents information about a supported framework
type FrameworkInfo struct {
	Name            string   `json:"name"`
	Language        string   `json:"language"`
	SupportedVersions []string `json:"supported_versions"`
	Description     string   `json:"description"`
}
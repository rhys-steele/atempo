package interfaces

// ScaffoldEngine defines the interface for project scaffolding operations
type ScaffoldEngine interface {
	// Run executes the complete scaffolding process
	Run(config *ScaffoldConfig) error
	
	// ValidateFramework validates that a framework and version combination is supported
	ValidateFramework(framework, version string) error
	
	// LoadTemplate loads template configuration from the specified path
	LoadTemplate(templatePath string) (*TemplateConfig, error)
	
	// ProcessTemplate processes template files with variable substitution
	ProcessTemplate(sourcePath, destPath string, variables map[string]string) error
	
	// RunInstaller executes the framework-specific installer
	RunInstaller(config *InstallerConfig, projectPath string) error
	
	// SetupServices sets up Docker services for the project
	SetupServices(projectPath string) error
}

// ScaffoldConfig represents the configuration for scaffolding a project
type ScaffoldConfig struct {
	Framework   string
	Version     string
	ProjectPath string
	ProjectName string
	EnableAI    bool
	AIManifest  string
	Interactive bool
}

// TemplateConfig represents the configuration of a project template
type TemplateConfig struct {
	Name      string                   `json:"name"`
	Framework string                   `json:"framework"`
	Language  string                   `json:"language"`
	Installer InstallerConfig          `json:"installer"`
	Services  map[string]ServiceConfig `json:"services"`
	Variables map[string]string        `json:"variables,omitempty"`
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
	Type        string            `json:"type"`
	Image       string            `json:"image,omitempty"`
	Dockerfile  string            `json:"dockerfile,omitempty"`
	Ports       []string          `json:"ports,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Volumes     []string          `json:"volumes,omitempty"`
	DependsOn   []string          `json:"depends_on,omitempty"`
}
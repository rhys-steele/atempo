package interfaces

// FrameworkDetector defines the interface for detecting project frameworks
type FrameworkDetector interface {
	// DetectFramework detects the framework of a project from its files
	DetectFramework(projectPath string) (string, error)
	
	// DetectFromCompose detects framework from docker-compose.yml
	DetectFromCompose(projectPath string) (string, error)
	
	// DetectFromFiles detects framework from project files
	DetectFromFiles(projectPath string) (string, error)
	
	// IsFramework checks if a path contains a specific framework
	IsFramework(projectPath, framework string) bool
	
	// GetSupportedFrameworks returns a list of supported frameworks
	GetSupportedFrameworks() []string
	
	// ValidateFramework validates that a framework is supported
	ValidateFramework(framework string) error
}

// FrameworkInfo represents information about a supported framework
type FrameworkInfo struct {
	Name              string   `json:"name"`
	Language          string   `json:"language"`
	SupportedVersions []string `json:"supported_versions"`
	Description       string   `json:"description"`
	DetectionFiles    []string `json:"detection_files"`
	DefaultVersion    string   `json:"default_version"`
}
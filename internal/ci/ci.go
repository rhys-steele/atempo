package ci

import (
	"embed"
	"time"
)

// CIProvider represents the CI/CD provider
type CIProvider string

const (
	ProviderGitHub CIProvider = "github"
	ProviderGitLab CIProvider = "gitlab"
)

// CIConfig represents the main CI configuration structure
type CIConfig struct {
	Provider      CIProvider              `json:"provider"`
	Framework     string                  `json:"framework"`        // laravel, django, express
	ProjectName   string                  `json:"project_name"`
	ProjectPath   string                  `json:"project_path"`     // Absolute path to project
	RepoURL       string                  `json:"repo_url"`
	Settings      map[string]interface{}  `json:"settings"`         // Provider-specific settings
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	LastRunAt     *time.Time             `json:"last_run_at,omitempty"`
	LastRunStatus string                 `json:"last_run_status,omitempty"` // success, failed, running
}

// Provider interface for extensibility
type Provider interface {
	Name() string
	ValidateConfig(config *CIConfig) error
	GenerateConfig(config *CIConfig, templateFS embed.FS) ([]byte, error)
	GetConfigFileName() string                              // .github/workflows/ci.yml, .gitlab-ci.yml
	GetConfigPath(projectPath string) string               // Full path where config should be written
	GetDefaultSettings(framework string) map[string]interface{}
	PromptForSettings(framework string) (map[string]interface{}, error)
	SupportedFrameworks() []string
}

// Runner interface for local execution
type Runner interface {
	Run(projectPath string, provider CIProvider) (*RunResult, error)
	Validate(projectPath string, provider CIProvider) (*ValidationResult, error)
	IsAvailable(provider CIProvider) bool                  // Check if Docker images are available
	GetRequiredImages(provider CIProvider) []string
	PullRequiredImages(provider CIProvider) error
}

// RunResult represents the result of a CI run
type RunResult struct {
	Provider      CIProvider    `json:"provider"`
	Status        string        `json:"status"`        // success, failed, timeout, cancelled
	Duration      time.Duration `json:"duration"`
	Output        string        `json:"output"`        // Full command output
	ExitCode      int          `json:"exit_code"`
	StartedAt     time.Time    `json:"started_at"`
	FinishedAt    time.Time    `json:"finished_at"`
	Steps         []StepResult `json:"steps"`         // Individual step results
	ArtifactPaths []string     `json:"artifact_paths"` // Generated artifacts
	LogPath       string       `json:"log_path"`      // Full log file location
}

// StepResult represents the result of an individual CI step
type StepResult struct {
	Name       string        `json:"name"`
	Status     string        `json:"status"`
	Duration   time.Duration `json:"duration"`
	Output     string        `json:"output"`
	ExitCode   int          `json:"exit_code"`
	StartedAt  time.Time    `json:"started_at"`
	FinishedAt time.Time    `json:"finished_at"`
}

// ValidationResult represents the result of CI configuration validation
type ValidationResult struct {
	Valid      bool     `json:"valid"`
	Errors     []string `json:"errors,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
	ConfigPath string   `json:"config_path"`
	Provider   CIProvider `json:"provider"`
}
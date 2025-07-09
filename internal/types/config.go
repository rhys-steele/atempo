package types

import (
	"os"
	"path/filepath"
	"time"
)

// Configuration represents the application configuration
type Configuration struct {
	// Registry settings
	RegistryPath        string `json:"registry_path"`
	AutoScanProjects    bool   `json:"auto_scan_projects"`
	CacheTimeout        int    `json:"cache_timeout_minutes"`
	
	// Docker settings
	DockerTimeout       int    `json:"docker_timeout_seconds"`
	DockerBakeSupport   bool   `json:"docker_bake_support"`
	
	// AI settings
	AIProvider          string `json:"ai_provider"`
	AIApiKey            string `json:"ai_api_key,omitempty"`
	AIEndpoint          string `json:"ai_endpoint,omitempty"`
	
	// UI settings
	UseColors           bool   `json:"use_colors"`
	ShowProgress        bool   `json:"show_progress"`
	VerboseLogging      bool   `json:"verbose_logging"`
	
	// Framework settings
	DefaultFramework    string            `json:"default_framework"`
	FrameworkVersions   map[string]string `json:"framework_versions"`
	
	// Security settings
	ValidateSSL         bool   `json:"validate_ssl"`
	
	// System settings
	LogLevel            string `json:"log_level"`
	LastUpdated         time.Time `json:"last_updated"`
}

// DefaultConfiguration returns the default configuration
func DefaultConfiguration() *Configuration {
	homeDir, _ := os.UserHomeDir()
	
	return &Configuration{
		RegistryPath:        filepath.Join(homeDir, ".atempo", "registry.json"),
		AutoScanProjects:    false,
		CacheTimeout:        5,
		DockerTimeout:       300,
		DockerBakeSupport:   true,
		AIProvider:          "claude",
		UseColors:           true,
		ShowProgress:        true,
		VerboseLogging:      false,
		DefaultFramework:    "laravel",
		FrameworkVersions:   map[string]string{
			"laravel": "11",
			"django":  "5.0",
		},
		ValidateSSL:         true,
		LogLevel:            "info",
		LastUpdated:         time.Now(),
	}
}
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"atempo/internal/logger"
	"atempo/internal/registry"
)

// ConfigManager handles CI configuration operations
type ConfigManager struct {
	projectPath string
	logger      *logger.Logger
	registry    *registry.Registry
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(projectPath string, logger *logger.Logger, registry *registry.Registry) *ConfigManager {
	return &ConfigManager{
		projectPath: projectPath,
		logger:      logger,
		registry:    registry,
	}
}

// Load loads the CI configuration from the project
func (cm *ConfigManager) Load() (*registry.CIConfig, error) {
	// Get project name from path
	projectName := filepath.Base(cm.projectPath)
	
	// Try to find project in registry
	for _, project := range cm.registry.Projects {
		if project.Path == cm.projectPath || project.Name == projectName {
			return project.CIConfig, nil
		}
	}
	
	return nil, fmt.Errorf("CI configuration not found for project at %s", cm.projectPath)
}

// Save saves the CI configuration to the registry
func (cm *ConfigManager) Save(config *registry.CIConfig) error {
	// Get project name from path
	projectName := filepath.Base(cm.projectPath)
	
	// Try to find project in registry
	for i, project := range cm.registry.Projects {
		if project.Path == cm.projectPath || project.Name == projectName {
			cm.registry.Projects[i].CIConfig = config
			cm.registry.Projects[i].CIStatus = "enabled"
			return cm.registry.SaveRegistry()
		}
	}
	
	return fmt.Errorf("project not found in registry for path %s", cm.projectPath)
}

// Exists checks if CI configuration exists for the project
func (cm *ConfigManager) Exists() bool {
	config, err := cm.Load()
	return err == nil && config != nil
}

// GetConfigPath returns the path where CI configuration should be stored
func (cm *ConfigManager) GetConfigPath() string {
	atempoDir := filepath.Join(cm.projectPath, ".atempo")
	return filepath.Join(atempoDir, "ci.json")
}

// Remove removes the CI configuration
func (cm *ConfigManager) Remove() error {
	// Get project name from path
	projectName := filepath.Base(cm.projectPath)
	
	// Try to find project in registry and remove CI config
	for i, project := range cm.registry.Projects {
		if project.Path == cm.projectPath || project.Name == projectName {
			cm.registry.Projects[i].CIConfig = nil
			cm.registry.Projects[i].CIStatus = "disabled"
			return cm.registry.SaveRegistry()
		}
	}
	
	return fmt.Errorf("project not found in registry for path %s", cm.projectPath)
}

// Backup creates a backup of the current CI configuration
func (cm *ConfigManager) Backup() error {
	// For now, this is a placeholder
	// In the future, this could create timestamped backups
	return nil
}

// EnsureProjectInRegistry ensures the project is registered
func (cm *ConfigManager) EnsureProjectInRegistry(projectName, framework, version string) error {
	// Check if project already exists
	for _, project := range cm.registry.Projects {
		if project.Path == cm.projectPath {
			return nil // Already exists
		}
	}
	
	// Add project to registry
	return cm.registry.AddProject(projectName, cm.projectPath, framework, version)
}
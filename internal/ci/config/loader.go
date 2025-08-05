package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"atempo/internal/registry"
	"atempo/internal/utils"
)

// LoadFromFile loads CI configuration from a file
func LoadFromFile(path string) (*registry.CIConfig, error) {
	if !utils.FileExists(path) {
		return nil, fmt.Errorf("configuration file not found: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var config registry.CIConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	return &config, nil
}

// SaveToFile saves CI configuration to a file
func SaveToFile(config *registry.CIConfig, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize configuration: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// ValidateConfigFile validates a CI configuration file
func ValidateConfigFile(path string) error {
	config, err := LoadFromFile(path)
	if err != nil {
		return err
	}

	return ValidateConfig(config)
}

// ValidateConfig validates a CI configuration structure
func ValidateConfig(config *registry.CIConfig) error {
	if config == nil {
		return fmt.Errorf("configuration is nil")
	}

	if config.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	if config.Framework == "" {
		return fmt.Errorf("framework is required")
	}

	if config.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Validate provider
	validProviders := []string{"github", "gitlab"}
	validProvider := false
	for _, provider := range validProviders {
		if config.Provider == provider {
			validProvider = true
			break
		}
	}
	if !validProvider {
		return fmt.Errorf("invalid provider '%s', must be one of: %v", config.Provider, validProviders)
	}

	// Validate framework
	validFrameworks := []string{"laravel", "django", "express", "lambda-node"}
	validFramework := false
	for _, framework := range validFrameworks {
		if config.Framework == framework {
			validFramework = true
			break
		}
	}
	if !validFramework {
		return fmt.Errorf("invalid framework '%s', must be one of: %v", config.Framework, validFrameworks)
	}

	return nil
}
package interfaces

import "atempo/internal/types"

// ConfigRepository defines the interface for configuration management
type ConfigRepository interface {
	// LoadConfig loads the application configuration
	LoadConfig() (*types.Configuration, error)
	
	// SaveConfig saves the application configuration
	SaveConfig(config *types.Configuration) error
	
	// GetSetting retrieves a specific configuration setting
	GetSetting(key string) (interface{}, error)
	
	// SetSetting updates a specific configuration setting
	SetSetting(key string, value interface{}) error
	
	// ResetToDefaults resets the configuration to default values
	ResetToDefaults() error
}
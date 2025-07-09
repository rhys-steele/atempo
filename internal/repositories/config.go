package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"atempo/internal/interfaces"
	"atempo/internal/types"
)

// FileConfigRepository implements ConfigRepository with file-based storage
type FileConfigRepository struct {
	configPath string
	cache      interfaces.CacheRepository
}

// NewFileConfigRepository creates a new file-based configuration repository
func NewFileConfigRepository(cache interfaces.CacheRepository) interfaces.ConfigRepository {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home is not available
		homeDir = "."
	}
	
	configPath := filepath.Join(homeDir, ".atempo", "config.json")
	
	return &FileConfigRepository{
		configPath: configPath,
		cache:      cache,
	}
}

// LoadConfig loads the configuration from storage with caching
func (r *FileConfigRepository) LoadConfig() (*types.Configuration, error) {
	// Check cache first
	if r.cache != nil {
		if cachedData, found := r.cache.Get("config"); found {
			if config, ok := cachedData.(*types.Configuration); ok {
				return config, nil
			}
		}
	}

	// Ensure the .atempo directory exists
	atempoDir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(atempoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create atempo directory: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		// Create default configuration
		config := types.DefaultConfiguration()
		
		// Save default config
		if err := r.SaveConfig(config); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config types.Configuration
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Validate and apply defaults for missing fields
	r.applyDefaults(&config)

	// Cache the loaded config
	if r.cache != nil {
		r.cache.Set("config", &config, 10*time.Minute)
	}

	return &config, nil
}

// SaveConfig saves the configuration to storage and invalidates cache
func (r *FileConfigRepository) SaveConfig(config *types.Configuration) error {
	// Ensure the .atempo directory exists
	atempoDir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(atempoDir, 0755); err != nil {
		return fmt.Errorf("failed to create atempo directory: %w", err)
	}

	// Update last modified time
	config.LastUpdated = time.Now()

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Write to file with proper permissions
	if err := os.WriteFile(r.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Update cache with new data
	if r.cache != nil {
		r.cache.Set("config", config, 10*time.Minute)
	}

	return nil
}

// GetSetting retrieves a specific configuration setting
func (r *FileConfigRepository) GetSetting(key string) (interface{}, error) {
	config, err := r.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "registry_path":
		return config.RegistryPath, nil
	case "auto_scan_projects":
		return config.AutoScanProjects, nil
	case "cache_timeout":
		return config.CacheTimeout, nil
	case "docker_timeout":
		return config.DockerTimeout, nil
	case "docker_bake_support":
		return config.DockerBakeSupport, nil
	case "ai_provider":
		return config.AIProvider, nil
	case "ai_api_key":
		return config.AIApiKey, nil
	case "ai_endpoint":
		return config.AIEndpoint, nil
	case "use_colors":
		return config.UseColors, nil
	case "show_progress":
		return config.ShowProgress, nil
	case "verbose_logging":
		return config.VerboseLogging, nil
	case "default_framework":
		return config.DefaultFramework, nil
	case "framework_versions":
		return config.FrameworkVersions, nil
	case "validate_ssl":
		return config.ValidateSSL, nil
	case "log_level":
		return config.LogLevel, nil
	default:
		return nil, fmt.Errorf("unknown configuration key: %s", key)
	}
}

// SetSetting updates a specific configuration setting
func (r *FileConfigRepository) SetSetting(key string, value interface{}) error {
	config, err := r.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "registry_path":
		if v, ok := value.(string); ok {
			config.RegistryPath = v
		} else {
			return fmt.Errorf("invalid type for registry_path: expected string")
		}
	case "auto_scan_projects":
		if v, ok := value.(bool); ok {
			config.AutoScanProjects = v
		} else {
			return fmt.Errorf("invalid type for auto_scan_projects: expected bool")
		}
	case "cache_timeout":
		if v, ok := value.(int); ok {
			config.CacheTimeout = v
		} else {
			return fmt.Errorf("invalid type for cache_timeout: expected int")
		}
	case "docker_timeout":
		if v, ok := value.(int); ok {
			config.DockerTimeout = v
		} else {
			return fmt.Errorf("invalid type for docker_timeout: expected int")
		}
	case "docker_bake_support":
		if v, ok := value.(bool); ok {
			config.DockerBakeSupport = v
		} else {
			return fmt.Errorf("invalid type for docker_bake_support: expected bool")
		}
	case "ai_provider":
		if v, ok := value.(string); ok {
			config.AIProvider = v
		} else {
			return fmt.Errorf("invalid type for ai_provider: expected string")
		}
	case "ai_api_key":
		if v, ok := value.(string); ok {
			config.AIApiKey = v
		} else {
			return fmt.Errorf("invalid type for ai_api_key: expected string")
		}
	case "ai_endpoint":
		if v, ok := value.(string); ok {
			config.AIEndpoint = v
		} else {
			return fmt.Errorf("invalid type for ai_endpoint: expected string")
		}
	case "use_colors":
		if v, ok := value.(bool); ok {
			config.UseColors = v
		} else {
			return fmt.Errorf("invalid type for use_colors: expected bool")
		}
	case "show_progress":
		if v, ok := value.(bool); ok {
			config.ShowProgress = v
		} else {
			return fmt.Errorf("invalid type for show_progress: expected bool")
		}
	case "verbose_logging":
		if v, ok := value.(bool); ok {
			config.VerboseLogging = v
		} else {
			return fmt.Errorf("invalid type for verbose_logging: expected bool")
		}
	case "default_framework":
		if v, ok := value.(string); ok {
			config.DefaultFramework = v
		} else {
			return fmt.Errorf("invalid type for default_framework: expected string")
		}
	case "validate_ssl":
		if v, ok := value.(bool); ok {
			config.ValidateSSL = v
		} else {
			return fmt.Errorf("invalid type for validate_ssl: expected bool")
		}
	case "log_level":
		if v, ok := value.(string); ok {
			config.LogLevel = v
		} else {
			return fmt.Errorf("invalid type for log_level: expected string")
		}
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return r.SaveConfig(config)
}

// ResetToDefaults resets the configuration to default values
func (r *FileConfigRepository) ResetToDefaults() error {
	config := types.DefaultConfiguration()
	return r.SaveConfig(config)
}

// applyDefaults ensures all configuration fields have default values
func (r *FileConfigRepository) applyDefaults(config *types.Configuration) {
	defaults := types.DefaultConfiguration()
	
	if config.RegistryPath == "" {
		config.RegistryPath = defaults.RegistryPath
	}
	if config.CacheTimeout == 0 {
		config.CacheTimeout = defaults.CacheTimeout
	}
	if config.DockerTimeout == 0 {
		config.DockerTimeout = defaults.DockerTimeout
	}
	if config.AIProvider == "" {
		config.AIProvider = defaults.AIProvider
	}
	if config.DefaultFramework == "" {
		config.DefaultFramework = defaults.DefaultFramework
	}
	if config.FrameworkVersions == nil {
		config.FrameworkVersions = defaults.FrameworkVersions
	}
	if config.LogLevel == "" {
		config.LogLevel = defaults.LogLevel
	}
}
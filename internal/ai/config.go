package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Config represents the AI configuration
type Config struct {
	Enabled         bool                      `json:"enabled"`
	CurrentProvider string                    `json:"current_provider"`
	Providers       map[string]ProviderConfig `json:"providers"`
	Preferences     Preferences               `json:"preferences"`
	Version         string                    `json:"version"`
}

// ProviderConfig represents configuration for an AI provider
type ProviderConfig struct {
	Name          string `json:"name"`
	Authenticated bool   `json:"authenticated"`
	DefaultModel  string `json:"default_model"`
	BaseURL       string `json:"base_url,omitempty"`
}

// Preferences represents user preferences for AI features
type Preferences struct {
	AutoGenerate     bool    `json:"auto_generate"`
	DefaultModel     string  `json:"default_model"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	ContextLength    int     `json:"context_length"`
	EnableAnalytics  bool    `json:"enable_analytics"`
	VerboseLogging   bool    `json:"verbose_logging"`
}

// ProjectInfo represents information about a detected project
type ProjectInfo struct {
	Name      string `json:"name"`
	Framework string `json:"framework"`
	Language  string `json:"language"`
	Version   string `json:"version"`
	Path      string `json:"path"`
}

// GetConfigPath returns the path to the AI configuration file
func GetConfigPath() string {
	configDir := getConfigDir()
	return filepath.Join(configDir, "ai-config.json")
}

// getConfigDir returns the configuration directory
func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return ".atempo"
	}
	return filepath.Join(homeDir, ".atempo")
}

// LoadConfig loads the AI configuration from disk
func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()
	
	// Create default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := getDefaultConfig()
		if err := SaveConfig(config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Ensure providers map is initialized
	if config.Providers == nil {
		config.Providers = make(map[string]ProviderConfig)
	}

	return &config, nil
}

// SaveConfig saves the AI configuration to disk
func SaveConfig(config *Config) error {
	configPath := GetConfigPath()
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ResetConfig resets the AI configuration to defaults
func ResetConfig() error {
	config := getDefaultConfig()
	return SaveConfig(config)
}

// getDefaultConfig returns the default AI configuration
func getDefaultConfig() *Config {
	return &Config{
		Enabled:         false,
		CurrentProvider: "",
		Providers:       make(map[string]ProviderConfig),
		Preferences: Preferences{
			AutoGenerate:     true,
			DefaultModel:     "",
			Temperature:      0.7,
			MaxTokens:        4096,
			ContextLength:    8192,
			EnableAnalytics:  false,
			VerboseLogging:   false,
		},
		Version: "1.0.0",
	}
}

// IsAIEnabled checks if AI features are enabled and authenticated
func IsAIEnabled() bool {
	config, err := LoadConfig()
	if err != nil {
		return false
	}

	if !config.Enabled {
		return false
	}

	if config.CurrentProvider == "" {
		return false
	}

	provider, exists := config.Providers[config.CurrentProvider]
	if !exists || !provider.Authenticated {
		return false
	}

	// Check if credentials exist
	_, err = GetCredential(config.CurrentProvider)
	return err == nil
}

// GetCurrentProvider returns the current AI provider configuration
func GetCurrentProvider() (*ProviderConfig, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if config.CurrentProvider == "" {
		return nil, fmt.Errorf("no AI provider configured")
	}

	provider, exists := config.Providers[config.CurrentProvider]
	if !exists {
		return nil, fmt.Errorf("current provider %s not found", config.CurrentProvider)
	}

	return &provider, nil
}

// DetectProject detects if the current directory is an Atempo project
func DetectProject(projectDir string) (*ProjectInfo, error) {
	// Check for atempo.json file
	atempoFile := filepath.Join(projectDir, "atempo.json")
	if _, err := os.Stat(atempoFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("not an Atempo project (atempo.json not found)")
	}

	// Read atempo.json to get project info
	data, err := os.ReadFile(atempoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read atempo.json: %w", err)
	}

	var atempoConfig map[string]interface{}
	if err := json.Unmarshal(data, &atempoConfig); err != nil {
		return nil, fmt.Errorf("failed to parse atempo.json: %w", err)
	}

	// Extract project information
	projectInfo := &ProjectInfo{
		Path: projectDir,
	}

	if name, ok := atempoConfig["name"].(string); ok {
		projectInfo.Name = name
	} else {
		projectInfo.Name = filepath.Base(projectDir)
	}

	if framework, ok := atempoConfig["framework"].(string); ok {
		projectInfo.Framework = framework
	}

	if language, ok := atempoConfig["language"].(string); ok {
		projectInfo.Language = language
	}

	// Try to detect version from various sources
	projectInfo.Version = detectProjectVersion(projectDir, projectInfo.Framework)

	return projectInfo, nil
}

// detectProjectVersion tries to detect the project version
func detectProjectVersion(projectDir, framework string) string {
	switch framework {
	case "laravel":
		return detectLaravelVersion(projectDir)
	case "django":
		return detectDjangoVersion(projectDir)
	default:
		return "unknown"
	}
}

// detectLaravelVersion detects Laravel version from composer.json
func detectLaravelVersion(projectDir string) string {
	composerFile := filepath.Join(projectDir, "src", "composer.json")
	if _, err := os.Stat(composerFile); os.IsNotExist(err) {
		return "unknown"
	}

	data, err := os.ReadFile(composerFile)
	if err != nil {
		return "unknown"
	}

	var composer map[string]interface{}
	if err := json.Unmarshal(data, &composer); err != nil {
		return "unknown"
	}

	if require, ok := composer["require"].(map[string]interface{}); ok {
		if laravel, ok := require["laravel/framework"].(string); ok {
			return laravel
		}
	}

	return "unknown"
}

// detectDjangoVersion detects Django version from requirements.txt
func detectDjangoVersion(projectDir string) string {
	requirementsFile := filepath.Join(projectDir, "src", "requirements.txt")
	if _, err := os.Stat(requirementsFile); os.IsNotExist(err) {
		return "unknown"
	}

	data, err := os.ReadFile(requirementsFile)
	if err != nil {
		return "unknown"
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), "django") {
			return strings.TrimSpace(line)
		}
	}

	return "unknown"
}

// GetCredentialsPath returns the path to the credentials directory
func GetCredentialsPath() string {
	configDir := getConfigDir()
	return filepath.Join(configDir, "credentials")
}

// SaveCredential saves an API key securely
func SaveCredential(provider, apiKey string) error {
	credentialsDir := GetCredentialsPath()
	
	// Create credentials directory if it doesn't exist
	if err := os.MkdirAll(credentialsDir, 0700); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	// Use a simple file-based storage for now
	// In production, this should use the OS keychain/credential store
	credentialFile := filepath.Join(credentialsDir, fmt.Sprintf("%s.key", provider))
	
	// Write API key to file with restricted permissions
	if err := os.WriteFile(credentialFile, []byte(apiKey), 0600); err != nil {
		return fmt.Errorf("failed to save credential: %w", err)
	}

	return nil
}

// GetCredential retrieves an API key securely
func GetCredential(provider string) (string, error) {
	credentialsDir := GetCredentialsPath()
	credentialFile := filepath.Join(credentialsDir, fmt.Sprintf("%s.key", provider))
	
	// Check if file exists
	if _, err := os.Stat(credentialFile); os.IsNotExist(err) {
		return "", fmt.Errorf("no credentials found for provider %s", provider)
	}

	// Read API key from file
	data, err := os.ReadFile(credentialFile)
	if err != nil {
		return "", fmt.Errorf("failed to read credential: %w", err)
	}

	return string(data), nil
}

// RemoveCredential removes stored credentials for a provider
func RemoveCredential(provider string) error {
	credentialsDir := GetCredentialsPath()
	credentialFile := filepath.Join(credentialsDir, fmt.Sprintf("%s.key", provider))
	
	// Remove file if it exists
	if _, err := os.Stat(credentialFile); err == nil {
		if err := os.Remove(credentialFile); err != nil {
			return fmt.Errorf("failed to remove credential: %w", err)
		}
	}

	return nil
}

// GetSystemInfo returns system information for debugging
func GetSystemInfo() map[string]string {
	return map[string]string{
		"os":           runtime.GOOS,
		"arch":         runtime.GOARCH,
		"go_version":   runtime.Version(),
		"config_path":  GetConfigPath(),
		"creds_path":   GetCredentialsPath(),
	}
}
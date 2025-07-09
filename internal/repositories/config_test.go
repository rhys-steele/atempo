package repositories

import (
	"os"
	"path/filepath"
	"testing"

	"atempo/internal/types"
)

func TestFileConfigRepository_LoadConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository with custom path
	cache := NewMemoryCacheRepository()
	repo := &FileConfigRepository{
		configPath: filepath.Join(tempDir, "config.json"),
		cache:      cache,
	}

	// Test loading non-existent config (should create default one)
	config, err := repo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config == nil {
		t.Fatal("Config should not be nil")
	}

	// Check default values
	if config.DefaultFramework != "laravel" {
		t.Errorf("Expected default framework 'laravel', got '%s'", config.DefaultFramework)
	}

	if config.DockerTimeout != 300 {
		t.Errorf("Expected docker timeout 300, got %d", config.DockerTimeout)
	}
}

func TestFileConfigRepository_SaveConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileConfigRepository{
		configPath: filepath.Join(tempDir, "config.json"),
		cache:      cache,
	}

	// Create test config
	config := types.DefaultConfiguration()
	config.DefaultFramework = "django"
	config.DockerTimeout = 600

	// Save config
	err = repo.SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Load config and verify changes
	loadedConfig, err := repo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConfig.DefaultFramework != "django" {
		t.Errorf("Expected framework 'django', got '%s'", loadedConfig.DefaultFramework)
	}

	if loadedConfig.DockerTimeout != 600 {
		t.Errorf("Expected docker timeout 600, got %d", loadedConfig.DockerTimeout)
	}
}

func TestFileConfigRepository_GetSetting(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileConfigRepository{
		configPath: filepath.Join(tempDir, "config.json"),
		cache:      cache,
	}

	// Test getting various settings
	testCases := []struct {
		key           string
		expectedType  string
		expectedValue interface{}
	}{
		{"default_framework", "string", "laravel"},
		{"docker_timeout", "int", 300},
		{"use_colors", "bool", true},
		{"auto_scan_projects", "bool", false},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			value, err := repo.GetSetting(tc.key)
			if err != nil {
				t.Fatalf("GetSetting failed for key '%s': %v", tc.key, err)
			}

			if value != tc.expectedValue {
				t.Errorf("Expected value %v for key '%s', got %v", tc.expectedValue, tc.key, value)
			}
		})
	}

	// Test getting non-existent setting
	_, err = repo.GetSetting("non-existent-key")
	if err == nil {
		t.Error("Expected error when getting non-existent setting")
	}
}

func TestFileConfigRepository_SetSetting(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileConfigRepository{
		configPath: filepath.Join(tempDir, "config.json"),
		cache:      cache,
	}

	// Test setting various settings
	testCases := []struct {
		key      string
		value    interface{}
		getKey   string
	}{
		{"default_framework", "django", "default_framework"},
		{"docker_timeout", 600, "docker_timeout"},
		{"use_colors", false, "use_colors"},
		{"auto_scan_projects", true, "auto_scan_projects"},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			// Set the setting
			err := repo.SetSetting(tc.key, tc.value)
			if err != nil {
				t.Fatalf("SetSetting failed for key '%s': %v", tc.key, err)
			}

			// Get and verify the setting
			value, err := repo.GetSetting(tc.getKey)
			if err != nil {
				t.Fatalf("GetSetting failed for key '%s': %v", tc.getKey, err)
			}

			if value != tc.value {
				t.Errorf("Expected value %v for key '%s', got %v", tc.value, tc.getKey, value)
			}
		})
	}

	// Test setting with wrong type
	err = repo.SetSetting("docker_timeout", "not-an-int")
	if err == nil {
		t.Error("Expected error when setting wrong type")
	}

	// Test setting non-existent key
	err = repo.SetSetting("non-existent-key", "value")
	if err == nil {
		t.Error("Expected error when setting non-existent key")
	}
}

func TestFileConfigRepository_ResetToDefaults(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileConfigRepository{
		configPath: filepath.Join(tempDir, "config.json"),
		cache:      cache,
	}

	// Change some settings
	err = repo.SetSetting("default_framework", "django")
	if err != nil {
		t.Fatalf("SetSetting failed: %v", err)
	}

	err = repo.SetSetting("docker_timeout", 600)
	if err != nil {
		t.Fatalf("SetSetting failed: %v", err)
	}

	// Reset to defaults
	err = repo.ResetToDefaults()
	if err != nil {
		t.Fatalf("ResetToDefaults failed: %v", err)
	}

	// Verify settings are back to defaults
	framework, err := repo.GetSetting("default_framework")
	if err != nil {
		t.Fatalf("GetSetting failed: %v", err)
	}

	if framework != "laravel" {
		t.Errorf("Expected framework 'laravel' after reset, got '%v'", framework)
	}

	timeout, err := repo.GetSetting("docker_timeout")
	if err != nil {
		t.Fatalf("GetSetting failed: %v", err)
	}

	if timeout != 300 {
		t.Errorf("Expected timeout 300 after reset, got %v", timeout)
	}
}

func TestFileConfigRepository_Caching(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository with cache
	cache := NewMemoryCacheRepository()
	repo := &FileConfigRepository{
		configPath: filepath.Join(tempDir, "config.json"),
		cache:      cache,
	}

	// Load config (should cache it)
	config1, err := repo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Load config again (should come from cache)
	config2, err := repo.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Both should have the same values
	if config1.DefaultFramework != config2.DefaultFramework {
		t.Error("Cache is not working correctly")
	}

	// Verify cache has the config
	if !cache.Has("config") {
		t.Error("Config should be cached")
	}
}

func TestDefaultConfiguration(t *testing.T) {
	config := types.DefaultConfiguration()

	// Test default values
	if config.DefaultFramework != "laravel" {
		t.Errorf("Expected default framework 'laravel', got '%s'", config.DefaultFramework)
	}

	if config.DockerTimeout != 300 {
		t.Errorf("Expected docker timeout 300, got %d", config.DockerTimeout)
	}

	if !config.UseColors {
		t.Error("Expected use_colors to be true by default")
	}

	if !config.ShowProgress {
		t.Error("Expected show_progress to be true by default")
	}

	if config.VerboseLogging {
		t.Error("Expected verbose_logging to be false by default")
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected log level 'info', got '%s'", config.LogLevel)
	}

	// Test framework versions
	if len(config.FrameworkVersions) == 0 {
		t.Error("Expected framework versions to be set")
	}

	if config.FrameworkVersions["laravel"] != "11" {
		t.Errorf("Expected Laravel version '11', got '%s'", config.FrameworkVersions["laravel"])
	}

	if config.FrameworkVersions["django"] != "5.0" {
		t.Errorf("Expected Django version '5.0', got '%s'", config.FrameworkVersions["django"])
	}
}
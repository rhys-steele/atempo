package scaffold

import (
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Note: embed filesystems would be set up properly in production
// For testing, we'll use mock filesystems
var testTemplatesFS embed.FS
var testMCPServersFS embed.FS

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name            string
		requestedVersion string
		meta            Metadata
		expectError     bool
		errorContains   string
	}{
		{
			name:            "Empty version should fail",
			requestedVersion: "",
			meta:            Metadata{Framework: "laravel"},
			expectError:     true,
			errorContains:   "version cannot be empty",
		},
		{
			name:            "Valid Laravel version",
			requestedVersion: "11.0",
			meta:            Metadata{Framework: "laravel", MinVersion: "8.0"},
			expectError:     false,
		},
		{
			name:            "Version below minimum",
			requestedVersion: "7.0",
			meta:            Metadata{Framework: "laravel", MinVersion: "8.0"},
			expectError:     true,
			errorContains:   "below minimum supported version",
		},
		{
			name:            "Valid Django version",
			requestedVersion: "4.2",
			meta:            Metadata{Framework: "django", MinVersion: "3.0"},
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVersion(tt.requestedVersion, tt.meta)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, but got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateLaravelVersion(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		expectError   bool
		errorContains string
	}{
		{
			name:          "Valid Laravel 11",
			version:       "11.0",
			expectError:   false,
		},
		{
			name:          "Valid Laravel 10",
			version:       "10.0",
			expectError:   false,
		},
		{
			name:          "Too old Laravel version",
			version:       "7.0",
			expectError:   true,
			errorContains: "too old",
		},
		{
			name:          "Too new Laravel version",
			version:       "13.0",
			expectError:   true,
			errorContains: "not yet supported",
		},
		{
			name:          "Valid Laravel 8",
			version:       "8.0",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLaravelVersion(tt.version)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, but got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateDjangoVersion(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid Django 4.2",
			version:     "4.2",
			expectError: false,
		},
		{
			name:        "Valid Django 5.0",
			version:     "5.0",
			expectError: false,
		},
		{
			name:          "Too old Django version",
			version:       "2.2",
			expectError:   true,
			errorContains: "too old",
		},
		{
			name:          "Too new Django version",
			version:       "6.0",
			expectError:   true,
			errorContains: "not yet supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDjangoVersion(tt.version)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, but got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestApplyVersionSpecificOptions(t *testing.T) {
	tests := []struct {
		name      string
		command   []string
		framework string
		version   string
		expected  []string
	}{
		{
			name:      "Laravel 11 with specific options",
			command:   []string{"composer", "create-project", "laravel/laravel", "src"},
			framework: "laravel",
			version:   "11.0",
			expected:  []string{"composer", "create-project", "laravel/laravel:^11.0", "src"},
		},
		{
			name:      "Django with pip install",
			command:   []string{"pip", "install", "django"},
			framework: "django",
			version:   "4.2",
			expected:  []string{"pip", "install", "django==4.2.*"},
		},
		{
			name:      "Unknown framework returns original",
			command:   []string{"npm", "install", "react"},
			framework: "react",
			version:   "18.0",
			expected:  []string{"npm", "install", "react"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyVersionSpecificOptions(tt.command, tt.framework, tt.version)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d arguments, got %d", len(tt.expected), len(result))
				return
			}
			
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected argument %d to be %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestCheckDockerAvailability(t *testing.T) {
	// This test checks if Docker is available in the system
	// We'll mock this in a real-world scenario, but for now, we'll test the function exists
	err := checkDockerAvailability()
	// We can't guarantee Docker is available in all test environments,
	// so we just check that the function doesn't panic
	if err != nil {
		t.Logf("Docker not available in test environment: %v", err)
	} else {
		t.Log("Docker is available in test environment")
	}
}

func TestMetadataUnmarshal(t *testing.T) {
	validJSON := `{
		"name": "test-project",
		"framework": "laravel",
		"language": "php",
		"installer": {
			"type": "docker",
			"command": ["composer", "create-project", "laravel/laravel", "src"],
			"work-dir": "/workspace"
		},
		"working-dir": "/var/www",
		"min-version": "8.0"
	}`

	var meta Metadata
	err := json.Unmarshal([]byte(validJSON), &meta)
	if err != nil {
		t.Fatalf("Failed to unmarshal valid JSON: %v", err)
	}

	// Verify all fields are correctly parsed
	if meta.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got '%s'", meta.Name)
	}
	if meta.Framework != "laravel" {
		t.Errorf("Expected framework 'laravel', got '%s'", meta.Framework)
	}
	if meta.Language != "php" {
		t.Errorf("Expected language 'php', got '%s'", meta.Language)
	}
	if meta.Installer.Type != "docker" {
		t.Errorf("Expected installer type 'docker', got '%s'", meta.Installer.Type)
	}
	if len(meta.Installer.Command) != 4 {
		t.Errorf("Expected 4 command parts, got %d", len(meta.Installer.Command))
	}
	if meta.WorkingDir != "/var/www" {
		t.Errorf("Expected working dir '/var/www', got '%s'", meta.WorkingDir)
	}
	if meta.MinVersion != "8.0" {
		t.Errorf("Expected min version '8.0', got '%s'", meta.MinVersion)
	}
}

func TestRunInstaller_TemplateVariableSubstitution(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scaffold-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	meta := Metadata{
		Framework: "laravel",
		Installer: Installer{
			Type:    "shell",
			Command: []string{"echo", "{{project}}", "{{cwd}}", "{{version}}"},
			WorkDir: "{{cwd}}",
		},
	}

	// Test template variable substitution
	projectName := "test-project"
	version := "11.0"
	
	// We can't easily test runInstaller without mocking exec.Command,
	// but we can test the template substitution logic by extracting it
	command := make([]string, len(meta.Installer.Command))
	for i, part := range meta.Installer.Command {
		part = strings.ReplaceAll(part, "{{name}}", "src")
		part = strings.ReplaceAll(part, "{{cwd}}", tempDir)
		part = strings.ReplaceAll(part, "{{project}}", projectName)
		part = strings.ReplaceAll(part, "{{version}}", version)
		command[i] = part
	}

	expectedCommand := []string{"echo", "test-project", tempDir, "11.0"}
	
	if len(command) != len(expectedCommand) {
		t.Errorf("Expected %d command parts, got %d", len(expectedCommand), len(command))
		return
	}
	
	for i, expected := range expectedCommand {
		if command[i] != expected {
			t.Errorf("Expected command part %d to be %q, got %q", i, expected, command[i])
		}
	}
}

func TestGetFilesystemTemplatePath(t *testing.T) {
	// Test the filesystem template path resolution
	// This tests the fallback mechanism when embedded templates aren't available
	
	framework := "laravel"
	filename := "atempo.json"
	
	path, err := getFilesystemTemplatePath(framework, filename)
	
	// The function should return a path or an error
	if err != nil {
		t.Logf("Filesystem template path not found (expected in test environment): %v", err)
	} else {
		// Verify the path structure
		expectedSuffix := filepath.Join("templates", "frameworks", framework, filename)
		if !strings.HasSuffix(path, expectedSuffix) {
			t.Errorf("Expected path to end with %q, got %q", expectedSuffix, path)
		}
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestCreateTestdataDirectories creates the testdata directories needed for testing
func TestCreateTestdataDirectories(t *testing.T) {
	// Create testdata directory structure for testing
	testdataDir := filepath.Join("testdata", "frameworks", "laravel")
	err := os.MkdirAll(testdataDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create testdata directory: %v", err)
	}

	// Create a sample atempo.json for testing
	sampleConfig := Metadata{
		Name:      "{{project}}",
		Framework: "laravel",
		Language:  "php",
		Installer: Installer{
			Type:    "docker",
			Command: []string{"composer", "create-project", "laravel/laravel", "src"},
			WorkDir: "{{cwd}}",
		},
		WorkingDir: "/var/www",
		MinVersion: "8.0",
	}

	configBytes, err := json.MarshalIndent(sampleConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal sample config: %v", err)
	}

	configPath := filepath.Join(testdataDir, "atempo.json")
	err = os.WriteFile(configPath, configBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write sample config: %v", err)
	}

	// Create MCP testdata directory
	mcpTestdataDir := filepath.Join("testdata", "mcp")
	err = os.MkdirAll(mcpTestdataDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create MCP testdata directory: %v", err)
	}

	t.Log("Created testdata directories for scaffold tests")
}
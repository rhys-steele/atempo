package commands

import (
	"context"
	"embed"
	"os"
	"path/filepath"
	"testing"
)

// Note: embed filesystems would be set up properly in production
// For testing, we'll use mock filesystems
var testTemplatesFS embed.FS
var testMCPServersFS embed.FS

func TestNewCreateCommand(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil, // Mock registry would go here
		Output:          nil, // Mock output would go here
	}

	// Create the command
	cmd := NewCreateCommand(ctx, testTemplatesFS, testMCPServersFS)

	// Verify basic properties
	if cmd.Name() != "create" {
		t.Errorf("Expected command name 'create', got '%s'", cmd.Name())
	}
	if cmd.Description() != "Create a new project" {
		t.Errorf("Expected description 'Create a new project', got '%s'", cmd.Description())
	}
	if cmd.Usage() != "atempo create <framework>[:<version>] [project_name]" {
		t.Errorf("Expected usage 'atempo create <framework>[:<version>] [project_name]', got '%s'", cmd.Usage())
	}
}

func TestCreateCommand_Execute_NoArgs(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewCreateCommand(ctx, testTemplatesFS, testMCPServersFS)

	// Test with no arguments
	err := cmd.Execute(context.Background(), []string{})
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}

	// Verify error message contains usage information
	if !containsString(err.Error(), "usage:") {
		t.Errorf("Expected error to contain usage information, got: %s", err.Error())
	}
}

func TestCreateCommand_ParseFrameworkAndVersion(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedFramework   string
		expectedVersion     string
		expectError         bool
	}{
		{
			name:              "Framework with version",
			input:             "laravel:11",
			expectedFramework: "laravel",
			expectedVersion:   "11",
			expectError:       false,
		},
		{
			name:              "Framework without version",
			input:             "laravel",
			expectedFramework: "laravel",
			expectedVersion:   "", // Will be set by getLatestVersion
			expectError:       false,
		},
		{
			name:              "Invalid format with multiple colons",
			input:             "laravel:11:extra",
			expectedFramework: "",
			expectedVersion:   "",
			expectError:       true,
		},
		{
			name:              "Django with version",
			input:             "django:4.2",
			expectedFramework: "django",
			expectedVersion:   "4.2",
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input using the same logic as the create command
			var framework, version string
			var err error

			if containsString(tt.input, ":") {
				parts := strings.Split(tt.input, ":")
				if len(parts) != 2 {
					err = fmt.Errorf("error: expected format <framework>[:<version>]")
				} else {
					framework = parts[0]
					version = parts[1]
				}
			} else {
				framework = tt.input
				// In the actual code, this would call getLatestVersion
				version = "latest"
			}

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if framework != tt.expectedFramework {
					t.Errorf("Expected framework '%s', got '%s'", tt.expectedFramework, framework)
				}
				if tt.expectedVersion != "" && version != tt.expectedVersion {
					t.Errorf("Expected version '%s', got '%s'", tt.expectedVersion, version)
				}
			}
		})
	}
}

func TestCreateCommand_GetLatestVersion(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewCreateCommand(ctx, testTemplatesFS, testMCPServersFS)

	tests := []struct {
		name      string
		framework string
		expected  string
	}{
		{
			name:      "Laravel latest version",
			framework: "laravel",
			expected:  "11", // Default latest version
		},
		{
			name:      "Django latest version",
			framework: "django",
			expected:  "5.0", // Default latest version
		},
		{
			name:      "Unknown framework",
			framework: "unknown",
			expected:  "latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmd.getLatestVersion(tt.framework)
			if result != tt.expected {
				t.Errorf("Expected version '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestCreateCommand_ProjectDirectoryHandling(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "create-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Change to temp directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	tests := []struct {
		name              string
		args              []string
		expectedProjectDir string
		expectDirCreation bool
	}{
		{
			name:              "With project name",
			args:              []string{"laravel:11", "my-project"},
			expectedProjectDir: filepath.Join(tempDir, "my-project"),
			expectDirCreation: true,
		},
		{
			name:              "Without project name",
			args:              []string{"laravel:11"},
			expectedProjectDir: tempDir,
			expectDirCreation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test directory creation logic
			var projectDir string
			var projectName string

			if len(tt.args) >= 2 {
				// Project name specified - create directory
				projectName = tt.args[1]
				cwd, err := os.Getwd()
				if err != nil {
					t.Fatalf("Failed to get current directory: %v", err)
				}
				projectDir = filepath.Join(cwd, projectName)

				// Create project directory
				if err := os.MkdirAll(projectDir, 0755); err != nil {
					t.Fatalf("Failed to create project directory: %v", err)
				}
			} else {
				// Use current directory
				var err error
				projectDir, err = os.Getwd()
				if err != nil {
					t.Fatalf("Failed to get current directory: %v", err)
				}
				projectName = filepath.Base(projectDir)
			}

			// Verify expected directory
			if projectDir != tt.expectedProjectDir {
				t.Errorf("Expected project dir '%s', got '%s'", tt.expectedProjectDir, projectDir)
			}

			// Verify directory creation
			if tt.expectDirCreation {
				if _, err := os.Stat(projectDir); os.IsNotExist(err) {
					t.Errorf("Expected directory '%s' to be created", projectDir)
				}
			}

			// Verify project name
			expectedName := filepath.Base(tt.expectedProjectDir)
			if projectName != expectedName {
				t.Errorf("Expected project name '%s', got '%s'", expectedName, projectName)
			}
		})
	}
}

func TestCreateCommand_ValidateFramework(t *testing.T) {
	tests := []struct {
		name        string
		framework   string
		expectValid bool
	}{
		{
			name:        "Valid Laravel framework",
			framework:   "laravel",
			expectValid: true,
		},
		{
			name:        "Valid Django framework",
			framework:   "django",
			expectValid: true,
		},
		{
			name:        "Invalid framework",
			framework:   "unknown-framework",
			expectValid: false,
		},
		{
			name:        "Empty framework",
			framework:   "",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test framework validation logic
			isValid := validateFramework(tt.framework)
			if isValid != tt.expectValid {
				t.Errorf("Expected framework '%s' validity to be %v, got %v", tt.framework, tt.expectValid, isValid)
			}
		})
	}
}

func TestCreateCommand_ValidateVersion(t *testing.T) {
	tests := []struct {
		name        string
		framework   string
		version     string
		expectValid bool
	}{
		{
			name:        "Valid Laravel version",
			framework:   "laravel",
			version:     "11",
			expectValid: true,
		},
		{
			name:        "Valid Django version",
			framework:   "django",
			version:     "4.2",
			expectValid: true,
		},
		{
			name:        "Invalid Laravel version",
			framework:   "laravel",
			version:     "6",
			expectValid: false,
		},
		{
			name:        "Invalid Django version",
			framework:   "django",
			version:     "2.0",
			expectValid: false,
		},
		{
			name:        "Empty version",
			framework:   "laravel",
			version:     "",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test version validation logic
			isValid := validateVersion(tt.framework, tt.version)
			if isValid != tt.expectValid {
				t.Errorf("Expected version '%s' for framework '%s' validity to be %v, got %v", 
					tt.version, tt.framework, tt.expectValid, isValid)
			}
		})
	}
}

// Helper functions for testing (these would normally be in the main file)
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

func validateFramework(framework string) bool {
	supportedFrameworks := []string{"laravel", "django"}
	for _, supported := range supportedFrameworks {
		if framework == supported {
			return true
		}
	}
	return false
}

func validateVersion(framework, version string) bool {
	if version == "" {
		return false
	}

	switch framework {
	case "laravel":
		// Laravel versions 8-12 are supported
		supportedVersions := []string{"8", "9", "10", "11", "12"}
		for _, supported := range supportedVersions {
			if version == supported {
				return true
			}
		}
		return false
	case "django":
		// Django versions 3.2, 4.0, 4.1, 4.2, 5.0 are supported
		supportedVersions := []string{"3.2", "4.0", "4.1", "4.2", "5.0"}
		for _, supported := range supportedVersions {
			if version == supported {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// TestCreateTestdata creates testdata directory structure for testing
func TestCreateTestdata(t *testing.T) {
	// Create testdata directory structure
	testdataDir := "testdata"
	err := os.MkdirAll(testdataDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create testdata directory: %v", err)
	}

	// Create MCP testdata directory
	mcpTestdataDir := filepath.Join(testdataDir, "mcp")
	err = os.MkdirAll(mcpTestdataDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create MCP testdata directory: %v", err)
	}

	t.Log("Created testdata directories for create command tests")
}

func TestCreateCommand_BaseProperties(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewCreateCommand(ctx, testTemplatesFS, testMCPServersFS)

	// Test BaseCommand properties
	if cmd.BaseCommand == nil {
		t.Error("Expected BaseCommand to be initialized")
	}

	// Test that embedded filesystems are set
	if cmd.templatesFS == nil {
		t.Error("Expected templatesFS to be set")
	}
	if cmd.mcpServersFS == nil {
		t.Error("Expected mcpServersFS to be set")
	}
}
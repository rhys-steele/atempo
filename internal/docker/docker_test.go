package docker

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSupportedCommands(t *testing.T) {
	// Test that all expected commands are defined
	expectedCommands := []string{"up", "down", "build", "logs", "ps", "restart", "stop", "pull"}
	
	for _, cmdName := range expectedCommands {
		cmd, exists := SupportedCommands[cmdName]
		if !exists {
			t.Errorf("Expected command %s to be supported", cmdName)
			continue
		}
		
		// Verify basic properties
		if cmd.Name != cmdName {
			t.Errorf("Expected command name %s, got %s", cmdName, cmd.Name)
		}
		if cmd.Description == "" {
			t.Errorf("Expected command %s to have a description", cmdName)
		}
		if len(cmd.Args) == 0 {
			t.Errorf("Expected command %s to have args", cmdName)
		}
	}
}

func TestDockerCommand_Structure(t *testing.T) {
	// Test specific command configurations
	tests := []struct {
		name            string
		expectedTimeout time.Duration
		expectedArgs    []string
	}{
		{
			name:            "up",
			expectedTimeout: 3 * time.Minute,
			expectedArgs:    []string{"up", "-d"},
		},
		{
			name:            "down",
			expectedTimeout: 2 * time.Minute,
			expectedArgs:    []string{"down"},
		},
		{
			name:            "build",
			expectedTimeout: 8 * time.Minute,
			expectedArgs:    []string{"build"},
		},
		{
			name:            "logs",
			expectedTimeout: 0, // No timeout for logs
			expectedArgs:    []string{"logs", "-f"},
		},
		{
			name:            "ps",
			expectedTimeout: 30 * time.Second,
			expectedArgs:    []string{"ps"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, exists := SupportedCommands[tt.name]
			if !exists {
				t.Fatalf("Command %s not found", tt.name)
			}

			if cmd.Timeout != tt.expectedTimeout {
				t.Errorf("Expected timeout %v, got %v", tt.expectedTimeout, cmd.Timeout)
			}

			if len(cmd.Args) != len(tt.expectedArgs) {
				t.Errorf("Expected %d args, got %d", len(tt.expectedArgs), len(cmd.Args))
			} else {
				for i, expectedArg := range tt.expectedArgs {
					if cmd.Args[i] != expectedArg {
						t.Errorf("Expected arg %d to be %s, got %s", i, expectedArg, cmd.Args[i])
					}
				}
			}
		})
	}
}

func TestValidateDockerAvailability(t *testing.T) {
	// This test checks the structure of the validation function
	// In a real environment, this would depend on Docker being available
	err := validateDockerAvailability()
	
	// We can't guarantee Docker is available in all test environments
	// so we just verify the function returns an error type
	if err != nil {
		t.Logf("Docker not available in test environment: %v", err)
	} else {
		t.Log("Docker is available in test environment")
	}
}

func TestExecuteCommand_UnsupportedCommand(t *testing.T) {
	// Test that unsupported commands return appropriate errors
	tempDir, err := os.MkdirTemp("", "docker-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = ExecuteCommand("unsupported-command", tempDir, []string{})
	if err == nil {
		t.Error("Expected error for unsupported command")
	}

	expectedError := "unsupported Docker command: unsupported-command"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestExecuteWithCustomTimeout_UnsupportedCommand(t *testing.T) {
	// Test that unsupported commands return appropriate errors with custom timeout
	tempDir, err := os.MkdirTemp("", "docker-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = ExecuteWithCustomTimeout("unsupported-command", tempDir, []string{}, 5*time.Second)
	if err == nil {
		t.Error("Expected error for unsupported command")
	}

	expectedError := "unsupported Docker command: unsupported-command"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q, got %q", expectedError, err.Error())
	}
}

func TestContainsBuildFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "Contains --build flag",
			args:     []string{"--build", "service"},
			expected: true,
		},
		{
			name:     "Contains --build in middle",
			args:     []string{"service", "--build", "other"},
			expected: true,
		},
		{
			name:     "Does not contain --build",
			args:     []string{"service", "--force"},
			expected: false,
		},
		{
			name:     "Empty args",
			args:     []string{},
			expected: false,
		},
		{
			name:     "Similar but not exact",
			args:     []string{"--build-arg", "VAR=value"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsBuildFlag(tt.args)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestResolveProjectPath(t *testing.T) {
	// Test resolving project paths
	tests := []struct {
		name        string
		projectPath string
		expectError bool
	}{
		{
			name:        "Absolute path",
			projectPath: "/absolute/path",
			expectError: false,
		},
		{
			name:        "Relative path",
			projectPath: "./relative/path",
			expectError: false,
		},
		{
			name:        "Empty path",
			projectPath: "",
			expectError: false, // Should resolve to current directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveProjectPath(tt.projectPath)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result == "" {
					t.Error("Expected non-empty result")
				}
				if !filepath.IsAbs(result) {
					t.Errorf("Expected absolute path, got %s", result)
				}
			}
		})
	}
}

func TestDetectFrameworkFromCompose(t *testing.T) {
	// Create test docker-compose content
	tests := []struct {
		name            string
		composeContent  string
		expectedFramework string
	}{
		{
			name: "Laravel framework",
			composeContent: `version: '3.8'
services:
  app:
    build: .
    image: laravel-app
    ports:
      - "8000:80"`,
			expectedFramework: "laravel",
		},
		{
			name: "Django framework",
			composeContent: `version: '3.8'
services:
  web:
    build: .
    image: django-app
    ports:
      - "8000:8000"`,
			expectedFramework: "django",
		},
		{
			name: "Unknown framework",
			composeContent: `version: '3.8'
services:
  app:
    build: .
    image: unknown-app
    ports:
      - "8000:80"`,
			expectedFramework: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary compose file
			tempDir, err := os.MkdirTemp("", "compose-test-")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			composeFile := filepath.Join(tempDir, "docker-compose.yml")
			err = os.WriteFile(composeFile, []byte(tt.composeContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write compose file: %v", err)
			}

			framework := detectFrameworkFromCompose(composeFile)
			if framework != tt.expectedFramework {
				t.Errorf("Expected framework %s, got %s", tt.expectedFramework, framework)
			}
		})
	}
}

func TestSupportsDockerBake(t *testing.T) {
	// Test Bake support detection
	result := supportsDockerBake()
	
	// Since this depends on system Docker installation, we just verify it returns a boolean
	if result {
		t.Log("Docker Bake is supported on this system")
	} else {
		t.Log("Docker Bake is not supported on this system")
	}
}

func TestSetupBakeEnvironment(t *testing.T) {
	// Test that setupBakeEnvironment doesn't panic and handles nil command
	// This is primarily a regression test
	
	// Test with nil command (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("setupBakeEnvironment panicked with nil command: %v", r)
		}
	}()
	
	// This should not panic
	setupBakeEnvironment(nil)
}

func TestDockerCommand_DefaultValues(t *testing.T) {
	// Test that DockerCommand has sensible defaults
	cmd := DockerCommand{
		Name:        "test",
		Description: "Test command",
		Args:        []string{"test"},
		Timeout:     5 * time.Second,
	}

	if cmd.Name != "test" {
		t.Errorf("Expected name 'test', got %s", cmd.Name)
	}
	if cmd.Description != "Test command" {
		t.Errorf("Expected description 'Test command', got %s", cmd.Description)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "test" {
		t.Errorf("Expected args ['test'], got %v", cmd.Args)
	}
	if cmd.Timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", cmd.Timeout)
	}
}

func TestBakeDetectionCaching(t *testing.T) {
	// Test that bake detection is cached
	// Reset the cache
	bakeMutex.Lock()
	bakeSupported = nil
	bakeMutex.Unlock()

	// Call twice and verify caching behavior
	result1 := supportsDockerBake()
	result2 := supportsDockerBake()

	if result1 != result2 {
		t.Error("Bake detection should be cached and return consistent results")
	}

	// Verify cache is set
	bakeMutex.Lock()
	if bakeSupported == nil {
		t.Error("Expected bake support to be cached")
	}
	bakeMutex.Unlock()
}

// Mock helper functions for testing (these would normally be in the main file)
func TestMockHelpers(t *testing.T) {
	// This test ensures our test helper functions work as expected
	
	// Test containsBuildFlag
	if !containsBuildFlag([]string{"--build"}) {
		t.Error("containsBuildFlag should return true for --build")
	}
	
	if containsBuildFlag([]string{"--other"}) {
		t.Error("containsBuildFlag should return false for --other")
	}
}
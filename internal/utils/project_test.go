package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveProjectPathFromArgs(t *testing.T) {
	// Setup test environment
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test project directory
	testProjectDir := filepath.Join(tempDir, "test-project")
	if err := os.MkdirAll(testProjectDir, 0755); err != nil {
		t.Fatalf("Failed to create test project directory: %v", err)
	}

	// Change to temp directory for testing
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedPath   string
		expectedName   string
		shouldError    bool
		errorContains  string
		setupFunc      func() error
		cleanupFunc    func() error
	}{
		{
			name:         "empty args should use current directory",
			args:         []string{},
			expectedPath: tempDir,
			expectedName: filepath.Base(tempDir),
			shouldError:  false,
		},
		{
			name:         "nil args should use current directory",
			args:         nil,
			expectedPath: tempDir,
			expectedName: filepath.Base(tempDir),
			shouldError:  false,
		},
		{
			name:         "empty first arg should use current directory",
			args:         []string{""},
			expectedPath: tempDir,
			expectedName: filepath.Base(tempDir),
			shouldError:  false,
		},
		{
			name:         "relative path should resolve to absolute",
			args:         []string{"test-project"},
			expectedPath: testProjectDir,
			expectedName: "test-project",
			shouldError:  false,
		},
		{
			name:         "absolute path should be used directly",
			args:         []string{testProjectDir},
			expectedPath: testProjectDir,
			expectedName: testProjectDir,
			shouldError:  false,
		},
		{
			name:          "non-existent directory should error",
			args:          []string{"non-existent"},
			shouldError:   true,
			errorContains: "project directory does not exist",
		},
		{
			name:          "non-existent absolute path should error",
			args:          []string{"/non/existent/path"},
			shouldError:   true,
			errorContains: "project directory does not exist",
		},
		{
			name:         "multiple args should use first arg",
			args:         []string{"test-project", "ignored"},
			expectedPath: testProjectDir,
			expectedName: "test-project",
			shouldError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				if err := tt.setupFunc(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			if tt.cleanupFunc != nil {
				defer func() {
					if err := tt.cleanupFunc(); err != nil {
						t.Errorf("Cleanup failed: %v", err)
					}
				}()
			}

			result, err := ResolveProjectPathFromArgs(tt.args)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
				return
			}

			// Convert expected path to absolute for comparison
			expectedAbsPath, err := filepath.Abs(tt.expectedPath)
			if err != nil {
				t.Fatalf("Failed to convert expected path to absolute: %v", err)
			}

			// Resolve symlinks for both paths to handle macOS /var -> /private/var
			expectedResolved, err := filepath.EvalSymlinks(expectedAbsPath)
			if err != nil {
				// If symlink resolution fails, use the original path
				expectedResolved = expectedAbsPath
				// For non-existent paths, try to resolve the parent directory and join the basename
				if dir := filepath.Dir(expectedAbsPath); dir != "." && dir != "/" {
					if resolvedDir, err := filepath.EvalSymlinks(dir); err == nil {
						expectedResolved = filepath.Join(resolvedDir, filepath.Base(expectedAbsPath))
					}
				}
			}

			resultResolved, err := filepath.EvalSymlinks(result.Path)
			if err != nil {
				// If symlink resolution fails, use the original path
				resultResolved = result.Path
			}

			if resultResolved != expectedResolved {
				t.Errorf("Expected path %s, got %s", expectedResolved, resultResolved)
			}

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}
		})
	}
}

func TestResolveProjectPath(t *testing.T) {
	// Setup test environment
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test project directory
	testProjectDir := filepath.Join(tempDir, "test-project")
	if err := os.MkdirAll(testProjectDir, 0755); err != nil {
		t.Fatalf("Failed to create test project directory: %v", err)
	}

	// Change to temp directory for testing
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	tests := []struct {
		name          string
		identifier    string
		expectedPath  string
		shouldError   bool
		errorContains string
	}{
		{
			name:         "empty identifier should use current directory",
			identifier:   "",
			expectedPath: tempDir,
			shouldError:  false,
		},
		{
			name:         "relative path should resolve to absolute",
			identifier:   "test-project",
			expectedPath: testProjectDir,
			shouldError:  false,
		},
		{
			name:         "absolute path should be used directly",
			identifier:   testProjectDir,
			expectedPath: testProjectDir,
			shouldError:  false,
		},
		{
			name:          "non-existent directory should error",
			identifier:    "non-existent",
			shouldError:   true,
			errorContains: "project directory does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolveProjectPath(tt.identifier)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Convert expected path to absolute for comparison
			expectedAbsPath, err := filepath.Abs(tt.expectedPath)
			if err != nil {
				t.Fatalf("Failed to convert expected path to absolute: %v", err)
			}

			// Resolve symlinks for both paths to handle macOS /var -> /private/var
			expectedResolved, err := filepath.EvalSymlinks(expectedAbsPath)
			if err != nil {
				// If symlink resolution fails, use the original path
				expectedResolved = expectedAbsPath
				// For non-existent paths, try to resolve the parent directory and join the basename
				if dir := filepath.Dir(expectedAbsPath); dir != "." && dir != "/" {
					if resolvedDir, err := filepath.EvalSymlinks(dir); err == nil {
						expectedResolved = filepath.Join(resolvedDir, filepath.Base(expectedAbsPath))
					}
				}
			}

			resultResolved, err := filepath.EvalSymlinks(result)
			if err != nil {
				// If symlink resolution fails (e.g., for non-existent paths), use the original path
				resultResolved = result
				// For non-existent paths, try to resolve the parent directory and join the basename
				if dir := filepath.Dir(result); dir != "." && dir != "/" {
					if resolvedDir, err := filepath.EvalSymlinks(dir); err == nil {
						resultResolved = filepath.Join(resolvedDir, filepath.Base(result))
					}
				}
			}

			if resultResolved != expectedResolved {
				t.Errorf("Expected path %s, got %s", expectedResolved, resultResolved)
			}
		})
	}
}

func TestResolveProjectPathFromArgsWithOffset(t *testing.T) {
	// Setup test environment
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test project directory
	testProjectDir := filepath.Join(tempDir, "test-project")
	if err := os.MkdirAll(testProjectDir, 0755); err != nil {
		t.Fatalf("Failed to create test project directory: %v", err)
	}

	// Change to temp directory for testing
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	tests := []struct {
		name         string
		args         []string
		offset       int
		expectedPath string
		expectedName string
		shouldError  bool
	}{
		{
			name:         "offset 0 should use first arg",
			args:         []string{"test-project"},
			offset:       0,
			expectedPath: testProjectDir,
			expectedName: "test-project",
			shouldError:  false,
		},
		{
			name:         "offset 1 should use second arg",
			args:         []string{"ignored", "test-project"},
			offset:       1,
			expectedPath: testProjectDir,
			expectedName: "test-project",
			shouldError:  false,
		},
		{
			name:         "offset beyond args should use current dir",
			args:         []string{"arg1"},
			offset:       2,
			expectedPath: tempDir,
			expectedName: filepath.Base(tempDir),
			shouldError:  false,
		},
		{
			name:         "offset with empty args should use current dir",
			args:         []string{},
			offset:       0,
			expectedPath: tempDir,
			expectedName: filepath.Base(tempDir),
			shouldError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolveProjectPathFromArgsWithOffset(tt.args, tt.offset)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
				return
			}

			// Convert expected path to absolute for comparison
			expectedAbsPath, err := filepath.Abs(tt.expectedPath)
			if err != nil {
				t.Fatalf("Failed to convert expected path to absolute: %v", err)
			}

			// Resolve symlinks for both paths to handle macOS /var -> /private/var
			expectedResolved, err := filepath.EvalSymlinks(expectedAbsPath)
			if err != nil {
				// If symlink resolution fails, use the original path
				expectedResolved = expectedAbsPath
				// For non-existent paths, try to resolve the parent directory and join the basename
				if dir := filepath.Dir(expectedAbsPath); dir != "." && dir != "/" {
					if resolvedDir, err := filepath.EvalSymlinks(dir); err == nil {
						expectedResolved = filepath.Join(resolvedDir, filepath.Base(expectedAbsPath))
					}
				}
			}

			resultResolved, err := filepath.EvalSymlinks(result.Path)
			if err != nil {
				// If symlink resolution fails, use the original path
				resultResolved = result.Path
			}

			if resultResolved != expectedResolved {
				t.Errorf("Expected path %s, got %s", expectedResolved, resultResolved)
			}

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}
		})
	}
}

func TestResolveCurrentProjectPath(t *testing.T) {
	// Setup test environment
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory for testing
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	result, err := ResolveCurrentProjectPath()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Errorf("Expected result but got nil")
		return
	}

	// Convert expected path to absolute for comparison
	expectedAbsPath, err := filepath.Abs(tempDir)
	if err != nil {
		t.Fatalf("Failed to convert expected path to absolute: %v", err)
	}

	// Resolve symlinks for both paths to handle macOS /var -> /private/var
	expectedResolved, err := filepath.EvalSymlinks(expectedAbsPath)
	if err != nil {
		// If symlink resolution fails, use the original path
		expectedResolved = expectedAbsPath
	}

	resultResolved, err := filepath.EvalSymlinks(result.Path)
	if err != nil {
		// If symlink resolution fails, use the original path
		resultResolved = result.Path
	}

	if resultResolved != expectedResolved {
		t.Errorf("Expected path %s, got %s", expectedResolved, resultResolved)
	}

	expectedName := filepath.Base(tempDir)
	if result.Name != expectedName {
		t.Errorf("Expected name %s, got %s", expectedName, result.Name)
	}
}

func TestResolveProjectIdentifier(t *testing.T) {
	// Setup test environment
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test project directory
	testProjectDir := filepath.Join(tempDir, "test-project")
	if err := os.MkdirAll(testProjectDir, 0755); err != nil {
		t.Fatalf("Failed to create test project directory: %v", err)
	}

	// Change to temp directory for testing
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	tests := []struct {
		name         string
		identifier   string
		expectedPath string
		shouldError  bool
	}{
		{
			name:         "empty identifier should use current directory",
			identifier:   "",
			expectedPath: tempDir,
			shouldError:  false,
		},
		{
			name:         "relative path should resolve to absolute",
			identifier:   "test-project",
			expectedPath: testProjectDir,
			shouldError:  false,
		},
		{
			name:         "absolute path should be used directly",
			identifier:   testProjectDir,
			expectedPath: testProjectDir,
			shouldError:  false,
		},
		{
			name:         "non-existent directory should still resolve path",
			identifier:   "non-existent",
			expectedPath: filepath.Join(tempDir, "non-existent"),
			shouldError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveProjectIdentifier(tt.identifier)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Convert expected path to absolute for comparison
			expectedAbsPath, err := filepath.Abs(tt.expectedPath)
			if err != nil {
				t.Fatalf("Failed to convert expected path to absolute: %v", err)
			}

			// Resolve symlinks for both paths to handle macOS /var -> /private/var
			expectedResolved, err := filepath.EvalSymlinks(expectedAbsPath)
			if err != nil {
				// If symlink resolution fails, use the original path
				expectedResolved = expectedAbsPath
				// For non-existent paths, try to resolve the parent directory and join the basename
				if dir := filepath.Dir(expectedAbsPath); dir != "." && dir != "/" {
					if resolvedDir, err := filepath.EvalSymlinks(dir); err == nil {
						expectedResolved = filepath.Join(resolvedDir, filepath.Base(expectedAbsPath))
					}
				}
			}

			resultResolved, err := filepath.EvalSymlinks(result)
			if err != nil {
				// If symlink resolution fails (e.g., for non-existent paths), use the original path
				resultResolved = result
				// For non-existent paths, try to resolve the parent directory and join the basename
				if dir := filepath.Dir(result); dir != "." && dir != "/" {
					if resolvedDir, err := filepath.EvalSymlinks(dir); err == nil {
						resultResolved = filepath.Join(resolvedDir, filepath.Base(result))
					}
				}
			}

			if resultResolved != expectedResolved {
				t.Errorf("Expected path %s, got %s", expectedResolved, resultResolved)
			}
		})
	}
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
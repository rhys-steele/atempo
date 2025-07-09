package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetRegistryPath(t *testing.T) {
	path, err := GetRegistryPath()
	if err != nil {
		t.Fatalf("Failed to get registry path: %v", err)
	}

	// Check that path ends with .atempo/registry.json
	if !filepath.IsAbs(path) {
		t.Errorf("Expected absolute path, got: %s", path)
	}

	expectedSuffix := filepath.Join(".atempo", "registry.json")
	if !endsWithPath(path, expectedSuffix) {
		t.Errorf("Expected path to end with %s, got: %s", expectedSuffix, path)
	}
}

func TestLoadRegistry_EmptyRegistry(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir, err := os.MkdirTemp("", "registry-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Load registry when no file exists
	registry, err := LoadRegistry()
	if err != nil {
		t.Fatalf("Failed to load empty registry: %v", err)
	}

	// Verify empty registry structure
	if len(registry.Projects) != 0 {
		t.Errorf("Expected empty projects slice, got %d projects", len(registry.Projects))
	}
	if registry.Version != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", registry.Version)
	}
}

func TestLoadRegistry_ExistingRegistry(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir, err := os.MkdirTemp("", "registry-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create .atempo directory
	atempoDir := filepath.Join(tempDir, ".atempo")
	err = os.MkdirAll(atempoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .atempo directory: %v", err)
	}

	// Create sample registry file
	sampleRegistry := Registry{
		Version: "1.0",
		Projects: []Project{
			{
				Name:         "test-project",
				Path:         "/path/to/test-project",
				Framework:    "laravel",
				Version:      "11.0",
				CreatedAt:    time.Now(),
				LastAccessed: time.Now(),
			},
		},
	}

	registryData, err := json.MarshalIndent(sampleRegistry, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal sample registry: %v", err)
	}

	registryPath := filepath.Join(atempoDir, "registry.json")
	err = os.WriteFile(registryPath, registryData, 0644)
	if err != nil {
		t.Fatalf("Failed to write sample registry: %v", err)
	}

	// Load the registry
	registry, err := LoadRegistry()
	if err != nil {
		t.Fatalf("Failed to load registry: %v", err)
	}

	// Verify loaded registry
	if len(registry.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(registry.Projects))
	}
	if registry.Projects[0].Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", registry.Projects[0].Name)
	}
	if registry.Projects[0].Framework != "laravel" {
		t.Errorf("Expected framework 'laravel', got '%s'", registry.Projects[0].Framework)
	}
	if registry.Projects[0].Version != "11.0" {
		t.Errorf("Expected version '11.0', got '%s'", registry.Projects[0].Version)
	}
}

func TestRegistry_AddProject(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir, err := os.MkdirTemp("", "registry-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create empty registry
	registry := &Registry{
		Version:  "1.0",
		Projects: []Project{},
	}

	// Add a project
	err = registry.AddProject("test-project", "/path/to/test-project", "laravel", "11.0")
	if err != nil {
		t.Fatalf("Failed to add project: %v", err)
	}

	// Verify project was added
	if len(registry.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(registry.Projects))
	}

	project := registry.Projects[0]
	if project.Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", project.Name)
	}
	if project.Framework != "laravel" {
		t.Errorf("Expected framework 'laravel', got '%s'", project.Framework)
	}
	if project.Version != "11.0" {
		t.Errorf("Expected version '11.0', got '%s'", project.Version)
	}
	if project.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if project.LastAccessed.IsZero() {
		t.Error("Expected LastAccessed to be set")
	}
}

func TestRegistry_AddProject_UpdateExisting(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir, err := os.MkdirTemp("", "registry-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create registry with existing project
	createdTime := time.Now().Add(-time.Hour)
	registry := &Registry{
		Version: "1.0",
		Projects: []Project{
			{
				Name:         "test-project",
				Path:         "/old/path",
				Framework:    "laravel",
				Version:      "10.0",
				CreatedAt:    createdTime,
				LastAccessed: createdTime,
			},
		},
	}

	// Update the project
	err = registry.AddProject("test-project", "/new/path", "laravel", "11.0")
	if err != nil {
		t.Fatalf("Failed to update project: %v", err)
	}

	// Verify project was updated, not duplicated
	if len(registry.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(registry.Projects))
	}

	project := registry.Projects[0]
	if project.Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", project.Name)
	}
	if project.Version != "11.0" {
		t.Errorf("Expected updated version '11.0', got '%s'", project.Version)
	}
	if project.CreatedAt != createdTime {
		t.Error("Expected CreatedAt to be preserved")
	}
	if !project.LastAccessed.After(createdTime) {
		t.Error("Expected LastAccessed to be updated")
	}
}

func TestRegistry_FindProject(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir, err := os.MkdirTemp("", "registry-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create registry with projects
	registry := &Registry{
		Version: "1.0",
		Projects: []Project{
			{
				Name:         "project1",
				Path:         "/path/to/project1",
				Framework:    "laravel",
				Version:      "11.0",
				CreatedAt:    time.Now(),
				LastAccessed: time.Now().Add(-time.Hour),
			},
			{
				Name:         "project2",
				Path:         "/path/to/project2",
				Framework:    "django",
				Version:      "4.2",
				CreatedAt:    time.Now(),
				LastAccessed: time.Now().Add(-time.Hour),
			},
		},
	}

	// Find existing project
	project, err := registry.FindProject("project1")
	if err != nil {
		t.Fatalf("Failed to find project: %v", err)
	}

	if project.Name != "project1" {
		t.Errorf("Expected project name 'project1', got '%s'", project.Name)
	}
	if project.Framework != "laravel" {
		t.Errorf("Expected framework 'laravel', got '%s'", project.Framework)
	}

	// Find non-existing project
	_, err = registry.FindProject("nonexistent")
	if err == nil {
		t.Error("Expected error when finding non-existing project")
	}
}

func TestRegistry_ListProjects(t *testing.T) {
	// Create registry with projects
	registry := &Registry{
		Version: "1.0",
		Projects: []Project{
			{
				Name:      "project1",
				Path:      "/path/to/project1",
				Framework: "laravel",
				Version:   "11.0",
			},
			{
				Name:      "project2",
				Path:      "/path/to/project2",
				Framework: "django",
				Version:   "4.2",
			},
		},
	}

	projects := registry.ListProjects()

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}

	// Verify projects are returned
	names := make(map[string]bool)
	for _, project := range projects {
		names[project.Name] = true
	}

	if !names["project1"] {
		t.Error("Expected project1 in list")
	}
	if !names["project2"] {
		t.Error("Expected project2 in list")
	}
}

func TestRegistry_RemoveProject(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir, err := os.MkdirTemp("", "registry-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create registry with projects
	registry := &Registry{
		Version: "1.0",
		Projects: []Project{
			{
				Name:      "project1",
				Path:      "/path/to/project1",
				Framework: "laravel",
				Version:   "11.0",
			},
			{
				Name:      "project2",
				Path:      "/path/to/project2",
				Framework: "django",
				Version:   "4.2",
			},
		},
	}

	// Remove existing project
	err = registry.RemoveProject("project1")
	if err != nil {
		t.Fatalf("Failed to remove project: %v", err)
	}

	// Verify project was removed
	if len(registry.Projects) != 1 {
		t.Errorf("Expected 1 project after removal, got %d", len(registry.Projects))
	}

	if registry.Projects[0].Name != "project2" {
		t.Errorf("Expected remaining project to be 'project2', got '%s'", registry.Projects[0].Name)
	}

	// Try to remove non-existing project
	err = registry.RemoveProject("nonexistent")
	if err == nil {
		t.Error("Expected error when removing non-existing project")
	}
}

func TestResolveProjectPath(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir, err := os.MkdirTemp("", "registry-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create registry with a project
	registry := &Registry{
		Version: "1.0",
		Projects: []Project{
			{
				Name:      "test-project",
				Path:      "/path/to/test-project",
				Framework: "laravel",
				Version:   "11.0",
			},
		},
	}

	// Save registry
	err = registry.SaveRegistry()
	if err != nil {
		t.Fatalf("Failed to save registry: %v", err)
	}

	// Test resolving by project name
	path, err := ResolveProjectPath("test-project")
	if err != nil {
		t.Fatalf("Failed to resolve project path: %v", err)
	}

	if path != "/path/to/test-project" {
		t.Errorf("Expected path '/path/to/test-project', got '%s'", path)
	}

	// Test resolving absolute path
	absPath := "/some/absolute/path"
	path, err = ResolveProjectPath(absPath)
	if err != nil {
		t.Fatalf("Failed to resolve absolute path: %v", err)
	}

	if path != absPath {
		t.Errorf("Expected path '%s', got '%s'", absPath, path)
	}

	// Test resolving empty string (should return current directory)
	path, err = ResolveProjectPath("")
	if err != nil {
		t.Fatalf("Failed to resolve empty path: %v", err)
	}

	if path == "" {
		t.Error("Expected non-empty path for current directory")
	}
}

func TestRegistry_SaveRegistry(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir, err := os.MkdirTemp("", "registry-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create registry
	registry := &Registry{
		Version: "1.0",
		Projects: []Project{
			{
				Name:      "test-project",
				Path:      "/path/to/test-project",
				Framework: "laravel",
				Version:   "11.0",
			},
		},
	}

	// Save registry
	err = registry.SaveRegistry()
	if err != nil {
		t.Fatalf("Failed to save registry: %v", err)
	}

	// Verify file was created
	registryPath := filepath.Join(tempDir, ".atempo", "registry.json")
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		t.Error("Registry file was not created")
	}

	// Verify file contents
	data, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("Failed to read registry file: %v", err)
	}

	var savedRegistry Registry
	err = json.Unmarshal(data, &savedRegistry)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved registry: %v", err)
	}

	if len(savedRegistry.Projects) != 1 {
		t.Errorf("Expected 1 project in saved registry, got %d", len(savedRegistry.Projects))
	}

	if savedRegistry.Projects[0].Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", savedRegistry.Projects[0].Name)
	}
}

// Helper function to check if path ends with expected suffix
func endsWithPath(path, suffix string) bool {
	// Normalize path separators
	normalizedPath := filepath.ToSlash(path)
	normalizedSuffix := filepath.ToSlash(suffix)
	
	return len(normalizedPath) >= len(normalizedSuffix) && 
		normalizedPath[len(normalizedPath)-len(normalizedSuffix):] == normalizedSuffix
}
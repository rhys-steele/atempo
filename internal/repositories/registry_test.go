package repositories

import (
	"os"
	"path/filepath"
	"testing"

	"atempo/internal/registry"
)

func TestFileRegistryRepository_LoadRegistry(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository with custom path
	cache := NewMemoryCacheRepository()
	repo := &FileRegistryRepository{
		registryPath: filepath.Join(tempDir, "registry.json"),
		cache:        cache,
	}

	// Test loading non-existent registry (should create empty one)
	reg, err := repo.LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}

	if reg == nil {
		t.Fatal("Registry should not be nil")
	}

	if len(reg.Projects) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(reg.Projects))
	}
}

func TestFileRegistryRepository_AddProject(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileRegistryRepository{
		registryPath: filepath.Join(tempDir, "registry.json"),
		cache:        cache,
	}

	// Create test project
	project := &registry.Project{
		Name:      "test-project",
		Path:      "/path/to/project",
		Framework: "laravel",
		Version:   "11",
		Status:    "stopped",
	}

	// Add project
	err = repo.AddProject(project)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	// Verify project was added
	reg, err := repo.LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}

	if len(reg.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(reg.Projects))
	}

	if reg.Projects[0].Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", reg.Projects[0].Name)
	}
}

func TestFileRegistryRepository_FindProject(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileRegistryRepository{
		registryPath: filepath.Join(tempDir, "registry.json"),
		cache:        cache,
	}

	// Create and add test project
	project := &registry.Project{
		Name:      "test-project",
		Path:      "/path/to/project",
		Framework: "laravel",
		Version:   "11",
		Status:    "stopped",
	}

	err = repo.AddProject(project)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	// Find existing project
	found, err := repo.FindProject("test-project")
	if err != nil {
		t.Fatalf("FindProject failed: %v", err)
	}

	if found.Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", found.Name)
	}

	// Try to find non-existent project
	_, err = repo.FindProject("non-existent")
	if err == nil {
		t.Error("Expected error when finding non-existent project")
	}
}

func TestFileRegistryRepository_RemoveProject(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileRegistryRepository{
		registryPath: filepath.Join(tempDir, "registry.json"),
		cache:        cache,
	}

	// Create and add test project
	project := &registry.Project{
		Name:      "test-project",
		Path:      "/path/to/project",
		Framework: "laravel",
		Version:   "11",
		Status:    "stopped",
	}

	err = repo.AddProject(project)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	// Remove project
	err = repo.RemoveProject("test-project")
	if err != nil {
		t.Fatalf("RemoveProject failed: %v", err)
	}

	// Verify project was removed
	reg, err := repo.LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}

	if len(reg.Projects) != 0 {
		t.Errorf("Expected 0 projects after removal, got %d", len(reg.Projects))
	}

	// Try to remove non-existent project
	err = repo.RemoveProject("non-existent")
	if err == nil {
		t.Error("Expected error when removing non-existent project")
	}
}

func TestFileRegistryRepository_ListProjects(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileRegistryRepository{
		registryPath: filepath.Join(tempDir, "registry.json"),
		cache:        cache,
	}

	// Add multiple test projects
	projects := []*registry.Project{
		{
			Name:      "project1",
			Path:      "/path/to/project1",
			Framework: "laravel",
			Version:   "11",
			Status:    "stopped",
		},
		{
			Name:      "project2",
			Path:      "/path/to/project2",
			Framework: "django",
			Version:   "5.0",
			Status:    "running",
		},
	}

	for _, project := range projects {
		err = repo.AddProject(project)
		if err != nil {
			t.Fatalf("AddProject failed: %v", err)
		}
	}

	// List projects
	projectList, err := repo.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	if len(projectList) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projectList))
	}
}

func TestFileRegistryRepository_UpdateProjectStatus(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	cache := NewMemoryCacheRepository()
	repo := &FileRegistryRepository{
		registryPath: filepath.Join(tempDir, "registry.json"),
		cache:        cache,
	}

	// Create and add test project
	project := &registry.Project{
		Name:      "test-project",
		Path:      "/path/to/project",
		Framework: "laravel",
		Version:   "11",
		Status:    "stopped",
	}

	err = repo.AddProject(project)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	// Update project status
	err = repo.UpdateProjectStatus("test-project")
	if err != nil {
		t.Fatalf("UpdateProjectStatus failed: %v", err)
	}

	// Try to update non-existent project
	err = repo.UpdateProjectStatus("non-existent")
	if err == nil {
		t.Error("Expected error when updating non-existent project")
	}
}

func TestFileRegistryRepository_Caching(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository with cache
	cache := NewMemoryCacheRepository()
	repo := &FileRegistryRepository{
		registryPath: filepath.Join(tempDir, "registry.json"),
		cache:        cache,
	}

	// Load registry (should cache it)
	reg1, err := repo.LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}

	// Load registry again (should come from cache)
	reg2, err := repo.LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}

	// Both should be the same instance from cache
	if len(reg1.Projects) != len(reg2.Projects) {
		t.Error("Cache is not working correctly")
	}

	// Verify cache has the registry
	if !cache.Has("registry") {
		t.Error("Registry should be cached")
	}
}

func TestDetectFrameworkFromPath(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atempo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name           string
		files          []string
		expectedResult string
	}{
		{
			name:           "Laravel project",
			files:          []string{"composer.json", "artisan"},
			expectedResult: "laravel",
		},
		{
			name:           "Django project",
			files:          []string{"manage.py"},
			expectedResult: "django",
		},
		{
			name:           "Node.js project",
			files:          []string{"package.json"},
			expectedResult: "nodejs",
		},
		{
			name:           "Unknown project",
			files:          []string{"random.txt"},
			expectedResult: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test directory
			testDir := filepath.Join(tempDir, tc.name)
			err := os.MkdirAll(testDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test dir: %v", err)
			}

			// Create test files
			for _, file := range tc.files {
				filePath := filepath.Join(testDir, file)
				err := os.WriteFile(filePath, []byte("test"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Test framework detection
			result := detectFrameworkFromPath(testDir)
			if result != tc.expectedResult {
				t.Errorf("Expected framework '%s', got '%s'", tc.expectedResult, result)
			}

			// Clean up test directory
			os.RemoveAll(testDir)
		})
	}
}
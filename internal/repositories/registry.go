package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"atempo/internal/interfaces"
	"atempo/internal/registry"
)

// FileRegistryRepository implements RegistryRepository with file-based storage
type FileRegistryRepository struct {
	registryPath string
	cache        interfaces.CacheRepository
}

// NewFileRegistryRepository creates a new file-based registry repository
func NewFileRegistryRepository(cache interfaces.CacheRepository) interfaces.RegistryRepository {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home is not available
		homeDir = "."
	}
	
	registryPath := filepath.Join(homeDir, ".atempo", "registry.json")
	
	return &FileRegistryRepository{
		registryPath: registryPath,
		cache:        cache,
	}
}

// LoadRegistry loads the project registry from storage with caching
func (r *FileRegistryRepository) LoadRegistry() (*registry.Registry, error) {
	// Check cache first
	if r.cache != nil {
		if cachedData, found := r.cache.Get("registry"); found {
			if reg, ok := cachedData.(*registry.Registry); ok {
				return reg, nil
			}
		}
	}

	// Ensure the .atempo directory exists
	atempoDir := filepath.Dir(r.registryPath)
	if err := os.MkdirAll(atempoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create atempo directory: %w", err)
	}

	// Check if registry file exists
	if _, err := os.Stat(r.registryPath); os.IsNotExist(err) {
		// Create empty registry
		reg := &registry.Registry{
			Projects: []registry.Project{},
		}
		
		// Cache the empty registry
		if r.cache != nil {
			r.cache.Set("registry", reg, 5*time.Minute)
		}
		
		return reg, nil
	}

	// Read registry file
	data, err := os.ReadFile(r.registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry file: %w", err)
	}

	// Parse JSON
	var reg registry.Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry JSON: %w", err)
	}

	// Cache the loaded registry
	if r.cache != nil {
		r.cache.Set("registry", &reg, 5*time.Minute)
	}

	return &reg, nil
}

// SaveRegistry saves the project registry to storage and invalidates cache
func (r *FileRegistryRepository) SaveRegistry(reg *registry.Registry) error {
	// Ensure the .atempo directory exists
	atempoDir := filepath.Dir(r.registryPath)
	if err := os.MkdirAll(atempoDir, 0755); err != nil {
		return fmt.Errorf("failed to create atempo directory: %w", err)
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry to JSON: %w", err)
	}

	// Write to file with proper permissions
	if err := os.WriteFile(r.registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry file: %w", err)
	}

	// Update cache with new data
	if r.cache != nil {
		r.cache.Set("registry", reg, 5*time.Minute)
	}

	return nil
}

// AddProject adds a new project to the registry
func (r *FileRegistryRepository) AddProject(project *registry.Project) error {
	reg, err := r.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Check if project already exists
	for i, p := range reg.Projects {
		if p.Name == project.Name {
			// Update existing project
			reg.Projects[i] = *project
			return r.SaveRegistry(reg)
		}
	}

	// Add new project
	project.CreatedAt = time.Now()
	project.LastAccessed = time.Now()
	reg.Projects = append(reg.Projects, *project)

	return r.SaveRegistry(reg)
}

// FindProject finds a project by name
func (r *FileRegistryRepository) FindProject(name string) (*registry.Project, error) {
	reg, err := r.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	for i, project := range reg.Projects {
		if project.Name == name {
			// Update last accessed time
			reg.Projects[i].LastAccessed = time.Now()
			
			// Save updated registry (async to avoid blocking)
			go func() {
				_ = r.SaveRegistry(reg)
			}()
			
			return &reg.Projects[i], nil
		}
	}

	return nil, fmt.Errorf("project '%s' not found in registry", name)
}

// RemoveProject removes a project from the registry
func (r *FileRegistryRepository) RemoveProject(name string) error {
	reg, err := r.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	for i, project := range reg.Projects {
		if project.Name == name {
			// Remove project from slice
			reg.Projects = append(reg.Projects[:i], reg.Projects[i+1:]...)
			return r.SaveRegistry(reg)
		}
	}

	return fmt.Errorf("project '%s' not found in registry", name)
}

// ListProjects returns all projects in the registry
func (r *FileRegistryRepository) ListProjects() ([]registry.Project, error) {
	reg, err := r.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	return reg.Projects, nil
}

// UpdateProjectStatus updates the status of a project
func (r *FileRegistryRepository) UpdateProjectStatus(name string) error {
	reg, err := r.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	for i, project := range reg.Projects {
		if project.Name == name {
			// Update project status (this would normally involve health checking)
			reg.Projects[i].LastAccessed = time.Now()
			
			// Here we would update status, ports, URLs, etc.
			// For now, just update the last accessed time
			
			return r.SaveRegistry(reg)
		}
	}

	return fmt.Errorf("project '%s' not found in registry", name)
}

// ScanForProjects scans a directory for Atempo projects
func (r *FileRegistryRepository) ScanForProjects(scanPath string) error {
	reg, err := r.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Walk the directory looking for docker-compose.yml files or atempo.json files
	err = filepath.Walk(scanPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}

		// Look for docker-compose.yml files that might be Atempo projects
		if info.Name() == "docker-compose.yml" {
			projectDir := filepath.Dir(path)
			projectName := filepath.Base(projectDir)
			
			// Check if project already exists in registry
			found := false
			for _, project := range reg.Projects {
				if project.Path == projectDir {
					found = true
					break
				}
			}
			
			if !found {
				// Try to detect framework
				framework := detectFrameworkFromPath(projectDir)
				if framework != "" {
					newProject := &registry.Project{
						Name:         projectName,
						Path:         projectDir,
						Framework:    framework,
						Status:       "stopped",
						CreatedAt:    time.Now(),
						LastAccessed: time.Now(),
					}
					
					reg.Projects = append(reg.Projects, *newProject)
				}
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	return r.SaveRegistry(reg)
}

// detectFrameworkFromPath attempts to detect the framework used in a project
func detectFrameworkFromPath(projectPath string) string {
	// Check for Laravel indicators
	if fileExists(filepath.Join(projectPath, "composer.json")) &&
		fileExists(filepath.Join(projectPath, "artisan")) {
		return "laravel"
	}

	// Check for Django indicators
	if fileExists(filepath.Join(projectPath, "manage.py")) ||
		fileExists(filepath.Join(projectPath, "requirements.txt")) {
		return "django"
	}

	// Check for Node.js indicators
	if fileExists(filepath.Join(projectPath, "package.json")) {
		return "nodejs"
	}

	return ""
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
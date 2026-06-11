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

const registryCacheTTL = 5 * time.Minute
const registryCacheKey = "registry"

// FileRegistryRepository stores the project registry as JSON in ~/.atempo/registry.json.
type FileRegistryRepository struct {
	path  string
	cache interfaces.CacheRepository
}

// NewFileRegistryRepository creates a new file-based registry repository.
func NewFileRegistryRepository(cache interfaces.CacheRepository) interfaces.RegistryRepository {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return &FileRegistryRepository{
		path:  filepath.Join(homeDir, ".atempo", "registry.json"),
		cache: cache,
	}
}

func (r *FileRegistryRepository) LoadRegistry() (*registry.Registry, error) {
	if r.cache != nil {
		if cached, ok := r.cache.Get(registryCacheKey); ok {
			if reg, ok := cached.(*registry.Registry); ok {
				return reg, nil
			}
		}
	}

	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return nil, fmt.Errorf("create atempo dir: %w", err)
	}

	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		reg := &registry.Registry{Version: "1", Projects: []registry.Project{}}
		r.cacheSet(reg)
		return reg, nil
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("read registry: %w", err)
	}

	var reg registry.Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("parse registry: %w", err)
	}

	r.cacheSet(&reg)
	return &reg, nil
}

func (r *FileRegistryRepository) SaveRegistry(reg *registry.Registry) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return fmt.Errorf("create atempo dir: %w", err)
	}

	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal registry: %w", err)
	}

	if err := os.WriteFile(r.path, data, 0644); err != nil {
		return fmt.Errorf("write registry: %w", err)
	}

	r.cacheSet(reg)
	return nil
}

func (r *FileRegistryRepository) AddProject(project *registry.Project) error {
	reg, err := r.LoadRegistry()
	if err != nil {
		return err
	}

	now := time.Now()
	for i, p := range reg.Projects {
		if p.Name == project.Name {
			reg.Projects[i] = *project
			reg.Projects[i].LastAccessed = now
			return r.SaveRegistry(reg)
		}
	}

	project.CreatedAt = now
	project.LastAccessed = now
	reg.Projects = append(reg.Projects, *project)
	return r.SaveRegistry(reg)
}

func (r *FileRegistryRepository) FindProject(name string) (*registry.Project, error) {
	reg, err := r.LoadRegistry()
	if err != nil {
		return nil, err
	}

	for i := range reg.Projects {
		if reg.Projects[i].Name == name {
			return &reg.Projects[i], nil
		}
	}

	return nil, fmt.Errorf("project %q not found", name)
}

func (r *FileRegistryRepository) RemoveProject(name string) error {
	reg, err := r.LoadRegistry()
	if err != nil {
		return err
	}

	for i, p := range reg.Projects {
		if p.Name == name {
			reg.Projects = append(reg.Projects[:i], reg.Projects[i+1:]...)
			return r.SaveRegistry(reg)
		}
	}

	return fmt.Errorf("project %q not found", name)
}

func (r *FileRegistryRepository) ListProjects() ([]registry.Project, error) {
	reg, err := r.LoadRegistry()
	if err != nil {
		return nil, err
	}
	return reg.Projects, nil
}

func (r *FileRegistryRepository) cacheSet(reg *registry.Registry) {
	if r.cache != nil {
		r.cache.Set(registryCacheKey, reg, registryCacheTTL)
	}
}

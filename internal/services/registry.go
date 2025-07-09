package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"atempo/internal/registry"
)

// RegistryService provides business operations for registry management with caching
type RegistryService interface {
	// LoadRegistry loads the registry from disk or cache
	LoadRegistry(ctx context.Context) (*registry.Registry, error)
	
	// SaveRegistry saves the registry to disk and updates cache
	SaveRegistry(ctx context.Context, reg *registry.Registry) error
	
	// AddProject adds a project to the registry
	AddProject(ctx context.Context, project *registry.Project) error
	
	// FindProject finds a project by name
	FindProject(ctx context.Context, name string) (*registry.Project, error)
	
	// RemoveProject removes a project from the registry
	RemoveProject(ctx context.Context, name string) error
	
	// ListProjects returns all projects in the registry
	ListProjects(ctx context.Context) ([]registry.Project, error)
	
	// UpdateProjectStatus updates the status of a project
	UpdateProjectStatus(ctx context.Context, name string) error
	
	// RefreshCache forces a refresh of the cached registry
	RefreshCache(ctx context.Context) error
	
	// ClearCache clears the cached registry
	ClearCache(ctx context.Context) error
}

// registryService implements RegistryService with caching
type registryService struct {
	mu           sync.RWMutex
	cachedRegistry *registry.Registry
	cacheTime     time.Time
	cacheTTL      time.Duration
}

// NewRegistryService creates a new RegistryService implementation with caching
func NewRegistryService() RegistryService {
	return &registryService{
		cacheTTL: 5 * time.Minute, // Cache for 5 minutes
	}
}

// LoadRegistry loads the registry from disk or cache
func (s *registryService) LoadRegistry(ctx context.Context) (*registry.Registry, error) {
	s.mu.RLock()
	// Check if we have a valid cached registry
	if s.cachedRegistry != nil && time.Since(s.cacheTime) < s.cacheTTL {
		s.mu.RUnlock()
		return s.cachedRegistry, nil
	}
	s.mu.RUnlock()

	// Need to load from disk
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring write lock
	if s.cachedRegistry != nil && time.Since(s.cacheTime) < s.cacheTTL {
		return s.cachedRegistry, nil
	}

	// Load from disk
	reg, err := registry.LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry from disk: %w", err)
	}

	// Update cache
	s.cachedRegistry = reg
	s.cacheTime = time.Now()

	return reg, nil
}

// SaveRegistry saves the registry to disk and updates cache
func (s *registryService) SaveRegistry(ctx context.Context, reg *registry.Registry) error {
	if err := reg.Save(); err != nil {
		return fmt.Errorf("failed to save registry to disk: %w", err)
	}

	// Update cache
	s.mu.Lock()
	s.cachedRegistry = reg
	s.cacheTime = time.Now()
	s.mu.Unlock()

	return nil
}

// AddProject adds a project to the registry
func (s *registryService) AddProject(ctx context.Context, project *registry.Project) error {
	reg, err := s.LoadRegistry(ctx)
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := reg.AddProject(project); err != nil {
		return fmt.Errorf("failed to add project: %w", err)
	}

	return s.SaveRegistry(ctx, reg)
}

// FindProject finds a project by name
func (s *registryService) FindProject(ctx context.Context, name string) (*registry.Project, error) {
	reg, err := s.LoadRegistry(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	project, err := reg.FindProject(name)
	if err != nil {
		return nil, fmt.Errorf("project '%s' not found: %w", name, err)
	}

	return project, nil
}

// RemoveProject removes a project from the registry
func (s *registryService) RemoveProject(ctx context.Context, name string) error {
	reg, err := s.LoadRegistry(ctx)
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := reg.RemoveProject(name); err != nil {
		return fmt.Errorf("failed to remove project: %w", err)
	}

	return s.SaveRegistry(ctx, reg)
}

// ListProjects returns all projects in the registry
func (s *registryService) ListProjects(ctx context.Context) ([]registry.Project, error) {
	reg, err := s.LoadRegistry(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	return reg.ListProjects(), nil
}

// UpdateProjectStatus updates the status of a project
func (s *registryService) UpdateProjectStatus(ctx context.Context, name string) error {
	reg, err := s.LoadRegistry(ctx)
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := reg.UpdateProjectStatus(name); err != nil {
		return fmt.Errorf("failed to update project status: %w", err)
	}

	return s.SaveRegistry(ctx, reg)
}

// RefreshCache forces a refresh of the cached registry
func (s *registryService) RefreshCache(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear cache
	s.cachedRegistry = nil
	s.cacheTime = time.Time{}

	// Force reload
	_, err := s.LoadRegistry(ctx)
	return err
}

// ClearCache clears the cached registry
func (s *registryService) ClearCache(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cachedRegistry = nil
	s.cacheTime = time.Time{}

	return nil
}
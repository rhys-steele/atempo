package services

import (
	"context"

	"atempo/internal/utils"
)

// pathResolverService implements PathResolverService
type pathResolverService struct{}

// NewPathResolverService creates a new PathResolverService implementation
func NewPathResolverService() PathResolverService {
	return &pathResolverService{}
}

// ResolveProjectFromArgs resolves a project identifier from command arguments
func (s *pathResolverService) ResolveProjectFromArgs(ctx context.Context, args []string) (*ProjectResolution, error) {
	resolution, err := utils.ResolveProjectPathFromArgs(args)
	if err != nil {
		return nil, err
	}
	
	return &ProjectResolution{
		Path: resolution.Path,
		Name: resolution.Name,
	}, nil
}

// ResolveProjectPath resolves a project identifier to an absolute path
func (s *pathResolverService) ResolveProjectPath(ctx context.Context, identifier string) (string, error) {
	return utils.ResolveProjectPath(identifier)
}

// GetCurrentProjectPath returns the current working directory as a project path
func (s *pathResolverService) GetCurrentProjectPath(ctx context.Context) (*ProjectResolution, error) {
	resolution, err := utils.ResolveCurrentProjectPath()
	if err != nil {
		return nil, err
	}
	
	return &ProjectResolution{
		Path: resolution.Path,
		Name: resolution.Name,
	}, nil
}
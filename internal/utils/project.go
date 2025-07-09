package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectResolution represents the result of resolving a project identifier
type ProjectResolution struct {
	Path string // The resolved absolute path to the project
	Name string // The project name (either from args or directory basename)
}

// ResolveProjectPathFromArgs resolves a project identifier from command-line arguments
// to an absolute path. This consolidates the common pattern found across multiple commands.
//
// The function handles the following cases:
// 1. If args is empty or first argument is empty: uses current working directory
// 2. If args[0] is provided: resolves as path (registry lookup handled by caller)
// 3. Always returns an absolute path
// 4. Validates that the resolved directory exists
//
// Parameters:
//   - args: Command-line arguments slice (typically from command execution)
//
// Returns:
//   - ProjectResolution: Contains the resolved path and project name
//   - error: If resolution fails (invalid path, directory doesn't exist, etc.)
func ResolveProjectPathFromArgs(args []string) (*ProjectResolution, error) {
	var identifier string
	
	// Extract identifier from args
	if len(args) > 0 && args[0] != "" {
		identifier = args[0]
	}
	
	// Resolve the project path
	resolvedPath, err := resolveProjectIdentifier(identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project: %w", err)
	}
	
	// Validate that the directory exists
	if !FileExists(resolvedPath) {
		return nil, fmt.Errorf("project directory does not exist: %s", resolvedPath)
	}
	
	// Determine project name
	var projectName string
	if identifier != "" {
		// If identifier was provided, use it as the project name
		// This allows for both project names and paths to be used
		projectName = identifier
	} else {
		// If no identifier, use the directory basename
		projectName = filepath.Base(resolvedPath)
	}
	
	return &ProjectResolution{
		Path: resolvedPath,
		Name: projectName,
	}, nil
}

// ResolveProjectPath resolves a single project identifier to an absolute path
// This is a simpler version that just returns the path without additional metadata
// Note: This version does not perform registry lookup to avoid circular imports
//
// Parameters:
//   - identifier: Project identifier (relative path, or absolute path)
//
// Returns:
//   - string: The resolved absolute path
//   - error: If resolution fails
func ResolveProjectPath(identifier string) (string, error) {
	resolvedPath, err := resolveProjectIdentifier(identifier)
	if err != nil {
		return "", err
	}
	
	// Validate that the directory exists
	if !FileExists(resolvedPath) {
		return "", fmt.Errorf("project directory does not exist: %s", resolvedPath)
	}
	
	return resolvedPath, nil
}

// resolveProjectIdentifier is the core resolution logic that handles:
// 1. Empty identifier -> current working directory
// 2. Path resolution (absolute or relative)
// Note: Registry lookup is handled by the registry package to avoid circular imports
func resolveProjectIdentifier(identifier string) (string, error) {
	// If empty, use current directory
	if identifier == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		return cwd, nil
	}

	// If absolute path, use directly
	if filepath.IsAbs(identifier) {
		return identifier, nil
	}

	// Convert relative path to absolute
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	return filepath.Join(cwd, identifier), nil
}

// ResolveProjectPathFromArgsWithOffset resolves project path from args with an offset
// This is useful for commands that have other arguments before the project identifier
//
// Parameters:
//   - args: Command-line arguments slice
//   - offset: Index in args where the project identifier is located
//
// Returns:
//   - ProjectResolution: Contains the resolved path and project name
//   - error: If resolution fails or offset is out of bounds
func ResolveProjectPathFromArgsWithOffset(args []string, offset int) (*ProjectResolution, error) {
	var projectArgs []string
	
	// Extract the project identifier from the specified offset
	if len(args) > offset {
		projectArgs = args[offset:]
	}
	
	return ResolveProjectPathFromArgs(projectArgs)
}

// ResolveCurrentProjectPath resolves the current working directory as a project path
// This is useful for commands that always operate on the current directory
//
// Returns:
//   - ProjectResolution: Contains the current directory path and name
//   - error: If current directory cannot be determined or doesn't exist
func ResolveCurrentProjectPath() (*ProjectResolution, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	
	// Validate that the directory exists (should always be true for cwd)
	if !FileExists(cwd) {
		return nil, fmt.Errorf("current directory does not exist: %s", cwd)
	}
	
	return &ProjectResolution{
		Path: cwd,
		Name: filepath.Base(cwd),
	}, nil
}
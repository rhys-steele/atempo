package interfaces

// ProjectResolver defines the interface for resolving project paths and identifiers
type ProjectResolver interface {
	// ResolveFromArgs resolves a project identifier from command arguments
	ResolveFromArgs(args []string) (string, error)
	
	// ResolveFromPath resolves a project identifier to an absolute path
	ResolveFromPath(identifier string) (string, error)
	
	// GetCurrentProjectPath returns the current working directory as a project path
	GetCurrentProjectPath() (string, error)
	
	// ValidatePath validates that a path exists and is accessible
	ValidatePath(path string) error
}
package interfaces

// StatusFormatter defines the interface for status display formatting
type StatusFormatter interface {
	// FormatStatus formats a status string with appropriate icon and color
	FormatStatus(status string) (icon, color string)
	
	// FormatProjectStatus formats project status for display
	FormatProjectStatus(status string) string
	
	// FormatServiceStatus formats service status for display
	FormatServiceStatus(status string) string
	
	// GetStatusIcon returns the appropriate icon for a status
	GetStatusIcon(status string) string
	
	// GetStatusColor returns the appropriate color for a status
	GetStatusColor(status string) string
}

// ProgressTracker defines the interface for tracking operation progress
type ProgressTracker interface {
	// StartStep starts a new step with the given description
	StartStep(description string)
	
	// CompleteStep marks the current step as completed
	CompleteStep()
	
	// FailStep marks the current step as failed
	FailStep(err error)
	
	// UpdateProgress updates the progress with current status
	UpdateProgress(current, total int, message string)
	
	// Finish completes the progress tracking
	Finish()
}
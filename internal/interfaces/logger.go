package interfaces

// Logger defines the interface for project logging operations
type Logger interface {
	// StartStep starts a new logging step
	StartStep(description string) LogStep
	
	// Log writes a message to the log
	Log(message string) error
	
	// LogError writes an error message to the log
	LogError(err error) error
	
	// LogInfo writes an info message to the log
	LogInfo(message string) error
	
	// LogWarning writes a warning message to the log
	LogWarning(message string) error
	
	// LogDebug writes a debug message to the log
	LogDebug(message string) error
	
	// Close closes the logger and flushes any remaining data
	Close() error
	
	// GetLogPath returns the path to the log file
	GetLogPath() string
}

// LogStep defines the interface for individual logging steps
type LogStep interface {
	// SetStatus sets the status of the current step
	SetStatus(status LogStatus) error
	
	// Log writes a message within this step
	Log(message string) error
	
	// LogError writes an error message within this step
	LogError(err error) error
	
	// Complete marks the step as completed
	Complete() error
	
	// Fail marks the step as failed
	Fail(err error) error
	
	// GetStatus returns the current status of the step
	GetStatus() LogStatus
}

// LogStatus represents the status of a logging step
type LogStatus int

const (
	// LogStatusPending indicates the step is pending
	LogStatusPending LogStatus = iota
	
	// LogStatusRunning indicates the step is currently running
	LogStatusRunning
	
	// LogStatusCompleted indicates the step completed successfully
	LogStatusCompleted
	
	// LogStatusFailed indicates the step failed
	LogStatusFailed
	
	// LogStatusSkipped indicates the step was skipped
	LogStatusSkipped
)

// String returns the string representation of the log status
func (s LogStatus) String() string {
	switch s {
	case LogStatusPending:
		return "pending"
	case LogStatusRunning:
		return "running"
	case LogStatusCompleted:
		return "completed"
	case LogStatusFailed:
		return "failed"
	case LogStatusSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}
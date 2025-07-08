package output

// OutputWriter provides an abstraction for output operations
type OutputWriter interface {
	// Success outputs a success message with optional details
	Success(message, details string)

	// Error outputs an error message with optional details
	Error(message, details string)

	// Info outputs an informational message
	Info(message string)

	// Warning outputs a warning message
	Warning(message string)

	// Debug outputs a debug message (may be suppressed based on verbosity)
	Debug(message string)

	// Progress outputs a progress message
	Progress(message string)

	// Printf provides formatted output (for compatibility)
	Printf(format string, args ...interface{})

	// Println provides line output (for compatibility)
	Println(args ...interface{})
}

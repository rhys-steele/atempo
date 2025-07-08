package output

import (
	"fmt"
)

// ConsoleWriter implements OutputWriter for console output
type ConsoleWriter struct {
	verbose bool
}

// NewConsoleWriter creates a new console output writer
func NewConsoleWriter(verbose bool) *ConsoleWriter {
	return &ConsoleWriter{
		verbose: verbose,
	}
}

// Success outputs a success message with optional details
func (c *ConsoleWriter) Success(message, details string) {
	fmt.Printf("✓ %s", message)
	if details != "" {
		fmt.Printf(": %s", details)
	}
	fmt.Println()
}

// Error outputs an error message with optional details
func (c *ConsoleWriter) Error(message, details string) {
	fmt.Printf("✗ %s", message)
	if details != "" {
		fmt.Printf(": %s", details)
	}
	fmt.Println()
}

// Info outputs an informational message
func (c *ConsoleWriter) Info(message string) {
	fmt.Printf("→ %s\n", message)
}

// Warning outputs a warning message
func (c *ConsoleWriter) Warning(message string) {
	fmt.Printf("! %s\n", message)
}

// Debug outputs a debug message (may be suppressed based on verbosity)
func (c *ConsoleWriter) Debug(message string) {
	if c.verbose {
		fmt.Printf("→ %s\n", message)
	}
}

// Progress outputs a progress message
func (c *ConsoleWriter) Progress(message string) {
	fmt.Printf("  %s\n", message)
}

// Printf provides formatted output (for compatibility)
func (c *ConsoleWriter) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Println provides line output (for compatibility)
func (c *ConsoleWriter) Println(args ...interface{}) {
	fmt.Println(args...)
}

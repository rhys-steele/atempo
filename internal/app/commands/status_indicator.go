package commands

import (
	"fmt"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
)

// Status types
type StatusType int

const (
	StatusThinking StatusType = iota
	StatusWorking
	StatusSuccess
	StatusError
	StatusInfo
	StatusWarning
)

// StatusIndicator provides Claude Code-style status feedback
type StatusIndicator struct {
	message   string
	startTime time.Time
}

// NewStatusIndicator creates a new status indicator
func NewStatusIndicator() *StatusIndicator {
	return &StatusIndicator{
		startTime: time.Now(),
	}
}

// Start displays a status with animated indicator
func (s *StatusIndicator) Start(status StatusType, message string) {
	s.message = message
	s.startTime = time.Now()
	
	symbol, color := s.getStatusSymbolAndColor(status)
	fmt.Printf("%s%s%s %s\n", color, symbol, ColorReset, message)
}

// Update shows progress with elapsed time
func (s *StatusIndicator) Update(status StatusType, message string) {
	elapsed := time.Since(s.startTime)
	symbol, color := s.getStatusSymbolAndColor(status)
	
	fmt.Printf("%s%s%s %s %s(%s)%s\n", 
		color, symbol, ColorReset, 
		message,
		ColorGray, s.formatDuration(elapsed), ColorReset)
}

// Success shows a success message with details
func (s *StatusIndicator) Success(message, details string) {
	elapsed := time.Since(s.startTime)
	fmt.Printf("%s⏺%s %s\n", ColorGreen, ColorReset, message)
	if details != "" {
		fmt.Printf("  %s⎿%s  %s %s(%s)%s\n", 
			ColorGray, ColorReset, 
			details,
			ColorGray, s.formatDuration(elapsed), ColorReset)
	}
}

// Error shows an error message with details
func (s *StatusIndicator) Error(message, details string) {
	elapsed := time.Since(s.startTime)
	fmt.Printf("%s⏺%s %s\n", ColorRed, ColorReset, message)
	if details != "" {
		fmt.Printf("  %s⎿%s  %s %s(%s)%s\n", 
			ColorGray, ColorReset, 
			details,
			ColorGray, s.formatDuration(elapsed), ColorReset)
	}
}

// Info shows an info message
func (s *StatusIndicator) Info(message string) {
	fmt.Printf("%s⏺%s %s\n", ColorBlue, ColorReset, message)
}

// Warning shows a warning message
func (s *StatusIndicator) Warning(message string) {
	fmt.Printf("%s⏺%s %s\n", ColorYellow, ColorReset, message)
}

// getStatusSymbolAndColor returns the appropriate symbol and color for status
func (s *StatusIndicator) getStatusSymbolAndColor(status StatusType) (string, string) {
	switch status {
	case StatusThinking:
		return "✶", ColorBlue
	case StatusWorking:
		return "⚡", ColorYellow
	case StatusSuccess:
		return "✓", ColorGreen
	case StatusError:
		return "✗", ColorRed
	case StatusInfo:
		return "⏺", ColorBlue
	case StatusWarning:
		return "⚠", ColorYellow
	default:
		return "•", ColorWhite
	}
}

// formatDuration formats elapsed time nicely
func (s *StatusIndicator) formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fm", d.Minutes())
}

// Global status indicator instance
var GlobalStatus = NewStatusIndicator()

// Helper functions for common status updates
func ShowThinking(message string) {
	GlobalStatus.Start(StatusThinking, message)
}

func ShowWorking(message string) {
	GlobalStatus.Start(StatusWorking, message)
}

func ShowSuccess(message, details string) {
	GlobalStatus.Success(message, details)
}

func ShowError(message, details string) {
	GlobalStatus.Error(message, details)
}

func ShowInfo(message string) {
	GlobalStatus.Info(message)
}

func ShowWarning(message string) {
	GlobalStatus.Warning(message)
}
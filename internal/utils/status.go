package utils

import (
	"fmt"
	"strings"
)

// ANSI color codes for status display
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

// StatusInfo represents the formatting information for a status
type StatusInfo struct {
	Icon  string
	Color string
	Label string
}

// StatusDisplay provides centralized status display functionality
type StatusDisplay struct {
	useColor bool
}

// NewStatusDisplay creates a new status display utility
func NewStatusDisplay() *StatusDisplay {
	return &StatusDisplay{
		useColor: true,
	}
}

// NewStatusDisplayNoColor creates a status display without color codes
func NewStatusDisplayNoColor() *StatusDisplay {
	return &StatusDisplay{
		useColor: false,
	}
}

// GetStatusInfo returns the complete status information (icon, color, label)
func (s *StatusDisplay) GetStatusInfo(status string) StatusInfo {
	switch status {
	case "running":
		return StatusInfo{
			Icon:  "✓",
			Color: s.colorOrEmpty(ColorGreen),
			Label: "running",
		}
	case "partial":
		return StatusInfo{
			Icon:  "⚠",
			Color: s.colorOrEmpty(ColorYellow),
			Label: "partial",
		}
	case "stopped", "no-docker", "no-services":
		return StatusInfo{
			Icon:  "✗",
			Color: s.colorOrEmpty(ColorRed),
			Label: "stopped",
		}
	case "docker-error":
		return StatusInfo{
			Icon:  "✗",
			Color: s.colorOrEmpty(ColorRed),
			Label: "error",
		}
	case "unknown":
		return StatusInfo{
			Icon:  "?",
			Color: s.colorOrEmpty(ColorGray),
			Label: "unknown",
		}
	default:
		return StatusInfo{
			Icon:  "?",
			Color: s.colorOrEmpty(ColorGray),
			Label: status,
		}
	}
}

// GetStatusIcon returns just the icon for a status
func (s *StatusDisplay) GetStatusIcon(status string) string {
	return s.GetStatusInfo(status).Icon
}

// GetStatusColor returns just the color for a status
func (s *StatusDisplay) GetStatusColor(status string) string {
	return s.GetStatusInfo(status).Color
}

// GetStatusLabel returns the normalized label for a status
func (s *StatusDisplay) GetStatusLabel(status string) string {
	return s.GetStatusInfo(status).Label
}

// FormatStatus returns a formatted status string with icon and color
func (s *StatusDisplay) FormatStatus(status string) string {
	info := s.GetStatusInfo(status)
	if s.useColor {
		return fmt.Sprintf("%s%s%s %s", info.Color, info.Icon, ColorReset, info.Label)
	}
	return fmt.Sprintf("%s %s", info.Icon, info.Label)
}

// FormatStatusWithIcon returns just the icon and status
func (s *StatusDisplay) FormatStatusWithIcon(status string) string {
	info := s.GetStatusInfo(status)
	return fmt.Sprintf("%s %s", info.Icon, info.Label)
}

// FormatColoredStatus returns a colored status without icon
func (s *StatusDisplay) FormatColoredStatus(status string) string {
	info := s.GetStatusInfo(status)
	if s.useColor {
		return fmt.Sprintf("%s%s%s", info.Color, info.Label, ColorReset)
	}
	return info.Label
}

// FormatServiceStatus formats service status for display
func (s *StatusDisplay) FormatServiceStatus(service, status string) string {
	info := s.GetStatusInfo(status)
	if s.useColor {
		return fmt.Sprintf("  %s%s%s %s: %s%s%s",
			info.Color, info.Icon, ColorReset,
			service,
			info.Color, info.Label, ColorReset)
	}
	return fmt.Sprintf("  %s %s: %s", info.Icon, service, info.Label)
}

// CountStatuses counts projects by status type
func (s *StatusDisplay) CountStatuses(statuses []string) (running, stopped, error int) {
	for _, status := range statuses {
		switch status {
		case "running", "partial":
			running++
		case "stopped", "no-docker", "no-services":
			stopped++
		case "docker-error", "unknown":
			error++
		}
	}
	return running, stopped, error
}

// FormatStatusSummary formats a summary of status counts
func (s *StatusDisplay) FormatStatusSummary(running, stopped, error int) string {
	parts := []string{}
	
	if running > 0 {
		if s.useColor {
			parts = append(parts, fmt.Sprintf("%s%d running%s", ColorGreen, running, ColorReset))
		} else {
			parts = append(parts, fmt.Sprintf("%d running", running))
		}
	}
	
	if stopped > 0 {
		if s.useColor {
			parts = append(parts, fmt.Sprintf("%s%d stopped%s", ColorRed, stopped, ColorReset))
		} else {
			parts = append(parts, fmt.Sprintf("%d stopped", stopped))
		}
	}
	
	if error > 0 {
		if s.useColor {
			parts = append(parts, fmt.Sprintf("%s%d errors%s", ColorRed, error, ColorReset))
		} else {
			parts = append(parts, fmt.Sprintf("%d errors", error))
		}
	}
	
	if len(parts) == 0 {
		return "No projects"
	}
	
	return strings.Join(parts, ", ")
}

// FormatProgressBar creates a simple progress bar
func (s *StatusDisplay) FormatProgressBar(current, total int, width int) string {
	if total == 0 {
		return strings.Repeat("─", width)
	}
	
	filled := int(float64(current) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}
	
	progress := strings.Repeat("█", filled) + strings.Repeat("─", width-filled)
	
	if s.useColor {
		return fmt.Sprintf("%s%s%s", ColorGreen, progress, ColorReset)
	}
	return progress
}

// colorOrEmpty returns the color if color mode is enabled, empty string otherwise
func (s *StatusDisplay) colorOrEmpty(color string) string {
	if s.useColor {
		return color
	}
	return ""
}

// Global status display instance
var GlobalStatusDisplay = NewStatusDisplay()

// Helper functions for common status formatting
func GetStatusIcon(status string) string {
	return GlobalStatusDisplay.GetStatusIcon(status)
}

func GetStatusColor(status string) string {
	return GlobalStatusDisplay.GetStatusColor(status)
}

func FormatStatus(status string) string {
	return GlobalStatusDisplay.FormatStatus(status)
}

func FormatStatusWithIcon(status string) string {
	return GlobalStatusDisplay.FormatStatusWithIcon(status)
}

func FormatColoredStatus(status string) string {
	return GlobalStatusDisplay.FormatColoredStatus(status)
}

func FormatServiceStatus(service, status string) string {
	return GlobalStatusDisplay.FormatServiceStatus(service, status)
}

func CountStatuses(statuses []string) (running, stopped, error int) {
	return GlobalStatusDisplay.CountStatuses(statuses)
}

func FormatStatusSummary(running, stopped, error int) string {
	return GlobalStatusDisplay.FormatStatusSummary(running, stopped, error)
}
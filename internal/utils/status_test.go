package utils

import (
	"strings"
	"testing"
)

func TestStatusDisplay_GetStatusInfo(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		status     string
		wantIcon   string
		wantColor  string
		wantLabel  string
	}{
		{"running", "✓", ColorGreen, "running"},
		{"partial", "⚠", ColorYellow, "partial"},
		{"stopped", "✗", ColorRed, "stopped"},
		{"no-docker", "✗", ColorRed, "stopped"},
		{"no-services", "✗", ColorRed, "stopped"},
		{"docker-error", "✗", ColorRed, "error"},
		{"unknown", "?", ColorGray, "unknown"},
		{"custom-status", "?", ColorGray, "custom-status"},
	}
	
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			info := display.GetStatusInfo(tt.status)
			
			if info.Icon != tt.wantIcon {
				t.Errorf("GetStatusInfo(%q).Icon = %q, want %q", tt.status, info.Icon, tt.wantIcon)
			}
			if info.Color != tt.wantColor {
				t.Errorf("GetStatusInfo(%q).Color = %q, want %q", tt.status, info.Color, tt.wantColor)
			}
			if info.Label != tt.wantLabel {
				t.Errorf("GetStatusInfo(%q).Label = %q, want %q", tt.status, info.Label, tt.wantLabel)
			}
		})
	}
}

func TestStatusDisplay_GetStatusIcon(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		status string
		want   string
	}{
		{"running", "✓"},
		{"partial", "⚠"},
		{"stopped", "✗"},
		{"docker-error", "✗"},
		{"unknown", "?"},
		{"custom", "?"},
	}
	
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := display.GetStatusIcon(tt.status)
			if got != tt.want {
				t.Errorf("GetStatusIcon(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay_GetStatusColor(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		status string
		want   string
	}{
		{"running", ColorGreen},
		{"partial", ColorYellow},
		{"stopped", ColorRed},
		{"docker-error", ColorRed},
		{"unknown", ColorGray},
	}
	
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := display.GetStatusColor(tt.status)
			if got != tt.want {
				t.Errorf("GetStatusColor(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay_FormatStatus(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		status string
		want   string
	}{
		{"running", ColorGreen + "✓" + ColorReset + " running"},
		{"partial", ColorYellow + "⚠" + ColorReset + " partial"},
		{"stopped", ColorRed + "✗" + ColorReset + " stopped"},
		{"docker-error", ColorRed + "✗" + ColorReset + " error"},
	}
	
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := display.FormatStatus(tt.status)
			if got != tt.want {
				t.Errorf("FormatStatus(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay_FormatStatusNoColor(t *testing.T) {
	display := NewStatusDisplayNoColor()
	
	tests := []struct {
		status string
		want   string
	}{
		{"running", "✓ running"},
		{"partial", "⚠ partial"},
		{"stopped", "✗ stopped"},
		{"docker-error", "✗ error"},
	}
	
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := display.FormatStatus(tt.status)
			if got != tt.want {
				t.Errorf("FormatStatus(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay_FormatStatusWithIcon(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		status string
		want   string
	}{
		{"running", "✓ running"},
		{"partial", "⚠ partial"},
		{"stopped", "✗ stopped"},
		{"docker-error", "✗ error"},
	}
	
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := display.FormatStatusWithIcon(tt.status)
			if got != tt.want {
				t.Errorf("FormatStatusWithIcon(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay_FormatColoredStatus(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		status string
		want   string
	}{
		{"running", ColorGreen + "running" + ColorReset},
		{"partial", ColorYellow + "partial" + ColorReset},
		{"stopped", ColorRed + "stopped" + ColorReset},
		{"docker-error", ColorRed + "error" + ColorReset},
	}
	
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := display.FormatColoredStatus(tt.status)
			if got != tt.want {
				t.Errorf("FormatColoredStatus(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay_FormatServiceStatus(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		service string
		status  string
		want    string
	}{
		{"app", "running", "  " + ColorGreen + "✓" + ColorReset + " app: " + ColorGreen + "running" + ColorReset},
		{"db", "stopped", "  " + ColorRed + "✗" + ColorReset + " db: " + ColorRed + "stopped" + ColorReset},
		{"redis", "partial", "  " + ColorYellow + "⚠" + ColorReset + " redis: " + ColorYellow + "partial" + ColorReset},
	}
	
	for _, tt := range tests {
		t.Run(tt.service+"_"+tt.status, func(t *testing.T) {
			got := display.FormatServiceStatus(tt.service, tt.status)
			if got != tt.want {
				t.Errorf("FormatServiceStatus(%q, %q) = %q, want %q", tt.service, tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay_CountStatuses(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		name     string
		statuses []string
		wantRunning int
		wantStopped int
		wantError   int
	}{
		{
			name:     "mixed statuses",
			statuses: []string{"running", "stopped", "partial", "docker-error", "no-docker", "unknown"},
			wantRunning: 2, // running + partial
			wantStopped: 2, // stopped + no-docker
			wantError:   2, // docker-error + unknown
		},
		{
			name:     "all running",
			statuses: []string{"running", "running", "partial"},
			wantRunning: 3,
			wantStopped: 0,
			wantError:   0,
		},
		{
			name:     "all stopped",
			statuses: []string{"stopped", "no-docker", "no-services"},
			wantRunning: 0,
			wantStopped: 3,
			wantError:   0,
		},
		{
			name:     "empty",
			statuses: []string{},
			wantRunning: 0,
			wantStopped: 0,
			wantError:   0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			running, stopped, error := display.CountStatuses(tt.statuses)
			if running != tt.wantRunning {
				t.Errorf("CountStatuses() running = %d, want %d", running, tt.wantRunning)
			}
			if stopped != tt.wantStopped {
				t.Errorf("CountStatuses() stopped = %d, want %d", stopped, tt.wantStopped)
			}
			if error != tt.wantError {
				t.Errorf("CountStatuses() error = %d, want %d", error, tt.wantError)
			}
		})
	}
}

func TestStatusDisplay_FormatStatusSummary(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		name    string
		running int
		stopped int
		error   int
		want    string
	}{
		{
			name:    "mixed",
			running: 2,
			stopped: 1,
			error:   1,
			want:    ColorGreen + "2 running" + ColorReset + ", " + ColorRed + "1 stopped" + ColorReset + ", " + ColorRed + "1 errors" + ColorReset,
		},
		{
			name:    "running only",
			running: 3,
			stopped: 0,
			error:   0,
			want:    ColorGreen + "3 running" + ColorReset,
		},
		{
			name:    "stopped only",
			running: 0,
			stopped: 2,
			error:   0,
			want:    ColorRed + "2 stopped" + ColorReset,
		},
		{
			name:    "none",
			running: 0,
			stopped: 0,
			error:   0,
			want:    "No projects",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := display.FormatStatusSummary(tt.running, tt.stopped, tt.error)
			if got != tt.want {
				t.Errorf("FormatStatusSummary(%d, %d, %d) = %q, want %q", tt.running, tt.stopped, tt.error, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay_FormatProgressBar(t *testing.T) {
	display := NewStatusDisplay()
	
	tests := []struct {
		name    string
		current int
		total   int
		width   int
		want    string
	}{
		{
			name:    "half progress",
			current: 50,
			total:   100,
			width:   10,
			want:    ColorGreen + "█████─────" + ColorReset,
		},
		{
			name:    "full progress",
			current: 100,
			total:   100,
			width:   10,
			want:    ColorGreen + "██████████" + ColorReset,
		},
		{
			name:    "no progress",
			current: 0,
			total:   100,
			width:   10,
			want:    ColorGreen + "──────────" + ColorReset,
		},
		{
			name:    "zero total",
			current: 0,
			total:   0,
			width:   10,
			want:    "──────────",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := display.FormatProgressBar(tt.current, tt.total, tt.width)
			if got != tt.want {
				t.Errorf("FormatProgressBar(%d, %d, %d) = %q, want %q", tt.current, tt.total, tt.width, got, tt.want)
			}
		})
	}
}

func TestGlobalHelperFunctions(t *testing.T) {
	// Test that global helper functions work correctly
	icon := GetStatusIcon("running")
	if icon != "✓" {
		t.Errorf("GetStatusIcon('running') = %q, want '✓'", icon)
	}
	
	color := GetStatusColor("running")
	if color != ColorGreen {
		t.Errorf("GetStatusColor('running') = %q, want %q", color, ColorGreen)
	}
	
	formatted := FormatStatus("running")
	expected := ColorGreen + "✓" + ColorReset + " running"
	if formatted != expected {
		t.Errorf("FormatStatus('running') = %q, want %q", formatted, expected)
	}
	
	// Test status counting
	running, stopped, errors := CountStatuses([]string{"running", "stopped", "partial"})
	if running != 2 || stopped != 1 || errors != 0 {
		t.Errorf("CountStatuses() = (%d, %d, %d), want (2, 1, 0)", running, stopped, errors)
	}
	
	// Test summary formatting
	summary := FormatStatusSummary(2, 1, 0)
	if !strings.Contains(summary, "2 running") || !strings.Contains(summary, "1 stopped") {
		t.Errorf("FormatStatusSummary(2, 1, 0) = %q, want to contain '2 running' and '1 stopped'", summary)
	}
}
package commands

import (
	"fmt"
	"strings"
	"time"
)

// ProgressTracker provides real-time progress updates with animated indicators
type ProgressTracker struct {
	currentStep   string
	totalSteps    int
	currentIndex  int
	startTime     time.Time
	stepStartTime time.Time
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(totalSteps int) *ProgressTracker {
	return &ProgressTracker{
		totalSteps: totalSteps,
		startTime:  time.Now(),
	}
}

// StartStep begins a new step with progress indicator
func (p *ProgressTracker) StartStep(stepIndex int, description string) {
	p.currentIndex = stepIndex
	p.currentStep = description
	p.stepStartTime = time.Now()
	
	// Start the step with a thinking indicator
	fmt.Printf("%s✶%s %s %s[%d/%d]%s\n", 
		ColorBlue, ColorReset, 
		description,
		ColorGray, stepIndex, p.totalSteps, ColorReset)
}

// UpdateStep provides a sub-step update within the current step
func (p *ProgressTracker) UpdateStep(subDescription string) {
	fmt.Printf("%s  ⚡%s %s\n", ColorYellow, ColorReset, subDescription)
}

// CompleteStep marks the current step as complete
func (p *ProgressTracker) CompleteStep(details string) {
	elapsed := time.Since(p.stepStartTime)
	fmt.Printf("%s⏺%s %s\n", ColorGreen, ColorReset, p.currentStep)
	if details != "" {
		fmt.Printf("  %s⎿%s  %s %s(%s)%s\n", 
			ColorGray, ColorReset,
			details,
			ColorGray, p.formatDuration(elapsed), ColorReset)
	}
	fmt.Println()
}

// ErrorStep marks the current step as failed
func (p *ProgressTracker) ErrorStep(errorMsg string) {
	elapsed := time.Since(p.stepStartTime)
	fmt.Printf("%s⏺%s %s\n", ColorRed, ColorReset, p.currentStep)
	fmt.Printf("  %s⎿%s  %s %s(%s)%s\n", 
		ColorGray, ColorReset,
		errorMsg,
		ColorGray, p.formatDuration(elapsed), ColorReset)
	fmt.Println()
}

// WarningStep shows a warning within the current step
func (p *ProgressTracker) WarningStep(warningMsg string) {
	fmt.Printf("%s  ⚠%s  %s\n", ColorYellow, ColorReset, warningMsg)
}

// ShowProgress displays overall progress
func (p *ProgressTracker) ShowProgress() {
	totalElapsed := time.Since(p.startTime)
	percentage := float64(p.currentIndex) / float64(p.totalSteps) * 100
	
	fmt.Printf("%sProgress: %.0f%% (%d/%d) • %s elapsed%s\n", 
		ColorGray, percentage, p.currentIndex, p.totalSteps, 
		p.formatDuration(totalElapsed), ColorReset)
}

// Complete marks the entire process as complete
func (p *ProgressTracker) Complete(projectName string) {
	totalElapsed := time.Since(p.startTime)
	fmt.Printf("%s🎉 Project '%s' created successfully!%s\n", ColorGreen, projectName, ColorReset)
	fmt.Printf("%s   Total time: %s%s\n\n", ColorGray, p.formatDuration(totalElapsed), ColorReset)
	
	// Show next steps
	fmt.Printf("%s💡 Next steps:%s\n", ColorBlue, ColorReset)
	fmt.Printf("   %satempo status%s           # Check project status\n", ColorCyan, ColorReset)
	fmt.Printf("   %satempo docker up%s        # Start development services\n", ColorCyan, ColorReset)
	fmt.Printf("   %satempo describe%s         # View project details\n", ColorCyan, ColorReset)
}


// formatDuration formats elapsed time nicely
func (p *ProgressTracker) formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fm", d.Minutes())
}

// Progress bar visualization
func (p *ProgressTracker) DrawProgressBar(width int) string {
	if p.totalSteps == 0 {
		return ""
	}
	
	filled := int(float64(p.currentIndex) / float64(p.totalSteps) * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := float64(p.currentIndex) / float64(p.totalSteps) * 100
	
	return fmt.Sprintf("%s[%s] %.0f%%%s", ColorBlue, bar, percentage, ColorReset)
}

// CreateSteps defines the standard steps for project creation
type CreateSteps struct {
	LoadTemplate      int
	InstallFramework  int
	CopyTemplateFiles int
	PostInstallSetup  int
	FinalizeProject   int
}

// StandardCreateSteps returns the standard step definitions
func StandardCreateSteps() CreateSteps {
	return CreateSteps{
		LoadTemplate:      1,
		InstallFramework:  2,
		CopyTemplateFiles: 3,
		PostInstallSetup:  4,
		FinalizeProject:   5,
	}
}
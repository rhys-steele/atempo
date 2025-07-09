package services

import (
	"context"
	"fmt"
	"strings"

	"atempo/internal/registry"
)

// StatusService provides business operations for status formatting and tracking
type StatusService interface {
	// FormatProjectStatus formats a project's status with appropriate styling
	FormatProjectStatus(ctx context.Context, project *registry.Project) (string, error)
	
	// FormatServiceStatus formats a service's status with appropriate styling
	FormatServiceStatus(ctx context.Context, service *registry.Service) (string, error)
	
	// GetStatusIcon returns the appropriate icon for a status
	GetStatusIcon(ctx context.Context, status string) string
	
	// GetStatusColor returns the appropriate color code for a status
	GetStatusColor(ctx context.Context, status string) string
	
	// FormatDashboard formats a complete project dashboard
	FormatDashboard(ctx context.Context, projects []registry.Project) (string, error)
	
	// FormatProjectSummary formats a summary of a single project
	FormatProjectSummary(ctx context.Context, project *registry.Project) (string, error)
}

// ProgressService provides business operations for progress tracking
type ProgressService interface {
	// StartProgress starts a new progress tracker with the given steps
	StartProgress(ctx context.Context, totalSteps int) (ProgressTracker, error)
	
	// UpdateProgress updates the progress tracker
	UpdateProgress(ctx context.Context, tracker ProgressTracker, currentStep int, message string) error
	
	// CompleteProgress completes the progress tracker
	CompleteProgress(ctx context.Context, tracker ProgressTracker) error
}

// ProgressTracker represents an active progress tracking session
type ProgressTracker interface {
	// StartStep starts a new step with the given message
	StartStep(message string)
	
	// CompleteStep completes the current step
	CompleteStep()
	
	// FailStep marks the current step as failed
	FailStep(err error)
	
	// GetCurrentStep returns the current step number
	GetCurrentStep() int
	
	// GetTotalSteps returns the total number of steps
	GetTotalSteps() int
	
	// IsComplete returns whether all steps are complete
	IsComplete() bool
}

// statusService implements StatusService
type statusService struct{}

// NewStatusService creates a new StatusService implementation
func NewStatusService() StatusService {
	return &statusService{}
}

// FormatProjectStatus formats a project's status with appropriate styling
func (s *statusService) FormatProjectStatus(ctx context.Context, project *registry.Project) (string, error) {
	if project == nil {
		return "", fmt.Errorf("project cannot be nil")
	}

	icon := s.GetStatusIcon(ctx, project.Status)
	color := s.GetStatusColor(ctx, project.Status)
	
	// Format: "✓ project-name (running) - Laravel 11"
	return fmt.Sprintf("%s%s %s (%s) - %s %s%s",
		color,
		icon,
		project.Name,
		project.Status,
		strings.Title(project.Framework),
		project.Version,
		getColorReset(),
	), nil
}

// FormatServiceStatus formats a service's status with appropriate styling
func (s *statusService) FormatServiceStatus(ctx context.Context, service *registry.Service) (string, error) {
	if service == nil {
		return "", fmt.Errorf("service cannot be nil")
	}

	icon := s.GetStatusIcon(ctx, service.Status)
	color := s.GetStatusColor(ctx, service.Status)
	
	// Format: "  ✓ app (running) - http://localhost:8000"
	status := fmt.Sprintf("  %s%s %s (%s)%s",
		color,
		icon,
		service.Name,
		service.Status,
		getColorReset(),
	)
	
	if service.URL != "" {
		status += fmt.Sprintf(" - %s", service.URL)
	}
	
	return status, nil
}

// GetStatusIcon returns the appropriate icon for a status
func (s *statusService) GetStatusIcon(ctx context.Context, status string) string {
	switch status {
	case "running", "active", "healthy", "up", "completed":
		return "✓"
	case "stopped", "inactive", "down", "exited":
		return "•"
	case "failed", "error", "unhealthy":
		return "✗"
	case "partial", "starting", "restarting":
		return "-"
	case "unknown", "pending":
		return "?"
	default:
		return "•"
	}
}

// GetStatusColor returns the appropriate color code for a status
func (s *statusService) GetStatusColor(ctx context.Context, status string) string {
	switch status {
	case "running", "active", "healthy", "up", "completed":
		return getGreenColor()
	case "stopped", "inactive", "down", "exited":
		return getGrayColor()
	case "failed", "error", "unhealthy":
		return getRedColor()
	case "partial", "starting", "restarting":
		return getYellowColor()
	case "unknown", "pending":
		return getCyanColor()
	default:
		return getDefaultColor()
	}
}

// FormatDashboard formats a complete project dashboard
func (s *statusService) FormatDashboard(ctx context.Context, projects []registry.Project) (string, error) {
	if len(projects) == 0 {
		return "No projects found.", nil
	}

	var output strings.Builder
	output.WriteString("Project Status Dashboard\n")
	output.WriteString("────────────────────────\n\n")

	for _, project := range projects {
		// Format project header
		projectStatus, err := s.FormatProjectStatus(ctx, &project)
		if err != nil {
			return "", fmt.Errorf("failed to format project status: %w", err)
		}
		output.WriteString(projectStatus + "\n")

		// Format services
		if len(project.Services) > 0 {
			for _, service := range project.Services {
				serviceStatus, err := s.FormatServiceStatus(ctx, &service)
				if err != nil {
					return "", fmt.Errorf("failed to format service status: %w", err)
				}
				output.WriteString(serviceStatus + "\n")
			}
		} else {
			output.WriteString("  No services configured\n")
		}
		
		output.WriteString("\n")
	}

	return output.String(), nil
}

// FormatProjectSummary formats a summary of a single project
func (s *statusService) FormatProjectSummary(ctx context.Context, project *registry.Project) (string, error) {
	if project == nil {
		return "", fmt.Errorf("project cannot be nil")
	}

	var output strings.Builder
	
	// Project header
	output.WriteString(fmt.Sprintf("Project: %s\n", project.Name))
	output.WriteString("────────────────\n\n")
	
	// Basic info
	output.WriteString(fmt.Sprintf("Framework: %s %s\n", strings.Title(project.Framework), project.Version))
	output.WriteString(fmt.Sprintf("Path: %s\n", project.Path))
	
	// Status with formatting
	icon := s.GetStatusIcon(ctx, project.Status)
	color := s.GetStatusColor(ctx, project.Status)
	output.WriteString(fmt.Sprintf("Status: %s%s %s%s\n", 
		color, icon, project.Status, getColorReset()))
	
	// Services
	if len(project.Services) > 0 {
		output.WriteString("\nServices:\n")
		for _, service := range project.Services {
			serviceStatus, err := s.FormatServiceStatus(ctx, &service)
			if err != nil {
				return "", fmt.Errorf("failed to format service status: %w", err)
			}
			output.WriteString(serviceStatus + "\n")
		}
	}

	return output.String(), nil
}

// progressTracker implements ProgressTracker
type progressTracker struct {
	totalSteps   int
	currentStep  int
	isComplete   bool
	stepMessages []string
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(totalSteps int) ProgressTracker {
	return &progressTracker{
		totalSteps:   totalSteps,
		currentStep:  0,
		isComplete:   false,
		stepMessages: make([]string, totalSteps),
	}
}

// StartStep starts a new step with the given message
func (p *progressTracker) StartStep(message string) {
	if p.currentStep < p.totalSteps {
		p.stepMessages[p.currentStep] = message
		fmt.Printf("Step %d/%d: %s\n", p.currentStep+1, p.totalSteps, message)
	}
}

// CompleteStep completes the current step
func (p *progressTracker) CompleteStep() {
	if p.currentStep < p.totalSteps {
		fmt.Printf("✓ Step %d/%d completed\n", p.currentStep+1, p.totalSteps)
		p.currentStep++
		
		if p.currentStep >= p.totalSteps {
			p.isComplete = true
		}
	}
}

// FailStep marks the current step as failed
func (p *progressTracker) FailStep(err error) {
	if p.currentStep < p.totalSteps {
		fmt.Printf("✗ Step %d/%d failed: %v\n", p.currentStep+1, p.totalSteps, err)
	}
}

// GetCurrentStep returns the current step number
func (p *progressTracker) GetCurrentStep() int {
	return p.currentStep
}

// GetTotalSteps returns the total number of steps
func (p *progressTracker) GetTotalSteps() int {
	return p.totalSteps
}

// IsComplete returns whether all steps are complete
func (p *progressTracker) IsComplete() bool {
	return p.isComplete
}

// Color helper functions following UI/UX guidelines (minimal, professional)
func getGreenColor() string   { return "\033[32m" }
func getRedColor() string     { return "\033[31m" }
func getYellowColor() string  { return "\033[33m" }
func getGrayColor() string    { return "\033[90m" }
func getCyanColor() string    { return "\033[36m" }
func getDefaultColor() string { return "\033[0m" }
func getColorReset() string   { return "\033[0m" }
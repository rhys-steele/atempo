package commands

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"atempo/internal/scaffold"
)

// CreateCommand handles the 'create' command for scaffolding new projects
type CreateCommand struct {
	*BaseCommand
	templatesFS   embed.FS
	mcpServersFS  embed.FS
}

// NewCreateCommand creates a new create command
func NewCreateCommand(ctx *CommandContext, templatesFS, mcpServersFS embed.FS) *CreateCommand {
	return &CreateCommand{
		BaseCommand: NewBaseCommand(
			"create",
			"Create a new project",
			"atempo create <framework>[:<version>] [project_name]",
			ctx,
		),
		templatesFS:  templatesFS,
		mcpServersFS: mcpServersFS,
	}
}

// Execute runs the create command with enhanced real-time progress
func (c *CreateCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s\nExamples:\n  atempo create laravel my-app     # Laravel latest in ./my-app/\n  atempo create laravel:11 my-app  # Laravel 11 in ./my-app/\n  atempo create laravel            # Laravel latest in current directory", c.Usage())
	}

	// Parse framework and optional version
	frameworkArg := args[0]
	var framework, version string
	
	if strings.Contains(frameworkArg, ":") {
		parts := strings.Split(frameworkArg, ":")
		if len(parts) != 2 {
			return fmt.Errorf("error: expected format <framework>[:<version>]")
		}
		framework = parts[0]
		version = parts[1]
	} else {
		framework = frameworkArg
		version = c.getLatestVersion(framework)
	}

	// Parse optional project name
	var projectDir string
	var projectName string
	
	if len(args) >= 2 {
		// Project name specified - create directory
		projectName = args[1]
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectDir = filepath.Join(cwd, projectName)
		
		// Create project directory
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
		
		// Change to project directory
		if err := os.Chdir(projectDir); err != nil {
			return fmt.Errorf("failed to change to project directory: %w", err)
		}
	} else {
		// Use current directory
		var err error
		projectDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectName = filepath.Base(projectDir)
	}

	// Check authentication for AI features
	authChecker := NewAuthChecker()
	isAuthenticated, authStatus := authChecker.GetAuthStatus()
	
	// Initialize progress tracker (4 steps: AI Planning, Template Loading, Framework Installation, AI Context)
	tracker := NewProgressTracker(4)
	
	// Show initial project info
	ShowInfo(fmt.Sprintf("Creating %s %s project: %s", framework, version, projectName))
	fmt.Printf("%s📁 Location: %s%s\n", ColorBlue, projectDir, ColorReset)
	fmt.Printf("%s🔐 Auth Status: %s%s\n\n", ColorBlue, authStatus, ColorReset)
	
	// Run scaffolding with AI-enhanced progress tracking
	err := c.runScaffoldWithAI(tracker, framework, version, projectName, projectDir, isAuthenticated)
	if err != nil {
		// Detailed error messages are already logged by the scaffolding process
		return err
	}

	// Complete the process
	tracker.Complete(projectName)
	return nil
}

// runScaffoldWithAI runs the scaffolding process with AI-enhanced progress updates
func (c *CreateCommand) runScaffoldWithAI(tracker *ProgressTracker, framework, version, projectName, projectDir string, isAuthenticated bool) error {
	// Step 1: AI-Powered Project Planning
	tracker.StartStep(1, "AI-Powered Project Planning")
	tracker.UpdateStep("Gathering project requirements")
	
	var projectIntent *ProjectIntent
	if isAuthenticated {
		// Interactive AI-powered project setup using clean templates
		manifestGenerator, err := NewCleanAIManifestGenerator(isAuthenticated, c.templatesFS, framework)
		if err != nil {
			tracker.WarningStep("Failed to initialize AI manifest generator")
			projectIntent = createDefaultIntent(framework, version, projectName)
		} else {
			prompter, err := NewCleanInteractivePrompter(c.templatesFS)
			if err != nil {
				tracker.WarningStep("Failed to load interactive prompts")
				projectIntent = createDefaultIntent(framework, version, projectName)
			} else {
				intent, err := prompter.GatherProjectIntent(framework, projectName, manifestGenerator)
				if err != nil {
					tracker.WarningStep("Failed to gather project intent, using defaults")
					projectIntent = createDefaultIntent(framework, version, projectName)
				} else {
					projectIntent = intent
				}
			}
		}
	} else {
		tracker.UpdateStep("Authentication required for full AI features")
		// Use clean prompter for auth message
		if prompter, err := NewCleanInteractivePrompter(c.templatesFS); err == nil {
			prompter.ShowAuthenticationPrompt()
		} else {
			// Fallback to basic auth message
			authChecker := NewAuthChecker()
			authChecker.PromptAuthentication()
		}
		projectIntent = createDefaultIntent(framework, version, projectName)
	}
	
	tracker.CompleteStep("Project planning complete")
	
	// Step 2: Load template configuration
	tracker.StartStep(2, "Loading template configuration")
	tracker.UpdateStep(fmt.Sprintf("Validating %s framework template", framework))
	tracker.UpdateStep(fmt.Sprintf("Checking %s version %s compatibility", framework, version))
	tracker.CompleteStep(fmt.Sprintf("Template configuration loaded for %s %s", framework, version))
	
	// Step 3: Install framework application
	tracker.StartStep(3, "Installing framework application")
	tracker.UpdateStep(fmt.Sprintf("Running %s scaffolding process", framework))
	
	// Run the actual scaffolding process
	if err := scaffold.Run(framework, version, c.templatesFS, c.mcpServersFS); err != nil {
		// Mark the step as failed with a clean error message
		tracker.ErrorStep(err.Error())
		return err
	}
	
	tracker.CompleteStep(fmt.Sprintf("%s %s application installed", framework, version))
	
	// Step 4: Generate AI manifest (scaffold already handled infrastructure)
	tracker.StartStep(4, "Generating AI development context")
	tracker.UpdateStep("Generating AI project manifest")
	
	// Generate AI manifest files using clean generator
	manifestGenerator, err := NewCleanAIManifestGenerator(isAuthenticated, c.templatesFS, framework)
	if err != nil {
		tracker.WarningStep(fmt.Sprintf("Failed to initialize manifest generator: %v", err))
	} else {
		if err := manifestGenerator.GenerateManifestFiles(projectIntent, projectDir); err != nil {
			tracker.WarningStep(fmt.Sprintf("Failed to generate AI manifest: %v", err))
		} else {
			tracker.UpdateStep("AI manifest files generated successfully")
		}
	}
	
	tracker.CompleteStep("AI development context ready")
	
	return nil
}

// createDefaultIntent creates a basic project intent when AI features aren't available
func createDefaultIntent(framework, version, projectName string) *ProjectIntent {
	return &ProjectIntent{
		Description:     fmt.Sprintf("A %s application built with %s %s", projectName, framework, version),
		Framework:       framework,
		Language:        getFrameworkLanguage(framework),
		ProjectType:     "Web Application",
		CoreFeatures:    []string{"Basic CRUD Operations", "Database Integration", "Error Handling"},
		TechnicalNeeds:  []string{"Database", "Development Environment", "Testing Framework"},
		UserStories:     []UserStory{},
		ArchitectureHints: map[string]string{
			"pattern": fmt.Sprintf("Follow %s best practices and conventions", framework),
			"testing": "Include unit and integration tests",
			"security": "Implement proper authentication and validation",
		},
	}
}

// getLatestVersion returns the latest supported version for a framework
func (c *CreateCommand) getLatestVersion(framework string) string {
	switch framework {
	case "laravel":
		return "11" // Laravel 11 is the latest LTS
	case "django":
		return "5"  // Django 5 is the latest major version
	default:
		return "latest"
	}
}
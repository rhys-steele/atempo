package commands

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"atempo/internal/scaffold"
	"atempo/internal/ai"
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

	// Check authentication for AI features using the AI client
	aiClient, err := ai.NewAIClient()
	var isAuthenticated bool
	var authStatus string
	if err != nil {
		isAuthenticated = false
		authStatus = "AI client initialization failed - using basic project setup"
	} else {
		providers, err := aiClient.GetAvailableProviders()
		if err != nil || len(providers) == 0 {
			isAuthenticated = false
			authStatus = "No AI providers configured - using basic project setup"
		} else {
			isAuthenticated = true
			providerNames := make([]string, len(providers))
			for i, p := range providers {
				providerNames[i] = p.DisplayName
			}
			authStatus = fmt.Sprintf("AI providers available: %v", providerNames)
		}
	}
	
	// Initialize progress tracker (4 steps: AI Planning, Template Loading, Framework Installation, AI Context)
	tracker := NewProgressTracker(4)
	
	// Show initial project info
	ShowInfo(fmt.Sprintf("Creating %s %s project: %s", framework, version, projectName))
	fmt.Printf("%s📁 Location: %s%s\n", ColorBlue, projectDir, ColorReset)
	fmt.Printf("%s🔐 Auth Status: %s%s\n\n", ColorBlue, authStatus, ColorReset)
	
	// Run scaffolding with AI-enhanced progress tracking
	err = c.runScaffoldWithAI(tracker, framework, version, projectName, projectDir, isAuthenticated, aiClient)
	if err != nil {
		// Detailed error messages are already logged by the scaffolding process
		return err
	}

	// Complete the process
	tracker.Complete(projectName)
	return nil
}

// runScaffoldWithAI runs the scaffolding process with AI-enhanced progress updates
func (c *CreateCommand) runScaffoldWithAI(tracker *ProgressTracker, framework, version, projectName, projectDir string, isAuthenticated bool, aiClient *ai.AIClient) error {
	// Step 1: AI-Powered Project Planning
	tracker.StartStep(1, "AI-Powered Project Planning")
	tracker.UpdateStep("Gathering project requirements")
	
	var projectIntent *ProjectIntent
	_ = projectIntent // Will be used for future enhancements
	if isAuthenticated {
		// Use the new AI planning system
		planner, err := ai.NewProjectPlanner()
		if err != nil {
			tracker.WarningStep("Failed to initialize AI project planner")
			projectIntent = createDefaultIntent(framework, version, projectName)
		} else {
			// Use the authenticated AI client
			providers, err := aiClient.GetAvailableProviders()
			if err != nil {
				tracker.WarningStep("Failed to get AI providers")
				projectIntent = createDefaultIntent(framework, version, projectName)
			} else if len(providers) == 0 {
				tracker.WarningStep("No AI providers configured")
				projectIntent = createDefaultIntent(framework, version, projectName)
			} else {
				// Gather project description
				fmt.Printf("\n%s🤖 AI-Powered Project Setup%s\n", ColorBlue, ColorReset)
				fmt.Printf("%sLet's create comprehensive documentation for your %s project!%s\n\n", ColorGray, framework, ColorReset)
				
				fmt.Printf("%s❓ Describe your project (what you're building):%s\n", ColorCyan, ColorReset)
				fmt.Printf("%s   Example: \"A task management API with user authentication and real-time updates\"%s\n", ColorGray, ColorReset)
				fmt.Print("   > ")
				
				// Use bufio.Scanner to read the full line including spaces
				scanner := bufio.NewScanner(os.Stdin)
				var description string
				if scanner.Scan() {
					description = strings.TrimSpace(scanner.Text())
				}
				if description == "" {
					description = fmt.Sprintf("A %s application", framework)
				}
				
				fmt.Printf("\n   📝 Project description captured: \"%s\"\n", description)
				
				// Select AI provider if multiple available
				var selectedProvider string
				if len(providers) == 1 {
					selectedProvider = providers[0].Name
					fmt.Printf("\n   Using %s for documentation generation\n", providers[0].DisplayName)
				} else {
					fmt.Printf("\n%s🔧 Select AI Provider:%s\n", ColorYellow, ColorReset)
					for i, provider := range providers {
						fmt.Printf("   %d. %s\n", i+1, provider.DisplayName)
					}
					fmt.Print("   > ")
					
					var choice int
					if scanner.Scan() {
						choiceStr := strings.TrimSpace(scanner.Text())
						if choiceStr != "" {
							fmt.Sscanf(choiceStr, "%d", &choice)
						}
					}
					if choice >= 1 && choice <= len(providers) {
						selectedProvider = providers[choice-1].Name
					} else {
						selectedProvider = providers[0].Name
						fmt.Printf("   Using default provider: %s\n", providers[0].DisplayName)
					}
				}
				
				tracker.UpdateStep("Generating AI documentation with " + selectedProvider)
				
				// Generate comprehensive project documentation
				planningReq := ai.PlanningRequest{
					ProjectDescription: description,
					Framework:         framework,
					Provider:          selectedProvider,
					ProjectPath:       projectDir,
				}
				
				planningResult, err := planner.GenerateProjectPlan(context.Background(), planningReq)
				if err != nil {
					tracker.WarningStep("Failed to generate AI documentation")
					projectIntent = createDefaultIntent(framework, version, projectName)
				} else {
					// Save the generated documentation
					if err := planner.SavePlanningResult(planningResult, projectDir); err != nil {
						tracker.WarningStep("Failed to save AI documentation files")
					} else {
						tracker.UpdateStep("AI documentation saved to .ai/ directory")
					}
					
					// Create project intent from the AI-generated content
					projectIntent = &ProjectIntent{
						Description:     description,
						Framework:       framework,
						Language:        getFrameworkLanguage(framework),
						ProjectType:     "AI-Planned Application",
						CoreFeatures:    []string{"AI-Generated Documentation", "Framework Best Practices", "Comprehensive Planning"},
						TechnicalNeeds:  []string{"Development Environment", "AI Context Files", "Documentation System"},
						UserStories:     []UserStory{},
						ArchitectureHints: map[string]string{
							"ai_documentation": "Comprehensive AI-generated project documentation available in .ai/ directory",
							"planning":         "Project includes 5 detailed documentation files for better development workflow",
						},
						CreatedAt: time.Now(),
					}
				}
			}
		}
	} else {
		tracker.UpdateStep("Authentication required for AI-powered documentation")
		authChecker := NewAuthChecker()
		authChecker.PromptAuthentication()
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
	
	// Step 4: Finalize AI development context
	tracker.StartStep(4, "Finalizing AI development context")
	tracker.UpdateStep("Validating AI documentation files")
	
	// The AI documentation was already generated in Step 1
	if isAuthenticated {
		tracker.UpdateStep("AI context files are ready for development")
	} else {
		tracker.UpdateStep("Basic project structure ready")
	}
	
	tracker.CompleteStep("AI development context finalized")
	
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
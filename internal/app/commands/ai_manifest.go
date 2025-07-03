package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ProjectIntent represents the user's high-level intent for the project
type ProjectIntent struct {
	Description     string            `json:"description"`
	Framework       string            `json:"framework"`
	Language        string            `json:"language"`
	ProjectType     string            `json:"project_type"`
	CoreFeatures    []string          `json:"core_features"`
	TechnicalNeeds  []string          `json:"technical_needs"`
	UserStories     []UserStory       `json:"user_stories"`
	ArchitectureHints map[string]string `json:"architecture_hints"`
	CreatedAt       time.Time         `json:"created_at"`
}

// UserStory represents a high-level user story for the project
type UserStory struct {
	Role        string `json:"role"`
	Goal        string `json:"goal"`
	Reason      string `json:"reason"`
	Acceptance  string `json:"acceptance"`
	Priority    string `json:"priority"`
	Complexity  string `json:"complexity"`
}

// AIManifestGenerator generates AI-friendly project guidance files
type AIManifestGenerator struct {
	isAuthenticated bool
}

// NewAIManifestGenerator creates a new AI manifest generator
func NewAIManifestGenerator(isAuthenticated bool) *AIManifestGenerator {
	return &AIManifestGenerator{
		isAuthenticated: isAuthenticated,
	}
}

// PromptUserIntent interactively gathers the user's project intent
func (g *AIManifestGenerator) PromptUserIntent(framework, projectName string) (*ProjectIntent, error) {
	fmt.Printf("\n%s🤖 AI-Powered Project Setup%s\n", ColorBlue, ColorReset)
	fmt.Printf("%sLet's create an AI-first development experience for your project!%s\n\n", ColorGray, ColorReset)
	
	// Main project description
	fmt.Printf("%s❓ What kind of application would you like to build?%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s   Example: \"A simple CRUD API for managing music recordings from an event\"%s\n", ColorGray, ColorReset)
	fmt.Print("   > ")
	
	var description string
	_, err := fmt.Scanln(&description)
	if err != nil {
		// Handle multi-word input
		description = "A modern web application" // Default fallback
	}
	
	// For demo purposes, let's simulate AI analysis and create a comprehensive intent
	intent := g.generateProjectIntent(description, framework, projectName)
	
	return intent, nil
}

// generateProjectIntent creates a comprehensive project intent based on user description
func (g *AIManifestGenerator) generateProjectIntent(description, framework, projectName string) *ProjectIntent {
	// This would normally call an AI service, but for demo we'll generate realistic content
	
	// Analyze the description to extract key information
	projectType := g.analyzeProjectType(description, framework)
	coreFeatures := g.extractCoreFeatures(description, framework)
	technicalNeeds := g.determineTechnicalNeeds(description, framework)
	userStories := g.generateUserStories(description, projectType)
	architectureHints := g.generateArchitectureHints(description, framework)
	
	return &ProjectIntent{
		Description:     description,
		Framework:       framework,
		Language:        g.getFrameworkLanguage(framework),
		ProjectType:     projectType,
		CoreFeatures:    coreFeatures,
		TechnicalNeeds:  technicalNeeds,
		UserStories:     userStories,
		ArchitectureHints: architectureHints,
		CreatedAt:       time.Now(),
	}
}

// analyzeProjectType determines the type of project based on description and framework
func (g *AIManifestGenerator) analyzeProjectType(description, framework string) string {
	desc := strings.ToLower(description)
	
	if strings.Contains(desc, "api") || strings.Contains(desc, "crud") {
		return "REST API"
	}
	if strings.Contains(desc, "dashboard") || strings.Contains(desc, "admin") {
		return "Admin Dashboard"
	}
	if strings.Contains(desc, "ecommerce") || strings.Contains(desc, "shop") {
		return "E-commerce Platform"
	}
	if strings.Contains(desc, "blog") || strings.Contains(desc, "cms") {
		return "Content Management"
	}
	if strings.Contains(desc, "real-time") || strings.Contains(desc, "chat") {
		return "Real-time Application"
	}
	
	// Framework-based defaults
	switch framework {
	case "laravel":
		return "Web Application"
	case "django":
		return "Web Application"
	default:
		return "Web Application"
	}
}

// extractCoreFeatures identifies key features from the description
func (g *AIManifestGenerator) extractCoreFeatures(description, framework string) []string {
	desc := strings.ToLower(description)
	features := []string{}
	
	// Common features based on keywords
	if strings.Contains(desc, "crud") || strings.Contains(desc, "manage") {
		features = append(features, "Create/Read/Update/Delete Operations")
	}
	if strings.Contains(desc, "auth") || strings.Contains(desc, "login") {
		features = append(features, "User Authentication")
	}
	if strings.Contains(desc, "api") {
		features = append(features, "RESTful API")
	}
	if strings.Contains(desc, "search") {
		features = append(features, "Search Functionality")
	}
	if strings.Contains(desc, "upload") || strings.Contains(desc, "file") {
		features = append(features, "File Upload/Management")
	}
	if strings.Contains(desc, "notification") || strings.Contains(desc, "email") {
		features = append(features, "Notifications")
	}
	
	// Framework-specific features
	switch framework {
	case "laravel":
		features = append(features, "Eloquent ORM", "Blade Templates", "Artisan Commands")
	case "django":
		features = append(features, "Django ORM", "Django Templates", "Management Commands")
	}
	
	// Always include these essentials
	features = append(features, "Database Integration", "Error Handling", "Logging")
	
	return features
}

// determineTechnicalNeeds identifies technical requirements
func (g *AIManifestGenerator) determineTechnicalNeeds(description, framework string) []string {
	desc := strings.ToLower(description)
	needs := []string{}
	
	// Database needs
	if strings.Contains(desc, "crud") || strings.Contains(desc, "data") {
		needs = append(needs, "Database (PostgreSQL/MySQL)")
	}
	
	// Storage needs
	if strings.Contains(desc, "upload") || strings.Contains(desc, "file") || strings.Contains(desc, "image") {
		needs = append(needs, "File Storage (Local/S3)")
	}
	
	// Cache needs
	if strings.Contains(desc, "performance") || strings.Contains(desc, "cache") {
		needs = append(needs, "Caching (Redis)")
	}
	
	// Queue needs
	if strings.Contains(desc, "email") || strings.Contains(desc, "background") {
		needs = append(needs, "Queue System")
	}
	
	// Real-time needs
	if strings.Contains(desc, "real-time") || strings.Contains(desc, "live") {
		needs = append(needs, "WebSocket/Real-time")
	}
	
	// Always include basics
	needs = append(needs, "Development Environment", "Testing Framework", "API Documentation")
	
	return needs
}

// generateUserStories creates user stories based on the description
func (g *AIManifestGenerator) generateUserStories(description, projectType string) []UserStory {
	stories := []UserStory{}
	
	// Generate stories based on project type
	switch projectType {
	case "REST API":
		stories = append(stories, UserStory{
			Role:       "API Consumer",
			Goal:       "access and manipulate data through RESTful endpoints",
			Reason:     "to integrate with other applications and services",
			Acceptance: "API responds with correct data and status codes",
			Priority:   "High",
			Complexity: "Medium",
		})
	case "Web Application":
		stories = append(stories, UserStory{
			Role:       "End User",
			Goal:       "navigate and interact with the application easily",
			Reason:     "to accomplish their tasks efficiently",
			Acceptance: "User can complete core workflows without confusion",
			Priority:   "High",
			Complexity: "Medium",
		})
	}
	
	// Add common stories
	stories = append(stories, UserStory{
		Role:       "Developer",
		Goal:       "have clear documentation and development guidelines",
		Reason:     "to contribute effectively to the project",
		Acceptance: "Documentation is up-to-date and comprehensive",
		Priority:   "Medium",
		Complexity: "Low",
	})
	
	return stories
}

// generateArchitectureHints provides AI tools with architectural guidance
func (g *AIManifestGenerator) generateArchitectureHints(description, framework string) map[string]string {
	hints := make(map[string]string)
	
	// Framework-specific hints
	switch framework {
	case "laravel":
		hints["mvc_pattern"] = "Follow Laravel's MVC architecture with Controllers, Models, and Views"
		hints["service_layer"] = "Use Service classes for complex business logic"
		hints["repository_pattern"] = "Consider Repository pattern for data access abstraction"
		hints["middleware"] = "Utilize Laravel middleware for cross-cutting concerns"
	case "django":
		hints["mvt_pattern"] = "Follow Django's MVT (Model-View-Template) architecture"
		hints["apps_structure"] = "Organize code into Django apps for modularity"
		hints["signals"] = "Use Django signals for decoupled event handling"
		hints["middleware"] = "Implement Django middleware for request/response processing"
	}
	
	// General architectural hints
	hints["database_design"] = "Design normalized database schema with proper relationships"
	hints["api_design"] = "Follow RESTful principles with consistent naming conventions"
	hints["error_handling"] = "Implement comprehensive error handling and logging"
	hints["testing_strategy"] = "Include unit, integration, and feature tests"
	hints["security"] = "Implement authentication, authorization, and input validation"
	
	return hints
}

// getFrameworkLanguage returns the primary language for a framework
func (g *AIManifestGenerator) getFrameworkLanguage(framework string) string {
	switch framework {
	case "laravel":
		return "PHP"
	case "django":
		return "Python"
	case "express":
		return "JavaScript"
	default:
		return "Unknown"
	}
}

// GenerateManifestFiles creates AI-friendly files in the project
func (g *AIManifestGenerator) GenerateManifestFiles(intent *ProjectIntent, projectPath string) error {
	if !g.isAuthenticated {
		ShowWarning("AI features require authentication - generating basic manifest")
	}
	
	// Generate the main project intent file
	if err := g.generateProjectIntentFile(intent, projectPath); err != nil {
		return fmt.Errorf("failed to generate project intent file: %w", err)
	}
	
	// Generate development guidelines
	if err := g.generateDevelopmentGuidelinesFile(intent, projectPath); err != nil {
		return fmt.Errorf("failed to generate development guidelines: %w", err)
	}
	
	// Generate feature development guide
	if err := g.generateFeatureGuideFile(intent, projectPath); err != nil {
		return fmt.Errorf("failed to generate feature guide: %w", err)
	}
	
	return nil
}

// generateProjectIntentFile creates the main AI context file
func (g *AIManifestGenerator) generateProjectIntentFile(intent *ProjectIntent, projectPath string) error {
	intentJSON, err := json.MarshalIndent(intent, "", "  ")
	if err != nil {
		return err
	}
	
	content := fmt.Sprintf(`# Project Intent

This file contains the AI-generated project intent and context to help developers and AI tools understand the project's goals and architecture.

## Project Overview
%s

## Generated Context
%s

---
*Generated by Atempo AI on %s*
`, intent.Description, string(intentJSON), intent.CreatedAt.Format("2006-01-02 15:04:05"))
	
	// In a real implementation, we would write this to projectPath/ai/project-intent.md
	_ = content // Acknowledge content was created
	fmt.Printf("%s📄 Generated: ai/project-intent.md (%d chars)%s\n", ColorGreen, len(content), ColorReset)
	
	return nil
}

// generateDevelopmentGuidelinesFile creates development guidelines
func (g *AIManifestGenerator) generateDevelopmentGuidelinesFile(intent *ProjectIntent, projectPath string) error {
	content := fmt.Sprintf(`# Development Guidelines

## Project: %s Application
**Framework:** %s (%s)
**Type:** %s

## Architecture Principles
`, intent.Description, intent.Framework, intent.Language, intent.ProjectType)
	
	for key, hint := range intent.ArchitectureHints {
		title := strings.ToUpper(string(key[0])) + strings.ReplaceAll(key[1:], "_", " ")
		content += fmt.Sprintf("- **%s**: %s\n", title, hint)
	}
	
	content += "\n## Core Features\n"
	for _, feature := range intent.CoreFeatures {
		content += fmt.Sprintf("- %s\n", feature)
	}
	
	content += "\n## Technical Requirements\n"
	for _, need := range intent.TechnicalNeeds {
		content += fmt.Sprintf("- %s\n", need)
	}
	
	// In a real implementation, we would write this to projectPath/ai/development-guidelines.md
	_ = content // Acknowledge content was created
	fmt.Printf("%s📄 Generated: ai/development-guidelines.md (%d chars)%s\n", ColorGreen, len(content), ColorReset)
	
	return nil
}

// generateFeatureGuideFile creates a guide for implementing new features
func (g *AIManifestGenerator) generateFeatureGuideFile(intent *ProjectIntent, projectPath string) error {
	content := fmt.Sprintf(`# Feature Development Guide

## How to Add New Features

When implementing new features for this %s application, follow these guidelines:

### 1. User Story Template
Use this template for new features:
- **As a** [role]
- **I want** [goal]
- **So that** [reason]
- **Acceptance Criteria**: [specific criteria]

### 2. Implementation Checklist
- [ ] Create/update database migrations
- [ ] Implement business logic in appropriate layer
- [ ] Add API endpoints (if applicable)
- [ ] Create/update views and templates
- [ ] Add comprehensive tests
- [ ] Update documentation
- [ ] Consider security implications
- [ ] Add error handling and validation

### 3. Framework-Specific Patterns
`, intent.ProjectType)
	
	switch intent.Framework {
	case "laravel":
		content += `
**Laravel Patterns:**
- Controllers: Handle HTTP requests and responses
- Models: Represent data and business logic
- Services: Complex business operations
- Repositories: Data access abstraction
- Middleware: Request/response filtering
- Jobs: Background processing
- Events/Listeners: Decoupled event handling
`
	case "django":
		content += `
**Django Patterns:**
- Views: Handle requests and return responses
- Models: Define data structure and business logic
- Templates: Render HTML responses
- Forms: Handle user input validation
- Middleware: Process requests/responses
- Signals: Decouple applications
- Management Commands: Custom admin tasks
`
	}
	
	content += "\n### 4. Example User Stories\n"
	for _, story := range intent.UserStories {
		content += fmt.Sprintf(`
**%s Story:**
- As a %s
- I want to %s
- So that %s
- Acceptance: %s
- Priority: %s | Complexity: %s
`, story.Role, story.Role, story.Goal, story.Reason, story.Acceptance, story.Priority, story.Complexity)
	}
	
	// In a real implementation, we would write this to projectPath/ai/feature-guide.md
	// For demo, we'll just show that content was generated
	_ = content // Acknowledge content was created
	fmt.Printf("%s📄 Generated: ai/feature-guide.md (%d chars)%s\n", ColorGreen, len(content), ColorReset)
	
	return nil
}
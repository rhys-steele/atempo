package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// InteractivePrompter handles sophisticated user input gathering
type InteractivePrompter struct {
	scanner *bufio.Scanner
}

// NewInteractivePrompter creates a new interactive prompter
func NewInteractivePrompter() *InteractivePrompter {
	return &InteractivePrompter{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// GatherProjectIntent collects comprehensive project information from the user
func (p *InteractivePrompter) GatherProjectIntent(framework, projectName string) (*ProjectIntent, error) {
	fmt.Printf("\n%s🚀 AI-Powered Project Setup%s\n", ColorBlue, ColorReset)
	fmt.Printf("═══════════════════════════════════════════════════════════\n")
	fmt.Printf("%sLet's create an intelligent, AI-first development experience!%s\n\n", ColorCyan, ColorReset)
	
	// Main project description
	description := p.promptProjectDescription()
	
	// Gather additional context
	projectType := p.promptProjectType(description, framework)
	additionalFeatures := p.promptAdditionalFeatures(framework)
	complexity := p.promptComplexity()
	
	// Generate AI-powered intent
	fmt.Printf("\n%s🤖 Generating AI project manifest...%s\n", ColorBlue, ColorReset)
	
	intent := &ProjectIntent{
		Description:     description,
		Framework:       framework,
		Language:        getFrameworkLanguage(framework),
		ProjectType:     projectType,
		CoreFeatures:    generateSmartFeatures(description, framework, additionalFeatures),
		TechnicalNeeds:  generateTechnicalNeeds(description, framework, complexity),
		UserStories:     generateSmartUserStories(description, projectType),
		ArchitectureHints: generateArchitectureHints(framework, complexity),
	}
	
	// Show preview
	p.showIntentPreview(intent)
	
	return intent, nil
}

// promptProjectDescription gets the main project description
func (p *InteractivePrompter) promptProjectDescription() string {
	fmt.Printf("%s❓ What kind of application would you like to build?%s\n", ColorYellow, ColorReset)
	fmt.Printf("%s   Be as descriptive as possible - this helps AI understand your vision%s\n", ColorGray, ColorReset)
	fmt.Println()
	fmt.Printf("%s   Examples:%s\n", ColorGray, ColorReset)
	fmt.Printf("%s   • \"A CRUD API for managing music recordings from live events\"%s\n", ColorGray, ColorReset)
	fmt.Printf("%s   • \"An e-commerce platform for selling handmade crafts\"%s\n", ColorGray, ColorReset)
	fmt.Printf("%s   • \"A project management dashboard for small teams\"%s\n", ColorGray, ColorReset)
	fmt.Println()
	fmt.Printf("   %s>%s ", ColorCyan, ColorReset)
	
	if p.scanner.Scan() {
		description := strings.TrimSpace(p.scanner.Text())
		if description == "" {
			return "A modern web application"
		}
		return description
	}
	
	return "A modern web application"
}

// promptProjectType asks for clarification on project type
func (p *InteractivePrompter) promptProjectType(description, framework string) string {
	suggestedType := analyzeProjectTypeFromDescription(description, framework)
	
	fmt.Printf("\n%s🎯 Project Type Classification%s\n", ColorBlue, ColorReset)
	fmt.Printf("   Based on your description, this looks like: %s%s%s\n", ColorGreen, suggestedType, ColorReset)
	fmt.Printf("   Is this correct? %s(y/n, or specify different type)%s: ", ColorGray, ColorReset)
	
	if p.scanner.Scan() {
		response := strings.TrimSpace(strings.ToLower(p.scanner.Text()))
		if response == "n" || response == "no" {
			fmt.Printf("   What type of project is this?: ")
			if p.scanner.Scan() {
				customType := strings.TrimSpace(p.scanner.Text())
				if customType != "" {
					return customType
				}
			}
		}
	}
	
	return suggestedType
}

// promptAdditionalFeatures asks about specific features
func (p *InteractivePrompter) promptAdditionalFeatures(framework string) []string {
	fmt.Printf("\n%s⚡ Additional Features%s\n", ColorYellow, ColorReset)
	fmt.Printf("   Would you like to include any of these common features?\n")
	fmt.Printf("   %s(Enter numbers separated by commas, or 'none')%s\n\n", ColorGray, ColorReset)
	
	options := []string{
		"User Authentication & Authorization",
		"File Upload & Management", 
		"Email Notifications",
		"Real-time Features (WebSockets)",
		"Search Functionality",
		"Admin Dashboard",
		"API Documentation (Swagger/OpenAPI)",
		"Background Job Processing",
		"Caching Layer",
		"Multi-language Support",
	}
	
	for i, option := range options {
		fmt.Printf("   %s%d.%s %s\n", ColorCyan, i+1, ColorReset, option)
	}
	
	fmt.Printf("\n   %s>%s ", ColorCyan, ColorReset)
	
	selected := []string{}
	if p.scanner.Scan() {
		input := strings.TrimSpace(p.scanner.Text())
		if input != "none" && input != "" {
			// Parse selected numbers
			for _, numStr := range strings.Split(input, ",") {
				numStr = strings.TrimSpace(numStr)
				if numStr >= "1" && numStr <= "10" {
					// Convert to actual index
					if idx := parseSimpleInt(numStr) - 1; idx >= 0 && idx < len(options) {
						selected = append(selected, options[idx])
					}
				}
			}
		}
	}
	
	return selected
}

// promptComplexity asks about project complexity
func (p *InteractivePrompter) promptComplexity() string {
	fmt.Printf("\n%s📊 Project Complexity%s\n", ColorBlue, ColorReset)
	fmt.Printf("   How complex do you expect this project to be?\n\n")
	fmt.Printf("   %s1.%s Simple - Basic CRUD, minimal features\n", ColorCyan, ColorReset)
	fmt.Printf("   %s2.%s Medium - Multiple features, some integrations\n", ColorCyan, ColorReset)
	fmt.Printf("   %s3.%s Complex - Advanced features, multiple integrations, high scalability\n", ColorCyan, ColorReset)
	fmt.Printf("\n   %s>%s ", ColorCyan, ColorReset)
	
	if p.scanner.Scan() {
		choice := strings.TrimSpace(p.scanner.Text())
		switch choice {
		case "1":
			return "Simple"
		case "2":
			return "Medium"
		case "3":
			return "Complex"
		default:
			return "Medium"
		}
	}
	
	return "Medium"
}

// showIntentPreview displays a preview of the generated project intent
func (p *InteractivePrompter) showIntentPreview(intent *ProjectIntent) {
	fmt.Printf("\n%s📋 Project Intent Summary%s\n", ColorGreen, ColorReset)
	fmt.Printf("═══════════════════════════════════════════════════════════\n")
	fmt.Printf("%sDescription:%s %s\n", ColorCyan, ColorReset, intent.Description)
	fmt.Printf("%sFramework:%s %s (%s)\n", ColorCyan, ColorReset, intent.Framework, intent.Language)
	fmt.Printf("%sType:%s %s\n", ColorCyan, ColorReset, intent.ProjectType)
	
	fmt.Printf("\n%sCore Features:%s\n", ColorCyan, ColorReset)
	for _, feature := range intent.CoreFeatures {
		fmt.Printf("  • %s\n", feature)
	}
	
	fmt.Printf("\n%sTechnical Needs:%s\n", ColorCyan, ColorReset)
	for _, need := range intent.TechnicalNeeds {
		fmt.Printf("  • %s\n", need)
	}
	
	fmt.Printf("\n%sUser Stories Generated:%s %d\n", ColorCyan, ColorReset, len(intent.UserStories))
	fmt.Printf("%sArchitecture Hints:%s %d guidelines\n", ColorCyan, ColorReset, len(intent.ArchitectureHints))
	
	fmt.Printf("\n%s✨ This manifest will help AI tools understand your project better!%s\n", ColorGreen, ColorReset)
}

// Helper functions

func analyzeProjectTypeFromDescription(description, framework string) string {
	desc := strings.ToLower(description)
	
	if strings.Contains(desc, "api") || strings.Contains(desc, "rest") || strings.Contains(desc, "endpoint") {
		return "REST API"
	}
	if strings.Contains(desc, "dashboard") || strings.Contains(desc, "admin") {
		return "Admin Dashboard"
	}
	if strings.Contains(desc, "ecommerce") || strings.Contains(desc, "shop") || strings.Contains(desc, "store") {
		return "E-commerce Platform"
	}
	if strings.Contains(desc, "blog") || strings.Contains(desc, "cms") || strings.Contains(desc, "content") {
		return "Content Management System"
	}
	if strings.Contains(desc, "real-time") || strings.Contains(desc, "chat") || strings.Contains(desc, "live") {
		return "Real-time Application"
	}
	if strings.Contains(desc, "crud") {
		return "CRUD Application"
	}
	
	return "Web Application"
}

func generateSmartFeatures(description, framework string, additionalFeatures []string) []string {
	features := []string{}
	desc := strings.ToLower(description)
	
	// Core features based on description
	if strings.Contains(desc, "crud") || strings.Contains(desc, "manage") {
		features = append(features, "Create/Read/Update/Delete Operations")
	}
	if strings.Contains(desc, "api") {
		features = append(features, "RESTful API Endpoints")
	}
	if strings.Contains(desc, "user") || strings.Contains(desc, "auth") {
		features = append(features, "User Management")
	}
	
	// Framework-specific features
	switch framework {
	case "laravel":
		features = append(features, "Eloquent ORM", "Blade Templates", "Laravel Validation")
	case "django":
		features = append(features, "Django ORM", "Django Forms", "Template System")
	}
	
	// Add user-selected features
	features = append(features, additionalFeatures...)
	
	// Always include essentials
	features = append(features, "Database Integration", "Error Handling", "Logging System")
	
	return features
}

func generateTechnicalNeeds(description, framework string, complexity string) []string {
	needs := []string{}
	desc := strings.ToLower(description)
	
	// Database
	needs = append(needs, "Database (PostgreSQL recommended)")
	
	// Based on description
	if strings.Contains(desc, "file") || strings.Contains(desc, "upload") || strings.Contains(desc, "image") {
		needs = append(needs, "File Storage Solution")
	}
	if strings.Contains(desc, "email") || strings.Contains(desc, "notification") {
		needs = append(needs, "Email Service Integration")
	}
	if strings.Contains(desc, "search") {
		needs = append(needs, "Search Engine (Elasticsearch/Database FTS)")
	}
	
	// Complexity-based needs
	switch complexity {
	case "Complex":
		needs = append(needs, "Caching Layer (Redis)", "Queue System", "Load Balancer", "CDN")
	case "Medium":
		needs = append(needs, "Caching Layer (Redis)", "Queue System")
	}
	
	// Always include
	needs = append(needs, "Development Environment", "Testing Framework", "CI/CD Pipeline")
	
	return needs
}

func generateSmartUserStories(description, projectType string) []UserStory {
	stories := []UserStory{}
	
	// Generate contextual user stories
	stories = append(stories, UserStory{
		Role:       "End User",
		Goal:       "interact with the application intuitively",
		Reason:     "to accomplish my tasks efficiently without confusion",
		Acceptance: "User can navigate and complete primary workflows",
		Priority:   "High",
		Complexity: "Medium",
	})
	
	stories = append(stories, UserStory{
		Role:       "Developer",
		Goal:       "have clear code organization and documentation",
		Reason:     "to maintain and extend the application effectively",
		Acceptance: "Code follows established patterns and is well-documented",
		Priority:   "Medium",
		Complexity: "Low",
	})
	
	// Type-specific stories
	if strings.Contains(strings.ToLower(projectType), "api") {
		stories = append(stories, UserStory{
			Role:       "API Consumer",
			Goal:       "access reliable and well-documented endpoints",
			Reason:     "to integrate with other systems seamlessly",
			Acceptance: "API returns consistent responses with proper status codes",
			Priority:   "High",
			Complexity: "Medium",
		})
	}
	
	return stories
}

func generateArchitectureHints(framework, complexity string) map[string]string {
	hints := make(map[string]string)
	
	// Framework-specific hints
	switch framework {
	case "laravel":
		hints["mvc_pattern"] = "Follow Laravel MVC with Controllers handling requests, Models for data, Views for presentation"
		hints["service_layer"] = "Use Service classes for complex business logic to keep controllers thin"
		hints["validation"] = "Utilize Laravel Form Requests for input validation and authorization"
		hints["eloquent"] = "Leverage Eloquent relationships and query optimization"
	case "django":
		hints["mvt_pattern"] = "Follow Django MVT with Views handling logic, Models for data, Templates for presentation"
		hints["apps_structure"] = "Organize functionality into focused Django apps for better modularity"
		hints["forms"] = "Use Django Forms for data validation and HTML generation"
		hints["querysets"] = "Optimize database queries using select_related and prefetch_related"
	}
	
	// Complexity-based hints
	switch complexity {
	case "Complex":
		hints["scalability"] = "Design for horizontal scaling with stateless services"
		hints["caching"] = "Implement multi-layer caching strategy"
		hints["monitoring"] = "Include comprehensive logging and monitoring"
	case "Medium":
		hints["performance"] = "Consider database indexing and query optimization"
		hints["testing"] = "Implement unit, integration, and feature tests"
	default:
		hints["simplicity"] = "Keep architecture simple and focus on core functionality"
	}
	
	// Universal hints
	hints["security"] = "Implement proper authentication, authorization, and input sanitization"
	hints["documentation"] = "Maintain up-to-date API documentation and code comments"
	
	return hints
}

func getFrameworkLanguage(framework string) string {
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

func parseSimpleInt(s string) int {
	switch s {
	case "1": return 1
	case "2": return 2
	case "3": return 3
	case "4": return 4
	case "5": return 5
	case "6": return 6
	case "7": return 7
	case "8": return 8
	case "9": return 9
	case "10": return 10
	default: return 0
	}
}
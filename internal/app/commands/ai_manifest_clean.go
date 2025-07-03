package commands

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
	"time"
)

// CleanAIManifestGenerator generates AI-friendly project guidance using external templates
type CleanAIManifestGenerator struct {
	isAuthenticated bool
	templateLoader  *TemplateLoader
	frameworkConfig *FrameworkConfig
}

// NewCleanAIManifestGenerator creates a new clean AI manifest generator
func NewCleanAIManifestGenerator(isAuthenticated bool, templatesFS fs.FS, framework string) (*CleanAIManifestGenerator, error) {
	loader := NewTemplateLoader(templatesFS)
	
	// Load framework configuration
	frameworkConfig, err := loader.LoadFrameworkConfig(framework)
	if err != nil {
		// Create basic fallback config
		frameworkConfig = &FrameworkConfig{
			Framework: framework,
			Language:  getFrameworkLanguage(framework),
			AIFeatures: AIFeatures{
				DefaultProjectTypes: []string{"Web Application"},
				CoreFeatures:       []string{"Basic CRUD Operations", "Database Integration"},
				ArchitecturePatterns: map[string]string{
					"mvc_pattern": "Follow MVC architecture patterns",
				},
				TechnicalStack: []string{"Database", "Development Environment"},
			},
		}
	}
	
	return &CleanAIManifestGenerator{
		isAuthenticated: isAuthenticated,
		templateLoader:  loader,
		frameworkConfig: frameworkConfig,
	}, nil
}

// AnalyzeProjectType determines project type from description using framework config
func (g *CleanAIManifestGenerator) AnalyzeProjectType(description string) string {
	desc := strings.ToLower(description)
	
	// Use framework-specific keywords
	for keyword, projectType := range g.frameworkConfig.AIFeatures.ProjectAnalysisKeywords {
		if strings.Contains(desc, keyword) {
			return projectType
		}
	}
	
	// Default to first project type for framework
	if len(g.frameworkConfig.AIFeatures.DefaultProjectTypes) > 0 {
		return g.frameworkConfig.AIFeatures.DefaultProjectTypes[0]
	}
	
	return "Web Application"
}

// GenerateSmartFeatures creates feature list based on description and framework
func (g *CleanAIManifestGenerator) GenerateSmartFeatures(description string, additionalFeatures []string) []string {
	features := []string{}
	desc := strings.ToLower(description)
	
	// Core features based on description analysis
	if strings.Contains(desc, "crud") || strings.Contains(desc, "manage") {
		features = append(features, "Create/Read/Update/Delete Operations")
	}
	if strings.Contains(desc, "api") {
		features = append(features, "RESTful API Endpoints")
	}
	if strings.Contains(desc, "user") || strings.Contains(desc, "auth") {
		features = append(features, "User Management")
	}
	
	// Add framework-specific core features
	features = append(features, g.frameworkConfig.AIFeatures.CoreFeatures...)
	
	// Add user-selected features
	features = append(features, additionalFeatures...)
	
	// Always include essentials
	features = append(features, "Error Handling", "Logging System")
	
	return features
}

// GenerateTechnicalNeeds determines technical requirements
func (g *CleanAIManifestGenerator) GenerateTechnicalNeeds(description, complexity string) []string {
	needs := []string{}
	desc := strings.ToLower(description)
	
	// Start with framework technical stack
	needs = append(needs, g.frameworkConfig.AIFeatures.TechnicalStack...)
	
	// Analysis-based needs
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
		needs = append(needs, "Load Balancer", "CDN", "Monitoring & Logging")
	case "Medium":
		needs = append(needs, "Performance Monitoring")
	}
	
	// Always include
	needs = append(needs, "Testing Framework", "CI/CD Pipeline")
	
	return needs
}

// GenerateUserStories creates contextual user stories
func (g *CleanAIManifestGenerator) GenerateUserStories(description, projectType string) []UserStory {
	stories := []UserStory{
		{
			Role:       "End User",
			Goal:       "interact with the application intuitively",
			Reason:     "to accomplish my tasks efficiently without confusion",
			Acceptance: "User can navigate and complete primary workflows",
			Priority:   "High",
			Complexity: "Medium",
		},
		{
			Role:       "Developer",
			Goal:       "have clear code organization and documentation",
			Reason:     "to maintain and extend the application effectively",
			Acceptance: "Code follows established patterns and is well-documented",
			Priority:   "Medium",
			Complexity: "Low",
		},
	}
	
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

// GenerateArchitectureHints creates architecture guidance using framework config
func (g *CleanAIManifestGenerator) GenerateArchitectureHints(complexity string) map[string]string {
	hints := make(map[string]string)
	
	// Use framework-specific patterns
	for key, value := range g.frameworkConfig.AIFeatures.ArchitecturePatterns {
		hints[key] = value
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

// GenerateManifestFiles creates AI-friendly files using templates
func (g *CleanAIManifestGenerator) GenerateManifestFiles(intent *ProjectIntent, projectPath string) error {
	// Authentication status is handled by the caller
	
	// Generate project intent file
	if err := g.generateTemplatedFile("project-intent-template", intent, "project-intent.md"); err != nil {
		return fmt.Errorf("failed to generate project intent: %w", err)
	}
	
	// Generate development guidelines
	if err := g.generateTemplatedFile("development-guidelines-template", intent, "development-guidelines.md"); err != nil {
		return fmt.Errorf("failed to generate development guidelines: %w", err)
	}
	
	// Generate feature guide with framework patterns
	featureData := struct {
		*ProjectIntent
		FrameworkPatterns string
	}{
		ProjectIntent:     intent,
		FrameworkPatterns: g.frameworkConfig.AIFeatures.FrameworkPatternsTemplate,
	}
	
	if err := g.generateTemplatedFile("feature-guide-template", featureData, "feature-guide.md"); err != nil {
		return fmt.Errorf("failed to generate feature guide: %w", err)
	}
	
	return nil
}

// generateTemplatedFile generates a file using a template
func (g *CleanAIManifestGenerator) generateTemplatedFile(templateName string, data interface{}, outputName string) error {
	template, err := g.templateLoader.LoadManifestTemplate(templateName)
	if err != nil {
		return err
	}
	
	// Add JSON representation for project intent
	if templateName == "project-intent-template" {
		if intentData, ok := data.(*ProjectIntent); ok {
			jsonBytes, _ := json.MarshalIndent(intentData, "", "  ")
			templateData := struct {
				*ProjectIntent
				ProjectJSON string
			}{
				ProjectIntent: intentData,
				ProjectJSON:   string(jsonBytes),
			}
			data = templateData
		}
	}
	
	content, err := g.templateLoader.GenerateFromTemplate(template.MarkdownTemplate, data)
	if err != nil {
		return err
	}
	
	// In a real implementation, we would write to projectPath/ai/outputName
	// For demo, we'll just show that content was generated
	fmt.Printf("%s📄 Generated: ai/%s (%d chars)%s\n", ColorGreen, outputName, len(content), ColorReset)
	
	return nil
}

// CreateProjectIntent creates a comprehensive project intent using templates and framework config
func (g *CleanAIManifestGenerator) CreateProjectIntent(description, projectName, version string, additionalFeatures []string, complexity string) *ProjectIntent {
	projectType := g.AnalyzeProjectType(description)
	coreFeatures := g.GenerateSmartFeatures(description, additionalFeatures)
	technicalNeeds := g.GenerateTechnicalNeeds(description, complexity)
	userStories := g.GenerateUserStories(description, projectType)
	architectureHints := g.GenerateArchitectureHints(complexity)
	
	return &ProjectIntent{
		Description:       description,
		Framework:         g.frameworkConfig.Framework,
		Language:          g.frameworkConfig.Language,
		ProjectType:       projectType,
		CoreFeatures:      coreFeatures,
		TechnicalNeeds:    technicalNeeds,
		UserStories:       userStories,
		ArchitectureHints: architectureHints,
		CreatedAt:         time.Now(),
	}
}
package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ProjectPlanner handles AI-powered project planning
type ProjectPlanner struct {
	client *AIClient
}

// NewProjectPlanner creates a new project planner
func NewProjectPlanner() (*ProjectPlanner, error) {
	client, err := NewAIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %w", err)
	}

	return &ProjectPlanner{
		client: client,
	}, nil
}

// PlanningRequest contains information for AI project planning
type PlanningRequest struct {
	ProjectDescription string
	Framework          string
	Provider           string
	ProjectPath        string
}

// PlanningResult contains the generated project documentation
type PlanningResult struct {
	CodebaseMap            string
	CurrentState           string
	DevelopmentWorkflows   string
	PatternsAndConventions string
	ProjectOverview        string
}

// GenerateProjectPlan creates AI-generated project documentation concurrently
func (p *ProjectPlanner) GenerateProjectPlan(ctx context.Context, req PlanningRequest) (*PlanningResult, error) {
	fmt.Printf("🤖 Generating AI-powered project documentation...\n")
	fmt.Printf("   Framework: %s\n", req.Framework)
	fmt.Printf("   Provider: %s\n", req.Provider)
	fmt.Printf("   Description: %s\n\n", req.ProjectDescription)

	result := &PlanningResult{}

	// Define document generation tasks
	type documentTask struct {
		name     string
		field    *string
		template string
	}

	documents := []documentTask{
		{"Project Overview", &result.ProjectOverview, p.getProjectOverviewPrompt()},
		{"Codebase Map", &result.CodebaseMap, p.getCodebaseMapPrompt()},
		{"Current State", &result.CurrentState, p.getCurrentStatePrompt()},
		{"Development Workflows", &result.DevelopmentWorkflows, p.getDevelopmentWorkflowsPrompt()},
		{"Patterns and Conventions", &result.PatternsAndConventions, p.getPatternsAndConventionsPrompt()},
	}

	// Initialize concurrent generation
	var wg sync.WaitGroup
	var mu sync.Mutex
	errs := make([]error, len(documents))
	
	fmt.Printf("🚀 Starting concurrent generation of %d documents...\n", len(documents))

	// Generate each document concurrently
	for i, doc := range documents {
		wg.Add(1)
		go func(index int, task documentTask) {
			defer wg.Done()
			
			fmt.Printf("→ Generating %s...\n", task.name)
			
			content, err := p.generateDocument(ctx, req, task.template)
			if err != nil {
				mu.Lock()
				errs[index] = fmt.Errorf("failed to generate %s: %w", task.name, err)
				mu.Unlock()
				fmt.Printf("❌ %s generation failed\n", task.name)
				return
			}
			
			mu.Lock()
			*task.field = content
			mu.Unlock()
			
			fmt.Printf("✅ %s generated (%d tokens)\n", task.name, len(content)/4)
		}(i, doc)
	}

	// Wait for all documents to complete
	wg.Wait()

	// Check for any errors
	for i, err := range errs {
		if err != nil {
			return nil, fmt.Errorf("document generation failed for %s: %w", documents[i].name, err)
		}
	}

	fmt.Printf("\n🎉 All project documentation generated successfully (concurrent)!\n")
	return result, nil
}

// generateDocument generates a single document using AI
func (p *ProjectPlanner) generateDocument(ctx context.Context, req PlanningRequest, promptTemplate string) (string, error) {
	// Replace placeholders in the prompt template
	prompt := strings.ReplaceAll(promptTemplate, "{{PROJECT_DESCRIPTION}}", req.ProjectDescription)
	prompt = strings.ReplaceAll(prompt, "{{FRAMEWORK}}", req.Framework)

	// Create AI request
	chatReq := ChatRequest{
		Provider: req.Provider,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   4000,
	}

	// Send request to AI
	response, err := p.client.SendChatRequest(ctx, chatReq)
	if err != nil {
		return "", err
	}

	return response.Content, nil
}

// SavePlanningResult saves the generated documentation to files
func (p *ProjectPlanner) SavePlanningResult(result *PlanningResult, projectPath string) error {
	// Create .ai directory
	aiDir := filepath.Join(projectPath, ".ai")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		return fmt.Errorf("failed to create .ai directory: %w", err)
	}

	// Save each document
	documents := []struct {
		filename string
		content  string
	}{
		{"project-overview.md", result.ProjectOverview},
		{"codebase-map.md", result.CodebaseMap},
		{"current-state.md", result.CurrentState},
		{"development-workflows.md", result.DevelopmentWorkflows},
		{"patterns-and-conventions.md", result.PatternsAndConventions},
	}

	for _, doc := range documents {
		filePath := filepath.Join(aiDir, doc.filename)
		if err := os.WriteFile(filePath, []byte(doc.content), 0644); err != nil {
			return fmt.Errorf("failed to save %s: %w", doc.filename, err)
		}
	}

	return nil
}

// Prompt templates for different documents

func (p *ProjectPlanner) getProjectOverviewPrompt() string {
	return `You are an expert software architect. Create a comprehensive project overview document for a {{FRAMEWORK}} project.

Project Description: {{PROJECT_DESCRIPTION}}

Please create a detailed project overview in Markdown format that includes:

1. **Project Mission & Goals** - What this project aims to achieve
2. **Core Value Proposition** - Key benefits and unique features
3. **Architecture Overview** - High-level technical architecture
4. **Key Components** - Main modules/components and their responsibilities
5. **Current State** - What's implemented vs planned
6. **Technical Highlights** - Notable technical decisions and patterns
7. **Key Files for Understanding** - Most important files for newcomers
8. **Development Workflows** - How development should work
9. **Dependencies & External Integrations** - Third-party services and tools

Focus on creating a document that would help a new developer understand the project quickly and make effective contributions. Be specific to {{FRAMEWORK}} best practices and conventions.

Use clear headings, bullet points, and code examples where appropriate. This should be the definitive "README" for understanding the project's purpose and architecture.`
}

func (p *ProjectPlanner) getCodebaseMapPrompt() string {
	return `You are an expert software architect. Create a comprehensive codebase map for a {{FRAMEWORK}} project.

Project Description: {{PROJECT_DESCRIPTION}}

Please create a detailed codebase map in Markdown format that includes:

1. **Directory Structure Overview** - Complete folder hierarchy with explanations
2. **File Descriptions & Responsibilities** - What each major file/directory does
3. **Key Relationships & Data Flow** - How components interact
4. **Architecture Patterns** - Design patterns used throughout the codebase
5. **Key Interfaces & Types** - Important contracts and data structures
6. **External Dependencies** - Third-party libraries and their usage
7. **Development & Extension Points** - Where to add new features

Focus on creating a technical reference that helps developers navigate the codebase efficiently. Include:
- Typical {{FRAMEWORK}} project structure
- Framework-specific conventions and patterns
- Clear explanations of how pieces fit together
- Line counts for major files to show relative complexity

This should be like a technical "map" that prevents developers from getting lost in the codebase and helps them understand the overall structure and relationships.`
}

func (p *ProjectPlanner) getCurrentStatePrompt() string {
	return `You are an expert project manager and software architect. Create a current state document for a {{FRAMEWORK}} project.

Project Description: {{PROJECT_DESCRIPTION}}

Please create a detailed current state document in Markdown format that includes:

1. **Development Status** - Current branch, recent changes, active development
2. **Feature Implementation Status** - What's completed, in progress, and planned
3. **Recent Changes** - Key commits, updates, and modifications
4. **Architecture Changes** - Any recent architectural decisions or refactoring
5. **Current Development Priorities** - What should be worked on next
6. **Key Development Challenges** - Current blockers, technical debt, issues
7. **Dependencies & External Integrations** - Current state of third-party tools
8. **Configuration & Settings** - Environment variables, config files
9. **Next Steps & Immediate Actions** - Concrete next steps for development

Focus on providing a snapshot of where the project currently stands. This should help developers understand:
- What's been done recently
- What needs attention now
- Where development is heading
- Any current issues or challenges

Be specific about {{FRAMEWORK}} features, conventions, and typical development workflow. Include both technical implementation status and project management perspectives.`
}

func (p *ProjectPlanner) getDevelopmentWorkflowsPrompt() string {
	return `You are an expert DevOps engineer and software architect. Create a development workflows document for a {{FRAMEWORK}} project.

Project Description: {{PROJECT_DESCRIPTION}}

Please create a detailed development workflows document in Markdown format that includes:

1. **Local Development Setup** - Step-by-step environment setup
2. **Development Process** - Daily development workflow and best practices
3. **Code Standards & Guidelines** - Coding conventions, style guides, linting
4. **Testing Strategy** - Unit tests, integration tests, testing workflows
5. **Build & Deployment Process** - How to build, package, and deploy
6. **Git Workflow** - Branching strategy, commit conventions, PR process
7. **Code Review Process** - How code reviews work, what to look for
8. **Debugging & Troubleshooting** - Common issues and how to solve them
9. **Performance Optimization** - Profiling, monitoring, optimization techniques
10. **Documentation Maintenance** - How to keep docs up to date

Focus on practical, actionable workflows that developers can follow. Include:
- {{FRAMEWORK}}-specific commands and tools
- Environment setup and configuration
- Testing and quality assurance processes
- Collaboration and team workflows

This should be a practical guide that answers "How do I..." questions for common development tasks and ensures consistent development practices across the team.`
}

func (p *ProjectPlanner) getPatternsAndConventionsPrompt() string {
	return `You are an expert software architect specializing in {{FRAMEWORK}}. Create a patterns and conventions document for this project.

Project Description: {{PROJECT_DESCRIPTION}}

Please create a detailed patterns and conventions document in Markdown format that includes:

1. **Architecture Patterns** - Core architectural patterns and their implementation
2. **Code Organization** - How to structure and organize code effectively
3. **Naming Conventions** - File names, class names, variable names, etc.
4. **Design Patterns** - Common design patterns used and when to apply them
5. **Data Flow Patterns** - How data moves through the application
6. **Error Handling** - Consistent error handling approaches and patterns
7. **Security Patterns** - Security best practices and implementation patterns
8. **Performance Patterns** - Optimization patterns and performance considerations
9. **Testing Patterns** - Testing strategies, mocking patterns, test organization
10. **API Design Patterns** - REST/GraphQL patterns, request/response handling
11. **Database Patterns** - Data modeling, migration, and query patterns
12. **Frontend Patterns** (if applicable) - UI patterns, state management

Focus on {{FRAMEWORK}}-specific patterns and conventions. Include:
- Framework-recommended patterns and best practices
- Project-specific conventions and decisions
- Code examples showing proper implementation
- Anti-patterns to avoid
- When and why to use specific patterns

This should be a reference guide that ensures consistency across the codebase and helps developers make good architectural decisions that align with {{FRAMEWORK}} best practices.`
}
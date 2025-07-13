package ai

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)


//go:embed prompts/create/*.md
var createPromptsFS embed.FS

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

// Prompt templates loaded from files

func (p *ProjectPlanner) getProjectOverviewPrompt() string {
	return p.loadPromptTemplate("project-overview.md")
}

func (p *ProjectPlanner) getCodebaseMapPrompt() string {
	return p.loadPromptTemplate("codebase-map.md")
}

func (p *ProjectPlanner) getCurrentStatePrompt() string {
	return p.loadPromptTemplate("current-state.md")
}

func (p *ProjectPlanner) getDevelopmentWorkflowsPrompt() string {
	return p.loadPromptTemplate("development-workflows.md")
}

func (p *ProjectPlanner) getPatternsAndConventionsPrompt() string {
	return p.loadPromptTemplate("patterns-and-conventions.md")
}

// loadPromptTemplate loads a prompt template from the embedded filesystem
func (p *ProjectPlanner) loadPromptTemplate(filename string) string {
	path := filepath.Join("prompts/create", filename)
	content, err := createPromptsFS.ReadFile(path)
	if err != nil {
		// Fallback to a basic template if file loading fails
		return fmt.Sprintf("Create documentation for a {{FRAMEWORK}} project: {{PROJECT_DESCRIPTION}}")
	}
	return string(content)
}

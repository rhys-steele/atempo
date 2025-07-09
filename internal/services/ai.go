package services

import (
	"context"
	"fmt"
	"time"

	"atempo/internal/ai"
	"atempo/internal/auth"
)

// AIService provides business operations for AI client management and interactions
type AIService interface {
	// CreateClient creates a new AI client with the specified provider
	CreateClient(ctx context.Context, providerName string, credentials *auth.Credentials) (ai.Client, error)
	
	// GetClient retrieves an existing AI client for a provider
	GetClient(ctx context.Context, providerName string) (ai.Client, error)
	
	// SendChatRequest sends a chat request to an AI provider
	SendChatRequest(ctx context.Context, providerName string, request *ai.ChatRequest) (*ai.ChatResponse, error)
	
	// GenerateManifest generates an AI project manifest
	GenerateManifest(ctx context.Context, providerName string, request *ManifestRequest) (*ManifestResponse, error)
	
	// AnalyzeProject analyzes a project and provides AI insights
	AnalyzeProject(ctx context.Context, providerName string, request *ProjectAnalysisRequest) (*ProjectAnalysisResponse, error)
	
	// ValidateProvider validates that a provider is available and configured
	ValidateProvider(ctx context.Context, providerName string) error
	
	// ListProviders returns all available AI providers
	ListProviders(ctx context.Context) ([]AIProviderInfo, error)
	
	// RefreshClient refreshes the AI client for a provider (useful for credential refresh)
	RefreshClient(ctx context.Context, providerName string) error
}

// ManifestRequest represents a request to generate an AI project manifest
type ManifestRequest struct {
	ProjectName  string            `json:"project_name"`
	Framework    string            `json:"framework"`
	Version      string            `json:"version"`
	Description  string            `json:"description,omitempty"`
	Features     []string          `json:"features,omitempty"`
	Complexity   string            `json:"complexity,omitempty"` // "simple", "moderate", "complex"
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// ManifestResponse represents the response from AI manifest generation
type ManifestResponse struct {
	ManifestPath string                 `json:"manifest_path"`
	Content      map[string]interface{} `json:"content"`
	Suggestions  []string               `json:"suggestions,omitempty"`
	Confidence   float64                `json:"confidence"`
}

// ProjectAnalysisRequest represents a request to analyze a project
type ProjectAnalysisRequest struct {
	ProjectPath string            `json:"project_path"`
	Framework   string            `json:"framework,omitempty"`
	Focus       []string          `json:"focus,omitempty"` // areas to focus on: "security", "performance", "architecture"
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ProjectAnalysisResponse represents the response from project analysis
type ProjectAnalysisResponse struct {
	Summary      string                 `json:"summary"`
	Issues       []AnalysisIssue        `json:"issues,omitempty"`
	Suggestions  []AnalysisSuggestion   `json:"suggestions,omitempty"`
	Metrics      map[string]interface{} `json:"metrics,omitempty"`
	Confidence   float64                `json:"confidence"`
}

// AnalysisIssue represents an issue found during project analysis
type AnalysisIssue struct {
	Type        string `json:"type"`        // "security", "performance", "architecture", etc.
	Severity    string `json:"severity"`    // "low", "medium", "high", "critical"
	Title       string `json:"title"`
	Description string `json:"description"`
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
	Suggestion  string `json:"suggestion,omitempty"`
}

// AnalysisSuggestion represents a suggestion for project improvement
type AnalysisSuggestion struct {
	Category    string `json:"category"`    // "optimization", "feature", "refactor", etc.
	Priority    string `json:"priority"`    // "low", "medium", "high"
	Title       string `json:"title"`
	Description string `json:"description"`
	Effort      string `json:"effort"`      // "small", "medium", "large"
}

// AIProviderInfo represents information about an AI provider
type AIProviderInfo struct {
	Name         string   `json:"name"`
	DisplayName  string   `json:"display_name"`
	Description  string   `json:"description"`
	Capabilities []string `json:"capabilities"`
	Available    bool     `json:"available"`
	Configured   bool     `json:"configured"`
}

// aiService implements AIService
type aiService struct {
	authService AuthService
	clients     map[string]ai.Client
}

// NewAIService creates a new AIService implementation
func NewAIService(authService AuthService) AIService {
	return &aiService{
		authService: authService,
		clients:     make(map[string]ai.Client),
	}
}

// CreateClient creates a new AI client with the specified provider
func (s *aiService) CreateClient(ctx context.Context, providerName string, credentials *auth.Credentials) (ai.Client, error) {
	if credentials == nil {
		return nil, fmt.Errorf("credentials cannot be nil")
	}
	
	if !credentials.IsValid() {
		return nil, fmt.Errorf("credentials are invalid")
	}
	
	// Create AI client
	client, err := ai.NewAIClient(providerName, credentials.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %w", err)
	}
	
	// Store client for reuse
	s.clients[providerName] = client
	
	return client, nil
}

// GetClient retrieves an existing AI client for a provider
func (s *aiService) GetClient(ctx context.Context, providerName string) (ai.Client, error) {
	// Check if client already exists
	if client, exists := s.clients[providerName]; exists {
		return client, nil
	}
	
	// Try to get credentials and create client
	credentials, err := s.authService.ValidateCredentials(ctx, providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to validate credentials for provider '%s': %w", providerName, err)
	}
	
	return s.CreateClient(ctx, providerName, credentials)
}

// SendChatRequest sends a chat request to an AI provider
func (s *aiService) SendChatRequest(ctx context.Context, providerName string, request *ai.ChatRequest) (*ai.ChatResponse, error) {
	client, err := s.GetClient(ctx, providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client: %w", err)
	}
	
	response, err := client.SendChatRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send chat request: %w", err)
	}
	
	return response, nil
}

// GenerateManifest generates an AI project manifest
func (s *aiService) GenerateManifest(ctx context.Context, providerName string, request *ManifestRequest) (*ManifestResponse, error) {
	// Prepare chat request for manifest generation
	chatRequest := &ai.ChatRequest{
		Messages: []ai.Message{
			{
				Role: "system",
				Content: "You are an expert software architect. Generate a comprehensive project manifest for the given requirements.",
			},
			{
				Role: "user",
				Content: s.buildManifestPrompt(request),
			},
		},
		MaxTokens:   2000,
		Temperature: 0.7,
	}
	
	response, err := s.SendChatRequest(ctx, providerName, chatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate manifest: %w", err)
	}
	
	// Process response into manifest format
	manifestResponse := &ManifestResponse{
		ManifestPath: fmt.Sprintf("manifest-%s-%d.json", request.ProjectName, time.Now().Unix()),
		Content: map[string]interface{}{
			"project_name": request.ProjectName,
			"framework":    request.Framework,
			"version":      request.Version,
			"generated_by": providerName,
			"generated_at": time.Now().UTC(),
			"ai_response":  response.Choices[0].Message.Content,
		},
		Confidence: 0.8, // Default confidence
	}
	
	return manifestResponse, nil
}

// AnalyzeProject analyzes a project and provides AI insights
func (s *aiService) AnalyzeProject(ctx context.Context, providerName string, request *ProjectAnalysisRequest) (*ProjectAnalysisResponse, error) {
	// Prepare chat request for project analysis
	chatRequest := &ai.ChatRequest{
		Messages: []ai.Message{
			{
				Role: "system",
				Content: "You are an expert software architect and code reviewer. Analyze the given project and provide insights.",
			},
			{
				Role: "user",
				Content: s.buildAnalysisPrompt(request),
			},
		},
		MaxTokens:   3000,
		Temperature: 0.3, // Lower temperature for more focused analysis
	}
	
	response, err := s.SendChatRequest(ctx, providerName, chatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze project: %w", err)
	}
	
	// Process response into analysis format
	analysisResponse := &ProjectAnalysisResponse{
		Summary:    response.Choices[0].Message.Content,
		Confidence: 0.7, // Default confidence
		Metrics: map[string]interface{}{
			"analyzed_at": time.Now().UTC(),
			"provider":    providerName,
		},
	}
	
	return analysisResponse, nil
}

// ValidateProvider validates that a provider is available and configured
func (s *aiService) ValidateProvider(ctx context.Context, providerName string) error {
	// Check if provider is registered with auth service
	_, err := s.authService.GetProvider(ctx, providerName)
	if err != nil {
		return fmt.Errorf("provider '%s' not registered: %w", providerName, err)
	}
	
	// Check if credentials are available and valid
	_, err = s.authService.ValidateCredentials(ctx, providerName)
	if err != nil {
		return fmt.Errorf("provider '%s' not properly configured: %w", providerName, err)
	}
	
	return nil
}

// ListProviders returns all available AI providers
func (s *aiService) ListProviders(ctx context.Context) ([]AIProviderInfo, error) {
	providers, err := s.authService.ListProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	
	var providerInfos []AIProviderInfo
	for _, providerName := range providers {
		info := AIProviderInfo{
			Name:        providerName,
			DisplayName: s.getDisplayName(providerName),
			Description: s.getProviderDescription(providerName),
			Capabilities: s.getProviderCapabilities(providerName),
			Available:   true,
		}
		
		// Check if provider is configured
		if err := s.ValidateProvider(ctx, providerName); err == nil {
			info.Configured = true
		}
		
		providerInfos = append(providerInfos, info)
	}
	
	return providerInfos, nil
}

// RefreshClient refreshes the AI client for a provider
func (s *aiService) RefreshClient(ctx context.Context, providerName string) error {
	// Remove existing client
	delete(s.clients, providerName)
	
	// Get fresh credentials
	credentials, err := s.authService.RefreshCredentials(ctx, providerName)
	if err != nil {
		return fmt.Errorf("failed to refresh credentials: %w", err)
	}
	
	// Create new client
	_, err = s.CreateClient(ctx, providerName, credentials)
	if err != nil {
		return fmt.Errorf("failed to create refreshed client: %w", err)
	}
	
	return nil
}

// Helper methods

// buildManifestPrompt builds a prompt for manifest generation
func (s *aiService) buildManifestPrompt(request *ManifestRequest) string {
	prompt := fmt.Sprintf("Generate a project manifest for:\n"+
		"Project: %s\n"+
		"Framework: %s %s\n",
		request.ProjectName, request.Framework, request.Version)
	
	if request.Description != "" {
		prompt += fmt.Sprintf("Description: %s\n", request.Description)
	}
	
	if len(request.Features) > 0 {
		prompt += fmt.Sprintf("Features: %v\n", request.Features)
	}
	
	if request.Complexity != "" {
		prompt += fmt.Sprintf("Complexity: %s\n", request.Complexity)
	}
	
	prompt += "\nProvide a comprehensive project manifest with architecture recommendations, suggested components, and implementation guidance."
	
	return prompt
}

// buildAnalysisPrompt builds a prompt for project analysis
func (s *aiService) buildAnalysisPrompt(request *ProjectAnalysisRequest) string {
	prompt := fmt.Sprintf("Analyze the project at: %s\n", request.ProjectPath)
	
	if request.Framework != "" {
		prompt += fmt.Sprintf("Framework: %s\n", request.Framework)
	}
	
	if len(request.Focus) > 0 {
		prompt += fmt.Sprintf("Focus areas: %v\n", request.Focus)
	}
	
	prompt += "\nProvide analysis including potential issues, improvement suggestions, and architectural recommendations."
	
	return prompt
}

// getDisplayName returns a human-readable display name for a provider
func (s *aiService) getDisplayName(providerName string) string {
	displayNames := map[string]string{
		"claude":  "Claude (Anthropic)",
		"openai":  "OpenAI GPT",
		"gemini":  "Google Gemini",
		"ollama":  "Ollama (Local)",
	}
	
	if displayName, exists := displayNames[providerName]; exists {
		return displayName
	}
	
	return providerName
}

// getProviderDescription returns a description for a provider
func (s *aiService) getProviderDescription(providerName string) string {
	descriptions := map[string]string{
		"claude": "Anthropic's Claude AI assistant with strong reasoning capabilities",
		"openai": "OpenAI's GPT models for natural language processing",
		"gemini": "Google's multimodal AI model",
		"ollama": "Local AI models running on your machine",
	}
	
	if description, exists := descriptions[providerName]; exists {
		return description
	}
	
	return "AI provider"
}

// getProviderCapabilities returns capabilities for a provider
func (s *aiService) getProviderCapabilities(providerName string) []string {
	capabilities := map[string][]string{
		"claude": {"chat", "analysis", "code-review", "manifest-generation"},
		"openai": {"chat", "analysis", "code-generation", "manifest-generation"},
		"gemini": {"chat", "analysis", "multimodal"},
		"ollama": {"chat", "local-processing"},
	}
	
	if caps, exists := capabilities[providerName]; exists {
		return caps
	}
	
	return []string{"chat"}
}
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"atempo/internal/auth"
)

// AIClient handles communication with AI providers
type AIClient struct {
	authService *auth.AuthService
}

// NewAIClient creates a new AI client
func NewAIClient() (*AIClient, error) {
	authService, err := auth.NewAuthService()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	return &AIClient{
		authService: authService,
	}, nil
}

// ProviderInfo represents information about an available AI provider
type ProviderInfo struct {
	Name          string
	DisplayName   string
	Description   string
	Authenticated bool
}

// GetAvailableProviders returns a list of authenticated AI providers
func (c *AIClient) GetAvailableProviders() ([]ProviderInfo, error) {
	providers := []ProviderInfo{}

	// Check Claude
	if c.authService.IsAuthenticated("claude") {
		providers = append(providers, ProviderInfo{
			Name:          "claude",
			DisplayName:   "Claude (Anthropic)",
			Description:   "Advanced AI assistant by Anthropic",
			Authenticated: true,
		})
	}

	// Check OpenAI
	if c.authService.IsAuthenticated("openai") {
		providers = append(providers, ProviderInfo{
			Name:          "openai",
			DisplayName:   "GPT (OpenAI)",
			Description:   "GPT models by OpenAI",
			Authenticated: true,
		})
	}

	return providers, nil
}

// ChatRequest represents a request to an AI provider
type ChatRequest struct {
	Provider    string
	Messages    []Message
	Temperature float64
	MaxTokens   int
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents a response from an AI provider
type ChatResponse struct {
	Content string
	Usage   UsageInfo
}

// UsageInfo represents token usage information
type UsageInfo struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

// SendChatRequest sends a chat request to the specified AI provider
func (c *AIClient) SendChatRequest(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	switch req.Provider {
	case "claude":
		return c.sendClaudeRequest(ctx, req)
	case "openai":
		return c.sendOpenAIRequest(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// sendClaudeRequest sends a request to Claude API
func (c *AIClient) sendClaudeRequest(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Get credentials
	creds, err := c.authService.GetCredentials("claude")
	if err != nil {
		return nil, fmt.Errorf("failed to get Claude credentials: %w", err)
	}

	// Prepare request body
	claudeReq := map[string]interface{}{
		"model":      "claude-3-5-sonnet-20241022",
		"max_tokens": req.MaxTokens,
		"messages":   req.Messages,
	}

	if req.Temperature > 0 {
		claudeReq["temperature"] = req.Temperature
	}

	jsonData, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", creds.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var claudeResp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract content
	var content string
	if len(claudeResp.Content) > 0 {
		content = claudeResp.Content[0].Text
	}

	return &ChatResponse{
		Content: content,
		Usage: UsageInfo{
			InputTokens:  claudeResp.Usage.InputTokens,
			OutputTokens: claudeResp.Usage.OutputTokens,
			TotalTokens:  claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		},
	}, nil
}

// sendOpenAIRequest sends a request to OpenAI API
func (c *AIClient) sendOpenAIRequest(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Get credentials
	creds, err := c.authService.GetCredentials("openai")
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI credentials: %w", err)
	}

	// Prepare request body
	openaiReq := map[string]interface{}{
		"model":      "gpt-4",
		"messages":   req.Messages,
		"max_tokens": req.MaxTokens,
	}

	if req.Temperature > 0 {
		openaiReq["temperature"] = req.Temperature
	}

	jsonData, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+creds.APIKey)

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract content
	var content string
	if len(openaiResp.Choices) > 0 {
		content = openaiResp.Choices[0].Message.Content
	}

	return &ChatResponse{
		Content: content,
		Usage: UsageInfo{
			InputTokens:  openaiResp.Usage.PromptTokens,
			OutputTokens: openaiResp.Usage.CompletionTokens,
			TotalTokens:  openaiResp.Usage.TotalTokens,
		},
	}, nil
}

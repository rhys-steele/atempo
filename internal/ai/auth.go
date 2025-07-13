package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TestClaudeAuth tests authentication with Claude/Anthropic
func TestClaudeAuth(apiKey string) error {
	// Create a simple request to test the API key
	url := "https://api.anthropic.com/v1/messages"
	
	requestBody := map[string]interface{}{
		"model":      "claude-3-haiku-20240307",
		"max_tokens": 10,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "Hello",
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// TestOpenAIAuth tests authentication with OpenAI
func TestOpenAIAuth(apiKey string) error {
	// Create a simple request to test the API key
	url := "https://api.openai.com/v1/models"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetAPIClient returns an authenticated API client for the current provider
func GetAPIClient() (*APIClient, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if !config.Enabled {
		return nil, fmt.Errorf("AI features are disabled")
	}

	if config.CurrentProvider == "" {
		return nil, fmt.Errorf("no AI provider configured")
	}

	provider, exists := config.Providers[config.CurrentProvider]
	if !exists || !provider.Authenticated {
		return nil, fmt.Errorf("AI provider not authenticated")
	}

	apiKey, err := GetCredential(config.CurrentProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	return &APIClient{
		Provider: provider,
		APIKey:   apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// APIClient represents an authenticated API client
type APIClient struct {
	Provider ProviderConfig
	APIKey   string
	Client   *http.Client
}

// GenerateText generates text using the configured AI provider
func (c *APIClient) GenerateText(prompt string, maxTokens int) (string, error) {
	switch c.Provider.Name {
	case "claude":
		return c.generateWithClaude(prompt, maxTokens)
	case "openai":
		return c.generateWithOpenAI(prompt, maxTokens)
	default:
		return "", fmt.Errorf("unsupported provider: %s", c.Provider.Name)
	}
}

// generateWithClaude generates text using Claude
func (c *APIClient) generateWithClaude(prompt string, maxTokens int) (string, error) {
	url := "https://api.anthropic.com/v1/messages"
	
	requestBody := map[string]interface{}{
		"model":      c.Provider.DefaultModel,
		"max_tokens": maxTokens,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract text from Claude response
	if content, ok := response["content"].([]interface{}); ok && len(content) > 0 {
		if block, ok := content[0].(map[string]interface{}); ok {
			if text, ok := block["text"].(string); ok {
				return text, nil
			}
		}
	}

	return "", fmt.Errorf("unexpected response format")
}

// generateWithOpenAI generates text using OpenAI
func (c *APIClient) generateWithOpenAI(prompt string, maxTokens int) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	
	requestBody := map[string]interface{}{
		"model":      c.Provider.DefaultModel,
		"max_tokens": maxTokens,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract text from OpenAI response
	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format")
}
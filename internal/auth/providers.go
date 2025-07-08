package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OpenAIProvider handles OpenAI API authentication
type OpenAIProvider struct{}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) Description() string {
	return "OpenAI API authentication using API key"
}

func (p *OpenAIProvider) RequiredFields() []string {
	return []string{"api_key"}
}

func (p *OpenAIProvider) Authenticate(ctx context.Context, options AuthOptions) (*Credentials, error) {
	apiKey := options.APIKey

	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for OpenAI authentication")
	}

	// Validate API key format
	if !strings.HasPrefix(apiKey, "sk-") {
		return nil, fmt.Errorf("invalid OpenAI API key format (should start with 'sk-')")
	}

	// Create a context for validation
	validationCtx := ctx
	if validationCtx == nil {
		validationCtx = context.Background()
	}

	// Test the API key by making a simple request
	if err := p.validateAPIKey(validationCtx, apiKey); err != nil {
		return nil, fmt.Errorf("API key validation failed: %w", err)
	}

	return &Credentials{
		Provider: p.Name(),
		APIKey:   apiKey,
		Metadata: map[string]string{
			"validated_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

func (p *OpenAIProvider) Validate(ctx context.Context, creds *Credentials) error {
	if creds == nil || creds.APIKey == "" {
		return fmt.Errorf("no API key found")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return p.validateAPIKey(ctx, creds.APIKey)
}

func (p *OpenAIProvider) validateAPIKey(ctx context.Context, apiKey string) error {
	// Create a simple request to OpenAI API to validate the key
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create validation request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", "atempo-cli/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate API key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API key validation returned status %d", resp.StatusCode)
	}

	return nil
}

// ClaudeProvider handles Anthropic Claude API authentication
type ClaudeProvider struct{}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider() *ClaudeProvider {
	return &ClaudeProvider{}
}

func (p *ClaudeProvider) Name() string {
	return "claude"
}

func (p *ClaudeProvider) Description() string {
	return "Anthropic Claude API authentication using API key"
}

func (p *ClaudeProvider) RequiredFields() []string {
	return []string{"api_key"}
}

func (p *ClaudeProvider) Authenticate(ctx context.Context, options AuthOptions) (*Credentials, error) {
	apiKey := options.APIKey

	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for Claude authentication")
	}

	// Validate API key format
	if !strings.HasPrefix(apiKey, "sk-ant-") {
		return nil, fmt.Errorf("invalid Claude API key format (should start with 'sk-ant-')")
	}

	// Create a context for validation
	validationCtx := ctx
	if validationCtx == nil {
		validationCtx = context.Background()
	}

	// Test the API key by making a simple request
	if err := p.validateAPIKey(validationCtx, apiKey); err != nil {
		return nil, fmt.Errorf("API key validation failed: %w", err)
	}

	return &Credentials{
		Provider: p.Name(),
		APIKey:   apiKey,
		Metadata: map[string]string{
			"validated_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

func (p *ClaudeProvider) Validate(ctx context.Context, creds *Credentials) error {
	if creds == nil || creds.APIKey == "" {
		return fmt.Errorf("no API key found")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return p.validateAPIKey(ctx, creds.APIKey)
}

func (p *ClaudeProvider) validateAPIKey(ctx context.Context, apiKey string) error {
	// Create a simple request to Claude API to validate the key
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.anthropic.com/v1/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create validation request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("User-Agent", "atempo-cli/1.0")
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate API key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API key validation returned status %d", resp.StatusCode)
	}

	return nil
}

// AtempoProvider handles Atempo platform authentication
type AtempoProvider struct{}

// NewAtempoProvider creates a new Atempo provider
func NewAtempoProvider() *AtempoProvider {
	return &AtempoProvider{}
}

func (p *AtempoProvider) Name() string {
	return "atempo"
}

func (p *AtempoProvider) Description() string {
	return "Atempo platform authentication"
}

func (p *AtempoProvider) RequiredFields() []string {
	return []string{} // Will be determined by the auth flow
}

func (p *AtempoProvider) Authenticate(ctx context.Context, options AuthOptions) (*Credentials, error) {
	// TODO: Implement Atempo platform authentication
	// This would typically involve OAuth2 flow or device authorization
	return nil, fmt.Errorf("atempo platform authentication not yet implemented")
}

func (p *AtempoProvider) Validate(ctx context.Context, creds *Credentials) error {
	// TODO: Implement Atempo credential validation
	return fmt.Errorf("atempo credential validation not yet implemented")
}

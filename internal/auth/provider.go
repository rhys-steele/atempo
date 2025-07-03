package auth

import (
	"context"
	"fmt"
	"time"
)

// Provider represents an authentication provider
type Provider interface {
	// Name returns the provider name (e.g., "openai", "claude", "atempo")
	Name() string
	
	// Authenticate performs the authentication flow for this provider
	Authenticate(ctx context.Context, options AuthOptions) (*Credentials, error)
	
	// Validate checks if existing credentials are still valid
	Validate(ctx context.Context, creds *Credentials) error
	
	// RequiredFields returns the fields required for authentication
	RequiredFields() []string
	
	// Description returns a human-readable description of the provider
	Description() string
}

// AuthOptions contains authentication options
type AuthOptions struct {
	APIKey     string            // API key for providers that use keys
	Interactive bool             // Whether to prompt user interactively
	Force      bool             // Force re-authentication even if credentials exist
	Metadata   map[string]string // Provider-specific metadata
}

// Credentials represents stored authentication credentials
type Credentials struct {
	Provider    string            `json:"provider"`
	APIKey      string            `json:"api_key,omitempty"`
	AccessToken string            `json:"access_token,omitempty"`
	RefreshToken string           `json:"refresh_token,omitempty"`
	ExpiresAt   int64            `json:"expires_at,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// IsValid checks if credentials are still valid
func (c *Credentials) IsValid() bool {
	if c == nil {
		return false
	}
	
	// Check if we have any form of credential
	if c.APIKey == "" && c.AccessToken == "" {
		return false
	}
	
	// Check expiration if set
	if c.ExpiresAt > 0 {
		return c.ExpiresAt > getCurrentUnixTime()
	}
	
	return true
}

// ProviderRegistry manages available authentication providers
type ProviderRegistry struct {
	providers map[string]Provider
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	registry := &ProviderRegistry{
		providers: make(map[string]Provider),
	}
	
	// Register built-in providers
	registry.Register(NewOpenAIProvider())
	registry.Register(NewClaudeProvider())
	registry.Register(NewAtempoProvider())
	
	return registry
}

// Register adds a provider to the registry
func (r *ProviderRegistry) Register(provider Provider) {
	r.providers[provider.Name()] = provider
}

// GetProvider returns a provider by name
func (r *ProviderRegistry) GetProvider(name string) (Provider, error) {
	provider, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("unknown authentication provider: %s", name)
	}
	return provider, nil
}

// ListProviders returns all available providers
func (r *ProviderRegistry) ListProviders() []Provider {
	providers := make([]Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers
}

// HasProvider checks if a provider exists
func (r *ProviderRegistry) HasProvider(name string) bool {
	_, exists := r.providers[name]
	return exists
}

// getCurrentUnixTime returns current Unix timestamp
func getCurrentUnixTime() int64 {
	return time.Now().Unix()
}
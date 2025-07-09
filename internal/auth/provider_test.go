package auth

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCredentials_IsValid(t *testing.T) {
	tests := []struct {
		name        string
		credentials *Credentials
		expected    bool
	}{
		{
			name:        "Nil credentials",
			credentials: nil,
			expected:    false,
		},
		{
			name: "Empty credentials",
			credentials: &Credentials{
				Provider: "test",
			},
			expected: false,
		},
		{
			name: "Valid API key",
			credentials: &Credentials{
				Provider: "openai",
				APIKey:   "sk-test-key",
			},
			expected: true,
		},
		{
			name: "Valid access token",
			credentials: &Credentials{
				Provider:    "claude",
				AccessToken: "test-token",
			},
			expected: true,
		},
		{
			name: "Expired credentials",
			credentials: &Credentials{
				Provider:    "openai",
				APIKey:      "sk-test-key",
				ExpiresAt:   time.Now().Add(-time.Hour).Unix(),
			},
			expected: false,
		},
		{
			name: "Valid non-expired credentials",
			credentials: &Credentials{
				Provider:    "openai",
				APIKey:      "sk-test-key",
				ExpiresAt:   time.Now().Add(time.Hour).Unix(),
			},
			expected: true,
		},
		{
			name: "No expiration set",
			credentials: &Credentials{
				Provider: "openai",
				APIKey:   "sk-test-key",
				ExpiresAt: 0,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.credentials.IsValid()
			if result != tt.expected {
				t.Errorf("Expected IsValid() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNewProviderRegistry(t *testing.T) {
	registry := NewProviderRegistry()

	// Verify registry is initialized
	if registry == nil {
		t.Error("Expected non-nil registry")
	}

	if registry.providers == nil {
		t.Error("Expected providers map to be initialized")
	}

	// Verify built-in providers are registered
	expectedProviders := []string{"openai", "claude", "atempo"}
	for _, providerName := range expectedProviders {
		provider, err := registry.GetProvider(providerName)
		if err != nil {
			t.Errorf("Expected provider %s to be registered, got error: %v", providerName, err)
		}
		if provider == nil {
			t.Errorf("Expected provider %s to be non-nil", providerName)
		}
	}
}

func TestProviderRegistry_Register(t *testing.T) {
	registry := NewProviderRegistry()
	
	// Create a mock provider
	mockProvider := &MockProvider{
		name:        "test-provider",
		description: "Test provider",
		requiredFields: []string{"api_key"},
	}

	// Register the provider
	registry.Register(mockProvider)

	// Verify it was registered
	provider, err := registry.GetProvider("test-provider")
	if err != nil {
		t.Errorf("Expected provider to be registered, got error: %v", err)
	}

	if provider.Name() != "test-provider" {
		t.Errorf("Expected provider name 'test-provider', got '%s'", provider.Name())
	}
}

func TestProviderRegistry_GetProvider(t *testing.T) {
	registry := NewProviderRegistry()

	// Test getting existing provider
	provider, err := registry.GetProvider("openai")
	if err != nil {
		t.Errorf("Expected no error getting existing provider, got: %v", err)
	}
	if provider == nil {
		t.Error("Expected non-nil provider")
	}

	// Test getting non-existing provider
	_, err = registry.GetProvider("non-existing")
	if err == nil {
		t.Error("Expected error getting non-existing provider")
	}

	expectedError := "unknown authentication provider: non-existing"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestProviderRegistry_ListProviders(t *testing.T) {
	registry := NewProviderRegistry()

	providers := registry.ListProviders()

	// Verify we have providers
	if len(providers) == 0 {
		t.Error("Expected at least one provider")
	}

	// Verify expected providers are present
	expectedProviders := []string{"openai", "claude", "atempo"}
	providerNames := make(map[string]bool)
	for _, provider := range providers {
		providerNames[provider.Name()] = true
	}

	for _, expectedName := range expectedProviders {
		if !providerNames[expectedName] {
			t.Errorf("Expected provider '%s' to be in list", expectedName)
		}
	}
}

func TestAuthOptions_Structure(t *testing.T) {
	// Test AuthOptions structure
	options := AuthOptions{
		APIKey:      "test-key",
		Interactive: true,
		Force:       false,
		Metadata: map[string]string{
			"endpoint": "https://api.example.com",
			"model":    "gpt-4",
		},
	}

	// Verify all fields are accessible
	if options.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", options.APIKey)
	}
	if !options.Interactive {
		t.Error("Expected Interactive to be true")
	}
	if options.Force {
		t.Error("Expected Force to be false")
	}
	if len(options.Metadata) != 2 {
		t.Errorf("Expected 2 metadata entries, got %d", len(options.Metadata))
	}
	if options.Metadata["endpoint"] != "https://api.example.com" {
		t.Errorf("Expected endpoint 'https://api.example.com', got '%s'", options.Metadata["endpoint"])
	}
}

func TestCredentials_Structure(t *testing.T) {
	// Test Credentials structure
	creds := Credentials{
		Provider:     "openai",
		APIKey:       "sk-test-key",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
		Metadata: map[string]string{
			"organization": "org-123",
			"project":      "proj-456",
		},
	}

	// Verify all fields are accessible
	if creds.Provider != "openai" {
		t.Errorf("Expected provider 'openai', got '%s'", creds.Provider)
	}
	if creds.APIKey != "sk-test-key" {
		t.Errorf("Expected API key 'sk-test-key', got '%s'", creds.APIKey)
	}
	if creds.AccessToken != "access-token" {
		t.Errorf("Expected access token 'access-token', got '%s'", creds.AccessToken)
	}
	if creds.RefreshToken != "refresh-token" {
		t.Errorf("Expected refresh token 'refresh-token', got '%s'", creds.RefreshToken)
	}
	if creds.ExpiresAt == 0 {
		t.Error("Expected ExpiresAt to be set")
	}
	if len(creds.Metadata) != 2 {
		t.Errorf("Expected 2 metadata entries, got %d", len(creds.Metadata))
	}
}

func TestGetCurrentUnixTime(t *testing.T) {
	// Test getCurrentUnixTime function
	currentTime := getCurrentUnixTime()
	
	// Verify it's a reasonable Unix timestamp
	if currentTime < 1000000000 { // Should be after 2001
		t.Error("Expected reasonable Unix timestamp")
	}
	
	// Verify it's close to current time
	actualTime := time.Now().Unix()
	if abs(currentTime-actualTime) > 10 { // Should be within 10 seconds
		t.Error("Expected current time to be close to actual time")
	}
}

// MockProvider is a mock implementation of Provider for testing
type MockProvider struct {
	name           string
	description    string
	requiredFields []string
	authError      error
	validateError  error
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Description() string {
	return m.description
}

func (m *MockProvider) RequiredFields() []string {
	return m.requiredFields
}

func (m *MockProvider) Authenticate(ctx context.Context, options AuthOptions) (*Credentials, error) {
	if m.authError != nil {
		return nil, m.authError
	}
	
	return &Credentials{
		Provider: m.name,
		APIKey:   options.APIKey,
	}, nil
}

func (m *MockProvider) Validate(ctx context.Context, creds *Credentials) error {
	if m.validateError != nil {
		return m.validateError
	}
	
	if creds == nil || creds.APIKey == "" {
		return fmt.Errorf("invalid credentials")
	}
	
	return nil
}

func TestMockProvider(t *testing.T) {
	// Test our mock provider
	mock := &MockProvider{
		name:           "test",
		description:    "Test provider",
		requiredFields: []string{"api_key"},
	}

	// Test Name()
	if mock.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", mock.Name())
	}

	// Test Description()
	if mock.Description() != "Test provider" {
		t.Errorf("Expected description 'Test provider', got '%s'", mock.Description())
	}

	// Test RequiredFields()
	fields := mock.RequiredFields()
	if len(fields) != 1 || fields[0] != "api_key" {
		t.Errorf("Expected required fields ['api_key'], got %v", fields)
	}

	// Test Authenticate()
	ctx := context.Background()
	options := AuthOptions{APIKey: "test-key"}
	creds, err := mock.Authenticate(ctx, options)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if creds.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", creds.APIKey)
	}

	// Test Validate()
	err = mock.Validate(ctx, creds)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test Validate() with invalid credentials
	err = mock.Validate(ctx, &Credentials{})
	if err == nil {
		t.Error("Expected error with invalid credentials")
	}
}

func TestProvider_Interface(t *testing.T) {
	// Test that MockProvider implements Provider interface
	var provider Provider = &MockProvider{
		name:           "test",
		description:    "Test provider",
		requiredFields: []string{"api_key"},
	}

	// Test interface methods
	if provider.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", provider.Name())
	}

	if provider.Description() != "Test provider" {
		t.Errorf("Expected description 'Test provider', got '%s'", provider.Description())
	}

	fields := provider.RequiredFields()
	if len(fields) != 1 || fields[0] != "api_key" {
		t.Errorf("Expected required fields ['api_key'], got %v", fields)
	}

	// Test Authenticate method exists
	ctx := context.Background()
	options := AuthOptions{APIKey: "test-key"}
	_, err := provider.Authenticate(ctx, options)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test Validate method exists
	creds := &Credentials{APIKey: "test-key"}
	err = provider.Validate(ctx, creds)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Helper function to calculate absolute value
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
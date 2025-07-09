package services

import (
	"context"
	"fmt"
	"time"

	"atempo/internal/auth"
)

// AuthService provides business operations for authentication and provider management
type AuthService interface {
	// RegisterProvider registers a new authentication provider
	RegisterProvider(ctx context.Context, provider auth.AuthProvider) error
	
	// GetProvider retrieves an authentication provider by name
	GetProvider(ctx context.Context, name string) (auth.AuthProvider, error)
	
	// ListProviders returns all registered authentication providers
	ListProviders(ctx context.Context) ([]string, error)
	
	// Authenticate authenticates with a specific provider
	Authenticate(ctx context.Context, providerName string, credentials map[string]string) (*auth.Credentials, error)
	
	// ValidateCredentials validates stored credentials
	ValidateCredentials(ctx context.Context, providerName string) (*auth.Credentials, error)
	
	// RefreshCredentials refreshes expired credentials if possible
	RefreshCredentials(ctx context.Context, providerName string) (*auth.Credentials, error)
	
	// StoreCredentials securely stores credentials for a provider
	StoreCredentials(ctx context.Context, providerName string, credentials *auth.Credentials) error
	
	// RemoveCredentials removes stored credentials for a provider
	RemoveCredentials(ctx context.Context, providerName string) error
	
	// IsAuthenticated checks if the user is authenticated with a provider
	IsAuthenticated(ctx context.Context, providerName string) (bool, error)
}

// CredentialsService provides business operations for secure credential storage
type CredentialsService interface {
	// Store securely stores credentials
	Store(ctx context.Context, key string, credentials *auth.Credentials) error
	
	// Load loads stored credentials
	Load(ctx context.Context, key string) (*auth.Credentials, error)
	
	// Delete removes stored credentials
	Delete(ctx context.Context, key string) error
	
	// List returns all stored credential keys
	List(ctx context.Context) ([]string, error)
	
	// Exists checks if credentials exist for a key
	Exists(ctx context.Context, key string) (bool, error)
}

// authService implements AuthService
type authService struct {
	providerRegistry *auth.ProviderRegistry
	credentialsStorage auth.CredentialsStorage
}

// NewAuthService creates a new AuthService implementation
func NewAuthService(credentialsStorage auth.CredentialsStorage) AuthService {
	return &authService{
		providerRegistry:   auth.NewProviderRegistry(),
		credentialsStorage: credentialsStorage,
	}
}

// RegisterProvider registers a new authentication provider
func (s *authService) RegisterProvider(ctx context.Context, provider auth.AuthProvider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}
	
	s.providerRegistry.Register(provider.Name(), provider)
	return nil
}

// GetProvider retrieves an authentication provider by name
func (s *authService) GetProvider(ctx context.Context, name string) (auth.AuthProvider, error) {
	provider, exists := s.providerRegistry.Get(name)
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}
	
	return provider, nil
}

// ListProviders returns all registered authentication providers
func (s *authService) ListProviders(ctx context.Context) ([]string, error) {
	return s.providerRegistry.List(), nil
}

// Authenticate authenticates with a specific provider
func (s *authService) Authenticate(ctx context.Context, providerName string, credentials map[string]string) (*auth.Credentials, error) {
	provider, err := s.GetProvider(ctx, providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}
	
	// Authenticate with the provider
	authCredentials, err := provider.Authenticate(ctx, credentials)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	
	// Store the credentials securely
	if err := s.StoreCredentials(ctx, providerName, authCredentials); err != nil {
		return nil, fmt.Errorf("failed to store credentials: %w", err)
	}
	
	return authCredentials, nil
}

// ValidateCredentials validates stored credentials
func (s *authService) ValidateCredentials(ctx context.Context, providerName string) (*auth.Credentials, error) {
	// Load stored credentials
	credentials, err := s.credentialsStorage.Load(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}
	
	// Check if credentials are valid
	if !credentials.IsValid() {
		return nil, fmt.Errorf("credentials are invalid or expired")
	}
	
	// Get provider for additional validation
	provider, err := s.GetProvider(ctx, providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}
	
	// Validate with provider
	if err := provider.ValidateCredentials(ctx, credentials); err != nil {
		return nil, fmt.Errorf("provider validation failed: %w", err)
	}
	
	return credentials, nil
}

// RefreshCredentials refreshes expired credentials if possible
func (s *authService) RefreshCredentials(ctx context.Context, providerName string) (*auth.Credentials, error) {
	// Load current credentials
	credentials, err := s.credentialsStorage.Load(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}
	
	// Get provider
	provider, err := s.GetProvider(ctx, providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}
	
	// Check if provider supports refresh
	refreshProvider, ok := provider.(auth.RefreshableProvider)
	if !ok {
		return nil, fmt.Errorf("provider '%s' does not support credential refresh", providerName)
	}
	
	// Refresh credentials
	newCredentials, err := refreshProvider.RefreshCredentials(ctx, credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh credentials: %w", err)
	}
	
	// Store refreshed credentials
	if err := s.StoreCredentials(ctx, providerName, newCredentials); err != nil {
		return nil, fmt.Errorf("failed to store refreshed credentials: %w", err)
	}
	
	return newCredentials, nil
}

// StoreCredentials securely stores credentials for a provider
func (s *authService) StoreCredentials(ctx context.Context, providerName string, credentials *auth.Credentials) error {
	if credentials == nil {
		return fmt.Errorf("credentials cannot be nil")
	}
	
	if !credentials.IsValid() {
		return fmt.Errorf("cannot store invalid credentials")
	}
	
	return s.credentialsStorage.Store(providerName, credentials)
}

// RemoveCredentials removes stored credentials for a provider
func (s *authService) RemoveCredentials(ctx context.Context, providerName string) error {
	return s.credentialsStorage.Delete(providerName)
}

// IsAuthenticated checks if the user is authenticated with a provider
func (s *authService) IsAuthenticated(ctx context.Context, providerName string) (bool, error) {
	// Try to validate credentials
	_, err := s.ValidateCredentials(ctx, providerName)
	if err != nil {
		// Try to refresh if validation fails
		if _, refreshErr := s.RefreshCredentials(ctx, providerName); refreshErr != nil {
			return false, nil // Not authenticated
		}
		return true, nil // Refreshed successfully
	}
	
	return true, nil // Already valid
}

// credentialsService implements CredentialsService
type credentialsService struct {
	storage auth.CredentialsStorage
}

// NewCredentialsService creates a new CredentialsService implementation
func NewCredentialsService(storage auth.CredentialsStorage) CredentialsService {
	return &credentialsService{
		storage: storage,
	}
}

// Store securely stores credentials
func (s *credentialsService) Store(ctx context.Context, key string, credentials *auth.Credentials) error {
	if credentials == nil {
		return fmt.Errorf("credentials cannot be nil")
	}
	
	if !credentials.IsValid() {
		return fmt.Errorf("cannot store invalid credentials")
	}
	
	return s.storage.Store(key, credentials)
}

// Load loads stored credentials
func (s *credentialsService) Load(ctx context.Context, key string) (*auth.Credentials, error) {
	credentials, err := s.storage.Load(key)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}
	
	// Check if credentials have expired
	if credentials.ExpiresAt != nil && time.Now().After(*credentials.ExpiresAt) {
		return nil, fmt.Errorf("credentials have expired")
	}
	
	return credentials, nil
}

// Delete removes stored credentials
func (s *credentialsService) Delete(ctx context.Context, key string) error {
	return s.storage.Delete(key)
}

// List returns all stored credential keys
func (s *credentialsService) List(ctx context.Context) ([]string, error) {
	// Note: This would need to be implemented in the storage interface
	// For now, return an error indicating it's not implemented
	return nil, fmt.Errorf("list operation not supported by current storage implementation")
}

// Exists checks if credentials exist for a key
func (s *credentialsService) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.storage.Load(key)
	if err != nil {
		return false, nil // Assume not exists if load fails
	}
	return true, nil
}

// Helper interfaces for extended provider capabilities

// RefreshableProvider defines a provider that can refresh credentials
type RefreshableProvider interface {
	auth.AuthProvider
	RefreshCredentials(ctx context.Context, credentials *auth.Credentials) (*auth.Credentials, error)
}

// ProviderWithValidation defines a provider that can validate credentials
type ProviderWithValidation interface {
	auth.AuthProvider
	ValidateCredentials(ctx context.Context, credentials *auth.Credentials) error
}
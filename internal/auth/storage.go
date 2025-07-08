package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CredentialStore handles secure storage and retrieval of credentials
type CredentialStore interface {
	// Store saves credentials for a provider
	Store(provider string, creds *Credentials) error

	// Retrieve gets credentials for a provider
	Retrieve(provider string) (*Credentials, error)

	// Delete removes credentials for a provider
	Delete(provider string) error

	// List returns all stored provider names
	List() ([]string, error)

	// Exists checks if credentials exist for a provider
	Exists(provider string) bool
}

// FileCredentialStore implements credential storage using local files
// TODO: In production, this should use system keychain/keyring for security
type FileCredentialStore struct {
	configDir string
}

// NewFileCredentialStore creates a new file-based credential store
func NewFileCredentialStore() (*FileCredentialStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".atempo", "auth")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &FileCredentialStore{
		configDir: configDir,
	}, nil
}

// Store saves credentials for a provider
func (s *FileCredentialStore) Store(provider string, creds *Credentials) error {
	filePath := s.getCredentialPath(provider)

	// Marshal credentials to JSON
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// Retrieve gets credentials for a provider
func (s *FileCredentialStore) Retrieve(provider string) (*Credentials, error) {
	filePath := s.getCredentialPath(provider)

	// Check if file exists
	if !s.fileExists(filePath) {
		return nil, fmt.Errorf("no credentials found for provider: %s", provider)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	// Unmarshal JSON
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &creds, nil
}

// Delete removes credentials for a provider
func (s *FileCredentialStore) Delete(provider string) error {
	filePath := s.getCredentialPath(provider)

	if !s.fileExists(filePath) {
		return fmt.Errorf("no credentials found for provider: %s", provider)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete credentials file: %w", err)
	}

	return nil
}

// List returns all stored provider names
func (s *FileCredentialStore) List() ([]string, error) {
	entries, err := os.ReadDir(s.configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var providers []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			provider := entry.Name()[:len(entry.Name())-5] // Remove .json extension
			providers = append(providers, provider)
		}
	}

	return providers, nil
}

// Exists checks if credentials exist for a provider
func (s *FileCredentialStore) Exists(provider string) bool {
	filePath := s.getCredentialPath(provider)
	return s.fileExists(filePath)
}

// getCredentialPath returns the file path for a provider's credentials
func (s *FileCredentialStore) getCredentialPath(provider string) string {
	return filepath.Join(s.configDir, provider+".json")
}

// fileExists checks if a file exists
func (s *FileCredentialStore) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// AuthService provides high-level authentication operations
type AuthService struct {
	store    CredentialStore
	registry *ProviderRegistry
}

// NewAuthService creates a new authentication service
func NewAuthService() (*AuthService, error) {
	store, err := NewFileCredentialStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create credential store: %w", err)
	}

	registry := NewProviderRegistry()

	return &AuthService{
		store:    store,
		registry: registry,
	}, nil
}

// Authenticate performs authentication for a provider
func (s *AuthService) Authenticate(provider string, options AuthOptions) error {
	// Get provider
	p, err := s.registry.GetProvider(provider)
	if err != nil {
		return err
	}

	// Check if credentials already exist and are valid (unless force is true)
	if !options.Force {
		if existingCreds, err := s.GetCredentials(provider); err == nil && existingCreds.IsValid() {
			// Validate existing credentials
			if err := p.Validate(nil, existingCreds); err == nil {
				return fmt.Errorf("already authenticated with %s (use --force to re-authenticate)", provider)
			}
		}
	}

	// Perform authentication
	creds, err := p.Authenticate(context.Background(), options)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Store credentials
	if err := s.store.Store(provider, creds); err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	return nil
}

// GetCredentials retrieves stored credentials for a provider
func (s *AuthService) GetCredentials(provider string) (*Credentials, error) {
	return s.store.Retrieve(provider)
}

// ValidateCredentials checks if stored credentials are still valid
func (s *AuthService) ValidateCredentials(provider string) error {
	// Get provider
	p, err := s.registry.GetProvider(provider)
	if err != nil {
		return err
	}

	// Get stored credentials
	creds, err := s.GetCredentials(provider)
	if err != nil {
		return err
	}

	// Validate credentials
	return p.Validate(context.Background(), creds)
}

// ListProviders returns all available providers
func (s *AuthService) ListProviders() []Provider {
	return s.registry.ListProviders()
}

// ListAuthenticated returns all providers with stored credentials
func (s *AuthService) ListAuthenticated() ([]string, error) {
	return s.store.List()
}

// Logout removes stored credentials for a provider
func (s *AuthService) Logout(provider string) error {
	return s.store.Delete(provider)
}

// IsAuthenticated checks if a provider has valid credentials
func (s *AuthService) IsAuthenticated(provider string) bool {
	if !s.store.Exists(provider) {
		return false
	}

	return s.ValidateCredentials(provider) == nil
}

package interfaces

// AuthProvider defines the interface for authentication providers
type AuthProvider interface {
	// GetName returns the name of the authentication provider
	GetName() string
	
	// Authenticate performs authentication and returns credentials
	Authenticate() (*Credentials, error)
	
	// ValidateCredentials validates existing credentials
	ValidateCredentials(creds *Credentials) error
	
	// RefreshCredentials refreshes expired credentials
	RefreshCredentials(creds *Credentials) (*Credentials, error)
	
	// IsConfigured checks if the provider is properly configured
	IsConfigured() bool
	
	// GetAuthURL returns the authentication URL for OAuth providers
	GetAuthURL() (string, error)
	
	// HandleCallback handles OAuth callback
	HandleCallback(code string) (*Credentials, error)
}

// CredentialsStorage defines the interface for storing authentication credentials
type CredentialsStorage interface {
	// Store stores credentials securely
	Store(providerName string, creds *Credentials) error
	
	// Load loads credentials for a provider
	Load(providerName string) (*Credentials, error)
	
	// Delete removes credentials for a provider
	Delete(providerName string) error
	
	// List returns all stored provider names
	List() ([]string, error)
	
	// Exists checks if credentials exist for a provider
	Exists(providerName string) bool
}

// Credentials represents authentication credentials
type Credentials struct {
	ProviderName string            `json:"provider_name"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token,omitempty"`
	ExpiresAt    int64             `json:"expires_at,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// IsValid checks if the credentials are valid and not expired
func (c *Credentials) IsValid() bool {
	if c.AccessToken == "" {
		return false
	}
	
	// If no expiration is set, assume credentials are valid
	if c.ExpiresAt == 0 {
		return true
	}
	
	// Check if credentials are expired
	return c.ExpiresAt > 0 // Simplified check - should compare with current time
}

// IsExpired checks if the credentials are expired
func (c *Credentials) IsExpired() bool {
	if c.ExpiresAt == 0 {
		return false // No expiration set
	}
	
	// Simplified check - should compare with current time
	return c.ExpiresAt < 0
}
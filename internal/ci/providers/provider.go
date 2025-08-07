package providers

import (
	"embed"
	"fmt"

	"atempo/internal/logger"
)

// Provider interface for extensibility (moved from ci package to avoid cycles)
type Provider interface {
	Name() string
	ValidateConfig(config *CIConfig) error
	GenerateConfig(config *CIConfig, templateFS embed.FS) ([]byte, error)
	GetConfigFileName() string                              // .github/workflows/ci.yml, .gitlab-ci.yml
	GetConfigPath(projectPath string) string               // Full path where config should be written
	GetDefaultSettings(framework string) map[string]interface{}
	PromptForSettings(framework string) (map[string]interface{}, error)
	SupportedFrameworks() []string
}

// ProviderRegistry manages CI providers for extensible provider management
type ProviderRegistry struct {
	providers map[CIProvider]Provider
	logger    *logger.Logger
}

// NewProviderRegistry creates a new provider registry with built-in providers
func NewProviderRegistry(logger *logger.Logger) *ProviderRegistry {
	registry := &ProviderRegistry{
		providers: make(map[CIProvider]Provider),
		logger:    logger,
	}

	// Register built-in providers
	registry.Register(NewGitHubProvider(logger))
	registry.Register(NewGitLabProvider(logger))

	return registry
}

// Register adds a provider to the registry
func (pr *ProviderRegistry) Register(provider Provider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	providerName := CIProvider(provider.Name())
	if _, exists := pr.providers[providerName]; exists {
		return fmt.Errorf("provider '%s' is already registered", provider.Name())
	}

	pr.providers[providerName] = provider
	
	// Note: Logger integration will be added in Phase 3

	return nil
}

// Get retrieves a provider by name
func (pr *ProviderRegistry) Get(name CIProvider) (Provider, error) {
	provider, exists := pr.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}
	return provider, nil
}

// List returns all available provider names
func (pr *ProviderRegistry) List() []CIProvider {
	providers := make([]CIProvider, 0, len(pr.providers))
	for name := range pr.providers {
		providers = append(providers, name)
	}
	return providers
}

// GetSupportedFrameworks returns supported frameworks for a provider
func (pr *ProviderRegistry) GetSupportedFrameworks(provider CIProvider) ([]string, error) {
	p, err := pr.Get(provider)
	if err != nil {
		return nil, err
	}
	return p.SupportedFrameworks(), nil
}

// IsFrameworkSupported checks if a framework is supported by a provider
func (pr *ProviderRegistry) IsFrameworkSupported(provider CIProvider, framework string) bool {
	frameworks, err := pr.GetSupportedFrameworks(provider)
	if err != nil {
		return false
	}

	for _, supportedFramework := range frameworks {
		if framework == supportedFramework {
			return true
		}
	}
	return false
}

// GetCompatibleProviders returns providers that support a given framework
func (pr *ProviderRegistry) GetCompatibleProviders(framework string) []CIProvider {
	var compatibleProviders []CIProvider

	for providerName := range pr.providers {
		if pr.IsFrameworkSupported(providerName, framework) {
			compatibleProviders = append(compatibleProviders, providerName)
		}
	}

	return compatibleProviders
}

// ValidateProviderFrameworkCombination validates that a provider supports a framework
func (pr *ProviderRegistry) ValidateProviderFrameworkCombination(provider CIProvider, framework string) error {
	if !pr.IsFrameworkSupported(provider, framework) {
		supportedFrameworks, _ := pr.GetSupportedFrameworks(provider)
		return fmt.Errorf("provider '%s' does not support framework '%s'. Supported frameworks: %v", 
			provider, framework, supportedFrameworks)
	}
	return nil
}
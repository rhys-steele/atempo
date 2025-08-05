package providers

import (
	"fmt"

	"atempo/internal/ci"
	"atempo/internal/logger"
)

// ProviderRegistry manages CI providers for extensible provider management
type ProviderRegistry struct {
	providers map[ci.CIProvider]ci.Provider
	logger    *logger.Logger
}

// NewProviderRegistry creates a new provider registry with built-in providers
func NewProviderRegistry(logger *logger.Logger) *ProviderRegistry {
	registry := &ProviderRegistry{
		providers: make(map[ci.CIProvider]ci.Provider),
		logger:    logger,
	}

	// Register built-in providers
	registry.Register(NewGitHubProvider(logger))
	registry.Register(NewGitLabProvider(logger))

	return registry
}

// Register adds a provider to the registry
func (pr *ProviderRegistry) Register(provider ci.Provider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	providerName := ci.CIProvider(provider.Name())
	if _, exists := pr.providers[providerName]; exists {
		return fmt.Errorf("provider '%s' is already registered", provider.Name())
	}

	pr.providers[providerName] = provider
	
	// Note: Logger integration will be added in Phase 3

	return nil
}

// Get retrieves a provider by name
func (pr *ProviderRegistry) Get(name ci.CIProvider) (ci.Provider, error) {
	provider, exists := pr.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}
	return provider, nil
}

// List returns all available provider names
func (pr *ProviderRegistry) List() []ci.CIProvider {
	providers := make([]ci.CIProvider, 0, len(pr.providers))
	for name := range pr.providers {
		providers = append(providers, name)
	}
	return providers
}

// GetSupportedFrameworks returns supported frameworks for a provider
func (pr *ProviderRegistry) GetSupportedFrameworks(provider ci.CIProvider) ([]string, error) {
	p, err := pr.Get(provider)
	if err != nil {
		return nil, err
	}
	return p.SupportedFrameworks(), nil
}

// IsFrameworkSupported checks if a framework is supported by a provider
func (pr *ProviderRegistry) IsFrameworkSupported(provider ci.CIProvider, framework string) bool {
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
func (pr *ProviderRegistry) GetCompatibleProviders(framework string) []ci.CIProvider {
	var compatibleProviders []ci.CIProvider

	for providerName := range pr.providers {
		if pr.IsFrameworkSupported(providerName, framework) {
			compatibleProviders = append(compatibleProviders, providerName)
		}
	}

	return compatibleProviders
}

// ValidateProviderFrameworkCombination validates that a provider supports a framework
func (pr *ProviderRegistry) ValidateProviderFrameworkCombination(provider ci.CIProvider, framework string) error {
	if !pr.IsFrameworkSupported(provider, framework) {
		supportedFrameworks, _ := pr.GetSupportedFrameworks(provider)
		return fmt.Errorf("provider '%s' does not support framework '%s'. Supported frameworks: %v", 
			provider, framework, supportedFrameworks)
	}
	return nil
}
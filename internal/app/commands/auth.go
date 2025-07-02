package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"atempo/internal/auth"
	"golang.org/x/term"
)

// AuthCommand handles authentication with various providers
type AuthCommand struct {
	*BaseCommand
	authService *auth.AuthService
}

// NewAuthCommand creates a new auth command
func NewAuthCommand(ctx *CommandContext) *AuthCommand {
	authService, err := auth.NewAuthService()
	if err != nil {
		// For now, create a command that will show the error
		// In production, this should be handled more gracefully
		authService = nil
	}

	return &AuthCommand{
		BaseCommand: NewBaseCommand(
			"auth",
			"Authenticate with AI providers and Atempo platform",
			"atempo auth [provider] [options]",
			ctx,
		),
		authService: authService,
	}
}

// Execute runs the auth command
func (c *AuthCommand) Execute(ctx context.Context, args []string) error {
	if c.authService == nil {
		return fmt.Errorf("authentication service unavailable")
	}

	// Parse subcommands and options
	if len(args) == 0 {
		return c.showAuthStatus()
	}

	subcommand := args[0]
	
	switch subcommand {
	case "login":
		return c.handleLogin(args[1:])
	case "logout":
		return c.handleLogout(args[1:])
	case "status":
		return c.showAuthStatus()
	case "list":
		return c.listProviders()
	case "validate":
		return c.validateCredentials(args[1:])
	default:
		// If it's not a subcommand, treat it as a provider name for login
		return c.handleLogin(args)
	}
}

// handleLogin performs authentication for a provider
func (c *AuthCommand) handleLogin(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: atempo auth login <provider>\nAvailable providers: %s", c.getProviderNames())
	}

	provider := args[0]
	
	// Parse flags
	options := auth.AuthOptions{
		Interactive: true,
	}
	
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--force", "-f":
			options.Force = true
		case "--api-key":
			if i+1 >= len(args) {
				return fmt.Errorf("--api-key requires a value")
			}
			options.APIKey = args[i+1]
			i++ // Skip the next argument
		case "--non-interactive":
			options.Interactive = false
		}
	}

	// Get provider info
	providers := c.authService.ListProviders()
	var selectedProvider auth.Provider
	for _, p := range providers {
		if p.Name() == provider {
			selectedProvider = p
			break
		}
	}
	
	if selectedProvider == nil {
		return fmt.Errorf("unknown provider: %s\nAvailable providers: %s", provider, c.getProviderNames())
	}

	fmt.Printf("🔐 Authenticating with %s\n", selectedProvider.Description())

	// Get API key if not provided and provider requires it
	requiredFields := selectedProvider.RequiredFields()
	if contains(requiredFields, "api_key") && options.APIKey == "" {
		if options.Interactive {
			var err error
			options.APIKey, err = c.promptForAPIKey(provider)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("API key required for %s. Use --api-key flag or run interactively", provider)
		}
	}

	// Perform authentication
	fmt.Printf("→ Validating credentials...\n")
	if err := c.authService.Authenticate(provider, options); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Printf("✅ Successfully authenticated with %s\n", provider)
	return nil
}

// handleLogout removes credentials for a provider
func (c *AuthCommand) handleLogout(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: atempo auth logout <provider>")
	}

	provider := args[0]
	
	if err := c.authService.Logout(provider); err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	fmt.Printf("✅ Logged out from %s\n", provider)
	return nil
}

// showAuthStatus displays authentication status for all providers
func (c *AuthCommand) showAuthStatus() error {
	fmt.Println("🔐 Authentication Status")
	fmt.Println(strings.Repeat("=", 40))

	providers := c.authService.ListProviders()
	authenticated, err := c.authService.ListAuthenticated()
	if err != nil {
		return fmt.Errorf("failed to list authenticated providers: %w", err)
	}

	authenticatedMap := make(map[string]bool)
	for _, provider := range authenticated {
		authenticatedMap[provider] = true
	}

	for _, provider := range providers {
		status := "❌ Not authenticated"
		if authenticatedMap[provider.Name()] {
			if c.authService.IsAuthenticated(provider.Name()) {
				status = "✅ Authenticated"
			} else {
				status = "⚠️  Invalid credentials"
			}
		}

		fmt.Printf("  %-12s %s\n", provider.Name(), status)
		fmt.Printf("               %s\n", provider.Description())
		fmt.Println()
	}

	fmt.Println("Commands:")
	fmt.Println("  atempo auth login <provider>     # Authenticate with a provider")
	fmt.Println("  atempo auth logout <provider>    # Remove credentials")
	fmt.Println("  atempo auth validate <provider>  # Test existing credentials")
	fmt.Println("  atempo auth list                 # List available providers")

	return nil
}

// listProviders shows all available authentication providers
func (c *AuthCommand) listProviders() error {
	fmt.Println("Available Authentication Providers:")
	fmt.Println()

	providers := c.authService.ListProviders()
	for _, provider := range providers {
		fmt.Printf("  %-12s %s\n", provider.Name(), provider.Description())
		
		requiredFields := provider.RequiredFields()
		if len(requiredFields) > 0 {
			fmt.Printf("               Required: %s\n", strings.Join(requiredFields, ", "))
		}
		fmt.Println()
	}

	return nil
}

// validateCredentials tests stored credentials for a provider
func (c *AuthCommand) validateCredentials(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: atempo auth validate <provider>")
	}

	provider := args[0]
	
	fmt.Printf("→ Validating %s credentials...\n", provider)
	
	if err := c.authService.ValidateCredentials(provider); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		return nil // Don't return error to avoid double error display
	}

	fmt.Printf("✅ %s credentials are valid\n", provider)
	return nil
}

// promptForAPIKey securely prompts for an API key
func (c *AuthCommand) promptForAPIKey(provider string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Printf("\nEnter your %s API key: ", strings.ToUpper(provider))
	
	// Try to read from terminal securely (hidden input)
	if term.IsTerminal(int(syscall.Stdin)) {
		keyBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", fmt.Errorf("failed to read API key: %w", err)
		}
		fmt.Println() // Add newline after hidden input
		return strings.TrimSpace(string(keyBytes)), nil
	}
	
	// Fallback to regular input if not a terminal
	fmt.Print("(input will be visible) ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	
	return strings.TrimSpace(input), nil
}

// getProviderNames returns a comma-separated list of provider names
func (c *AuthCommand) getProviderNames() string {
	providers := c.authService.ListProviders()
	names := make([]string, len(providers))
	for i, provider := range providers {
		names[i] = provider.Name()
	}
	return strings.Join(names, ", ")
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
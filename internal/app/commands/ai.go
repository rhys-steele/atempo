package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"atempo/internal/ai"
	"atempo/internal/auth"
	"atempo/internal/logger"
	"golang.org/x/term"
)

// AICommand handles AI-related operations
type AICommand struct {
	logger *logger.Logger
}

// NewAICommand creates a new AI command instance
func NewAICommand() *AICommand {
	return &AICommand{}
}

// Name returns the command name
func (c *AICommand) Name() string {
	return "ai"
}

// Description returns the command description
func (c *AICommand) Description() string {
	return "Manage AI features and authentication"
}

// Usage returns the command usage information
func (c *AICommand) Usage() string {
	return `ai <subcommand> [options]

Available subcommands:
  enable              Enable AI features globally
  disable             Disable AI features globally
  auth <provider>     Authenticate with AI provider (claude, openai)
  unauth [provider]   Remove authentication for provider
  status              Show AI configuration and authentication status
  configure           Configure AI preferences
  providers           List available AI providers
  models              List available models for current provider
  test                Test AI connection and authentication
  generate context    Generate AI context for current project
  update context      Update existing AI context
  init               Initialize AI features for existing project
  analyze            Analyze current project structure
  suggest            Get AI suggestions for improvements
  validate           Validate project against best practices
  reset              Reset AI configuration to defaults

Examples:
  atempo ai enable
  atempo ai auth claude
  atempo ai status
  atempo ai generate context`
}

// Execute runs the AI command
func (c *AICommand) Execute(ctx context.Context, args []string) error {
	// Create logger for this command
	log, err := logger.NewQuiet("ai-command")
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer log.Close()
	c.logger = log

	if len(args) == 0 {
		return fmt.Errorf("no subcommand specified\n\n%s", c.Usage())
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "enable":
		return c.handleEnable(subArgs)
	case "disable":
		return c.handleDisable(subArgs)
	case "auth":
		return c.handleAuth(subArgs)
	case "unauth":
		return c.handleUnauth(subArgs)
	case "status":
		return c.handleStatus(subArgs)
	case "configure":
		return c.handleConfigure(subArgs)
	case "providers":
		return c.handleProviders(subArgs)
	case "models":
		return c.handleModels(subArgs)
	case "test":
		return c.handleTest(subArgs)
	case "generate":
		return c.handleGenerate(subArgs)
	case "update":
		return c.handleUpdate(subArgs)
	case "init":
		return c.handleInit(subArgs)
	case "analyze":
		return c.handleAnalyze(subArgs)
	case "suggest":
		return c.handleSuggest(subArgs)
	case "validate":
		return c.handleValidate(subArgs)
	case "reset":
		return c.handleReset(subArgs)
	default:
		return fmt.Errorf("unknown subcommand: %s\n\n%s", subcommand, c.Usage())
	}
}

// handleEnable enables AI features globally
func (c *AICommand) handleEnable(args []string) error {
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	config.Enabled = true
	if err := ai.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save AI config: %w", err)
	}

	fmt.Println("✓ AI features enabled globally")
	return nil
}

// handleDisable disables AI features globally
func (c *AICommand) handleDisable(args []string) error {
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	config.Enabled = false
	if err := ai.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save AI config: %w", err)
	}

	fmt.Println("✓ AI features disabled globally")
	return nil
}

// handleAuth handles AI provider authentication
func (c *AICommand) handleAuth(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("provider required. Available providers: claude, openai")
	}

	provider := strings.ToLower(args[0])
	switch provider {
	case "claude":
		return c.authClaude()
	case "openai":
		return c.authOpenAI()
	default:
		return fmt.Errorf("unknown provider: %s. Available providers: claude, openai", provider)
	}
}

// authClaude handles Claude authentication
func (c *AICommand) authClaude() error {
	fmt.Print("Enter your Claude API key: ")
	
	// Use secure input for API key
	var apiKey string
	if term.IsTerminal(int(syscall.Stdin)) {
		keyBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
		fmt.Println() // Add newline after hidden input
		apiKey = strings.TrimSpace(string(keyBytes))
	} else {
		if _, err := fmt.Scanln(&apiKey); err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
	}

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Use the modern auth system
	authService, err := auth.NewAuthService()
	if err != nil {
		return fmt.Errorf("failed to create auth service: %w", err)
	}

	// Authenticate using the modern system
	authOptions := auth.AuthOptions{
		APIKey:      apiKey,
		Force:       true,
		Interactive: false,
	}

	if err := authService.Authenticate("claude", authOptions); err != nil {
		return fmt.Errorf("failed to authenticate with Claude: %w", err)
	}

	// Update configuration
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	config.CurrentProvider = "claude"
	config.Providers["claude"] = ai.ProviderConfig{
		Name:         "claude",
		Authenticated: true,
		DefaultModel: "claude-3-haiku-20240307",
	}

	if err := ai.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save AI config: %w", err)
	}

	fmt.Println("✓ Successfully authenticated with Claude")
	return nil
}

// authOpenAI handles OpenAI authentication
func (c *AICommand) authOpenAI() error {
	fmt.Print("Enter your OpenAI API key: ")
	
	// Use secure input for API key
	var apiKey string
	if term.IsTerminal(int(syscall.Stdin)) {
		keyBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
		fmt.Println() // Add newline after hidden input
		apiKey = strings.TrimSpace(string(keyBytes))
	} else {
		if _, err := fmt.Scanln(&apiKey); err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
	}

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Use the modern auth system
	authService, err := auth.NewAuthService()
	if err != nil {
		return fmt.Errorf("failed to create auth service: %w", err)
	}

	// Authenticate using the modern system
	authOptions := auth.AuthOptions{
		APIKey:      apiKey,
		Force:       true,
		Interactive: false,
	}

	if err := authService.Authenticate("openai", authOptions); err != nil {
		return fmt.Errorf("failed to authenticate with OpenAI: %w", err)
	}

	// Update configuration
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	config.CurrentProvider = "openai"
	config.Providers["openai"] = ai.ProviderConfig{
		Name:         "openai",
		Authenticated: true,
		DefaultModel: "gpt-4",
	}

	if err := ai.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save AI config: %w", err)
	}

	fmt.Println("✓ Successfully authenticated with OpenAI")
	return nil
}

// handleUnauth removes authentication for a provider
func (c *AICommand) handleUnauth(args []string) error {
	provider := "all"
	if len(args) > 0 {
		provider = strings.ToLower(args[0])
	}

	// Use the modern auth system
	authService, err := auth.NewAuthService()
	if err != nil {
		return fmt.Errorf("failed to create auth service: %w", err)
	}

	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	if provider == "all" {
		// Remove all credentials
		for providerName := range config.Providers {
			if err := authService.Logout(providerName); err != nil {
				fmt.Printf("Warning: failed to remove %s credentials: %v\n", providerName, err)
			}
			delete(config.Providers, providerName)
		}
		config.CurrentProvider = ""
		fmt.Println("✓ Removed all AI provider authentication")
	} else {
		// Remove specific provider
		if err := authService.Logout(provider); err != nil {
			return fmt.Errorf("failed to remove credentials: %w", err)
		}
		delete(config.Providers, provider)
		if config.CurrentProvider == provider {
			config.CurrentProvider = ""
		}
		fmt.Printf("✓ Removed %s authentication\n", provider)
	}

	if err := ai.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save AI config: %w", err)
	}

	return nil
}

// handleStatus shows AI configuration and authentication status
func (c *AICommand) handleStatus(args []string) error {
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	fmt.Printf("AI Features: %s\n", getStatusIndicator(config.Enabled))
	fmt.Printf("Current Provider: %s\n", getProviderName(config.CurrentProvider))
	fmt.Printf("Configuration File: %s\n", ai.GetConfigPath())

	if len(config.Providers) == 0 {
		fmt.Println("\nNo AI providers configured")
		fmt.Println("Use 'atempo ai auth <provider>' to authenticate")
		return nil
	}

	fmt.Println("\nAuthenticated Providers:")
	for name, provider := range config.Providers {
		status := "✗ Not authenticated"
		if provider.Authenticated {
			status = "✓ Authenticated"
		}
		current := ""
		if name == config.CurrentProvider {
			current = " (current)"
		}
		fmt.Printf("  %s: %s%s\n", provider.Name, status, current)
		if provider.DefaultModel != "" {
			fmt.Printf("    Default Model: %s\n", provider.DefaultModel)
		}
	}

	return nil
}

// handleConfigure configures AI preferences
func (c *AICommand) handleConfigure(args []string) error {
	fmt.Println("AI Configuration")
	fmt.Println("This feature is coming soon!")
	return nil
}

// handleProviders lists available AI providers
func (c *AICommand) handleProviders(args []string) error {
	fmt.Println("Available AI Providers:")
	fmt.Println("  claude  - Anthropic Claude (Claude 3 Haiku, Sonnet, Opus)")
	fmt.Println("  openai  - OpenAI GPT (GPT-4, GPT-3.5)")
	fmt.Println("\nUse 'atempo ai auth <provider>' to authenticate")
	return nil
}

// handleModels lists available models for current provider
func (c *AICommand) handleModels(args []string) error {
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	if config.CurrentProvider == "" {
		return fmt.Errorf("no AI provider configured. Use 'atempo ai auth <provider>' first")
	}

	provider := config.Providers[config.CurrentProvider]
	fmt.Printf("Available models for %s:\n", provider.Name)

	switch provider.Name {
	case "claude":
		fmt.Println("  claude-3-haiku-20240307   - Fast, efficient model")
		fmt.Println("  claude-3-sonnet-20240229  - Balanced performance")
		fmt.Println("  claude-3-opus-20240229    - Most capable model")
	case "openai":
		fmt.Println("  gpt-4                     - Most capable model")
		fmt.Println("  gpt-4-turbo               - Fast, cost-effective")
		fmt.Println("  gpt-3.5-turbo             - Fast, affordable")
	default:
		fmt.Println("  Model information not available")
	}

	return nil
}

// handleTest tests AI connection and authentication
func (c *AICommand) handleTest(args []string) error {
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	if !config.Enabled {
		fmt.Println("✗ AI features are disabled")
		fmt.Println("Use 'atempo ai enable' to enable AI features")
		return nil
	}

	if config.CurrentProvider == "" {
		fmt.Println("✗ No AI provider configured")
		fmt.Println("Use 'atempo ai auth <provider>' to authenticate")
		return nil
	}

	provider := config.Providers[config.CurrentProvider]
	fmt.Printf("Testing %s connection...\n", provider.Name)

	// Use the modern auth system
	authService, err := auth.NewAuthService()
	if err != nil {
		return fmt.Errorf("failed to create auth service: %w", err)
	}

	// Test the credentials using the modern auth system
	if err := authService.ValidateCredentials(provider.Name); err != nil {
		return fmt.Errorf("✗ %s authentication failed: %w", provider.Name, err)
	}

	fmt.Printf("✓ %s connection successful\n", provider.Name)
	return nil
}

// handleGenerate generates AI context for current project
func (c *AICommand) handleGenerate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("specify what to generate. Available: context")
	}

	switch args[0] {
	case "context":
		return c.generateContext()
	default:
		return fmt.Errorf("unknown generate target: %s", args[0])
	}
}

// generateContext generates AI context for current project
func (c *AICommand) generateContext() error {
	// Get current working directory
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if this is an Atempo project
	projectInfo, err := ai.DetectProject(projectDir)
	if err != nil {
		return fmt.Errorf("failed to detect project: %w", err)
	}

	// Check if AI is enabled and authenticated
	if !ai.IsAIEnabled() {
		return fmt.Errorf("AI features are not enabled or authenticated. Run 'atempo ai enable' and 'atempo ai auth <provider>'")
	}

	fmt.Printf("Generating AI context for %s project...\n", projectInfo.Framework)

	// Generate dynamic AI context
	if err := ai.GenerateAIContext(projectDir, projectInfo.Name, projectInfo.Framework, projectInfo.Language, projectInfo.Version, true); err != nil {
		return fmt.Errorf("failed to generate AI context: %w", err)
	}

	fmt.Println("✓ AI context generated successfully")
	fmt.Println("Use 'claude /setup' to load the context in Claude")
	return nil
}

// handleUpdate updates existing AI context
func (c *AICommand) handleUpdate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("specify what to update. Available: context")
	}

	switch args[0] {
	case "context":
		return c.generateContext() // Same as generate for now
	default:
		return fmt.Errorf("unknown update target: %s", args[0])
	}
}

// handleInit initializes AI features for existing project
func (c *AICommand) handleInit(args []string) error {
	fmt.Println("Initializing AI features for existing project...")
	fmt.Println("This feature is coming soon!")
	return nil
}

// handleAnalyze analyzes current project structure
func (c *AICommand) handleAnalyze(args []string) error {
	fmt.Println("Analyzing project structure...")
	fmt.Println("This feature is coming soon!")
	return nil
}

// handleSuggest gets AI suggestions for improvements
func (c *AICommand) handleSuggest(args []string) error {
	fmt.Println("Getting AI suggestions for project improvements...")
	fmt.Println("This feature is coming soon!")
	return nil
}

// handleValidate validates project against best practices
func (c *AICommand) handleValidate(args []string) error {
	fmt.Println("Validating project against best practices...")
	fmt.Println("This feature is coming soon!")
	return nil
}

// handleReset resets AI configuration to defaults
func (c *AICommand) handleReset(args []string) error {
	fmt.Print("This will reset all AI configuration and remove all credentials. Continue? (y/N): ")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("Operation cancelled")
		return nil
	}

	// Use the modern auth system
	authService, err := auth.NewAuthService()
	if err != nil {
		return fmt.Errorf("failed to create auth service: %w", err)
	}

	// Remove all credentials
	config, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load AI config: %w", err)
	}

	for providerName := range config.Providers {
		if err := authService.Logout(providerName); err != nil {
			fmt.Printf("Warning: failed to remove %s credentials: %v\n", providerName, err)
		}
	}

	// Reset configuration
	if err := ai.ResetConfig(); err != nil {
		return fmt.Errorf("failed to reset AI config: %w", err)
	}

	fmt.Println("✓ AI configuration reset to defaults")
	return nil
}

// Helper functions

func getStatusIndicator(enabled bool) string {
	if enabled {
		return "✓ Enabled"
	}
	return "✗ Disabled"
}

func getProviderName(provider string) string {
	if provider == "" {
		return "None"
	}
	return provider
}
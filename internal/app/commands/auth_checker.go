package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

// AuthChecker handles authentication status for AI features
type AuthChecker struct{}

// NewAuthChecker creates a new authentication checker
func NewAuthChecker() *AuthChecker {
	return &AuthChecker{}
}

// IsAuthenticated checks if user is authenticated for AI features
func (a *AuthChecker) IsAuthenticated() bool {
	// Check for auth token file (simplified for demo)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	
	tokenPath := filepath.Join(homeDir, ".atempo", "auth.token")
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		return false
	}
	
	return true
}

// GetAuthStatus returns a user-friendly auth status message
func (a *AuthChecker) GetAuthStatus() (bool, string) {
	isAuth := a.IsAuthenticated()
	if isAuth {
		return true, "Authenticated - AI features enabled"
	}
	return false, "Not authenticated - using basic project setup"
}

// PromptAuthentication suggests authentication to the user
func (a *AuthChecker) PromptAuthentication() {
	fmt.Printf("\n%s🔐 Authentication Required for AI Features%s\n", ColorYellow, ColorReset)
	fmt.Printf("   AI-powered project manifests require authentication.\n")
	fmt.Printf("   Run %s'atempo auth'%s to connect your AI provider.\n\n", ColorCyan, ColorReset)
	fmt.Printf("   %sProceeding with basic project setup...%s\n", ColorGray, ColorReset)
}
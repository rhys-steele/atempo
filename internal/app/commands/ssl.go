package commands

import (
	"context"
	"fmt"

	"atempo/internal/docker"
)

// SSLCommand handles SSL certificate management for local HTTPS
type SSLCommand struct {
	*BaseCommand
}

// NewSSLCommand creates a new SSL command
func NewSSLCommand(ctx *CommandContext) *SSLCommand {
	return &SSLCommand{
		BaseCommand: NewBaseCommand(
			"ssl",
			"Manage SSL certificates for local HTTPS development",
			"atempo ssl <subcommand>",
			ctx,
		),
	}
}

// Execute runs the SSL command
func (c *SSLCommand) Execute(ctx context.Context, args []string) error {
	sslManager := docker.NewSSLManager()

	if len(args) == 0 {
		return sslManager.Status()
	}

	switch args[0] {
	case "status":
		return sslManager.Status()
	case "setup":
		return sslManager.Setup()
	case "renew":
		return sslManager.Renew()
	case "help", "--help", "-h":
		c.showHelp()
		return nil
	default:
		return fmt.Errorf("unknown SSL command: %s\n\nRun 'atempo ssl help' for available commands", args[0])
	}
}

// showHelp displays help information for the SSL command
func (c *SSLCommand) showHelp() {
	fmt.Println("SSL Certificate Management")
	fmt.Println("─────────────────────────")
	fmt.Println()
	fmt.Println("Manage SSL certificates for local HTTPS development with wildcard")
	fmt.Println("self-signed certificates for *.test domains.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  atempo ssl <subcommand>")
	fmt.Println()
	fmt.Println("AVAILABLE COMMANDS:")
	fmt.Println("  status    Show SSL certificate status and information")
	fmt.Println("  setup     Generate wildcard SSL certificate for *.test domains")
	fmt.Println("  renew     Renew the existing wildcard certificate")
	fmt.Println("  help      Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  atempo ssl status       # Check current certificate status")
	fmt.Println("  atempo ssl setup        # Generate wildcard certificate")
	fmt.Println("  atempo ssl renew        # Renew expiring certificate")
	fmt.Println()
	fmt.Println("ABOUT:")
	fmt.Println("  The SSL command generates self-signed wildcard certificates")
	fmt.Println("  for *.test domains, enabling HTTPS for all local projects.")
	fmt.Println("  Certificates are valid for 1 year and stored in ~/.atempo/ssl/")
}
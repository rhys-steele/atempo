package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"atempo/internal/docker"
)

// DNSCommand handles DNS configuration and diagnostics
type DNSCommand struct {
	*BaseCommand
}

// NewDNSCommand creates a new DNS command
func NewDNSCommand(ctx *CommandContext) *DNSCommand {
	return &DNSCommand{
		BaseCommand: NewBaseCommand(
			"dns",
			"Manage DNS configuration for custom domains",
			"atempo dns <subcommand>",
			ctx,
		),
	}
}

// Execute runs the DNS command
func (c *DNSCommand) Execute(ctx context.Context, args []string) error {
	dnsService := docker.NewDNSService()

	if len(args) == 0 {
		return dnsService.Status()
	}

	switch args[0] {
	case "status":
		return dnsService.Status()
	case "setup":
		return dnsService.Setup()
	case "start":
		if err := dnsService.Start(); err != nil {
			return err
		}
		fmt.Println("✓ DNS service started")
		return nil
	case "stop":
		if err := dnsService.Stop(); err != nil {
			return err
		}
		fmt.Println("✓ DNS service stopped")
		return nil
	case "test":
		domain := "example.local"
		if len(args) > 1 {
			domain = args[1]
		}
		return c.testDNS(domain)
	default:
		return fmt.Errorf("unknown DNS command: %s\n\nAvailable commands:\n  status  - Check DNS configuration\n  setup   - Interactive DNS setup\n  start   - Start DNS service\n  stop    - Stop DNS service\n  test    - Test DNS resolution", args[0])
	}
}

// testDNS tests DNS resolution for a specific domain
func (c *DNSCommand) testDNS(domain string) error {
	fmt.Printf("Testing DNS resolution for: %s\n", domain)

	cmd := exec.Command("nslookup", domain, "127.0.0.1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("✗ DNS test failed: %v\n", err)
		return err
	}

	fmt.Printf("✓ DNS resolution successful for %s\n", domain)
	return nil
}

// quickDNSTest performs a quick DNS resolution test
func (c *DNSCommand) quickDNSTest() bool {
	// Test multiple methods to be more reliable

	// Method 1: Try nslookup with port 5353
	cmd := exec.Command("nslookup", "-port=5353", "test.local", "127.0.0.1")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if cmd.Run() == nil {
		return true
	}

	// Method 2: Try dig with port 5353
	cmd = exec.Command("dig", "@127.0.0.1", "-p", "5353", "test.local", "+short")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if cmd.Run() == nil {
		return true
	}

	// Method 3: Test if we can at least connect to the DNS port
	cmd = exec.Command("nc", "-z", "127.0.0.1", "5353")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if cmd.Run() == nil {
		// Port is open - DNS service is running but may not have proper config
		return true
	}

	return false
}

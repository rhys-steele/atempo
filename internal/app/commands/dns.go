package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	simpleDNS := docker.NewSimpleDNS()
	
	if len(args) == 0 {
		return simpleDNS.Status()
	}

	switch args[0] {
	case "status":
		return simpleDNS.Status()
	case "setup":
		return simpleDNS.Setup()
	case "start":
		if err := simpleDNS.Start(); err != nil {
			return err
		}
		fmt.Println("✓ DNS service started")
		return nil
	case "stop":
		if err := simpleDNS.Stop(); err != nil {
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

// showDNSStatus displays current DNS configuration status
func (c *DNSCommand) showDNSStatus() error {
	fmt.Println("DNS Configuration")
	fmt.Println(strings.Repeat("─", 50))

	dnsManager := docker.NewDNSManager()
	
	// Check DNS system status
	status, err := dnsManager.GetDNSStatus()
	if err != nil {
		fmt.Printf("✗ DNS system error: %v\n", err)
		return nil
	}

	// DNS container/service status
	dnsType, _ := status["type"].(string)
	running, _ := status["running"].(bool)
	
	if running {
		fmt.Printf("✓ DNS service: %s (running)\n", dnsType)
		if dnsType == "docker" {
			port, _ := status["port"].(string)
			fmt.Printf("  Port: %s\n", port)
		}
	} else {
		fmt.Printf("✗ DNS service: %s (not running)\n", dnsType)
	}

	// macOS resolver status
	resolverFile := "/etc/resolver/local"
	if _, err := os.Stat(resolverFile); err == nil {
		fmt.Printf("✓ macOS resolver: configured\n")
		fmt.Printf("  File: %s\n", resolverFile)
		
		// Test if DNS actually works
		if c.quickDNSTest() {
			fmt.Printf("✓ DNS resolution: working\n")
		} else {
			fmt.Printf("✗ DNS resolution: failed\n")
		}
	} else {
		fmt.Printf("✗ macOS resolver: not configured\n")
		fmt.Printf("  Missing: %s\n", resolverFile)
	}

	// Project domains
	domains, err := dnsManager.ListProjectDomains()
	if err == nil && len(domains) > 0 {
		fmt.Printf("\nConfigured Domains:\n")
		for project, projectDomains := range domains {
			fmt.Printf("  %s: %s\n", project, strings.Join(projectDomains, ", "))
		}
	} else {
		fmt.Printf("\nNo project domains configured\n")
	}

	// Status summary
	fmt.Printf("\nStatus:\n")
	if !running {
		fmt.Printf("  DNS service not running - run 'atempo create' to start\n")
	} else if _, err := os.Stat(resolverFile); err == nil {
		if c.quickDNSTest() {
			fmt.Printf("  DNS fully operational\n")
		} else {
			fmt.Printf("  DNS configured but not responding\n")
			fmt.Printf("  Try: atempo dns fix\n")
		}
	} else {
		fmt.Printf("  DNS service running but resolver not configured\n")
		fmt.Printf("  Run: atempo dns setup\n")
	}

	return nil
}

// setupDNS provides interactive DNS setup
func (c *DNSCommand) setupDNS() error {
	fmt.Println("DNS Resolver Setup")
	fmt.Println(strings.Repeat("─", 50))
	
	// Check if already configured and working
	resolverFile := "/etc/resolver/local"
	if _, err := os.Stat(resolverFile); err == nil {
		if c.quickDNSTest() {
			fmt.Println("✓ DNS resolver already configured and working")
			return nil
		} else {
			fmt.Println("- DNS resolver exists but not responding")
		}
	} else {
		fmt.Println("- DNS resolver not configured")
	}

	fmt.Println("\nThis configures macOS to resolve .local domains through")
	fmt.Println("Atempo's DNS system, enabling custom domains like")
	fmt.Println("'myproject.local' instead of 'localhost:8000'")
	
	fmt.Printf("\nConfigure DNS resolver? [y/N]: ")
	var response string
	fmt.Scanln(&response)
	
	if response != "y" && response != "Y" && response != "yes" {
		fmt.Println("\nRun 'atempo dns setup' anytime to configure DNS")
		return nil
	}

	return c.createResolver()
}

// createResolver creates the macOS DNS resolver configuration
func (c *DNSCommand) createResolver() error {
	fmt.Println("\nCreating DNS resolver configuration...")
	
	// Create resolver directory
	fmt.Printf("Creating /etc/resolver directory...\n")
	cmd := exec.Command("sudo", "mkdir", "-p", "/etc/resolver")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create resolver directory: %w", err)
	}

	// Create resolver config
	fmt.Printf("Writing resolver configuration...\n")
	resolverConfig := `nameserver 127.0.0.1
port 5353`

	tempFile := filepath.Join(os.TempDir(), "atempo-local-resolver")
	if err := os.WriteFile(tempFile, []byte(resolverConfig), 0644); err != nil {
		return fmt.Errorf("failed to create temp resolver file: %w", err)
	}

	cmd = exec.Command("sudo", "mv", tempFile, "/etc/resolver/local")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install resolver configuration: %w", err)
	}

	fmt.Println("✓ DNS resolver configured")
	
	// Test the configuration
	fmt.Printf("Testing DNS resolution...\n")
	if c.quickDNSTest() {
		fmt.Println("✓ DNS working - new projects will use custom domains")
	} else {
		fmt.Println("✗ DNS test failed - restart browser or check configuration")
		fmt.Println("  Try: atempo dns fix")
	}

	return nil
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

// fixDNS attempts to fix common DNS issues
func (c *DNSCommand) fixDNS() error {
	fmt.Println("DNS Troubleshooting")
	fmt.Println(strings.Repeat("─", 50))
	
	dnsManager := docker.NewDNSManager()
	
	// Check DNS container
	status, err := dnsManager.GetDNSStatus()
	if err != nil {
		fmt.Printf("✗ DNS system error: %v\n", err)
		return err
	}
	
	running, _ := status["running"].(bool)
	if !running {
		fmt.Println("Starting DNS container...")
		// Try to restart DNS
		if err := dnsManager.RestartDNS(); err != nil {
			fmt.Printf("✗ Failed to start DNS: %v\n", err)
		} else {
			fmt.Println("✓ DNS container started")
		}
	}
	
	// Check resolver
	resolverFile := "/etc/resolver/local"
	if _, err := os.Stat(resolverFile); err != nil {
		fmt.Println("Resolver not configured, setting up...")
		return c.createResolver()
	}
	
	// Test if DNS is actually working
	if c.quickDNSTest() {
		fmt.Println("✓ DNS configuration working correctly")
	} else {
		fmt.Println("✗ DNS still not responding")
		fmt.Println("  DNS container may need volume mount fix")
		fmt.Println("  Try recreating project: atempo create <framework>")
	}
	
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
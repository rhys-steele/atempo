package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"atempo/internal/docker"
	"atempo/internal/registry"
)

// openProjectInBrowser opens the project in a web browser
func (r *CommandRegistry) openProjectInBrowser(projectName string, args []string) error {
	// Load registry to get project details
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}

	// Find the project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project not found: %s", projectName)
	}

	// Update project status to get current URLs and services
	err = reg.UpdateProjectStatus(projectName)
	if err != nil {
		return fmt.Errorf("failed to update project status: %w", err)
	}

	// Reload to get updated project info
	project, err = reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("failed to reload project: %w", err)
	}

	// Check if project has running services
	if project.Status == "stopped" || project.Status == "no-docker" || project.Status == "no-services" {
		return fmt.Errorf("project '%s' is not running. Start it with: %s up", projectName, projectName)
	}

	var targetURL string

	if len(args) == 0 {
		// Try to use DNS URL first, then fall back to port-based URLs
		dnsURL := r.getDNSURL(projectName)
		if dnsURL != "" {
			targetURL = dnsURL
			ShowInfo(fmt.Sprintf("Opening main application: %s", targetURL))
		} else {
			// Fallback to port-based URLs
			if len(project.URLs) == 0 {
				return fmt.Errorf("no web URLs found for project '%s'. Make sure services are running and have exposed web ports", projectName)
			}
			targetURL = r.findBestPortURL(project)
			if targetURL == "" {
				targetURL = project.URLs[0] // Last resort
			}
			ShowInfo(fmt.Sprintf("Opening main application: %s (DNS not configured)", targetURL))
		}
	} else {
		// Open specific service
		serviceName := args[0]
		found := false

		for _, service := range project.Services {
			if service.Name == serviceName {
				if service.URL == "" {
					return fmt.Errorf("service '%s' doesn't have a web URL (no exposed web ports)", serviceName)
				}
				targetURL = service.URL
				found = true
				break
			}
		}

		if !found {
			// List available services for user
			var availableServices []string
			for _, service := range project.Services {
				if service.URL != "" {
					availableServices = append(availableServices, service.Name)
				}
			}

			if len(availableServices) == 0 {
				return fmt.Errorf("no services with web URLs found for project '%s'", projectName)
			}

			return fmt.Errorf("service '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", "))
		}

		ShowInfo(fmt.Sprintf("Opening service '%s': %s", serviceName, targetURL))
	}

	// Open URL in default browser
	return openURL(targetURL)
}

// getDNSURL returns the DNS-based URL for a project if available and working
func (r *CommandRegistry) getDNSURL(projectName string) string {
	// Check if DNS is configured for this project using simple DNS
	if r.isDNSConfiguredForProject(projectName) {
		domain := fmt.Sprintf("%s.local", projectName)

		// Quick test if DNS resolution actually works
		if r.testDNSResolution(domain) {
			return fmt.Sprintf("http://%s", domain)
		}
	}

	return ""
}

// testDNSResolution quickly tests if a domain resolves (with timeout)
func (r *CommandRegistry) testDNSResolution(domain string) bool {
	// Try a quick nslookup with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, "nslookup", domain, "127.0.0.1")
	cmd.Stdout = nil // Suppress output
	cmd.Stderr = nil // Suppress errors

	return cmd.Run() == nil
}

// isDNSConfiguredForProject checks if DNS is configured for a project using simple DNS
func (r *CommandRegistry) isDNSConfiguredForProject(projectName string) bool {
	dnsService := docker.NewDNSService()

	// Check if DNS service is running
	if !dnsService.IsRunning() {
		return false
	}

	// Check if project has DNS configuration file in simple DNS
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	dnsFile := filepath.Join(homeDir, ".atempo", "dns", "projects", fmt.Sprintf("%s.dns", projectName))
	_, err = os.Stat(dnsFile)
	return err == nil
}

// findBestPortURL finds the best port-based URL from project URLs
func (r *CommandRegistry) findBestPortURL(project *registry.Project) string {
	// Look for main web service URLs first
	for _, url := range project.URLs {
		// Parse the URL to get the port
		if strings.Contains(url, ":8000") || strings.Contains(url, ":8001") || strings.Contains(url, ":8080") {
			return url
		}
	}

	// Return first available URL
	if len(project.URLs) > 0 {
		return project.URLs[0]
	}

	return ""
}

// openURL opens a URL in the default browser (cross-platform)
func openURL(url string) error {
	var cmd *exec.Cmd

	// Try platform-specific commands
	if err := exec.Command("which", "open").Run(); err == nil {
		// macOS
		cmd = exec.Command("open", url)
	} else if err := exec.Command("which", "xdg-open").Run(); err == nil {
		// Linux
		cmd = exec.Command("xdg-open", url)
	} else if err := exec.Command("where", "cmd").Run(); err == nil {
		// Windows
		cmd = exec.Command("cmd", "/c", "start", url)
	} else {
		return fmt.Errorf("unable to find browser command (tried: open, xdg-open, cmd)")
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open URL in browser: %w", err)
	}

	// Don't wait for browser to close
	go func() {
		cmd.Wait()
	}()

	ShowSuccess("Browser opened", url)
	return nil
}

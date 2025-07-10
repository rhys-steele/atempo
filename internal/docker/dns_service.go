package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DNSService manages a single container with dnsmasq + nginx for local development
type DNSService struct {
	configDir string
}

const (
	DNSContainerName = "atempo-dns"
	DNSImage         = "nginx:alpine"
	DNSPort          = "5353"
	HTTPPort         = "80"
	DNSNetworkName   = "atempo-net"
	DNSStaticIP      = "172.21.0.53"
	DNSDomain        = "test"  // Use .test instead of .local to avoid mDNS conflicts
)

// NewDNSService creates a new DNS service manager
func NewDNSService() *DNSService {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".atempo", "dns")

	return &DNSService{
		configDir: configDir,
	}
}

// Setup performs one-time DNS setup
func (s *DNSService) Setup() error {
	fmt.Println("DNS Setup")
	fmt.Println(strings.Repeat("─", 50))

	// Always validate and fix resolver configuration
	resolverFile := fmt.Sprintf("/etc/resolver/%s", DNSDomain)
	expectedConfig := fmt.Sprintf("nameserver %s\n", DNSStaticIP)
	needsSetup := true

	// Check if resolver is correctly configured
	if contents, err := os.ReadFile(resolverFile); err == nil {
		if string(contents) == expectedConfig {
			fmt.Println("✓ DNS resolver already configured")
			needsSetup = false
		} else {
			fmt.Printf("⚠ DNS resolver misconfigured (contains: %q, expected: %q)\n", 
				strings.TrimSpace(string(contents)), strings.TrimSpace(expectedConfig))
		}
	} else {
		fmt.Println("DNS resolver not found")
	}

	// Check DNS service status
	if s.IsRunning() {
		fmt.Println("✓ DNS service running")
	} else {
		fmt.Println("DNS service not running")
		needsSetup = true
	}

	// If everything is correctly configured, exit early
	if !needsSetup {
		return nil
	}

	// Auto-fix resolver configuration without prompting
	fmt.Println("\nFixing DNS resolver configuration...")
	return s.createResolver()
}

// createResolver creates the macOS DNS resolver configuration
func (s *DNSService) createResolver() error {
	resolverFile := fmt.Sprintf("/etc/resolver/%s", DNSDomain)
	expectedConfig := fmt.Sprintf("nameserver %s\n", DNSStaticIP)

	// Check if resolver is already correctly configured
	if contents, err := os.ReadFile(resolverFile); err == nil {
		if string(contents) == expectedConfig {
			fmt.Println("✓ DNS resolver already configured")
			return nil
		}
		fmt.Println("Updating DNS resolver configuration...")
	} else {
		fmt.Println("Creating DNS resolver configuration...")
	}

	// Create resolver directory
	cmd := exec.Command("sudo", "mkdir", "-p", "/etc/resolver")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create resolver directory: %w", err)
	}

	// Create resolver config - use static IP instead of port mapping
	tempFile := filepath.Join(os.TempDir(), "atempo-resolver")
	if err := os.WriteFile(tempFile, []byte(expectedConfig), 0644); err != nil {
		return fmt.Errorf("failed to create resolver config: %w", err)
	}

	cmd = exec.Command("sudo", "mv", tempFile, resolverFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install resolver: %w", err)
	}

	fmt.Println("✓ DNS resolver configured")
	
	// Flush DNS cache
	fmt.Println("Flushing DNS cache...")
	exec.Command("sudo", "dscacheutil", "-flushcache").Run()
	exec.Command("sudo", "killall", "-HUP", "mDNSResponder").Run()
	fmt.Println("✓ DNS cache flushed")

	// Start DNS service
	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to start DNS service: %w", err)
	}

	fmt.Println("✓ DNS service started")
	fmt.Println("✓ Setup complete - new projects will use custom domains")

	return nil
}

// Start starts the DNS container
func (s *DNSService) Start() error {
	if s.IsRunning() {
		return nil // Already running
	}

	// Create Docker network if it doesn't exist
	if err := s.createNetwork(); err != nil {
		return fmt.Errorf("failed to create Docker network: %w", err)
	}

	// Check for port conflicts and clean up
	if err := s.handlePortConflicts(); err != nil {
		return fmt.Errorf("failed to handle port conflicts: %w", err)
	}

	// Create config directories
	if err := s.createConfigDirectories(); err != nil {
		return err
	}

	// Remove any existing container
	s.remove()

	// Create startup script that runs both services
	startupScript := `#!/bin/sh
set -e
echo "Installing dnsmasq..."
apk add --no-cache dnsmasq

echo "Setting up nginx config..."
cat > /etc/nginx/nginx.conf << 'EOF'
events { worker_connections 1024; }
http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    # Include all project configs
    include /etc/atempo/projects/*.nginx;
    
    # Default server (fallback)
    server {
        listen 80 default_server;
        return 404;
    }
}
EOF

echo "Starting services..."
nginx &
exec dnsmasq --conf-file=/etc/atempo/dnsmasq.conf --no-daemon
`

	scriptFile := filepath.Join(s.configDir, "startup.sh")
	if err := os.WriteFile(scriptFile, []byte(startupScript), 0755); err != nil {
		return fmt.Errorf("failed to create startup script: %w", err)
	}

	// Create and start container with custom network and static IP
	cmd := exec.Command("docker", "run", "-d",
		"--name", DNSContainerName,
		"--network", DNSNetworkName,
		"--ip", DNSStaticIP,
		"-v", fmt.Sprintf("%s:/etc/atempo", s.configDir),
		"--restart", "unless-stopped",
		DNSImage,
		"/etc/atempo/startup.sh")

	return cmd.Run()
}

// Stop stops the DNS container
func (s *DNSService) Stop() error {
	if !s.IsRunning() {
		return nil
	}

	cmd := exec.Command("docker", "stop", DNSContainerName)
	if err := cmd.Run(); err != nil {
		return err
	}

	s.remove()
	return nil
}

// IsRunning checks if the DNS container is running
func (s *DNSService) IsRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", DNSContainerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), DNSContainerName)
}

// flushDNSCache flushes the local DNS cache (macOS)
func (s *DNSService) flushDNSCache() {
	// Flush DNS cache without requiring sudo password
	// Use exec.Command with /bin/sh to run multiple commands
	cmd := exec.Command("/bin/sh", "-c", "dscacheutil -flushcache 2>/dev/null || true; killall -HUP mDNSResponder 2>/dev/null || true")
	cmd.Run() // Ignore errors - cache flushing is best effort
}

// createNetwork creates the Docker network if it doesn't exist
func (s *DNSService) createNetwork() error {
	// Check if network already exists
	cmd := exec.Command("docker", "network", "ls", "--filter", fmt.Sprintf("name=%s", DNSNetworkName), "--format", "{{.Name}}")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check network: %w", err)
	}
	
	if strings.Contains(string(output), DNSNetworkName) {
		return nil // Network already exists
	}

	// Create network with custom subnet
	cmd = exec.Command("docker", "network", "create",
		"--driver", "bridge",
		"--subnet", "172.21.0.0/24",
		DNSNetworkName)
	
	return cmd.Run()
}

// reloadDNS gracefully reloads dnsmasq configuration without restarting container
func (s *DNSService) reloadDNS() error {
	if !s.IsRunning() {
		return fmt.Errorf("DNS container is not running")
	}

	// Send HUP signal to dnsmasq to reload configuration
	cmd := exec.Command("docker", "exec", DNSContainerName, "pkill", "-HUP", "dnsmasq")
	if err := cmd.Run(); err != nil {
		// If graceful reload fails, fall back to restart
		fmt.Println("⚠ Graceful reload failed, restarting container...")
		return s.restart()
	}

	// Reload nginx configuration
	cmd = exec.Command("docker", "exec", DNSContainerName, "nginx", "-s", "reload")
	if err := cmd.Run(); err != nil {
		fmt.Println("⚠ Nginx reload failed, restarting container...")
		return s.restart()
	}

	// Flush local DNS cache to ensure new domains resolve immediately
	s.flushDNSCache()

	return nil
}

// AddProject adds DNS configuration for a project
func (s *DNSService) AddProject(projectName string, services map[string]int) error {
	// Create DNS config - point to DNS container IP where nginx proxy runs
	dnsConfig := fmt.Sprintf("address=/%s.%s/%s\n", projectName, DNSDomain, DNSStaticIP)
	for serviceName := range services {
		if serviceName != "app" && serviceName != "webserver" {
			dnsConfig += fmt.Sprintf("address=/%s.%s.%s/%s\n", serviceName, projectName, DNSDomain, DNSStaticIP)
		}
	}

	dnsFile := filepath.Join(s.configDir, "projects", fmt.Sprintf("%s.dns", projectName))
	if err := os.WriteFile(dnsFile, []byte(dnsConfig), 0644); err != nil {
		return fmt.Errorf("failed to write DNS config: %w", err)
	}

	// Create nginx config
	nginxConfig := s.generateNginxConfig(projectName, services)
	nginxFile := filepath.Join(s.configDir, "projects", fmt.Sprintf("%s.nginx", projectName))
	if err := os.WriteFile(nginxFile, []byte(nginxConfig), 0644); err != nil {
		return fmt.Errorf("failed to write nginx config: %w", err)
	}

	// Gracefully reload dnsmasq configuration
	return s.reloadDNS()
}

// RemoveProject removes DNS configuration for a project
func (s *DNSService) RemoveProject(projectName string) error {
	dnsFile := filepath.Join(s.configDir, "projects", fmt.Sprintf("%s.dns", projectName))
	nginxFile := filepath.Join(s.configDir, "projects", fmt.Sprintf("%s.nginx", projectName))

	os.Remove(dnsFile)
	os.Remove(nginxFile)

	return s.reloadDNS()
}

// Status returns DNS system status
func (s *DNSService) Status() error {
	fmt.Println("DNS Configuration")
	fmt.Println(strings.Repeat("─", 50))

	// Check service status
	if s.IsRunning() {
		fmt.Println("✓ DNS service: running")
	} else {
		fmt.Println("✗ DNS service: not running")
	}

	// Check resolver
	resolverFile := "/etc/resolver/local"
	if _, err := os.Stat(resolverFile); err == nil {
		fmt.Println("✓ Resolver: configured")
	} else {
		fmt.Println("✗ Resolver: not configured")
	}

	// List projects
	projects, err := s.listProjects()
	if err != nil || len(projects) == 0 {
		fmt.Println("\nNo project domains configured")
	} else {
		fmt.Printf("\nActive Domains:\n")
		for _, project := range projects {
			fmt.Printf("  %s.%s\n", project, DNSDomain)
		}
	}

	return nil
}

// createConfigDirectories creates the DNS configuration structure
func (s *DNSService) createConfigDirectories() error {
	dirs := []string{
		s.configDir,
		filepath.Join(s.configDir, "projects"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create main dnsmasq config
	dnsmasqConfig := `# Atempo DNS Configuration
listen-address=0.0.0.0
port=53
bind-interfaces
no-hosts
cache-size=1000
domain-needed
bogus-priv
log-queries
log-facility=/var/log/dnsmasq.log
conf-dir=/etc/atempo/projects,*.dns
`

	configFile := filepath.Join(s.configDir, "dnsmasq.conf")
	return os.WriteFile(configFile, []byte(dnsmasqConfig), 0644)
}

// generateNginxConfig generates nginx configuration for a project
func (s *DNSService) generateNginxConfig(projectName string, services map[string]int) string {
	config := fmt.Sprintf(`# Nginx config for %s
server {
    listen 80;
    server_name %s.%s;
    
    location / {
        proxy_pass http://host.docker.internal:%d;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

`, projectName, projectName, DNSDomain, s.getMainPort(services))

	// Add service subdomains
	for serviceName, port := range services {
		if serviceName != "app" && serviceName != "webserver" {
			config += fmt.Sprintf(`server {
    listen 80;
    server_name %s.%s.%s;
    
    location / {
        proxy_pass http://host.docker.internal:%d;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

`, serviceName, projectName, DNSDomain, port)
		}
	}

	return config
}

// getMainPort returns the main web port for a project
func (s *DNSService) getMainPort(services map[string]int) int {
	// Look for main web service
	if port, exists := services["webserver"]; exists {
		return port
	}
	if port, exists := services["app"]; exists {
		return port
	}

	// Return first port as fallback
	for _, port := range services {
		return port
	}

	return 8000 // Default fallback
}

// restart restarts the DNS container
func (s *DNSService) restart() error {
	if s.IsRunning() {
		s.Stop()
	}
	return s.Start()
}

// remove removes the DNS container
func (s *DNSService) remove() {
	exec.Command("docker", "rm", "-f", DNSContainerName).Run()
}

// listProjects lists configured projects
func (s *DNSService) listProjects() ([]string, error) {
	projectsDir := filepath.Join(s.configDir, "projects")
	files, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, err
	}

	var projects []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".dns") {
			project := strings.TrimSuffix(file.Name(), ".dns")
			projects = append(projects, project)
		}
	}

	return projects, nil
}

// handlePortConflicts checks for and handles conflicts with existing nginx proxy
func (s *DNSService) handlePortConflicts() error {
	// Check if nginx proxy is running
	cmd := exec.Command("docker", "ps", "--filter", "name=atempo-nginx-proxy", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return nil // No conflict if command fails
	}

	if strings.Contains(string(output), "atempo-nginx-proxy") {
		fmt.Println("Found existing nginx proxy - stopping to avoid conflicts...")

		// Stop the existing nginx proxy
		cmd = exec.Command("docker", "stop", "atempo-nginx-proxy")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to stop existing nginx proxy: %w", err)
		}

		// Remove the container
		cmd = exec.Command("docker", "rm", "atempo-nginx-proxy")
		cmd.Run() // Ignore errors

		fmt.Println("✓ Existing nginx proxy stopped")
	}

	return nil
}

package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SimpleDNS manages a single container with dnsmasq + nginx for local development
type SimpleDNS struct {
	configDir string
}

const (
	SimpleDNSContainerName = "atempo-dns"
	SimpleDNSImage         = "nginx:alpine"
	DNSPort                = "5353"
	HTTPPort               = "80"
)

// NewSimpleDNS creates a new simple DNS manager
func NewSimpleDNS() *SimpleDNS {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".atempo", "dns")
	
	return &SimpleDNS{
		configDir: configDir,
	}
}

// Setup performs one-time DNS setup
func (s *SimpleDNS) Setup() error {
	fmt.Println("DNS Setup")
	fmt.Println(strings.Repeat("─", 50))
	
	// Check if already configured
	resolverFile := "/etc/resolver/local"
	if _, err := os.Stat(resolverFile); err == nil {
		fmt.Println("✓ DNS resolver already configured")
		if s.IsRunning() {
			fmt.Println("✓ DNS service running")
			return nil
		}
	}
	
	fmt.Println("This configures macOS to resolve .local domains")
	fmt.Println("through Atempo's DNS system.")
	fmt.Printf("\nConfigure DNS resolver? [y/N]: ")
	
	var response string
	fmt.Scanln(&response)
	
	if response != "y" && response != "Y" && response != "yes" {
		fmt.Println("DNS setup cancelled")
		return nil
	}
	
	return s.createResolver()
}

// createResolver creates the macOS DNS resolver configuration
func (s *SimpleDNS) createResolver() error {
	fmt.Println("Creating DNS resolver...")
	
	// Create resolver directory
	cmd := exec.Command("sudo", "mkdir", "-p", "/etc/resolver")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create resolver directory: %w", err)
	}
	
	// Create resolver config
	resolverConfig := `nameserver 127.0.0.1
port 5353`
	
	tempFile := filepath.Join(os.TempDir(), "atempo-resolver")
	if err := os.WriteFile(tempFile, []byte(resolverConfig), 0644); err != nil {
		return fmt.Errorf("failed to create resolver config: %w", err)
	}
	
	cmd = exec.Command("sudo", "mv", tempFile, "/etc/resolver/local")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install resolver: %w", err)
	}
	
	fmt.Println("✓ DNS resolver configured")
	
	// Start DNS service
	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to start DNS service: %w", err)
	}
	
	fmt.Println("✓ DNS service started")
	fmt.Println("✓ Setup complete - new projects will use custom domains")
	
	return nil
}

// Start starts the DNS container
func (s *SimpleDNS) Start() error {
	if s.IsRunning() {
		return nil // Already running
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

	// Create and start container
	cmd := exec.Command("docker", "run", "-d",
		"--name", SimpleDNSContainerName,
		"-p", fmt.Sprintf("%s:53/udp", DNSPort),
		"-p", fmt.Sprintf("%s:53/tcp", DNSPort),
		"-p", fmt.Sprintf("%s:80", HTTPPort),
		"-v", fmt.Sprintf("%s:/etc/atempo", s.configDir),
		"--restart", "unless-stopped",
		SimpleDNSImage,
		"/etc/atempo/startup.sh")
	
	return cmd.Run()
}

// Stop stops the DNS container
func (s *SimpleDNS) Stop() error {
	if !s.IsRunning() {
		return nil
	}
	
	cmd := exec.Command("docker", "stop", SimpleDNSContainerName)
	if err := cmd.Run(); err != nil {
		return err
	}
	
	s.remove()
	return nil
}

// IsRunning checks if the DNS container is running
func (s *SimpleDNS) IsRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", SimpleDNSContainerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), SimpleDNSContainerName)
}

// AddProject adds DNS configuration for a project
func (s *SimpleDNS) AddProject(projectName string, services map[string]int) error {
	// Create DNS config
	dnsConfig := fmt.Sprintf("address=/%s.local/127.0.0.1\n", projectName)
	for serviceName := range services {
		if serviceName != "app" && serviceName != "webserver" {
			dnsConfig += fmt.Sprintf("address=/%s.%s.local/127.0.0.1\n", serviceName, projectName)
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
	
	// Restart container to reload configs
	return s.restart()
}

// RemoveProject removes DNS configuration for a project
func (s *SimpleDNS) RemoveProject(projectName string) error {
	dnsFile := filepath.Join(s.configDir, "projects", fmt.Sprintf("%s.dns", projectName))
	nginxFile := filepath.Join(s.configDir, "projects", fmt.Sprintf("%s.nginx", projectName))
	
	os.Remove(dnsFile)
	os.Remove(nginxFile)
	
	return s.restart()
}

// Status returns DNS system status
func (s *SimpleDNS) Status() error {
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
			fmt.Printf("  %s.local\n", project)
		}
	}
	
	return nil
}

// createConfigDirectories creates the DNS configuration structure
func (s *SimpleDNS) createConfigDirectories() error {
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
no-hosts
cache-size=1000
domain-needed
bogus-priv
conf-dir=/etc/atempo/projects,*.dns
`
	
	configFile := filepath.Join(s.configDir, "dnsmasq.conf")
	return os.WriteFile(configFile, []byte(dnsmasqConfig), 0644)
}

// generateNginxConfig generates nginx configuration for a project
func (s *SimpleDNS) generateNginxConfig(projectName string, services map[string]int) string {
	config := fmt.Sprintf(`# Nginx config for %s
server {
    listen 80;
    server_name %s.local;
    
    location / {
        proxy_pass http://host.docker.internal:%d;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

`, projectName, projectName, s.getMainPort(services))
	
	// Add service subdomains
	for serviceName, port := range services {
		if serviceName != "app" && serviceName != "webserver" {
			config += fmt.Sprintf(`server {
    listen 80;
    server_name %s.%s.local;
    
    location / {
        proxy_pass http://host.docker.internal:%d;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

`, serviceName, projectName, port)
		}
	}
	
	return config
}

// getMainPort returns the main web port for a project
func (s *SimpleDNS) getMainPort(services map[string]int) int {
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
func (s *SimpleDNS) restart() error {
	if s.IsRunning() {
		s.Stop()
	}
	return s.Start()
}

// remove removes the DNS container
func (s *SimpleDNS) remove() {
	exec.Command("docker", "rm", "-f", SimpleDNSContainerName).Run()
}

// listProjects lists configured projects
func (s *SimpleDNS) listProjects() ([]string, error) {
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
func (s *SimpleDNS) handlePortConflicts() error {
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
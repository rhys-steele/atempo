package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DNSContainerManager manages a Docker-based DNSmasq container
type DNSContainerManager struct {
	mutex sync.RWMutex
}

const (
	// DNS container configuration
	DNSContainerName = "atempo-dnsmasq"
	DNSContainerPort = "5353" // Use non-standard port to avoid conflicts
	DNSConfigDir     = "dnsmasq"
	DNSImage         = "strm/dnsmasq:latest"
)

// NewDNSContainerManager creates a new DNS container manager
func NewDNSContainerManager() *DNSContainerManager {
	return &DNSContainerManager{}
}

// GetDNSConfigDir returns the DNSmasq configuration directory
func (dcm *DNSContainerManager) GetDNSConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".atempo", DNSConfigDir), nil
}

// GetDNSConfDir returns the DNSmasq conf.d directory
func (dcm *DNSContainerManager) GetDNSConfDir() (string, error) {
	configDir, err := dcm.GetDNSConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "conf.d"), nil
}

// EnsureDNSContainer ensures the DNSmasq container is running
func (dcm *DNSContainerManager) EnsureDNSContainer() error {
	dcm.mutex.Lock()
	defer dcm.mutex.Unlock()

	return dcm.ensureDNSContainerUnsafe()
}

// ensureDNSContainerUnsafe ensures the DNSmasq container is running without acquiring mutex
func (dcm *DNSContainerManager) ensureDNSContainerUnsafe() error {
	// Create configuration directories
	if err := dcm.createConfigDirectories(); err != nil {
		return fmt.Errorf("failed to create DNS config directories: %w", err)
	}

	// Check if DNS container is already running
	if dcm.isDNSContainerRunning() {
		return nil
	}

	// Start DNS container
	if err := dcm.startDNSContainer(); err != nil {
		return fmt.Errorf("failed to start DNS container: %w", err)
	}

	// Wait for container to be ready
	if err := dcm.waitForDNSContainer(); err != nil {
		return fmt.Errorf("DNS container failed to start properly: %w", err)
	}

	return nil
}

// createConfigDirectories creates DNSmasq configuration directories
func (dcm *DNSContainerManager) createConfigDirectories() error {
	configDir, err := dcm.GetDNSConfigDir()
	if err != nil {
		return err
	}
	
	confDir, err := dcm.GetDNSConfDir()
	if err != nil {
		return err
	}

	dirs := []string{configDir, confDir}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create default dnsmasq.conf if it doesn't exist
	dnsmasqConfPath := filepath.Join(configDir, "dnsmasq.conf")
	if _, err := os.Stat(dnsmasqConfPath); os.IsNotExist(err) {
		if err := dcm.createDefaultDNSmasqConfig(dnsmasqConfPath); err != nil {
			return fmt.Errorf("failed to create default dnsmasq config: %w", err)
		}
	}

	return nil
}

// createDefaultDNSmasqConfig creates the default DNSmasq configuration
func (dcm *DNSContainerManager) createDefaultDNSmasqConfig(configPath string) error {
	config := `# Atempo DNSmasq Configuration
# Listen on all interfaces inside the container
listen-address=0.0.0.0
port=53

# Don't read /etc/hosts
no-hosts

# Read additional configuration files from conf.d directory
conf-dir=/etc/dnsmasq.d

# Log queries for debugging (can be disabled in production)
log-queries

# Cache size
cache-size=1000

# Don't forward plain names (without domain)
domain-needed

# Never forward addresses in the non-routed address spaces
bogus-priv

# Enable DHCP authoritative mode (safe for local DNS)
dhcp-authoritative
`

	return os.WriteFile(configPath, []byte(config), 0644)
}

// isDNSContainerRunning checks if the DNS container is running
func (dcm *DNSContainerManager) isDNSContainerRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", DNSContainerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), DNSContainerName)
}

// startDNSContainer starts the DNSmasq container
func (dcm *DNSContainerManager) startDNSContainer() error {
	configDir, err := dcm.GetDNSConfigDir()
	if err != nil {
		return err
	}
	
	confDir, err := dcm.GetDNSConfDir()
	if err != nil {
		return err
	}

	// Remove existing container if it exists (but not running)
	dcm.removeExistingContainer()

	// Create dummy config file if conf.d is empty (dnsmasq needs at least one file)
	dummyConfFile := filepath.Join(confDir, "00-default.conf")
	if _, err := os.Stat(dummyConfFile); os.IsNotExist(err) {
		dummyConfig := "# Default Atempo DNS configuration\n# Additional project configs will be added here\n"
		if err := os.WriteFile(dummyConfFile, []byte(dummyConfig), 0644); err != nil {
			return fmt.Errorf("failed to create dummy config: %w", err)
		}
	}

	cmd := exec.Command("docker", "run", "-d",
		"--name", DNSContainerName,
		"--network", NginxProxyNetwork,
		"-p", fmt.Sprintf("%s:53/udp", DNSContainerPort),
		"-p", fmt.Sprintf("%s:53/tcp", DNSContainerPort),
		"-v", fmt.Sprintf("%s/dnsmasq.conf:/etc/dnsmasq.conf:ro", configDir),
		"-v", fmt.Sprintf("%s:/etc/dnsmasq.d:ro", confDir),
		"--restart", "unless-stopped",
		"--cap-add", "NET_ADMIN",
		DNSImage,
		"--conf-file=/etc/dnsmasq.conf")
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start DNS container: %w", err)
	}

	return nil
}

// removeExistingContainer removes an existing stopped container
func (dcm *DNSContainerManager) removeExistingContainer() {
	// Check if container exists (running or stopped)
	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", DNSContainerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil || !strings.Contains(string(output), DNSContainerName) {
		return // Container doesn't exist
	}

	// Stop and remove container
	exec.Command("docker", "stop", DNSContainerName).Run()
	exec.Command("docker", "rm", DNSContainerName).Run()
}

// waitForDNSContainer waits for the DNS container to be ready
func (dcm *DNSContainerManager) waitForDNSContainer() error {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		if dcm.isDNSContainerHealthy() {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("DNS container did not become healthy within 30 seconds")
}

// isDNSContainerHealthy checks if the DNS container is healthy
func (dcm *DNSContainerManager) isDNSContainerHealthy() bool {
	// Check if container is running
	if !dcm.isDNSContainerRunning() {
		return false
	}

	// Test DNS resolution from within the container
	cmd := exec.Command("docker", "exec", DNSContainerName, "nslookup", "localhost")
	return cmd.Run() == nil
}

// AddProjectConfig adds DNSmasq configuration for a project
func (dcm *DNSContainerManager) AddProjectConfig(projectName string, domains []string) error {
	dcm.mutex.Lock()
	defer dcm.mutex.Unlock()

	// Ensure DNS container is running (use unsafe version since we already have the mutex)
	if err := dcm.ensureDNSContainerUnsafe(); err != nil {
		return fmt.Errorf("failed to ensure DNS container: %w", err)
	}

	// Generate configuration content
	var config strings.Builder
	config.WriteString(fmt.Sprintf("# Atempo DNS configuration for %s\n", projectName))
	
	for _, domain := range domains {
		config.WriteString(fmt.Sprintf("address=/%s/127.0.0.1\n", domain))
	}

	// Write configuration file
	confDir, err := dcm.GetDNSConfDir()
	if err != nil {
		return err
	}
	
	configFile := filepath.Join(confDir, fmt.Sprintf("%s.conf", projectName))
	if err := os.WriteFile(configFile, []byte(config.String()), 0644); err != nil {
		return fmt.Errorf("failed to write DNS config: %w", err)
	}

	// Restart DNS container to reload configuration
	return dcm.restartDNSContainer()
}

// RemoveProjectConfig removes DNSmasq configuration for a project
func (dcm *DNSContainerManager) RemoveProjectConfig(projectName string) error {
	dcm.mutex.Lock()
	defer dcm.mutex.Unlock()

	confDir, err := dcm.GetDNSConfDir()
	if err != nil {
		return err
	}
	
	configFile := filepath.Join(confDir, fmt.Sprintf("%s.conf", projectName))
	
	// Remove config file
	if err := os.Remove(configFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove DNS config: %w", err)
	}

	// Restart DNS container to reload configuration
	return dcm.restartDNSContainer()
}

// restartDNSContainer restarts the DNS container
func (dcm *DNSContainerManager) restartDNSContainer() error {
	if !dcm.isDNSContainerRunning() {
		return dcm.startDNSContainer()
	}

	// Force stop and remove, then recreate to ensure volume mounts are fresh
	dcm.removeExistingContainer()
	return dcm.startDNSContainer()
}

// StopDNSContainer stops the DNS container
func (dcm *DNSContainerManager) StopDNSContainer() error {
	dcm.mutex.Lock()
	defer dcm.mutex.Unlock()

	if !dcm.isDNSContainerRunning() {
		return nil // Already stopped
	}

	cmd := exec.Command("docker", "stop", DNSContainerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop DNS container: %w", err)
	}

	cmd = exec.Command("docker", "rm", DNSContainerName)
	return cmd.Run()
}

// GetDNSContainerStatus returns the status of the DNS container
func (dcm *DNSContainerManager) GetDNSContainerStatus() (bool, error) {
	running := dcm.isDNSContainerRunning()
	
	if !running {
		return false, nil
	}

	healthy := dcm.isDNSContainerHealthy()
	return healthy, nil
}

// IsDockerAvailable checks if Docker is available
func (dcm *DNSContainerManager) IsDockerAvailable() bool {
	cmd := exec.Command("docker", "version")
	return cmd.Run() == nil
}

// GetDNSContainerLogs gets logs from the DNS container
func (dcm *DNSContainerManager) GetDNSContainerLogs(lines int) (string, error) {
	if !dcm.isDNSContainerRunning() {
		return "", fmt.Errorf("DNS container is not running")
	}

	cmd := exec.Command("docker", "logs", "--tail", fmt.Sprintf("%d", lines), DNSContainerName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}

	return string(output), nil
}

// GetInfrastructureConfigDir returns the infrastructure configuration directory
func (dcm *DNSContainerManager) GetInfrastructureConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".atempo", "infrastructure"), nil
}

// SetupInfrastructure creates the global Atempo infrastructure
func (dcm *DNSContainerManager) SetupInfrastructure() error {
	infraDir, err := dcm.GetInfrastructureConfigDir()
	if err != nil {
		return err
	}

	// Create infrastructure directory
	if err := os.MkdirAll(infraDir, 0755); err != nil {
		return fmt.Errorf("failed to create infrastructure directory: %w", err)
	}

	// Copy infrastructure docker-compose template
	templatePath, err := dcm.getInfrastructureTemplatePath()
	if err != nil {
		return fmt.Errorf("failed to get infrastructure template path: %w", err)
	}

	composePath := filepath.Join(infraDir, "docker-compose.yml")
	
	// Read template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read infrastructure template: %w", err)
	}

	// Write compose file
	if err := os.WriteFile(composePath, templateContent, 0644); err != nil {
		return fmt.Errorf("failed to write infrastructure compose file: %w", err)
	}

	return nil
}

// getInfrastructureTemplatePath returns the path to the infrastructure template
func (dcm *DNSContainerManager) getInfrastructureTemplatePath() (string, error) {
	// Get current working directory and navigate to project root
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// Navigate to project root (look for go.mod)
	projectRoot := wd
	for {
		goModPath := filepath.Join(projectRoot, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return "", fmt.Errorf("could not find project root (go.mod)")
		}
		projectRoot = parent
	}
	
	return filepath.Join(projectRoot, "templates", "infrastructure", "atempo-infrastructure.yml"), nil
}

// StartInfrastructure starts the global Atempo infrastructure
func (dcm *DNSContainerManager) StartInfrastructure() error {
	infraDir, err := dcm.GetInfrastructureConfigDir()
	if err != nil {
		return err
	}

	// Ensure infrastructure is set up
	if err := dcm.SetupInfrastructure(); err != nil {
		return fmt.Errorf("failed to setup infrastructure: %w", err)
	}

	// Start services with docker-compose
	cmd := exec.Command("docker-compose", "-f", filepath.Join(infraDir, "docker-compose.yml"), "up", "-d")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start infrastructure: %w", err)
	}

	return nil
}

// StopInfrastructure stops the global Atempo infrastructure
func (dcm *DNSContainerManager) StopInfrastructure() error {
	infraDir, err := dcm.GetInfrastructureConfigDir()
	if err != nil {
		return err
	}

	composePath := filepath.Join(infraDir, "docker-compose.yml")
	
	// Check if compose file exists
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return nil // Nothing to stop
	}

	// Stop services with docker-compose
	cmd := exec.Command("docker-compose", "-f", composePath, "down")
	return cmd.Run()
}
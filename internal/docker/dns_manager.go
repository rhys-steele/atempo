package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// DNSManager handles local DNS routing for projects
type DNSManager struct {
	mutex            sync.RWMutex
	containerManager *DNSContainerManager
	useDocker        bool
}

// ProjectDNS represents DNS configuration for a project
type ProjectDNS struct {
	ProjectName string            `json:"project_name"`
	Domain      string            `json:"domain"`      // e.g., "myapp.local"
	Services    map[string]string `json:"services"`    // service -> subdomain mapping
	TLD         string            `json:"tld"`         // top-level domain (.local, .dev, etc.)
}

const (
	// Default TLD for local development
	DefaultTLD = "local"
	
	// dnsmasq configuration directory
	DnsmasqConfigDir = "/opt/homebrew/etc/dnsmasq.d"
	DnsmasqAltConfigDir = "/usr/local/etc/dnsmasq.d"
	
	// Resolver directory for macOS
	ResolverDir = "/etc/resolver"
)

// NewDNSManager creates a new DNS manager
func NewDNSManager() *DNSManager {
	containerManager := NewDNSContainerManager()
	useDocker := containerManager.IsDockerAvailable()
	
	return &DNSManager{
		containerManager: containerManager,
		useDocker:        useDocker,
	}
}

// SetupProjectDNS configures DNS routing for a project with nginx proxy integration
func (dm *DNSManager) SetupProjectDNS(projectName string, servicePorts map[string]map[int]int) (*ProjectDNS, error) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	domain := fmt.Sprintf("%s.%s", projectName, DefaultTLD)
	
	projectDNS := &ProjectDNS{
		ProjectName: projectName,
		Domain:      domain,
		Services:    make(map[string]string),
		TLD:         DefaultTLD,
	}

	// Create nginx proxy manager
	nginxProxy := NewNginxProxyManager()
	
	// Prepare service mappings for nginx proxy
	var serviceMappings []ServiceMapping
	
	// Create service subdomains and mappings
	for serviceName, portMapping := range servicePorts {
		if isWebService(serviceName, portMapping) {
			// For services, use service.project.local format
			var subdomain string
			if IsMainWebService(serviceName) {
				// Main web service gets both project.local and service.project.local
				subdomain = domain
				projectDNS.Services[serviceName] = subdomain
				
				// Also add service-specific subdomain
				serviceSubdomain := fmt.Sprintf("%s.%s", serviceName, domain)
				projectDNS.Services[serviceName+"_service"] = serviceSubdomain
			} else {
				// Other services get service.project.local
				subdomain = fmt.Sprintf("%s.%s", serviceName, domain)
				projectDNS.Services[serviceName] = subdomain
			}
			
			// Find the main port for this service
			var mainPort int
			for internalPort, externalPort := range portMapping {
				if internalPort == 80 || internalPort == 443 || internalPort == 8000 {
					mainPort = externalPort
					break
				}
			}
			
			// If no standard web port found, use the first port
			if mainPort == 0 {
				for _, externalPort := range portMapping {
					mainPort = externalPort
					break
				}
			}
			
			if mainPort > 0 {
				serviceMappings = append(serviceMappings, ServiceMapping{
					ServiceName:  serviceName,
					Domain:       subdomain,
					Port:         mainPort,
					InternalPort: 80, // This will be the internal port inside nginx
					IsMain:       IsMainWebService(serviceName),
				})
			}
		}
	}

	// Configure nginx proxy with service mappings
	if len(serviceMappings) > 0 {
		if err := nginxProxy.AddProjectConfig(projectName, serviceMappings); err != nil {
			return nil, fmt.Errorf("failed to configure nginx proxy: %w", err)
		}
		
		// Connect project to proxy network
		if err := nginxProxy.ConnectProjectToProxy(projectName); err != nil {
			fmt.Printf("Warning: failed to connect project to proxy network: %v\n", err)
		}
	}

	// Configure system DNS (points domains to nginx proxy on localhost:80)
	if err := dm.configureDNS(projectDNS); err != nil {
		return nil, fmt.Errorf("failed to configure DNS: %w", err)
	}

	if err := dm.configureResolver(projectDNS); err != nil {
		return nil, fmt.Errorf("failed to configure resolver: %w", err)
	}

	return projectDNS, nil
}

// configureDNS sets up DNS configuration for the project (Docker or local)
func (dm *DNSManager) configureDNS(projectDNS *ProjectDNS) error {
	// Collect all domains for this project
	domains := []string{projectDNS.Domain}
	for _, subdomain := range projectDNS.Services {
		domains = append(domains, subdomain)
	}

	// Try Docker DNSmasq first if available
	if dm.useDocker {
		if err := dm.containerManager.AddProjectConfig(projectDNS.ProjectName, domains); err != nil {
			fmt.Printf("Warning: Docker DNS setup failed, falling back to local: %v\n", err)
			return dm.configureDNSmasqLocal(projectDNS)
		}
		return nil
	}

	// Fall back to local DNSmasq
	return dm.configureDNSmasqLocal(projectDNS)
}

// configureDNSmasqLocal sets up local dnsmasq configuration for the project
func (dm *DNSManager) configureDNSmasqLocal(projectDNS *ProjectDNS) error {
	// Find dnsmasq config directory
	configDir := DnsmasqConfigDir
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = DnsmasqAltConfigDir
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			return fmt.Errorf("dnsmasq not found. Install with: brew install dnsmasq")
		}
	}

	configFile := filepath.Join(configDir, fmt.Sprintf("%s.conf", projectDNS.ProjectName))
	
	// Create configuration content
	var config strings.Builder
	config.WriteString(fmt.Sprintf("# Atempo DNS configuration for %s\n", projectDNS.ProjectName))
	config.WriteString(fmt.Sprintf("address=/%s/127.0.0.1\n", projectDNS.Domain))
	
	for _, subdomain := range projectDNS.Services {
		config.WriteString(fmt.Sprintf("address=/%s/127.0.0.1\n", subdomain))
	}

	// Write configuration file
	if err := os.WriteFile(configFile, []byte(config.String()), 0644); err != nil {
		return fmt.Errorf("failed to write dnsmasq config: %w", err)
	}

	// Restart dnsmasq
	return dm.restartDNSmasq()
}

// configureResolver sets up macOS resolver for the TLD
func (dm *DNSManager) configureResolver(projectDNS *ProjectDNS) error {
	resolverFile := filepath.Join(ResolverDir, projectDNS.TLD)
	
	// Check if resolver already exists and is compatible
	if dm.isResolverConfigured(resolverFile) {
		return nil // Already configured
	}

	// Determine the DNS port based on whether we're using Docker
	port := "53"
	if dm.useDocker {
		port = DNSContainerPort
	}

	resolverConfig := fmt.Sprintf(`# Atempo DNS resolver
nameserver 127.0.0.1
port %s
`, port)

	// Write resolver configuration (requires sudo)
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("atempo-resolver-%s", projectDNS.TLD))
	if err := os.WriteFile(tempFile, []byte(resolverConfig), 0644); err != nil {
		return fmt.Errorf("failed to create temp resolver file: %w", err)
	}

	// Move to resolver directory with sudo
	cmd := exec.Command("sudo", "mv", tempFile, resolverFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure resolver (requires sudo): %w", err)
	}

	return nil
}

// isResolverConfigured checks if the resolver is already properly configured
func (dm *DNSManager) isResolverConfigured(resolverFile string) bool {
	content, err := os.ReadFile(resolverFile)
	if err != nil {
		return false // File doesn't exist
	}

	expectedPort := "53"
	if dm.useDocker {
		expectedPort = DNSContainerPort
	}

	// Check if the resolver file contains the correct port
	return strings.Contains(string(content), fmt.Sprintf("port %s", expectedPort))
}

// restartDNSmasq restarts the dnsmasq service
func (dm *DNSManager) restartDNSmasq() error {
	// Try different methods to restart dnsmasq
	commands := [][]string{
		{"brew", "services", "restart", "dnsmasq"},
		{"sudo", "brew", "services", "restart", "dnsmasq"},
		{"sudo", "killall", "-HUP", "dnsmasq"},
	}

	var lastErr error
	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	return fmt.Errorf("failed to restart dnsmasq: %w", lastErr)
}

// RemoveProjectDNS removes DNS configuration for a project
func (dm *DNSManager) RemoveProjectDNS(projectName string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Try Docker DNS removal first if using Docker
	if dm.useDocker {
		if err := dm.containerManager.RemoveProjectConfig(projectName); err != nil {
			fmt.Printf("Warning: Docker DNS removal failed, falling back to local: %v\n", err)
			return dm.removeProjectDNSLocal(projectName)
		}
		return nil
	}

	// Fall back to local DNSmasq removal
	return dm.removeProjectDNSLocal(projectName)
}

// removeProjectDNSLocal removes local dnsmasq configuration for a project
func (dm *DNSManager) removeProjectDNSLocal(projectName string) error {
	// Remove dnsmasq configuration
	configDirs := []string{DnsmasqConfigDir, DnsmasqAltConfigDir}
	
	for _, configDir := range configDirs {
		configFile := filepath.Join(configDir, fmt.Sprintf("%s.conf", projectName))
		if _, err := os.Stat(configFile); err == nil {
			if err := os.Remove(configFile); err != nil {
				return fmt.Errorf("failed to remove dnsmasq config: %w", err)
			}
		}
	}

	// Restart dnsmasq to apply changes
	return dm.restartDNSmasq()
}

// GetProjectDomain returns the primary domain for a project
func (dm *DNSManager) GetProjectDomain(projectName string) string {
	return fmt.Sprintf("%s.%s", projectName, DefaultTLD)
}

// GetServiceDomain returns the domain for a specific service
func (dm *DNSManager) GetServiceDomain(projectName, serviceName string) string {
	domain := dm.GetProjectDomain(projectName)
	return fmt.Sprintf("%s.%s", serviceName, domain)
}

// isWebService determines if a service serves web traffic
func isWebService(serviceName string, portMapping map[int]int) bool {
	// Check service name
	webServices := []string{"web", "webserver", "nginx", "apache", "app", "frontend", "ui"}
	serviceLower := strings.ToLower(serviceName)
	
	for _, webService := range webServices {
		if strings.Contains(serviceLower, webService) {
			return true
		}
	}

	// Check if it exposes common web ports
	webPorts := []int{80, 443, 8000, 8080, 3000, 4000, 5000, 9000}
	for internalPort := range portMapping {
		for _, webPort := range webPorts {
			if internalPort == webPort {
				return true
			}
		}
	}

	return false
}

// CheckDNSInstallation checks if DNS (Docker or local) is installed and configured
func (dm *DNSManager) CheckDNSInstallation() error {
	// Check Docker DNS first if available
	if dm.useDocker {
		if dm.containerManager.IsDockerAvailable() {
			return nil // Docker is available, no additional setup needed
		}
		return fmt.Errorf("Docker is not available for DNS container")
	}

	// Fall back to local DNSmasq check
	return dm.checkLocalDNSmasqInstallation()
}

// checkLocalDNSmasqInstallation checks if local dnsmasq is installed and configured
func (dm *DNSManager) checkLocalDNSmasqInstallation() error {
	// Check if dnsmasq is installed
	if err := exec.Command("which", "dnsmasq").Run(); err != nil {
		return fmt.Errorf("dnsmasq not installed. Install with: brew install dnsmasq")
	}

	// Check if config directory exists
	configDirs := []string{DnsmasqConfigDir, DnsmasqAltConfigDir}
	found := false
	
	for _, configDir := range configDirs {
		if _, err := os.Stat(configDir); err == nil {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("dnsmasq config directory not found. Ensure dnsmasq is properly configured")
	}

	return nil
}

// CheckDNSmasqInstallation checks if dnsmasq is installed and configured (deprecated, use CheckDNSInstallation)
func (dm *DNSManager) CheckDNSmasqInstallation() error {
	return dm.CheckDNSInstallation()
}

// SetupDNS provides instructions for DNS setup (Docker or local)
func (dm *DNSManager) SetupDNS() error {
	fmt.Println("Setting up DNS for local development...")
	
	// Check if already configured
	if err := dm.CheckDNSInstallation(); err == nil {
		if dm.useDocker {
			fmt.Println("✓ Docker DNS container is available and configured")
		} else {
			fmt.Println("✓ Local dnsmasq is already installed and configured")
		}
		return nil
	}

	// Show setup instructions based on availability
	if dm.containerManager.IsDockerAvailable() {
		fmt.Println(`
Docker DNS Setup (Recommended):
✓ Docker is available - Atempo will automatically manage DNS using containers.

No additional setup required! Atempo will:
1. Create and manage a DNSmasq container automatically
2. Configure DNS routing for project domains
3. Integrate with the nginx proxy system

This approach requires:
- Docker installed and running
- One-time sudo access for macOS resolver configuration

Benefits:
- No local DNSmasq installation needed
- Automatic container management
- Isolated DNS configuration
`)
	} else {
		fmt.Println(`
Local DNSmasq Setup:
Docker is not available. To enable local DNS routing, please run:

1. Install dnsmasq:
   brew install dnsmasq

2. Start dnsmasq service:
   sudo brew services start dnsmasq

3. Configure dnsmasq to start on boot:
   brew services start dnsmasq

4. Create resolver directory (if it doesn't exist):
   sudo mkdir -p /etc/resolver
`)
	}

	fmt.Println(`
This will allow atempo projects to use custom domains like:
- myproject.local
- web.myproject.local
- api.myproject.local

Note: Initial setup requires sudo access for resolver configuration.
`)

	return fmt.Errorf("DNS setup required")
}

// SetupDNSmasq provides instructions for dnsmasq setup (deprecated, use SetupDNS)
func (dm *DNSManager) SetupDNSmasq() error {
	return dm.SetupDNS()
}

// GenerateNginxProxy generates an nginx proxy configuration for routing domains to ports
func (dm *DNSManager) GenerateNginxProxy(projectName string, servicePorts map[string]map[int]int) (string, error) {
	var config strings.Builder
	
	config.WriteString(fmt.Sprintf("# Nginx proxy configuration for %s\n", projectName))
	config.WriteString("# This configuration routes custom domains to specific ports\n\n")

	for serviceName, portMapping := range servicePorts {
		if !isWebService(serviceName, portMapping) {
			continue
		}

		domain := dm.GetServiceDomain(projectName, serviceName)
		
		// Find the web port for this service
		var port int
		for internalPort, externalPort := range portMapping {
			if internalPort == 80 || internalPort == 8000 || internalPort == 3000 {
				port = externalPort
				break
			}
		}

		if port == 0 {
			// Use the first available port
			for _, externalPort := range portMapping {
				port = externalPort
				break
			}
		}

		if port > 0 {
			config.WriteString(fmt.Sprintf(`server {
    listen 80;
    server_name %s;
    
    location / {
        proxy_pass http://localhost:%d;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

`, domain, port))
		}
	}

	return config.String(), nil
}

// ListProjectDomains returns all configured domains for projects
func (dm *DNSManager) ListProjectDomains() (map[string][]string, error) {
	domains := make(map[string][]string)
	
	// Read dnsmasq config files
	configDirs := []string{DnsmasqConfigDir, DnsmasqAltConfigDir}
	
	for _, configDir := range configDirs {
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			continue
		}
		
		files, err := filepath.Glob(filepath.Join(configDir, "*.conf"))
		if err != nil {
			continue
		}
		
		for _, file := range files {
			projectName := strings.TrimSuffix(filepath.Base(file), ".conf")
			
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			
			var projectDomains []string
			lines := strings.Split(string(content), "\n")
			
			for _, line := range lines {
				if strings.HasPrefix(line, "address=/") && strings.Contains(line, "/127.0.0.1") {
					domain := strings.TrimPrefix(line, "address=/")
					domain = strings.TrimSuffix(domain, "/127.0.0.1")
					if domain != "" {
						projectDomains = append(projectDomains, domain)
					}
				}
			}
			
			if len(projectDomains) > 0 {
				domains[projectName] = projectDomains
			}
		}
	}
	
	return domains, nil
}

// GetDNSStatus returns the status of the DNS system (Docker or local)
func (dm *DNSManager) GetDNSStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})
	
	if dm.useDocker {
		status["type"] = "docker"
		status["container_name"] = DNSContainerName
		status["port"] = DNSContainerPort
		
		running, err := dm.containerManager.GetDNSContainerStatus()
		status["running"] = running
		if err != nil {
			status["error"] = err.Error()
		}
		
		if running {
			// Get container logs for debugging
			logs, err := dm.containerManager.GetDNSContainerLogs(10)
			if err == nil && logs != "" {
				status["recent_logs"] = logs
			}
		}
	} else {
		status["type"] = "local"
		status["port"] = "53"
		
		// Check if local dnsmasq is running
		cmd := exec.Command("pgrep", "dnsmasq")
		running := cmd.Run() == nil
		status["running"] = running
		
		// Check installation
		err := dm.checkLocalDNSmasqInstallation()
		if err != nil {
			status["error"] = err.Error()
		}
	}
	
	return status, nil
}

// StopDNS stops the DNS system (Docker container or local service)
func (dm *DNSManager) StopDNS() error {
	if dm.useDocker {
		return dm.containerManager.StopDNSContainer()
	}
	
	// For local DNSmasq, we don't stop the service as it may be used by other apps
	fmt.Println("Note: Local DNSmasq service not stopped as it may be used by other applications")
	return nil
}

// RestartDNS restarts the DNS system
func (dm *DNSManager) RestartDNS() error {
	if dm.useDocker {
		return dm.containerManager.restartDNSContainer()
	}
	
	return dm.restartDNSmasq()
}

// SetUseDocker allows manual override of Docker usage
func (dm *DNSManager) SetUseDocker(useDocker bool) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	dm.useDocker = useDocker
}
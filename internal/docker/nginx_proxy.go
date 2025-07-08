package docker

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

var nginxTemplates embed.FS

// NginxProxyManager manages a shared nginx reverse proxy for all projects
type NginxProxyManager struct {
	mutex sync.RWMutex
}

// ServiceMapping represents a service to port mapping
type ServiceMapping struct {
	ServiceName  string `json:"service_name"`
	Domain       string `json:"domain"`
	Port         int    `json:"port"`
	InternalPort int    `json:"internal_port"`
	IsMain       bool   `json:"is_main"`
}

// ProjectTemplateData represents data for project nginx template
type ProjectTemplateData struct {
	ProjectName string           `json:"project_name"`
	Services    []ServiceMapping `json:"services"`
}

const (
	// Nginx proxy configuration
	NginxProxyContainer = "atempo-nginx-proxy"
	NginxProxyNetwork   = "atempo-proxy-network"
)

// NewNginxProxyManager creates a new nginx proxy manager
func NewNginxProxyManager() *NginxProxyManager {
	return &NginxProxyManager{}
}

// GetNginxConfigDir returns the nginx configuration directory
func (npm *NginxProxyManager) GetNginxConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".atempo", "nginx"), nil
}

// GetNginxSitesDir returns the nginx sites-enabled directory
func (npm *NginxProxyManager) GetNginxSitesDir() (string, error) {
	configDir, err := npm.GetNginxConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "sites-enabled"), nil
}

// EnsureNginxProxy ensures the nginx proxy container is running
func (npm *NginxProxyManager) EnsureNginxProxy() error {
	npm.mutex.Lock()
	defer npm.mutex.Unlock()

	return npm.ensureNginxProxyUnsafe()
}

// ensureNginxProxyUnsafe ensures the nginx proxy container is running without acquiring mutex
func (npm *NginxProxyManager) ensureNginxProxyUnsafe() error {
	// Create nginx config directories
	if err := npm.createConfigDirectories(); err != nil {
		return fmt.Errorf("failed to create nginx config directories: %w", err)
	}

	// Check if nginx proxy is already running
	if npm.isNginxProxyRunning() {
		return nil
	}

	// Create docker network for proxy
	if err := npm.createProxyNetwork(); err != nil {
		return fmt.Errorf("failed to create proxy network: %w", err)
	}

	// Start nginx proxy container
	if err := npm.startNginxProxy(); err != nil {
		return fmt.Errorf("failed to start nginx proxy: %w", err)
	}

	return nil
}

// createConfigDirectories creates nginx configuration directories
func (npm *NginxProxyManager) createConfigDirectories() error {
	configDir, err := npm.GetNginxConfigDir()
	if err != nil {
		return err
	}
	
	sitesDir, err := npm.GetNginxSitesDir()
	if err != nil {
		return err
	}

	dirs := []string{configDir, sitesDir}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create default nginx.conf if it doesn't exist
	nginxConfPath := filepath.Join(configDir, "nginx.conf")
	if _, err := os.Stat(nginxConfPath); os.IsNotExist(err) {
		if err := npm.createDefaultNginxConfig(nginxConfPath); err != nil {
			return fmt.Errorf("failed to create default nginx config: %w", err)
		}
	}

	return nil
}

// createDefaultNginxConfig creates the main nginx configuration from template
func (npm *NginxProxyManager) createDefaultNginxConfig(configPath string) error {
	// Get the template path relative to the project root
	templatePath, err := npm.getNginxTemplatePath("nginx.conf.tmpl")
	if err != nil {
		return fmt.Errorf("failed to get template path: %w", err)
	}

	// Read template file
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read nginx template: %w", err)
	}

	// For the main nginx.conf, we don't need template variables, so just write it directly
	return os.WriteFile(configPath, tmplContent, 0644)
}

// getNginxTemplatePath returns the path to a nginx template file
func (npm *NginxProxyManager) getNginxTemplatePath(templateName string) (string, error) {
	// Get current working directory and go up to project root
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
	
	return filepath.Join(projectRoot, "templates", "nginx", templateName), nil
}

// isNginxProxyRunning checks if the nginx proxy container is running
func (npm *NginxProxyManager) isNginxProxyRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", NginxProxyContainer), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), NginxProxyContainer)
}

// createProxyNetwork creates the shared proxy network
func (npm *NginxProxyManager) createProxyNetwork() error {
	// Check if network already exists
	cmd := exec.Command("docker", "network", "ls", "--filter", fmt.Sprintf("name=%s", NginxProxyNetwork), "--format", "{{.Name}}")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), NginxProxyNetwork) {
		return nil // Network already exists
	}

	// Create the network
	cmd = exec.Command("docker", "network", "create", NginxProxyNetwork)
	return cmd.Run()
}

// startNginxProxy starts the nginx proxy container
func (npm *NginxProxyManager) startNginxProxy() error {
	configDir, err := npm.GetNginxConfigDir()
	if err != nil {
		return err
	}
	
	sitesDir, err := npm.GetNginxSitesDir()
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "run", "-d",
		"--name", NginxProxyContainer,
		"--network", NginxProxyNetwork,
		"-p", "80:80",
		"-p", "443:443",
		"-v", fmt.Sprintf("%s/nginx.conf:/etc/nginx/nginx.conf:ro", configDir),
		"-v", fmt.Sprintf("%s:/etc/nginx/sites-enabled:ro", sitesDir),
		"--restart", "unless-stopped",
		"nginx:alpine")
	
	return cmd.Run()
}

// AddProjectConfig adds nginx configuration for a project using templates
func (npm *NginxProxyManager) AddProjectConfig(projectName string, serviceMappings []ServiceMapping) error {
	npm.mutex.Lock()
	defer npm.mutex.Unlock()

	// Ensure nginx proxy is running (use unsafe version since we already have the mutex)
	if err := npm.ensureNginxProxyUnsafe(); err != nil {
		return fmt.Errorf("failed to ensure nginx proxy: %w", err)
	}

	// Generate nginx config for the project using templates
	configContent, err := npm.generateProjectConfigFromTemplate(projectName, serviceMappings)
	if err != nil {
		return fmt.Errorf("failed to generate project config: %w", err)
	}

	// Write config file
	sitesDir, err := npm.GetNginxSitesDir()
	if err != nil {
		return err
	}
	
	configFile := filepath.Join(sitesDir, fmt.Sprintf("%s.conf", projectName))
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write nginx config: %w", err)
	}

	// Reload nginx configuration
	return npm.reloadNginx()
}

// generateProjectConfigFromTemplate generates nginx configuration using the template
func (npm *NginxProxyManager) generateProjectConfigFromTemplate(projectName string, serviceMappings []ServiceMapping) (string, error) {
	// Get the template path
	templatePath, err := npm.getNginxTemplatePath("project.conf.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to get template path: %w", err)
	}

	// Load project template
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read project template: %w", err)
	}

	// Parse template
	tmpl, err := template.New("project").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	templateData := ProjectTemplateData{
		ProjectName: projectName,
		Services:    serviceMappings,
	}

	// Execute template
	var result strings.Builder
	if err := tmpl.Execute(&result, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// RemoveProjectConfig removes nginx configuration for a project
func (npm *NginxProxyManager) RemoveProjectConfig(projectName string) error {
	npm.mutex.Lock()
	defer npm.mutex.Unlock()

	sitesDir, err := npm.GetNginxSitesDir()
	if err != nil {
		return err
	}
	
	configFile := filepath.Join(sitesDir, fmt.Sprintf("%s.conf", projectName))
	
	// Remove config file
	if err := os.Remove(configFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove nginx config: %w", err)
	}

	// Reload nginx configuration
	return npm.reloadNginx()
}

// reloadNginx reloads the nginx configuration
func (npm *NginxProxyManager) reloadNginx() error {
	cmd := exec.Command("docker", "exec", NginxProxyContainer, "nginx", "-s", "reload")
	return cmd.Run()
}

// StopNginxProxy stops the nginx proxy container
func (npm *NginxProxyManager) StopNginxProxy() error {
	npm.mutex.Lock()
	defer npm.mutex.Unlock()

	cmd := exec.Command("docker", "stop", NginxProxyContainer)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop nginx proxy: %w", err)
	}

	cmd = exec.Command("docker", "rm", NginxProxyContainer)
	return cmd.Run()
}

// ConnectProjectToProxy connects a project's network to the proxy network
func (npm *NginxProxyManager) ConnectProjectToProxy(projectName string) error {
	projectNetwork := fmt.Sprintf("%s-network", projectName)
	
	// Check if project network exists
	cmd := exec.Command("docker", "network", "ls", "--filter", fmt.Sprintf("name=%s", projectNetwork), "--format", "{{.Name}}")
	output, err := cmd.Output()
	if err != nil || !strings.Contains(string(output), projectNetwork) {
		return fmt.Errorf("project network %s not found", projectNetwork)
	}

	// Connect networks (this allows communication between proxy and project containers)
	cmd = exec.Command("docker", "network", "connect", projectNetwork, NginxProxyContainer)
	
	// Ignore error if already connected
	err = cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("failed to connect project to proxy: %w", err)
	}

	return nil
}

// GetProxyStatus returns the status of the nginx proxy
func (npm *NginxProxyManager) GetProxyStatus() (bool, error) {
	running := npm.isNginxProxyRunning()
	
	if !running {
		return false, nil
	}

	// Check if nginx is healthy
	cmd := exec.Command("docker", "exec", NginxProxyContainer, "nginx", "-t")
	err := cmd.Run()
	
	return err == nil, err
}

// IsMainWebService determines if a service is the main web service
func IsMainWebService(serviceName string) bool {
	mainServices := []string{"webserver", "nginx", "apache", "web", "app", "frontend"}
	serviceLower := strings.ToLower(serviceName)
	
	for _, mainService := range mainServices {
		if strings.Contains(serviceLower, mainService) {
			return true
		}
	}
	return false
}
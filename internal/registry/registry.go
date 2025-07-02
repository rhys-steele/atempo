package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"atempo/internal/utils"
)

// Project represents a registered Atempo project
type Project struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Framework    string    `json:"framework"`
	Version      string    `json:"version"`
	CreatedAt    time.Time `json:"created_at"`
	LastAccessed time.Time `json:"last_accessed"`
	
	// Enhanced project state
	Status       string    `json:"status"`         // running/stopped/healthy/unhealthy
	Ports        []Port    `json:"ports"`
	URLs         []string  `json:"urls"`
	GitBranch    string    `json:"git_branch,omitempty"`
	GitStatus    string    `json:"git_status,omitempty"`
	Services     []Service `json:"services"`
}

// Port represents a port mapping for a service
type Port struct {
	Service   string `json:"service"`
	Internal  int    `json:"internal"`
	External  int    `json:"external"`
	Protocol  string `json:"protocol"`
}

// Service represents a Docker service with its status
type Service struct {
	Name    string `json:"name"`
	Status  string `json:"status"`  // running/stopped/healthy/unhealthy
	URL     string `json:"url,omitempty"`
}

// Registry manages the mapping of project names to paths
type Registry struct {
	Projects []Project `json:"projects"`
	Version  string    `json:"version"`
}

// GetRegistryPath returns the path to the registry file
func GetRegistryPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	atempoDir := filepath.Join(homeDir, ".atempo")
	if err := os.MkdirAll(atempoDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create atempo directory: %w", err)
	}

	return filepath.Join(atempoDir, "registry.json"), nil
}

// LoadRegistry loads the project registry from disk
func LoadRegistry() (*Registry, error) {
	registryPath, err := GetRegistryPath()
	if err != nil {
		return nil, err
	}

	// If registry doesn't exist, return empty registry
	if !utils.FileExists(registryPath) {
		return &Registry{
			Projects: []Project{},
			Version:  "1.0",
		}, nil
	}

	data, err := os.ReadFile(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	var registry Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	return &registry, nil
}

// SaveRegistry saves the project registry to disk
func (r *Registry) SaveRegistry() error {
	registryPath, err := GetRegistryPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize registry: %w", err)
	}

	return os.WriteFile(registryPath, data, 0644)
}

// AddProject adds a new project to the registry
func (r *Registry) AddProject(name, path, framework, version string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Check if project name already exists
	for i, project := range r.Projects {
		if project.Name == name {
			// Update existing project
			r.Projects[i] = Project{
				Name:         name,
				Path:         absPath,
				Framework:    framework,
				Version:      version,
				CreatedAt:    project.CreatedAt,
				LastAccessed: time.Now(),
			}
			return r.SaveRegistry()
		}
	}

	// Add new project
	project := Project{
		Name:         name,
		Path:         absPath,
		Framework:    framework,
		Version:      version,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}

	r.Projects = append(r.Projects, project)
	return r.SaveRegistry()
}

// FindProject finds a project by name
func (r *Registry) FindProject(name string) (*Project, error) {
	for i, project := range r.Projects {
		if project.Name == name {
			// Update last accessed time
			r.Projects[i].LastAccessed = time.Now()
			r.SaveRegistry() // Save updated access time
			return &r.Projects[i], nil
		}
	}

	return nil, fmt.Errorf("project '%s' not found in registry", name)
}

// ListProjects returns all registered projects
func (r *Registry) ListProjects() []Project {
	return r.Projects
}

// RemoveProject removes a project from the registry
func (r *Registry) RemoveProject(name string) error {
	for i, project := range r.Projects {
		if project.Name == name {
			r.Projects = append(r.Projects[:i], r.Projects[i+1:]...)
			return r.SaveRegistry()
		}
	}

	return fmt.Errorf("project '%s' not found in registry", name)
}

// ResolveProjectPath resolves a project identifier to an absolute path
// The identifier can be:
// - A project name (from registry)
// - A relative path
// - An absolute path
func ResolveProjectPath(identifier string) (string, error) {
	// If empty, use current directory
	if identifier == "" {
		return os.Getwd()
	}

	// Try to find by project name first
	registry, err := LoadRegistry()
	if err == nil {
		if project, err := registry.FindProject(identifier); err == nil {
			return project.Path, nil
		}
	}

	// If not found in registry, treat as path
	if filepath.IsAbs(identifier) {
		return identifier, nil
	}

	// Convert relative path to absolute
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	return filepath.Join(cwd, identifier), nil
}

// ScanForProjects scans a directory for Atempo projects and adds them to registry
func (r *Registry) ScanForProjects(scanPath string) error {
	return filepath.Walk(scanPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		if info.IsDir() && info.Name() == "atempo.json" {
			return nil // Skip directories named atempo.json
		}

		if !info.IsDir() && info.Name() == "atempo.json" {
			// Found a atempo.json file
			projectPath := filepath.Dir(path)
			if err := r.addProjectFromAtempoJson(projectPath); err != nil {
				fmt.Printf("Warning: Failed to add project from %s: %v\n", projectPath, err)
			}
		}

		return nil
	})
}

// addProjectFromAtempoJson reads a atempo.json file and adds the project to registry
func (r *Registry) addProjectFromAtempoJson(projectPath string) error {
	atempoJsonPath := filepath.Join(projectPath, "atempo.json")
	
	data, err := os.ReadFile(atempoJsonPath)
	if err != nil {
		return err
	}

	var config struct {
		Name      string `json:"name"`
		Framework string `json:"framework"`
		Version   string `json:"version"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// Generate a unique name if name is not set or is a template
	name := config.Name
	if name == "" || strings.Contains(name, "{{") {
		name = filepath.Base(projectPath)
	}

	return r.AddProject(name, projectPath, config.Framework, config.Version)
}

// CleanupInvalidProjects removes projects with non-existent paths
func (r *Registry) CleanupInvalidProjects() error {
	validProjects := []Project{}

	for _, project := range r.Projects {
		if utils.FileExists(project.Path) {
			validProjects = append(validProjects, project)
		} else {
			fmt.Printf("Removing invalid project: %s (path no longer exists: %s)\n", project.Name, project.Path)
		}
	}

	r.Projects = validProjects
	return r.SaveRegistry()
}

// UpdateProjectStatus updates the status and health information for a project
func (r *Registry) UpdateProjectStatus(name string) error {
	for i, project := range r.Projects {
		if project.Name == name {
			// Check project status
			status, services, ports, urls := r.checkProjectHealth(project.Path)
			
			r.Projects[i].Status = status
			r.Projects[i].Services = services
			r.Projects[i].Ports = ports
			r.Projects[i].URLs = urls
			r.Projects[i].LastAccessed = time.Now()
			
			// Update Git information if in a Git repository
			gitBranch, gitStatus := r.getGitInfo(project.Path)
			r.Projects[i].GitBranch = gitBranch
			r.Projects[i].GitStatus = gitStatus
			
			return r.SaveRegistry()
		}
	}
	
	return fmt.Errorf("project '%s' not found", name)
}

// UpdateAllProjectsStatus updates status for all registered projects
func (r *Registry) UpdateAllProjectsStatus() error {
	for i := range r.Projects {
		status, services, ports, urls := r.checkProjectHealth(r.Projects[i].Path)
		
		r.Projects[i].Status = status
		r.Projects[i].Services = services
		r.Projects[i].Ports = ports
		r.Projects[i].URLs = urls
		
		// Update Git information if in a Git repository
		gitBranch, gitStatus := r.getGitInfo(r.Projects[i].Path)
		r.Projects[i].GitBranch = gitBranch
		r.Projects[i].GitStatus = gitStatus
	}
	
	return r.SaveRegistry()
}

// checkProjectHealth checks the health of a project by examining Docker services
func (r *Registry) checkProjectHealth(projectPath string) (string, []Service, []Port, []string) {
	var services []Service
	var ports []Port
	var urls []string
	var overallStatus string = "stopped"

	// Check if docker-compose.yml exists
	composePath := filepath.Join(projectPath, "docker-compose.yml")
	if !utils.FileExists(composePath) {
		return "no-docker", services, ports, urls
	}

	// Run docker-compose ps to get service status
	cmd := exec.Command("docker-compose", "ps", "--format", "json")
	cmd.Dir = projectPath
	output, err := cmd.Output()
	if err != nil {
		return "docker-error", services, ports, urls
	}

	// Parse docker-compose ps output
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	runningServices := 0
	totalServices := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var serviceData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &serviceData); err != nil {
			continue
		}

		totalServices++
		serviceName := serviceData["Service"].(string)
		state := serviceData["State"].(string)
		
		// Determine service status
		var serviceStatus string
		switch state {
		case "running":
			serviceStatus = "running"
			runningServices++
		case "exited":
			serviceStatus = "stopped"
		default:
			serviceStatus = "unhealthy"
		}

		services = append(services, Service{
			Name:   serviceName,
			Status: serviceStatus,
		})

		// Extract port information if service is running
		if serviceStatus == "running" && serviceData["Publishers"] != nil {
			if publishers, ok := serviceData["Publishers"].([]interface{}); ok {
				for _, pub := range publishers {
					if pubMap, ok := pub.(map[string]interface{}); ok {
						if targetPort, ok := pubMap["TargetPort"].(float64); ok {
							if publishedPort, ok := pubMap["PublishedPort"].(float64); ok {
								ports = append(ports, Port{
									Service:  serviceName,
									Internal: int(targetPort),
									External: int(publishedPort),
									Protocol: "tcp",
								})

								// Generate URL for web services
								if isWebPort(int(publishedPort)) {
									url := fmt.Sprintf("http://localhost:%d", int(publishedPort))
									urls = append(urls, url)
									
									// Update service with URL
									for i := range services {
										if services[i].Name == serviceName {
											services[i].URL = url
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Determine overall status
	if totalServices == 0 {
		overallStatus = "no-services"
	} else if runningServices == totalServices {
		overallStatus = "running"
	} else if runningServices > 0 {
		overallStatus = "partial"
	} else {
		overallStatus = "stopped"
	}

	return overallStatus, services, ports, urls
}

// getGitInfo retrieves Git branch and status information
func (r *Registry) getGitInfo(projectPath string) (string, string) {
	// Check if it's a Git repository
	gitDir := filepath.Join(projectPath, ".git")
	if !utils.FileExists(gitDir) {
		return "", ""
	}

	// Get current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = projectPath
	branchOutput, err := cmd.Output()
	if err != nil {
		return "", ""
	}
	branch := strings.TrimSpace(string(branchOutput))

	// Get status
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = projectPath
	statusOutput, err := cmd.Output()
	if err != nil {
		return branch, ""
	}

	statusLines := strings.Split(strings.TrimSpace(string(statusOutput)), "\n")
	if len(statusLines) == 1 && statusLines[0] == "" {
		return branch, "clean"
	}

	// Count modifications
	modifications := 0
	for _, line := range statusLines {
		if strings.TrimSpace(line) != "" {
			modifications++
		}
	}

	if modifications > 0 {
		return branch, fmt.Sprintf("%d changes", modifications)
	}

	return branch, "clean"
}

// isWebPort determines if a port is typically used for web services
func isWebPort(port int) bool {
	webPorts := []int{80, 443, 3000, 4000, 5000, 8000, 8080, 8443, 9000}
	for _, webPort := range webPorts {
		if port == webPort {
			return true
		}
	}
	return false
}
package docker

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"atempo/internal/registry"
)

// PortManager handles dynamic port allocation for projects
type PortManager struct {
	mutex sync.RWMutex
}

// PortAllocation represents allocated ports for a project
type PortAllocation struct {
	ProjectName string             `json:"project_name"`
	Ports       map[string]int     `json:"ports"` // service_name:internal_port -> external_port
	Reserved    bool               `json:"reserved"`
}

// PortRegistry stores all port allocations
type PortRegistry struct {
	Allocations map[string]PortAllocation `json:"allocations"` // project_name -> allocation
	NextPort    int                       `json:"next_port"`
}

const (
	// Port range for dynamic allocation (avoiding common system ports)
	PortRangeStart = 3000
	PortRangeEnd   = 65535
	
	// Common web service ports to prioritize
	DefaultWebPort = 80
	DefaultHTTPSPort = 443
)

// NewPortManager creates a new port manager
func NewPortManager() *PortManager {
	return &PortManager{}
}

// GetPortRegistryPath returns the path to the port registry file
func GetPortRegistryPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	atempoDir := filepath.Join(homeDir, ".atempo")
	if err := os.MkdirAll(atempoDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create atempo directory: %w", err)
	}

	return filepath.Join(atempoDir, "ports.json"), nil
}

// LoadPortRegistry loads the port registry from disk
func (pm *PortManager) LoadPortRegistry() (*PortRegistry, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	registryPath, err := GetPortRegistryPath()
	if err != nil {
		return nil, err
	}

	// If registry doesn't exist, return empty registry
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		return &PortRegistry{
			Allocations: make(map[string]PortAllocation),
			NextPort:    PortRangeStart,
		}, nil
	}

	data, err := os.ReadFile(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read port registry: %w", err)
	}

	var registry PortRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse port registry: %w", err)
	}

	return &registry, nil
}

// SavePortRegistry saves the port registry to disk
func (pm *PortManager) SavePortRegistry(registry *PortRegistry) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	registryPath, err := GetPortRegistryPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize port registry: %w", err)
	}

	return os.WriteFile(registryPath, data, 0644)
}

// AllocatePortsForProject allocates ports for a project's services
func (pm *PortManager) AllocatePortsForProject(projectName string, servicePorts map[string][]int) (map[string]map[int]int, error) {
	registry, err := pm.LoadPortRegistry()
	if err != nil {
		return nil, err
	}

	// Check if project already has allocation
	if allocation, exists := registry.Allocations[projectName]; exists {
		return pm.convertAllocationToMapping(allocation, servicePorts), nil
	}

	// Allocate new ports
	allocation := PortAllocation{
		ProjectName: projectName,
		Ports:       make(map[string]int),
		Reserved:    true,
	}

	result := make(map[string]map[int]int)

	for serviceName, internalPorts := range servicePorts {
		serviceMapping := make(map[int]int)
		
		for _, internalPort := range internalPorts {
			externalPort, err := pm.findAvailablePort(registry, internalPort)
			if err != nil {
				return nil, fmt.Errorf("failed to allocate port for %s:%d: %w", serviceName, internalPort, err)
			}
			
			key := fmt.Sprintf("%s:%d", serviceName, internalPort)
			allocation.Ports[key] = externalPort
			serviceMapping[internalPort] = externalPort
		}
		
		result[serviceName] = serviceMapping
	}

	// Save allocation
	registry.Allocations[projectName] = allocation
	if err := pm.SavePortRegistry(registry); err != nil {
		return nil, fmt.Errorf("failed to save port allocation: %w", err)
	}

	return result, nil
}

// findAvailablePort finds an available port, preferring the suggested port if available
func (pm *PortManager) findAvailablePort(registry *PortRegistry, suggestedPort int) (int, error) {
	// Try the suggested port first (for common services like web:80 -> 8000)
	if pm.isPortAvailable(suggestedPort, registry) && suggestedPort >= PortRangeStart {
		return suggestedPort, nil
	}

	// Try common mappings for web services
	if suggestedPort == DefaultWebPort {
		commonWebPorts := []int{8000, 8001, 8002, 8003, 8004, 8005, 8080, 8081, 8082}
		for _, port := range commonWebPorts {
			if pm.isPortAvailable(port, registry) {
				return port, nil
			}
		}
	}

	// Sequential search from next available port
	for port := registry.NextPort; port <= PortRangeEnd; port++ {
		if pm.isPortAvailable(port, registry) {
			// Update next port to continue from here
			registry.NextPort = port + 1
			return port, nil
		}
	}

	return 0, fmt.Errorf("no available ports in range %d-%d", PortRangeStart, PortRangeEnd)
}

// isPortAvailable checks if a port is available (not in use by system or other projects)
func (pm *PortManager) isPortAvailable(port int, registry *PortRegistry) bool {
	// Check if port is already allocated to another project
	for _, allocation := range registry.Allocations {
		for _, allocatedPort := range allocation.Ports {
			if allocatedPort == port {
				return false
			}
		}
	}

	// Check if port is actually available on the system
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	conn.Close()
	
	return true
}

// convertAllocationToMapping converts stored allocation to the expected format
func (pm *PortManager) convertAllocationToMapping(allocation PortAllocation, servicePorts map[string][]int) map[string]map[int]int {
	result := make(map[string]map[int]int)

	for serviceName, internalPorts := range servicePorts {
		serviceMapping := make(map[int]int)
		
		for _, internalPort := range internalPorts {
			key := fmt.Sprintf("%s:%d", serviceName, internalPort)
			if externalPort, exists := allocation.Ports[key]; exists {
				serviceMapping[internalPort] = externalPort
			}
		}
		
		if len(serviceMapping) > 0 {
			result[serviceName] = serviceMapping
		}
	}

	return result
}

// ReleasePortsForProject releases all ports allocated to a project
func (pm *PortManager) ReleasePortsForProject(projectName string) error {
	registry, err := pm.LoadPortRegistry()
	if err != nil {
		return err
	}

	delete(registry.Allocations, projectName)
	
	return pm.SavePortRegistry(registry)
}

// GetProjectPorts returns the port allocation for a specific project
func (pm *PortManager) GetProjectPorts(projectName string) (*PortAllocation, error) {
	registry, err := pm.LoadPortRegistry()
	if err != nil {
		return nil, err
	}

	if allocation, exists := registry.Allocations[projectName]; exists {
		return &allocation, nil
	}

	return nil, fmt.Errorf("no port allocation found for project: %s", projectName)
}

// ListAllAllocations returns all current port allocations
func (pm *PortManager) ListAllAllocations() (map[string]PortAllocation, error) {
	registry, err := pm.LoadPortRegistry()
	if err != nil {
		return nil, err
	}

	return registry.Allocations, nil
}

// GetProjectURLs generates URLs for a project's web services
func (pm *PortManager) GetProjectURLs(projectName string) ([]string, error) {
	allocation, err := pm.GetProjectPorts(projectName)
	if err != nil {
		return nil, err
	}

	var urls []string
	webPorts := []int{}

	// Find web service ports
	for key, externalPort := range allocation.Ports {
		// Look for common web service ports
		if isWebServicePort(key, externalPort) {
			webPorts = append(webPorts, externalPort)
		}
	}

	// Sort ports to ensure consistent ordering
	sort.Ints(webPorts)

	// Generate URLs
	for _, port := range webPorts {
		urls = append(urls, fmt.Sprintf("http://localhost:%d", port))
	}

	return urls, nil
}

// isWebServicePort determines if a port mapping represents a web service
func isWebServicePort(serviceKey string, externalPort int) bool {
	// Common web service indicators
	webServices := []string{"web:", "webserver:", "nginx:", "apache:", "app:80", "frontend:"}
	
	for _, indicator := range webServices {
		if len(serviceKey) > len(indicator) && serviceKey[:len(indicator)] == indicator {
			return true
		}
	}

	// Common web ports
	webPorts := []int{80, 443, 8000, 8080, 3000, 4000, 5000, 9000}
	for _, webPort := range webPorts {
		if externalPort == webPort {
			return true
		}
	}

	return false
}

// CleanupOrphanedAllocations removes allocations for projects that no longer exist
func (pm *PortManager) CleanupOrphanedAllocations() error {
	portRegistry, err := pm.LoadPortRegistry()
	if err != nil {
		return err
	}

	// Get list of current projects
	projectRegistry, err := registry.LoadRegistry()
	if err != nil {
		return err
	}

	currentProjects := make(map[string]bool)
	for _, project := range projectRegistry.ListProjects() {
		currentProjects[project.Name] = true
	}

	// Remove allocations for non-existent projects
	modified := false
	for projectName := range portRegistry.Allocations {
		if !currentProjects[projectName] {
			delete(portRegistry.Allocations, projectName)
			modified = true
		}
	}

	if modified {
		return pm.SavePortRegistry(portRegistry)
	}

	return nil
}
package compose

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"atempo/internal/docker"
	"gopkg.in/yaml.v3"
)

// AtempoConfig represents the enhanced atempo.json structure
type AtempoConfig struct {
	Name      string             `json:"name"`
	Framework string             `json:"framework"`
	Language  string             `json:"language"`
	Services  map[string]Service `json:"services"`
	Volumes   map[string]Volume  `json:"volumes,omitempty"`
	Networks  map[string]Network `json:"networks,omitempty"`
	Version   string             `json:"version,omitempty"`
}

// Service represents a Docker service definition
type Service struct {
	Type        string            `json:"type"` // "image" or "build"
	Image       string            `json:"image,omitempty"`
	Dockerfile  string            `json:"dockerfile,omitempty"`
	Context     string            `json:"context,omitempty"`
	Command     interface{}       `json:"command,omitempty"` // string or []string
	WorkingDir  string            `json:"working_dir,omitempty"`
	Ports       []string          `json:"ports,omitempty"`
	Volumes     []string          `json:"volumes,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	DependsOn   []string          `json:"depends_on,omitempty"`
	Restart     string            `json:"restart,omitempty"`
	Networks    []string          `json:"networks,omitempty"`
}

// Volume represents a Docker volume definition
type Volume struct {
	Driver       string            `json:"driver,omitempty"`
	DriverOpts   map[string]string `json:"driver_opts,omitempty"`
	External     bool              `json:"external,omitempty"`
	ExternalName string            `json:"external_name,omitempty"`
}

// Network represents a Docker network definition
type Network struct {
	Driver     string            `json:"driver,omitempty"`
	DriverOpts map[string]string `json:"driver_opts,omitempty"`
	External   bool              `json:"external,omitempty"`
}

// DockerCompose represents the docker-compose.yml structure
type DockerCompose struct {
	Version  string                 `yaml:"version"`
	Services map[string]interface{} `yaml:"services"`
	Volumes  map[string]interface{} `yaml:"volumes,omitempty"`
	Networks map[string]interface{} `yaml:"networks,omitempty"`
}

// LoadAtempoConfig loads and parses the atempo.json file
func LoadAtempoConfig(projectPath string) (*AtempoConfig, error) {
	atempoJsonPath := filepath.Join(projectPath, "atempo.json")

	data, err := os.ReadFile(atempoJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read atempo.json: %w", err)
	}

	var config AtempoConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse atempo.json: %w", err)
	}

	return &config, nil
}

// GenerateDockerCompose generates a docker-compose.yml from atempo.json
func GenerateDockerCompose(projectPath string) error {
	return GenerateDockerComposeWithDynamicPorts(projectPath)
}

// GenerateDockerComposeWithDynamicPorts generates a docker-compose.yml with dynamic port allocation
func GenerateDockerComposeWithDynamicPorts(projectPath string) error {
	config, err := LoadAtempoConfig(projectPath)
	if err != nil {
		return err
	}

	// Extract project name from config name or use directory name
	projectName := config.Name
	if projectName == "" {
		projectName = filepath.Base(projectPath)
	}

	// Initialize port manager
	portManager := docker.NewPortManager()

	// Collect service port requirements
	servicePorts := make(map[string][]int)
	for serviceName, service := range config.Services {
		var ports []int
		for _, portMapping := range service.Ports {
			if internalPort := extractInternalPort(portMapping); internalPort > 0 {
				ports = append(ports, internalPort)
			}
		}
		if len(ports) > 0 {
			servicePorts[serviceName] = ports
		}
	}

	// Allocate ports dynamically
	allocatedPorts, err := portManager.AllocatePortsForProject(projectName, servicePorts)
	if err != nil {
		return fmt.Errorf("failed to allocate ports: %w", err)
	}

	// Setup DNS routing with simplified system
	dnsService := docker.NewDNSService()

	// Convert allocated ports to simple map for DNS system
	services := make(map[string]int)
	for serviceName, portMapping := range allocatedPorts {
		// Get the main external port for this service
		// For mailhog, prioritize the web interface port (8025) over SMTP port (1025)
		if serviceName == "mailhog" {
			// Look for web interface port (8025) first
			for internalPort, externalPort := range portMapping {
				if internalPort == 8025 {
					services[serviceName] = externalPort
					break
				}
			}
			// If no web interface port found, fall back to first port
			if _, exists := services[serviceName]; !exists {
				for _, externalPort := range portMapping {
					services[serviceName] = externalPort
					break
				}
			}
		} else {
			// For other services, take first port as main port
			for _, externalPort := range portMapping {
				services[serviceName] = externalPort
				break
			}
		}
	}

	if err := dnsService.AddProject(projectName, services); err != nil {
		// DNS setup is optional - continue without it
		// Log the error but don't fail the entire process
		// This is a business logic function - let the caller handle user messaging
	}

	compose := &DockerCompose{
		Version:  "3.8",
		Services: make(map[string]interface{}),
		Volumes:  make(map[string]interface{}),
		Networks: make(map[string]interface{}),
	}

	// Convert services with dynamic ports
	for serviceName, service := range config.Services {
		dockerService := convertServiceWithDynamicPorts(service, serviceName, projectName, config.Framework, allocatedPorts[serviceName])
		compose.Services[serviceName] = dockerService
	}

	// Convert volumes
	for volumeName, volume := range config.Volumes {
		compose.Volumes[volumeName] = convertVolume(volume)
	}

	// Convert networks
	for networkName, network := range config.Networks {
		compose.Networks[networkName] = convertNetwork(network)
	}

	// Add project-specific network if none specified
	if len(compose.Networks) == 0 {
		networkName := fmt.Sprintf("%s-network", projectName)
		compose.Networks[networkName] = map[string]interface{}{
			"driver": "bridge",
		}

		// Add network to all services
		for _, serviceInterface := range compose.Services {
			if serviceMap, ok := serviceInterface.(map[string]interface{}); ok {
				serviceMap["networks"] = []string{networkName}
			}
		}
	}

	// Write docker-compose.yml
	composePath := filepath.Join(projectPath, "docker-compose.yml")
	if err := writeDockerCompose(compose, composePath); err != nil {
		return err
	}

	// Generate access information
	if err := generateAccessInfo(projectPath, projectName, allocatedPorts, services); err != nil {
		// Access info generation is optional - continue without it
		// This is a business logic function - let the caller handle user messaging
	}

	return nil
}

// convertService converts a Atempo service to Docker Compose service
func convertService(service Service, serviceName, projectName, framework string) map[string]interface{} {
	dockerService := make(map[string]interface{})

	// Handle build vs image
	if service.Type == "build" {
		// Generate project-specific image name
		imageName := fmt.Sprintf("%s-%s-%s", projectName, framework, serviceName)
		dockerService["image"] = imageName

		if service.Context != "" {
			dockerService["build"] = map[string]interface{}{
				"context":    service.Context,
				"dockerfile": service.Dockerfile,
			}
		} else {
			dockerService["build"] = map[string]interface{}{
				"context":    ".",
				"dockerfile": service.Dockerfile,
			}
		}
	} else if service.Image != "" {
		dockerService["image"] = service.Image
	}

	// Add container name with project prefix
	dockerService["container_name"] = fmt.Sprintf("%s-%s", projectName, serviceName)

	// Add restart policy
	if service.Restart != "" {
		dockerService["restart"] = service.Restart
	} else {
		dockerService["restart"] = "unless-stopped"
	}

	// Add optional fields
	if service.Command != nil {
		dockerService["command"] = service.Command
	}

	if service.WorkingDir != "" {
		dockerService["working_dir"] = service.WorkingDir
	}

	if len(service.Ports) > 0 {
		dockerService["ports"] = service.Ports
	}

	if len(service.Volumes) > 0 {
		dockerService["volumes"] = service.Volumes
	}

	if len(service.Environment) > 0 {
		dockerService["environment"] = service.Environment
	}

	if len(service.DependsOn) > 0 {
		dockerService["depends_on"] = service.DependsOn
	}

	if len(service.Networks) > 0 {
		dockerService["networks"] = service.Networks
	}

	return dockerService
}

// convertVolume converts a Atempo volume to Docker Compose volume
func convertVolume(volume Volume) map[string]interface{} {
	dockerVolume := make(map[string]interface{})

	if volume.Driver != "" {
		dockerVolume["driver"] = volume.Driver
	}

	if len(volume.DriverOpts) > 0 {
		dockerVolume["driver_opts"] = volume.DriverOpts
	}

	if volume.External {
		external := make(map[string]interface{})
		external["external"] = true
		if volume.ExternalName != "" {
			external["name"] = volume.ExternalName
		}
		return external
	}

	return dockerVolume
}

// convertNetwork converts a Atempo network to Docker Compose network
func convertNetwork(network Network) map[string]interface{} {
	dockerNetwork := make(map[string]interface{})

	if network.Driver != "" {
		dockerNetwork["driver"] = network.Driver
	}

	if len(network.DriverOpts) > 0 {
		dockerNetwork["driver_opts"] = network.DriverOpts
	}

	if network.External {
		dockerNetwork["external"] = true
	}

	return dockerNetwork
}

// getServiceKey extracts a service identifier for container naming
func getServiceKey(service map[string]interface{}) string {
	if image, ok := service["image"].(string); ok {
		// Extract service name from image (e.g., "mysql:8.0" -> "mysql")
		parts := strings.Split(image, ":")
		if len(parts) > 0 {
			imageParts := strings.Split(parts[0], "/")
			return imageParts[len(imageParts)-1]
		}
	}
	return "service"
}

// writeDockerCompose writes the Docker Compose structure to a YAML file
func writeDockerCompose(compose *DockerCompose, filePath string) error {
	data, err := yaml.Marshal(compose)
	if err != nil {
		return fmt.Errorf("failed to marshal docker-compose: %w", err)
	}

	// Add header comment
	header := "# Generated by Atempo from atempo.json\n# Do not edit this file directly - modify atempo.json and run 'atempo reconfigure'\n\n"
	content := header + string(data)

	return os.WriteFile(filePath, []byte(content), 0644)
}

// AddService adds a new service to atempo.json
func AddService(projectPath, serviceName string, service Service) error {
	config, err := LoadAtempoConfig(projectPath)
	if err != nil {
		return err
	}

	if config.Services == nil {
		config.Services = make(map[string]Service)
	}

	config.Services[serviceName] = service

	return saveAtempoConfig(config, projectPath)
}

// RemoveService removes a service from atempo.json
func RemoveService(projectPath, serviceName string) error {
	config, err := LoadAtempoConfig(projectPath)
	if err != nil {
		return err
	}

	delete(config.Services, serviceName)

	return saveAtempoConfig(config, projectPath)
}

// saveAtempoConfig saves the atempo.json file
func saveAtempoConfig(config *AtempoConfig, projectPath string) error {
	atempoJsonPath := filepath.Join(projectPath, "atempo.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal atempo.json: %w", err)
	}

	return os.WriteFile(atempoJsonPath, data, 0644)
}

// AddPredefinedService adds a common service (minio, elasticsearch, etc.)
func AddPredefinedService(projectPath, serviceType string) error {
	service, exists := GetPredefinedService(serviceType)
	if !exists {
		return fmt.Errorf("unknown service type: %s", serviceType)
	}

	return AddService(projectPath, serviceType, service)
}

// GetPredefinedService returns predefined service configurations
func GetPredefinedService(serviceType string) (Service, bool) {
	services := map[string]Service{
		"minio": {
			Type:  "image",
			Image: "minio/minio",
			Ports: []string{"9000:9000", "9001:9001"},
			Command: []string{
				"server", "/data", "--console-address", ":9001",
			},
			Environment: map[string]string{
				"MINIO_ROOT_USER":     "minioadmin",
				"MINIO_ROOT_PASSWORD": "minioadmin",
			},
			Volumes: []string{"minio_data:/data"},
		},
		"elasticsearch": {
			Type:  "image",
			Image: "elasticsearch:8.8.0",
			Ports: []string{"9200:9200"},
			Environment: map[string]string{
				"discovery.type":         "single-node",
				"xpack.security.enabled": "false",
				"ES_JAVA_OPTS":           "-Xms512m -Xmx512m",
			},
			Volumes: []string{"elasticsearch_data:/usr/share/elasticsearch/data"},
		},
		"rabbitmq": {
			Type:  "image",
			Image: "rabbitmq:3-management",
			Ports: []string{"5672:5672", "15672:15672"},
			Environment: map[string]string{
				"RABBITMQ_DEFAULT_USER": "admin",
				"RABBITMQ_DEFAULT_PASS": "admin",
			},
			Volumes: []string{"rabbitmq_data:/var/lib/rabbitmq"},
		},
		"mongodb": {
			Type:  "image",
			Image: "mongo:6",
			Ports: []string{"27017:27017"},
			Environment: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "admin",
				"MONGO_INITDB_ROOT_PASSWORD": "admin",
			},
			Volumes: []string{"mongodb_data:/data/db"},
		},
	}

	service, exists := services[serviceType]
	return service, exists
}

// ListPredefinedServices returns available predefined services
func ListPredefinedServices() []string {
	return []string{"minio", "elasticsearch", "rabbitmq", "mongodb"}
}

// extractInternalPort extracts the internal port from a port mapping string
func extractInternalPort(portMapping string) int {
	parts := strings.Split(portMapping, ":")
	if len(parts) >= 2 {
		if port, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
			return port
		}
	}
	return 0
}

// convertServiceWithDynamicPorts converts a service with dynamically allocated ports
func convertServiceWithDynamicPorts(service Service, serviceName, projectName, framework string, portMapping map[int]int) map[string]interface{} {
	dockerService := make(map[string]interface{})

	// Handle build vs image
	if service.Type == "build" {
		// Generate project-specific image name
		imageName := fmt.Sprintf("%s-%s-%s", projectName, framework, serviceName)
		dockerService["image"] = imageName

		if service.Context != "" {
			dockerService["build"] = map[string]interface{}{
				"context":    service.Context,
				"dockerfile": service.Dockerfile,
			}
		} else {
			dockerService["build"] = map[string]interface{}{
				"context":    ".",
				"dockerfile": service.Dockerfile,
			}
		}
	} else if service.Image != "" {
		dockerService["image"] = service.Image
	}

	// Add container name with project prefix
	dockerService["container_name"] = fmt.Sprintf("%s-%s", projectName, serviceName)

	// Add restart policy
	if service.Restart != "" {
		dockerService["restart"] = service.Restart
	} else {
		dockerService["restart"] = "unless-stopped"
	}

	// Add optional fields
	if service.Command != nil {
		dockerService["command"] = service.Command
	}

	if service.WorkingDir != "" {
		dockerService["working_dir"] = service.WorkingDir
	}

	// Handle dynamic ports
	if len(service.Ports) > 0 && portMapping != nil {
		var dynamicPorts []string
		for _, originalPort := range service.Ports {
			internalPort := extractInternalPort(originalPort)
			if externalPort, exists := portMapping[internalPort]; exists {
				dynamicPorts = append(dynamicPorts, fmt.Sprintf("%d:%d", externalPort, internalPort))
			} else {
				// Keep original mapping if no dynamic allocation
				dynamicPorts = append(dynamicPorts, originalPort)
			}
		}
		dockerService["ports"] = dynamicPorts
	}

	if len(service.Volumes) > 0 {
		dockerService["volumes"] = service.Volumes
	}

	if len(service.Environment) > 0 {
		dockerService["environment"] = service.Environment
	}

	if len(service.DependsOn) > 0 {
		dockerService["depends_on"] = service.DependsOn
	}

	if len(service.Networks) > 0 {
		dockerService["networks"] = service.Networks
	}

	return dockerService
}

// generateAccessInfo creates a file with access information for the project
func generateAccessInfo(projectPath, projectName string, allocatedPorts map[string]map[int]int, services map[string]int) error {
	var info strings.Builder

	info.WriteString(fmt.Sprintf("# Access Information for %s\n\n", projectName))
	info.WriteString("## Service URLs\n\n")

	// Generate URLs for each service
	for serviceName, portMapping := range allocatedPorts {
		if len(portMapping) > 0 {
			info.WriteString(fmt.Sprintf("### %s\n", serviceName))

			// Port-based access
			for internalPort, externalPort := range portMapping {
				info.WriteString(fmt.Sprintf("- Port-based: http://localhost:%d (internal port %d)\n", externalPort, internalPort))
			}

			// Domain-based access (if DNS is configured)
			if serviceName == "webserver" || serviceName == "app" {
				info.WriteString(fmt.Sprintf("- Domain-based: http://%s.local\n", projectName))
			} else {
				info.WriteString(fmt.Sprintf("- Domain-based: http://%s.%s.local\n", serviceName, projectName))
			}

			info.WriteString("\n")
		}
	}

	// Write to file
	accessFile := filepath.Join(projectPath, ".atempo-access.md")
	return os.WriteFile(accessFile, []byte(info.String()), 0644)
}

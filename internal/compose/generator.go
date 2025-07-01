package compose

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// AtempoConfig represents the enhanced atempo.json structure
type AtempoConfig struct {
	Name      string                 `json:"name"`
	Framework string                 `json:"framework"`
	Language  string                 `json:"language"`
	Services  map[string]Service     `json:"services"`
	Volumes   map[string]Volume      `json:"volumes,omitempty"`
	Networks  map[string]Network     `json:"networks,omitempty"`
	Version   string                 `json:"version,omitempty"`
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
	config, err := LoadAtempoConfig(projectPath)
	if err != nil {
		return err
	}

	compose := &DockerCompose{
		Version:  "3.8",
		Services: make(map[string]interface{}),
		Volumes:  make(map[string]interface{}),
		Networks: make(map[string]interface{}),
	}

	// Convert services
	for serviceName, service := range config.Services {
		dockerService := convertService(service, config.Framework)
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

	// Add default network if none specified
	if len(compose.Networks) == 0 {
		compose.Networks[config.Framework] = map[string]interface{}{
			"driver": "bridge",
		}
		
		// Add network to all services
		for _, serviceInterface := range compose.Services {
			if serviceMap, ok := serviceInterface.(map[string]interface{}); ok {
				serviceMap["networks"] = []string{config.Framework}
			}
		}
	}

	// Write docker-compose.yml
	composePath := filepath.Join(projectPath, "docker-compose.yml")
	return writeDockerCompose(compose, composePath)
}

// convertService converts a Atempo service to Docker Compose service
func convertService(service Service, framework string) map[string]interface{} {
	dockerService := make(map[string]interface{})

	// Handle build vs image
	if service.Type == "build" {
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

	// Add container name
	if containerName, ok := dockerService["container_name"]; !ok || containerName == "" {
		dockerService["container_name"] = fmt.Sprintf("%s-%s", framework, getServiceKey(dockerService))
	}

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
				"ES_JAVA_OPTS":          "-Xms512m -Xmx512m",
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
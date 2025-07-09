package compose

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"gopkg.in/yaml.v3"
)

func TestLoadAtempoConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "compose-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create sample atempo.json
	sampleConfig := AtempoConfig{
		Name:      "test-project",
		Framework: "laravel",
		Language:  "php",
		Services: map[string]Service{
			"app": {
				Type:       "build",
				Dockerfile: "Dockerfile",
				Context:    ".",
				Ports:      []string{"80:80"},
				Environment: map[string]string{
					"APP_ENV": "local",
				},
			},
			"db": {
				Type:  "image",
				Image: "mysql:8.0",
				Ports: []string{"3306:3306"},
				Environment: map[string]string{
					"MYSQL_ROOT_PASSWORD": "password",
				},
			},
		},
		Volumes: map[string]Volume{
			"app-data": {
				Driver: "local",
			},
		},
		Networks: map[string]Network{
			"app-network": {
				Driver: "bridge",
			},
		},
	}

	// Write atempo.json
	configData, err := json.MarshalIndent(sampleConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	configPath := filepath.Join(tempDir, "atempo.json")
	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test loading config
	config, err := LoadAtempoConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config
	if config.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got '%s'", config.Name)
	}
	if config.Framework != "laravel" {
		t.Errorf("Expected framework 'laravel', got '%s'", config.Framework)
	}
	if config.Language != "php" {
		t.Errorf("Expected language 'php', got '%s'", config.Language)
	}

	// Verify services
	if len(config.Services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(config.Services))
	}

	appService, exists := config.Services["app"]
	if !exists {
		t.Error("Expected 'app' service to exist")
	} else {
		if appService.Type != "build" {
			t.Errorf("Expected app service type 'build', got '%s'", appService.Type)
		}
		if appService.Dockerfile != "Dockerfile" {
			t.Errorf("Expected app service dockerfile 'Dockerfile', got '%s'", appService.Dockerfile)
		}
	}

	dbService, exists := config.Services["db"]
	if !exists {
		t.Error("Expected 'db' service to exist")
	} else {
		if dbService.Type != "image" {
			t.Errorf("Expected db service type 'image', got '%s'", dbService.Type)
		}
		if dbService.Image != "mysql:8.0" {
			t.Errorf("Expected db service image 'mysql:8.0', got '%s'", dbService.Image)
		}
	}

	// Verify volumes
	if len(config.Volumes) != 1 {
		t.Errorf("Expected 1 volume, got %d", len(config.Volumes))
	}

	// Verify networks
	if len(config.Networks) != 1 {
		t.Errorf("Expected 1 network, got %d", len(config.Networks))
	}
}

func TestLoadAtempoConfig_FileNotFound(t *testing.T) {
	// Create temporary directory without atempo.json
	tempDir, err := os.MkdirTemp("", "compose-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Try to load non-existent config
	_, err = LoadAtempoConfig(tempDir)
	if err == nil {
		t.Error("Expected error when loading non-existent config")
	}
}

func TestLoadAtempoConfig_InvalidJSON(t *testing.T) {
	// Create temporary directory with invalid JSON
	tempDir, err := os.MkdirTemp("", "compose-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write invalid JSON
	configPath := filepath.Join(tempDir, "atempo.json")
	err = os.WriteFile(configPath, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	// Try to load invalid config
	_, err = LoadAtempoConfig(tempDir)
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

func TestConvertService(t *testing.T) {
	tests := []struct {
		name        string
		service     Service
		serviceName string
		projectName string
		framework   string
		expected    map[string]interface{}
	}{
		{
			name: "Image-based service",
			service: Service{
				Type:  "image",
				Image: "nginx:alpine",
				Ports: []string{"80:80"},
				Environment: map[string]string{
					"NGINX_HOST": "localhost",
				},
			},
			serviceName: "web",
			projectName: "test-project",
			framework:   "laravel",
			expected: map[string]interface{}{
				"image": "nginx:alpine",
				"ports": []string{"80:80"},
				"environment": map[string]string{
					"NGINX_HOST": "localhost",
				},
			},
		},
		{
			name: "Build-based service",
			service: Service{
				Type:       "build",
				Dockerfile: "Dockerfile",
				Context:    ".",
				Ports:      []string{"8000:8000"},
				Volumes:    []string{"./src:/app"},
				WorkingDir: "/app",
			},
			serviceName: "app",
			projectName: "test-project",
			framework:   "django",
			expected: map[string]interface{}{
				"build": map[string]interface{}{
					"context":    ".",
					"dockerfile": "Dockerfile",
				},
				"ports":       []string{"8000:8000"},
				"volumes":     []string{"./src:/app"},
				"working_dir": "/app",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertService(tt.service, tt.serviceName, tt.projectName, tt.framework)
			
			// Compare key fields
			if tt.service.Type == "image" {
				if result["image"] != tt.expected["image"] {
					t.Errorf("Expected image %v, got %v", tt.expected["image"], result["image"])
				}
			} else if tt.service.Type == "build" {
				build, exists := result["build"]
				if !exists {
					t.Error("Expected build configuration to exist")
				} else {
					buildMap, ok := build.(map[string]interface{})
					if !ok {
						t.Error("Expected build to be a map")
					} else {
						expectedBuild := tt.expected["build"].(map[string]interface{})
						if buildMap["context"] != expectedBuild["context"] {
							t.Errorf("Expected context %v, got %v", expectedBuild["context"], buildMap["context"])
						}
						if buildMap["dockerfile"] != expectedBuild["dockerfile"] {
							t.Errorf("Expected dockerfile %v, got %v", expectedBuild["dockerfile"], buildMap["dockerfile"])
						}
					}
				}
			}

			// Check ports
			if tt.service.Ports != nil {
				ports, exists := result["ports"]
				if !exists {
					t.Error("Expected ports to exist")
				} else {
					portSlice, ok := ports.([]string)
					if !ok {
						t.Error("Expected ports to be string slice")
					} else {
						expectedPorts := tt.expected["ports"].([]string)
						if len(portSlice) != len(expectedPorts) {
							t.Errorf("Expected %d ports, got %d", len(expectedPorts), len(portSlice))
						}
					}
				}
			}
		})
	}
}

func TestConvertVolume(t *testing.T) {
	tests := []struct {
		name     string
		volume   Volume
		expected map[string]interface{}
	}{
		{
			name: "Local volume",
			volume: Volume{
				Driver: "local",
			},
			expected: map[string]interface{}{
				"driver": "local",
			},
		},
		{
			name: "External volume",
			volume: Volume{
				External:     true,
				ExternalName: "shared-volume",
			},
			expected: map[string]interface{}{
				"external": true,
				"name":     "shared-volume",
			},
		},
		{
			name: "Volume with driver options",
			volume: Volume{
				Driver: "local",
				DriverOpts: map[string]string{
					"type": "tmpfs",
				},
			},
			expected: map[string]interface{}{
				"driver": "local",
				"driver_opts": map[string]string{
					"type": "tmpfs",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertVolume(tt.volume)
			
			// Compare driver
			if tt.volume.Driver != "" {
				if result["driver"] != tt.expected["driver"] {
					t.Errorf("Expected driver %v, got %v", tt.expected["driver"], result["driver"])
				}
			}

			// Compare external
			if tt.volume.External {
				if result["external"] != tt.expected["external"] {
					t.Errorf("Expected external %v, got %v", tt.expected["external"], result["external"])
				}
			}

			// Compare external name
			if tt.volume.ExternalName != "" {
				if result["name"] != tt.expected["name"] {
					t.Errorf("Expected name %v, got %v", tt.expected["name"], result["name"])
				}
			}

			// Compare driver options
			if tt.volume.DriverOpts != nil {
				opts, exists := result["driver_opts"]
				if !exists {
					t.Error("Expected driver_opts to exist")
				} else {
					optsMap, ok := opts.(map[string]string)
					if !ok {
						t.Error("Expected driver_opts to be string map")
					} else {
						expectedOpts := tt.expected["driver_opts"].(map[string]string)
						for key, value := range expectedOpts {
							if optsMap[key] != value {
								t.Errorf("Expected driver_opts[%s] = %v, got %v", key, value, optsMap[key])
							}
						}
					}
				}
			}
		})
	}
}

func TestConvertNetwork(t *testing.T) {
	tests := []struct {
		name     string
		network  Network
		expected map[string]interface{}
	}{
		{
			name: "Bridge network",
			network: Network{
				Driver: "bridge",
			},
			expected: map[string]interface{}{
				"driver": "bridge",
			},
		},
		{
			name: "External network",
			network: Network{
				External: true,
			},
			expected: map[string]interface{}{
				"external": true,
			},
		},
		{
			name: "Network with driver options",
			network: Network{
				Driver: "bridge",
				DriverOpts: map[string]string{
					"subnet": "172.20.0.0/16",
				},
			},
			expected: map[string]interface{}{
				"driver": "bridge",
				"driver_opts": map[string]string{
					"subnet": "172.20.0.0/16",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertNetwork(tt.network)
			
			// Compare driver
			if tt.network.Driver != "" {
				if result["driver"] != tt.expected["driver"] {
					t.Errorf("Expected driver %v, got %v", tt.expected["driver"], result["driver"])
				}
			}

			// Compare external
			if tt.network.External {
				if result["external"] != tt.expected["external"] {
					t.Errorf("Expected external %v, got %v", tt.expected["external"], result["external"])
				}
			}

			// Compare driver options
			if tt.network.DriverOpts != nil {
				opts, exists := result["driver_opts"]
				if !exists {
					t.Error("Expected driver_opts to exist")
				} else {
					optsMap, ok := opts.(map[string]string)
					if !ok {
						t.Error("Expected driver_opts to be string map")
					} else {
						expectedOpts := tt.expected["driver_opts"].(map[string]string)
						for key, value := range expectedOpts {
							if optsMap[key] != value {
								t.Errorf("Expected driver_opts[%s] = %v, got %v", key, value, optsMap[key])
							}
						}
					}
				}
			}
		})
	}
}

func TestExtractInternalPort(t *testing.T) {
	tests := []struct {
		name        string
		portMapping string
		expected    int
	}{
		{
			name:        "Standard port mapping",
			portMapping: "8080:80",
			expected:    80,
		},
		{
			name:        "Port mapping with host IP",
			portMapping: "127.0.0.1:8080:80",
			expected:    80,
		},
		{
			name:        "Port mapping with protocol",
			portMapping: "8080:80/tcp",
			expected:    80,
		},
		{
			name:        "Invalid port mapping",
			portMapping: "invalid",
			expected:    0,
		},
		{
			name:        "Empty port mapping",
			portMapping: "",
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractInternalPort(tt.portMapping)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestWriteDockerCompose(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "compose-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create sample compose structure
	compose := &DockerCompose{
		Version: "3.8",
		Services: map[string]interface{}{
			"app": map[string]interface{}{
				"image": "nginx:alpine",
				"ports": []string{"80:80"},
			},
		},
		Volumes: map[string]interface{}{
			"app-data": map[string]interface{}{
				"driver": "local",
			},
		},
		Networks: map[string]interface{}{
			"app-network": map[string]interface{}{
				"driver": "bridge",
			},
		},
	}

	// Write compose file
	composePath := filepath.Join(tempDir, "docker-compose.yml")
	err = writeDockerCompose(compose, composePath)
	if err != nil {
		t.Fatalf("Failed to write compose file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		t.Error("Compose file was not created")
	}

	// Read and verify file contents
	data, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("Failed to read compose file: %v", err)
	}

	var result DockerCompose
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal compose file: %v", err)
	}

	// Verify version
	if result.Version != "3.8" {
		t.Errorf("Expected version '3.8', got '%s'", result.Version)
	}

	// Verify services exist
	if len(result.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(result.Services))
	}

	// Verify volumes exist
	if len(result.Volumes) != 1 {
		t.Errorf("Expected 1 volume, got %d", len(result.Volumes))
	}

	// Verify networks exist
	if len(result.Networks) != 1 {
		t.Errorf("Expected 1 network, got %d", len(result.Networks))
	}
}

func TestService_StructureValidation(t *testing.T) {
	// Test Service structure validation
	service := Service{
		Type:        "build",
		Dockerfile:  "Dockerfile",
		Context:     ".",
		Command:     "npm start",
		WorkingDir:  "/app",
		Ports:       []string{"3000:3000"},
		Volumes:     []string{"./src:/app/src"},
		Environment: map[string]string{"NODE_ENV": "development"},
		DependsOn:   []string{"db"},
		Restart:     "unless-stopped",
		Networks:    []string{"app-network"},
	}

	// Verify all fields are accessible
	if service.Type != "build" {
		t.Errorf("Expected type 'build', got '%s'", service.Type)
	}
	if service.Dockerfile != "Dockerfile" {
		t.Errorf("Expected dockerfile 'Dockerfile', got '%s'", service.Dockerfile)
	}
	if service.Context != "." {
		t.Errorf("Expected context '.', got '%s'", service.Context)
	}
	if service.Command != "npm start" {
		t.Errorf("Expected command 'npm start', got '%v'", service.Command)
	}
	if service.WorkingDir != "/app" {
		t.Errorf("Expected working dir '/app', got '%s'", service.WorkingDir)
	}
	if len(service.Ports) != 1 {
		t.Errorf("Expected 1 port, got %d", len(service.Ports))
	}
	if len(service.Volumes) != 1 {
		t.Errorf("Expected 1 volume, got %d", len(service.Volumes))
	}
	if len(service.Environment) != 1 {
		t.Errorf("Expected 1 environment variable, got %d", len(service.Environment))
	}
	if len(service.DependsOn) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(service.DependsOn))
	}
	if service.Restart != "unless-stopped" {
		t.Errorf("Expected restart 'unless-stopped', got '%s'", service.Restart)
	}
	if len(service.Networks) != 1 {
		t.Errorf("Expected 1 network, got %d", len(service.Networks))
	}
}
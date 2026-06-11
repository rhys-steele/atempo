package compose

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the full structure of an atempo.json file.
type Config struct {
	Name      string             `json:"name"`
	Framework string             `json:"framework"`
	Language  string             `json:"language"`
	Version   string             `json:"version,omitempty"`
	Services  map[string]Service `json:"services"`
	Volumes   map[string]Volume  `json:"volumes,omitempty"`
	Networks  map[string]Network `json:"networks,omitempty"`
}

// Service is a single Docker service entry in atempo.json.
type Service struct {
	Type        string            `json:"type"`                  // "build" or "image"
	Image       string            `json:"image,omitempty"`
	Dockerfile  string            `json:"dockerfile,omitempty"`
	Context     string            `json:"context,omitempty"`
	Command     interface{}       `json:"command,omitempty"`     // string or []string
	WorkingDir  string            `json:"working_dir,omitempty"`
	Ports       []string          `json:"ports,omitempty"`
	Volumes     []string          `json:"volumes,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	DependsOn   []string          `json:"depends_on,omitempty"`
	Restart     string            `json:"restart,omitempty"`
	Networks    []string          `json:"networks,omitempty"`
}

// Volume is a named Docker volume definition.
type Volume struct {
	Driver     string            `json:"driver,omitempty"`
	DriverOpts map[string]string `json:"driver_opts,omitempty"`
	External   bool              `json:"external,omitempty"`
}

// Network is a named Docker network definition.
type Network struct {
	Driver     string            `json:"driver,omitempty"`
	DriverOpts map[string]string `json:"driver_opts,omitempty"`
	External   bool              `json:"external,omitempty"`
}

// PortAssignment records the resolved host port for one container port.
type PortAssignment struct {
	Service  string
	Internal int
	External int
}

// LoadConfig reads and parses atempo.json from projectPath.
func LoadConfig(projectPath string) (*Config, error) {
	data, err := os.ReadFile(filepath.Join(projectPath, "atempo.json"))
	if err != nil {
		return nil, fmt.Errorf("read atempo.json: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse atempo.json: %w", err)
	}
	return &cfg, nil
}

// SaveConfig writes cfg back to atempo.json in projectPath.
func SaveConfig(cfg *Config, projectPath string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal atempo.json: %w", err)
	}
	return os.WriteFile(filepath.Join(projectPath, "atempo.json"), data, 0644)
}

// AssignPorts probes for free ports for every mapped port in cfg, updates cfg
// in place, and returns the full list of assignments.
//
// The preferred host port is read from the left side of each "host:container"
// mapping in atempo.json.  If that port is busy the next free one is used.
// The updated port strings are written back into cfg so SaveConfig persists them.
func AssignPorts(cfg *Config) ([]PortAssignment, error) {
	var assignments []PortAssignment

	for svcName, svc := range cfg.Services {
		var newPorts []string
		for _, mapping := range svc.Ports {
			preferred, internal, err := parsePortMapping(mapping)
			if err != nil {
				return nil, fmt.Errorf("service %s: %w", svcName, err)
			}

			external, err := freePort(preferred)
			if err != nil {
				return nil, fmt.Errorf("service %s port %d: %w", svcName, internal, err)
			}

			newPorts = append(newPorts, fmt.Sprintf("%d:%d", external, internal))
			assignments = append(assignments, PortAssignment{
				Service:  svcName,
				Internal: internal,
				External: external,
			})
		}
		svc.Ports = newPorts
		cfg.Services[svcName] = svc
	}

	return assignments, nil
}

// GenerateDockerCompose reads atempo.json from projectPath, assigns free ports,
// writes the updated atempo.json, then writes docker-compose.yml.
// It returns the port assignments so callers can display them to the user.
func GenerateDockerCompose(projectPath string) ([]PortAssignment, error) {
	cfg, err := LoadConfig(projectPath)
	if err != nil {
		return nil, err
	}

	projectName := cfg.Name
	if projectName == "" || strings.Contains(projectName, "{{") {
		projectName = filepath.Base(projectPath)
	}

	assignments, err := AssignPorts(cfg)
	if err != nil {
		return nil, err
	}

	// Persist the resolved ports back into atempo.json.
	if err := SaveConfig(cfg, projectPath); err != nil {
		return nil, fmt.Errorf("save atempo.json with port assignments: %w", err)
	}

	compose := buildCompose(cfg, projectName)

	if err := writeCompose(compose, filepath.Join(projectPath, "docker-compose.yml")); err != nil {
		return nil, err
	}

	return assignments, nil
}

// buildCompose converts a Config into the docker-compose.yml data structure.
func buildCompose(cfg *Config, projectName string) map[string]interface{} {
	services := make(map[string]interface{})
	for name, svc := range cfg.Services {
		services[name] = convertService(svc, name, projectName, cfg.Framework)
	}

	volumes := make(map[string]interface{})
	for name, v := range cfg.Volumes {
		volumes[name] = convertVolume(v)
	}

	networks := make(map[string]interface{})
	for name, n := range cfg.Networks {
		networks[name] = convertNetwork(n)
	}

	// Add a project-scoped bridge network when none are declared.
	if len(networks) == 0 {
		netName := projectName + "-network"
		networks[netName] = map[string]interface{}{"driver": "bridge"}
		for name, svcIface := range services {
			if m, ok := svcIface.(map[string]interface{}); ok {
				m["networks"] = []string{netName}
				services[name] = m
			}
		}
	}

	compose := map[string]interface{}{
		"services": services,
	}
	if len(volumes) > 0 {
		compose["volumes"] = volumes
	}
	if len(networks) > 0 {
		compose["networks"] = networks
	}
	return compose
}

func convertService(svc Service, name, projectName, framework string) map[string]interface{} {
	m := map[string]interface{}{
		"container_name": fmt.Sprintf("%s-%s", projectName, name),
		"restart":        orDefault(svc.Restart, "unless-stopped"),
	}

	if svc.Type == "build" {
		ctx := orDefault(svc.Context, ".")
		m["build"] = map[string]interface{}{
			"context":    ctx,
			"dockerfile": svc.Dockerfile,
		}
		m["image"] = fmt.Sprintf("%s-%s-%s", projectName, framework, name)
	} else if svc.Image != "" {
		m["image"] = svc.Image
	}

	if svc.Command != nil {
		m["command"] = svc.Command
	}
	if svc.WorkingDir != "" {
		m["working_dir"] = svc.WorkingDir
	}
	if len(svc.Ports) > 0 {
		m["ports"] = svc.Ports
	}
	if len(svc.Volumes) > 0 {
		m["volumes"] = svc.Volumes
	}
	if len(svc.Environment) > 0 {
		m["environment"] = svc.Environment
	}
	if len(svc.DependsOn) > 0 {
		m["depends_on"] = svc.DependsOn
	}
	if len(svc.Networks) > 0 {
		m["networks"] = svc.Networks
	}
	return m
}

func convertVolume(v Volume) map[string]interface{} {
	if v.External {
		return map[string]interface{}{"external": true}
	}
	m := map[string]interface{}{}
	if v.Driver != "" {
		m["driver"] = v.Driver
	}
	if len(v.DriverOpts) > 0 {
		m["driver_opts"] = v.DriverOpts
	}
	return m
}

func convertNetwork(n Network) map[string]interface{} {
	m := map[string]interface{}{}
	if n.External {
		m["external"] = true
		return m
	}
	if n.Driver != "" {
		m["driver"] = n.Driver
	}
	if len(n.DriverOpts) > 0 {
		m["driver_opts"] = n.DriverOpts
	}
	return m
}

func writeCompose(compose map[string]interface{}, path string) error {
	data, err := yaml.Marshal(compose)
	if err != nil {
		return fmt.Errorf("marshal docker-compose.yml: %w", err)
	}
	header := "# Generated by Atempo — edit atempo.json and run 'atempo reconfigure' to update\n\n"
	return os.WriteFile(path, []byte(header+string(data)), 0644)
}

// ─── service management ───────────────────────────────────────────────────────

// AddService appends a service to atempo.json and regenerates docker-compose.yml.
func AddService(projectPath, name string, svc Service) error {
	cfg, err := LoadConfig(projectPath)
	if err != nil {
		return err
	}
	if cfg.Services == nil {
		cfg.Services = make(map[string]Service)
	}
	cfg.Services[name] = svc
	if err := SaveConfig(cfg, projectPath); err != nil {
		return err
	}
	_, err = GenerateDockerCompose(projectPath)
	return err
}

// AddPredefinedService adds a known service type by name.
func AddPredefinedService(projectPath, serviceType string) error {
	svc, ok := predefinedService(serviceType)
	if !ok {
		return fmt.Errorf("unknown service %q — available: %s",
			serviceType, strings.Join(ListPredefinedServices(), ", "))
	}
	return AddService(projectPath, serviceType, svc)
}

// ListPredefinedServices returns the names of available predefined services.
func ListPredefinedServices() []string {
	return []string{"minio", "elasticsearch", "rabbitmq", "mongodb"}
}

func predefinedService(t string) (Service, bool) {
	m := map[string]Service{
		"minio": {
			Type:  "image",
			Image: "minio/minio",
			Ports: []string{"9000:9000", "9001:9001"},
			Command: []string{"server", "/data", "--console-address", ":9001"},
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
			Type:    "image",
			Image:   "rabbitmq:3-management",
			Ports:   []string{"5672:5672", "15672:15672"},
			Volumes: []string{"rabbitmq_data:/var/lib/rabbitmq"},
			Environment: map[string]string{
				"RABBITMQ_DEFAULT_USER": "admin",
				"RABBITMQ_DEFAULT_PASS": "admin",
			},
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
	svc, ok := m[t]
	return svc, ok
}

// ─── helpers ──────────────────────────────────────────────────────────────────

// parsePortMapping splits "host:container" or "container" into (preferred, internal).
func parsePortMapping(mapping string) (preferred, internal int, err error) {
	parts := strings.Split(mapping, ":")
	switch len(parts) {
	case 1:
		p, e := strconv.Atoi(parts[0])
		if e != nil {
			return 0, 0, fmt.Errorf("invalid port %q", mapping)
		}
		return p, p, nil
	case 2:
		host, e1 := strconv.Atoi(parts[0])
		cont, e2 := strconv.Atoi(parts[1])
		if e1 != nil || e2 != nil {
			return 0, 0, fmt.Errorf("invalid port mapping %q", mapping)
		}
		return host, cont, nil
	default:
		return 0, 0, fmt.Errorf("invalid port mapping %q", mapping)
	}
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

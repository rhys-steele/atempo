package registry

import "time"

// Project represents a registered Atempo project.
type Project struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Framework    string    `json:"framework"`
	Version      string    `json:"version"`
	Status       string    `json:"status"` // running | stopped | healthy | unhealthy
	CreatedAt    time.Time `json:"created_at"`
	LastAccessed time.Time `json:"last_accessed"`

	Ports    []Port    `json:"ports"`
	URLs     []string  `json:"urls"`
	Services []Service `json:"services"`

	GitBranch string `json:"git_branch,omitempty"`
	GitStatus string `json:"git_status,omitempty"`
}

// Port represents a host:container port mapping for a service.
type Port struct {
	Service  string `json:"service"`
	External int    `json:"external"`
	Internal int    `json:"internal"`
	Protocol string `json:"protocol"`
}

// Service tracks a Docker service and its runtime status.
type Service struct {
	Name   string `json:"name"`
	Status string `json:"status"` // running | stopped | healthy | unhealthy
	URL    string `json:"url,omitempty"`
}

// Registry is the top-level document stored in ~/.atempo/registry.json.
type Registry struct {
	Version  string    `json:"version"`
	Projects []Project `json:"projects"`
}

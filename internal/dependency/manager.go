package dependency

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Status represents the availability of a dependency.
type Status int

const (
	StatusAvailable Status = iota
	StatusMissing
	StatusNotRunning   // binary found but daemon not running
	StatusUnauthorised // daemon running but permission denied
)

// Dependency describes a system tool Atempo requires or recommends.
type Dependency struct {
	Name     string
	Required bool
	install  map[string]string // runtime.GOOS → install hint
}

// Check is the result of probing a single dependency.
type Check struct {
	Dep     Dependency
	Status  Status
	Version string
	Message string // human-readable problem description
	Hint    string // how to fix it
}

// OK reports whether the dependency is available.
func (c *Check) OK() bool { return c.Status == StatusAvailable }

// Manager checks system dependencies required to run Atempo.
type Manager struct{}

func NewManager() *Manager { return &Manager{} }

// CheckAll probes every dependency and returns all results.
func (m *Manager) CheckAll() []Check {
	return []Check{
		m.checkDocker(),
		m.checkDockerCompose(),
		m.checkGit(),
	}
}

// Validate returns an error if any required dependency is unavailable.
func (m *Manager) Validate() error {
	var missing []string
	for _, c := range m.CheckAll() {
		if c.Dep.Required && !c.OK() {
			missing = append(missing, c.Dep.Name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("required dependencies missing: %s — run 'atempo doctor' for details",
			strings.Join(missing, ", "))
	}
	return nil
}

// checkDocker checks that the Docker CLI is present and the daemon is reachable.
func (m *Manager) checkDocker() Check {
	dep := Dependency{
		Name:     "Docker",
		Required: true,
		install: map[string]string{
			"darwin":  "Install Docker Desktop: https://docs.docker.com/desktop/mac/",
			"windows": "Install Docker Desktop: https://docs.docker.com/desktop/windows/",
			"linux":   "Install Docker Engine: https://docs.docker.com/engine/install/",
		},
	}

	if _, err := exec.LookPath("docker"); err != nil {
		return Check{
			Dep:     dep,
			Status:  StatusMissing,
			Message: "docker binary not found",
			Hint:    installHint(dep),
		}
	}

	// Verify the daemon is reachable.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "info", "--format", "{{.ServerVersion}}")
	out, err := cmd.Output()
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "permission denied") {
			return Check{
				Dep:     dep,
				Status:  StatusUnauthorised,
				Message: "Docker daemon not accessible (permission denied)",
				Hint:    dockerPermHint(),
			}
		}
		return Check{
			Dep:     dep,
			Status:  StatusNotRunning,
			Message: "Docker daemon is not running",
			Hint:    dockerStartHint(),
		}
	}

	return Check{
		Dep:     dep,
		Status:  StatusAvailable,
		Version: "Docker " + strings.TrimSpace(string(out)),
	}
}

// checkDockerCompose checks that the Docker Compose v2 plugin is available.
func (m *Manager) checkDockerCompose() Check {
	dep := Dependency{
		Name:     "Docker Compose v2",
		Required: true,
		install: map[string]string{
			"darwin":  "Update Docker Desktop — Compose v2 is bundled",
			"windows": "Update Docker Desktop — Compose v2 is bundled",
			"linux":   "Install the compose plugin: https://docs.docker.com/compose/install/linux/",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// v2 is `docker compose version`, not `docker-compose --version`
	cmd := exec.CommandContext(ctx, "docker", "compose", "version", "--short")
	out, err := cmd.Output()
	if err != nil {
		return Check{
			Dep:     dep,
			Status:  StatusMissing,
			Message: "Docker Compose v2 plugin not found (need 'docker compose', not 'docker-compose')",
			Hint:    installHint(dep),
		}
	}

	return Check{
		Dep:     dep,
		Status:  StatusAvailable,
		Version: "Compose v" + strings.TrimSpace(string(out)),
	}
}

// checkGit checks that git is installed (recommended, not required).
func (m *Manager) checkGit() Check {
	dep := Dependency{
		Name:     "Git",
		Required: false,
		install: map[string]string{
			"darwin":  "xcode-select --install",
			"windows": "winget install Git.Git",
			"linux":   "apt-get install git  (or equivalent)",
		},
	}

	if _, err := exec.LookPath("git"); err != nil {
		return Check{
			Dep:     dep,
			Status:  StatusMissing,
			Message: "git not found — project initialisation will be skipped",
			Hint:    installHint(dep),
		}
	}

	cmd := exec.Command("git", "--version")
	out, _ := cmd.Output()

	return Check{
		Dep:     dep,
		Status:  StatusAvailable,
		Version: strings.TrimPrefix(strings.TrimSpace(string(out)), "git version "),
	}
}

func installHint(dep Dependency) string {
	if hint, ok := dep.install[runtime.GOOS]; ok {
		return hint
	}
	return "see https://docs.docker.com for installation instructions"
}

func dockerPermHint() string {
	switch runtime.GOOS {
	case "linux":
		return "add your user to the docker group: sudo usermod -aG docker $USER && newgrp docker"
	default:
		return "ensure Docker Desktop is running and has the necessary permissions"
	}
}

func dockerStartHint() string {
	switch runtime.GOOS {
	case "darwin":
		return "start Docker Desktop"
	case "windows":
		return "start Docker Desktop from the Start menu"
	case "linux":
		return "sudo systemctl start docker"
	default:
		return "start the Docker daemon"
	}
}

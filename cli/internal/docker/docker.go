package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ServiceInfo is a row from `docker compose ps --format json`.
type ServiceInfo struct {
	Name    string
	Service string
	Status  string // e.g. "running", "exited"
	Health  string // e.g. "healthy", "unhealthy", ""
	Ports   string
}

// Up runs `docker compose up -d` in projectPath.
func Up(projectPath string) error {
	return run(projectPath, 3*time.Minute, "up", "--build", "-d")
}

// Down runs `docker compose down` in projectPath.
func Down(projectPath string) error {
	return run(projectPath, 2*time.Minute, "down")
}

// Logs streams `docker compose logs -f` in projectPath, inheriting stdio.
// It blocks until the user interrupts (Ctrl-C) or the process exits.
func Logs(projectPath string, services []string) error {
	args := append([]string{"compose", "logs", "-f"}, services...)
	cmd := exec.Command("docker", args...)
	cmd.Dir = composeDir(projectPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// LogsWriter streams `docker compose logs -f` writing each line to w.
func LogsWriter(projectPath string, services []string, w io.Writer) error {
	args := append([]string{"compose", "logs", "-f"}, services...)
	cmd := exec.Command("docker", args...)
	cmd.Dir = composeDir(projectPath)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

// PS returns the running service list for the compose project in projectPath.
// Uses `docker compose ps --format json` for structured, panic-safe parsing.
func PS(projectPath string) ([]ServiceInfo, error) {
	dir := composeDir(projectPath)
	if dir == "" {
		return nil, fmt.Errorf("no docker-compose.yml found in %s", projectPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", "ps", "--format", "json")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker compose ps: %w", err)
	}

	return parsePS(out)
}

// ─── internal ─────────────────────────────────────────────────────────────────

func run(projectPath string, timeout time.Duration, args ...string) error {
	dir := composeDir(projectPath)
	if dir == "" {
		return fmt.Errorf("no docker-compose.yml found in %s — run 'atempo reconfigure' to generate it", projectPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmdArgs := append([]string{"compose"}, args...)
	cmd := exec.CommandContext(ctx, "docker", cmdArgs...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timed out after %v", timeout)
		}
		return err
	}
	return nil
}

// composeDir returns the directory containing docker-compose.yml, checking both
// the project root and infra/docker/ subdirectory. Returns "" if not found.
func composeDir(projectPath string) string {
	candidates := []string{
		filepath.Join(projectPath, "docker-compose.yml"),
		filepath.Join(projectPath, "infra", "docker", "docker-compose.yml"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return filepath.Dir(p)
		}
	}
	return ""
}

// psRow matches the JSON output of `docker compose ps --format json`.
type psRow struct {
	Name    string `json:"Name"`
	Service string `json:"Service"`
	State   string `json:"State"`
	Health  string `json:"Health"`
}

// parsePS handles both a JSON array and newline-delimited JSON objects,
// since Docker Compose changed its output format across versions.
func parsePS(data []byte) ([]ServiceInfo, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return nil, nil
	}

	var rows []psRow

	if strings.HasPrefix(trimmed, "[") {
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("parse docker compose ps output: %w", err)
		}
	} else {
		// Newline-delimited JSON (older Compose versions).
		for _, line := range strings.Split(trimmed, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var r psRow
			if err := json.Unmarshal([]byte(line), &r); err != nil {
				continue
			}
			rows = append(rows, r)
		}
	}

	out := make([]ServiceInfo, len(rows))
	for i, r := range rows {
		out[i] = ServiceInfo{
			Name:    r.Name,
			Service: r.Service,
			Status:  strings.ToLower(r.State),
			Health:  strings.ToLower(r.Health),
		}
	}
	return out, nil
}

package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Installer defines how a framework should be installed.
// This includes the command to run and the working directory context.
type Installer struct {
	Type    string   `json:"type"`     // e.g., "composer", "docker", "shell"
	Command []string `json:"command"`  // Full command with args (supports templating)
	WorkDir string   `json:"work-dir"` // Directory to run the command in
}

// Metadata describes a Steele template's configuration,
// including language, installer, and framework compatibility.
type Metadata struct {
	Framework   string    `json:"framework"`    // e.g., "laravel"
	Language    string    `json:"language"`     // e.g., "php"
	Installer   Installer `json:"installer"`    // How to scaffold the source code
	WorkingDir  string    `json:"working-dir"`  // Expected project root path in container, e.g., /var/www
	MinVersion  string    `json:"min-version"`  // Minimum supported version (semantic)
}

// Run executes the scaffolding process for the given framework and version.
// It loads the template's `steele.json`, performs template substitution,
// and runs the specified install command (e.g., composer create-project).
func Run(framework string, version string) error {
	templatePath := filepath.Join("templates", framework)
	metaPath := filepath.Join(templatePath, "steele.json")

	// Load steele.json from the template directory
	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		return fmt.Errorf("could not read steele.json: %w", err)
	}

	// Parse the metadata JSON into a structured object
	var meta Metadata
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return fmt.Errorf("invalid steele.json: %w", err)
	}

	// Get the current working directory (user's target project root)
	projectDir, _ := os.Getwd()
	projectName := filepath.Base(projectDir)

	// Perform template variable substitution in the command
	command := make([]string, len(meta.Installer.Command))
	for i, part := range meta.Installer.Command {
		part = strings.ReplaceAll(part, "{{name}}", "src")
		part = strings.ReplaceAll(part, "{{cwd}}", projectDir)
		part = strings.ReplaceAll(part, "{{project}}", projectName)
		command[i] = part
	}

	// Prepare the executable command
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("→ Running:", strings.Join(command, " "))
	return cmd.Run()
}

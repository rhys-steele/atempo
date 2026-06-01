package scaffold

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"atempo/internal/compose"
	"atempo/internal/interfaces"
	"atempo/internal/registry"
)

// Installer describes how to run the framework installer command.
type Installer struct {
	Type    string   `json:"type"`    // "docker" or "template" (template = no installer, just copy files)
	Command []string `json:"command"` // full command with {{projectPath}} and {{version}} placeholders
}

// Metadata is the parsed atempo.json for a framework template.
type Metadata struct {
	Name      string    `json:"name"`
	Framework string    `json:"framework"`
	Language  string    `json:"language"`
	Version   string    `json:"version"`
	Installer Installer `json:"installer"`
}

// Options controls what Run does.
type Options struct {
	Framework   string
	Version     string
	ProjectPath string // absolute path to the project directory (must already exist)
	ProjectName string // display name / registry name
	TemplatesFS fs.FS  // embedded templates FS rooted at the binary
	Registry    interfaces.RegistryRepository
}

// Run scaffolds a new project. The directory at opts.ProjectPath must already exist.
// It returns the port assignments so callers can display them.
func Run(opts Options) ([]compose.PortAssignment, error) {
	meta, err := loadMetadata(opts.TemplatesFS, opts.Framework)
	if err != nil {
		return nil, err
	}

	if err := runInstaller(meta, opts); err != nil {
		return nil, fmt.Errorf("installer: %w", err)
	}

	if err := copyTemplateFiles(opts.TemplatesFS, opts.Framework, opts.ProjectPath, opts.ProjectName, opts.Version); err != nil {
		return nil, fmt.Errorf("copy templates: %w", err)
	}

	if err := postInstall(meta.Framework, opts.ProjectPath); err != nil {
		return nil, fmt.Errorf("post-install: %w", err)
	}

	assignments, err := registerAndGenerate(opts)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return assignments, nil
}

// loadMetadata reads and parses the framework's atempo.json from the embedded FS.
func loadMetadata(templatesFS fs.FS, framework string) (Metadata, error) {
	path := fmt.Sprintf("frameworks/%s/atempo.json", framework)
	data, err := fs.ReadFile(templatesFS, path)
	if err != nil {
		return Metadata{}, fmt.Errorf("no template for framework %q — run 'atempo doctor' to verify your installation", framework)
	}
	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return Metadata{}, fmt.Errorf("parse template config for %s: %w", framework, err)
	}
	return meta, nil
}

// runInstaller executes the framework installer if the installer type is "docker".
// For "template" type there is nothing to run — the source files are included in the template.
func runInstaller(meta Metadata, opts Options) error {
	if meta.Installer.Type != "docker" {
		return nil
	}
	if len(meta.Installer.Command) == 0 {
		return nil
	}

	cmd := applyPlaceholders(meta.Installer.Command, opts)
	cmd = applyVersionConstraints(cmd, meta.Framework, opts.Version)

	// Each element is a separate argument — no shell interpolation.
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Dir = opts.ProjectPath
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("command %q failed: %w", cmd[0], err)
	}
	return nil
}

// applyPlaceholders replaces template variables in each command token.
func applyPlaceholders(command []string, opts Options) []string {
	result := make([]string, len(command))
	majorVersion := majorVersion(opts.Version)
	for i, part := range command {
		part = strings.ReplaceAll(part, "{{projectPath}}", opts.ProjectPath)
		part = strings.ReplaceAll(part, "{{project}}", opts.ProjectName)
		part = strings.ReplaceAll(part, "{{version}}", opts.Version)
		part = strings.ReplaceAll(part, "{{majorVersion}}", majorVersion)
		result[i] = part
	}
	return result
}

// copyTemplateFiles walks the framework's template directory in the embedded FS
// and copies every file to projectPath, expanding {{project}} and {{version}} tokens in content.
func copyTemplateFiles(templatesFS fs.FS, framework, projectPath, projectName, version string) error {
	root := fmt.Sprintf("frameworks/%s", framework)
	return fs.WalkDir(templatesFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}
		// Skip the top-level atempo.json — that gets written by registerAndGenerate.
		if rel == "atempo.json" {
			return nil
		}

		dst := filepath.Join(projectPath, filepath.FromSlash(rel))

		if d.IsDir() {
			return os.MkdirAll(dst, 0755)
		}

		data, err := fs.ReadFile(templatesFS, path)
		if err != nil {
			return fmt.Errorf("read template %s: %w", path, err)
		}

		content := expandTokens(string(data), projectName, version)

		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
		}
		return os.WriteFile(dst, []byte(content), 0644)
	})
}

// applyVersionConstraints adds framework-specific version constraints to installer commands.
// For Laravel: appends ":^{major}.0" to the package name if version is specified.
func applyVersionConstraints(command []string, framework, version string) []string {
	if version == "" {
		return command
	}
	switch framework {
	case "laravel":
		major := majorVersion(version)
		for i, arg := range command {
			if arg == "laravel/laravel" {
				command[i] = fmt.Sprintf("laravel/laravel:^%s.0", major)
				break
			}
		}
	}
	return command
}

// expandTokens replaces the simple template variables used in template files.
func expandTokens(content, projectName, version string) string {
	content = strings.ReplaceAll(content, "{{project}}", projectName)
	content = strings.ReplaceAll(content, "{{version}}", version)
	content = strings.ReplaceAll(content, "{{majorVersion}}", majorVersion(version))
	return content
}

// postInstall runs framework-specific setup after the installer and template copy.
func postInstall(framework, projectPath string) error {
	switch framework {
	case "laravel":
		return postInstallLaravel(projectPath)
	}
	return nil
}

// postInstallLaravel copies .env.example → .env and patches DB/Redis config for Docker.
func postInstallLaravel(projectPath string) error {
	srcDir := filepath.Join(projectPath, "src")
	envExample := filepath.Join(srcDir, ".env.example")
	envFile := filepath.Join(srcDir, ".env")

	if _, err := os.Stat(envExample); err != nil {
		return nil // .env.example not present yet — installer may not have run
	}
	if _, err := os.Stat(envFile); err == nil {
		return nil // .env already exists
	}

	data, err := os.ReadFile(envExample)
	if err != nil {
		return fmt.Errorf("read .env.example: %w", err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "DB_HOST=127.0.0.1", "DB_HOST=mysql")
	env = strings.ReplaceAll(env, "DB_USERNAME=root", "DB_USERNAME=laravel")
	env = strings.ReplaceAll(env, "DB_PASSWORD=", "DB_PASSWORD=laravel")
	if !strings.Contains(env, "REDIS_HOST=") {
		env += "\nREDIS_HOST=redis\n"
	}

	return os.WriteFile(envFile, []byte(env), 0644)
}

// registerAndGenerate writes atempo.json (with assigned ports), registers in registry,
// and generates docker-compose.yml. Returns the port assignments.
func registerAndGenerate(opts Options) ([]compose.PortAssignment, error) {
	// Load the template's atempo.json, expand tokens, and write to projectPath.
	templatePath := fmt.Sprintf("frameworks/%s/atempo.json", opts.Framework)
	data, err := fs.ReadFile(opts.TemplatesFS, templatePath)
	if err != nil {
		return nil, fmt.Errorf("read atempo.json template: %w", err)
	}
	expanded := expandTokens(string(data), opts.ProjectName, opts.Version)
	atempoJSONPath := filepath.Join(opts.ProjectPath, "atempo.json")
	if err := os.WriteFile(atempoJSONPath, []byte(expanded), 0644); err != nil {
		return nil, fmt.Errorf("write atempo.json: %w", err)
	}

	// Assign free ports, write updated atempo.json, and generate docker-compose.yml.
	assignments, err := compose.GenerateDockerCompose(opts.ProjectPath)
	if err != nil {
		return nil, fmt.Errorf("generate docker-compose.yml: %w", err)
	}

	// Register in the project registry.
	if opts.Registry != nil {
		proj := &registry.Project{
			Name:      opts.ProjectName,
			Path:      opts.ProjectPath,
			Framework: opts.Framework,
			Version:   opts.Version,
			Status:    "stopped",
			CreatedAt: time.Now(),
		}
		if err := opts.Registry.AddProject(proj); err != nil {
			return nil, fmt.Errorf("register project: %w", err)
		}
	}

	return assignments, nil
}

// majorVersion returns the major version string from a semver string.
// "11.3" → "11", "3" → "3", "" → "latest".
func majorVersion(version string) string {
	if version == "" {
		return "latest"
	}
	parts := strings.SplitN(version, ".", 2)
	return parts[0]
}

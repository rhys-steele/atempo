package scaffold

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"steele/internal/compose"
	"steele/internal/registry"
	"steele/internal/utils"
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
	Name        string    `json:"name"`         // Project name template
	Framework   string    `json:"framework"`    // e.g., "laravel"
	Language    string    `json:"language"`     // e.g., "php"
	Installer   Installer `json:"installer"`    // How to scaffold the source code
	WorkingDir  string    `json:"working-dir"`  // Expected project root path in container, e.g., /var/www
	MinVersion  string    `json:"min-version"`  // Minimum supported version (semantic)
}

// Run executes the scaffolding process for the given framework and version.
// It loads the template's `steele.json`, performs template substitution,
// runs the specified install command, and copies template files.
func Run(framework string, version string, templatesFS, mcpServersFS embed.FS) error {
	// Load steele.json (try embedded first, fallback to filesystem)
	var metaBytes []byte
	var err error
	
	// Try embedded first
	embeddedPath := fmt.Sprintf("templates/%s/steele.json", framework)
	metaBytes, err = templatesFS.ReadFile(embeddedPath)
	if err != nil {
		// Fallback to filesystem - find templates relative to binary location
		filesystemPath, pathErr := getFilesystemTemplatePath(framework, "steele.json")
		if pathErr != nil {
			return fmt.Errorf("could not locate steele.json for %s: %w", framework, pathErr)
		}
		metaBytes, err = os.ReadFile(filesystemPath)
		if err != nil {
			return fmt.Errorf("could not read steele.json for %s: %w", framework, err)
		}
	}

	// Parse the metadata JSON into a structured object
	var meta Metadata
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return fmt.Errorf("invalid steele.json: %w", err)
	}

	// Validate version compatibility
	if err := validateVersion(version, meta); err != nil {
		return fmt.Errorf("version validation failed: %w", err)
	}

	// Get the current working directory (user's target project root)
	projectDir, _ := os.Getwd()
	projectName := filepath.Base(projectDir)

	// Step 1: Run the framework installer (e.g., composer create-project)
	fmt.Println("🚀 Installing", framework, version, "application...")
	if err := runInstaller(meta, projectDir, projectName, version); err != nil {
		return fmt.Errorf("installer failed: %w", err)
	}

	// Step 2: Copy template files (AI context, Docker setup, etc.)
	fmt.Println("📁 Copying template files...")
	if err := copyTemplateFiles(projectDir, meta.Framework, templatesFS, mcpServersFS); err != nil {
		return fmt.Errorf("failed to copy template files: %w", err)
	}

	// Step 3: Run post-installation setup
	fmt.Println("⚙️  Running post-installation setup...")
	if err := runPostInstall(meta, projectDir); err != nil {
		return fmt.Errorf("post-installation failed: %w", err)
	}

	// Step 4: Register project and generate docker-compose
	fmt.Println("📝 Registering project and generating docker-compose...")
	if err := finalizeProject(meta, projectDir, projectName, version); err != nil {
		fmt.Printf("⚠️  Warning: %v\n", err)
	}

	return nil
}

// runInstaller executes the framework installation command
func runInstaller(meta Metadata, projectDir, projectName, version string) error {
	// Perform template variable substitution in the command
	command := make([]string, len(meta.Installer.Command))
	for i, part := range meta.Installer.Command {
		part = strings.ReplaceAll(part, "{{name}}", "src")
		part = strings.ReplaceAll(part, "{{cwd}}", projectDir)
		part = strings.ReplaceAll(part, "{{project}}", projectName)
		part = strings.ReplaceAll(part, "{{version}}", version)
		command[i] = part
	}

	// Add version-specific logic for different frameworks
	command = applyVersionSpecificOptions(command, meta.Framework, version)

	// Prepare the executable command
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("→ Running:", strings.Join(command, " "))
	return cmd.Run()
}

// validateVersion checks if the requested version is compatible with the template
func validateVersion(requestedVersion string, meta Metadata) error {
	if requestedVersion == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Check against minimum version
	if meta.MinVersion != "" {
		if utils.CompareVersions(requestedVersion, meta.MinVersion) < 0 {
			return fmt.Errorf("version %s is below minimum supported version %s for %s", 
				requestedVersion, meta.MinVersion, meta.Framework)
		}
	}

	// Framework-specific version validation
	switch meta.Framework {
	case "laravel":
		return validateLaravelVersion(requestedVersion)
	case "django":
		return validateDjangoVersion(requestedVersion)
	}

	return nil
}

// validateLaravelVersion checks Laravel-specific version constraints
func validateLaravelVersion(version string) error {
	// Laravel version constraints
	majorVersion := utils.ParseVersionPart(strings.Split(version, ".")[0])
	
	if majorVersion < 8 {
		return fmt.Errorf("Laravel version %s is too old (minimum supported: 8.0)", version)
	}
	
	if majorVersion > 12 {
		return fmt.Errorf("Laravel version %s is not yet supported (maximum: 12.x)", version)
	}

	return nil
}

// validateDjangoVersion checks Django-specific version constraints
func validateDjangoVersion(version string) error {
	// Django version constraints
	majorVersion := utils.ParseVersionPart(strings.Split(version, ".")[0])
	
	if majorVersion < 4 {
		return fmt.Errorf("Django version %s is too old (minimum supported: 4.0)", version)
	}
	
	if majorVersion > 6 {
		return fmt.Errorf("Django version %s is not yet supported (maximum: 6.x)", version)
	}

	return nil
}

// applyVersionSpecificOptions modifies the installation command based on framework and version
func applyVersionSpecificOptions(command []string, framework, version string) []string {
	switch framework {
	case "laravel":
		return applyLaravelVersionOptions(command, version)
	case "django":
		return applyDjangoVersionOptions(command, version)
	}
	
	return command
}

// applyLaravelVersionOptions adds Laravel version-specific installation options
func applyLaravelVersionOptions(command []string, version string) []string {
	// For Laravel, we need to specify the exact version constraint
	// Find the package name in the command and add version constraint
	for i, arg := range command {
		if arg == "laravel/laravel" {
			// Add version constraint: laravel/laravel:^11.0 for version 11
			majorVersion := strings.Split(version, ".")[0]
			command[i] = fmt.Sprintf("laravel/laravel:^%s.0", majorVersion)
			break
		}
	}
	
	return command
}

// applyDjangoVersionOptions adds Django version-specific installation options
func applyDjangoVersionOptions(command []string, version string) []string {
	// Django doesn't need version-specific startproject options
	// Version is controlled by the Django package installed in requirements.txt
	return command
}

// copyTemplateFiles copies AI context, Docker setup, and other template files (embedded or filesystem)
func copyTemplateFiles(projectDir, framework string, templatesFS, mcpServersFS embed.FS) error {
	// Copy AI context directory
	aiDstPath := filepath.Join(projectDir, "ai")
	fmt.Println("→ Copying AI context files...")
	
	// Try embedded first, fallback to filesystem
	embeddedAiPath := fmt.Sprintf("templates/%s/ai", framework)
	if err := copyEmbeddedDir(templatesFS, embeddedAiPath, aiDstPath); err != nil {
		// Fallback to filesystem
		aiSrcPath, pathErr := getFilesystemTemplateDir(framework, "ai")
		if pathErr == nil {
			if err := utils.CopyDir(aiSrcPath, aiDstPath); err != nil {
				return fmt.Errorf("failed to copy AI context: %w", err)
			}
		}
	}

	// Copy MCP server for the framework
	if err := copyMCPServer(framework, projectDir, mcpServersFS); err != nil {
		fmt.Printf("⚠️  Failed to copy MCP server: %v\n", err)
		// Don't fail the entire setup if MCP server copy fails
	}

	// Copy infrastructure directory (Docker setup)
	infraDstPath := filepath.Join(projectDir, "infra")
	fmt.Println("→ Copying Docker infrastructure...")
	
	// Try embedded first, fallback to filesystem
	embeddedInfraPath := fmt.Sprintf("templates/%s/infra", framework)
	if err := copyEmbeddedDir(templatesFS, embeddedInfraPath, infraDstPath); err != nil {
		// Fallback to filesystem
		infraSrcPath, pathErr := getFilesystemTemplateDir(framework, "infra")
		if pathErr == nil {
			if err := utils.CopyDir(infraSrcPath, infraDstPath); err != nil {
				return fmt.Errorf("failed to copy infrastructure: %w", err)
			}
		}
	}

	// Copy README.md
	readmeDstPath := filepath.Join(projectDir, "README.md")
	
	// Try embedded first, fallback to filesystem
	embeddedReadmePath := fmt.Sprintf("templates/%s/README.md", framework)
	if err := copyEmbeddedFile(templatesFS, embeddedReadmePath, readmeDstPath); err != nil {
		// Fallback to filesystem
		readmeSrcPath, pathErr := getFilesystemTemplatePath(framework, "README.md")
		if pathErr == nil {
			fmt.Println("→ Copying project README...")
			if err := utils.CopyFile(readmeSrcPath, readmeDstPath); err != nil {
				return fmt.Errorf("failed to copy README: %w", err)
			}
		}
	} else {
		fmt.Println("→ Copying project README...")
	}

	return nil
}

// copyMCPServer copies the framework-specific MCP server to the project
func copyMCPServer(framework, projectDir string, mcpServersFS embed.FS) error {
	embeddedMcpPath := fmt.Sprintf("mcp-servers/%s", framework)
	mcpDstPath := filepath.Join(projectDir, "ai", "mcp-server")

	fmt.Println("→ Copying MCP server for", framework, "...")
	
	// Try embedded first, fallback to filesystem
	if err := copyEmbeddedDir(mcpServersFS, embeddedMcpPath, mcpDstPath); err != nil {
		// Fallback to filesystem - find MCP servers relative to binary location
		filesystemMcpPath, pathErr := getFilesystemMCPPath(framework)
		if pathErr != nil {
			return fmt.Errorf("MCP server for %s not found: %w", framework, pathErr)
		}
		if err := utils.CopyDir(filesystemMcpPath, mcpDstPath); err != nil {
			return fmt.Errorf("failed to copy MCP server from filesystem: %w", err)
		}
	}

	// Install npm dependencies
	fmt.Println("→ Installing MCP server dependencies...")
	cmd := exec.Command("npm", "install")
	cmd.Dir = mcpDstPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install MCP server dependencies: %w", err)
	}

	fmt.Println("✅ MCP server ready!")
	fmt.Printf("   Add this to your Claude Code MCP settings:\n")
	fmt.Printf("   {\n")
	fmt.Printf("     \"steele-%s\": {\n", framework)
	fmt.Printf("       \"command\": \"node\",\n")
	fmt.Printf("       \"args\": [\"ai/mcp-server/index.js\"],\n")
	fmt.Printf("       \"cwd\": \"%s\"\n", projectDir)
	fmt.Printf("     }\n")
	fmt.Printf("   }\n")

	return nil
}

// finalizeProject registers the project and generates docker-compose.yml
func finalizeProject(meta Metadata, projectDir, projectName, version string) error {
	// Resolve project name from template
	resolvedName := meta.Name
	if resolvedName == "" || strings.Contains(resolvedName, "{{") {
		resolvedName = strings.ReplaceAll(resolvedName, "{{project}}", projectName)
		if resolvedName == "" || strings.Contains(resolvedName, "{{") {
			resolvedName = projectName
		}
	}

	// Register project in registry
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := reg.AddProject(resolvedName, projectDir, meta.Framework, version); err != nil {
		return fmt.Errorf("failed to register project: %w", err)
	}

	fmt.Printf("→ Project '%s' registered\n", resolvedName)

	// Generate docker-compose.yml from steele.json if it has services defined
	steeleJsonPath := filepath.Join(projectDir, "steele.json")
	if utils.FileExists(steeleJsonPath) {
		fmt.Println("→ Generating docker-compose.yml from steele.json...")
		if err := compose.GenerateDockerCompose(projectDir); err != nil {
			return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
		}
		fmt.Println("✅ docker-compose.yml generated")
	}

	return nil
}

// runPostInstall handles framework-specific setup after installation
func runPostInstall(meta Metadata, projectDir string) error {
	// Set up Laravel environment file
	if meta.Framework == "laravel" {
		return setupLaravel(projectDir)
	}

	// Set up Django environment
	if meta.Framework == "django" {
		return setupDjango(projectDir)
	}

	return nil
}

// setupLaravel performs Laravel-specific post-installation setup
func setupLaravel(projectDir string) error {
	srcDir := filepath.Join(projectDir, "src")
	
	// Copy .env.example to .env
	envExample := filepath.Join(srcDir, ".env.example")
	envFile := filepath.Join(srcDir, ".env")
	
	if utils.FileExists(envExample) && !utils.FileExists(envFile) {
		fmt.Println("→ Creating .env file from .env.example...")
		if err := utils.CopyFile(envExample, envFile); err != nil {
			return fmt.Errorf("failed to create .env file: %w", err)
		}
	}

	// Update .env with Docker database configuration
	fmt.Println("→ Configuring Laravel for Docker environment...")
	if err := updateLaravelEnv(envFile); err != nil {
		return fmt.Errorf("failed to update .env: %w", err)
	}

	// Check if Docker is available and start services
	if err := startDockerServices(projectDir); err != nil {
		fmt.Println("⚠️  Docker not available or failed to start services")
		fmt.Println("   Run 'docker-compose up -d' manually to start the development environment")
		return nil // Don't fail the entire setup if Docker isn't available
	}

	// Run Laravel setup commands
	return runLaravelSetup(projectDir)
}

// updateLaravelEnv updates the .env file with Docker-specific configuration
func updateLaravelEnv(envFile string) error {
	// Read current .env content
	content, err := os.ReadFile(envFile)
	if err != nil {
		return err
	}

	envContent := string(content)
	
	// Update database configuration for Docker
	envContent = strings.ReplaceAll(envContent, "DB_HOST=127.0.0.1", "DB_HOST=mysql")
	envContent = strings.ReplaceAll(envContent, "DB_DATABASE=laravel", "DB_DATABASE=laravel")
	envContent = strings.ReplaceAll(envContent, "DB_USERNAME=root", "DB_USERNAME=laravel")
	envContent = strings.ReplaceAll(envContent, "DB_PASSWORD=", "DB_PASSWORD=laravel")
	
	// Add Redis configuration
	if !strings.Contains(envContent, "REDIS_HOST=") {
		envContent += "\nREDIS_HOST=redis\n"
	}
	
	return os.WriteFile(envFile, []byte(envContent), 0644)
}

// startDockerServices attempts to start Docker services
func startDockerServices(projectDir string) error {
	fmt.Println("→ Starting Docker services...")
	
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// runLaravelSetup runs essential Laravel setup commands in Docker
func runLaravelSetup(projectDir string) error {
	commands := [][]string{
		{"docker-compose", "exec", "-T", "app", "composer", "install"},
		{"docker-compose", "exec", "-T", "app", "php", "artisan", "key:generate"},
		{"docker-compose", "exec", "-T", "app", "php", "artisan", "migrate", "--force"},
	}

	for _, command := range commands {
		fmt.Println("→ Running:", strings.Join(command, " "))
		
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Command failed: %s\n", strings.Join(command, " "))
			fmt.Println("   You may need to run this manually after Docker services are ready")
			continue // Continue with other commands
		}
	}

	return nil
}

// setupDjango performs Django-specific post-installation setup
func setupDjango(projectDir string) error {
	srcDir := filepath.Join(projectDir, "src")
	
	// Copy and update requirements.txt from Docker template
	requirementsSrc := filepath.Join(projectDir, "infra", "docker", "requirements.txt")
	requirementsDst := filepath.Join(srcDir, "requirements.txt")
	
	if utils.FileExists(requirementsSrc) {
		fmt.Println("→ Copying and updating requirements.txt...")
		if err := copyAndUpdateRequirements(requirementsSrc, requirementsDst, projectDir); err != nil {
			return fmt.Errorf("failed to copy requirements.txt: %w", err)
		}
	}

	// Check if Docker is available and start services
	if err := startDockerServices(projectDir); err != nil {
		fmt.Println("⚠️  Docker not available or failed to start services")
		fmt.Println("   Run 'docker-compose up -d' manually to start the development environment")
		return nil // Don't fail the entire setup if Docker isn't available
	}

	// Run Django setup commands
	return runDjangoSetup(projectDir)
}

// copyAndUpdateRequirements copies requirements.txt and updates Django version
func copyAndUpdateRequirements(src, dst, projectDir string) error {
	// Read the template requirements.txt
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read requirements template: %w", err)
	}

	// Get the requested Django version from the project name or context
	version := extractVersionFromProject(projectDir)
	
	if version != "" {
		// Update Django version in requirements
		reqContent := string(content)
		majorVersion := strings.Split(version, ".")[0]
		nextMajorVersion := fmt.Sprintf("%d", utils.ParseVersionPart(majorVersion)+1)
		
		// Replace Django version constraint
		oldConstraint := "Django>=5.0,<6.0"
		newConstraint := fmt.Sprintf("Django>=%s.0,<%s.0", majorVersion, nextMajorVersion)
		reqContent = strings.ReplaceAll(reqContent, oldConstraint, newConstraint)
		
		content = []byte(reqContent)
	}

	// Write the updated requirements.txt
	return os.WriteFile(dst, content, 0644)
}

// extractVersionFromProject attempts to extract the Django version from project context
func extractVersionFromProject(projectDir string) string {
	// For now, return default version
	// TODO: Pass version parameter through the entire call stack
	return "5"
}

// runDjangoSetup runs essential Django setup commands in Docker
func runDjangoSetup(projectDir string) error {
	commands := [][]string{
		{"docker-compose", "exec", "-T", "web", "pip", "install", "-r", "requirements.txt"},
		{"docker-compose", "exec", "-T", "web", "python", "manage.py", "migrate"},
		{"docker-compose", "exec", "-T", "web", "python", "manage.py", "collectstatic", "--noinput"},
	}

	for _, command := range commands {
		fmt.Println("→ Running:", strings.Join(command, " "))
		
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Command failed: %s\n", strings.Join(command, " "))
			fmt.Println("   You may need to run this manually after Docker services are ready")
			continue // Continue with other commands
		}
	}

	fmt.Println("🎉 Django setup complete!")
	fmt.Println("   → Django app: http://localhost:8000")
	fmt.Println("   → Django admin: http://localhost:8000/admin")
	fmt.Println("   → Mailhog: http://localhost:8025")
	fmt.Println("   → Create superuser: docker-compose exec web python manage.py createsuperuser")

	return nil
}


// copyEmbeddedFile copies a single file from embedded filesystem to local filesystem
func copyEmbeddedFile(fsys embed.FS, srcPath, dstPath string) error {
	// Read file from embedded filesystem
	data, err := fsys.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded file %s: %w", srcPath, err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Write file to local filesystem
	if err := os.WriteFile(dstPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dstPath, err)
	}

	return nil
}

// copyEmbeddedDir recursively copies a directory from embedded filesystem to local filesystem
func copyEmbeddedDir(fsys embed.FS, srcPath, dstPath string) error {
	// Create destination directory
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Walk through embedded directory
	return fs.WalkDir(fsys, srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Calculate destination path
		destPath := filepath.Join(dstPath, relPath)

		if d.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, 0755)
		} else {
			// Copy file
			return copyEmbeddedFile(fsys, path, destPath)
		}
	})
}

// getFilesystemTemplatePath finds template files relative to the binary location
func getFilesystemTemplatePath(framework, filename string) (string, error) {
	// Get the path to the current executable
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	execDir := filepath.Dir(executable)
	
	// Look for templates in the same directory as the binary
	templatePath := filepath.Join(execDir, "templates", framework, filename)
	if utils.FileExists(templatePath) {
		return templatePath, nil
	}
	
	// If not found, try looking in the parent directory (development mode)
	parentTemplateDir := filepath.Join(filepath.Dir(execDir), "templates", framework, filename)
	if utils.FileExists(parentTemplateDir) {
		return parentTemplateDir, nil
	}
	
	// Try looking in the current working directory (fallback)
	cwd, _ := os.Getwd()
	cwdTemplateDir := filepath.Join(cwd, "templates", framework, filename)
	if utils.FileExists(cwdTemplateDir) {
		return cwdTemplateDir, nil
	}
	
	return "", fmt.Errorf("template file %s not found for framework %s", filename, framework)
}

// getFilesystemTemplateDir finds template directories relative to the binary location
func getFilesystemTemplateDir(framework, subdir string) (string, error) {
	// Get the path to the current executable
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	execDir := filepath.Dir(executable)
	
	// Look for templates in the same directory as the binary
	templatePath := filepath.Join(execDir, "templates", framework, subdir)
	if utils.FileExists(templatePath) {
		return templatePath, nil
	}
	
	// If not found, try looking in the parent directory (development mode)
	parentTemplateDir := filepath.Join(filepath.Dir(execDir), "templates", framework, subdir)
	if utils.FileExists(parentTemplateDir) {
		return parentTemplateDir, nil
	}
	
	// Try looking in the current working directory (fallback)
	cwd, _ := os.Getwd()
	cwdTemplateDir := filepath.Join(cwd, "templates", framework, subdir)
	if utils.FileExists(cwdTemplateDir) {
		return cwdTemplateDir, nil
	}
	
	return "", fmt.Errorf("template directory %s not found for framework %s", subdir, framework)
}

// getFilesystemMCPPath finds MCP server directories relative to the binary location
func getFilesystemMCPPath(framework string) (string, error) {
	// Get the path to the current executable
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	execDir := filepath.Dir(executable)
	
	// Look for MCP servers in the same directory as the binary
	mcpPath := filepath.Join(execDir, "mcp-servers", framework)
	if utils.FileExists(mcpPath) {
		return mcpPath, nil
	}
	
	// If not found, try looking in the parent directory (development mode)
	parentMCPDir := filepath.Join(filepath.Dir(execDir), "mcp-servers", framework)
	if utils.FileExists(parentMCPDir) {
		return parentMCPDir, nil
	}
	
	// Try looking in the current working directory (fallback)
	cwd, _ := os.Getwd()
	cwdMCPDir := filepath.Join(cwd, "mcp-servers", framework)
	if utils.FileExists(cwdMCPDir) {
		return cwdMCPDir, nil
	}
	
	return "", fmt.Errorf("MCP server directory not found for framework %s", framework)
}
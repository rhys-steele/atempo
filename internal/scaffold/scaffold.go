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

	"atempo/internal/compose"
	"atempo/internal/logger"
	"atempo/internal/mcp"
	"atempo/internal/registry"
	"atempo/internal/utils"
)

// Installer defines how a framework should be installed.
// This includes the command to run and the working directory context.
type Installer struct {
	Type    string   `json:"type"`     // e.g., "composer", "docker", "shell"
	Command []string `json:"command"`  // Full command with args (supports templating)
	WorkDir string   `json:"work-dir"` // Directory to run the command in
}

// Metadata describes a Atempo template's configuration,
// including language, installer, and framework compatibility.
type Metadata struct {
	Name       string    `json:"name"`        // Project name template
	Framework  string    `json:"framework"`   // e.g., "laravel"
	Language   string    `json:"language"`    // e.g., "php"
	Installer  Installer `json:"installer"`   // How to scaffold the source code
	WorkingDir string    `json:"working-dir"` // Expected project root path in container, e.g., /var/www
	MinVersion string    `json:"min-version"` // Minimum supported version (semantic)
}

// Run executes the scaffolding process for the given framework and version.
// It loads the template's `atempo.json`, performs template substitution,
// runs the specified install command, and copies template files.
func Run(framework string, version string, templatesFS, mcpServersFS embed.FS) error {
	// Get the current working directory (user's target project root)
	projectDir, _ := os.Getwd()
	projectName := filepath.Base(projectDir)

	// Create quiet logger for this project (progress shown by caller)
	log, err := logger.NewQuiet(projectName)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer log.Close()

	// Log file location is only shown in verbose mode or on error

	// Step 1: Load and validate template configuration
	loadStep := log.StartStep("Loading template configuration")
	// Load atempo.json (try embedded first, fallback to filesystem)
	var metaBytes []byte

	// Try embedded first
	embeddedPath := fmt.Sprintf("templates/frameworks/%s/atempo.json", framework)
	metaBytes, readErr := templatesFS.ReadFile(embeddedPath)
	if readErr != nil {
		// Fallback to filesystem - find templates relative to binary location
		filesystemPath, pathErr := getFilesystemTemplatePath(framework, "atempo.json")
		if pathErr != nil {
			log.ErrorStep(loadStep, fmt.Errorf("could not locate atempo.json for %s: %w", framework, pathErr))
			return fmt.Errorf("could not locate atempo.json for %s: %w", framework, pathErr)
		}
		metaBytes, readErr = os.ReadFile(filesystemPath)
		if readErr != nil {
			log.ErrorStep(loadStep, fmt.Errorf("could not read atempo.json for %s: %w", framework, readErr))
			return fmt.Errorf("could not read atempo.json for %s: %w", framework, readErr)
		}
	}

	// Parse the metadata JSON into a structured object
	var meta Metadata
	if parseErr := json.Unmarshal(metaBytes, &meta); parseErr != nil {
		log.ErrorStep(loadStep, fmt.Errorf("invalid atempo.json: %w", parseErr))
		return fmt.Errorf("invalid atempo.json: %w", parseErr)
	}

	// Validate version compatibility
	if validateErr := validateVersion(version, meta); validateErr != nil {
		log.ErrorStep(loadStep, fmt.Errorf("version validation failed: %w", validateErr))
		return fmt.Errorf("version validation failed: %w", validateErr)
	}

	log.CompleteStep(loadStep)

	// Step 2: Run the framework installer (e.g., composer create-project)
	installStep := log.StartStep(fmt.Sprintf("Installing %s %s application", framework, version))
	if err := runInstaller(log, installStep, meta, projectDir, projectName, version); err != nil {
		log.ErrorStep(installStep, err)
		return fmt.Errorf("installer failed: %w", err)
	}
	log.CompleteStep(installStep)

	// Step 3: Copy template files (AI context, Docker setup, etc.)
	copyStep := log.StartStep("Copying template files")
	if err := copyTemplateFiles(log, copyStep, projectDir, projectName, meta.Framework, version, templatesFS, mcpServersFS); err != nil {
		log.ErrorStep(copyStep, err)
		return fmt.Errorf("failed to copy template files: %w", err)
	}
	log.CompleteStep(copyStep)

	// Step 4: Run post-installation setup
	postStep := log.StartStep("Running post-installation setup")
	if err := runPostInstall(log, postStep, meta, projectDir); err != nil {
		log.ErrorStep(postStep, err)
		return fmt.Errorf("post-installation failed: %w", err)
	}
	log.CompleteStep(postStep)

	// Step 5: Register project and generate docker-compose
	finalStep := log.StartStep("Registering project and generating docker-compose")
	if err := finalizeProject(log, finalStep, meta, projectDir, projectName, version); err != nil {
		log.WarningStep(finalStep, err.Error())
	} else {
		log.CompleteStep(finalStep)
	}

	log.PrintSummary()
	return nil
}

// runInstaller executes the framework installation command
func runInstaller(log *logger.Logger, step *logger.Step, meta Metadata, projectDir, projectName, version string) error {
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

	// Check if Docker is required and available
	if meta.Installer.Type == "docker" && command[0] == "docker" {
		if err := checkDockerAvailability(); err != nil {
			return fmt.Errorf("Docker is required but not available: %w", err)
		}
	}

	// Prepare the executable command
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = projectDir

	// Use logger to capture command output
	return log.RunCommand(step, cmd)
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
func copyTemplateFiles(log *logger.Logger, step *logger.Step, projectDir, projectName, framework, version string, templatesFS, mcpServersFS embed.FS) error {
	// Copy AI context directory
	aiDstPath := filepath.Join(projectDir, "ai")

	// Try embedded first, fallback to filesystem
	embeddedAiPath := fmt.Sprintf("templates/frameworks/%s/ai", framework)
	if err := copyEmbeddedDirWithContext(templatesFS, embeddedAiPath, aiDstPath, projectName, projectDir, version); err != nil {
		// Fallback to filesystem
		aiSrcPath, pathErr := getFilesystemTemplateDir(framework, "ai")
		if pathErr == nil {
			if err := copyFilesystemDirWithContext(aiSrcPath, aiDstPath, projectName, projectDir, version); err != nil {
				return fmt.Errorf("failed to copy AI context: %w", err)
			}
		}
	}

	// Copy MCP server for the framework
	if err := copyMCPServer(log, step, framework, projectDir, mcpServersFS); err != nil {
		log.WarningStep(step, fmt.Sprintf("Failed to copy MCP server: %v", err))
		// Don't fail the entire setup if MCP server copy fails
	}

	// Copy infrastructure directory (Docker setup)
	infraDstPath := filepath.Join(projectDir, "infra")

	// Try embedded first, fallback to filesystem
	embeddedInfraPath := fmt.Sprintf("templates/frameworks/%s/infra", framework)
	if err := copyEmbeddedDirWithContext(templatesFS, embeddedInfraPath, infraDstPath, projectName, projectDir, version); err != nil {
		// Fallback to filesystem
		infraSrcPath, pathErr := getFilesystemTemplateDir(framework, "infra")
		if pathErr == nil {
			if err := copyFilesystemDirWithContext(infraSrcPath, infraDstPath, projectName, projectDir, version); err != nil {
				return fmt.Errorf("failed to copy infrastructure: %w", err)
			}
		}
	}

	// Copy README.md
	readmeDstPath := filepath.Join(projectDir, "README.md")

	// Try embedded first, fallback to filesystem
	embeddedReadmePath := fmt.Sprintf("templates/frameworks/%s/README.md", framework)
	if err := copyEmbeddedFileWithContext(templatesFS, embeddedReadmePath, readmeDstPath, projectName, projectDir, version); err != nil {
		// Fallback to filesystem
		readmeSrcPath, pathErr := getFilesystemTemplatePath(framework, "README.md")
		if pathErr == nil {
			if err := copyFilesystemFileWithContext(readmeSrcPath, readmeDstPath, projectName, projectDir, version); err != nil {
				return fmt.Errorf("failed to copy README: %w", err)
			}
		}
	}

	// Copy atempo.json for Docker Compose and DNS setup
	atempoJsonDstPath := filepath.Join(projectDir, "atempo.json")

	// Try embedded first, fallback to filesystem
	embeddedAtempoJsonPath := fmt.Sprintf("templates/frameworks/%s/atempo.json", framework)
	if err := copyEmbeddedFileWithContext(templatesFS, embeddedAtempoJsonPath, atempoJsonDstPath, projectName, projectDir, version); err != nil {
		// Fallback to filesystem
		atempoJsonSrcPath, pathErr := getFilesystemTemplatePath(framework, "atempo.json")
		if pathErr != nil {
			return fmt.Errorf("failed to find atempo.json template: embedded error: %v, filesystem error: %w", err, pathErr)
		}
		if err := copyFilesystemFileWithContext(atempoJsonSrcPath, atempoJsonDstPath, projectName, projectDir, version); err != nil {
			return fmt.Errorf("failed to copy atempo.json from filesystem: %w", err)
		}
	}

	return nil
}

// copyMCPServer discovers and installs the best available MCP server for the framework
func copyMCPServer(log *logger.Logger, step *logger.Step, framework, projectDir string, mcpServersFS embed.FS) error {
	mcpDstPath := filepath.Join(projectDir, "ai", "mcp-server")

	// Discover available MCP servers
	discovery, err := mcp.DiscoverMCPServers(framework)
	if err != nil {
		return fmt.Errorf("failed to discover MCP servers: %w", err)
	}

	var selectedServer mcp.MCPServer
	var serverType string

	// Prefer official servers, then community, then generated
	if len(discovery.Official) > 0 {
		selectedServer = discovery.Official[0]
		serverType = "official"
	} else if len(discovery.Community) > 0 {
		selectedServer = discovery.Community[0]
		serverType = "community"
	} else if discovery.Generated != nil {
		selectedServer = *discovery.Generated
		serverType = "generated"
	} else {
		return fmt.Errorf("no MCP servers available for %s", framework)
	}

	// Install the selected server
	if serverType == "generated" {
		// Use template-based generation
		projectInfo := mcp.ProjectInfo{
			Name:      filepath.Base(projectDir),
			Framework: framework,
			Path:      projectDir,
		}

		if err := mcp.GenerateServerFromTemplate(selectedServer, projectInfo, mcpDstPath); err != nil {
			return fmt.Errorf("failed to generate MCP server: %w", err)
		}
	} else {
		// Install official/community server
		if err := mcp.InstallMCPServer(selectedServer, projectDir); err != nil {
			return fmt.Errorf("failed to install MCP server: %w", err)
		}
	}

	// Install npm dependencies
	cmd := exec.Command("npm", "install")
	cmd.Dir = mcpDstPath

	if err := log.RunCommand(step, cmd); err != nil {
		return fmt.Errorf("failed to install MCP server dependencies: %w", err)
	}

	return nil
}

// finalizeProject registers the project and generates docker-compose.yml
func finalizeProject(log *logger.Logger, step *logger.Step, meta Metadata, projectDir, projectName, version string) error {
	// Resolve project name from template
	resolvedName := meta.Name
	if resolvedName == "" || strings.Contains(resolvedName, "{{") {
		resolvedName = strings.ReplaceAll(resolvedName, "{{project}}", projectName)
		if resolvedName == "" || strings.Contains(resolvedName, "{{") {
			resolvedName = projectName
		}
	}
	
	// Use the basename of the project directory as the registry name
	// This ensures projects work with simple names even if created with paths
	registryName := filepath.Base(projectDir)

	// Register project in registry
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := reg.AddProject(registryName, projectDir, meta.Framework, version); err != nil {
		return fmt.Errorf("failed to register project: %w", err)
	}

	// Generate docker-compose.yml from atempo.json if it has services defined
	atempoJsonPath := filepath.Join(projectDir, "atempo.json")
	if utils.FileExists(atempoJsonPath) {
		if err := compose.GenerateDockerCompose(projectDir); err != nil {
			return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
		}
	}

	return nil
}

// runPostInstall handles framework-specific setup after installation
func runPostInstall(log *logger.Logger, step *logger.Step, meta Metadata, projectDir string) error {
	// Set up Laravel environment file
	if meta.Framework == "laravel" {
		return setupLaravel(log, step, projectDir)
	}

	// Set up Django environment
	if meta.Framework == "django" {
		return setupDjango(log, step, projectDir)
	}

	return nil
}

// setupLaravel performs Laravel-specific post-installation setup
func setupLaravel(log *logger.Logger, step *logger.Step, projectDir string) error {
	srcDir := filepath.Join(projectDir, "src")

	// Copy .env.example to .env
	envExample := filepath.Join(srcDir, ".env.example")
	envFile := filepath.Join(srcDir, ".env")

	if utils.FileExists(envExample) && !utils.FileExists(envFile) {
		if err := utils.CopyFile(envExample, envFile); err != nil {
			return fmt.Errorf("failed to create .env file: %w", err)
		}
	}

	// Update .env with Docker database configuration
	if err := updateLaravelEnv(envFile); err != nil {
		return fmt.Errorf("failed to update .env: %w", err)
	}

	// Check if Docker is available and start services
	if err := startDockerServices(log, step, projectDir); err != nil {
		log.WarningStep(step, "Docker not available or failed to start services - run 'docker-compose up -d' manually")
		return nil // Don't fail the entire setup if Docker isn't available
	}

	// Run Laravel setup commands
	return runLaravelSetup(log, step, projectDir)
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
func startDockerServices(log *logger.Logger, step *logger.Step, projectDir string) error {
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = projectDir

	return log.RunCommand(step, cmd)
}

// runLaravelSetup runs essential Laravel setup commands in Docker
func runLaravelSetup(log *logger.Logger, step *logger.Step, projectDir string) error {
	commands := [][]string{
		{"docker-compose", "exec", "-T", "app", "composer", "install"},
		{"docker-compose", "exec", "-T", "app", "php", "artisan", "key:generate"},
		{"docker-compose", "exec", "-T", "app", "php", "artisan", "migrate", "--force"},
	}

	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Dir = projectDir

		if err := log.RunCommand(step, cmd); err != nil {
			log.WarningStep(step, fmt.Sprintf("Command failed: %s - you may need to run this manually", strings.Join(command, " ")))
			continue // Continue with other commands
		}
	}

	return nil
}

// setupDjango performs Django-specific post-installation setup
func setupDjango(log *logger.Logger, step *logger.Step, projectDir string) error {
	srcDir := filepath.Join(projectDir, "src")

	// Copy and update requirements.txt from Docker template
	requirementsSrc := filepath.Join(projectDir, "infra", "docker", "requirements.txt")
	requirementsDst := filepath.Join(srcDir, "requirements.txt")

	if utils.FileExists(requirementsSrc) {
		if err := copyAndUpdateRequirements(requirementsSrc, requirementsDst, projectDir); err != nil {
			return fmt.Errorf("failed to copy requirements.txt: %w", err)
		}
	}

	// Check if Docker is available and start services
	if err := startDockerServices(log, step, projectDir); err != nil {
		log.WarningStep(step, "Docker not available or failed to start services - run 'docker-compose up -d' manually")
		return nil // Don't fail the entire setup if Docker isn't available
	}

	// Run Django setup commands
	return runDjangoSetup(log, step, projectDir)
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
func runDjangoSetup(log *logger.Logger, step *logger.Step, projectDir string) error {
	commands := [][]string{
		{"docker-compose", "exec", "-T", "web", "pip", "install", "-r", "requirements.txt"},
		{"docker-compose", "exec", "-T", "web", "python", "manage.py", "migrate"},
		{"docker-compose", "exec", "-T", "web", "python", "manage.py", "collectstatic", "--noinput"},
	}

	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Dir = projectDir

		if err := log.RunCommand(step, cmd); err != nil {
			log.WarningStep(step, fmt.Sprintf("Command failed: %s - you may need to run this manually", strings.Join(command, " ")))
			continue // Continue with other commands
		}
	}

	return nil
}

// processTemplateContent processes template variables in content
func processTemplateContent(content string, projectName, projectDir, version string) string {
	content = strings.ReplaceAll(content, "{{project}}", projectName)
	content = strings.ReplaceAll(content, "{{name}}", "src")
	content = strings.ReplaceAll(content, "{{cwd}}", projectDir)
	content = strings.ReplaceAll(content, "{{version}}", version)
	return content
}

// copyEmbeddedFile copies a single file from embedded filesystem to local filesystem with template processing
func copyEmbeddedFile(fsys embed.FS, srcPath, dstPath string) error {
	return copyEmbeddedFileWithContext(fsys, srcPath, dstPath, "", "", "")
}

// copyEmbeddedFileWithContext copies a file with template variable processing
func copyEmbeddedFileWithContext(fsys embed.FS, srcPath, dstPath, projectName, projectDir, version string) error {
	// Read file from embedded filesystem
	data, err := fsys.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded file %s: %w", srcPath, err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Process template variables if context is provided
	var processedData []byte
	if projectName != "" {
		content := string(data)
		processedContent := processTemplateContent(content, projectName, projectDir, version)
		processedData = []byte(processedContent)
	} else {
		processedData = data
	}

	// Write file to local filesystem
	if err := os.WriteFile(dstPath, processedData, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dstPath, err)
	}

	return nil
}

// copyEmbeddedDir recursively copies a directory from embedded filesystem to local filesystem
func copyEmbeddedDir(fsys embed.FS, srcPath, dstPath string) error {
	return copyEmbeddedDirWithContext(fsys, srcPath, dstPath, "", "", "")
}

// copyEmbeddedDirWithContext recursively copies a directory with template variable processing
func copyEmbeddedDirWithContext(fsys embed.FS, srcPath, dstPath, projectName, projectDir, version string) error {
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
			// Copy file with template processing
			return copyEmbeddedFileWithContext(fsys, path, destPath, projectName, projectDir, version)
		}
	})
}

// copyFilesystemDirWithContext copies a directory from filesystem with template processing
func copyFilesystemDirWithContext(srcPath, dstPath, projectName, projectDir, version string) error {
	// Create destination directory
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Walk through filesystem directory
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
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

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, 0755)
		} else {
			// Copy file with template processing
			return copyFilesystemFileWithContext(path, destPath, projectName, projectDir, version)
		}
	})
}

// copyFilesystemFileWithContext copies a file from filesystem with template processing
func copyFilesystemFileWithContext(srcPath, dstPath, projectName, projectDir, version string) error {
	// Read file from filesystem
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", srcPath, err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Process template variables
	var processedData []byte
	if projectName != "" {
		content := string(data)
		processedContent := processTemplateContent(content, projectName, projectDir, version)
		processedData = []byte(processedContent)
	} else {
		processedData = data
	}

	// Write file to destination
	if err := os.WriteFile(dstPath, processedData, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dstPath, err)
	}

	return nil
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
	templatePath := filepath.Join(execDir, "templates", "frameworks", framework, filename)
	if utils.FileExists(templatePath) {
		return templatePath, nil
	}

	// If not found, try looking in the parent directory (development mode)
	parentTemplateDir := filepath.Join(filepath.Dir(execDir), "templates", "frameworks", framework, filename)
	if utils.FileExists(parentTemplateDir) {
		return parentTemplateDir, nil
	}

	// Try looking in the current working directory (fallback)
	cwd, _ := os.Getwd()
	cwdTemplateDir := filepath.Join(cwd, "templates", "frameworks", framework, filename)
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
	templatePath := filepath.Join(execDir, "templates", "frameworks", framework, subdir)
	if utils.FileExists(templatePath) {
		return templatePath, nil
	}

	// If not found, try looking in the parent directory (development mode)
	parentTemplateDir := filepath.Join(filepath.Dir(execDir), "templates", "frameworks", framework, subdir)
	if utils.FileExists(parentTemplateDir) {
		return parentTemplateDir, nil
	}

	// Try looking in the current working directory (fallback)
	cwd, _ := os.Getwd()
	cwdTemplateDir := filepath.Join(cwd, "templates", "frameworks", framework, subdir)
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

// checkDockerAvailability verifies that Docker is installed and running
func checkDockerAvailability() error {
	// First check if docker command is available
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker command not found in PATH")
	}
	
	// Check if Docker daemon is running by running a simple command
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker daemon is not running")
	}
	
	return nil
}

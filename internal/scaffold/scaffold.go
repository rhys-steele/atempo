package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	Framework   string    `json:"framework"`    // e.g., "laravel"
	Language    string    `json:"language"`     // e.g., "php"
	Installer   Installer `json:"installer"`    // How to scaffold the source code
	WorkingDir  string    `json:"working-dir"`  // Expected project root path in container, e.g., /var/www
	MinVersion  string    `json:"min-version"`  // Minimum supported version (semantic)
}

// Run executes the scaffolding process for the given framework and version.
// It loads the template's `steele.json`, performs template substitution,
// runs the specified install command, and copies template files.
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
	if err := copyTemplateFiles(templatePath, projectDir); err != nil {
		return fmt.Errorf("failed to copy template files: %w", err)
	}

	// Step 3: Run post-installation setup
	fmt.Println("⚙️  Running post-installation setup...")
	if err := runPostInstall(meta, projectDir); err != nil {
		return fmt.Errorf("post-installation failed: %w", err)
	}

	return nil
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

// copyTemplateFiles copies AI context, Docker setup, and other template files
func copyTemplateFiles(templatePath, projectDir string) error {
	// Copy AI context directory
	aiSrc := filepath.Join(templatePath, "ai")
	aiDst := filepath.Join(projectDir, "ai")
	if utils.FileExists(aiSrc) {
		fmt.Println("→ Copying AI context files...")
		if err := utils.CopyDir(aiSrc, aiDst); err != nil {
			return fmt.Errorf("failed to copy AI context: %w", err)
		}
	}

	// Copy infrastructure directory (Docker setup)
	infraSrc := filepath.Join(templatePath, "infra")
	infraDst := filepath.Join(projectDir, "infra")
	if utils.FileExists(infraSrc) {
		fmt.Println("→ Copying Docker infrastructure...")
		if err := utils.CopyDir(infraSrc, infraDst); err != nil {
			return fmt.Errorf("failed to copy infrastructure: %w", err)
		}
	}

	// Copy README.md
	readmeSrc := filepath.Join(templatePath, "README.md")
	readmeDst := filepath.Join(projectDir, "README.md")
	if utils.FileExists(readmeSrc) {
		fmt.Println("→ Copying project README...")
		if err := utils.CopyFile(readmeSrc, readmeDst); err != nil {
			return fmt.Errorf("failed to copy README: %w", err)
		}
	}

	return nil
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

// copyAndUpdateRequirements copies requirements.txt and updates Django version
func copyAndUpdateRequirements(src, dst, projectDir string) error {
	// Read the template requirements.txt
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read requirements template: %w", err)
	}

	// Get the requested Django version from the project name or context
	// For now, we'll extract it from the command that was run
	// This is a bit hacky but works for the current implementation
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
// This is a temporary solution until we pass version through the call stack
func extractVersionFromProject(projectDir string) string {
	// For now, return default version
	// TODO: Pass version parameter through the entire call stack
	return "5"
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

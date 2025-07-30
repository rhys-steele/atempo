package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FrameworkInfo represents metadata about a framework
type FrameworkInfo struct {
	Name            string
	Language        string
	LatestVersion   string
	SupportedVersions []string
	ConfigFiles     []string
	DetectionFiles  []string
}

// FrameworkDetector provides framework detection and parsing utilities
type FrameworkDetector struct {
	frameworks map[string]FrameworkInfo
}

// NewFrameworkDetector creates a new framework detector
func NewFrameworkDetector() *FrameworkDetector {
	return &FrameworkDetector{
		frameworks: map[string]FrameworkInfo{
			"laravel": {
				Name:            "Laravel",
				Language:        "php",
				LatestVersion:   "11",
				SupportedVersions: []string{"8", "9", "10", "11"},
				ConfigFiles:     []string{"composer.json", "artisan"},
				DetectionFiles:  []string{"composer.json", "artisan", "app/Http/Kernel.php"},
			},
			"django": {
				Name:            "Django",
				Language:        "python",
				LatestVersion:   "5",
				SupportedVersions: []string{"3", "4", "5"},
				ConfigFiles:     []string{"requirements.txt", "manage.py"},
				DetectionFiles:  []string{"manage.py", "requirements.txt", "settings.py"},
			},
			"nextjs": {
				Name:            "Next.js",
				Language:        "javascript",
				LatestVersion:   "14",
				SupportedVersions: []string{"12", "13", "14"},
				ConfigFiles:     []string{"package.json", "next.config.js"},
				DetectionFiles:  []string{"package.json", "next.config.js", "pages/", "app/"},
			},
			"react": {
				Name:            "React",
				Language:        "javascript",
				LatestVersion:   "18",
				SupportedVersions: []string{"16", "17", "18"},
				ConfigFiles:     []string{"package.json"},
				DetectionFiles:  []string{"package.json", "public/index.html", "src/App.js"},
			},
			"vue": {
				Name:            "Vue.js",
				Language:        "javascript",
				LatestVersion:   "3",
				SupportedVersions: []string{"2", "3"},
				ConfigFiles:     []string{"package.json", "vue.config.js"},
				DetectionFiles:  []string{"package.json", "vue.config.js", "src/main.js"},
			},
			"nuxt": {
				Name:            "Nuxt.js",
				Language:        "javascript",
				LatestVersion:   "3",
				SupportedVersions: []string{"2", "3"},
				ConfigFiles:     []string{"package.json", "nuxt.config.js"},
				DetectionFiles:  []string{"package.json", "nuxt.config.js", "pages/", "layouts/"},
			},
			"express": {
				Name:            "Express",
				Language:        "javascript",
				LatestVersion:   "4.18.0",
				SupportedVersions: []string{"4.18.0", "4.19.0", "5.0.0"},
				ConfigFiles:     []string{"package.json"},
				DetectionFiles:  []string{"package.json", "app.js", "server.js"},
			},
			"fastapi": {
				Name:            "FastAPI",
				Language:        "python",
				LatestVersion:   "0.104",
				SupportedVersions: []string{"0.100", "0.104"},
				ConfigFiles:     []string{"requirements.txt", "pyproject.toml"},
				DetectionFiles:  []string{"main.py", "app.py", "requirements.txt"},
			},
			"rails": {
				Name:            "Rails",
				Language:        "ruby",
				LatestVersion:   "7",
				SupportedVersions: []string{"6", "7"},
				ConfigFiles:     []string{"Gemfile", "config/application.rb"},
				DetectionFiles:  []string{"Gemfile", "config/application.rb", "app/controllers/"},
			},
			"spring": {
				Name:            "Spring Boot",
				Language:        "java",
				LatestVersion:   "3",
				SupportedVersions: []string{"2", "3"},
				ConfigFiles:     []string{"pom.xml", "build.gradle"},
				DetectionFiles:  []string{"pom.xml", "build.gradle", "src/main/java/"},
			},
			"dotnet": {
				Name:            ".NET",
				Language:        "csharp",
				LatestVersion:   "8",
				SupportedVersions: []string{"6", "7", "8"},
				ConfigFiles:     []string{"*.csproj", "appsettings.json"},
				DetectionFiles:  []string{"*.csproj", "Program.cs", "Startup.cs"},
			},
		},
	}
}

// ParseFrameworkArg parses a framework argument in the format "framework:version"
func (fd *FrameworkDetector) ParseFrameworkArg(arg string) (framework, version string, err error) {
	if strings.Contains(arg, ":") {
		parts := strings.Split(arg, ":")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("expected format <framework>[:<version>], got: %s", arg)
		}
		framework = strings.TrimSpace(parts[0])
		version = strings.TrimSpace(parts[1])
	} else {
		framework = strings.TrimSpace(arg)
		version = fd.GetLatestVersion(framework)
	}
	
	// Validate framework exists
	if !fd.IsValidFramework(framework) {
		return "", "", fmt.Errorf("unknown framework: %s", framework)
	}
	
	return framework, version, nil
}

// GetLatestVersion returns the latest supported version for a framework
func (fd *FrameworkDetector) GetLatestVersion(framework string) string {
	if info, exists := fd.frameworks[framework]; exists {
		return info.LatestVersion
	}
	return "latest"
}

// IsValidFramework checks if a framework is supported
func (fd *FrameworkDetector) IsValidFramework(framework string) bool {
	_, exists := fd.frameworks[framework]
	return exists
}

// GetFrameworkInfo returns detailed information about a framework
func (fd *FrameworkDetector) GetFrameworkInfo(framework string) (FrameworkInfo, bool) {
	info, exists := fd.frameworks[framework]
	return info, exists
}

// GetSupportedFrameworks returns a list of all supported frameworks
func (fd *FrameworkDetector) GetSupportedFrameworks() []string {
	frameworks := make([]string, 0, len(fd.frameworks))
	for name := range fd.frameworks {
		frameworks = append(frameworks, name)
	}
	return frameworks
}

// DetectFramework attempts to detect the framework used in a project directory
func (fd *FrameworkDetector) DetectFramework(projectPath string) (string, error) {
	// First, try to detect from atempo.json
	if framework, err := fd.detectFromAtempoJSON(projectPath); err == nil {
		return framework, nil
	}
	
	// Then try to detect from docker-compose.yml
	if framework, err := fd.detectFromDockerCompose(projectPath); err == nil {
		return framework, nil
	}
	
	// Finally, try to detect from project files
	return fd.detectFromProjectFiles(projectPath)
}

// detectFromAtempoJSON attempts to detect framework from atempo.json
func (fd *FrameworkDetector) detectFromAtempoJSON(projectPath string) (string, error) {
	atempoFile := filepath.Join(projectPath, "atempo.json")
	if _, err := os.Stat(atempoFile); os.IsNotExist(err) {
		return "", fmt.Errorf("atempo.json not found")
	}
	
	data, err := os.ReadFile(atempoFile)
	if err != nil {
		return "", fmt.Errorf("failed to read atempo.json: %w", err)
	}
	
	var config struct {
		Framework string `json:"framework"`
	}
	
	if err := json.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("failed to parse atempo.json: %w", err)
	}
	
	if config.Framework == "" {
		return "", fmt.Errorf("no framework specified in atempo.json")
	}
	
	return config.Framework, nil
}

// detectFromDockerCompose attempts to detect framework from docker-compose.yml
func (fd *FrameworkDetector) detectFromDockerCompose(projectPath string) (string, error) {
	composeFile := filepath.Join(projectPath, "docker-compose.yml")
	if _, err := os.Stat(composeFile); os.IsNotExist(err) {
		return "", fmt.Errorf("docker-compose.yml not found")
	}
	
	data, err := os.ReadFile(composeFile)
	if err != nil {
		return "", fmt.Errorf("failed to read docker-compose.yml: %w", err)
	}
	
	content := string(data)
	
	// Look for framework-specific patterns in docker-compose.yml
	if strings.Contains(content, "composer") && strings.Contains(content, "artisan") {
		return "laravel", nil
	}
	if strings.Contains(content, "manage.py") || strings.Contains(content, "django") {
		return "django", nil
	}
	if strings.Contains(content, "next") || strings.Contains(content, "Next.js") {
		return "nextjs", nil
	}
	if strings.Contains(content, "npm") && strings.Contains(content, "react") {
		return "react", nil
	}
	
	return "", fmt.Errorf("could not detect framework from docker-compose.yml")
}

// detectFromProjectFiles attempts to detect framework from project files
func (fd *FrameworkDetector) detectFromProjectFiles(projectPath string) (string, error) {
	// Check each framework's detection files
	for frameworkName, info := range fd.frameworks {
		if fd.hasDetectionFiles(projectPath, info.DetectionFiles) {
			return frameworkName, nil
		}
	}
	
	return "", fmt.Errorf("could not detect framework from project files")
}

// hasDetectionFiles checks if the project has the required detection files
func (fd *FrameworkDetector) hasDetectionFiles(projectPath string, detectionFiles []string) bool {
	foundFiles := 0
	requiredFiles := len(detectionFiles)
	
	for _, file := range detectionFiles {
		fullPath := filepath.Join(projectPath, file)
		if _, err := os.Stat(fullPath); err == nil {
			foundFiles++
		}
	}
	
	// Require at least half of the detection files to be present
	return foundFiles >= (requiredFiles+1)/2
}

// ValidateVersion checks if a version is supported for a framework
func (fd *FrameworkDetector) ValidateVersion(framework, version string) error {
	info, exists := fd.frameworks[framework]
	if !exists {
		return fmt.Errorf("unknown framework: %s", framework)
	}
	
	if version == "latest" {
		return nil // "latest" is always valid
	}
	
	// Check if version is in supported versions
	for _, supportedVersion := range info.SupportedVersions {
		if version == supportedVersion {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported version %s for framework %s. Supported versions: %s", 
		version, framework, strings.Join(info.SupportedVersions, ", "))
}

// GetFrameworkLanguage returns the programming language for a given framework
func GetFrameworkLanguage(framework string) string {
	detector := NewFrameworkDetector()
	if info, exists := detector.GetFrameworkInfo(framework); exists {
		return info.Language
	}
	return "unknown"
}

// Global framework detector instance
var GlobalFrameworkDetector = NewFrameworkDetector()

// Helper functions for common framework operations
func ParseFrameworkArg(arg string) (framework, version string, err error) {
	return GlobalFrameworkDetector.ParseFrameworkArg(arg)
}

func GetLatestVersion(framework string) string {
	return GlobalFrameworkDetector.GetLatestVersion(framework)
}

func IsValidFramework(framework string) bool {
	return GlobalFrameworkDetector.IsValidFramework(framework)
}

func DetectFramework(projectPath string) (string, error) {
	return GlobalFrameworkDetector.DetectFramework(projectPath)
}

func ValidateVersion(framework, version string) error {
	return GlobalFrameworkDetector.ValidateVersion(framework, version)
}

func GetSupportedFrameworks() []string {
	return GlobalFrameworkDetector.GetSupportedFrameworks()
}

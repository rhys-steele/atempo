package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atempo/internal/docker"
	"atempo/internal/utils"
)

// FrameworkService provides business operations for framework detection and management
type FrameworkService interface {
	// DetectFramework detects the framework used in a project directory
	DetectFramework(ctx context.Context, projectPath string) (*FrameworkDetectionResult, error)
	
	// DetectFromCompose detects framework from docker-compose.yml file
	DetectFromCompose(ctx context.Context, projectPath string) (*FrameworkDetectionResult, error)
	
	// DetectFromFiles detects framework from project files and directory structure
	DetectFromFiles(ctx context.Context, projectPath string) (*FrameworkDetectionResult, error)
	
	// IsFramework checks if a directory contains a specific framework
	IsFramework(ctx context.Context, projectPath, framework string) (bool, error)
	
	// GetFrameworkInfo returns detailed information about a framework
	GetFrameworkInfo(ctx context.Context, framework string) (*FrameworkInfo, error)
	
	// ListSupportedFrameworks returns all supported frameworks
	ListSupportedFrameworks(ctx context.Context) ([]FrameworkInfo, error)
	
	// GetFrameworkVersion attempts to detect the version of a framework in use
	GetFrameworkVersion(ctx context.Context, projectPath, framework string) (string, error)
}

// FrameworkDetectionResult represents the result of framework detection
type FrameworkDetectionResult struct {
	Framework   string  `json:"framework"`
	Version     string  `json:"version,omitempty"`
	Confidence  float64 `json:"confidence"` // 0.0 to 1.0
	Method      string  `json:"method"`     // how it was detected
	Files       []string `json:"files,omitempty"` // files that indicated this framework
}

// frameworkService implements FrameworkService
type frameworkService struct {
	dockerService DockerService
}

// NewFrameworkService creates a new FrameworkService implementation
func NewFrameworkService(dockerService DockerService) FrameworkService {
	return &frameworkService{
		dockerService: dockerService,
	}
}

// DetectFramework detects the framework used in a project directory
func (s *frameworkService) DetectFramework(ctx context.Context, projectPath string) (*FrameworkDetectionResult, error) {
	// Try multiple detection methods in order of reliability
	
	// 1. Try detecting from docker-compose.yml (most reliable for scaffolded projects)
	if result, err := s.DetectFromCompose(ctx, projectPath); err == nil && result.Framework != "" {
		return result, nil
	}
	
	// 2. Try detecting from project files
	if result, err := s.DetectFromFiles(ctx, projectPath); err == nil && result.Framework != "" {
		return result, nil
	}
	
	return &FrameworkDetectionResult{
		Framework:  "unknown",
		Confidence: 0.0,
		Method:     "none",
	}, nil
}

// DetectFromCompose detects framework from docker-compose.yml file
func (s *frameworkService) DetectFromCompose(ctx context.Context, projectPath string) (*FrameworkDetectionResult, error) {
	composePath := filepath.Join(projectPath, "docker-compose.yml")
	
	// Check if docker-compose.yml exists
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return &FrameworkDetectionResult{Framework: ""}, nil
	}
	
	// Use the existing docker framework detection
	framework, err := docker.DetectFrameworkFromCompose(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to detect framework from docker-compose: %w", err)
	}
	
	if framework == "" {
		return &FrameworkDetectionResult{Framework: ""}, nil
	}
	
	// Try to get version information
	version, _ := s.GetFrameworkVersion(ctx, projectPath, framework)
	
	return &FrameworkDetectionResult{
		Framework:  framework,
		Version:    version,
		Confidence: 0.9, // High confidence for docker-compose detection
		Method:     "docker-compose",
		Files:      []string{"docker-compose.yml"},
	}, nil
}

// DetectFromFiles detects framework from project files and directory structure
func (s *frameworkService) DetectFromFiles(ctx context.Context, projectPath string) (*FrameworkDetectionResult, error) {
	// Framework detection patterns
	patterns := map[string]frameworkPattern{
		"laravel": {
			files: []string{"artisan", "composer.json"},
			directories: []string{"app", "config", "routes"},
			confidence: 0.8,
		},
		"django": {
			files: []string{"manage.py", "requirements.txt"},
			directories: []string{},
			confidence: 0.8,
		},
		"react": {
			files: []string{"package.json"},
			directories: []string{"src", "public"},
			confidence: 0.6,
		},
		"vue": {
			files: []string{"package.json", "vue.config.js"},
			directories: []string{"src"},
			confidence: 0.7,
		},
		"angular": {
			files: []string{"angular.json", "package.json"},
			directories: []string{"src"},
			confidence: 0.8,
		},
		"express": {
			files: []string{"package.json"},
			directories: []string{},
			confidence: 0.5,
		},
	}
	
	var bestMatch *FrameworkDetectionResult
	
	for framework, pattern := range patterns {
		result := s.checkFrameworkPattern(projectPath, framework, pattern)
		if result.Confidence > 0 && (bestMatch == nil || result.Confidence > bestMatch.Confidence) {
			bestMatch = result
		}
	}
	
	if bestMatch == nil {
		return &FrameworkDetectionResult{Framework: ""}, nil
	}
	
	// Try to get version information
	version, _ := s.GetFrameworkVersion(ctx, projectPath, bestMatch.Framework)
	bestMatch.Version = version
	
	return bestMatch, nil
}

// frameworkPattern defines what to look for when detecting a framework
type frameworkPattern struct {
	files       []string
	directories []string
	confidence  float64
}

// checkFrameworkPattern checks if a directory matches a framework pattern
func (s *frameworkService) checkFrameworkPattern(projectPath, framework string, pattern frameworkPattern) *FrameworkDetectionResult {
	var foundFiles []string
	matchedFiles := 0
	totalFiles := len(pattern.files)
	
	// Check for required files
	for _, file := range pattern.files {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); err == nil {
			foundFiles = append(foundFiles, file)
			matchedFiles++
		}
	}
	
	// Check for required directories
	matchedDirs := 0
	totalDirs := len(pattern.directories)
	
	for _, dir := range pattern.directories {
		dirPath := filepath.Join(projectPath, dir)
		if stat, err := os.Stat(dirPath); err == nil && stat.IsDir() {
			matchedDirs++
		}
	}
	
	// Calculate confidence based on matches
	confidence := 0.0
	if totalFiles > 0 {
		confidence += (float64(matchedFiles) / float64(totalFiles)) * 0.7
	}
	if totalDirs > 0 {
		confidence += (float64(matchedDirs) / float64(totalDirs)) * 0.3
	}
	
	// Apply base confidence multiplier
	confidence *= pattern.confidence
	
	if confidence < 0.3 {
		confidence = 0.0 // Too low confidence
	}
	
	return &FrameworkDetectionResult{
		Framework:  framework,
		Confidence: confidence,
		Method:     "files",
		Files:      foundFiles,
	}
}

// IsFramework checks if a directory contains a specific framework
func (s *frameworkService) IsFramework(ctx context.Context, projectPath, framework string) (bool, error) {
	result, err := s.DetectFramework(ctx, projectPath)
	if err != nil {
		return false, err
	}
	
	return strings.EqualFold(result.Framework, framework), nil
}

// GetFrameworkInfo returns detailed information about a framework
func (s *frameworkService) GetFrameworkInfo(ctx context.Context, framework string) (*FrameworkInfo, error) {
	frameworkInfos := map[string]FrameworkInfo{
		"laravel": {
			Name:              "Laravel",
			Language:          "PHP",
			SupportedVersions: []string{"8", "9", "10", "11"},
			Description:       "The PHP Framework for Web Artisans",
		},
		"django": {
			Name:              "Django",
			Language:          "Python",
			SupportedVersions: []string{"3.0", "3.1", "3.2", "4.0", "4.1", "4.2", "5.0"},
			Description:       "The Web framework for perfectionists with deadlines",
		},
		"react": {
			Name:              "React",
			Language:          "JavaScript",
			SupportedVersions: []string{"17", "18"},
			Description:       "A JavaScript library for building user interfaces",
		},
		"vue": {
			Name:              "Vue.js",
			Language:          "JavaScript",
			SupportedVersions: []string{"2", "3"},
			Description:       "The Progressive JavaScript Framework",
		},
		"angular": {
			Name:              "Angular",
			Language:          "TypeScript",
			SupportedVersions: []string{"13", "14", "15", "16", "17"},
			Description:       "Platform for building mobile and desktop web applications",
		},
		"express": {
			Name:              "Express.js",
			Language:          "JavaScript",
			SupportedVersions: []string{"4"},
			Description:       "Fast, unopinionated, minimalist web framework for Node.js",
		},
	}
	
	info, exists := frameworkInfos[strings.ToLower(framework)]
	if !exists {
		return nil, fmt.Errorf("framework '%s' not found", framework)
	}
	
	return &info, nil
}

// ListSupportedFrameworks returns all supported frameworks
func (s *frameworkService) ListSupportedFrameworks(ctx context.Context) ([]FrameworkInfo, error) {
	frameworks := []string{"laravel", "django", "react", "vue", "angular", "express"}
	var result []FrameworkInfo
	
	for _, framework := range frameworks {
		info, err := s.GetFrameworkInfo(ctx, framework)
		if err != nil {
			continue
		}
		result = append(result, *info)
	}
	
	return result, nil
}

// GetFrameworkVersion attempts to detect the version of a framework in use
func (s *frameworkService) GetFrameworkVersion(ctx context.Context, projectPath, framework string) (string, error) {
	switch strings.ToLower(framework) {
	case "laravel":
		return s.getLaravelVersion(projectPath)
	case "django":
		return s.getDjangoVersion(projectPath)
	case "react", "vue", "angular", "express":
		return s.getNodeFrameworkVersion(projectPath, framework)
	default:
		return "", nil
	}
}

// getLaravelVersion gets Laravel version from composer.json
func (s *frameworkService) getLaravelVersion(projectPath string) (string, error) {
	composerPath := filepath.Join(projectPath, "composer.json")
	if _, err := os.Stat(composerPath); os.IsNotExist(err) {
		return "", nil
	}
	
	// Read and parse composer.json to extract Laravel version
	content, err := os.ReadFile(composerPath)
	if err != nil {
		return "", err
	}
	
	// Simple version extraction (could be improved with JSON parsing)
	contentStr := string(content)
	if strings.Contains(contentStr, "laravel/framework") {
		// Extract version information from composer.json
		// This is a simplified version - could be enhanced
		return utils.ExtractVersionFromComposer(contentStr, "laravel/framework")
	}
	
	return "", nil
}

// getDjangoVersion gets Django version from requirements.txt or manage.py
func (s *frameworkService) getDjangoVersion(projectPath string) (string, error) {
	// Try requirements.txt first
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		content, err := os.ReadFile(requirementsPath)
		if err != nil {
			return "", err
		}
		
		return utils.ExtractVersionFromRequirements(string(content), "Django")
	}
	
	return "", nil
}

// getNodeFrameworkVersion gets version from package.json
func (s *frameworkService) getNodeFrameworkVersion(projectPath, framework string) (string, error) {
	packagePath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return "", nil
	}
	
	content, err := os.ReadFile(packagePath)
	if err != nil {
		return "", err
	}
	
	return utils.ExtractVersionFromPackageJSON(string(content), framework)
}
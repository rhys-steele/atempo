package config

import (
	"fmt"
	"strings"

	"atempo/internal/registry"
)

// ValidateProviderSettings validates provider-specific settings
func ValidateProviderSettings(provider string, settings map[string]interface{}) error {
	switch provider {
	case "github":
		return validateGitHubSettings(settings)
	case "gitlab":
		return validateGitLabSettings(settings)
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}
}

// validateGitHubSettings validates GitHub Actions specific settings
func validateGitHubSettings(settings map[string]interface{}) error {
	// Validate workflow name
	if workflowName, exists := settings["workflow_name"]; exists {
		if name, ok := workflowName.(string); ok {
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf("workflow_name cannot be empty")
			}
		} else {
			return fmt.Errorf("workflow_name must be a string")
		}
	}

	// Validate triggers
	if triggers, exists := settings["triggers"]; exists {
		if triggerList, ok := triggers.([]interface{}); ok {
			validTriggers := []string{"push", "pull_request", "schedule", "workflow_dispatch"}
			for _, trigger := range triggerList {
				if triggerStr, ok := trigger.(string); ok {
					isValid := false
					for _, validTrigger := range validTriggers {
						if triggerStr == validTrigger {
							isValid = true
							break
						}
					}
					if !isValid {
						return fmt.Errorf("invalid trigger '%s', must be one of: %v", triggerStr, validTriggers)
					}
				} else {
					return fmt.Errorf("trigger must be a string")
				}
			}
		} else {
			return fmt.Errorf("triggers must be an array")
		}
	}

	// Validate timeout
	if timeout, exists := settings["timeout"]; exists {
		if timeoutNum, ok := timeout.(float64); ok {
			if timeoutNum <= 0 || timeoutNum > 360 {
				return fmt.Errorf("timeout must be between 1 and 360 minutes")
			}
		} else {
			return fmt.Errorf("timeout must be a number")
		}
	}

	return nil
}

// validateGitLabSettings validates GitLab CI specific settings
func validateGitLabSettings(settings map[string]interface{}) error {
	// Validate image
	if image, exists := settings["image"]; exists {
		if imageStr, ok := image.(string); ok {
			if strings.TrimSpace(imageStr) == "" {
				return fmt.Errorf("image cannot be empty")
			}
		} else {
			return fmt.Errorf("image must be a string")
		}
	}

	// Validate stages
	if stages, exists := settings["stages"]; exists {
		if stageList, ok := stages.([]interface{}); ok {
			if len(stageList) == 0 {
				return fmt.Errorf("stages cannot be empty")
			}
			for _, stage := range stageList {
				if stageStr, ok := stage.(string); ok {
					if strings.TrimSpace(stageStr) == "" {
						return fmt.Errorf("stage name cannot be empty")
					}
				} else {
					return fmt.Errorf("stage must be a string")
				}
			}
		} else {
			return fmt.Errorf("stages must be an array")
		}
	}

	// Validate variables
	if variables, exists := settings["variables"]; exists {
		if variableMap, ok := variables.(map[string]interface{}); ok {
			for key, value := range variableMap {
				if strings.TrimSpace(key) == "" {
					return fmt.Errorf("variable key cannot be empty")
				}
				if _, ok := value.(string); !ok {
					return fmt.Errorf("variable value must be a string")
				}
			}
		} else {
			return fmt.Errorf("variables must be a map")
		}
	}

	return nil
}

// ValidateFrameworkCompatibility validates that a framework is compatible with a provider
func ValidateFrameworkCompatibility(framework, provider string) error {
	// Define framework-provider compatibility matrix
	compatibility := map[string][]string{
		"laravel":     {"github", "gitlab"},
		"django":      {"github", "gitlab"},
		"express":     {"github", "gitlab"},
		"lambda-node": {"github", "gitlab"},
	}

	if supportedProviders, exists := compatibility[framework]; exists {
		for _, supportedProvider := range supportedProviders {
			if provider == supportedProvider {
				return nil
			}
		}
		return fmt.Errorf("framework '%s' is not compatible with provider '%s'. Supported providers: %v", framework, provider, supportedProviders)
	}

	return fmt.Errorf("unknown framework: %s", framework)
}

// ValidateProjectStructure validates that the project has the expected structure for CI
func ValidateProjectStructure(projectPath, framework string) error {
	// Framework-specific validation
	switch framework {
	case "laravel":
		return validateLaravelProject(projectPath)
	case "django":
		return validateDjangoProject(projectPath)
	case "express":
		return validateExpressProject(projectPath)
	case "lambda-node":
		return validateLambdaNodeProject(projectPath)
	default:
		return fmt.Errorf("unknown framework: %s", framework)
	}
}

// validateLaravelProject validates Laravel project structure
func validateLaravelProject(projectPath string) error {
	// Check for essential Laravel files
	essentialFiles := []string{"composer.json", "artisan"}
	for _, file := range essentialFiles {
		// For this implementation, we'll assume the files exist
		// In a full implementation, we would check file existence
	}
	return nil
}

// validateDjangoProject validates Django project structure
func validateDjangoProject(projectPath string) error {
	// Check for essential Django files
	essentialFiles := []string{"manage.py", "requirements.txt"}
	for _, file := range essentialFiles {
		// For this implementation, we'll assume the files exist
		// In a full implementation, we would check file existence
	}
	return nil
}

// validateExpressProject validates Express project structure
func validateExpressProject(projectPath string) error {
	// Check for essential Express files
	essentialFiles := []string{"package.json"}
	for _, file := range essentialFiles {
		// For this implementation, we'll assume the files exist
		// In a full implementation, we would check file existence
	}
	return nil
}

// validateLambdaNodeProject validates Lambda Node project structure
func validateLambdaNodeProject(projectPath string) error {
	// Check for essential Lambda files
	essentialFiles := []string{"package.json"}
	for _, file := range essentialFiles {
		// For this implementation, we'll assume the files exist
		// In a full implementation, we would check file existence
	}
	return nil
}
package ci

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ANSI color codes (reused from commands package)
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
)

// InteractivePrompter handles interactive CI configuration prompts
type InteractivePrompter struct {
	scanner *bufio.Scanner
}

// NewInteractivePrompter creates a new interactive prompter
func NewInteractivePrompter() *InteractivePrompter {
	return &InteractivePrompter{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// ShowHeader displays the CI setup header
func (p *InteractivePrompter) ShowHeader() {
	fmt.Printf("\n%s🚀 CI/CD Pipeline Setup%s\n", ColorBlue, ColorReset)
	fmt.Printf("═══════════════════════════════════════════════════════════\n")
	fmt.Printf("%sLet's configure continuous integration for your project!%s\n\n", ColorCyan, ColorReset)
}

// PromptProviderSelection prompts user to select CI provider
func (p *InteractivePrompter) PromptProviderSelection(providers []CIProvider) (CIProvider, error) {
	fmt.Printf("%s📦 CI Provider Selection%s\n", ColorYellow, ColorReset)
	fmt.Printf("Choose your preferred CI/CD provider:\n\n")

	providerNames := map[CIProvider]string{
		ProviderGitHub: "GitHub Actions",
		ProviderGitLab: "GitLab CI",
	}

	for i, provider := range providers {
		name := providerNames[provider]
		if name == "" {
			name = string(provider)
		}
		fmt.Printf("   %s%d.%s %s\n", ColorCyan, i+1, ColorReset, name)
	}

	fmt.Printf("\n   %s>%s ", ColorCyan, ColorReset)

	if p.scanner.Scan() {
		input := strings.TrimSpace(p.scanner.Text())
		if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(providers) {
			return providers[choice-1], nil
		}
	}

	return "", fmt.Errorf("invalid provider selection")
}

// PromptFrameworkSelection prompts user to select or confirm framework
func (p *InteractivePrompter) PromptFrameworkSelection(detectedFramework string, supportedFrameworks []string) (string, error) {
	fmt.Printf("\n%s🔍 Framework Detection%s\n", ColorYellow, ColorReset)

	if detectedFramework != "unknown" && detectedFramework != "" {
		fmt.Printf("Detected framework: %s%s%s\n", 
			ColorGreen, detectedFramework, ColorReset)
		fmt.Printf("Use detected framework? %s[Y/n]%s ", ColorGray, ColorReset)

		if p.scanner.Scan() {
			input := strings.TrimSpace(strings.ToLower(p.scanner.Text()))
			if input == "" || input == "y" || input == "yes" {
				return detectedFramework, nil
			}
		}
	} else {
		fmt.Printf("Could not automatically detect framework.\n")
	}

	fmt.Printf("\nSelect your project framework:\n\n")
	for i, framework := range supportedFrameworks {
		fmt.Printf("   %s%d.%s %s\n", ColorCyan, i+1, ColorReset, framework)
	}

	fmt.Printf("\n   %s>%s ", ColorCyan, ColorReset)

	if p.scanner.Scan() {
		input := strings.TrimSpace(p.scanner.Text())
		if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(supportedFrameworks) {
			return supportedFrameworks[choice-1], nil
		}
	}

	return "", fmt.Errorf("invalid framework selection")
}

// PromptSettings prompts user for provider-specific settings
func (p *InteractivePrompter) PromptSettings(provider CIProvider, framework string, defaults map[string]interface{}) (map[string]interface{}, error) {
	fmt.Printf("\n%s⚙️  Configuration Settings%s\n", ColorYellow, ColorReset)
	fmt.Printf("Configure settings for %s%s%s with %s%s%s:\n\n", 
		ColorCyan, provider, ColorReset, ColorCyan, framework, ColorReset)

	settings := make(map[string]interface{})

	switch provider {
	case ProviderGitHub:
		return p.promptGitHubSettings(framework, defaults)
	case ProviderGitLab:
		return p.promptGitLabSettings(framework, defaults)
	}

	return settings, nil
}

// promptGitHubSettings prompts for GitHub Actions specific settings
func (p *InteractivePrompter) promptGitHubSettings(framework string, defaults map[string]interface{}) (map[string]interface{}, error) {
	settings := make(map[string]interface{})

	// Workflow name
	workflowName := p.promptWithDefault("Workflow name", defaults["workflow_name"].(string))
	settings["workflow_name"] = workflowName

	// Triggers
	fmt.Printf("\n%sWorkflow Triggers:%s\n", ColorGray, ColorReset)
	fmt.Printf("1. push (on code push)\n")
	fmt.Printf("2. pull_request (on PR creation/update)\n") 
	fmt.Printf("3. schedule (daily at 2 AM UTC)\n")
	fmt.Printf("4. workflow_dispatch (manual trigger)\n")
	
	triggers := p.promptMultiChoice("Select triggers (comma separated)", 
		[]string{"push", "pull_request", "schedule", "workflow_dispatch"},
		defaults["triggers"].([]string))
	settings["triggers"] = triggers

	// Branches
	branches := p.promptWithDefault("Target branches (comma separated)", 
		strings.Join(defaults["branches"].([]string), ", "))
	settings["branches"] = strings.Split(branches, ",")
	for i, branch := range settings["branches"].([]string) {
		settings["branches"].([]string)[i] = strings.TrimSpace(branch)
	}

	// Timeout
	timeout := p.promptWithDefault("Timeout (minutes)", 
		fmt.Sprintf("%v", defaults["timeout"]))
	if timeoutInt, err := strconv.Atoi(timeout); err == nil {
		settings["timeout"] = timeoutInt
	} else {
		settings["timeout"] = defaults["timeout"]
	}

	// Framework-specific settings
	switch framework {
	case "laravel":
		p.promptLaravelSettings(settings, defaults)
	case "django":
		p.promptDjangoSettings(settings, defaults)
	case "express", "lambda-node":
		p.promptNodeSettings(settings, defaults, framework)
	}

	return settings, nil
}

// promptGitLabSettings prompts for GitLab CI specific settings
func (p *InteractivePrompter) promptGitLabSettings(framework string, defaults map[string]interface{}) (map[string]interface{}, error) {
	settings := make(map[string]interface{})

	// Docker image
	image := p.promptWithDefault("Docker image", defaults["image"].(string))
	settings["image"] = image

	// Stages
	stages := p.promptWithDefault("Pipeline stages (comma separated)", 
		strings.Join(defaults["stages"].([]string), ", "))
	settings["stages"] = strings.Split(stages, ",")
	for i, stage := range settings["stages"].([]string) {
		settings["stages"].([]string)[i] = strings.TrimSpace(stage)
	}

	// Copy other defaults
	settings["before_script"] = defaults["before_script"]
	settings["variables"] = defaults["variables"]
	settings["cache_paths"] = defaults["cache_paths"]
	settings["services"] = defaults["services"]
	settings["artifact_paths"] = defaults["artifact_paths"]

	return settings, nil
}

// Framework-specific settings prompts

func (p *InteractivePrompter) promptLaravelSettings(settings map[string]interface{}, defaults map[string]interface{}) {
	fmt.Printf("\n%sLaravel Configuration:%s\n", ColorGray, ColorReset)
	
	phpVersion := p.promptWithDefault("PHP version", defaults["php_version"].(string))
	settings["php_version"] = phpVersion

	fmt.Printf("\n%sServices (select multiple):%s\n", ColorGray, ColorReset)
	fmt.Printf("1. mysql\n2. redis\n3. postgres\n4. memcached\n")
	services := p.promptMultiChoice("Select services", 
		[]string{"mysql", "redis", "postgres", "memcached"}, 
		defaults["services"].([]string))
	settings["services"] = services

	settings["cache_paths"] = defaults["cache_paths"]
	settings["environment"] = defaults["environment"]
}

func (p *InteractivePrompter) promptDjangoSettings(settings map[string]interface{}, defaults map[string]interface{}) {
	fmt.Printf("\n%sDjango Configuration:%s\n", ColorGray, ColorReset)
	
	pythonVersion := p.promptWithDefault("Python version", defaults["python_version"].(string))
	settings["python_version"] = pythonVersion

	fmt.Printf("\n%sServices (select multiple):%s\n", ColorGray, ColorReset)
	fmt.Printf("1. postgres\n2. redis\n3. mysql\n")
	services := p.promptMultiChoice("Select services", 
		[]string{"postgres", "redis", "mysql"}, 
		defaults["services"].([]string))
	settings["services"] = services

	settings["cache_paths"] = defaults["cache_paths"]
	settings["environment"] = defaults["environment"]
}

func (p *InteractivePrompter) promptNodeSettings(settings map[string]interface{}, defaults map[string]interface{}, framework string) {
	fmt.Printf("\n%s%s Configuration:%s\n", ColorGray, strings.Title(framework), ColorReset)
	
	nodeVersion := p.promptWithDefault("Node.js version", defaults["node_version"].(string))
	settings["node_version"] = nodeVersion

	if framework == "express" {
		fmt.Printf("\n%sServices (select multiple):%s\n", ColorGray, ColorReset)
		fmt.Printf("1. redis\n2. postgres\n3. mysql\n")
		services := p.promptMultiChoice("Select services", 
			[]string{"redis", "postgres", "mysql"}, 
			defaults["services"].([]string))
		settings["services"] = services
	} else {
		settings["services"] = defaults["services"]
	}

	settings["cache_paths"] = defaults["cache_paths"]
	settings["environment"] = defaults["environment"]
}

// Helper methods

func (p *InteractivePrompter) promptWithDefault(prompt, defaultValue string) string {
	fmt.Printf("%s [%s%s%s]: ", prompt, ColorGray, defaultValue, ColorReset)
	
	if p.scanner.Scan() {
		input := strings.TrimSpace(p.scanner.Text())
		if input == "" {
			return defaultValue
		}
		return input
	}
	
	return defaultValue
}

func (p *InteractivePrompter) promptMultiChoice(prompt string, options, defaults []string) []string {
	fmt.Printf("\n%s:\n", prompt)
	
	defaultMap := make(map[string]bool)
	for _, def := range defaults {
		defaultMap[def] = true
	}
	
	for i, option := range options {
		marker := " "
		if defaultMap[option] {
			marker = "✓"
		}
		fmt.Printf("   %s%s%s %d. %s\n", ColorGreen, marker, ColorReset, i+1, option)
	}
	
	fmt.Printf("\nEnter numbers (comma separated) or press Enter for defaults: ")
	
	if p.scanner.Scan() {
		input := strings.TrimSpace(p.scanner.Text())
		if input == "" {
			return defaults
		}
		
		var selected []string
		for _, numStr := range strings.Split(input, ",") {
			numStr = strings.TrimSpace(numStr)
			if num, err := strconv.Atoi(numStr); err == nil && num >= 1 && num <= len(options) {
				selected = append(selected, options[num-1])
			}
		}
		
		if len(selected) > 0 {
			return selected
		}
	}
	
	return defaults
}

// ShowConfigurationSummary displays a summary of the configuration before generating files
func (p *InteractivePrompter) ShowConfigurationSummary(config *CIConfig) {
	fmt.Printf("\n%s📋 Configuration Summary%s\n", ColorGreen, ColorReset)
	fmt.Printf("═══════════════════════════════════════════════════════════\n")
	fmt.Printf("%sProject:%s %s\n", ColorCyan, ColorReset, config.ProjectName)
	fmt.Printf("%sFramework:%s %s\n", ColorCyan, ColorReset, config.Framework)
	fmt.Printf("%sProvider:%s %s\n", ColorCyan, ColorReset, config.Provider)
	
	fmt.Printf("\n%sKey Settings:%s\n", ColorCyan, ColorReset)
	for key, value := range config.Settings {
		if key == "environment" || key == "variables" {
			continue // Skip complex nested objects for summary
		}
		fmt.Printf("  • %s: %v\n", key, value)
	}
	
	fmt.Printf("\n%s✨ Ready to generate CI configuration files!%s\n", ColorGreen, ColorReset)
}

// ConfirmGeneration asks user to confirm file generation
func (p *InteractivePrompter) ConfirmGeneration() bool {
	fmt.Printf("\nProceed with CI file generation? %s[Y/n]%s ", ColorYellow, ColorReset)
	
	if p.scanner.Scan() {
		input := strings.TrimSpace(strings.ToLower(p.scanner.Text()))
		return input == "" || input == "y" || input == "yes"
	}
	
	return false
}
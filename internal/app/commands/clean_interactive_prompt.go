package commands

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// CleanInteractivePrompter handles user input using external JSON templates
type CleanInteractivePrompter struct {
	scanner *bufio.Scanner
	prompts *InteractivePrompts
	loader  *TemplateLoader
}

// NewCleanInteractivePrompter creates a new clean interactive prompter
func NewCleanInteractivePrompter(templatesFS fs.FS) (*CleanInteractivePrompter, error) {
	loader := NewTemplateLoader(templatesFS)

	prompts, err := loader.LoadInteractivePrompts()
	if err != nil {
		// Fallback to basic prompts if loading fails
		prompts = &InteractivePrompts{
			UIElements: UIElements{
				Header:    "🚀 AI-Powered Project Setup",
				Separator: "═══════════════════════════════════════════════════════════",
				Subtitle:  "Let's create an intelligent, AI-first development experience!",
				AuthRequired: AuthRequiredElement{
					Title:    "🔐 Authentication Required for AI Features",
					Message:  "AI-powered project manifests require authentication.",
					Action:   "Run 'atempo auth' to connect your AI provider",
					Fallback: "Proceeding with basic project setup...",
				},
			},
		}
	}

	return &CleanInteractivePrompter{
		scanner: bufio.NewScanner(os.Stdin),
		prompts: prompts,
		loader:  loader,
	}, nil
}

// GatherProjectIntent collects project information using JSON-defined prompts
func (p *CleanInteractivePrompter) GatherProjectIntent(framework, projectName string, manifestGenerator *CleanAIManifestGenerator) (*ProjectIntent, error) {
	p.showHeader()

	// Get project description
	description := p.promptProjectDescription()

	// Get additional features
	additionalFeatures := p.promptAdditionalFeatures()

	// Get complexity
	complexity := p.promptComplexity()

	// Generate AI-powered intent using the clean generator
	fmt.Printf("\n%s🤖 Generating AI project manifest...%s\n", ColorBlue, ColorReset)

	intent := manifestGenerator.CreateProjectIntent(description, projectName, "", additionalFeatures, complexity)

	// Show preview
	p.showIntentPreview(intent)

	return intent, nil
}

// showHeader displays the styled header using JSON config
func (p *CleanInteractivePrompter) showHeader() {
	fmt.Printf("\n%s%s%s\n", ColorBlue, p.prompts.UIElements.Header, ColorReset)
	fmt.Printf("%s\n", p.prompts.UIElements.Separator)
	fmt.Printf("%s%s%s\n\n", ColorCyan, p.prompts.UIElements.Subtitle, ColorReset)
}

// promptProjectDescription gets the main project description using JSON config
func (p *CleanInteractivePrompter) promptProjectDescription() string {
	prompt := p.prompts.Prompts.ProjectDescription

	fmt.Printf("%s❓ %s%s\n", ColorYellow, prompt.Question, ColorReset)
	fmt.Printf("%s   %s%s\n", ColorGray, prompt.Subtitle, ColorReset)
	fmt.Println()
	fmt.Printf("%s   Examples:%s\n", ColorGray, ColorReset)

	for _, example := range prompt.Examples {
		fmt.Printf("%s   • \"%s\"%s\n", ColorGray, example, ColorReset)
	}

	fmt.Println()
	fmt.Printf("   %s>%s ", ColorCyan, ColorReset)

	if p.scanner.Scan() {
		description := strings.TrimSpace(p.scanner.Text())
		if description == "" {
			return "A modern web application"
		}
		return description
	}

	return "A modern web application"
}

// promptAdditionalFeatures asks about additional features using JSON config
func (p *CleanInteractivePrompter) promptAdditionalFeatures() []string {
	prompt := p.prompts.Prompts.AdditionalFeatures

	fmt.Printf("\n%s⚡ %s%s\n", ColorYellow, prompt.Question, ColorReset)
	fmt.Printf("   %s%s%s\n\n", ColorGray, prompt.Subtitle, ColorReset)

	for i, option := range prompt.Options {
		fmt.Printf("   %s%d.%s %s\n", ColorCyan, i+1, ColorReset, option)
	}

	fmt.Printf("\n   %s>%s ", ColorCyan, ColorReset)

	selected := []string{}
	if p.scanner.Scan() {
		input := strings.TrimSpace(p.scanner.Text())
		if input != "none" && input != "" {
			// Parse selected numbers
			for _, numStr := range strings.Split(input, ",") {
				numStr = strings.TrimSpace(numStr)
				if num, err := strconv.Atoi(numStr); err == nil && num >= 1 && num <= len(prompt.Options) {
					selected = append(selected, prompt.Options[num-1])
				}
			}
		}
	}

	return selected
}

// promptComplexity asks about project complexity using JSON config
func (p *CleanInteractivePrompter) promptComplexity() string {
	prompt := p.prompts.Prompts.Complexity

	fmt.Printf("\n%s📊 %s%s\n", ColorBlue, prompt.Question, ColorReset)
	fmt.Println()

	for _, option := range prompt.Options {
		fmt.Printf("   %s%s.%s %s - %s\n", ColorCyan, option.Key, ColorReset, option.Label, option.Description)
	}

	fmt.Printf("\n   %s>%s ", ColorCyan, ColorReset)

	if p.scanner.Scan() {
		choice := strings.TrimSpace(p.scanner.Text())
		for _, option := range prompt.Options {
			if choice == option.Key {
				return option.Label
			}
		}
	}

	return "Medium"
}

// showIntentPreview displays a preview of the generated project intent
func (p *CleanInteractivePrompter) showIntentPreview(intent *ProjectIntent) {
	fmt.Printf("\n%s📋 Project Intent Summary%s\n", ColorGreen, ColorReset)
	fmt.Printf("═══════════════════════════════════════════════════════════\n")
	fmt.Printf("%sDescription:%s %s\n", ColorCyan, ColorReset, intent.Description)
	fmt.Printf("%sFramework:%s %s (%s)\n", ColorCyan, ColorReset, intent.Framework, intent.Language)
	fmt.Printf("%sType:%s %s\n", ColorCyan, ColorReset, intent.ProjectType)

	fmt.Printf("\n%sCore Features:%s\n", ColorCyan, ColorReset)
	for _, feature := range intent.CoreFeatures {
		fmt.Printf("  • %s\n", feature)
	}

	fmt.Printf("\n%sTechnical Needs:%s\n", ColorCyan, ColorReset)
	for _, need := range intent.TechnicalNeeds {
		fmt.Printf("  • %s\n", need)
	}

	fmt.Printf("\n%sUser Stories Generated:%s %d\n", ColorCyan, ColorReset, len(intent.UserStories))
	fmt.Printf("%sArchitecture Hints:%s %d guidelines\n", ColorCyan, ColorReset, len(intent.ArchitectureHints))

	fmt.Printf("\n%s✨ This manifest will help AI tools understand your project better!%s\n", ColorGreen, ColorReset)
}

// ShowAuthenticationPrompt displays auth required message using JSON config
func (p *CleanInteractivePrompter) ShowAuthenticationPrompt() {
	auth := p.prompts.UIElements.AuthRequired

	fmt.Printf("\n%s%s%s\n", ColorYellow, auth.Title, ColorReset)
	fmt.Printf("   %s\n", auth.Message)
	fmt.Printf("   %s%s%s.\n\n", ColorCyan, auth.Action, ColorReset)
	fmt.Printf("   %s%s%s\n", ColorGray, auth.Fallback, ColorReset)
}

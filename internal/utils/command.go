package utils

import (
	"fmt"
	"strings"
)

// CommandFactory provides utilities for creating and managing CLI commands
type CommandFactory struct{}

// NewCommandFactory creates a new command factory
func NewCommandFactory() *CommandFactory {
	return &CommandFactory{}
}

// FormatUsage formats a usage string with consistent styling
func (cf *CommandFactory) FormatUsage(commandName string, args ...string) string {
	if len(args) == 0 {
		return fmt.Sprintf("atempo %s", commandName)
	}
	
	var formattedArgs []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") {
			// Optional argument
			formattedArgs = append(formattedArgs, arg)
		} else if strings.HasPrefix(arg, "<") && strings.HasSuffix(arg, ">") {
			// Required argument
			formattedArgs = append(formattedArgs, arg)
		} else {
			// Literal argument
			formattedArgs = append(formattedArgs, arg)
		}
	}
	
	return fmt.Sprintf("atempo %s %s", commandName, strings.Join(formattedArgs, " "))
}

// FormatUsageWithExamples formats usage with examples
func (cf *CommandFactory) FormatUsageWithExamples(commandName string, args []string, examples []string) string {
	usage := cf.FormatUsage(commandName, args...)
	
	if len(examples) == 0 {
		return usage
	}
	
	var exampleStrings []string
	for _, example := range examples {
		exampleStrings = append(exampleStrings, fmt.Sprintf("  %s", example))
	}
	
	return fmt.Sprintf("%s\nExamples:\n%s", usage, strings.Join(exampleStrings, "\n"))
}

// FormatDescription formats a command description with consistent styling
func (cf *CommandFactory) FormatDescription(description string) string {
	// Ensure description doesn't end with a period (for consistency)
	description = strings.TrimSuffix(description, ".")
	return description
}

// ValidateCommandName validates a command name follows conventions
func (cf *CommandFactory) ValidateCommandName(name string) error {
	if name == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	
	if strings.Contains(name, " ") {
		return fmt.Errorf("command name cannot contain spaces")
	}
	
	if strings.Contains(name, "-") {
		// Allow hyphens but validate they're not at start/end
		if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
			return fmt.Errorf("command name cannot start or end with hyphen")
		}
	}
	
	// Check for valid characters (alphanumeric and hyphens only)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("command name contains invalid character: %c", char)
		}
	}
	
	return nil
}

// CreateStandardUsage creates a standard usage string for common command patterns
func (cf *CommandFactory) CreateStandardUsage(commandName string, pattern CommandPattern, args ...string) string {
	switch pattern {
	case PatternSimple:
		return cf.FormatUsage(commandName)
	case PatternWithOptionalArgs:
		return cf.FormatUsage(commandName, args...)
	case PatternWithRequiredArgs:
		return cf.FormatUsage(commandName, args...)
	case PatternWithProjectContext:
		return cf.FormatUsage(commandName, "[project]")
	case PatternWithFrameworkSpec:
		return cf.FormatUsage(commandName, "<framework>:<version>", "[project-name]")
	case PatternDockerPassthrough:
		return cf.FormatUsage(commandName, "<docker-command>", "[args...]")
	default:
		return cf.FormatUsage(commandName, args...)
	}
}

// CommandPattern represents common command patterns
type CommandPattern int

const (
	PatternSimple CommandPattern = iota
	PatternWithOptionalArgs
	PatternWithRequiredArgs
	PatternWithProjectContext
	PatternWithFrameworkSpec
	PatternDockerPassthrough
)

// CommandMetadata contains metadata for command creation
type CommandMetadata struct {
	Name        string
	Description string
	Pattern     CommandPattern
	Args        []string
	Examples    []string
}

// CreateUsageFromMetadata creates usage string from command metadata
func (cf *CommandFactory) CreateUsageFromMetadata(metadata CommandMetadata) string {
	var args []string
	if len(metadata.Args) > 0 {
		args = metadata.Args
	} else {
		// Use standard args based on pattern
		switch metadata.Pattern {
		case PatternWithProjectContext:
			args = []string{"[project]"}
		case PatternWithFrameworkSpec:
			args = []string{"<framework>:<version>", "[project-name]"}
		case PatternDockerPassthrough:
			args = []string{"<docker-command>", "[args...]"}
		}
	}
	
	if len(metadata.Examples) > 0 {
		return cf.FormatUsageWithExamples(metadata.Name, args, metadata.Examples)
	}
	
	return cf.CreateStandardUsage(metadata.Name, metadata.Pattern, args...)
}

// ValidateCommandMetadata validates command metadata
func (cf *CommandFactory) ValidateCommandMetadata(metadata CommandMetadata) error {
	if err := cf.ValidateCommandName(metadata.Name); err != nil {
		return fmt.Errorf("invalid command name: %w", err)
	}
	
	if metadata.Description == "" {
		return fmt.Errorf("command description cannot be empty")
	}
	
	return nil
}

// CommonCommandDescriptions provides standard descriptions for common commands
var CommonCommandDescriptions = map[string]string{
	"create":      "Create a new project with specified framework",
	"projects":    "Show all projects with their status",
	"docker":      "Execute Docker commands for current project",
	"shell":       "Enter interactive shell mode",
	"describe":    "Show detailed project information",
	"remove":      "Remove project from registry",
	"logs":        "View project logs",
	"stop":        "Stop project services",
	"reconfigure": "Reconfigure project settings",
	"add-service": "Add new service to project",
	"test":        "Run project tests",
	"audit":       "Run comprehensive codebase audit using Claude Code",
}

// GetStandardDescription returns a standard description for common commands
func (cf *CommandFactory) GetStandardDescription(commandName string) string {
	if desc, exists := CommonCommandDescriptions[commandName]; exists {
		return desc
	}
	return ""
}

// CreateDockerSubcommandUsage creates usage for docker subcommands
func (cf *CommandFactory) CreateDockerSubcommandUsage(subcommand string, args ...string) string {
	baseUsage := cf.FormatUsage("docker", subcommand)
	if len(args) > 0 {
		return fmt.Sprintf("%s %s", baseUsage, strings.Join(args, " "))
	}
	return baseUsage
}

// CreateProjectSpecificUsage creates usage for commands that work with specific projects
func (cf *CommandFactory) CreateProjectSpecificUsage(commandName string, projectName string, args ...string) string {
	usage := fmt.Sprintf("%s %s", projectName, commandName)
	if len(args) > 0 {
		usage = fmt.Sprintf("%s %s", usage, strings.Join(args, " "))
	}
	return usage
}

// ValidateArgumentPattern validates argument patterns in usage strings
func (cf *CommandFactory) ValidateArgumentPattern(pattern string) error {
	// Check for balanced brackets
	openBrackets := strings.Count(pattern, "[")
	closeBrackets := strings.Count(pattern, "]")
	if openBrackets != closeBrackets {
		return fmt.Errorf("unbalanced square brackets in pattern: %s", pattern)
	}
	
	// Check for balanced angle brackets
	openAngle := strings.Count(pattern, "<")
	closeAngle := strings.Count(pattern, ">")
	if openAngle != closeAngle {
		return fmt.Errorf("unbalanced angle brackets in pattern: %s", pattern)
	}
	
	return nil
}

// SuggestCommandName suggests a command name based on description
func (cf *CommandFactory) SuggestCommandName(description string) string {
	words := strings.Fields(strings.ToLower(description))
	if len(words) == 0 {
		return "command"
	}
	
	// Take first word and clean it
	name := words[0]
	name = strings.TrimSuffix(name, "s") // Remove plural
	
	// Handle common patterns
	switch name {
	case "show", "display":
		if len(words) > 1 {
			secondWord := strings.TrimSuffix(words[1], "s")
			return secondWord
		}
		return "show"
	case "list":
		if len(words) > 1 {
			secondWord := strings.TrimSuffix(words[1], "s")
			return secondWord
		}
		return "list"
	case "create", "make", "generate":
		return "create"
	case "remove", "delete", "destroy":
		return "remove"
	case "start", "run", "execute":
		return "start"
	case "stop", "halt", "terminate":
		return "stop"
	default:
		return name
	}
}

// Global command factory instance
var GlobalCommandFactory = NewCommandFactory()

// Helper functions for common command factory operations
func FormatUsage(commandName string, args ...string) string {
	return GlobalCommandFactory.FormatUsage(commandName, args...)
}

func FormatUsageWithExamples(commandName string, args []string, examples []string) string {
	return GlobalCommandFactory.FormatUsageWithExamples(commandName, args, examples)
}

func ValidateCommandName(name string) error {
	return GlobalCommandFactory.ValidateCommandName(name)
}

func CreateStandardUsage(commandName string, pattern CommandPattern, args ...string) string {
	return GlobalCommandFactory.CreateStandardUsage(commandName, pattern, args...)
}

func GetStandardDescription(commandName string) string {
	return GlobalCommandFactory.GetStandardDescription(commandName)
}

func CreateDockerSubcommandUsage(subcommand string, args ...string) string {
	return GlobalCommandFactory.CreateDockerSubcommandUsage(subcommand, args...)
}

func CreateProjectSpecificUsage(commandName string, projectName string, args ...string) string {
	return GlobalCommandFactory.CreateProjectSpecificUsage(commandName, projectName, args...)
}
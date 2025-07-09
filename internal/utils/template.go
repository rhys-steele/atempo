package utils

import (
	"regexp"
	"strings"
	"text/template"
)

// TemplateProcessor provides utilities for template processing and validation
type TemplateProcessor struct{}

// NewTemplateProcessor creates a new template processor
func NewTemplateProcessor() *TemplateProcessor {
	return &TemplateProcessor{}
}

// TemplateVariable represents a template variable
type TemplateVariable struct {
	Name    string
	Pattern string
	Value   string
}

// Common template variables
var CommonTemplateVariables = []TemplateVariable{
	{Name: "project", Pattern: "{{project}}", Value: ""},
	{Name: "cwd", Pattern: "{{cwd}}", Value: ""},
	{Name: "version", Pattern: "{{version}}", Value: ""},
	{Name: "name", Pattern: "{{name}}", Value: ""},
	{Name: "framework", Pattern: "{{framework}}", Value: ""},
	{Name: "language", Pattern: "{{language}}", Value: ""},
}

// HasTemplateVariables checks if content contains any template variables
func (tp *TemplateProcessor) HasTemplateVariables(content string) bool {
	// Check for mustache-style variables {{variable}}
	mustacheRegex := regexp.MustCompile(`\{\{[^}]+\}\}`)
	return mustacheRegex.MatchString(content)
}

// GetTemplateVariables extracts all template variables from content
func (tp *TemplateProcessor) GetTemplateVariables(content string) []string {
	// Extract mustache-style variables {{variable}}
	mustacheRegex := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := mustacheRegex.FindAllStringSubmatch(content, -1)
	
	variables := make([]string, 0) // Initialize empty slice instead of nil
	seen := make(map[string]bool)
	
	for _, match := range matches {
		if len(match) > 1 {
			variable := strings.TrimSpace(match[1])
			if !seen[variable] {
				variables = append(variables, variable)
				seen[variable] = true
			}
		}
	}
	
	return variables
}

// ValidateTemplateVariables validates that all template variables are known
func (tp *TemplateProcessor) ValidateTemplateVariables(content string) []string {
	variables := tp.GetTemplateVariables(content)
	unknownVariables := make([]string, 0) // Initialize empty slice instead of nil
	
	knownVariables := make(map[string]bool)
	for _, commonVar := range CommonTemplateVariables {
		knownVariables[commonVar.Name] = true
	}
	
	for _, variable := range variables {
		if !knownVariables[variable] {
			unknownVariables = append(unknownVariables, variable)
		}
	}
	
	return unknownVariables
}

// SubstituteVariables replaces template variables with their values
func (tp *TemplateProcessor) SubstituteVariables(content string, variables map[string]string) string {
	result := content
	
	for key, value := range variables {
		// Handle both {{key}} and {{key}} patterns
		patterns := []string{
			"{{" + key + "}}",
			"{{ " + key + " }}",
		}
		
		for _, pattern := range patterns {
			result = strings.ReplaceAll(result, pattern, value)
		}
	}
	
	return result
}

// ProcessTemplate processes a template with Go's text/template package
func (tp *TemplateProcessor) ProcessTemplate(templateContent string, data interface{}) (string, error) {
	// Create template with helper functions
	tmpl := template.New("template").Funcs(template.FuncMap{
		"title": strings.Title,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
		"replace": func(old, new, str string) string {
			return strings.ReplaceAll(str, old, new)
		},
		"contains": func(substr, str string) bool {
			return strings.Contains(str, substr)
		},
		"hasPrefix": func(prefix, str string) bool {
			return strings.HasPrefix(str, prefix)
		},
		"hasSuffix": func(suffix, str string) bool {
			return strings.HasSuffix(str, suffix)
		},
		"join": func(sep string, strs []string) string {
			return strings.Join(strs, sep)
		},
		"split": func(sep, str string) []string {
			return strings.Split(str, sep)
		},
	})
	
	var err error
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// IsValidTemplateVariable checks if a variable name is valid
func (tp *TemplateProcessor) IsValidTemplateVariable(name string) bool {
	// Variable names should be alphanumeric with underscores
	validNameRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	return validNameRegex.MatchString(name)
}

// GetVariablePattern returns the template pattern for a variable name
func (tp *TemplateProcessor) GetVariablePattern(name string) string {
	return "{{" + name + "}}"
}

// ExtractVariableFromPattern extracts the variable name from a pattern
func (tp *TemplateProcessor) ExtractVariableFromPattern(pattern string) string {
	// Remove {{ and }} from pattern
	if strings.HasPrefix(pattern, "{{") && strings.HasSuffix(pattern, "}}") {
		return strings.TrimSpace(pattern[2 : len(pattern)-2])
	}
	return ""
}

// CreateVariableMap creates a map of template variables with their values
func (tp *TemplateProcessor) CreateVariableMap(values map[string]string) map[string]string {
	result := make(map[string]string)
	
	// Add common variables with empty defaults
	for _, commonVar := range CommonTemplateVariables {
		result[commonVar.Name] = ""
	}
	
	// Override with provided values
	for key, value := range values {
		result[key] = value
	}
	
	return result
}

// ValidateTemplateContent validates template content for common issues
func (tp *TemplateProcessor) ValidateTemplateContent(content string) []string {
	issues := make([]string, 0) // Initialize empty slice instead of nil
	
	// Check for unbalanced braces
	openCount := strings.Count(content, "{{")
	closeCount := strings.Count(content, "}}")
	if openCount != closeCount {
		issues = append(issues, "unbalanced template braces")
	}
	
	// Check for nested braces
	nestedRegex := regexp.MustCompile(`\{\{[^}]*\{\{`)
	if nestedRegex.MatchString(content) {
		issues = append(issues, "nested template braces detected")
	}
	
	// Check for unclosed braces
	unclosedRegex := regexp.MustCompile(`\{\{[^}]*$`)
	if unclosedRegex.MatchString(content) {
		issues = append(issues, "unclosed template braces")
	}
	
	// Check for empty variables
	emptyVarRegex := regexp.MustCompile(`\{\{\s*\}\}`)
	if emptyVarRegex.MatchString(content) {
		issues = append(issues, "empty template variables")
	}
	
	// Check for invalid variable names
	variables := tp.GetTemplateVariables(content)
	for _, variable := range variables {
		if !tp.IsValidTemplateVariable(variable) {
			issues = append(issues, "invalid variable name: "+variable)
		}
	}
	
	return issues
}

// GetTemplateInfo returns information about a template
func (tp *TemplateProcessor) GetTemplateInfo(content string) TemplateInfo {
	variables := tp.GetTemplateVariables(content)
	unknownVariables := tp.ValidateTemplateVariables(content)
	issues := tp.ValidateTemplateContent(content)
	
	return TemplateInfo{
		HasVariables:     len(variables) > 0,
		Variables:        variables,
		UnknownVariables: unknownVariables,
		Issues:           issues,
		IsValid:          len(issues) == 0,
	}
}

// TemplateInfo contains information about a template
type TemplateInfo struct {
	HasVariables     bool
	Variables        []string
	UnknownVariables []string
	Issues           []string
	IsValid          bool
}

// Global template processor instance
var GlobalTemplateProcessor = NewTemplateProcessor()

// Helper functions for common template operations
func HasTemplateVariables(content string) bool {
	return GlobalTemplateProcessor.HasTemplateVariables(content)
}

func GetTemplateVariables(content string) []string {
	return GlobalTemplateProcessor.GetTemplateVariables(content)
}

func ValidateTemplateVariables(content string) []string {
	return GlobalTemplateProcessor.ValidateTemplateVariables(content)
}

func SubstituteVariables(content string, variables map[string]string) string {
	return GlobalTemplateProcessor.SubstituteVariables(content, variables)
}

func ProcessTemplate(templateContent string, data interface{}) (string, error) {
	return GlobalTemplateProcessor.ProcessTemplate(templateContent, data)
}

func ValidateTemplateContent(content string) []string {
	return GlobalTemplateProcessor.ValidateTemplateContent(content)
}

func GetTemplateInfo(content string) TemplateInfo {
	return GlobalTemplateProcessor.GetTemplateInfo(content)
}

func CreateVariableMap(values map[string]string) map[string]string {
	return GlobalTemplateProcessor.CreateVariableMap(values)
}
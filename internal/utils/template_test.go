package utils

import (
	"reflect"
	"testing"
)

func TestTemplateProcessor_HasTemplateVariables(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "content with variables",
			content: "Hello {{name}}, welcome to {{project}}",
			want:    true,
		},
		{
			name:    "content without variables",
			content: "Hello world, welcome to the project",
			want:    false,
		},
		{
			name:    "content with spaced variables",
			content: "Hello {{ name }}, welcome to {{ project }}",
			want:    true,
		},
		{
			name:    "empty content",
			content: "",
			want:    false,
		},
		{
			name:    "content with single brace",
			content: "Hello {name}, welcome to the project",
			want:    false,
		},
		{
			name:    "content with incomplete variables",
			content: "Hello {{name, welcome to project}}",
			want:    true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processor.HasTemplateVariables(tt.content)
			if got != tt.want {
				t.Errorf("HasTemplateVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_GetTemplateVariables(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "simple variables",
			content: "Hello {{name}}, welcome to {{project}}",
			want:    []string{"name", "project"},
		},
		{
			name:    "variables with spaces",
			content: "Hello {{ name }}, welcome to {{ project }}",
			want:    []string{"name", "project"},
		},
		{
			name:    "duplicate variables",
			content: "{{name}} says hello to {{name}} again",
			want:    []string{"name"},
		},
		{
			name:    "no variables",
			content: "Hello world",
			want:    []string{},
		},
		{
			name:    "empty content",
			content: "",
			want:    []string{},
		},
		{
			name:    "mixed variables",
			content: "Project: {{project}}, Version: {{version}}, Path: {{cwd}}",
			want:    []string{"project", "version", "cwd"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processor.GetTemplateVariables(tt.content)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTemplateVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_ValidateTemplateVariables(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "all known variables",
			content: "{{project}} {{version}} {{name}}",
			want:    []string{},
		},
		{
			name:    "unknown variables",
			content: "{{project}} {{unknown}} {{invalid}}",
			want:    []string{"unknown", "invalid"},
		},
		{
			name:    "mixed known and unknown",
			content: "{{project}} {{version}} {{custom}}",
			want:    []string{"custom"},
		},
		{
			name:    "no variables",
			content: "plain text",
			want:    []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processor.ValidateTemplateVariables(tt.content)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateTemplateVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_SubstituteVariables(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name      string
		content   string
		variables map[string]string
		want      string
	}{
		{
			name:    "simple substitution",
			content: "Hello {{name}}, welcome to {{project}}",
			variables: map[string]string{
				"name":    "John",
				"project": "Atempo",
			},
			want: "Hello John, welcome to Atempo",
		},
		{
			name:    "variables with spaces",
			content: "Hello {{ name }}, welcome to {{ project }}",
			variables: map[string]string{
				"name":    "John",
				"project": "Atempo",
			},
			want: "Hello John, welcome to Atempo",
		},
		{
			name:    "partial substitution",
			content: "{{name}} works on {{project}} using {{language}}",
			variables: map[string]string{
				"name":    "John",
				"project": "Atempo",
			},
			want: "John works on Atempo using {{language}}",
		},
		{
			name:    "no variables",
			content: "plain text",
			variables: map[string]string{
				"name": "John",
			},
			want: "plain text",
		},
		{
			name:    "empty variables",
			content: "{{name}} {{project}}",
			variables: map[string]string{},
			want:    "{{name}} {{project}}",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processor.SubstituteVariables(tt.content, tt.variables)
			if got != tt.want {
				t.Errorf("SubstituteVariables() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_ProcessTemplate(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name         string
		templateContent string
		data         interface{}
		want         string
		wantErr      bool
	}{
		{
			name:         "simple template",
			templateContent: "Hello {{.Name}}, welcome to {{.Project}}",
			data:         map[string]string{"Name": "John", "Project": "Atempo"},
			want:         "Hello John, welcome to Atempo",
			wantErr:      false,
		},
		{
			name:         "template with helper functions",
			templateContent: "Hello {{.Name | upper}}, welcome to {{.Project | title}}",
			data:         map[string]string{"Name": "john", "Project": "atempo"},
			want:         "Hello JOHN, welcome to Atempo",
			wantErr:      false,
		},
		{
			name:         "template with conditions",
			templateContent: "{{if .Name}}Hello {{.Name}}{{else}}Hello guest{{end}}",
			data:         map[string]string{"Name": "John"},
			want:         "Hello John",
			wantErr:      false,
		},
		{
			name:         "template with range",
			templateContent: "Features: {{range .Features}}{{.}} {{end}}",
			data:         map[string]interface{}{"Features": []string{"auth", "db", "api"}},
			want:         "Features: auth db api ",
			wantErr:      false,
		},
		{
			name:         "invalid template",
			templateContent: "{{.Name",
			data:         map[string]string{"Name": "John"},
			want:         "",
			wantErr:      true,
		},
		{
			name:         "template with custom functions",
			templateContent: "{{.Text | replace \" \" \"_\"}}",
			data:         map[string]string{"Text": "hello world"},
			want:         "hello_world",
			wantErr:      false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processor.ProcessTemplate(tt.templateContent, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProcessTemplate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_IsValidTemplateVariable(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name     string
		variable string
		want     bool
	}{
		{
			name:     "valid simple variable",
			variable: "name",
			want:     true,
		},
		{
			name:     "valid variable with underscore",
			variable: "project_name",
			want:     true,
		},
		{
			name:     "valid variable with numbers",
			variable: "version2",
			want:     true,
		},
		{
			name:     "valid variable starting with underscore",
			variable: "_private",
			want:     true,
		},
		{
			name:     "invalid variable with space",
			variable: "project name",
			want:     false,
		},
		{
			name:     "invalid variable with hyphen",
			variable: "project-name",
			want:     false,
		},
		{
			name:     "invalid variable starting with number",
			variable: "2version",
			want:     false,
		},
		{
			name:     "invalid variable with special characters",
			variable: "project!",
			want:     false,
		},
		{
			name:     "empty variable",
			variable: "",
			want:     false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processor.IsValidTemplateVariable(tt.variable)
			if got != tt.want {
				t.Errorf("IsValidTemplateVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_ValidateTemplateContent(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "valid template",
			content: "Hello {{name}}, welcome to {{project}}",
			want:    []string{},
		},
		{
			name:    "unbalanced braces",
			content: "Hello {{name}}, welcome to {{project}",
			want:    []string{"unbalanced template braces"},
		},
		{
			name:    "nested braces",
			content: "Hello {{name {{project}}}}",
			want:    []string{"nested template braces detected", "invalid variable name: name {{project"},
		},
		{
			name:    "unclosed braces",
			content: "Hello {{name",
			want:    []string{"unbalanced template braces", "unclosed template braces"},
		},
		{
			name:    "empty variables",
			content: "Hello {{}}, welcome to {{project}}",
			want:    []string{"empty template variables"},
		},
		{
			name:    "invalid variable names",
			content: "Hello {{2name}}, welcome to {{project-name}}",
			want:    []string{"invalid variable name: 2name", "invalid variable name: project-name"},
		},
		{
			name:    "multiple issues",
			content: "Hello {{name}}, welcome to {{project} and {{}}",
			want:    []string{"unbalanced template braces", "empty template variables"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processor.ValidateTemplateContent(tt.content)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateTemplateContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_GetTemplateInfo(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name    string
		content string
		want    TemplateInfo
	}{
		{
			name:    "valid template with variables",
			content: "Hello {{name}}, welcome to {{project}}",
			want: TemplateInfo{
				HasVariables:     true,
				Variables:        []string{"name", "project"},
				UnknownVariables: []string{},
				Issues:           []string{},
				IsValid:          true,
			},
		},
		{
			name:    "template with unknown variables",
			content: "Hello {{name}}, welcome to {{custom}}",
			want: TemplateInfo{
				HasVariables:     true,
				Variables:        []string{"name", "custom"},
				UnknownVariables: []string{"custom"},
				Issues:           []string{},
				IsValid:          true,
			},
		},
		{
			name:    "invalid template",
			content: "Hello {{name}}, welcome to {{project}",
			want: TemplateInfo{
				HasVariables:     true,
				Variables:        []string{"name"},
				UnknownVariables: []string{},
				Issues:           []string{"unbalanced template braces"},
				IsValid:          false,
			},
		},
		{
			name:    "plain text",
			content: "Hello world",
			want: TemplateInfo{
				HasVariables:     false,
				Variables:        []string{},
				UnknownVariables: []string{},
				Issues:           []string{},
				IsValid:          true,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processor.GetTemplateInfo(tt.content)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTemplateInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_CreateVariableMap(t *testing.T) {
	processor := NewTemplateProcessor()
	
	tests := []struct {
		name   string
		values map[string]string
		want   map[string]string
	}{
		{
			name: "with custom values",
			values: map[string]string{
				"project": "Atempo",
				"version": "1.0",
				"custom":  "value",
			},
			want: map[string]string{
				"project":   "Atempo",
				"cwd":       "",
				"version":   "1.0",
				"name":      "",
				"framework": "",
				"language":  "",
				"custom":    "value",
			},
		},
		{
			name:   "empty values",
			values: map[string]string{},
			want: map[string]string{
				"project":   "",
				"cwd":       "",
				"version":   "",
				"name":      "",
				"framework": "",
				"language":  "",
			},
		},
		{
			name: "partial values",
			values: map[string]string{
				"name": "MyProject",
			},
			want: map[string]string{
				"project":   "",
				"cwd":       "",
				"version":   "",
				"name":      "MyProject",
				"framework": "",
				"language":  "",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processor.CreateVariableMap(tt.values)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateVariableMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobalTemplateProcessorFunctions(t *testing.T) {
	// Test global helper functions
	content := "Hello {{name}}, welcome to {{project}}"
	
	hasVars := HasTemplateVariables(content)
	if !hasVars {
		t.Error("HasTemplateVariables() should return true for content with variables")
	}
	
	vars := GetTemplateVariables(content)
	expected := []string{"name", "project"}
	if !reflect.DeepEqual(vars, expected) {
		t.Errorf("GetTemplateVariables() = %v, want %v", vars, expected)
	}
	
	unknownVars := ValidateTemplateVariables(content)
	if len(unknownVars) != 0 {
		t.Errorf("ValidateTemplateVariables() = %v, want empty slice", unknownVars)
	}
	
	substituted := SubstituteVariables(content, map[string]string{"name": "John", "project": "Atempo"})
	expectedSubstituted := "Hello John, welcome to Atempo"
	if substituted != expectedSubstituted {
		t.Errorf("SubstituteVariables() = %q, want %q", substituted, expectedSubstituted)
	}
	
	issues := ValidateTemplateContent(content)
	if len(issues) != 0 {
		t.Errorf("ValidateTemplateContent() = %v, want empty slice", issues)
	}
	
	info := GetTemplateInfo(content)
	if !info.HasVariables || !info.IsValid {
		t.Errorf("GetTemplateInfo() should indicate valid template with variables")
	}
	
	varMap := CreateVariableMap(map[string]string{"name": "test"})
	if varMap["name"] != "test" {
		t.Error("CreateVariableMap() should include custom values")
	}
	if varMap["project"] != "" {
		t.Error("CreateVariableMap() should include default empty values")
	}
	
	// Test ProcessTemplate
	result, err := ProcessTemplate("Hello {{.Name}}", map[string]string{"Name": "World"})
	if err != nil {
		t.Errorf("ProcessTemplate() unexpected error: %v", err)
	}
	if result != "Hello World" {
		t.Errorf("ProcessTemplate() = %q, want 'Hello World'", result)
	}
}
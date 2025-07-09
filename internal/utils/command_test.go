package utils

import (
	"testing"
)

func TestCommandFactory_FormatUsage(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name        string
		commandName string
		args        []string
		want        string
	}{
		{
			name:        "simple command",
			commandName: "status",
			args:        []string{},
			want:        "atempo status",
		},
		{
			name:        "command with optional arg",
			commandName: "status",
			args:        []string{"[project]"},
			want:        "atempo status [project]",
		},
		{
			name:        "command with required arg",
			commandName: "create",
			args:        []string{"<framework>:<version>"},
			want:        "atempo create <framework>:<version>",
		},
		{
			name:        "command with multiple args",
			commandName: "create",
			args:        []string{"<framework>:<version>", "[project-name]"},
			want:        "atempo create <framework>:<version> [project-name]",
		},
		{
			name:        "command with literal arg",
			commandName: "docker",
			args:        []string{"up"},
			want:        "atempo docker up",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.FormatUsage(tt.commandName, tt.args...)
			if got != tt.want {
				t.Errorf("FormatUsage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommandFactory_FormatUsageWithExamples(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name        string
		commandName string
		args        []string
		examples    []string
		want        string
	}{
		{
			name:        "no examples",
			commandName: "status",
			args:        []string{},
			examples:    []string{},
			want:        "atempo status",
		},
		{
			name:        "with examples",
			commandName: "create",
			args:        []string{"<framework>:<version>", "[project-name]"},
			examples:    []string{"atempo create laravel:11", "atempo create django:5 my-project"},
			want:        "atempo create <framework>:<version> [project-name]\nExamples:\n  atempo create laravel:11\n  atempo create django:5 my-project",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.FormatUsageWithExamples(tt.commandName, tt.args, tt.examples)
			if got != tt.want {
				t.Errorf("FormatUsageWithExamples() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommandFactory_FormatDescription(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "description without period",
			description: "Show project status",
			want:        "Show project status",
		},
		{
			name:        "description with period",
			description: "Show project status.",
			want:        "Show project status",
		},
		{
			name:        "empty description",
			description: "",
			want:        "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.FormatDescription(tt.description)
			if got != tt.want {
				t.Errorf("FormatDescription() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommandFactory_ValidateCommandName(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name    string
		cmdName string
		wantErr bool
	}{
		{
			name:    "valid simple name",
			cmdName: "status",
			wantErr: false,
		},
		{
			name:    "valid name with hyphen",
			cmdName: "add-service",
			wantErr: false,
		},
		{
			name:    "valid name with numbers",
			cmdName: "v1-status",
			wantErr: false,
		},
		{
			name:    "empty name",
			cmdName: "",
			wantErr: true,
		},
		{
			name:    "name with spaces",
			cmdName: "add service",
			wantErr: true,
		},
		{
			name:    "name starting with hyphen",
			cmdName: "-status",
			wantErr: true,
		},
		{
			name:    "name ending with hyphen",
			cmdName: "status-",
			wantErr: true,
		},
		{
			name:    "name with invalid characters",
			cmdName: "status!",
			wantErr: true,
		},
		{
			name:    "name with underscore",
			cmdName: "add_service",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.ValidateCommandName(tt.cmdName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommandName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandFactory_CreateStandardUsage(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name        string
		commandName string
		pattern     CommandPattern
		args        []string
		want        string
	}{
		{
			name:        "simple pattern",
			commandName: "status",
			pattern:     PatternSimple,
			args:        []string{},
			want:        "atempo status",
		},
		{
			name:        "project context pattern",
			commandName: "logs",
			pattern:     PatternWithProjectContext,
			args:        []string{},
			want:        "atempo logs [project]",
		},
		{
			name:        "framework spec pattern",
			commandName: "create",
			pattern:     PatternWithFrameworkSpec,
			args:        []string{},
			want:        "atempo create <framework>:<version> [project-name]",
		},
		{
			name:        "docker passthrough pattern",
			commandName: "docker",
			pattern:     PatternDockerPassthrough,
			args:        []string{},
			want:        "atempo docker <docker-command> [args...]",
		},
		{
			name:        "with optional args pattern",
			commandName: "test",
			pattern:     PatternWithOptionalArgs,
			args:        []string{"[flags]"},
			want:        "atempo test [flags]",
		},
		{
			name:        "with required args pattern",
			commandName: "remove",
			pattern:     PatternWithRequiredArgs,
			args:        []string{"<project>"},
			want:        "atempo remove <project>",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.CreateStandardUsage(tt.commandName, tt.pattern, tt.args...)
			if got != tt.want {
				t.Errorf("CreateStandardUsage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommandFactory_ValidateCommandMetadata(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name     string
		metadata CommandMetadata
		wantErr  bool
	}{
		{
			name: "valid metadata",
			metadata: CommandMetadata{
				Name:        "status",
				Description: "Show project status",
				Pattern:     PatternSimple,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			metadata: CommandMetadata{
				Name:        "",
				Description: "Show project status",
				Pattern:     PatternSimple,
			},
			wantErr: true,
		},
		{
			name: "empty description",
			metadata: CommandMetadata{
				Name:        "status",
				Description: "",
				Pattern:     PatternSimple,
			},
			wantErr: true,
		},
		{
			name: "invalid name",
			metadata: CommandMetadata{
				Name:        "status!",
				Description: "Show project status",
				Pattern:     PatternSimple,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.ValidateCommandMetadata(tt.metadata)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommandMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandFactory_CreateUsageFromMetadata(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name     string
		metadata CommandMetadata
		want     string
	}{
		{
			name: "metadata without examples",
			metadata: CommandMetadata{
				Name:        "status",
				Description: "Show project status",
				Pattern:     PatternSimple,
			},
			want: "atempo status",
		},
		{
			name: "metadata with examples",
			metadata: CommandMetadata{
				Name:        "create",
				Description: "Create new project",
				Pattern:     PatternWithFrameworkSpec,
				Examples:    []string{"atempo create laravel:11", "atempo create django:5 my-project"},
			},
			want: "atempo create <framework>:<version> [project-name]\nExamples:\n  atempo create laravel:11\n  atempo create django:5 my-project",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.CreateUsageFromMetadata(tt.metadata)
			if got != tt.want {
				t.Errorf("CreateUsageFromMetadata() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommandFactory_GetStandardDescription(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name        string
		commandName string
		want        string
	}{
		{
			name:        "known command",
			commandName: "status",
			want:        "Show project dashboard with health status",
		},
		{
			name:        "unknown command",
			commandName: "unknown",
			want:        "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.GetStandardDescription(tt.commandName)
			if got != tt.want {
				t.Errorf("GetStandardDescription() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommandFactory_CreateDockerSubcommandUsage(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name       string
		subcommand string
		args       []string
		want       string
	}{
		{
			name:       "simple subcommand",
			subcommand: "up",
			args:       []string{},
			want:       "atempo docker up",
		},
		{
			name:       "subcommand with args",
			subcommand: "exec",
			args:       []string{"<service>", "<command>"},
			want:       "atempo docker exec <service> <command>",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.CreateDockerSubcommandUsage(tt.subcommand, tt.args...)
			if got != tt.want {
				t.Errorf("CreateDockerSubcommandUsage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommandFactory_CreateProjectSpecificUsage(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name        string
		commandName string
		projectName string
		args        []string
		want        string
	}{
		{
			name:        "simple project command",
			commandName: "status",
			projectName: "my-project",
			args:        []string{},
			want:        "my-project status",
		},
		{
			name:        "project command with args",
			commandName: "docker",
			projectName: "my-project",
			args:        []string{"up", "-d"},
			want:        "my-project docker up -d",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.CreateProjectSpecificUsage(tt.commandName, tt.projectName, tt.args...)
			if got != tt.want {
				t.Errorf("CreateProjectSpecificUsage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommandFactory_ValidateArgumentPattern(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{
			name:    "valid pattern with square brackets",
			pattern: "[optional]",
			wantErr: false,
		},
		{
			name:    "valid pattern with angle brackets",
			pattern: "<required>",
			wantErr: false,
		},
		{
			name:    "valid pattern with mixed brackets",
			pattern: "<required> [optional]",
			wantErr: false,
		},
		{
			name:    "unbalanced square brackets",
			pattern: "[optional",
			wantErr: true,
		},
		{
			name:    "unbalanced angle brackets",
			pattern: "<required",
			wantErr: true,
		},
		{
			name:    "empty pattern",
			pattern: "",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.ValidateArgumentPattern(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateArgumentPattern() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandFactory_SuggestCommandName(t *testing.T) {
	factory := NewCommandFactory()
	
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "create command",
			description: "Create a new project",
			want:        "create",
		},
		{
			name:        "show command",
			description: "Show project status",
			want:        "project",
		},
		{
			name:        "list command",
			description: "List all projects",
			want:        "all",
		},
		{
			name:        "remove command",
			description: "Remove project from registry",
			want:        "remove",
		},
		{
			name:        "start command",
			description: "Start project services",
			want:        "start",
		},
		{
			name:        "stop command",
			description: "Stop project services",
			want:        "stop",
		},
		{
			name:        "generic command",
			description: "Manage project configuration",
			want:        "manage",
		},
		{
			name:        "empty description",
			description: "",
			want:        "command",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := factory.SuggestCommandName(tt.description)
			if got != tt.want {
				t.Errorf("SuggestCommandName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGlobalCommandFactoryFunctions(t *testing.T) {
	// Test global helper functions
	usage := FormatUsage("test", "[args]")
	expected := "atempo test [args]"
	if usage != expected {
		t.Errorf("FormatUsage() = %q, want %q", usage, expected)
	}
	
	err := ValidateCommandName("test-command")
	if err != nil {
		t.Errorf("ValidateCommandName() unexpected error: %v", err)
	}
	
	err = ValidateCommandName("invalid command")
	if err == nil {
		t.Error("ValidateCommandName() expected error for invalid name, got nil")
	}
	
	desc := GetStandardDescription("create")
	if desc != "Create a new project with specified framework" {
		t.Errorf("GetStandardDescription() = %q, want standard description", desc)
	}
	
	dockerUsage := CreateDockerSubcommandUsage("up", "-d")
	expectedDockerUsage := "atempo docker up -d"
	if dockerUsage != expectedDockerUsage {
		t.Errorf("CreateDockerSubcommandUsage() = %q, want %q", dockerUsage, expectedDockerUsage)
	}
	
	projectUsage := CreateProjectSpecificUsage("status", "my-project")
	expectedProjectUsage := "my-project status"
	if projectUsage != expectedProjectUsage {
		t.Errorf("CreateProjectSpecificUsage() = %q, want %q", projectUsage, expectedProjectUsage)
	}
}
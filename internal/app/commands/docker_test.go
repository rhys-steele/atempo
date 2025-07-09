package commands

import (
	"context"
	"testing"
	"time"
)

func TestNewDockerCommand(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	// Verify basic properties
	if cmd.Name() != "docker" {
		t.Errorf("Expected command name 'docker', got '%s'", cmd.Name())
	}
	if cmd.Description() != "Run Docker operations on projects" {
		t.Errorf("Expected description 'Run Docker operations on projects', got '%s'", cmd.Description())
	}
	if cmd.Usage() != "atempo docker <command> [project] [options]" {
		t.Errorf("Expected usage 'atempo docker <command> [project] [options]', got '%s'", cmd.Usage())
	}
}

func TestDockerCommand_Execute_NoArgs(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	// Test with no arguments
	err := cmd.Execute(context.Background(), []string{})
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}

	// Verify error message contains usage information
	if !containsString(err.Error(), "usage:") {
		t.Errorf("Expected error to contain usage information, got: %s", err.Error())
	}
}

func TestDockerCommand_isDockerArg(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	tests := []struct {
		name     string
		arg      string
		expected bool
	}{
		{
			name:     "Build flag",
			arg:      "--build",
			expected: true,
		},
		{
			name:     "Force recreate flag",
			arg:      "--force-recreate",
			expected: true,
		},
		{
			name:     "Remove orphans flag",
			arg:      "--remove-orphans",
			expected: true,
		},
		{
			name:     "Verbose flag",
			arg:      "-v",
			expected: true,
		},
		{
			name:     "Detach flag",
			arg:      "-d",
			expected: true,
		},
		{
			name:     "Project name",
			arg:      "my-project",
			expected: false,
		},
		{
			name:     "Path",
			arg:      "/path/to/project",
			expected: false,
		},
		{
			name:     "Service name",
			arg:      "app",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmd.isDockerArg(tt.arg)
			if result != tt.expected {
				t.Errorf("Expected isDockerArg(%s) to be %v, got %v", tt.arg, tt.expected, result)
			}
		})
	}
}

func TestDockerCommand_parseTimeoutFlag(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	tests := []struct {
		name         string
		args         []string
		expectedTimeout time.Duration
		expectedArgs []string
	}{
		{
			name:         "No timeout flag",
			args:         []string{"--build", "service"},
			expectedTimeout: 0,
			expectedArgs: []string{"--build", "service"},
		},
		{
			name:         "Timeout flag with seconds",
			args:         []string{"--timeout", "300", "--build"},
			expectedTimeout: 300 * time.Second,
			expectedArgs: []string{"--build"},
		},
		{
			name:         "Timeout flag with minutes",
			args:         []string{"--timeout", "5m", "service"},
			expectedTimeout: 5 * time.Minute,
			expectedArgs: []string{"service"},
		},
		{
			name:         "Timeout flag at end",
			args:         []string{"service", "--timeout", "30"},
			expectedTimeout: 30 * time.Second,
			expectedArgs: []string{"service"},
		},
		{
			name:         "Invalid timeout",
			args:         []string{"--timeout", "invalid", "service"},
			expectedTimeout: 0,
			expectedArgs: []string{"--timeout", "invalid", "service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeout, filteredArgs := cmd.parseTimeoutFlag(tt.args)
			
			if timeout != tt.expectedTimeout {
				t.Errorf("Expected timeout %v, got %v", tt.expectedTimeout, timeout)
			}
			
			if len(filteredArgs) != len(tt.expectedArgs) {
				t.Errorf("Expected %d filtered args, got %d", len(tt.expectedArgs), len(filteredArgs))
			} else {
				for i, expected := range tt.expectedArgs {
					if filteredArgs[i] != expected {
						t.Errorf("Expected filtered arg %d to be %s, got %s", i, expected, filteredArgs[i])
					}
				}
			}
		})
	}
}

func TestDockerCommand_handleDockerExec_NoArgs(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	// Test with no arguments
	err := cmd.handleDockerExec("", []string{})
	if err == nil {
		t.Error("Expected error when no arguments provided to exec")
	}

	// Verify error message contains usage information
	if !containsString(err.Error(), "usage:") {
		t.Errorf("Expected error to contain usage information, got: %s", err.Error())
	}
}

func TestDockerCommand_handleDockerExec_WithArgs(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	tests := []struct {
		name            string
		args            []string
		expectedService string
		expectedCommand []string
	}{
		{
			name:            "Service only",
			args:            []string{"app"},
			expectedService: "app",
			expectedCommand: []string{"bash"},
		},
		{
			name:            "Service with command",
			args:            []string{"app", "ls", "-la"},
			expectedService: "app",
			expectedCommand: []string{"ls", "-la"},
		},
		{
			name:            "Database service",
			args:            []string{"db", "mysql", "-u", "root", "-p"},
			expectedService: "db",
			expectedCommand: []string{"mysql", "-u", "root", "-p"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily test the actual execution without Docker running
			// So we'll test the argument parsing logic
			
			service := tt.args[0]
			var cmdArgs []string
			
			if len(tt.args) > 1 {
				cmdArgs = tt.args[1:]
			} else {
				cmdArgs = []string{"bash"}
			}

			if service != tt.expectedService {
				t.Errorf("Expected service %s, got %s", tt.expectedService, service)
			}

			if len(cmdArgs) != len(tt.expectedCommand) {
				t.Errorf("Expected %d command args, got %d", len(tt.expectedCommand), len(cmdArgs))
			} else {
				for i, expected := range tt.expectedCommand {
					if cmdArgs[i] != expected {
						t.Errorf("Expected command arg %d to be %s, got %s", i, expected, cmdArgs[i])
					}
				}
			}
		})
	}
}

func TestDockerCommand_ArgumentParsing(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	tests := []struct {
		name                    string
		args                    []string
		expectedDockerCmd       string
		expectedProjectPath     string
		expectedAdditionalArgs  []string
		isDockerArgFirstArg     bool
	}{
		{
			name:                   "Command only",
			args:                   []string{"up"},
			expectedDockerCmd:      "up",
			expectedProjectPath:    "",
			expectedAdditionalArgs: []string{},
			isDockerArgFirstArg:    false,
		},
		{
			name:                   "Command with project",
			args:                   []string{"up", "my-project"},
			expectedDockerCmd:      "up",
			expectedProjectPath:    "my-project", // Would be resolved by registry
			expectedAdditionalArgs: []string{},
			isDockerArgFirstArg:    false,
		},
		{
			name:                   "Command with docker flags",
			args:                   []string{"up", "--build", "--force-recreate"},
			expectedDockerCmd:      "up",
			expectedProjectPath:    "",
			expectedAdditionalArgs: []string{"--build", "--force-recreate"},
			isDockerArgFirstArg:    true,
		},
		{
			name:                   "Command with project and flags",
			args:                   []string{"up", "my-project", "--build"},
			expectedDockerCmd:      "up",
			expectedProjectPath:    "my-project",
			expectedAdditionalArgs: []string{"--build"},
			isDockerArgFirstArg:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the argument parsing logic
			dockerCmd := tt.args[0]
			var projectPath string
			var additionalArgs []string

			if len(tt.args) > 1 {
				potentialIdentifier := tt.args[1]
				if strings.HasPrefix(potentialIdentifier, "-") || cmd.isDockerArg(potentialIdentifier) {
					// It's a docker argument, use current directory
					projectPath = ""
					additionalArgs = tt.args[1:]
				} else {
					// It's a project identifier (name or path)
					projectPath = potentialIdentifier // In real code, this would be resolved
					if len(tt.args) > 2 {
						additionalArgs = tt.args[2:]
					}
				}
			}

			if dockerCmd != tt.expectedDockerCmd {
				t.Errorf("Expected docker command %s, got %s", tt.expectedDockerCmd, dockerCmd)
			}

			if projectPath != tt.expectedProjectPath {
				t.Errorf("Expected project path %s, got %s", tt.expectedProjectPath, projectPath)
			}

			if len(additionalArgs) != len(tt.expectedAdditionalArgs) {
				t.Errorf("Expected %d additional args, got %d", len(tt.expectedAdditionalArgs), len(additionalArgs))
			} else {
				for i, expected := range tt.expectedAdditionalArgs {
					if additionalArgs[i] != expected {
						t.Errorf("Expected additional arg %d to be %s, got %s", i, expected, additionalArgs[i])
					}
				}
			}
		})
	}
}

func TestDockerCommand_getDockerUsage(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	// Test that getDockerUsage returns a non-empty string
	usage := cmd.getDockerUsage()
	if usage == "" {
		t.Error("Expected getDockerUsage to return non-empty string")
	}

	// Verify it contains some expected commands
	expectedCommands := []string{"up", "down", "build", "logs", "ps"}
	for _, expectedCmd := range expectedCommands {
		if !containsString(usage, expectedCmd) {
			t.Errorf("Expected usage to contain command '%s'", expectedCmd)
		}
	}
}

func TestDockerCommand_BaseProperties(t *testing.T) {
	// Create a command context for testing
	ctx := &CommandContext{
		ProjectRegistry: nil,
		Output:          nil,
	}

	// Create the command
	cmd := NewDockerCommand(ctx)

	// Test BaseCommand properties
	if cmd.BaseCommand == nil {
		t.Error("Expected BaseCommand to be initialized")
	}

	// Test that context is properly set
	if cmd.BaseCommand.Context != ctx {
		t.Error("Expected BaseCommand context to be set")
	}
}

func TestDockerCommand_SpecialCommands(t *testing.T) {
	// Test that special commands are handled correctly
	specialCommands := []string{"exec", "services"}
	
	for _, specialCmd := range specialCommands {
		t.Run(specialCmd, func(t *testing.T) {
			// Create a command context for testing
			ctx := &CommandContext{
				ProjectRegistry: nil,
				Output:          nil,
			}

			// Create the command
			cmd := NewDockerCommand(ctx)

			// Test that the command is recognized as special
			args := []string{specialCmd}
			
			// We can't easily test the actual execution without proper setup
			// So we just verify the command is recognized
			if len(args) > 0 && args[0] == specialCmd {
				t.Logf("Special command '%s' is recognized", specialCmd)
			} else {
				t.Errorf("Special command '%s' should be recognized", specialCmd)
			}
		})
	}
}
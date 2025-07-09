package commands

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"atempo/internal/docker"
	"atempo/internal/mcp"
	"atempo/internal/registry"
	"atempo/internal/utils"
)

// TestCommand runs tests for a project using framework-specific test commands
type TestCommand struct {
	*BaseCommand
}

// NewTestCommand creates a new test command
func NewTestCommand(ctx *CommandContext) *TestCommand {
	return &TestCommand{
		BaseCommand: NewBaseCommand(
			"test",
			utils.GetStandardDescription("test"),
			utils.CreateStandardUsage("test", utils.PatternWithOptionalArgs, "[project]", "[suite]"),
			ctx,
		),
	}
}

// Execute runs the test command
func (c *TestCommand) Execute(ctx context.Context, args []string) error {
	var projectPath string
	var testSuite string

	// Parse arguments to distinguish between project path and test suite
	if len(args) > 0 {
		firstArg := args[0]

		// If it contains path separators, treat as project path
		if strings.Contains(firstArg, "/") || strings.Contains(firstArg, "\\") {
			// Resolve as project path
			resolution, err := utils.ResolveProjectPathFromArgs(args[:1])
			if err != nil {
				return fmt.Errorf("failed to resolve project path: %w", err)
			}
			projectPath = resolution.Path

			// Second arg could be test suite
			if len(args) > 1 {
				testSuite = args[1]
			}
		} else {
			// Try to find as registered project first
			reg, err := registry.LoadRegistry()
			if err != nil {
				return fmt.Errorf("failed to load registry: %w", err)
			}

			_, err = reg.FindProject(firstArg)
			if err == nil {
				// Found as registered project
				resolution, err := utils.ResolveProjectPathFromArgs(args[:1])
				if err != nil {
					return fmt.Errorf("failed to resolve project: %w", err)
				}
				projectPath = resolution.Path

				// Second arg could be test suite
				if len(args) > 1 {
					testSuite = args[1]
				}
			} else {
				// Not a registered project, treat as test suite and use current directory
				resolution, err := utils.ResolveCurrentProjectPath()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				projectPath = resolution.Path
				testSuite = firstArg
			}
		}
	} else {
		// No args, use current directory
		resolution, err := utils.ResolveCurrentProjectPath()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectPath = resolution.Path
	}

	// Check if project has atempo.json
	atempoJSONPath := filepath.Join(projectPath, "atempo.json")
	if !utils.FileExists(atempoJSONPath) {
		return fmt.Errorf("no atempo.json found in %s - this doesn't appear to be an Atempo project", projectPath)
	}

	// Try to use MCP server for framework-agnostic testing
	fmt.Printf("→ Initializing MCP server for testing...\n")

	mcpClient, err := mcp.NewMCPClient(projectPath)
	if err != nil {
		// Fallback to legacy approach if MCP server not available
		fmt.Printf("⚠ MCP server not available, falling back to legacy testing: %v\n", err)
		return c.runLegacyTest(projectPath, testSuite)
	}

	// Start the MCP server
	if err := mcpClient.Start(); err != nil {
		fmt.Printf("⚠ Failed to start MCP server, falling back to legacy testing: %v\n", err)
		return c.runLegacyTest(projectPath, testSuite)
	}
	defer mcpClient.Close()

	fmt.Printf("→ Running tests via MCP server...\n")
	if testSuite != "" {
		fmt.Printf("→ Test suite: %s\n", testSuite)
	}

	// Execute tests through MCP server
	err = mcpClient.RunTests(testSuite)
	if err != nil {
		fmt.Printf("✗ Tests failed: %v\n", err)
		return err
	}

	fmt.Printf("✓ Tests completed successfully!\n")
	return nil
}

// runLegacyTest executes tests using the legacy approach (fallback when MCP not available)
func (c *TestCommand) runLegacyTest(projectPath, testSuite string) error {
	// Get framework information
	framework, err := c.detectFramework(projectPath)
	if err != nil {
		return fmt.Errorf("failed to detect framework: %w", err)
	}

	// Get framework-specific test command
	testCommand, containerName, err := c.getTestCommand(framework, testSuite)
	if err != nil {
		return fmt.Errorf("failed to get test command for %s: %w", framework, err)
	}

	fmt.Printf("→ Running tests for %s project...\n", framework)
	if testSuite != "" {
		fmt.Printf("→ Test suite: %s\n", testSuite)
	}
	fmt.Printf("→ Command: %s\n", testCommand)

	// Execute the test command in the appropriate container
	err = c.runTestInContainer(projectPath, containerName, testCommand)
	if err != nil {
		fmt.Printf("✗ Tests failed: %v\n", err)
		return err
	}

	fmt.Printf("✓ Tests completed successfully!\n")
	return nil
}

// detectFramework detects the framework of a project
func (c *TestCommand) detectFramework(projectPath string) (string, error) {
	// Use centralized framework detection
	framework, err := utils.DetectFramework(projectPath)
	if err == nil {
		return framework, nil
	}

	// Fallback to file-based detection
	srcPath := filepath.Join(projectPath, "src")

	// Check for Laravel
	if utils.FileExists(filepath.Join(srcPath, "artisan")) ||
		utils.FileExists(filepath.Join(srcPath, "composer.json")) {
		return "laravel", nil
	}

	// Check for Django
	if utils.FileExists(filepath.Join(srcPath, "manage.py")) ||
		utils.FileExists(filepath.Join(srcPath, "requirements.txt")) {
		return "django", nil
	}

	return "", fmt.Errorf("unknown framework - unable to detect Laravel or Django")
}

// getTestCommand returns the appropriate test command and container for the framework
func (c *TestCommand) getTestCommand(framework, testSuite string) (string, string, error) {
	switch framework {
	case "laravel":
		containerName := "app"
		if testSuite == "" {
			return "php artisan test", containerName, nil
		}
		// Laravel supports test filtering
		return fmt.Sprintf("php artisan test --filter=%s", testSuite), containerName, nil

	case "django":
		containerName := "web"
		if testSuite == "" {
			return "python manage.py test", containerName, nil
		}
		// Django supports app-specific testing
		return fmt.Sprintf("python manage.py test %s", testSuite), containerName, nil

	default:
		return "", "", fmt.Errorf("unsupported framework: %s", framework)
	}
}

// runTestInContainer executes the test command in the appropriate Docker container
func (c *TestCommand) runTestInContainer(projectPath, containerName, testCommand string) error {
	// Use docker-compose exec to run the test command in the running container
	cmdArgs := []string{"sh", "-c", testCommand}

	// Use the existing docker exec infrastructure
	return docker.ExecuteExecCommand(containerName, projectPath, cmdArgs)
}
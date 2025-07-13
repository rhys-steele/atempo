package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"atempo/internal/utils"
)

// AuditCommand runs a comprehensive codebase audit using Claude Code
type AuditCommand struct {
	*BaseCommand
	authChecker *AuthChecker
}

// NewAuditCommand creates a new audit command
func NewAuditCommand(ctx *CommandContext) *AuditCommand {
	return &AuditCommand{
		BaseCommand: NewBaseCommand(
			"audit",
			utils.GetStandardDescription("audit"),
			utils.FormatUsage("audit"),
			ctx,
		),
		authChecker: NewAuthChecker(),
	}
}

// Execute runs the audit command
func (c *AuditCommand) Execute(ctx context.Context, args []string) error {
	// Handle help requests
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help") {
		fmt.Printf("Usage: %s\n", c.Usage())
		fmt.Printf("Description: %s\n\n", c.Description())
		fmt.Printf("This command runs a comprehensive codebase audit using Claude Code.\n")
		fmt.Printf("It requires:\n")
		fmt.Printf("  • Authentication with AI providers (run 'atempo auth')\n")
		fmt.Printf("  • Claude Code CLI installed (https://claude.ai/code)\n")
		fmt.Printf("  • An audit prompt template in .claude/commands/audit.md\n\n")
		fmt.Printf("Example:\n")
		fmt.Printf("  atempo audit                 # Run audit in current project\n")
		fmt.Printf("  my-project audit             # Run audit in specific project\n")
		return nil
	}

	// Check if AI features are enabled
	isAuthenticated, authMessage := c.authChecker.GetAuthStatus()
	if !isAuthenticated {
		fmt.Printf("✗ %s\n", authMessage)
		fmt.Printf("  AI audit features require authentication.\n")
		fmt.Printf("  Run %s'atempo auth'%s to enable AI features.\n\n", ColorCyan, ColorReset)
		return fmt.Errorf("authentication required for audit feature")
	}

	fmt.Printf("✓ %s\n", authMessage)
	fmt.Printf("\nRunning comprehensive codebase audit...\n")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if audit prompt file exists
	auditPromptPath := filepath.Join(cwd, ".claude", "commands", "audit.md")
	if _, err := os.Stat(auditPromptPath); os.IsNotExist(err) {
		return fmt.Errorf("audit prompt not found at %s\nThis command should be run from a project created with Atempo", auditPromptPath)
	}

	// Check if claude command is available
	if !c.isClaudeAvailable() {
		return fmt.Errorf("claude command not found\nPlease install Claude Code CLI: https://claude.ai/code")
	}

	// Execute claude audit command
	fmt.Printf("⎿ Running: claude /audit (in %s)\n", cwd)
	
	cmd := exec.CommandContext(ctx, "claude", "/audit")
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("audit command failed: %w", err)
	}

	fmt.Printf("\n✓ Audit completed successfully\n")
	return nil
}

// isClaudeAvailable checks if the claude command is available in PATH
func (c *AuditCommand) isClaudeAvailable() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}
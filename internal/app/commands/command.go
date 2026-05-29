package commands

import (
	"context"
	"fmt"
	"strings"

	"atempo/internal/interfaces"
)

// ANSI colour codes used across all commands.
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
	ColorBold   = "\033[1m"
)

// Command is the interface every CLI command must implement.
type Command interface {
	Execute(ctx context.Context, args []string) error
	Name() string
	Description() string
	Usage() string
	ShowHelp() string
}

// CommandContext holds shared dependencies injected into every command.
type CommandContext struct {
	Registry interfaces.RegistryRepository
}

// BaseCommand provides name, description, usage, and help handling that all
// concrete commands embed.
type BaseCommand struct {
	name        string
	description string
	usage       string
	helpText    string
	ctx         *CommandContext
}

func newBase(name, description, usage string, ctx *CommandContext) *BaseCommand {
	return &BaseCommand{
		name:        name,
		description: description,
		usage:       usage,
		ctx:         ctx,
	}
}

func (c *BaseCommand) Name() string        { return c.name }
func (c *BaseCommand) Description() string { return c.description }
func (c *BaseCommand) Usage() string       { return c.usage }

func (c *BaseCommand) ShowHelp() string {
	if c.helpText != "" {
		return c.helpText
	}
	return fmt.Sprintf("Usage: %s\n\n%s\n", c.usage, c.description)
}

func (c *BaseCommand) setHelp(text string) { c.helpText = text }

// handleHelp checks whether the first arg is a help flag and prints help if so.
// Returns true if help was shown (caller should return nil immediately).
func (c *BaseCommand) handleHelp(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch strings.ToLower(args[0]) {
	case "help", "--help", "-h":
		fmt.Print(c.ShowHelp())
		return true
	}
	return false
}

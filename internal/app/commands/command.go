package commands

import (
	"context"
)

// Command represents a CLI command that can be executed
type Command interface {
	// Execute runs the command with the given arguments
	Execute(ctx context.Context, args []string) error

	// Name returns the command name (e.g., "start", "docker")
	Name() string

	// Description returns a brief description of the command
	Description() string

	// Usage returns usage information for the command
	Usage() string
}

// CommandContext provides shared dependencies for commands
type CommandContext struct {
	// Add shared dependencies here as we refactor
	// For now, keeping it simple
}

// BaseCommand provides common functionality for all commands
type BaseCommand struct {
	name        string
	description string
	usage       string
	ctx         *CommandContext
}

// NewBaseCommand creates a new base command
func NewBaseCommand(name, description, usage string, ctx *CommandContext) *BaseCommand {
	return &BaseCommand{
		name:        name,
		description: description,
		usage:       usage,
		ctx:         ctx,
	}
}

// Name returns the command name
func (c *BaseCommand) Name() string {
	return c.name
}

// Description returns the command description
func (c *BaseCommand) Description() string {
	return c.description
}

// Usage returns the command usage
func (c *BaseCommand) Usage() string {
	return c.usage
}

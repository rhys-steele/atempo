package commands

// This file contains placeholder constructors for commands not yet implemented.
// Each will be replaced by its own file as phases complete.

import "context"

// --- create ---

type CreateCommand struct{ *BaseCommand }

func NewCreateCommand(ctx *CommandContext, _ interface{}) *CreateCommand {
	return &CreateCommand{newBase("create", "Scaffold a new project", "atempo create <framework>[:<version>] [name]", ctx)}
}
func (c *CreateCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 7
}


// --- add-service ---

type AddServiceCommand struct{ *BaseCommand }

func NewAddServiceCommand(ctx *CommandContext) *AddServiceCommand {
	return &AddServiceCommand{newBase("add-service", "Add a service to the project", "atempo add-service <service>", ctx)}
}
func (c *AddServiceCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 8
}

// --- shell ---

type ShellCommand struct{ *BaseCommand }

func NewShellCommand(ctx *CommandContext, _ *CommandRegistry) *ShellCommand {
	return &ShellCommand{newBase("shell", "Start interactive shell", "atempo shell", ctx)}
}
func (c *ShellCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 8
}

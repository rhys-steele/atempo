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

// --- projects ---

type ProjectsCommand struct{ *BaseCommand }

func NewProjectsCommand(ctx *CommandContext) *ProjectsCommand {
	return &ProjectsCommand{newBase("projects", "List registered projects", "atempo projects", ctx)}
}
func (c *ProjectsCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 4
}

// --- describe ---

type DescribeCommand struct{ *BaseCommand }

func NewDescribeCommand(ctx *CommandContext) *DescribeCommand {
	return &DescribeCommand{newBase("describe", "Show project details", "atempo describe [project]", ctx)}
}
func (c *DescribeCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 4
}

// --- start ---

type StartCommand struct{ *BaseCommand }

func NewStartCommand(ctx *CommandContext) *StartCommand {
	return &StartCommand{newBase("start", "Start project containers", "atempo start [project]", ctx)}
}
func (c *StartCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 6
}

// --- stop ---

type StopCommand struct{ *BaseCommand }

func NewStopCommand(ctx *CommandContext) *StopCommand {
	return &StopCommand{newBase("stop", "Stop project containers", "atempo stop [project]", ctx)}
}
func (c *StopCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 6
}

// --- logs ---

type LogsCommand struct{ *BaseCommand }

func NewLogsCommand(ctx *CommandContext) *LogsCommand {
	return &LogsCommand{newBase("logs", "Stream container logs", "atempo logs [project]", ctx)}
}
func (c *LogsCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 6
}

// --- status ---

type StatusCommand struct{ *BaseCommand }

func NewStatusCommand(ctx *CommandContext) *StatusCommand {
	return &StatusCommand{newBase("status", "Show service health", "atempo status [project]", ctx)}
}
func (c *StatusCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 6
}

// --- add-service ---

type AddServiceCommand struct{ *BaseCommand }

func NewAddServiceCommand(ctx *CommandContext) *AddServiceCommand {
	return &AddServiceCommand{newBase("add-service", "Add a service to the project", "atempo add-service <service>", ctx)}
}
func (c *AddServiceCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 8
}

// --- reconfigure ---

type ReconfigureCommand struct{ *BaseCommand }

func NewReconfigureCommand(ctx *CommandContext) *ReconfigureCommand {
	return &ReconfigureCommand{newBase("reconfigure", "Regenerate docker-compose.yml", "atempo reconfigure [project]", ctx)}
}
func (c *ReconfigureCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 5
}

// --- remove ---

type RemoveCommand struct{ *BaseCommand }

func NewRemoveCommand(ctx *CommandContext) *RemoveCommand {
	return &RemoveCommand{newBase("remove", "Remove a project from the registry", "atempo remove <project>", ctx)}
}
func (c *RemoveCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 4
}

// --- shell ---

type ShellCommand struct{ *BaseCommand }

func NewShellCommand(ctx *CommandContext, _ *CommandRegistry) *ShellCommand {
	return &ShellCommand{newBase("shell", "Start interactive shell", "atempo shell", ctx)}
}
func (c *ShellCommand) Execute(_ context.Context, _ []string) error {
	return nil // Phase 8
}

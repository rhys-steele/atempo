package commands

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"atempo/internal/interfaces"
)

// CommandRegistry routes CLI arguments to the correct Command implementation.
type CommandRegistry struct {
	commands map[string]Command
	ctx      *CommandContext
}

// NewCommandRegistry wires up all V1 commands with shared dependencies.
func NewCommandRegistry(repo interfaces.RegistryRepository, templatesFS embed.FS) *CommandRegistry {
	ctx := &CommandContext{Registry: repo}

	r := &CommandRegistry{
		commands: make(map[string]Command),
		ctx:      ctx,
	}

	r.register(NewCreateCommand(ctx, templatesFS))
	r.register(NewDoctorCommand(ctx))
	r.register(NewProjectsCommand(ctx))
	r.register(NewDescribeCommand(ctx))
	r.register(NewStartCommand(ctx))
	r.register(NewStopCommand(ctx))
	r.register(NewLogsCommand(ctx))
	r.register(NewStatusCommand(ctx))
	r.register(NewAddServiceCommand(ctx))
	r.register(NewReconfigureCommand(ctx))
	r.register(NewRemoveCommand(ctx))
	r.register(NewShellCommand(ctx, r))

	return r
}

func (r *CommandRegistry) register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Execute dispatches a command by name, falling back to project-name routing.
func (r *CommandRegistry) Execute(ctx context.Context, name string, args []string) error {
	if cmd, ok := r.commands[name]; ok {
		return cmd.Execute(ctx, args)
	}

	// Allow `atempo <project> <subcommand>` shorthand
	if r.isProjectName(name) {
		if len(args) == 0 {
			return fmt.Errorf("usage: atempo %s <command>  (try: start, stop, logs, status)", name)
		}
		return r.routeProjectCommand(ctx, name, args[0], args[1:])
	}

	return r.unknownCommand(name)
}

// routeProjectCommand maps `atempo <project> <sub>` to the equivalent global command.
func (r *CommandRegistry) routeProjectCommand(ctx context.Context, project, sub string, args []string) error {
	mapping := map[string]string{
		"start": "start", "up": "start",
		"stop": "stop", "down": "stop",
		"logs": "logs",
		"status": "status",
		"describe": "describe", "info": "describe",
		"reconfigure": "reconfigure", "reconfig": "reconfigure",
	}

	target, ok := mapping[sub]
	if !ok {
		return fmt.Errorf("unknown project command %q — try: start, stop, logs, status, describe, reconfigure", sub)
	}

	cmd := r.commands[target]
	return cmd.Execute(ctx, append([]string{project}, args...))
}

// HasCommand reports whether name is a registered command.
func (r *CommandRegistry) HasCommand(name string) bool {
	_, ok := r.commands[name]
	return ok
}

// isProjectName checks the registry to see if name is a known project.
func (r *CommandRegistry) isProjectName(name string) bool {
	if r.ctx.Registry == nil {
		return false
	}
	_, err := r.ctx.Registry.FindProject(name)
	return err == nil
}

// ShowUsage prints the top-level help.
func (r *CommandRegistry) ShowUsage() {
	fmt.Printf(`
%s█████╗ ████████╗███████╗███╗   ███╗██████╗  ██████╗%s
%s██╔══██╗╚══██╔══╝██╔════╝████╗ ████║██╔══██╗██╔═══██╗%s
%s███████║   ██║   █████╗  ██╔████╔██║██████╔╝██║   ██║%s
%s██╔══██║   ██║   ██╔══╝  ██║╚██╔╝██║██╔═══╝ ██║   ██║%s
%s██║  ██║   ██║   ███████╗██║ ╚═╝ ██║██║     ╚██████╔╝%s
%s╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝     ╚═╝╚═╝      ╚═════╝%s

`, ColorPurple, ColorReset, ColorPurple, ColorReset, ColorPurple, ColorReset,
		ColorPurple, ColorReset, ColorPurple, ColorReset, ColorPurple, ColorReset)

	fmt.Printf("%sUsage:%s atempo <command> [arguments]\n\n", ColorBold, ColorReset)

	// Print commands grouped in a fixed display order
	order := []string{
		"create",
		"start", "stop", "logs", "status",
		"projects", "describe",
		"add-service", "reconfigure", "remove",
		"doctor",
		"shell",
	}

	fmt.Printf("%sCommands:%s\n", ColorBold, ColorReset)
	for _, name := range order {
		if cmd, ok := r.commands[name]; ok {
			fmt.Printf("  %s%-16s%s %s\n", ColorCyan, name, ColorReset, cmd.Description())
		}
	}

	fmt.Printf(`
%sExamples:%s
  atempo create laravel my-app        Laravel (latest) in ./my-app/
  atempo create laravel:12 my-app     Laravel 12 in ./my-app/
  atempo create fastapi my-api        FastAPI in ./my-api/
  atempo create bun my-site           Bun in ./my-site/
  atempo create python my-worker      Plain Python in ./my-worker/
  atempo start my-app                 Start project containers
  atempo logs my-app                  Stream container logs
  atempo projects                     List all registered projects
  atempo <project> start              Shorthand: start a project by name

Run %satempo <command> --help%s for command-specific usage.
`, ColorBold, ColorReset, ColorYellow, ColorReset)
}

// IsHelpCommand reports whether arg is a help flag.
func IsHelpCommand(arg string) bool {
	switch strings.ToLower(arg) {
	case "help", "--help", "-h":
		return true
	}
	return false
}

// unknownCommand returns a helpful error with fuzzy-matched suggestions.
func (r *CommandRegistry) unknownCommand(name string) error {
	names := make([]string, 0, len(r.commands))
	for n := range r.commands {
		names = append(names, n)
	}

	suggestions := similarCommands(name, names)
	if len(suggestions) > 0 {
		return fmt.Errorf("unknown command %q\n\nDid you mean?\n  atempo %s\n\nRun 'atempo --help' for all commands.",
			name, strings.Join(suggestions, "\n  atempo "))
	}
	return fmt.Errorf("unknown command %q — run 'atempo --help' for available commands", name)
}

// similarCommands returns up to 3 commands with the highest similarity to input.
func similarCommands(input string, commands []string) []string {
	type scored struct {
		name  string
		score int
	}

	input = strings.ToLower(input)
	var results []scored

	for _, cmd := range commands {
		s := similarity(input, strings.ToLower(cmd))
		if s > 0 {
			results = append(results, scored{cmd, s})
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].score > results[j].score })

	out := make([]string, 0, 3)
	for i, r := range results {
		if i >= 3 {
			break
		}
		out = append(out, r.name)
	}
	return out
}

func similarity(a, b string) int {
	if a == b {
		return 100
	}
	score := 0
	if strings.HasPrefix(b, a) || strings.HasPrefix(a, b) {
		score += 50
	}
	if strings.Contains(b, a) || strings.Contains(a, b) {
		score += 30
	}
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] == b[i] {
			score += 5
		} else {
			break
		}
	}
	diff := len(a) - len(b)
	if diff < 0 {
		diff = -diff
	}
	if diff <= 2 {
		score += 10 - diff*3
	}
	return score
}

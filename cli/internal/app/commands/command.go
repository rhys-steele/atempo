package commands

import (
	"bytes"
	"context"
	"encoding/json"
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

// JSONCommand is an optional interface for commands that support --json output.
// The registry checks for this and routes --json invocations here automatically.
type JSONCommand interface {
	ExecuteJSON(ctx context.Context, args []string) error
}

// StreamCommand is an optional interface for commands that support --stream-json output.
// The registry checks for this and routes --stream-json invocations here automatically.
type StreamCommand interface {
	ExecuteStream(ctx context.Context, args []string) error
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

// emitEvent writes a single newline-delimited JSON event to stdout.
// Used by commands running in --stream-json mode.
func emitEvent(t, data string) {
	type ev struct {
		Type string `json:"type"`
		Data string `json:"data,omitempty"`
	}
	b, _ := json.Marshal(ev{Type: t, Data: data})
	fmt.Println(string(b))
}

// lineWriter is an io.Writer that splits writes into lines and calls emit for each one.
// Used to wrap installer/docker output as stream events.
type lineWriter struct {
	emit func(string)
	buf  []byte
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	for {
		idx := bytes.IndexByte(w.buf, '\n')
		if idx < 0 {
			break
		}
		line := strings.TrimRight(string(w.buf[:idx]), "\r")
		if line != "" {
			w.emit(line)
		}
		w.buf = w.buf[idx+1:]
	}
	return len(p), nil
}

func (w *lineWriter) flush() {
	if rest := strings.TrimRight(string(w.buf), "\r\n"); rest != "" {
		w.emit(rest)
	}
	w.buf = nil
}

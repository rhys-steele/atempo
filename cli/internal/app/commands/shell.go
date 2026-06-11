package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
)

type ShellCommand struct {
	*BaseCommand
	registry *CommandRegistry
}

func NewShellCommand(ctx *CommandContext, registry *CommandRegistry) *ShellCommand {
	cmd := &ShellCommand{
		BaseCommand: newBase(
			"shell",
			"Start interactive shell",
			"atempo shell",
			ctx,
		),
		registry: registry,
	}
	cmd.setHelp(`Usage: atempo shell

Starts an interactive Atempo shell. Type any Atempo command without the
'atempo' prefix. Type 'exit' or press Ctrl-D to quit.

Examples inside the shell:
  projects
  create fastapi my-api
  start my-app
  logs my-app
  exit
`)
	return cmd
}

func (c *ShellCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	fmt.Printf("\n%sAtempo Shell%s  (type 'exit' or Ctrl-D to quit)\n\n",
		ColorBold, ColorReset)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          c.prompt(),
		HistoryLimit:    200,
		AutoComplete:    c.completer(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return fmt.Errorf("start shell: %w", err)
	}
	defer rl.Close()

	for {
		// Refresh prompt each iteration so cd updates are reflected.
		rl.SetPrompt(c.prompt())

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if line == "" {
				break // Ctrl-C on empty line → exit
			}
			continue // Ctrl-C mid-line → clear and continue
		}
		if err == io.EOF {
			fmt.Println()
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			break
		}

		parts := strings.Fields(line)
		name := parts[0]
		cmdArgs := parts[1:]

		if IsHelpCommand(name) {
			c.registry.ShowUsage()
			continue
		}

		// cd is a shell builtin — must be handled directly so the process cwd updates.
		if name == "cd" {
			if err := changeDir(cmdArgs); err != nil {
				fmt.Fprintf(os.Stderr, "%scd:%s %v\n\n", ColorRed, ColorReset, err)
			}
			continue
		}

		// Atempo commands take priority; fall back to the system shell for everything else.
		if c.registry.HasCommand(name) {
			if err := c.registry.Execute(context.Background(), name, cmdArgs); err != nil {
				fmt.Fprintf(os.Stderr, "%sError:%s %v\n\n", ColorRed, ColorReset, err)
			}
		} else {
			if err := runShellCommand(line); err != nil {
				fmt.Fprintf(os.Stderr, "%s%v%s\n\n", ColorRed, err, ColorReset)
			}
		}
	}

	fmt.Printf("%sBye%s\n\n", ColorGray, ColorReset)
	return nil
}

// completer builds a tab completer. The first word completes against atempo
// command names and common shell commands. Subsequent words use filesystem
// path completion.
func (c *ShellCommand) completer() readline.AutoCompleter {
	// PcItemDynamic receives the FULL line typed so far. We extract the last
	// word ourselves to find the partial path, then return full-path candidates.
	// readline compares each candidate against the remaining text after the
	// matched command prefix, so candidates must preserve the user's format
	// (e.g. ~/foo not /Users/name/foo) to produce correct suffix diffs.
	fileDynamic := readline.PcItemDynamic(func(fullLine string) []string {
		var partial string
		if !strings.HasSuffix(fullLine, " ") {
			fields := strings.Fields(fullLine)
			if len(fields) > 0 {
				partial = fields[len(fields)-1]
			}
		}
		return fsCompletions(partial)
	})

	var items []readline.PrefixCompleterInterface
	for name := range c.registry.commands {
		items = append(items, readline.PcItem(name, fileDynamic))
	}
	for _, name := range []string{"cd", "ls", "cat", "grep", "git", "docker", "exit"} {
		items = append(items, readline.PcItem(name, fileDynamic))
	}
	return readline.NewPrefixCompleter(items...)
}

// fsCompletions returns filesystem candidates that match partial, preserving
// the user's original prefix format (e.g. ~/ stays as ~/, not expanded).
// Directories get a trailing "/" so the user can keep tabbing deeper.
func fsCompletions(partial string) []string {
	// Separate the directory prefix the user typed from the entry name prefix.
	typedDir, namePrefix := filepath.Split(partial)

	// Expand ~ for the actual filesystem read, but keep typedDir for candidates.
	readDir := typedDir
	switch {
	case partial == "~" || partial == "~/":
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		entries, err := os.ReadDir(home)
		if err != nil {
			return nil
		}
		var matches []string
		for _, e := range entries {
			candidate := "~/" + e.Name()
			if e.IsDir() {
				candidate += "/"
			}
			matches = append(matches, candidate)
		}
		return matches
	case strings.HasPrefix(typedDir, "~/"):
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		readDir = filepath.Join(home, typedDir[2:])
	case typedDir == "":
		readDir = "."
	}

	entries, err := os.ReadDir(readDir)
	if err != nil {
		return nil
	}

	var matches []string
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), namePrefix) {
			continue
		}
		candidate := typedDir + e.Name()
		if e.IsDir() {
			candidate += "/"
		}
		matches = append(matches, candidate)
	}
	return matches
}

// prompt builds a context-aware shell prompt.
func (c *ShellCommand) prompt() string {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("%satempo%s › ", ColorPurple, ColorReset)
	}

	dir := filepath.Base(cwd)
	branch := gitBranch(cwd)

	if branch != "" {
		return fmt.Sprintf("%satempo%s %s%s%s %s(%s)%s › ",
			ColorPurple, ColorReset,
			ColorBold, dir, ColorReset,
			ColorCyan, branch, ColorReset)
	}
	return fmt.Sprintf("%satempo%s %s%s%s › ",
		ColorPurple, ColorReset,
		ColorBold, dir, ColorReset)
}

// changeDir handles the `cd` builtin by updating the process working directory.
func changeDir(args []string) error {
	var target string
	switch len(args) {
	case 0:
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		target = home
	case 1:
		target = args[0]
	default:
		return fmt.Errorf("too many arguments")
	}

	// Expand ~ manually since exec doesn't go through a shell.
	if strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		target = filepath.Join(home, target[2:])
	}

	return os.Chdir(target)
}

// runShellCommand executes an arbitrary command string through the system shell,
// wiring stdin/stdout/stderr to the terminal so interactive programs work.
func runShellCommand(line string) error {
	var sh, flag string
	if runtime.GOOS == "windows" {
		sh, flag = "cmd", "/C"
	} else {
		sh = os.Getenv("SHELL")
		if sh == "" {
			sh = "/bin/sh"
		}
		flag = "-c"
	}

	cmd := exec.Command(sh, flag, line)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// gitBranch returns the current git branch name, or "" if not in a repo.
func gitBranch(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	branch := strings.TrimSpace(string(out))
	if branch == "HEAD" {
		return "" // detached HEAD — not useful to show
	}
	return branch
}

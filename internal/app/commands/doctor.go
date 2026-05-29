package commands

import (
	"context"
	"fmt"
	"runtime"

	"atempo/internal/dependency"
)

type DoctorCommand struct{ *BaseCommand }

func NewDoctorCommand(ctx *CommandContext) *DoctorCommand {
	cmd := &DoctorCommand{newBase(
		"doctor",
		"Check system dependencies",
		"atempo doctor [--verbose]",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo doctor [--verbose]

Checks that all system dependencies required by Atempo are installed and
accessible. Reports the version of each tool when --verbose is set.

Dependencies checked:
  Docker          (required) — container runtime
  Docker Compose  (required) — v2 plugin ('docker compose', not 'docker-compose')
  Git             (optional) — used to initialise project repos
`)
	return cmd
}

func (c *DoctorCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	verbose := hasFlag(args, "--verbose", "-v")

	fmt.Printf("\n%sSystem Diagnostics%s\n", ColorBold, ColorReset)
	fmt.Printf("Platform: %s/%s   Go: %s\n\n", runtime.GOOS, runtime.GOARCH, runtime.Version())

	mgr := dependency.NewManager()
	checks := mgr.CheckAll()

	allOK := true
	for _, ch := range checks {
		c.printCheck(ch, verbose)
		if ch.Dep.Required && !ch.OK() {
			allOK = false
		}
	}

	fmt.Println()
	if allOK {
		fmt.Printf("%s✓ All required dependencies are available.%s\n\n", ColorGreen, ColorReset)
		fmt.Printf("Next step: %satempo create laravel my-app%s\n\n", ColorCyan, ColorReset)
	} else {
		fmt.Printf("%s✗ One or more required dependencies are missing.%s\n", ColorRed, ColorReset)
		fmt.Printf("Fix the issues above then run %satempo doctor%s again.\n\n", ColorCyan, ColorReset)
	}

	return nil
}

func (c *DoctorCommand) printCheck(ch dependency.Check, verbose bool) {
	if ch.OK() {
		version := ""
		if verbose && ch.Version != "" {
			version = fmt.Sprintf(" %s(%s)%s", ColorGray, ch.Version, ColorReset)
		}
		tag := "required"
		if !ch.Dep.Required {
			tag = "optional"
		}
		fmt.Printf("  %s✓%s %-22s %s(%s)%s%s\n",
			ColorGreen, ColorReset,
			ch.Dep.Name,
			ColorGray, tag, ColorReset,
			version,
		)
		return
	}

	symbol := ColorRed + "✗" + ColorReset
	if !ch.Dep.Required {
		symbol = ColorYellow + "!" + ColorReset
	}

	fmt.Printf("  %s %-22s %s\n", symbol, ch.Dep.Name, ch.Message)
	if ch.Hint != "" {
		fmt.Printf("    %s→ %s%s\n", ColorGray, ch.Hint, ColorReset)
	}
}

func hasFlag(args []string, flags ...string) bool {
	for _, a := range args {
		for _, f := range flags {
			if a == f {
				return true
			}
		}
	}
	return false
}

package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DescribeCommand struct{ *BaseCommand }

func NewDescribeCommand(ctx *CommandContext) *DescribeCommand {
	cmd := &DescribeCommand{newBase(
		"describe",
		"Show project details",
		"atempo describe [project|path]",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo describe [project|path]

Shows details for an Atempo project: framework, status, port mappings,
services, and quick-action commands.

Arguments:
  project   Registered project name (from 'atempo projects')
  path      Explicit path to a project directory
  (none)    Uses the current directory

Examples:
  atempo describe                  Describe project in current directory
  atempo describe my-app           Describe registered project 'my-app'
  atempo describe /path/to/proj    Describe project at an explicit path
`)
	return cmd
}

func (c *DescribeCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	// Resolve which project to describe.
	name, projectPath, err := c.resolve(args)
	if err != nil {
		return err
	}

	// Pull extra metadata from atempo.json if present.
	meta := readAtempoJSON(projectPath)

	// Registry entry (may not exist for path-only resolution).
	var status, gitBranch string
	var ports []portRow
	var services []serviceRow

	if name != "" {
		if p, err := c.ctx.Registry.FindProject(name); err == nil {
			status = p.Status
			gitBranch = p.GitBranch
			for _, pt := range p.Ports {
				if pt.External != 0 {
					ports = append(ports, portRow{pt.Service, pt.External, pt.Internal})
				}
			}
			for _, svc := range p.Services {
				services = append(services, serviceRow{svc.Name, svc.Status})
			}
		}
	}

	// Display.
	fmt.Printf("\n%s%s%s\n", ColorBold, name, ColorReset)
	fmt.Println(strings.Repeat("─", 50))

	if meta.Framework != "" {
		fw := meta.Framework
		if meta.Version != "" {
			fw += " " + meta.Version
		}
		fmt.Printf("  Framework   %s\n", fw)
	}

	if status != "" {
		lbl, col := statusLabel(status)
		fmt.Printf("  Status      %s%s%s\n", col, lbl, ColorReset)
	}

	if gitBranch != "" {
		fmt.Printf("  Branch      %s\n", gitBranch)
	}

	fmt.Printf("  Path        %s\n", projectPath)

	if len(services) > 0 {
		fmt.Printf("\n%sServices%s\n", ColorBold, ColorReset)
		for _, s := range services {
			icon, col := serviceIcon(s.status)
			fmt.Printf("  %s%s%s %s\n", col, icon, ColorReset, s.name)
		}
	}

	if len(ports) > 0 {
		fmt.Printf("\n%sPorts%s\n", ColorBold, ColorReset)
		for _, p := range ports {
			fmt.Printf("  %-20s localhost:%d → %d\n", p.service, p.external, p.internal)
		}
	}

	fmt.Printf("\n%sQuick actions%s\n", ColorBold, ColorReset)
	if status == "stopped" || status == "" {
		fmt.Printf("  atempo start %s\n", name)
	} else {
		fmt.Printf("  atempo stop %s\n", name)
		fmt.Printf("  atempo logs %s\n", name)
	}
	fmt.Printf("  cd %s\n\n", projectPath)

	return nil
}

func (c *DescribeCommand) ExecuteJSON(_ context.Context, args []string) error {
	name, _, err := c.resolve(args)
	if err != nil {
		return err
	}
	p, err := c.ctx.Registry.FindProject(name)
	if err != nil {
		return fmt.Errorf("project %q not found: %w", name, err)
	}
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

// resolve returns (name, path) for the target project.
func (c *DescribeCommand) resolve(args []string) (string, string, error) {
	if len(args) == 0 {
		// Current directory.
		cwd, err := os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("get working directory: %w", err)
		}
		if !hasAtempoJSON(cwd) {
			return "", "", fmt.Errorf("no atempo.json found in current directory\nRun 'atempo create <framework>' to scaffold a project here")
		}
		return filepath.Base(cwd), cwd, nil
	}

	target := args[0]

	// Explicit path (absolute or starts with ./).
	if filepath.IsAbs(target) || strings.HasPrefix(target, "./") || strings.HasPrefix(target, "../") {
		abs, err := filepath.Abs(target)
		if err != nil {
			return "", "", fmt.Errorf("resolve path: %w", err)
		}
		if !hasAtempoJSON(abs) {
			return "", "", fmt.Errorf("no atempo.json found in %s", abs)
		}
		return filepath.Base(abs), abs, nil
	}

	// Registered project name.
	if p, err := c.ctx.Registry.FindProject(target); err == nil {
		return p.Name, p.Path, nil
	}

	// Fall back to treating it as a local directory name.
	if hasAtempoJSON(target) {
		abs, _ := filepath.Abs(target)
		return filepath.Base(abs), abs, nil
	}

	return "", "", fmt.Errorf("project %q not found in registry and no atempo.json at that path\nRun 'atempo projects' to see registered projects", target)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

type portRow struct {
	service          string
	external, internal int
}

type serviceRow struct{ name, status string }

type atempoMeta struct {
	Framework string `json:"framework"`
	Version   string `json:"version"`
	Name      string `json:"name"`
}

func readAtempoJSON(dir string) atempoMeta {
	data, err := os.ReadFile(filepath.Join(dir, "atempo.json"))
	if err != nil {
		return atempoMeta{}
	}
	var m atempoMeta
	_ = json.Unmarshal(data, &m)
	return m
}

func hasAtempoJSON(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "atempo.json"))
	return err == nil
}

func serviceIcon(status string) (string, string) {
	switch status {
	case "running", "healthy":
		return "✓", ColorGreen
	case "stopped":
		return "✗", ColorGray
	default:
		return "!", ColorYellow
	}
}

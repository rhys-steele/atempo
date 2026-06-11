package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"atempo/internal/docker"
)

type StatusCommand struct{ *BaseCommand }

func NewStatusCommand(ctx *CommandContext) *StatusCommand {
	cmd := &StatusCommand{newBase(
		"status",
		"Show service health",
		"atempo status [project]",
		ctx,
	)}
	cmd.setHelp(`Usage: atempo status [project]

Shows the live status of all Docker services in the project by running
'docker compose ps'.

Arguments:
  project   Registered project name or path (uses current directory if omitted)

Examples:
  atempo status             Status of project in current directory
  atempo status my-app      Status of registered project 'my-app'
`)
	return cmd
}

type serviceJSON struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (c *StatusCommand) Execute(_ context.Context, args []string) error {
	if c.handleHelp(args) {
		return nil
	}

	projectPath, projectName, err := resolveProject(args, c.ctx)
	if err != nil {
		return err
	}

	services, err := docker.PS(projectPath)
	if err != nil {
		return fmt.Errorf("status %s: %w", projectName, err)
	}

	fmt.Printf("\n%s%s%s\n", ColorBold, projectName, ColorReset)
	fmt.Println(strings.Repeat("─", 50))

	if len(services) == 0 {
		fmt.Printf("  %sNo containers running%s\n", ColorGray, ColorReset)
		fmt.Printf("  Run %satempo start %s%s to start the project\n\n",
			ColorCyan, projectName, ColorReset)
		return nil
	}

	for _, svc := range services {
		icon, col := serviceStatusIcon(svc.Status, svc.Health)
		label := formatServiceStatus(svc.Status, svc.Health)
		name := svc.Service
		if name == "" {
			name = svc.Name
		}
		fmt.Printf("  %s%s%s %-20s %s%s%s\n",
			col, icon, ColorReset,
			name,
			col, label, ColorReset,
		)
	}
	fmt.Println()
	return nil
}

func (c *StatusCommand) ExecuteJSON(_ context.Context, args []string) error {
	projectPath, _, err := resolveProject(args, c.ctx)
	if err != nil {
		return err
	}
	services, err := docker.PS(projectPath)
	if err != nil {
		return err
	}
	out := make([]serviceJSON, len(services))
	for i, s := range services {
		name := s.Service
		if name == "" {
			name = s.Name
		}
		status := s.Status
		if s.Health != "" && s.Health != "none" {
			status = s.Health
		}
		out[i] = serviceJSON{Name: name, Status: status}
	}
	b, _ := json.Marshal(out)
	fmt.Println(string(b))
	return nil
}

func serviceStatusIcon(status, health string) (string, string) {
	if health == "healthy" {
		return "✓", ColorGreen
	}
	if health == "unhealthy" {
		return "✗", ColorRed
	}
	switch status {
	case "running":
		return "✓", ColorGreen
	case "exited", "stopped":
		return "✗", ColorRed
	default:
		return "○", ColorYellow
	}
}

func formatServiceStatus(status, health string) string {
	if health != "" && health != "none" {
		return status + " (" + health + ")"
	}
	return status
}

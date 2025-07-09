package commands

import (
	"context"
	"fmt"

	"atempo/internal/compose"
	"atempo/internal/utils"
)

// AddServiceCommand adds a predefined service to a project
type AddServiceCommand struct {
	*BaseCommand
}

// NewAddServiceCommand creates a new add-service command
func NewAddServiceCommand(ctx *CommandContext) *AddServiceCommand {
	return &AddServiceCommand{
		BaseCommand: NewBaseCommand(
			"add-service",
			utils.GetStandardDescription("add-service"),
			utils.CreateStandardUsage("add-service", utils.PatternWithRequiredArgs, "<service_type>", "[project]"),
			ctx,
		),
	}
}

// Execute runs the add-service command
func (c *AddServiceCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: atempo add-service <service_type> [project]")
		fmt.Println("\nAvailable services:")
		for _, service := range compose.ListPredefinedServices() {
			fmt.Printf("  %s\n", service)
		}
		return fmt.Errorf("service type required")
	}

	serviceType := args[0]
	// Resolve project path from arguments (service is args[0], project is args[1])
	projectArgs := []string{}
	if len(args) > 1 {
		projectArgs = args[1:]
	}
	resolution, err := utils.ResolveProjectPathFromArgs(projectArgs)
	if err != nil {
		return err
	}
	projectPath := resolution.Path

	fmt.Printf("→ Adding %s service to project...\n", serviceType)

	if err := compose.AddPredefinedService(projectPath, serviceType); err != nil {
		return fmt.Errorf("failed to add service: %w", err)
	}

	fmt.Printf("✓ %s service added to atempo.json\n", serviceType)
	fmt.Println("Run 'atempo reconfigure' to update docker-compose.yml")
	return nil
}
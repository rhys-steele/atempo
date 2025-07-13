package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"atempo/internal/registry"
	"atempo/internal/utils"
)

// executeProjectCommand handles project-specific commands
func (r *CommandRegistry) executeProjectCommand(ctx context.Context, projectName, command string, args []string) error {
	// Map project commands to existing global commands
	switch command {
	case "up", "start":
		// Execute docker up for this project
		dockerCmd := r.commands["docker"]
		return dockerCmd.Execute(ctx, append([]string{"up", projectName}, args...))

	case "down", "stop":
		// Execute docker down for this project
		dockerCmd := r.commands["docker"]
		return dockerCmd.Execute(ctx, append([]string{"down", projectName}, args...))

	case "status":
		// Execute projects command for this project (replaces old status command)
		projectsCmd := r.commands["projects"]
		return projectsCmd.Execute(ctx, append([]string{projectName}, args...))

	case "logs":
		// Execute logs for this project
		logsCmd := r.commands["logs"]
		return logsCmd.Execute(ctx, append([]string{projectName}, args...))

	case "describe", "info":
		// Execute describe for this project
		describeCmd := r.commands["describe"]
		return describeCmd.Execute(ctx, append([]string{projectName}, args...))

	case "shell", "bash", "exec":
		// Execute shell access for this project
		dockerCmd := r.commands["docker"]
		return dockerCmd.Execute(ctx, append([]string{"bash", projectName}, args...))

	case "reconfigure", "reconfig":
		// Execute reconfigure for this project
		reconfigCmd := r.commands["reconfigure"]
		return reconfigCmd.Execute(ctx, append([]string{projectName}, args...))

	case "code":
		// Open project in VS Code
		return r.openProjectInVSCode(projectName)

	case "cd":
		// Change directory to project (note: this only works in shell session)
		return r.changeToProjectDirectory(projectName)

	case "delete", "remove":
		// Delete project files and remove from registry
		return r.deleteProject(projectName)

	case "open":
		// Open project or specific service in browser
		return r.openProjectInBrowser(projectName, args)

	case "audit":
		// Execute audit for this project
		auditCmd := r.commands["audit"]
		
		// Load registry to get project details
		reg, err := registry.LoadRegistry()
		if err != nil {
			return fmt.Errorf("failed to load project registry: %w", err)
		}

		// Find the specified project
		project, err := reg.FindProject(projectName)
		if err != nil {
			return fmt.Errorf("project '%s' not found in registry", projectName)
		}

		// Change to project directory temporarily
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		
		if err := os.Chdir(project.Path); err != nil {
			return fmt.Errorf("failed to change to project directory: %w", err)
		}

		return auditCmd.Execute(ctx, args)

	default:
		return fmt.Errorf("unknown project command: %s. Available: up, down, status, logs, describe, shell, reconfigure, code, cd, delete, open, audit", command)
	}
}

// openProjectInVSCode opens the specified project in VS Code
func (r *CommandRegistry) openProjectInVSCode(projectName string) error {
	// Load registry to get project details
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}

	// Find the project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project not found: %s", projectName)
	}

	// Check if VS Code is installed
	codePath, err := exec.LookPath("code")
	if err != nil {
		return fmt.Errorf("VS Code (code command) not found. Please install VS Code and ensure 'code' is in your PATH")
	}

	// Open the project directory in VS Code
	cmd := exec.Command(codePath, project.Path)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open project in VS Code: %w", err)
	}

	// Don't wait for VS Code to close - let it run in background
	go func() {
		cmd.Wait()
	}()

	return nil
}

// changeToProjectDirectory changes the current working directory to the project path
func (r *CommandRegistry) changeToProjectDirectory(projectName string) error {
	// Load registry to get project details
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}

	// Find the specified project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project '%s' not found in registry", projectName)
	}

	// Change to the project directory
	if err := os.Chdir(project.Path); err != nil {
		return fmt.Errorf("failed to change to project directory %s: %w", project.Path, err)
	}

	fmt.Printf("  ⎿ Changed to project directory: %s\n", project.Path)
	return nil
}

// deleteProject deletes project files and removes from registry
func (r *CommandRegistry) deleteProject(projectName string) error {
	// Load registry to get project details
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load project registry: %w", err)
	}

	// Find the specified project
	project, err := reg.FindProject(projectName)
	if err != nil {
		return fmt.Errorf("project '%s' not found in registry", projectName)
	}

	// Show confirmation prompt with detailed info
	fmt.Printf("! Are you sure you want to delete project '%s'?\n\n", projectName)
	fmt.Printf("  Path: %s\n", project.Path)
	fmt.Printf("  Framework: %s %s\n", project.Framework, project.Version)
	fmt.Printf("  Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("\n  This will:\n")
	fmt.Printf("  • Move project files to Trash\n")
	fmt.Printf("  • Remove project from atempo registry\n")
	fmt.Printf("  • This action cannot be undone!\n\n")
	fmt.Print("Type 'delete' to confirm, or anything else to cancel: ")

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "delete" {
		fmt.Println("✗ Cancelled - project not deleted.")
		return nil
	}

	// Move to trash using macOS 'trash' command or fallback
	err = utils.MoveToTrash(project.Path)
	if err != nil {
		return fmt.Errorf("failed to move project to trash: %w", err)
	}

	// Remove project from registry
	err = reg.RemoveProject(projectName)
	if err != nil {
		return fmt.Errorf("failed to remove project from registry: %w", err)
	}

	ShowSuccess(fmt.Sprintf("Project '%s' deleted successfully!", projectName), "Files moved to Trash, removed from registry")

	return nil
}

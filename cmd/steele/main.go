package main

import (
	"fmt"
	"os"
	"strings"

	"steele/internal/scaffold"
)

// main is the entry point for the Steele CLI.
// Steele is a simple command-line tool for scaffolding developer projects
// using AI-first principles and an MCP-ready context architecture.
func main() {
	// Ensure correct usage: steele start <framework>:<version>
	if len(os.Args) < 3 || os.Args[1] != "start" {
		fmt.Println("Usage: steele start <framework>:<version>")
		os.Exit(1)
	}

	// Extract framework and version
	arg := os.Args[2]
	parts := strings.Split(arg, ":")
	if len(parts) != 2 {
		fmt.Println("Error: expected format <framework>:<version>")
		os.Exit(1)
	}

	framework := parts[0]
	version := parts[1]

	// Trigger the scaffold process
	err := scaffold.Run(framework, version)
	if err != nil {
		fmt.Printf("Scaffold error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Project scaffolding complete.")
}

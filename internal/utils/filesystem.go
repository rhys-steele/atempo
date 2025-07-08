package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// MoveToTrash moves a directory or file to the trash
func MoveToTrash(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Try using 'trash' command if available (brew install trash)
	if _, err := exec.LookPath("trash"); err == nil {
		cmd := exec.Command("trash", path)
		return cmd.Run()
	}

	// Fallback: use 'mv' to move to ~/.Trash (macOS)
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	trashDir := filepath.Join(home, ".Trash")

	// Ensure .Trash directory exists
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("failed to create .Trash directory: %w", err)
	}

	basename := filepath.Base(path)

	// Generate unique name if file already exists in trash
	targetPath := filepath.Join(trashDir, basename)
	if _, err := os.Stat(targetPath); err == nil {
		timestamp := time.Now().Format("20060102-150405")
		targetPath = filepath.Join(trashDir, fmt.Sprintf("%s-%s", basename, timestamp))
	}

	return os.Rename(path, targetPath)
}

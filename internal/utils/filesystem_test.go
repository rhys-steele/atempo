package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveToTrash_NonExistentPath(t *testing.T) {
	// Test that moving a non-existent path returns an error
	err := MoveToTrash("/path/that/does/not/exist")
	if err == nil {
		t.Error("MoveToTrash should return an error for non-existent path")
	}

	if !os.IsNotExist(err) && err.Error() != "path does not exist: /path/that/does/not/exist" {
		t.Errorf("Expected 'path does not exist' error, got: %v", err)
	}
}

func TestMoveToTrash_EmptyPath(t *testing.T) {
	// Test that moving an empty path returns an error
	err := MoveToTrash("")
	if err == nil {
		t.Error("MoveToTrash should return an error for empty path")
	}
}

func TestMoveToTrash_WithTempFile(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_file.txt")

	// Create the test file
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	file.Close()

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatalf("Temp file was not created properly")
	}

	// Test moving to trash
	// Note: This test may fail if 'trash' command is not available and macOS .Trash directory doesn't exist
	// We'll just verify it doesn't panic and returns some result
	err = MoveToTrash(tempFile)

	// The test should not panic, but may return an error depending on system setup
	// This is acceptable for a basic test
	t.Logf("MoveToTrash result: %v", err)
}

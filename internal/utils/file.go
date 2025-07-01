package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDir recursively copies a directory from src to dst.
// It creates the destination directory if it doesn't exist.
func CopyDir(src, dst string) error {
	// Get info about the source
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}

	// Create the destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Read the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := CopyDir(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to copy subdirectory %s: %w", entry.Name(), err)
			}
		} else {
			// Copy files
			if err := CopyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// CopyFile copies a single file from src to dst.
func CopyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy file contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Set file permissions
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// FileExists checks if a file or directory exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CompareVersions compares two semantic version strings.
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
// Supports simple semantic versioning (e.g., "10.0", "11.5", "5.1")
func CompareVersions(v1, v2 string) int {
	// Simple version comparison for major.minor format
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	
	// Pad shorter version with zeros
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}
	
	for len(parts1) < maxLen {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < maxLen {
		parts2 = append(parts2, "0")
	}
	
	for i := 0; i < maxLen; i++ {
		num1 := ParseVersionPart(parts1[i])
		num2 := ParseVersionPart(parts2[i])
		
		if num1 < num2 {
			return -1
		}
		if num1 > num2 {
			return 1
		}
	}
	
	return 0
}

// ParseVersionPart converts a version part string to an integer
func ParseVersionPart(part string) int {
	// Remove any non-numeric characters and convert to int
	numStr := ""
	for _, r := range part {
		if r >= '0' && r <= '9' {
			numStr += string(r)
		}
	}
	
	if numStr == "" {
		return 0
	}
	
	// Simple integer parsing
	result := 0
	for _, r := range numStr {
		result = result*10 + int(r-'0')
	}
	
	return result
}
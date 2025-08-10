package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// MaxFileSize limits file operations to 10MB
	MaxFileSize = 10 * 1024 * 1024 // 10MB
	// MaxPathLength prevents extremely long paths
	MaxPathLength = 4096
)

// validatePath checks if a path is safe for file operations
func validatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if len(path) > MaxPathLength {
		return fmt.Errorf("path too long (max %d characters)", MaxPathLength)
	}

	// Clean the path to resolve any .., ., etc.
	cleanPath := filepath.Clean(path)
	
	// Get absolute path to check for directory traversal
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	absWd, err := filepath.Abs(wd)
	if err != nil {
		return fmt.Errorf("failed to resolve working directory: %w", err)
	}

	// Check if the path tries to escape the working directory
	if !strings.HasPrefix(absPath, absWd) {
		return fmt.Errorf("path '%s' attempts to access files outside working directory", path)
	}

	// Prevent access to sensitive files/directories
	restrictedPaths := []string{
		"/etc/passwd",
		"/etc/shadow",
		"/proc",
		"/sys",
		".ssh",
		".git",
		"node_modules",
	}

	lowerPath := strings.ToLower(absPath)
	for _, restricted := range restrictedPaths {
		if strings.Contains(lowerPath, strings.ToLower(restricted)) {
			return fmt.Errorf("access to '%s' is restricted", path)
		}
	}

	return nil
}

// validateFileSize checks if a file size is within acceptable limits
func validateFileSize(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		// If file doesn't exist, it's okay for write operations
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if info.Size() > MaxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", 
			info.Size(), MaxFileSize)
	}

	return nil
}

// validateWriteOperation checks if a write operation is safe
func validateWriteOperation(path string, contentSize int) error {
	if err := validatePath(path); err != nil {
		return err
	}

	if contentSize > MaxFileSize {
		return fmt.Errorf("content size %d bytes exceeds maximum allowed size of %d bytes", 
			contentSize, MaxFileSize)
	}

	// Check if directory is writable
	dir := filepath.Dir(path)
	if dir != "." {
		// Try to create directory if it doesn't exist
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create directory '%s': %w", dir, err)
		}
	}

	return nil
}

// validateReadOperation checks if a read operation is safe
func validateReadOperation(path string) error {
	if err := validatePath(path); err != nil {
		return err
	}

	if err := validateFileSize(path); err != nil {
		return err
	}

	// Check if file is readable
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file '%s' does not exist", path)
		}
		return fmt.Errorf("cannot access file '%s': %w", path, err)
	}

	return nil
}
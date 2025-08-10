package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vivesm/GOSS-CLI/agentic-cli/openai"
)

// CreateFilesystemTools returns a slice of filesystem MCP tools
func CreateFilesystemTools() []openai.Tool {
	return []openai.Tool{
		{
			Type: "function",
			Function: openai.ToolFunction{
				Name:        "read_file",
				Description: "Read the contents of a file",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the file to read",
						},
					},
					"required": []string{"path"},
				},
				Handler: readFileHandler,
			},
		},
		{
			Type: "function",
			Function: openai.ToolFunction{
				Name:        "write_file",
				Description: "Write content to a file",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the file to write",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "Content to write to the file",
						},
					},
					"required": []string{"path", "content"},
				},
				Handler: writeFileHandler,
			},
		},
		{
			Type: "function",
			Function: openai.ToolFunction{
				Name:        "list_directory",
				Description: "List files and directories in a given path",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the directory to list",
						},
					},
					"required": []string{"path"},
				},
				Handler: listDirectoryHandler,
			},
		},
		{
			Type: "function",
			Function: openai.ToolFunction{
				Name:        "search_files",
				Description: "Search for files matching a pattern",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Directory path to search in",
						},
						"pattern": map[string]interface{}{
							"type":        "string",
							"description": "File name pattern to search for",
						},
					},
					"required": []string{"path", "pattern"},
				},
				Handler: searchFilesHandler,
			},
		},
		{
			Type: "function",
			Function: openai.ToolFunction{
				Name:        "create_directory",
				Description: "Create a new directory",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the directory to create",
						},
					},
					"required": []string{"path"},
				},
				Handler: createDirectoryHandler,
			},
		},
	}
}

func readFileHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path must be a string")
	}

	// Security validation
	if err := validateReadOperation(path); err != nil {
		return "", fmt.Errorf("security validation failed: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return string(content), nil
}

func writeFileHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path must be a string")
	}

	content, ok := args["content"].(string)
	if !ok {
		return "", fmt.Errorf("content must be a string")
	}

	// Security validation
	if err := validateWriteOperation(path, len(content)); err != nil {
		return "", fmt.Errorf("security validation failed: %w", err)
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return fmt.Sprintf("Successfully wrote %d bytes to %s", len(content), path), nil
}

func listDirectoryHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path must be a string")
	}

	// Security validation  
	if err := validatePath(path); err != nil {
		return "", fmt.Errorf("security validation failed: %w", err)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Contents of %s:\n", path))

	for _, entry := range entries {
		if entry.IsDir() {
			result.WriteString(fmt.Sprintf("[DIR]  %s/\n", entry.Name()))
		} else {
			info, err := entry.Info()
			if err != nil {
				result.WriteString(fmt.Sprintf("[FILE] %s (size unknown)\n", entry.Name()))
			} else {
				result.WriteString(fmt.Sprintf("[FILE] %s (%d bytes)\n", entry.Name(), info.Size()))
			}
		}
	}

	return result.String(), nil
}

func searchFilesHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	basePath, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path must be a string")
	}

	pattern, ok := args["pattern"].(string)
	if !ok {
		return "", fmt.Errorf("pattern must be a string")
	}

	// Security validation
	if err := validatePath(basePath); err != nil {
		return "", fmt.Errorf("security validation failed: %w", err)
	}

	var matches []string

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking even if we can't access some files
		}

		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to search files: %w", err)
	}

	if len(matches) == 0 {
		return fmt.Sprintf("No files matching pattern '%s' found in %s", pattern, basePath), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d files matching pattern '%s':\n", len(matches), pattern))
	for _, match := range matches {
		result.WriteString(fmt.Sprintf("- %s\n", match))
	}

	return result.String(), nil
}

func createDirectoryHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path must be a string")
	}

	// Security validation
	if err := validatePath(path); err != nil {
		return "", fmt.Errorf("security validation failed: %w", err)
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	return fmt.Sprintf("Successfully created directory: %s", path), nil
}

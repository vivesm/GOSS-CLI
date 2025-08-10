package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config contains the application configuration data and methods.
// This consolidates the previous Configuration and ApplicationData structs.
type Config struct {
	filePath      string                 // Path to the configuration file
	SystemPrompts map[string]string      `json:"SystemPrompts"`
	History       map[string]interface{} `json:"History"`
}

// NewConfig returns a new Config from a JSON file.
// If the file doesn't exist, it creates a default configuration.
func NewConfig(filePath string) (*Config, error) {
	config := &Config{
		filePath:      filePath,
		SystemPrompts: getDefaultSystemPrompts(),
		History:       make(map[string]interface{}),
	}

	// Try to load existing configuration
	if err := config.Load(); err != nil {
		// If file doesn't exist, create it with defaults
		if os.IsNotExist(err) {
			if saveErr := config.Save(); saveErr != nil {
				return nil, fmt.Errorf("failed to create default config file: %w", saveErr)
			}
			return config, nil
		}
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate the loaded configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Load reads and parses the configuration from the file.
func (c *Config) Load() error {
	file, err := os.Open(c.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}

// Save writes the configuration to the file atomically.
// This replaces the previous Flush/Write methods.
func (c *Config) Save() error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(c.filePath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write to temporary file first (atomic operation)
	tempFile := c.filePath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary config file: %w", err)
	}
	defer func() {
		file.Close()
		os.Remove(tempFile) // Clean up temp file on error
	}()

	// Encode JSON with formatting
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode configuration: %w", err)
	}

	// Sync to disk and close
	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync config file: %w", err)
	}
	file.Close()

	// Atomically replace the original file
	if err := os.Rename(tempFile, c.filePath); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// Write is an alias for Save for backward compatibility.
func (c *Config) Write() error {
	return c.Save()
}

// Backup creates a backup of the current configuration file.
func (c *Config) Backup() error {
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		return nil // No file to backup
	}

	backupPath := c.filePath + ".backup"
	return copyFile(c.filePath, backupPath)
}

// Validate checks the configuration for validity.
func (c *Config) Validate() error {
	if err := c.ValidateSystemPrompts(); err != nil {
		return err
	}
	return c.ValidateHistory()
}

// ValidateSystemPrompts ensures system prompts are valid.
func (c *Config) ValidateSystemPrompts() error {
	if c.SystemPrompts == nil {
		return fmt.Errorf("SystemPrompts cannot be nil")
	}

	if len(c.SystemPrompts) == 0 {
		return fmt.Errorf("at least one system prompt is required")
	}

	for name, prompt := range c.SystemPrompts {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("system prompt name cannot be empty")
		}
		if strings.TrimSpace(prompt) == "" {
			return fmt.Errorf("system prompt '%s' cannot be empty", name)
		}
	}

	return nil
}

// ValidateHistory ensures history data is valid.
func (c *Config) ValidateHistory() error {
	if c.History == nil {
		c.History = make(map[string]interface{}) // Auto-fix nil history
		return nil
	}

	for key, value := range c.History {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("history key cannot be empty")
		}

		// Ensure history values can be marshaled to JSON
		if _, err := json.Marshal(value); err != nil {
			return fmt.Errorf("invalid history entry '%s': %w", key, err)
		}
	}

	return nil
}

// getDefaultSystemPrompts returns the default system prompts.
func getDefaultSystemPrompts() map[string]string {
	return map[string]string{
		"Assistant":  "You are a helpful AI assistant with access to filesystem and web search tools. When users ask for current information (like weather, news, prices), ALWAYS use web search first. For web searches: try multiple specific queries if the first doesn't give direct answers. For weather queries, search for 'current temperature [location]' or '[location] weather now'. Extract information from search result descriptions and be persistent - if one search doesn't work, try different keywords.",
		"Developer":  "You are an expert software developer with access to filesystem and web search tools. Help with coding tasks, debugging, and software development questions. Use file operations to read code, search for patterns, and create or modify files as needed. When users ask about current technologies, APIs, or documentation, use web search to get the latest information.",
		"Researcher": "You are a research assistant with web search capabilities. Help find information online and use filesystem tools to organize research findings into files. Be thorough with web searches - try multiple search queries with different keywords to get comprehensive results. Always verify information by searching from multiple angles.",
		"Writer":     "You are a writing assistant that can help with document creation and editing. Use filesystem tools to read, write, and organize documents. When writing about current events or factual information, use web search to verify accuracy and get up-to-date details.",
	}
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

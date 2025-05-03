package config

import (
	"fmt"
	"path/filepath"
)

// Config holds the MCP server configuration
type Config struct {
	AllowedDirectories []string
}

// NewConfig creates a new Config with validated allowed directories
func NewConfig(dirs []string) (*Config, error) {
	if len(dirs) == 0 {
		return nil, fmt.Errorf("at least one allowed directory must be provided")
	}

	// Convert all paths to absolute and clean them
	cleanDirs := make([]string, len(dirs))
	for i, dir := range dirs {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			return nil, fmt.Errorf("invalid directory path %q: %w", dir, err)
		}
		cleanDirs[i] = filepath.Clean(absPath)
	}

	return &Config{
		AllowedDirectories: cleanDirs,
	}, nil
}

// IsPathAllowed checks if a given path is within any of the allowed directories
func (c *Config) IsPathAllowed(path string) bool {
	if !filepath.IsAbs(path) {
		return false
	}

	cleanPath := filepath.Clean(path)
	for _, dir := range c.AllowedDirectories {
		// Use HasPrefix with an extra separator to ensure we're matching complete path segments
		if filepath.HasPrefix(cleanPath, dir+string(filepath.Separator)) || cleanPath == dir {
			return true
		}
	}
	return false
}

package config

import (
	"fmt"
	"path/filepath"
)

type Config struct {
	AllowedDirectories []string
}

func NewConfig(dirs []string) (*Config, error) {
	if len(dirs) == 0 {
		return nil, fmt.Errorf("at least one allowed directory must be provided")
	}

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

func (c *Config) IsPathAllowed(path string) bool {
	if !filepath.IsAbs(path) {
		return false
	}

	cleanPath := filepath.Clean(path)
	for _, dir := range c.AllowedDirectories {
		if cleanPath == dir {
			return true
		}

		dirWithSep := dir + string(filepath.Separator)
		rel := filepath.Clean(cleanPath[len(dirWithSep):])
		if len(cleanPath) > len(dirWithSep) &&
			cleanPath[:len(dirWithSep)] == dirWithSep &&
			!filepath.IsAbs(rel) {
			return true
		}
	}
	return false
}

package tool

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielsuguimoto/go-filesystem-mcp/config"
	"github.com/mark3labs/mcp-go/mcp"
)

type ListDirectoryHandler struct {
	cfg *config.Config
}

func NewListDirectoryHandler(cfg *config.Config) *ListDirectoryHandler {
	return &ListDirectoryHandler{cfg: cfg}
}

func (h *ListDirectoryHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("invalid arguments for list_directory")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	if !h.cfg.IsPathAllowed(absPath) {
		return nil, fmt.Errorf("path %q is outside of allowed directories", absPath)
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var formattedEntries []string
	for _, entry := range entries {
		prefix := "[FILE]"
		if entry.IsDir() {
			prefix = "[DIR]"
		}
		formattedEntries = append(formattedEntries, fmt.Sprintf("%s %s", prefix, entry.Name()))
	}

	return mcp.NewToolResultText(strings.Join(formattedEntries, "\n")), nil
}

package tool

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielsuguimoto/go-filesystem-mcp/config"
	"github.com/mark3labs/mcp-go/mcp"
)

type CreateDirectoryHandler struct {
	cfg *config.Config
}

func NewCreateDirectoryHandler(cfg *config.Config) *CreateDirectoryHandler {
	return &CreateDirectoryHandler{cfg: cfg}
}

func (h *CreateDirectoryHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'path' argument")
	}

	// Convert to absolute path if not already
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Check if path is allowed
	if !h.cfg.IsPathAllowed(absPath) {
		return nil, fmt.Errorf("path %q is outside of allowed directories", absPath)
	}

	if _, err := os.Stat(absPath); err == nil {
		return mcp.NewToolResultText(fmt.Sprintf("Directory created: %s", absPath)), nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to check directory: %w", err)
	}

	recursive := "true"
	if recursiveVal, ok := request.Params.Arguments["recursive"].(string); ok {
		recursive = recursiveVal
	}

	if recursive == "true" {
		err = os.MkdirAll(absPath, os.FileMode(0777))
	} else {
		err = os.Mkdir(absPath, os.FileMode(0777))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Directory created: %s", absPath)), nil
}

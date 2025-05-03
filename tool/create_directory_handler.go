package tool

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
)

func CreateDirectoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'path' argument")
	}

	if _, err := os.Stat(path); err == nil {
		return mcp.NewToolResultText(fmt.Sprintf("Directory created: %s", path)), nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to check directory: %w", err)
	}

	recursive := "true"
	if recursiveVal, ok := request.Params.Arguments["recursive"].(string); ok {
		recursive = recursiveVal
	}

	var err error
	if recursive == "true" {
		err = os.MkdirAll(path, os.FileMode(0777))
	} else {
		err = os.Mkdir(path, os.FileMode(0777))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Directory created: %s", path)), nil
}

package tool

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielsuguimoto/go-filesystem-mcp/config"
	"github.com/mark3labs/mcp-go/mcp"
)

type ReadFileHandler struct {
	cfg *config.Config
}

func NewReadFileHandler(cfg *config.Config) *ReadFileHandler {
	return &ReadFileHandler{cfg: cfg}
}

func (h *ReadFileHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("invalid arguments for read_file")
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

	var fromLine, toLine int64 = 0, -1

	if from, ok := request.Params.Arguments["from"].(float64); ok {
		fromLine = int64(from)
	}

	if to, ok := request.Params.Arguments["to"].(float64); ok {
		toLine = int64(to)
	}

	if fromLine > 0 || toLine >= 0 {
		file, err := os.Open(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var lines []string
		var currentLine int64 = 0
		start := fromLine
		end := toLine
		if start < 0 {
			start = 0
		}
		if end < 0 {
			end = 1<<63 - 1
		}
		for scanner.Scan() {
			if currentLine >= start && currentLine <= end {
				lines = append(lines, scanner.Text())
			}
			currentLine++
			if currentLine > end {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return mcp.NewToolResultText(string(content)), nil
}

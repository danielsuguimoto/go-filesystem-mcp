package tool

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func ReadFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("invalid arguments for read_file")
	}

	var fromLine, toLine int64 = 0, -1

	if from, ok := request.Params.Arguments["from"].(float64); ok {
		fromLine = int64(from)
	}

	if to, ok := request.Params.Arguments["to"].(float64); ok {
		toLine = int64(to)
	}

	if fromLine > 0 || toLine >= 0 {
		file, err := os.Open(path)
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
			if currentLine > end {
				break
			}
			currentLine++
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return mcp.NewToolResultText(string(content)), nil
}

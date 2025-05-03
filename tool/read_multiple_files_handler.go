package tool

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)

func ReadMultipleFilesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pathsInterface, ok := request.Params.Arguments["paths"].([]interface{})
	if !ok {
		return nil, errors.New("missing required argument: paths")
	}

	paths := make([]string, len(pathsInterface))
	errChan := make(chan error, len(pathsInterface))
	var castingWg sync.WaitGroup

	// Process paths concurrently
	for i, p := range pathsInterface {
		castingWg.Add(1)
		go func(index int, path interface{}) {
			defer castingWg.Done()
			str, ok := path.(string)
			if !ok {
				errChan <- fmt.Errorf("path at index %d must be a string, got %T", index, path)
				return
			}
			paths[index] = str
		}(i, p)
	}

	// Wait for all goroutines to finish
	castingWg.Wait()
	close(errChan)

	// Check for any errors
	if err := <-errChan; err != nil {
		return nil, err
	}

	results := make([]string, len(paths))
	type result struct {
		idx    int
		output string
	}
	resCh := make(chan result, len(paths))
	var readingWg sync.WaitGroup

	for i, path := range paths {
		readingWg.Add(1)
		go func(idx int, p string) {
			defer readingWg.Done()
			content, err := os.ReadFile(p)
			if err != nil {
				resCh <- result{idx, formatResult(p, err.Error(), "")}
				return
			}
			resCh <- result{idx, formatResult(p, "", string(content))}
		}(i, path)
	}

	readingWg.Wait()
	close(resCh)

	for r := range resCh {
		results[r.idx] = r.output
	}

	return mcp.NewToolResultText(joinResults(results)), nil
}

func formatResult(path string, errMsg string, content string) string {
	if errMsg != "" {
		return fmt.Sprintf("%v: Error – %s\n", path, errMsg)
	}
	return fmt.Sprintf("%v:\n%s\n", path, content)
}

func joinResults(results []string) string {
	return strings.Join(results, "\n---\n")
}

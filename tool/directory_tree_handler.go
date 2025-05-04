package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/danielsuguimoto/go-filesystem-mcp/config"
	"github.com/mark3labs/mcp-go/mcp"
)

type TreeEntry struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Children []TreeEntry `json:"children,omitempty"`
}

type DirectoryTreeHandler struct {
	cfg *config.Config
}

func NewDirectoryTreeHandler(cfg *config.Config) *DirectoryTreeHandler {
	return &DirectoryTreeHandler{cfg: cfg}
}

func (h *DirectoryTreeHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("invalid arguments for directory_tree")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	if !h.cfg.IsPathAllowed(absPath) {
		return nil, fmt.Errorf("path %q is outside of allowed directories", absPath)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path %q is not a directory", absPath)
	}

	tree, err := h.buildTree(ctx, absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build directory tree: %w", err)
	}

	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tree to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func (h *DirectoryTreeHandler) buildTree(ctx context.Context, dirPath string) ([]TreeEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	result := make([]TreeEntry, 0, len(entries))

	resultChan := make(chan TreeEntry, len(entries))
	errChan := make(chan error, len(entries))

	for _, entry := range entries {
		wg.Add(1)
		go func(entry os.DirEntry) {
			defer wg.Done()

			entryPath := filepath.Join(dirPath, entry.Name())
			treeEntry := TreeEntry{
				Name: entry.Name(),
				Type: "file",
			}

			if entry.IsDir() {
				treeEntry.Type = "directory"
				children, err := h.buildTree(ctx, entryPath)
				if err != nil {
					errChan <- err
					return
				}
				treeEntry.Children = children
			}

			resultChan <- treeEntry
		}(entry)
	}
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	for {
		select {
		case err, ok := <-errChan:
			if !ok {
				continue
			}
			return nil, err
		case entry, ok := <-resultChan:
			if !ok {
				return result, nil
			}
			mu.Lock()
			result = append(result, entry)
			mu.Unlock()
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

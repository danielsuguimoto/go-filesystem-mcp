package main

import (
	"fmt"

	"github.com/danielsuguimoto/go-filesystem-mcp/tool"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Filesystem MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	addTools(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func addTools(s *server.MCPServer) {
	addReadFileTool(s)
	addReadMultipleFilesTool(s)
}

func addReadFileTool(s *server.MCPServer) {
	readFileTool := mcp.NewTool(
		"read_file",
		mcp.WithDescription("Read the contents of a file from the file system. You can specify optional from and to parameters to read a specific range of lines. If these are not provided, the complete file contents are returned. Handles various text encodings and provides detailed error messages if the file cannot be read. Use this tool when you need to examine the contents of a single file. Only works within allowed directories."),
		mcp.WithString(
			"path",
			mcp.Required(),
			mcp.Description("The file path to read from. Must be an absolute path to a file that exists on the filesystem."),
		),
		mcp.WithNumber(
			"from",
			mcp.Description("The starting line number to read from. If not specified, the file is read from the beginning."),
		),
		mcp.WithNumber(
			"to",
			mcp.Description("The ending line number to read to. If not specified, the file is read to the end."),
		),
	)

	s.AddTool(readFileTool, tool.ReadFileHandler)
}

func addReadMultipleFilesTool(s *server.MCPServer) {
	readMultipleFilesTool := mcp.NewTool(
		"read_multiple_files",
		mcp.WithDescription("Read the contents of multiple files simultaneously. This is more efficient than reading files one by one when you need to analyze or compare multiple files. Each file's content is returned with its path as a reference. Failed reads for individual files won't stop the entire operation. Only works within allowed directories."),
		mcp.WithArray(
			"paths",
			mcp.Required(),
			mcp.Description("An array of paths to files. All of them must be absolute paths to files that exist on the filesystem."),
		),
	)

	s.AddTool(readMultipleFilesTool, tool.ReadMultipleFilesHandler)
}

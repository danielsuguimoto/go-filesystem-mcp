package main

import (
	"fmt"
	"os"

	"github.com/danielsuguimoto/go-filesystem-mcp/config"
	"github.com/danielsuguimoto/go-filesystem-mcp/tool"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-filesystem-mcp <directory> [additional_directories...]")
		fmt.Println("Example: go-filesystem-mcp /path/to/project /path/to/another")
		fmt.Println("\nAt least one directory must be specified to allow filesystem operations.")
		return
	}

	allowedDirs := os.Args[1:]
	cfg, err := config.NewConfig(allowedDirs)
	if err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		return
	}

	s := server.NewMCPServer(
		"Filesystem MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	addTools(s, cfg)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func addTools(s *server.MCPServer, cfg *config.Config) {
	addReadFileTool(s, cfg)
	addReadMultipleFilesTool(s, cfg)
	addCreateDirectoryTool(s, cfg)
	addListDirectoryTool(s, cfg)
	addDirectoryTreeTool(s, cfg)
}

func addReadFileTool(s *server.MCPServer, cfg *config.Config) {
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

	handler := tool.NewReadFileHandler(cfg)
	s.AddTool(readFileTool, handler.Handle)
}

func addReadMultipleFilesTool(s *server.MCPServer, cfg *config.Config) {
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

func addCreateDirectoryTool(s *server.MCPServer, cfg *config.Config) {
	createDirTool := mcp.NewTool(
		"create_directory",
		mcp.WithDescription("Create a new directory or ensure a directory exists. Can create multiple nested directories in one operation. If the directory already exists, this operation will succeed silently. Perfect for setting up directory structures for projects or ensuring required paths exist. Only works within allowed directories."),
		mcp.WithString(
			"path",
			mcp.Required(),
			mcp.Description("The directory path to create. Must be an absolute path."),
		),
		mcp.WithString(
			"recursive",
			mcp.Description("Create parent directories if they don't exist. Default is 'true'."),
		),
	)

	handler := tool.NewCreateDirectoryHandler(cfg)
	s.AddTool(createDirTool, handler.Handle)
}

func addListDirectoryTool(s *server.MCPServer, cfg *config.Config) {
	listDirectoryTool := mcp.NewTool(
		"list_directory",
		mcp.WithDescription("Get a detailed listing of all files and directories in a specified path. Results clearly distinguish between files and directories with [FILE] and [DIR] prefixes. This tool is essential for understanding directory structure and finding specific files within a directory. Only works within allowed directories."),
		mcp.WithString(
			"path",
			mcp.Required(),
			mcp.Description("The path to list contents of. Must be an absolute path to a directory that exists on the filesystem."),
		),
	)

	handler := tool.NewListDirectoryHandler(cfg)
	s.AddTool(listDirectoryTool, handler.Handle)
}

func addDirectoryTreeTool(s *server.MCPServer, cfg *config.Config) {
	directoryTreeTool := mcp.NewTool(
		"directory_tree",
		mcp.WithDescription("Get a recursive tree view of files and directories as a JSON structure. Each entry includes 'name', 'type' (file/directory), and 'children' for directories. Files have no children array, while directories always have a children array (which may be empty). The output is formatted with 2-space indentation for readability. Only works within allowed directories."),
		mcp.WithString(
			"path",
			mcp.Required(),
			mcp.Description("The directory path to get the tree structure for. Must be an absolute path to a directory that exists on the filesystem."),
		),
	)

	handler := tool.NewDirectoryTreeHandler(cfg)
	s.AddTool(directoryTreeTool, handler.Handle)
}

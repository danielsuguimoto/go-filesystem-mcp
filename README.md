# Go Filesystem MCP Server

Go implementation of Model Context Protocol (MCP) for filesystem operations. Built using the [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) library.

## Features

- Read files (full or partial content)
- Read multiple files simultaneously
- Create directories
- List directory contents
- Get directory tree structure
- Sandboxed operations (only within allowed directories)

## API

### Tools

- **read_file**
  - Read contents of a file with optional line range
  - Inputs:
    - `path` (string): File path to read (absolute path required)
    - `from` (number): Optional. Start line number for reading
    - `to` (number): Optional. End line number
  - Reads file contents with proper encoding handling

- **read_multiple_files**
  - Read multiple files in a single operation
  - Input: `paths` (string[]): List of absolute file paths
  - Individual file read failures don't affect other files

- **create_directory**
  - Create a new directory or ensure it exists
  - Inputs:
    - `path` (string): Directory path to create (absolute path required)
    - `recursive` (string): Create parent directories if missing (default: "true")
  - Creates nested directories in one operation

- **list_directory**
  - List directory contents with [FILE] or [DIR] prefixes
  - Input: `path` (string): Directory path to list (absolute path required)
  - Returns detailed listing of files and directories

- **directory_tree**
  - Get recursive tree view of directory structure
  - Input: `path` (string): Directory path to analyze (absolute path required)
  - Returns JSON structure with:
    - File/directory names
    - Type indicators
    - Nested children for directories

## Usage with AI Assistants

Add this to your MCP configuration file (e.g. `claude_desktop_config.json`, `mcp_config.json`):

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--mount",
        "type=bind,src=/path/to/dir,dst=/path/to/dir",
        "danielsuguimoto/go-filesystem-mcp",
        "/path/to/dir"
      ]
    }
  }
}
```

Note: You need to:
1. Replace `/path/to/dir` with the directories you want to allow access to
2. Add multiple `--mount` args if you need to access multiple directories
3. List all allowed directories at the end of the args array

## License

This project is licensed under the MIT License. See the LICENSE file for details.

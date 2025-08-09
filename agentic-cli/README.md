# GOSS CLI

A command-line interface for chat with local LLMs using LM Studio with MCP (Model Context Protocol) tools for file operations and web search.

## Features

🤖 **Local AI Chat**: Uses LM Studio models with function calling capabilities
🔧 **MCP Tools**: Built-in filesystem and web search tools
📁 **File Operations**: Read, write, list, search, and create files/directories
🌐 **Web Search**: Search the web using Brave Search API (with fallback to DuckDuckGo)
💬 **Chat Interface**: Clean terminal-based chat like original Gemini CLI
📚 **History Management**: Save/load conversation history
🔄 **Model Switching**: Switch between different local models
✨ **Tested & Verified**: All MCP tools tested and working

## Prerequisites

1. **Go 1.21+**: Install from [golang.org](https://golang.org)
2. **LM Studio**: Install and run with Local Server enabled
3. **Function-calling Model**: Load a model that supports function calling (e.g., Mistral, CodeLlama)

## Installation

```bash
# Clone and build
git clone <repository-url>
cd agentic-cli
go mod tidy
go build -o bin/goss cmd/goss/main.go
```

## Usage

### Start LM Studio
1. Open LM Studio
2. Load a function-calling capable model
3. Go to Local Server tab
4. Click "Start Server"
5. Note the server URL (default: http://localhost:1234/v1)

### Run the CLI

```bash
# Basic usage
./bin/goss

# With custom LM Studio endpoint
./bin/goss --base-url http://localhost:1234/v1

# With specific model
./bin/goss --model "mistral-7b-instruct"

# With API key (if required)
export LMSTUDIO_API_KEY=your-key-here
./bin/goss
```

## MCP Tools Available

### Filesystem Tools (✅ Tested & Working)
- `read_file`: Read file contents
- `write_file`: Write content to files  
- `list_directory`: List directory contents with file sizes
- `search_files`: Search for files by pattern (supports wildcards like *.go)
- `create_directory`: Create directories

### Web Search Tools (✅ Tested & Working)
- `web_search`: Search the web using Brave Search API (with DuckDuckGo fallback)
- Secure API key loading from `.env.brave.api` file
- Real-time search results with descriptions and URLs

## Example Interactions

```
> Can you read the package.json file and tell me about the project?

[AI uses read_file tool to read package.json and analyzes it]

> Search for information about "MCP protocol" online

[AI uses web_search tool to find information]

> Create a new file called "summary.md" with a summary of what we learned

[AI uses write_file tool to create the file]
```

## System Commands

The CLI supports system commands prefixed with `!`:

- `!help` - Show help
- `!m` - Model operations (switch model, show info, list tools)
- `!h` - History operations (save, load, clear)
- `!p` - Select system prompts
- `!i` - Toggle input mode
- `!q` - Quit

## Configuration

Create a `gemini_cli_config.json` file:

```json
{
  "SystemPrompts": {
    "Assistant": "You are a helpful AI assistant with access to filesystem and web search tools.",
    "Developer": "You are an expert software developer. Use the available tools to help with coding tasks."
  },
  "History": {}
}
```

## Architecture

### Core Components

1. **OpenAI Client** (`openai/client.go`)
   - HTTP client for LM Studio API
   - Tool definition and execution
   - Chat completions with function calling

2. **MCP Tools** (`mcp/`)
   - `filesystem.go` - File system operations
   - `websearch.go` - Web search capabilities

3. **Agentic Session** (`agentic/chat_session.go`)
   - Manages conversation with tool calling
   - Handles tool execution loop
   - Maintains conversation history

4. **CLI Interface** (`internal/handler/agentic_*.go`)
   - Terminal interface
   - System commands
   - Response formatting

### Tool Execution Flow

1. User sends message
2. LM Studio processes with available tools
3. If model wants to call tools:
   - Extract tool calls from response
   - Execute tools locally
   - Send tool results back to model
   - Get final response
4. Display formatted response to user

## Development

### Adding New Tools

1. Create tool definition in `mcp/` directory
2. Implement `ToolHandler` function
3. Add tool to session in `agentic/chat_session.go`

Example:
```go
func createMyTool() openai.Tool {
    return openai.Tool{
        Type: "function",
        Function: openai.ToolFunction{
            Name: "my_tool",
            Description: "Does something useful",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "param": map[string]interface{}{
                        "type": "string",
                        "description": "Parameter description",
                    },
                },
                "required": []string{"param"},
            },
            Handler: myToolHandler,
        },
    }
}

func myToolHandler(ctx context.Context, args map[string]interface{}) (string, error) {
    param := args["param"].(string)
    // Do something with param
    return "Result", nil
}
```

### Testing

```bash
# Unit tests
go test ./...

# Test MCP tools directly
go run test_filesystem_detailed.go  # Test all filesystem operations
go run test_mcp.go                  # Test MCP integration

# Manual testing
./bin/goss
> !m  # Test model operations  
> !h  # Test history operations
> List files in current directory  # Test filesystem tools
> Search for "golang tutorials"    # Test web search
> Create a file called hello.txt with "Hello World"  # Test file creation
```

## Troubleshooting

### "Connection refused" errors
- Ensure LM Studio Local Server is running
- Check the base URL (default: http://localhost:1234/v1)
- Verify the model is loaded and ready

### "Tool not found" errors
- Make sure you're using a function-calling capable model
- Check LM Studio model supports tool calling
- Try models like Mistral, CodeLlama, or others with function calling

### Build errors
- Ensure Go 1.21+ is installed
- Run `go mod tidy` to update dependencies
- Check for any missing imports

## Performance Tips

1. **Use streaming models** for faster responses
2. **Limit tool complexity** for quicker execution
3. **Cache results** when possible
4. **Use specific file paths** instead of broad searches

## Security Notes

- File operations are limited to the current working directory by default
- Web searches use public APIs only
- No sensitive data is transmitted to external services
- All processing happens locally via LM Studio

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

MIT License - see LICENSE file for details.
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GOSS CLI is a command-line interface for chatting with local LLMs using LM Studio and MCP (Model Context Protocol) tools. It provides a local AI assistant with **real-time streaming responses** and **configurable thinking tokens**, similar to LM Studio's built-in interface. The assistant can perform file operations and web searches through function calling.

**Key Features:**
- ✅ **Real-time streaming responses** - See AI responses as they're generated
- ✅ **Configurable thinking tokens** - View model reasoning process (off/low/med/high)
- ✅ **MCP tool integration** - File operations and web search capabilities
- ✅ **Conversation persistence** - Save and load chat history
- ✅ **Multiple system prompts** - Assistant, Developer, Researcher, Writer modes

**All MCP tools and streaming functionality have been tested and verified working.**

## Development Commands

### Build and Install
```bash
# Quick install system-wide (recommended)
make install

# Manual build
go mod tidy
go build -ldflags "-X main.version=$(date +%Y%m%d)" -o bin/gossai cmd/goss/main.go

# Run locally (before install)
./bin/gossai

# Run system-wide (after make install) 
gossai

# Run with custom LM Studio endpoint
gossai --base-url http://localhost:1234/v1

# Run with specific model
gossai --model "mistral-7b-instruct"
```

### Testing
```bash
# Unit tests
go test ./...

# Test specific packages
go test ./agentic
go test ./openai  
go test ./mcp

# Manual testing
gossai  # or ./bin/gossai for local build

# Legacy test script (needs updating to use gossai)
./test_gemini.sh

# Interactive testing (within gossai CLI)
> !m  # Test model operations
> !h  # Test history operations  
> List files in current directory  # Test filesystem tools
> Search for "golang tutorials"    # Test web search
> Create a test file with some content  # Test file creation
```

### Prerequisites
- Go 1.21+ installed
- LM Studio running with Local Server enabled on http://localhost:1234
- Function-calling capable model loaded in LM Studio (e.g., openai/gpt-oss-20b, Mistral, CodeLlama)

### Installation Options

**Option 1: Quick Install (Recommended)**
```bash
git clone <repository-url>
cd agentic-cli
make install
```
This installs the `gossai` command system-wide in `~/.local/bin`.

**Option 2: Manual Build**
```bash
git clone <repository-url>
cd agentic-cli  
go mod tidy
go build -o bin/gossai cmd/goss/main.go
./bin/gossai
```

**Uninstall**
```bash
make uninstall  # Removes gossai
```

## Architecture

### Core Components

**Entry Point** (`cmd/goss/main.go`):
- Cobra CLI setup with command-line flags
- Initializes agentic chat session with configuration
- Default model: `openai/gpt-oss-20b`
- Default temperature: 0.3

**Agentic Chat Session** (`agentic/chat_session.go`):
- Manages conversation with MCP tool calling
- Handles tool execution loop with LM Studio  
- Maintains conversation history with thread safety
- Implements the core agentic loop: message -> tool calls -> tool results -> response

**OpenAI Client** (`openai/client.go`):
- HTTP client for LM Studio's OpenAI-compatible API
- Tool definition and execution framework
- Chat completions with function calling support
- 30-second timeout for HTTP requests

**MCP Tools** (`mcp/`):
- `filesystem.go` - File system operations (read, write, list, search, create directories)
- `websearch.go` - Web search using Brave Search API with DuckDuckGo fallback

**CLI Interface** (`internal/`):
- `chat/agentic_chat.go` - Main chat handler and terminal interface
- `handler/` - System command handlers for !commands
- `terminal/` - Terminal utilities (colors, spinner, prompt, I/O)
- `config/` - Configuration management for system prompts and history

### Tool Execution Flow

1. User sends message to chat session
2. Message added to conversation history  
3. LM Studio processes request with available tools
4. If model requests tool calls:
   - Extract tool calls from response
   - Execute tools locally via MCP handlers
   - Add tool results to conversation history
   - Continue loop until final response (no more tool calls)
5. Display formatted response with tool usage indicators

### System Commands

The CLI supports these system commands (prefix with `!`):
- `!help` - Show all available commands and help
- `!m` - Model operations (show info, list available tools, select model)
- `!h` - History operations (save, load, clear, delete all)
- `!p` - Select system prompts (Assistant, Developer, etc.)
- `!t` - Temperature adjustment (0.0-1.0, show, set, reset)
- `!stream` - Toggle streaming responses on/off
- `!thinking [level]` - Set thinking level (off/low/med/high)
- `!show-thinking` - Toggle thinking token visibility
- `!i` - Toggle input mode (single-line vs multi-line)
- `!q` - Quit the application

### Configuration

**Default config file**: `goss_config.json` (customizable with `--config` flag)

```json
{
  "SystemPrompts": {
    "Assistant": "You are a helpful AI assistant with access to filesystem and web search tools.",
    "Developer": "You are an expert software developer. Use the available tools to help with coding tasks."
  },
  "History": {},
  "Streaming": {
    "enabled": true,
    "showThinking": false,
    "thinkingLevel": "med"
  }
}
```

**Streaming Configuration:**
- `enabled` - Toggle streaming responses (true/false)
- `showThinking` - Display thinking tokens (true/false)  
- `thinkingLevel` - Thinking token budget:
  - `"off"` - No thinking tokens (0 tokens)
  - `"low"` - Minimal reasoning (50 tokens)
  - `"med"` - Standard reasoning (200 tokens)
  - `"high"` - Detailed reasoning (500 tokens)

**Environment variables**:
- `LMSTUDIO_API_KEY` - Optional API key for LM Studio

**Command-line flags**:
- `--model, -m` - Generative model name (default: openai/gpt-oss-20b)
- `--base-url, -b` - LM Studio API base URL (default: http://localhost:1234/v1)
- `--config, -c` - Configuration file path (default: goss_config.json)
- `--multiline` - Read input as multi-line string
- `--term, -t` - Multi-line input terminator (default: $)
- `--style, -s` - Markdown format style (auto, ascii, dark, light, pink, notty, dracula)
- `--wrap, -w` - Line length for response word wrapping (default: 80)

## Adding New MCP Tools

1. Create tool definition in `mcp/` directory
2. Implement tool handler function with signature: `func(context.Context, map[string]interface{}) (string, error)`
3. Add tool creation function returning `[]openai.Tool`
4. Register tools in `agentic/chat_session.go` NewChatSession function

Example:
```go
func CreateMyTool() []openai.Tool {
    return []openai.Tool{{
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
    }}
}

func myToolHandler(ctx context.Context, args map[string]interface{}) (string, error) {
    param := args["param"].(string)
    // Implement tool logic here
    return "Tool result", nil
}
```

## File Structure

- `cmd/goss/main.go` - Application entry point using Cobra CLI
- `agentic/` - Core agentic chat session logic and tool execution
- `openai/` - OpenAI-compatible API client with tool support
- `mcp/` - Model Context Protocol tool implementations  
- `internal/chat/` - Chat interface and terminal interaction
- `internal/handler/` - System command handlers and response formatting
- `internal/config/` - Configuration file management
- `internal/terminal/` - Terminal utilities (colors, spinner, prompt, I/O)
- `internal/cli/` - CLI command utilities

## Code Quality

### Testing Strategy
- Unit tests in `*_test.go` files alongside source code
- Test coverage for core components (agentic, openai packages)
- Integration testing via interactive CLI session
- Legacy test script needs updating to use `gossai` binary

### Development Workflow
```bash
# Format code
make fmt

# Vet code for issues  
make vet

# Run all quality checks
make check  # Runs fmt + vet + test

# Clean build artifacts
make clean

# Development cycle
make dev    # Builds and runs
```

## Troubleshooting

**Connection Issues**:
- Verify LM Studio Local Server is running on port 1234
- Ensure function-calling model is loaded and ready
- Check base URL configuration with `--base-url` flag

**Tool Execution Issues**:
- Confirm model supports function calling (try openai/gpt-oss-20b)
- Verify tool definitions are properly registered in chat session
- Check tool handler error messages in conversation history

**Build Issues**:
- Ensure Go 1.21+ is installed
- Run `go mod tidy` to update dependencies
- Use `make clean` then `make build` for fresh build

## Quick Usage Summary

After running `make install`:

```bash
# Basic chat
gossai

# With specific model  
gossai --model "mistral-7b-instruct"

# Show help and version
gossai --help
gossai --version

# Example interactions
> List files in current directory
> Search for "MCP protocol documentation"  
> Create a file called notes.txt with today's tasks
> !m  # Show model info and available tools
> !stream  # Toggle streaming responses on/off
> !thinking high  # Enable detailed reasoning display
> !show-thinking  # Toggle thinking token visibility
> !t 0.7  # Adjust temperature for more creative responses
> !q  # Quit
```

System commands within the CLI:
- `!help` - Show all commands
- `!m` - Model operations (info, available tools)
- `!h` - History management (save, load, clear)  
- `!p` - System prompts selection
- `!stream` - Toggle streaming responses
- `!thinking <level>` - Set thinking level (off/low/med/high)
- `!show-thinking` - Toggle thinking token visibility
- `!t <value>` - Temperature adjustment (0.0-1.0)
- `!i` - Toggle input mode
- `!q` - Quit
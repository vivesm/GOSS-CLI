# Usage Guide

## About GOSS CLI

GOSS CLI is a command-line interface for chatting with local LLMs using LM Studio and MCP tools.

**Key Features:**
- Chat with local LLMs via LM Studio
- MCP tools for file operations and web search
- Function calling support for advanced interactions
- Clean terminal interface with system commands

## How to Run GOSS CLI

From this directory (`/Users/melvin/Developer/GitHub/GOSS-CLI/agentic-cli`):

```bash
# Run directly
./bin/gossai

# With options
./bin/gossai --model "openai/gpt-oss-20b" --base-url "http://localhost:1234/v1"

# Check version
./bin/gossai --version

# Get help
./bin/gossai --help
```

## Prerequisites

1. **Start LM Studio**
   - Open LM Studio application
   - Load the `openai/gpt-oss-20b` model (or another function-calling model)
   - Go to "Local Server" tab
   - Click "Start Server"
   - Verify it's running on `http://localhost:1234`

2. **Run the CLI**
   ```bash
   cd /Users/melvin/Developer/GitHub/GOSS-CLI/agentic-cli
   ./bin/gossai
   ```

## Features

- Chat with local LLMs via LM Studio
- MCP tools for:
  - File operations (read, write, list, search, create)
  - Web search via DuckDuckGo
- System commands:
  - `!m` - Model info and operations
  - `!h` - History management
  - `!p` - Prompt selection
  - `!q` - Quit

## Troubleshooting

If you get "command not found", make sure you're in the right directory and use `./bin/gossai` (with the `./` prefix).

If you want to install it system-wide:
```bash
make install  # Installs to ~/.local/bin as gossai
```

Then you can run it from anywhere as just `gossai`.
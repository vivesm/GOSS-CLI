# Usage Guide

## You Have Two Different Programs

1. **System-wide `gemini`** - A different CLI tool (shows the ASCII banner)
   - Located somewhere in your PATH
   - Different project with different features

2. **This Project's `gemini`** - The LM Studio + MCP tools version
   - Located at: `./bin/gemini` 
   - Works with LM Studio and openai/gpt-oss-20b model
   - Includes MCP tools for file operations and web search

## How to Run This Project's Gemini

From this directory (`/Users/melvin/Developer/GitHub/GOSS-CLI/agentic-cli`):

```bash
# Run directly
./bin/gemini

# With options
./bin/gemini --model "openai/gpt-oss-20b" --base-url "http://localhost:1234/v1"

# Check version
./bin/gemini --version

# Get help
./bin/gemini --help
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
   ./bin/gemini
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

If you get "command not found", make sure you're in the right directory and use `./bin/gemini` (with the `./` prefix).

If you want to install it system-wide:
```bash
make install  # Installs to /usr/local/bin
```

Then you can run it from anywhere as just `gemini` (but it might conflict with your other gemini tool).
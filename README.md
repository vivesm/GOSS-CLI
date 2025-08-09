# GOSS-CLI

Universal CLI for local and remote AI models. Works with LM Studio, Ollama, OpenAI, LocalAI, and any OpenAI-compatible API.

## Quick Start

```bash
# Install
npm install

# LM Studio (default)
./bin/goss "Hello, world!"

# Ollama
./bin/goss --provider ollama --model llama2 "Hello!"

# OpenAI
export OPENAI_API_KEY=sk-...
./bin/goss --provider openai --model gpt-4 "Hello!"

# Interactive chat
./bin/goss chat
```

## Supported Providers

### LM Studio (Default)
```bash
# Start LM Studio with Local Server enabled
./bin/goss "Hello!"
./bin/goss --model codellama-7b "Write a function"
```

### Ollama
```bash
# Start Ollama: ollama serve
./bin/goss --provider ollama --model llama2 "Hello!"
./bin/goss --api-base http://localhost:11434 "Test"
```

### OpenAI
```bash
export OPENAI_API_KEY=sk-your-key
./bin/goss --provider openai "Hello!"
./bin/goss --provider openai --model gpt-4 "Complex task"
```

### LocalAI
```bash
# Start LocalAI server
./bin/goss --api-base http://localhost:8080/v1 "Hello!"
```

### Custom OpenAI-Compatible
```bash
./bin/goss --api-base https://your-api.com/v1 --model your-model "Hello!"
```

## Installation

```bash
git clone <this-repo>
cd GOSS-CLI
npm install
cp .env.example .env  # Optional: set defaults
```

## Usage

### Single Prompt Mode
```bash
# Basic usage
./bin/goss "What is the capital of France?"

# With options
./bin/goss --model gpt-oss-20b --temperature 0.8 "Explain quantum computing"

# Disable streaming
./bin/goss --no-stream "Quick math: 2+2"

# Custom API endpoint
./bin/goss --api-base http://localhost:5000/v1 "Hello"
```

### Interactive Chat Mode
```bash
# Start interactive chat
./bin/goss chat

# With custom model
./bin/goss --model codellama-7b chat

# Debug mode (shows API requests/responses)
./bin/goss --debug chat
```

### New Features

#### Save Conversations
```bash
# Save conversation to timestamped file in logs/
./bin/goss --save "Important question"
./bin/goss --save chat

# Output: logs/conversation_2024-01-15T10-30-45.txt
```

#### Context Files
```bash
# Pre-load conversation history (APPENDS to conversation)
./bin/goss --context-file previous-chat.txt "Follow-up question"

# Continue a saved conversation in interactive mode
./bin/goss --context-file logs/conversation_2024-01-15_10-30-45.txt chat

# Context files are loaded BEFORE your prompt/chat, maintaining conversation flow
```

**Context File Behavior:**
- Context files are **appended** to conversation history (not replaced)
- Supports saved conversation format or simple "User:/Assistant:" format
- Path resolution: relative to current directory, works cross-platform
- Invalid/missing files show warning but don't stop execution

#### Auto Model Detection
```bash
# Lists available models if wrong model specified
./bin/goss --model wrong-model "Test"
# Warning: Model 'wrong-model' not found in available models.
# Available models:
#   - gpt-oss-20b
#   - codellama-7b
```

### Configuration

Configure via environment variables (`.env` file):
```bash
PROVIDER=lmstudio              # Provider type
API_BASE=http://localhost:1234/v1
MODEL=gpt-oss-20b
TEMPERATURE=0.7
MAX_TOKENS=2048
OPENAI_API_KEY=sk-...          # For OpenAI provider
```

Or use command-line flags (overrides env vars):
- `--provider <name>`: Provider (lmstudio, ollama, openai, localai)
- `--api-base <url>`: API endpoint URL
- `--model <name>`: Model name
- `--temperature <num>`: Generation temperature (0-2)
- `--max-tokens <n>`: Maximum tokens to generate (1-32000)
- `--debug`: Enable debug logging
- `--no-stream`: Disable streaming responses
- `--save`: Save conversation to timestamped files in logs/
- `--context-file <path>`: Pre-load conversation (appends to history)

### Commands

```bash
# Interactive chat mode
./bin/goss chat

# List available models from provider
./bin/goss list-models
./bin/goss models  # short alias

# Single prompt mode (default when no command specified)
./bin/goss "Your question here"
```

## Development

```bash
# Run tests
npm test

# Development mode with debug output
npm run dev "Test prompt"

# Direct execution
node bin/goss --debug "Hello"
```

## Troubleshooting

### LM Studio Connection Issues
1. Ensure LM Studio is running
2. Check that "Local Server" is enabled in settings
3. Verify the API endpoint (default: `http://localhost:1234/v1`)
4. Test with: `curl http://localhost:1234/v1/models`

### Common Errors
- **ECONNREFUSED**: LM Studio server is not running or wrong port
- **Model not found**: Check model name matches LM Studio's loaded model
- **Timeout**: Increase timeout or reduce `max_tokens`

## Features

- ✅ **Multiple Providers**: LM Studio, Ollama, OpenAI, LocalAI, and any OpenAI-compatible API
- ✅ **Smart Detection**: Auto-detects provider based on API endpoint
- ✅ **Streaming Responses**: Real-time output for better UX
- ✅ **Interactive Chat**: Full conversation history with context
- ✅ **Save Conversations**: Export chats to timestamped files
- ✅ **Context Loading**: Resume previous conversations
- ✅ **Model Validation**: Lists available models when incorrect
- ✅ **Flexible Config**: Environment variables, CLI flags, or both
- ✅ **Debug Mode**: See full API requests/responses
- ✅ **Error Handling**: Helpful messages for common issues
- ✅ **No Lock-in**: Works with local or cloud models

## License

ISC
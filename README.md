# GOSS-CLI

Gemini-like CLI using local GPT-OSS via LM Studio (OpenAI-compatible API).

## Quick Start

```bash
# Start LM Studio -> enable Local Server (default: http://localhost:1234/v1)
cp .env.example .env
npm install
./bin/goss --debug "Say hello"
./bin/goss chat
```

## Installation

1. **Install LM Studio**: Download from [lmstudio.ai](https://lmstudio.ai)
2. **Load a model**: Download and load a GPT-OSS model (e.g., `gpt-oss-20b`)
3. **Enable Local Server**: In LM Studio settings, enable the local server
4. **Clone and install**:
   ```bash
   git clone <this-repo>
   cd GOSS-CLI
   npm install
   cp .env.example .env
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

### Configuration

Configure via environment variables (`.env` file):
```bash
API_BASE=http://localhost:1234/v1
MODEL=gpt-oss-20b
TEMPERATURE=0.7
MAX_TOKENS=2048
```

Or use command-line flags (overrides env vars):
- `--api-base <url>`: API endpoint URL
- `--model <name>`: Model name
- `--temperature <num>`: Generation temperature (0-1)
- `--max-tokens <n>`: Maximum tokens to generate
- `--debug`: Enable debug logging
- `--no-stream`: Disable streaming responses

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

- ✅ OpenAI-compatible API format
- ✅ Streaming responses for better UX
- ✅ Interactive chat with conversation history
- ✅ Environment variable configuration
- ✅ Debug mode for troubleshooting
- ✅ Colored terminal output
- ✅ Error handling with helpful messages
- ✅ No API keys required (local only)

## License

ISC
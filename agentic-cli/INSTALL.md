# Installation Guide

This guide covers different ways to install and run the Gemini CLI.

## Prerequisites

### Required
- **Go 1.21+** - Install from [golang.org](https://golang.org)
- **LM Studio** - Install and run with Local Server enabled
- **Function-calling Model** - Load a model that supports function calling (e.g., Mistral, CodeLlama, GPT-4, etc.)

### Optional
- **Docker** - For containerized deployment
- **Make** - For build automation (included on most Unix systems)

## Installation Methods

### Method 1: Build from Source (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd agentic-cli

# Build using the build script
./build.sh

# Or build using Make
make build

# Or build manually
go mod tidy
go build -ldflags "-X main.version=$(date +%Y%m%d)" -o bin/gemini cmd/gemini/main.go
```

### Method 2: Using Make (Development)

```bash
# Install dependencies and build
make deps
make build

# Development mode (build and run)
make dev

# Install to system PATH
make install

# Clean build artifacts
make clean
```

### Method 3: Docker

```bash
# Build Docker image
docker build -t gemini .

# Run with Docker Compose
docker-compose up

# Or run directly
docker run -it --rm \
  -e LMSTUDIO_BASE_URL=http://host.docker.internal:1234/v1 \
  -v $(pwd)/gemini_cli_config.json:/home/appuser/gemini_cli_config.json \
  gemini
```

## Setup LM Studio

1. **Download and Install LM Studio**
   - Visit [lmstudio.ai](https://lmstudio.ai)
   - Download for your operating system
   - Install and launch

2. **Download a Function-Calling Model**
   - Go to the "Chat" or "Models" section
   - Search for models that support function calling:
     - `microsoft/DialoGPT-medium`
     - `mistralai/Mistral-7B-Instruct-v0.1`
     - `codellama/CodeLlama-7b-Instruct-hf`
     - `openai/gpt-3.5-turbo` (if available)
     - `openai/gpt-oss-20b` (recommended default)
   - Download your preferred model

3. **Start Local Server**
   - Go to "Local Server" tab
   - Select your downloaded model
   - Click "Start Server"
   - Note the server URL (typically `http://localhost:1234/v1`)
   - Ensure "Load model in GPU" is enabled for better performance

4. **Verify Server is Running**
   ```bash
   curl http://localhost:1234/v1/models
   ```

## Configuration

### Environment Variables

Create a `.env` file (copy from `.env.example`):
```bash
cp .env.example .env
```

Edit `.env`:
```bash
# Optional: LM Studio API key (usually not needed for local)
LMSTUDIO_API_KEY=

# Optional: Override default base URL
LMSTUDIO_BASE_URL=http://localhost:1234/v1

# Optional: Default model
DEFAULT_MODEL=openai/gpt-oss-20b
```

### Configuration File

Create configuration file (copy from example):
```bash
cp gemini_cli_config.json.example gemini_cli_config.json
```

Edit `gemini_cli_config.json` to customize system prompts and settings.

## First Run

### Basic Usage
```bash
# Run the CLI
./bin/gemini

# With custom settings
./bin/gemini --base-url http://localhost:1234/v1 --model "mistral-7b-instruct"

# Show help
./bin/gemini --help
```

### Test the Installation
```bash
# Start the CLI
./bin/gemini

# Try these commands in the CLI:
> Hello! Can you help me?
> !m                           # Show model info
> !h                           # Show history commands  
> List files in current directory
> Search for "golang tutorials" online
> Create a file called test.txt with "Hello World"
> !q                           # Quit
```

## Verification

### Check Installation
```bash
# Verify binary exists
ls -la bin/gemini

# Check version
./bin/gemini --version

# Test connection to LM Studio
curl http://localhost:1234/v1/models
```

### Run Tests
```bash
# Unit tests
make test

# Or manually
go test ./...
```

## Troubleshooting

### Build Issues

**"Go not found"**
- Install Go 1.21+ from [golang.org](https://golang.org)
- Add Go to your PATH

**"Permission denied"**
```bash
chmod +x build.sh
chmod +x bin/gemini
```

### Runtime Issues

**"Connection refused"**
- Ensure LM Studio Local Server is running
- Check the base URL (default: `http://localhost:1234/v1`)
- Verify the model is loaded in LM Studio

**"Tool not found" errors**
- Ensure you're using a function-calling capable model
- Try different models like Mistral or CodeLlama
- Check LM Studio model compatibility

**"Model not loaded"**
- Load a model in LM Studio before starting the server
- Wait for the model to fully load (check LM Studio status)
- Try restarting LM Studio

### Docker Issues

**Can't connect to LM Studio from container**
- Use `host.docker.internal:1234` instead of `localhost:1234`
- On Linux, use `--network host` or `network_mode: host`
- Ensure LM Studio allows external connections

## Uninstall

```bash
# Remove system installation
make uninstall

# Remove build artifacts
make clean

# Remove Docker images
docker rmi gemini
```

## Next Steps

- Read the [README.md](README.md) for usage examples
- Check [CLAUDE.md](CLAUDE.md) for development guidance  
- Explore the MCP tools for file operations and web search
- Customize system prompts in the configuration file
# Changelog

All notable changes to GOSS-CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-15

### 🎉 Initial Release

GOSS-CLI is a universal command-line interface for local and remote AI models, supporting multiple providers with a consistent interface.

### ✨ Features

#### Multi-Provider Support
- **LM Studio** - Local GPT-OSS models (default)
- **Ollama** - Local open-source models
- **OpenAI** - GPT-3.5, GPT-4, and other OpenAI models
- **LocalAI** - Self-hosted OpenAI-compatible API
- **Custom APIs** - Any OpenAI-compatible endpoint

#### Core Functionality
- **Interactive Chat Mode** - Full conversation history with streaming responses
- **Single Prompt Mode** - Quick one-off questions
- **Streaming Support** - Real-time response generation
- **Model Management** - List and validate available models

#### Quality-of-Life Features
- **Conversation Persistence** - Save chats to timestamped files (`--save`)
- **Context Loading** - Resume conversations from files (`--context-file`)
- **Auto Model Detection** - Lists available models when invalid model specified
- **Smart Provider Detection** - Auto-detects provider based on API endpoint
- **Debug Mode** - Verbose logging for troubleshooting (`--debug`)

#### Developer Experience
- **Cross-Platform** - Works on macOS, Linux, and Windows
- **Flexible Configuration** - Environment variables, CLI flags, or both
- **Error Handling** - Graceful handling of connection issues and malformed responses
- **Comprehensive Testing** - 18 test cases covering core functionality

### 🚀 Commands

```bash
# Interactive chat
goss chat

# Single prompt
goss "Your question here"

# List available models
goss list-models

# Save conversation
goss --save chat

# Resume conversation
goss --context-file logs/conversation_2024-01-15_10-30-45.txt chat
```

### ⚙️ Configuration

**Environment Variables:**
- `PROVIDER` - Provider type (lmstudio, ollama, openai, localai)
- `API_BASE` - API endpoint URL
- `MODEL` - Default model name
- `TEMPERATURE` - Generation temperature (0-2)
- `MAX_TOKENS` - Maximum tokens (1-32000)
- `OPENAI_API_KEY` - For OpenAI provider

**CLI Flags:**
- `--provider <name>` - Override provider
- `--api-base <url>` - Override API endpoint
- `--model <name>` - Override model
- `--temperature <num>` - Override temperature
- `--max-tokens <n>` - Override max tokens
- `--debug` - Enable verbose logging
- `--no-stream` - Disable streaming
- `--save` - Save conversation to logs/
- `--context-file <path>` - Pre-load conversation

### 🛠️ Technical Details

- **Node.js** 16+ required
- **Dependencies** - axios, commander, inquirer, chalk, ora, dotenv
- **License** - MIT
- **Testing** - Vitest with comprehensive mocks
- **Architecture** - Modular provider system for easy extensibility

### 📦 Installation

```bash
# Via npm
npm install -g goss-cli

# Via git
git clone https://github.com/your-username/GOSS-CLI.git
cd GOSS-CLI
npm install
npm link
```

---

*This is the first stable release of GOSS-CLI. Future versions will add more providers, enhanced features, and community-requested improvements.*
# 🚀 GOSS-CLI v1.0.0 - Universal AI Model CLI

## 🎉 First Stable Release

GOSS-CLI is now production-ready! A universal command-line interface that works with **LM Studio**, **Ollama**, **OpenAI**, **LocalAI**, and any **OpenAI-compatible API**.

### ✨ What's New in v1.0.0

🔥 **Multi-Provider Support** - Switch seamlessly between local and cloud models  
💾 **Conversation Persistence** - Save and resume chats with `--save` and `--context-file`  
🎯 **Smart Model Detection** - Auto-lists available models when you specify an invalid one  
🛡️ **Battle-Tested** - Comprehensive error handling and 18+ test cases  
🌍 **Cross-Platform** - Works identically on macOS, Linux, and Windows  

### 🚀 Quick Start

```bash
# Install globally
npm install -g goss-cli

# Use with any provider
goss "Explain quantum computing"                    # LM Studio (default)
goss --provider ollama --model llama2 "Code review"  # Ollama
goss --provider openai --model gpt-4 "Complex task" # OpenAI
```

### 📋 Full Feature List

- **5 Providers**: LM Studio, Ollama, OpenAI, LocalAI, Custom APIs
- **Interactive Chat**: Full conversation history with streaming
- **Single Prompts**: Quick one-off questions
- **Save Conversations**: Timestamped logs in `logs/` directory
- **Resume Chats**: Load previous conversations with `--context-file`
- **Model Management**: List available models with `goss list-models`
- **Smart Detection**: Auto-detects provider from API endpoint
- **Debug Mode**: Verbose logging for troubleshooting
- **Flexible Config**: Environment variables, CLI flags, or both
- **Error Recovery**: Graceful handling of connection drops and malformed JSON
- **Stream Timeouts**: Prevents hanging on stuck connections

### 🛠️ Installation Options

```bash
# npm (recommended)
npm install -g goss-cli

# Git clone
git clone https://github.com/your-username/GOSS-CLI.git
cd GOSS-CLI && npm install && npm link

# Download release assets below ⬇️
```

### 📚 Documentation

- **Full README**: Complete setup guide and examples
- **Provider Guides**: Specific instructions for each AI provider
- **Configuration**: All environment variables and CLI flags
- **Troubleshooting**: Common issues and solutions

### 🔗 Links

- **npm Package**: https://www.npmjs.com/package/goss-cli
- **Documentation**: [README.md](https://github.com/your-username/GOSS-CLI#readme)
- **Report Issues**: [GitHub Issues](https://github.com/your-username/GOSS-CLI/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/GOSS-CLI/discussions)

### 🙏 Acknowledgments

Special thanks to the communities around **LM Studio**, **Ollama**, and the broader **local AI** ecosystem for inspiration and testing feedback.

### 🚀 What's Next?

- **v1.1**: Additional providers (Anthropic Claude, Cohere, etc.)
- **Enhanced Features**: Configuration profiles, conversation search, model switching mid-chat
- **Community Requests**: [Share your ideas!](https://github.com/your-username/GOSS-CLI/discussions)

---

**Ready to unify your AI CLI experience?**

⭐ **Star the repo** if GOSS-CLI is useful to you!  
🗣️ **Share** in AI communities (r/LocalLLaMA, LM Studio Discord, etc.)  
🤝 **Contribute** new providers or improvements
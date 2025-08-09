# 🚀 GOSS-CLI v1.0.0 Launch

## Universal CLI for Local & Remote AI Models

**GOSS-CLI** is a battle-tested, cross-platform command-line interface that works with **any AI model provider**. Born from the need for a unified CLI that works with both local and cloud models.

### 🎯 **What Makes It Special**

✅ **Multi-Provider** - Works with LM Studio, Ollama, OpenAI, LocalAI, and any OpenAI-compatible API  
✅ **Smart Detection** - Auto-detects your provider based on endpoint  
✅ **Conversation Memory** - Save and resume conversations  
✅ **Battle-Tested** - 18+ test cases, comprehensive error handling  
✅ **Cross-Platform** - macOS, Linux, Windows support  
✅ **Zero Lock-in** - Switch between local and cloud models seamlessly  

### 🚀 **Quick Start**

```bash
# Install globally via npm
npm install -g goss-cli

# LM Studio (default - just works!)
goss "Explain quantum computing in simple terms"

# Ollama 
goss --provider ollama --model llama2 "Write a Python function"

# OpenAI
export OPENAI_API_KEY=sk-your-key
goss --provider openai --model gpt-4 "Complex reasoning task"

# Interactive chat mode
goss chat
```

### 🔥 **Key Features**

| Feature | Description | Command |
|---------|-------------|---------|
| **Multi-Provider** | Works with 5+ providers | `--provider ollama` |
| **Save Conversations** | Timestamped conversation logs | `--save` |
| **Resume Chats** | Load previous conversations | `--context-file logs/chat.txt` |
| **List Models** | See available models | `goss list-models` |
| **Debug Mode** | Troubleshoot API issues | `--debug` |
| **Streaming** | Real-time responses | Default behavior |

### 🏗️ **Architecture Highlights**

- **Provider System** - Modular design, easy to extend
- **Error Recovery** - Graceful handling of connection drops
- **Config Validation** - Helpful warnings for invalid settings
- **Cross-Platform Paths** - Works identically on all OS

### 🎯 **Perfect For**

- **AI Hobbyists** running local models (LM Studio, Ollama)
- **Developers** who want consistent CLI across providers  
- **Privacy-conscious users** preferring local inference
- **Teams** switching between local dev and cloud prod
- **Automation** scripts that need reliable AI interaction

### 📦 **Installation Options**

```bash
# npm (recommended)
npm install -g goss-cli

# Git clone
git clone https://github.com/your-username/GOSS-CLI.git
cd GOSS-CLI && npm install && npm link

# Download release
curl -L https://github.com/your-username/GOSS-CLI/releases/latest/download/goss-cli.tar.gz
```

### 🌟 **Example Use Cases**

```bash
# Code review with local model
goss --save "Review this Python code: $(cat script.py)"

# Continuous conversation
goss --context-file previous-chat.txt "Follow up on that suggestion"

# Switch providers mid-project
goss --provider openai --model gpt-4 "Complex analysis task"
goss --provider ollama --model codellama "Simple code generation"

# Automation-friendly
echo "Summarize this log" | goss --no-stream > summary.txt
```

### 🔗 **Links**

- **GitHub**: [https://github.com/your-username/GOSS-CLI](https://github.com/your-username/GOSS-CLI)
- **npm**: [https://www.npmjs.com/package/goss-cli](https://www.npmjs.com/package/goss-cli)
- **Documentation**: Full README with all providers and options
- **Issues**: Report bugs or request features

### 💬 **Community**

Found this useful? 
- ⭐ Star the repo
- 🗣️ Share in LM Studio Discord, Ollama forums, r/LocalLLaMA
- 🐛 Report issues or suggest features
- 🤝 Contribute new providers or improvements

---

**Built for the community, by developers who actually use local AI models daily.**

*Ready to unify your AI CLI experience? Try GOSS-CLI today!*
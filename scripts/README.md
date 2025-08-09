# 🧪 GOSS-CLI Test Scripts

## Smoke Tests

Pre-launch validation that covers all critical functionality in under 2 minutes.

### Quick Run

```bash
# From repo root
bash scripts/smoke.sh

# Or via npm
npm run test:smoke

# Run all tests (unit + smoke)
npm run test:all
```

### What It Tests

| Test | Description | Validates |
|------|-------------|-----------|
| **CLI Present** | Executable exists and is runnable | Basic installation |
| **Help Text** | `--help` shows correct output | CLI framework working |
| **Non-Stream JSON** | Deterministic JSON response | API integration |
| **Streaming** | Real-time response chunks | Stream handling |
| **Save Transcript** | `--save` creates log files | File persistence |
| **Context File** | `--context-file` loads history | Conversation continuity |
| **List Models** | Shows available models | Provider integration |
| **Invalid Model** | Helpful error for bad model | Error handling |
| **Provider Override** | `--provider` flag works | Multi-provider support |
| **Connection Error** | Clear message for unreachable API | Network error handling |
| **Debug Logging** | `--debug` shows internal info | Troubleshooting |

### Prerequisites

**One of these providers must be running:**
- LM Studio with Local Server enabled
- Ollama with `ollama serve`
- OpenAI API key in `OPENAI_API_KEY` env var
- LocalAI or custom OpenAI-compatible server

### Platform Support

#### Linux/macOS
```bash
bash scripts/smoke.sh
```

#### Windows (PowerShell)
```powershell
.\scripts\smoke.ps1
```

#### Windows (Git Bash)
```bash
bash scripts/smoke.sh
```

### Customization

**Timeout Adjustment:**
```bash
# Edit smoke.sh and change timeout values
timeout 30 ./bin/goss ...  # Increase for slower systems
```

**Provider-Specific Testing:**
```bash
# Test specific provider
./bin/goss --provider ollama --debug 'test'

# Force different endpoint
./bin/goss --api-base http://localhost:8080/v1 'test'
```

### Troubleshooting

#### ❌ "No provider running"
- Start LM Studio and enable Local Server
- Or start Ollama: `ollama serve`
- Or set `OPENAI_API_KEY` environment variable

#### ❌ "Timeout errors"
- Increase timeout values in smoke.sh
- Check provider is loaded with models
- Verify network connectivity

#### ❌ "Permission denied"
- Make script executable: `chmod +x scripts/smoke.sh`
- Check file system permissions for logs directory

#### ❌ "Context file test fails"
- Verify `/tmp` directory is writable
- On Windows, check `%TEMP%` permissions

### CI Integration

**GitHub Actions:**
```yaml
- name: Run smoke tests
  run: |
    # Start test provider (e.g., mock server)
    npm run test:smoke
```

**Local Pre-Publish:**
```bash
# Full validation before npm publish
npm run test:all
git push origin main --tags
npm publish --access public
```

### Adding New Tests

1. Add test function to `smoke.sh`
2. Follow the pattern: `command && pass "description" || fail "error"`
3. Use timeouts for external calls
4. Clean up temporary files
5. Update this README

### Performance

- **Total time**: ~2 minutes with running provider
- **Network calls**: 8-10 API requests
- **Disk usage**: <1MB (log files created/cleaned)
- **Memory**: Minimal (CLI process only)

---

**These tests ensure GOSS-CLI works correctly before every release! 🚀**
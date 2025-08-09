# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GOSS-CLI (GPT-OSS CLI) is a modified version of the Gemini CLI that routes inference requests to a locally hosted GPT-OSS model (e.g., gpt-oss-20b) running in LM Studio instead of Google's Gemini API. This enables offline, privacy-focused AI assistance through a familiar CLI interface.

## Project Status

This is a new project in the planning phase. The Product Requirements Document (PRD) is complete, but implementation has not yet begun.

## Key Technical Requirements

### API Integration
- **Target Endpoint**: LM Studio's OpenAI-compatible API at `http://localhost:1234/v1/chat/completions`
- **Model**: Default to `gpt-oss-20b`, configurable via `--model` flag
- **Authentication**: No API keys required (local-only)

### Request/Response Format Conversion

Convert from Gemini format:
```json
{
  "contents": [
    { "role": "user", "parts": [{ "text": "Hello!" }] }
  ]
}
```

To OpenAI format:
```json
{
  "model": "gpt-oss-20b",
  "messages": [
    { "role": "user", "content": "Hello!" }
  ]
}
```

### Response Parsing
- **Gemini path**: `response.candidates[0].content.parts[0].text`
- **OpenAI path**: `response.choices[0].message.content`

## Implementation Approach

1. Fork the original Gemini CLI repository
2. Locate and modify the API client layer (likely in `client.js` or `api.ts`)
3. Implement request converter and response mapper functions
4. Add `--api-base` flag for endpoint customization
5. Remove Google API key dependencies
6. Test against LM Studio with gpt-oss-20b model

## CLI Configuration

### Environment Variables
- `GOSS_API_BASE`: Override default LM Studio endpoint
- `GOSS_MODEL`: Default model name

### CLI Flags
- `--api-base <url>`: Specify LM Studio endpoint
- `--model <name>`: Specify model name

## Prerequisites

- LM Studio installed and running with API server enabled
- Node.js (version compatible with original Gemini CLI)
- A compatible GPT-OSS model loaded in LM Studio

## Error Handling

When LM Studio is not reachable:
1. Check if LM Studio is running
2. Verify API server is enabled in LM Studio settings
3. Confirm the endpoint URL matches LM Studio's configuration
4. Provide clear error messages with troubleshooting steps

## Testing Strategy

1. Basic connectivity test with LM Studio
2. Simple prompt/response validation
3. Multi-turn conversation support
4. Error case handling (offline server, malformed responses)
5. Performance benchmarking vs original Gemini CLI

## Key Files to Create/Modify

- API client module (conversion logic)
- Configuration handler (endpoint, model selection)
- Error handling utilities
- CLI argument parser updates
- Documentation and setup guide

## Success Metrics

- CLI runs without Google API key
- All inference happens locally through LM Studio
- Response time within 1-3 seconds (hardware dependent)
- Maintains original Gemini CLI user experience
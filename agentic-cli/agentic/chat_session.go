package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/vivesm/GOSS-CLI/agentic-cli/mcp"
	"github.com/vivesm/GOSS-CLI/agentic-cli/openai"
)

const (
	DefaultModel = "openai/gpt-oss-20b"
	// MaxHistorySize limits conversation history to prevent memory issues
	MaxHistorySize = 50
	// MaxContextTokens approximates max context length
	MaxContextTokens = 4000
)

// ChatSession represents an agentic chat session with MCP tools
type ChatSession struct {
	ctx         context.Context
	client      *openai.Client
	model       string
	temperature float64
	maxTokens   int
	history     []openai.Message

	mu sync.Mutex
}

// SessionConfig holds configuration for the chat session
type SessionConfig struct {
	BaseURL     string
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
}

// NewChatSession creates a new agentic chat session
func NewChatSession(ctx context.Context, config SessionConfig) (*ChatSession, error) {
	if config.Model == "" {
		config.Model = DefaultModel
	}

	client := openai.NewClient(config.BaseURL, config.APIKey)

	// Add MCP tools
	filesystemTools := mcp.CreateFilesystemTools()
	for _, tool := range filesystemTools {
		client.AddTool(tool)
	}

	websearchTools := mcp.CreateWebSearchTools()
	for _, tool := range websearchTools {
		client.AddTool(tool)
	}

	// Set defaults if not provided
	temperature := config.Temperature
	if temperature == 0 {
		temperature = 0.3 // Default focused temperature
	}
	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2048 // Default max tokens
	}

	session := &ChatSession{
		ctx:         ctx,
		client:      client,
		model:       config.Model,
		temperature: temperature,
		maxTokens:   maxTokens,
		history:     make([]openai.Message, 0),
	}

	// Add default system message for better tool usage
	session.SetSystemMessage(`You are an AI assistant with MANDATORY tool usage requirements.

CRITICAL INSTRUCTIONS - YOU MUST FOLLOW THESE:

1. TOOL USAGE IS REQUIRED for these queries:
   - ANY question about current events, news, or recent information → USE web_search tool
   - ANY weather-related query → USE web_search tool  
   - ANY request to read, write, list, or search files → USE appropriate filesystem tool
   - ANY request for current information that you don't have → USE web_search tool

2. NEVER respond with "I couldn't find" or "I'm sorry" without first using tools.

3. When you receive a query like:
   - "search for X" → You MUST call web_search with query="X"
   - "what's the weather in Y" → You MUST call web_search with query="weather in Y today"
   - "list files" → You MUST call list_directory
   - "read file X" → You MUST call read_file

4. Available tools you MUST use:
   - web_search: For ANY current information, news, weather, or real-time data
   - read_file, write_file, list_directory, search_files, create_directory: For file operations

5. INCORRECT BEHAVIOR (DO NOT DO THIS):
   - Saying "I couldn't retrieve" without calling tools
   - Apologizing for not finding information without searching
   - Suggesting the user check elsewhere without trying tools first

Remember: You have tools available. USE THEM. Do not respond without attempting to use relevant tools first.`)

	return session, nil
}

// SetTemperature updates the temperature setting
func (s *ChatSession) SetTemperature(temperature float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clamp temperature to valid range (OpenAI standard is 0.0-1.0)
	if temperature < 0.0 {
		temperature = 0.0
	}
	if temperature > 1.0 {
		temperature = 1.0
	}

	s.temperature = temperature
}

// GetTemperature returns the current temperature setting
func (s *ChatSession) GetTemperature() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.temperature
}

// SetMaxTokens updates the max tokens setting
func (s *ChatSession) SetMaxTokens(maxTokens int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if maxTokens < 1 {
		maxTokens = 1
	}
	if maxTokens > 8192 {
		maxTokens = 8192
	}

	s.maxTokens = maxTokens
}

// GetMaxTokens returns the current max tokens setting
func (s *ChatSession) GetMaxTokens() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.maxTokens
}

// SetSystemMessage sets or updates the system message
func (s *ChatSession) SetSystemMessage(content string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing system message if present
	if len(s.history) > 0 && s.history[0].Role == "system" {
		s.history = s.history[1:]
	}

	// Add new system message at the beginning
	systemMsg := openai.Message{
		Role:    "system",
		Content: content,
	}
	s.history = append([]openai.Message{systemMsg}, s.history...)
}

// SendMessage sends a message and handles tool calls
func (s *ChatSession) SendMessage(input string) (*AgenticResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add user message to history
	userMsg := openai.Message{
		Role:    "user",
		Content: input,
	}
	s.history = append(s.history, userMsg)
	
	// Trim history to prevent memory issues
	s.trimHistory()

	maxIterations := 10 // Prevent infinite loops
	iteration := 0

	for iteration < maxIterations {
		iteration++
		// Create chat completion request
		req := openai.ChatCompletionRequest{
			Model:       s.model,
			Messages:    s.history,
			Temperature: s.temperature,
			MaxTokens:   s.maxTokens,
		}

		// Send request to LM Studio
		resp, err := s.client.CreateChatCompletion(s.ctx, req)
		if err != nil {
			return nil, fmt.Errorf("chat completion failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("no response choices returned")
		}

		choice := resp.Choices[0]
		assistantMsg := choice.Message

		// Add assistant message to history
		s.history = append(s.history, assistantMsg)

		// If there are tool calls, execute them
		if len(assistantMsg.ToolCalls) > 0 {
			err := s.executeToolCalls(assistantMsg.ToolCalls)
			if err != nil {
				return nil, fmt.Errorf("tool execution failed: %w", err)
			}
			// Continue the loop to get the final response
			continue
		}

		// No tool calls, return the response
		return &AgenticResponse{
			Content:      assistantMsg.Content,
			ToolCalls:    len(assistantMsg.ToolCalls) > 0,
			FinishReason: choice.FinishReason,
			Usage:        resp.Usage,
		}, nil
	}

	// If we exit the loop without a response, return the last assistant message if available
	if len(s.history) > 0 {
		for i := len(s.history) - 1; i >= 0; i-- {
			if s.history[i].Role == "assistant" {
				return &AgenticResponse{
					Content:      fmt.Sprintf("Response reached maximum iterations (%d). Last response: %s", maxIterations, s.history[i].Content),
					ToolCalls:    false,
					FinishReason: "max_iterations",
					Usage:        openai.Usage{},
				}, nil
			}
		}
	}

	return &AgenticResponse{
		Content:      fmt.Sprintf("Response reached maximum iterations (%d) without completion", maxIterations),
		ToolCalls:    false,
		FinishReason: "max_iterations",
		Usage:        openai.Usage{},
	}, nil
}

// executeToolCalls executes tool calls and adds results to history
func (s *ChatSession) executeToolCalls(toolCalls []openai.ToolCall) error {
	for _, toolCall := range toolCalls {
		result, err := s.client.ExecuteTool(s.ctx, toolCall)
		if err != nil {
			result = fmt.Sprintf("Error executing tool %s: %v", toolCall.Function.Name, err)
		}

		// Add tool result to history
		toolResultMsg := openai.Message{
			Role:       "tool",
			Content:    result,
			ToolCallID: toolCall.ID,
		}
		s.history = append(s.history, toolResultMsg)
	}
	return nil
}

// GetHistory returns the current conversation history
func (s *ChatSession) GetHistory() []openai.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return a copy to prevent external modification
	history := make([]openai.Message, len(s.history))
	copy(history, s.history)
	return history
}

// SetHistory sets the conversation history
func (s *ChatSession) SetHistory(history []openai.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.history = make([]openai.Message, len(history))
	copy(s.history, history)
}

// ClearHistory clears the conversation history
func (s *ChatSession) ClearHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.history = make([]openai.Message, 0)
}

// trimHistory keeps conversation history within reasonable limits
func (s *ChatSession) trimHistory() {
	if len(s.history) <= MaxHistorySize {
		return
	}

	// Keep system message if present
	systemMsg := []openai.Message{}
	startIdx := 0
	if len(s.history) > 0 && s.history[0].Role == "system" {
		systemMsg = s.history[:1]
		startIdx = 1
	}

	// Keep the most recent messages within the limit
	keepCount := MaxHistorySize - len(systemMsg)
	if keepCount > 0 {
		startKeep := len(s.history) - keepCount
		if startKeep < startIdx {
			startKeep = startIdx
		}
		s.history = append(systemMsg, s.history[startKeep:]...)
	} else {
		s.history = systemMsg
	}
}

// ListModels returns available models
func (s *ChatSession) ListModels() ([]string, error) {
	return s.client.ListModels(s.ctx)
}

// SetModel changes the active model
func (s *ChatSession) SetModel(model string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.model = model
}

// GetModel returns the current model
func (s *ChatSession) GetModel() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.model
}

// ModelInfo returns information about the current model
func (s *ChatSession) ModelInfo() (string, error) {
	info := map[string]interface{}{
		"name":        s.model,
		"type":        "OpenAI Compatible",
		"base_url":    s.client.BaseURL,
		"tools_count": len(s.client.Tools),
		"tools":       s.getToolNames(),
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal model info: %w", err)
	}

	return string(jsonData), nil
}

// getToolNames returns the names of available tools
func (s *ChatSession) getToolNames() []string {
	var names []string
	for _, tool := range s.client.Tools {
		names = append(names, tool.Function.Name)
	}
	return names
}

// Close closes the chat session (placeholder for compatibility)
func (s *ChatSession) Close() error {
	// Nothing to close for HTTP client
	return nil
}

// AgenticResponse represents a response from the agentic chat session
type AgenticResponse struct {
	Content      string       `json:"content"`
	ToolCalls    bool         `json:"tool_calls"`
	FinishReason string       `json:"finish_reason"`
	Usage        openai.Usage `json:"usage"`
}

// FormatResponse formats the response for display
func (r *AgenticResponse) FormatResponse() string {
	var result strings.Builder

	result.WriteString(r.Content)

	if r.ToolCalls {
		result.WriteString("\n\n---\n")
		result.WriteString("🔧 *Generated using MCP tools (filesystem & web search)*")
	}

	return result.String()
}

package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/reugn/gemini-cli/mcp"
	"github.com/reugn/gemini-cli/openai"
)

const DefaultModel = "openai/gpt-oss-20b"

// ChatSession represents an agentic chat session with MCP tools
type ChatSession struct {
	ctx     context.Context
	client  *openai.Client
	model   string
	history []openai.Message
	
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
	
	return &ChatSession{
		ctx:     ctx,
		client:  client,
		model:   config.Model,
		history: make([]openai.Message, 0),
	}, nil
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
	
	for {
		// Create chat completion request
		req := openai.ChatCompletionRequest{
			Model:       s.model,
			Messages:    s.history,
			Temperature: 0.7,
			MaxTokens:   2048,
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
	Content      string           `json:"content"`
	ToolCalls    bool             `json:"tool_calls"`
	FinishReason string           `json:"finish_reason"`
	Usage        openai.Usage     `json:"usage"`
}

// FormatResponse formats the response for display
func (r *AgenticResponse) FormatResponse() string {
	var result strings.Builder
	
	result.WriteString(r.Content)
	
	if r.ToolCalls {
		result.WriteString("\n\n[Tools were used to generate this response]")
	}
	
	return result.String()
}
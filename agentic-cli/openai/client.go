package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client represents an OpenAI-compatible API client
type Client struct {
	BaseURL    string
	APIKey     string
	httpClient *http.Client
	Tools      []Tool
}

// NewClient creates a new OpenAI-compatible client
func NewClient(baseURL, apiKey string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Tools: make([]Tool, 0),
	}
}

// AddTool adds an MCP tool to the client
func (c *Client) AddTool(tool Tool) {
	c.Tools = append(c.Tools, tool)
}

// ChatCompletionRequest represents the request structure for chat completions
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	Tools       []Tool    `json:"tools,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// ToolCall represents a function call
type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents a function call
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool represents an available tool/function
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents the function definition
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Handler     ToolHandler            `json:"-"`
}

// ToolHandler is a function that executes a tool
type ToolHandler func(ctx context.Context, args map[string]interface{}) (string, error)

// ChatCompletionResponse represents the response from chat completions
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamingChoice represents a streaming completion choice
type StreamingChoice struct {
	Index int          `json:"index"`
	Delta StreamingDelta `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

// StreamingDelta represents the delta content in streaming response
type StreamingDelta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	Reasoning string     `json:"reasoning,omitempty"` // LM Studio thinking tokens
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ChatCompletionStreamResponse represents a streaming response chunk
type ChatCompletionStreamResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []StreamingChoice `json:"choices"`
}

// StreamCallback is called for each streaming chunk
type StreamCallback func(chunk ChatCompletionStreamResponse) error

// ThinkingLevel represents different levels of thinking detail
type ThinkingLevel struct {
	Name   string `json:"name"`   // "off", "low", "med", "high"
	Budget int    `json:"budget"` // Token budget for thinking
}

// GetThinkingBudget returns the token budget for a thinking level
func GetThinkingBudget(level string) int {
	budgets := map[string]int{
		"off":  0,
		"low":  50,
		"med":  200,
		"high": 500,
	}
	if budget, exists := budgets[level]; exists {
		return budget
	}
	return 200 // Default to medium
}

// CreateChatCompletion sends a chat completion request
func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Add tools to request if available
	if len(c.Tools) > 0 {
		req.Tools = c.Tools
		// Debug: Log that tools are being sent
		if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
			fmt.Fprintf(os.Stderr, "[DEBUG] Sending %d tools with request\n", len(c.Tools))
			for _, tool := range c.Tools {
				fmt.Fprintf(os.Stderr, "[DEBUG] Tool: %s - %s\n", tool.Function.Name, tool.Function.Description)
			}
		}
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	
	// Debug: Log the full request if in debug mode
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Request to %s\n", c.BaseURL+"chat/completions")
		if len(req.Messages) > 0 {
			lastMsg := req.Messages[len(req.Messages)-1]
			fmt.Fprintf(os.Stderr, "[DEBUG] Last message role: %s, content: %.100s...\n", 
				lastMsg.Role, lastMsg.Content)
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// Debug: Log response details
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		if len(chatResp.Choices) > 0 {
			choice := chatResp.Choices[0]
			fmt.Fprintf(os.Stderr, "[DEBUG] Response received\n")
			fmt.Fprintf(os.Stderr, "[DEBUG] Tool calls in response: %d\n", len(choice.Message.ToolCalls))
			if len(choice.Message.ToolCalls) > 0 {
				for _, tc := range choice.Message.ToolCalls {
					fmt.Fprintf(os.Stderr, "[DEBUG] Tool call: %s\n", tc.Function.Name)
				}
			} else if choice.Message.Content != "" {
				fmt.Fprintf(os.Stderr, "[DEBUG] Text response (no tools): %.100s...\n", choice.Message.Content)
			}
		}
	}

	return &chatResp, nil
}

// ExecuteTool executes a tool call and returns the result
func (c *Client) ExecuteTool(ctx context.Context, toolCall ToolCall) (string, error) {
	// Find the tool handler
	var handler ToolHandler
	for _, tool := range c.Tools {
		if tool.Function.Name == toolCall.Function.Name {
			handler = tool.Function.Handler
			break
		}
	}

	if handler == nil {
		return "", fmt.Errorf("tool not found: %s", toolCall.Function.Name)
	}

	// Parse arguments
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return "", fmt.Errorf("parse tool arguments: %w", err)
	}

	// Execute tool
	return handler(ctx, args)
}

// CreateChatCompletionStream sends a streaming chat completion request
func (c *Client) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest, thinkingLevel string, showThinking bool, callback StreamCallback) error {
	// Force streaming mode
	req.Stream = true
	
	// Add tools to request if available
	if len(c.Tools) > 0 {
		req.Tools = c.Tools
	}
	
	// LM Studio doesn't need extra thinking configuration - it provides reasoning automatically
	// Just make a standard streaming request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal streaming request: %w", err)
	}
	
	return c.performStreamingRequest(ctx, reqBody, callback)
}

func (c *Client) performStreamingRequest(ctx context.Context, reqBody []byte, callback StreamCallback) error {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("create streaming request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Cache-Control", "no-cache")
	
	if c.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	
	// Debug: log request details
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Sending streaming request to %s\n", httpReq.URL)
		fmt.Fprintf(os.Stderr, "[DEBUG] Headers: %v\n", httpReq.Header)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send streaming request: %w", err)
	}
	defer resp.Body.Close()
	
	// Debug: log response details
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Streaming response status: %s\n", resp.Status)
		fmt.Fprintf(os.Stderr, "[DEBUG] Response headers: %v\n", resp.Header)
		fmt.Fprintf(os.Stderr, "[DEBUG] Content-Type: %s\n", resp.Header.Get("Content-Type"))
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("streaming API error (%d): %s", resp.StatusCode, string(body))
	}
	
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Starting to process streaming response...\n")
	}
	
	return c.processStreamingResponse(ctx, resp.Body, callback)
}

func (c *Client) processStreamingResponse(ctx context.Context, body io.Reader, callback StreamCallback) error {
	scanner := bufio.NewScanner(body)
	
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		line := scanner.Text()
		
		// Debug mode: log all lines
		if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
			fmt.Fprintf(os.Stderr, "[DEBUG] SSE Line: %q\n", line)
		}
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		
		// Check for end of stream
		if strings.TrimSpace(line) == "data: [DONE]" {
			if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
				fmt.Fprintf(os.Stderr, "[DEBUG] Stream completed with [DONE]\n")
			}
			break
		}
		
		// Parse SSE data
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			
			// Skip empty data
			if strings.TrimSpace(data) == "" {
				continue
			}
			
			var chunk ChatCompletionStreamResponse
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				// Debug malformed JSON in debug mode
				if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
					fmt.Fprintf(os.Stderr, "[DEBUG] Failed to parse streaming chunk: %s, Error: %s\n", data, err)
				}
				continue // Skip malformed chunks
			}
			
			// Debug mode: log successful parse
			if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
				fmt.Fprintf(os.Stderr, "[DEBUG] Parsed chunk, choices: %d\n", len(chunk.Choices))
			}
			
			// Call the callback with the parsed chunk
			if err := callback(chunk); err != nil {
				return fmt.Errorf("streaming callback error: %w", err)
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading streaming response: %w", err)
	}
	
	return nil
}

// ListModels returns available models (placeholder implementation)
func (c *Client) ListModels(ctx context.Context) ([]string, error) {
	return []string{"openai/gpt-oss-20b", "meta-llama/llama-3.2-3b-instruct"}, nil
}

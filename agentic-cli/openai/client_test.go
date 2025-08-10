package openai

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:1234"
	apiKey := "test-key"

	client := NewClient(baseURL, apiKey)

	if client == nil {
		t.Fatal("Client should not be nil")
	}

	expectedURL := "http://localhost:1234/"
	if client.BaseURL != expectedURL {
		t.Errorf("Expected baseURL %s, got %s", expectedURL, client.BaseURL)
	}

	if client.APIKey != apiKey {
		t.Errorf("Expected apiKey %s, got %s", apiKey, client.APIKey)
	}

	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}

	if len(client.Tools) != 0 {
		t.Errorf("Expected 0 tools initially, got %d", len(client.Tools))
	}
}

func TestNewClientWithTrailingSlash(t *testing.T) {
	baseURL := "http://localhost:1234/"
	client := NewClient(baseURL, "")

	if client.BaseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.BaseURL)
	}
}

func TestAddTool(t *testing.T) {
	client := NewClient("http://localhost:1234", "")

	tool := Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "test_tool",
			Description: "A test tool",
			Parameters:  make(map[string]interface{}),
		},
	}

	client.AddTool(tool)

	if len(client.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(client.Tools))
	}

	if client.Tools[0].Function.Name != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got %s", client.Tools[0].Function.Name)
	}
}

func TestChatCompletionRequest(t *testing.T) {
	req := ChatCompletionRequest{
		Model:       "test-model",
		Messages:    []Message{{Role: "user", Content: "Hello"}},
		Temperature: 0.7,
		MaxTokens:   100,
		Stream:      false,
	}

	if req.Model != "test-model" {
		t.Errorf("Expected model 'test-model', got %s", req.Model)
	}

	if len(req.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(req.Messages))
	}

	if req.Messages[0].Content != "Hello" {
		t.Errorf("Expected message content 'Hello', got %s", req.Messages[0].Content)
	}

	if req.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", req.Temperature)
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Test message",
	}

	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got %s", msg.Role)
	}

	if msg.Content != "Test message" {
		t.Errorf("Expected content 'Test message', got %s", msg.Content)
	}
}

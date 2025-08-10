package agentic

import (
	"context"
	"testing"

	"github.com/vivesm/GOSS-CLI/agentic-cli/openai"
)

func TestNewChatSession(t *testing.T) {
	config := SessionConfig{
		BaseURL:     "http://localhost:1234/v1",
		APIKey:      "",
		Model:       "test-model",
		Temperature: 0.7,
		MaxTokens:   2048,
	}

	session, err := NewChatSession(context.Background(), config)
	if err != nil {
		t.Fatalf("NewChatSession failed: %v", err)
	}

	if session == nil {
		t.Fatal("Session should not be nil")
	}

	if session.model != config.Model {
		t.Errorf("Expected model %s, got %s", config.Model, session.model)
	}

	if len(session.history) != 0 {
		t.Errorf("Expected empty history, got %d messages", len(session.history))
	}
}

func TestNewChatSessionWithDefaultModel(t *testing.T) {
	config := SessionConfig{
		BaseURL:     "http://localhost:1234/v1",
		APIKey:      "",
		Model:       "", // Empty model should use default
		Temperature: 0.7,
		MaxTokens:   2048,
	}

	session, err := NewChatSession(context.Background(), config)
	if err != nil {
		t.Fatalf("NewChatSession failed: %v", err)
	}

	if session.model != DefaultModel {
		t.Errorf("Expected default model %s, got %s", DefaultModel, session.model)
	}
}

func TestChatSessionGetModel(t *testing.T) {
	config := SessionConfig{
		BaseURL: "http://localhost:1234/v1",
		Model:   "test-model",
	}

	session, err := NewChatSession(context.Background(), config)
	if err != nil {
		t.Fatalf("NewChatSession failed: %v", err)
	}

	model := session.GetModel()
	if model != "test-model" {
		t.Errorf("Expected model 'test-model', got %s", model)
	}
}

func TestChatSessionSetModel(t *testing.T) {
	config := SessionConfig{
		BaseURL: "http://localhost:1234/v1",
		Model:   "initial-model",
	}

	session, err := NewChatSession(context.Background(), config)
	if err != nil {
		t.Fatalf("NewChatSession failed: %v", err)
	}

	newModel := "new-model"
	session.SetModel(newModel)

	if session.GetModel() != newModel {
		t.Errorf("Expected model %s, got %s", newModel, session.GetModel())
	}
}

func TestChatSessionHistoryOperations(t *testing.T) {
	config := SessionConfig{
		BaseURL: "http://localhost:1234/v1",
		Model:   "test-model",
	}

	session, err := NewChatSession(context.Background(), config)
	if err != nil {
		t.Fatalf("NewChatSession failed: %v", err)
	}

	// Test initial empty history
	history := session.GetHistory()
	if len(history) != 0 {
		t.Errorf("Expected empty history, got %d messages", len(history))
	}

	// Test setting history
	testHistory := []openai.Message{
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there!"},
	}
	session.SetHistory(testHistory)

	history = session.GetHistory()
	if len(history) != 2 {
		t.Errorf("Expected 2 messages in history, got %d", len(history))
	}

	if history[0].Content != "Hello" {
		t.Errorf("Expected first message content 'Hello', got %s", history[0].Content)
	}

	// Test clearing history
	session.ClearHistory()
	history = session.GetHistory()
	if len(history) != 0 {
		t.Errorf("Expected empty history after clear, got %d messages", len(history))
	}
}

func TestAgenticResponseFormatResponse(t *testing.T) {
	// Test response without tool calls
	response := &AgenticResponse{
		Content:      "Hello, world!",
		ToolCalls:    false,
		FinishReason: "stop",
		Usage:        openai.Usage{},
	}

	formatted := response.FormatResponse()
	expected := "Hello, world!"
	if formatted != expected {
		t.Errorf("Expected %q, got %q", expected, formatted)
	}

	// Test response with tool calls
	response.ToolCalls = true
	formatted = response.FormatResponse()
	expected = "Hello, world!\n\n[Tools were used to generate this response]"
	if formatted != expected {
		t.Errorf("Expected %q, got %q", expected, formatted)
	}
}
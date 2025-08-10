package handler

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/config"
)

// AgenticQuery processes queries to agentic models with MCP tools.
// It implements the MessageHandler interface.
type AgenticQuery struct {
	*IO
	session  *agentic.ChatSession
	renderer *glamour.TermRenderer
	config   *config.Config
}

var _ MessageHandler = (*AgenticQuery)(nil)

// NewAgenticQuery returns a new AgenticQuery message handler.
func NewAgenticQuery(io *IO, session *agentic.ChatSession, config *config.Config, opts RendererOptions) (*AgenticQuery, error) {
	renderer, err := opts.newTermRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate terminal renderer: %w", err)
	}

	return &AgenticQuery{
		IO:       io,
		session:  session,
		renderer: renderer,
		config:   config,
	}, nil
}

// Handle processes the chat message with agentic capabilities
func (h *AgenticQuery) Handle(message string) (Response, bool) {
	// Check if streaming is enabled
	if h.config.Streaming.Enabled {
		return h.handleStreaming(message)
	}
	return h.handleNonStreaming(message)
}

// handleStreaming processes the message with real-time streaming
func (h *AgenticQuery) handleStreaming(message string) (Response, bool) {
	// Debug output
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] handleStreaming called for message: %s\n", message)
	}
	
	// Clear any spinner since we're streaming  
	// Note: Skip spinner.Stop() as it causes hangs - let process exit handle cleanup
	// h.terminal.Spinner.Stop()
	
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Skipped spinner stop, creating callback...\n")
	}
	
	streamCallback := func(content string, isThinking bool) error {
		if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
			fmt.Fprintf(os.Stderr, "[DEBUG] Stream callback: isThinking=%v, content=%.50s...\n", isThinking, content)
		}
		if isThinking && h.config.Streaming.ShowThinking {
			// Show thinking tokens in a different color/style
			h.terminal.Write("\033[90m" + content + "\033[0m") // Gray text for thinking
		} else if !isThinking {
			// Show regular response tokens
			h.terminal.Write(content)
		}
		return nil
	}
	
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] About to call SendMessageStream...\n")
	}
	
	_, err := h.session.SendMessageStream(message, h.config.Streaming.ThinkingLevel, h.config.Streaming.ShowThinking, streamCallback)
	if err != nil {
		if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
			fmt.Fprintf(os.Stderr, "[DEBUG] SendMessageStream returned error: %v\n", err)
		}
		return newErrorResponse(err), false
	}
	
	if debugMode := os.Getenv("GOSS_DEBUG"); debugMode != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] SendMessageStream completed successfully\n")
	}
	
	// Add a newline after streaming
	h.terminal.Write("\n")
	
	// Since we already showed the content during streaming, we don't need to render again
	// Just return success status
	return dataResponse(""), false
}

// handleNonStreaming processes the message with traditional spinner approach
func (h *AgenticQuery) handleNonStreaming(message string) (Response, bool) {
	h.terminal.Spinner.Start()
	defer h.terminal.Spinner.Stop()

	response, err := h.session.SendMessage(message)
	if err != nil {
		return newErrorResponse(err), false
	}

	// Format the response content
	content := response.FormatResponse()

	// Render markdown
	rendered, err := h.renderer.Render(content)
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to format response: %w", err)), false
	}

	return dataResponse(rendered), false
}

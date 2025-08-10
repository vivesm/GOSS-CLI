package handler

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
)

// AgenticQuery processes queries to agentic models with MCP tools.
// It implements the MessageHandler interface.
type AgenticQuery struct {
	*IO
	session  *agentic.ChatSession
	renderer *glamour.TermRenderer
}

var _ MessageHandler = (*AgenticQuery)(nil)

// NewAgenticQuery returns a new AgenticQuery message handler.
func NewAgenticQuery(io *IO, session *agentic.ChatSession, opts RendererOptions) (*AgenticQuery, error) {
	renderer, err := opts.newTermRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate terminal renderer: %w", err)
	}

	return &AgenticQuery{
		IO:       io,
		session:  session,
		renderer: renderer,
	}, nil
}

// Handle processes the chat message with agentic capabilities
func (h *AgenticQuery) Handle(message string) (Response, bool) {
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

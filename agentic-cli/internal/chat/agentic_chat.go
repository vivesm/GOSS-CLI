package chat

import (
	"github.com/reugn/gemini-cli/agentic"
	"github.com/reugn/gemini-cli/internal/config"
	"github.com/reugn/gemini-cli/internal/handler"
	"github.com/reugn/gemini-cli/internal/terminal"
)

// NewAgentic returns a new Chat with agentic capabilities
func NewAgentic(
	user string, session *agentic.ChatSession,
	configuration *config.Configuration, opts *Opts,
) (*Chat, error) {
	terminalIOConfig := &terminal.IOConfig{
		User:           user,
		Multiline:      opts.Multiline,
		LineTerminator: opts.LineTerminator,
	}

	terminalIO, err := terminal.NewIO(terminalIOConfig)
	if err != nil {
		return nil, err
	}

	// Create agentic query handler
	agenticIO := handler.NewIO(terminalIO, terminalIO.Prompt.Gemini)
	agenticHandler, err := handler.NewAgenticQuery(agenticIO, session, opts.rendererOptions())
	if err != nil {
		return nil, err
	}

	// Create agentic system handler
	systemIO := handler.NewIO(terminalIO, terminalIO.Prompt.Cli)
	systemHandler, err := handler.NewAgenticSystem(systemIO, session, configuration,
		opts.GenerativeModel, opts.rendererOptions())
	if err != nil {
		return nil, err
	}

	return &Chat{
		io:            terminalIO,
		geminiHandler: agenticHandler,  // Re-use same field name for compatibility
		systemHandler: systemHandler,
	}, nil
}
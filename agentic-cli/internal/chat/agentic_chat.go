package chat

import (
	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/config"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/handler"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/terminal"
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
	agenticIO := handler.NewIO(terminalIO, terminalIO.Prompt.Goss)
	agenticHandler, err := handler.NewAgenticQuery(agenticIO, session, opts.rendererOptions())
	if err != nil {
		return nil, err
	}

	// Create agentic system handler
	systemIO := handler.NewIO(terminalIO, terminalIO.Prompt.System)
	systemHandler, err := handler.NewSystem(systemIO, session, configuration,
		opts.GenerativeModel, opts.rendererOptions())
	if err != nil {
		return nil, err
	}

	return &Chat{
		io:            terminalIO,
		gossHandler:   agenticHandler, // Updated field name for GOSS
		systemHandler: systemHandler,
	}, nil
}

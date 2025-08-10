package handler

import (
	"github.com/reugn/gemini-cli/agentic"
	"github.com/reugn/gemini-cli/internal/cli"
	"github.com/reugn/gemini-cli/internal/config"
)

// AgenticSystem processes chat system commands for agentic sessions
// It implements the MessageHandler interface.
type AgenticSystem struct {
	*IO
	handlers map[string]MessageHandler
}

var _ MessageHandler = (*AgenticSystem)(nil)

// NewAgenticSystem returns a new AgenticSystem command handler.
func NewAgenticSystem(io *IO, session *agentic.ChatSession, configuration *config.Configuration,
	modelName string, rendererOptions RendererOptions) (*AgenticSystem, error) {
	helpCommandHandler, err := NewHelpCommand(io, rendererOptions)
	if err != nil {
		return nil, err
	}

	handlers := map[string]MessageHandler{
		cli.SystemCmdHelp:            helpCommandHandler,
		cli.SystemCmdQuit:            NewQuitCommand(io),
		cli.SystemCmdSelectPrompt:    NewAgenticPromptCommand(io, session, configuration.Data),
		cli.SystemCmdSelectInputMode: NewInputModeCommand(io),
		cli.SystemCmdModel:           NewAgenticModelCommand(io, session, modelName),
		cli.SystemCmdHistory:         NewAgenticHistoryCommand(io, session, configuration),
		cli.SystemCmdTemperature:     NewAgenticTemperatureCommand(io, session),
	}

	return &AgenticSystem{
		IO:       io,
		handlers: handlers,
	}, nil
}

// Handle processes the chat system command.
func (s *AgenticSystem) Handle(message string) (Response, bool) {
	commandName, err := cli.ExtractSystemCommandName(message)
	if err != nil {
		return newErrorResponse(err), false
	}
	handler, found := s.handlers[commandName]
	if !found {
		return newErrorResponse(cli.ErrInvalidSystemCommand), false
	}
	return handler.Handle(message)
}
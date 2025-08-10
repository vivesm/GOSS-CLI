package handler

import (
	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/cli"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/config"
)

// System processes chat system commands for sessions
// It implements the MessageHandler interface.
type System struct {
	BaseCommand
	handlers map[string]MessageHandler
}

var _ MessageHandler = (*System)(nil)

// NewSystem returns a new System command handler.
func NewSystem(io *IO, session *agentic.ChatSession, configuration *config.Config,
	modelName string, rendererOptions RendererOptions) (*System, error) {
	helpCommandHandler, err := NewHelpCommand(io, rendererOptions)
	if err != nil {
		return nil, err
	}

	handlers := map[string]MessageHandler{
		cli.SystemCmdHelp:            helpCommandHandler,
		cli.SystemCmdQuit:            NewQuitCommand(io),
		cli.SystemCmdSelectPrompt:    NewPromptCommand(io, session, configuration),
		cli.SystemCmdSelectInputMode: NewInputModeCommand(io),
		cli.SystemCmdModel:           NewModelCommand(io, session, modelName),
		cli.SystemCmdHistory:         NewHistoryCommand(io, session, configuration),
		cli.SystemCmdTemperature:     NewTemperatureCommand(io, session),
		cli.SystemCmdStream:          NewStreamCommand(io, configuration),
		cli.SystemCmdThinking:        NewThinkingCommand(io, configuration),
		cli.SystemCmdShowThinking:    NewShowThinkingCommand(io, configuration),
	}

	return &System{
		BaseCommand: NewBaseCommand(io),
		handlers:    handlers,
	}, nil
}

// Handle processes the chat system command.
func (s *System) Handle(message string) (Response, bool) {
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

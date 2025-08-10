package handler

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/config"
)

var inputModeOptions = []string{
	"Single-line",
	"Multi-line",
}

// =============================================================================
// PROMPT COMMAND
// =============================================================================

// PromptCommand handles system prompt operations for sessions
type PromptCommand struct {
	BaseCommand
	session *agentic.ChatSession
	config  *config.Config
}

var _ MessageHandler = (*PromptCommand)(nil)

// NewPromptCommand returns a new PromptCommand
func NewPromptCommand(io *IO, session *agentic.ChatSession, config *config.Config) *PromptCommand {
	return &PromptCommand{
		BaseCommand: NewBaseCommand(io),
		session:     session,
		config:      config,
	}
}

// Handle processes system prompt commands
func (p *PromptCommand) Handle(_ string) (Response, bool) {
	if p.config.SystemPrompts == nil || len(p.config.SystemPrompts) == 0 {
		return dataResponse("No system prompts configured. Add them to your configuration file."), false
	}

	var prompts []string
	var promptMap = make(map[string]string)

	for name, prompt := range p.config.SystemPrompts {
		prompts = append(prompts, name)
		promptMap[name] = prompt
	}

	prompt := promptui.Select{
		Label: "Select system prompt",
		Items: prompts,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return newErrorResponse(err), false
	}

	selectedPrompt := promptMap[result]

	// For agentic sessions, we add the system prompt as a system message
	// Clear history and start with system prompt
	p.session.ClearHistory()

	// Note: In a real implementation, you might want to implement system message support
	// For now, we'll just inform the user that the prompt was selected
	return dataResponse(fmt.Sprintf("System prompt selected: %s\n\nPrompt: %s\n\nNote: System prompts are informational only in this implementation.", result, selectedPrompt)), false
}

// =============================================================================
// INPUT MODE COMMAND
// =============================================================================

// InputModeCommand processes the chat input mode system command.
// It implements the MessageHandler interface.
type InputModeCommand struct {
	BaseCommand
}

var _ MessageHandler = (*InputModeCommand)(nil)

// NewInputModeCommand returns a new InputModeCommand.
func NewInputModeCommand(io *IO) *InputModeCommand {
	return &InputModeCommand{
		BaseCommand: NewBaseCommand(io),
	}
}

// Handle processes the chat input mode system command.
func (h *InputModeCommand) Handle(_ string) (Response, bool) {
	defer h.IO.terminal.Write(h.IO.terminalPrompt)
	multiline, err := h.selectInputMode()
	if err != nil {
		return newErrorResponse(err), false
	}

	if h.IO.terminal.Config.Multiline == multiline {
		// the same input mode is selected
		return dataResponse(unchangedMessage), false
	}

	h.IO.terminal.Config.Multiline = multiline
	h.IO.terminal.SetUserPrompt()
	if h.IO.terminal.Config.Multiline {
		// disable history for multi-line messages since it is
		// unusable for future requests
		h.IO.terminal.Reader.HistoryDisable()
	} else {
		h.IO.terminal.Reader.HistoryEnable()
	}

	mode := inputModeOptions[modeIndex(h.IO.terminal.Config.Multiline)]
	return dataResponse(fmt.Sprintf("Switched to %q input mode.", mode)), false
}

// selectInputMode returns true if multiline input is selected;
// otherwise, it returns false.
func (h *InputModeCommand) selectInputMode() (bool, error) {
	prompt := promptui.Select{
		Label:        "Select input mode",
		HideSelected: true,
		Items:        inputModeOptions,
		CursorPos:    modeIndex(h.IO.terminal.Config.Multiline),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == inputModeOptions[1], nil
}

func modeIndex(b bool) int {
	if b {
		return 1
	}
	return 0
}

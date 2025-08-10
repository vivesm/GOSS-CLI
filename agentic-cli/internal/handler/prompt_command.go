package handler

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/config"
)

// PromptCommand handles system prompt operations for sessions
type PromptCommand struct {
	BaseCommand
	session *agentic.ChatSession
	data    *config.ApplicationData
}

var _ MessageHandler = (*PromptCommand)(nil)

// NewPromptCommand returns a new PromptCommand
func NewPromptCommand(io *IO, session *agentic.ChatSession, data *config.ApplicationData) *PromptCommand {
	return &PromptCommand{
		BaseCommand: NewBaseCommand(io),
		session:     session,
		data:        data,
	}
}

// Handle processes system prompt commands
func (p *PromptCommand) Handle(_ string) (Response, bool) {
	if p.data.SystemPrompts == nil || len(p.data.SystemPrompts) == 0 {
		return dataResponse("No system prompts configured. Add them to your configuration file."), false
	}

	var prompts []string
	var promptMap = make(map[string]string)

	for name, prompt := range p.data.SystemPrompts {
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

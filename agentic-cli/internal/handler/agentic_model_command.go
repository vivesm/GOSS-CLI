package handler

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/reugn/gemini-cli/agentic"
)

// AgenticModelCommand handles model operations for agentic sessions
type AgenticModelCommand struct {
	*IO
	session   *agentic.ChatSession
	modelName string
}

var _ MessageHandler = (*AgenticModelCommand)(nil)

// NewAgenticModelCommand returns a new AgenticModelCommand
func NewAgenticModelCommand(io *IO, session *agentic.ChatSession, modelName string) *AgenticModelCommand {
	return &AgenticModelCommand{
		IO:        io,
		session:   session,
		modelName: modelName,
	}
}

// Handle processes model-related commands
func (m *AgenticModelCommand) Handle(_ string) (Response, bool) {
	items := []string{
		"Select model",
		"Show model information",
		"List tools available",
	}

	prompt := promptui.Select{
		Label: "Model operations",
		Items: items,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return newErrorResponse(err), false
	}

	switch index {
	case 0:
		return m.selectModel(), false
	case 1:
		return m.showModelInfo(), false
	case 2:
		return m.listTools(), false
	default:
		return newErrorResponse(fmt.Errorf("invalid selection")), false
	}
}

func (m *AgenticModelCommand) selectModel() Response {
	models, err := m.session.ListModels()
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to list models: %w", err))
	}

	if len(models) == 0 {
		return dataResponse("No models available")
	}

	// Find current model index
	selectedIndex := 0
	currentModel := m.session.GetModel()
	for i, model := range models {
		if model == currentModel {
			selectedIndex = i
			break
		}
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "{{ . | cyan }}",
		Inactive: "{{ . }}",
		Selected: "{{ . | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select model",
		Items:     models,
		Templates: templates,
		CursorPos: selectedIndex,
	}

	index, result, err := prompt.Run()
	if err != nil {
		return newErrorResponse(err)
	}

	if index == selectedIndex {
		return dataResponse(unchangedMessage)
	}

	m.session.SetModel(result)
	m.modelName = result
	return dataResponse(fmt.Sprintf("Selected model: %s", result))
}

func (m *AgenticModelCommand) showModelInfo() Response {
	modelInfo, err := m.session.ModelInfo()
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to get model info: %w", err))
	}

	return dataResponse(modelInfo)
}

func (m *AgenticModelCommand) listTools() Response {
	modelInfo, err := m.session.ModelInfo()
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to get model info: %w", err))
	}

	return dataResponse(fmt.Sprintf("Available MCP Tools:\n%s", modelInfo))
}
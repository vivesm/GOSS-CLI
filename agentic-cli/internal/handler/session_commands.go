package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/config"
	"github.com/vivesm/GOSS-CLI/agentic-cli/openai"
)

// =============================================================================
// MODEL COMMAND
// =============================================================================

// ModelCommand handles model operations for sessions
type ModelCommand struct {
	BaseCommand
	session   *agentic.ChatSession
	modelName string
}

var _ MessageHandler = (*ModelCommand)(nil)

// NewModelCommand returns a new ModelCommand
func NewModelCommand(io *IO, session *agentic.ChatSession, modelName string) *ModelCommand {
	return &ModelCommand{
		BaseCommand: NewBaseCommand(io),
		session:     session,
		modelName:   modelName,
	}
}

// Handle processes model-related commands
func (m *ModelCommand) Handle(_ string) (Response, bool) {
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

func (m *ModelCommand) selectModel() Response {
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

func (m *ModelCommand) showModelInfo() Response {
	modelInfo, err := m.session.ModelInfo()
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to get model info: %w", err))
	}

	return dataResponse(modelInfo)
}

func (m *ModelCommand) listTools() Response {
	modelInfo, err := m.session.ModelInfo()
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to get model info: %w", err))
	}

	return dataResponse(fmt.Sprintf("Available MCP Tools:\n%s", modelInfo))
}

// =============================================================================
// TEMPERATURE COMMAND
// =============================================================================

// TemperatureCommand processes temperature control system commands.
// It implements the MessageHandler interface.
type TemperatureCommand struct {
	BaseCommand
	session *agentic.ChatSession
}

var _ MessageHandler = (*TemperatureCommand)(nil)

// NewTemperatureCommand returns a new TemperatureCommand.
func NewTemperatureCommand(io *IO, session *agentic.ChatSession) *TemperatureCommand {
	return &TemperatureCommand{
		BaseCommand: NewBaseCommand(io),
		session:     session,
	}
}

// Handle processes the temperature control command.
func (tc *TemperatureCommand) Handle(message string) (Response, bool) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		return tc.showTemperatureHelp(), false
	}

	subcommand := parts[1]

	switch subcommand {
	case "show", "get":
		return tc.showTemperature(), false

	case "set":
		if len(parts) < 3 {
			return dataResponse("❌ Usage: !t set <value>\nExample: !t set 0.7"), false
		}

		tempStr := parts[2]
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			return dataResponse(fmt.Sprintf("❌ Invalid temperature value: %s\nTemperature must be a number between 0.0 and 2.0", tempStr)), false
		}

		if temp < 0.0 || temp > 2.0 {
			return dataResponse("❌ Temperature must be between 0.0 and 2.0\n• 0.0-0.3: Very focused, deterministic\n• 0.4-0.7: Balanced\n• 0.8-2.0: More creative, random"), false
		}

		oldTemp := tc.session.GetTemperature()
		tc.session.SetTemperature(temp)

		return dataResponse(fmt.Sprintf("✅ Temperature updated: %.2f → %.2f\n%s",
			oldTemp, temp, tc.getTemperatureDescription(temp))), false

	case "reset":
		tc.session.SetTemperature(0.3)
		return dataResponse("✅ Temperature reset to default: 0.3 (focused reasoning)"), false

	default:
		return tc.showTemperatureHelp(), false
	}
}

func (tc *TemperatureCommand) showTemperature() Response {
	temp := tc.session.GetTemperature()
	maxTokens := tc.session.GetMaxTokens()

	response := fmt.Sprintf("🌡️  **Current Settings:**\n")
	response += fmt.Sprintf("• Temperature: %.2f %s\n", temp, tc.getTemperatureDescription(temp))
	response += fmt.Sprintf("• Max Tokens: %d\n", maxTokens)

	return dataResponse(response)
}

func (tc *TemperatureCommand) showTemperatureHelp() Response {
	help := `🌡️  **Temperature Control Commands:**

**Usage:**
• !t show        - Show current temperature
• !t set <value> - Set temperature (0.0-2.0)
• !t reset       - Reset to default (0.3)

**Temperature Guide:**
• 0.0-0.3 - Very focused, deterministic responses
• 0.4-0.7 - Balanced creativity and focus  
• 0.8-2.0 - More creative and random responses

**Examples:**
• !t set 0.1 - Maximum focus for analysis
• !t set 0.7 - Balanced for general use
• !t set 1.0 - More creative for brainstorming`

	return dataResponse(help)
}

func (tc *TemperatureCommand) getTemperatureDescription(temp float64) string {
	if temp <= 0.3 {
		return "(🎯 focused)"
	} else if temp <= 0.7 {
		return "(⚖️ balanced)"
	} else {
		return "(🎨 creative)"
	}
}

// =============================================================================
// HISTORY COMMAND
// =============================================================================

// HistoryCommand handles history operations for sessions
type HistoryCommand struct {
	BaseCommand
	session *agentic.ChatSession
	config  *config.Config
}

var _ MessageHandler = (*HistoryCommand)(nil)

// NewHistoryCommand returns a new HistoryCommand
func NewHistoryCommand(io *IO, session *agentic.ChatSession, config *config.Config) *HistoryCommand {
	return &HistoryCommand{
		BaseCommand: NewBaseCommand(io),
		session:     session,
		config:      config,
	}
}

// Handle processes history-related commands
func (h *HistoryCommand) Handle(_ string) (Response, bool) {
	items := []string{
		"Clear history",
		"Save history",
		"Load history",
		"Delete all history records",
	}

	prompt := promptui.Select{
		Label: "History operations",
		Items: items,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return newErrorResponse(err), false
	}

	switch index {
	case 0:
		return h.clearHistory(), false
	case 1:
		return h.saveHistory(), false
	case 2:
		return h.loadHistory(), false
	case 3:
		return h.deleteAllHistory(), false
	default:
		return newErrorResponse(fmt.Errorf("invalid selection")), false
	}
}

func (h *HistoryCommand) clearHistory() Response {
	h.session.ClearHistory()
	return dataResponse("History cleared")
}

func (h *HistoryCommand) saveHistory() Response {
	history := h.session.GetHistory()
	if len(history) == 0 {
		return dataResponse("No history to save")
	}

	// Generate a timestamp-based key
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	key := fmt.Sprintf("session_%s", timestamp)

	// Convert to JSON
	historyJSON, err := json.Marshal(history)
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to serialize history: %w", err))
	}

	// Save to configuration
	if h.config.History == nil {
		h.config.History = make(map[string]interface{})
	}
	h.config.History[key] = string(historyJSON)

	// Write configuration
	if err := h.config.Save(); err != nil {
		return newErrorResponse(fmt.Errorf("failed to save history: %w", err))
	}

	return dataResponse(fmt.Sprintf("History saved as: %s", key))
}

func (h *HistoryCommand) loadHistory() Response {
	if h.config.History == nil || len(h.config.History) == 0 {
		return dataResponse("No saved history available")
	}

	// Get list of saved histories
	var keys []string
	for key := range h.config.History {
		keys = append(keys, key)
	}

	prompt := promptui.Select{
		Label: "Select history to load",
		Items: keys,
	}

	_, selectedKey, err := prompt.Run()
	if err != nil {
		return newErrorResponse(err)
	}

	// Load the selected history
	historyData, exists := h.config.History[selectedKey]
	if !exists {
		return newErrorResponse(fmt.Errorf("history not found: %s", selectedKey))
	}

	historyJSON, ok := historyData.(string)
	if !ok {
		return newErrorResponse(fmt.Errorf("invalid history format"))
	}

	var history []openai.Message
	if err := json.Unmarshal([]byte(historyJSON), &history); err != nil {
		return newErrorResponse(fmt.Errorf("failed to deserialize history: %w", err))
	}

	// Set the loaded history
	h.session.SetHistory(history)

	return dataResponse(fmt.Sprintf("History loaded: %s (%d messages)", selectedKey, len(history)))
}

func (h *HistoryCommand) deleteAllHistory() Response {
	if h.config.History == nil || len(h.config.History) == 0 {
		return dataResponse("No history records to delete")
	}

	count := len(h.config.History)

	// Confirm deletion
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Delete all %d history records? (y/N)", count),
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil || result != "y" {
		return dataResponse("Deletion cancelled")
	}

	// Clear history
	h.config.History = make(map[string]interface{})

	// Write configuration
	if err := h.config.Save(); err != nil {
		return newErrorResponse(fmt.Errorf("failed to delete history: %w", err))
	}

	return dataResponse(fmt.Sprintf("Deleted %d history records", count))
}

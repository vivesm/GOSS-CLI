package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/config"
	"github.com/vivesm/GOSS-CLI/agentic-cli/openai"
)

// HistoryCommand handles history operations for sessions
type HistoryCommand struct {
	BaseCommand
	session       *agentic.ChatSession
	configuration *config.Configuration
}

var _ MessageHandler = (*HistoryCommand)(nil)

// NewHistoryCommand returns a new HistoryCommand
func NewHistoryCommand(io *IO, session *agentic.ChatSession, configuration *config.Configuration) *HistoryCommand {
	return &HistoryCommand{
		BaseCommand:   NewBaseCommand(io),
		session:       session,
		configuration: configuration,
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
	return dataResponse("Chat history cleared")
}

func (h *HistoryCommand) saveHistory() Response {
	history := h.session.GetHistory()
	if len(history) == 0 {
		return dataResponse("No history to save")
	}

	// Convert to JSON for storage
	historyJSON, err := json.Marshal(history)
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to marshal history: %w", err))
	}

	// Generate a unique key based on timestamp
	key := fmt.Sprintf("session_%d", time.Now().Unix())

	// Save to configuration
	if h.configuration.Data.History == nil {
		h.configuration.Data.History = make(map[string]interface{})
	}
	h.configuration.Data.History[key] = string(historyJSON)

	// Write configuration to file
	if err := h.configuration.Write(); err != nil {
		return newErrorResponse(fmt.Errorf("failed to save history: %w", err))
	}

	return dataResponse(fmt.Sprintf("History saved as: %s", key))
}

func (h *HistoryCommand) loadHistory() Response {
	if h.configuration.Data.History == nil || len(h.configuration.Data.History) == 0 {
		return dataResponse("No saved history records found")
	}

	// Create list of available history records
	var keys []string
	for key := range h.configuration.Data.History {
		keys = append(keys, key)
	}

	prompt := promptui.Select{
		Label: "Select history record to load",
		Items: keys,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return newErrorResponse(err)
	}

	// Load the selected history
	historyJSON, ok := h.configuration.Data.History[result].(string)
	if !ok {
		return newErrorResponse(fmt.Errorf("invalid history record format"))
	}

	var history []openai.Message
	if err := json.Unmarshal([]byte(historyJSON), &history); err != nil {
		return newErrorResponse(fmt.Errorf("failed to unmarshal history: %w", err))
	}

	h.session.SetHistory(history)
	return dataResponse(fmt.Sprintf("Loaded history: %s (%d messages)", result, len(history)))
}

func (h *HistoryCommand) deleteAllHistory() Response {
	if h.configuration.Data.History == nil || len(h.configuration.Data.History) == 0 {
		return dataResponse("No history records to delete")
	}

	count := len(h.configuration.Data.History)

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
	h.configuration.Data.History = make(map[string]interface{})

	// Write configuration
	if err := h.configuration.Write(); err != nil {
		return newErrorResponse(fmt.Errorf("failed to delete history: %w", err))
	}

	return dataResponse(fmt.Sprintf("Deleted %d history records", count))
}

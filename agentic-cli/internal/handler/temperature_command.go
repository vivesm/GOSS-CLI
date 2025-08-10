package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vivesm/GOSS-CLI/agentic-cli/agentic"
)

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

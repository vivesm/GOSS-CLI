package chat

import (
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/handler"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/terminal"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/terminal/color"
)

// Chat handles the interactive exchange of messages between user and model.
type Chat struct {
	io *terminal.IO

	gossHandler   handler.MessageHandler
	systemHandler handler.MessageHandler
}

// Start starts the chat.
func (c *Chat) Start() {
	defer func() {
		// Ensure spinner is stopped before any output
		if c.io.Spinner != nil {
			c.io.Spinner.Stop()
		}
		
		// Print goodbye message
		c.io.Write("\n")
		c.io.Write("╭─────────────────────────────────╮\n")
		c.io.Write("│  👋 Thanks for using GOSS AI!  │\n")
		c.io.Write("│     Have a wonderful day!       │\n")
		c.io.Write("╰─────────────────────────────────╯\n\n")
		
		// Properly close the IO to release readline resources
		c.io.Close()
	}()

	c.showWelcomeMessage()

	for {
		c.io.SetUserPrompt()
		input := c.io.Read()

		if input == "" {
			continue
		}

		// Handle system commands (prefixed with !)
		if isSystemCommand(input) {
			response, exit := c.systemHandler.Handle(input)
			c.printResponse(response)
			if exit {
				break
			}
		} else {
			// Handle regular agentic queries
			response, _ := c.gossHandler.Handle(input)
			c.printResponse(response)
		}
	}
}

// printResponse prints a response to the terminal
func (c *Chat) printResponse(response handler.Response) {
	c.io.Write(response.String())
}

// showWelcomeMessage displays an enhanced welcome message
func (c *Chat) showWelcomeMessage() {
	c.io.Write("\n")
	c.io.Write(color.Cyan("╭─────────────────────────────────────────────────────────────────╮\n"))
	c.io.Write(color.Cyan("│                                                                 │\n"))
	c.io.Write(color.Cyan("│   ") + color.Green("🎯 GOSS AI Assistant") + color.Cyan("                                         │\n"))
	c.io.Write(color.Cyan("│   ") + color.Green("──────────────────────") + color.Cyan("                                       │\n"))
	c.io.Write(color.Cyan("│                                                                 │\n"))
	c.io.Write(color.Cyan("│   ") + color.Yellow("🔧 MCP Tools Available:") + color.Cyan("                                      │\n"))
	c.io.Write(color.Cyan("│      • ") + color.White("📁 File Operations (read, write, list, search)") + color.Cyan("         │\n"))
	c.io.Write(color.Cyan("│      • ") + color.White("🔍 Web Search (current information, weather, etc.)") + color.Cyan("     │\n"))
	c.io.Write(color.Cyan("│                                                                 │\n"))
	c.io.Write(color.Cyan("│   ") + color.Yellow("💡 System Commands:") + color.Cyan("                                          │\n"))
	c.io.Write(color.Cyan("│      ") + color.Blue("!help") + color.White("  - Show all commands    ") + color.Blue("!m") + color.White(" - Model info") + color.Cyan("            │\n"))
	c.io.Write(color.Cyan("│      ") + color.Blue("!h") + color.White("     - History management   ") + color.Blue("!t") + color.White(" - Temperature") + color.Cyan("           │\n"))
	c.io.Write(color.Cyan("│      ") + color.Blue("!q") + color.White("     - Quit application") + color.Cyan("                                │\n"))
	c.io.Write(color.Cyan("│                                                                 │\n"))
	c.io.Write(color.Cyan("╰─────────────────────────────────────────────────────────────────╯\n"))
	c.io.Write("\n")
	c.io.Write(color.Green("Ready to assist! ") + "Ask me anything or try:\n")
	c.io.Write(color.Gray("• \"List files in current directory\"\n"))
	c.io.Write(color.Gray("• \"Search for recent AI news\"\n"))
	c.io.Write(color.Gray("• \"Create a file with today's tasks\"\n\n"))
}

// isSystemCommand checks if the input is a system command
func isSystemCommand(input string) bool {
	return len(input) > 0 && input[0] == '!'
}

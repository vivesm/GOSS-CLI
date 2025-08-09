package chat

import (
	"github.com/reugn/gemini-cli/internal/handler"
	"github.com/reugn/gemini-cli/internal/terminal"
)

// Chat handles the interactive exchange of messages between user and model.
type Chat struct {
	io *terminal.IO

	geminiHandler handler.MessageHandler
	systemHandler handler.MessageHandler
}

// Start starts the chat.
func (c *Chat) Start() {
	defer func() {
		c.io.Spinner.Stop()
		c.io.Write("Have a nice day!\n")
	}()

	c.io.Write("🤖 Agentic CLI - Ready to assist!\n\n")

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
			response, _ := c.geminiHandler.Handle(input)
			c.printResponse(response)
		}
	}
}

// printResponse prints a response to the terminal
func (c *Chat) printResponse(response handler.Response) {
	c.io.Write(response.String())
}

// isSystemCommand checks if the input is a system command
func isSystemCommand(input string) bool {
	return len(input) > 0 && input[0] == '!'
}
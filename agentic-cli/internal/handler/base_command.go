package handler

// BaseCommand provides common functionality for all command handlers.
// It embeds the IO struct that was previously duplicated across all commands.
type BaseCommand struct {
	*IO
}

// NewBaseCommand creates a new BaseCommand with the provided IO.
// This eliminates the need for each command to embed *IO directly.
func NewBaseCommand(io *IO) BaseCommand {
	return BaseCommand{IO: io}
}

// TerminalPrompt returns the terminal prompt from the embedded IO.
// Commands can override this method to provide custom prompts.
func (b *BaseCommand) TerminalPrompt() string {
	return b.IO.TerminalPrompt()
}

package handler

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/vivesm/GOSS-CLI/agentic-cli/internal/cli"
)

// =============================================================================
// HELP COMMAND
// =============================================================================

// HelpCommand handles the help system command request.
type HelpCommand struct {
	BaseCommand
	renderer *glamour.TermRenderer
}

var _ MessageHandler = (*HelpCommand)(nil)

// NewHelpCommand returns a new HelpCommand.
func NewHelpCommand(io *IO, opts RendererOptions) (*HelpCommand, error) {
	renderer, err := opts.newTermRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate terminal renderer: %w", err)
	}

	return &HelpCommand{
		BaseCommand: NewBaseCommand(io),
		renderer:    renderer,
	}, nil
}

// Handle processes the help system command.
func (h *HelpCommand) Handle(_ string) (Response, bool) {
	var b strings.Builder
	b.WriteString("# System commands\n")
	b.WriteString("Use a command prefixed with an exclamation mark (e.g., `!h`).\n")
	fmt.Fprintf(&b, "* `%s` - Select the generative model system prompt.\n", cli.SystemCmdSelectPrompt)
	fmt.Fprintf(&b, "* `%s` - Select from a list of generative model operations.\n", cli.SystemCmdModel)
	fmt.Fprintf(&b, "* `%s` - Select from a list of chat history operations.\n", cli.SystemCmdHistory)
	fmt.Fprintf(&b, "* `%s` - Control temperature settings (focus vs creativity).\n", cli.SystemCmdTemperature)
	fmt.Fprintf(&b, "* `%s` - Toggle the input mode.\n", cli.SystemCmdSelectInputMode)
	fmt.Fprintf(&b, "* `%s` - Exit the application.\n", cli.SystemCmdQuit)

	rendered, err := h.renderer.Render(b.String())
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to format instructions: %w", err)), false
	}

	return dataResponse(rendered), false
}

// =============================================================================
// QUIT COMMAND
// =============================================================================

// QuitCommand processes the chat quit system command.
// It implements the MessageHandler interface.
type QuitCommand struct {
	BaseCommand
}

var _ MessageHandler = (*QuitCommand)(nil)

// NewQuitCommand returns a new QuitCommand.
func NewQuitCommand(io *IO) *QuitCommand {
	return &QuitCommand{
		BaseCommand: NewBaseCommand(io),
	}
}

// Handle processes the chat quit command.
func (h *QuitCommand) Handle(_ string) (Response, bool) {
	return dataResponse("Exiting goss..."), true
}

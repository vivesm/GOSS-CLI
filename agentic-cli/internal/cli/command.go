package cli

import (
	"errors"
	"strings"
)

const (
	SystemCmdPrefix          = "!"
	SystemCmdHelp            = "help"
	SystemCmdQuit            = "q"
	SystemCmdSelectPrompt    = "p"
	SystemCmdSelectInputMode = "i"
	SystemCmdModel           = "m"
	SystemCmdHistory         = "h"
	SystemCmdTemperature     = "t"
)

var ErrInvalidSystemCommand = errors.New("invalid system command")

// ExtractSystemCommandName extracts the command name from a system command input
func ExtractSystemCommandName(input string) (string, error) {
	if !strings.HasPrefix(input, SystemCmdPrefix) {
		return "", ErrInvalidSystemCommand
	}

	command := strings.TrimPrefix(input, SystemCmdPrefix)
	command = strings.TrimSpace(command)

	// Extract just the command name (first word)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", ErrInvalidSystemCommand
	}

	return parts[0], nil
}

package terminal

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"
	"github.com/reugn/gemini-cli/internal/terminal/color"
)

const (
	geminiUser = "goss"
	cliUser    = "system"
)

type Prompt struct {
	User              string
	UserMultiline     string
	UserMultilineNext string
	Goss              string  // AI assistant prompt
	System            string  // System command prompt
}

type promptColor struct {
	user   func(string) string
	goss   func(string) string
	system func(string) string
}

func newPromptColor() *promptColor {
	if termenv.HasDarkBackground() {
		return &promptColor{
			user:   color.Cyan,
			goss:   color.Green,
			system: color.Yellow,
		}
	}
	return &promptColor{
		user:   color.Blue,
		goss:   color.Green,
		system: color.Magenta,
	}
}

func NewPrompt(currentUser string) *Prompt {
	maxLength := maxLength(currentUser, geminiUser, cliUser)
	pc := newPromptColor()
	return &Prompt{
		User:              pc.user(buildPrompt(currentUser, '>', maxLength)),
		UserMultiline:     pc.user(buildPrompt(currentUser, '#', maxLength)),
		UserMultilineNext: pc.user(buildPrompt(strings.Repeat(" ", len(currentUser)), '>', maxLength)),
		Goss:              pc.goss(buildPrompt(geminiUser, '>', maxLength)),
		System:            pc.system(buildPrompt(cliUser, '>', maxLength)),
	}
}

func maxLength(strings ...string) int {
	var maxLength int
	for _, s := range strings {
		length := len(s)
		if maxLength < length {
			maxLength = length
		}
	}
	return maxLength
}

func buildPrompt(user string, p byte, length int) string {
	// Create a more modern, clean prompt style
	if user == "goss" {
		return "🤖 "
	}
	if user == "system" {
		return "⚙️  "
	}
	return fmt.Sprintf("%s%c ", user, p)
}

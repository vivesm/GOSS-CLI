package main

import (
	"context"
	"os"
	"os/user"

	"github.com/reugn/gemini-cli/agentic"
	"github.com/reugn/gemini-cli/internal/chat"
	"github.com/reugn/gemini-cli/internal/config"
	"github.com/spf13/cobra"
)

const (
	version           = "0.4.0"
	apiKeyEnv         = "LMSTUDIO_API_KEY" //nolint:gosec
	defaultConfigPath = "goss_config.json"
	defaultBaseURL    = "http://localhost:1234/v1"
)

func run() int {
	rootCmd := &cobra.Command{
		Use:     "goss",
		Short:   "Chat with local LLMs using MCP tools",
		Version: version,
	}

	var opts chat.Opts
	var configPath string
	var baseURL string
	rootCmd.Flags().StringVarP(&opts.GenerativeModel, "model", "m", agentic.DefaultModel,
		"generative model name")
	rootCmd.Flags().BoolVar(&opts.Multiline, "multiline", false,
		"read input as a multi-line string")
	rootCmd.Flags().StringVarP(&opts.LineTerminator, "term", "t", "$",
		"multi-line input terminator")
	rootCmd.Flags().StringVarP(&opts.StylePath, "style", "s", "auto",
		"markdown format style (ascii, dark, light, pink, notty, dracula)")
	rootCmd.Flags().IntVarP(&opts.WordWrap, "wrap", "w", 80,
		"line length for response word wrapping")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", defaultConfigPath,
		"path to configuration file in JSON format")
	rootCmd.Flags().StringVarP(&baseURL, "base-url", "b", defaultBaseURL,
		"LM Studio API base URL")

	rootCmd.RunE = func(_ *cobra.Command, _ []string) error {
		configuration, err := config.NewConfiguration(configPath)
		if err != nil {
			return err
		}

		// Create agentic chat session
		sessionConfig := agentic.SessionConfig{
			BaseURL:     baseURL,
			APIKey:      os.Getenv(apiKeyEnv), // Optional for LM Studio
			Model:       opts.GenerativeModel,
			Temperature: 0.3, // Default focused temperature, changeable with !t
			MaxTokens:   2048,
		}
		
		chatSession, err := agentic.NewChatSession(context.Background(), sessionConfig)
		if err != nil {
			return err
		}

		chatHandler, err := chat.NewAgentic(getCurrentUser(), chatSession, configuration, &opts)
		if err != nil {
			return err
		}
		chatHandler.Start()

		return chatSession.Close()
	}

	err := rootCmd.Execute()
	if err != nil {
		return 1
	}
	return 0
}

func getCurrentUser() string {
	currentUser, err := user.Current()
	if err != nil {
		return "user"
	}
	return currentUser.Username
}

func main() {
	// start the application
	os.Exit(run())
}
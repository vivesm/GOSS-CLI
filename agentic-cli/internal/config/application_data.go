package config

// ApplicationData contains the configuration data for the agentic CLI.
type ApplicationData struct {
	SystemPrompts map[string]string      `json:"SystemPrompts"`
	History       map[string]interface{} `json:"History"`
}

// newDefaultApplicationData returns a new ApplicationData with default values.
func newDefaultApplicationData() *ApplicationData {
	return &ApplicationData{
		SystemPrompts: map[string]string{
			"Assistant": "You are a helpful AI assistant with access to filesystem and web search tools. Use the tools when needed to provide accurate and helpful responses.",
			"Developer": "You are an expert software developer with access to filesystem and web search tools. Help with coding tasks, debugging, and software development questions. Use file operations to read code, search for patterns, and create or modify files as needed.",
			"Researcher": "You are a research assistant with web search capabilities. Help find information online and use filesystem tools to organize research findings into files.",
			"Writer": "You are a writing assistant that can help with document creation and editing. Use filesystem tools to read, write, and organize documents.",
		},
		History: make(map[string]interface{}),
	}
}
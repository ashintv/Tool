package agents

// CommanderConfig holds configuration settings for the agent commander.
// Deprecated: Moving to interface-based configs for better testability.
type CommanderConfig struct {
	//u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/agent-ws"}
	WsServerHOST string // WebSocket server URL
	Path         string // WebSocket server path
	Servername   string // Name of the server

}

// GetDefaultCommander returns a CommanderConfig with default values for local development.
// Deprecated: Moving to interface-based configs for better testability.
func GetDefaultCommander() *CommanderConfig {
	DefaultCommander := CommanderConfig{
		WsServerHOST: "localhost:8080",
		Path:         "/agent",
		Servername:   "sentinal-agent",
	}
	return &DefaultCommander
}

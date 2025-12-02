package agents

type CommanderConfig struct {
	//u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/agent-ws"}
	WsServerHOST string // WebSocket server URL
	Path         string // WebSocket server path
	Servername   string // Name of the server

}

func GetDefaultCommander() *CommanderConfig {
	DefaultCommander := CommanderConfig{
		WsServerHOST: "localhost:8080",
		Path:         "/agent-ws",
		Servername:   "sentinal-agent",
	}
	return &DefaultCommander
}

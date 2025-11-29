package lib


type CommandType string
const (
	// general container commands that do not need container id
	LIST_CONTAINERS = "list_container"
	CREATE_CONTAINER = "create_container"

	// container specific commands that need container id
	DELETE_CONTAINER = "delete_container"
	STOP_CONTAINER = "stop_container"
	START_CONTAINER = "start_container"
	RESTART_CONTAINER = "restart_container"

)



type StartConatianerParams struct {
	Image        string `json:"image" binding:"required"`
	ExposedPorts int    `json:"exposed_ports"`
	Protocol     string `json:"protocol"`
	HostIP       string `json:"host_ip"`
	HostPort     string `json:"host_port"`
	Name         string `json:"name" binding:"required"`
}
type ListContainsersParams struct {
	Size    bool `json:"size"`
    All     bool `json:"all"`
    Latest  bool `json:"latest"`
    Since   string `json:"since"`
    Before  string `json:"before"`
    Limit   int `json:"limit"`
}
type Params struct {
	StartParams StartConatianerParams `json:"start_params"`
}

//TODO: replace
type Command struct {
	ContainerID string
	MachineID string
	CMD CommandType
	Stream bool
	Params Params
}


// Factory method to create Command instances
// If ContainerID is empty, it indicates a general command
// that does not target a specific container
func GetCommand(ContainerID string ,MachineID string , command CommandType , stream bool, params Params) Command{
	if ContainerID == "" {
		return Command{
			ContainerID: "",
			MachineID: MachineID,
			CMD: command,
			Stream: stream,
			Params: params,
		}
	}

	return Command{
		ContainerID: ContainerID,
		MachineID: MachineID,
		CMD: command,
		Stream: stream,
		Params: params,
	}
}


package lib

type CommandType string

const (
	// general container commands that do not need container id
	LIST_CONTAINERS  = "list_container"
	CREATE_CONTAINER = "create_container"

	// container specific commands that need container id
	DELETE_CONTAINER    = "delete_container"
	STOP_CONTAINER      = "stop_container"
	START_NEW_CONTAINER = "start_new_container"
	RESTART_CONTAINER   = "restart_container"
	START_CONTAINER     = "start_container"
)

type StartNewContainerConfig struct {
	Image        string `json:"image" binding:"required"`
	ExposedPorts int    `json:"exposed_ports"`
	Protocol     string `json:"protocol"`
	HostIP       string `json:"host_ip"`
	HostPort     string `json:"host_port"`
	Name         string `json:"name" binding:"required"`
}

type ListContainersConfig struct {
	Size bool `json:"size"`
	All  bool `json:"all"`
}

type Params struct {
	StartNewContainerConfig *StartNewContainerConfig `json:"start,omitempty"`
	ListContainersConfig    *ListContainersConfig    `json:"list,omitempty"`
	ContainerID             string                  `json:"container_id"`
}

// TODO: replace
type Command struct {
	ContainerID string
	MachineID   string
	CMD         CommandType
	Stream      bool
	Params      Params
}

// Factory method to create Command instances
// If ContainerID is empty, it indicates a general command
// that does not target a specific container
func GetCommand(ContainerID string, MachineID string, command CommandType, stream bool, params Params) Command {
	if ContainerID == "" {
		return Command{
			ContainerID: "",
			MachineID:   MachineID,
			CMD:         command,
			Stream:      stream,
			Params:      params,
		}
	}

	return Command{
		ContainerID: ContainerID,
		MachineID:   MachineID,
		CMD:         command,
		Stream:      stream,
		Params:      params,
	}
}

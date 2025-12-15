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



// StartNewContainerConfig defines configuration
// required to create and start a new Docker container
type StartNewContainerConfig struct {

	// Image is the Docker image name
	// Example: "redis", "nginx:latest"
	Image string `json:"image" binding:"required"`

	// ContainerPort is the port the application
	// listens on INSIDE the container
	// Example: Redis -> 6379, Nginx -> 80
	ContainerPort int `json:"container_port"`

	// HostPort is the port exposed on the HOST machine
	// Example: localhost:30001 -> container:6379
	HostPort int `json:"host_port"`

	// Protocol defines transport protocol
	// Allowed values: "tcp", "udp"
	// Default: "tcp"
	Protocol string `json:"protocol"`

	// HostIP defines the host interface to bind
	// Default: "0.0.0.0"
	HostIP string `json:"host_ip"`

	// Name is the Docker container name
	Name string `json:"name" binding:"required"`
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

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


//TODO: replace
type Command struct {
	ContainerID string
	MachineID string
	command CommandType
}

// Factory method to create Command instances
// If ContainerID is empty, it indicates a general command
// that does not target a specific container
func GetCommand(ContainerID string ,MachineID string , command CommandType) Command{
	if ContainerID == "" {
		return Command{
			ContainerID: "",
			MachineID: MachineID,
			command: command,
		}
	}

	return Command{
		ContainerID,
		MachineID,
		command,
	}
}




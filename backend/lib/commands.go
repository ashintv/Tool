package lib

type CommandType string
const (
	LIST_CONTAINERS = "list_container"
	DELETE_CONTAINER = "delete_container"
)


//TODO: replace
type Command struct {
	ContainerID string
	MachineID string
	command CommandType
}

func GetCommand(ContainerID string ,MachineID string , command CommandType) Command{
	return Command{
		ContainerID,
		MachineID,
		command,
	}
}



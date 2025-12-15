package main

import (
	"aetrix/observer/internals/agent"
	"context"

	"github.com/docker/docker/client"
)

// a seprate light weight routine
// a worker is goroutine
// a worker will find save all running containers
// a new  container start then in next fetch worker will add these into list ,
// if a container not found or stoped from previos list
// record the issue
// start a rotine for restart (with retry logic)
// report the error , outcone through mail service



func main() {
	// Initialize and start the agent
	config := agents.GetDefaultCommander()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	agent := agents.NewCommander(cli , config)

	// Run the agent
	agent.Run(context.Background())
}

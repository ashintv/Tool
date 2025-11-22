package services

import (
	"testing"
	"github.com/docker/docker/client"
)


type Config  struct{
	cli  *client.Client
}

func NewTestConfig(cli *client.Client) *Config {
	return &Config{
		cli: cli,
	}
}

func PullImage(T *testing.T){
	
}

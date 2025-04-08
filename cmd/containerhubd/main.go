package main

import (
	"github.com/shashanksingh24/ContainerHub/pkg/rpc"
)

func main() {
	rpc.StartServer("/tmp/containerhub.sock")
}

package main

import (
	"containerhub/pkg/rpc"
)

func main() {
	rpc.StartServer("/var/run/containerhub.sock")
}

package main

import (
	"github.com/skycoin/skycoin/src/mesh2/nodemanager"
)

func main() {
	rpcInstance := nodemanager.NewRPC()
	rpcInstance.Serve()
}

package main

import (
	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	rpcInstance := nodemanager.NewRPC()
	rpcInstance.Serve()
}

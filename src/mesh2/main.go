package main

import (
	"github.com/skycoin/skycoin/src/mesh2/meshrpc"
)

func main() {
	rpcInstance := meshrpc.NewRPC()
	rpcInstance.Serve()
}

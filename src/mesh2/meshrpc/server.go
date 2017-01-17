package meshrpc

import (
	//	"encoding/gob"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/skycoin/skycoin/src/mesh2/nodemanager"
)

type RPC struct {
}

func NewRPC() *RPC {
	newRPC := &RPC{}
	return newRPC
}

func (self *RPC) Serve() {
	port := os.Getenv("MESH_RPC_PORT")
	nm := nodemanager.NewNodeManager()
	//registerTypes()
	receiver := new(RPCReceiver)
	receiver.NodeManager = nm
	err := rpc.Register(receiver)
	if err != nil {
		panic(err)
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	log.Println("Serving RPC on port", port)
	err = http.Serve(l, nil)
	if err != nil {
		panic(err)
	}
}

/*
func registerTypes() {
	gob.Register([]byte{})
}
*/

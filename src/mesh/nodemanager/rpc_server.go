package nodemanager

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type RPC struct {
}

const DEFAULT_PORT = "1234"

func NewRPC() *RPC {
	newRPC := &RPC{}
	return newRPC
}

func (self *RPC) Serve() {
	port := os.Getenv("MESH_RPC_PORT")
	if port == "" {
		log.Println("No MESH_RPC_PORT environmental variable is found, assigning default port value:", DEFAULT_PORT)
		port = DEFAULT_PORT
	}
	nm, _ := newNodeManager("demo.network", "127.0.0.1:5999")
	receiver := new(RPCReceiver)
	receiver.NodeManager = nm
	receiver.cmPort = 5000
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

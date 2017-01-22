package app

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/nodemanager"
)

type Server struct {
	app
	Handle func([]byte) []byte
}

func NewServer(meshnet *nodemanager.NodeManager, address cipher.PubKey, handle func([]byte) []byte) (*Server, error) {
	server := &Server{}
	server.Register(meshnet, address)
	server.Handle = handle
	err := meshnet.AssignConsumer(address, server)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (self *Server) Consume(_ uint32, request []byte, responseChannel chan<- []byte) {
	response := self.Handle(request) // user defined
	go func() { responseChannel <- response }()
}

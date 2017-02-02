package app

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Server struct {
	app
	Handle func([]byte) []byte
}

func NewServer(network messages.Network, address cipher.PubKey, handle func([]byte) []byte) (*Server, error) {
	server := &Server{}
	server.register(network, address)
	server.Handle = handle
	err := network.Register(address, server)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (self *Server) Consume(_ uint32, request []byte, responseChannel chan<- []byte) {
	response := self.Handle(request) // user defined
	go func() { responseChannel <- response }()
}

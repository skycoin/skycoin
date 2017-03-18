package app

import (
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Server struct {
	app
	Accepted int // remove, testing purposes
}

func NewServer(meshnet messages.Network, address cipher.PubKey, handle func([]byte) []byte) (*Server, error) {
	server := &Server{}
	server.lock = &sync.Mutex{}
	server.register(meshnet, address)
	server.lock = &sync.Mutex{}
	server.timeout = time.Duration(messages.GetConfig().AppTimeout)
	server.handle = handle

	conn, err := meshnet.NewConnection(address)
	if err != nil {
		return nil, err
	}
	server.connection = conn

	err = meshnet.Register(address, server)
	if err != nil {
		return nil, err
	}
	return server, nil
}

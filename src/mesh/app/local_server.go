package app

import (
	"sync"

	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
)

type Server struct {
	app
}

func BrandNewServer(host, meshnet string, handle func([]byte) []byte) (*Server, error) {

	server := newServer(handle)

	conn, err := node.ConnectToMeshnet(host, meshnet)
	if err != nil {
		return nil, err
	}
	server.register(conn)

	return server, nil
}

func NewServer(conn messages.Connection, handle func([]byte) []byte) *Server {

	server := newServer(handle)

	server.register(conn)

	return server
}

func newServer(handle func([]byte) []byte) *Server {
	server := &Server{}
	server.lock = &sync.Mutex{}
	server.timeout = APP_TIMEOUT
	server.handle = handle
	return server
}

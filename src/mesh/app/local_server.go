package app

import (
	"sync"

	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
)

type Server struct {
	app
}

func BrandNewServer(appId messages.AppId, host, meshnet string, handle func([]byte) []byte) (*Server, error) {

	server := newServer(appId, handle)

	node, err := node.CreateAndConnectNode(host, meshnet)
	if err != nil {
		return nil, err
	}

	err = server.RegisterAtNode(node)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func NewServer(appId messages.AppId, node messages.NodeInterface, handle func([]byte) []byte) (*Server, error) {

	server := newServer(appId, handle)

	err := server.RegisterAtNode(node)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (self *Server) RegisterAtNode(node messages.NodeInterface) error {
	err := node.RegisterApp(self)
	if err != nil {
		return err
	}
	self.node = node
	return nil
}

func (self *Server) Consume(appMsg *messages.AppMessage) {

	sequence := appMsg.Sequence
	go func() {
		responsePayload := self.handle(appMsg.Payload)
		response := &messages.AppMessage{
			sequence,
			responsePayload,
		}
		responseSerialized := messages.Serialize(messages.MsgAppMessage, response)
		self.send(responseSerialized)
	}()
}

func newServer(appId messages.AppId, handle func([]byte) []byte) *Server {
	server := &Server{}
	server.id = appId
	server.lock = &sync.Mutex{}
	server.timeout = APP_TIMEOUT
	server.handle = handle
	return server
}

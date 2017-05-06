package app

import (
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type app struct {
	ProxyAddress     string
	id               messages.AppId
	node             messages.NodeInterface
	handle           func([]byte) []byte
	timeout          time.Duration
	sequence         uint32
	connection       messages.Connection
	responseChannels map[uint32]chan messages.AppResponse
	lock             *sync.Mutex
}

var APP_TIMEOUT = 100000 * time.Duration(time.Millisecond)

func (self *app) Id() messages.AppId {
	return self.id
}

func (self *app) RegisterAtNode(node messages.NodeInterface) error {
	err := node.RegisterApp(self)
	if err != nil {
		return err
	}
	self.node = node
	return nil
}

func (self *app) Connect(appId messages.AppId, address cipher.PubKey) error {
	_, err := self.node.Dial(address, self.id, appId)
	return err
}

func (self *app) Consume(appMsg *messages.AppMessage) {

	sequence := appMsg.Sequence
	responseChannel, err := self.getResponseChannel(sequence)
	if err != nil {
		responseChannel <- messages.AppResponse{nil, err}
	} else {
		responseChannel <- messages.AppResponse{appMsg.Payload, nil}
	}
}

func (self *app) AssignConnection(conn messages.Connection) {
	self.connection = conn
}

func (self *app) Shutdown() {
	if self.node != nil {
		self.node.Shutdown()
	}
}

func (self *app) send(msg []byte) {

	conn := self.connection
	if conn != nil {
		conn.Send(msg)
	}
}

func (self *app) getResponseChannel(sequence uint32) (chan messages.AppResponse, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	ch, ok := self.responseChannels[sequence]
	if !ok {
		return nil, messages.ERR_NO_CLIENT_RESPONSE_CHANNEL
	}
	return ch, nil
}

func (self *app) setResponseChannel(responseChannel chan messages.AppResponse) uint32 {
	self.lock.Lock()
	defer self.lock.Unlock()

	sequence := self.sequence
	self.sequence++
	self.responseChannels[sequence] = responseChannel
	return sequence
}

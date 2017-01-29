package nodemanager

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Connection struct {
	id messages.ConnectionId
	//	consumer        *Application
	nm              *NodeManager
	status          uint8
	nodeAttached    cipher.PubKey
	routeId         messages.RouteId
	backRouteId     messages.RouteId
	incomingChannel chan []byte
	closingChannel  chan bool
	sequence        uint32
}

const (
	DISCONNECTED = iota
	CONNECTED
)

func (nm *NodeManager) NewConnectionWithRoutes(nodeAttached cipher.PubKey, routeId, backRouteId messages.RouteId) (messages.Connection, error) {
	conn, err := newConnection(nm, nodeAttached)
	if err != nil {
		return nil, err
	}
	conn.routeId = routeId
	conn.backRouteId = backRouteId
	conn.status = CONNECTED
	return conn, nil
}

func (nm *NodeManager) NewConnection(nodeAttached, nodeTo cipher.PubKey) (messages.Connection, error) {
	conn, err := newConnection(nm, nodeAttached)
	if err != nil {
		return nil, err
	}
	routeId, backRouteId, err := nm.findRoute(nodeAttached, nodeTo)
	if err != nil {
		return nil, err
	}
	conn.routeId = routeId
	conn.backRouteId = backRouteId
	conn.status = CONNECTED
	return conn, nil
}

func newConnection(nm *NodeManager, nodeAttached cipher.PubKey) (*Connection, error) {
	id := messages.RandConnectionId()
	_, err := nm.getNodeById(nodeAttached)
	if err != nil {
		return nil, err
	}
	conn := &Connection{
		id:           id,
		nm:           nm,
		status:       DISCONNECTED,
		nodeAttached: nodeAttached,
	}
	conn.incomingChannel = make(chan []byte, 1024)
	conn.closingChannel = make(chan bool, 1024)
	return conn, nil
}

func (self *Connection) Send(msg []byte) (uint32, error) {
	if self.status != CONNECTED {
		return 0, errors.ERR_DISCONNECTED
	}
	requestMessage := messages.RequestMessage{
		BackRoute: self.backRouteId,
		Sequence:  self.sequence,
		Payload:   msg,
	}
	requestSerialized := messages.Serialize(messages.MsgRequestMessage, requestMessage)
	inRouteMessage := messages.InRouteMessage{
		messages.NIL_TRANSPORT,
		self.routeId,
		requestSerialized,
	}
	msgSerialized := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)
	node, err := self.nm.getNodeById(self.nodeAttached)
	if err != nil {
		return 0, err
	}
	node.InjectTransportMessage(msgSerialized)
	self.sequence++
	return self.sequence - 1, nil
}

func (self *Connection) receivingLoop() error {
	for self.status == CONNECTED {
		select {
		case data := <-self.incomingChannel: // accept from meshnet(node)
			fmt.Println("Data received", string(data)) // pass to server/client
		case <-self.closingChannel:
			self.Close()
			break
		}
	}
	return errors.ERR_DISCONNECTED
}

func (self *Connection) consume(msg []byte) error {
	self.incomingChannel <- msg
	return nil
}

func (self *Connection) GetStatus() uint8 {
	return self.status
}

func (self *Connection) Close() {
	close(self.incomingChannel)
	close(self.closingChannel)
	self.status = DISCONNECTED
}

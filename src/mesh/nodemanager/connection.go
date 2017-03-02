package nodemanager

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Connection struct {
	id             messages.ConnectionId
	nm             *NodeManager
	status         uint8
	nodeAttached   cipher.PubKey
	routeId        messages.RouteId
	backRouteId    messages.RouteId
	closingChannel chan bool
	sequence       uint32
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
	conn.closingChannel = make(chan bool)
	return conn, nil
}

func (self *Connection) Send(msg []byte) (uint32, error) {

	//	messages.RegisterEvent("Conn.Send start")

	if self.status != CONNECTED {

		//		messages.RegisterEvent("Conn.Send return (DISCONNECT)")

		return 0, errors.ERR_DISCONNECTED
	}
	requestMessage := messages.RequestMessage{
		BackRoute: self.backRouteId,
		Sequence:  self.sequence,
		Payload:   msg,
	}

	//	messages.RegisterEvent("conn.Send serializing")

	requestSerialized := messages.Serialize(messages.MsgRequestMessage, requestMessage)

	//	messages.RegisterEvent("conn.Send serialized")

	inRouteMessage := messages.InRouteMessage{
		messages.NIL_TRANSPORT,
		self.routeId,
		requestSerialized,
	}
	//	msgSerialized := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)
	node, err := self.nm.getNodeById(self.nodeAttached)
	if err != nil {

		//		messages.RegisterEvent("Conn.Send return (No Node By ID)")

		return 0, err
	}
	//	messages.RegisterEvent("Conn.Send passes to InjectTransportMesage")
	node.InjectTransportMessage(&inRouteMessage)
	self.sequence++

	//	messages.RegisterEvent("Conn.Send return (SUCCESS)")

	return self.sequence - 1, nil
}

func (self *Connection) GetStatus() uint8 {
	return self.status
}

func (self *Connection) Close() {
	//close(self.closingChannel)
	self.status = DISCONNECTED
}

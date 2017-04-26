package messages

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type NodeInterface interface {
	Dial(cipher.PubKey, AppId, AppId) (Connection, error)
	Id() cipher.PubKey
	InjectTransportMessage(*InRouteMessage)
	InjectCongestionPacket(*CongestionPacket)
	GetTransportToNode(cipher.PubKey) (TransportInterface, error)
	GetConnection(ConnectionId) Connection
	ConnectedTo(cipher.PubKey) bool
	RegisterApp(Consumer) error
	Shutdown()
}

type TransportInterface interface {
	Id() TransportId
	PacketsSent() uint32
	PacketsConfirmed() uint32
}

type Consumer interface {
	Id() AppId
	Consume(*AppMessage)
	AssignConnection(Connection)
}

type Network interface {
	Shutdown()
}

type Connection interface {
	Send([]byte) error
	Status() uint8
	Id() ConnectionId
}

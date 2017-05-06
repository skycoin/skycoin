package messages

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type NodeInterface interface {
	Id() cipher.PubKey
	ConnectDirectly(cipher.PubKey) error
	Dial(cipher.PubKey, AppId, AppId) (Connection, error)
	InjectTransportMessage(*InRouteMessage)
	InjectCongestionPacket(*CongestionPacket)
	GetTransportToNode(cipher.PubKey) (TransportInterface, error)
	GetConnection(ConnectionId) Connection
	ConnectedTo(cipher.PubKey) bool
	RegisterApp(Consumer) error
	Shutdown()
	TalkToViscript(uint32, uint32)
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
	TalkToViscript(uint32, uint32)
	Shutdown()
}

type Connection interface {
	Send([]byte) error
	Status() uint8
	Id() ConnectionId
}

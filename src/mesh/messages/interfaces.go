package messages

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type NodeInterface interface {
	GetId() cipher.PubKey
	InjectTransportMessage(*InRouteMessage)
	InjectCongestionPacket(*CongestionPacket)
	GetTransportToNode(cipher.PubKey) (TransportInterface, error)
	GetConnection() Connection
	ConnectedTo(cipher.PubKey) bool
	Shutdown()
}

type TransportInterface interface {
	GetId() TransportId
	GetPacketsSent() uint32
	GetPacketsConfirmed() uint32
}

type Consumer interface {
	Consume([]byte)
}

type Network interface {
	Shutdown()
}

type Connection interface {
	Address() cipher.PubKey
	Send([]byte) error
	Dial(cipher.PubKey) error
	AssignConsumer(Consumer)
	Shutdown()
	GetStatus() uint8
}

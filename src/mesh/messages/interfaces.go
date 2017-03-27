package messages

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type NodeInterface interface {
	GetId() cipher.PubKey
	GetPeer() *Peer
	InjectTransportMessage(*InRouteMessage)
	InjectConnectionMessage(*InRouteMessage)
	InjectCongestionPacket(*CongestionPacket)
	SetTransport(TransportId, TransportInterface)
	ConnectedTo(NodeInterface) bool
}

type TransportInterface interface {
	InjectNodeMessage(*InRouteMessage)
}

type Consumer interface {
	Consume([]byte)
}

type User interface {
	Use([]byte)
}

type Network interface {
	NewConnection(cipher.PubKey) (Connection, error)
	Connect(cipher.PubKey, cipher.PubKey) error
	Register(cipher.PubKey, Consumer) error
}

type Connection interface {
	Send([]byte)
	Use([]byte)
}

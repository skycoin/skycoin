package messages

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type NodeInterface interface {
	GetId() cipher.PubKey
	InjectTransportMessage(msg []byte)
	SetTransport(TransportId, TransportInterface)
	ConnectedTo(NodeInterface) bool
}

type TransportInterface interface {
	InjectNodeMessage([]byte)
}

type Consumer interface {
	Consume(uint32, []byte, chan<- []byte) // number of message, what to consume and channel for accepting responses
}

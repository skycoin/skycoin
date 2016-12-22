package messages

import ()

type NodeInterface interface {
	InjectTransportMessage([]byte)
}

type TransportInterface interface {
	InjectNodeMessage([]byte)
}

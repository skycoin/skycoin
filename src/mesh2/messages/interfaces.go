package messages

import ()

type NodeInterface interface {
	InjectTransportMessage(transportId TransportId, msg []byte)
}

type TransportInterface interface {
	InjectNodeMessage([]byte)
}

//later add "transport status" struct

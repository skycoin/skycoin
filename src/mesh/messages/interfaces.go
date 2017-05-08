package messages

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type NodeInterface interface {
	Id() cipher.PubKey
	ConnectDirectly(cipher.PubKey) error
	AppTalkAddr() string
	Shutdown()
	TalkToViscript(uint32, uint32)
}

type NodeInTransport interface {
	InjectTransportMessage(*InRouteMessage)
	InjectCongestionPacket(*CongestionPacket)
}

type TransportInterface interface {
	Id() TransportId
	PacketsSent() uint32
	PacketsConfirmed() uint32
}

type Network interface {
	Addr() string
	TalkToViscript(uint32, uint32)
	Shutdown()
}

type Connection interface {
	Id() ConnectionId
	Send([]byte) error
	Status() uint8
}

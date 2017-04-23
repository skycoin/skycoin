package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
)

const (
	MsgInRouteMessage             = iota // Transport -> Node
	MsgOutRouteMessage                   // Node -> Transport
	MsgTransportDatagramTransfer         // Transport -> Transport, simulating sending packet over network
	MsgTransportDatagramACK              // Transport -> Transport, simulating ACK for packet
	MsgConnectionMessage                 // Connection -> Connection
	MsgConnectionAck                     // Connection -> Connection
	MsgProxyMessage                      // Application -> Application
	MsgAppMessage                        // Application -> Application
	MsgCongestionPacket                  // Transport -> Transport
	MsgInControlMessage                  // Transport -> Node, control message
	MsgOutControlMessage                 // Node -> Transport, control message
	MsgCloseChannelControlMessage        // Node -> Control channel, close control channel
	MsgAddRouteCM                        // Node -> Control channel, add new route
	MsgRemoveRouteCM                     // Node -> Control channel, remove route
	MsgRegisterNodeCM                    // Node -> NodeManager
	MsgRegisterNodeCMAck                 // NodeManager -> Node
	MsgAssignPortCM                      // NodeManager -> Node
	MsgTransportCreateCM                 // NodeManager -> Node
	MsgTransportTickCM                   // NodeManager -> Node
	MsgTransportShutdownCM               // NodeManager -> Node
	MsgOpenUDPCM                         // NodeManager -> Node
	MsgCommonCMAck                       // Node -> NodeManager, NodeManager -> Node
	MsgConnectCM                         // Node -> NodeManager
	MsgAssignRouteCM                     // NodeManager -> Node
	MsgConnectionOnCM                    // NodeManager -> Node
	MsgShutdownCM                        // NodeManager -> Node
	//MessageMouseScroll        // 1
	//MessageMouseButton        // 2
	//MessageCharacter
	//MessageKey
)

func GetMessageType(message []byte) uint16 {
	var value uint16
	rBuf := bytes.NewReader(message[0:2])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	}
	//	value = (uint16)(message[0])
	return value
}

//Node Messages

// Transport -> Node Messages

//message received by node, from transport
//message comes in by a channel
type InRouteMessage struct {
	TransportId TransportId //who sent it
	RouteId     RouteId     //the incoming route
	Datagram    []byte      //length prefixed message
}

// Node -> Transport Messages

//message node, writes to the channel of the transport
type OutRouteMessage struct {
	RouteId          RouteId //the incoming route
	Datagram         []byte  //length prefixed message
	ResponseRequired bool
}

type CongestionPacket struct {
	TransportId TransportId //who sent it
	Congestion  bool        // true - increase throttle, false - decrease
}

// Transport -> Transport

//simulates one end of a transport, sending data to other end of the pair
type TransportDatagramTransfer struct {
	//put seq number for confirmation/ACK
	RouteId          RouteId
	Sequence         uint32 //sequential sequence number of ACK
	Datagram         []byte
	ResponseRequired bool
}

type TransportDatagramACK struct {
	LowestSequence uint32 //ACK anything below this SEQ number
	Bitarray       uint32 //ACK packets at LowestSequence + Bit offset, if equal to 1
}

type InControlMessage struct {
	ChannelId      ChannelId
	Sequence       uint32
	PayloadMessage []byte
}

type AddRouteCM struct {
	IncomingTransportId TransportId
	OutgoingTransportId TransportId
	IncomingRouteId     RouteId
	OutgoingRouteId     RouteId
}

type RemoveRouteCM struct {
	RouteId RouteId
}

type ConnectionMessage struct {
	Sequence uint32
	Order    uint32
	Total    uint32
	Payload  []byte
}

type ConnectionAck struct {
	Sequence uint32
}

type AppMessage struct {
	Sequence         uint32
	ResponseRequired bool
	Payload          []byte
}

type AppResponse struct {
	Response []byte
	Err      error
}

type ProxyMessage struct {
	Data       []byte
	RemoteAddr string
	NeedClose  bool
}

// ==================== control messages ========================

type RegisterNodeCM struct {
	Host    string
	Connect bool
}

type RegisterNodeCMAck struct {
	Ok                bool
	NodeId            cipher.PubKey
	MaxBuffer         uint64
	MaxPacketSize     uint32
	TimeUnit          uint32
	SendInterval      uint32
	ConnectionTimeout uint32
}

type AssignPortCM struct {
	Port uint32
}

type TransportCreateCM struct {
	Id                TransportId
	PairId            TransportId
	PairedNodeId      cipher.PubKey
	MaxBuffer         uint64
	TimeUnit          uint32
	TransportTimeout  uint32
	SimulateDelay     bool
	MaxSimulatedDelay uint32
	RetransmitLimit   uint32
}

type TransportTickCM struct {
	Id TransportId
}

type TransportShutdownCM struct {
	Id TransportId
}

type OpenUDPCM struct {
	Id    TransportId
	PeerA Peer
	PeerB Peer
}

type CommonCMAck struct {
	Ok bool
}

type ConnectCM struct {
	From cipher.PubKey
	To   cipher.PubKey
}

type AssignRouteCM struct {
	RouteId RouteId
}

type ConnectionOnCM struct {
	NodeId cipher.PubKey
}

type ShutdownCM struct {
	NodeId cipher.PubKey
}

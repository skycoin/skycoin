package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
)

const (
	MsgInRouteMessage            = iota // Transport -> Node
	MsgOutRouteMessage                  // Node -> Transport
	MsgTransportDatagramTransfer        // Transport -> Transport, simulating sending packet over network
	MsgTransportDatagramACK             // Transport -> Transport, simulating ACK for packet
	MsgCongestionPacket                 // Transport -> Transport

	MsgConnectionMessage // Connection -> Connection
	MsgConnectionAck     // Connection -> Connection

	MsgProxyMessage // Application -> Application
	MsgAppMessage   // Application -> Application

	MsgInControlMessage           // Transport -> Node, control message
	MsgOutControlMessage          // Node -> Transport, control message
	MsgCloseChannelControlMessage // Node -> Control channel, close control channel
	MsgAddRouteCM                 // Node -> Control channel, add new route
	MsgRemoveRouteCM              // Node -> Control channel, remove route
	MsgRegisterNodeCM             // Node -> NodeManager
	MsgRegisterNodeCMAck          // NodeManager -> Node
	MsgAssignPortCM               // NodeManager -> Node
	MsgTransportCreateCM          // NodeManager -> Node
	MsgTransportTickCM            // NodeManager -> Node
	MsgTransportShutdownCM        // NodeManager -> Node
	MsgOpenUDPCM                  // NodeManager -> Node
	MsgCommonCMAck                // Node -> NodeManager, NodeManager -> Node
	MsgConnectDirectlyCM          // Node -> NodeManager
	MsgConnectDirectlyCMAck       // NodeManager -> Node
	MsgConnectWithRouteCM         // Node -> NodeManager
	MsgConnectWithRouteCMAck      // NodeManager -> Node
	MsgAssignConnectionCM         // NodeManager -> Node
	MsgConnectionOnCM             // NodeManager -> Node
	MsgShutdownCM                 // NodeManager -> Node

	MsgNodeAppMessage      // Application -> Node
	MsgNodeAppResponse     // Node -> Application
	MsgSendFromAppMessage  // Application -> Node
	MsgRegisterAppMessage  // Application -> Node
	MsgConnectToAppMessage // Application -> Node
	MsgAssignConnectionNAM // Node -> Application

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
	RouteId  RouteId //the incoming route
	Datagram []byte  //length prefixed message
}

type CongestionPacket struct {
	TransportId TransportId //who sent it
	Congestion  bool        // true - increase throttle, false - decrease
}

// Transport -> Transport

//simulates one end of a transport, sending data to other end of the pair
type TransportDatagramTransfer struct {
	//put seq number for confirmation/ACK
	RouteId  RouteId
	Sequence uint32 //sequential sequence number of ACK
	Datagram []byte
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
	Sequence     uint32
	ConnectionId ConnectionId
	Order        uint32
	Total        uint32
	Payload      []byte
}

type ConnectionAck struct {
	Sequence     uint32
	ConnectionId ConnectionId
}

type AppMessage struct {
	Sequence uint32
	Payload  []byte
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
	Hostname string
	Host     string
	Connect  bool
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

type ConnectDirectlyCM struct {
	Sequence uint32
	From     cipher.PubKey
	To       string
}

type ConnectDirectlyCMAck struct {
	Sequence uint32
	Ok       bool
}

type ConnectWithRouteCM struct {
	Sequence  uint32
	AppIdFrom AppId
	AppIdTo   AppId
	From      cipher.PubKey
	To        string
}

type ConnectWithRouteCMAck struct {
	Sequence     uint32
	Ok           bool
	ConnectionId ConnectionId
}

type AssignConnectionCM struct {
	ConnectionId ConnectionId
	RouteId      RouteId
	AppId        AppId
}

type ConnectionOnCM struct {
	NodeId       cipher.PubKey
	ConnectionId ConnectionId
}

type ShutdownCM struct {
	NodeId cipher.PubKey
}

type UserCommand struct {
	Sequence uint32
	AppId    uint32
	Payload  []byte
}

type NodeAppMessage struct {
	Sequence uint32
	AppId    AppId
	Payload  []byte
}

type NodeAppResponse struct {
	Sequence uint32
}

type SendFromAppMessage struct {
	ConnectionId ConnectionId
	Payload      []byte
}

type RegisterAppMessage struct {
}

type AssignConnectionNAM struct {
	ConnectionId ConnectionId
}

type ConnectToAppMessage struct {
	Address string
	AppFrom AppId
	AppTo   AppId
}

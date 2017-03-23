package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	MsgInRouteMessage             = iota // Transport -> Node
	MsgOutRouteMessage                   // Node -> Transport
	MsgTransportDatagramTransfer         // Transport -> Transport, simulating sending packet over network
	MsgTransportDatagramACK              // Transport -> Transport, simulating ACK for packet
	MsgInControlMessage                  // Transport -> Node, control message
	MsgOutControlMessage                 // Node -> Transport, control message
	MsgCloseChannelControlMessage        // Node -> Control channel, close control channel
	MsgAddRouteControlMessage            // Node -> Control channel, add new route
	MsgRemoveRouteControlMessage         // Node -> Control channel, remove route
	MsgConnectionMessage                 // Connection -> Connection
	MsgConnectionAck                     // Connection -> Connection
	MsgAppMessage                        // Application -> Application
	MsgCongestionPacket                  // Transport -> Transport
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
	ChannelId       ChannelId
	PayloadMessage  []byte
	ResponseChannel chan bool
}

//type CreateChannelControlMessage struct {
//}

type AddRouteControlMessage struct {
	IncomingTransportId TransportId
	OutgoingTransportId TransportId
	IncomingRouteId     RouteId
	OutgoingRouteId     RouteId
}

type RemoveRouteControlMessage struct {
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

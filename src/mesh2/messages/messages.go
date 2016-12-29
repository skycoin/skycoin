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
	MsgTransportDatagramTransfer        //Transport -> Transport, simulating sending packet over network
	MsgTransportDatagramACK             //Transport -> Transport, simulating ACK for packet
//MessageMouseScroll        // 1
//MessageMouseButton        // 2
//MessageCharacter
//MessageKey
)

func GetMessageType(message []byte) uint8 {
	var value uint8
	rBuf := bytes.NewReader(message[4:5])
	err := binary.Read(rBuf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("binary.Read failed: ", err)
	} else {
		//fmt.Printf("from byte buffer, %s: %d\n", s, value)
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

// Transport -> Transport

//simulates one end of a transport, sending data to other end of the pair
type TransportDatagramTransfer struct {
	//put seq number for confirmation/ACK
	Sequence uint32 //sequential sequence number of ACK
	Datagram []byte
}

type TransportDatagramACK struct {
	LowestSequence uint32 //ACK anything below this SEQ number
	Bitarray       uint32 //ACK packets at LowestSequence + Bit offset, if equal to 1
}

type CreateChannelControlMessage struct {
}

type AddRouteControlMessage struct {
	NodeId  cipher.PubKey
	RouteId RouteId
}

type ExtendRouteControlMessage struct {
	NodeId  cipher.PubKey
	RouteId RouteId
}

type RemoveRouteControlMessage struct {
	RouteId RouteId
}

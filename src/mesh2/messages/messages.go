package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	//"math/rand"
)

const (
	MsgInRouteMessage  = iota // Transport -> Node
	MsgOutRouteMessage        // Node -> Transport
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

package transport

import (
	"fmt"
	"github.com/skycoin/skycoin/src/mesh2/messages"
)

//This is stub transport
type Transport struct {
	Id              messages.TransportId
	IncomingChannel chan ([]byte)

	AttachedNode *messages.NodeInterface //node the transport is attached to
}

//are created by the factories
func (self *Transport) NewTransportStub() {
	fmt.Printf("Created Transport:")
	self.IncomingChannel = make(chan []byte, 1024)
	self.Id = messages.RandTransportId()
}

func (self *Transport) Shutdown() {
	close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *Transport) Tick() {
	//process incoming messages
	for msg := range self.IncomingChannel {
		//process our incoming messages
		//fmt.Println(msg)

		switch messages.GetMessageType(msg) {

		//InRouteMessage is the only message coming in to node from transports
		case messages.MsgOutRouteMessage:

			var m1 messages.OutRouteMessage
			messages.Deserialize(msg, m1)
			//get message and put into the queue to be sent out
		}

	}
}

//inject an incoming message from the transport
func (self *Transport) InjectNodeMessage(msg []byte) {
	self.AttachedNode.InjectTransportMessage(self.transportId, msg)
}

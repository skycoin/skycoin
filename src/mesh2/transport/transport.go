package transport

import (
	"fmt"
	"github.com/skycoin/skycoin/src/mesh2/messages"
)

//Stub Transport

//TODO:
// - implement simulated "delay" for transport
// - implement simulated out of order packet delivery
// - implement simulated packet drop
// - implement real UDP trasport
// - implement status "connected/disconnected"
// - implement ACKs

// TODO:
// pendingOut channel may need to be on the transport factory itself
// more efficient to push to central location than to pull/poll the transport list
// TODO:
// - may be more efficient to replace pending out, with an array on TransportFactory
// - or with array on the transport (who is responsible for processing ACKs?)

//This is stub transport
type Transport struct {
	Id              messages.TransportId
	IncomingChannel chan ([]byte)
	PendingOut      chan ([]byte) //messages to send to other end of transport
	//Note: PendingOut channel may need to be on the transport_factory

	AttachedNode messages.NodeInterface //node the transport is attached to

	StubPair *Transport //this is the other transport stub pair
}

//are created by the factories
func (self *Transport) NewTransportStub() {
	fmt.Printf("Created Transport:\n")
	self.IncomingChannel = make(chan []byte, 1024)
	//self.PendingOut = make(chan []byte, 1024)
	self.Id = messages.RandTransportId()
}

func (self *Transport) Shutdown() {
	close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *Transport) Tick() {
	//process incoming messages
	fmt.Println("tick")
	for msg := range self.IncomingChannel {
		//process our incoming messages
		fmt.Printf("\ntransport with id %d gets message %d\n\n", self.Id, msg)

		switch messages.GetMessageType(msg) {

		//OutRouteMessage is the only message coming in to node from transports
		//Node -> Transport message
		case messages.MsgOutRouteMessage:

			var m1 messages.OutRouteMessage
			messages.Deserialize(msg, &m1)
			//get message and put into the queue to be sent out
			//prime message for transit between the two transport ends
			var m1b messages.TransportDatagramTransfer
			m1b.Datagram = m1.Datagram
			b1 := messages.Serialize(messages.MsgTransportDatagramTransfer, &m1b)
			self.PendingOut <- b1 //push to queue, to be transferred

		//Transport -> Transport messag
		case messages.MsgTransportDatagramTransfer:
			var m2 messages.OutRouteMessage
			messages.Deserialize(msg, m2)

		}
	}
}

//inject an incoming message from the transport
func (self *Transport) InjectNodeMessage(msg []byte) {
	self.AttachedNode.InjectTransportMessage(self.Id, msg)
}

//message from stub to stub
//used internally by transport factory
func (self *Transport) SendMessageToStubPair(msg []byte) {
	self.StubPair.IncomingChannel <- msg
}

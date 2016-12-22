package node

import (
	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/transport"

	"fmt"
)

//A Node has a map of route rewriting rules
//A Node has a control channel for setting and modifying the route rewrite rules
//A Node has a list of transports

//Route rewriting rules
//-nodes receive messages on a route
//-nodes look up the route in a table and if it has a rewrite rule, rewrites the route
// and forwards it to the transport

type Node struct {
	IncomingChannel chan ([]byte)

	Transports           map[messages.TransportId]*transport.Transport
	RouteForwardingRules map[messages.RouteId]RouteRule
}

type RouteRule struct {
	IncomingTransport messages.TransportId
	OutgoingTransport messages.TransportId
	IncomingRoute     messages.RouteId
	OutgoingRoute     messages.RouteId
}

func (self *Node) New() {
	self.IncomingChannel = make(chan []byte, 1024)
	self.RouteForwardingRules = make(map[messages.RouteId]RouteRule)
}

func (self *Node) Shutdown() {
	close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *Node) Tick() {
	//process incoming messages
	for msg := range self.IncomingChannel {
		//process our incoming messages
		//fmt.Println(msg)

		switch messsages.GetMessageType(msg) {

		//InRouteMessage is the only message coming in to node from transports
		case messages.InRouteMessage:

			var m1 messages.InRouteMessage
			messages.Deserialize(msg, m1)

			//look in route table
			routeId := m1.RouteId
			transportId := m1.TransportId //who is is from

			//check that transport exists
			if _, ok := self.Transports[transportId]; !ok {
				log.Printf("Received message From Transport that does not exist")
			}

			//check if route exists
			if _, ok := self.Transports[transportId]; !ok {
				log.Printf("Received message From Transport that does not exist")
			}

			//case MessageMouseScroll:
			//s("MessageMouseScroll", message)
			//showFloat64("X Offset", message)
			//showFloat64("Y Offset", message)

			//case MessageMouseButton:

		}

	}
}

//inject an incoming message from the transport
func (self *Node) InjectTransportMessage([]byte) {

}

package node

import (
	errmsg "github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/transport"

	"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"sync"
)

//A Node has a map of route rewriting rules
//A Node has a control channel for setting and modifying the route rewrite rules
//A Node has a list of transports

//Route rewriting rules
//-nodes receive messages on a route
//-nodes look up the route in a table and if it has a rewrite rule, rewrites the route
// and forwards it to the transport

type Node struct {
	Id                     cipher.PubKey
	IncomingChannel        chan *messages.InRouteMessage
	incomingControlChannel chan messages.InControlMessage

	Transports map[messages.TransportId]*transport.Transport
	//	transportChannel chan *TransportRequest

	RouteForwardingRules map[messages.RouteId]*RouteRule
	//routeChannel         chan *RouteRequest

	lock *sync.Mutex

	consumer messages.Consumer

	controlChannels map[messages.ChannelId]*ControlChannel

	Host string
	Port uint32
}

type RouteRule struct {
	IncomingTransport messages.TransportId
	OutgoingTransport messages.TransportId
	IncomingRoute     messages.RouteId
	OutgoingRoute     messages.RouteId
}

func NewNode() *Node {
	node := new(Node)
	node.Id = createPubKey()
	node.IncomingChannel = make(chan *messages.InRouteMessage, 1024)
	node.incomingControlChannel = make(chan messages.InControlMessage, 256)
	node.Transports = make(map[messages.TransportId]*transport.Transport)
	node.RouteForwardingRules = make(map[messages.RouteId]*RouteRule)
	node.controlChannels = make(map[messages.ChannelId]*ControlChannel)
	//node.routeChannel = make(chan *RouteRequest)
	//	node.transportChannel = make(chan *TransportRequest)
	node.lock = &sync.Mutex{}
	if messages.IsDebug() {
		fmt.Printf("Created Node %s\n", node.Id.Hex())
	}
	node.Tick()
	return node
}

func createPubKey() cipher.PubKey {
	pub, _ := cipher.GenerateKeyPair()
	return pub
}

func (self *Node) Shutdown() {
	//close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *Node) Tick() {
	//go self.serveTransports()
	//go self.serveRoutes()
	backChannel := make(chan bool, 1)
	self.runCycles(backChannel)
	<-backChannel
}

//move node forward on tick, process events
func (self *Node) runCycles(backChannel chan bool) {
	//process incoming messages
	go self.handleIncomingTransportMessages() //pop them off the channel
	go self.handleIncomingControlMessages()   //pop them off the channel
	backChannel <- true
}

func (self *Node) handleIncomingTransportMessages() {
	for msg := range self.IncomingChannel {
		//		messages.RegisterEvent("Node accepted message through incoming channel")
		if messages.IsDebug() {
			fmt.Printf("\nnode with id %s accepting a message\n\n", self.Id.Hex())
		}
		//process our incoming messages
		//InRouteMessage is the only message coming in to node from transports
		if messages.IsDebug() {
			fmt.Println("InRouteMessage", msg)
		}
		self.handleInRouteMessage(msg)
	}
}

func (self *Node) handleIncomingControlMessages() {
	for msg := range self.incomingControlChannel {
		self.handleControlMessage(msg.ChannelId, msg.PayloadMessage)
		msg.ResponseChannel <- true
	}
}

func (self *Node) handleInRouteMessage(m1 *messages.InRouteMessage) {
	//look in route table

	//	messages.RegisterEvent("handleInRouteMessage start")

	routeId := m1.RouteId
	transportId := m1.TransportId //who is is from
	//check that transport exists
	if _, err := self.getTransport(transportId); err != nil && transportId != messages.NIL_TRANSPORT {
		log.Printf("Node %s received message From Transport that does not exist\n", self.Id.Hex())
	}
	//check if route exists
	routeRule, err := self.getRoute(routeId)
	if err != nil {

		//		messages.RegisterEvent("handleInRouteMessage error (route doesn't exist)")

		log.Fatalf("Node %s received route message for route that does not exist\n", self.Id.Hex())
	}
	//check that incoming transport is correct
	if transportId != routeRule.IncomingTransport {

		//		messages.RegisterEvent("handleInRouteMessage error (route doesn't match transport)")

		log.Panic("Node: incoming route does not match the transport id it should be received from")
	}
	if routeId != routeRule.IncomingRoute {

		//		messages.RegisterEvent("handleInRouteMessage error (route doesn't match route rule)")

		log.Panic("Node: impossible, incoming route id does not match route rule id")
	}
	//construct new packet
	outgoingRouteId := routeRule.OutgoingRoute
	outgoingTransportId := routeRule.OutgoingTransport //find transport to resend datagram
	datagram := m1.Datagram

	if outgoingRouteId == messages.NIL_ROUTE && outgoingTransportId == messages.NIL_TRANSPORT {

		//		messages.RegisterEvent("handleInRouteMessage passes to consume")

		self.consume(datagram)

		//		messages.RegisterEvent("handleInRouteMessage ends consume")

	} else {

		//		messages.RegisterEvent("handleInRouteMessage starts creating OutRouteMessage")

		var out messages.OutRouteMessage //replace inRoute, with outRoute, using rule
		out.RouteId = outgoingRouteId
		out.Datagram = datagram
		//serialize message, with prefix
		//		b1 := messages.Serialize(messages.MsgOutRouteMessage, out)

		if outgoingTransport, err := self.getTransport(outgoingTransportId); err == nil {

			//			messages.RegisterEvent("handleInRouteMessage passes OutRouteMessage to node")

			outgoingTransport.GetFromNode(out) //inject message to transport

			//			messages.RegisterEvent("handleInRouteMessage passed OutRouteMessage to node")

		}
	}
}

//inject an incoming message from the transport
func (self *Node) InjectTransportMessage(msg *messages.InRouteMessage) {
	//	messages.RegisterEvent("InjectTransportMessage")
	self.IncomingChannel <- msg //push message to channel
}

func (self *Node) InjectControlMessage(msg messages.InControlMessage) {
	msg.ResponseChannel = make(chan bool)
	self.incomingControlChannel <- msg

	<-msg.ResponseChannel
}

func (self *Node) GetId() cipher.PubKey {
	return self.Id
}

func (self *Node) GetPeer() *messages.Peer {
	peer := &messages.Peer{self.Host, self.Port}
	return peer
}

func (self *Node) GetTransportById(id messages.TransportId) (*transport.Transport, error) {
	tr, err := self.getTransport(id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Node %d doesn't contain transport with id %d", self.Id, id))
	}
	return tr, nil
}

func (self *Node) GetTransportToNode(nodeId cipher.PubKey) (*transport.Transport, error) {
	for _, transport := range self.Transports {
		if nodeId == transport.StubPair.AttachedNode.GetId() {
			return transport, nil
		}
	}
	return nil, errmsg.ERR_NO_TRANSPORT_TO_NODE
}

func (self *Node) ConnectedTo(other messages.NodeInterface) bool {
	_, err := self.GetTransportToNode(other.GetId())
	return err == nil
}

func (self *Node) SetTransport(id messages.TransportId, tr messages.TransportInterface) {
	self.setTransport(id, tr.(*transport.Transport))
}

func (self *Node) consume(msg []byte) {
	if messages.IsDebug() {
		fmt.Printf("\nNode %s consumed message %d\n\n", self.Id.Hex(), msg)
	}

	//	messages.RegisterEvent("starting consume")

	go func() {
		switch messages.GetMessageType(msg) {

		case messages.MsgRequestMessage:

			//			messages.RegisterEvent("recognized consume request: deserializing")

			requestMessage := &messages.RequestMessage{}
			err := messages.Deserialize(msg, requestMessage)
			if err != nil {
				// send wrong request format
			} else {

				//				messages.RegisterEvent("passing to consumeRequest")

				go self.consumeRequest(requestMessage)
			}

		case messages.MsgResponseMessage:

			//			messages.RegisterEvent("recognized consume response: deserializing")

			responseMessage := &messages.ResponseMessage{}
			err := messages.Deserialize(msg, responseMessage)
			if err != nil {
				// send wrong request format
			} else {

				//				messages.RegisterEvent("passing to consumeResponse")

				go self.consumeResponse(responseMessage)
			}

		default:
			fmt.Println("Incorrect message type:", messages.GetMessageType(msg))
		}
	}()
}

func (self *Node) consumeRequest(requestMessage *messages.RequestMessage) {

	//	messages.RegisterEvent("consumeRequest start")

	if self.consumer == nil {
		fmt.Printf("\nNo consumer registered at node %s\n\n", self.Id.Hex())
		return
	}
	backRoute := requestMessage.BackRoute
	sequence := requestMessage.Sequence
	payload := requestMessage.Payload
	responseChannel := make(chan []byte)

	//	messages.RegisterEvent("consumeRequest passing to consumer.Consume")

	go self.consumer.Consume(sequence, payload, responseChannel)

	//	messages.RegisterEvent("consumeRequest passed to consumer.Consume")

	responsePayload := <-responseChannel

	//	messages.RegisterEvent("consumeRequest got response from consumer.Consume: sending response")

	go self.sendResponse(backRoute, sequence, responsePayload)

	//	messages.RegisterEvent("consumeRequest finish")

}

func (self *Node) consumeResponse(responseMessage *messages.ResponseMessage) {

	//	messages.RegisterEvent("consumeResponse start")

	if self.consumer == nil {
		fmt.Printf("\nNo consumer registered at node %s\n\n", self.Id.Hex())
		return
	}
	sequence := responseMessage.Sequence
	payload := responseMessage.Payload

	//	messages.RegisterEvent("consumeResponse passed to consumer.Consume")

	self.consumer.Consume(sequence, payload, nil)

	//	messages.RegisterEvent("consumeResponse finish")

}

func (self *Node) sendResponse(backRoute messages.RouteId, sequence uint32, responsePayload []byte) {

	//	messages.RegisterEvent("sendResponse start")

	responseMessage := messages.ResponseMessage{sequence, responsePayload}

	//	messages.RegisterEvent("sendResponse serializing")

	responseSerialized := messages.Serialize(messages.MsgResponseMessage, responseMessage)

	//	messages.RegisterEvent("sendResponse finished serializing")

	_, err := self.getRoute(backRoute)
	if err != nil {
		fmt.Println("Wrong back route ID:", backRoute)
		return
	}
	inRouteMessage := messages.InRouteMessage{
		TransportId: messages.NIL_TRANSPORT,
		RouteId:     backRoute,
		Datagram:    responseSerialized,
	}
	//	inRouteS := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)

	//	messages.RegisterEvent("sendResponse passes to InjectTransportMessage")

	go self.InjectTransportMessage(&inRouteMessage)

	//	messages.RegisterEvent("sendResponse finish")

}

func (self *Node) AssignConsumer(consumer messages.Consumer) {
	self.consumer = consumer
}

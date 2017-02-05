package node

import (
	errmsg "github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/transport"

	"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
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
	IncomingChannel        chan []byte
	incomingControlChannel chan messages.InControlMessage

	Transports           map[messages.TransportId]*transport.Transport
	RouteForwardingRules map[messages.RouteId]*RouteRule
	consumer             messages.Consumer

	controlChannels map[messages.ChannelId]*ControlChannel
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
	node.IncomingChannel = make(chan []byte, 1024)
	node.incomingControlChannel = make(chan messages.InControlMessage, 1024)
	node.Transports = make(map[messages.TransportId]*transport.Transport)
	node.RouteForwardingRules = make(map[messages.RouteId]*RouteRule)
	node.controlChannels = make(map[messages.ChannelId]*ControlChannel)
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
	close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *Node) Tick() {
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
		if messages.IsDebug() {
			fmt.Printf("\nnode with id %s accepting a message with type %d\n\n", self.Id.Hex(), messages.GetMessageType(msg))
		}
		//process our incoming messages
		switch messages.GetMessageType(msg) {
		//InRouteMessage is the only message coming in to node from transports
		case messages.MsgInRouteMessage:
			var m1 messages.InRouteMessage
			messages.Deserialize(msg, &m1)
			if messages.IsDebug() {
				fmt.Println("InRouteMessage", m1)
			}
			self.handleInRouteMessage(m1)
		default:
			fmt.Println("wrong type", messages.GetMessageType(msg))
		}
	}
}

func (self *Node) handleIncomingControlMessages() {
	for msg := range self.incomingControlChannel {
		self.handleControlMessage(msg.ChannelId, msg.PayloadMessage)
		msg.ResponseChannel <- true
	}
}

func (self *Node) handleInRouteMessage(m1 messages.InRouteMessage) {
	//look in route table
	routeId := m1.RouteId
	transportId := m1.TransportId //who is is from
	//check that transport exists
	if _, ok := self.Transports[transportId]; !ok && transportId != messages.NIL_TRANSPORT {
		log.Printf("Node %s received message From Transport that does not exist\n", self.Id.Hex())
	}
	//check if route exists
	routeRule, ok := self.RouteForwardingRules[routeId]
	if !ok {
		log.Printf("Node %s received route message for route that does not exist\n", self.Id.Hex())
		panic("")
	}
	//check that incoming transport is correct
	if transportId != routeRule.IncomingTransport {
		log.Panic("Node: incoming route does not match the transport id it should be received from")
	}
	if routeId != routeRule.IncomingRoute {
		log.Panic("Node: impossible, incoming route id does not match route rule id")
	}
	//construct new packet
	outgoingRouteId := routeRule.OutgoingRoute
	outgoingTransportId := routeRule.OutgoingTransport //find transport to resend datagram
	datagram := m1.Datagram

	if outgoingRouteId == messages.NIL_ROUTE && outgoingTransportId == messages.NIL_TRANSPORT {
		self.consume(datagram)
	} else {
		var out messages.OutRouteMessage //replace inRoute, with outRoute, using rule
		out.RouteId = outgoingRouteId
		out.Datagram = datagram
		//serialize message, with prefix
		b1 := messages.Serialize(messages.MsgOutRouteMessage, out)

		if outgoingTransport, found := self.Transports[outgoingTransportId]; found {
			outgoingTransport.IncomingChannel <- b1 //inject message to transport
		}
	}
}

//inject an incoming message from the transport
func (self *Node) InjectTransportMessage(msg []byte) {
	self.IncomingChannel <- msg //push message to channel
}

func (self *Node) InjectControlMessage(msg messages.InControlMessage) {
	msg.ResponseChannel = make(chan bool, 1024)
	self.incomingControlChannel <- msg

	<-msg.ResponseChannel
}

func (self *Node) GetId() cipher.PubKey {
	return self.Id
}

func (self *Node) GetTransportById(id messages.TransportId) (*transport.Transport, error) {
	tr, ok := self.Transports[id]
	if !ok {
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
	self.Transports[id] = tr.(*transport.Transport)
}

func (self *Node) consume(msg []byte) {
	if messages.IsDebug() {
		fmt.Printf("\nNode %s consumed message %d\n\n", self.Id.Hex(), msg)
	}

	go func() {
		switch messages.GetMessageType(msg) {

		case messages.MsgRequestMessage:
			requestMessage := &messages.RequestMessage{}
			err := messages.Deserialize(msg, requestMessage)
			if err != nil {
				// send wrong request format
			} else {
				self.consumeRequest(requestMessage)
			}

		case messages.MsgResponseMessage:
			responseMessage := &messages.ResponseMessage{}
			err := messages.Deserialize(msg, responseMessage)
			if err != nil {
				// send wrong request format
			} else {
				self.consumeResponse(responseMessage)
			}

		default:
			fmt.Println("Incorrect message type:", messages.GetMessageType(msg))
		}
	}()
}

func (self *Node) consumeRequest(requestMessage *messages.RequestMessage) {
	if self.consumer == nil {
		fmt.Printf("\nNo consumer registered at node %s\n\n", self.Id.Hex())
		return
	}
	backRoute := requestMessage.BackRoute
	sequence := requestMessage.Sequence
	payload := requestMessage.Payload
	responseChannel := make(chan []byte)
	self.consumer.Consume(sequence, payload, responseChannel)
	responsePayload := <-responseChannel
	self.sendResponse(backRoute, sequence, responsePayload)
}

func (self *Node) consumeResponse(responseMessage *messages.ResponseMessage) {
	if self.consumer == nil {
		fmt.Printf("\nNo consumer registered at node %s\n\n", self.Id.Hex())
		return
	}
	sequence := responseMessage.Sequence
	payload := responseMessage.Payload
	self.consumer.Consume(sequence, payload, nil)
}

func (self *Node) sendResponse(backRoute messages.RouteId, sequence uint32, responsePayload []byte) {
	responseMessage := messages.ResponseMessage{sequence, responsePayload}
	responseSerialized := messages.Serialize(messages.MsgResponseMessage, responseMessage)
	_, ok := self.RouteForwardingRules[backRoute]
	if !ok {
		fmt.Println("Wrong back route ID:", backRoute)
		return
	}
	inRouteMessage := messages.InRouteMessage{
		TransportId: messages.NIL_TRANSPORT,
		RouteId:     backRoute,
		Datagram:    responseSerialized,
	}
	inRouteS := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)
	self.InjectTransportMessage(inRouteS)
}

func (self *Node) AssignConsumer(consumer messages.Consumer) {
	self.consumer = consumer
}

package node

import (
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/transport"

	"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"sync"
	"time"
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
	incomingControlChannel chan messages.InControlMessage
	incomingFromTransport  chan *messages.InRouteMessage
	incomingFromConnection chan *messages.InRouteMessage
	congestionChannel      chan *messages.CongestionPacket

	Transports map[messages.TransportId]*transport.Transport

	RouteForwardingRules map[messages.RouteId]*RouteRule

	lock *sync.Mutex

	user messages.User

	controlChannels map[messages.ChannelId]*ControlChannel

	Host string
	Port uint32

	Congested          bool
	connectionThrottle uint32
	throttle           uint32

	Sent      uint32
	Responses uint32

	Ticks int
}

type RouteRule struct {
	IncomingTransport messages.TransportId
	OutgoingTransport messages.TransportId
	IncomingRoute     messages.RouteId
	OutgoingRoute     messages.RouteId
}

var config = messages.GetConfig()

func NewNode() *Node {
	node := new(Node)
	node.Id = createPubKey()
	node.incomingFromTransport = make(chan *messages.InRouteMessage, config.MaxBuffer)
	node.incomingFromConnection = make(chan *messages.InRouteMessage, config.MaxBuffer*2)
	node.incomingControlChannel = make(chan messages.InControlMessage, 256)
	node.congestionChannel = make(chan *messages.CongestionPacket, 1024)
	node.Transports = make(map[messages.TransportId]*transport.Transport)
	node.RouteForwardingRules = make(map[messages.RouteId]*RouteRule)
	node.controlChannels = make(map[messages.ChannelId]*ControlChannel)
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
}

//move node forward on tick, process events
func (self *Node) Tick() {
	backChannel := make(chan bool, 32)
	self.runCycles(backChannel)
	<-backChannel
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
	return nil, messages.ERR_NO_TRANSPORT_TO_NODE
}

func (self *Node) ConnectedTo(other messages.NodeInterface) bool {
	_, err := self.GetTransportToNode(other.GetId())
	return err == nil
}

func (self *Node) SetTransport(id messages.TransportId, tr messages.TransportInterface) {
	self.setTransport(id, tr.(*transport.Transport))
}

func (self *Node) AssignUser(user messages.User) {
	self.user = user
}

//inject an incoming message from the connection
func (self *Node) InjectConnectionMessage(msg *messages.InRouteMessage) {
	c := cap(self.incomingFromConnection)
	for {
		if len(self.incomingFromConnection) < c {
			for {
				if len(self.incomingFromConnection) > c/2 {
					self.connectionThrottle = messages.Increase(self.connectionThrottle)
					time.Sleep(time.Duration(self.connectionThrottle) * config.TimeUnit)
					self.Congested = true
				} else {
					if len(self.incomingFromConnection) < c/4 {
						self.Congested = false
						self.connectionThrottle = messages.Decrease(self.connectionThrottle)
					}
					self.incomingFromConnection <- msg
					break
				}
			}
			break
		} else {
			fmt.Println("node is congested by connection packets")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

//inject an incoming message from the transport
func (self *Node) InjectTransportMessage(msg *messages.InRouteMessage) {
	go self.handleInRouteMessage(msg)
}

func (self *Node) InjectControlMessage(msg messages.InControlMessage) {
	msg.ResponseChannel = make(chan bool, 32)
	self.incomingControlChannel <- msg

	<-msg.ResponseChannel
}

func (self *Node) InjectCongestionPacket(msg *messages.CongestionPacket) {
	self.congestionChannel <- msg
}

//move node forward on tick, process events
func (self *Node) runCycles(backChannel chan bool) {
	//process incoming messages
	go self.handleCongestionMessages()         //pop them off the channel
	go self.handleIncomingConnectionMessages() //pop them off the channel
	go self.handleIncomingControlMessages()    //pop them off the channel
	backChannel <- true
}

func (self *Node) handleIncomingTransportMessages() {
	for msg := range self.incomingFromTransport {

		if messages.IsDebug() {
			fmt.Printf("\nnode with id %s accepting a message\n\n", self.Id.Hex())
			fmt.Println("InRouteMessage", msg)
		}

		go self.handleInRouteMessage(msg)
	}
}

func (self *Node) handleIncomingConnectionMessages() {
	for msg := range self.incomingFromConnection {

		if messages.IsDebug() {
			fmt.Printf("\nnode with id %s accepting a message\n\n", self.Id.Hex())
			fmt.Println("InRouteMessage", msg)
		}

		go self.handleInRouteMessage(msg)
	}
}

func (self *Node) handleIncomingControlMessages() {
	for msg := range self.incomingControlChannel {
		self.handleControlMessage(msg.ChannelId, msg.PayloadMessage)
		msg.ResponseChannel <- true
	}
}

func (self *Node) handleCongestionMessages() {
	for msg := range self.congestionChannel {
		if msg.Congestion {
			self.throttle = messages.Increase(self.throttle)
			fmt.Println("node throttle increased")
		} else {
			self.throttle = messages.Decrease(self.throttle)
		}
	}
}

func (self *Node) handleInRouteMessage(m1 *messages.InRouteMessage) {
	self.Ticks++
	//look in route table

	routeId := m1.RouteId
	transportId := m1.TransportId //who is is from
	//check that transport exists
	if _, err := self.getTransport(transportId); err != nil && transportId != messages.NIL_TRANSPORT {
		log.Printf("Node %s received message From Transport that does not exist\n", self.Id.Hex())
	}
	//check if route exists
	routeRule, err := self.getRoute(routeId)
	if err != nil {

		log.Fatalf("Node %s received route message for route that does not exist\n", self.Id.Hex())
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

		go self.out(datagram)

	} else {

		var out messages.OutRouteMessage //replace inRoute, with outRoute, using rule
		out.RouteId = outgoingRouteId
		out.Datagram = datagram
		//serialize message, with prefix

		if outgoingTransport, err := self.getTransport(outgoingTransportId); err == nil {

			if self.throttle > 0 {
				time.Sleep(time.Duration(self.throttle) * config.TimeUnit)
			}
			outgoingTransport.GetFromNode(out) //inject message to transport

		}
	}
}

func (self *Node) out(msg []byte) {
	if self.user != nil {
		self.user.Use(msg)
	}
}

func (self *Node) GetTicks() int {
	ticks := 0
	for _, tr := range self.Transports {
		ticks += tr.Ticks
	}
	ticks += self.Ticks
	return ticks
}

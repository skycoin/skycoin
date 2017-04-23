package node

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/transport"

	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
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
	incomingControlChannel chan *CM
	incomingFromTransport  chan *messages.InRouteMessage
	incomingFromConnection chan *messages.InRouteMessage
	congestionChannel      chan *messages.CongestionPacket

	responseChannels map[uint32]chan bool

	transports        map[messages.TransportId]*transport.Transport
	transportsByNodes map[cipher.PubKey]*transport.Transport

	routeForwardingRules map[messages.RouteId]*messages.RouteRule

	lock *sync.Mutex

	connection *Connection

	controlChannels             map[messages.ChannelId]*ControlChannel
	closeControlMessagesChannel chan bool

	host string
	port uint32

	controlConn *net.UDPConn
	nmAddr      net.Addr

	congested          bool
	connectionThrottle uint32
	throttle           uint32

	sequence uint32

	ticks int

	maxBuffer     uint64
	maxPacketSize uint32
	timeUnit      time.Duration
}

type CM struct {
	msg      *messages.InControlMessage
	respChan chan bool
}

const (
	CM_PORT = 5998
)

var (
	CONTROL_TIMEOUT = 10000 * time.Millisecond
)

func ConnectToMeshnet(host, nmAddr string) (messages.Connection, error) { // maybe make connection not to random but to exact node (by pubkey)
	n, err := CreateAndConnectNode(host, nmAddr)
	if err != nil {
		return nil, err
	}
	return n.GetConnection(), nil
}

func CreateNode(host, nmAddr string) (messages.NodeInterface, error) { // public for test reasons
	node, err := createAndRegisterNode(host, nmAddr, false)
	return node, err
}

func CreateAndConnectNode(host, nmAddr string) (messages.NodeInterface, error) { // public for test reasons
	node, err := createAndRegisterNode(host, nmAddr, true)
	return node, err
}

func (self *Node) Shutdown() {
	self.controlConn.Close()
	self.closeControlMessagesChannel <- true
	for _, tr := range self.transports {
		tr.Shutdown()
	}
}

func (self *Node) GetConnection() messages.Connection {
	return self.connection
}

func createAndRegisterNode(fullhost, nmAddr string, connect bool) (messages.NodeInterface, error) {

	hostData := strings.Split(fullhost, ":")
	if len(hostData) != 2 {
		return nil, messages.ERR_INCORRECT_HOST
	}

	host, portStr := hostData[0], hostData[1]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	node := newNode(host)

	controlConn, err := node.openUDPforCM(port, nmAddr)
	if err != nil {
		return nil, err
	}
	node.controlConn = controlConn
	go node.receiveControlMessages()

	err = node.sendRegisterNodeToServer(fullhost, connect) //**** send registration request to nm
	if err != nil {
		return nil, err
	}
	if messages.IsDebug() {
		fmt.Printf("Created Node %s\n", node.Id.Hex())
	}
	return node, nil
}

func newNode(host string) *Node {
	node := new(Node)
	node.lock = &sync.Mutex{}
	node.incomingControlChannel = make(chan *CM, 256)
	node.congestionChannel = make(chan *messages.CongestionPacket, 1024)
	node.responseChannels = make(map[uint32]chan bool)
	node.transports = make(map[messages.TransportId]*transport.Transport)
	node.transportsByNodes = make(map[cipher.PubKey]*transport.Transport)
	node.routeForwardingRules = make(map[messages.RouteId]*messages.RouteRule)
	node.controlChannels = make(map[messages.ChannelId]*ControlChannel)
	node.addZeroControlChannel()
	node.host = host
	node.closeControlMessagesChannel = make(chan bool)
	node.Tick()
	return node
}

func newLocalNode() *Node { // for tests
	return newNode(messages.LOCALHOST)
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

func (self *Node) GetTransportToNode(nodeId cipher.PubKey) (messages.TransportInterface, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	tr, ok := self.transportsByNodes[nodeId]
	if !ok {
		return nil, messages.ERR_NO_TRANSPORT_TO_NODE
	}
	return tr, nil
}

func (self *Node) ConnectedTo(nodeId cipher.PubKey) bool {
	_, err := self.GetTransportToNode(nodeId)
	return err == nil
}

func (self *Node) InjectCongestionPacket(msg *messages.CongestionPacket) {
	self.congestionChannel <- msg
}

//inject an incoming message from the transport
func (self *Node) InjectTransportMessage(msg *messages.InRouteMessage) {
	go self.handleInRouteMessage(msg)
}

//inject an incoming message from the connection
func (self *Node) injectConnectionMessage(msg *messages.InRouteMessage) {
	c := cap(self.incomingFromConnection)
	for {
		if len(self.incomingFromConnection) < c {
			for {
				if len(self.incomingFromConnection) > c/2 {
					self.connectionThrottle = messages.Increase(self.connectionThrottle)
					time.Sleep(time.Duration(self.connectionThrottle) * self.timeUnit)
					self.congested = true
				} else {
					if len(self.incomingFromConnection) < c/4 {
						self.congested = false
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

func (self *Node) injectControlMessage(msg *messages.InControlMessage) {
	respChan := make(chan bool)
	cm := &CM{
		msg,
		respChan,
	}
	self.incomingControlChannel <- cm
	<-respChan
}

//move node forward on tick, process events
func (self *Node) runCycles(backChannel chan bool) {
	//process incoming messages
	go self.handleCongestionMessages()      //pop them off the channel
	go self.handleIncomingControlMessages() //pop them off the channel
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
	for cm := range self.incomingControlChannel {
		msg := cm.msg
		respChan := cm.respChan
		self.handleControlMessage(messages.ChannelId(0), msg)
		respChan <- true
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
	self.ticks++
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

		log.Printf("Node %s received route message for route that does not exist\n", self.Id.Hex())
		return
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
				time.Sleep(time.Duration(self.throttle) * self.timeUnit)
			}
			outgoingTransport.GetFromNode(out) //inject message to transport

		}
	}
}

func (self *Node) out(msg []byte) {
	if self.connection != nil {
		self.connection.use(msg)
	}
}

func (self *Node) register(ack *messages.RegisterNodeCMAck) {
	self.Id = ack.NodeId
	self.timeUnit = time.Duration(ack.TimeUnit) * time.Microsecond
	self.maxBuffer = ack.MaxBuffer
	self.incomingFromTransport = make(chan *messages.InRouteMessage, self.maxBuffer)
	self.incomingFromConnection = make(chan *messages.InRouteMessage, self.maxBuffer*2)
	self.maxPacketSize = ack.MaxPacketSize
	self.connection = self.newConnection()
	self.connection.sendInterval = time.Duration(ack.SendInterval) * time.Microsecond
	self.connection.timeout = time.Duration(ack.ConnectionTimeout) * time.Millisecond
	self.connection.address = self.Id
	go self.handleIncomingConnectionMessages() //pop them off the channel
}

func (self *Node) getResponseChannel(sequence uint32) (chan bool, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	responseChannel, ok := self.responseChannels[sequence]
	return responseChannel, ok
}

func (self *Node) setResponseChannel(sequence uint32, responseChannel chan bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.responseChannels[sequence] = responseChannel
}

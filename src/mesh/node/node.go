package node

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/transport"

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
	id       cipher.PubKey
	hostname string

	incomingControlChannel  chan *CM
	incomingFromTransport   chan *messages.InRouteMessage
	incomingFromConnections chan *messages.InRouteMessage
	congestionChannel       chan *messages.CongestionPacket

	responseChannels map[uint32]chan bool

	transports        map[messages.TransportId]*transport.Transport
	transportsByNodes map[cipher.PubKey]*transport.Transport

	routeForwardingRules map[messages.RouteId]*messages.RouteRule

	lock *sync.Mutex

	connections       map[messages.ConnectionId]*Connection
	sendInterval      time.Duration
	connectionTimeout time.Duration

	controlChannels             map[messages.ChannelId]*ControlChannel
	closeControlMessagesChannel chan bool

	host string
	port uint32

	controlConn *net.UDPConn
	serverAddrs []net.Addr

	congested          bool
	connectionThrottle uint32
	throttle           uint32

	sequence uint32

	ticks uint32

	maxBuffer     uint64
	maxPacketSize uint32
	timeUnit      time.Duration

	connectResponseSequence uint32
	connectResponseChannels map[uint32]chan bool

	connectionResponseSequence uint32
	connectionResponseChannels map[uint32]chan messages.ConnectionId

	appTalkPort string
	appSequence uint32
	appConns    map[string]net.Conn

	viscriptServer *NodeViscriptServer
}

type CM struct {
	msg      *messages.InControlMessage
	respChan chan bool
}

var (
	CONTROL_TIMEOUT = 10000 * time.Millisecond
)

func CreateNode(nodeConfig *NodeConfig) (messages.NodeInterface, error) { // public for test reasons
	node, err := createAndRegisterNode(nodeConfig, false)
	return node, err
}

func CreateAndConnectNode(nodeConfig *NodeConfig) (messages.NodeInterface, error) { // public for test reasons
	node, err := createAndRegisterNode(nodeConfig, true)
	return node, err
}

func createAndRegisterNode(nodeConfig *NodeConfig, connect bool) (messages.NodeInterface, error) {

	fullhost := nodeConfig.ClientAddr
	serverAddrs := nodeConfig.ServerAddrs
	appTalkPort := strconv.Itoa(nodeConfig.AppTalkPort)
	hostname := nodeConfig.Hostname

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

	node.appTalkPort = appTalkPort
	node.hostname = hostname

	for _, serverAddr := range serverAddrs {
		node.addServer(serverAddr)
	}

	controlConn, err := node.openUDPforCM(port)
	if err != nil {
		panic(err)
		return nil, err
	}
	node.controlConn = controlConn

	go node.receiveControlMessages()

	go node.listenForApps()

	err = node.sendRegisterNodeToServer(hostname, fullhost, connect)
	if err != nil {
		return nil, err
	}
	if messages.IsDebug() {
		log.Printf("Created Node %s\n", node.id.Hex())
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
	node.connections = make(map[messages.ConnectionId]*Connection)
	node.connectResponseChannels = make(map[uint32]chan bool)
	node.connectionResponseChannels = make(map[uint32]chan messages.ConnectionId)
	node.addZeroControlChannel()
	node.host = host
	node.closeControlMessagesChannel = make(chan bool)
	node.appConns = make(map[string]net.Conn)
	node.Tick()
	return node
}

func newLocalNode() *Node { // for tests
	return newNode(messages.LOCALHOST)
}

func (self *Node) Shutdown() {

	if self.viscriptServer != nil {
		self.viscriptServer.Shutdown()
	}

	for _, appConn := range self.appConns {
		appConn.Close()
	}

	for _, tr := range self.transports {
		tr.Shutdown()
	}

	self.controlConn.Close()
	self.closeControlMessagesChannel <- true
}

func (self *Node) Dial(address string, appIdFrom, appIdTo messages.AppId) (messages.Connection, error) {
	connId, err := self.sendConnectWithRouteToServer(address, appIdFrom, appIdTo)
	if err != nil {
		return nil, err
	}

	self.lock.Lock()

	conn, ok := self.connections[connId]
	self.lock.Unlock()
	if !ok {
		return nil, messages.ERR_CONNECTION_DOESNT_EXIST
	}

	return conn, nil
}

func (self *Node) ConnectDirectly(address string) error {
	return self.sendConnectDirectlyToServer(address)
}

func (self *Node) GetConnection(id messages.ConnectionId) messages.Connection {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.connections[id]
}

func (self *Node) AppTalkAddr() string {
	return self.host + ":" + self.appTalkPort
}

//move node forward on tick, process events
func (self *Node) Tick() {
	backChannel := make(chan bool, 32)
	self.runCycles(backChannel)
	<-backChannel
}

func (self *Node) GetTicks() uint32 {
	ticks := self.ticks
	for _, tr := range self.transports {
		ticks += tr.GetTicks()
	}
	return ticks
}

func (self *Node) Id() cipher.PubKey {
	return self.id
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
	c := cap(self.incomingFromConnections)
	for {
		if len(self.incomingFromConnections) < c {
			for {
				if len(self.incomingFromConnections) > c/2 {
					self.connectionThrottle = messages.Increase(self.connectionThrottle)
					time.Sleep(time.Duration(self.connectionThrottle) * self.timeUnit)
					self.congested = true
				} else {
					if len(self.incomingFromConnections) < c/4 {
						self.congested = false
						self.connectionThrottle = messages.Decrease(self.connectionThrottle)
					}
					self.incomingFromConnections <- msg
					break
				}
			}
			break
		} else {
			log.Println("node is congested by connection packets")
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
			log.Printf("\nnode with id %s accepting a message\n\n", self.id.Hex())
			log.Println("InRouteMessage", msg)
		}

		go self.handleInRouteMessage(msg)
	}
}

func (self *Node) handleIncomingConnectionMessages() {
	for msg := range self.incomingFromConnections {

		if messages.IsDebug() {
			log.Printf("\nnode with id %s accepting a message\n\n", self.id.Hex())
			log.Println("InRouteMessage", msg)
		}

		go self.handleInRouteMessage(msg)
	}
}

func (self *Node) handleIncomingControlMessages() {
	for cm := range self.incomingControlChannel {
		msg := cm.msg
		respChan := cm.respChan
		err := self.handleControlMessage(messages.ChannelId(0), msg)
		if err != nil {
			log.Println(err)
		}
		respChan <- true
	}
}

func (self *Node) handleCongestionMessages() {
	for msg := range self.congestionChannel {
		if msg.Congestion {
			self.throttle = messages.Increase(self.throttle)
			log.Println("node throttle increased")
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
		log.Printf("Node %s received message From Transport that does not exist\n", self.id.Hex())
	}
	//check if route exists
	routeRule, err := self.getRoute(routeId)
	if err != nil {

		log.Printf("Node %s received route message for route that does not exist\n", self.id.Hex())
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

	switch messages.GetMessageType(msg) {

	case messages.MsgConnectionMessage:
		connMsg := messages.ConnectionMessage{}
		err := messages.Deserialize(msg, &connMsg)
		if err != nil {
			log.Println("wrong connection message", msg)
			return
		}

		self.lock.Lock()
		conn, ok := self.connections[connMsg.ConnectionId]
		self.lock.Unlock()

		if ok {
			go conn.handleConnectionMessage(&connMsg)
		} else {
			log.Println("no connection with id", connMsg.ConnectionId)
		}

	case messages.MsgConnectionAck:
		connAck := messages.ConnectionAck{}
		err := messages.Deserialize(msg, &connAck)
		if err != nil {
			log.Println("wrong connection Ack", msg)
			return
		}

		self.lock.Lock()
		conn, ok := self.connections[connAck.ConnectionId]
		self.lock.Unlock()

		if ok {
			go conn.receiveAck(connAck.Sequence)
		}

	default:
		log.Println("wrong connection message", msg)
	}
}

func (self *Node) register(ack *messages.RegisterNodeCMAck) {
	self.id = ack.NodeId
	self.timeUnit = time.Duration(ack.TimeUnit) * time.Microsecond
	self.maxBuffer = ack.MaxBuffer
	self.incomingFromTransport = make(chan *messages.InRouteMessage, self.maxBuffer)
	self.incomingFromConnections = make(chan *messages.InRouteMessage, self.maxBuffer*2)
	self.maxPacketSize = ack.MaxPacketSize
	self.sendInterval = time.Duration(ack.SendInterval) * time.Microsecond
	self.connectionTimeout = time.Duration(ack.ConnectionTimeout) * time.Millisecond
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

func (self *Node) setConnectionOn(connId messages.ConnectionId) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if conn, ok := self.connections[connId]; ok {
		conn.status = CONNECTED
	}
}

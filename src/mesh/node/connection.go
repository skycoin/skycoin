package node

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Connection struct {
	id               messages.ConnectionId
	status           uint8
	nodeAttached     *Node
	routeId          messages.RouteId
	sequence         uint32
	ackChannels      map[uint32]chan bool
	incomingMessages map[uint32]map[uint32][]byte
	incomingCounter  map[uint32]uint32
	lock             *sync.Mutex
	throttle         uint32
	consumer         net.Conn
	errChan          chan error

	packetSize   int
	sendInterval time.Duration
	timeout      time.Duration
}

const (
	DISCONNECTED uint8 = 0
	CONNECTED    uint8 = 1
	CONNECTING   uint8 = 2
)

func (node *Node) newConnection(connId messages.ConnectionId, routeId messages.RouteId, appId messages.AppId) (*Connection, error) {

	conn := &Connection{
		id:               connId,
		routeId:          routeId,
		status:           CONNECTING,
		nodeAttached:     node,
		lock:             &sync.Mutex{},
		ackChannels:      make(map[uint32]chan bool),
		errChan:          make(chan error),
		incomingMessages: make(map[uint32]map[uint32][]byte),
		incomingCounter:  make(map[uint32]uint32),
	}

	node.lock.Lock()
	appIdStr := string(appId)
	if appIdStr != "" {
		appConn, ok := node.appConns[appIdStr]
		if ok {
			conn.consumer = appConn
			msg := messages.AssignConnectionNAM{connId}
			msgS := messages.Serialize(messages.MsgAssignConnectionNAM, msg)
			err := sendToAppConn(appConn, msgS)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, messages.ERR_APP_DOESNT_EXIST
		}
	}
	node.connections[connId] = conn
	node.lock.Unlock()

	conn.packetSize = int(node.maxPacketSize / 2)

	go conn.receiveErrors()
	conn.sendInterval = node.sendInterval
	conn.timeout = node.connectionTimeout
	return conn, nil
}

func (self *Connection) Send(msg []byte) error {
	packets := self.split(msg)
	total := len(packets)

	ackChannel := make(chan bool, 32)
	sequence := self.setAckChannel(ackChannel)

	for order, packet := range packets {
		time.Sleep(self.sendInterval)
		self.sendPacket(packet, sequence, uint32(order), uint32(total))
	}

	select {
	case <-ackChannel:
		return nil
	case <-time.After(time.Duration(self.timeout) * time.Millisecond):
		self.errChan <- messages.ERR_CONN_TIMEOUT
		return messages.ERR_CONN_TIMEOUT
	}
}

func (self *Connection) Id() messages.ConnectionId {
	return self.id
}

func (self *Connection) Status() uint8 {
	return self.status
}

func (self *Connection) Close() {
	self.status = DISCONNECTED
}

func (self *Connection) split(msg []byte) [][]byte {
	packetSize := self.packetSize
	msgSize := len(msg)
	packets := [][]byte{}
	num := (msgSize-1)/packetSize + 1
	var start, end int
	for i := 0; i < num; i++ {
		start = i * packetSize
		if i < num-1 {
			end = (i + 1) * packetSize
		} else {
			end = msgSize
		}
		packet := msg[start:end]
		packets = append(packets, packet)
	}
	return packets
}

func (self *Connection) sendPacket(msg []byte, sequence, order, total uint32) {

	if self.status != CONNECTED {

		self.errChan <- messages.ERR_DISCONNECTED
	}

	connMessage := messages.ConnectionMessage{
		Sequence:     sequence,
		ConnectionId: self.id,
		Order:        order,
		Total:        total,
		Payload:      msg,
	}

	connMessageSerialized := messages.Serialize(messages.MsgConnectionMessage, connMessage)
	inRouteMessage := messages.InRouteMessage{
		messages.NIL_TRANSPORT,
		self.routeId,
		connMessageSerialized,
	}

	self.sendToNode(&inRouteMessage)
}

func (self *Connection) sendToNode(inRouteMessage *messages.InRouteMessage) {

	node := self.nodeAttached

	if node.congested {
		self.throttle = messages.Increase(self.throttle)
		time.Sleep(time.Duration(self.throttle) * self.nodeAttached.timeUnit)
	} else {
		self.throttle = messages.Decrease(self.throttle)
	}
	node.injectConnectionMessage(inRouteMessage)
}

func (self *Connection) receiveErrors() {
	for err := range self.errChan {
		if err != nil {
			log.Printf("Disconnect connection %d because of error %s\n", self.id, err.Error())
			self.status = DISCONNECTED
		}
	}
}

func (self *Connection) handleConnectionMessage(connMsg *messages.ConnectionMessage) {
	sequence := connMsg.Sequence
	order := connMsg.Order
	total := connMsg.Total
	payload := connMsg.Payload

	if _, err := self.getIncomingMessages(sequence); err != nil {
		self.createIncomingMessages(sequence)
	}

	self.setIncomingMessages(sequence, order, payload)
	self.increaseIncomingCounter(sequence)
	if self.getIncomingCounter(sequence) >= total {
		go self.sendAck(sequence)
		fullMessage := self.assemble(sequence, int(total))
		err := self.sendToConsumer(fullMessage)
		if err != nil {
			log.Println("sending message to consumer is failed:", err)
		}
	}
}

func (self *Connection) assemble(sequence uint32, total int) []byte {
	result := []byte{}
	for i := 0; i < total; i++ {
		packet := self.incomingMessages[sequence][uint32(i)]
		result = append(result, packet...)
	}
	return result
}

func (self *Connection) sendToConsumer(msg []byte) error {
	if self.consumer != nil {
		err := sendToAppConn(self.consumer, msg)
		return err
	} else {
		return nil
	}
}

func (self *Connection) sendAck(sequence uint32) {
	connAck := messages.ConnectionAck{
		Sequence:     sequence,
		ConnectionId: self.id,
	}

	connAckSerialized := messages.Serialize(messages.MsgConnectionAck, connAck)

	inRouteMessage := messages.InRouteMessage{
		messages.NIL_TRANSPORT,
		self.routeId,
		connAckSerialized,
	}
	go self.sendToNode(&inRouteMessage)
}

func (self *Connection) receiveAck(sequence uint32) {
	ackChannel, err := self.getAckChannel(sequence)
	if err != nil {
		log.Printf("no response channel with id %d is found\n", sequence)
		return
	}
	ackChannel <- true
}

func (self *Connection) getAckChannel(sequence uint32) (chan bool, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	ch, ok := self.ackChannels[sequence]
	if !ok {
		return nil, messages.ERR_NO_CLIENT_RESPONSE_CHANNEL
	}
	return ch, nil
}

func (self *Connection) setAckChannel(ackChannel chan bool) uint32 {
	self.lock.Lock()
	defer self.lock.Unlock()

	sequence := self.sequence
	self.sequence++
	self.ackChannels[sequence] = ackChannel
	return sequence
}

func (self *Connection) getIncomingMessages(sequence uint32) (map[uint32][]byte, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	im, ok := self.incomingMessages[sequence]
	if !ok {
		return nil, messages.ERR_NO_CLIENT_RESPONSE_CHANNEL
	}
	return im, nil
}

func (self *Connection) createIncomingMessages(sequence uint32) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.incomingMessages[sequence] = make(map[uint32][]byte)
	return
}

func (self *Connection) setIncomingMessages(sequence, order uint32, payload []byte) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.incomingMessages[sequence][order] = payload
	return
}

func (self *Connection) getIncomingCounter(sequence uint32) uint32 {
	self.lock.Lock()
	defer self.lock.Unlock()

	counter, _ := self.incomingCounter[sequence]
	return counter
}

func (self *Connection) increaseIncomingCounter(sequence uint32) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.incomingCounter[sequence]++
}

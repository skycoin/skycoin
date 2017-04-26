package node

import (
	"fmt"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type Connection struct {
	address          cipher.PubKey
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
	consumer         messages.Consumer
	errChan          chan error

	packetSize   int
	sendInterval time.Duration
	timeout      time.Duration
}

const (
	DISCONNECTED uint8 = 0
	CONNECTED    uint8 = 1
)

func (node *Node) newConnection() *Connection {
	id := messages.RandConnectionId()

	conn := &Connection{
		id:               id,
		status:           DISCONNECTED,
		nodeAttached:     node,
		lock:             &sync.Mutex{},
		ackChannels:      make(map[uint32]chan bool),
		errChan:          make(chan error),
		incomingMessages: make(map[uint32]map[uint32][]byte),
		incomingCounter:  make(map[uint32]uint32),
	}

	conn.packetSize = int(node.maxPacketSize / 2)

	go conn.receiveErrors()
	node.connection = conn
	return conn
}

func (self *Connection) Address() cipher.PubKey {
	return self.address
}

func (self *Connection) Shutdown() {
	self.nodeAttached.Shutdown()
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

func (self *Connection) Dial(address cipher.PubKey) error {
	if self.nodeAttached != nil {
		return self.nodeAttached.sendConnectToServer(address)
	} else {
		return nil
	}
}

func (self *Connection) AssignConsumer(consumer messages.Consumer) {
	self.consumer = consumer
}

func (self *Connection) GetStatus() uint8 {
	return self.status
}

func (self *Connection) Close() {
	self.status = DISCONNECTED
}

func (self *Connection) use(msg []byte) {

	switch messages.GetMessageType(msg) {

	case messages.MsgConnectionMessage:
		connMsg := messages.ConnectionMessage{}
		err := messages.Deserialize(msg, &connMsg)
		if err != nil {
			fmt.Println("wrong connection message", msg)
			return
		}

		go self.handleConnectionMessage(&connMsg)

	case messages.MsgConnectionAck:
		connAck := messages.ConnectionAck{}
		err := messages.Deserialize(msg, &connAck)
		if err != nil {
			fmt.Println("wrong connection Ack", msg)
			return
		}

		go self.receiveAck(connAck.Sequence)
	}
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
		Sequence: sequence,
		Order:    order,
		Total:    total,
		Payload:  msg,
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

	node := self.nodeAttached //**** should connection know it's node, not id? yes!

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
			fmt.Printf("Disconnect connection %d because of error %s\n", self.id, err.Error())
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
		self.sendToConsumer(fullMessage)
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

func (self *Connection) sendToConsumer(msg []byte) {
	if self.consumer != nil {
		self.consumer.Consume(msg)
	}
}

func (self *Connection) sendAck(sequence uint32) {
	connAck := messages.ConnectionAck{
		Sequence: sequence,
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
		fmt.Printf("no response channel with id %d is found\n", sequence)
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

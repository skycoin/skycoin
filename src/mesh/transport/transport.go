package transport

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

// Transport

//TODO:
// - implement simulated out of order packet delivery
// - implement simulated packet drop

type Transport struct {
	id messages.TransportId

	AttachedNode messages.NodeInTransport //node the transport is attached to

	pair             messages.TransportId //this is the other transport pair
	packetsSent      uint32
	packetsConfirmed uint32 // last confirmed ack

	status uint8

	simulateDelay     bool //
	maxSimulatedDelay int  // stub for testing

	incomingFromPair chan ([]byte)
	incomingFromNode chan messages.OutRouteMessage
	pendingOut       chan (*Job) //messages to send to other end of transport
	ackChannels      map[uint32]chan bool
	errChan          chan error

	throttle        uint32 // delay to send to pair
	pendingThrottle uint32

	dispatcher *Dispatcher

	timeout         uint32
	retransmitLimit int
	timeUnit        time.Duration

	nodeCongestion      bool //are we congested?
	transportCongestion bool

	udp *UDPConfig

	lock *sync.Mutex

	ticks uint32
}

const (
	DISCONNECTED = iota
	CONNECTED
)

func init() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
}

//are created by the factories
func CreateTransportFromMessage(msg *messages.TransportCreateCM) *Transport {
	tr := Transport{}
	tr.id = msg.Id
	tr.pair = msg.PairId
	maxBuffer := msg.MaxBuffer
	tr.incomingFromPair = make(chan []byte, maxBuffer*2)
	tr.incomingFromNode = make(chan messages.OutRouteMessage, maxBuffer*2)
	tr.pendingOut = make(chan *Job, maxBuffer*2)
	tr.ackChannels = make(map[uint32]chan bool)
	tr.status = DISCONNECTED
	tr.simulateDelay = msg.SimulateDelay
	if msg.SimulateDelay {
		tr.maxSimulatedDelay = int(msg.MaxSimulatedDelay)
	}
	tr.timeout = msg.TransportTimeout
	tr.retransmitLimit = int(msg.RetransmitLimit)
	tr.timeUnit = time.Duration(msg.TimeUnit) * time.Microsecond
	tr.lock = &sync.Mutex{}
	if messages.IsDebug() {
		fmt.Printf("Created Transport: %d\n", tr.id)
	}
	tr.dispatcher = NewDispatcher(&tr, maxWorkers)
	return &tr
}

func (self *Transport) Id() messages.TransportId {
	return self.id
}

func (self *Transport) PacketsSent() uint32 {
	return self.packetsSent
}

func (self *Transport) PacketsConfirmed() uint32 {
	return self.packetsConfirmed
}

func (self *Transport) Shutdown() {
	if self.udp != nil {
		self.udp.closeConn()
	}
}

//move node forward on tick, process events
func (self *Transport) Tick() {
	self.dispatcher.Run()
	go self.receiveFromPair() // receiving messages
	go self.receiveErrors()
	go self.broadcastCongestion()
	self.udp.Tick() // run udp listen
}

func (self *Transport) GetFromNode(msg messages.OutRouteMessage) {
	self.ticks++
	go self.sendTransportDatagramTransfer(&msg)
}

func (self *Transport) GetTicks() uint32 {
	return self.ticks
}

//inject an incoming message from the transport
func (self *Transport) injectNodeMessage(msg *messages.InRouteMessage) {
	if self.AttachedNode != nil {
		go self.AttachedNode.InjectTransportMessage(msg)
	}
}

func (self *Transport) sendTransportDatagramTransfer(msg *messages.OutRouteMessage) {
	//get message and put into the queue to be sent out
	//prime message for transit between the two transport ends

	var m1b messages.TransportDatagramTransfer
	m1b.Datagram = msg.Datagram
	m1b.RouteId = msg.RouteId

	c := cap(self.pendingOut)
	for {
		if len(self.pendingOut) < c {
			for {
				if len(self.pendingOut) >= c/2 {
					self.nodeCongestion = true
					self.pendingThrottle = messages.Increase(self.pendingThrottle)
					time.Sleep(time.Duration(self.pendingThrottle) * 10 * time.Microsecond)
					// send congestion message to node
				} else {
					self.pendingOut <- &Job{&m1b}
					self.pendingThrottle = messages.Decrease(self.pendingThrottle)
					if len(self.pendingOut) < c/4 {
						self.nodeCongestion = false
					}
					break
				}
			}
			break
		} else {
			self.nodeCongestion = true
			fmt.Println("pending out is overloaded")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

//used internally by transport factory

func (self *Transport) sendPacket(msg *messages.TransportDatagramTransfer) {

	ackChannel := make(chan bool, 32)
	self.setAckChannel(msg.Sequence, ackChannel)

	msgS := messages.Serialize(messages.MsgTransportDatagramTransfer, *msg)

	retransmits := 0
	for {
		if self.status == DISCONNECTED {
			self.errChan <- messages.ERR_DISCONNECTED
			return
		}

		self.sendMessageToPair(msgS)

		select {
		case <-ackChannel:
			if messages.IsDebug() {
				fmt.Printf("message %d is successfully sent, attempt %d\n", msg, retransmits+1)
			}
			return

		case <-time.After(time.Duration(self.timeout) * time.Millisecond):
			retransmits++
			if retransmits >= self.retransmitLimit {
				self.errChan <- messages.ERR_TRANSPORT_TIMEOUT
			}
			fmt.Printf("message %d will be sent again, attempt %d\n", msg, retransmits+1)
		}
	}
}

func (self *Transport) sendMessageToPair(msg []byte) {

	if self.throttle > 0 {
		fmt.Printf("transport is throttling %d milliseconds\n", self.throttle)
		time.Sleep(time.Duration(self.throttle) * self.timeUnit)
	}
	go self.udp.send(msg)
}

func (self *Transport) getFromUDP(msg []byte) {
	c := cap(self.incomingFromPair)
	throttle := uint32(0)
	for {
		if len(self.incomingFromPair) < c {
			for {
				if len(self.incomingFromPair) > c/2 {
					self.transportCongestion = true
					throttle = messages.Increase(throttle)
					time.Sleep(time.Duration(throttle) * self.timeUnit)
				} else {
					if len(self.incomingFromPair) < c/4 {
						self.transportCongestion = false
					}
					self.incomingFromPair <- msg
					throttle = messages.Decrease(throttle)
					break
				}
			}
			break
		} else {
			self.transportCongestion = true
			fmt.Println("incomingFromPair is overloaded:", len(self.incomingFromPair), c, self.id)
			time.Sleep(100 * self.timeUnit)
		}
	}
	//	go self.checkTransportCongestion()
}

func (self *Transport) receiveFromPair() {
	for m0 := range self.incomingFromPair {
		self.ticks++

		if self.status == DISCONNECTED {
			break
		}

		msg := m0
		go self.handleReceived(msg)
	}
}

func (self *Transport) handleReceived(msg []byte) {

	//		go self.checkTransportCongestion()

	//process our incoming messages
	if messages.IsDebug() {
		fmt.Printf("\ntransport with id %d gets message %d\n\n", self.id, msg)
	}

	switch messages.GetMessageType(msg) {

	//Transport -> Transport messag
	case messages.MsgTransportDatagramTransfer:

		var m2 messages.TransportDatagramTransfer
		err := messages.Deserialize(msg, &m2)
		if err != nil {
			fmt.Printf("Cannot deserialize transport datagram: %s %s\n", err.Error(), string(msg[2:]))
		} else {
			self.acceptAndSendAck(&msg, &m2)
		}

	case messages.MsgTransportDatagramACK:

		var m3 messages.TransportDatagramACK
		err := messages.Deserialize(msg, &m3)
		if err != nil {
			fmt.Printf("Cannot deserialize transport datagram ACK: %s\n", err.Error())
		} else {
			self.receiveAck(&m3)
		}
	/*
		case messages.MsgCongestionPacket:
				fmt.Println("congestionpacket")
				var m4 messages.CongestionPacket
				err := messages.Deserialize(msg, &m4)
				if err != nil {
					fmt.Printf("Cannot deserialize congestion packet: %s\n", err.Error())
				} else {
					self.handleCongestion(&m4)
				}
	*/
	default:
		fmt.Println("incorrect message type for transport input")
	}
}

func (self *Transport) acceptAndSendAck(msg *[]byte, m2 *messages.TransportDatagramTransfer) {

	routeId := m2.RouteId
	sequence := m2.Sequence
	datagram := m2.Datagram

	go func() {
		msgToNode := messages.InRouteMessage{self.id, routeId, datagram}
		self.injectNodeMessage(&msgToNode)
	}()

	if self.simulateDelay {
		time.Sleep(time.Duration(rand.Intn(self.maxSimulatedDelay)) * time.Millisecond)
	} // simulating delay, testing purposes!

	go func() {
		ackMsg := messages.TransportDatagramACK{sequence, 0}

		ackSerialized := messages.Serialize(messages.MsgTransportDatagramACK, ackMsg)

		self.sendMessageToPair(ackSerialized)
	}()
}

func (self *Transport) receiveAck(m3 *messages.TransportDatagramACK) {

	lowestSequence := m3.LowestSequence
	ackChannel, err := self.getAckChannel(lowestSequence)
	if err != nil {
		panic(err)
	}

	ackChannel <- true

	if lowestSequence > self.packetsConfirmed {
		self.packetsConfirmed = lowestSequence
	}

	if messages.IsDebug() {
		fmt.Printf("transport %d sent %d packets and got %d acks\n", self.id, self.packetsSent, self.packetsConfirmed)
	}
}

func (self *Transport) OpenUDPConn(peer, pair *messages.Peer) error {
	udp, err := openUDPConn(self, peer, pair)
	if err != nil {
		return err
	}
	self.udp = udp
	self.status = CONNECTED
	return nil
}

func (self *Transport) receiveErrors() {
	for err := range self.errChan {
		if err != nil {
			fmt.Printf("error:%s\n", err.Error())
			self.status = DISCONNECTED
			panic(err)
		}
	}
}

func (self *Transport) getAckChannel(sequence uint32) (chan bool, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	ackChannel, ok := self.ackChannels[sequence]
	if !ok {
		return nil, messages.ERR_NO_TRANSPORT_ACK_CHANNEL
	}
	return ackChannel, nil
}

func (self *Transport) setAckChannel(sequence uint32, ackChannel chan bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.ackChannels[sequence] = ackChannel
}

func (self *Transport) deleteAckChannel(sequence uint32) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, ok := self.ackChannels[sequence]
	if !ok {
		return messages.ERR_NO_TRANSPORT_ACK_CHANNEL
	}

	delete(self.ackChannels, sequence)

	return nil
}

package transport

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

//Stub Transport

//TODO:
// - implement simulated "delay" for transport +
// - implement simulated out of order packet delivery
// - implement simulated packet drop
// - implement real UDP trasport +
// - implement status "connected/disconnected" +
// - implement ACKs +

// TODO:
// pendingOut channel may need to be on the transport factory itself
// more efficient to push to central location than to pull/poll the transport list
// TODO:
// - may be more efficient to replace pending out, with an array on TransportFactory
// - or with array on the transport (who is responsible for processing ACKs?)

//This is stub transport
type Transport struct {
	Id               messages.TransportId
	incomingFromPair chan ([]byte)
	incomingFromNode chan messages.OutRouteMessage
	pendingOut       chan (*Job) //messages to send to other end of transport
	//Note: pendingOut channel may need to be on the transport_factory
	ackChannels map[uint32]chan bool
	errChan     chan error

	AttachedNode messages.NodeInterface //node the transport is attached to

	StubPair         *Transport //this is the other transport stub pair
	PacketsSent      uint32
	PacketsConfirmed uint32 // last confirmed ack

	Status uint8

	SimulateDelay     bool //
	MaxSimulatedDelay int  // stub for testing

	throttle        uint32 // delay to send to stub pair
	pendingThrottle uint32

	dispatcher *Dispatcher

	nodeCongestion      bool //are we congested?
	transportCongestion bool

	udp *UDPConfig

	lock *sync.Mutex

	Ticks int
}

const (
	DISCONNECTED = iota
	CONNECTED
)

var (
	config              *messages.ConfigStruct
	SIMULATE_DELAY      bool
	MAX_SIMULATED_DELAY int
	TIMEOUT             uint32
	RETRANSMIT_LIMIT    int
	MAX_BUFFER          uint64
)

func init() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	config = messages.GetConfig()
	SIMULATE_DELAY = config.SimulateDelay
	MAX_SIMULATED_DELAY = config.MaxSimulatedDelay
	TIMEOUT = config.TransportTimeout // time for ack waiting
	RETRANSMIT_LIMIT = config.RetransmitLimit
	MAX_BUFFER = config.MaxBuffer
}

//are created by the factories
func newTransportStub() *Transport {
	tr := Transport{}
	//	tr.maxBuffer = MAX_BUFFER
	tr.incomingFromPair = make(chan []byte, MAX_BUFFER*2)
	tr.incomingFromNode = make(chan messages.OutRouteMessage, MAX_BUFFER*2)
	tr.pendingOut = make(chan *Job, MAX_BUFFER*2)
	tr.ackChannels = make(map[uint32]chan bool)
	tr.Id = messages.RandTransportId()
	tr.Status = DISCONNECTED
	tr.SimulateDelay = SIMULATE_DELAY
	if SIMULATE_DELAY {
		tr.MaxSimulatedDelay = int(MAX_SIMULATED_DELAY)
	}
	tr.lock = &sync.Mutex{}
	if messages.IsDebug() {
		fmt.Printf("Created Transport: %d\n", tr.Id)
	}
	tr.dispatcher = NewDispatcher(&tr, maxWorkers)
	return &tr
}

func (self *Transport) Shutdown(wg *sync.WaitGroup) {
	self.udp.closeConn()
	wg.Done()
}

//move node forward on tick, process events
func (self *Transport) Tick() {
	self.dispatcher.Run()
	go self.receiveFromPair() // receiving messages
	go self.receiveErrors()
	//go self.sendFromPending() // for testing purposes
	//process incoming messages
	//go self.receiveFromNode() // receiving messages
	go self.broadcastCongestion()
	self.udp.Tick() // run udp listen
}

//inject an incoming message from the transport
func (self *Transport) InjectNodeMessage(msg *messages.InRouteMessage) {
	if self.AttachedNode != nil {
		go self.AttachedNode.InjectTransportMessage(msg)
	}
}

func (self *Transport) GetFromNode(msg messages.OutRouteMessage) {
	self.Ticks++
	go self.sendTransportDatagramTransfer(&msg)
}

func (self *Transport) sendTransportDatagramTransfer(msg *messages.OutRouteMessage) {
	//get message and put into the queue to be sent out
	//prime message for transit between the two transport ends

	var m1b messages.TransportDatagramTransfer
	m1b.Datagram = msg.Datagram
	m1b.RouteId = msg.RouteId
	m1b.IsResponse = msg.IsResponse

	//	self.pendingOut <- &Job{&m1b}

	c := cap(self.pendingOut)
	for {
		if len(self.pendingOut) < c {
			for {
				if len(self.pendingOut) >= c/2 {
					self.nodeCongestion = true
					//					fmt.Println("node congestion became true", self.Id)
					self.pendingThrottle = messages.Increase(self.pendingThrottle)
					time.Sleep(time.Duration(self.pendingThrottle) * 10 * time.Microsecond)
					// send congestion message to node
				} else {
					self.pendingOut <- &Job{&m1b}
					self.pendingThrottle = messages.Decrease(self.pendingThrottle)
					if len(self.pendingOut) < c/4 {
						self.nodeCongestion = false
						//						fmt.Println("node congestion became false", self.Id)
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

//message from stub to stub
//used internally by transport factory

func (self *Transport) sendPacket(msg *messages.TransportDatagramTransfer) {

	ackChannel := make(chan bool, 32)
	self.setAckChannel(msg.Sequence, ackChannel)

	msgS := messages.Serialize(messages.MsgTransportDatagramTransfer, *msg)

	retransmits := 0
	for {
		if self.Status == DISCONNECTED {
			self.errChan <- messages.ERR_DISCONNECTED
			return
		}

		self.sendMessageToStubPair(msgS)
		/*
			if err != nil {
				self.errChan <- err
				return
			}
		*/

		//		fmt.Println("WAITING FOR ACK #", msg.Sequence, self.Id)
		select {
		case <-ackChannel:
			if messages.IsDebug() {
				fmt.Printf("msg %d is successfully sent, attempt %d\n", msg, retransmits+1)
			}
			return

		case <-time.After(time.Duration(TIMEOUT) * time.Millisecond):
			retransmits++
			if retransmits >= RETRANSMIT_LIMIT {
				self.errChan <- messages.ERR_TRANSPORT_TIMEOUT
			}
			fmt.Printf("msg %d will be sent again, attempt %d\n", msg, retransmits+1)
		}
	}
}

func (self *Transport) sendMessageToStubPair(msg []byte) {
	//	self.StubPair.incomingFromPair <- msg
	//	return

	if self.throttle > 0 {
		fmt.Printf("transport is throttling %d milliseconds\n", self.throttle)
		time.Sleep(time.Duration(self.throttle) * config.TimeUnit)
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
					time.Sleep(time.Duration(throttle) * config.TimeUnit)
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
			fmt.Println("incomingFromPair is overloaded:", len(self.incomingFromPair), c, self.Id)
			time.Sleep(100 * config.TimeUnit)
		}
	}
	//go self.handleReceived(msg)
	//	go self.checkTransportCongestion()
}

func (self *Transport) receiveFromPair() {
	for m0 := range self.incomingFromPair {
		self.Ticks++

		if self.Status == DISCONNECTED {
			break
		}

		msg := m0
		//		fmt.Println("Incoming: ", string(msg))
		go self.handleReceived(msg)
	}
}

func (self *Transport) handleReceived(msg []byte) {

	//		go self.checkTransportCongestion()

	//process our incoming messages
	if messages.IsDebug() {
		fmt.Printf("\ntransport with id %d gets message %d\n\n", self.Id, msg)
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
		msgToNode := messages.InRouteMessage{self.Id, routeId, datagram}
		self.InjectNodeMessage(&msgToNode)
	}()

	if self.SimulateDelay {
		time.Sleep(time.Duration(rand.Intn(self.MaxSimulatedDelay)) * time.Millisecond)
	} // simulating delay, testing purposes!

	go func() {
		ackMsg := messages.TransportDatagramACK{sequence, 0}

		ackSerialized := messages.Serialize(messages.MsgTransportDatagramACK, ackMsg)

		self.sendMessageToStubPair(ackSerialized)
		/*
			if err != nil {
				self.errChan <- err
			}
		*/

	}()
}

func (self *Transport) receiveAck(m3 *messages.TransportDatagramACK) {

	lowestSequence := m3.LowestSequence
	ackChannel, err := self.getAckChannel(lowestSequence)
	if err != nil {
		panic(err)
	}

	ackChannel <- true

	if lowestSequence > self.PacketsConfirmed {
		self.PacketsConfirmed = lowestSequence
	}

	if messages.IsDebug() {
		fmt.Printf("transport %d sent %d packets and got %d acks\n", self.Id, self.PacketsSent, self.PacketsConfirmed)
	}
}

func (self *Transport) openUDPConn(peer, pair *messages.Peer) error {
	udp, err := openUDPConn(self, peer, pair)
	if err != nil {
		return err
	}
	self.udp = udp
	return nil
}

func (self *Transport) receiveErrors() {
	for err := range self.errChan {
		if err != nil {
			fmt.Printf("error:%s\n", err.Error())
			self.Status, self.StubPair.Status = DISCONNECTED, DISCONNECTED
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

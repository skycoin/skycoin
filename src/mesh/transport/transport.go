package transport

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/errors"
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
	Id              messages.TransportId
	incomingChannel chan ([]byte)
	pendingOut      chan (messages.TransportDatagramTransfer) //messages to send to other end of transport
	//Note: pendingOut channel may need to be on the transport_factory
	ackChannels map[uint32]chan bool

	AttachedNode messages.NodeInterface //node the transport is attached to

	StubPair         *Transport //this is the other transport stub pair
	PacketsSent      uint32
	PacketsConfirmed uint32 // last confirmed ack

	Status uint8

	MaxSimulatedDelay int // stub for testing

	udp *UDPConfig
}

const (
	DISCONNECTED = iota
	CONNECTED
)

var (
	config              *messages.ConfigStruct
	MAX_SIMULATED_DELAY int
	TIMEOUT             uint32
	RETRANSMIT_LIMIT    int
)

func init() {
	config = messages.GetConfig()
	MAX_SIMULATED_DELAY = config.MaxSimulatedDelay
	TIMEOUT = config.TransportTimeout // time for ack waiting
	RETRANSMIT_LIMIT = config.RetransmitLimit
}

//are created by the factories
func newTransportStub() *Transport {
	tr := &Transport{}
	tr.incomingChannel = make(chan []byte, 1024)
	tr.pendingOut = make(chan messages.TransportDatagramTransfer, 1024)
	tr.ackChannels = make(map[uint32]chan bool)
	tr.Id = messages.RandTransportId()
	tr.Status = DISCONNECTED
	tr.MaxSimulatedDelay = int(MAX_SIMULATED_DELAY)
	if messages.IsDebug() {
		fmt.Printf("Created Transport: %d\n", tr.Id)
	}
	return tr
}

func (self *Transport) Shutdown(wg *sync.WaitGroup) {
	close(self.incomingChannel)
	self.udp.closeConn()
	wg.Done()
}

func (self *Transport) openConn(peer, pair *messages.Peer) error {
	udp, err := openConn(self, peer, pair)
	if err != nil {
		return err
	}
	self.udp = udp
	return nil
}

//move node forward on tick, process events
func (self *Transport) Tick() {
	go self.sendFromPending() // for testing purposes
	//process incoming messages
	go self.receiveIncoming() // receiving messages
	self.udp.Tick()           // run udp listen
}

func (self *Transport) receiveIncoming() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	for msg := range self.incomingChannel {
		if self.Status == DISCONNECTED {
			break
		}
		//process our incoming messages
		if messages.IsDebug() {
			fmt.Printf("\ntransport with id %d gets message %d\n\n", self.Id, msg)
		}

		switch messages.GetMessageType(msg) {

		//OutRouteMessage is the only message coming in to transports from node
		//Node -> Transport message
		case messages.MsgOutRouteMessage:

			var m1 messages.OutRouteMessage
			err := messages.Deserialize(msg, &m1)
			if err != nil {
				fmt.Printf("Cannot deserialize outroute message: %s\n", err.Error())
			}
			self.sendTransportDatagramTransfer(&m1)

		//Transport -> Transport messag
		case messages.MsgTransportDatagramTransfer:

			var m2 messages.TransportDatagramTransfer
			err := messages.Deserialize(msg, &m2)
			if err != nil {
				fmt.Printf("Cannot deserialize transport datagram: %s\n", err.Error())
			}
			err = self.acceptAndSendAck(&msg, &m2)
			if err != nil {
				fmt.Printf("transport %d isn't responding, error:%s\n", self.StubPair.Id, err.Error())
				self.Status, self.StubPair.Status = DISCONNECTED, DISCONNECTED
			}

		case messages.MsgTransportDatagramACK:
			err := self.receiveAck(msg)
			if err != nil {
				fmt.Printf("Incorrect ack message: %s\n", err.Error())
			}

		default:
			fmt.Println("incorrect message type for transport input")
		}
	}
}

func (self *Transport) sendTransportDatagramTransfer(msg *messages.OutRouteMessage) {
	//get message and put into the queue to be sent out
	//prime message for transit between the two transport ends
	self.PacketsSent++
	var m1b messages.TransportDatagramTransfer
	m1b.Datagram = msg.Datagram
	m1b.Sequence = self.PacketsSent
	m1b.RouteId = msg.RouteId

	self.pendingOut <- m1b //push to queue, to be transferred
}

func (self *Transport) acceptAndSendAck(msg *[]byte, m2 *messages.TransportDatagramTransfer) error {
	routeId := m2.RouteId
	sequence := m2.Sequence
	datagram := m2.Datagram

	msgToNode := messages.InRouteMessage{self.Id, routeId, datagram}
	serialized := messages.Serialize(messages.MsgInRouteMessage, msgToNode)
	self.InjectNodeMessage(serialized)

	time.Sleep(time.Duration(rand.Intn(self.MaxSimulatedDelay)) * time.Millisecond) // simulating delay, testing purposes!

	ackMsg := messages.TransportDatagramACK{sequence, 0}
	ackSerialized := messages.Serialize(messages.MsgTransportDatagramACK, ackMsg)
	err := self.udp.send(ackSerialized)
	return err
}

func (self *Transport) receiveAck(msg []byte) error {
	var m3 messages.TransportDatagramACK
	err := messages.Deserialize(msg, &m3)
	if err != nil {
		return err
	} else {
		lowestSequence := m3.LowestSequence
		ackChannel := self.ackChannels[lowestSequence]
		ackChannel <- true
		if lowestSequence > self.PacketsConfirmed {
			self.PacketsConfirmed = lowestSequence
		}
		if messages.IsDebug() {
			fmt.Printf("transport %d sent %d packets and got %d acks\n", self.Id, self.PacketsSent, self.PacketsConfirmed)
		}
	}
	return nil
}

func (self *Transport) sendFromPending() {
	for msg := range self.pendingOut {
		b1 := messages.Serialize(messages.MsgTransportDatagramTransfer, msg)
		sequence := msg.Sequence
		ackChannel := make(chan bool, 1024)
		self.ackChannels[sequence] = ackChannel
		err := self.sendPacket(b1, ackChannel)
		if err != nil {
			fmt.Printf("transport %d isn't responding, error:%s\n", self.StubPair.Id, err.Error())
			self.Status, self.StubPair.Status = DISCONNECTED, DISCONNECTED
		}
		close(ackChannel)
		delete(self.ackChannels, sequence)
	}
}

func (self *Transport) sendPacket(msg []byte, ackChannel chan bool) error {
	retransmits := 0
	for {
		if self.Status == DISCONNECTED {
			return errors.ERR_DISCONNECTED
		}
		err := self.sendMessageToStubPair(msg)
		if err != nil {
			return err
		}
		select {
		case <-ackChannel:
			if messages.IsDebug() {
				fmt.Printf("msg %d is successfully sent, attempt %d\n", msg, retransmits+1)
			}
			return nil
		case <-time.After(time.Duration(TIMEOUT) * time.Millisecond):
			retransmits++
			if retransmits >= RETRANSMIT_LIMIT {
				return errors.ERR_TIMEOUT
			}
			fmt.Printf("msg %d will be sent again, attempt %d\n", msg, retransmits+1)
		}
	}
}

//inject an incoming message from the transport
func (self *Transport) InjectNodeMessage(msg []byte) {
	if self.AttachedNode != nil {
		self.AttachedNode.InjectTransportMessage(msg)
	}
}

func (self *Transport) GetFromNode(msg []byte) {
	self.incomingChannel <- msg
}

//message from stub to stub
//used internally by transport factory
/*
func (self *Transport) sendMessageToStubPair(msg []byte) {
	self.StubPair.incomingChannel <- msg
}
*/

func (self *Transport) sendMessageToStubPair(msg []byte) error {
	return self.udp.send(msg)
}

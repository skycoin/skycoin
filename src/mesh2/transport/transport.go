package transport

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

//Stub Transport

//TODO:
// - implement simulated "delay" for transport +
// - implement simulated out of order packet delivery
// - implement simulated packet drop
// - implement real UDP trasport
// - implement status "connected/disconnected" +
// - implement ACKs +

// TODO:
// pendingOut channel may need to be on the transport factory itself
// more efficient to push to central location than to pull/poll the transport list
// TODO:
// - may be more efficient to replace pending out, with an array on TransportFactory
// - or with array on the transport (who is responsible for processing ACKs?)

const (
	TIMEOUT   uint32 = 1000 // time for ack waiting
	CONNECTED        = iota
	DISCONNECTED
)

//This is stub transport
type Transport struct {
	Id              messages.TransportId
	IncomingChannel chan ([]byte)
	PendingOut      chan ([]byte) //messages to send to other end of transport
	//Note: PendingOut channel may need to be on the transport_factory
	ackChannels map[string]chan ([]byte)

	AttachedNode messages.NodeInterface //node the transport is attached to

	StubPair         *Transport //this is the other transport stub pair
	PacketsSent      uint32
	PacketsConfirmed uint32 // last confirmed ack

	Status uint8

	MaxSimulatedDelay int // stub for testing
}

//are created by the factories
func (self *Transport) NewTransportStub() {
	self.IncomingChannel = make(chan []byte, 1024)
	self.PendingOut = make(chan []byte, 1024)
	self.ackChannels = make(map[string]chan []byte)
	self.Id = messages.RandTransportId()
	self.Status = DISCONNECTED
	self.MaxSimulatedDelay = 1000
	fmt.Printf("Created Transport: %d\n", self.Id)
}

func (self *Transport) Shutdown() {
	close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *Transport) Tick() {
	//process incoming messages
	go self.sendFromPending() // for testing purposes
	go self.receiveIncoming() // receiving messages
}

func (self *Transport) receiveIncoming() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	for msg := range self.IncomingChannel {
		if self.Status == DISCONNECTED {
			break
		}
		//process our incoming messages
		fmt.Printf("\ntransport with id %d gets message %d\n\n", self.Id, msg)

		switch messages.GetMessageType(msg) {

		//OutRouteMessage is the only message coming in to transports from node
		//Node -> Transport message
		case messages.MsgOutRouteMessage:

			var m1 messages.OutRouteMessage
			err := messages.Deserialize(msg, &m1)
			if err != nil {
				panic(err)
			}
			self.sendTransportDatagramTransfer(&m1)

		//Transport -> Transport messag
		case messages.MsgTransportDatagramTransfer:

			var m2 messages.TransportDatagramTransfer
			err := messages.Deserialize(msg, &m2)
			if err != nil {
				panic(err)
			}
			self.sendAck(&msg, &m2)

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

	b1 := messages.Serialize(messages.MsgTransportDatagramTransfer, m1b)
	self.PendingOut <- b1 //push to queue, to be transferred
}

func (self *Transport) sendAck(msg *[]byte, m2 *messages.TransportDatagramTransfer) {
	routeId := m2.RouteId
	sequence := m2.Sequence
	datagram := m2.Datagram

	msgToNode := messages.InRouteMessage{self.Id, routeId, datagram}
	serialized := messages.Serialize(messages.MsgInRouteMessage, msgToNode)
	self.InjectNodeMessage(serialized)

	time.Sleep(time.Duration(rand.Intn(self.MaxSimulatedDelay)) * time.Millisecond) // simulating delay, testing purposes!

	ackMsg := messages.TransportDatagramACK{sequence, 0}
	ackSerialized := messages.Serialize(messages.MsgTransportDatagramACK, ackMsg)
	ackChannel := self.StubPair.ackChannels[string(*msg)]
	ackChannel <- ackSerialized
}

func (self *Transport) receiveAck(msg []byte) {
	var m3 messages.TransportDatagramACK
	err := messages.Deserialize(msg, &m3)
	if err != nil {
		fmt.Println("incorrect ACK message")
	} else {
		lowestSequence := m3.LowestSequence
		if lowestSequence > self.PacketsConfirmed {
			self.PacketsConfirmed = lowestSequence
		}
		fmt.Printf("transport %d sent %d packets and got %d acks\n", self.Id, self.PacketsSent, self.PacketsConfirmed)
	}
}

func (self *Transport) sendFromPending() {
	for msg := range self.PendingOut {
		ackChannel := make(chan []byte, 1024)
		self.ackChannels[string(msg)] = ackChannel
		self.SendMessageToStubPair(msg)
		select {
		case ack := <-ackChannel:
			self.receiveAck(ack)
		case <-time.After(time.Duration(TIMEOUT) * time.Millisecond):
			fmt.Printf("transport %d isn't responding\n", self.StubPair.Id)
			self.Status, self.StubPair.Status = DISCONNECTED, DISCONNECTED
			break
		}
		close(ackChannel)
		delete(self.ackChannels, string(msg))
	}
}

//inject an incoming message from the transport
func (self *Transport) InjectNodeMessage(msg []byte) {
	if self.AttachedNode != nil {
		self.AttachedNode.InjectTransportMessage(self.Id, msg)
	}
}

//message from stub to stub
//used internally by transport factory
func (self *Transport) SendMessageToStubPair(msg []byte) {
	self.StubPair.IncomingChannel <- msg
}

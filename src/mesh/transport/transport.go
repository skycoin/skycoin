package transport

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/serialize"
)

type TransportConfig struct {
	MyPeerID                        cipher.PubKey
	PhysicalReceivedChannelLength   int
	ExpireMessagesInterval          time.Duration
	RememberMessageReceivedDuration time.Duration

	// If an ACK is not received for this long, the message is dropped
	RetransmitDuration time.Duration
}

type SendMessage struct {
	MessageID domain.MessageID
	FromPeer  cipher.PubKey
	Contents  []byte
}

type ReplyMessage struct {
	MessageID domain.MessageID
}

type messageSentState struct {
	toPeer      cipher.PubKey
	serialized  []byte
	expiryTime  time.Time
	receivedAck bool
}

// Wraps Transport, but adds store-and-forward
type Transport struct {
	config            TransportConfig
	physicalTransport ITransport
	outputChannel     chan []byte
	serializer        *serialize.Serializer

	lock             *sync.Mutex
	messagesSent     map[domain.MessageID]messageSentState
	messagesReceived map[domain.MessageID]time.Time
	nextMsgId        uint32

	physicalReceived chan []byte
	closing          chan bool
	closeWait        *sync.WaitGroup
}

func NewTransport(physicalTransport ITransport, config TransportConfig) *Transport {
	transport := &Transport{
		config,
		physicalTransport,
		nil,
		serialize.NewSerializer(),
		&sync.Mutex{},
		make(map[domain.MessageID]messageSentState),
		make(map[domain.MessageID]time.Time),
		1000,
		make(chan []byte, config.PhysicalReceivedChannelLength),
		make(chan bool, 10),
		&sync.WaitGroup{},
	}

	transport.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, SendMessage{})
	transport.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, ReplyMessage{})

	go transport.processReceivedLoop()
	go transport.expireMessagesLoop()
	go transport.retransmitLoop()

	transport.physicalTransport.SetReceiveChannel(transport.physicalReceived)
	return transport
}

func (self *Transport) processReceivedLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()

	for len(self.closing) == 0 {
		select {
		case physicalMsg, ok := <-self.physicalReceived:
			{
				if ok {
					self.processPhysicalMessage(physicalMsg)
				}
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}
func (self *Transport) doRetransmits() {
	self.lock.Lock()
	defer self.lock.Unlock()
	for _, state := range self.messagesSent {
		if !state.receivedAck {
			go self.physicalTransport.SendMessage(state.toPeer, state.serialized)
		}
	}
}

func (self *Transport) retransmitLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.config.RetransmitDuration):
			{
				self.doRetransmits()
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Transport) expireMessagesLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.config.ExpireMessagesInterval):
			{
				self.expireMessages()
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Transport) expireMessages() {
	timeNow := time.Now()
	self.lock.Lock()
	defer self.lock.Unlock()
	lastMessagesSent := self.messagesSent
	self.messagesSent = make(map[domain.MessageID]messageSentState)
	for messageID, messageSentState := range lastMessagesSent {
		if timeNow.Before(messageSentState.expiryTime) {
			self.messagesSent[messageID] = messageSentState
		}
	}

	lastReceived := self.messagesReceived
	self.messagesReceived = make(map[domain.MessageID]time.Time)
	for id, expiry := range lastReceived {
		if timeNow.Before(expiry) {
			self.messagesReceived[id] = expiry
		}
	}
}

func (self *Transport) sendAck(message SendMessage) {
	reply := ReplyMessage{message.MessageID}
	serialized := self.serializer.SerializeMessage(reply)
	go self.physicalTransport.SendMessage(message.FromPeer, serialized)
}

func (self *Transport) processSend(message SendMessage) {
	self.sendAck(message)

	self.lock.Lock()
	defer self.lock.Unlock()
	_, alreadyReceived := self.messagesReceived[message.MessageID]
	if !alreadyReceived {
		self.outputChannel <- message.Contents
		self.messagesReceived[message.MessageID] = time.Now().Add(self.config.RememberMessageReceivedDuration)
	}
}

func (self *Transport) processReply(message ReplyMessage) {
	self.lock.Lock()
	defer self.lock.Unlock()
	state, exists := self.messagesSent[message.MessageID]
	if !exists {
		fmt.Fprintf(os.Stderr, "Received ack for unknown sent message %v\n", message.MessageID)
		return
	}
	state.receivedAck = true
	self.messagesSent[message.MessageID] = state

	// Test
	if !self.messagesSent[message.MessageID].receivedAck {
		panic("Test error")
	}
}

func (self *Transport) processPhysicalMessage(physicalMessage []byte) {
	message, err := self.serializer.UnserializeMessage(physicalMessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Deserialization error %v\n", err)
		return
	}
	messageType := reflect.TypeOf(message)

	if messageType == reflect.TypeOf(SendMessage{}) {
		send := message.(SendMessage)
		self.processSend(send)
	} else if messageType == reflect.TypeOf(ReplyMessage{}) {
		reply := message.(ReplyMessage)
		self.processReply(reply)
	} else {
		panic("Internal error: unknown message type")
	}
}

func (self *Transport) newMessageID() domain.MessageID {
	return (domain.MessageID)(uuid.NewV4())
}

func (self *Transport) SendMessage(toPeer cipher.PubKey, contents []byte) error {
	messageID := self.newMessageID()
	sendMessage := SendMessage{messageID, self.config.MyPeerID, contents}
	sendSerialized := self.serializer.SerializeMessage(sendMessage)
	state := messageSentState{toPeer,
		sendSerialized,
		time.Now().Add(self.config.RetransmitDuration),
		false}
	err := self.physicalTransport.SendMessage(toPeer, sendSerialized)
	if err == nil {
		self.lock.Lock()
		defer self.lock.Unlock()
		self.messagesSent[messageID] = state
	}
	return err
}

func (self *Transport) SetReceiveChannel(received chan []byte) {
	self.outputChannel = received
}

func (self *Transport) Close() error {
	for i := 0; i < 10; i++ {
		self.closing <- true
	}
	self.closeWait.Wait()
	return self.physicalTransport.Close()
}

func (self *Transport) SetCrypto(crypto ITransportCrypto) {
	self.physicalTransport.SetCrypto(crypto)
}

func (self *Transport) ConnectedToPeer(peer cipher.PubKey) bool {
	return self.physicalTransport.ConnectedToPeer(peer)
}

func (self *Transport) GetConnectedPeers() []cipher.PubKey {
	return self.physicalTransport.GetConnectedPeers()
}

func (self *Transport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	empty := SendMessage{}
	emptySerialized := self.serializer.SerializeMessage(empty)
	if (uint)(len(emptySerialized)) >= self.physicalTransport.GetMaximumMessageSizeToPeer(peer) {
		return 0
	}
	return self.physicalTransport.GetMaximumMessageSizeToPeer(peer) - (uint)(len(emptySerialized))
}

func (self *Transport) debug_countMapItems() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return len(self.messagesSent) + len(self.messagesReceived)
}

// Create TransportConfig to the node.
func CreateTransportConfig(pubKey cipher.PubKey) TransportConfig {
	config := TransportConfig{}
	config.MyPeerID = pubKey
	config.PhysicalReceivedChannelLength = 100
	config.ExpireMessagesInterval = 5 * time.Minute
	config.RememberMessageReceivedDuration = 1 * time.Minute
	config.RetransmitDuration = 1 * time.Minute

	return config
}

package transport

import (
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/serialize"
)

// transport statuses
const (
	DISCONNECTED uint32 = iota
	CONNECTED
	SENDING
	RECEIVING
	REPLYING
	ACKWAITING
	TIMEOUT
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
	MessageID  domain.MessageID
	FromPeerID cipher.PubKey
	Contents   []byte
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

	id		uint32
	config            TransportConfig
	physicalTransport ITransport
	output            chan []byte
	serializer        *serialize.Serializer
	metadata	[]byte

	status		uint32

	lock             *sync.Mutex
	messagesSent     map[domain.MessageID]messageSentState
	messagesReceived map[domain.MessageID]time.Time
	nextMsgId        uint32

	packetIsSent	time.Time
	latency		uint64
	packetsCount	uint32
	packetsSent	uint32
	packetsReceived	uint32
	packetsRetransmissions	uint32

	physicalReceived chan []byte
	closing          chan bool
	closeWait        *sync.WaitGroup
}

var currentID uint32 = 0

func NewTransport(physicalTransport ITransport, config TransportConfig) *Transport {
	transport := &Transport{
		currentID,
		config,
		physicalTransport,
		nil,
		serialize.NewSerializer(),
		[]byte{},
		DISCONNECTED,
		&sync.Mutex{},
		make(map[domain.MessageID]messageSentState),
		make(map[domain.MessageID]time.Time),
		1000,
		time.Time{},
		0,
		0,
		0,
		0,
		0,
		make(chan []byte, config.PhysicalReceivedChannelLength),
		make(chan bool, 10),
		&sync.WaitGroup{},
	}

	currentID++

	transport.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, SendMessage{})
	transport.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, ReplyMessage{})

	go transport.processReceivedLoop()
	go transport.expireMessagesLoop()
	go transport.retransmitLoop()

	transport.physicalTransport.SetReceiveChannel(transport.physicalReceived)
	return transport
}

func (self *Transport) ConnectedToPeer(peer cipher.PubKey) bool {
	return self.physicalTransport.ConnectedToPeer(peer)
}

func (self *Transport) GetConnectedPeer() cipher.PubKey {
	return self.physicalTransport.GetConnectedPeer()
}

func (self *Transport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	empty := SendMessage{}
	emptySerialized := self.serializer.SerializeMessage(empty)
	if (uint)(len(emptySerialized)) >= self.physicalTransport.GetMaximumMessageSizeToPeer(peer) {
		return 0
	}
	return self.physicalTransport.GetMaximumMessageSizeToPeer(peer) - (uint)(len(emptySerialized))
}

func (self *Transport) Close() error {
	for i := 0; i < 10; i++ {
		self.closing <- true
	}
	self.closeWait.Wait()
	self.status = DISCONNECTED
	return self.physicalTransport.Close()
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

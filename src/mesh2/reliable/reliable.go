
package reliable

import(
	"os"
	"fmt"
	"time"
	"sync"
	"reflect"
)

import(
	"github.com/skycoin/skycoin/src/mesh2/transport"
	"github.com/skycoin/skycoin/src/mesh2/serialize"
	"github.com/skycoin/skycoin/src/cipher")

import ("github.com/satori/go.uuid")

type ReliableTransportConfig struct {
	MyPeerId						cipher.PubKey
	PhysicalReceivedChannelLength 	int
	ExpireMessagesInterval			time.Duration
	RememberMessageReceivedDuration	time.Duration

	// If an ACK is not received for this long, the message is dropped
	RetransmitDuration				time.Duration
}

// 0 is not nil
type reliableId uuid.UUID

type ReliableSend struct {
	MsgId    reliableId
	FromPeer cipher.PubKey
	Contents []byte
}

type ReliableReply struct {
	MsgId    reliableId
}

type messageSentState struct {
	toPeer       cipher.PubKey
	serialized   []byte
	expiryTime   time.Time
	receivedAck  bool
}

// Wraps Transport, but adds store-and-forward
type ReliableTransport struct {
	config              ReliableTransportConfig
	physicalTransport 	transport.Transport
	outputChannel 		chan []byte
    serializer 			*serialize.Serializer

	lock 				*sync.Mutex
	messagesSent        map[reliableId]messageSentState
	messagesReceived    map[reliableId]time.Time
	nextMsgId			uint32

	physicalReceived    chan []byte
	closing 			chan bool
	closeWait           *sync.WaitGroup
}

func NewReliableTransport(physicalTransport transport.Transport, config ReliableTransportConfig) *ReliableTransport {
	ret := &ReliableTransport{
		config,
		physicalTransport,
		nil,
		serialize.NewSerializer(),
		&sync.Mutex{},
		make(map[reliableId]messageSentState),
		make(map[reliableId]time.Time),
		1000,
		make(chan[]byte, config.PhysicalReceivedChannelLength),
		make(chan bool, 10),
		&sync.WaitGroup{},
	}

	ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, ReliableSend{})
	ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, ReliableReply{})

	go ret.processReceivedLoop()
	go ret.expireMessagesLoop()
	go ret.retransmitLoop()

	ret.physicalTransport.SetReceiveChannel(ret.physicalReceived)
	return ret
}

func (self*ReliableTransport) processReceivedLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()

	for len(self.closing) == 0 {
		select {
			case physicalMsg, ok := <- self.physicalReceived: {
				if ok {
					self.processPhysicalMessage(physicalMsg)
				}
			}
			case <-self.closing: {
				return
			}
		}
	}
}
func (self*ReliableTransport) doRetransmits() {
	self.lock.Lock()
	defer self.lock.Unlock()
	for _, state := range(self.messagesSent) {
		if !state.receivedAck {
			go self.physicalTransport.SendMessage(state.toPeer, state.serialized)
		}
	}
}

func (self*ReliableTransport) retransmitLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()
	for len(self.closing) == 0 {
		select {
			case <-time.After(self.config.RetransmitDuration): {
				self.doRetransmits()
			}
			case <-self.closing: {
				return
			}
		}
	}
}

func (self*ReliableTransport) expireMessagesLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()
	for len(self.closing) == 0 {
		select {
			case <-time.After(self.config.ExpireMessagesInterval): {
				self.expireMessages()
			}
			case <-self.closing: {
				return
			}
		}
	}
}

func (self*ReliableTransport) expireMessages() {
	time_now := time.Now()
	self.lock.Lock()
	defer self.lock.Unlock()
	lastSent := self.messagesSent
	self.messagesSent = make(map[reliableId]messageSentState)
	for id, state := range(lastSent) {
		if time_now.Before(state.expiryTime) {
			self.messagesSent[id] = state
		}
	}

	lastReceived := self.messagesReceived
	self.messagesReceived = make(map[reliableId]time.Time)
	for id, expiry := range(lastReceived) {
		if time_now.Before(expiry) {
			self.messagesReceived[id] = expiry
		}
	}
}

func (self*ReliableTransport) sendAck(msg ReliableSend) {
	reply := ReliableReply{msg.MsgId}
	serialized := self.serializer.SerializeMessage(reply)
	go self.physicalTransport.SendMessage(msg.FromPeer, serialized)
}

func (self*ReliableTransport) processSend(msg ReliableSend) {
	self.sendAck(msg)

	self.lock.Lock()
	defer self.lock.Unlock()
	_, alreadyReceived := self.messagesReceived[msg.MsgId]
	if !alreadyReceived {
		self.outputChannel <- msg.Contents
		self.messagesReceived[msg.MsgId] = time.Now().Add(self.config.RememberMessageReceivedDuration)
	}
}

func (self*ReliableTransport) processReply(msg ReliableReply) {
	self.lock.Lock()
	defer self.lock.Unlock()
	state, exists := self.messagesSent[msg.MsgId]
	if !exists {
        fmt.Fprintf(os.Stderr, "Received ack for unknown sent message %v\n", msg.MsgId)
        return
	}
	state.receivedAck = true
	self.messagesSent[msg.MsgId] = state

	// Test
	if !self.messagesSent[msg.MsgId].receivedAck {
		panic("Test error")
	}
}

func (self*ReliableTransport) processPhysicalMessage(physicalMsg []byte) {
    msg, deserialize_error := self.serializer.UnserializeMessage(physicalMsg)
    if deserialize_error != nil {
        fmt.Fprintf(os.Stderr, "Deserialization error %v\n", deserialize_error)
        return
    }
    msg_type := reflect.TypeOf(msg) 

	if msg_type == reflect.TypeOf(ReliableSend{}) {
        send := msg.(ReliableSend)
        self.processSend(send)
    } else if msg_type == reflect.TypeOf(ReliableReply{}) {
    	reply := msg.(ReliableReply)
    	self.processReply(reply)
    } else {
    	panic("Internal error: unknown message type")
    }
}

func (self*ReliableTransport) newMsgId() reliableId {
	return (reliableId)(uuid.NewV4())
}

func (self*ReliableTransport) SendMessage(toPeer cipher.PubKey, contents []byte) error {
	msgId := self.newMsgId() 
	sendMsg := ReliableSend{msgId, self.config.MyPeerId, contents}
	sendSerialized := self.serializer.SerializeMessage(sendMsg)
	state := messageSentState{toPeer, 
							  sendSerialized, 
							  time.Now().Add(self.config.RetransmitDuration),
							  false,}
	send_error := self.physicalTransport.SendMessage(toPeer, sendSerialized)
	if send_error == nil {
		self.lock.Lock()
		defer self.lock.Unlock()
		self.messagesSent[msgId] = state
	}
	return send_error
}

func (self*ReliableTransport) SetReceiveChannel(received chan []byte) {
	self.outputChannel = received
}

func (self*ReliableTransport) Close() error {
	for i := 0;i < 10;i++ {
		self.closing <- true
	}
	self.closeWait.Wait()
	return self.physicalTransport.Close()
}

func (self*ReliableTransport) SetCrypto(crypto transport.TransportCrypto) {
	self.physicalTransport.SetCrypto(crypto)
}

func (self*ReliableTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	return self.physicalTransport.ConnectedToPeer(peer)
}

func (self*ReliableTransport) GetConnectedPeers() []cipher.PubKey {
	return self.physicalTransport.GetConnectedPeers()
}

func (self*ReliableTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	empty := ReliableSend{}
	emptySerialized := self.serializer.SerializeMessage(empty)
	if (uint)(len(emptySerialized)) >= self.physicalTransport.GetMaximumMessageSizeToPeer(peer) {
		return 0
	}
	return self.physicalTransport.GetMaximumMessageSizeToPeer(peer) - (uint)(len(emptySerialized))
}

func (self*ReliableTransport) debug_countMapItems() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return len(self.messagesSent) + len(self.messagesReceived)
}

func (self*ReliableTransport) IsReliable() bool {
	return true
}


package transport

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
)

type StubTransport struct {
	testing          *testing.T
	maxMessageSize   uint
	messagesReceived chan []byte
	stubbedPeers     map[cipher.PubKey]*StubTransport
	lock             *sync.Mutex
	closeWait        *sync.WaitGroup
	ignoreSend       bool
	amReliable       bool
	messageBuffer    []QueuedMessage
	numMessagesSent  int32
	crypto           TransportCrypto
}

type QueuedMessage struct {
	toPeer *StubTransport
	msg    []byte
}

func NewStubTransport(testing *testing.T,
	maxMessageSize uint) *StubTransport {
	ret := &StubTransport{
		testing,
		maxMessageSize,
		nil,
		make(map[cipher.PubKey]*StubTransport),
		&sync.Mutex{},
		&sync.WaitGroup{},
		false,
		false,
		nil,
		0,
		nil,
	}
	return ret
}

func (self *StubTransport) Close() error {
	return nil
}

// Call before adding to node
func (self *StubTransport) AddStubbedPeer(key cipher.PubKey, peer *StubTransport) {
	self.stubbedPeers[key] = peer
}

func (self *StubTransport) getMessageBuffer() (retMessages []QueuedMessage) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.messageBuffer
}

func (self *StubTransport) SendMessage(toPeer cipher.PubKey, msg []byte) error {
	peer, exists := self.stubbedPeers[toPeer]
	if exists {
		msg_encd := msg
		if self.crypto != nil {
			peerKey := []byte{}
			if peer.crypto != nil {
				peerKey = peer.crypto.GetKey()
			}
			msg_encd = self.crypto.Encrypt(msg, peerKey)
		}
		if (uint)(len(msg)) > self.maxMessageSize {
			return errors.New(fmt.Sprintf("Message too large: %v > %v\n", len(msg), self.maxMessageSize))
		}
		if self.crypto != nil {
			msg = self.crypto.Decrypt(msg_encd)
		}
		if !self.ignoreSend {
			messageBuffer := self.getMessageBuffer()
			if messageBuffer == nil {
				peer.messagesReceived <- msg
				atomic.AddInt32(&self.numMessagesSent, 1)
			} else {
				self.lock.Lock()
				defer self.lock.Unlock()
				self.messageBuffer = append(self.messageBuffer, QueuedMessage{peer, msg})
			}
		}
		return nil
	}
	return errors.New("No stubbed transport for this peer")
}

func (self *StubTransport) SetIgnoreSendStatus(status bool) {
	self.ignoreSend = status
}

func (self *StubTransport) SetAmReliable(status bool) {
	self.amReliable = status
}

func (self *StubTransport) StartBuffer() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.messageBuffer = make([]QueuedMessage, 0)
}

func (self *StubTransport) consumeBuffer() (retMessages []QueuedMessage) {
	self.lock.Lock()
	defer self.lock.Unlock()
	retMessages = self.messageBuffer
	self.messageBuffer = nil
	return
}

func (self *StubTransport) StopAndConsumeBuffer(reorder bool, dropCount int) {
	messages := self.consumeBuffer()
	messages = messages[dropCount:]
	if reorder {
		for i := range messages {
			j := rand.Intn(i + 1)
			messages[i], messages[j] = messages[j], messages[i]
		}
	}
	for _, queued := range messages {
		queued.toPeer.messagesReceived <- queued.msg
		atomic.AddInt32(&self.numMessagesSent, 1)
	}
}

func (self *StubTransport) SetReceiveChannel(received chan []byte) {
	self.messagesReceived = received
}

func (self *StubTransport) SetCrypto(crypto TransportCrypto) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.crypto = crypto
}

func (self *StubTransport) GetConnectedPeers() []cipher.PubKey {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []cipher.PubKey{}
	for key, _ := range self.stubbedPeers {
		ret = append(ret, key)
	}
	return ret
}

func (self *StubTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	_, exists := self.stubbedPeers[peer]
	return exists
}

func (self *StubTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	return self.maxMessageSize
}

func (self *StubTransport) IsReliable() bool {
	return self.amReliable
}

func (self *StubTransport) CountNumMessagesSent() int {
	return (int)(atomic.LoadInt32(&self.numMessagesSent))
}

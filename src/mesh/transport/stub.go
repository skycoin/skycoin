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
	Testing          *testing.T
	MaxMessageSize   uint
	MessagesReceived chan []byte
	StubbedKey       *cipher.PubKey
	StubbedPeer      *StubTransport
	Lock             *sync.Mutex
	CloseWait        *sync.WaitGroup
	IgnoreSend       bool
	MessageBuffer    []QueuedMessage
	NumMessagesSent  int32
	Crypto           ITransportCrypto
}

type QueuedMessage struct {
	TransportToPeer *StubTransport
	messageContent  []byte
}

func NewStubTransport(testing *testing.T, maxMessageSize uint) *StubTransport {
	stub := &StubTransport{
		Testing:          testing,
		MaxMessageSize:   maxMessageSize,
		MessagesReceived: nil,
		StubbedKey:       &cipher.PubKey{},
		StubbedPeer:      &StubTransport{},
		Lock:             &sync.Mutex{},
		CloseWait:        &sync.WaitGroup{},
		IgnoreSend:       false,
		MessageBuffer:    nil,
		NumMessagesSent:  0,
		Crypto:           nil,
	}
	return stub
}

func (self *StubTransport) Close() error {
	return nil
}

// Call before adding to node
func (self *StubTransport) SetStubbedPeer(key cipher.PubKey, peer *StubTransport) {
	self.StubbedKey = &key
	self.StubbedPeer = peer
	return
}

func (self *StubTransport) getMessageBuffer() []QueuedMessage {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	return self.MessageBuffer
}

func (self *StubTransport) SendMessage(toPeer cipher.PubKey, message []byte, retChan chan error) error {
	var retErr error = nil
	if toPeer != *self.StubbedKey {
		retErr = errors.New("No such peer in stub")
		if retChan != nil {
			retChan <- retErr
		}
		return retErr
	}
	peer := self.StubbedPeer

	messageEncrypted := message
	if self.Crypto != nil {
		peerKey := []byte{}
		if peer.Crypto != nil {
			peerKey = peer.Crypto.GetKey()
		}
		messageEncrypted = self.Crypto.Encrypt(message, peerKey)
	}
	if (uint)(len(message)) > self.MaxMessageSize {
		retErr = errors.New(fmt.Sprintf("Message too large: %v > %v\n", len(message), self.MaxMessageSize))
		if retChan != nil {
			retChan <- retErr
		}
		return retErr
	}
	if self.Crypto != nil {
		message = self.Crypto.Decrypt(messageEncrypted)
	}
	if !self.IgnoreSend {
		messageBuffer := self.getMessageBuffer()
		if messageBuffer == nil {
			peer.MessagesReceived <- message
			atomic.AddInt32(&self.NumMessagesSent, 1)
		} else {
			self.Lock.Lock()
			defer self.Lock.Unlock()
			self.MessageBuffer = append(self.MessageBuffer, QueuedMessage{peer, message})
		}
	}
	if retChan != nil {
		retChan <- nil
	}
	return nil
	
	retErr = errors.New("No stubbed transport for this peer")
	if retChan != nil {
		retChan <- retErr
	}
	return retErr
}

func (self *StubTransport) SetIgnoreSendStatus(status bool) {
	self.IgnoreSend = status
	return
}

func (self *StubTransport) StartBuffer() {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	self.MessageBuffer = make([]QueuedMessage, 0)
	return
}

func (self *StubTransport) consumeBuffer() (retMessages []QueuedMessage) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	retMessages = self.MessageBuffer
	self.MessageBuffer = nil
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
		fmt.Println(queued.TransportToPeer)
		queued.TransportToPeer.MessagesReceived <- queued.messageContent
		fmt.Println(".fds;ljflsdk")
		atomic.AddInt32(&self.NumMessagesSent, 1)
	}
	return
}

func (self *StubTransport) SetReceiveChannel(received chan []byte) {
	fmt.Println("Setting receive channel:", received)
	self.MessagesReceived = received
	return
}

func (self *StubTransport) SetCrypto(crypto ITransportCrypto) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	self.Crypto = crypto
	return
}

func (self *StubTransport) GetConnectedPeer() cipher.PubKey {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	return *self.StubbedKey
}

func (self *StubTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	return peer == *self.StubbedKey

}

func (self *StubTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	return self.MaxMessageSize
}

func (self *StubTransport) CountNumMessagesSent() int {
	return (int)(atomic.LoadInt32(&self.NumMessagesSent))
}

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

func (s *StubTransport) Close() error {
	return nil
}

// Call before adding to node
func (s *StubTransport) SetStubbedPeer(key cipher.PubKey, peer *StubTransport) {
	s.StubbedKey = &key
	s.StubbedPeer = peer
}

func (s *StubTransport) getMessageBuffer() []QueuedMessage {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.MessageBuffer
}

func (s *StubTransport) SendMessage(toPeer cipher.PubKey, message []byte, retChan chan error) error {
	var retErr error = nil
	if toPeer != *s.StubbedKey {
		retErr = errors.New("No such peer in stub")
		if retChan != nil {
			retChan <- retErr
		}
		return retErr
	}
	peer := s.StubbedPeer

	messageEncrypted := message
	if s.Crypto != nil {
		peerKey := []byte{}
		if peer.Crypto != nil {
			peerKey = peer.Crypto.GetKey()
		}
		messageEncrypted = s.Crypto.Encrypt(message, peerKey)
	}
	if (uint)(len(message)) > s.MaxMessageSize {
		retErr = errors.New(fmt.Sprintf("Message too large: %v > %v\n", len(message), s.MaxMessageSize))
		if retChan != nil {
			retChan <- retErr
		}
		return retErr
	}
	if s.Crypto != nil {
		message = s.Crypto.Decrypt(messageEncrypted)
	}
	if !s.IgnoreSend {
		messageBuffer := s.getMessageBuffer()
		if messageBuffer == nil {
			peer.MessagesReceived <- message
			atomic.AddInt32(&s.NumMessagesSent, 1)
		} else {
			s.Lock.Lock()
			defer s.Lock.Unlock()
			s.MessageBuffer = append(s.MessageBuffer, QueuedMessage{peer, message})
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

func (s *StubTransport) SetIgnoreSendStatus(status bool) {
	s.IgnoreSend = status
}

func (s *StubTransport) StartBuffer() {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.MessageBuffer = make([]QueuedMessage, 0)
}

func (s *StubTransport) consumeBuffer() (retMessages []QueuedMessage) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	retMessages = s.MessageBuffer
	s.MessageBuffer = nil
	return
}

func (s *StubTransport) StopAndConsumeBuffer(reorder bool, dropCount int) {
	messages := s.consumeBuffer()
	messages = messages[dropCount:]
	if reorder {
		for i := range messages {
			j := rand.Intn(i + 1)
			messages[i], messages[j] = messages[j], messages[i]
		}
	}
	for _, queued := range messages {
		queued.TransportToPeer.MessagesReceived <- queued.messageContent
		atomic.AddInt32(&s.NumMessagesSent, 1)
	}
}

func (s *StubTransport) SetReceiveChannel(received chan []byte) {
	s.MessagesReceived = received
}

func (s *StubTransport) SetCrypto(crypto ITransportCrypto) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Crypto = crypto
}

func (s *StubTransport) GetConnectedPeer() cipher.PubKey {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return *s.StubbedKey
}

func (s *StubTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return peer == *s.StubbedKey

}

func (s *StubTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	return s.MaxMessageSize
}

func (s *StubTransport) CountNumMessagesSent() int {
	return (int)(atomic.LoadInt32(&s.NumMessagesSent))
}

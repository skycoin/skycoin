package transport

import(
	"fmt"
	"sync"
	"errors"
	"testing")

import("github.com/skycoin/skycoin/src/cipher")

type StubTransport struct {
	testing *testing.T
	maxMessageSize uint
	messagesReceived chan []byte
	stubbedPeers map[cipher.PubKey]*StubTransport
    lock *sync.Mutex
    closeWait *sync.WaitGroup
    ignoreSend bool
}

func NewStubTransport(testing *testing.T, 
					  maxMessageSize uint) (*StubTransport) {
	ret := &StubTransport{
		testing,
		maxMessageSize,
		nil,
		make(map[cipher.PubKey]*StubTransport),
		&sync.Mutex{},
		&sync.WaitGroup{},
		false,
	}
	return ret
}
func (self*StubTransport) Close() error {
	return nil
}
// Call before adding to node
func (self*StubTransport) AddStubbedPeer(key cipher.PubKey, peer *StubTransport) {
	self.stubbedPeers[key] = peer
}
func (self*StubTransport) SendMessage(toPeer cipher.PubKey, msg []byte) error {
	if (uint)(len(msg)) > self.maxMessageSize {
		return errors.New(fmt.Sprintf("Message too large: %v > %v\n", len(msg), self.maxMessageSize))
	}
	peer, exists := self.stubbedPeers[toPeer]
	if exists {
		if !self.ignoreSend {
			peer.messagesReceived <- msg
		}
		return nil
	}
	return errors.New("No stubbed transport for this peer")
}
func (self*StubTransport) SetIgnoreSendStatus(status bool) {
	self.ignoreSend = status
}
func (self*StubTransport) SetReceiveChannel(received chan []byte) {
	self.messagesReceived = received
}
func (self*StubTransport) SetCrypto(crypto TransportCrypto) {
	panic("crypto unsupported")
}
func (self*StubTransport) GetConnectedPeers() []cipher.PubKey {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []cipher.PubKey{}
	for key, _ := range(self.stubbedPeers) {
		ret = append(ret, key)
	}
	return ret
}
func (self*StubTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	_, exists := self.stubbedPeers[peer]
	return exists
}
func (self*StubTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	return self.maxMessageSize
}
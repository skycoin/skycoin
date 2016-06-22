package mesh

import(
	"sync"
	"testing")

import("github.com/skycoin/skycoin/src/cipher"
	   "github.com/stretchr/testify/assert")

type StubTransport struct {
	testing *testing.T
	maxMessageSize uint
	messagesSent chan TransportMessage
	MessagesReceived chan TransportMessage
	connectedTo map[cipher.PubKey]bool
    lock *sync.Mutex
}

func NewStubTransport(testing *testing.T, maxMessageSize uint, sentMessages chan TransportMessage) (*StubTransport) {
	ret := &StubTransport{
		testing,
		maxMessageSize,
		sentMessages, 	// MessagesSent
		nil,
		make(map[cipher.PubKey]bool),
		&sync.Mutex{},
	}
	return ret
}
func (self*StubTransport) Close() error {
	return nil
}
func (self*StubTransport) SendMessage(msg TransportMessage) error {
	self.messagesSent <- msg
	return nil
}
func (self*StubTransport) SetReceiveChannel(received chan TransportMessage) {
	self.MessagesReceived = received
}
func (self*StubTransport) SetCrypto(crypto TransportCrypto) {
}
func (self*StubTransport) GetConnectedPeers() []cipher.PubKey {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []cipher.PubKey{}
	for key, _ := range(self.connectedTo) {
		ret = append(ret, key)
	}
	return ret
}
func (self*StubTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	_, exists := self.connectedTo[peer]
	return exists
}
func (self*StubTransport) ConnectToPeer(peer cipher.PubKey, connectInfo string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	assert.Equal(self.testing, "foo", connectInfo)
	self.connectedTo[peer] = true
	return nil
}
func (self*StubTransport) DisconnectFromPeer(peer cipher.PubKey) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.connectedTo, peer)
}
func (self*StubTransport) GetTransportConnectInfo() string {
	return "foo"
}
func (self*StubTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	return self.maxMessageSize
}
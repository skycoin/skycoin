package mesh

import(
	"testing")

import("github.com/skycoin/skycoin/src/cipher"
	   "github.com/stretchr/testify/assert")

type StubTransport struct {
	testing *testing.T
	maxMessageSize uint
	messagesSent chan TransportMessage
	MessagesReceived chan TransportMessage
}

func NewStubTransport(testing *testing.T, maxMessageSize uint, sentMessages chan TransportMessage) (*StubTransport) {
	ret := &StubTransport{
		testing,
		maxMessageSize,
		sentMessages, 	// MessagesSent
		nil,
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
func (self*StubTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	return true
}
func (self*StubTransport) ConnectToPeer(peer cipher.PubKey, connectInfo string) error {
	assert.Equal(self.testing, "foo", connectInfo)
	return nil
}
func (self*StubTransport) DisconnectFromPeer(peer cipher.PubKey) {
}
func (self*StubTransport) GetTransportConnectInfo() string {
	return "foo"
}
func (self*StubTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	return self.maxMessageSize
}
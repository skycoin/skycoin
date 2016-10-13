package physical

import (
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

var staticTestConfig UDPConfig = UDPConfig{
	TransportConfig: TransportConfig{
		SendChannelLength: 8,
	},
	DatagramLength:  512,
	LocalAddress:    "",
	NumListenPorts:  5,
	ListenPortMin:   10300, // If 0, STUN is used
	ExternalAddress: "127.0.0.1",
	StunEndpoints:   nil, // STUN servers to try for NAT traversal
}

func TestBindStaticPorts(t *testing.T) {
	transport, err := NewUDPTransport(staticTestConfig)
	assert.Nil(t, err)
	assert.NotNil(t, transport)
	defer transport.Close()
}

func TestClose(t *testing.T) {
	transport, err := NewUDPTransport(staticTestConfig)
	assert.Nil(t, err)
	assert.NotNil(t, transport)
	defer transport.Close()
	time.Sleep(3 * time.Second)
}

func TestBindSTUNPorts(t *testing.T) {
	config := UDPConfig{
		TransportConfig: TransportConfig{
			SendChannelLength: 8,
		},
		DatagramLength:  512,
		LocalAddress:    "",
		NumListenPorts:  5,
		ListenPortMin:   0, // If 0, STUN is used
		ExternalAddress: "127.0.0.1",
		StunEndpoints:   []string{"stun1.voiceeclipse.net:3478"}, // STUN servers to try for NAT traversal
	}
	transport, err := NewUDPTransport(config)
	assert.Nil(t, err)
	assert.NotNil(t, transport)
	defer transport.Close()
}

func SetupAB(encrypt bool, t *testing.T) (
	*UDPTransport, cipher.PubKey,
	*UDPTransport, cipher.PubKey) {
	transportA, err := NewUDPTransport(staticTestConfig)
	assert.Nil(t, err)
	assert.NotNil(t, transportA)

	configB := staticTestConfig
	configB.ListenPortMin = 10400

	transportB, err := NewUDPTransport(configB)
	assert.Nil(t, err)
	assert.NotNil(t, transportB)

	tc := &TestCryptoStruct{}
	transportA.SetCrypto(tc)
	transportB.SetCrypto(tc)

	testKeyA := cipher.NewPubKey([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	testKeyB := cipher.NewPubKey([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	testKeyC := cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Nil(t, transportA.ConnectToPeer(testKeyB, transportB.GetTransportConnectInfo()))
	assert.True(t, transportA.ConnectedToPeer(testKeyB))
	assert.False(t, transportA.ConnectedToPeer(testKeyC))
	assert.Nil(t, transportB.ConnectToPeer(testKeyA, transportA.GetTransportConnectInfo()))
	assert.True(t, transportB.ConnectedToPeer(testKeyA))
	assert.False(t, transportB.ConnectedToPeer(testKeyC))

	return transportA, testKeyA, transportB, testKeyB
}

func TestSendDatagram(t *testing.T) {
	transportA, keyA, transportB, keyB := SetupAB(false, t)
	defer transportA.Close()
	defer transportB.Close()

	testKeyC := cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Zero(t, transportA.GetMaximumMessageSizeToPeer(keyA))
	assert.NotZero(t, transportA.GetMaximumMessageSizeToPeer(keyB))
	assert.NotZero(t, transportB.GetMaximumMessageSizeToPeer(keyA))
	assert.Zero(t, transportB.GetMaximumMessageSizeToPeer(keyB))
	assert.Zero(t, transportA.GetMaximumMessageSizeToPeer(testKeyC))
	assert.Zero(t, transportB.GetMaximumMessageSizeToPeer(testKeyC))

	sendBytesA := []byte{66, 44, 33, 2, 123, 100, 22}
	sendBytesB := []byte{23, 33, 12, 88, 43, 120}

	assert.Nil(t, transportA.SendMessage(keyB, sendBytesA))
	assert.Nil(t, transportB.SendMessage(keyA, sendBytesB))

	chanA := make(chan []byte, 10)
	chanB := make(chan []byte, 10)

	transportA.SetReceiveChannel(chanA)
	transportB.SetReceiveChannel(chanB)

	gotA := false
	gotB := false

	for !gotA || !gotB {
		select {
		case msg_a := <-chanA:
			{
				assert.Equal(t, sendBytesB, msg_a)
				gotA = true
				break
			}
		case msg_b := <-chanB:
			{
				assert.Equal(t, sendBytesA, msg_b)
				gotB = true
				break
			}
		case <-time.After(5 * time.Second):
			panic("Test timed out")
		}
	}
}

type TestCryptoStruct struct {
}

func (self *TestCryptoStruct) GetKey() []byte {
	return []byte{44, 23}
}

func (self *TestCryptoStruct) Encrypt(data []byte, key []byte) []byte {
	if len(key) != 2 || key[0] != 44 || key[1] != 23 {
		panic("Wrong Key")
	}
	ret := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		ret[i] = data[i] + 1
	}
	return ret
}

func (self *TestCryptoStruct) Decrypt(data []byte) []byte {
	ret := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		ret[i] = data[i] - 1
	}
	return ret
}

func TestCrypto(t *testing.T) {
	transportA, _, transportB, keyB := SetupAB(true, t)
	defer transportA.Close()
	defer transportB.Close()

	send_bytes := []byte{66, 44, 33, 2, 123, 100, 22}

	assert.Nil(t, transportA.SendMessage(keyB, send_bytes))

	chanB := make(chan []byte, 10)
	transportB.SetReceiveChannel(chanB)

	select {
	case msg_b := <-chanB:
		{
			assert.Equal(t, send_bytes, msg_b)
		}
	case <-time.After(5 * time.Second):
		panic("Test timed out")
	}
}

func TestDisconnect(t *testing.T) {
	transportA, _, transportB, keyB := SetupAB(false, t)
	defer transportA.Close()
	defer transportB.Close()

	assert.Equal(t, []cipher.PubKey{keyB}, transportA.GetConnectedPeers())
	transportA.DisconnectFromPeer(keyB)
	assert.False(t, transportA.ConnectedToPeer(keyB))
	assert.Zero(t, transportA.GetMaximumMessageSizeToPeer(keyB))
	assert.Equal(t, []cipher.PubKey{}, transportA.GetConnectedPeers())
}

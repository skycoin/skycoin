package mesh

import(
	"time"
	"testing")

import(
    "github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert")

var staticTestConfig UDPConfig = UDPConfig {
		TransportConfig {
			8, // SendChannelLength uint32
		},
		512, // DatagramLength	uint64
		"", // LocalAddress string 	// "" for default

		5, // NumListenPorts uint16
		10300, // ListenPortMin uint16		// If 0, STUN is used
		"127.0.0.1", // ExternalAddress
		nil, // StunEndpoints []string		// STUN servers to try for NAT traversal
	}

func TestBindStaticPorts(t *testing.T) {
	transport, error := NewUDPTransport(staticTestConfig)
	assert.Nil(t, error)
	assert.NotNil(t, transport)
	defer transport.Close()
}

func TestClose(t *testing.T) {
	transport, error := NewUDPTransport(staticTestConfig)
	assert.Nil(t, error)
	assert.NotNil(t, transport)
	defer transport.Close()
	time.Sleep(3 * time.Second)
}

func TestBindSTUNPorts(t *testing.T) {
	config := UDPConfig {
		TransportConfig {
			8, // SendChannelLength uint32
		},
		512, // DatagramLength	uint64
		"", // LocalAddress string 	// "" for default

		5, // NumListenPorts uint16
		0, // ListenPortMin uint16		// If 0, STUN is used
		"127.0.0.1", // ExternalAddress
		[]string{"stun1.voiceeclipse.net:3478"}, // StunEndpoints []string		// STUN servers to try for NAT traversal
	}
	transport, error := NewUDPTransport(config)
	assert.Nil(t, error)
	assert.NotNil(t, transport)
	defer transport.Close()
}

func SetupAB(t *testing.T) (
				*UDPTransport, cipher.PubKey,
				*UDPTransport, cipher.PubKey) {
	transport_a, error := NewUDPTransport(staticTestConfig)
	assert.Nil(t, error)
	assert.NotNil(t, transport_a)

	config_b := staticTestConfig
	config_b.ListenPortMin = 10400

	transport_b, error := NewUDPTransport(config_b)
	assert.Nil(t, error)
	assert.NotNil(t, transport_b)

	test_key_a := cipher.NewPubKey([]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	test_key_b := cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	test_key_c := cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	assert.Nil(t, transport_a.ConnectToPeer(test_key_b, transport_b.GetTransportConnectInfo()))
	assert.True(t, transport_a.ConnectedToPeer(test_key_b))
	assert.False(t, transport_a.ConnectedToPeer(test_key_c))
	assert.Nil(t, transport_b.ConnectToPeer(test_key_a, transport_a.GetTransportConnectInfo()))
	assert.True(t, transport_b.ConnectedToPeer(test_key_a))
	assert.False(t, transport_b.ConnectedToPeer(test_key_c))

	return transport_a, test_key_a, transport_b, test_key_b
}

func TestSendDatagram(t *testing.T) {
	transport_a , key_a, transport_b, key_b := SetupAB(t)
	defer transport_a.Close()
	defer transport_b.Close()

	test_key_c := cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	assert.Zero(t, transport_a.GetMaximumMessageSizeToPeer(key_a))
	assert.NotZero(t, transport_a.GetMaximumMessageSizeToPeer(key_b))
	assert.NotZero(t, transport_b.GetMaximumMessageSizeToPeer(key_a))
	assert.Zero(t, transport_b.GetMaximumMessageSizeToPeer(key_b))
	assert.Zero(t, transport_a.GetMaximumMessageSizeToPeer(test_key_c))
	assert.Zero(t, transport_b.GetMaximumMessageSizeToPeer(test_key_c))

	send_bytes_a := []byte{66,44,33,2,123,100,22}
	send_bytes_b := []byte{23,33,12,88,43,120}

	assert.Nil(t, transport_a.SendMessage(TransportMessage{key_b, send_bytes_a}))
	assert.Nil(t, transport_b.SendMessage(TransportMessage{key_a, send_bytes_b}))

	chan_a := make(chan TransportMessage, 10)
	chan_b := make(chan TransportMessage, 10)

	transport_a.SetReceiveChannel(chan_a)
	transport_b.SetReceiveChannel(chan_b)

	got_a := false
	got_b := false

	for !got_a || !got_b {
		select {
			case msg_a := <- chan_a: {
				assert.Equal(t, key_a, msg_a.DestPeer)
				assert.Equal(t, send_bytes_b, msg_a.Contents)
				got_a = true
				break
			}
			case msg_b := <- chan_b: {
				assert.Equal(t, key_b, msg_b.DestPeer)
				assert.Equal(t, send_bytes_a, msg_b.Contents)
				got_b = true
				break
			}
			case <-time.After(5*time.Second):
				panic("Test timed out")
		}
	}
}

type TestCryptoStruct struct {
}

func (self*TestCryptoStruct) Encrypt(data[]byte)[]byte {
	ret := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		ret[i] = data[i] + 1
	}
	return ret
}

func (self*TestCryptoStruct) Decrypt(data[]byte)[]byte {
	ret := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		ret[i] = data[i] - 1
	}
	return ret
}

func TestCrypto(t *testing.T) {
	transport_a, _, transport_b, key_b := SetupAB(t)
	defer transport_a.Close()
	defer transport_b.Close()

	send_bytes := []byte{66,44,33,2,123,100,22}

	tc := &TestCryptoStruct{}
	transport_a.SetCrypto(tc)
	transport_b.SetCrypto(tc)

	assert.Nil(t, transport_a.SendMessage(TransportMessage{key_b, send_bytes}))

	chan_b := make(chan TransportMessage, 10)
	transport_b.SetReceiveChannel(chan_b)

	select {
		case msg_b := <- chan_b: {
			assert.Equal(t, key_b, msg_b.DestPeer)
			assert.Equal(t, send_bytes, msg_b.Contents)
		}
		case <-time.After(5*time.Second):
			panic("Test timed out")
	}
}

func TestDisconnect(t *testing.T) {
	transport_a, _, transport_b, key_b := SetupAB(t)
	defer transport_a.Close()
	defer transport_b.Close()

	transport_a.DisconnectFromPeer(key_b)
	assert.False(t, transport_a.ConnectedToPeer(key_b))
	assert.Zero(t, transport_a.GetMaximumMessageSizeToPeer(key_b))
}



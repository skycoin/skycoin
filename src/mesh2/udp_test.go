package mesh

import(
	"time"
	"testing")

import(
	"github.com/stretchr/testify/assert")



func TestBindStaticPorts(t *testing.T) {
	config := UDPConfig {
		TransportConfig {
			8, // SendChannelLength uint32
			8, // ReceiveChannelLength uint32
		},
		512, // DatagramLength	uint64
		"", // LocalAddress string 	// "" for default

		5, // NumListenPorts uint16
		10300, // ListenPortMin uint16		// If 0, STUN is used
		nil, // StunEndpoints []string		// STUN servers to try for NAT traversal
	}
	transport, error := NewUDPTransport(config)
	assert.Nil(t, error)
	assert.NotNil(t, transport)
	defer transport.Close()
    time.Sleep(1 * time.Second)
}

func TestBindSTUNPorts(t *testing.T) {
	config := UDPConfig {
		TransportConfig {
			8, // SendChannelLength uint32
			8, // ReceiveChannelLength uint32
		},
		512, // DatagramLength	uint64
		"", // LocalAddress string 	// "" for default

		5, // NumListenPorts uint16
		0, // ListenPortMin uint16		// If 0, STUN is used
		[]string{"stun1.voiceeclipse.net:3478"}, // StunEndpoints []string		// STUN servers to try for NAT traversal
	}
	transport, error := NewUDPTransport(config)
	assert.Nil(t, error)
	assert.NotNil(t, transport)
	defer transport.Close()
    time.Sleep(1 * time.Second)
}



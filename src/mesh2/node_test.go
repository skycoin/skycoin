package mesh

import(
	"time"
	"testing")

import(
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert")

// TODO: Test GetTransports(), RemoveTransport(), GetConnectedPeers()

func SetupNode(t *testing.T, sentMessages chan TransportMessage) (*Node, *StubTransport) {
	transport := NewStubTransport(t, 512, sentMessages)
	node, error := NewNode(NodeConfig{
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
		})
	assert.Nil(t, error)
	node.AddTransport(transport)
	return node, transport
}

func SetupNodes(n uint, t *testing.T) (nodes []*Node, to_close chan TransportMessage) {
	nodes = make([]*Node, n)
	transports := make([]*StubTransport, n)
	to_close = make(chan TransportMessage, 20)
	sentMessages := make(chan TransportMessage, 20)
	for i := (uint)(0); i < n; i++ {
		nodes[i], transports[i] = SetupNode(t, sentMessages)
		nodes[i].AddTransport(transports[i])
	}
	go func() {
		for {
			msg, more := <-sentMessages
			if more {
				for i := (uint)(0); i < n; i++ {
					transports[i].MessagesReceived <- msg
				}
			} else {
				return
			}
		}
	}()
	return nodes, sentMessages
}

func TestEstablishRoute(t *testing.T) {
	nodes, to_close := SetupNodes(2, t)
	defer close(to_close)

	test_key_a := cipher.NewPubKey([]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	test_key_b := cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

	route := uuid.NewV4()

	established_ch := make(chan bool, 1)
	nodes[0].SetRouteStatusCallback(func(routeId uuid.UUID, ready bool, establishedToHopIdx int) {
		if routeId == route && ready {
			established_ch <- true
		}
	})
	assert.Nil(t, nodes[0].AddRoute(route, []cipher.PubKey{test_key_a, test_key_b}))

	select {
		case <- established_ch: {
			break
		}
		case <-time.After(5*time.Second):
			panic("Test timed out")
	}
}

// Packet loss test
// Multiple transport test
// UDP Test
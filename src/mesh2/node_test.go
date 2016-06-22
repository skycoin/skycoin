package mesh

import(
	"testing")

import(
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert")

func TestManageTransports(t *testing.T) {
	transport_a := NewStubTransport(t, 512, nil)
	transport_b := NewStubTransport(t, 512, nil)
	node, error := NewNode(NodeConfig{
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
		})
	assert.Nil(t, error)
	assert.Equal(t, []Transport{}, node.GetTransports())
	node.AddTransport(transport_a)
	assert.Equal(t, []Transport{transport_a}, node.GetTransports())
	node.AddTransport(transport_b)
	assert.Equal(t, []Transport{transport_a, transport_b}, node.GetTransports())
	node.RemoveTransport(transport_a)
	assert.Equal(t, []Transport{transport_b}, node.GetTransports())
	node.RemoveTransport(transport_b)
	assert.Equal(t, []Transport{}, node.GetTransports())
}

func TestConnectedPeers(t *testing.T) {
	transport_a := NewStubTransport(t, 512, nil)
	transport_b := NewStubTransport(t, 512, nil)
	node, error := NewNode(NodeConfig{
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
		})
	assert.Nil(t, error)
	peer_a := cipher.NewPubKey([]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	peer_b := cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	peer_c := cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	assert.Nil(t, transport_a.ConnectToPeer(peer_a, "foo"))
	assert.Nil(t, transport_a.ConnectToPeer(peer_b, "foo"))
	assert.Nil(t, transport_b.ConnectToPeer(peer_c, "foo"))

	assert.Equal(t, []cipher.PubKey{}, node.GetConnectedPeers())
	node.AddTransport(transport_a)
	assert.Equal(t, []cipher.PubKey{peer_a, peer_b}, node.GetConnectedPeers())
	node.AddTransport(transport_b)
	assert.Equal(t, []cipher.PubKey{peer_a, peer_b, peer_c}, node.GetConnectedPeers())
	assert.True(t, transport_a.ConnectedToPeer(peer_a))
	transport_a.DisconnectFromPeer(peer_a)
	assert.False(t, transport_a.ConnectedToPeer(peer_a))
	assert.Equal(t, []cipher.PubKey{peer_b, peer_c}, node.GetConnectedPeers())
	node.RemoveTransport(transport_a)
	assert.Equal(t, []cipher.PubKey{peer_c}, node.GetConnectedPeers())
	node.RemoveTransport(transport_b)
	assert.Equal(t, []cipher.PubKey{}, node.GetConnectedPeers())
}

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

func TestSendDirect(t *testing.T) {
	nodes, to_close := SetupNodes(3, t)
	defer close(to_close)
	defer func() {
		for _, node := range(nodes) {
			node.Close()
		}
	}()

	test_key_b := cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	test_key_c := cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

	route := uuid.NewV4()
	assert.Nil(t, nodes[0].AddRoute(route, test_key_b))
	assert.Nil(t, nodes[0].ExtendRoute(route, test_key_c))
}

/*
func TestEstablishRoute(t *testing.T) {
	nodes, to_close := SetupNodes(3, t)
	defer close(to_close)
	defer func() {
		for _, node := range(nodes) {
			node.Close()
		}
	}()

	test_key_b := cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	test_key_c := cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

	route := uuid.NewV4()
	assert.Nil(t, nodes[0].AddRoute(route, test_key_b))
	assert.Nil(t, nodes[0].ExtendRoute(route, test_key_c))
}
*/
/*
func TestSendThruRoute(t *testing.T) {
}
*/

/*
func TestSendReply(t *testing.T) {

}
*/

// Packet loss test
// Multiple transport test
// UDP Test
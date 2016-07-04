package mesh

import(
	"time"
	"testing"
	"sort")

import(
    "github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert")

func sortPubKeys(pubKeys []cipher.PubKey) ([]cipher.PubKey) {
	var ret cipher.PubKeySlice = pubKeys
	sort.Sort(ret)
	return ret
}

func TestManageTransports(t *testing.T) {
	transport_a := NewStubTransport(t, 512)
	transport_b := NewStubTransport(t, 512)
	node, error := NewNode(NodeConfig{
			cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}),
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
			time.Minute,
			10 * time.Second,
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
	transport_a := NewStubTransport(t, 512)
	transport_b := NewStubTransport(t, 512)
	node, error := NewNode(NodeConfig{
			cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}),
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
			time.Minute,
			10 * time.Second,
		})
	assert.Nil(t, error)
	peer_a := cipher.NewPubKey([]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	peer_b := cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	peer_c := cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	transport_a.AddStubbedPeer(peer_a, nil)
	transport_a.AddStubbedPeer(peer_b, nil)
	transport_b.AddStubbedPeer(peer_c, nil)

	assert.False(t, node.ConnectedToPeer(peer_a))
	assert.False(t, node.ConnectedToPeer(peer_b))
	assert.False(t, node.ConnectedToPeer(peer_c))
	assert.Equal(t, []cipher.PubKey{}, sortPubKeys(node.GetConnectedPeers()))
	node.AddTransport(transport_a)
	assert.Equal(t, []cipher.PubKey{peer_a, peer_b}, sortPubKeys(node.GetConnectedPeers()))
	assert.True(t, node.ConnectedToPeer(peer_a))
	assert.True(t, node.ConnectedToPeer(peer_b))
	assert.False(t, node.ConnectedToPeer(peer_c))

	node.AddTransport(transport_b)
	assert.Equal(t, []cipher.PubKey{peer_a, peer_b, peer_c}, sortPubKeys(node.GetConnectedPeers()))
	assert.True(t, node.ConnectedToPeer(peer_a))
	assert.True(t, node.ConnectedToPeer(peer_b))
	assert.True(t, node.ConnectedToPeer(peer_c))
	assert.True(t, transport_a.ConnectedToPeer(peer_a))
	node.RemoveTransport(transport_a)
	assert.False(t, node.ConnectedToPeer(peer_a))
	assert.False(t, node.ConnectedToPeer(peer_b))
	assert.True(t, node.ConnectedToPeer(peer_c))

	assert.Equal(t, []cipher.PubKey{peer_c}, sortPubKeys(node.GetConnectedPeers()))
	node.RemoveTransport(transport_b)
	assert.Equal(t, []cipher.PubKey{}, sortPubKeys(node.GetConnectedPeers()))
	assert.False(t, node.ConnectedToPeer(peer_a))
	assert.False(t, node.ConnectedToPeer(peer_b))
	assert.False(t, node.ConnectedToPeer(peer_c))
}

func SetupNode(t *testing.T) (*Node, *StubTransport) {
	transport := NewStubTransport(t, 512)
	newPubKey, _ := cipher.GenerateKeyPair()
	node, error := NewNode(NodeConfig{
			newPubKey,
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
			time.Minute,
			10 * time.Second,
		})
	assert.Nil(t, error)
	node.AddTransport(transport)
	return node, transport
}

// Nodes each have one transport
// All nodes receive all other nodes' messages, but stub transport filters
func SetupNodes(n uint, connections [][]int, t *testing.T) (nodes []*Node, to_close chan []byte) {
	nodes = make([]*Node, n)
	transports := make([]*StubTransport, n)
	to_close = make(chan []byte, 20)
	sentMessages := make(chan []byte, 20)
	for i := (uint)(0); i < n; i++ {
		nodes[i], transports[i] = SetupNode(t)
		nodes[i].AddTransport(transports[i])
	}
	for i := (uint)(0); i < n; i++ {
		transport_from := transports[i]
		for j := (uint)(0); j < n; j++ {
			transport_to := transports[j]
			if connections[i][j] != 0 {
				transport_from.AddStubbedPeer(nodes[j].GetConfig().PubKey, transport_to)
			}
		}
	}
	return nodes, sentMessages
}

func sendDirect(t *testing.T, reliable bool) {
	connections  := [][]int{
		[]int{1,1,},
		[]int{1,1,},
	}
	nodes, to_close := SetupNodes(2, connections, t)
	defer close(to_close)
	defer func() {
		for _, node := range(nodes) {
			node.Close()
		}
	}()

	contents := []byte{4,66,7,44,33}

	received := make(chan MeshMessage, 10)
	nodes[1].SetReceiveChannel(received)

	test_key_b := nodes[1].GetConfig().PubKey
	send_err, route_id := nodes[0].SendMessageToPeer(test_key_b, contents, reliable, time.Second)
	assert.Nil(t, send_err)
	assert.Zero(t, route_id)

	select {
		case recvd := <- received: {
			assert.Zero(t, recvd.RouteId)
			assert.Equal(t, contents, recvd.Contents)
		}
		case <-time.After(5*time.Second):
			panic("Test timed out")
	}
}

// Route expiry test
// Packet loss test
// Multiple transport test

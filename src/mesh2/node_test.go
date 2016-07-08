package mesh

import(
	"time"
	"testing"
	"sort")

import(
	"github.com/skycoin/skycoin/src/mesh2/transport"
    "github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert")

func sortPubKeys(pubKeys []cipher.PubKey) ([]cipher.PubKey) {
	var ret cipher.PubKeySlice = pubKeys
	sort.Sort(ret)
	return ret
}

func TestManageTransports(t *testing.T) {
	transport_a := transport.NewStubTransport(t, 512)
	transport_b := transport.NewStubTransport(t, 512)
	test_key_a := cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	node, error := NewNode(NodeConfig{
			test_key_a,
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
			time.Minute,
			10 * time.Second,
			time.Second,
			time.Second,
			2*time.Second,
			100, // Transport message channel length
		})
	assert.Nil(t, error)
	assert.Equal(t, []transport.Transport{}, node.GetTransports())
	node.AddTransport(transport_a)
	assert.Equal(t, []transport.Transport{transport_a}, node.GetTransports())
	node.AddTransport(transport_b)
	assert.Equal(t, []transport.Transport{transport_a, transport_b}, node.GetTransports())
	node.RemoveTransport(transport_a)
	assert.Equal(t, []transport.Transport{transport_b}, node.GetTransports())
	node.RemoveTransport(transport_b)
	assert.Equal(t, []transport.Transport{}, node.GetTransports())
}

func TestConnectedPeers(t *testing.T) {
	transport_a := transport.NewStubTransport(t, 512)
	transport_b := transport.NewStubTransport(t, 512)
	test_key_a := cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	node, error := NewNode(NodeConfig{
			test_key_a,
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
			time.Minute,
			10 * time.Second,
			time.Second,
			time.Second,
			2*time.Second,
			100, // Transport message channel length
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

func SetupNode(t *testing.T,
			   newPubKey cipher.PubKey) (node *Node, 
							  unreliableTransport *transport.StubTransport,
							  reliableTransport *transport.StubTransport) {
	unreliableTransport = transport.NewStubTransport(t, 512)
	reliableTransport = transport.NewStubTransport(t, 512)
	var error error
	node, error = NewNode(NodeConfig{
			newPubKey,
			[32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, },
			time.Minute,
			10 * time.Second,
			time.Second,
			time.Second,
			2*time.Second,
			100, // Transport message channel length
		})
	assert.Nil(t, error)
	node.AddTransport(unreliableTransport)
	node.AddTransport(reliableTransport)
	return
}

// Nodes each have one transport
// All nodes receive all other nodes' messages, but stub transport filters
func SetupNodes(n uint, connections [][]int, t *testing.T) (nodes []*Node, to_close chan []byte, 
															unreliableTransports []*transport.StubTransport, 
															reliableTransports []*transport.StubTransport) {
	nodes = make([]*Node, n)
	unreliableTransports = make([]*transport.StubTransport, n)
	reliableTransports = make([]*transport.StubTransport, n)
	to_close = make(chan []byte, 20)
	sentMessages := make(chan []byte, 20)
	for i := (uint)(0); i < n; i++ {
		pubKey := cipher.PubKey{}
		pubKey[0] = (byte)(i + 1)
		nodes[i], unreliableTransports[i], reliableTransports[i] = SetupNode(t, pubKey)
		unreliableTransports[i].SetAmReliable(false)
		reliableTransports[i].SetAmReliable(true)
	}
	for i := (uint)(0); i < n; i++ {
		transport_from := unreliableTransports[i]
		for j := (uint)(0); j < n; j++ {
			transport_to := unreliableTransports[j]
			if connections[i][j] != 0 {
				transport_from.AddStubbedPeer(nodes[j].GetConfig().PubKey, transport_to)
			}
		}
	}
	for i := (uint)(0); i < n; i++ {
		transport_from := reliableTransports[i]
		for j := (uint)(0); j < n; j++ {
			transport_to := reliableTransports[j]
			if connections[i][j] != 0 {
				transport_from.AddStubbedPeer(nodes[j].GetConfig().PubKey, transport_to)
			}
		}
	}	
	return nodes, sentMessages, unreliableTransports, reliableTransports
}

func TestDeleteRoute(t *testing.T) {
	// todo
}

func sendTest(t *testing.T, nPeers int, reliable bool, dropFirst bool, reorder bool, sendBack bool, contents []byte) {
	if nPeers < 2 {
		panic("Fewer than 2 peers doesn't make sense")
	}

	allConnections  := make([][]int, 0)
	for from_idx := 0; from_idx < nPeers; from_idx++ {
		toConnections := make([]int, 0)
		for i := 0; i < nPeers; i++ {
			toConnections = append(toConnections, 0)
		}

		for to_idx := from_idx - 1; to_idx <= from_idx + 1; to_idx++ {
			if to_idx >= 0 && to_idx != from_idx && to_idx < nPeers {
				toConnections[to_idx] = 1
			}
		}
		allConnections = append(allConnections, toConnections)
	}
	nodes, to_close, unreliableTransport, reliableTransport := SetupNodes((uint)(nPeers), allConnections, t)
	defer close(to_close)
	defer func() {
		for _, node := range(nodes) {
			node.Close()
		}
	}()

	if dropFirst {
		for _, unreliableTransport := range(unreliableTransport) {
			unreliableTransport.SetAmReliable(false)
		}
	}

	received := make(chan MeshMessage, 10)
	nodes[nPeers-1].SetReceiveChannel(received)

	terminating_id := nodes[nPeers-1].GetConfig().PubKey

	addedRouteId := RouteId{}
	addedRouteId[0] = 22
	assert.Nil(t, nodes[0].AddRoute(addedRouteId, nodes[1].GetConfig().PubKey))

	for extendIdx := 2; extendIdx < nPeers; extendIdx++ {
		assert.Nil(t, nodes[0].ExtendRoute(addedRouteId, nodes[extendIdx].GetConfig().PubKey, time.Second))
	}

	var replyTo ReplyTo

	for dropFirstIdx := 0; dropFirstIdx<2; dropFirstIdx++ {
		shouldReceive := true
		if dropFirst && dropFirstIdx == 0 {
			shouldReceive = false
		}

		for _, unreliableTransport := range(unreliableTransport) {
			unreliableTransport.StartBuffer()
			unreliableTransport.SetIgnoreSendStatus(!shouldReceive)
		}
		for _, reliableTransport := range(reliableTransport) {
			reliableTransport.StartBuffer()
		}

		send_err, route_id := nodes[0].SendMessageToPeer(terminating_id, contents, reliable)
		assert.Nil(t, send_err)
		assert.Equal(t, addedRouteId, route_id)

		for _, unreliableTransport := range(unreliableTransport) {
			unreliableTransport.StopAndConsumeBuffer(reorder)
		}

		for _, reliableTransport := range(reliableTransport) {
			reliableTransport.StopAndConsumeBuffer(reorder)
		}

		if shouldReceive {
			select {
				case recvd := <- received: {
					replyTo = recvd.ReplyTo
					assert.Equal(t, addedRouteId, recvd.ReplyTo.routeId)
					assert.Equal(t, contents, recvd.Contents)
				}
				case <-time.After(5*time.Second):
					panic("Test timed out")
			}
		} else {
			select {
				case <- received: {
					panic("Should not receive")
				}
				case <-time.After(5*time.Second): {
					break
				}
			}
		}
	}

	if sendBack {
		back_received := make(chan MeshMessage, 10)
		nodes[0].SetReceiveChannel(back_received)
		replyContents := []byte{6,44,2,1,1,1,1,2}
		assert.Nil(t, nodes[nPeers-1].SendMessageBackThruRoute(replyTo, replyContents, reliable))
		select {
			case recvd_back := <- back_received: {
				assert.Equal(t, replyContents, recvd_back.Contents)
			}
			case <-time.After(5*time.Second):
				panic("Test timed out")
		}
	}
}

func TestSendDirectUnreliably(t *testing.T) {
	contents := []byte{4,66,7,44,33}
	sendTest(t, 2, false, false, false, false, contents)
}

func TestSendDirectUnreliablyNegative(t *testing.T) {
	contents := []byte{4,66,7,44,33}
	sendTest(t, 2, false, true, false, false, contents)
}

func TestSendDirectReliably(t *testing.T) {
	contents := []byte{4,66,7,44,33}
	sendTest(t, 2, true, false, false, false, contents)
}

func TestSendLongMessage(t *testing.T) {
	contents := []byte{}
	for i := 0; i < 25670 ; i++ {
		contents = append(contents, (byte)(i))
	}
	sendTest(t, 2, false, false, false, false, contents)
}

func TestSendLongMessageWithReorder(t *testing.T) {
	contents := []byte{}
	for i := 0; i < 25670 ; i++ {
		contents = append(contents, (byte)(i))
	}
	sendTest(t, 2, false, false, true, false, contents)
}

func TestLongRoute(t *testing.T) {
	contents := []byte{4,66,7,44,33}
	sendTest(t, 5, true, false, false, false, contents)
}

func TestShortSendBack(t *testing.T) {
	contents := []byte{1,44,2,22,11,22}
	sendTest(t, 2, true, false, false, true, contents)
}

func TestMediumSendBack(t *testing.T) {
	contents := []byte{1,44,2,22,11,22}
	sendTest(t, 3, true, false, false, true, contents)
}

func TestLongSendBack(t *testing.T) {
	contents := []byte{1,44,2,22,11,22}
	sendTest(t, 5, true, false, false, true, contents)
}

// Refragment!

// Reorder messages with establish
// Establish route and send unreliable
/// TODO! Needs to pass a reliable flag in base?


// Send back test
// Send back long route
// Refragment test

// Expire old routes, messages test
// Route expiry test
// Packet loss test
// Multiple transport test
// Threading test


package mesh

import (
	"sort"
	"testing"
	"time"
)

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/transport/transport"
	"github.com/stretchr/testify/assert"
)

func sortPubKeys(pubKeys []cipher.PubKey) []cipher.PubKey {
	var ret cipher.PubKeySlice = pubKeys
	sort.Sort(ret)
	return ret
}

func TestManageTransports(t *testing.T) {
	transport_a := transport.NewStubTransport(t, 512)
	transport_b := transport.NewStubTransport(t, 512)
	test_key_a := cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	nodeConfig := NodeConfig{
		PubKey:                        test_key_a,
		ChaCha20Key:                   [32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3},
		MaximumForwardingDuration:     time.Minute,
		RefreshRouteDuration:          10 * time.Second,
		ExpireMessagesInterval:        time.Second,
		ExpireRoutesInterval:          time.Second,
		TimeToAssembleMessage:         2 * time.Second,
		TransportMessageChannelLength: 100, // Transport message channel length
	}
	node, error := NewNode(nodeConfig)
	assert.Nil(t, error)
	assert.Equal(t, []transport.Transport{}, node.GetTransports())
	node.AddTransport(transport_a, nodeConfig.ChaCha20Key)
	assert.Equal(t, []transport.Transport{transport_a}, node.GetTransports())
	node.AddTransport(transport_b, nodeConfig.ChaCha20Key)
	assert.Equal(t, []transport.Transport{transport_a, transport_b}, node.GetTransports())
	node.RemoveTransport(transport_a)
	assert.Equal(t, []transport.Transport{transport_b}, node.GetTransports())
	node.RemoveTransport(transport_b)
	assert.Equal(t, []transport.Transport{}, node.GetTransports())
}

func TestConnectedPeers(t *testing.T) {
	transport_a := transport.NewStubTransport(t, 512)
	transport_b := transport.NewStubTransport(t, 512)
	test_key_a := cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	nodeConfig := NodeConfig{
		PubKey:                        test_key_a,
		ChaCha20Key:                   [32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3},
		MaximumForwardingDuration:     time.Minute,
		RefreshRouteDuration:          10 * time.Second,
		ExpireMessagesInterval:        time.Second,
		ExpireRoutesInterval:          time.Second,
		TimeToAssembleMessage:         2 * time.Second,
		TransportMessageChannelLength: 100, // Transport message channel length
	}
	node, error := NewNode(nodeConfig)
	assert.Nil(t, error)
	peer_a := cipher.NewPubKey([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	peer_b := cipher.NewPubKey([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	peer_c := cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	transport_a.AddStubbedPeer(peer_a, nil)
	transport_a.AddStubbedPeer(peer_b, nil)
	transport_b.AddStubbedPeer(peer_c, nil)

	assert.False(t, node.ConnectedToPeer(peer_a))
	assert.False(t, node.ConnectedToPeer(peer_b))
	assert.False(t, node.ConnectedToPeer(peer_c))
	assert.Equal(t, []cipher.PubKey{}, sortPubKeys(node.GetConnectedPeers()))
	node.AddTransport(transport_a, nodeConfig.ChaCha20Key)
	assert.Equal(t, []cipher.PubKey{peer_a, peer_b}, sortPubKeys(node.GetConnectedPeers()))
	assert.True(t, node.ConnectedToPeer(peer_a))
	assert.True(t, node.ConnectedToPeer(peer_b))
	assert.False(t, node.ConnectedToPeer(peer_c))

	node.AddTransport(transport_b, nodeConfig.ChaCha20Key)
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
	maxDatagramLength uint,
	newPubKey cipher.PubKey) (node *Node,
	unreliableTransport *transport.StubTransport,
	reliableTransport *transport.StubTransport) {
	unreliableTransport = transport.NewStubTransport(t, maxDatagramLength)
	reliableTransport = transport.NewStubTransport(t, maxDatagramLength)
	var error error
	nodeConfig := NodeConfig{
		PubKey:                        newPubKey,
		ChaCha20Key:                   [32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3},
		MaximumForwardingDuration:     time.Minute,
		RefreshRouteDuration:          time.Second,
		ExpireMessagesInterval:        time.Second,
		ExpireRoutesInterval:          time.Second,
		TimeToAssembleMessage:         2 * time.Second,
		TransportMessageChannelLength: 100, // Transport message channel length
	}
	node, error = NewNode(nodeConfig)
	assert.Nil(t, error)
	node.AddTransport(unreliableTransport, nodeConfig.ChaCha20Key)
	node.AddTransport(reliableTransport, nodeConfig.ChaCha20Key)
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
	maxDatagramLengths := []uint{512, 450, 1000, 150, 200}
	for i := (uint)(0); i < n; i++ {
		pubKey := cipher.PubKey{}
		pubKey[0] = (byte)(i + 1)
		nodes[i], unreliableTransports[i], reliableTransports[i] = SetupNode(t, maxDatagramLengths[i%((uint)(len(maxDatagramLengths)))], pubKey)
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

func sendTest(t *testing.T, nPeers int, reliable bool, dropFirst bool, reorder bool, sendBack bool, contents []byte) {
	if nPeers < 2 {
		panic("Fewer than 2 peers doesn't make sense")
	}

	allConnections := make([][]int, 0)
	for from_idx := 0; from_idx < nPeers; from_idx++ {
		toConnections := make([]int, 0)
		for i := 0; i < nPeers; i++ {
			toConnections = append(toConnections, 0)
		}

		for to_idx := from_idx - 1; to_idx <= from_idx+1; to_idx++ {
			if to_idx >= 0 && to_idx != from_idx && to_idx < nPeers {
				toConnections[to_idx] = 1
			}
		}
		allConnections = append(allConnections, toConnections)
	}
	nodes, to_close, unreliableTransport, reliableTransport := SetupNodes((uint)(nPeers), allConnections, t)
	defer close(to_close)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()

	if dropFirst {
		for _, unreliableTransport := range unreliableTransport {
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

	for dropFirstIdx := 0; dropFirstIdx < 2; dropFirstIdx++ {
		shouldReceive := true
		if dropFirst && dropFirstIdx == 0 {
			shouldReceive = false
		}

		for _, unreliableTransport := range unreliableTransport {
			unreliableTransport.StartBuffer()
			unreliableTransport.SetIgnoreSendStatus(!shouldReceive)
		}
		for _, reliableTransport := range reliableTransport {
			reliableTransport.StartBuffer()
		}

		send_err, route_id := nodes[0].SendMessageToPeer(terminating_id, contents, reliable)
		assert.Nil(t, send_err)
		assert.Equal(t, addedRouteId, route_id)

		for _, unreliableTransport := range unreliableTransport {
			unreliableTransport.StopAndConsumeBuffer(reorder, 0)
		}

		for _, reliableTransport := range reliableTransport {
			reliableTransport.StopAndConsumeBuffer(reorder, 0)
		}

		if shouldReceive {
			select {
			case recvd := <-received:
				{
					replyTo = recvd.ReplyTo
					assert.Equal(t, addedRouteId, recvd.ReplyTo.routeId)
					assert.Equal(t, contents, recvd.Contents)
				}
			case <-time.After(5 * time.Second):
				panic("Test timed out")
			}
		} else {
			select {
			case <-received:
				{
					panic("Should not receive")
				}
			case <-time.After(5 * time.Second):
				{
					break
				}
			}
		}
	}

	if sendBack {
		back_received := make(chan MeshMessage, 10)
		nodes[0].SetReceiveChannel(back_received)
		replyContents := []byte{6, 44, 2, 1, 1, 1, 1, 2}
		assert.Nil(t, nodes[nPeers-1].SendMessageBackThruRoute(replyTo, replyContents, reliable))
		select {
		case recvd_back := <-back_received:
			{
				assert.Equal(t, replyContents, recvd_back.Contents)
			}
		case <-time.After(10 * time.Second):
			panic("Test timed out")
		}
	}
}

func TestSendDirectUnreliablyPositive(t *testing.T) {
	contents := []byte{4, 66, 7, 44, 33}
	sendTest(t, 2, false, false, false, false, contents)
}

func TestSendDirectUnreliablyNegative(t *testing.T) {
	contents := []byte{4, 66, 7, 44, 33}
	sendTest(t, 2, false, true, false, false, contents)
}

func TestSendDirectReliably(t *testing.T) {
	contents := []byte{4, 66, 7, 44, 33}
	sendTest(t, 2, true, false, false, false, contents)
}

func TestSendLongMessage(t *testing.T) {
	contents := []byte{}
	for i := 0; i < 25670; i++ {
		contents = append(contents, (byte)(i))
	}
	sendTest(t, 2, false, false, false, false, contents)
}

func TestSendLongMessageWithReorder(t *testing.T) {
	contents := []byte{}
	for i := 0; i < 25670; i++ {
		contents = append(contents, (byte)(i))
	}
	sendTest(t, 2, false, false, true, false, contents)
}

func TestLongRoute(t *testing.T) {
	contents := []byte{4, 66, 7, 44, 33}
	sendTest(t, 5, true, false, false, false, contents)
}

func TestShortSendBack(t *testing.T) {
	contents := []byte{1, 44, 2, 22, 11, 22}
	sendTest(t, 2, true, false, false, true, contents)
}

func TestMediumSendBack(t *testing.T) {
	contents := []byte{1, 44, 2, 22, 11, 22}
	sendTest(t, 3, true, false, false, true, contents)
}

func TestLongSendBack(t *testing.T) {
	contents := []byte{1, 44, 2, 22, 11, 22}
	sendTest(t, 5, true, false, false, true, contents)
}

// Refragmentation test (sendTest varies the datagram length)
func TestLongSendLongMessage(t *testing.T) {
	contents := []byte{}
	for i := 0; i < 25670; i++ {
		contents = append(contents, (byte)(i))
	}
	sendTest(t, 5, true, false, false, false, contents)
}

func TestSendThruRoute(t *testing.T) {
	allConnections := [][]int{
		[]int{0, 1},
		[]int{1, 0},
	}
	nodes, to_close, _, _ := SetupNodes((uint)(2), allConnections, t)
	defer close(to_close)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()
	received := make(chan MeshMessage, 10)
	nodes[1].SetReceiveChannel(received)
	contents := []byte{1, 44, 2, 22, 11, 22}
	addedRouteId := RouteId{}
	addedRouteId[0] = 55
	addedRouteId[1] = 4
	assert.Nil(t, nodes[0].AddRoute(addedRouteId, nodes[1].GetConfig().PubKey))
	assert.Nil(t, nodes[0].SendMessageThruRoute(addedRouteId, contents, true))

	select {
	case recvd := <-received:
		{
			assert.Equal(t, contents, recvd.Contents)
		}
	case <-time.After(5 * time.Second):
		panic("Test timed out")
	}
}

func TestRouteExpiry(t *testing.T) {
	allConnections := [][]int{
		[]int{0, 1, 0},
		[]int{1, 0, 1},
		[]int{0, 1, 0},
	}

	nodes, to_close, _, reliableTransports := SetupNodes((uint)(3), allConnections, t)
	defer close(to_close)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()
	addedRouteId := RouteId{}
	addedRouteId[0] = 55
	addedRouteId[1] = 4
	assert.Nil(t, nodes[0].AddRoute(addedRouteId, nodes[1].GetConfig().PubKey))
	{
		lastConfirmed, err := nodes[0].GetRouteLastConfirmed(addedRouteId)
		assert.Nil(t, err)
		assert.Zero(t, lastConfirmed.Unix())
	}
	assert.Nil(t, nodes[0].ExtendRoute(addedRouteId, nodes[2].GetConfig().PubKey, time.Second))
	assert.NotZero(t, nodes[1].debug_countRoutes())
	var afterExtendConfirmedTime time.Time
	{
		lastConfirmed, err := nodes[0].GetRouteLastConfirmed(addedRouteId)
		assert.Nil(t, err)
		afterExtendConfirmedTime = lastConfirmed
	}
	time.Sleep(5 * time.Second)
	assert.NotZero(t, nodes[1].debug_countRoutes())
	var afterWaitConfirmedTime time.Time
	{
		lastConfirmed, err := nodes[0].GetRouteLastConfirmed(addedRouteId)
		assert.Nil(t, err)
		afterWaitConfirmedTime = lastConfirmed
	}
	// Don't allow refreshes to get thru
	reliableTransports[0].SetIgnoreSendStatus(true)
	time.Sleep(5 * time.Second)
	var afterIgnoreConfirmedTime time.Time
	{
		lastConfirmed, err := nodes[0].GetRouteLastConfirmed(addedRouteId)
		assert.Nil(t, err)
		afterIgnoreConfirmedTime = lastConfirmed
	}
	assert.Zero(t, nodes[1].debug_countRoutes())
	assert.NotZero(t, afterExtendConfirmedTime)
	assert.NotZero(t, afterWaitConfirmedTime)
	assert.NotEqual(t, afterExtendConfirmedTime, afterWaitConfirmedTime)
	assert.Equal(t, afterWaitConfirmedTime, afterIgnoreConfirmedTime)
}

func TestDeleteRoute(t *testing.T) {
	allConnections := [][]int{
		[]int{0, 1, 0},
		[]int{1, 0, 1},
		[]int{0, 1, 0},
	}

	nodes, to_close, _, _ := SetupNodes((uint)(3), allConnections, t)
	defer close(to_close)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()
	addedRouteId := RouteId{}
	addedRouteId[0] = 55
	addedRouteId[1] = 4
	assert.Nil(t, nodes[0].AddRoute(addedRouteId, nodes[1].GetConfig().PubKey))
	assert.Nil(t, nodes[0].ExtendRoute(addedRouteId, nodes[2].GetConfig().PubKey, time.Second))
	time.Sleep(5 * time.Second)
	assert.NotZero(t, nodes[0].debug_countRoutes())
	assert.NotZero(t, nodes[1].debug_countRoutes())
	assert.Nil(t, nodes[0].DeleteRoute(addedRouteId))
	time.Sleep(1 * time.Second)
	assert.Zero(t, nodes[0].debug_countRoutes())
	assert.Zero(t, nodes[1].debug_countRoutes())
}

func TestMessageExpiry(t *testing.T) {
	allConnections := [][]int{
		[]int{0, 1},
		[]int{1, 0},
	}
	nodes, to_close, _, reliableTransports := SetupNodes((uint)(2), allConnections, t)
	defer close(to_close)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()
	addedRouteId := RouteId{}
	addedRouteId[0] = 66

	contents := []byte{}
	for i := 0; i < 25670; i++ {
		contents = append(contents, (byte)(i))
	}

	assert.Nil(t, nodes[0].AddRoute(addedRouteId, nodes[1].GetConfig().PubKey))

	reliableTransports[0].StartBuffer()
	assert.Nil(t, nodes[0].SendMessageThruRoute(addedRouteId, contents, true))
	// Drop ten, so the message will never be reassembled
	reliableTransports[0].StopAndConsumeBuffer(true, 10)

	time.Sleep(1 * time.Second)
	assert.NotZero(t, nodes[1].debug_countMessages())
	time.Sleep(10 * time.Second)
	assert.Zero(t, nodes[1].debug_countMessages())
}

func TestLongRouteUnreliable(t *testing.T) {
	allConnections := [][]int{
		[]int{0, 1, 0},
		[]int{1, 0, 1},
		[]int{0, 1, 0},
	}

	nodes, to_close, unreliableTransports, reliableTransports := SetupNodes((uint)(3), allConnections, t)
	defer close(to_close)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()
	received := make(chan MeshMessage, 10)
	nodes[2].SetReceiveChannel(received)
	addedRouteId := RouteId{}
	addedRouteId[0] = 77
	assert.Nil(t, nodes[0].AddRoute(addedRouteId, nodes[1].GetConfig().PubKey))
	assert.Nil(t, nodes[0].ExtendRoute(addedRouteId, nodes[2].GetConfig().PubKey, time.Second))

	contents := []byte{2, 3, 44, 22, 11, 3, 3, 3, 3, 5}

	assert.Nil(t, nodes[0].SendMessageThruRoute(addedRouteId, contents, false))

	select {
	case recvd := <-received:
		{
			assert.Equal(t, contents, recvd.Contents)
		}
	case <-time.After(5 * time.Second):
		panic("Test timed out")
	}

	assert.NotZero(t, unreliableTransports[0].CountNumMessagesSent())
	assert.NotZero(t, reliableTransports[0].CountNumMessagesSent())
	assert.NotZero(t, unreliableTransports[1].CountNumMessagesSent())
	assert.NotZero(t, reliableTransports[1].CountNumMessagesSent())
	assert.Zero(t, unreliableTransports[2].CountNumMessagesSent())
	// ACKs going back don't count
	assert.Zero(t, reliableTransports[2].CountNumMessagesSent())
}

// Tests TODO

// Establish route and send unreliable

// Packet loss test
// Multiple transport test
// Threading test

package mesh

import (
	"sort"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/transport"
	"github.com/stretchr/testify/assert"
)

func sortPubKeys(pubKeys []cipher.PubKey) []cipher.PubKey {
	var keys cipher.PubKeySlice = pubKeys
	sort.Sort(keys)
	return keys
}

func TestManageTransports(t *testing.T) {
	transportA := transport.NewStubTransport(t, 512)
	transportB := transport.NewStubTransport(t, 512)
	testKeyA := cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	nodeConfig := domain.NodeConfig{
		PubKey: testKeyA,
		//ChaCha20Key:                   [32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3},
		MaximumForwardingDuration:     time.Minute,
		RefreshRouteDuration:          10 * time.Second,
		ExpireMessagesInterval:        time.Second,
		ExpireRoutesInterval:          time.Second,
		TimeToAssembleMessage:         2 * time.Second,
		TransportMessageChannelLength: 100, // Transport message channel length
	}
	node, err := NewNode(nodeConfig)
	assert.Nil(t, err)
	assert.Equal(t, []transport.ITransport{}, node.GetTransports())
	node.AddTransport(transportA)
	assert.Equal(t, []transport.ITransport{transportA}, node.GetTransports())
	node.AddTransport(transportB)
	assert.Equal(t, []transport.ITransport{transportA, transportB}, node.GetTransports())
	node.RemoveTransport(transportA)
	assert.Equal(t, []transport.ITransport{transportB}, node.GetTransports())
	node.RemoveTransport(transportB)
	assert.Equal(t, []transport.ITransport{}, node.GetTransports())
}

func TestConnectedPeers(t *testing.T) {
	transportA := transport.NewStubTransport(t, 512)
	transportB := transport.NewStubTransport(t, 512)
	testKeyA := cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	nodeConfig := domain.NodeConfig{
		PubKey: testKeyA,
		//ChaCha20Key:                   [32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3},
		MaximumForwardingDuration:     time.Minute,
		RefreshRouteDuration:          10 * time.Second,
		ExpireMessagesInterval:        time.Second,
		ExpireRoutesInterval:          time.Second,
		TimeToAssembleMessage:         2 * time.Second,
		TransportMessageChannelLength: 100, // Transport message channel length
	}
	node, err := NewNode(nodeConfig)
	assert.Nil(t, err)
	peerA := cipher.NewPubKey([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	peerB := cipher.NewPubKey([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	peerC := cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	transportA.AddStubbedPeer(peerA, nil)
	transportA.AddStubbedPeer(peerB, nil)
	transportB.AddStubbedPeer(peerC, nil)

	assert.False(t, node.ConnectedToPeer(peerA))
	assert.False(t, node.ConnectedToPeer(peerB))
	assert.False(t, node.ConnectedToPeer(peerC))
	assert.Equal(t, []cipher.PubKey{}, sortPubKeys(node.GetConnectedPeers()))
	node.AddTransport(transportA)
	assert.Equal(t, []cipher.PubKey{peerA, peerB}, sortPubKeys(node.GetConnectedPeers()))
	assert.True(t, node.ConnectedToPeer(peerA))
	assert.True(t, node.ConnectedToPeer(peerB))
	assert.False(t, node.ConnectedToPeer(peerC))

	node.AddTransport(transportB)
	assert.Equal(t, []cipher.PubKey{peerA, peerB, peerC}, sortPubKeys(node.GetConnectedPeers()))
	assert.True(t, node.ConnectedToPeer(peerA))
	assert.True(t, node.ConnectedToPeer(peerB))
	assert.True(t, node.ConnectedToPeer(peerC))
	assert.True(t, transportA.ConnectedToPeer(peerA))
	node.RemoveTransport(transportA)
	assert.False(t, node.ConnectedToPeer(peerA))
	assert.False(t, node.ConnectedToPeer(peerB))
	assert.True(t, node.ConnectedToPeer(peerC))

	assert.Equal(t, []cipher.PubKey{peerC}, sortPubKeys(node.GetConnectedPeers()))
	node.RemoveTransport(transportB)
	assert.Equal(t, []cipher.PubKey{}, sortPubKeys(node.GetConnectedPeers()))
	assert.False(t, node.ConnectedToPeer(peerA))
	assert.False(t, node.ConnectedToPeer(peerB))
	assert.False(t, node.ConnectedToPeer(peerC))
}

func SetupNode(t *testing.T,
	maxDatagramLength uint,
	newPubKey cipher.PubKey) (node *Node,
	stubTransport *transport.StubTransport) {
	stubTransport = transport.NewStubTransport(t, maxDatagramLength)
	var err error
	nodeConfig := domain.NodeConfig{
		PubKey: newPubKey,
		//ChaCha20Key:                   [32]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3},
		MaximumForwardingDuration:     time.Minute,
		RefreshRouteDuration:          time.Second,
		ExpireMessagesInterval:        time.Second,
		ExpireRoutesInterval:          time.Second,
		TimeToAssembleMessage:         2 * time.Second,
		TransportMessageChannelLength: 100,
	}
	node, err = NewNode(nodeConfig)
	assert.Nil(t, err)
	node.AddTransport(stubTransport)
	return
}

// Nodes each have one transport
// All nodes receive all other nodes' messages, but stub transport filters
func SetupNodes(n uint, connections [][]int, t *testing.T) (nodes []*Node, to_close chan []byte,
	transports []*transport.StubTransport) {
	nodes = make([]*Node, n)
	transports = make([]*transport.StubTransport, n)
	to_close = make(chan []byte, 20)
	sentMessages := make(chan []byte, 20)
	maxDatagramLengths := []uint{512, 450, 1000, 150, 200}
	for i := (uint)(0); i < n; i++ {
		pubKey := cipher.PubKey{}
		pubKey[0] = (byte)(i + 1)
		nodes[i], transports[i] = SetupNode(t, maxDatagramLengths[i%((uint)(len(maxDatagramLengths)))], pubKey)
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
	return nodes, sentMessages, transports
}

func sendTest(t *testing.T, nPeers int, dropFirst bool, reorder bool, sendBack bool, contents []byte) {
	if nPeers < 2 {
		panic("Fewer than 2 peers doesn't make sense")
	}

	allConnections := make([][]int, 0)
	for fromIndex := 0; fromIndex < nPeers; fromIndex++ {
		toConnections := make([]int, 0)
		for i := 0; i < nPeers; i++ {
			toConnections = append(toConnections, 0)
		}

		for toIndex := fromIndex - 1; toIndex <= fromIndex+1; toIndex++ {
			if toIndex >= 0 && toIndex != fromIndex && toIndex < nPeers {
				toConnections[toIndex] = 1
			}
		}
		allConnections = append(allConnections, toConnections)
	}
	nodes, toClose, transports := SetupNodes((uint)(nPeers), allConnections, t)
	defer close(toClose)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()

	receivedMessages := make(chan domain.MeshMessage, 10)
	nodes[nPeers-1].SetReceiveChannel(receivedMessages)

	terminatingID := nodes[nPeers-1].GetConfig().PubKey

	addedRouteID := domain.RouteID{}
	addedRouteID[0] = 22
	assert.Nil(t, nodes[0].AddRoute(addedRouteID, nodes[1].GetConfig().PubKey))

	for extendIdx := 2; extendIdx < nPeers; extendIdx++ {
		assert.Nil(t, nodes[0].ExtendRoute(addedRouteID, nodes[extendIdx].GetConfig().PubKey, time.Second))
	}

	var replyTo domain.ReplyTo

	for dropFirstIdx := 0; dropFirstIdx < 2; dropFirstIdx++ {
		shouldReceive := true
		if dropFirst && dropFirstIdx == 0 {
			shouldReceive = false
		}

		for _, transportToPeer := range transports {
			transportToPeer.StartBuffer()
		}

		err, routeID := nodes[0].SendMessageToPeer(terminatingID, contents)
		assert.Nil(t, err)
		assert.Equal(t, addedRouteID, routeID)

		for _, transportToPeer := range transports {
			transportToPeer.StopAndConsumeBuffer(reorder, 0)
		}

		if shouldReceive {
			select {
			case receivedMessage := <-receivedMessages:
				{
					replyTo = receivedMessage.ReplyTo
					assert.Equal(t, addedRouteID, receivedMessage.ReplyTo.RouteID)
					assert.Equal(t, contents, receivedMessage.Contents)
				}
			case <-time.After(5 * time.Second):
				panic("Test timed out")
			}
		} else {
			select {
			case <-receivedMessages:
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
		backReceivedMessages := make(chan domain.MeshMessage, 10)
		nodes[0].SetReceiveChannel(backReceivedMessages)
		replyContents := []byte{6, 44, 2, 1, 1, 1, 1, 2}
		assert.Nil(t, nodes[nPeers-1].SendMessageBackThruRoute(replyTo, replyContents))
		select {
		case receivedBack := <-backReceivedMessages:
			{
				assert.Equal(t, replyContents, receivedBack.Contents)
			}
		case <-time.After(10 * time.Second):
			panic("Test timed out")
		}
	}
}

func TestSendDirect(t *testing.T) {
	contents := []byte{4, 66, 7, 44, 33}
	sendTest(t, 2, false, false, false, contents)
}

func TestSendLongMessage(t *testing.T) {
	contents := []byte{}
	for i := 0; i < 25670; i++ {
		contents = append(contents, (byte)(i))
	}
	sendTest(t, 2, false, false, false, contents)
}

func TestSendLongMessageWithReorder(t *testing.T) {
	contents := []byte{}
	for i := 0; i < 25670; i++ {
		contents = append(contents, (byte)(i))
	}
	sendTest(t, 2, false, true, false, contents)
}

func TestLongRoute(t *testing.T) {
	contents := []byte{4, 66, 7, 44, 33}
	sendTest(t, 5, false, false, false, contents)
}

func TestShortSendBack(t *testing.T) {
	contents := []byte{1, 44, 2, 22, 11, 22}
	sendTest(t, 2, false, false, true, contents)
}

func TestMediumSendBack(t *testing.T) {
	contents := []byte{1, 44, 2, 22, 11, 22}
	sendTest(t, 3, false, false, true, contents)
}

func TestLongSendBack(t *testing.T) {
	contents := []byte{1, 44, 2, 22, 11, 22}
	sendTest(t, 5, false, false, true, contents)
}

// Refragmentation test (sendTest varies the datagram length)
func TestLongSendLongMessage(t *testing.T) {
	contents := []byte{}
	for i := 0; i < 25670; i++ {
		contents = append(contents, (byte)(i))
	}
	sendTest(t, 5, false, false, false, contents)
}

func TestSendThruRoute(t *testing.T) {
	allConnections := [][]int{
		[]int{0, 1},
		[]int{1, 0},
	}
	nodes, toClose, _ := SetupNodes((uint)(2), allConnections, t)
	defer close(toClose)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()
	receivedMessages := make(chan domain.MeshMessage, 10)
	nodes[1].SetReceiveChannel(receivedMessages)
	contents := []byte{1, 44, 2, 22, 11, 22}
	addedRouteID := domain.RouteID{}
	addedRouteID[0] = 55
	addedRouteID[1] = 4
	assert.Nil(t, nodes[0].AddRoute(addedRouteID, nodes[1].GetConfig().PubKey))
	assert.Nil(t, nodes[0].SendMessageThruRoute(addedRouteID, contents))

	select {
	case receivedMessage := <-receivedMessages:
		{
			assert.Equal(t, contents, receivedMessage.Contents)
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

	nodes, toClose, transports := SetupNodes((uint)(3), allConnections, t)
	defer close(toClose)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()

	addedRouteID := domain.RouteID{}
	addedRouteID[0] = 55
	addedRouteID[1] = 4

	assert.Nil(t, nodes[0].AddRoute(addedRouteID, nodes[1].GetConfig().PubKey))
	{
		lastConfirmed, err := nodes[0].GetRouteLastConfirmed(addedRouteID)
		assert.Nil(t, err)
		assert.Zero(t, lastConfirmed.Unix())
	}
	assert.Nil(t, nodes[0].ExtendRoute(addedRouteID, nodes[2].GetConfig().PubKey, time.Second))
	assert.NotZero(t, nodes[1].debug_countRoutes())

	var afterExtendConfirmedTime time.Time
	{
		lastConfirmed, err := nodes[0].GetRouteLastConfirmed(addedRouteID)
		assert.Nil(t, err)
		afterExtendConfirmedTime = lastConfirmed
	}

	time.Sleep(5 * time.Second)
	assert.NotZero(t, nodes[1].debug_countRoutes())
	var afterWaitConfirmedTime time.Time
	{
		lastConfirmed, err := nodes[0].GetRouteLastConfirmed(addedRouteID)
		assert.Nil(t, err)
		afterWaitConfirmedTime = lastConfirmed
	}

	// Don't allow refreshes to get thru
	transports[0].SetIgnoreSendStatus(true)
	time.Sleep(5 * time.Second)
	var afterIgnoreConfirmedTime time.Time
	{
		lastConfirmed, err := nodes[0].GetRouteLastConfirmed(addedRouteID)
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

	nodes, toClose, _ := SetupNodes((uint)(3), allConnections, t)
	defer close(toClose)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()
	addedRouteID := domain.RouteID{}
	addedRouteID[0] = 55
	addedRouteID[1] = 4
	assert.Nil(t, nodes[0].AddRoute(addedRouteID, nodes[1].GetConfig().PubKey))
	assert.Nil(t, nodes[0].ExtendRoute(addedRouteID, nodes[2].GetConfig().PubKey, time.Second))
	time.Sleep(5 * time.Second)
	assert.NotZero(t, nodes[0].debug_countRoutes())
	assert.NotZero(t, nodes[1].debug_countRoutes())
	assert.Nil(t, nodes[0].DeleteRoute(addedRouteID))
	time.Sleep(1 * time.Second)
	assert.Zero(t, nodes[0].debug_countRoutes())
	assert.Zero(t, nodes[1].debug_countRoutes())
}

func TestMessageExpiry(t *testing.T) {
	allConnections := [][]int{
		[]int{0, 1},
		[]int{1, 0},
	}
	nodes, toClose, transports := SetupNodes((uint)(2), allConnections, t)
	defer close(toClose)
	defer func() {
		for _, node := range nodes {
			node.Close()
		}
	}()
	addedRouteID := domain.RouteID{}
	addedRouteID[0] = 66

	contents := []byte{}
	for i := 0; i < 25670; i++ {
		contents = append(contents, (byte)(i))
	}

	assert.Nil(t, nodes[0].AddRoute(addedRouteID, nodes[1].GetConfig().PubKey))

	transports[0].StartBuffer()
	assert.Nil(t, nodes[0].SendMessageThruRoute(addedRouteID, contents))
	// Drop ten, so the message will never be reassembled
	transports[0].StopAndConsumeBuffer(true, 10)

	time.Sleep(1 * time.Second)
	assert.NotZero(t, nodes[1].debug_countMessages())
	time.Sleep(10 * time.Second)
	assert.Zero(t, nodes[1].debug_countMessages())
}

// Tests TODO

// Establish route and send unreliable

// Packet loss test
// Multiple transport test
// Threading test

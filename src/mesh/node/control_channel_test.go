package mesh

import (
	"fmt"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/transport"
	"github.com/stretchr/testify/assert"
)

// TODO refactor and use SetupNode() instead
func SetupNode2(t *testing.T, newPubKey cipher.PubKey) *Node {
	nodeConfig := NodeConfig{
		PubKey: newPubKey,
		MaximumForwardingDuration:     time.Minute,
		RefreshRouteDuration:          10 * time.Second,
		ExpireRoutesInterval:          10 * time.Second,
		TransportMessageChannelLength: 100,
	}
	node, err := NewNode(nodeConfig)
	assert.Nil(t, err)
	return node
}

// TODO refactor and use SetupNodes() instead
func SetupNodes2(t *testing.T, n uint, connections [][]int) ([]*Node, []*transport.StubTransport) {
	nodes := make([]*Node, n)
	transports := []*transport.StubTransport{}
	for i := (uint)(0); i < n; i++ {
		pubKey := cipher.PubKey{}
		pubKey[0] = (byte)(i + 1)
		nodes[i] = SetupNode2(t, pubKey)
	}

	for i := (uint)(0); i < n; i++ {
		transportsFrom := []*transport.StubTransport{}
		for j := (uint)(0); j < n; j++ {
			if connections[i][j] != 0 {
				transportFrom := transport.NewStubTransport(t, 512)
				transportTo := transport.NewStubTransport(t, 512)
				transportFrom.SetStubbedPeer(nodes[j].GetConfig().PubKey, transportTo)
				transportFrom.MessagesReceived = nodes[i].transportsMessagesReceived
				transportTo.MessagesReceived = nodes[j].transportsMessagesReceived
				transportsFrom = append(transportsFrom, transportFrom)
				nodes[i].AddTransport(transportFrom)
			}
		}
		transports = append(transports, transportsFrom...)
	}
	return nodes, transports
}

func TestSendOneDirection(t *testing.T) {

	connections := [][]int{
		[]int{0, 1, 1},
		[]int{1, 0, 1},
		[]int{0, 1, 0},
	}

	nodes, transports := SetupNodes2(t, 3, connections)

	receivedMessagesA := make(chan domain.MeshMessage, 10)
	nodes[0].SetReceiveChannel(receivedMessagesA)

	receivedMessagesC := make(chan domain.MeshMessage, 10)
	nodes[2].SetReceiveChannel(receivedMessagesC)

	A2BRouteID := domain.RouteID{}
	A2BRouteID[0] = 1
	assert.Nil(t, nodes[0].AddRoute(A2BRouteID, nodes[1].GetConfig().PubKey))
	assert.Equal(t, 1, nodes[0].DebugCountRoutes())

	A2CRouteID := domain.RouteID{}
	A2CRouteID[0] = 2
	assert.Nil(t, nodes[0].AddRoute(A2CRouteID, nodes[2].GetConfig().PubKey))
	assert.Equal(t, 2, nodes[0].DebugCountRoutes())

	multiHopRouteID := domain.RouteID{}
	multiHopRouteID[0] = 3
	assert.Nil(t, nodes[0].AddRoute(multiHopRouteID, nodes[1].GetConfig().PubKey)) // first hop A->B
	assert.Equal(t, 3, nodes[0].DebugCountRoutes())

	// Create new Control Channel in Node B
	transports[0].StartBuffer()
	err := nodes[0].SendSetControlChannelMessage(A2BRouteID)
	assert.Nil(t, err)
	transports[0].StopAndConsumeBuffer(true, 0)

	receivedMessage := <-receivedMessagesA
	nodeBChannelID, _ := uuid.FromBytes(receivedMessage.Contents)
	fmt.Println("Channel ID: " + nodeBChannelID.String())

	// Setup route B->A
	transports[0].StartBuffer()
	err = nodes[0].SendSetRouteControlMessage(nodeBChannelID, A2BRouteID, multiHopRouteID, nodes[2].Config.PubKey)
	assert.Nil(t, err)
	transports[0].StopAndConsumeBuffer(true, 0)
	receivedMessage = <-receivedMessagesA
	fmt.Println(receivedMessage.Contents)

	// Send user message
	transports[2].StartBuffer()
	contents := []byte{4, 66, 7, 44, 33}
	err = nodes[0].SendMessageThruRoute(multiHopRouteID, contents)
	assert.Nil(t, err)
	transports[2].StopAndConsumeBuffer(true, 0)
	receivedMessage = <-receivedMessagesC
	fmt.Println(receivedMessage.Contents)

	assert.Equal(t, contents, receivedMessage.Contents)
}

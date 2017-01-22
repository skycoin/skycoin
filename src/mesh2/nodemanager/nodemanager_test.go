package nodemanager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/node"
)

func TestAddingNodes(t *testing.T) {

	nm := NewNodeManager()
	assert.Len(t, nm.NodeList, 0, "Error expected 0 nodes")
	nm.AddNewNode()
	assert.Len(t, nm.NodeList, 1, "Error expected 1 nodes")
	nm.AddNewNode()
	nm.AddNewNode()
	nm.AddNewNode()
	nm.AddNewNode()
	assert.Len(t, nm.NodeList, 5, "Error expected 5 nodes")
}

func TestConnectTwoNodes(t *testing.T) {

	nm := NewNodeManager()
	id1 := nm.AddNewNode()
	id2 := nm.AddNewNode()
	node1, err := nm.GetNodeById(id1)
	assert.Nil(t, err)
	node2, err := nm.GetNodeById(id2)
	assert.Nil(t, err)
	assert.Len(t, nm.TransportFactoryList, 0, "Should be 0 TransportFactory")
	tf := nm.ConnectNodeToNode(id1, id2)
	assert.Len(t, nm.TransportFactoryList, 1, "Should be 1 TransportFactory")
	assert.True(t, node1.ConnectedTo(node2))
	assert.True(t, node2.ConnectedTo(node1))
	t1, t2 := tf.GetTransports()
	assert.Len(t, node1.Transports, 1, "Error expected 1 transport")
	assert.Len(t, node2.Transports, 1, "Error expected 1 transport")
	assert.Equal(t, t1.Id, t2.StubPair.Id)
	assert.Equal(t, t2.Id, t1.StubPair.Id)
	tr1, err := node1.GetTransportToNode(id2)
	assert.Nil(t, err)
	assert.Equal(t, tr1.StubPair.Id, t2.Id)
}

func TestNetwork(t *testing.T) {
	n := 20
	nm := NewNodeManager()
	nm.CreateNodeList(n)
	assert.Len(t, nm.NodeIdList, n, fmt.Sprintf("Should be %d nodes", n))

	nm.Tick()

	nm.ConnectAll()

	time.Sleep(1 * time.Second)

	node0, err := nm.GetNodeById(nm.NodeIdList[0])
	if err != nil {
		panic(err)
	}

	initRoute := &node.RouteRule{}
	for _, initRoute = range node0.RouteForwardingRules {
		break
	}

	inRouteMessage := messages.InRouteMessage{(messages.TransportId)(0), initRoute.IncomingRoute, []byte{'t', 'e', 's', 't'}}
	serialized := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)
	node0.IncomingChannel <- serialized
	time.Sleep(10 * time.Second)
	for _, tf := range nm.TransportFactoryList {
		t0 := tf.TransportList[0]
		t1 := tf.TransportList[1]
		assert.Equal(t, (uint32)(1), t0.PacketsSent)
		assert.Equal(t, (uint32)(1), t0.PacketsConfirmed)
		assert.Equal(t, (uint32)(0), t1.PacketsSent)
		assert.Equal(t, (uint32)(0), t1.PacketsConfirmed)
	}
}

func TestBuildRoute(t *testing.T) {
	n := 100
	m := 5
	nm := NewNodeManager()
	nm.CreateNodeList(n)

	nodes := []cipher.PubKey{}

	for i := 0; i < m; i++ {
		nodenum := rand.Intn(n)
		nodeId := nm.NodeIdList[nodenum]
		nodes = append(nodes, nodeId)
	}

	for i := 0; i < m-1; i++ {
		nm.ConnectNodeToNode(nodes[i], nodes[i+1])
	}

	nm.Tick()

	routes := nm.BuildRoute(nodes)
	time.Sleep(100 * time.Millisecond)
	assert.Len(t, routes, m, fmt.Sprintf("Should be %d routes", m))
}

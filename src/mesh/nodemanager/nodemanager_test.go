package nodemanager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

func TestAddingNodes(t *testing.T) {
	nm := NewNodeManager()
	assert.Len(t, nm.nodeList, 0, "Error expected 0 nodes")
	nm.AddNewNode()
	assert.Len(t, nm.nodeList, 1, "Error expected 1 nodes")
	nm.AddNewNode()
	nm.AddNewNode()
	nm.AddNewNode()
	nm.AddNewNode()
	assert.Len(t, nm.nodeList, 5, "Error expected 5 nodes")
}

func TestConnectTwoNodes(t *testing.T) {
	nm := NewNodeManager()
	id1 := nm.AddNewNode()
	id2 := nm.AddNewNode()
	node1, err := nm.GetNodeById(id1)
	assert.Nil(t, err)
	node2, err := nm.GetNodeById(id2)
	assert.Nil(t, err)
	assert.Len(t, nm.transportFactoryList, 0, "Should be 0 TransportFactory")
	tf := nm.ConnectNodeToNode(id1, id2)
	assert.Len(t, nm.transportFactoryList, 1, "Should be 1 TransportFactory")
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
	assert.Len(t, nm.nodeIdList, n, fmt.Sprintf("Should be %d nodes", n))

	nm.Tick()
	initRoute, err := nm.ConnectAll()
	assert.Nil(t, err)

	node0, err := nm.GetNodeById(nm.nodeIdList[0])
	if err != nil {
		panic(err)
	}

	fmt.Println(initRoute, node0)
	inRouteMessage := messages.InRouteMessage{messages.NIL_TRANSPORT, initRoute, []byte{'t', 'e', 's', 't'}}
	serialized := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)
	node0.IncomingChannel <- serialized
	time.Sleep(10 * time.Second)
	for _, tf := range nm.transportFactoryList {
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
		nodeId := nm.nodeIdList[nodenum]
		nodes = append(nodes, nodeId)
	}

	for i := 0; i < m-1; i++ {
		nm.ConnectNodeToNode(nodes[i], nodes[i+1])
	}

	nm.Tick()

	routes, err := nm.buildRoute(nodes)
	assert.Nil(t, err)
	assert.Len(t, routes, m, fmt.Sprintf("Should be %d routes", m))
}

func TestFindRoute(t *testing.T) {
	nm := NewNodeManager()
	nodeList := nm.CreateNodeList(10)
	/*
		  1-2-3-4   long route
		 /	 \
		0---5-----9 short route, which should be selected
		 \ /     /
		  6_7_8_/   medium route
	*/
	nm.Tick()
	nm.ConnectNodeToNode(nodeList[0], nodeList[1]) // making long route
	nm.ConnectNodeToNode(nodeList[1], nodeList[2])
	nm.ConnectNodeToNode(nodeList[2], nodeList[3])
	nm.ConnectNodeToNode(nodeList[3], nodeList[4])
	nm.ConnectNodeToNode(nodeList[4], nodeList[9])
	nm.ConnectNodeToNode(nodeList[0], nodeList[5]) // making short route
	nm.ConnectNodeToNode(nodeList[5], nodeList[9])
	nm.ConnectNodeToNode(nodeList[0], nodeList[6]) // make medium route, then findRoute should select the short one
	nm.ConnectNodeToNode(nodeList[6], nodeList[7])
	nm.ConnectNodeToNode(nodeList[7], nodeList[8])
	nm.ConnectNodeToNode(nodeList[8], nodeList[9])
	nm.ConnectNodeToNode(nodeList[5], nodeList[6]) // just for

	nm.RebuildRoutes()

	nodeFrom, nodeTo := nodeList[0], nodeList[9]
	nodes, found := nm.routeGraph.findRoute(nodeFrom, nodeTo)
	assert.True(t, found)
	assert.Len(t, nodes, 3, "Should be 3 nodes")
}

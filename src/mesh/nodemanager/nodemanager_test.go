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
	messages.SetDebugLogLevel()
	nm := newNodeManager()
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
	messages.SetDebugLogLevel()
	nm := newNodeManager()
	id1 := nm.AddNewNode()
	id2 := nm.AddNewNode()
	node1, err := nm.getNodeById(id1)
	assert.Nil(t, err)
	node2, err := nm.getNodeById(id2)
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
	messages.SetDebugLogLevel()
	n := 20
	nm := newNodeManager()
	nm.createNodeList(n)
	assert.Len(t, nm.nodeIdList, n, fmt.Sprintf("Should be %d nodes", n))

	initRoute, err := nm.connectAllAndBuildRoute()
	assert.Nil(t, err)

	node0, err := nm.getNodeById(nm.nodeIdList[0])
	if err != nil {
		panic(err)
	}

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
	messages.SetInfoLogLevel()
	n := 100
	m := 5
	nm := newNodeManager()
	nm.createNodeList(n)

	nodes := []cipher.PubKey{}

	for i := 0; i < m; i++ {
		nodenum := rand.Intn(n)
		nodeId := nm.nodeIdList[nodenum]
		nodes = append(nodes, nodeId)
	}

	for i := 0; i < m-1; i++ {
		nm.ConnectNodeToNode(nodes[i], nodes[i+1])
	}

	routes, err := nm.buildRouteOneSide(nodes)
	assert.Nil(t, err)
	assert.Len(t, routes, m, fmt.Sprintf("Should be %d routes", m))
}

func TestFindRoute(t *testing.T) {
	messages.SetDebugLogLevel()
	nm := newNodeManager()
	nodeList := nm.createNodeList(10)
	/*
		  1-2-3-4   long route
		 /	 \
		0---5-----9 short route, which should be selected
		 \ /     /
		  6_7_8_/   medium route
	*/
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

	nm.rebuildRoutes()

	nodeFrom, nodeTo := nodeList[0], nodeList[9]
	nodes, found := nm.routeGraph.findRoute(nodeFrom, nodeTo)
	assert.True(t, found)
	assert.Len(t, nodes, 3, "Should be 3 nodes")
}

func TestConnection(t *testing.T) {
	messages.SetDebugLogLevel()
	n := 4
	nm := newNodeManager()

	nodes := nm.createNodeList(n)
	nm.connectAll()

	node0 := nodes[0]
	route, backRoute, err := nm.buildRoute(nodes)
	assert.Nil(t, err)
	conn0, err := nm.NewConnectionWithRoutes(node0, route, backRoute)
	assert.Nil(t, err)
	payload := []byte{'t', 'e', 's', 't'}
	msg := messages.RequestMessage{
		0,
		backRoute,
		payload,
	}
	msgS := messages.Serialize(messages.MsgRequestMessage, msg)
	sequence, err := conn0.Send(msgS)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), sequence)
	time.Sleep(time.Duration(n) * time.Second)
}

func TestAddAndConnect2Nodes(t *testing.T) {
	messages.SetDebugLogLevel()
	nm := newNodeManager()

	pubkey0 := nm.AddAndConnect()
	pubkey1 := nm.AddAndConnect()

	assert.Len(t, nm.nodeIdList, 2)
	assert.True(t, nm.connected(pubkey0, pubkey1))
}

func TestRandomNetwork100Nodes(t *testing.T) {
	messages.SetInfoLogLevel()
	n := 100
	nm := newNodeManager()

	nodes := nm.CreateRandomNetwork(n)

	assert.Len(t, nm.nodeIdList, n)
	assert.Equal(t, nm.nodeIdList, nodes)
	assert.True(t, nm.routeExists(nodes[0], nodes[n-1]))
}

func TestSendThroughRandomNetworks(t *testing.T) {
	messages.SetDebugLogLevel()
	nm := newNodeManager()
	lens := []int{2, 5, 10} // sizes of different networks which will be tested

	for _, n := range lens {
		nodes := nm.CreateRandomNetwork(n)

		node0 := nodes[0]
		node1 := nodes[len(nodes)-1]
		conn0, err := nm.NewConnection(node0, node1)
		assert.Nil(t, err)
		msg := []byte{'t', 'e', 's', 't'}
		sequence, err := conn0.Send(msg)
		assert.Nil(t, err)
		assert.Equal(t, uint32(0), sequence)
		time.Sleep(time.Duration(n) * time.Second)
	}
}

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
	fmt.Println("TestAddingNodes")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	assert.Len(t, nm.nodeList, 0, "Error expected 0 nodes")
	nm.AddNewNodeStub()
	assert.Len(t, nm.nodeList, 1, "Error expected 1 nodes")
	nm.AddNewNodeStub()
	nm.AddNewNodeStub()
	nm.AddNewNodeStub()
	nm.AddNewNodeStub()
	assert.Len(t, nm.nodeList, 5, "Error expected 5 nodes")
	fmt.Println("TestAddingNodes end")
}

func TestConnectTwoNodes(t *testing.T) {
	fmt.Println("TestConnectTwoNodes")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	id1 := nm.AddNewNodeStub()
	id2 := nm.AddNewNodeStub()
	node1, err := nm.getNodeById(id1)
	assert.Nil(t, err)
	node2, err := nm.getNodeById(id2)
	assert.Nil(t, err)
	assert.Len(t, nm.transportFactoryList, 0, "Should be 0 TransportFactory")
	tf, err := nm.ConnectNodeToNode(id1, id2)
	assert.Nil(t, err)
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
	fmt.Println("TestConnectTwoNodes end")
}

func TestNetwork(t *testing.T) {
	fmt.Println("TestNetwork")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	n := 20

	nm.createNodeList(n)
	assert.Len(t, nm.nodeIdList, n, fmt.Sprintf("Should be %d nodes", n))

	initRoute, err := nm.connectAllAndBuildRoute()
	assert.Nil(t, err)

	node0, err := nm.getNodeById(nm.nodeIdList[0])
	if err != nil {
		panic(err)
	}

	inRouteMessage := messages.InRouteMessage{messages.NIL_TRANSPORT, initRoute, []byte{'t', 'e', 's', 't'}}
	//	serialized := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)
	node0.InjectTransportMessage(&inRouteMessage)
	time.Sleep(10 * time.Second)
	for _, tf := range nm.transportFactoryList {
		t0 := tf.TransportList[0]
		t1 := tf.TransportList[1]
		assert.Equal(t, (uint32)(1), t0.PacketsSent)
		assert.Equal(t, (uint32)(1), t0.PacketsConfirmed)
		assert.Equal(t, (uint32)(0), t1.PacketsSent)
		assert.Equal(t, (uint32)(0), t1.PacketsConfirmed)
	}
	fmt.Println("TestNetwork end")
}

func TestBuildRoute(t *testing.T) {
	fmt.Println("TestBuildRoute")
	messages.SetInfoLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	n := 100
	m := 5

	nm.createNodeList(n)

	nodes := []cipher.PubKey{}

	for i := 0; i < m; i++ {
		nodenum := rand.Intn(n)
		nodeId := nm.nodeIdList[nodenum]
		nodes = append(nodes, nodeId)
	}

	for i := 0; i < m-1; i++ {
		_, err := nm.ConnectNodeToNode(nodes[i], nodes[i+1])
		assert.Nil(t, err)
	}

	routes, err := nm.buildRouteForward(nodes)
	assert.Nil(t, err)
	assert.Len(t, routes, m, fmt.Sprintf("Should be %d routes", m))
	fmt.Println("TestBuildRoute end")
}

func TestFindRoute(t *testing.T) {
	fmt.Println("TestFindRoute")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

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
	routes, err := nm.findRouteForward(nodeFrom, nodeTo)
	assert.Nil(t, err)
	assert.Len(t, routes, 3, "Should be 3 routes")
	fmt.Println("TestFindRoute end")
}

/*
func TestConnection(t *testing.T) {
	fmt.Println("TestConnection")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	n := 4

	nodes := nm.createNodeList(n)
	nm.connectAll()

	node0 := nodes[0]
	time.Sleep(2 * time.Second)
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
	fmt.Println("TestConnection end")
}
*/
func TestAddAndConnect2Nodes(t *testing.T) {
	fmt.Println("TestAddAndConnect")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	pubkey0 := nm.AddAndConnectStub()
	pubkey1 := nm.AddAndConnectStub()

	assert.Len(t, nm.nodeIdList, 2)
	assert.True(t, nm.connected(pubkey0, pubkey1))
	fmt.Println("TestAddAndConnect end")
}

func TestRandomNetwork100Nodes(t *testing.T) {
	fmt.Println("TestRandomNetwork100Nodes")
	messages.SetInfoLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	n := 100

	nodes := nm.CreateRandomNetwork(n)

	assert.Len(t, nm.nodeIdList, n)
	assert.Equal(t, nm.nodeIdList, nodes)
	assert.True(t, nm.routeExists(nodes[0], nodes[n-1]))
	fmt.Println("TestRandomNetwork100Nodes end")
}

/*
func TestSendThroughRandomNetworks(t *testing.T) {
	fmt.Println("TestSendThroughRandomNetworks")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

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
	fmt.Println("TestSendThroughRandomNetworks end")
}
*/

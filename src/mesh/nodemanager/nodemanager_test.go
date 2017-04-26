package nodemanager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
)

func TestMessagingServer(t *testing.T) {
	return
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	msgSrv := nm.msgServer
	assert.NotNil(t, msgSrv)

	config := messages.GetConfig()

	host := net.ParseIP(config.MsgSrvHost)
	port := int(config.MsgSrvPort)
	msgSrvAddr := &net.UDPAddr{IP: host, Port: port}
	assert.Equal(t, *msgSrvAddr, *msgSrv.addr)
}

func TestRegisterNode(t *testing.T) {
	return

	nm := newNodeManager()
	defer nm.Shutdown()

	assert.Len(t, nm.nodeList, 0)

	n, err := node.CreateNode(messages.LOCALHOST+":5992", messages.LOCALHOST+":5999")
	assert.Nil(t, err)
	defer n.Shutdown()

	assert.Len(t, nm.nodeList, 1)
	assert.Equal(t, n.Id(), nm.nodeIdList[0])
}

func TestConnectNodes(t *testing.T) {
	return
	fmt.Println("")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	n0, err := node.CreateNode(messages.LOCALHOST+":5992", messages.LOCALHOST+":5999")
	assert.Nil(t, err)
	defer n0.Shutdown()

	n1, err := node.CreateNode(messages.LOCALHOST+":5993", messages.LOCALHOST+":5999")
	assert.Nil(t, err)
	defer n1.Shutdown()

	assert.Len(t, nm.nodeList, 2)

	_, err = nm.ConnectNodeToNode(n0.Id(), n1.Id())
	assert.Nil(t, err)

	assert.True(t, n0.ConnectedTo(n1.Id()))
	assert.True(t, n1.ConnectedTo(n0.Id()))

	tf := nm.transportFactoryList[0]
	t0, t1 := tf.getTransports()
	assert.Equal(t, t0.id, t1.pair.id)
	assert.Equal(t, t1.id, t0.pair.id)
}

func TestNetwork(t *testing.T) {
	return
	fmt.Println("TestNetwork")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	q := 20

	nodes := node.CreateNodeList(q)
	assert.Len(t, nodes, q, fmt.Sprintf("Should be %d nodes", q))
	assert.Len(t, nm.nodeIdList, q, fmt.Sprintf("Should be %d nodes", q))
	initRoute, err := nm.connectAllAndBuildRoute()
	assert.Nil(t, err)

	node0 := nodes[0]

	inRouteMessage := messages.InRouteMessage{messages.NIL_TRANSPORT, initRoute, []byte{'t', 'e', 's', 't'}}
	node0.InjectTransportMessage(&inRouteMessage)
	time.Sleep(10 * time.Second)
	for i := 0; i < q-1; i++ {
		n0 := nodes[i]
		n1 := nodes[i+1]
		t0, err := n0.GetTransportToNode(n1.Id())
		assert.Nil(t, err)
		t1, err := n1.GetTransportToNode(n0.Id())
		assert.Nil(t, err)
		assert.Equal(t, uint32(1), t0.PacketsSent())
		assert.Equal(t, uint32(1), t0.PacketsConfirmed())
		assert.Equal(t, uint32(0), t1.PacketsSent())
		assert.Equal(t, uint32(0), t1.PacketsConfirmed())
	}

	node.ShutdownAll(nodes)

	fmt.Println("TestNetwork end")
}

func TestBuildRoute(t *testing.T) {
	return
	fmt.Println("TestBuildRoute")
	messages.SetInfoLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	n := 100
	m := 5

	allNodes := node.CreateNodeList(n)

	nodes := []cipher.PubKey{}

	for i := 0; i < m; i++ {
		nodenum := rand.Intn(n)
		node := allNodes[nodenum]
		nodes = append(nodes, node.Id())
	}

	for i := 0; i < m-1; i++ {
		_, err := nm.ConnectNodeToNode(nodes[i], nodes[i+1])
		assert.Nil(t, err)
	}

	routes, err := nm.buildRouteForward(nodes)
	assert.Nil(t, err)
	assert.Len(t, routes, m)

	node.ShutdownAll(allNodes)
	fmt.Println("TestBuildRoute end")
}

func TestFindRoute(t *testing.T) {
	return
	fmt.Println("TestFindRoute")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	nodes := node.CreateNodeList(10)

	nodeList := []cipher.PubKey{}
	for _, n := range nodes {
		nodeList = append(nodeList, n.Id())
	}

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

	node.ShutdownAll(nodes)
	fmt.Println("TestFindRoute end")
}

func TestAddAndConnect2Nodes(t *testing.T) {
	return
	fmt.Println("TestAddAndConnect")
	messages.SetDebugLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	n0, err := node.CreateAndConnectNode(messages.LOCALHOST+":5992", messages.LOCALHOST+":5999")
	assert.Nil(t, err)
	defer n0.Shutdown()

	n1, err := node.CreateAndConnectNode(messages.LOCALHOST+":5993", messages.LOCALHOST+":5999")
	assert.Nil(t, err)
	defer n1.Shutdown()

	assert.Len(t, nm.nodeIdList, 2)
	assert.True(t, nm.connected(n0.Id(), n1.Id()))

	fmt.Println("TestAddAndConnect end")
}

func TestRandomNetwork100Nodes(t *testing.T) {
	return
	fmt.Println("TestRandomNetwork100Nodes")
	messages.SetInfoLogLevel()

	nm := newNodeManager()
	defer nm.Shutdown()

	n := 100

	nodes := nm.CreateRandomNetwork(n)

	nodeIds := []cipher.PubKey{}

	for _, node := range nodes {
		nodeIds = append(nodeIds, node.Id())
	}

	assert.Len(t, nm.nodeIdList, n)
	assert.Equal(t, nm.nodeIdList, nodeIds)
	assert.True(t, nm.routeExists(nodeIds[0], nodeIds[n-1]))

	node.ShutdownAll(nodes)
	fmt.Println("TestRandomNetwork100Nodes end")
}

func TestSendThroughRandomNetworks(t *testing.T) {
	fmt.Println("TestSendThroughRandomNetworks")
	messages.SetDebugLogLevel()

	lens := []int{2, 5, 10} // sizes of different networks which will be tested

	for _, n := range lens {

		nm := newNodeManager()

		nodes := nm.CreateRandomNetwork(n)

		n0 := nodes[0]
		n1 := nodes[len(nodes)-1]
		conn0, err := n0.Dial(n1.Id(), messages.AppId([]byte{}), messages.AppId([]byte{}))
		connId := conn0.Id()
		if err != nil {
			panic(err)
		}
		conn1 := n1.GetConnection(connId)
		assert.Equal(t, conn0.Status(), CONNECTED)
		assert.Equal(t, conn1.Status(), CONNECTED)
		fmt.Println(conn0.Id(), conn1.Id())
		msg := []byte{'t', 'e', 's', 't'}
		err = conn0.Send(msg)
		assert.Nil(t, err)
		time.Sleep(time.Duration(n) * time.Second)

		node.ShutdownAll(nodes)
		nm.Shutdown()
		time.Sleep(time.Duration(n) * time.Millisecond)
	}
	fmt.Println("TestSendThroughRandomNetworks end")
}

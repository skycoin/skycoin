package node_manager

import (
	"fmt"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/node"
	"github.com/stretchr/testify/assert"
)

func TestAddingNodes(t *testing.T) {

	nm := NewNodeManager()
	assert.Len(t, nm.NodeList.nodes, 0, "Error expected 0 nodes")
	nm.AddNode()
	assert.Len(t, nm.NodeList.nodes, 1, "Error expected 1 nodes")
	nm.AddNode()
	nm.AddNode()
	nm.AddNode()
	nm.AddNode()
	assert.Len(t, nm.NodeList.nodes, 5, "Error expected 5 nodes")
}

func TestConnectTwoNodes(t *testing.T) {

	nm := NewNodeManager()
	id1 := nm.AddNode()
	id2 := nm.AddNode()
	node1, err := nm.GetNodeById(id1)
	assert.Nil(t, err)
	node2, err := nm.GetNodeById(id2)
	assert.Nil(t, err)
	tid1, tid2 := nm.ConnectNodeToNode(id1, id2)
	assert.Len(t, node1.Transports, 1, "Error expected 1 transport")
	assert.Len(t, node2.Transports, 1, "Error expected 1 transport")
	assert.Equal(t, node1.Transports[tid1].Id, node2.Transports[tid2].StubPair.Id)
	assert.Equal(t, node2.Transports[tid2].Id, node1.Transports[tid1].StubPair.Id)
	tr1, err := node1.GetTransportToNode(id2)
	assert.Nil(t, err)
	assert.Equal(t, tr1.StubPair.Id, tid2)
}

func TestSendMessage(t *testing.T) {

	nm := NewNodeManager()
	id1 := nm.AddNode()
	id2 := nm.AddNode()
	node1, err := nm.GetNodeById(id1)
	assert.Nil(t, err)
	node2, err := nm.GetNodeById(id2)
	assert.Nil(t, err)
	tid1, tid2 := nm.ConnectNodeToNode(id1, id2)
	routeId1 := messages.RandRouteId()
	routeId2 := messages.RandRouteId()
	route1 := node.RouteRule{tid1, tid1, routeId1, routeId2}
	route2 := node.RouteRule{tid2, tid2, routeId2, routeId1} //endless cycle
	node1.RouteForwardingRules[routeId1] = &route1
	node2.RouteForwardingRules[routeId2] = &route2
	nm.Tick()
	time.Sleep(1 * time.Second)

	for tid1 := range node1.Transports {
		fmt.Println("node1 tr", tid1)
	}

	for rid1, routeRule1 := range node1.RouteForwardingRules {
		fmt.Println("node1 rr", rid1, routeRule1)
	}

	for tid2 := range node2.Transports {
		fmt.Println("node2 tr", tid2)
	}

	for rid2, routeRule2 := range node2.RouteForwardingRules {
		fmt.Println("node2 rr", rid2, routeRule2)
	}

	inRouteMessage := messages.InRouteMessage{tid1, routeId1, []byte{'t', 'e', 's', 't'}}
	serialized := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)
	node1.IncomingChannel <- serialized
	time.Sleep(1 * time.Second)
}

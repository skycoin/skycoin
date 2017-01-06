package node

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/transport"
)

func TestCreateControlChannel(t *testing.T) {
	node := NewNode()
	assert.Len(t, node.controlChannels, 0, "Should be 0 control channels")
	node.AddControlChannel()
	assert.Len(t, node.controlChannels, 1, "Should be 1 control channels")
	fmt.Println("--------------------\n")
}

func TestRemoveControlChannel(t *testing.T) {
	node := NewNode()
	ccid := node.AddControlChannel()
	assert.Len(t, node.controlChannels, 1, "Should be 1 control channels")
	node.CloseControlChannel(ccid)
	assert.Len(t, node.controlChannels, 0, "Should be 0 control channels")
	fmt.Println("--------------------\n")
}

func TestAddRoute(t *testing.T) {
	node1, node2 := NewNode(), NewNode()
	assert.Len(t, node1.RouteForwardingRules, 0, "Should be 0 routes")

	ccid := node1.AddControlChannel()
	node1.Tick() // run channels consuming

	tf := transport.NewTransportFactory()
	tf.ConnectNodeToNode(node1, node2)
	tr1, _ := tf.GetTransports()
	tid1 := tr1.Id

	routeId := messages.RandRouteId()

	msg := messages.AddRouteControlMessage{node2.Id, routeId}
	msgS := messages.Serialize(messages.MsgAddRouteControlMessage, msg)

	controlMessage := messages.InControlMessage{ccid, msgS}

	node1.InjectControlMessage(controlMessage)
	time.Sleep(1 * time.Millisecond)

	assert.Len(t, node1.RouteForwardingRules, 1, "Should be 1 routes")
	assert.Equal(t, node1.RouteForwardingRules[routeId].OutgoingRoute, routeId)
	assert.Equal(t, node1.RouteForwardingRules[routeId].OutgoingTransport, tid1)

	fmt.Println("--------------------\n")
}

func TestRemoveRoute(t *testing.T) {
	node1, node2 := NewNode(), NewNode()

	ccid := node1.AddControlChannel()
	node1.Tick() // run channels consuming

	tf := transport.NewTransportFactory()
	tf.ConnectNodeToNode(node1, node2)

	routeId := messages.RandRouteId()

	msg := messages.AddRouteControlMessage{node2.Id, routeId}
	msgS := messages.Serialize(messages.MsgAddRouteControlMessage, msg)

	controlMessage := messages.InControlMessage{ccid, msgS}

	node1.InjectControlMessage(controlMessage)

	msg2 := messages.RemoveRouteControlMessage{routeId}
	msgS2 := messages.Serialize(messages.MsgRemoveRouteControlMessage, msg2)

	controlMessage = messages.InControlMessage{ccid, msgS2}
	node1.InjectControlMessage(controlMessage)
	time.Sleep(1 * time.Millisecond)

	assert.Len(t, node1.RouteForwardingRules, 0, "Should be 0 routes")

	fmt.Println("--------------------\n")
}

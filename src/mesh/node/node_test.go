package node

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

func TestCreateControlChannel(t *testing.T) {
	messages.SetDebugLogLevel()
	node := NewNode()
	assert.Len(t, node.controlChannels, 0, "Should be 0 control channels")
	node.AddControlChannel()
	assert.Len(t, node.controlChannels, 1, "Should be 1 control channels")
	fmt.Println("--------------------\n")
}

func TestRemoveControlChannel(t *testing.T) {
	messages.SetDebugLogLevel()
	node := NewNode()
	ccid := node.AddControlChannel()
	assert.Len(t, node.controlChannels, 1, "Should be 1 control channels")
	node.CloseControlChannel(ccid)
	assert.Len(t, node.controlChannels, 0, "Should be 0 control channels")
	fmt.Println("--------------------\n")
}

func TestAddRoute(t *testing.T) {
	messages.SetDebugLogLevel()
	node1, node2 := NewNode(), NewNode()
	assert.Len(t, node1.RouteForwardingRules, 0, "Should be 0 routes")

	ccid := node1.AddControlChannel()

	tf := transport.NewTransportFactory()
	tf.ConnectNodeToNode(node1, node2)
	tr1, _ := tf.GetTransports()
	tid1 := tr1.Id

	incomingRouteId := messages.RandRouteId()
	outgoingRouteId := messages.RandRouteId()

	msg := messages.AddRouteControlMessage{
		messages.NIL_TRANSPORT,
		tr1.Id,
		incomingRouteId,
		outgoingRouteId,
	}
	msgS := messages.Serialize(messages.MsgAddRouteControlMessage, msg)

	controlMessage := messages.InControlMessage{ccid, msgS, nil}

	node1.InjectControlMessage(controlMessage)

	assert.Len(t, node1.RouteForwardingRules, 1, "Should be 1 routes")
	assert.Equal(t, node1.RouteForwardingRules[incomingRouteId].IncomingRoute, incomingRouteId)
	assert.Equal(t, node1.RouteForwardingRules[incomingRouteId].OutgoingRoute, outgoingRouteId)
	assert.Equal(t, node1.RouteForwardingRules[incomingRouteId].OutgoingTransport, tid1)

	fmt.Println("--------------------\n")
}

func TestRemoveRoute(t *testing.T) {
	messages.SetDebugLogLevel()
	node1, node2 := NewNode(), NewNode()

	ccid := node1.AddControlChannel()

	tf := transport.NewTransportFactory()
	tf.ConnectNodeToNode(node1, node2)

	routeId := messages.RandRouteId()

	msg := messages.AddRouteControlMessage{}
	msg.IncomingRouteId = routeId
	msgS := messages.Serialize(messages.MsgAddRouteControlMessage, msg)

	controlMessage := messages.InControlMessage{ccid, msgS, nil}

	node1.InjectControlMessage(controlMessage)

	msg2 := messages.RemoveRouteControlMessage{routeId}
	msgS2 := messages.Serialize(messages.MsgRemoveRouteControlMessage, msg2)

	controlMessage = messages.InControlMessage{ccid, msgS2, nil}
	node1.InjectControlMessage(controlMessage)

	assert.Len(t, node1.RouteForwardingRules, 0, "Should be 0 routes")

	fmt.Println("--------------------\n")
}

package node

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

func TestCreateNode(t *testing.T) {
	node := newLocalNode()
	assert.Len(t, node.controlChannels, 1)
	assert.Equal(t, cap(node.incomingControlChannel), 256)
	assert.Equal(t, cap(node.congestionChannel), 1024)
	assert.Equal(t, node.host, messages.LOCALHOST)
	assert.NotNil(t, node.lock)
}

func TestCreateControlChannel(t *testing.T) {
	messages.SetDebugLogLevel()
	node := newLocalNode()
	assert.Len(t, node.controlChannels, 1, "Should be 1 control channels")
	node.addControlChannel()
	assert.Len(t, node.controlChannels, 2, "Should be 2 control channels")
	fmt.Println("--------------------\n")
}

func TestRemoveControlChannel(t *testing.T) {
	messages.SetDebugLogLevel()
	node := newLocalNode()
	ccid := node.addControlChannel()
	assert.Len(t, node.controlChannels, 2, "Should be 2 control channels")
	node.closeControlChannel(ccid)
	assert.Len(t, node.controlChannels, 1, "Should be 1 control channels")
	fmt.Println("--------------------\n")
}

func TestAddRoute(t *testing.T) {
	messages.SetDebugLogLevel()
	node1 := newLocalNode()
	assert.Len(t, node1.routeForwardingRules, 0, "Should be 0 routes")

	incomingRouteId := messages.RandRouteId()
	outgoingRouteId := messages.RandRouteId()
	trId := messages.RandTransportId()

	msg := messages.AddRouteCM{
		messages.NIL_TRANSPORT,
		trId,
		incomingRouteId,
		outgoingRouteId,
	}
	msgS := messages.Serialize(messages.MsgAddRouteCM, msg)

	controlMessage := &messages.InControlMessage{messages.ChannelId(0), 0, msgS}

	node1.injectControlMessage(controlMessage)

	assert.Len(t, node1.routeForwardingRules, 1)
	assert.Equal(t, node1.routeForwardingRules[incomingRouteId].IncomingRoute, incomingRouteId)
	assert.Equal(t, node1.routeForwardingRules[incomingRouteId].OutgoingRoute, outgoingRouteId)
	assert.Equal(t, node1.routeForwardingRules[incomingRouteId].OutgoingTransport, trId)

	fmt.Println("--------------------\n")
}

func TestRemoveRoute(t *testing.T) {
	messages.SetDebugLogLevel()
	node1 := newLocalNode()
	assert.Len(t, node1.routeForwardingRules, 0, "Should be 0 routes")

	incomingRouteId := messages.RandRouteId()
	outgoingRouteId := messages.RandRouteId()
	trId := messages.RandTransportId()

	msg := messages.AddRouteCM{
		messages.NIL_TRANSPORT,
		trId,
		incomingRouteId,
		outgoingRouteId,
	}
	msgS := messages.Serialize(messages.MsgAddRouteCM, msg)

	controlMessage := &messages.InControlMessage{messages.ChannelId(0), 0, msgS}

	node1.injectControlMessage(controlMessage)

	msg0 := messages.RemoveRouteCM{
		incomingRouteId,
	}
	msgS = messages.Serialize(messages.MsgRemoveRouteCM, msg0)

	controlMessage = &messages.InControlMessage{messages.ChannelId(0), 0, msgS}

	node1.injectControlMessage(controlMessage)

	assert.Len(t, node1.routeForwardingRules, 0)

	fmt.Println("--------------------\n")
}

func TestRegisterAckAccept(t *testing.T) {
	messages.SetDebugLogLevel()

	node := newLocalNode()
	pubKey, _ := cipher.GenerateKeyPair()
	registerAck := &messages.RegisterNodeCMAck{
		NodeId:            pubKey,
		TimeUnit:          1000,
		MaxBuffer:         512,
		MaxPacketSize:     1024,
		SendInterval:      10,
		ConnectionTimeout: 10000,
	}
	node.register(registerAck)
	assert.Equal(t, node.id, pubKey)
	assert.Equal(t, node.timeUnit, 1000*time.Microsecond)
	assert.Equal(t, node.maxPacketSize, uint32(1024))
	assert.Equal(t, node.maxBuffer, uint64(512))
}

func TestConnectionCreate(t *testing.T) {
	messages.SetDebugLogLevel()

	node := newLocalNode()
	pubKey, _ := cipher.GenerateKeyPair()
	node.id = pubKey

	registerAck := &messages.RegisterNodeCMAck{
		NodeId:            pubKey,
		SendInterval:      10,
		ConnectionTimeout: 10000,
	}
	node.register(registerAck)

	assert.Len(t, node.connections, 0)

	conn, err := node.newConnection(messages.RandConnectionId(), messages.RandRouteId(), messages.AppId([]byte{}))
	assert.Nil(t, err)
	assert.Len(t, node.connections, 1)

	assert.Equal(t, pubKey, conn.nodeAttached.id)
	assert.Equal(t, CONNECTING, conn.status)
	assert.NotNil(t, conn.lock)
	assert.NotNil(t, conn.errChan)
	assert.Len(t, conn.ackChannels, 0)
	assert.Len(t, conn.incomingMessages, 0)
	assert.Len(t, conn.incomingCounter, 0)
	assert.Equal(t, conn.timeout, 10000*time.Millisecond)
	assert.Equal(t, conn.sendInterval, 10*time.Microsecond)
}

func TestConnectionMessage(t *testing.T) {
	messages.SetDebugLogLevel()

	node := newLocalNode()
	assert.Equal(t, uint32(0), node.ticks)
	registerAck := &messages.RegisterNodeCMAck{
		MaxBuffer:         512,
		ConnectionTimeout: 10000,
	}
	node.register(registerAck)

	conn, err := node.newConnection(messages.RandConnectionId(), messages.RandRouteId(), messages.AppId([]byte{}))
	assert.Nil(t, err)

	inRouteMessage := messages.InRouteMessage{}
	conn.sendToNode(&inRouteMessage)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, uint32(1), node.ticks)
}

func TestTransportCreate(t *testing.T) {
	messages.SetDebugLogLevel()

	node := newLocalNode()
	assert.Len(t, node.transports, 0)

	pubKey, _ := cipher.GenerateKeyPair()
	node.id = pubKey

	trId := messages.RandTransportId()
	transportCreateMessage := &messages.TransportCreateCM{
		Id: trId,
	}
	node.setTransportFromMessage(transportCreateMessage)
	assert.Len(t, node.transports, 1)
	assert.Equal(t, trId, node.transports[trId].Id())
	assert.Equal(t, pubKey, node.transports[trId].AttachedNode.(*Node).Id())
}

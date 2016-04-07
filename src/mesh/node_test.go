package mesh

import(
	"time"
	"reflect"
	"testing"
	//"fmt"
	)

import (
	"github.com/stretchr/testify/assert"
    "github.com/skycoin/skycoin/src/cipher"
)

var test_key1 = cipher.NewPubKey([]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
var test_key2 = cipher.NewPubKey([]byte{2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
var test_key3 = cipher.NewPubKey([]byte{3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
var test_key4 = cipher.NewPubKey([]byte{4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

func init() {
}

func TestEstablishRoutes(t *testing.T) {
	test_route := RouteConfig{[]cipher.PubKey{test_key2, test_key3, test_key4}}
	established_ch := make(chan EstablishedRoute)
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{test_key2},
		1024,
		1024,
		[]RouteConfig{test_route},
		time.Hour,
		time.Hour,
		// RouteEstablishedCB
		func(route EstablishedRoute) {
			established_ch <- route
		},
	}
	node := NewNode(test_config)

	go node.Run()

	// Establish 2 -> 3
	{
		outgoing := <- node.MessagesOut
		assert.Equal(t, outgoing.ConnectedPeerPubKey, test_key2)
		msg := outgoing.Message
		assert.Equal(t, reflect.TypeOf(msg), reflect.TypeOf(EstablishRouteMessage{}))
		establish_msg := msg.(EstablishRouteMessage)
		assert.Equal(t, establish_msg.ToPubKey, test_key3)

		// Reply
		node.MessagesIn <- EstablishRouteReplyMessage{
			OperationReply{OperationMessage{Message{0}, establish_msg.MsgId}, true, ""}, 
			55, "secret_abc"}
	}
	// Establish 3 -> 4
	{
		outgoing := <- node.MessagesOut
		assert.Equal(t, outgoing.ConnectedPeerPubKey, test_key2)
		msg := outgoing.Message
		assert.Equal(t, reflect.TypeOf(msg), reflect.TypeOf(EstablishRouteMessage{}))
		establish_msg := msg.(EstablishRouteMessage)
		assert.Equal(t, establish_msg.ToPubKey, test_key4)

		// Reply
		node.MessagesIn <- EstablishRouteReplyMessage{
			OperationReply{OperationMessage{Message{55}, establish_msg.MsgId}, true, ""}, 
			120, "secret_xyz"}
	}
	// Set rewrite 2 -> 3
	{
		outgoing := <- node.MessagesOut
		assert.Equal(t, outgoing.ConnectedPeerPubKey, test_key2)
		msg := outgoing.Message
		assert.Equal(t, reflect.TypeOf(msg), reflect.TypeOf(RouteRewriteMessage{}))
		rewrite_msg := msg.(RouteRewriteMessage)
		assert.Equal(t, rewrite_msg.Secret, "secret_abc")
		assert.Equal(t, rewrite_msg.RewriteSendId, uint32(120))

		// Reply
		node.MessagesIn <- OperationReply{OperationMessage{Message{0}, rewrite_msg.MsgId}, true, ""}
	}
	route_established := <- established_ch

	assert.Equal(t, route_established.SendId, uint32(55))
	assert.Equal(t, route_established.ConnectedPeer, test_key2)
}

// Establish route test, rewrite test

// Retransmit, stop retransmitting
// Timeout
// Out of channel route
// Backward send


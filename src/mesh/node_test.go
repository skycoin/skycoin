package mesh

import(
	"time"
	"reflect"
	"testing"
	//"fmt"
	)

import (
	"github.com/stretchr/testify/assert"
    "github.com/satori/go.uuid"
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
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg := outgoing.Message
		assert.Equal(t, reflect.TypeOf(EstablishRouteMessage{}), reflect.TypeOf(msg))
		establish_msg := msg.(EstablishRouteMessage)
		assert.Equal(t, test_key3, establish_msg.ToPubKey)

		// Reply
		node.MessagesIn <- PhysicalMessage{
			test_key2,
			EstablishRouteReplyMessage{
				OperationReply{OperationMessage{Message{0, true}, establish_msg.MsgId}, true, ""}, 
				55, "secret_abc"}}
	}
	// Establish 3 -> 4
	{
		outgoing := <- node.MessagesOut
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg := outgoing.Message
		assert.Equal(t, reflect.TypeOf(EstablishRouteMessage{}), reflect.TypeOf(msg))
		establish_msg := msg.(EstablishRouteMessage)
		assert.Equal(t, test_key4, establish_msg.ToPubKey)

		// Reply
		node.MessagesIn <- PhysicalMessage{
			test_key3,
			EstablishRouteReplyMessage{
				OperationReply{OperationMessage{Message{55, true}, establish_msg.MsgId}, true, ""}, 
				120, "secret_xyz"}}
	}
	// Set rewrite 2 -> 3
	{
		outgoing := <- node.MessagesOut
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg := outgoing.Message
		assert.Equal(t, reflect.TypeOf(RouteRewriteMessage{}), reflect.TypeOf(msg))
		rewrite_msg := msg.(RouteRewriteMessage)
		assert.Equal(t, "secret_abc", rewrite_msg.Secret)
		assert.Equal(t, uint32(120), rewrite_msg.RewriteSendId)

		// Reply
		node.MessagesIn <-  PhysicalMessage{
			test_key2,
			OperationReply{OperationMessage{Message{0, true}, rewrite_msg.MsgId}, true, ""}}
	}
	route_established := <- established_ch

	assert.Equal(t, uint32(55), route_established.SendId)
	assert.Equal(t, test_key2, route_established.ConnectedPeer)
}

func TestReceiveMessage(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{},
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		time.Hour,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()
	sample_data := []byte{3, 5, 10, 1, 2, 3}
	to_send := SendMessage{Message{0, false}, sample_data}
	node.MessagesIn <- PhysicalMessage{test_key4, to_send}
	
	select {
		case contents_recvd := <-node.MeshMessagesIn:
			assert.Equal(t, sample_data, contents_recvd)
		case <-node.MessagesOut:
			assert.Fail(t, "Node tried to route message with SendId=0")
	}
}

func TestSendMessageToPeer(t *testing.T) {
	test_route := RouteConfig{[]cipher.PubKey{test_key2}}
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{test_key2},
		1024,
		1024,
		[]RouteConfig{test_route},
		time.Hour,
		time.Hour,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	sample_data := []byte{3, 5, 10, 1, 2, 3}
	
	// Have to wait for route to be established
	go func() {
		for {
			node.SendMessage(0, sample_data)
			time.Sleep(time.Second/4)
		}
	}()

	select {
		case <-node.MeshMessagesIn: {
			assert.Fail(t, "Node received message with SendId != 0")
		}
		case send_message := <-node.MessagesOut: {
			assert.Equal(t, PhysicalMessage{test_key2, SendMessage{Message{0, false}, sample_data}}, send_message)
		}
	}
}

func TestRouteMessage(t *testing.T) {
	test_route := RouteConfig{[]cipher.PubKey{test_key2, test_key3}}
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{test_key2},
		1024,
		1024,
		[]RouteConfig{test_route},
		time.Hour,
		time.Hour,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	// Establish 2 -> 3
	{
		outgoing := <- node.MessagesOut
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg := outgoing.Message
		assert.Equal(t, reflect.TypeOf(EstablishRouteMessage{}), reflect.TypeOf(msg))
		establish_msg := msg.(EstablishRouteMessage)
		assert.Equal(t, test_key3, establish_msg.ToPubKey)

		// Reply
		node.MessagesIn <- PhysicalMessage{test_key2, EstablishRouteReplyMessage{
			OperationReply{OperationMessage{Message{0, true}, establish_msg.MsgId}, true, ""}, 
			11, "secret_abc"}}
	}

	sample_data := []byte{3, 5, 10, 1, 2, 3}

	// Have to wait for route to be established
	go func() {
		for {
			node.SendMessage(0, sample_data)
			time.Sleep(time.Second/4)
		}
	}()

	select {
		case <-node.MeshMessagesIn: {
			assert.Fail(t, "Node received message with SendId != 0")
		}
		case send_message := <-node.MessagesOut: {
			assert.Equal(t, PhysicalMessage{test_key2, SendMessage{Message{11, false}, sample_data}}, send_message)
		}
	}
}

func TestRouteAndRewriteMessage(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{},
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		time.Hour,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	// EstablishRouteMessage
	var establish_reply EstablishRouteReplyMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- 
			PhysicalMessage{
				test_key1,
				EstablishRouteMessage{
					OperationMessage{Message{0, false}, msgId},
					test_key3,
					time.Hour,
				},
			}
		select {
			case route_reply := <-node.MessagesOut: {
				assert.Equal(t, reflect.TypeOf(EstablishRouteReplyMessage{}), reflect.TypeOf(route_reply.Message))
				establish_reply = route_reply.Message.(EstablishRouteReplyMessage)
				newSendId := establish_reply.NewSendId
				newSecret := establish_reply.Secret
				assert.Equal(t, 
					PhysicalMessage{test_key1, 
									EstablishRouteReplyMessage{
										OperationReply{OperationMessage{Message{0, true}, msgId}, true, ""},
										newSendId,
										newSecret,
										}},
					route_reply)
			}
		}
	}

	// Test route without rewrite
	{
		test_contents := []byte{3,7,1,2,3}
		node.MessagesIn <- PhysicalMessage{test_key2,
							 	SendMessage{Message{establish_reply.NewSendId, false}, test_contents}}
		select {
			case physical_msg := <-node.MessagesOut: {
				assert.Equal(t, reflect.TypeOf(SendMessage{}), reflect.TypeOf(physical_msg.Message))
				assert.Equal(t, 
							 PhysicalMessage{test_key3,
							 	SendMessage{Message{0, false}, test_contents}},
							 physical_msg)
			}
		}
	}
/*
	// RouteRewriteMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- 
			PhysicalMessage{
				test_key1,
				RouteRewriteMessage{
					OperationMessage{Message{0, false}, msgId},
					establish_reply.Secret,
					155,
				},
			}
		select {
			case rewrite_reply := <-node.MessagesOut: {
				assert.Equal(t, reflect.TypeOf(OperationReply{}), reflect.TypeOf(rewrite_reply.Message))
				assert.Equal(t, 
					PhysicalMessage{test_key1, 
									OperationReply{OperationMessage{Message{0, true}, msgId}, true, ""},
									},
					rewrite_reply)
			}

		}
	}
*/
	// Test message route and rewrite
	// ...
}

// Rewrite unknown route test
// Routes have distinct indices test

// Send messages thru chain of nodes



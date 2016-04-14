package mesh

import(
	"time"
	"reflect"
	"testing"
	"fmt"
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
var test_key5 = cipher.NewPubKey([]byte{5,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

func init() {
}

func TestEstablishRoutes(t *testing.T) {
	established_ch := make(chan EstablishedRoute)
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{test_key2},
		1024,
		1024,
		[]RouteConfig{RouteConfig{[]cipher.PubKey{test_key2, test_key3, test_key4}}},
		time.Hour,
		0,	// No retransmits
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

		msg := outgoing.Message
		assert.Equal(t, reflect.TypeOf(EstablishRouteMessage{}), reflect.TypeOf(msg))
		establish_msg := msg.(EstablishRouteMessage)
		newDuration := establish_msg.DurationHint
		newMsgId := establish_msg.MsgId
		assert.Equal(t, 
					 PhysicalMessage {
					 	 outgoing.ConnectedPeerPubKey,
						 EstablishRouteMessage {
						 	OperationMessage {
						 		Message {
						 			0,
	    							false,
						 		},
						 		newMsgId,
						 	},
						 	test_key3,
	    					newDuration,
						 },
						},
					 outgoing)

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
		newDuration := establish_msg.DurationHint
		newMsgId := establish_msg.MsgId

		assert.Equal(t, 
					 PhysicalMessage {
					 	 outgoing.ConnectedPeerPubKey,
						 EstablishRouteMessage {
						 	OperationMessage {
						 		Message {
						 			55,
	    							false,
						 		},
						 		newMsgId,
						 	},
						 	test_key4,
	    					newDuration,
						 },
						},
					 outgoing)

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
		newMsgId := rewrite_msg.MsgId

		assert.Equal(t, 
			 PhysicalMessage {
			 	 outgoing.ConnectedPeerPubKey,
				 RouteRewriteMessage {
				 	OperationMessage {
				 		Message {
				 			55,
							false,
				 		},
				 		newMsgId,
				 	},
				 	"secret_abc",
				 	120,
				 },
				},
			 outgoing)

		// Reply
		node.MessagesIn <-  PhysicalMessage{
			test_key2,
			OperationReply{OperationMessage{Message{0, true}, rewrite_msg.MsgId}, true, ""}}
	}
	// Check establish callback
	{
		route_established := <- established_ch

		assert.Equal(t, uint32(55), route_established.SendId)
		assert.Equal(t, test_key2, route_established.ConnectedPeer)
	}
}

func TestReceiveMessage(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{},
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0,	// No retransmits
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	// Default forward
	{
		sample_data := []byte{3, 5, 10, 1, 2, 3}
		to_send := SendMessage{Message{0, false}, sample_data}
		node.MessagesIn <- PhysicalMessage{test_key1, to_send}
		
		select {
			case contents_recvd := <-node.MeshMessagesIn:
				assert.Equal(t, sample_data, contents_recvd)
			case <-node.MessagesOut:
				assert.Fail(t, "Node tried to route message when it should have received it")
		}
	}

	// Default backward
	sample_data := []byte{3, 5, 10, 1, 2, 3}
	{
		to_send := SendMessage{Message{0, true}, sample_data}
		node.MessagesIn <- PhysicalMessage{test_key1, to_send}
		
		select {
			case contents_recvd := <-node.MeshMessagesIn:
				assert.Equal(t, sample_data, contents_recvd)
			case <-node.MessagesOut:
				assert.Fail(t, "Node tried to route message when it should have received it")
		}
	}
}

func TestSendMessageToPeer(t *testing.T) {
	test_route := RouteConfig{[]cipher.PubKey{test_key2}}
	route_established := make(chan bool)
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{test_key2},
		1024,
		1024,
		[]RouteConfig{test_route},
		time.Hour,
		0,	// No retransmits
		// RouteEstablishedCB
		func(route EstablishedRoute) {
			route_established <- true
		},
	}
	node := NewNode(test_config)
	go node.Run()

	sample_data := []byte{3, 5, 10, 1, 2, 3}
	
	<-route_established

	node.SendMessage(0, sample_data)

	select {
		case <-node.MeshMessagesIn: {
			assert.Fail(t, "Node received message when it should have forwarded it")
		}
		case send_message := <-node.MessagesOut: {
			assert.Equal(t, PhysicalMessage{test_key2, SendMessage{Message{0, false}, sample_data}}, send_message)
		}
	}
}

func TestRouteMessage(t *testing.T) {
	test_route := RouteConfig{[]cipher.PubKey{test_key2, test_key3}}
	route_established := make(chan bool)
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{test_key2},
		1024,
		1024,
		[]RouteConfig{test_route},
		time.Hour,
		0,	// No retransmits
		// RouteEstablishedCB
		func(route EstablishedRoute) {
			route_established <- true
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
		node.MessagesIn <- PhysicalMessage{test_key2, EstablishRouteReplyMessage{
			OperationReply{OperationMessage{Message{0, true}, establish_msg.MsgId}, true, ""}, 
			11, "secret_abc"}}
	}
	<-route_established

	sample_data := []byte{3, 5, 10, 1, 2, 3}

	node.SendMessage(0, sample_data)

	select {
		case <-node.MeshMessagesIn: {
			assert.Fail(t, "Node received message when it should have forwarded it")
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
		0,	// No retransmits
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

	// Test message route and rewrite
	{
		test_contents := []byte{10,7,1,128,35}
		node.MessagesIn <- PhysicalMessage{test_key2,
							 	SendMessage{Message{establish_reply.NewSendId, false}, test_contents}}
		select {
			case physical_msg := <-node.MessagesOut: {
				assert.Equal(t, reflect.TypeOf(SendMessage{}), reflect.TypeOf(physical_msg.Message))
				assert.Equal(t, 
							 PhysicalMessage{test_key3,
							 	SendMessage{Message{155, false}, test_contents}},
							 physical_msg)
			}
		}
	}
}

func TestRewriteUnknownRoute(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{},
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0,	// No retransmits
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	// RouteRewriteMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- 
			PhysicalMessage{
				test_key1,
				RouteRewriteMessage{
					OperationMessage{Message{0, false}, msgId},
					"unknown",
					122,
				},
			}
		select {
			case rewrite_reply := <-node.MessagesOut: {
				assert.Equal(t, reflect.TypeOf(OperationReply{}), reflect.TypeOf(rewrite_reply.Message))
				reply := rewrite_reply.Message.(OperationReply)
				assert.Equal(t, 
					PhysicalMessage{test_key1, 
									OperationReply{OperationMessage{Message{0, true}, msgId}, false, reply.Error},
									},
					rewrite_reply)
			}

		}
	}
}

func TestRoutesHaveDifferentSendIds(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		[]cipher.PubKey{test_key2},
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0,	// No retransmits
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	var got_send_id uint32 = 0
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
				establish_reply := route_reply.Message.(EstablishRouteReplyMessage)
				newSendId := establish_reply.NewSendId
				newSecret := establish_reply.Secret
				assert.NotEqual(t, 0, newSendId)
				assert.Equal(t, 
					PhysicalMessage{test_key1, 
									EstablishRouteReplyMessage{
										OperationReply{OperationMessage{Message{0, true}, msgId}, true, ""},
										newSendId,
										newSecret,
										}},
					route_reply)
				got_send_id = newSendId
			}
		}
	}

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
				establish_reply := route_reply.Message.(EstablishRouteReplyMessage)
				newSendId := establish_reply.NewSendId
				newSecret := establish_reply.Secret
				assert.NotEqual(t, 0, newSendId)
				assert.Equal(t, 
					PhysicalMessage{test_key1, 
									EstablishRouteReplyMessage{
										OperationReply{OperationMessage{Message{0, true}, msgId}, true, ""},
										newSendId,
										newSecret,
										}},
					route_reply)
				assert.NotEqual(t, got_send_id, newSendId)
			}
		}
	}
}

func RouteRequestAndRewrite(t *testing.T, node *Node, rewriteId uint32) uint32 {
	var forwardSendId uint32 = 0
	var secret string = ""

	// EstablishRouteMessage
	var establish_reply EstablishRouteReplyMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- 
			PhysicalMessage{
				test_key2,
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
				forwardSendId = establish_reply.NewSendId
				secret = establish_reply.Secret
			}
		}
	}

	// RouteRewriteMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- 
			PhysicalMessage{
				test_key2,
				RouteRewriteMessage{
					OperationMessage{Message{forwardSendId, false}, msgId},
					secret,
					rewriteId,
				},
			}
//fmt.Printf("OperationMessage forwardSendId %v\n", forwardSendId)
		select {
			case rewrite_reply := <-node.MessagesOut: {
				assert.Equal(t, reflect.TypeOf(OperationReply{}), reflect.TypeOf(rewrite_reply.Message))
			}

		}
	}

	return forwardSendId
}

func TestBackwardRoute(t *testing.T) {
	// key1 <-> key 2 <-> key 3
	var test_config = NodeConfig{
		test_key2,
		[]cipher.PubKey{test_key1, test_key3},
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0,	// No retransmits
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	rewriteId := uint32(330)
	forwardSendId := RouteRequestAndRewrite(t, node, rewriteId)

	sample_data := []byte{3, 5, 10, 1, 2, 3}

	// Have to wait for route to be established
	node.MessagesIn <-
		PhysicalMessage{
			test_key2,
			SendMessage {
				Message {
					rewriteId,
					true,
				},
				sample_data,
			},
		}

	select {
		case <-node.MeshMessagesIn: {
			assert.Fail(t, "Node received message when it should have forwarded it")
		}
		case send_message := <-node.MessagesOut: {
			assert.Equal(t, PhysicalMessage{test_key2, SendMessage{Message{forwardSendId, true}, sample_data}}, send_message)
		}
	}
}

func TestInternodeCommunication(t *testing.T) {
	// TODO: Try different numbers of hops
	do_not_want_route_established := make(chan bool)
	want_route_established := make(chan bool)
	do_not_want_cb := func(route EstablishedRoute) {
		do_not_want_route_established <- true
	}
	nodes_by_key := make(map[cipher.PubKey]*Node)
	nodes := []*Node {
		NewNode(
			NodeConfig{
			test_key1,
			[]cipher.PubKey{test_key2},
			1024,
			1024,
			[]RouteConfig{RouteConfig{[]cipher.PubKey{test_key2, test_key3, test_key4, test_key5}}},
			time.Hour,
			0,	// No retransmits
			// RouteEstablishedCB
			func(route EstablishedRoute) {
				want_route_established <- true
			},
		}),
		NewNode(
			NodeConfig{
			test_key2,
			[]cipher.PubKey{test_key1,test_key3},
			1024,
			1024,
			[]RouteConfig{},
			time.Hour,
			0,	// No retransmits
			// RouteEstablishedCB
			do_not_want_cb,
		}),
		NewNode(
			NodeConfig{
			test_key3,
			[]cipher.PubKey{test_key2,test_key4},
			1024,
			1024,
			[]RouteConfig{},
			time.Hour,
			0,	// No retransmits
			// RouteEstablishedCB
			do_not_want_cb,
		}),
		NewNode(
			NodeConfig{
			test_key4,
			[]cipher.PubKey{test_key3,test_key5},
			1024,
			1024,
			[]RouteConfig{},
			time.Hour,
			0,	// No retransmits
			// RouteEstablishedCB
			do_not_want_cb,
		}),
		NewNode(
			NodeConfig{
			test_key5,
			[]cipher.PubKey{test_key4},
			1024,
			1024,
			[]RouteConfig{},
			time.Hour,
			0,	// No retransmits
			// RouteEstablishedCB
			do_not_want_cb,
		}),
	}

	for nodeIdx, _ := range nodes {
		var node *Node = nodes[nodeIdx]
		nodes_by_key[node.Config.MyPubKey] = node
	}
	for nodeIdx, _ := range nodes {
		var node *Node = nodes[nodeIdx]
		go func() {
			for {
				messageToSend := <- node.MessagesOut
	fmt.Printf("messageToSend to %v %v: %v\n", messageToSend.ConnectedPeerPubKey[0], reflect.TypeOf(messageToSend.Message).Name(), messageToSend)
				sendToNode, nodeExists := nodes_by_key[messageToSend.ConnectedPeerPubKey]
				assert.True(t, nodeExists)
				sendToNode.MessagesIn <- PhysicalMessage{node.Config.MyPubKey, messageToSend.Message}
//	fmt.Printf("messageToSend to %v %v: %v\n", messageToSend.ConnectedPeerPubKey[0], reflect.TypeOf(messageToSend.Message).Name(), messageToSend)
			}
		}()
		go node.Run()
	}

	// Wait for route established
	<-want_route_established
	fmt.Printf("Established!\n")
	assert.Equal(t, 0, do_not_want_route_established)
	for _, node := range nodes {
		assert.Equal(t, 0, node.MeshMessagesIn)
	}

for {time.Sleep(time.Second)}
	// TODO: Send message
}




package mesh

import (
	"reflect"
	"testing"
	"time"
)

import (
	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

var test_key1 = cipher.NewPubKey([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
var test_key2 = cipher.NewPubKey([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
var test_key3 = cipher.NewPubKey([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
var test_key4 = cipher.NewPubKey([]byte{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
var test_key5 = cipher.NewPubKey([]byte{5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

func init() {
}

func TestEstablishRoutes(t *testing.T) {
	established_ch := make(chan EstablishedRoute)
	establishment_ch := make(chan []int, 10)
	var test_config = NodeConfig{
		test_key1,
		1024,
		1024,
		[]RouteConfig{RouteConfig{"TestRoute", []cipher.PubKey{test_key2, test_key3, test_key4}}},
		time.Hour,
		0, // No retransmits
		// RouteEstablishmentCallback
		func(RouteIdx int, HopIdx int) {
			establishment_ch <- []int{RouteIdx, HopIdx}
		},
		// RouteEstablishedCB
		func(route EstablishedRoute) {
			established_ch <- route
		},
	}
	node := NewNode(test_config)

	go node.Run()

	// Establish 2 -> 3
	{
		outgoing := <-node.MessagesOut

		msg, error1 := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error1)
		assert.Equal(t, reflect.TypeOf(EstablishRouteMessage{}), reflect.TypeOf(msg))
		establish_msg := msg.(EstablishRouteMessage)
		newDuration := establish_msg.DurationHint
		newMsgId := establish_msg.MsgId
		unserialized_contents, error2 := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error2)
		assert.Equal(t,
			EstablishRouteMessage{
				OperationMessage{
					Message{
						0,
						false,
						0,
					},
					newMsgId,
				},
				test_key3,
				0,
				newDuration,
			},
			unserialized_contents)

		// Reply
		node.MessagesIn <- PhysicalMessage{
			test_key2,
			node.Serializer.SerializeMessage(EstablishRouteReplyMessage{
				OperationReply{OperationMessage{Message{0, true, 0}, establish_msg.MsgId}, true, ""},
				55, "secret_abc"})}
	}
	// Establish 3 -> 4
	{
		outgoing := <-node.MessagesOut

		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg, error1 := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error1)
		assert.Equal(t, reflect.TypeOf(EstablishRouteMessage{}), reflect.TypeOf(msg))
		establish_msg := msg.(EstablishRouteMessage)
		newDuration := establish_msg.DurationHint
		newMsgId := establish_msg.MsgId
		unserialized_contents, error2 := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error2)

		assert.Equal(t,
			EstablishRouteMessage{
				OperationMessage{
					Message{
						55,
						false,
						0,
					},
					newMsgId,
				},
				test_key4,
				55,
				newDuration,
			},
			unserialized_contents)
		// Reply
		node.MessagesIn <- PhysicalMessage{
			test_key3,
			node.Serializer.SerializeMessage(EstablishRouteReplyMessage{
				OperationReply{OperationMessage{Message{0, true, 0}, establish_msg.MsgId}, true, ""},
				120, "secret_xyz"})}
	}
	// Set rewrite 2 -> 3
	{
		outgoing := <-node.MessagesOut
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg, error1 := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error1)
		assert.Equal(t, reflect.TypeOf(RouteRewriteMessage{}), reflect.TypeOf(msg))
		rewrite_msg := msg.(RouteRewriteMessage)
		newMsgId := rewrite_msg.MsgId

		assert.Equal(t,
			RouteRewriteMessage{
				OperationMessage{
					Message{
						55,
						false,
						0,
					},
					newMsgId,
				},
				"secret_abc",
				120,
			},
			msg)

		// Reply
		node.MessagesIn <- PhysicalMessage{
			test_key2,
			node.Serializer.SerializeMessage(OperationReply{OperationMessage{Message{0, true, 0}, rewrite_msg.MsgId}, true, ""})}
	}
	// Reply to ping
	{
		outgoing := <-node.MessagesOut
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg, error1 := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error1)
		assert.Equal(t, reflect.TypeOf(PingMessage{}), reflect.TypeOf(msg))
		ping_msg := msg.(PingMessage)

		// Reply
		node.MessagesIn <- PhysicalMessage{
			test_key2,
			node.Serializer.SerializeMessage(OperationReply{OperationMessage{Message{0, true, 0}, ping_msg.MsgId}, true, ""})}
	}
	// Check establish callback
	{
		route_established := <-established_ch

		assert.Equal(t, uint32(55), route_established.SendId)
		assert.Equal(t, test_key2, route_established.ConnectedPeer)
	}
}

func TestReceiveMessage(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0, // No retransmits
		// RouteEstablishmentCallback
		nil,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	// Default forward
	{
		sample_data := []byte{3, 5, 10, 1, 2, 3}
		to_send := SendMessage{Message{0, false, 11}, sample_data}
		node.MessagesIn <- PhysicalMessage{test_key1, node.Serializer.SerializeMessage(to_send)}

		select {
		case contents_recvd := <-node.MeshMessagesIn:
			assert.Equal(t, test_key1, contents_recvd.ConnectedPeer)
			assert.NotZero(t, contents_recvd.SendId)
			assert.Equal(t, sample_data, contents_recvd.Contents)
		case <-node.MessagesOut:
			assert.Fail(t, "Node tried to route message when it should have received it")
		}
	}

	// Default backward
	sample_data := []byte{3, 5, 10, 1, 2, 3}
	{
		to_send := SendMessage{Message{0, true, 11}, sample_data}
		node.MessagesIn <- PhysicalMessage{test_key1, node.Serializer.SerializeMessage(to_send)}

		select {
		case contents_recvd := <-node.MeshMessagesIn:
			assert.Equal(t, test_key1, contents_recvd.ConnectedPeer)
			assert.NotZero(t, contents_recvd.SendId)
			assert.Equal(t, sample_data, contents_recvd.Contents)
		case <-node.MessagesOut:
			assert.Fail(t, "Node tried to route message when it should have received it")
		}
	}
}

func TestSendMessageToPeer(t *testing.T) {
	test_route := RouteConfig{"Test Route", []cipher.PubKey{test_key2}}
	route_established := make(chan bool)
	var test_config = NodeConfig{
		test_key1,
		1024,
		1024,
		[]RouteConfig{test_route},
		time.Hour,
		0, // No retransmits
		// RouteEstablishmentCallback
		nil,
		// RouteEstablishedCB
		func(route EstablishedRoute) {
			route_established <- true
		},
	}
	node := NewNode(test_config)
	go node.Run()

	// Reply to ping
	{
		outgoing := <-node.MessagesOut
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg, error1 := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error1)
		assert.Equal(t, reflect.TypeOf(PingMessage{}), reflect.TypeOf(msg))
		ping_msg := msg.(PingMessage)

		// Reply
		node.MessagesIn <- PhysicalMessage{
			test_key2,
			node.Serializer.SerializeMessage(OperationReply{OperationMessage{Message{0, true, 0}, ping_msg.MsgId}, true, ""})}
	}

	sample_data := []byte{3, 5, 10, 1, 2, 3}

	<-route_established

	node.SendMessage(0, sample_data)

	select {
	case <-node.MeshMessagesIn:
		{
			assert.Fail(t, "Node received message when it should have forwarded it")
		}
	case send_message := <-node.MessagesOut:
		{
			unserialized_contents, error := node.Serializer.UnserializeMessage(send_message.Contents)
			assert.Nil(t, error)
			assert.Equal(t, SendMessage{Message{0, false, 0}, sample_data}, unserialized_contents)
		}
	}
}

func TestRouteMessage(t *testing.T) {
	test_route := RouteConfig{"TestRoute", []cipher.PubKey{test_key2, test_key3}}
	route_established := make(chan bool)
	var test_config = NodeConfig{
		test_key1,
		1024,
		1024,
		[]RouteConfig{test_route},
		time.Hour,
		0, // No retransmits
		// RouteEstablishmentCallback
		nil,
		// RouteEstablishedCB
		func(route EstablishedRoute) {
			route_established <- true
		},
	}
	node := NewNode(test_config)
	go node.Run()

	// Establish 2 -> 3
	{
		outgoing := <-node.MessagesOut
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg, error := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error)
		assert.Equal(t, reflect.TypeOf(EstablishRouteMessage{}), reflect.TypeOf(msg))
		establish_msg := msg.(EstablishRouteMessage)
		assert.Equal(t, test_key3, establish_msg.ToPubKey)

		// Reply
		node.MessagesIn <- PhysicalMessage{test_key2, node.Serializer.SerializeMessage(EstablishRouteReplyMessage{
			OperationReply{OperationMessage{Message{0, true, 0}, establish_msg.MsgId}, true, ""},
			11, "secret_abc"})}
	}

	// Reply to ping
	{
		outgoing := <-node.MessagesOut
		assert.Equal(t, test_key2, outgoing.ConnectedPeerPubKey)
		msg, error1 := node.Serializer.UnserializeMessage(outgoing.Contents)
		assert.Nil(t, error1)
		assert.Equal(t, reflect.TypeOf(PingMessage{}), reflect.TypeOf(msg))
		ping_msg := msg.(PingMessage)

		// Reply
		node.MessagesIn <- PhysicalMessage{
			test_key2,
			node.Serializer.SerializeMessage(OperationReply{OperationMessage{Message{0, true, 0}, ping_msg.MsgId}, true, ""})}
	}

	<-route_established

	sample_data := []byte{3, 5, 10, 1, 2, 3}

	node.SendMessage(0, sample_data)

	select {
	case <-node.MeshMessagesIn:
		{
			assert.Fail(t, "Node received message when it should have forwarded it")
		}
	case send_message := <-node.MessagesOut:
		{
			unserialized_contents, error := node.Serializer.UnserializeMessage(send_message.Contents)
			assert.Nil(t, error)
			assert.Equal(t, test_key2, send_message.ConnectedPeerPubKey)
			assert.Equal(t, SendMessage{Message{11, false, 0}, sample_data}, unserialized_contents)
		}
	}
}

func TestRouteAndRewriteMessage(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0, // No retransmits
		// RouteEstablishmentCallback
		nil,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	// EstablishRouteMessage
	var establish_reply EstablishRouteReplyMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- PhysicalMessage{
			test_key1,
			node.Serializer.SerializeMessage(EstablishRouteMessage{
				OperationMessage{Message{0, false, 0}, msgId},
				test_key3,
				0,
				time.Hour,
			}),
		}
		select {
		case route_reply := <-node.MessagesOut:
			{
				unserialized_contents, error := node.Serializer.UnserializeMessage(route_reply.Contents)
				assert.Nil(t, error)
				assert.Equal(t, reflect.TypeOf(EstablishRouteReplyMessage{}), reflect.TypeOf(unserialized_contents))
				establish_reply = unserialized_contents.(EstablishRouteReplyMessage)
				newSendId := establish_reply.NewSendId
				newSecret := establish_reply.Secret
				assert.Equal(t, test_key1, route_reply.ConnectedPeerPubKey)
				assert.Equal(t,
					EstablishRouteReplyMessage{
						OperationReply{OperationMessage{Message{0, true, 0}, msgId}, true, ""},
						newSendId,
						newSecret,
					},
					establish_reply)
			}
		}
	}

	// Test route without rewrite
	{
		test_contents := []byte{3, 7, 1, 2, 3}
		node.MessagesIn <- PhysicalMessage{test_key2,
			node.Serializer.SerializeMessage(
				SendMessage{Message{establish_reply.NewSendId, false, 0}, test_contents})}
		select {
		case physical_msg := <-node.MessagesOut:
			{
				unserialized_contents, error := node.Serializer.UnserializeMessage(physical_msg.Contents)
				assert.Nil(t, error)
				assert.Equal(t, reflect.TypeOf(SendMessage{}), reflect.TypeOf(unserialized_contents))
				assert.Equal(t,
					SendMessage{Message{0, false, establish_reply.NewSendId}, test_contents},
					unserialized_contents)
			}
		}
	}

	// RouteRewriteMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- PhysicalMessage{
			test_key1,
			node.Serializer.SerializeMessage(RouteRewriteMessage{
				OperationMessage{Message{0, false, 0}, msgId},
				establish_reply.Secret,
				155,
			}),
		}
		select {
		case rewrite_reply := <-node.MessagesOut:
			{
				assert.Equal(t, test_key1, rewrite_reply.ConnectedPeerPubKey)
				unserialized_contents, error := node.Serializer.UnserializeMessage(rewrite_reply.Contents)
				assert.Nil(t, error)
				assert.Equal(t, reflect.TypeOf(OperationReply{}), reflect.TypeOf(unserialized_contents))
				assert.Equal(t,
					OperationReply{OperationMessage{Message{0, true, 0}, msgId}, true, ""},
					unserialized_contents)
			}

		}
	}

	// Test message route and rewrite
	{
		test_contents := []byte{10, 7, 1, 128, 35}
		node.MessagesIn <- PhysicalMessage{test_key2,
			node.Serializer.SerializeMessage(SendMessage{Message{establish_reply.NewSendId, false, 0}, test_contents})}
		select {
		case physical_msg := <-node.MessagesOut:
			{
				assert.Equal(t, test_key3, physical_msg.ConnectedPeerPubKey)
				unserialized_contents, error := node.Serializer.UnserializeMessage(physical_msg.Contents)
				assert.Nil(t, error)
				assert.Equal(t, reflect.TypeOf(SendMessage{}), reflect.TypeOf(unserialized_contents))
				assert.Equal(t,
					SendMessage{Message{155, false, establish_reply.NewSendId}, test_contents},
					unserialized_contents)
			}
		}
	}
}

func TestRewriteUnknownRoute(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0, // No retransmits
		// RouteEstablishmentCallback
		nil,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	// RouteRewriteMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- PhysicalMessage{
			test_key1,
			node.Serializer.SerializeMessage(RouteRewriteMessage{
				OperationMessage{Message{0, false, 0}, msgId},
				"unknown",
				122,
			}),
		}
		select {
		case rewrite_reply := <-node.MessagesOut:
			{
				assert.Equal(t, test_key1, rewrite_reply.ConnectedPeerPubKey)
				unserialized_contents, error := node.Serializer.UnserializeMessage(rewrite_reply.Contents)
				assert.Nil(t, error)
				assert.Equal(t, reflect.TypeOf(OperationReply{}), reflect.TypeOf(unserialized_contents))
				reply := unserialized_contents.(OperationReply)
				assert.Equal(t,
					OperationReply{OperationMessage{Message{0, true, 0}, msgId}, false, reply.Error},
					unserialized_contents)
			}

		}
	}
}

func TestRoutesHaveDifferentSendIds(t *testing.T) {
	var test_config = NodeConfig{
		test_key1,
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0, // No retransmits
		// RouteEstablishmentCallback
		nil,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	var got_send_id uint32 = 0
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- PhysicalMessage{
			test_key1,
			node.Serializer.SerializeMessage(EstablishRouteMessage{
				OperationMessage{Message{0, false, 0}, msgId},
				test_key3,
				0,
				time.Hour,
			}),
		}
		select {
		case route_reply := <-node.MessagesOut:
			{
				assert.Equal(t, test_key1, route_reply.ConnectedPeerPubKey)
				unserialized_contents, error := node.Serializer.UnserializeMessage(route_reply.Contents)
				assert.Nil(t, error)
				assert.Equal(t, reflect.TypeOf(EstablishRouteReplyMessage{}), reflect.TypeOf(unserialized_contents))
				establish_reply := unserialized_contents.(EstablishRouteReplyMessage)
				newSendId := establish_reply.NewSendId
				newSecret := establish_reply.Secret
				assert.NotEqual(t, 0, newSendId)
				assert.Equal(t,
					EstablishRouteReplyMessage{
						OperationReply{OperationMessage{Message{0, true, 0}, msgId}, true, ""},
						newSendId,
						newSecret,
					},
					unserialized_contents)
				got_send_id = newSendId
			}
		}
	}

	{
		msgId := uuid.NewV4()
		node.MessagesIn <- PhysicalMessage{
			test_key1,
			node.Serializer.SerializeMessage(EstablishRouteMessage{
				OperationMessage{Message{0, false, 0}, msgId},
				test_key3,
				0,
				time.Hour,
			}),
		}
		select {
		case route_reply := <-node.MessagesOut:
			{
				assert.Equal(t, test_key1, route_reply.ConnectedPeerPubKey)
				unserialized_contents, error := node.Serializer.UnserializeMessage(route_reply.Contents)
				assert.Nil(t, error)
				assert.Equal(t, reflect.TypeOf(EstablishRouteReplyMessage{}), reflect.TypeOf(unserialized_contents))
				establish_reply := unserialized_contents.(EstablishRouteReplyMessage)
				newSendId := establish_reply.NewSendId
				newSecret := establish_reply.Secret
				assert.NotEqual(t, 0, newSendId)
				assert.Equal(t,
					EstablishRouteReplyMessage{
						OperationReply{OperationMessage{Message{0, true, 0}, msgId}, true, ""},
						newSendId,
						newSecret,
					},
					unserialized_contents)
				assert.NotEqual(t, got_send_id, newSendId)
			}
		}
	}
}

func TestBackwardRoute(t *testing.T) {
	// key1 <-> key 2 <-> key 3
	var test_config = NodeConfig{
		test_key2,
		1024,
		1024,
		[]RouteConfig{},
		time.Hour,
		0, // No retransmits
		// RouteEstablishmentCallback
		nil,
		// RouteEstablishedCB
		nil,
	}
	node := NewNode(test_config)
	go node.Run()

	backwardRewriteId := uint32(1001)
	var forwardSendId uint32 = 0

	// EstablishRouteMessage
	{
		msgId := uuid.NewV4()
		node.MessagesIn <- PhysicalMessage{
			test_key2,
			node.Serializer.SerializeMessage(EstablishRouteMessage{
				OperationMessage{Message{0, false, 0}, msgId},
				test_key3,
				backwardRewriteId,
				time.Hour,
			}),
		}
		select {
		case route_reply := <-node.MessagesOut:
			{
				assert.Equal(t, test_key2, route_reply.ConnectedPeerPubKey)
				unserialized_contents, error := node.Serializer.UnserializeMessage(route_reply.Contents)
				assert.Nil(t, error)
				assert.Equal(t, reflect.TypeOf(EstablishRouteReplyMessage{}), reflect.TypeOf(unserialized_contents))
				establish_reply := unserialized_contents.(EstablishRouteReplyMessage)
				forwardSendId = establish_reply.NewSendId
			}
		}
	}

	sample_data := []byte{3, 5, 10, 1, 2, 3}

	node.MessagesIn <- PhysicalMessage{
		test_key2,
		node.Serializer.SerializeMessage(SendMessage{
			Message{
				forwardSendId,
				true,
				0,
			},
			sample_data,
		}),
	}

	select {
	case <-node.MeshMessagesIn:
		{
			assert.Fail(t, "Node received message when it should have forwarded it")
		}
	case send_message := <-node.MessagesOut:
		{
			assert.Equal(t, test_key2, send_message.ConnectedPeerPubKey)
			unserialized_contents, error := node.Serializer.UnserializeMessage(send_message.Contents)
			assert.Nil(t, error)
			assert.Equal(t,
				SendMessage{Message{backwardRewriteId, true, forwardSendId}, sample_data},
				unserialized_contents)
		}
	}
}

func TestInternodeCommunication(t *testing.T) {
	do_not_want_route_established := make(chan bool)
	want_route_established := make(chan bool)
	do_not_want_cb := func(route EstablishedRoute) {
		do_not_want_route_established <- true
	}

	nodes_configs := [][]*Node{
		[]*Node{
			NewNode(
				NodeConfig{
					test_key1,
					1024,
					1024,
					[]RouteConfig{RouteConfig{"TestRoute", []cipher.PubKey{test_key2}}},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					func(route EstablishedRoute) {
						want_route_established <- true
					},
				}),
			NewNode(
				NodeConfig{
					test_key2,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
		},
		[]*Node{
			NewNode(
				NodeConfig{
					test_key1,
					1024,
					1024,
					[]RouteConfig{RouteConfig{"TestRoute", []cipher.PubKey{test_key2, test_key3}}},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					func(route EstablishedRoute) {
						want_route_established <- true
					},
				}),
			NewNode(
				NodeConfig{
					test_key2,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
			NewNode(
				NodeConfig{
					test_key3,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
		},
		[]*Node{
			NewNode(
				NodeConfig{
					test_key1,
					1024,
					1024,
					[]RouteConfig{RouteConfig{"TestRoute", []cipher.PubKey{test_key2, test_key3, test_key4}}},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					func(route EstablishedRoute) {
						want_route_established <- true
					},
				}),
			NewNode(
				NodeConfig{
					test_key2,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
			NewNode(
				NodeConfig{
					test_key3,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
			NewNode(
				NodeConfig{
					test_key4,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
		},
		[]*Node{
			NewNode(
				NodeConfig{
					test_key1,
					1024,
					1024,
					[]RouteConfig{RouteConfig{"TestRoute", []cipher.PubKey{test_key2, test_key3, test_key4, test_key5}}},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					func(route EstablishedRoute) {
						want_route_established <- true
					},
				}),
			NewNode(
				NodeConfig{
					test_key2,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
			NewNode(
				NodeConfig{
					test_key3,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
			NewNode(
				NodeConfig{
					test_key4,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
			NewNode(
				NodeConfig{
					test_key5,
					1024,
					1024,
					[]RouteConfig{},
					time.Hour,
					0, // No retransmits
					// RouteEstablishmentCallback
					nil,
					// RouteEstablishedCB
					do_not_want_cb,
				}),
		},
	}

	for _, nodes := range nodes_configs {
		t.Logf("Testing with %v nodes\n", len(nodes))
		for {
			if len(do_not_want_route_established) == 0 {
				break
			}
			<-do_not_want_route_established
		}
		for {
			if len(want_route_established) == 0 {
				break
			}
			<-want_route_established
		}
		nodes_by_key := make(map[cipher.PubKey]*Node)
		for nodeIdx, _ := range nodes {
			var node *Node = nodes[nodeIdx]
			nodes_by_key[node.Config.MyPubKey] = node
		}
		for nodeIdx, _ := range nodes {
			var node *Node = nodes[nodeIdx]
			go func() {
				for {
					messageToSend := <-node.MessagesOut
					sendToNode, nodeExists := nodes_by_key[messageToSend.ConnectedPeerPubKey]
					assert.True(t, nodeExists)
					sendToNode.MessagesIn <- PhysicalMessage{node.Config.MyPubKey, messageToSend.Contents}
				}
			}()
			go node.Run()
		}

		// Wait for route established
		<-want_route_established
		assert.Equal(t, 0, len(do_not_want_route_established))
		for _, node := range nodes {
			assert.Equal(t, 0, len(node.MeshMessagesIn))
		}

		// Send
		sample_data := []byte{50, 10, 1, 2, 3}
		nodes[0].SendMessage(0, sample_data)

		var received_mesh_msg MeshMessage

		select {
		case received_mesh_msg = <-nodes[len(nodes)-1].MeshMessagesIn:
			{
				route_id := received_mesh_msg.SendId
				if len(nodes) > 2 {
					assert.NotZero(t, route_id)
				} else {
					assert.Zero(t, route_id)
				}
				assert.Equal(t, MeshMessage{route_id, nodes[len(nodes)-2].Config.MyPubKey, sample_data}, received_mesh_msg)
			}
		}

		// Reply
		sample_reply_data := []byte{5, 7, 3, 2, 2}
		nodes[len(nodes)-1].SendReply(received_mesh_msg, sample_reply_data)

		select {
		case received_mesh_msg := <-nodes[0].MeshMessagesIn:
			{
				route_id := received_mesh_msg.SendId
				if len(nodes) > 2 {
					assert.NotZero(t, route_id)
				} else {
					assert.Zero(t, route_id)
				}
				assert.Equal(t, MeshMessage{route_id, test_key2, sample_reply_data}, received_mesh_msg)
			}
		}
	}
}

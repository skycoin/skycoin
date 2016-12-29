package node

import (
	"github.com/satori/go.uuid"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

type ControlChannel struct {
	Id              uuid.UUID
	IncomingChannel chan []byte
}

func NewControlChannel() *ControlChannel {
	c := ControlChannel{
		Id:              uuid.NewV4(),
		IncomingChannel: make(chan []byte),
	}
	return &c
}

func (c *ControlChannel) HandleMessage(handledNode *Node, msg []byte) error {

	switch messages.GetMessageType(msg) {

	case messages.MsgCreateChannelControlMessage:
		controlChannel := NewControlChannel()
		handledNode.AddControlChannel(controlChannel)
		return nil

	case messages.MsgAddRouteControlMessage:
		var m1 messages.AddRouteControlMessage
		messages.Deserialize(msg, m1)
		routeId := m1.RouteId
		nodeToAdd := m1.NodeId
		return handledNode.addRoute(nodeToAdd, routeId)

	case messages.MsgExtendRouteControlMessage:
		//do something
		//var m1 messages.ExtendRouteControlMessage
		//messages.Deserialize(msg, m1)
		//routeId := m1.RouteId
		//nodeToExtend := m1.NodeId
		//return handledNode.extendRoute(nodeToAdd, routeId)
		return nil

	case messages.MsgRemoveRouteControlMessage:
		var m1 messages.RemoveRouteControlMessage
		messages.Deserialize(msg, m1)
		routeId := m1.RouteId
		return handledNode.removeRoute(routeId)
	}

	return nil
}

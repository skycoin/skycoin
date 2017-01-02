package node

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

type ControlChannel struct {
	Id uuid.UUID
}

func NewControlChannel() *ControlChannel {
	c := ControlChannel{
		Id: uuid.NewV4(),
	}
	return &c
}

func (c *ControlChannel) HandleMessage(handledNode *Node, msg []byte) (interface{}, error) {
	switch messages.GetMessageType(msg) {

	case messages.MsgCreateChannelControlMessage:
		channelID := handledNode.AddControlChannel()
		return channelID, nil

	case messages.MsgAddRouteControlMessage:
		fmt.Println("adding route")
		var m1 messages.AddRouteControlMessage
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			panic(err)
		}
		routeId := m1.RouteId
		nodeToAdd := m1.NodeId
		return nil, handledNode.addRoute(nodeToAdd, routeId)

	case messages.MsgExtendRouteControlMessage:
		//var m1 messages.ExtendRouteControlMessage
		//messages.Deserialize(msg, &m1)
		//routeId := m1.RouteId
		//nodeToExtend := m1.NodeId
		//return handledNode.extendRoute(nodeToAdd, routeId)
		return nil, nil

	case messages.MsgRemoveRouteControlMessage:
		fmt.Println("removing route")
		var m1 messages.RemoveRouteControlMessage
		messages.Deserialize(msg, &m1)
		routeId := m1.RouteId
		return nil, handledNode.removeRoute(routeId)
	}

	return nil, errors.New("Unknown message type")
}

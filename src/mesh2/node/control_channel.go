package node

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

type ControlChannel struct {
	Id messages.ChannelId
}

func NewControlChannel() *ControlChannel {
	c := ControlChannel{
		Id: messages.RandChannelId(),
	}
	return &c
}

func (c *ControlChannel) HandleMessage(handledNode *Node, msg []byte) (interface{}, error) {
	switch messages.GetMessageType(msg) {
	/*
		case messages.MsgCreateChannelControlMessage:
			channelID := handledNode.AddControlChannel()
			return channelID, nil
	*/
	case messages.MsgAddRouteControlMessage:
		fmt.Println("adding route")
		var m1 messages.AddRouteControlMessage
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			panic(err)
		}
		routeRule := RouteRule{
			m1.IncomingTransportId,
			m1.OutgoingTransportId,
			m1.IncomingRouteId,
			m1.OutgoingRouteId,
		}
		return nil, handledNode.addRoute(&routeRule)

	case messages.MsgRemoveRouteControlMessage:
		fmt.Println("removing route")
		var m1 messages.RemoveRouteControlMessage
		messages.Deserialize(msg, &m1)
		routeId := m1.RouteId
		return nil, handledNode.removeRoute(routeId)
	}

	return nil, errors.New("Unknown message type")
}

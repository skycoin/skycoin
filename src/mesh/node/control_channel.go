package node

import (
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type ControlChannel struct {
	id messages.ChannelId
}

func newControlChannel() *ControlChannel {
	c := ControlChannel{
		id: messages.RandChannelId(),
	}
	return &c
}

func (c *ControlChannel) handleMessage(handledNode *Node, msg []byte) error {
	switch messages.GetMessageType(msg) {
	case messages.MsgAddRouteControlMessage:
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
		err = handledNode.addRoute(&routeRule)
		return err

	case messages.MsgRemoveRouteControlMessage:
		var m1 messages.RemoveRouteControlMessage
		messages.Deserialize(msg, &m1)
		routeId := m1.RouteId
		err := handledNode.removeRoute(routeId)
		return err
	}

	return messages.ERR_UNKNOWN_MESSAGE_TYPE
}

package mesh_rpc

import (
	"fmt"

	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/node"
)

type RPCReceiver struct {
}

func (receiver *RPCReceiver) CreateControlChannel(args []interface{}, result *[]byte) error {
	handledNode := args[0].(*node.Node)
	controlChannel := node.NewControlChannel()
	fmt.Println("NODE IS", handledNode)
	fmt.Println("NODE controlChannels are", handledNode.ControlChannels)
	handledNode.AddControlChannel(controlChannel)
	res := controlChannel.Id
	*result = messages.Serialize((uint16)(0), res)
	return nil
}

/*
func (receiver *RPCReceiver) AddRoute(args []interface{}, result smth) error {
	nodeFrom := args[0].(cipher.PubKey)
	nodeTo := args[1].(cipher.PubKey)
	routeId := args[2].(RouteId)

}

func (receiver *RPCReceiver) ExtendRoute(args []interface{}, result smth) error {
	nodeFrom := args[0].(cipher.PubKey)
	nodeTo := args[1].(cipher.PubKey)
	routeId := args[2].(RouteId)

}

func (receiver *RPCReceiver) RemoveRoute(args []interface{}, result smth) error {
	node := args[0].(cipher.PubKey)
	routeId := args[2].(RouteId)

}
*/
/*
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
*/

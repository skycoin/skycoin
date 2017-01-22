package nodemanager

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/node"
)

func (self *NodeManager) BuildRoute(nodes []cipher.PubKey) []messages.RouteId {

	n := len(nodes)

	routeIds := make([]messages.RouteId, n)

	for i, _ := range nodes {
		routeIds[i] = messages.RandRouteId()
	}

	for i, currentNodeId := range nodes {

		currentNode, err := self.GetNodeById(currentNodeId)
		if err != nil {
			fmt.Println(err)
			return []messages.RouteId{}
		}

		incomingRoute := routeIds[i]

		var prevNodeId, nextNodeId cipher.PubKey
		var prevNode *node.Node
		var incomingTransport, outgoingTransport messages.TransportId
		var outgoingRoute messages.RouteId

		// if it is the first node in the route, there is no incoming transport
		if i == 0 {
			incomingTransport = (messages.TransportId)(0)
		} else {
			prevNodeId = nodes[i-1]
			prevNode, _ = self.GetNodeById(prevNodeId)
			incomingTransportObj, err := prevNode.GetTransportToNode(currentNodeId)
			if err != nil {
				fmt.Println(err)
				return []messages.RouteId{}
			}
			incomingTransport = incomingTransportObj.StubPair.Id
		}

		// if it is the last node in the route, there is no outgoing transport and outgoing route
		if i == len(nodes)-1 {
			outgoingTransport = (messages.TransportId)(0)
			outgoingRoute = (messages.RouteId)(0)
		} else {
			outgoingRoute = routeIds[i+1]
			nextNodeId = nodes[i+1]
			outgoingTransportObj, err := currentNode.GetTransportToNode(nextNodeId)
			if err != nil {
				fmt.Println(err)
				return []messages.RouteId{}
			}
			outgoingTransport = outgoingTransportObj.Id
		}

		msg := messages.AddRouteControlMessage{
			incomingTransport,
			outgoingTransport,
			incomingRoute,
			outgoingRoute,
		}
		msgS := messages.Serialize(messages.MsgAddRouteControlMessage, msg)

		ccid := currentNode.AddControlChannel()
		controlMessage := messages.InControlMessage{ccid, msgS}

		currentNode.InjectControlMessage(controlMessage)
	}

	return routeIds
}

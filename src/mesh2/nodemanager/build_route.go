package nodemanager

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/errors"
	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/node"
)

func (self *NodeManager) BuildRoute(nodes []cipher.PubKey) (route, backRoute messages.RouteId, err error) {
	route, err = self.getFirstRoute(nodes)
	if err == nil {
		for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 { //side effect but nodes aren't used after this anyway
			nodes[i], nodes[j] = nodes[j], nodes[i]
		}
		backRoute, err = self.getFirstRoute(nodes)
	}
	return
}

func (self *NodeManager) FindRoute(nodeFrom, nodeTo cipher.PubKey) (routeId, backRouteId messages.RouteId, err error) {
	nodes, found := self.routeGraph.findRoute(nodeFrom, nodeTo)
	if !found {
		return messages.NIL_ROUTE, messages.NIL_ROUTE, errors.ERR_NOROUTE
	}
	routeId, backRouteId, err = self.BuildRoute(nodes)
	return
}

func (self *NodeManager) getFirstRoute(nodes []cipher.PubKey) (messages.RouteId, error) {
	routes := self.buildRoute(nodes)
	if len(routes) < 0 {
		return messages.NIL_ROUTE, errors.ERR_NOROUTE
	} else {
		return routes[0], nil
	}
}

func (self *NodeManager) buildRoute(nodes []cipher.PubKey) []messages.RouteId {

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
			incomingTransport = messages.NIL_TRANSPORT
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
			outgoingTransport = messages.NIL_TRANSPORT
			outgoingRoute = messages.NIL_ROUTE
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

func (self *NodeManager) rebuildRouteGraph() {
	self.routeGraph.clear()
	for _, nodeFrom := range self.nodeList {
		nodeFromId := nodeFrom.GetId()
		for _, transport := range nodeFrom.Transports {
			nodeToId := transport.StubPair.AttachedNode.GetId()
			self.routeGraph.addDirectRoute(nodeFromId, nodeToId, 1) // weight is always 1 because so far all routes are equal! Change this if needed
		}
	}
}

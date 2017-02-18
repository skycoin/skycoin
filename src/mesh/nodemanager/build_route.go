package nodemanager

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
)

func (self *NodeManager) findRoute(nodeFrom, nodeTo cipher.PubKey) (routeId, backRouteId messages.RouteId, err error) {
	nodes, err := self.routeGraph.findRoute(nodeFrom, nodeTo)
	if err != nil {
		return messages.NIL_ROUTE, messages.NIL_ROUTE, err
	}
	routeId, backRouteId, err = self.buildRoute(nodes)
	return
}

func (self *NodeManager) findRouteForward(nodeFrom, nodeTo cipher.PubKey) (routes []messages.RouteId, err error) {
	nodes, err := self.routeGraph.findRoute(nodeFrom, nodeTo)
	if err != nil {
		return nil, err
	}
	routes, err = self.buildRouteForward(nodes)
	return
}

func (self *NodeManager) buildRoute(nodes []cipher.PubKey) (route, backRoute messages.RouteId, err error) {
	route, err = self.getFirstRouteForward(nodes)
	if err != nil {
		return
	}
	backRoute, err = self.getFirstRouteBack(nodes)
	return
}

func (self *NodeManager) getFirstRouteForward(nodes []cipher.PubKey) (messages.RouteId, error) {
	return self.getFirstRoute(nodes, true)
}

func (self *NodeManager) getFirstRouteBack(nodes []cipher.PubKey) (messages.RouteId, error) {
	return self.getFirstRoute(nodes, false)
}

func (self *NodeManager) getFirstRoute(nodes []cipher.PubKey, forward bool) (messages.RouteId, error) {
	routes, err := self.buildRouteOneSide(nodes, forward)
	if err != nil {
		return messages.NIL_ROUTE, err
	}
	if len(routes) < 1 {
		return messages.NIL_ROUTE, errors.ERR_NO_ROUTE
	}
	return routes[0], nil
}

func (self *NodeManager) buildRouteForward(nodes []cipher.PubKey) ([]messages.RouteId, error) {
	return self.buildRouteOneSide(nodes, true)
}

func (self *NodeManager) buildRouteBackward(nodes []cipher.PubKey) ([]messages.RouteId, error) {
	return self.buildRouteOneSide(nodes, false)
}

func (self *NodeManager) buildRouteOneSide(nodes []cipher.PubKey, forward bool) ([]messages.RouteId, error) {

	n := len(nodes)

	routeIds := make([]messages.RouteId, n)

	var startIndex, endIndex, next int

	if forward {
		startIndex = 0
		endIndex = n
		next = 1
	} else {
		startIndex = n - 1
		endIndex = -1
		next = -1
	}

	for i := range routeIds {
		routeIds[i] = messages.RandRouteId()
	}

	for i := startIndex; i != endIndex; i += next {

		currentNodeId := nodes[i]

		currentNode, err := self.getNodeById(currentNodeId)
		if err != nil {
			return []messages.RouteId{}, err
		}

		var routeIndex int
		if forward {
			routeIndex = i
		} else {
			routeIndex = startIndex - i
		}

		incomingRoute := routeIds[routeIndex]

		var prevNodeId, nextNodeId cipher.PubKey
		var prevNode *node.Node
		var incomingTransport, outgoingTransport messages.TransportId
		var outgoingRoute messages.RouteId

		// if it is the first node in the route, there is no incoming transport
		if i == startIndex {
			incomingTransport = messages.NIL_TRANSPORT
		} else {
			prevNodeId = nodes[i-next]
			prevNode, _ = self.getNodeById(prevNodeId)
			incomingTransportObj, err := prevNode.GetTransportToNode(currentNodeId)
			if err != nil {
				return []messages.RouteId{}, err
			}
			incomingTransport = incomingTransportObj.StubPair.Id
		}

		// if it is the last node in the route, there is no outgoing transport and outgoing route
		if i == endIndex-next {
			outgoingTransport = messages.NIL_TRANSPORT
			outgoingRoute = messages.NIL_ROUTE
		} else {
			outgoingRoute = routeIds[routeIndex+1]
			nextNodeId = nodes[i+next]
			outgoingTransportObj, err := currentNode.GetTransportToNode(nextNodeId)
			if err != nil {
				return []messages.RouteId{}, err
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
		controlMessage := messages.InControlMessage{ccid, msgS, nil}

		currentNode.InjectControlMessage(controlMessage)
	}

	return routeIds, nil
}

func (self *NodeManager) rebuildRoutes() {
	self.routeGraph.clear()
	for _, nodeFrom := range self.nodeList {
		nodeFromId := nodeFrom.GetId()
		for _, transport := range nodeFrom.Transports {
			nodeToId := transport.StubPair.AttachedNode.GetId()
			self.routeGraph.addDirectRoute(nodeFromId, nodeToId, 1) // weight is always 1 because so far all routes are equal! Change this if needed
		}
	}
}

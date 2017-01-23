package nodemanager

//methods for testing purposes only

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/node"
)

func (self *NodeManager) CreateNodeList(n int) []cipher.PubKey {
	nodes := []cipher.PubKey{}
	for i := 0; i < n; i++ {
		nodeId := self.AddNewNode()
		nodes = append(nodes, nodeId)
	}
	return nodes
}

func (self *NodeManager) ConnectAll() {

	n := len(self.nodeIdList)

	// form transportFactories between the nodes sequentially
	for i := 0; i < n-1; i++ {
		id1, id2 := self.nodeIdList[i], self.nodeIdList[i+1]
		self.ConnectNodeToNode(id1, id2)
	}

	// create rules for building a route from the first node to the last one
	rules := []*node.RouteRule{}
	for i := 0; i < n; i++ {
		ruleId := messages.RandRouteId()
		rules = append(rules, &node.RouteRule{IncomingRoute: ruleId})
	}
	for i := 0; i < n; i++ {
		routeRule := rules[i]
		nodeId := self.nodeIdList[i]
		if i > 0 {
			prevNodeId := self.nodeIdList[i-1]
			incomingTransport, err := self.nodeList[prevNodeId].GetTransportToNode(nodeId)
			if err != nil {
				panic(err)
			}
			routeRule.IncomingTransport = incomingTransport.StubPair.Id
		} else {
			routeRule.IncomingTransport = messages.NIL_TRANSPORT
		}
		if i < n-1 {
			nextNodeId := self.nodeIdList[i+1]
			outgoingTransport, err := self.nodeList[nodeId].GetTransportToNode(nextNodeId)
			if err != nil {
				panic(err)
			}
			routeRule.OutgoingTransport = outgoingTransport.Id
			routeRule.OutgoingRoute = rules[i+1].IncomingRoute
		} else {
			routeRule.OutgoingTransport = messages.NIL_TRANSPORT
			routeRule.OutgoingRoute = messages.NIL_ROUTE
		}
		addRouteMessage := messages.AddRouteControlMessage{
			routeRule.IncomingTransport,
			routeRule.OutgoingTransport,
			routeRule.IncomingRoute,
			routeRule.OutgoingRoute,
		}
		serialized := messages.Serialize(messages.MsgAddRouteControlMessage, addRouteMessage)

		node0 := self.nodeList[nodeId]

		ccid := node0.AddControlChannel()
		controlMessage := messages.InControlMessage{ccid, serialized}
		node0.InjectControlMessage(controlMessage)
	}
}

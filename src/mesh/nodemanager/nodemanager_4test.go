package nodemanager

//methods for testing purposes only

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

func (self *NodeManager) CreateSequenceOfNodes(n int) (cipher.PubKey, cipher.PubKey) {
	/*
		This function creates a network with sequentially chained n nodes like 0-1-2-3-4-5-6-7-8-9 and returns the addresses of the first and last node
	*/

	nodeList := self.createNodeList(n)
	self.connectAll()
	self.rebuildRoutes()
	firstNode, lastNode := nodeList[0], nodeList[len(nodeList)-1]
	return firstNode, lastNode
}

func (self *NodeManager) CreateSequenceOfNodesAndBuildRoutes(n int) (cipher.PubKey, cipher.PubKey, messages.RouteId, messages.RouteId) {
	/*
		This function creates a network with sequentially chained n nodes like 0-1-2-3-4-5-6-7-8-9, builds route between the first and the last nodes in a chainand returns the addresses of them, a route from the first to the last one and a back route from the last to the first one
	*/

	nodeList := self.createNodeList(n)
	self.connectAll()

	route, backRoute, err := self.buildRoute(nodeList)
	if err != nil {
		panic(err)
	}
	serverNode, clientNode := nodeList[0], nodeList[len(nodeList)-1]
	return clientNode, serverNode, route, backRoute
}

func (self *NodeManager) CreateThreeRoutes() (cipher.PubKey, cipher.PubKey) {
	nodeList := self.createNodeList(10)
	/*
		  1-2-3-4
		 /	 \
		0----5----9
		 \	 /
		  6_7_8_/
	*/
	self.ConnectNodeToNode(nodeList[0], nodeList[1])
	self.ConnectNodeToNode(nodeList[1], nodeList[2])
	self.ConnectNodeToNode(nodeList[2], nodeList[3])
	self.ConnectNodeToNode(nodeList[3], nodeList[4])
	self.ConnectNodeToNode(nodeList[4], nodeList[9])
	self.ConnectNodeToNode(nodeList[0], nodeList[5])
	self.ConnectNodeToNode(nodeList[5], nodeList[9])
	self.ConnectNodeToNode(nodeList[0], nodeList[6])
	self.ConnectNodeToNode(nodeList[6], nodeList[7])
	self.ConnectNodeToNode(nodeList[7], nodeList[8])
	self.ConnectNodeToNode(nodeList[8], nodeList[9])

	self.rebuildRoutes()

	clientNode, serverNode := nodeList[0], nodeList[9]
	return clientNode, serverNode
}

func (self *NodeManager) createNodeList(n int) []cipher.PubKey {
	nodes := []cipher.PubKey{}
	for i := 0; i < n; i++ {
		nodeId := self.AddNewNode()
		nodes = append(nodes, nodeId)
	}
	return nodes
}

func (self *NodeManager) connectAll() {

	n := len(self.nodeIdList)

	for i := 0; i < n-1; i++ {
		id1, id2 := self.nodeIdList[i], self.nodeIdList[i+1]
		self.ConnectNodeToNode(id1, id2)
	}
	return
}

func (self *NodeManager) connectAllAndBuildRoute() (messages.RouteId, error) {

	self.connectAll()

	initRoute, err := self.getFirstRoute(self.nodeIdList)
	return initRoute, err
}

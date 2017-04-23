package nodemanager

//methods for testing purposes only

import (
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
)

func (self *NodeManager) CreateRandomNetwork(n int) []messages.NodeInterface {
	nodes := []messages.NodeInterface{}

	for i := 0; i < n; i++ {
		node, err := node.CreateAndConnectNode(messages.LOCALHOST+":"+strconv.Itoa(5000+i), messages.LOCALHOST+":5999")
		if err != nil {
			panic(err)
		}
		nodes = append(nodes, node)
	}
	self.rebuildRoutes()
	return nodes
}

func (self *NodeManager) CreateSequenceOfNodes(n int) (messages.Connection, messages.Connection) {
	/*
		This function creates a network with sequentially chained n nodes like 0-1-2-3-4-5-6-7-8-9 and returns the addresses of the first and last node
	*/

	nodeList := node.CreateNodeList(n)
	self.connectAll()
	self.rebuildRoutes()
	firstNode, lastNode := nodeList[0], nodeList[len(nodeList)-1]
	return firstNode.GetConnection(), lastNode.GetConnection()
}

func (self *NodeManager) CreateSequenceOfNodesAndBuildRoutes(n int) (cipher.PubKey, cipher.PubKey, messages.RouteId, messages.RouteId) {
	/*
		This function creates a network with sequentially chained n nodes like 0-1-2-3-4-5-6-7-8-9, builds route between the first and the last nodes in a chainand returns the addresses of them, a route from the first to the last one and a back route from the last to the first one
	*/

	node.CreateNodeList(n)
	self.connectAll()

	nodeList := self.nodeIdList

	route, backRoute, err := self.buildRoute(nodeList)
	if err != nil {
		panic(err)
	}
	clientNode, serverNode := nodeList[0], nodeList[len(nodeList)-1]
	return clientNode, serverNode, route, backRoute
}

func (self *NodeManager) CreateThreeRoutes() (messages.Connection, messages.Connection) {
	nodes := node.CreateNodeList(10)
	nodeList := self.nodeIdList
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

	clientNode, serverNode := nodes[0], nodes[9]
	return clientNode.GetConnection(), serverNode.GetConnection()
}

/* This function creates a network of n nodes randomly connected to each other */

func (self *NodeManager) connectAll() error {

	n := len(self.nodeIdList)

	for i := 0; i < n-1; i++ {
		id1, id2 := self.nodeIdList[i], self.nodeIdList[i+1]
		_, err := self.ConnectNodeToNode(id1, id2)
		if err != nil {
			panic(err)
			return err
		}
	}
	return nil
}

func (self *NodeManager) connectAllAndBuildRoute() (messages.RouteId, error) {

	err := self.connectAll()
	if err != nil {
		return messages.NIL_ROUTE, err
	}

	initRoute, err := self.getFirstRouteForward(self.nodeIdList)
	return initRoute, err
}

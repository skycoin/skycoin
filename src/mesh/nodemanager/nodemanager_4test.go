package nodemanager

//methods for testing purposes only

import (
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
)

func (self *NodeManager) CreateRandomNetwork(n, startPort int) []messages.NodeInterface {

	/* This function creates a network of n nodes randomly connected to each other */

	nodes := []messages.NodeInterface{}

	for i := 0; i < n; i++ {
		node, err := node.CreateAndConnectNode(&node.NodeConfig{"127.0.0.1:" + strconv.Itoa(startPort+i), []string{"127.0.0.1:5999"}, startPort + n + i, ""})
		if err != nil {
			panic(err)
		}
		nodes = append(nodes, node)
	}
	self.rebuildRoutes()
	return nodes
}

func (self *NodeManager) CreateSequenceOfNodes(n, startPort int) (messages.NodeInterface, messages.NodeInterface) {
	/*
		This function creates a network with sequentially chained n nodes like 0-1-2-3-4-5-6-7-8-9 and returns the first and last node
	*/

	nodeList := node.CreateNodeList(n, startPort)
	self.connectAll()
	self.rebuildRoutes()
	firstNode, lastNode := nodeList[0], nodeList[len(nodeList)-1]
	return firstNode, lastNode
}

func (self *NodeManager) CreateSequenceOfNodesAndBuildRoutes(n, startPort int) (cipher.PubKey, cipher.PubKey, messages.RouteId, messages.RouteId) {
	/*
		This function creates a network with sequentially chained n nodes like 0-1-2-3-4-5-6-7-8-9, builds route between the first and the last nodes in a chainand returns the addresses of them, a route from the first to the last one and a back route from the last to the first one
	*/

	node.CreateNodeList(n, startPort)
	self.connectAll()

	nodeList := self.nodeIdList

	route, backRoute, err := self.buildRoute(nodeList)
	if err != nil {
		panic(err)
	}
	clientNode, serverNode := nodeList[0], nodeList[len(nodeList)-1]
	return clientNode, serverNode, route, backRoute
}

func (self *NodeManager) CreateThreeRoutes(startPort int) (messages.NodeInterface, messages.NodeInterface) {
	nodes := node.CreateNodeList(10, startPort)
	nodeList := self.nodeIdList
	/*
		  1-2-3-4
		 /	 \
		0----5----9
		 \	 /
		  6_7_8_/
	*/
	self.connectNodeToNode(nodeList[0], nodeList[1])
	self.connectNodeToNode(nodeList[1], nodeList[2])
	self.connectNodeToNode(nodeList[2], nodeList[3])
	self.connectNodeToNode(nodeList[3], nodeList[4])
	self.connectNodeToNode(nodeList[4], nodeList[9])
	self.connectNodeToNode(nodeList[0], nodeList[5])
	self.connectNodeToNode(nodeList[5], nodeList[9])
	self.connectNodeToNode(nodeList[0], nodeList[6])
	self.connectNodeToNode(nodeList[6], nodeList[7])
	self.connectNodeToNode(nodeList[7], nodeList[8])
	self.connectNodeToNode(nodeList[8], nodeList[9])

	self.rebuildRoutes()

	clientNode, serverNode := nodes[0], nodes[9]
	return clientNode, serverNode
}

func (self *NodeManager) connectAll() error {

	n := len(self.nodeIdList)

	for i := 0; i < n-1; i++ {
		id1, id2 := self.nodeIdList[i], self.nodeIdList[i+1]
		_, err := self.connectNodeToNode(id1, id2)
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

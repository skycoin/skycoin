package nodemanager

import (
	"fmt"
	"math/rand"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

//contains a list of nodes
//calls tick on the list of nodes

//contains transport_mananger / transport_factory
//calls ticket methods on the transport factory
type NodeManager struct {
	nodeIdList           []cipher.PubKey
	nodeList             map[cipher.PubKey]*node.Node
	transportFactoryList []*transport.TransportFactory
	routeGraph           *RouteGraph
}

func NewNetwork() *NodeManager {
	nm := newNodeManager()
	return nm
}

func (self *NodeManager) AddNewNode() cipher.PubKey {
	nodeToAdd := node.NewNode()
	self.addNode(nodeToAdd)
	return nodeToAdd.Id
}

func (self *NodeManager) AddAndConnect() cipher.PubKey {
	nodeToAdd := node.NewNode()
	self.addNode(nodeToAdd)
	id := nodeToAdd.Id
	if len(self.nodeIdList) >= 2 {
		self.connectRandomly(id)
	}
	return id
}

func (self *NodeManager) CreateRandomNetwork(n int) []cipher.PubKey {
	nodes := []cipher.PubKey{}
	for i := 0; i < n; i++ {
		nodes = append(nodes, self.AddAndConnect())
	}
	self.rebuildRoutes()
	return nodes
}

func (self *NodeManager) connectRandomly(node0 cipher.PubKey) {
	var node1 cipher.PubKey
	for {
		node1 = self.getRandomNode()
		if node0 != node1 {
			break
		}
	}
	self.ConnectNodeToNode(node0, node1)

}

func (self *NodeManager) routeExists(pubkey0, pubkey1 cipher.PubKey) bool {
	_, exists := self.routeGraph.findRoute(pubkey0, pubkey1)
	return exists
}

func (self *NodeManager) Register(address cipher.PubKey, consumer messages.Consumer) error {
	node0, err := self.getNodeById(address)
	if err != nil {
		return err
	}
	node0.AssignConsumer(consumer)
	return nil
}

func (self *NodeManager) Tick() {
}

func newNodeManager() *NodeManager {
	nm := new(NodeManager)
	nm.nodeList = make(map[cipher.PubKey]*node.Node)
	nm.transportFactoryList = []*transport.TransportFactory{}
	nm.routeGraph = newGraph()
	return nm
}

func (self *NodeManager) getNodeById(id cipher.PubKey) (*node.Node, error) {
	result, found := self.nodeList[id]
	if !found {
		return &node.Node{}, errors.ERR_NODE_NOT_FOUND
	}
	return result, nil
}

func (self *NodeManager) addNode(nodeToAdd *node.Node) {
	id := nodeToAdd.Id
	self.nodeList[id] = nodeToAdd
	self.nodeIdList = append(self.nodeIdList, id)
}

func (self *NodeManager) getRandomNode() cipher.PubKey {
	list := self.nodeIdList
	max := len(list)
	index := rand.Intn(max)
	randomNode := list[index]
	return randomNode
}

func (self *NodeManager) connected(pubkey0, pubkey1 cipher.PubKey) bool {
	node0, err := self.getNodeById(pubkey0)
	if err != nil {
		return false
	}
	node1, err := self.getNodeById(pubkey1)
	if err != nil {
		return false
	}

	return node0.ConnectedTo(node1) && node1.ConnectedTo(node0)
}

func (self *NodeManager) ConnectNodeToNode(idA, idB cipher.PubKey) *transport.TransportFactory {

	if idA == idB {
		fmt.Println("Cannot connect node to itself")
		return &transport.TransportFactory{}
	}
	nodes := self.nodeList
	nodeA, found := nodes[idA]
	if !found {
		fmt.Println("Cannot find node with ID", idA)
		return &transport.TransportFactory{}
	}
	nodeB, found := nodes[idB]
	if !found {
		fmt.Println("Cannot find node with ID", idB)
		return &transport.TransportFactory{}
	}

	if nodeA.ConnectedTo(nodeB) || nodeB.ConnectedTo(nodeA) {
		fmt.Println("Nodes already connected")
		return &transport.TransportFactory{}
	}

	tf := transport.NewTransportFactory()
	tf.ConnectNodeToNode(nodeA, nodeB)
	self.transportFactoryList = append(self.transportFactoryList, tf)
	go tf.Tick()
	return tf
}

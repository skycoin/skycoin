package nodemanager

import (
	"fmt"

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

func (self *NodeManager) connectNodeToNode(idA, idB cipher.PubKey) *transport.TransportFactory {

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

package nodemanager

import (
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
	portDelivery         *PortDelivery
}

func NewNetwork() *NodeManager {
	nm := newNodeManager()
	return nm
}

func (self *NodeManager) AddNewNodeStub() cipher.PubKey {
	return self.AddNewNode(messages.LOCALHOST)
}

func (self *NodeManager) AddAndConnectStub() cipher.PubKey {
	return self.AddAndConnect(messages.LOCALHOST)
}

func (self *NodeManager) AddNewNode(host string) cipher.PubKey {
	nodeToAdd := self.newNode(host)
	return nodeToAdd.Id
}

func (self *NodeManager) AddAndConnect(host string) cipher.PubKey {
	id := self.AddNewNode(host)
	if len(self.nodeIdList) >= 2 {
		self.connectRandomly(id)
	}
	return id
}

func (self *NodeManager) newNode(host string) *node.Node {
	newNode := node.NewNode()

	newNode.Host = host

	self.addNode(newNode)
	return newNode
}

func (self *NodeManager) ConnectNodeToNode(idA, idB cipher.PubKey) (*transport.TransportFactory, error) {

	if idA == idB {
		return nil, errors.ERR_CONNECTED_TO_ITSELF
	}
	nodes := self.nodeList
	nodeA, found := nodes[idA]
	if !found {
		return nil, errors.ERR_NODE_NOT_FOUND
	}
	nodeB, found := nodes[idB]
	if !found {
		return nil, errors.ERR_NODE_NOT_FOUND
	}

	if nodeA.ConnectedTo(nodeB) || nodeB.ConnectedTo(nodeA) {
		return nil, errors.ERR_ALREADY_CONNECTED
	}

	nodeA.Port = self.portDelivery.Get(nodeA.Host)
	nodeB.Port = self.portDelivery.Get(nodeB.Host)

	tf := transport.NewTransportFactory()
	err := tf.ConnectNodeToNode(nodeA, nodeB)
	if err != nil {
		return nil, err
	}

	self.transportFactoryList = append(self.transportFactoryList, tf)
	tf.Tick()
	return tf, nil
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

func (self *NodeManager) Shutdown() {
	for _, tf := range self.transportFactoryList {
		tf.Shutdown()
	}
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

func newNodeManager() *NodeManager {
	nm := new(NodeManager)
	nm.nodeList = make(map[cipher.PubKey]*node.Node)
	nm.transportFactoryList = []*transport.TransportFactory{}
	nm.routeGraph = newGraph()
	nm.portDelivery = newPortDelivery()
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

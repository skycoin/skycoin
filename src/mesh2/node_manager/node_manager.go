package node_manager

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/node"
	"github.com/skycoin/skycoin/src/mesh2/transport"
)

//contains a list of nodes
//calls tick on the list of nodes

//contains transport_mananger / transport_factory
//calls ticket methods on the transport factory
type NodeManager struct {
	NodeList             *NodeListT
	TransportFactoryList []*transport.TransportFactory
}

type NodeListT struct {
	nodes map[cipher.PubKey]*node.Node
}

func NewNodeManager() *NodeManager {
	nm := new(NodeManager)
	nm.NodeList = &NodeListT{nodes: make(map[cipher.PubKey]*node.Node)}
	nm.TransportFactoryList = []*transport.TransportFactory{}
	return nm
}

func (self *NodeManager) GetNodeById(id cipher.PubKey) (*node.Node, error) {
	result, found := self.NodeList.nodes[id]
	if !found {
		return &node.Node{}, errors.New("Node not found")
	}
	return result, nil
}

func (self *NodeManager) AddNode() cipher.PubKey {
	nodeToAdd := node.NewNode()
	id := nodeToAdd.Id
	self.NodeList.nodes[id] = nodeToAdd
	return id
}

func (self *NodeManager) Tick() {
	self.NodeList.Tick()
}

func (self *NodeListT) Tick() {
	for _, node := range self.nodes {
		node.Tick()
	}
}

func (self *NodeManager) ConnectNodeToNode(idA, idB cipher.PubKey) *transport.TransportFactory {
	if idA == idB {
		fmt.Println("Cannot connect node to itself")
		return &transport.TransportFactory{}
	}
	nodes := self.NodeList.nodes
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

	tf := transport.NewTransportFactory()
	transportA, transportB := tf.CreateStubTransportPair()
	transportA.AttachedNode = nodeA
	tidA := transportA.Id
	transportB.AttachedNode = nodeB
	tidB := transportB.Id
	nodeA.Transports[tidA] = transportA
	nodeB.Transports[tidB] = transportB
	go tf.Tick()
	return tf
}

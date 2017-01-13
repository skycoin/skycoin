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
	NodeIdList           []cipher.PubKey
	NodeList             map[cipher.PubKey]*node.Node
	TransportFactoryList []*transport.TransportFactory
}

func NewNodeManager() *NodeManager {
	nm := new(NodeManager)
	nm.NodeList = make(map[cipher.PubKey]*node.Node)
	nm.TransportFactoryList = []*transport.TransportFactory{}
	return nm
}

func (self *NodeManager) GetNodeById(id cipher.PubKey) (*node.Node, error) {
	result, found := self.NodeList[id]
	if !found {
		return &node.Node{}, errors.New("Node not found")
	}
	return result, nil
}

func (self *NodeManager) AddNewNode() cipher.PubKey {
	nodeToAdd := node.NewNode()
	self.AddNode(nodeToAdd)
	return nodeToAdd.Id
}

func (self *NodeManager) AddNode(nodeToAdd *node.Node) {
	id := nodeToAdd.Id
	self.NodeList[id] = nodeToAdd
	self.NodeIdList = append(self.NodeIdList, id)
}

func (self *NodeManager) Tick() {
	for _, node := range self.NodeList {
		node.Tick()
	}
}

func (self *NodeManager) ConnectNodeToNode(idA, idB cipher.PubKey) *transport.TransportFactory {
	if idA == idB {
		fmt.Println("Cannot connect node to itself")
		return &transport.TransportFactory{}
	}
	nodes := self.NodeList
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
	tf.ConnectNodeToNode(nodeA, nodeB)
	self.TransportFactoryList = append(self.TransportFactoryList, tf)
	go tf.Tick()
	return tf
}

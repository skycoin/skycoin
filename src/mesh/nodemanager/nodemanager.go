package nodemanager

import (
	"math/rand"
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

//contains a list of nodes
//calls tick on the list of nodes

//contains transport_mananger / transport_factory
//calls ticket methods on the transport factory
type NodeManager struct {
	connectionList       map[cipher.PubKey]*Connection
	nodeIdList           []cipher.PubKey
	nodeList             map[cipher.PubKey]*node.Node
	transportFactoryList []*transport.TransportFactory
	routeGraph           *RouteGraph
	portDelivery         *PortDelivery
	lock                 *sync.Mutex
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

func (self *NodeManager) ConnectNodeToNode(idA, idB cipher.PubKey) (*transport.TransportFactory, error) {

	if idA == idB {
		return nil, messages.ERR_CONNECTED_TO_ITSELF
	}

	nodeA, err := self.getNodeById(idA)
	if err != nil {
		return nil, err
	}
	nodeB, err := self.getNodeById(idB)
	if err != nil {
		return nil, err
	}

	if nodeA.ConnectedTo(nodeB) || nodeB.ConnectedTo(nodeA) {
		return nil, messages.ERR_ALREADY_CONNECTED
	}

	nodeA.Port = self.portDelivery.Get(nodeA.Host)
	nodeB.Port = self.portDelivery.Get(nodeB.Host)

	tf := transport.NewTransportFactory()
	err = tf.ConnectNodeToNode(nodeA, nodeB)
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
	conn := self.connectionList[address]
	conn.AssignConsumer(consumer)
	node0.AssignUser(conn)
	return nil
}

func (self *NodeManager) Tick() {
}

func (self *NodeManager) Shutdown() {
	for _, tf := range self.transportFactoryList {
		tf.Shutdown()
	}
}

func newNodeManager() *NodeManager {
	nm := new(NodeManager)
	nm.nodeList = make(map[cipher.PubKey]*node.Node)
	nm.transportFactoryList = []*transport.TransportFactory{}
	nm.routeGraph = newGraph()
	nm.portDelivery = newPortDelivery()
	nm.connectionList = make(map[cipher.PubKey]*Connection)
	nm.lock = &sync.Mutex{}
	return nm
}

func (self *NodeManager) newNode(host string) *node.Node {
	newNode := node.NewNode()

	newNode.Host = host

	self.addNode(newNode)
	return newNode
}

func (self *NodeManager) addNode(nodeToAdd *node.Node) {
	id := nodeToAdd.Id
	self.lock.Lock()
	self.nodeList[id] = nodeToAdd
	self.nodeIdList = append(self.nodeIdList, id)
	self.lock.Unlock()
}

func (self *NodeManager) getNodeById(id cipher.PubKey) (*node.Node, error) {
	self.lock.Lock()
	result, found := self.nodeList[id]
	self.lock.Unlock()

	if !found {
		return &node.Node{}, messages.ERR_NODE_NOT_FOUND
	}
	return result, nil
}

func (self *NodeManager) GetAllNodes() map[cipher.PubKey]*node.Node {
	return self.nodeList
}

func (self *NodeManager) GetNodeById(id cipher.PubKey) (*node.Node, error) {
	n, err := self.getNodeById(id)
	return n, err
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
	_, err := self.routeGraph.findRoute(pubkey0, pubkey1)
	return err == nil
}

func (self *NodeManager) GetTicks() int {
	ticks := 0
	for _, n := range self.nodeList {
		ticks += n.GetTicks()
	}
	return ticks
}

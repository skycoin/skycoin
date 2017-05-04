package nodemanager

// make nodemanager an app?

import (
	"math/rand"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

//contains a list of nodes
//calls tick on the list of nodes

//contains transport_mananger / transport_factory
//calls ticket methods on the transport factory
type NodeManager struct {
	nodeIdList           []cipher.PubKey
	nodeList             map[cipher.PubKey]*NodeRecord
	transportFactoryList []*TransportFactory
	nodesByTransport     map[messages.TransportId]cipher.PubKey
	routeGraph           *RouteGraph
	portDelivery         *PortDelivery
	msgServer            *MsgServer
	lock                 *sync.Mutex
}

var config = messages.GetConfig()

func NewNetwork() *NodeManager {
	nm := newNodeManager()
	return nm
}

func (self *NodeManager) addNewNode(host string) (cipher.PubKey, error) { //**** will be called by messaging server, response will be the reply
	nodeToAdd, err := self.newNode(host)
	if err != nil {
		return cipher.PubKey{}, err
	}
	return nodeToAdd.id, nil
}

func (self *NodeManager) addAndConnect(host string) (cipher.PubKey, error) { //**** will be called by messaging server, response will be the reply
	id, err := self.addNewNode(host)
	if err != nil {
		return cipher.PubKey{}, err
	}
	if len(self.nodeIdList) >= 2 {
		self.connectRandomly(id)
	}
	return id, nil
}

func (self *NodeManager) connectNodeToNode(idA, idB cipher.PubKey) (*TransportFactory, error) {

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

	nodeA.port = self.portDelivery.Get(nodeA.host)
	portACM := messages.AssignPortCM{nodeA.port}
	portACMS := messages.Serialize(messages.MsgAssignPortCM, portACM)

	nodeB.port = self.portDelivery.Get(nodeB.host)
	portBCM := messages.AssignPortCM{nodeB.port}
	portBCMS := messages.Serialize(messages.MsgAssignPortCM, portBCM)

	err = nodeA.sendToNode(portACMS)
	if err != nil {
		return nil, err
	}

	err = nodeB.sendToNode(portBCMS)
	if err != nil {
		return nil, err
	}

	tf := newTransportFactory()
	err = tf.connectNodeToNode(nodeA, nodeB)
	if err != nil {
		panic(err)
		return nil, err
	}

	self.transportFactoryList = append(self.transportFactoryList, tf)
	tf.tick()
	return tf, nil
}

func (self *NodeManager) connectWithRoute(nodeFromId, nodeToId cipher.PubKey, appIdFrom, appIdTo messages.AppId) (messages.ConnectionId, error) {

	connectionId := messages.RandConnectionId()

	routeId, backRouteId, err := self.findRoute(nodeFromId, nodeToId)
	if err != nil {
		return messages.ConnectionId(0), err
	}

	assignConnectionFrom := messages.AssignConnectionCM{
		connectionId,
		routeId,
		appIdFrom,
	}
	assignConnectionFromS := messages.Serialize(messages.MsgAssignConnectionCM, assignConnectionFrom)

	assignConnectionTo := messages.AssignConnectionCM{
		connectionId,
		backRouteId,
		appIdTo,
	}
	assignConnectionToS := messages.Serialize(messages.MsgAssignConnectionCM, assignConnectionTo)

	nodeFrom, err := self.getNodeById(nodeFromId)
	if err != nil {
		return messages.ConnectionId(0), err
	}

	nodeTo, err := self.getNodeById(nodeToId)
	if err != nil {
		return messages.ConnectionId(0), err
	}

	err = nodeFrom.sendToNode(assignConnectionFromS)
	if err != nil {
		return messages.ConnectionId(0), err
	}

	err = nodeTo.sendToNode(assignConnectionToS)
	if err != nil {
		return messages.ConnectionId(0), err
	}

	connectionFrom := messages.ConnectionOnCM{
		nodeFrom.id,
		connectionId,
	}

	connectionFromS := messages.Serialize(messages.MsgConnectionOnCM, connectionFrom)
	err = nodeFrom.sendToNode(connectionFromS)

	if err != nil {
		return messages.ConnectionId(0), err
	}

	connectionTo := messages.ConnectionOnCM{
		nodeTo.id,
		connectionId,
	}

	connectionToS := messages.Serialize(messages.MsgConnectionOnCM, connectionTo)

	err = nodeTo.sendToNode(connectionToS)
	if err != nil {
		return messages.ConnectionId(0), err
	}

	return connectionId, nil
}

func (self *NodeManager) Tick() {
}

func (self *NodeManager) Shutdown() {
	for _, n := range self.nodeList {
		n.shutdown()
	}
	self.msgServer.shutdown()
	time.Sleep(1 * time.Millisecond)
}

func newNodeManager() *NodeManager {
	nm := new(NodeManager)
	nm.nodeList = make(map[cipher.PubKey]*NodeRecord)
	nm.transportFactoryList = []*TransportFactory{}
	nm.routeGraph = newGraph()
	nm.portDelivery = newPortDelivery()
	msgServer, err := newMsgServer(nm)
	if err != nil {
		panic(err)
	}
	nm.msgServer = msgServer
	nm.lock = &sync.Mutex{}
	return nm
}

func (self *NodeManager) getNodeById(id cipher.PubKey) (*NodeRecord, error) { // resolve it
	result, found := self.nodeList[id]

	if !found {
		return &NodeRecord{}, messages.ERR_NODE_NOT_FOUND
	}
	return result, nil
}

func (self *NodeManager) GetAllNodes() map[cipher.PubKey]*NodeRecord {
	return self.nodeList
}

func (self *NodeManager) GetNodeById(id cipher.PubKey) (*NodeRecord, error) {
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
	self.connectNodeToNode(node0, node1)

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

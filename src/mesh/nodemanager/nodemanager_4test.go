package nodemanager

//methods for testing purposes only

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

func (self *NodeManager) CreateNodeList(n int) []cipher.PubKey {
	nodes := []cipher.PubKey{}
	for i := 0; i < n; i++ {
		nodeId := self.AddNewNode()
		nodes = append(nodes, nodeId)
	}
	return nodes
}

func (self *NodeManager) ConnectAll() (messages.RouteId, error) {

	n := len(self.nodeIdList)

	for i := 0; i < n-1; i++ {
		id1, id2 := self.nodeIdList[i], self.nodeIdList[i+1]
		self.ConnectNodeToNode(id1, id2)
	}

	initRoute, err := self.getFirstRoute(self.nodeIdList)
	return initRoute, err
}

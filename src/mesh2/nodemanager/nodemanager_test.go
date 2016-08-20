package nodemanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNodeList(t *testing.T) {
	nodeManager := &NodeManager{}
	nodeManager.CreateNodeConfigList(4)
	assert.Len(t, nodeManager.ConfigList, 4, "Error expected 4 nodes")
	assert.Len(t, nodeManager.ConfigList[0].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 1")
}

func TestConnectNodes(t *testing.T) {
	nodeManager := &NodeManager{}
	nodeManager.CreateNodeConfigList(5)
	assert.Len(t, nodeManager.ConfigList, 5, "Error expected 5 nodes")
	assert.Len(t, nodeManager.ConfigList[0].PeersToConnect, 0, "Error expected 0 PeersToConnect from Node 1")
	nodeManager.ConnectNodes()
	assert.Len(t, nodeManager.ConfigList[0].PeersToConnect, 1, "Error expected 1 PeersToConnect from Node 1")
	assert.Len(t, nodeManager.ConfigList[1].PeersToConnect, 2, "Error expected 2 PeersToConnect from Node 1")
	assert.Len(t, nodeManager.ConfigList[4].PeersToConnect, 1, "Error expected 1 PeersToConnect from Node 1")
}

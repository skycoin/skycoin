package node_manager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddingNodes(t *testing.T) {

	nm := NewNodeManager()
	assert.Len(t, nm.NodeList.nodes, 0, "Error expected 0 nodes")
	nm.AddNode()
	assert.Len(t, nm.NodeList.nodes, 1, "Error expected 1 nodes")
	nm.AddNode()
	nm.AddNode()
	nm.AddNode()
	nm.AddNode()
	assert.Len(t, nm.NodeList.nodes, 5, "Error expected 5 nodes")
}

func TestConnectTwoNodes(t *testing.T) {

	nm := NewNodeManager()
	id1 := nm.AddNode()
	id2 := nm.AddNode()
	node1, err := nm.GetNodeById(id1)
	assert.Nil(t, err)
	node2, err := nm.GetNodeById(id2)
	assert.Nil(t, err)
	tf := nm.ConnectNodeToNode(id1, id2)
	t1, t2 := tf.GetTransports()
	assert.Len(t, node1.Transports, 1, "Error expected 1 transport")
	assert.Len(t, node2.Transports, 1, "Error expected 1 transport")
	assert.Equal(t, t1.Id, t2.StubPair.Id)
	assert.Equal(t, t2.Id, t1.StubPair.Id)
	tr1, err := node1.GetTransportToNode(id2)
	assert.Nil(t, err)
	assert.Equal(t, tr1.StubPair.Id, t2.Id)
}

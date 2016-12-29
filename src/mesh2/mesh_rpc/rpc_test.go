package mesh_rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skycoin/src/mesh2/node"
)

func TestRPC(t *testing.T) {
	rpcInstance := NewRPC()
	rpcInstance.Start()
	node1 := node.NewNode()
	assert.Equal(t, 0, node1.NumControlChannels())
	rpcInstance.CreateControlChannel(node1)
	assert.Equal(t, 1, node1.NumControlChannels())
}

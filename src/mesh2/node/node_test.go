package node

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

func TestCreateControlChannel(t *testing.T) {
	node := NewNode()
	ccid := uuid.UUID{}

	msg := messages.CreateChannelControlMessage{}
	node.HandleControlMessage(ccid, msg)
	assert.Len(t, node.ControlChannels, 2, "Should be 2 control channels")
}

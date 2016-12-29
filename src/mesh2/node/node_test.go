package node

import (
	"fmt"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

func TestCreateControlChannel(t *testing.T) {
	panic("ok")
	node := NewNode()
	ccid := uuid.UUID{}

	msg := messages.CreateChannelControlMessage{}
	node.HandleControlMessage(ccid, msg)
	assert.Len(t, node.controlChannels, 2, "Should be 2 control channels")
}

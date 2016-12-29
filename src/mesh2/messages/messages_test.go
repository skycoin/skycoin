package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skycoin/src/cipher"
)

func TestSerialize(t *testing.T) {

	routeId := RandRouteId()
	nodeId, _ := cipher.GenerateKeyPair()
	msg := AddRouteControlMessage{nodeId, routeId}
	serialized := Serialize((uint16)(MsgAddRouteControlMessage), msg)
	msg1 := AddRouteControlMessage{}
	err := Deserialize(serialized, &msg1)
	assert.Nil(t, err)
	assert.Equal(t, msg.NodeId, msg1.NodeId)
	assert.Equal(t, msg.RouteId, msg1.RouteId)

	msg2 := RemoveRouteControlMessage{routeId}
	serialized = Serialize((uint16)(MsgRemoveRouteControlMessage), msg2)
	msg3 := RemoveRouteControlMessage{}
	err = Deserialize(serialized, &msg3)
	assert.Nil(t, err)
	assert.Equal(t, msg2.RouteId, msg3.RouteId)
}

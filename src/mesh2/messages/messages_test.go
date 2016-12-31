package messages

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
)

func TestSerialize(t *testing.T) {

	routeId := RandRouteId()
	transportId := RandTransportId()
	datagram := []byte{'t', 'e', 's', 't'}

	msg := InRouteMessage{transportId, routeId, datagram}
	serialized := Serialize((uint16)(MsgInRouteMessage), msg)
	msg1 := InRouteMessage{}
	err := Deserialize(serialized, &msg1)
	assert.Nil(t, err)
	assert.Equal(t, msg.TransportId, msg1.TransportId)
	assert.Equal(t, msg.RouteId, msg1.RouteId)
	assert.Equal(t, msg.Datagram, msg1.Datagram)

	sequence := (uint32)(rand.Intn(65536))
	msg2 := TransportDatagramTransfer{sequence, datagram}
	serialized = Serialize((uint16)(MsgTransportDatagramTransfer), msg2)
	msg3 := TransportDatagramTransfer{}
	err = Deserialize(serialized, &msg3)
	assert.Nil(t, err)
	assert.Equal(t, msg2.Sequence, msg3.Sequence)
	assert.Equal(t, msg2.Datagram, msg3.Datagram)

	nodeId, _ := cipher.GenerateKeyPair()
	msg4 := AddRouteControlMessage{nodeId, routeId}
	serialized = Serialize((uint16)(MsgAddRouteControlMessage), msg4)
	msg5 := AddRouteControlMessage{}
	err = Deserialize(serialized, &msg5)
	assert.Nil(t, err)
	assert.Equal(t, msg4.NodeId, msg5.NodeId)
	assert.Equal(t, msg4.RouteId, msg5.RouteId)
}

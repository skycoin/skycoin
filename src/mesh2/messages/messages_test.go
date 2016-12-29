package messages

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
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
}

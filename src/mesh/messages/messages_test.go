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
	msg2 := TransportDatagramTransfer{RandRouteId(), sequence, datagram}
	serialized = Serialize((uint16)(MsgTransportDatagramTransfer), msg2)
	msg3 := TransportDatagramTransfer{}
	err = Deserialize(serialized, &msg3)
	assert.Nil(t, err)
	assert.Equal(t, msg2.RouteId, msg3.RouteId)
	assert.Equal(t, msg2.Sequence, msg3.Sequence)
	assert.Equal(t, msg2.Datagram, msg3.Datagram)

	route1Id := RandRouteId()
	transport1Id := RandTransportId()

	msg4 := AddRouteCM{transportId, transport1Id, routeId, route1Id}
	serialized = Serialize((uint16)(MsgAddRouteCM), msg4)
	msg5 := AddRouteCM{}
	err = Deserialize(serialized, &msg5)
	assert.Nil(t, err)
	assert.Equal(t, msg4.IncomingTransportId, msg5.IncomingTransportId)
	assert.Equal(t, msg4.IncomingRouteId, msg5.IncomingRouteId)
	assert.Equal(t, msg4.OutgoingTransportId, msg5.OutgoingTransportId)
	assert.Equal(t, msg4.OutgoingRouteId, msg5.OutgoingRouteId)
}

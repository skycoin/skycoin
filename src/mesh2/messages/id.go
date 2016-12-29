package messages

import (
	"math/rand"
)

//Node: RouteId
type RouteId uint64

func RandRouteId() RouteId {
	return (RouteId)(rand.Int63())
}

//Transport: TransportId

type TransportId uint64

func RandTransportId() TransportId {
	return (TransportId)(rand.Int63())
}

type ControlChannelId uint64

func RandCCId() ControlChannelId {
	return (ControlChannelId)(rand.Int63())
}

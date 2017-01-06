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

type ChannelId uint64

func RandChannelId() ChannelId {
	return (ChannelId)(rand.Int63())
}

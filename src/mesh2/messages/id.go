package messages

import (
	"math/rand"
)

//Node
type RouteId uint64

func RandRouteId() TransportId {
	return (TransportId)(rand.Int63())
}

//Transport

type TransportId uint64

func RandTransportId() TransportId {
	return (TransportId)(rand.Int63())
}

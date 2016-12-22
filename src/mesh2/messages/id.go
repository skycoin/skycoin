package messages

import (
	"math/rand"
)

//Node: NodeId

type NodeId []byte

//Node: RouteId
type RouteId uint64

func RandRouteId() TransportId {
	return (TransportId)(rand.Int63())
}

//Transport: TransportId

type TransportId uint64

func RandTransportId() TransportId {
	return (TransportId)(rand.Int63())
}

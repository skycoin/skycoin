package messages

type RouteRule struct {
	IncomingTransport TransportId
	OutgoingTransport TransportId
	IncomingRoute     RouteId
	OutgoingRoute     RouteId
}

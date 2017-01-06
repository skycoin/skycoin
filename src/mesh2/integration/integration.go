package main

import (
	"time"

	"github.com/skycoin/skycoin/src/mesh2/messages"
	"github.com/skycoin/skycoin/src/mesh2/node"
	nodemanager "github.com/skycoin/skycoin/src/mesh2/node_manager"
)

func main() {

	nm := nodemanager.NewNodeManager()

	id1, id2, id3 := nm.AddNode(), nm.AddNode(), nm.AddNode()
	node1, err := nm.GetNodeById(id1)
	if err != nil {
		panic(err)
	}
	node2, err := nm.GetNodeById(id2)
	if err != nil {
		panic(err)
	}
	node3, err := nm.GetNodeById(id3)
	if err != nil {
		panic(err)
	}

	tf12 := nm.ConnectNodeToNode(id1, id2) // transport pair between nodes 1 and 2
	t12, t21 := tf12.GetTransports()
	tid12, tid21 := t12.Id, t21.Id

	tf13 := nm.ConnectNodeToNode(id1, id3) // transport pair between nodes 1 and 3
	t13, t31 := tf13.GetTransports()
	tid13, tid31 := t13.Id, t31.Id

	tf23 := nm.ConnectNodeToNode(id2, id3) // transport pair between nodes 2 and 3
	t23, t32 := tf23.GetTransports()
	tid23, tid32 := t23.Id, t32.Id

	t13.MaxSimulatedDelay = 1100 // a possible point of failure, going through the transport can lead to packet drop

	routeId123, routeId231, routeId312 := messages.RandRouteId(), messages.RandRouteId(), messages.RandRouteId()

	route123 := node.RouteRule{tid21, tid23, routeId123, routeId231} // route from 1 to 3 through 2
	route231 := node.RouteRule{tid32, tid31, routeId231, routeId312} // route from 2 to 1 through 3
	route312 := node.RouteRule{tid13, tid12, routeId312, routeId123} // route from 3 to 2 through 1

	node1.RouteForwardingRules[routeId312] = &route312
	node2.RouteForwardingRules[routeId123] = &route123
	node3.RouteForwardingRules[routeId231] = &route231

	nm.Tick()

	inRouteMessage := messages.InRouteMessage{tid21, routeId123, []byte{'t', 'e', 's', 't'}}
	serialized := messages.Serialize(messages.MsgInRouteMessage, inRouteMessage)
	node2.IncomingChannel <- serialized
	time.Sleep(10 * time.Second)
}

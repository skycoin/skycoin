package nodemanager

import (
	"github.com/skycoin/skycoin/src/cipher"
)

/*
	Service for finding routes
	- records the network topology
*/

//TODO:
//	Add network topology service
//  - keep network topology graph
//  - find distance between nodes and return multihop routes

// Find routes from the connections from a node
// One hop direct routes
// WTF does this do?
func (self *NodeManager) FindRoute(pubKey1 cipher.PubKey) {
	config1 := self.ConfigList[pubKey1]
	for _, v := range config1.PeerToPeers {
		pubKey2 := v.Peer
		route := Route{}
		route.SourceNode = pubKey1
		route.TargetNode = pubKey2
		route.Weight = 1
		route.RoutesToEstablish = append(route.RoutesToEstablish, pubKey2)

		routeKey := RouteKey{SourceNode: pubKey1, TargetNode: pubKey2}
		self.Routes[routeKey] = route
		self.FindIndirectRoutes(route)
	}
}

// Find indirect routes from a route
// WTF does this do?
func (self *NodeManager) FindIndirectRoutes(route Route) {
	config1 := self.ConfigList[route.TargetNode]
	for _, v := range config1.PeerToPeers {
		pubKey2 := v.Peer
		route2 := Route{}
		route2.SourceNode = route.SourceNode
		route2.TargetNode = pubKey2
		route2.RoutesToEstablish = append(route2.RoutesToEstablish, route.RoutesToEstablish...)
		route2.RoutesToEstablish = append(route2.RoutesToEstablish, pubKey2)
		route2.Weight = route.Weight + 1

		routeKey := RouteKey{SourceNode: route.SourceNode, TargetNode: pubKey2}
		if _, ok := self.Routes[routeKey]; !ok {

			self.Routes[routeKey] = route
			self.Routes[routeKey] = route
			self.FindIndirectRoutes(route)
		}
	}
}

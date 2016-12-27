package nodemanager

import (
	"github.com/satori/go.uuid"

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

func (s *NodeManager) FindRoutes(pubKey1, pubKey2 cipher.PubKey) bool {
	_, found := s.FindRoute(pubKey1, pubKey2)
	if found {
		_, found = s.FindRoute(pubKey2, pubKey1)
	}
	return found
}

func (s *NodeManager) FindRoute(pubKey1, pubKey2 cipher.PubKey) (*RouteConfig, bool) {
	config1 := s.ConfigList[pubKey1]
	routeKey := RouteKey{SourceNode: pubKey1, TargetNode: pubKey2}
	//	if existing, found := s.Routes[routeKey]; found { return existing, true } //??? do we need this cache and should we add expiration conditions?
	peers, found := s.RouteGraph.FindRoute(pubKey1, pubKey2)
	if !found {
		return nil, false
	}
	routeConfig := &RouteConfig{Peers: peers}
	routeConfig.RouteID = uuid.NewV4()
	config1.RoutesConfigsToEstablish = append(config1.RoutesConfigsToEstablish, *routeConfig)
	s.Routes[routeKey] = routeConfig
	return routeConfig, true
}

func (s *NodeManager) RebuildRouteGraph() {
	s.RouteGraph.Clear()
	for _, config := range s.ConfigList {
		nodeFrom := config.NodeConfig.PubKey
		for _, peerToPeer := range config.PeerToPeers {
			s.RouteGraph.AddDirectRoute(nodeFrom, peerToPeer.Peer, 1) // weight is always 1 because so far all routes are equal! Change this if needed
		}
	}
}

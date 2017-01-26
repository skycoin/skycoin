package node

import (
	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

//A Node has a map of route rewriting rules
//A Node has a control channel for setting and modifying the route rewrite rules
//A Node has a list of transports

//Route rewriting rules
//-nodes receive messages on a route
//-nodes look up the route in a table and if it has a rewrite rule, rewrites the route
// and forwards it to the transport

func (self *Node) addRoute(routeRule *RouteRule) error {

	routeId := routeRule.IncomingRoute

	if _, ok := self.RouteForwardingRules[routeId]; ok {
		err := errors.ERR_ROUTE_EXISTS
		return err
	}

	self.RouteForwardingRules[routeId] = routeRule
	return nil
}

func (self *Node) removeRoute(routeId messages.RouteId) error {
	if _, ok := self.RouteForwardingRules[routeId]; !ok {
		err := errors.ERR_ROUTE_DOESNT_EXIST
		return err
	}

	delete(self.RouteForwardingRules, routeId)
	return nil
}

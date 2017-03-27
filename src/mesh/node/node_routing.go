package node

import (
	"github.com/skycoin/skycoin/src/mesh/messages"
)

func (self *Node) getRoute(routeId messages.RouteId) (*RouteRule, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	routeRule, ok := self.RouteForwardingRules[routeId]
	if !ok {
		return nil, messages.ERR_ROUTE_DOESNT_EXIST
	}
	return routeRule, nil
}

func (self *Node) addRoute(routeRule *RouteRule) error {
	routeId := routeRule.IncomingRoute

	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.RouteForwardingRules[routeId]; ok {
		return messages.ERR_ROUTE_EXISTS
	}
	self.RouteForwardingRules[routeId] = routeRule
	return nil
}

func (self *Node) removeRoute(routeId messages.RouteId) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.RouteForwardingRules[routeId]; !ok {
		return messages.ERR_ROUTE_DOESNT_EXIST
	}

	delete(self.RouteForwardingRules, routeId)
	return nil
}

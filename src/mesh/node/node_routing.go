package node

import (
	"github.com/skycoin/skycoin/src/mesh/messages"
)

func (self *Node) getRoute(routeId messages.RouteId) (*messages.RouteRule, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	routeRule, ok := self.routeForwardingRules[routeId]
	if !ok {
		return nil, messages.ERR_ROUTE_DOESNT_EXIST
	}
	return routeRule, nil
}

func (self *Node) addRoute(routeRule *messages.RouteRule) error {
	routeId := routeRule.IncomingRoute

	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.routeForwardingRules[routeId]; ok {
		return messages.ERR_ROUTE_EXISTS
	}
	self.routeForwardingRules[routeId] = routeRule
	return nil
}

func (self *Node) removeRoute(routeId messages.RouteId) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.routeForwardingRules[routeId]; !ok {
		return messages.ERR_ROUTE_DOESNT_EXIST
	}

	delete(self.routeForwardingRules, routeId)
	return nil
}

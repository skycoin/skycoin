package node

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/messages"

	"errors"
	"fmt"
)

//A Node has a map of route rewriting rules
//A Node has a control channel for setting and modifying the route rewrite rules
//A Node has a list of transports

//Route rewriting rules
//-nodes receive messages on a route
//-nodes look up the route in a table and if it has a rewrite rule, rewrites the route
// and forwards it to the transport

func (self *Node) addRoute(nodeTo cipher.PubKey, routeId messages.RouteId) error {
	if _, ok := self.RouteForwardingRules[routeId]; ok {
		err := errors.New("Route already exists")
		fmt.Println(err)
		return err
	}

	outgoingTransport, err := self.GetTransportToNode(nodeTo)
	if err != nil {
		err := errors.New("No transport to node")
		fmt.Println(err)
		return err
	}

	routeRule := &RouteRule{
		IncomingTransport: (messages.TransportId)(0),
		IncomingRoute:     (messages.RouteId)(0),
		OutgoingTransport: outgoingTransport.Id,
		OutgoingRoute:     routeId,
	}
	self.RouteForwardingRules[routeId] = routeRule
	return nil
}

func (self *Node) removeRoute(routeId messages.RouteId) error {
	if _, ok := self.RouteForwardingRules[routeId]; !ok {
		err := errors.New("Route doesn't exist")
		fmt.Println(err)
		return err
	}

	delete(self.RouteForwardingRules, routeId)
	return nil
}

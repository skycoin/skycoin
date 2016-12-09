package mesh

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

var NilRouteID domain.RouteID = (domain.RouteID)(uuid.Nil)

// toPeer must be the public key of a connected peer
func (self *Node) AddRoute(routeID domain.RouteID, peerID cipher.PubKey) error {
	_, routeExists := self.safelyGetRoute(routeID)
	if routeExists {
		return errors.New(fmt.Sprintf("Route %v already exists\n", routeID))
	}
	transportToPeer := self.safelyGetTransportToPeer(peerID)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", peerID))
	}

	// Add locally to routes for backward termination
	self.lock.Lock()
	defer self.lock.Unlock()
	self.routes[routeID] =
		domain.Route{
			ForwardToPeerID:   peerID,
			ForwardToRouteID:  routeID,
			BackwardToPeerID:  cipher.PubKey{},
			BackwardToRouteID: NilRouteID,
			// Route lifetime: never dies until route is removed
			ExpiryTime: time.Unix(0, 0),
		}

	self.localRoutesByTerminatingPeer[peerID] = routeID
	self.localRoutes[routeID] = domain.LocalRoute{
		LastForwardingPeerID: self.Config.PubKey,
		TerminatingPeerID:    peerID,
		LastHopRouteID:       NilRouteID,
		LastConfirmed:        time.Unix(0, 0),
	}
	return nil
}

func (self *Node) DeleteRoute(routeID domain.RouteID) (err error) {
	route, routeExists := self.safelyGetRoute(routeID)
	if !routeExists {
		return errors.New(fmt.Sprintf("Cannot delete route %v, doesn't exist\n", routeID))
	}

	err = self.sendDeleteRoute(routeID, route)

	self.lock.Lock()
	defer self.lock.Unlock()
	localRoute, localExists := self.localRoutes[routeID]

	delete(self.routes, routeID)
	delete(self.routeExtensionsAwaitingConfirm, routeID)

	if localExists {
		delete(self.localRoutes, routeID)
		delete(self.localRoutesByTerminatingPeer, localRoute.TerminatingPeerID)
	}
	return err
}

// toPeer must be the public key of a peer connected to the current last node in this route
// Blocks until the set route reply is received or the timeout is reached
func (self *Node) ExtendRoute(routeID domain.RouteID, toPeer cipher.PubKey, timeout time.Duration) (err error) {
	message, directPeer, confirm, extendError := self.extendRouteWithoutSending(routeID, toPeer)
	if extendError != nil {
		return extendError
	}
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v from %v\n", directPeer, self.Config.PubKey))
	}
	serialized := self.serializer.SerializeMessage(message)
	err = transportToPeer.SendMessage(directPeer, serialized, nil)
	if err != nil {
		return err
	}
	select {
	case <-confirm:
		{
			break
		}
	case <-time.After(timeout):
		{
			// Still need to remove from confirm map
			err = &TimeoutError{}
		}
	}
	self.lock.Lock()
	defer self.lock.Unlock()

	delete(self.routeExtensionsAwaitingConfirm, routeID)
	return
}

func (self *Node) GetRouteLastConfirmed(routeID domain.RouteID) (time.Time, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	localRoute, found := self.localRoutes[routeID]
	if !found {
		return time.Unix(0, 0), errors.New(fmt.Sprintf("Unknown route id: %v", routeID))
	}
	return localRoute.LastConfirmed, nil
}

func (self *Node) findRouteToPeer(toPeer cipher.PubKey) (directPeerID cipher.PubKey, localID, sendID domain.RouteID, transport transport.ITransport, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	localRouteID, ok := self.localRoutesByTerminatingPeer[toPeer]
	if ok {
		route, ok := self.routes[localRouteID]
		if !ok {
			panic("Local route struct does not exists")
		}
		directPeerID = route.ForwardToPeerID
		localID = localRouteID
		sendID = route.ForwardToRouteID
	} else {
		return cipher.PubKey{}, NilRouteID, NilRouteID, nil, errors.New(fmt.Sprintf("No route to peer: %v", toPeer))
	}
	transport = self.GetTransportToPeer(directPeerID)
	if transport == nil {
		return cipher.PubKey{}, NilRouteID, NilRouteID, nil,
			errors.New(fmt.Sprintf("No route or transport to peer %v\n", toPeer))
	}
	return
}

func (self *Node) extendRouteWithoutSending(routeID domain.RouteID, toPeer cipher.PubKey) (message domain.SetRouteMessage, directPeer cipher.PubKey, waitConfirm chan bool, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, alreadyExtending := self.routeExtensionsAwaitingConfirm[routeID]
	if alreadyExtending {
		return domain.SetRouteMessage{}, cipher.PubKey{}, nil, errors.New("Cannot extend route more than once at the same time")
	}

	newHopId := routeID

	localRoute, localRouteExists := self.localRoutes[routeID]
	if !localRouteExists {
		return domain.SetRouteMessage{}, cipher.PubKey{}, nil, errors.New(fmt.Sprintf("ExtendRoute called on unknown route: %v", routeID))
	}

	route, routeExists := self.routes[routeID]
	if !routeExists {
		panic("Internal consistency error 8JUL2016544")
	}

	directPeer = route.ForwardToPeerID

	sendBase := domain.MessageBase{
		SendRouteID: route.ForwardToRouteID,
		SendBack:    false,
		FromPeerID:  self.Config.PubKey,
		Nonce:       GenerateNonce(),
	}

	newTermMessage := domain.SetRouteMessage{
		MessageBase:                sendBase,
		SetRouteID:                 routeID,
		ConfirmRouteID:             routeID,
		ForwardToPeerID:            toPeer,
		ForwardRewriteSendRouteID:  routeID,
		BackwardToPeerID:           localRoute.LastForwardingPeerID,
		BackwardRewriteSendRouteID: localRoute.LastHopRouteID,
		DurationHint:               3 * self.Config.RefreshRouteDuration,
	}

	delete(self.localRoutesByTerminatingPeer, localRoute.TerminatingPeerID)
	self.localRoutes[routeID] = domain.LocalRoute{
		LastForwardingPeerID: localRoute.TerminatingPeerID,
		TerminatingPeerID:    toPeer,
		LastHopRouteID:       newHopId,
		LastConfirmed:        localRoute.LastConfirmed,
	}
	self.localRoutesByTerminatingPeer[toPeer] = routeID

	updatedRoute := route
	updatedRoute.ForwardToRouteID = newHopId
	self.routes[routeID] = updatedRoute
	confirmChan := make(chan bool, 1)
	self.routeExtensionsAwaitingConfirm[routeID] = confirmChan

	return newTermMessage, directPeer, confirmChan, nil
}

func (self *Node) sendDeleteRoute(routeID domain.RouteID, route domain.Route) error {
	sendBase := domain.MessageBase{
		SendRouteID: route.ForwardToRouteID,
		SendBack:    false,
		FromPeerID:  self.Config.PubKey,
		Nonce:       GenerateNonce(),
	}
	message := domain.DeleteRouteMessage{
		MessageBase: sendBase,
	}

	directPeer := route.ForwardToPeerID
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("2No transport to peer %v from %v\n", directPeer, self.Config.PubKey))
	}
	serialized := self.serializer.SerializeMessage(message)
	err := transportToPeer.SendMessage(directPeer, serialized, nil)
	if err != nil {
		return err
	}

	return nil
}

func (self *Node) expireOldRoutesLoop() {
	self.closeGroup.Add(1)
	defer self.closeGroup.Done()
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.Config.ExpireRoutesInterval):
			{
				self.expireOldRoutes()
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Node) refreshRoutesLoop() {
	self.closeGroup.Add(1)
	defer self.closeGroup.Done()
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.Config.RefreshRouteDuration):
			{
				self.refreshRoutes()
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Node) safelyGetForwarding(msg interface{}) (SendBack bool, route domain.Route, doForward bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	messageBase := getMessageBase(msg)

	routeFound, foundRoute := self.routes[messageBase.SendRouteID]

	doForward = foundRoute

	if messageBase.SendBack {
		if routeFound.BackwardToPeerID == (cipher.PubKey{}) {
			doForward = false
		}
	} else {
		if routeFound.ForwardToPeerID == (cipher.PubKey{}) {
			doForward = false
		}
	}

	if doForward {
		return messageBase.SendBack, routeFound, doForward
	} else {
		return false, domain.Route{}, doForward
	}
}

func (self *Node) safelyGetRoute(routeID domain.RouteID) (route domain.Route, ok bool) {
	route, ok = self.routes[routeID]
	return route, ok
}

func (self *Node) expireOldRoutes() {
	timeNow := time.Now()
	self.lock.Lock()
	defer self.lock.Unlock()

	lastMessages := self.routes
	self.routes = make(map[domain.RouteID]domain.Route)
	// Last refresh of time.Unix(0,0) means it lives forever
	for id, route := range lastMessages {
		if (route.ExpiryTime == time.Unix(0, 0)) || timeNow.Before(route.ExpiryTime) {
			self.routes[id] = route
		}
	}
}

func (self *Node) refreshRoute(routeId domain.RouteID) {
	route, routeFound := self.safelyGetRoute(routeId)
	if !routeFound {
		fmt.Fprintf(os.Stderr, "Cannot refresh unknown route: %v\n", routeId)
		return
	}
	base := domain.MessageBase{
		SendRouteID: route.ForwardToRouteID,
		SendBack:    false,
		FromPeerID:  self.Config.PubKey,
		Nonce:       GenerateNonce(),
	}
	directPeer := route.ForwardToPeerID
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v\n", directPeer)
		return
	}
	message := domain.RefreshRouteMessage{
		MessageBase:     base,
		DurationHint:    3 * self.Config.RefreshRouteDuration,
		ConfirmRoutedID: routeId,
	}
	serialized := self.serializer.SerializeMessage(message)
	err := transportToPeer.SendMessage(directPeer, serialized, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Serialization error %v\n", err)
		return
	}
}

func (self *Node) refreshRoutes() {
	self.lock.Lock()
	localRoutes := self.localRoutes
	self.lock.Unlock()

	for routeID := range localRoutes {
		self.refreshRoute(routeID)
	}
}

func (self *Node) clipExpiryTime(hint time.Time) time.Time {
	maxTime := time.Now().Add(self.Config.MaximumForwardingDuration)
	if hint.Unix() > maxTime.Unix() {
		return maxTime
	}
	return hint
}

func getMessageBase(msg interface{}) domain.MessageBase {
	messageType := reflect.TypeOf(msg)

	switch messageType {
	case reflect.TypeOf(domain.UserMessage{}):
		return (msg.(domain.UserMessage)).MessageBase

	case reflect.TypeOf(domain.SetRouteMessage{}):
		return (msg.(domain.SetRouteMessage)).MessageBase

	case reflect.TypeOf(domain.RefreshRouteMessage{}):
		return (msg.(domain.RefreshRouteMessage)).MessageBase

	case reflect.TypeOf(domain.DeleteRouteMessage{}):
		return (msg.(domain.DeleteRouteMessage)).MessageBase

	case reflect.TypeOf(domain.SetRouteReply{}):
		return (msg.(domain.SetRouteReply)).MessageBase
	}

	debug.PrintStack()
	panic(fmt.Sprintf("Internal error: getMessageBase incomplete (%v)", messageType))
}

func (self *Node) DebugCountRoutes() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return len(self.routes)
}

type TimeoutError struct {
}

func (self *TimeoutError) Error() string {
	return "Timeout"
}

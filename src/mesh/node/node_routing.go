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
func (s *Node) AddRoute(routeID domain.RouteID, peerID cipher.PubKey) error {
	_, routeExists := s.safelyGetRoute(routeID)
	if routeExists {
		return errors.New(fmt.Sprintf("Route %v already exists\n", routeID))
	}
	transportToPeer := s.safelyGetTransportToPeer(peerID)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", peerID))
	}

	// Add locally to routes for backward termination
	s.lock.Lock()
	defer s.lock.Unlock()
	s.routes[routeID] =
		domain.Route{
			ForwardToPeerID:   peerID,
			ForwardToRouteID:  routeID,
			BackwardToPeerID:  cipher.PubKey{},
			BackwardToRouteID: NilRouteID,
			// Route lifetime: never dies until route is removed
			ExpiryTime: time.Unix(0, 0),
		}

	s.localRoutesByTerminatingPeer[peerID] = routeID
	s.localRoutes[routeID] = domain.LocalRoute{
		LastForwardingPeerID: s.Config.PubKey,
		TerminatingPeerID:    peerID,
		LastHopRouteID:       NilRouteID,
		LastConfirmed:        time.Unix(0, 0),
	}
	return nil
}

func (s *Node) DeleteRoute(routeID domain.RouteID) (err error) {
	route, routeExists := s.safelyGetRoute(routeID)
	if !routeExists {
		return errors.New(fmt.Sprintf("Cannot delete route %v, doesn't exist\n", routeID))
	}

	err = s.sendDeleteRoute(routeID, route)

	s.lock.Lock()
	defer s.lock.Unlock()
	localRoute, localExists := s.localRoutes[routeID]

	delete(s.routes, routeID)
	delete(s.routeExtensionsAwaitingConfirm, routeID)

	if localExists {
		delete(s.localRoutes, routeID)
		delete(s.localRoutesByTerminatingPeer, localRoute.TerminatingPeerID)
	}
	return err
}

// toPeer must be the public key of a peer connected to the current last node in this route
// Blocks until the set route reply is received or the timeout is reached
func (s *Node) ExtendRoute(routeID domain.RouteID, toPeer cipher.PubKey, timeout time.Duration) (err error) {
	message, directPeer, confirm, extendError := s.extendRouteWithoutSending(routeID, toPeer)
	if extendError != nil {
		return extendError
	}
	transportToPeer := s.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v from %v\n", directPeer, s.Config.PubKey))
	}
	serialized := s.serializer.SerializeMessage(message)
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
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.routeExtensionsAwaitingConfirm, routeID)
	return
}

func (s *Node) ExtendRouteSimple(routeID domain.RouteID, toPeer cipher.PubKey) error {
	message, directPeer, _, extendError := s.extendRouteWithoutSending(routeID, toPeer)
	if extendError != nil {
		return extendError
	}
	transportToPeer := s.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v from %v\n", directPeer, s.Config.PubKey))
	}
	serialized := s.serializer.SerializeMessage(message)
	err := transportToPeer.SendMessage(directPeer, serialized, nil)
	if err != nil {
		return err
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.routeExtensionsAwaitingConfirm, routeID)
	return extendError
}

func (s *Node) GetRouteLastConfirmed(routeID domain.RouteID) (time.Time, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	localRoute, found := s.localRoutes[routeID]
	if !found {
		return time.Unix(0, 0), errors.New(fmt.Sprintf("Unknown route id: %v", routeID))
	}
	return localRoute.LastConfirmed, nil
}

func (s *Node) findRouteToPeer(toPeer cipher.PubKey) (directPeerID cipher.PubKey, localID, sendID domain.RouteID, transport transport.ITransport, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	localRouteID, ok := s.localRoutesByTerminatingPeer[toPeer]
	if ok {
		route, ok := s.routes[localRouteID]
		if !ok {
			panic("Local route struct does not exists")
		}
		directPeerID = route.ForwardToPeerID
		localID = localRouteID
		sendID = route.ForwardToRouteID
	} else {
		return cipher.PubKey{}, NilRouteID, NilRouteID, nil, errors.New(fmt.Sprintf("No route to peer: %v", toPeer))
	}
	transport = s.GetTransportToPeer(directPeerID)
	if transport == nil {
		return cipher.PubKey{}, NilRouteID, NilRouteID, nil,
			errors.New(fmt.Sprintf("No route or transport to peer %v\n", toPeer))
	}
	return
}

func (s *Node) extendRouteWithoutSending(routeID domain.RouteID, toPeer cipher.PubKey) (message domain.SetRouteMessage, directPeer cipher.PubKey, waitConfirm chan bool, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, alreadyExtending := s.routeExtensionsAwaitingConfirm[routeID]
	if alreadyExtending {
		return domain.SetRouteMessage{}, cipher.PubKey{}, nil, errors.New("Cannot extend route more than once at the same time")
	}

	newHopId := routeID

	localRoute, localRouteExists := s.localRoutes[routeID]
	if !localRouteExists {
		return domain.SetRouteMessage{}, cipher.PubKey{}, nil, errors.New(fmt.Sprintf("ExtendRoute called on unknown route: %v", routeID))
	}

	route, routeExists := s.routes[routeID]
	if !routeExists {
		panic("Internal consistency error 8JUL2016544")
	}

	directPeer = route.ForwardToPeerID

	sendBase := domain.MessageBase{
		SendRouteID: route.ForwardToRouteID,
		SendBack:    false,
		FromPeerID:  s.Config.PubKey,
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
		DurationHint:               3 * s.Config.RefreshRouteDuration,
	}

	delete(s.localRoutesByTerminatingPeer, localRoute.TerminatingPeerID)
	s.localRoutes[routeID] = domain.LocalRoute{
		LastForwardingPeerID: localRoute.TerminatingPeerID,
		TerminatingPeerID:    toPeer,
		LastHopRouteID:       newHopId,
		LastConfirmed:        localRoute.LastConfirmed,
	}
	s.localRoutesByTerminatingPeer[toPeer] = routeID

	updatedRoute := route
	updatedRoute.ForwardToRouteID = newHopId
	s.routes[routeID] = updatedRoute
	confirmChan := make(chan bool, 1)
	s.routeExtensionsAwaitingConfirm[routeID] = confirmChan

	return newTermMessage, directPeer, confirmChan, nil
}

func (s *Node) sendDeleteRoute(routeID domain.RouteID, route domain.Route) error {
	sendBase := domain.MessageBase{
		SendRouteID: route.ForwardToRouteID,
		SendBack:    false,
		FromPeerID:  s.Config.PubKey,
		Nonce:       GenerateNonce(),
	}
	message := domain.DeleteRouteMessage{
		MessageBase: sendBase,
	}

	directPeer := route.ForwardToPeerID
	transportToPeer := s.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("2No transport to peer %v from %v\n", directPeer, s.Config.PubKey))
	}
	serialized := s.serializer.SerializeMessage(message)
	err := transportToPeer.SendMessage(directPeer, serialized, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Node) expireOldRoutesLoop() {
	s.closeGroup.Add(1)
	defer s.closeGroup.Done()
	for len(s.closing) == 0 {
		select {
		case <-time.After(s.Config.ExpireRoutesInterval):
			{
				s.expireOldRoutes()
			}
		case <-s.closing:
			{
				return
			}
		}
	}
}

func (s *Node) refreshRoutesLoop() {
	s.closeGroup.Add(1)
	defer s.closeGroup.Done()
	for len(s.closing) == 0 {
		select {
		case <-time.After(s.Config.RefreshRouteDuration):
			{
				s.refreshRoutes()
			}
		case <-s.closing:
			{
				return
			}
		}
	}
}

func (s *Node) safelyGetForwarding(msg interface{}) (SendBack bool, route domain.Route, doForward bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	messageBase := getMessageBase(msg)

	routeFound, foundRoute := s.routes[messageBase.SendRouteID]

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

func (s *Node) safelyGetRoute(routeID domain.RouteID) (route domain.Route, ok bool) {
	route, ok = s.routes[routeID]
	return route, ok
}

func (s *Node) expireOldRoutes() {
	timeNow := time.Now()
	s.lock.Lock()
	defer s.lock.Unlock()

	lastMessages := s.routes
	s.routes = make(map[domain.RouteID]domain.Route)
	// Last refresh of time.Unix(0,0) means it lives forever
	for id, route := range lastMessages {
		if (route.ExpiryTime == time.Unix(0, 0)) || timeNow.Before(route.ExpiryTime) {
			s.routes[id] = route
		}
	}
}

func (s *Node) refreshRoute(routeId domain.RouteID) {
	route, routeFound := s.safelyGetRoute(routeId)
	if !routeFound {
		fmt.Fprintf(os.Stderr, "Cannot refresh unknown route: %v\n", routeId)
		return
	}
	base := domain.MessageBase{
		SendRouteID: route.ForwardToRouteID,
		SendBack:    false,
		FromPeerID:  s.Config.PubKey,
		Nonce:       GenerateNonce(),
	}
	directPeer := route.ForwardToPeerID
	transportToPeer := s.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v\n", directPeer)
		return
	}
	message := domain.RefreshRouteMessage{
		MessageBase:     base,
		DurationHint:    3 * s.Config.RefreshRouteDuration,
		ConfirmRoutedID: routeId,
	}
	serialized := s.serializer.SerializeMessage(message)
	err := transportToPeer.SendMessage(directPeer, serialized, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Serialization error %v\n", err)
		return
	}
}

func (s *Node) refreshRoutes() {
	s.lock.Lock()
	localRoutes := s.localRoutes
	s.lock.Unlock()

	for routeID := range localRoutes {
		s.refreshRoute(routeID)
	}
}

func (s *Node) clipExpiryTime(hint time.Time) time.Time {
	maxTime := time.Now().Add(s.Config.MaximumForwardingDuration)
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

func (s *Node) DebugCountRoutes() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.routes)
}

type TimeoutError struct {
}

func (s *TimeoutError) Error() string {
	return "Timeout"
}

func (s *Node) GetAllRoutes() map[domain.RouteID]domain.Route {
	return s.routes
}

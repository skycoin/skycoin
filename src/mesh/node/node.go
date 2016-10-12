package mesh

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/node/connection"
	"github.com/skycoin/skycoin/src/mesh/serialize"
	"github.com/skycoin/skycoin/src/mesh/transport"
	"gopkg.in/op/go-logging.v1"
)

type TimeoutError struct {
}

var NilRouteID domain.RouteID = (domain.RouteID)(uuid.Nil)

type rewritableMessage interface {
	Rewrite(newSendRouteID domain.RouteID) rewritableMessage
}

type Node struct {
	config                     domain.NodeConfig
	outputMessagesReceived     chan domain.MeshMessage
	transportsMessagesReceived chan []byte
	serializer                 *serialize.Serializer
	//myCrypto                   transport.TransportCrypto

	lock       *sync.Mutex
	closeGroup *sync.WaitGroup
	closing    chan bool

	transports                     map[transport.ITransport]bool
	routes                         map[domain.RouteID]domain.Route
	routeExtensionsAwaitingConfirm map[domain.RouteID]chan bool
	localRoutesByTerminatingPeer   map[cipher.PubKey]domain.RouteID
	localRoutes                    map[domain.RouteID]domain.LocalRoute
	messagesBeingAssembled         map[domain.MessageID]*domain.MessageUnderAssembly
}

func (self *TimeoutError) Error() string {
	return "Timeout"
}

var logger = logging.MustGetLogger("node")

func NewNode(config domain.NodeConfig) (*Node, error) {
	ret := &Node{
		config:                     config,
		outputMessagesReceived:     nil,                                                     // received
		transportsMessagesReceived: make(chan []byte, config.TransportMessageChannelLength), // received
		serializer:                 serialize.NewSerializer(),
		lock:                       &sync.Mutex{}, // Lock
		closeGroup:                 &sync.WaitGroup{},
		closing:                    make(chan bool, 10),
		transports:                 make(map[transport.ITransport]bool),
		messagesBeingAssembled:     make(map[domain.MessageID]*domain.MessageUnderAssembly),
		routes:                     make(map[domain.RouteID]domain.Route),
		localRoutesByTerminatingPeer:   make(map[cipher.PubKey]domain.RouteID),
		localRoutes:                    make(map[domain.RouteID]domain.LocalRoute),
		routeExtensionsAwaitingConfirm: make(map[domain.RouteID]chan bool),
		//myCrypto:                   &ChaChaCrypto{config.ChaCha20Key},
	}
	ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, domain.UserMessage{})
	ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, domain.SetRouteMessage{})
	ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{3}, domain.RefreshRouteMessage{})
	ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{4}, domain.DeleteRouteMessage{})
	ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{5}, domain.SetRouteReply{})

	go ret.processIncomingMessagesLoop()
	go ret.expireOldRoutesLoop()
	go ret.expireOldMessagesLoop()
	go ret.refreshRoutesLoop()

	return ret, nil
}

// Returns nil if reassembly didn't happen (incomplete message)
func (self *Node) reassembleUserMessage(msgIn domain.UserMessage) []byte {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, assembledExists := self.messagesBeingAssembled[msgIn.MessageID]
	if !assembledExists {
		beingAssembled := &domain.MessageUnderAssembly{
			Fragments:   make(map[uint64]domain.UserMessage),
			SendRouteID: msgIn.SendRouteID,
			SendBack:    msgIn.SendBack,
			Count:       msgIn.Count,
			Dropped:     false,
			ExpiryTime:  time.Now().Add(self.config.TimeToAssembleMessage),
		}
		self.messagesBeingAssembled[msgIn.MessageID] = beingAssembled
	}

	beingAssembled, _ := self.messagesBeingAssembled[msgIn.MessageID]

	if beingAssembled.Dropped {
		return nil
	}

	if beingAssembled.Count != msgIn.Count {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different total counts!\n", msgIn.MessageID)
		beingAssembled.Dropped = true
		return nil
	}

	if beingAssembled.SendRouteID != msgIn.SendRouteID {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different send ids!\n", msgIn.SendRouteID)
		beingAssembled.Dropped = true
		return nil
	}

	if beingAssembled.SendBack != msgIn.SendBack {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different send directions!\n", msgIn.SendRouteID)
		beingAssembled.Dropped = true
		return nil
	}

	_, messageExists := beingAssembled.Fragments[msgIn.Index]
	if messageExists {
		fmt.Fprintf(os.Stderr, "Fragment %v of message %v is duplicated, dropping message\n", msgIn.Index, msgIn.MessageID)
		return nil
	}

	beingAssembled.Fragments[msgIn.Index] = msgIn
	if (uint64)(len(beingAssembled.Fragments)) == beingAssembled.Count {
		delete(self.messagesBeingAssembled, msgIn.MessageID)
		reassembled := []byte{}
		for i := (uint64)(0); i < beingAssembled.Count; i++ {
			reassembled = append(reassembled, beingAssembled.Fragments[i].Contents...)
		}
		return reassembled
	}

	return nil
}

func generateNonce() [4]byte {
	ret := make([]byte, 4)
	n, err := rand.Read(ret)
	if n != 4 {
		panic("rand.Read() failed")
	}
	if err != nil {
		panic(err)
	}
	ret_b := [4]byte{0, 0, 0, 0}
	copy(ret_b[:], ret[:])
	return ret_b
}

func getMessageBase(msg interface{}) domain.MessageBase {
	msg_type := reflect.TypeOf(msg)

	if msg_type == reflect.TypeOf(domain.UserMessage{}) {
		return (msg.(domain.UserMessage)).MessageBase
	} else if msg_type == reflect.TypeOf(domain.SetRouteMessage{}) {
		return (msg.(domain.SetRouteMessage)).MessageBase
	} else if msg_type == reflect.TypeOf(domain.RefreshRouteMessage{}) {
		return (msg.(domain.RefreshRouteMessage)).MessageBase
	} else if msg_type == reflect.TypeOf(domain.DeleteRouteMessage{}) {
		return (msg.(domain.DeleteRouteMessage)).MessageBase
	} else if msg_type == reflect.TypeOf(domain.SetRouteReply{}) {
		return (msg.(domain.SetRouteReply)).MessageBase
	}
	debug.PrintStack()
	panic(fmt.Sprintf("Internal error: getMessageBase incomplete (%v)", msg_type))
}

func rewriteMessage(msg interface{}, newBase domain.MessageBase) interface{} {
	msg_type := reflect.TypeOf(msg)
	newBase.Nonce = generateNonce()

	if msg_type == reflect.TypeOf(domain.UserMessage{}) {
		ret := (msg.(domain.UserMessage))
		ret.MessageBase = newBase
		return ret
	} else if msg_type == reflect.TypeOf(domain.SetRouteMessage{}) {
		ret := (msg.(domain.SetRouteMessage))
		ret.MessageBase = newBase
		return ret
	} else if msg_type == reflect.TypeOf(domain.RefreshRouteMessage{}) {
		ret := (msg.(domain.RefreshRouteMessage))
		ret.MessageBase = newBase
		return ret
	} else if msg_type == reflect.TypeOf(domain.DeleteRouteMessage{}) {
		ret := (msg.(domain.DeleteRouteMessage))
		ret.MessageBase = newBase
		return ret
	} else if msg_type == reflect.TypeOf(domain.SetRouteReply{}) {
		ret := (msg.(domain.SetRouteReply))
		ret.MessageBase = newBase
		return ret
	}
	panic("Internal error: rewriteMessage incomplete")
}

func (self *Node) safelyGetForwarding(msg interface{}) (SendBack bool, route domain.Route, doForward bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	messageBase := getMessageBase(msg)
	routeFound, foundRoute := self.routes[messageBase.SendRouteID]

	doForward = foundRoute

	if messageBase.SendBack {
		if routeFound.BackwardToPeer == (cipher.PubKey{}) {
			doForward = false
		}
	} else {
		if routeFound.ForwardToPeer == (cipher.PubKey{}) {
			doForward = false
		}
	}

	if doForward {
		return messageBase.SendBack, routeFound, doForward
	} else {
		return false, domain.Route{}, doForward
	}
}

func (self *Node) safelyGetRoute(routeID domain.RouteID) (route domain.Route, exists bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	route, exists = self.routes[routeID]
	return
}

func (self *Node) safelyGetRewriteBase(msg interface{}) (forwardTo cipher.PubKey, base domain.MessageBase, doForward bool) {
	// sendBack
	sendBack, route, foundRoute := self.safelyGetForwarding(msg)
	if !foundRoute {
		return cipher.PubKey{}, domain.MessageBase{}, false
	}
	forwardTo = route.ForwardToPeer
	rewriteTo := route.ForwardRewriteSendRouteID
	if sendBack {
		forwardTo = route.BackwardToPeer
		rewriteTo = route.BackwardRewriteSendRouteID
	}
	if forwardTo == (cipher.PubKey{}) {
		return cipher.PubKey{}, domain.MessageBase{}, false
	}
	newBase :=
		domain.MessageBase{
			SendRouteID: rewriteTo,
			SendBack:    sendBack,
			FromPeer:    self.config.PubKey,
			Nonce:       generateNonce(),
		}
	return forwardTo, newBase, true
}

func (self *Node) forwardMessage(msg interface{}) bool {
	forwardTo, newBase, doForward := self.safelyGetRewriteBase(msg)
	if !doForward {
		return false
	}
	// Rewrite
	rewritten := rewriteMessage(msg, newBase)
	transportToPeer := self.safelyGetTransportToPeer(forwardTo)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport found for forwarded message from %v to %v, dropping\n", self.config.PubKey, forwardTo)
		return true
	}

	serialized := self.serializer.SerializeMessage(rewritten)
	err := transportToPeer.SendMessage(forwardTo, serialized)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to send forwarded message, dropping\n")
		return true
	}

	// Forward, not receive
	return true
}

func (self *Node) processUserMessage(msgIn domain.UserMessage) {
	reassembled := self.reassembleUserMessage(msgIn)
	// Not finished reassembling yet
	if reassembled == nil {
		return
	}
	directPeer, forwardBase, doForward := self.safelyGetRewriteBase(msgIn)
	if doForward {
		transportToPeer := self.safelyGetTransportToPeer(directPeer)
		if transportToPeer == nil {
			fmt.Fprintf(os.Stderr, "No transport to peer %v from %v, dropping\n", directPeer, self.config.PubKey)
			return
		}
		// Forward reassembled message, not individual pieces. This is done because of the need for refragmentation
		fragments := connection.ConnectionManager.FragmentMessage(reassembled, directPeer, transportToPeer, forwardBase)
		for _, fragment := range fragments {
			serialized := self.serializer.SerializeMessage(fragment)
			err := transportToPeer.SendMessage(directPeer, serialized)
			if err != nil {
				fmt.Fprint(os.Stderr, "Failed to send forwarded message, dropping\n")
				return
			}
		}
	} else {
		self.outputMessagesReceived <- domain.MeshMessage{domain.ReplyTo{msgIn.SendRouteID, msgIn.FromPeer}, reassembled}
	}
}

func (self *Node) sendSetRouteReply(base domain.MessageBase, confirmId domain.RouteID) {
	reply := domain.SetRouteReply{
		MessageBase: domain.MessageBase{
			SendRouteID: base.SendRouteID,
			SendBack:    true,
			FromPeer:    self.config.PubKey,
			Nonce:       generateNonce(),
		},
		ConfirmRouteID: confirmId,
	}
	transportToPeer := self.safelyGetTransportToPeer(base.FromPeer)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v from %v\n", base.FromPeer, self.config.PubKey)
		return
	}
	serialized := self.serializer.SerializeMessage(reply)
	err := transportToPeer.SendMessage(base.FromPeer, serialized)
	if err != nil {
		return
	}
}

func (self *Node) processSetRouteMessage(msg domain.SetRouteMessage) {
	if msg.SendBack {
		fmt.Fprintf(os.Stderr, "Invalid SetRouteMessage received, dropping: %v\n", msg)
		return
	}
	self.lock.Lock()
	defer self.lock.Unlock()

	if msg.SetRouteID == NilRouteID {
		fmt.Fprintf(os.Stderr, "Invalid SetRouteMessage received, dropping: %v\n", msg)
		return
	}
	self.routes[msg.SetRouteID] =
		domain.Route{
			ForwardToPeer:              msg.ForwardToPeer,
			ForwardRewriteSendRouteID:  msg.ForwardRewriteSendRouteID,
			BackwardToPeer:             msg.BackwardToPeer,
			BackwardRewriteSendRouteID: msg.BackwardRewriteSendRouteID,
			ExpiryTime:                 self.clipExpiryTime(time.Now().Add(msg.DurationHint)),
		}

	// Don't block to send reply
	go self.sendSetRouteReply(msg.MessageBase, msg.ConfirmRouteID)
}

func (self *Node) processSetRouteReplyMessage(msg domain.SetRouteReply) {
	self.lock.Lock()
	defer self.lock.Unlock()
	confirmChan, foundChan := self.routeExtensionsAwaitingConfirm[msg.ConfirmRouteID]
	if foundChan {
		confirmChan <- true
	}
	localRoute, foundLocal := self.localRoutes[msg.ConfirmRouteID]
	if foundLocal {
		localRoute.LastConfirmed = time.Now()
		self.localRoutes[msg.ConfirmRouteID] = localRoute
	}
}

func (self *Node) clipExpiryTime(hint time.Time) time.Time {
	maxTime := time.Now().Add(self.config.MaximumForwardingDuration)
	if hint.Unix() > maxTime.Unix() {
		return maxTime
	}
	return hint
}

func (self *Node) processRefreshRouteMessage(msg domain.RefreshRouteMessage, forwarded bool) {
	if forwarded {
		self.lock.Lock()
		defer self.lock.Unlock()
		route, exists := self.routes[msg.SendRouteID]
		if !exists {
			fmt.Fprintf(os.Stderr, "Refresh sent for unknown route: %v\n", msg.SendRouteID)
			return
		}
		route.ExpiryTime = self.clipExpiryTime(time.Now().Add(msg.DurationHint))
		self.routes[msg.SendRouteID] = route
	} else {
		// Don't block to send reply
		go self.sendSetRouteReply(msg.MessageBase, msg.ConfirmRoutedID)
	}
}

func (self *Node) processDeleteRouteMessage(msg domain.DeleteRouteMessage) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.routes, msg.SendRouteID)
}

func (self *Node) processMessage(serialized []byte) {
	msg, deserialize_error := self.serializer.UnserializeMessage(serialized)
	if deserialize_error != nil {
		fmt.Fprintf(os.Stderr, "Deserialization error %v\n", deserialize_error)
		return
	}

	msg_type := reflect.TypeOf(msg)
	// User messages have fragmentation to deal with
	if msg_type == reflect.TypeOf(domain.UserMessage{}) {
		self.processUserMessage(msg.(domain.UserMessage))
	} else {
		forwardedMessage := self.forwardMessage(msg)
		if !forwardedMessage {
			// Receive or forward. Refragment on forward!
			if msg_type == reflect.TypeOf(domain.SetRouteMessage{}) {
				self.processSetRouteMessage(msg.(domain.SetRouteMessage))
			} else if msg_type == reflect.TypeOf(domain.SetRouteReply{}) {
				self.processSetRouteReplyMessage(msg.(domain.SetRouteReply))
			}
		} else {
			if msg_type == reflect.TypeOf(domain.DeleteRouteMessage{}) {
				self.processDeleteRouteMessage(msg.(domain.DeleteRouteMessage))
			}
		}

		if msg_type == reflect.TypeOf(domain.RefreshRouteMessage{}) {
			self.processRefreshRouteMessage(msg.(domain.RefreshRouteMessage), forwardedMessage)
		}
	}
}

func (self *Node) expireOldMessages() {
	time_now := time.Now()
	self.lock.Lock()
	defer self.lock.Unlock()

	lastMessages := self.messagesBeingAssembled
	self.messagesBeingAssembled = make(map[domain.MessageID]*domain.MessageUnderAssembly)
	for id, msg := range lastMessages {
		if time_now.Before(msg.ExpiryTime) {
			self.messagesBeingAssembled[id] = msg
		}
	}
}

func (self *Node) expireOldRoutes() {
	time_now := time.Now()
	self.lock.Lock()
	defer self.lock.Unlock()

	lastMessages := self.routes
	self.routes = make(map[domain.RouteID]domain.Route)
	// Last refresh of time.Unix(0,0) means it lives forever
	for id, route := range lastMessages {
		if (route.ExpiryTime == time.Unix(0, 0)) || time_now.Before(route.ExpiryTime) {
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
		SendRouteID: route.ForwardRewriteSendRouteID,
		SendBack:    false,
		FromPeer:    self.config.PubKey,
		Nonce:       generateNonce(),
	}
	directPeer := route.ForwardToPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v\n", directPeer)
		return
	}
	message := domain.RefreshRouteMessage{
		MessageBase:     base,
		DurationHint:    3 * self.config.RefreshRouteDuration,
		ConfirmRoutedID: routeId,
	}
	serialized := self.serializer.SerializeMessage(message)
	err := transportToPeer.SendMessage(directPeer, serialized)
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

func (self *Node) expireOldMessagesLoop() {
	self.closeGroup.Add(1)
	defer self.closeGroup.Done()
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.config.ExpireMessagesInterval):
			{
				self.expireOldMessages()
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Node) processIncomingMessagesLoop() {
	self.closeGroup.Add(1)
	defer self.closeGroup.Done()
	for len(self.closing) == 0 {
		select {
		case msg, ok := <-self.transportsMessagesReceived:
			{
				if ok {
					self.processMessage(msg)
				}
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Node) expireOldRoutesLoop() {
	self.closeGroup.Add(1)
	defer self.closeGroup.Done()
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.config.ExpireRoutesInterval):
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
		case <-time.After(self.config.RefreshRouteDuration):
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

// Waits for transports to close
func (self *Node) Close() error {
	for i := 0; i < 10; i++ {
		self.closing <- true
	}
	close(self.transportsMessagesReceived)
	self.closeGroup.Wait()
	return nil
}

func (self *Node) GetConfig() domain.NodeConfig {
	return self.config
}

// Node takes ownership of the transport, and will call Close() when it is closed
func (self *Node) AddTransport(transportNode transport.ITransport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	//chaCha20Key := &transport.ChaChaCrypto{}
	//chaCha20Key.SetKey(chaChaKey)
	//transportNode.SetCrypto(chaCha20Key)
	transportNode.SetReceiveChannel(self.transportsMessagesReceived)
	self.transports[transportNode] = true
}

func (self *Node) RemoveTransport(transport transport.ITransport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.transports, transport)
}

func (self *Node) GetTransports() []transport.ITransport {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []transport.ITransport{}
	for nodeTransport := range self.transports {
		ret = append(ret, nodeTransport)
	}
	return ret
}

func (self *Node) GetConnectedPeers() []cipher.PubKey {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []cipher.PubKey{}
	for nodeTransport := range self.transports {
		peers := nodeTransport.GetConnectedPeers()
		ret = append(ret, peers...)
	}
	return ret
}

func (self *Node) ConnectedToPeer(peer cipher.PubKey) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	for nodeTransport := range self.transports {
		if nodeTransport.ConnectedToPeer(peer) {
			return true
		}
	}
	return false
}

// Message order is not preserved
func (self *Node) SetReceiveChannel(received chan domain.MeshMessage) {
	self.outputMessagesReceived = received
}

// toPeer must be the public key of a connected peer
func (self *Node) AddRoute(routeID domain.RouteID, toPeer cipher.PubKey) error {
	_, routeExists := self.safelyGetRoute(routeID)
	if routeExists {
		return errors.New(fmt.Sprintf("Route %v already exists\n", routeID))
	}

	transportToPeer := self.safelyGetTransportToPeer(toPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", toPeer))
	}

	// Add locally to routesById for backward termination
	self.lock.Lock()
	defer self.lock.Unlock()
	self.routes[routeID] =
		domain.Route{
			ForwardToPeer:              toPeer,
			ForwardRewriteSendRouteID:  routeID,
			BackwardToPeer:             cipher.PubKey{},
			BackwardRewriteSendRouteID: NilRouteID,
			// Route lifetime: never dies until route is removed
			ExpiryTime: time.Unix(0, 0),
		}

	self.localRoutesByTerminatingPeer[toPeer] = routeID
	self.localRoutes[routeID] = domain.LocalRoute{
		LastForwardingPeer: self.config.PubKey,
		TerminatingPeer:    toPeer,
		LastHopRouteID:     NilRouteID,
		LastConfirmed:      time.Unix(0, 0),
	}
	return nil
}

func (self *Node) sendDeleteRoute(routeID domain.RouteID, route domain.Route) error {
	sendBase := domain.MessageBase{
		SendRouteID: route.ForwardRewriteSendRouteID,
		SendBack:    false,
		FromPeer:    self.config.PubKey,
		Nonce:       generateNonce(),
	}
	message := domain.DeleteRouteMessage{
		MessageBase: sendBase,
	}

	directPeer := route.ForwardToPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v from %v\n", directPeer, self.config.PubKey))
	}
	serialized := self.serializer.SerializeMessage(message)
	err := transportToPeer.SendMessage(directPeer, serialized)
	if err != nil {
		return err
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
		delete(self.localRoutesByTerminatingPeer, localRoute.TerminatingPeer)
	}
	return err
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

	directPeer = route.ForwardToPeer

	sendBase := domain.MessageBase{
		SendRouteID: route.ForwardRewriteSendRouteID,
		SendBack:    false,
		FromPeer:    self.config.PubKey,
		Nonce:       generateNonce(),
	}

	newTermMessage := domain.SetRouteMessage{
		MessageBase:                sendBase,
		SetRouteID:                 routeID,
		ConfirmRouteID:             routeID,
		ForwardToPeer:              toPeer,
		ForwardRewriteSendRouteID:  routeID,
		BackwardToPeer:             localRoute.LastForwardingPeer,
		BackwardRewriteSendRouteID: localRoute.LastHopRouteID,
		DurationHint:               3 * self.config.RefreshRouteDuration,
	}
	delete(self.localRoutesByTerminatingPeer, localRoute.TerminatingPeer)
	self.localRoutes[routeID] = domain.LocalRoute{
		LastForwardingPeer: localRoute.TerminatingPeer,
		TerminatingPeer:    toPeer,
		LastHopRouteID:     newHopId,
		LastConfirmed:      localRoute.LastConfirmed,
	}
	self.localRoutesByTerminatingPeer[toPeer] = routeID

	updatedRoute := route
	updatedRoute.ForwardRewriteSendRouteID = newHopId
	self.routes[routeID] = updatedRoute
	confirmChan := make(chan bool, 1)
	self.routeExtensionsAwaitingConfirm[routeID] = confirmChan

	return newTermMessage, directPeer, confirmChan, nil
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
		return errors.New(fmt.Sprintf("No transport to peer %v from %v\n", directPeer, self.config.PubKey))
	}
	serialized := self.serializer.SerializeMessage(message)
	err = transportToPeer.SendMessage(directPeer, serialized)
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

func (self *Node) unsafelyGetTransportToPeer(peerKey cipher.PubKey) transport.ITransport {
	for transportToPeer := range self.transports {
		// TODO: Choose transport more intelligently
		if transportToPeer.ConnectedToPeer(peerKey) {
			return transportToPeer
		}
	}
	return nil
}

func (self *Node) safelyGetTransportToPeer(peerKey cipher.PubKey) transport.ITransport {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.unsafelyGetTransportToPeer(peerKey)
}

func (self *Node) findRouteToPeer(toPeer cipher.PubKey) (directPeer cipher.PubKey, localId, sendId domain.RouteID, transport transport.ITransport, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	localRouteId, localRouteExists := self.localRoutesByTerminatingPeer[toPeer]
	if localRouteExists {
		route, routeStructExists := self.routes[localRouteId]
		if !routeStructExists {
			panic("Local route struct does not exists")
		}
		directPeer = route.ForwardToPeer
		localId = localRouteId
		sendId = route.ForwardRewriteSendRouteID
	} else {
		return cipher.PubKey{}, NilRouteID, NilRouteID, nil, errors.New(fmt.Sprintf("No route to peer: %v", toPeer))
	}
	transport = self.unsafelyGetTransportToPeer(directPeer)
	if transport == nil {
		return cipher.PubKey{}, NilRouteID, NilRouteID, nil,
			errors.New(fmt.Sprintf("No route or transport to peer %v\n", toPeer))
	}
	return
}

// Chooses a route automatically. Sends directly without a route if connected to that peer
func (self *Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte) (err error, routeId domain.RouteID) {
	directPeer, localRouteId, sendId, transportToPeer, err := self.findRouteToPeer(toPeer)
	if err != nil {
		return err, NilRouteID
	}
	base := domain.MessageBase{
		SendRouteID: sendId,
		SendBack:    false,
		FromPeer:    self.config.PubKey,
		Nonce:       generateNonce(),
	}
	messages := connection.ConnectionManager.FragmentMessage(contents, directPeer, transportToPeer, base)
	for _, message := range messages {
		serialized := self.serializer.SerializeMessage(message)
		err := transportToPeer.SendMessage(directPeer, serialized)
		if err != nil {
			return err, NilRouteID
		}
	}
	return nil, localRouteId
}

// Blocks until message is confirmed received
func (self *Node) SendMessageThruRoute(routeId domain.RouteID, contents []byte) error {
	route, routeFound := self.safelyGetRoute(routeId)
	if !routeFound {
		return errors.New("Route not found")
	}

	base := domain.MessageBase{
		SendRouteID: route.ForwardRewriteSendRouteID,
		SendBack:    false,
		FromPeer:    self.config.PubKey,
		Nonce:       generateNonce(),
	}
	directPeer := route.ForwardToPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", directPeer))
	}
	messages := connection.ConnectionManager.FragmentMessage(contents, directPeer, transportToPeer, base)
	for _, message := range messages {
		serialized := self.serializer.SerializeMessage(message)
		fmt.Fprintln(os.Stdout, "Send Message")
		err := transportToPeer.SendMessage(directPeer, serialized)
		if err != nil {
			return err
		}
	}
	return nil
}

// Blocks until message is confirmed received
func (self *Node) SendMessageBackThruRoute(replyTo domain.ReplyTo, contents []byte) error {
	directPeer := replyTo.FromPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No route or transport to peer %v\n", directPeer))
	}
	base := domain.MessageBase{
		SendRouteID: replyTo.RouteID,
		SendBack:    true,
		FromPeer:    self.config.PubKey,
		Nonce:       generateNonce(),
	}
	messages := connection.ConnectionManager.FragmentMessage(contents, directPeer, transportToPeer, base)
	for _, message := range messages {
		serialized := self.serializer.SerializeMessage(message)
		err := transportToPeer.SendMessage(directPeer, serialized)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Node) debug_countRoutes() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return len(self.routes)
}

func (self *Node) debug_countMessages() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return len(self.messagesBeingAssembled)
}

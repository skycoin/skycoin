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
	"github.com/skycoin/skycoin/src/mesh/transport/protocol"
	"github.com/skycoin/skycoin/src/mesh/transport/transport"
	"gopkg.in/op/go-logging.v1"
)

type TimeoutError struct {
}

type ReplyTo struct {
	routeId  domain.RouteId
	fromPeer cipher.PubKey
}

type MeshMessage struct {
	ReplyTo  ReplyTo
	Contents []byte
}

var NilRouteId domain.RouteId = (domain.RouteId)(uuid.Nil)

type rewriteableMessage interface {
	Rewrite(newSendId domain.RouteId) rewriteableMessage
}

type Node struct {
	config                     domain.NodeConfig
	outputMessagesReceived     chan MeshMessage
	transportsMessagesReceived chan []byte
	serializer                 *serialize.Serializer
	//myCrypto                   transport.TransportCrypto

	lock       *sync.Mutex
	closeGroup *sync.WaitGroup
	closing    chan bool

	transports                     map[transport.Transport]bool
	messagesBeingAssembled         map[domain.MessageId]*domain.MessageUnderAssembly
	routesById                     map[domain.RouteId]domain.Route
	localRoutesByTerminatingPeer   map[cipher.PubKey]domain.RouteId
	localRoutesById                map[domain.RouteId]domain.LocalRoute
	routeExtensionsAwaitingConfirm map[domain.RouteId]chan bool
}

func (self *TimeoutError) Error() string {
	return "Timeout"
}

var logger = logging.MustGetLogger("node")

func NewNode(config domain.NodeConfig) (*Node, error) {
	ret := &Node{
		config:                         config,
		outputMessagesReceived:         nil,                                                     // received
		transportsMessagesReceived:     make(chan []byte, config.TransportMessageChannelLength), // received
		serializer:                     serialize.NewSerializer(),
		lock:                           &sync.Mutex{}, // Lock
		closeGroup:                     &sync.WaitGroup{},
		closing:                        make(chan bool, 10),
		transports:                     make(map[transport.Transport]bool),
		messagesBeingAssembled:         make(map[domain.MessageId]*domain.MessageUnderAssembly),
		routesById:                     make(map[domain.RouteId]domain.Route),
		localRoutesByTerminatingPeer:   make(map[cipher.PubKey]domain.RouteId),
		localRoutesById:                make(map[domain.RouteId]domain.LocalRoute),
		routeExtensionsAwaitingConfirm: make(map[domain.RouteId]chan bool),
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
func (self *Node) reassembleUserMessage(msgIn domain.UserMessage) (contents []byte) {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, assembledExists := self.messagesBeingAssembled[msgIn.MessageId]
	if !assembledExists {
		beingAssembled := &domain.MessageUnderAssembly{
			Fragments:  make(map[uint64]domain.UserMessage),
			SendId:     msgIn.SendId,
			SendBack:   msgIn.SendBack,
			Count:      msgIn.Count,
			Dropped:    false,
			ExpiryTime: time.Now().Add(self.config.TimeToAssembleMessage),
		}
		self.messagesBeingAssembled[msgIn.MessageId] = beingAssembled
	}

	beingAssembled, _ := self.messagesBeingAssembled[msgIn.MessageId]

	if beingAssembled.Dropped {
		return nil
	}

	if beingAssembled.Count != msgIn.Count {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different total counts!\n", msgIn.MessageId)
		beingAssembled.Dropped = true
		return nil
	}

	if beingAssembled.SendId != msgIn.SendId {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different send ids!\n", msgIn.SendId)
		beingAssembled.Dropped = true
		return nil
	}

	if beingAssembled.SendBack != msgIn.SendBack {
		fmt.Fprintf(os.Stderr, "Fragments of message %v have different send directions!\n", msgIn.SendId)
		beingAssembled.Dropped = true
		return nil
	}

	_, messageExists := beingAssembled.Fragments[msgIn.Index]
	if messageExists {
		fmt.Fprintf(os.Stderr, "Fragment %v of message %v is duplicated, dropping message\n", msgIn.Index, msgIn.MessageId)
		return nil
	}

	beingAssembled.Fragments[msgIn.Index] = msgIn
	if (uint64)(len(beingAssembled.Fragments)) == beingAssembled.Count {
		delete(self.messagesBeingAssembled, msgIn.MessageId)
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

func getMessageBase(msg interface{}) (base domain.MessageBase) {
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

func rewriteMessage(msg interface{}, newBase domain.MessageBase) (rewritten interface{}) {
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

func (self *Node) safelyGetForwarding(msg interface{}) (sendBack bool, route domain.Route, doForward bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	messageBase := getMessageBase(msg)
	routeFound, foundRoute := self.routesById[messageBase.SendId]

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

func (self *Node) safelyGetRoute(id domain.RouteId) (route domain.Route, exists bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	route, exists = self.routesById[id]
	return
}

func (self *Node) safelyGetRewriteBase(msg interface{}) (forwardTo cipher.PubKey, base domain.MessageBase, doForward bool) {
	// sendBack
	sendBack, route, foundRoute := self.safelyGetForwarding(msg)
	if !foundRoute {
		return cipher.PubKey{}, domain.MessageBase{}, false
	}
	forwardTo = route.ForwardToPeer
	rewriteTo := route.ForwardRewriteSendId
	if sendBack {
		forwardTo = route.BackwardToPeer
		rewriteTo = route.BackwardRewriteSendId
	}
	if forwardTo == (cipher.PubKey{}) {
		return cipher.PubKey{}, domain.MessageBase{}, false
	}
	newBase :=
		domain.MessageBase{
			SendId:   rewriteTo,
			SendBack: sendBack,
			FromPeer: self.config.PubKey,
			Reliably: base.Reliably,
			Nonce:    generateNonce(),
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
	transportToPeer := self.safelyGetTransportToPeer(forwardTo, newBase.Reliably)
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
		transportToPeer := self.safelyGetTransportToPeer(directPeer, msgIn.Reliably)
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
		self.outputMessagesReceived <- MeshMessage{ReplyTo{msgIn.SendId, msgIn.FromPeer}, reassembled}
	}
}

func (self *Node) sendSetRouteReply(base domain.MessageBase, confirmId domain.RouteId) {
	reply := domain.SetRouteReply{
		MessageBase: domain.MessageBase{
			SendId:   base.SendId,
			SendBack: true,
			FromPeer: self.config.PubKey,
			Reliably: true,
			Nonce:    generateNonce(),
		},
		ConfirmId: confirmId,
	}
	transportToPeer := self.safelyGetTransportToPeer(base.FromPeer, true)
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

	if msg.SetRouteId == NilRouteId {
		fmt.Fprintf(os.Stderr, "Invalid SetRouteMessage received, dropping: %v\n", msg)
		return
	}
	self.routesById[msg.SetRouteId] =
		domain.Route{
			ForwardToPeer:         msg.ForwardToPeer,
			ForwardRewriteSendId:  msg.ForwardRewriteSendId,
			BackwardToPeer:        msg.BackwardToPeer,
			BackwardRewriteSendId: msg.BackwardRewriteSendId,
			ExpiryTime:            self.clipExpiryTime(time.Now().Add(msg.DurationHint)),
		}

	// Don't block to send reply
	go self.sendSetRouteReply(msg.MessageBase, msg.ConfirmId)
}

func (self *Node) processSetRouteReplyMessage(msg domain.SetRouteReply) {
	self.lock.Lock()
	defer self.lock.Unlock()
	confirmChan, foundChan := self.routeExtensionsAwaitingConfirm[msg.ConfirmId]
	if foundChan {
		confirmChan <- true
	}
	localRoute, foundLocal := self.localRoutesById[msg.ConfirmId]
	if foundLocal {
		localRoute.LastConfirmed = time.Now()
		self.localRoutesById[msg.ConfirmId] = localRoute
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
		route, exists := self.routesById[msg.SendId]
		if !exists {
			fmt.Fprintf(os.Stderr, "Refresh sent for unknown route: %v\n", msg.SendId)
			return
		}
		route.ExpiryTime = self.clipExpiryTime(time.Now().Add(msg.DurationHint))
		self.routesById[msg.SendId] = route
	} else {
		// Don't block to send reply
		go self.sendSetRouteReply(msg.MessageBase, msg.ConfirmId)
	}
}

func (self *Node) processDeleteRouteMessage(msg domain.DeleteRouteMessage) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.routesById, msg.SendId)
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
	self.messagesBeingAssembled = make(map[domain.MessageId]*domain.MessageUnderAssembly)
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

	lastMessages := self.routesById
	self.routesById = make(map[domain.RouteId]domain.Route)
	// Last refresh of time.Unix(0,0) means it lives forever
	for id, route := range lastMessages {
		if (route.ExpiryTime == time.Unix(0, 0)) || time_now.Before(route.ExpiryTime) {
			self.routesById[id] = route
		}
	}
}

func (self *Node) refreshRoute(routeId domain.RouteId) {
	route, routeFound := self.safelyGetRoute(routeId)
	if !routeFound {
		fmt.Fprintf(os.Stderr, "Cannot refresh unknown route: %v\n", routeId)
		return
	}
	reliably := true
	base := domain.MessageBase{
		SendId:   route.ForwardRewriteSendId,
		SendBack: false,
		FromPeer: self.config.PubKey,
		Reliably: reliably,
		Nonce:    generateNonce(),
	}
	directPeer := route.ForwardToPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer, reliably)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v\n", directPeer)
		return
	}
	message := domain.RefreshRouteMessage{
		MessageBase:  base,
		DurationHint: 3 * self.config.RefreshRouteDuration,
		ConfirmId:    routeId,
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
	localRoutes := self.localRoutesById
	self.lock.Unlock()

	for id := range localRoutes {
		self.refreshRoute(id)
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
func (self *Node) AddTransport(transportNode transport.Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	//chaCha20Key := &transport.ChaChaCrypto{}
	//chaCha20Key.SetKey(chaChaKey)
	//transportNode.SetCrypto(chaCha20Key)
	transportNode.SetReceiveChannel(self.transportsMessagesReceived)
	self.transports[transportNode] = true
}

func (self *Node) RemoveTransport(transport transport.Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.transports, transport)
}

func (self *Node) GetTransports() []transport.Transport {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []transport.Transport{}
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
func (self *Node) SetReceiveChannel(received chan MeshMessage) {
	self.outputMessagesReceived = received
}

// toPeer must be the public key of a connected peer
// Requires reliable transport (for now)
func (self *Node) AddRoute(id domain.RouteId, toPeer cipher.PubKey) error {
	_, routeExists := self.safelyGetRoute(id)
	if routeExists {
		return errors.New(fmt.Sprintf("Route %v already exists\n", id))
	}

	transportToPeer := self.safelyGetTransportToPeer(toPeer, true)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", toPeer))
	}

	// Add locally to routesById for backward termination
	self.lock.Lock()
	defer self.lock.Unlock()
	self.routesById[id] =
		domain.Route{
			ForwardToPeer:         toPeer,
			ForwardRewriteSendId:  id,
			BackwardToPeer:        cipher.PubKey{},
			BackwardRewriteSendId: NilRouteId,
			// Route lifetime: never dies until route is removed
			ExpiryTime: time.Unix(0, 0),
		}

	self.localRoutesByTerminatingPeer[toPeer] = id
	self.localRoutesById[id] = domain.LocalRoute{
		LastForwardingPeer: self.config.PubKey,
		TerminatingPeer:    toPeer,
		LastHopId:          NilRouteId,
		LastConfirmed:      time.Unix(0, 0),
	}
	return nil
}

func (self *Node) sendDeleteRoute(id domain.RouteId, route domain.Route) error {
	sendBase := domain.MessageBase{
		SendId:   route.ForwardRewriteSendId,
		SendBack: false,
		FromPeer: self.config.PubKey,
		Reliably: true, // Reliable
		Nonce:    generateNonce(),
	}
	message := domain.DeleteRouteMessage{
		sendBase,
	}

	directPeer := route.ForwardToPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer, true)
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

func (self *Node) DeleteRoute(id domain.RouteId) (err error) {
	route, routeExists := self.safelyGetRoute(id)
	if !routeExists {
		return errors.New(fmt.Sprintf("Cannot delete route %v, doesn't exist\n", id))
	}

	err = self.sendDeleteRoute(id, route)

	self.lock.Lock()
	defer self.lock.Unlock()
	localRoute, localExists := self.localRoutesById[id]

	delete(self.routesById, id)
	delete(self.routeExtensionsAwaitingConfirm, id)

	if localExists {
		delete(self.localRoutesById, id)
		delete(self.localRoutesByTerminatingPeer, localRoute.TerminatingPeer)
	}
	return err
}

func (self *Node) extendRouteWithoutSending(id domain.RouteId, toPeer cipher.PubKey) (message domain.SetRouteMessage, directPeer cipher.PubKey, waitConfirm chan bool, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, alreadyExtending := self.routeExtensionsAwaitingConfirm[id]
	if alreadyExtending {
		return domain.SetRouteMessage{}, cipher.PubKey{}, nil, errors.New("Cannot extend route more than once at the same time")
	}

	newHopId := id

	localRoute, localRouteExists := self.localRoutesById[id]
	if !localRouteExists {
		return domain.SetRouteMessage{}, cipher.PubKey{}, nil, errors.New(fmt.Sprintf("ExtendRoute called on unknown route: %v", id))
	}

	route, routeExists := self.routesById[id]
	if !routeExists {
		panic("Internal consistency error 8JUL2016544")
	}

	directPeer = route.ForwardToPeer

	sendBase := domain.MessageBase{
		SendId:   route.ForwardRewriteSendId,
		SendBack: false,
		FromPeer: self.config.PubKey,
		Reliably: true,
		Nonce:    generateNonce(),
	}

	newTermMessage := domain.SetRouteMessage{
		MessageBase:           sendBase,
		SetRouteId:            id,
		ConfirmId:             id,
		ForwardToPeer:         toPeer,
		ForwardRewriteSendId:  id,
		BackwardToPeer:        localRoute.LastForwardingPeer,
		BackwardRewriteSendId: localRoute.LastHopId,
		DurationHint:          3 * self.config.RefreshRouteDuration,
	}
	delete(self.localRoutesByTerminatingPeer, localRoute.TerminatingPeer)
	self.localRoutesById[id] = domain.LocalRoute{
		LastForwardingPeer: localRoute.TerminatingPeer,
		TerminatingPeer:    toPeer,
		LastHopId:          newHopId,
		LastConfirmed:      localRoute.LastConfirmed,
	}
	self.localRoutesByTerminatingPeer[toPeer] = id

	updatedRoute := route
	updatedRoute.ForwardRewriteSendId = newHopId
	self.routesById[id] = updatedRoute
	confirmChan := make(chan bool, 1)
	self.routeExtensionsAwaitingConfirm[id] = confirmChan

	return newTermMessage, directPeer, confirmChan, nil
}

// toPeer must be the public key of a peer connected to the current last node in this route
// Blocks until the set route reply is received or the timeout is reached
func (self *Node) ExtendRoute(id domain.RouteId, toPeer cipher.PubKey, timeout time.Duration) (err error) {
	message, directPeer, confirm, extendError := self.extendRouteWithoutSending(id, toPeer)
	if extendError != nil {
		return extendError
	}
	transportToPeer := self.safelyGetTransportToPeer(directPeer, true)
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

	delete(self.routeExtensionsAwaitingConfirm, id)
	return
}

func (self *Node) GetRouteLastConfirmed(id domain.RouteId) (time.Time, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	localRoute, found := self.localRoutesById[id]
	if !found {
		return time.Unix(0, 0), errors.New(fmt.Sprintf("Unknown route id: %v", id))
	}
	return localRoute.LastConfirmed, nil
}

func (self *Node) unsafelyGetTransportToPeer(peerKey cipher.PubKey, reliably bool) transport.Transport {
	// If unreliable, prefer unreliable transport
	if !reliably {
		for transportToPeer := range self.transports {
			// TODO: Choose transport more intelligently
			if transportToPeer.ConnectedToPeer(peerKey) && !transportToPeer.IsReliable() {
				return transportToPeer
			}
		}
	}
	for transportToPeer := range self.transports {
		// TODO: Choose transport more intelligently
		if transportToPeer.ConnectedToPeer(peerKey) && ((!reliably) || transportToPeer.IsReliable()) {
			return transportToPeer
		}
	}
	return nil
}

func (self *Node) safelyGetTransportToPeer(peerKey cipher.PubKey, reliably bool) transport.Transport {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.unsafelyGetTransportToPeer(peerKey, reliably)
}

func (self *Node) findRouteToPeer(toPeer cipher.PubKey, reliably bool) (directPeer cipher.PubKey, localId, sendId domain.RouteId, transport transport.Transport, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	localRouteId, localRouteExists := self.localRoutesByTerminatingPeer[toPeer]
	if localRouteExists {
		route, routeStructExists := self.routesById[localRouteId]
		if !routeStructExists {
			panic("Local route struct does not exists")
		}
		directPeer = route.ForwardToPeer
		localId = localRouteId
		sendId = route.ForwardRewriteSendId
	} else {
		return cipher.PubKey{}, NilRouteId, NilRouteId, nil, errors.New(fmt.Sprintf("No route to peer: %v", toPeer))
	}
	transport = self.unsafelyGetTransportToPeer(directPeer, reliably)
	if transport == nil {
		return cipher.PubKey{}, NilRouteId, NilRouteId, nil,
			errors.New(fmt.Sprintf("No route or transport to peer %v\n", toPeer))
	}
	return
}

// Chooses a route automatically. Sends directly without a route if connected to that peer
// if reliably is false, unreliable transport is preferred, but reliable is chosen if it's the only option
// if reliably is true, reliable transport only is used
func (self *Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte, reliably bool) (err error, routeId domain.RouteId) {
	directPeer, localRouteId, sendId, transportToPeer, err := self.findRouteToPeer(toPeer, reliably)
	if err != nil {
		return err, NilRouteId
	}
	base := domain.MessageBase{
		SendId:   sendId,
		SendBack: false,
		FromPeer: self.config.PubKey,
		Reliably: reliably,
		Nonce:    generateNonce(),
	}
	messages := connection.ConnectionManager.FragmentMessage(contents, directPeer, transportToPeer, base)
	for _, message := range messages {
		serialized := self.serializer.SerializeMessage(message)
		err := transportToPeer.SendMessage(directPeer, serialized)
		if err != nil {
			return err, NilRouteId
		}
	}
	return nil, localRouteId
}

// Blocks until message is confirmed received if reliably is true
func (self *Node) SendMessageThruRoute(routeId domain.RouteId, contents []byte, reliably bool) error {
	route, routeFound := self.safelyGetRoute(routeId)
	if !routeFound {
		return errors.New("Route not found")
	}

	base := domain.MessageBase{
		SendId:   route.ForwardRewriteSendId,
		SendBack: false,
		FromPeer: self.config.PubKey,
		Reliably: reliably,
		Nonce:    generateNonce(),
	}
	directPeer := route.ForwardToPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer, reliably)
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

// Blocks until message is confirmed received if reliably is true
func (self *Node) SendMessageBackThruRoute(replyTo ReplyTo, contents []byte, reliably bool) error {
	directPeer := replyTo.fromPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer, reliably)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No route or transport to peer %v\n", directPeer))
	}
	base := domain.MessageBase{
		SendId:   replyTo.routeId,
		SendBack: true,
		FromPeer: self.config.PubKey,
		Reliably: reliably,
		Nonce:    generateNonce(),
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
	return len(self.routesById)
}

func (self *Node) debug_countMessages() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return len(self.messagesBeingAssembled)
}

type TestConfig struct {
	Reliable protocol.ReliableTransportConfig
	Udp      protocol.UDPConfig
	Node     domain.NodeConfig

	PeersToConnect    []domain.Peer
	RoutesToEstablish []domain.RouteConfig
	MessagesToSend    []domain.MessageToSend
	MessagesToReceive []domain.MessageToReceive
}

func (self *TestConfig) AddPeerToConnect(addr string, config *TestConfig) {
	peerToConnect := domain.Peer{}
	peerToConnect.Peer = config.Node.PubKey
	peerToConnect.Info = protocol.CreateUDPCommConfig(addr, nil)
	self.PeersToConnect = append(self.PeersToConnect, peerToConnect)
}

func (self *TestConfig) AddRouteToEstablish(config *TestConfig) {
	routeToEstablish := domain.RouteConfig{}
	routeToEstablish.Id = uuid.NewV4()
	routeToEstablish.Peers = append(routeToEstablish.Peers, config.Node.PubKey)
	self.RoutesToEstablish = append(self.RoutesToEstablish, routeToEstablish)
}

func (self *TestConfig) AddPeerToRoute(indexRoute int, config *TestConfig) {
	self.RoutesToEstablish[indexRoute].Peers = append(self.RoutesToEstablish[indexRoute].Peers, config.Node.PubKey)
}

func (self *TestConfig) AddMessageToSend(thruRoute uuid.UUID, message string, reliably bool) {
	messageToSend := domain.MessageToSend{}
	messageToSend.ThruRoute = thruRoute
	messageToSend.Contents = []byte(message)
	messageToSend.Reliably = reliably
	self.MessagesToSend = append(self.MessagesToSend, messageToSend)
}

func (self *TestConfig) AddMessageToReceive(messageReceive, messageReply string, replyReliably bool) {
	messageToReceive := domain.MessageToReceive{}
	messageToReceive.Contents = []byte(messageReceive)
	messageToReceive.Reply = []byte(messageReply)
	messageToReceive.ReplyReliably = replyReliably
	self.MessagesToReceive = append(self.MessagesToReceive, messageToReceive)
}

// Add transport to Node
func (self *Node) AddTransportToNode(config TestConfig) {
	udpTransport := protocol.CreateNewUDPTransport(config.Udp)

	// Connect
	for _, connectTo := range config.PeersToConnect {
		connectError := udpTransport.ConnectToPeer(connectTo.Peer, connectTo.Info)
		if connectError != nil {
			panic(connectError)
		}
	}

	// Reliable transport closes UDPTransport
	reliableTransport := protocol.NewReliableTransport(udpTransport, config.Reliable)
	//defer reliableTransport.Close()

	self.AddTransport(reliableTransport)
}

// Add Routes to Node
func (self *Node) AddRoutesToEstablish(config TestConfig) {
	// Setup route
	for _, routeConfig := range config.RoutesToEstablish {
		if len(routeConfig.Peers) == 0 {
			continue
		}
		addRouteErr := self.AddRoute((domain.RouteId)(routeConfig.Id), routeConfig.Peers[0])
		if addRouteErr != nil {
			panic(addRouteErr)
		}
		for peer := 1; peer < len(routeConfig.Peers); peer++ {
			extendErr := self.ExtendRoute((domain.RouteId)(routeConfig.Id), routeConfig.Peers[peer], 5*time.Second)
			if extendErr != nil {
				panic(extendErr)
			}
		}
	}
}

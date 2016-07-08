package mesh

import(
	"os"
	"fmt"
	"time"
    "sync"
    "errors"
    "reflect"
    "runtime/debug"
    "gopkg.in/op/go-logging.v1")

import(
	"github.com/skycoin/skycoin/src/mesh2/transport"
	"github.com/skycoin/skycoin/src/mesh2/serialize"
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid")

type NodeConfig struct {
	PubKey		    			cipher.PubKey
	ChaCha20Key	    			[32]byte
	MaximumForwardingDuration	time.Duration
	RefreshRouteDuration		time.Duration
	ExpireMessagesInterval      time.Duration
	ExpireRoutesInterval		time.Duration
	TimeToAssembleMessage		time.Duration
	TransportMessageChannelLength int
}

func min(a, b uint64) uint64 {
    if a < b {
        return a
    }
    return b
}

type RouteId uuid.UUID
type messageId uuid.UUID

var NilRouteId RouteId = (RouteId)(uuid.Nil)

type rewriteableMessage interface {
    Rewrite(newSendId RouteId) rewriteableMessage
}

type MeshMessage struct {
    RouteId       RouteId
    Contents      []byte
}

type Route struct {
	// Forward should never be cipher.PubKey{}
	forwardToPeer 			cipher.PubKey
	forwardRewriteSendId 	RouteId

	backwardToPeer 			cipher.PubKey
	backwardRewriteSendId 	RouteId

	// time.Unix(0,0) means it lives forever
	lastRefresh				time.Time
}

type MessageUnderAssembly struct {
	fragments 				map[uint64]UserMessage
	sendId    				RouteId
	sendBack                bool
	count                   uint64
	dropped                 bool
	firstFragmentReceived	time.Time
}

type LocalRoute struct {
	terminatingPeer cipher.PubKey
	lastHopId       RouteId
	termSetRoute    SetRouteOp
}

type Node struct {
	config 						NodeConfig
    outputMessagesReceived 		chan MeshMessage
    transportsMessagesReceived 	chan []byte
	serializer 					*serialize.Serializer

    lock *sync.Mutex
    closeGroup *sync.WaitGroup
	closing 	chan bool

    transports 						map[transport.Transport]bool
    messagesBeingAssembled          map[messageId]*MessageUnderAssembly
    routesById                      map[RouteId]Route
    localRoutesByTerminatingPeer	map[cipher.PubKey]RouteId
    localRoutesById                 map[RouteId]LocalRoute
}

// Fields must be public (capital first letter) for encoder
type MessageBase struct {
	// If RouteId is unknown, but not cipher.PubKey{}, then the message should be received here
	//  the RouteId can be used to reply back thru the route
    SendId RouteId
    SendBack bool
}

type UserMessage struct {
	MessageBase
	MessageId messageId
	Index     uint64
	Count     uint64
	Contents  []byte
}

type SetRouteMessage struct {
	MessageBase
	SetRouteId     			RouteId
	ForwardToPeer 			cipher.PubKey
	ForwardRewriteSendId 	RouteId
	BackwardToPeer 			cipher.PubKey
	BackwardRewriteSendId 	RouteId
    DurationHint   			time.Duration
}

// Refreshes the route as it passes thru it
type RefreshRouteMessage struct {
	MessageBase
    DurationHint   time.Duration
}

// Deletes the route as it passes thru it
type DeleteRouteMessage struct {
	MessageBase
}

type TimeoutError struct {
}

func (self*TimeoutError) Error() string {
	return "Timeout"
}

var logger = logging.MustGetLogger("node")

// TODO: Transport crypto test

func NewNode(config NodeConfig) (*Node, error) {
	ret := &Node{
		config,
		nil,			// received
		make(chan []byte, config.TransportMessageChannelLength),			// received		
		serialize.NewSerializer(),
		&sync.Mutex{},	// Lock
		&sync.WaitGroup{},
		make(chan bool, 10),
		make(map[transport.Transport]bool),
		make(map[messageId]*MessageUnderAssembly),
		make(map[RouteId]Route),
		make(map[cipher.PubKey]RouteId),
		make(map[RouteId]LocalRoute),
	}
    ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, UserMessage{})
    ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, SetRouteMessage{})
    ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{3}, RefreshRouteMessage{})
    ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{4}, DeleteRouteMessage{})

    go ret.processIncomingMessagesLoop()
    go ret.expireOldRoutesLoop()
    go ret.expireOldMessagesLoop()

	return ret, nil
}

// Returns nil if reassembly didn't happen (incomplete message)
func (self*Node) reassembleUserMessage(msgIn UserMessage) (contents []byte) {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, routeExists := self.routesById[msgIn.SendId]
	if !routeExists {
        fmt.Fprintf(os.Stderr, "Message %v is to unknown route id %v, ignoring\n", msgIn.MessageId, msgIn.SendId)
		return nil
	}

	_, assembledExists := self.messagesBeingAssembled[msgIn.MessageId]
	if !assembledExists {
		beingAssembled := &MessageUnderAssembly{
			make(map[uint64]UserMessage),
			msgIn.SendId,
			msgIn.SendBack,
			msgIn.Count, 
			false, 
			time.Now().Add(self.config.TimeToAssembleMessage),
		}
		self.messagesBeingAssembled[msgIn.MessageId] = beingAssembled
	}

	beingAssembled, _ := self.messagesBeingAssembled[msgIn.MessageId]

	if beingAssembled.dropped {
		return nil
	}

	if beingAssembled.count != msgIn.Count {
        fmt.Fprintf(os.Stderr, "Fragments of message %v have different total counts!\n", msgIn.MessageId)
        beingAssembled.dropped = true
		return nil
	}

	if beingAssembled.sendId != msgIn.SendId {
        fmt.Fprintf(os.Stderr, "Fragments of message %v have different send ids!\n", msgIn.SendId)
        beingAssembled.dropped = true
		return nil
	}

	if beingAssembled.sendBack != msgIn.SendBack {
        fmt.Fprintf(os.Stderr, "Fragments of message %v have different send directions!\n", msgIn.SendId)
        beingAssembled.dropped = true
		return nil
	}

	_, messageExists := beingAssembled.fragments[msgIn.Index]
	if messageExists {
        fmt.Fprintf(os.Stderr, "Fragment %v of message %v is duplicated, dropping message\n", msgIn.Index, msgIn.MessageId)
		return nil
	}

	beingAssembled.fragments[msgIn.Index] = msgIn
	if (uint64)(len(beingAssembled.fragments)) == beingAssembled.count {
		delete(self.messagesBeingAssembled, msgIn.MessageId)
		reassembled := []byte{}
		for i := (uint64)(0); i < beingAssembled.count; i++ {
			reassembled = append(reassembled, beingAssembled.fragments[i].Contents...)
		}
		return reassembled
	}

	return nil
}

func getMessageBase(msg interface{}) (base MessageBase) {
    msg_type := reflect.TypeOf(msg) 

	if msg_type == reflect.TypeOf(UserMessage{}) {
		return (msg.(UserMessage)).MessageBase
	} else if msg_type == reflect.TypeOf(SetRouteMessage{}) {
		return (msg.(SetRouteMessage)).MessageBase
	} else if msg_type == reflect.TypeOf(RefreshRouteMessage{}) {
		return (msg.(RefreshRouteMessage)).MessageBase
	} else if msg_type == reflect.TypeOf(DeleteRouteMessage{}) {
		return (msg.(DeleteRouteMessage)).MessageBase
	}
	debug.PrintStack()
	panic(fmt.Sprintf("Internal error: getMessageBase incomplete (%v)", msg_type))
}

func rewriteMessage(msg interface{}, newBase MessageBase) (rewritten interface{}) {
    msg_type := reflect.TypeOf(msg) 

	if msg_type == reflect.TypeOf(UserMessage{}) {
		ret := (msg.(UserMessage))
		ret.MessageBase = newBase
		return ret
	} else if msg_type == reflect.TypeOf(SetRouteMessage{}) {
		ret := (msg.(SetRouteMessage))
		ret.MessageBase = newBase
		return ret
	} else if msg_type == reflect.TypeOf(RefreshRouteMessage{}) {
		ret := (msg.(RefreshRouteMessage))
		ret.MessageBase = newBase
		return ret
	} else if msg_type == reflect.TypeOf(DeleteRouteMessage{}) {
		ret := (msg.(DeleteRouteMessage))
		ret.MessageBase = newBase
		return ret
	}
	panic("Internal error: rewriteMessage incomplete")
}

func (self*Node) safelyGetForwarding(msg interface{}) (sendBack bool, route Route, doForward bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	messageBase := getMessageBase(msg)
	routeFound, foundRoute := self.routesById[messageBase.SendId]

	doForward = foundRoute

	if messageBase.SendBack {
		if routeFound.backwardToPeer == (cipher.PubKey{}) {
			doForward = false
		}
	} else {
		if routeFound.forwardToPeer == (cipher.PubKey{}) {
			doForward = false
		}
	}

	if doForward {
		return messageBase.SendBack, routeFound, doForward
	} else {
		return false, Route{}, doForward
	}
}

func (self*Node) safelyGetRoute(id RouteId) (route Route, exists bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	route, exists = self.routesById[id]
	return
}

func (self*Node) safelyGetRewriteBase(msg interface{}) (forwardTo cipher.PubKey, base MessageBase, doForward bool) {
	// sendBack
	sendBack, route, foundRoute := self.safelyGetForwarding(msg)
	if !foundRoute {
		return cipher.PubKey{}, MessageBase{}, false
	}
	forwardTo = route.forwardToPeer
	rewriteTo := route.forwardRewriteSendId
	if sendBack {
		forwardTo = route.backwardToPeer
		rewriteTo = route.backwardRewriteSendId
	}
	if forwardTo == (cipher.PubKey{}) {
		return cipher.PubKey{}, MessageBase{}, false
	}
	newBase := 
		MessageBase {
			rewriteTo,
			sendBack,
		}
	return forwardTo, newBase, true
}

func (self*Node) forwardMessage(msg interface{}) bool {
	forwardTo, newBase, doForward := self.safelyGetRewriteBase(msg)
	if !doForward {
		return false
	}
	// Rewrite
	rewritten := rewriteMessage(msg, newBase)
fmt.Fprintf(os.Stderr, "fwd %v base %v -> rewritten %v\n", reflect.TypeOf(msg), getMessageBase(msg), rewritten)
	transport := self.safelyGetTransportToPeer(forwardTo, true)
	if transport == nil {
        fmt.Fprintf(os.Stderr, "No transport found for forwarded message from %v to %v, dropping\n", self.config.PubKey, forwardTo)
        return true
   	}

	serialized := self.serializer.SerializeMessage(rewritten)
	send_error := transport.SendMessage(forwardTo, serialized)
	if send_error != nil {
        fmt.Fprintf(os.Stderr, "Failed to send forwarded message, dropping\n")
        return true
	}

	// Forward, not receive
	return true
}

func (self*Node) processUserMessage(msgIn UserMessage) {
	reassembled := self.reassembleUserMessage(msgIn)
	// Not finished reassembling yet
	if reassembled == nil {
		return
	}

	directPeer, forwardBase, doForward := self.safelyGetRewriteBase(msgIn)
	if doForward {
fmt.Fprintf(os.Stderr, "fwd usermsg from %v to %v\n", self.config.PubKey, directPeer)
		transport := self.safelyGetTransportToPeer(directPeer, true)
		if transport == nil {
			fmt.Fprintf(os.Stderr, "No transport to peer %v from %v, dropping\n", directPeer, self.config.PubKey)
			return
		}
		// Forward reassembled message, not individual pieces. This is done because of the need for refragmentation
		fragments := self.fragmentMessage(reassembled, directPeer, transport, forwardBase)
		for _, fragment := range(fragments) {
				serialized := self.serializer.SerializeMessage(fragment)
				send_error := transport.SendMessage(directPeer, serialized)
				if send_error != nil {
			        fmt.Fprintf(os.Stderr, "Failed to send forwarded message, dropping\n")
			        return
				}
		}
	} else {
		self.outputMessagesReceived <- MeshMessage{msgIn.SendId, reassembled}
	}
}

func (self*Node) processSetRouteMessage(msg SetRouteMessage) {
fmt.Fprintf(os.Stderr, "peer (%v) processSetRouteMessage: %v\n", self.config.PubKey, msg)
	if msg.SendBack {
		fmt.Fprintf(os.Stderr, "Invalid SetRouteMessage received, dropping: %v\n", msg)
		return
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	
	op := msg.OnReceive
	if op.SetRouteId == NilRouteId {
		fmt.Fprintf(os.Stderr, "Invalid SetRouteOp received, dropping: %v\n", op)
		return
	}
	self.routesById[op.SetRouteId] = 
		Route{
			op.ForwardToPeer,
			op.ForwardRewriteSendId,
			op.BackwardToPeer,
			op.BackwardRewriteSendId,
			time.Now(),
		}
}

func (self*Node) processMessage(serialized []byte) {
    msg, deserialize_error := self.serializer.UnserializeMessage(serialized)
    if deserialize_error != nil {
        fmt.Fprintf(os.Stderr, "Deserialization error %v\n", deserialize_error)
        return
    }

    msg_type := reflect.TypeOf(msg) 
fmt.Fprintf(os.Stderr, "peer %v got %v\n", self.config.PubKey, msg_type)
    // User messages have fragmentation to deal with
    if msg_type == reflect.TypeOf(UserMessage{}) {
		self.processUserMessage(msg.(UserMessage))
    } else {
	    if !self.forwardMessage(msg) {
			// Receive or forward. Refragment on forward!
			if msg_type == reflect.TypeOf(SetRouteMessage{}) {
				self.processSetRouteMessage(msg.(SetRouteMessage))
			} else {
		        fmt.Fprintf(os.Stderr, "Unknown message type received\n")
		        return
			}
		}
	}
}

func (self*Node) expireOldMessages() {
	fmt.Fprintf(os.Stderr, "TODO: expireOldMessages\n")
}

func (self*Node) expireOldMessagesLoop() {
	for len(self.closing) == 0 {
		select {
			case <-time.After(self.config.ExpireMessagesInterval): {
				self.expireOldMessages()
			}
			case <-self.closing: {
				return
			}
		}
	}
}

func (self*Node) processIncomingMessagesLoop() {
   	for len(self.closing) == 0 {
		select {
			case msg, ok := <- self.transportsMessagesReceived: {
				if ok {
					self.processMessage(msg)
				}
			}
			case <- self.closing: {
				return
			}
		}
	}
}

func (self*Node) expireOldRoutes() {
	// Last refresh of time.Unix(0,0) means it lives forever
	fmt.Fprintf(os.Stderr, "TODO: expireOldRoutes\n")

}

func (self*Node) expireOldRoutesLoop() {
	for len(self.closing) == 0 {
		select {
			case <-time.After(self.config.ExpireRoutesInterval): {
				self.expireOldRoutes()
			}
			case <-self.closing: {
				return
			}
		}
	}
}

// Waits for transports to close
func (self*Node) Close() error {
	close(self.transportsMessagesReceived)
	self.closeGroup.Wait()
	return nil
}

func (self*Node) GetConfig() NodeConfig {
	return self.config
}

// Node takes ownership of the transport, and will call Close() when it is closed
func (self*Node) AddTransport(transport transport.Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	transport.SetReceiveChannel(self.transportsMessagesReceived)
	self.transports[transport] = true
}

func (self*Node) RemoveTransport(transport transport.Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.transports, transport)
}

func (self*Node) GetTransports() ([]transport.Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []transport.Transport{}
	for transport, _ := range(self.transports) {
		ret = append(ret, transport)
	}
	return ret
}

func (self*Node) GetConnectedPeers() ([]cipher.PubKey) {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []cipher.PubKey{}
	for transport, _ := range(self.transports) {
		peers := transport.GetConnectedPeers()
		ret = append(ret, peers...)
	}
	return ret
}

func (self*Node) ConnectedToPeer(peer cipher.PubKey) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	for transport, _ := range(self.transports) {
		if transport.ConnectedToPeer(peer) {
			return true
		}
	}
	return false
}

// Message order is not preserved
func  (self*Node) SetReceiveChannel(received chan MeshMessage) {
	self.outputMessagesReceived = received
}

// toPeer must be the public key of a connected peer
// Requires reliable transport (for now)
func (self*Node) AddRoute(id RouteId, toPeer cipher.PubKey) error {	
	_, routeExists := self.safelyGetRoute(id)
	if routeExists {
		return errors.New(fmt.Sprintf("Rotue %v already exists\n", id))
	}

	transport := self.safelyGetTransportToPeer(toPeer, true)
	if transport == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", toPeer))
	}
fmt.Fprintf(os.Stderr, "AddRoute %v\n", toPeer)
	// Send message
	setOp := SetRouteOp{
					id,
					// ForwardToPeer
					cipher.PubKey{},
					NilRouteId,
					// BackwardToPeer
					self.config.PubKey,
					NilRouteId,
					// Route lifetime hint
				    3*self.config.RefreshRouteDuration,
			    }
	message :=
		SetRouteMessage{
			MessageBase{
				NilRouteId,
				false,
			},
			setOp,
			false,
			SetRouteOp{},
		}

	serialized := self.serializer.SerializeMessage(message)
	send_error := transport.SendMessage(toPeer, serialized)
	if send_error != nil {
		return send_error
	}

	// Add locally to routesById for backward termination
	self.lock.Lock()
	defer self.lock.Unlock()
	self.routesById[id] = 
		Route{
			toPeer,
			id,
			cipher.PubKey{},
			NilRouteId,
			// Route lifetime: never dies until route is removed
		    time.Unix(0,0),
		}

	self.localRoutesByTerminatingPeer[toPeer] = id
	self.localRoutesById[id] = LocalRoute{toPeer, id, setOp}
	return nil
}

func (self*Node) extendRouteWithoutSending(id RouteId, toPeer cipher.PubKey) (message SetRouteMessage, directPeer cipher.PubKey, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	newHopId := (RouteId)(uuid.NewV4())

	localRoute, localRouteExists := self.localRoutesById[id]
	if !localRouteExists {
		return SetRouteMessage{}, cipher.PubKey{}, errors.New(fmt.Sprintf("ExtendRoute called on unknown route: %v", id))
	}

	route, routeExists := self.routesById[id]
	if !routeExists {
		panic("Internal consistency error 8JUL2016544")
	}

	directPeer = route.forwardToPeer
	termMsg := localRoute.termSetRoute

	sendBase := MessageBase{
		route.forwardRewriteSendId,
		false,
	}

	newTermOp := SetRouteOp{
		// SetRouteId
		newHopId,
		// ForwardToPeer
		cipher.PubKey{},
		NilRouteId,
		// BackwardToPeer
		localRoute.terminatingPeer,
		termMsg.SetRouteId,
		// Route lifetime hint
	    3*self.config.RefreshRouteDuration,
	}

	message = SetRouteMessage {
		sendBase,
		SetRouteOp{
				// SetRouteId
				localRoute.lastHopId,
				// ForwardToPeer
				toPeer,
				// ForwardRewriteSendId
				newHopId,

				// Unchanged...
				termMsg.BackwardToPeer,
				termMsg.BackwardRewriteSendId,
			    termMsg.DurationHint,
		},
		true,
		newTermOp,
	}

	delete(self.localRoutesByTerminatingPeer, localRoute.terminatingPeer)
	self.localRoutesById[id] = LocalRoute{toPeer, newHopId, newTermOp}
	self.localRoutesByTerminatingPeer[toPeer] = id

	return message, directPeer, nil
}

// toPeer must be the public key of a peer connected to the current last node in this route
// Blocks until the operation is completed
func (self*Node) ExtendRoute(id RouteId, toPeer cipher.PubKey) error {
	message, directPeer, extendError := self.extendRouteWithoutSending(id, toPeer)
	if extendError != nil {
		return extendError
	}
	transport := self.safelyGetTransportToPeer(directPeer, true)
	if transport == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v from %v\n", directPeer, self.config.PubKey))
	}
	serialized := self.serializer.SerializeMessage(message)
	send_error := transport.SendMessage(directPeer, serialized)
	if send_error != nil {
		return send_error
	}
fmt.Fprintf(os.Stderr, "extended directPeer %v, toPeer %v, message %v (directPeer %v)\n", directPeer, toPeer, message, directPeer)
	return nil
}

func (self*Node) RemoveRoute(id RouteId) (error) {
	// routesById
	// localRoutesByTerminatingPeer
	// localRoutesById
	return errors.New("todo")
}

func (self*Node) getMaximumContentLength(toPeer cipher.PubKey, transport transport.Transport) uint64 {	
	transportSize := transport.GetMaximumMessageSizeToPeer(toPeer)
	empty := UserMessage{}
	emptySerialized := self.serializer.SerializeMessage(empty)
	if (uint)(len(emptySerialized)) >= transportSize {
		return 0
	}
	return (uint64)(transportSize) - (uint64)(len(emptySerialized))
}

func (self*Node) fragmentMessage(fullContents []byte, toPeer cipher.PubKey, transport transport.Transport, base MessageBase) []UserMessage {
	ret_noCount := make([]UserMessage, 0)
	maxContentLength := self.getMaximumContentLength(toPeer, transport)
	remainingBytes := fullContents[:]
	messageId := (messageId)(uuid.NewV4())
	for len(remainingBytes) > 0 {
		nBytesThisMessage := min(maxContentLength, (uint64)(len(remainingBytes)))
		bytesThisMessage := remainingBytes[:nBytesThisMessage]
		remainingBytes = remainingBytes[nBytesThisMessage:]
		message := UserMessage {
			base,
			messageId,
			(uint64)(len(ret_noCount)),
			0,
			bytesThisMessage,
		}
		ret_noCount = append(ret_noCount, message)
	}
	ret := make([]UserMessage, 0)
	for _, message := range(ret_noCount) {
		message.Count = (uint64)(len(ret_noCount))
		ret = append(ret, message)
	}
	return ret
}

func (self*Node) unsafelyGetTransportToPeer(peerKey cipher.PubKey, reliably bool) transport.Transport {
	// If unreliable, prefer unreliable transport
	if !reliably {
		for transport, _ := range(self.transports) {
			// TODO: Choose transport more intelligently
			if transport.ConnectedToPeer(peerKey) && !transport.IsReliable() {
				return transport
			}
		}	
	}
	for transport, _ := range(self.transports) {
		// TODO: Choose transport more intelligently
		if transport.ConnectedToPeer(peerKey) && ((!reliably) || transport.IsReliable()) {
			return transport
		}
	}
	return nil
}

func (self*Node) safelyGetTransportToPeer(peerKey cipher.PubKey, reliably bool) transport.Transport {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.unsafelyGetTransportToPeer(peerKey, reliably)
}

func (self*Node) findRouteToPeer(toPeer cipher.PubKey, reliably bool) (directPeer cipher.PubKey, routeId RouteId, transport transport.Transport, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	localRouteId, localRouteExists := self.localRoutesByTerminatingPeer[toPeer]
	if localRouteExists {
		route, routeStructExists := self.routesById[localRouteId]
		if !routeStructExists {
			panic("Local route struct does not exists")
		}
		directPeer = route.forwardToPeer
		routeId = route.forwardRewriteSendId
	} else {
		return cipher.PubKey{}, NilRouteId, nil, errors.New(fmt.Sprintf("No route to peer: %v", toPeer))
	}
	transport = self.unsafelyGetTransportToPeer(directPeer, reliably)
	if transport == nil {
		return cipher.PubKey{}, NilRouteId, nil, 
			    errors.New(fmt.Sprintf("No route or transport to peer %v\n", toPeer))
	}
	return
}

// Chooses a route automatically. Sends directly without a route if connected to that peer
// if reliably is false, unreliable transport is preferred, but reliable is chosen if it's the only option
// if reliably is true, reliable transport only is used
func (self*Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte, reliably bool) (err error, routeId RouteId) {
	directPeer, routeId, transport, error := self.findRouteToPeer(toPeer, reliably)
	if error != nil {
		return error, NilRouteId
	}
	base := MessageBase{
		routeId,
		false,		// Sending forward
	}
	messages := self.fragmentMessage(contents, directPeer, transport, base)
	for _, message := range(messages) {
		serialized := self.serializer.SerializeMessage(message)
		send_error := transport.SendMessage(directPeer, serialized)
		if send_error != nil {
			return send_error, NilRouteId
		}
	}
	return nil, routeId
}

// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageThruRoute(route_id RouteId, contents []byte, reliably bool) (error) {
//fragmentMessage()
	return errors.New("todo")
}

// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageBackThruRoute(route_id RouteId, contents []byte, reliably bool) (error) {
//fragmentMessage()
	return errors.New("todo")
}



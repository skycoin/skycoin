package mesh

import(
	"time"
    "sync"
    "errors"
    "reflect"
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
}

// Fields must be public (capital first letter) for encoder
type MessageBase struct {
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
	panic("Internal error: getMessageBase incomplete")
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

func (self*Node) forwardMessage(msg interface{}) bool {
	// sendBack
	_, _, foundRoute := self.safelyGetForwarding(msg)
	if !foundRoute {
		return false
	}

/*
	forwardTo := route.forwardToPeer
	if msgIn.SendBack {
		forwardTo = route.backwardToPeer
	}
*/

	panic("TODO: forward message\n")
	return true
}

func (self*Node) processUserMessage(msgIn UserMessage) {
	reassembled := self.reassembleUserMessage(msgIn)
	// Not finished reassembling yet
	if reassembled == nil {
		return
	}
	self.lock.Lock()
	_, routeExists := self.routesById[msgIn.SendId]
	if !routeExists {
		self.lock.Unlock()
		fmt.Fprintf(os.Stderr, "Dropping message %v to unknown route id %v\n", msgIn.MessageId, msgIn.SendId)
		return
	}
	self.lock.Unlock()

	self.outputMessagesReceived <- MeshMessage{msgIn.SendId, reassembled}
}

func (self*Node) processSetRouteMessage(msg SetRouteMessage) {
	if msg.SetRouteId == NilRouteId || msg.SendBack {
		fmt.Fprintf(os.Stderr, "Invalid SetRouteMessage received, dropping: %v\n", msg)
		return
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	self.routesById[msg.SetRouteId] = 
		Route{
			msg.ForwardToPeer,
			msg.ForwardRewriteSendId,
			msg.BackwardToPeer,
			msg.BackwardRewriteSendId,
			time.Now(),
		}
}

func (self*Node) processMessage(serialized []byte) {
    msg, deserialize_error := self.serializer.UnserializeMessage(serialized)
    if deserialize_error != nil {
        fmt.Fprintf(os.Stderr, "Deserialization error %v\n", deserialize_error)
        return
    }

    if !self.forwardMessage(msg) {
		// Receive or forward. Refragment on forward!
	    msg_type := reflect.TypeOf(msg) 

		if msg_type == reflect.TypeOf(UserMessage{}) {
			self.processUserMessage(msg.(UserMessage))
		} else if msg_type == reflect.TypeOf(SetRouteMessage{}) {
			self.processSetRouteMessage(msg.(SetRouteMessage))
		} else {
	        fmt.Fprintf(os.Stderr, "Unknown message type received\n")
	        return
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
func (self*Node) AddRoute(id RouteId, toPeer cipher.PubKey) error {	
	_, routeExists := self.safelyGetRoute(id)
	if routeExists {
		return errors.New(fmt.Sprintf("Rotue %v already exists\n", id))
	}

	transport := self.safelyGetTransportToPeer(toPeer)
	if transport == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", toPeer))
	}

	// Send message
	message :=
		SetRouteMessage{
			MessageBase{
				NilRouteId,
				false,
			},
			id,
			// ForwardToPeer
			cipher.PubKey{},
			NilRouteId,
			// BackwardToPeer
			cipher.PubKey{},
			NilRouteId,
			// Route lifetime hint
		    3*self.config.RefreshRouteDuration,
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
	return nil
}

// toPeer must be the public key of a peer connected to the current last node in this route
// Blocks until the operation is completed
func (self*Node) ExtendRoute(id RouteId, toPeer cipher.PubKey) error {
	// localRoutesByTerminatingPeer
	return errors.New("todo")
}

func (self*Node) RemoveRoute(id RouteId) (error) {
	// routesById
	// localRoutesByTerminatingPeer
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

func (self*Node) unsafelyGetTransportToPeer(peerKey cipher.PubKey) transport.Transport {
	for transport, _ := range(self.transports) {
		// TODO: Choose transport more intelligently
		if transport.ConnectedToPeer(peerKey) {
			return transport
		}
	}
	return nil
}

func (self*Node) safelyGetTransportToPeer(peerKey cipher.PubKey) transport.Transport {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.unsafelyGetTransportToPeer(peerKey)
}

func (self*Node) findRouteToPeer(toPeer cipher.PubKey, reliably bool) (directPeer cipher.PubKey, routeId RouteId, transport transport.Transport, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	localRouteIdFound, localRouteExists := self.localRoutesByTerminatingPeer[toPeer]
	if localRouteExists {
		route, routeStructExists := self.routesById[localRouteIdFound]
		if !routeStructExists {
			panic("Local route struct does not exists")
		}
		directPeer = route.forwardToPeer
		routeId = route.forwardRewriteSendId
	} else {
		return cipher.PubKey{}, NilRouteId, nil, errors.New(fmt.Sprintf("No route to peer: %v", toPeer))
	}
	transport = self.unsafelyGetTransportToPeer(directPeer)
	if transport == nil {
		return cipher.PubKey{}, NilRouteId, nil, 
			    errors.New(fmt.Sprintf("No route or transport to peer %v\n", toPeer))
	}
	return
}

// Chooses a route automatically. Sends directly without a route if connected to that peer
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
func (self*Node) SendMessageThruRoute(route_id RouteId, contents []byte, reliably bool,) (error) {
//fragmentMessage()
	return errors.New("todo")
}

// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageBackThruRoute(route_id RouteId, contents []byte, reliably bool) (error) {
//fragmentMessage()
	return errors.New("todo")
}



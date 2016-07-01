package mesh

import(
"os"
	"fmt"
	"time"
    "sync"
    "errors"
    "reflect"
    "gopkg.in/op/go-logging.v1")

import(
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid")

type NodeConfig struct {
	PubKey		cipher.PubKey
	ChaCha20Key	[32]byte
	DefaultRetransmitTime	time.Duration
	AckKeepTime				time.Duration
}

type LocalRouteId uuid.UUID
type RouteId uint32
type messageId uuid.UUID

var NilRouteId RouteId = 0

type MeshMessage struct {
    RouteId       RouteId
    Contents      []byte
}

type LocalRoute struct {
	connectedPeer cipher.PubKey
	routeId       RouteId
}

type Node struct {
	config 						NodeConfig
    outputMessagesReceived 		chan MeshMessage
    transportsMessagesReceived 	chan []byte
	serializer *Serializer

    lock *sync.Mutex
    closeGroup *sync.WaitGroup

    transports 						map[Transport]bool
    messagesAwaitingConfirmation 	map[messageId]chan ReliableMessageAck
    routesById						map[LocalRoute]LocalRouteId
}

// Fields must be public (capital first letter) for encoder
type MessageBase struct {
    SendId RouteId
    SendBack bool
}

type UnreliableUserMessage struct {
	MessageBase
	Contents  	[]byte
}

type ReliableMessage struct {
	MessageBase
	MsgId		messageId
}

type ReliableMessageAck struct {
	MessageBase
	MsgId       messageId
}

// UnreliableMessage
type EstablishForwardingMessage struct {
	MessageBase
	// toPubKey = 0 for terminate (receive on this hop)
    ToPubKey cipher.PubKey
    FromPubKey cipher.PubKey
    BackwardRewriteSendId messageId
    DurationHint time.Duration
}

type EstablishForwardingReply struct {
	ReliableMessage
    
    NewSendId uint32
    // Secret sent back in AddRewriteMessage
    Secret string
}

type AddRewriteMessage struct {
	ReliableMessage

    // Secret from EstablishRouteReply
    Secret string
    RewriteSendId uint32
}

type ReliableUserMessage struct {
	ReliableMessage
	Contents 	[]byte	
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
		make(chan []byte),			// received		
		NewSerializer(),
		&sync.Mutex{},	// Lock
		&sync.WaitGroup{},
		make(map[Transport]bool),
		make(map[messageId]chan ReliableMessageAck),
		make(map[LocalRoute]LocalRouteId),
	}
    ret.serializer.RegisterMessageForSerialization(MessagePrefix{1}, ReliableUserMessage{})
    ret.serializer.RegisterMessageForSerialization(MessagePrefix{2}, ReliableMessageAck{})
    ret.serializer.RegisterMessageForSerialization(MessagePrefix{3}, UnreliableUserMessage{})
	go func() {
		ret.closeGroup.Add(1)
		defer ret.closeGroup.Done()
		for {
			msg, more := <- ret.transportsMessagesReceived
			if !more {
				break
			}
			ret.processMessage(msg)
		}
	}()
	return ret, nil
}

// Waits for transports to close
func (self*Node) Close() error {
	close(self.transportsMessagesReceived)
	self.closeGroup.Wait()
	return nil
}

func (self*Node) ForwardOrReceive(msg MessageBase, contents []byte) {
	// TODO: Check routes
fmt.Fprintf(os.Stderr, "ForwardOrReceive %v\n", msg)
	routeId := NilRouteId
	self.outputMessagesReceived <- MeshMessage{routeId, contents}
}

func (self*Node) expireMessageConfirm(msgId messageId, deadline time.Time) {
	// TODO
	fmt.Fprintf(os.Stderr, "expireMessageConfirm\n")
}

func (self*Node) SendAck(msg ReliableMessage) {
panic("todo")
/*
	if msg.SendId != NilRouteId {
		panic("Routes")
	}
	// Reply all the way back
	msg_out := ReliableMessageAck{
		MessageBase{NilRouteId, true},
		msg.msgId,
	}
	serialized := self.serializer.SerializeMessage(msg_out)
	//self.sendMessageUnreliably(toPeer cipher.PubKey, NilRouteId, backward bool, serialized []byte)
	self.expireMessageConfirm(msg.MsgId, time.Now().Add(self.config.AckKeepTime))
*/
}

func (self*Node) processMessage(datagram []byte) {
	msg_i, deserialize_error := self.serializer.UnserializeMessage(datagram)
	if deserialize_error != nil {
		logger.Debug("Deserialization error %v msg_i %v\n", deserialize_error, msg_i)
		return
	}
	msg_type := reflect.TypeOf(msg_i)
	if msg_type == reflect.TypeOf(ReliableMessageAck{}) {
		msg := msg_i.(ReliableMessageAck)
		self.lock.Lock()
		conf_chan, exists := self.messagesAwaitingConfirmation[msg.MsgId]
		self.lock.Unlock()
		if exists {
			conf_chan <- msg	
		}
	} else if msg_type == reflect.TypeOf(UnreliableUserMessage{}) {
		msg := msg_i.(UnreliableUserMessage)
		self.ForwardOrReceive(msg.MessageBase, msg.Contents)
	} else if msg_type == reflect.TypeOf(ReliableUserMessage{}) {
		msg := msg_i.(ReliableUserMessage)
		self.ForwardOrReceive(msg.MessageBase, msg.Contents)
		// Send ack
		self.SendAck(msg.ReliableMessage)
	}
}

func (self*Node) GetConfig() NodeConfig {
	return self.config
}

func (self*Node) safelyGetTransportToPeer(peerKey cipher.PubKey) Transport {
	self.lock.Lock()
	defer self.lock.Unlock()
	for transport, _ := range(self.transports) {
		// TODO: Choose transport more intelligently
		if transport.ConnectedToPeer(peerKey) {
			return transport
		}
	}
	return nil
}

// Node takes ownership of the transport, and will call Close() when it is closed
func (self*Node) AddTransport(transport Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	transport.SetReceiveChannel(self.transportsMessagesReceived)
	self.transports[transport] = true
}

func (self*Node) RemoveTransport(transport Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.transports, transport)
}

func (self*Node) GetTransports() ([]Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []Transport{}
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

// toPeer must be the public key of a connected peer
func (self*Node) AddRoute(id LocalRouteId, toPeer cipher.PubKey) error {
//Direct, go thru transports
	return errors.New("todo")
}

// toPeer must be the public key of a peer connected to the current last node in this route
// Blocks until the operation is completed
func (self*Node) ExtendRoute(id LocalRouteId, toPeer cipher.PubKey) error {
// blocks waiting
	return errors.New("todo")
}

func (self*Node) RemoveRoute(id LocalRouteId) (error) {
	return errors.New("todo")
}

// Returns route thru which it was sent
func (self*Node) sendMessageUnreliably(toPeer cipher.PubKey, routeId RouteId, backward bool, serialized []byte) (err error) {
	if routeId == NilRouteId {
		transport := self.safelyGetTransportToPeer(toPeer)
		// Send directly
		if transport == nil {
			return errors.New(fmt.Sprintf("No transport to peer: %v", toPeer))
		}
		return transport.SendMessage(toPeer, serialized)
	}
	// TODO: Send thru route
	return errors.New("todo: routes")
}

func (self*Node) getRetransmitTimeForRoute(toPeer cipher.PubKey, routeFound RouteId) time.Duration {
	return self.config.DefaultRetransmitTime
}

// Chooses a route automatically. Sends directly without a route if connected to that peer. 
// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte, reliably bool, timeout time.Duration) (err error, routeId RouteId) {
	msgId := (messageId)(uuid.NewV4())
	var confirm_channel chan ReliableMessageAck
	var serialized []byte
	if reliably {
		// TODO: Route
		msg_out := ReliableUserMessage{
					ReliableMessage{
						MessageBase{NilRouteId, false},
						msgId},
					contents}
		serialized = self.serializer.SerializeMessage(msg_out)
		confirm_channel = make(chan ReliableMessageAck, 1)
		self.lock.Lock()
		self.messagesAwaitingConfirmation[msgId] = confirm_channel
		self.expireMessageConfirm(msgId, time.Now().Add(self.config.AckKeepTime))
		defer self.lock.Unlock()
		defer func() {
			defer self.lock.Unlock()
			self.lock.Lock()
			delete(self.messagesAwaitingConfirmation, msgId)
		}()
	} else {
		// TODO: Route
		msg_out := UnreliableUserMessage{
			MessageBase{NilRouteId, false},
			contents,
		}
		serialized = self.serializer.SerializeMessage(msg_out)
	}
	for {
		// TODO: Find route
		routeFound := NilRouteId
		send_err := self.sendMessageUnreliably(toPeer, routeFound, false, serialized)
		if send_err != nil {
			return send_err, NilRouteId
		}
		if !reliably {
			return nil, routeFound
		}
		retransmitDuration := self.getRetransmitTimeForRoute(toPeer, routeFound)
		select {
			case <-confirm_channel:
				return nil, routeFound
			case <-time.After(timeout):
				return &TimeoutError{}, NilRouteId
			case <-time.After(retransmitDuration):
				// Continue loop
		}
	}
}

// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageThruRoute(route_id RouteId, contents []byte, reliably bool, deadline time.Time) (error) {
	return errors.New("todo")
}

// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageBackThruRoute(replyRoute RouteId, contents []byte, reliably bool, deadline time.Time) (error) {
	return errors.New("todo")
}

// Message order is not preserved
func  (self*Node) SetReceiveChannel(received chan MeshMessage) {
	self.outputMessagesReceived = received
}



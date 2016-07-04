package mesh

import(
/*
"os"
	"fmt"
	*/
	"time"
    "sync"
    "errors"
//    "reflect"
    "gopkg.in/op/go-logging.v1")

import(
	"github.com/skycoin/skycoin/src/mesh2/transport"
	"github.com/skycoin/skycoin/src/mesh2/serialize"
	"github.com/skycoin/skycoin/src/mesh2/reliable"
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid")

type NodeConfig struct {
	PubKey		    			cipher.PubKey
	ChaCha20Key	    			[32]byte
	MaximumForwardingDuration	time.Duration
	RefreshRouteDuration		time.Duration
}

type LocalRouteId uuid.UUID
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

type LocalRoute struct {
	connectedPeer cipher.PubKey
	routeId       RouteId
}

type Route struct {
	forwardToPeer 			cipher.PubKey
	forwardRewriteSendId 	RouteId

	backwardToPeer 			cipher.PubKey
	backwardRewriteSendId 	RouteId
}

type Node struct {
	config 						NodeConfig
    outputMessagesReceived 		chan MeshMessage
    transportsMessagesReceived 	chan []byte
	serializer 					*serialize.Serializer

    lock *sync.Mutex
    closeGroup *sync.WaitGroup

    transports 						map[transport.Transport]bool
    reliableTransports				map[transport.Transport]reliable.ReliableTransport
    localRoutesById					map[LocalRouteId]LocalRoute
    routesById                      map[messageId]Route
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
	SetRouteId     RouteId
	ForwardToPeer  cipher.PubKey
	BackwardToPeer cipher.PubKey
    DurationHint   time.Duration
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
		make(chan []byte),			// received		
		serialize.NewSerializer(),
		&sync.Mutex{},	// Lock
		&sync.WaitGroup{},
		make(map[transport.Transport]bool),
		make(map[transport.Transport]reliable.ReliableTransport),
		make(map[LocalRouteId]LocalRoute),
		make(map[messageId]Route),
	}
    ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, UserMessage{})
    ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, SetRouteMessage{})
    ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{3}, RefreshRouteMessage{})
    ret.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{4}, DeleteRouteMessage{})
	return ret, nil
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

func (self*Node) safelyGetTransportToPeer(peerKey cipher.PubKey) transport.Transport {
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

func (self*Node) getMaximumContentLength(transport transport.Transport) uint64 {
	return 0
}

// Chooses a route automatically. Sends directly without a route if connected to that peer. 
func (self*Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte, reliably bool, timeout time.Duration) (err error, routeId RouteId) {
//fragmentMessage()
	return nil, NilRouteId
}

// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageThruRoute(route_id RouteId, contents []byte, reliably bool, deadline time.Time) (error) {
//fragmentMessage()
	return errors.New("todo")
}

// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageBackThruRoute(replyRoute RouteId, contents []byte, reliably bool, deadline time.Time) (error) {
//fragmentMessage()
	return errors.New("todo")
}




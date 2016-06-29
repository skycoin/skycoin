package mesh

import(
	"time"
    "sync"
    "errors")

import(
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid")

type NodeConfig struct {
	PubKey		cipher.PubKey
	ChaCha20Key	[32]byte
}

type RouteId uint32

type MeshMessage struct {
    RouteId       uint32
    Contents      []byte
}

type Node struct {
	config NodeConfig
    outputMessagesReceived chan MeshMessage
    transportsMessagesReceived chan TransportMessage

    lock *sync.Mutex
    closeGroup *sync.WaitGroup

    transports map[Transport]bool
}

// TODO: Reliable / unreliable messages
// TODO: Store and forward
// TODO: Congestion control for reliable
// TODO: Transport crypto test

func NewNode(config NodeConfig) (*Node, error) {
	ret := &Node{
		config,
		nil,			// received
		make(chan TransportMessage),			// received
		&sync.Mutex{},	// Lock
		&sync.WaitGroup{},
		make(map[Transport]bool),
	}
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

func (self*Node) processMessage(msg TransportMessage) {
	// TODO: Reliability etc
	self.outputMessagesReceived <- MeshMessage{0, msg.Contents}
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
func (self*Node) AddRoute(id uuid.UUID, toPeer cipher.PubKey) error {
//Direct, go thru transports
	return errors.New("todo")
}

// toPeer must be the public key of a peer connected to the current last node in this route
// Blocks until the operation is completed
func (self*Node) ExtendRoute(id uuid.UUID, toPeer cipher.PubKey) error {
// blocks waiting
	return errors.New("todo")
}

func (self*Node) RemoveRoute(id uuid.UUID) (error) {
	return errors.New("todo")
}

// Chooses a route automatically. Sends directly without a route if connected to that peer. 
// Blocks until message is confirmed received if reliably is true
	// TODO: reliably, deadline
func (self*Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte, reliably bool, deadline time.Time) (err error, routeId uuid.UUID) {
	transport_msg := TransportMessage{toPeer, contents}
	transport := self.safelyGetTransportToPeer(toPeer)
	// Send directly
	if transport != nil {
		// TODO: reliably, deadline
		transport.SendMessage(transport_msg)
		return nil, uuid.Nil
	}
	return errors.New("todo"), uuid.Nil
}

// Blocks until message is confirmed received if reliably is true
func (self*Node) SendMessageThruRoute(route_id uuid.UUID, contents []byte, reliably bool, deadline time.Time) (error) {
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



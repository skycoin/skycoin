package mesh

import(
    "sync"
    "errors")

import(
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid")

type NodeConfig struct {
	ChaCha20Key	[32]byte
}

type BackRoute struct {
	sendId          uint32
    connectedPeer   cipher.PubKey
}

type MeshMessage struct {
    BackRoute
    Contents        []byte
}

type Node struct {
    messagesReceived chan MeshMessage

    lock *sync.Mutex
    transports map[Transport]bool
}

// TODO: Reliable / unreliable messages
// TODO: Store and forward
// TODO: Congestion control for reliable

func NewNode(config NodeConfig) (*Node, error) {
	ret := &Node{
		nil,			// received
		&sync.Mutex{},	// Lock
		make(map[Transport]bool),
	}
	return ret, nil
}

// Waits for transports to close
func (self*Node) Close() error {
	return nil
}

// Node takes ownership of the transport, and will call Close() when it is closed
func (self*Node) AddTransport(transport Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()
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
func (self*Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte, reliably bool) (err error, routeId uuid.UUID) {
	return errors.New("todo"), uuid.NewV4()
}

func (self*Node) SendMessageThruRoute(route_id uuid.UUID, contents []byte, reliably bool) (error) {
	return errors.New("todo")
}

func (self*Node) SendMessageBackThruRoute(replyRoute BackRoute, contents []byte, reliably bool) (error) {
	return errors.New("todo")
}

// Message order is not preserved
func  (self*Node) SetReceiveChannel(received chan MeshMessage) {
	self.messagesReceived = received
}



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

type RouteStatusCallback func(routeId uuid.UUID, ready bool, establishedToHopIdx int)

type Node struct {
	routeStatusCB RouteStatusCallback
    lock *sync.Mutex
    transports map[Transport]bool
}

// TODO: Reliable / unreliable messages
// TODO: Store and forward
// TODO: Congestion control for reliable

func NewNode(config NodeConfig) (*Node, error) {
	ret := &Node{
		func(routeId uuid.UUID, ready bool, establishedToHopIdx int) {},
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
	return nil
}

func (self*Node) GetConnectedPeers() ([]cipher.PubKey) {
	return nil
}

func (self*Node) AddRoute(id uuid.UUID, peerPubKeys []cipher.PubKey) (error) {
	return errors.New("todo")
}

func (self*Node) RemoveRoute(id uuid.UUID) (error) {
	return errors.New("todo")
}

func (self*Node) SetRouteStatusCallback(cb RouteStatusCallback) {
	self.routeStatusCB = cb
}

// Chooses a route automatically
func (self*Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte) (err error, routeId uuid.UUID) {
	return errors.New("todo"), uuid.NewV4()
}

func (self*Node) SendMessageThruRoute(route_id uuid.UUID, contents []byte) (error) {
	return errors.New("todo")
}

func (self*Node) SendMessageBackThruRoute(replyRoute BackRoute, contents []byte) (error) {
	return errors.New("todo")
}



package mesh

import (
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/serialize"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

//var logger = logging.MustGetLogger("node")

type Node struct {
	Config                     NodeConfig
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
}

type NodeConfig struct {
	PubKey                        cipher.PubKey
	MaximumForwardingDuration     time.Duration
	RefreshRouteDuration          time.Duration
	ExpireRoutesInterval          time.Duration
	TransportMessageChannelLength int
	//ChaCha20Key                   [32]byte
}

func NewNode(config NodeConfig) (*Node, error) {
	node := &Node{
		Config:                     config,
		outputMessagesReceived:     nil,                                                     // received
		transportsMessagesReceived: make(chan []byte, config.TransportMessageChannelLength), // received
		serializer:                 serialize.NewSerializer(),
		lock:                       &sync.Mutex{}, // Lock
		closeGroup:                 &sync.WaitGroup{},
		closing:                    make(chan bool, 10),
		transports:                 make(map[transport.ITransport]bool),
		routes:                     make(map[domain.RouteID]domain.Route),
		localRoutesByTerminatingPeer:   make(map[cipher.PubKey]domain.RouteID),
		localRoutes:                    make(map[domain.RouteID]domain.LocalRoute),
		routeExtensionsAwaitingConfirm: make(map[domain.RouteID]chan bool),
		//myCrypto:                   &ChaChaCrypto{config.ChaCha20Key},
	}
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, domain.UserMessage{})
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, domain.SetRouteMessage{})
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{3}, domain.RefreshRouteMessage{})
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{4}, domain.DeleteRouteMessage{})
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{5}, domain.SetRouteReply{})

	go node.processIncomingMessagesLoop()
	go node.expireOldRoutesLoop()
	go node.refreshRoutesLoop()

	return node, nil
}

func (self *Node) GetConfig() NodeConfig {
	return self.Config
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

// Waits for transports to close
func (self *Node) Close() error {
	for i := 0; i < 10; i++ {
		self.closing <- true
	}
	close(self.transportsMessagesReceived)
	self.closeGroup.Wait()
	return nil
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

func (self *Node) GetTransportToPeer(peerKey cipher.PubKey) transport.ITransport {
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
	return self.GetTransportToPeer(peerKey)
}

func (self *Node) GetMaximumContentLength(peerID cipher.PubKey, emptySerializedMessage []byte) uint64 {
	trans := self.GetTransportToPeer(peerID)
	transportSize := trans.GetMaximumMessageSizeToPeer(peerID)
	if (uint)(len(emptySerializedMessage)) >= transportSize {
		return 0
	}
	return (uint64)(transportSize) - (uint64)(len(emptySerializedMessage))
}

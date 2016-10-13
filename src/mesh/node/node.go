package mesh

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/serialize"
	"github.com/skycoin/skycoin/src/mesh/transport"
	//"gopkg.in/op/go-logging.v1"
)

//var logger = logging.MustGetLogger("node")

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

func NewNode(config domain.NodeConfig) (*Node, error) {
	node := &Node{
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
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, domain.UserMessage{})
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, domain.SetRouteMessage{})
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{3}, domain.RefreshRouteMessage{})
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{4}, domain.DeleteRouteMessage{})
	node.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{5}, domain.SetRouteReply{})

	go node.processIncomingMessagesLoop()
	go node.expireOldRoutesLoop()
	go node.expireOldMessagesLoop()
	go node.refreshRoutesLoop()

	return node, nil
}

func (self *Node) GetConfig() domain.NodeConfig {
	return self.config
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

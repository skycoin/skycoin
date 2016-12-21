package nodemanager

import (
//	"strconv"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	mesh "github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/transport"
	//"github.com/skycoin/skycoin/src/mesh/transport/physical"
)

type TestConfig struct {
	TransportConfig transport.TransportConfig
	NodeConfig      mesh.NodeConfig

	PeersToConnect           []Peer
	PeerToPeers              map[string]*Peer
	RoutesConfigsToEstablish []RouteConfig
	MessagesToSend           []MessageToSend
	MessagesToReceive        []MessageToReceive
	ExternalAddress          string
	StartPort                int
	Port                     int
}

type RouteConfig struct {
	RouteID uuid.UUID
	Peers   []cipher.PubKey
}

type Peer struct {
	Peer cipher.PubKey
	Info string
}

type MessageToSend struct {
	ThruRoute uuid.UUID
	Contents  []byte
}

type MessageToReceive struct {
	Contents []byte
	Reply    []byte
}

func (self *TestConfig) AddPeerToConnect(config *TestConfig) {

	peerToConnect := makePeer(config.NodeConfig.PubKey, config.ExternalAddress, config.Port)
	ownPeer := makePeer(self.NodeConfig.PubKey, self.ExternalAddress, self.Port)

	self.PeerToPeers[ownPeer.Info] = peerToConnect
}

func (self *TestConfig) AddRouteToEstablish(config *TestConfig) {
	routeConfigToEstablish := RouteConfig{}
	routeConfigToEstablish.RouteID = uuid.NewV4()
	routeConfigToEstablish.Peers = append(routeConfigToEstablish.Peers, config.NodeConfig.PubKey)
	self.RoutesConfigsToEstablish = append(self.RoutesConfigsToEstablish, routeConfigToEstablish)
}

func (self *TestConfig) AddPeerToRoute(indexRoute int, config *TestConfig) {
	self.RoutesConfigsToEstablish[indexRoute].Peers = append(self.RoutesConfigsToEstablish[indexRoute].Peers, config.NodeConfig.PubKey)
}

func (self *TestConfig) AddMessageToSendThruRoute(thruRouteID uuid.UUID, message string) {
	messageToSend := MessageToSend{}
	messageToSend.ThruRoute = thruRouteID
	messageToSend.Contents = []byte(message)
	self.MessagesToSend = append(self.MessagesToSend, messageToSend)
}

func (self *TestConfig) AddMessageToSend(config *TestConfig, message string) {
	messageToSend := MessageToSend{}
	pubKey := config.NodeConfig.PubKey

	for _, routeConfig := range(self.RoutesConfigsToEstablish) {
		thruRouteID := routeConfig.RouteID
		peers := routeConfig.Peers
		if peers[len(peers) - 1] == pubKey {
			messageToSend.ThruRoute = thruRouteID
			messageToSend.Contents = []byte(message)
			self.MessagesToSend = append(self.MessagesToSend, messageToSend)
			break
		}
	}
}

func (self *TestConfig) AddMessageToReceive(messageReceive, messageReply string) {
	messageToReceive := MessageToReceive{}
	messageToReceive.Contents = []byte(messageReceive)
	messageToReceive.Reply = []byte(messageReply)
	self.MessagesToReceive = append(self.MessagesToReceive, messageToReceive)
}

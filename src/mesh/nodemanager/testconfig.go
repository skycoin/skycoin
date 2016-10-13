package nodemanager

import (
	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/transport"
	"github.com/skycoin/skycoin/src/mesh/transport/physical"
)

type TestConfig struct {
	TransportConfig transport.TransportConfig
	UDPConfig       physical.UDPConfig
	NodeConfig      domain.NodeConfig

	PeersToConnect           []Peer
	RoutesConfigsToEstablish []RouteConfig
	MessagesToSend           []MessageToSend
	MessagesToReceive        []MessageToReceive
}

type RouteConfig struct {
	ID    uuid.UUID
	Peers []cipher.PubKey
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

func (self *TestConfig) AddPeerToConnect(addr string, config *TestConfig) {
	peerToConnect := Peer{}
	peerToConnect.Peer = config.NodeConfig.PubKey
	peerToConnect.Info = physical.CreateUDPCommConfig(addr, nil)
	self.PeersToConnect = append(self.PeersToConnect, peerToConnect)
}

func (self *TestConfig) AddRouteToEstablish(config *TestConfig) {
	routeConfigToEstablish := RouteConfig{}
	routeConfigToEstablish.ID = uuid.NewV4()
	routeConfigToEstablish.Peers = append(routeConfigToEstablish.Peers, config.NodeConfig.PubKey)
	self.RoutesConfigsToEstablish = append(self.RoutesConfigsToEstablish, routeConfigToEstablish)
}

func (self *TestConfig) AddPeerToRoute(indexRoute int, config *TestConfig) {
	self.RoutesConfigsToEstablish[indexRoute].Peers = append(self.RoutesConfigsToEstablish[indexRoute].Peers, config.NodeConfig.PubKey)
}

func (self *TestConfig) AddMessageToSend(thruRoute uuid.UUID, message string) {
	messageToSend := MessageToSend{}
	messageToSend.ThruRoute = thruRoute
	messageToSend.Contents = []byte(message)
	self.MessagesToSend = append(self.MessagesToSend, messageToSend)
}

func (self *TestConfig) AddMessageToReceive(messageReceive, messageReply string) {
	messageToReceive := MessageToReceive{}
	messageToReceive.Contents = []byte(messageReceive)
	messageToReceive.Reply = []byte(messageReply)
	self.MessagesToReceive = append(self.MessagesToReceive, messageToReceive)
}

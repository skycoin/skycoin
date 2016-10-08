package nodemanager

import (
	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/transport"
	"github.com/skycoin/skycoin/src/mesh/transport/physical"
)

type TestConfig struct {
	Reliable   transport.ReliableTransportConfig
	UDPConfig  physical.UDPConfig
	NodeConfig domain.NodeConfig

	PeersToConnect           []domain.Peer
	RoutesConfigsToEstablish []domain.RouteConfig
	MessagesToSend           []domain.MessageToSend
	MessagesToReceive        []domain.MessageToReceive
}

func (self *TestConfig) AddPeerToConnect(addr string, config *TestConfig) {
	peerToConnect := domain.Peer{}
	peerToConnect.Peer = config.NodeConfig.PubKey
	peerToConnect.Info = physical.CreateUDPCommConfig(addr, nil)
	self.PeersToConnect = append(self.PeersToConnect, peerToConnect)
}

func (self *TestConfig) AddRouteToEstablish(config *TestConfig) {
	routeConfigToEstablish := domain.RouteConfig{}
	routeConfigToEstablish.ID = uuid.NewV4()
	routeConfigToEstablish.Peers = append(routeConfigToEstablish.Peers, config.NodeConfig.PubKey)
	self.RoutesConfigsToEstablish = append(self.RoutesConfigsToEstablish, routeConfigToEstablish)
}

func (self *TestConfig) AddPeerToRoute(indexRoute int, config *TestConfig) {
	self.RoutesConfigsToEstablish[indexRoute].Peers = append(self.RoutesConfigsToEstablish[indexRoute].Peers, config.NodeConfig.PubKey)
}

func (self *TestConfig) AddMessageToSend(thruRoute uuid.UUID, message string, reliably bool) {
	messageToSend := domain.MessageToSend{}
	messageToSend.ThruRoute = thruRoute
	messageToSend.Contents = []byte(message)
	messageToSend.Reliably = reliably
	self.MessagesToSend = append(self.MessagesToSend, messageToSend)
}

func (self *TestConfig) AddMessageToReceive(messageReceive, messageReply string, replyReliably bool) {
	messageToReceive := domain.MessageToReceive{}
	messageToReceive.Contents = []byte(messageReceive)
	messageToReceive.Reply = []byte(messageReply)
	messageToReceive.ReplyReliably = replyReliably
	self.MessagesToReceive = append(self.MessagesToReceive, messageToReceive)
}

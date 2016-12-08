package nodemanager

import (
"fmt"
	"strconv"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	mesh "github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/transport"
	"github.com/skycoin/skycoin/src/mesh/transport/physical"
)

type TestConfig struct {
	TransportConfig transport.TransportConfig
	UDPConfigs       []physical.UDPConfig
	NodeConfig      mesh.NodeConfig

	PeersToConnect           []Peer
	PeerToPeers		map[string]*Peer
	RoutesConfigsToEstablish []RouteConfig
	MessagesToSend           []MessageToSend
	MessagesToReceive        []MessageToReceive
	ExternalAddress		string
	Port			int
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

func (self *TestConfig) AddPeerToConnect(addr string, config *TestConfig) {
	peerToConnect := Peer{}
	peerToConnect.Peer = config.NodeConfig.PubKey
	peerToConnect.Info = physical.CreateUDPCommConfig(addr, nil)
	self.PeersToConnect = append(self.PeersToConnect, peerToConnect)
}

func (self *TestConfig) AddPeersToConnectNew(configData *ConfigData) {
	ownPubKey := self.NodeConfig.PubKey
	ownAddress := self.ExternalAddress
	for _, transportData := range(configData.Transports) {
		addrIncoming := ownAddress + ":" + strconv.Itoa(transportData.IncomingPort)
		addrOutgoing := transportData.OutgoingAddress + ":" + strconv.Itoa(transportData.OutgoingPort)

		fmt.Println(addrIncoming, addrOutgoing)

		peerToConnect := Peer{}
		peerToConnect.Peer = cipher.PubKey{}
		peerToConnect.Info = physical.CreateUDPCommConfig(addrOutgoing, nil)
		self.PeersToConnect = append(self.PeersToConnect, peerToConnect)

		ownPeer := Peer{}
		ownPeer.Peer = ownPubKey
		ownPeer.Info = physical.CreateUDPCommConfig(addrIncoming, nil)

		self.PeerToPeers[ownPeer.Info] = &peerToConnect
		fmt.Printf("Adding %s to %s\n", peerToConnect.Info, ownPeer.Info)
	}
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

func (self *TestConfig) AddMessageToSend(thruRouteID uuid.UUID, message string) {
	messageToSend := MessageToSend{}
	messageToSend.ThruRoute = thruRouteID
	messageToSend.Contents = []byte(message)
	self.MessagesToSend = append(self.MessagesToSend, messageToSend)
}

func (self *TestConfig) AddMessageToReceive(messageReceive, messageReply string) {
	messageToReceive := MessageToReceive{}
	messageToReceive.Contents = []byte(messageReceive)
	messageToReceive.Reply = []byte(messageReply)
	self.MessagesToReceive = append(self.MessagesToReceive, messageToReceive)
}

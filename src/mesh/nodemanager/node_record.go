package nodemanager

import (
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type NodeRecord struct {
	id cipher.PubKey
	nm *NodeManager

	transports map[messages.TransportId]*TransportRecord

	transportsByNodes map[cipher.PubKey]*TransportRecord

	host string
	port uint32

	ticks int

	routeForwardingRules map[messages.RouteId]*messages.RouteRule

	lock *sync.Mutex
}

func (self *NodeManager) newNode(host, hostname string) (*NodeRecord, error) {
	node := new(NodeRecord)
	id := createPubKey()
	node.id = id
	node.nm = self
	node.transports = make(map[messages.TransportId]*TransportRecord)
	node.transportsByNodes = make(map[cipher.PubKey]*TransportRecord)
	node.routeForwardingRules = make(map[messages.RouteId]*messages.RouteRule)
	node.lock = &sync.Mutex{}

	if messages.IsDebug() {
		log.Printf("Created NodeRecord %s\n", node.id.Hex())
	}

	hostData := strings.Split(host, ":")
	if len(hostData) != 2 {
		return nil, messages.ERR_INCORRECT_HOST
	}

	ipStr, portStr := hostData[0], hostData[1]
	node.host = ipStr

	ip := net.ParseIP(ipStr)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	if hostname != "" {
		err = self.dnsServer.register(id, hostname)
		if err != nil {
			return nil, err
		}
	}

	nodeAddr := &net.UDPAddr{IP: ip, Port: port}

	self.nodeList[id] = node
	self.nodeIdList = append(self.nodeIdList, id)

	self.msgServer.lock.Lock()
	self.msgServer.nodeAddrs[id] = nodeAddr
	self.msgServer.lock.Unlock()

	return node, nil
}

func (self *NodeRecord) getTransportToNode(nodeId cipher.PubKey) (*TransportRecord, error) {
	self.lock.Lock()
	transport, ok := self.transportsByNodes[nodeId]
	self.lock.Unlock()
	if !ok {
		return nil, messages.ERR_NO_TRANSPORT_TO_NODE
	}
	return transport, nil
}

func (self *NodeRecord) connectedTo(other *NodeRecord) bool {
	_, err := self.getTransportToNode(other.id)
	return err == nil
}

func (self *NodeRecord) getPeer() *messages.Peer {
	peer := &messages.Peer{self.host, self.port}
	return peer
}

func (self *NodeRecord) shutdown() {
	shutdownCM := messages.ShutdownCM{self.id}
	shutdownCMS := messages.Serialize(messages.MsgShutdownCM, shutdownCM)
	self.nm.msgServer.sendNoWait(self.id, shutdownCMS)
}

func (self *NodeRecord) setTransport(id, pairId messages.TransportId, tr *TransportRecord) error {
	pairedNodeId := tr.pair.attachedNode.id
	if self == nil {
		panic("self is nil")
	}
	self.lock.Lock()
	self.transports[id] = tr
	self.transportsByNodes[pairedNodeId] = tr
	self.lock.Unlock()

	createCM := messages.TransportCreateCM{
		Id:                id,
		PairId:            pairId,
		PairedNodeId:      pairedNodeId,
		MaxBuffer:         config.MaxBuffer,
		TimeUnit:          uint32(config.TimeUnitNum),
		TransportTimeout:  config.TransportTimeout,
		SimulateDelay:     config.SimulateDelay,
		MaxSimulatedDelay: uint32(config.MaxSimulatedDelay),
		RetransmitLimit:   uint32(config.RetransmitLimit),
	}

	createCMS := messages.Serialize(messages.MsgTransportCreateCM, createCM)
	err := self.sendToNode(createCMS)
	return err
}

func (self *NodeRecord) getTicks() int {
	ticks := 0
	for _, tr := range self.transports {
		ticks += tr.ticks
	}
	ticks += self.ticks
	return ticks
}

func (self *NodeRecord) sendToNode(msg []byte) error {
	err := self.nm.msgServer.sendMessage(self.id, msg)
	return err
}

func createPubKey() cipher.PubKey {
	pub, _ := cipher.GenerateKeyPair()
	return pub
}

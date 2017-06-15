package nodemanager

import (
	"fmt"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

//use to spawn transports
type TransportFactory struct {
	transportList []*TransportRecord
}

const (
	DISCONNECTED uint8 = 0
	CONNECTED    uint8 = 1
)

//func NewTransportFactory(network messages.Network) *TransportFactory {
func newTransportFactory() *TransportFactory {
	tf := new(TransportFactory)
	if messages.IsDebug() {
		fmt.Printf("Created Transport Factory\n")
	}
	return tf
}

func (self *TransportFactory) shutdown() {
	transports := self.transportList
	for _, tr := range transports {
		id := tr.id
		shutdownCM := messages.TransportShutdownCM{id}
		shutdownCMS := messages.Serialize(messages.MsgTransportShutdownCM, shutdownCM)
		tr.attachedNode.sendToNode(shutdownCMS)
	}
}

//move node forward on tick, process events
func (self *TransportFactory) tick() {
	//call tick on the transport
	for _, tr := range self.transportList {
		id := tr.id
		tickCM := messages.TransportTickCM{id}
		tickCMS := messages.Serialize(messages.MsgTransportTickCM, tickCM)
		tr.attachedNode.sendToNode(tickCMS)
	}
}

func (self *TransportFactory) createTransportPair() (*TransportRecord, *TransportRecord) {
	a, b := newTransportRecord(), newTransportRecord()
	a.pair, b.pair = b, a
	self.transportList = []*TransportRecord{a, b}
	return a, b
}

func (self *TransportFactory) connectNodeToNode(nodeA, nodeB *NodeRecord) error {
	if nodeA.connectedTo(nodeB) || nodeB.connectedTo(nodeA) {
		return messages.ERR_ALREADY_CONNECTED
	}

	peerA := nodeA.getPeer()
	peerB := nodeB.getPeer()

	transportA, transportB := self.createTransportPair()

	tidA := transportA.id
	tidB := transportB.id

	transportA.attachedNode = nodeA
	transportB.attachedNode = nodeB

	err := nodeA.setTransport(tidA, tidB, transportA)
	if err != nil {
		return err
	}

	err = nodeB.setTransport(tidB, tidA, transportB)
	if err != nil {
		return err
	}

	err = transportA.openUDPConn(peerA, peerB)
	if err != nil {
		return err
	}

	err = transportB.openUDPConn(peerB, peerA)
	if err != nil {
		return err
	}

	transportA.status, transportB.status = CONNECTED, CONNECTED
	return nil
}

func (self *TransportFactory) getTransports() (*TransportRecord, *TransportRecord) {
	list := self.transportList
	if len(list) < 2 {
		return nil, nil
	}
	return list[0], list[1]
}

func (self *TransportFactory) getTransportIDs() []messages.TransportId {
	list := self.transportList
	if len(list) < 2 {
		return []messages.TransportId{messages.NIL_TRANSPORT, messages.NIL_TRANSPORT}
	}
	return []messages.TransportId{list[0].id, list[1].id}
}

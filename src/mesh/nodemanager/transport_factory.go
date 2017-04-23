package nodemanager

import (
	"fmt"
	"sync"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

//use to spawn transports
type TransportFactory struct {
	//	network       messages.Network
	transportList []*TransportRecord
}

const (
	DISCONNECTED uint8 = 0
	CONNECTED    uint8 = 1
)

//func NewTransportFactory(network messages.Network) *TransportFactory {
func newTransportFactory() *TransportFactory {
	tf := new(TransportFactory)
	//	tf.network = network
	if messages.IsDebug() {
		fmt.Printf("Created Transport Factory\n")
	}
	return tf
}

func (self *TransportFactory) shutdown() {
	transports := self.transportList
	var wg sync.WaitGroup
	wg.Add(len(transports))
	for _, tr := range transports {
		//		go tr.Shutdown(&wg) //**** send shutdown command to transport owner
		id := tr.id
		shutdownCM := messages.TransportShutdownCM{id}
		shutdownCMS := messages.Serialize(messages.MsgTransportShutdownCM, shutdownCM)
		tr.attachedNode.sendToNode(shutdownCMS)
	}
	//	wg.Wait() - add waiting for ACKs!
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
	// add waiting for ACKs!
}

func (self *TransportFactory) createStubTransportPair() (*TransportRecord, *TransportRecord) {
	a, b := newTransportRecord(), newTransportRecord()
	a.pair, b.pair = b, a
	self.transportList = []*TransportRecord{a, b}
	return a, b
}

func (self *TransportFactory) connectNodeToNode(nodeA, nodeB *NodeRecord) error {
	if nodeA.ConnectedTo(nodeB) || nodeB.ConnectedTo(nodeA) {
		return messages.ERR_ALREADY_CONNECTED
	} //**** we need only to check for paired nodes

	peerA := nodeA.GetPeer()
	peerB := nodeB.GetPeer()

	transportA, transportB := self.createStubTransportPair()

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

	err = transportA.openUDPConn(peerA, peerB) //**** send command to open udp connection and wait for an answer
	if err != nil {
		return err
	}

	err = transportB.openUDPConn(peerB, peerA) //**** send command to open udp connection and wait for an answer
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

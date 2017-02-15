package transport

import (
	"fmt"
	"sync"

	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

//use to spawn transports
type TransportFactory struct {
	TransportList []*Transport
}

func NewTransportFactory() *TransportFactory {
	tf := new(TransportFactory)
	if messages.IsDebug() {
		fmt.Printf("Created Transport Factory\n")
	}
	return tf
}

func (self *TransportFactory) Shutdown() {
	transports := self.TransportList
	var wg sync.WaitGroup
	wg.Add(len(transports))
	for _, tr := range transports {
		go tr.Shutdown(&wg)
	}
	wg.Wait()
}

//move node forward on tick, process events
func (self *TransportFactory) Tick() {
	//call tick on the transport
	for _, t := range self.TransportList {
		t.Tick()
	}
	/*
		//If this does not work
		//- then force transports to push to a transport factory incoming channel
		for _, t := range self.TransportList {
			//check each transport for data?
			for len(t.pendingOut) > 0 { //len will expose the number of elements in the channels buffer
				var b messages.TransportDatagramTransfer
				t.pendingOut <- b                 //the channel data
				t.sendMessageToStubPair([]byte{}) //the transport now has the data
			}
		}
	*/
}

//implement/fix
func (self *TransportFactory) createStubTransportPair() (*Transport, *Transport) {
	a, b := newTransportStub(), newTransportStub()
	a.StubPair, b.StubPair = b, a
	self.TransportList = []*Transport{a, b}
	return a, b
}

func (self *TransportFactory) connectPeers(peerA, peerB *messages.Peer) (*Transport, *Transport, error) {
	transportA, transportB := self.createStubTransportPair()

	err := transportA.openConn(peerA, peerB)
	if err != nil {
		return nil, nil, err
	}
	err = transportB.openConn(peerB, peerA)
	if err != nil {
		panic(err)
		return nil, nil, err
	}
	transportA.Status, transportB.Status = CONNECTED, CONNECTED
	return transportA, transportB, nil
}

func (self *TransportFactory) ConnectNodeToNode(nodeA, nodeB messages.NodeInterface) error {
	if nodeA.ConnectedTo(nodeB) || nodeB.ConnectedTo(nodeA) {
		return errors.ERR_ALREADY_CONNECTED
	}

	peerA := nodeA.GetPeer()
	peerB := nodeB.GetPeer()

	transportA, transportB, err := self.connectPeers(peerA, peerB)
	if err != nil {
		return err
	}

	transportA.AttachedNode = nodeA
	tidA := transportA.Id
	transportB.AttachedNode = nodeB
	tidB := transportB.Id

	nodeA.SetTransport(tidA, transportA)
	nodeB.SetTransport(tidB, transportB)

	return nil
}

func (self *TransportFactory) GetTransports() (*Transport, *Transport) {
	list := self.TransportList
	if len(list) < 2 {
		return nil, nil
	}
	return list[0], list[1]
}

func (self *TransportFactory) GetTransportIDs() []messages.TransportId {
	list := self.TransportList
	if len(list) < 2 {
		return []messages.TransportId{messages.NIL_TRANSPORT, messages.NIL_TRANSPORT}
	}
	return []messages.TransportId{list[0].Id, list[1].Id}
}

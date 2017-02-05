package transport

import (
	//"fmt"

	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

//use to spawn transports
type TransportFactory struct {
	TransportList []*Transport
}

func NewTransportFactory() *TransportFactory {
	tf := new(TransportFactory)
	//fmt.Printf("Created Transport Factory\n")
	return tf
}

func (self *TransportFactory) Shutdown() {
	for _, tr := range self.TransportList {
		tr.Shutdown()
	}
}

//move node forward on tick, process events
func (self *TransportFactory) Tick() {
	//call tick on the transport
	for _, t := range self.TransportList {
		t.Tick()
	}

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
}

//implement/fix
//Implement the nodes the transports are attached to
func (self *TransportFactory) CreateStubTransportPair() (*Transport, *Transport) {
	a, b := &Transport{}, &Transport{}
	a.newTransportStub()
	b.newTransportStub()
	a.StubPair, b.StubPair = b, a
	a.Status, b.Status = CONNECTED, CONNECTED
	self.TransportList = []*Transport{a, b}
	return a, b
}

func (self *TransportFactory) ConnectNodeToNode(nodeA, nodeB messages.NodeInterface) error {
	if nodeA.ConnectedTo(nodeB) || nodeB.ConnectedTo(nodeA) {
		return errors.ERR_ALREADY_CONNECTED
	}
	transportA, transportB := self.CreateStubTransportPair()
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

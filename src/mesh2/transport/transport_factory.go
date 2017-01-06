package transport

import (
	"fmt"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

//use to spawn transports
type TransportFactory struct {
	TransportList []*Transport
}

func NewTransportFactory() *TransportFactory {
	tf := new(TransportFactory)
	fmt.Printf("Created Transport Factory\n")
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
		for len(t.PendingOut) > 0 { //len will expose the number of elements in the channels buffer
			var b []byte
			t.PendingOut <- b          //the channel data
			t.SendMessageToStubPair(b) //the transport now has the data
		}
	}
}

//implement/fix
//Implement the nodes the transports are attached to
func (self *TransportFactory) CreateStubTransportPair() (*Transport, *Transport) {
	a, b := &Transport{}, &Transport{}
	a.NewTransportStub()
	b.NewTransportStub()
	a.StubPair, b.StubPair = b, a
	a.Status, b.Status = CONNECTED, CONNECTED
	self.TransportList = []*Transport{a, b}
	return a, b
}

func (self *TransportFactory) ConnectNodeToNode(nodeA, nodeB messages.NodeInterface) {
	transportA, transportB := self.CreateStubTransportPair()
	transportA.AttachedNode = nodeA
	tidA := transportA.Id
	transportB.AttachedNode = nodeB
	tidB := transportB.Id
	nodeA.SetTransport(tidA, transportA)
	nodeB.SetTransport(tidB, transportB)
}

func (self *TransportFactory) GetTransports() (*Transport, *Transport) {
	list := self.TransportList
	if len(list) < 2 {
		return nil, nil
	}
	return list[0], list[1]
}

package transport

import (
	"fmt"
)

//use to spawn transports
type TransportFactory struct {
	TransportList []Transport
}

func NewTransportFactory() *TransportFactory {
	tf := new(TransportFactory)
	fmt.Printf("Created Transport Factory")
	return tf
}

func (self *TransportFactory) Shutdown() {
	//close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *TransportFactory) Tick() {
	//call tick on the transport
	for _, t := range self.TransportList {
		t.Tick()
	}

	for _, t := range self.TransportList {
		t.Tick()
	}

	//If this does not work
	//- then force transports to push to a transport factory incoming channel
	for _, t := range self.TransportList {
		//check each transport for data?
		for len(t.PendingOut) > 0 {
			var b []byte
			t.PendingOut <- b          //the channel data
			t.SendMessageToStubPair(b) //the transport now has the data
		}
	}
}

//implement/fix
func (self *TransportFactory) CreateStubTransportPair() (Transport, Transport) {
	var a Transport
	var b Transport
	a.NewTransportStub()
	b.NewTransportStub()
	a.StubPair = &b
	b.StubPair = &a
	return a, b
}

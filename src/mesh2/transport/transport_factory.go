package transport

import (
	"fmt"
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
	//close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *TransportFactory) Tick() {
	/*	fmt.Println("ticking  tf")
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
			for len(t.PendingOut) > 0 { //len will expose the number of elements in the channels buffer
				var b []byte
				t.PendingOut <- b          //the channel data
				t.SendMessageToStubPair(b) //the transport now has the data
			}
		}*/
}

//implement/fix
//Implement the nodes the transports are attached to
func (self *TransportFactory) CreateStubTransportPair() (*Transport, *Transport) {
	a, b := &Transport{}, &Transport{}
	a.NewTransportStub()
	b.NewTransportStub()
	a.StubPair, b.StubPair = b, a
	a.Status, b.Status = CONNECTED, CONNECTED
	self.TransportList = append(self.TransportList, a)
	self.TransportList = append(self.TransportList, b)
	return a, b
}

func (self *TransportFactory) GetTransports() (*Transport, *Transport) {
	list := self.TransportList
	if len(list) < 2 {
		return nil, nil
	}
	return list[0], list[1]
}

package transport

import ()

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
	close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *TransportFactory) Tick() {
	//call tick on the transport
	for t := range self.TransportList {
		t.Tick()
	}
}

func (self *TransportFactory) CreateStubTransportPair() (Transport, Transport) {

}

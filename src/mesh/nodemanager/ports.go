package nodemanager

import (
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type PortDelivery struct {
	content        map[string]uint32
	requestChannel chan Request
	startPort      uint32
}

type Request struct {
	host            string
	responseChannel chan uint32
}

func newPortDelivery() *PortDelivery {
	portDelivery := &PortDelivery{}
	portDelivery.content = map[string]uint32{}
	portDelivery.requestChannel = make(chan Request, 1024)
	portDelivery.startPort = messages.GetConfig().StartPort
	go portDelivery.deliver()
	return portDelivery
}

func (self *PortDelivery) Get(host string) uint32 {
	responseChannel := make(chan uint32, 1024)
	request := Request{host, responseChannel}
	self.requestChannel <- request
	response := <-responseChannel
	return response
}

func (self *PortDelivery) deliver() uint32 {
	for {
		select {
		case request := <-self.requestChannel:
			host := request.host
			port := self.content[host]
			if port == 0 {
				port = self.startPort
			}
			self.content[host] = port + 1
			request.responseChannel <- port
		}
	}
}

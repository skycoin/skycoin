package nodemanager

import (
	"sync"
)

type PortDelivery struct {
	port      map[string]uint32
	lock      *sync.Mutex
	startPort uint32
}

func newPortDelivery() *PortDelivery {
	portDelivery := &PortDelivery{}
	portDelivery.port = map[string]uint32{}
	portDelivery.startPort = config.StartPort
	portDelivery.lock = &sync.Mutex{}
	return portDelivery
}

func (self *PortDelivery) get(host string) uint32 {
	self.lock.Lock()
	defer self.lock.Unlock()
	port := self.port[host]
	if port == 0 {
		port = self.startPort
	}
	self.port[host] = port + 1
	return port
}

package node

import (
	"github.com/skycoin/skycoin/src/mesh/errors"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

func (self *Node) getTransport(transportId messages.TransportId) (*transport.Transport, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	tr, ok := self.Transports[transportId]
	if !ok {
		return nil, errors.ERR_TRANSPORT_DOESNT_EXIST
	}
	return tr, nil
}

func (self *Node) setTransport(transportId messages.TransportId, tr *transport.Transport) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.Transports[transportId] = tr
}

func (self *Node) removeTransport(transportId messages.TransportId) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.Transports[transportId]; !ok {
		return errors.ERR_TRANSPORT_DOESNT_EXIST
	}

	delete(self.Transports, transportId)
	return nil
}

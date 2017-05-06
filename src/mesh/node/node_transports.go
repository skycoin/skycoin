package node

import (
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

func (self *Node) getTransport(transportId messages.TransportId) (*transport.Transport, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	tr, ok := self.transports[transportId]
	if !ok {
		return nil, messages.ERR_TRANSPORT_DOESNT_EXIST
	}
	return tr, nil
}

func (self *Node) setTransportFromMessage(msg *messages.TransportCreateCM) {
	tr := transport.CreateTransportFromMessage(msg)
	tr.AttachedNode = self

	self.lock.Lock()
	defer self.lock.Unlock()

	self.transports[tr.Id()] = tr
	self.transportsByNodes[msg.PairedNodeId] = tr
}

package nodemanager

import (
	"log"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type TransportRecord struct {
	id           messages.TransportId
	pair         *TransportRecord //this is the other transport pair
	ticks        int
	attachedNode *NodeRecord //node the transport is attached to
	status       uint8
}

//are created by the factories
func newTransportRecord() *TransportRecord {
	tr := TransportRecord{}
	tr.id = messages.RandTransportId()
	tr.status = DISCONNECTED
	if messages.IsDebug() {
		log.Printf("Created TransportRecord: %d\n", tr.id)
	}
	return &tr
}

func (self *TransportRecord) openUDPConn(peerA, peerB *messages.Peer) error {
	openMsg := messages.OpenUDPCM{
		self.id,
		*peerA,
		*peerB,
	}
	openMsgS := messages.Serialize(messages.MsgOpenUDPCM, openMsg)
	err := self.attachedNode.sendToNode(openMsgS)
	return err
}

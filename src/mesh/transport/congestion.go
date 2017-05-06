package transport

import (
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

func (self *Transport) broadcastCongestion() {
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		go self.sendCongestion()
		<-ticker.C
	}
}

func (self *Transport) sendCongestion() {
	if self.AttachedNode != nil {
		packetToNode := messages.CongestionPacket{
			self.id,
			self.nodeCongestion,
		}

		go self.AttachedNode.InjectCongestionPacket(&packetToNode)
	}
}

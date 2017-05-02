package nodemanager

import (
	"log"

	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/viscript"
)

type NMViscriptServer struct {
	viscript.ViscriptServer
	nm *NodeManager
}

func (self *NodeManager) TalkToViscript(sequence, appId uint32) {
	vs := &NMViscriptServer{nm: self}
	vs.Init(sequence, appId)
}

func (self *NMViscriptServer) handleUserCommand(uc *messages.UserCommand) {
	log.Println("command received:", uc)
	sequence := uc.Sequence
	appId := uc.AppId
	message := uc.Payload

	switch messages.GetMessageType(message) {

	case messages.MsgPing:
		ack := &messages.PingAck{}
		ackS := messages.Serialize(messages.MsgPingAck, ack)
		self.SendAck(ackS, sequence, appId)

	case messages.MsgResourceUsage:
		cpu, memory, err := self.GetResources()
		if err == nil {
			ack := &messages.ResourceUsageAck{
				cpu,
				memory,
			}
			ackS := messages.Serialize(messages.MsgResourceUsageAck, ack)
			self.SendAck(ackS, sequence, appId)
		}

	case messages.MsgUserShutdown:
		self.nm.Shutdown()
		ack := &messages.UserShutdownAck{}
		ackS := messages.Serialize(messages.MsgUserShutdownAck, ack)
		self.SendAck(ackS, sequence, appId)
		panic("goodbye")

	default:
		log.Println("Unknown user command:", message)
	}
}

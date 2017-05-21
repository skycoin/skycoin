package nodemanager

import (
	"log"

	vsmsg "github.com/corpusc/viscript/msg"

	"github.com/skycoin/skycoin/src/mesh/viscript"
)

type NMViscriptServer struct {
	viscript.ViscriptServer
	nm *NodeManager
}

func (self *NodeManager) TalkToViscript(sequence, appId uint32) {
	vs := &NMViscriptServer{nm: self}
	self.viscriptServer = vs
	vs.Init(sequence, appId)
}

func (self *NMViscriptServer) handleUserCommand(uc *vsmsg.MessageUserCommand) {
	log.Println("command received:", uc)
	sequence := uc.Sequence
	appId := uc.AppId
	message := uc.Payload

	switch vsmsg.GetType(message) {

	case vsmsg.TypePing:
		ack := &vsmsg.MessagePingAck{}
		ackS := vsmsg.Serialize(vsmsg.TypePingAck, ack)
		self.SendAck(ackS, sequence, appId)

	case vsmsg.TypeResourceUsage:
		cpu, memory, err := self.GetResources()
		if err == nil {
			ack := &vsmsg.MessageResourceUsageAck{
				cpu,
				memory,
			}
			ackS := vsmsg.Serialize(vsmsg.TypeResourceUsageAck, ack)
			self.SendAck(ackS, sequence, appId)
		}

	case vsmsg.TypeShutdown:
		self.nm.Shutdown()
		ack := &vsmsg.MessageShutdownAck{}
		ackS := vsmsg.Serialize(vsmsg.TypeShutdownAck, ack)
		self.SendAck(ackS, sequence, appId)
		panic("goodbye")

	default:
		log.Println("Unknown user command:", message)
	}
}

package app

import (
	"log"

	vsmsg "github.com/corpusc/viscript/msg"

	"github.com/skycoin/skycoin/src/mesh/viscript"
)

type AppViscriptServer struct {
	viscript.ViscriptServer
	app *app
}

func (self *app) TalkToViscript(sequence, appId uint32) {
	vs := &AppViscriptServer{app: self}
	self.viscriptServer = vs
	vs.Init(sequence, appId)
}

func (self *AppViscriptServer) handleUserCommand(uc *vsmsg.MessageUserCommand) {
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
		self.app.Shutdown()
		ack := &vsmsg.MessageShutdownAck{}
		ackS := vsmsg.Serialize(vsmsg.TypeShutdownAck, ack)
		self.SendAck(ackS, sequence, appId)
		panic("goodbye")

	default:
		log.Println("Unknown user command:", message)
	}
}

package node

import (
	"log"

	vsmsg "github.com/corpusc/viscript/msg"

	"github.com/skycoin/skycoin/src/mesh/viscript"
)

type NodeViscriptServer struct {
	viscript.ViscriptServer
	node *Node
}

func (self *Node) TalkToViscript(sequence, appId uint32) {
	vs := &NodeViscriptServer{node: self}
	self.viscriptServer = vs
	vs.Init(sequence, appId)
}

func (self *NodeViscriptServer) handleUserCommand(uc *vsmsg.MessageUserCommand) {
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

	case vsmsg.TypeConnectDirectly:
		ucd := &vsmsg.MessageConnectDirectly{}
		err := vsmsg.Deserialize(message, ucd)
		if err == nil {
			address := ucd.Address
			err = self.node.ConnectDirectly(address)
			if err != nil {
				ack := &vsmsg.MessageConnectDirectlyAck{}
				ackS := vsmsg.Serialize(vsmsg.TypeConnectDirectlyAck, ack)
				self.SendAck(ackS, sequence, appId)
			}
		}

	case vsmsg.TypeShutdown:
		self.node.Shutdown()
		ack := &vsmsg.MessageShutdownAck{}
		ackS := vsmsg.Serialize(vsmsg.TypeShutdownAck, ack)
		self.SendAck(ackS, sequence, appId)
		panic("goodbye")

	default:
		log.Println("Unknown user command:", message)
	}
}

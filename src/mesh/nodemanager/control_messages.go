package nodemanager

import (
	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

func (self *NodeManager) handleControlMessage(cm *messages.InControlMessage) {
	sequence := cm.Sequence
	msg := cm.PayloadMessage

	switch messages.GetMessageType(msg) {
	case messages.MsgConnectCM:
		m1 := messages.ConnectCM{}
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			log.Println(err)
			return
		}
		from := m1.From
		to := m1.To
		connSequence := m1.Sequence
		appIdFrom := m1.AppIdFrom
		appIdTo := m1.AppIdTo
		connId, err := self.connect(from, to, appIdFrom, appIdTo)
		if err != nil {
			log.Println(err)
			self.sendConnectAck(from, sequence, connSequence, false, messages.ConnectionId(0))
			return
		}
		self.sendConnectAck(from, sequence, connSequence, true, connId)

	case messages.MsgRegisterNodeCM:
		m1 := messages.RegisterNodeCM{}
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			log.Println(err)
			return
		}
		host := m1.Host
		connect := m1.Connect
		var nodeId cipher.PubKey
		if !connect {
			id, err := self.addNewNode(host)
			if err == nil {
				nodeId = id
			} else {
				return
			}
		} else {
			id, err := self.addAndConnect(host)
			if err == nil {
				nodeId = id
			} else {
				return
			}
		}
		self.sendRegisterAck(sequence, nodeId)

	case messages.MsgCommonCMAck:
		m1 := messages.CommonCMAck{}
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			log.Println(err)
			return
		}
		self.msgServer.getResponse(sequence, &m1)
	}
}

func (self *NodeManager) sendRegisterAck(sequence uint32, nodeId cipher.PubKey) {
	ack := messages.RegisterNodeCMAck{
		Ok:                true,
		NodeId:            nodeId,
		MaxBuffer:         config.MaxBuffer,
		MaxPacketSize:     uint32(config.MaxPacketSize),
		TimeUnit:          uint32(config.TimeUnitNum),
		SendInterval:      config.SendIntervalNum,
		ConnectionTimeout: config.ConnectionTimeout,
	}
	ackS := messages.Serialize(messages.MsgRegisterNodeCMAck, ack)
	self.msgServer.sendAck(sequence, nodeId, ackS)
}

func (self *NodeManager) sendConnectAck(nodeId cipher.PubKey, sequence, connSequence uint32, ok bool, connectionId messages.ConnectionId) {
	ack := messages.ConnectCMAck{
		Sequence:     connSequence,
		Ok:           ok,
		ConnectionId: connectionId,
	}
	ackS := messages.Serialize(messages.MsgConnectCMAck, ack)
	self.msgServer.sendAck(sequence, nodeId, ackS)
}

func (self *NodeManager) sendTrueCommonAck(sequence uint32, nodeId cipher.PubKey) {
	ack := &messages.CommonCMAck{
		Ok: true,
	}
	self.sendCommonAck(sequence, nodeId, ack)
}

func (self *NodeManager) sendFalseCommonAck(sequence uint32, nodeId cipher.PubKey) {
	ack := &messages.CommonCMAck{
		Ok: false,
	}
	self.sendCommonAck(sequence, nodeId, ack)
}

func (self *NodeManager) sendCommonAck(sequence uint32, nodeId cipher.PubKey, ack *messages.CommonCMAck) {
	ackS := messages.Serialize(messages.MsgCommonCMAck, ack)
	self.msgServer.sendAck(sequence, nodeId, ackS)
}

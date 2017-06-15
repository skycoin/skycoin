package node

import (
	"github.com/skycoin/skycoin/src/mesh/messages"
	"log"
)

type ControlChannel struct {
	id messages.ChannelId
}

func newControlChannel() *ControlChannel {
	c := ControlChannel{
		id: messages.RandChannelId(),
	}
	return &c
}

func (c *ControlChannel) handleMessage(handledNode *Node, sequence uint32, msg []byte) error {

	switch messages.GetMessageType(msg) {
	case messages.MsgAddRouteCM:
		var m1 messages.AddRouteCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		routeRule := messages.RouteRule{
			m1.IncomingTransportId,
			m1.OutgoingTransportId,
			m1.IncomingRouteId,
			m1.OutgoingRouteId,
		}
		err = handledNode.addRoute(&routeRule)
		if err == nil {
			handledNode.sendTrueAckToServer(sequence)
		} else {
			handledNode.sendFalseAckToServer(sequence)
		}
		return err

	case messages.MsgRemoveRouteCM:
		var m1 messages.RemoveRouteCM
		messages.Deserialize(msg, &m1)
		routeId := m1.RouteId
		err := handledNode.removeRoute(routeId)
		if err == nil {
			handledNode.sendTrueAckToServer(sequence)
		} else {
			handledNode.sendFalseAckToServer(sequence)
		}
		return err

	case messages.MsgAssignPortCM:
		var m1 messages.AssignPortCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		handledNode.port = m1.Port
		handledNode.sendTrueAckToServer(sequence)
		return nil

	case messages.MsgTransportCreateCM:
		var m1 messages.TransportCreateCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		handledNode.setTransportFromMessage(&m1)
		handledNode.sendTrueAckToServer(sequence)
		return nil

	case messages.MsgTransportTickCM:
		var m1 messages.TransportTickCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		trId := m1.Id
		tr, err := handledNode.getTransport(trId)
		if err == nil {
			tr.Tick()
			handledNode.sendTrueAckToServer(sequence)
		} else {
			handledNode.sendFalseAckToServer(sequence)
		}
		return nil

	case messages.MsgTransportShutdownCM:
		var m1 messages.TransportShutdownCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		trId := m1.Id
		tr, err := handledNode.getTransport(trId)
		if err != nil {
			return err
		}
		if err == nil {
			tr.Shutdown()
			handledNode.sendTrueAckToServer(sequence)
		} else {
			handledNode.sendFalseAckToServer(sequence)
		}
		return err

	case messages.MsgOpenUDPCM:
		var m1 messages.OpenUDPCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			panic(err)
		}
		trId := m1.Id
		tr, err := handledNode.getTransport(trId)
		if err != nil {
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		err = tr.OpenUDPConn(&m1.PeerA, &m1.PeerB)
		if err == nil {
			handledNode.sendTrueAckToServer(sequence)
		} else {
			handledNode.sendFalseAckToServer(sequence)
		}
		return err

	case messages.MsgAssignConnectionCM:
		var m1 messages.AssignConnectionCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		handledNode.sendTrueAckToServer(sequence)
		routeId := m1.RouteId
		connId := m1.ConnectionId
		appId := m1.AppId
		_, err = handledNode.newConnection(connId, routeId, appId)
		return err

	case messages.MsgConnectionOnCM:
		var m1 messages.ConnectionOnCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			log.Println("deserialization failed")
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		nodeId := m1.NodeId
		if nodeId != handledNode.id {
			log.Println("wrong id")
			handledNode.sendFalseAckToServer(sequence)
			return err
		}
		connId := m1.ConnectionId
		handledNode.sendTrueAckToServer(sequence)
		handledNode.setConnectionOn(connId)
		return nil

	case messages.MsgRegisterNodeCMAck:
		var ack messages.RegisterNodeCMAck
		err := messages.Deserialize(msg, &ack)
		if err != nil {
			return err
		}

		responseChannel, ok := handledNode.getResponseChannel(sequence)
		if ok {
			handledNode.register(&ack)
			responseChannel <- ack.Ok
		}
		return nil

	case messages.MsgConnectDirectlyCMAck:
		var ack messages.ConnectDirectlyCMAck
		err := messages.Deserialize(msg, &ack)
		if err != nil {
			return err
		}

		responseChannel, ok := handledNode.getResponseChannel(sequence)
		if ok {
			responseChannel <- ack.Ok
			connectResponseChannel, ok0 := handledNode.connectResponseChannels[ack.Sequence]
			if ok0 {
				connectResponseChannel <- true
			}
		}
		return nil

	case messages.MsgConnectWithRouteCMAck:
		var ack messages.ConnectWithRouteCMAck
		err := messages.Deserialize(msg, &ack)
		if err != nil {
			return err
		}

		responseChannel, ok := handledNode.getResponseChannel(sequence)
		if ok {
			responseChannel <- ack.Ok
			connectionResponseChannel, ok0 := handledNode.connectionResponseChannels[ack.Sequence]
			if ok0 {
				connectionResponseChannel <- ack.ConnectionId
			}
		}
		return nil

	case messages.MsgShutdownCM:
		var m1 messages.ShutdownCM
		err := messages.Deserialize(msg, &m1)
		if err != nil {
			return err
		}

		if m1.NodeId == handledNode.id {
			handledNode.Shutdown()
		}

		return nil

	case messages.MsgCommonCMAck:
		var ack messages.CommonCMAck
		err := messages.Deserialize(msg, &ack)
		if err != nil {
			return err
		}

		responseChannel, ok := handledNode.getResponseChannel(sequence)
		if ok {
			responseChannel <- ack.Ok
		}
		return nil

	default:
		log.Println("Incorrect message type:", msg)
	}

	handledNode.sendFalseAckToServer(sequence)
	return messages.ERR_UNKNOWN_MESSAGE_TYPE
}

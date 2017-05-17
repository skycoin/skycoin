package node

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

func (self *Node) addControlChannel() messages.ChannelId {

	channel := newControlChannel()

	self.lock.Lock()
	defer self.lock.Unlock()

	self.controlChannels[channel.id] = channel
	return channel.id
}

func (self *Node) addZeroControlChannel() {

	channel := newControlChannel()
	channel.id = messages.ChannelId(0)

	self.lock.Lock()
	defer self.lock.Unlock()

	self.controlChannels[channel.id] = channel
	return
}

func (self *Node) closeControlChannel(channelID messages.ChannelId) error {

	if _, ok := self.controlChannels[channelID]; !ok {
		return errors.New("Control channel not found")
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	delete(self.controlChannels, channelID)
	return nil
}

func (self *Node) handleControlMessage(_ messages.ChannelId, inControlMsg *messages.InControlMessage) error {

	channelID := messages.ChannelId(0)

	self.lock.Lock()
	channel, ok := self.controlChannels[channelID]
	self.lock.Unlock()
	if !ok {
		return errors.New("Control channel not found")
	}

	sequence := inControlMsg.Sequence
	msg := inControlMsg.PayloadMessage
	err := channel.handleMessage(self, sequence, msg)
	return err
}

func (self *Node) sendTrueAckToServer(sequence uint32) error {
	ack := &messages.CommonCMAck{true}
	return self.sendAckToServer(sequence, ack)
}

func (self *Node) sendFalseAckToServer(sequence uint32) error {
	ack := &messages.CommonCMAck{false}
	return self.sendAckToServer(sequence, ack)
}

func (self *Node) sendAckToServer(sequence uint32, ack *messages.CommonCMAck) error {
	ackS := messages.Serialize(messages.MsgCommonCMAck, *ack)
	return self.sendToServer(sequence, ackS)
}

func (self *Node) sendRegisterNodeToServer(hostname, host string, connect bool) error {
	msg := messages.RegisterNodeCM{hostname, host, connect}
	msgS := messages.Serialize(messages.MsgRegisterNodeCM, msg)
	err := self.sendMessageToServer(msgS)
	return err
}

func (self *Node) sendConnectDirectlyToServer(nodeToId string) error {
	responseChannel := make(chan bool)

	self.lock.Lock()
	connectSequence := self.connectResponseSequence
	self.connectResponseSequence++
	self.connectResponseChannels[connectSequence] = responseChannel
	self.lock.Unlock()

	msg := messages.ConnectDirectlyCM{connectSequence, self.id, nodeToId}
	msgS := messages.Serialize(messages.MsgConnectDirectlyCM, msg)

	err := self.sendMessageToServer(msgS)
	if err != nil {
		return err
	}

	select {
	case <-responseChannel:
		return nil
	case <-time.After(CONTROL_TIMEOUT):
		return messages.ERR_MSG_SRV_TIMEOUT
	}
}

func (self *Node) sendConnectWithRouteToServer(nodeToId string, appIdFrom, appIdTo messages.AppId) (messages.ConnectionId, error) {
	responseChannel := make(chan messages.ConnectionId)

	self.lock.Lock()
	connSequence := self.connectionResponseSequence
	self.connectionResponseSequence++
	self.connectionResponseChannels[connSequence] = responseChannel
	self.lock.Unlock()

	msg := messages.ConnectWithRouteCM{connSequence, appIdFrom, appIdTo, self.id, nodeToId}
	msgS := messages.Serialize(messages.MsgConnectWithRouteCM, msg)

	err := self.sendMessageToServer(msgS)
	if err != nil {
		return messages.ConnectionId(0), err
	}

	select {
	case connId := <-responseChannel:
		return connId, nil
	case <-time.After(CONTROL_TIMEOUT):
		return messages.ConnectionId(0), messages.ERR_MSG_SRV_TIMEOUT
	}
}

func (self *Node) sendMessageToServer(msg []byte) error {
	sequence := self.sequence
	self.sequence++

	responseChannel := make(chan bool)
	self.setResponseChannel(sequence, responseChannel)

	err := self.sendToServer(sequence, msg)
	if err != nil {
		return err
	}

	select {
	case ok := <-responseChannel:
		if ok {
			return nil
		} else {
			return messages.ERR_REGISTER_NODE_FAILED
		}
	case <-time.After(CONTROL_TIMEOUT):
		return messages.ERR_MSG_SRV_TIMEOUT
	}
}

func (self *Node) sendToServer(sequence uint32, msg []byte) error {
	if len(self.serverAddrs) == 0 {
		return nil
	}

	inControl := messages.InControlMessage{
		messages.ChannelId(0),
		sequence,
		msg,
	}
	inControlS := messages.Serialize(messages.MsgInControlMessage, inControl)
	_, err := self.controlConn.WriteTo(inControlS, self.serverAddrs[0])
	return err
}

func (self *Node) openUDPforCM(port int) (*net.UDPConn, error) {
	host := net.ParseIP(messages.LOCALHOST)
	connAddr := &net.UDPAddr{IP: host, Port: port}

	conn, err := net.ListenUDP("udp", connAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (self *Node) addServer(serverAddrStr string) {

	nmData := strings.Split(serverAddrStr, ":")
	nmHostStr := nmData[0]
	nmPort := 5999
	if len(nmData) > 1 {
		port, err := strconv.Atoi(nmData[1])
		if err == nil {
			nmPort = port
		}
	}
	nmHost := net.ParseIP(nmHostStr)
	serverAddr := &net.UDPAddr{IP: nmHost, Port: nmPort}
	self.serverAddrs = append(self.serverAddrs, serverAddr)
}

func (self *Node) receiveControlMessages() {
	go_on := true
	go func() {
		for go_on {

			buffer := make([]byte, 1024)

			n, _, err := self.controlConn.ReadFrom(buffer)

			if err != nil {
				break
			} else {
				if n == 0 {
					continue
				}
				cm := messages.InControlMessage{}
				err := messages.Deserialize(buffer[:n], &cm)
				if err != nil {
					log.Println("Incorrect InControlMessage:", buffer[:n])
					continue
				}
				go self.injectControlMessage(&cm)
			}
		}
	}()
	<-self.closeControlMessagesChannel
	go_on = false
}

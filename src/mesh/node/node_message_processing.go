package mesh

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/mesh/domain"
)

// Message order is not preserved
func (self *Node) SetReceiveChannel(received chan domain.MeshMessage) {
	self.outputMessagesReceived = received
}

func (self *Node) processIncomingMessagesLoop() {
	self.closeGroup.Add(1)
	defer self.closeGroup.Done()
	for len(self.closing) == 0 {
		select {
		case msg, ok := <-self.transportsMessagesReceived:
			{
				if ok {
					self.processMessage(msg)
				}
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Node) processMessage(serialized []byte) {
	message, err := self.serializer.UnserializeMessage(serialized)
	fmt.Fprintf(os.Stdout, "Peer %d: Incoming message!\n", self.Config.PubKey[:3])
	fmt.Fprintf(os.Stdout, "Incoming message %v\n", serialized)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Deserialization error %v %v\n", serialized, err)
		return
	}

	messageType := reflect.TypeOf(message)

	if messageType == reflect.TypeOf(domain.SetControlChannelMessage{}) {
		self.processSetControlChannelMessage(message.(domain.SetControlChannelMessage))
		return
	}
	if messageType == reflect.TypeOf(domain.SetControlChannelResponseMessage{}) {
		self.processSetControlChannelResponseMessage(message.(domain.SetControlChannelResponseMessage))
		return
	}
	if messageType == reflect.TypeOf(domain.SetRouteControlMessage{}) {
		self.processSetRouteControlMessage(message.(domain.SetRouteControlMessage))
		return
	}
	if messageType == reflect.TypeOf(domain.ResponseMessage{}) {
		self.processResponseControlMessage(message.(domain.ResponseMessage))
		return
	}

	if messageType == reflect.TypeOf(domain.UserMessage{}) {
		self.processUserMessage(message.(domain.UserMessage))
		return
	}

	forwardedMessage := self.forwardMessage(message)
	if !forwardedMessage {
		if messageType == reflect.TypeOf(domain.SetRouteMessage{}) {
			self.processSetRouteMessage(message.(domain.SetRouteMessage))
			return
		} else if messageType == reflect.TypeOf(domain.SetRouteReply{}) {
			self.processSetRouteReplyMessage(message.(domain.SetRouteReply))
			return
		}
	} else {
		if messageType == reflect.TypeOf(domain.DeleteRouteMessage{}) {
			self.processDeleteRouteMessage(message.(domain.DeleteRouteMessage))
			return
		}
	}

	if messageType == reflect.TypeOf(domain.RefreshRouteMessage{}) {
		self.processRefreshRouteMessage(message.(domain.RefreshRouteMessage), forwardedMessage)
	}
}

func (self *Node) processResponseControlMessage(message domain.ResponseMessage) {
	fmt.Fprintf(os.Stdout, "Peer %d: Processing 'ResponseMessage' message from peer %d \n", self.Config.PubKey[:3], message.FromPeerID[:3])
	result := []byte{0}
	if message.Result {
		result = []byte{1}
	}
	self.outputMessagesReceived <- domain.MeshMessage{
		Contents: result,
	}
}

func (self *Node) processSetRouteControlMessage(message domain.SetRouteControlMessage) {
	fmt.Fprintf(os.Stdout, "Peer %d: Processing 'SetRouteMessage' message from peer %d \n", self.Config.PubKey[:3], message.FromPeerID[:3])

	err := self.HandleControlMessage(message.ChannelID, message)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Can't handle SetRoute message %s \n", err)
		return
	}

	replyMessage := domain.ResponseMessage{
		FromPeerID: self.Config.PubKey,
		RequestID:  message.RequestID,
		Result:     true,
	}

	transportToPeer := self.safelyGetTransportToPeer(message.FromPeerID)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v from %v, dropping\n", message.FromPeerID, self.Config.PubKey)
		return
	}

	serialized := self.serializer.SerializeMessage(replyMessage)
	err = transportToPeer.SendMessage(message.FromPeerID, serialized, nil)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stdout, "Peer %d: Sending 'ResponseMessage' message to peer %d! \n", self.Config.PubKey[:3], message.FromPeerID[:3])
}

func (self *Node) processSetControlChannelResponseMessage(message domain.SetControlChannelResponseMessage) {
	fmt.Fprintf(os.Stdout, "Peer %d: Processing 'SetControlChannelResponse' message from peer %d \n", self.Config.PubKey[:3], message.FromPeerID[:3])
	fmt.Fprintf(os.Stdout, "Peer %d: New channel ID %d \n", self.Config.PubKey[:3], message.ChannelID)
	self.outputMessagesReceived <- domain.MeshMessage{
		Contents: message.ChannelID[:],
	}
}

func (self *Node) processSetControlChannelMessage(message domain.SetControlChannelMessage) {
	fmt.Fprintf(os.Stdout, "Peer %d: Processing 'SetControlChannel' message from peer %d! \n", self.Config.PubKey[:3], message.FromPeerID[:3])

	controlChannel := NewControlChannel()
	self.AddControlChannel(controlChannel)

	replyMessage := domain.SetControlChannelResponseMessage{
		FromPeerID: self.Config.PubKey,
		ChannelID:  controlChannel.ID,
	}

	transportToPeer := self.safelyGetTransportToPeer(message.FromPeerID)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v from %v, dropping\n", message.FromPeerID, self.Config.PubKey)
		return
	}

	serialized := self.serializer.SerializeMessage(replyMessage)
	err := transportToPeer.SendMessage(message.FromPeerID, serialized, nil)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stdout, "Peer %d: Sending 'SetControlChannelResponse' message to peer %d! \n", self.Config.PubKey[:3], message.FromPeerID[:3])
}

func (self *Node) processUserMessage(incomingMessage domain.UserMessage) {
	fmt.Fprintf(os.Stdout, "Peer %d: Processing 'UserMessage' from peer %d! \n", self.Config.PubKey[:3], incomingMessage.FromPeerID[:3])
	directPeerID, forwardBase, doForward := self.safelyGetRewriteBase(incomingMessage)
	fmt.Fprintf(os.Stdout, "Peer %d: Do forward: %v\n", self.Config.PubKey[:3], doForward)

	if doForward {
		transportToPeer := self.safelyGetTransportToPeer(directPeerID)
		if transportToPeer == nil {
			fmt.Fprintf(os.Stderr, "No transport to peer %v from %v, dropping\n", directPeerID, self.Config.PubKey)
			return
		}

		message := domain.UserMessage{
			MessageBase: forwardBase,
			MessageID:   (domain.MessageID)(uuid.NewV4()),
			Index:       0,
			Count:       1,
			Contents:    incomingMessage.Contents,
		}
		serialized := self.serializer.SerializeMessage(message)
		err := transportToPeer.SendMessage(directPeerID, serialized, nil)
		if err != nil {
			fmt.Fprint(os.Stdout, "Failed to send forwarded message, dropping\n")
			return
		}
	} else {
		self.outputMessagesReceived <- domain.MeshMessage{
			ReplyTo: domain.ReplyTo{
				incomingMessage.SendRouteID,
				incomingMessage.FromPeerID,
			},
			Contents: incomingMessage.Contents,
		}
	}
}

func (self *Node) processSetRouteMessage(msg domain.SetRouteMessage) {
	if msg.SendBack {
		fmt.Fprintf(os.Stderr, "Invalid SetRouteMessage received, dropping: %v\n", msg)
		return
	}
	self.lock.Lock()
	defer self.lock.Unlock()

	if msg.SetRouteID == NilRouteID {
		fmt.Fprintf(os.Stderr, "Invalid SetRouteMessage received, dropping: %v\n", msg)
		return
	}
	self.routes[msg.SetRouteID] =
		domain.Route{
			ForwardToPeerID:   msg.ForwardToPeerID,
			ForwardToRouteID:  msg.ForwardRewriteSendRouteID,
			BackwardToPeerID:  msg.BackwardToPeerID,
			BackwardToRouteID: msg.BackwardRewriteSendRouteID,
			ExpiryTime:        self.clipExpiryTime(time.Now().Add(msg.DurationHint)),
		}

	// Don't block to send reply
	go self.sendSetRouteReply(msg.MessageBase, msg.ConfirmRouteID)
}

func (self *Node) processSetRouteReplyMessage(msg domain.SetRouteReply) {
	self.lock.Lock()
	defer self.lock.Unlock()
	confirmChan, foundChan := self.routeExtensionsAwaitingConfirm[msg.ConfirmRouteID]
	if foundChan {
		confirmChan <- true
	}
	localRoute, foundLocal := self.localRoutes[msg.ConfirmRouteID]
	if foundLocal {
		localRoute.LastConfirmed = time.Now()
		self.localRoutes[msg.ConfirmRouteID] = localRoute
	}
}

func (self *Node) processRefreshRouteMessage(msg domain.RefreshRouteMessage, forwarded bool) {
	if forwarded {
		self.lock.Lock()
		defer self.lock.Unlock()
		route, exists := self.routes[msg.SendRouteID]
		if !exists {
			fmt.Fprintf(os.Stderr, "Refresh sent for unknown route: %v\n", msg.SendRouteID)
			return
		}
		route.ExpiryTime = self.clipExpiryTime(time.Now().Add(msg.DurationHint))
		self.routes[msg.SendRouteID] = route
	} else {
		// Don't block to send reply
		go self.sendSetRouteReply(msg.MessageBase, msg.ConfirmRoutedID)
	}
}

func (self *Node) processDeleteRouteMessage(msg domain.DeleteRouteMessage) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.routes, msg.SendRouteID)
}

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
	fmt.Fprintf(os.Stdout, "Incoming message %v\n", serialized)
	fmt.Printf("Incoming message %v\n", serialized)
	fmt.Println("Incoming message %v\n", serialized)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Deserialization error %v %v\n", serialized, err)
		return
	}

	messageType := reflect.TypeOf(message)
	// User messages have fragmentation to deal with
	if messageType == reflect.TypeOf(domain.UserMessage{}) {
		self.processUserMessage(message.(domain.UserMessage))
	} else {
		forwardedMessage := self.forwardMessage(message)
		if !forwardedMessage {
			// Receive or forward. Refragment on forward!
			if messageType == reflect.TypeOf(domain.SetRouteMessage{}) {
				self.processSetRouteMessage(message.(domain.SetRouteMessage))
			} else if messageType == reflect.TypeOf(domain.SetRouteReply{}) {
				self.processSetRouteReplyMessage(message.(domain.SetRouteReply))
			}
		} else {
			if messageType == reflect.TypeOf(domain.DeleteRouteMessage{}) {
				self.processDeleteRouteMessage(message.(domain.DeleteRouteMessage))
			}
		}

		if messageType == reflect.TypeOf(domain.RefreshRouteMessage{}) {
			self.processRefreshRouteMessage(message.(domain.RefreshRouteMessage), forwardedMessage)
		}
	}
}

func (self *Node) processUserMessage(incomingMessage domain.UserMessage) {

	directPeerID, forwardBase, doForward := self.safelyGetRewriteBase(incomingMessage)
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
			fmt.Fprint(os.Stderr, "Failed to send forwarded message, dropping\n")
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
			ForwardToPeerID:            msg.ForwardToPeerID,
			ForwardRewriteSendRouteID:  msg.ForwardRewriteSendRouteID,
			BackwardToPeerID:           msg.BackwardToPeerID,
			BackwardRewriteSendRouteID: msg.BackwardRewriteSendRouteID,
			ExpiryTime:                 self.clipExpiryTime(time.Now().Add(msg.DurationHint)),
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

package mesh

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/node/connection"
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
	if err != nil {
		fmt.Fprintf(os.Stderr, "Deserialization error %v\n", err)
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

func (self *Node) processUserMessage(msgIn domain.UserMessage) {
	reassembled := self.reassembleUserMessage(msgIn)
	// Not finished reassembling yet
	if reassembled == nil {
		return
	}
	directPeer, forwardBase, doForward := self.safelyGetRewriteBase(msgIn)
	if doForward {
		transportToPeer := self.safelyGetTransportToPeer(directPeer)
		if transportToPeer == nil {
			fmt.Fprintf(os.Stderr, "No transport to peer %v from %v, dropping\n", directPeer, self.config.PubKey)
			return
		}
		// Forward reassembled message, not individual pieces. This is done because of the need for refragmentation
		fragments := connection.ConnectionManager.FragmentMessage(reassembled, directPeer, transportToPeer, forwardBase)
		for _, fragment := range fragments {
			serialized := self.serializer.SerializeMessage(fragment)
			err := transportToPeer.SendMessage(directPeer, serialized)
			if err != nil {
				fmt.Fprint(os.Stderr, "Failed to send forwarded message, dropping\n")
				return
			}
		}
	} else {
		self.outputMessagesReceived <- domain.MeshMessage{domain.ReplyTo{msgIn.SendRouteID, msgIn.FromPeerID}, reassembled}
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

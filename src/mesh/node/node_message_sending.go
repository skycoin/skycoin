package mesh

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/node/connection"
)

// Chooses a route automatically. Sends directly without a route if connected to that peer
func (self *Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte) (error, domain.RouteID) {
	directPeer, localRouteID, sendRoutID, transportToPeer, err := self.findRouteToPeer(toPeer)
	if err != nil {
		return err, NilRouteID
	}
	messageBase := domain.MessageBase{
		SendRouteID: sendRoutID,
		SendBack:    false,
		FromPeerID:  self.config.PubKey,
		Nonce:       generateNonce(),
	}
	messages := connection.ConnectionManager.FragmentMessage(contents, directPeer, transportToPeer, messageBase)
	for _, message := range messages {
		serialized := self.serializer.SerializeMessage(message)
		err := transportToPeer.SendMessage(directPeer, serialized)
		if err != nil {
			return err, NilRouteID
		}
	}
	return nil, localRouteID
}

// Blocks until message is confirmed received
func (self *Node) SendMessageThruRoute(routeId domain.RouteID, contents []byte) error {
	route, routeFound := self.safelyGetRoute(routeId)
	if !routeFound {
		return errors.New("Route not found")
	}

	base := domain.MessageBase{
		SendRouteID: route.ForwardRewriteSendRouteID,
		SendBack:    false,
		FromPeerID:  self.config.PubKey,
		Nonce:       generateNonce(),
	}
	directPeer := route.ForwardToPeerID
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", directPeer))
	}
	messages := connection.ConnectionManager.FragmentMessage(contents, directPeer, transportToPeer, base)
	for _, message := range messages {
		serialized := self.serializer.SerializeMessage(message)
		fmt.Fprintln(os.Stdout, "Send Message")
		err := transportToPeer.SendMessage(directPeer, serialized)
		if err != nil {
			return err
		}
	}
	return nil
}

// Blocks until message is confirmed received
func (self *Node) SendMessageBackThruRoute(replyTo domain.ReplyTo, contents []byte) error {
	directPeer := replyTo.FromPeer
	transportToPeer := self.safelyGetTransportToPeer(directPeer)
	if transportToPeer == nil {
		return errors.New(fmt.Sprintf("No route or transport to peer %v\n", directPeer))
	}
	base := domain.MessageBase{
		SendRouteID: replyTo.RouteID,
		SendBack:    true,
		FromPeerID:  self.config.PubKey,
		Nonce:       generateNonce(),
	}
	messages := connection.ConnectionManager.FragmentMessage(contents, directPeer, transportToPeer, base)
	for _, message := range messages {
		serialized := self.serializer.SerializeMessage(message)
		err := transportToPeer.SendMessage(directPeer, serialized)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Node) expireOldMessages() {
	time_now := time.Now()
	self.lock.Lock()
	defer self.lock.Unlock()

	lastMessages := self.messagesBeingAssembled
	self.messagesBeingAssembled = make(map[domain.MessageID]*domain.MessageUnderAssembly)
	for id, msg := range lastMessages {
		if time_now.Before(msg.ExpiryTime) {
			self.messagesBeingAssembled[id] = msg
		}
	}
}

func (self *Node) expireOldMessagesLoop() {
	self.closeGroup.Add(1)
	defer self.closeGroup.Done()
	for len(self.closing) == 0 {
		select {
		case <-time.After(self.config.ExpireMessagesInterval):
			{
				self.expireOldMessages()
			}
		case <-self.closing:
			{
				return
			}
		}
	}
}

func (self *Node) sendSetRouteReply(base domain.MessageBase, confirmId domain.RouteID) {
	replyMessage := domain.SetRouteReply{
		MessageBase: domain.MessageBase{
			SendRouteID: base.SendRouteID,
			SendBack:    true,
			FromPeerID:  self.config.PubKey,
			Nonce:       generateNonce(),
		},
		ConfirmRouteID: confirmId,
	}
	transportToPeer := self.safelyGetTransportToPeer(base.FromPeerID)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v from %v\n", base.FromPeerID, self.config.PubKey)
		return
	}
	serialized := self.serializer.SerializeMessage(replyMessage)
	err := transportToPeer.SendMessage(base.FromPeerID, serialized)
	if err != nil {
		return
	}
}

func (self *Node) forwardMessage(msg interface{}) bool {
	forwardTo, newBase, doForward := self.safelyGetRewriteBase(msg)
	if !doForward {
		return false
	}
	// Rewrite
	rewritten := rewriteMessage(msg, newBase)
	transportToPeer := self.safelyGetTransportToPeer(forwardTo)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport found for forwarded message from %v to %v, dropping\n", self.config.PubKey, forwardTo)
		return true
	}

	serialized := self.serializer.SerializeMessage(rewritten)
	err := transportToPeer.SendMessage(forwardTo, serialized)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to send forwarded message, dropping\n")
		return true
	}

	// Forward, not receive
	return true
}

func (self *Node) safelyGetRewriteBase(msg interface{}) (forwardTo cipher.PubKey, base domain.MessageBase, doForward bool) {
	// sendBack
	sendBack, route, foundRoute := self.safelyGetForwarding(msg)
	if !foundRoute {
		return cipher.PubKey{}, domain.MessageBase{}, false
	}
	forwardTo = route.ForwardToPeerID
	rewriteTo := route.ForwardRewriteSendRouteID
	if sendBack {
		forwardTo = route.BackwardToPeerID
		rewriteTo = route.BackwardRewriteSendRouteID
	}
	if forwardTo == (cipher.PubKey{}) {
		return cipher.PubKey{}, domain.MessageBase{}, false
	}
	newBase :=
		domain.MessageBase{
			SendRouteID: rewriteTo,
			SendBack:    sendBack,
			FromPeerID:  self.config.PubKey,
			Nonce:       generateNonce(),
		}
	return forwardTo, newBase, true
}

func rewriteMessage(message interface{}, newBase domain.MessageBase) interface{} {
	messageType := reflect.TypeOf(message)
	newBase.Nonce = generateNonce()

	switch messageType {
	case reflect.TypeOf(domain.UserMessage{}):
		newMessage := (message.(domain.UserMessage))
		newMessage.MessageBase = newBase
		return newMessage

	case reflect.TypeOf(domain.SetRouteMessage{}):
		newMessage := (message.(domain.SetRouteMessage))
		newMessage.MessageBase = newBase
		return newMessage

	case reflect.TypeOf(domain.RefreshRouteMessage{}):
		newMessage := (message.(domain.RefreshRouteMessage))
		newMessage.MessageBase = newBase
		return newMessage

	case reflect.TypeOf(domain.DeleteRouteMessage{}):
		newMessage := (message.(domain.DeleteRouteMessage))
		newMessage.MessageBase = newBase
		return newMessage

	case reflect.TypeOf(domain.SetRouteReply{}):
		newMessage := (message.(domain.SetRouteReply))
		newMessage.MessageBase = newBase
		return newMessage
	}

	panic("Internal error: rewriteMessage incomplete")
}

func (self *Node) debug_countMessages() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return len(self.messagesBeingAssembled)
}

func generateNonce() [4]byte {
	ret := make([]byte, 4)
	n, err := rand.Read(ret)
	if n != 4 {
		panic("rand.Read() failed")
	}
	if err != nil {
		panic(err)
	}
	ret_b := [4]byte{0, 0, 0, 0}
	copy(ret_b[:], ret[:])
	return ret_b
}

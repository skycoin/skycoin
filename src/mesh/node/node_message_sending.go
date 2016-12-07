package mesh

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
)

// Chooses a route automatically. Sends directly without a route if connected to that peer
func (self *Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte) (error, domain.RouteID) {
	directPeerID, localRouteID, sendRoutID, transport, err := self.findRouteToPeer(toPeer)
	if err != nil {
		return err, NilRouteID
	}

	message := domain.UserMessage{
		MessageBase: domain.MessageBase{
			SendRouteID: sendRoutID,
			SendBack:    false,
			FromPeerID:  self.Config.PubKey,
			Nonce:       GenerateNonce(),
		},
		MessageID: (domain.MessageID)(uuid.NewV4()),
		Index:     0,
		Count:     1,
		Contents:  contents,
	}

	serialized := self.serializer.SerializeMessage(message)
	err = transport.SendMessage(directPeerID, serialized, nil)
	if err != nil {
		return err, NilRouteID
	}
	return nil, localRouteID
}

// Blocks until message is confirmed received
func (self *Node) SendMessageThruRoute(routeID domain.RouteID, contents []byte) error {
	route, ok := self.safelyGetRoute(routeID)
	if !ok {
		return errors.New("Route not found")
	}

	directPeerID := route.ForwardToPeerID
	transport := self.safelyGetTransportToPeer(directPeerID)
	if transport == nil {
		return errors.New(fmt.Sprintf("No transport to peer %v\n", directPeerID))
	}

	message := domain.UserMessage{
		MessageBase: domain.MessageBase{
			SendRouteID: route.ForwardRewriteSendRouteID,
			SendBack:    false,
			FromPeerID:  self.Config.PubKey,
			Nonce:       GenerateNonce(),
		},
		MessageID: (domain.MessageID)(uuid.NewV4()),
		Index:     0,
		Count:     1,
		Contents:  contents,
	}

	serialized := self.serializer.SerializeMessage(message)
	//	fmt.Fprintln(os.Stdout, "Send Message")
	err := transport.SendMessage(directPeerID, serialized, nil)
	if err != nil {
		return err
	}

	return nil
}

// Blocks until message is confirmed received
func (self *Node) SendMessageBackThruRoute(replyTo domain.ReplyTo, contents []byte) error {
	directPeerID := replyTo.FromPeerID
	transport := self.safelyGetTransportToPeer(directPeerID)
	if transport == nil {
		return errors.New(fmt.Sprintf("No route or transport to peer %v\n", directPeerID))
	}

	message := domain.UserMessage{
		MessageBase: domain.MessageBase{
			SendRouteID: replyTo.RouteID,
			SendBack:    true,
			FromPeerID:  self.Config.PubKey,
			Nonce:       GenerateNonce(),
		},
		MessageID: (domain.MessageID)(uuid.NewV4()),
		Index:     0,
		Count:     1,
		Contents:  contents,
	}

	serialized := self.serializer.SerializeMessage(message)
	err := transport.SendMessage(directPeerID, serialized, nil)
	if err != nil {
		return err
	}
	return nil
}

func (self *Node) sendSetRouteReply(base domain.MessageBase, confirmID domain.RouteID) {
	replyMessage := domain.SetRouteReply{
		MessageBase: domain.MessageBase{
			SendRouteID: base.SendRouteID,
			SendBack:    true,
			FromPeerID:  self.Config.PubKey,
			Nonce:       GenerateNonce(),
		},
		ConfirmRouteID: confirmID,
	}
	transportToPeer := self.safelyGetTransportToPeer(base.FromPeerID)
	if transportToPeer == nil {
		fmt.Fprintf(os.Stderr, "No transport to peer %v from %v\n", base.FromPeerID, self.Config.PubKey)
		return
	}
	serialized := self.serializer.SerializeMessage(replyMessage)
	err := transportToPeer.SendMessage(base.FromPeerID, serialized, nil)
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
		fmt.Fprintf(os.Stderr, "No transport found for forwarded message from %v to %v, dropping\n", self.Config.PubKey, forwardTo)
		return true
	}

	serialized := self.serializer.SerializeMessage(rewritten)
	err := transportToPeer.SendMessage(forwardTo, serialized, nil)
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
			FromPeerID:  self.Config.PubKey,
			Nonce:       GenerateNonce(),
		}
	return forwardTo, newBase, true
}

func rewriteMessage(message interface{}, newBase domain.MessageBase) interface{} {
	messageType := reflect.TypeOf(message)
	newBase.Nonce = GenerateNonce()

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

func GenerateNonce() [4]byte {
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

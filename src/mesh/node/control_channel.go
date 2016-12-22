package mesh

import (
	"reflect"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/mesh/domain"
)

type ControlChannel struct {
	ID uuid.UUID
}

func NewControlChannel() *ControlChannel {
	c := ControlChannel{
		ID: uuid.NewV4(),
	}
	return &c
}

func (c *ControlChannel) HandleMessage(node *Node, message interface{}) error {

	messageType := reflect.TypeOf(message)

	if messageType == reflect.TypeOf(domain.SetRouteControlMessage{}) {
		err := processSetRouteMessage(node, message.(domain.SetRouteControlMessage))
		if err != nil {
			return err
		}
	}

	return nil
}

//must return a confirmation message
func processSetRouteMessage(node *Node, message domain.SetRouteControlMessage) error {
	return node.AddRoute(message.ForwardRouteID, message.ForwardPeerID)
}

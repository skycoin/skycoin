package mesh

import (
	"errors"
	"fmt"

	"github.com/satori/go.uuid"
)

func (self *Node) AddControlChannel(channel *ControlChannel) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.controlChannels[channel.ID] = channel
}

func (self *Node) RemoveControlChannel(channelID uuid.UUID) {
	self.lock.Lock()
	defer self.lock.Unlock()

	delete(self.controlChannels, channelID)
}

func (self *Node) HandleControlMessage(channelID uuid.UUID, message interface{}) error {

	channel, ok := self.controlChannels[channelID]
	if !ok {
		return errors.New(fmt.Sprintf("Control channel %s not found", channelID))
	}

	err := channel.HandleMessage(self, message)
	if err != nil {
		return err
	}

	return nil
}

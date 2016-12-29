package node

import (
	"errors"
	"fmt"

	"github.com/satori/go.uuid"
)

func (self *Node) AddControlChannel(channel *ControlChannel) {
	//self.lock.Lock()
	//defer self.lock.Unlock()

	self.controlChannels[channel.Id] = channel
}

func (self *Node) RemoveControlChannel(channelID uuid.UUID) {
	//self.lock.Lock()
	//defer self.lock.Unlock()

	delete(self.controlChannels, channelID)
}

func (self *Node) HandleControlMessage(channelID uuid.UUID, message interface{}) error {

	channel, ok := self.controlChannels[channelID]
	if !ok {
		return errors.New(fmt.Sprintf("Control channel %s not found", channelID))
	}

	err := channel.HandleMessage(self, message.([]byte))
	if err != nil {
		return err
	}

	return nil
}

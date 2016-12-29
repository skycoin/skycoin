package node

import (
	"github.com/satori/go.uuid"
)

func (self *Node) AddControlChannel(channel *ControlChannel) {
	//self.lock.Lock()
	//defer self.lock.Unlock()

	self.ControlChannels[channel.Id] = channel
}

func (self *Node) RemoveControlChannel(channelID uuid.UUID) {
	//self.lock.Lock()
	//defer self.lock.Unlock()

	delete(self.ControlChannels, channelID)
}

func (self *Node) NumControlChannels() int { // for testing purposes
	return len(self.ControlChannels)
}

/*
func (self *Node) HandleControlMessage(channelID uuid.UUID, message interface{}) error {

	channel, ok := self.ControlChannels[channelID]
	if !ok {
		return errors.New(fmt.Sprintf("Control channel %s not found", channelID))
	}

	err := channel.HandleMessage(self, message.([]byte))
	if err != nil {
		return err
	}

	return nil
}
*/

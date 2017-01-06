package node

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

func (self *Node) AddControlChannel() messages.ChannelId {
	//self.lock.Lock()
	//defer self.lock.Unlock()

	channel := NewControlChannel()

	self.controlChannels[channel.Id] = channel
	return channel.Id
}

func (self *Node) CloseControlChannel(channelID messages.ChannelId) error {
	//self.lock.Lock()
	//defer self.lock.Unlock()

	if _, ok := self.controlChannels[channelID]; !ok {
		return errors.New(fmt.Sprintf("Control channel %s not found", channelID))
	}

	delete(self.controlChannels, channelID)
	return nil
}

func (self *Node) HandleControlMessage(channelID messages.ChannelId, message []byte) (interface{}, error) {

	channel, ok := self.controlChannels[channelID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Control channel %s not found", channelID))
	}

	return channel.HandleMessage(self, message)
}

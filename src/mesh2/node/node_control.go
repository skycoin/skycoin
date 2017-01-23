package node

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/mesh2/messages"
)

func (self *Node) AddControlChannel() messages.ChannelId {
	//self.lock.Lock()
	//defer self.lock.Unlock()

	channel := newControlChannel()

	self.controlChannels[channel.id] = channel
	return channel.id
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

func (self *Node) handleControlMessage(channelID messages.ChannelId, message []byte) (interface{}, error) {

	channel, ok := self.controlChannels[channelID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Control channel %s not found", channelID))
	}

	return channel.handleMessage(self, message)
}

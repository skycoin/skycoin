package node

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

func (self *Node) AddControlChannel() messages.ChannelId {

	channel := newControlChannel()

	self.lock.Lock()
	defer self.lock.Unlock()

	self.controlChannels[channel.id] = channel
	return channel.id
}

func (self *Node) CloseControlChannel(channelID messages.ChannelId) error {

	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.controlChannels[channelID]; !ok {
		return errors.New(fmt.Sprintf("Control channel %s not found", channelID))
	}

	delete(self.controlChannels, channelID)
	return nil
}

func (self *Node) handleControlMessage(channelID messages.ChannelId, message []byte) error {

	self.lock.Lock()

	channel, ok := self.controlChannels[channelID]
	self.lock.Unlock()
	if !ok {
		return errors.New(fmt.Sprintf("Control channel %s not found", channelID))
	}

	return channel.handleMessage(self, message)
}

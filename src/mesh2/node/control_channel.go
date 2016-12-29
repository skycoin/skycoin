package node

import (
	"github.com/satori/go.uuid"
)

type ControlChannel struct {
	Id              uuid.UUID
	IncomingChannel chan []byte
}

func NewControlChannel() *ControlChannel {
	c := ControlChannel{
		Id:              uuid.NewV4(),
		IncomingChannel: make(chan []byte),
	}
	return &c
}

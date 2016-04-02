package main

import (
    "github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/mesh")

import "fmt"

// Can't add Handle functions to out-of-package types, so create wrappers
type SendMessageWrapper struct {
    mesh.SendMessage
}
var SendMessagePrefix = gnet.MessagePrefix{0,0,0,1}
func (self *SendMessageWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    fmt.Println("SendMessage")
    return nil
}


type EstablishRouteMessageWrapper struct {
    mesh.EstablishRouteMessage
}
var EstablishRouteMessagePrefix = gnet.MessagePrefix{0,0,0,2}
func (self *EstablishRouteMessageWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    fmt.Println("EstablishRouteMessage")
    return nil
}


type EstablishRouteReplyMessageWrapper struct {
    mesh.EstablishRouteReplyMessage
}
var EstablishRouteReplyMessagePrefix = gnet.MessagePrefix{0,0,0,3}
func (self *EstablishRouteReplyMessageWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    fmt.Println("EstablishRouteReplyMessage")
    return nil
}


type SetRouteRewriteIdMessageWrapper struct {
    mesh.SetRouteRewriteIdMessage
}
var SetRouteRewriteIdMessagePrefix = gnet.MessagePrefix{0,0,0,4}
func (self *SetRouteRewriteIdMessageWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    fmt.Println("SetRouteRewriteIdMessage %i", self.RewriteSendId)
    return nil
}

func RegisterTCPMessages() {
    gnet.RegisterMessage(SendMessagePrefix, SendMessageWrapper{})
    gnet.RegisterMessage(EstablishRouteMessagePrefix, EstablishRouteMessageWrapper{})
    gnet.RegisterMessage(EstablishRouteReplyMessagePrefix, EstablishRouteReplyMessageWrapper{})
    gnet.RegisterMessage(SetRouteRewriteIdMessagePrefix, SetRouteRewriteIdMessageWrapper{})
}
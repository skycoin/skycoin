package main

import (
    "github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/mesh")

import ("reflect")

// Can't add Handle functions to out-of-package types, so create wrappers
type SendMessageWrapper struct {
    mesh.SendMessage
}
var SendMessagePrefix = gnet.MessagePrefix{0,0,0,1}
func (self *SendMessageWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    var node_impl = (x).(*mesh.Node)
    node_impl.MessagesIn <- self.SendMessage
    return nil
}

type EstablishRouteMessageWrapper struct {
    mesh.EstablishRouteMessage
}
var EstablishRouteMessagePrefix = gnet.MessagePrefix{0,0,0,2}
func (self *EstablishRouteMessageWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    var node_impl = (x).(*mesh.Node)
    node_impl.MessagesIn <- self.EstablishRouteMessage
    return nil
}


type EstablishRouteReplyMessageWrapper struct {
    mesh.EstablishRouteReplyMessage
}
var EstablishRouteReplyMessagePrefix = gnet.MessagePrefix{0,0,0,3}
func (self *EstablishRouteReplyMessageWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    var node_impl = (x).(*mesh.Node)
    node_impl.MessagesIn <- self.EstablishRouteReplyMessage
    return nil
}


type RouteRewriteMessageWrapper struct {
    mesh.RouteRewriteMessage
}
var RouteRewriteMessagePrefix = gnet.MessagePrefix{0,0,0,4}
func (self *RouteRewriteMessageWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    var node_impl = (x).(*mesh.Node)
    node_impl.MessagesIn <- self.RouteRewriteMessage
    return nil
}

type OperationReplyWrapper struct {
    mesh.OperationReply
}
var OperationReplyPrefix = gnet.MessagePrefix{0,0,0,5}
func (self *OperationReplyWrapper) Handle(context *gnet.MessageContext, x interface{}) error {
    var node_impl = (x).(*mesh.Node)
    node_impl.MessagesIn <- self.OperationReply
    return nil
}

func RegisterTCPMessages() {
    gnet.RegisterMessage(SendMessagePrefix, SendMessageWrapper{})
    gnet.RegisterMessage(EstablishRouteMessagePrefix, EstablishRouteMessageWrapper{})
    gnet.RegisterMessage(EstablishRouteReplyMessagePrefix, EstablishRouteReplyMessageWrapper{})
    gnet.RegisterMessage(RouteRewriteMessagePrefix, RouteRewriteMessageWrapper{})
    gnet.RegisterMessage(OperationReplyPrefix, OperationReplyWrapper{})
}

func WrapMessage(msg interface{}) gnet.Message {
    msg_type := reflect.TypeOf(msg)
    if msg_type == reflect.TypeOf(mesh.SendMessage{}) {
        wrapped := &SendMessageWrapper{msg.(mesh.SendMessage)}
        return wrapped
    } else if msg_type == reflect.TypeOf(mesh.EstablishRouteMessage{}) {
        wrapped := &EstablishRouteMessageWrapper{msg.(mesh.EstablishRouteMessage)}
        return wrapped
    } else if msg_type == reflect.TypeOf(mesh.EstablishRouteReplyMessage{}) {
        wrapped := &EstablishRouteReplyMessageWrapper{msg.(mesh.EstablishRouteReplyMessage)}
        return wrapped
    } else if msg_type == reflect.TypeOf(mesh.RouteRewriteMessage{}) {
        wrapped := &RouteRewriteMessageWrapper{msg.(mesh.RouteRewriteMessage)}
        return wrapped
    } else if msg_type == reflect.TypeOf(mesh.OperationReply{}) {
        wrapped := &OperationReplyWrapper{msg.(mesh.OperationReply)}
        return wrapped
    }
    panic("Unknown message type passed to WrapMessage()")
}

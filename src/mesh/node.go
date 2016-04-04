package mesh

import (
    "time"
)

import (
    "github.com/skycoin/skycoin/src/cipher"
)

type SendMessage struct {
    SendId uint32
    Message []byte
}

type OperationMessage struct {
    MsgIndex int64
}

type OperationReply struct {
    OperationMessage
    Success bool
    Error string
}

type EstablishRouteMessage struct {
    OperationMessage
    ToPubKey cipher.PubKey
    DurationHint time.Duration
}

type EstablishRouteReplyMessage struct {
    OperationReply
    SendId uint32
    // Secret be sent back in SetRouteRewriteIdMessage
    Secret string
}

type SetRouteRewriteIdMessage struct {
    // Secret from EstablishRouteReply
    Secret string
    RewriteSendId uint32
}

type Node struct {
    MessagesOut chan interface{}
    MessagesIn chan interface{} 
}


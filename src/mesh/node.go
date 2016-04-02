package mesh

import (
    //"github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/cipher"
 //   "github.com/skycoin/encoder"
/*
    "encoding/json"
    "os"
    "log"
    "io/ioutil"
    "flag"
    */
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
//    MessagesIn chan 
}


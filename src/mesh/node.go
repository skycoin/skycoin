package mesh

import (
    "time"

    "fmt"
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

type QueryConnectedPeersMessage struct {
    OperationMessage
}

type QueryConnectedPeersReply struct {
    OperationReply
    PubKey []cipher.PubKey
}

type SetRouteRewriteIdMessage struct {
    // Secret from EstablishRouteReply
    Secret string
    RewriteSendId uint32
}

type OutgoingMessage struct {
    ConnectedPeerPubKey cipher.PubKey
    Message interface{}
}

type Node struct {
    ConnectedPeers []cipher.PubKey
    MessagesOut chan OutgoingMessage
    MessagesIn chan interface{}
}

type Config struct {
    MyPubKey cipher.PubKey
    MessagesOutSize int
    MessagesInSize int
}

func NewNode(config Config) *Node {
    ret := &Node{}
    ret.MessagesOut = make(chan OutgoingMessage, config.MessagesOutSize)
    ret.MessagesIn = make(chan interface{}, config.MessagesInSize)
    return ret
}

// Blocks
func (self *Node) Run() {
    // TODO: Establish routes

    // Test
    go func() {
        for {
            var test OutgoingMessage
            key := cipher.NewPubKey([]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
            test.ConnectedPeerPubKey = key
            test.Message = SendMessage{5, []byte{13, 15, 17}}
            self.MessagesOut <- test
            time.Sleep(time.Second / 2)
        fmt.Printf("sending %v\n", test)
        }
    }()

    for{
        msg := <- self.MessagesIn
        fmt.Printf("msg %v\n", msg)
    }
}



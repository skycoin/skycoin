package mesh

import (
    "time"
    "fmt"
    "sync"
    "reflect"
    "gopkg.in/op/go-logging.v1"
)

import (
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid"
)

var logger = logging.MustGetLogger("node")

type Message struct {
    SendId uint32
}

type SendMessage struct {
    Message
    Contents []byte
}

type OperationMessage struct {
    Message
    MsgId uuid.UUID
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
    // Secret be sent back in ActivateRouteMessage
    Secret string
}

type QueryConnectedPeersMessage struct {
    OperationMessage
}

type QueryConnectedPeersReply struct {
    OperationReply
    PubKey []cipher.PubKey
}

type RouteRewriteMessage struct {
    OperationMessage
    // Secret from EstablishRouteReply
    Secret string
    RewriteSendId uint32
}

type OutgoingMessage struct {
    ConnectedPeerPubKey cipher.PubKey
    Message interface{}
}

type EstablishedRoute struct {
    SendId          uint32
    ConnectedPeer   cipher.PubKey
}

type Node struct {
    Config NodeConfig
    MessagesOut chan OutgoingMessage
    MessagesIn chan interface{}
    EstablishedRoutesByIndex map[int]EstablishedRoute
    
    RetransmitQueue map[uuid.UUID]OutgoingMessage
    OperationsAwaitingReply map[uuid.UUID]chan interface{} 
    OperationsAwaitingReplyLock *sync.Mutex
}

type RouteEstablishedCallback func(route EstablishedRoute)

type NodeConfig struct {
    MyPubKey cipher.PubKey
    ConnectedPeers []cipher.PubKey
    MessagesOutSize int
    MessagesInSize int
    Routes []RouteConfig
    OperationTimeout time.Duration
    RetransmitInterval time.Duration
    RouteEstablishedCB RouteEstablishedCallback
}

type RouteConfig struct {
    PeerPubKeys []cipher.PubKey
}

func NewNode(config NodeConfig) *Node {
    ret := &Node{}
    ret.Config = config
    ret.MessagesOut = make(chan OutgoingMessage, config.MessagesOutSize)
    ret.MessagesIn = make(chan interface{}, config.MessagesInSize)
    ret.OperationsAwaitingReplyLock = &sync.Mutex{}
    ret.EstablishedRoutesByIndex = make(map[int]EstablishedRoute)
    ret.RetransmitQueue = make(map[uuid.UUID]OutgoingMessage)
    ret.OperationsAwaitingReply = make(map[uuid.UUID]chan interface{})
    return ret
}

func (self *Node) sendAndConfirmOperation(connected cipher.PubKey, operation_id uuid.UUID, msg interface{}) interface{} {
    reply_channel := make(chan interface{})
    self.OperationsAwaitingReplyLock.Lock()
    self.OperationsAwaitingReply[operation_id] = reply_channel
    outgoing := OutgoingMessage{connected, msg}
    self.RetransmitQueue[operation_id] = outgoing
    self.OperationsAwaitingReplyLock.Unlock()

    self.MessagesOut <- outgoing

    select {
        case reply := <- reply_channel:
            return reply
        case <- time.After(self.Config.OperationTimeout):
            logger.Debug(fmt.Sprintf("Operation %v timed out", self.Config.OperationTimeout))
    }
    return nil
}

func (self *Node) establishRoute(route_idx int, route RouteConfig) {
    if len(route.PeerPubKeys) == 0 {
        logger.Debug("Empty route passed to establishRoute()")
        return
    }

    var new_route = EstablishedRoute{0, route.PeerPubKeys[0]}
//    var last_secret = ""

    for peer_idx, _ /*peer_key*/ := range route.PeerPubKeys {
        // For each peer, we establish a route from this one to the next one
        // Nothing needs to be established for the destination
        if(peer_idx == (len(route.PeerPubKeys) - 1)) {
            break
        }

        // Establish forwarding
        establish_id := uuid.NewV4()
        route_message := 
            EstablishRouteMessage{
                OperationMessage {
                    Message {
                        new_route.SendId,
                    },
                    establish_id,
                },
                route.PeerPubKeys[peer_idx + 1],
                time.Hour,
            }
        reply := self.sendAndConfirmOperation(new_route.ConnectedPeer,
                                                establish_id,
                                                route_message)
        if reply != nil {
            route_reply := reply.(EstablishRouteReplyMessage)
            if peer_idx == 0 {
                new_route.SendId = route_reply.SendId;
            }
            logger.Debug("TODO: route reply %v", route_reply)
        }

/*
    // Activate the hop established before this one
    uint32 send_id = 0;

    if peer_idx < (len(route.PeerPubKeys) - 1) {
        send_id = route_reply.SendId
    }

    go sendAndConfirmOperation(connected_peer_key,
                                 RouteRewriteMessage{last_secret,
                                                      send_id)
            last_secret = route_reply.Secret
*/
    }

    self.EstablishedRoutesByIndex[route_idx] = new_route
    if(self.Config.RouteEstablishedCB != nil) {
        self.Config.RouteEstablishedCB(new_route)
    }
}

func (self *Node) onOperationReply(msg interface{}, operation_id uuid.UUID) {
    self.OperationsAwaitingReplyLock.Lock()
    ch := self.OperationsAwaitingReply[operation_id]
    delete(self.OperationsAwaitingReply, operation_id)
    self.OperationsAwaitingReplyLock.Unlock()

    ch <- msg
}

// Blocks
func (self *Node) Run() {
    // Retransmit loop
    go func() {
        for {
            self.OperationsAwaitingReplyLock.Lock()
            for _, outgoing := range self.RetransmitQueue {
                self.MessagesOut <- outgoing
            }
            self.OperationsAwaitingReplyLock.Unlock()
            time.Sleep(self.Config.RetransmitInterval)
        }
    }()

    // Establish routes asynchronously
    for route_idx, route := range self.Config.Routes {
        go self.establishRoute(route_idx, route)
    }

    // Incoing messages loop
    for{
        msg := <- self.MessagesIn
        msg_type := reflect.TypeOf(msg)

        if msg_type == reflect.TypeOf(OperationReply{}) {
            self.onOperationReply(msg, (msg.(OperationReply)).OperationMessage.MsgId)
        } else if msg_type == reflect.TypeOf(EstablishRouteReplyMessage{}) {
            self.onOperationReply(msg, (msg.(EstablishRouteReplyMessage)).OperationReply.OperationMessage.MsgId)
        } else if msg_type == reflect.TypeOf(QueryConnectedPeersReply{}) {
            self.onOperationReply(msg, (msg.(QueryConnectedPeersReply)).OperationReply.OperationMessage.MsgId)
        }

        // TODO
    }
}



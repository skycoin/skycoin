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
    SendBack bool
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
    NewSendId uint32
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

type PhysicalMessage struct {
    ConnectedPeerPubKey cipher.PubKey
    Message rewriteableMessage
}

type EstablishedRoute struct {
    SendId          uint32
    ConnectedPeer   cipher.PubKey
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

type Node struct {
    // These are the messages contents received from other nodes in the network
    // The Node fills this channel
    MeshMessagesIn chan[]byte

    // These are the messages to be sent across physical links
    // The Node fills this channel
    MessagesOut chan PhysicalMessage

    // These are the messages received across physical links
    // The Node empties this channel
    MessagesIn chan PhysicalMessage

    // TODO: Make private
    Config NodeConfig
    EstablishedRoutesByIndex map[int]EstablishedRoute
    RetransmitQueue map[uuid.UUID]PhysicalMessage
    OperationsAwaitingReply map[uuid.UUID]chan rewriteableMessage

    // TODO: Randomize route indexes?
    NextRouteIdx uint32
    ForwardPeerIdsBySendId map[uint32]cipher.PubKey
    SendIdsBySecret map[string]uint32
    ForwardRewriteBySendId map[uint32]uint32

    Lock *sync.Mutex
}

type rewriteableMessage interface {
    Rewrite(newSendId uint32) rewriteableMessage
}

func (msg SendMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.SendId = newSendId
    return msg
}

func (msg OperationReply) Rewrite(newSendId uint32) rewriteableMessage {
    msg.SendId = newSendId
    return msg
}

func (msg EstablishRouteMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.SendId = newSendId
    return msg
}

func (msg EstablishRouteReplyMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.SendId = newSendId
    return msg
}

func (msg QueryConnectedPeersMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.SendId = newSendId
    return msg
}

func (msg QueryConnectedPeersReply) Rewrite(newSendId uint32) rewriteableMessage {
    msg.SendId = newSendId
    return msg
}

func (msg RouteRewriteMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.SendId = newSendId
    return msg
}

type RouteConfig struct {
    PeerPubKeys []cipher.PubKey
}

func NewNode(config NodeConfig) *Node {
    ret := &Node{}
    ret.Config = config
    ret.MessagesOut = make(chan PhysicalMessage, config.MessagesOutSize)
    ret.MessagesIn = make(chan PhysicalMessage, config.MessagesInSize)
    ret.MeshMessagesIn = make(chan[]byte)
    ret.Lock = &sync.Mutex{}
    ret.EstablishedRoutesByIndex = make(map[int]EstablishedRoute)
    ret.RetransmitQueue = make(map[uuid.UUID]PhysicalMessage)
    ret.OperationsAwaitingReply = make(map[uuid.UUID]chan rewriteableMessage)
    ret.NextRouteIdx = 1;
    ret.ForwardPeerIdsBySendId = make(map[uint32]cipher.PubKey)
    ret.SendIdsBySecret = make(map[string]uint32)
    ret.ForwardRewriteBySendId = make(map[uint32]uint32)
    return ret
}

func (self *Node) sendAndConfirmOperation(connected cipher.PubKey, operation_id uuid.UUID, msg rewriteableMessage) rewriteableMessage {
    reply_channel := make(chan rewriteableMessage)
    self.Lock.Lock()
    self.OperationsAwaitingReply[operation_id] = reply_channel
    outgoing := PhysicalMessage{connected, msg}
    self.RetransmitQueue[operation_id] = outgoing
    self.Lock.Unlock()

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
    var last_secret = ""

    for peer_idx, _ := range route.PeerPubKeys {
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
                        false,
                    },
                    establish_id,
                },
                route.PeerPubKeys[peer_idx + 1],
                time.Hour,
            }

        establish_reply := self.sendAndConfirmOperation(new_route.ConnectedPeer,
                                                        establish_id,
                                                        route_message)
        if establish_reply != nil {
            route_reply := establish_reply.(EstablishRouteReplyMessage)
            if peer_idx == 0 {
                new_route.SendId = route_reply.NewSendId
            } else {
                rewrite_id := uuid.NewV4()
                rewrite_message := 
                    RouteRewriteMessage{
                        OperationMessage {
                            Message {
                                new_route.SendId,
                                false,
                            },
                            rewrite_id,
                        },
                        last_secret,
                        route_reply.NewSendId,
                    }
                rewrite_reply := self.sendAndConfirmOperation(new_route.ConnectedPeer,
                                                                rewrite_id,
                                                                rewrite_message)
                if rewrite_reply == nil {
                    return
                }
            }
            last_secret = route_reply.Secret
        } else {
            return
        }
    }

    self.Lock.Lock();
    self.EstablishedRoutesByIndex[route_idx] = new_route
    self.Lock.Unlock();

    if(self.Config.RouteEstablishedCB != nil) {
        self.Config.RouteEstablishedCB(new_route)
    }
}

func (self *Node) onOperationReply(msg rewriteableMessage, operation_id uuid.UUID) {
    self.Lock.Lock()
    ch := self.OperationsAwaitingReply[operation_id]
    delete(self.OperationsAwaitingReply, operation_id)
    self.Lock.Unlock()

    ch <- msg
}

func (self *Node) onSendMessage(msg SendMessage) {
    self.MeshMessagesIn <- msg.Contents
}

func (self *Node) onRouteRequest(msg EstablishRouteMessage, peerFrom cipher.PubKey) {
    new_secret := fmt.Sprintf("%v", uuid.NewV4())

    self.Lock.Lock()
    route_id := self.NextRouteIdx
    self.NextRouteIdx += 1
    self.ForwardPeerIdsBySendId[route_id] = msg.ToPubKey
    self.SendIdsBySecret[new_secret] = route_id
    // TODO: Use duration hint
    self.Lock.Unlock()
    
    self.MessagesOut <- 
        PhysicalMessage {
            peerFrom,
            EstablishRouteReplyMessage{
                OperationReply {
                    OperationMessage {
                        Message {
                            // SendId
                            msg.SendId,
                            // SendBack
                            true,
                        },
                        msg.MsgId,
                    },
                    // Success
                    true,
                    // String
                    "",
                },
                // NewSendId
                route_id,
                // Secret
                new_secret,
            },
        }
}

func (self *Node) forwardMessage(SendId uint32, msg rewriteableMessage) bool {
    if SendId == 0 {
        return false
    }

    // Get routing info
    self.Lock.Lock()
    var rewrite_id uint32 = 0
    forward_peer, forward_exists := self.ForwardPeerIdsBySendId[SendId]
    if !forward_exists {
        logger.Debug("Asked to forward on unknown route %v", SendId)
        return false
    }
    rewrite_id, _ = self.ForwardRewriteBySendId[SendId]
    self.Lock.Unlock()

    // Special routing for RouteRewriteMessage
    msg_type := reflect.TypeOf(msg)
    if msg_type == reflect.TypeOf(RouteRewriteMessage{}) && rewrite_id == 0 {
        // Do not forward: interpret this message on this hop
        return false
    }
    
    rewritten_msg := msg.(rewriteableMessage) 

    rewritten_msg = rewritten_msg.Rewrite(rewrite_id)
    self.MessagesOut <- PhysicalMessage{forward_peer, rewritten_msg}

    return true
}

func (self *Node) onRewriteRequest(msg RouteRewriteMessage, peerFrom cipher.PubKey) {
    self.Lock.Lock()
    // Check secret
    sendId, exists := self.SendIdsBySecret[msg.Secret]
    if exists {
        delete(self.SendIdsBySecret, msg.Secret)

        self.ForwardRewriteBySendId[sendId] = msg.RewriteSendId

        self.MessagesOut <- 
            PhysicalMessage {
                peerFrom,
                OperationReply {
                    OperationMessage {
                        Message {
                            // SendId
                            msg.SendId,
                            // SendBack
                            true,
                        },
                        msg.MsgId,
                    },
                    true,
                    "",
                },
            }
    } else {
        self.MessagesOut <- 
            PhysicalMessage {
                peerFrom,
                OperationReply {
                    OperationMessage {
                        Message {
                            // SendId
                            msg.SendId,
                            // SendBack
                            true,
                        },
                        msg.MsgId,
                    },
                    false,
                    "Unknown secret",
                },
            }
    }
    self.Lock.Unlock()
}

func (self *Node) SendMessage(route_idx int, contents []byte) {
    route, exists := self.EstablishedRoutesByIndex[route_idx]
    if !exists {
        return
    }
    self.MessagesOut <- PhysicalMessage{route.ConnectedPeer, 
                                        SendMessage{Message{route.SendId, false}, contents}}
}

// Blocks
func (self *Node) Run() {
    // Retransmit loop
    if self.Config.RetransmitInterval > 0 {
        go func() {
            for {
                self.Lock.Lock()
                for _, outgoing := range self.RetransmitQueue {
                    self.MessagesOut <- outgoing
                }
                self.Lock.Unlock()
                time.Sleep(self.Config.RetransmitInterval)
            }
        }()
    }

    // Establish routes asynchronously
    for route_idx, route := range self.Config.Routes {
        go self.establishRoute(route_idx, route)
    }

    // Incoing messages loop
    for{
        msg_physical := <- self.MessagesIn
        msg := msg_physical.Message
        msg_type := reflect.TypeOf(msg)

        if msg_type == reflect.TypeOf(OperationReply{}) {
            self.onOperationReply(msg, (msg.(OperationReply)).OperationMessage.MsgId)
        } else if msg_type == reflect.TypeOf(EstablishRouteReplyMessage{}) {
            self.onOperationReply(msg, (msg.(EstablishRouteReplyMessage)).OperationReply.OperationMessage.MsgId)
        } else if msg_type == reflect.TypeOf(QueryConnectedPeersReply{}) {
            self.onOperationReply(msg, (msg.(QueryConnectedPeersReply)).OperationReply.OperationMessage.MsgId)
        } else if msg_type == reflect.TypeOf(SendMessage{}) {
            send_message := msg.(SendMessage)
            if !self.forwardMessage(send_message.SendId, msg) {
                self.onSendMessage(send_message)
            }
        } else if msg_type == reflect.TypeOf(EstablishRouteMessage{}) {
            establish_message := msg.(EstablishRouteMessage)
            if !self.forwardMessage(establish_message.SendId, msg) {
                self.onRouteRequest(establish_message, msg_physical.ConnectedPeerPubKey)
            }
        } else if msg_type == reflect.TypeOf(RouteRewriteMessage{}) {
            rewrite_message := msg.(RouteRewriteMessage)
            if !self.forwardMessage(rewrite_message.SendId, msg) {
                self.onRewriteRequest(rewrite_message, msg_physical.ConnectedPeerPubKey)
            }
        }
    }
}



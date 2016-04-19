package mesh

import (
    "time"
    "fmt"
    "sync"
    "reflect"
    "gopkg.in/op/go-logging.v1"
    "math"
    "math/rand"
)

import (
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/satori/go.uuid"
)

var logger = logging.MustGetLogger("node")

type Message struct {
    SendId uint32
    SendBack bool
    LastSendId uint32
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
    BackwardRewriteSendId uint32
    DurationHint time.Duration
}

type EstablishRouteReplyMessage struct {
    OperationReply
    NewSendId uint32
    // Secret be sent back in ActivateRouteMessage
    Secret string
}

type RouteRewriteMessage struct {
    OperationMessage
    // Secret from EstablishRouteReply
    Secret string
    RewriteSendId uint32
}

type PhysicalMessage struct {
    ConnectedPeerPubKey cipher.PubKey
    Contents []byte
}

type EstablishedRoute struct {
    SendId          uint32
    ConnectedPeer   cipher.PubKey
}

type RouteEstablishedCallback func(route EstablishedRoute)

type NodeConfig struct {
    MyPubKey cipher.PubKey
    MessagesOutSize int
    MessagesInSize int
    Routes []RouteConfig
    OperationTimeout time.Duration
    RetransmitInterval time.Duration
    RouteEstablishedCB RouteEstablishedCallback
}

type MeshMessage struct {
    SendId          uint32
    ConnectedPeer   cipher.PubKey
    Contents        []byte
}

type Node struct {
    // These are the messages contents received from other nodes in the network
    // The Node fills this channel
    MeshMessagesIn chan MeshMessage

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
    BackwardPeerIdsBySendId map[uint32]cipher.PubKey
    BackwardSendIdsBySendId map[uint32]uint32

    // TODO: Randomize route indexes?
   // NextRouteIdx uint32
    ForwardPeerIdsBySendId map[uint32]cipher.PubKey
    SendIdsBySecret map[string]uint32
    ForwardRewriteBySendId map[uint32]uint32

    Lock *sync.Mutex

    serializer *Serializer
}

type rewriteableMessage interface {
    Rewrite(newSendId uint32) rewriteableMessage
}

func (msg SendMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.LastSendId = msg.SendId
    msg.SendId = newSendId
    return msg
}

func (msg OperationReply) Rewrite(newSendId uint32) rewriteableMessage {
    msg.LastSendId = msg.SendId
    msg.SendId = newSendId
    return msg
}

func (msg EstablishRouteMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.LastSendId = msg.SendId
    msg.SendId = newSendId
    return msg
}

func (msg EstablishRouteReplyMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.LastSendId = msg.SendId
    msg.SendId = newSendId
    return msg
}

func (msg RouteRewriteMessage) Rewrite(newSendId uint32) rewriteableMessage {
    msg.LastSendId = msg.SendId
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
    ret.MeshMessagesIn = make(chan MeshMessage)
    ret.Lock = &sync.Mutex{}
    ret.EstablishedRoutesByIndex = make(map[int]EstablishedRoute)
    ret.RetransmitQueue = make(map[uuid.UUID]PhysicalMessage)
    ret.OperationsAwaitingReply = make(map[uuid.UUID]chan rewriteableMessage)
    ret.ForwardPeerIdsBySendId = make(map[uint32]cipher.PubKey)
    ret.SendIdsBySecret = make(map[string]uint32)
    ret.ForwardRewriteBySendId = make(map[uint32]uint32)
    ret.BackwardPeerIdsBySendId = make(map[uint32]cipher.PubKey)
    ret.BackwardSendIdsBySendId = make(map[uint32]uint32)

    ret.serializer = NewSerializer()
    ret.serializer.RegisterMessageForSerialization(messagePrefix{1}, SendMessage{})
    ret.serializer.RegisterMessageForSerialization(messagePrefix{2}, EstablishRouteMessage{})
    ret.serializer.RegisterMessageForSerialization(messagePrefix{3}, EstablishRouteReplyMessage{})
    ret.serializer.RegisterMessageForSerialization(messagePrefix{4}, RouteRewriteMessage{})
    ret.serializer.RegisterMessageForSerialization(messagePrefix{5}, OperationReply{})
    return ret
}

func init() {
}

func (self *Node) sendAndConfirmOperation(connected cipher.PubKey, operation_id uuid.UUID, msg rewriteableMessage) rewriteableMessage {
    reply_channel := make(chan rewriteableMessage)
    self.Lock.Lock()
    self.OperationsAwaitingReply[operation_id] = reply_channel
    outgoing := PhysicalMessage{connected, self.serializer.SerializeMessage(msg)}
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
    var prevSendId uint32 = 0

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
                        0,
                    },
                    establish_id,
                },
                route.PeerPubKeys[peer_idx + 1],
                prevSendId,
                time.Hour,
            }

        establish_reply := self.sendAndConfirmOperation(new_route.ConnectedPeer,
                                                        establish_id,
                                                        route_message)
        if establish_reply != nil {
            route_reply := establish_reply.(EstablishRouteReplyMessage)
            prevSendId = route_reply.NewSendId
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
                                0,
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

func (self *Node) onSendMessage(msg MeshMessage) {
    self.MeshMessagesIn <- msg
}

func (self *Node) onRouteRequest(msg EstablishRouteMessage, peerFrom cipher.PubKey) {
    new_secret := fmt.Sprintf("%v", uuid.NewV4())

    var route_id uint32 = 0
    self.Lock.Lock()
    if len(self.ForwardPeerIdsBySendId) == (math.MaxUint32 - 2) {
        panic("Too many routes in table to generate a new ID")
    }
    for {
        route_id = rand.Uint32()
        if route_id == 0 {
            continue
        }
        _, exists := self.ForwardPeerIdsBySendId[route_id]
        if !exists {
            break
        }
    }

    self.ForwardPeerIdsBySendId[route_id] = msg.ToPubKey
    self.SendIdsBySecret[new_secret] = route_id

    self.BackwardPeerIdsBySendId[route_id] = peerFrom
    self.BackwardSendIdsBySendId[route_id] = msg.BackwardRewriteSendId

    // TODO: Use duration hint
    self.Lock.Unlock()
    
    self.MessagesOut <- 
        PhysicalMessage {
            peerFrom,
            self.serializer.SerializeMessage(EstablishRouteReplyMessage{
                OperationReply {
                    OperationMessage {
                        Message {
                            // SendId
                            msg.LastSendId,
                            // SendBack
                            true,
                            // LastSendId
                            0,
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
            }),
        }
}

func (self *Node) forwardMessage(SendId uint32, SendBack bool, msg rewriteableMessage) bool {
    if SendId == 0 {
        return false
    }

    var rewrite_id uint32 = 0
    var rewrite_peer cipher.PubKey

    if !SendBack {
        // Get routing info
        self.Lock.Lock()
        var forward_exists bool = false
        rewrite_peer, forward_exists = self.ForwardPeerIdsBySendId[SendId]
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
    } else {
        // Which peer does it go to? 
        var back_exists bool = false
        rewrite_peer, back_exists = self.BackwardPeerIdsBySendId[SendId]
        if !back_exists  {
            logger.Debug("Asked to backward route on unknown route %v", SendId)
            return false
        }

        var id_exists bool = false
        rewrite_id, id_exists = self.BackwardSendIdsBySendId[SendId]
        if !id_exists {
            panic("Internal inconsistency: BackwardPeerIdsBySendId has key but BackwardSendIdsBySendId doesn't")
        }
    }
    
    rewritten_msg := msg.(rewriteableMessage) 

    rewritten_msg = rewritten_msg.Rewrite(rewrite_id)
    self.MessagesOut <- PhysicalMessage{rewrite_peer, self.serializer.SerializeMessage(rewritten_msg)}

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
                self.serializer.SerializeMessage(OperationReply {
                    OperationMessage {
                        Message {
                            // SendId
                            msg.LastSendId,
                            // SendBack
                            true,
                            // LastSendId
                            0,
                        },
                        msg.MsgId,
                    },
                    true,
                    "",
                }),
            }
    } else {
        self.MessagesOut <- 
            PhysicalMessage {
                peerFrom,
                self.serializer.SerializeMessage(OperationReply {
                    OperationMessage {
                        Message {
                            // SendId
                            msg.LastSendId,
                            // SendBack
                            true,
                            // LastSendId
                            0,
                        },
                        msg.MsgId,
                    },
                    false,
                    "Unknown secret",
                }),
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
                                        self.serializer.SerializeMessage(SendMessage{Message{route.SendId, false, 0}, contents})}
}

func (self *Node) SendReply(to MeshMessage, contents []byte) {
    self.MessagesOut <- PhysicalMessage{to.ConnectedPeer, 
                                        self.serializer.SerializeMessage(SendMessage{Message{to.SendId, true, 0}, contents})}
}

// Blocks
func (self *Node) Run() {
    // Retransmit loop
    if self.Config.RetransmitInterval > 0 {
        // TODO: Re-enable retransmits when they are needed, or maybe replace with store and forward
        panic("Retransmits are disabled for now, having only been partially implemented. ")
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

    // Incoming messages loop
    for{
        msg_physical := <- self.MessagesIn
        msg, deserialize_error := self.serializer.UnserializeMessage(msg_physical.Contents)
        if deserialize_error != nil {
            logger.Debug("Deserialization error %v\n", deserialize_error)
            continue
        }
        msg_type := reflect.TypeOf(msg) 

        if msg_type == reflect.TypeOf(OperationReply{}) {
            reply := msg.(OperationReply)
            if !self.forwardMessage(reply.SendId, reply.SendBack, reply) {
                self.onOperationReply(reply, (msg.(OperationReply)).OperationMessage.MsgId)
            }
        } else if msg_type == reflect.TypeOf(EstablishRouteReplyMessage{}) {
            reply := msg.(EstablishRouteReplyMessage)
            if !self.forwardMessage(reply.SendId, reply.SendBack, reply) {
                self.onOperationReply(reply, (msg.(EstablishRouteReplyMessage)).OperationReply.OperationMessage.MsgId)
            }
        } else if msg_type == reflect.TypeOf(SendMessage{}) {
            send_message := msg.(SendMessage)
            if !self.forwardMessage(send_message.SendId, send_message.SendBack, send_message) {
                self.onSendMessage(MeshMessage{send_message.LastSendId, msg_physical.ConnectedPeerPubKey, send_message.Contents})
            }
        } else if msg_type == reflect.TypeOf(EstablishRouteMessage{}) {
            establish_message := msg.(EstablishRouteMessage)
            if !self.forwardMessage(establish_message.SendId, establish_message.SendBack, establish_message) {
                self.onRouteRequest(establish_message, msg_physical.ConnectedPeerPubKey)
            }
        } else if msg_type == reflect.TypeOf(RouteRewriteMessage{}) {
            rewrite_message := msg.(RouteRewriteMessage)
            if !self.forwardMessage(rewrite_message.SendId, rewrite_message.SendBack, rewrite_message) {
                self.onRewriteRequest(rewrite_message, msg_physical.ConnectedPeerPubKey)
            }
        }
    }
}



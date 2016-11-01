#mesh

package domain
--------------------------

package domain // import "github.com/skycoin/src/mesh/domain"


== struct

type AddNodeMessage struct { ... }
#    Add a new node to the network

type DeleteRouteMessage struct { ... }
#    Deletes the route as it passes thru it

type LocalRoute struct { ... }

type MeshMessage struct { ... }

type MessageBase struct { ... }
#   If RouteId is unknown, but not cipher.PubKey{}, then the message should be received here the RouteId can be used to reply back thru the route
#   Fields must be public (capital first letter) for encoder


type MessageID uuid.UUID

type MessageUnderAssembly struct { ... }

type NodeConfig struct { ... }

type RefreshRouteMessage struct { ... }
#       Refreshes the route as it passes thru it

type ReplyTo struct { ... }

type Route struct { ... }
# Forward should never be cipher.PubKey{}
# time.Unix(0,0) means it lives forever

type RouteID uuid.UUID

type SetRouteMessage struct { ... }

type SetRouteReply struct { ... }
#  This allows ExtendRoute() to block so that messages aren't lost while a route is not yet established

type UserMessage struct { ... }


package gui
-----------------

package gui // import "github.com/skycoin/src/mesh/gui"

== struct

type ConfigWithID struct { ... }
#    struct for nodeAddTransportHandler

type TransportWithID struct { ... }
#    struct for nodeRemoveTransportHandler



== function

func LaunchWebInterface(host, staticDir string, nm *nodemanager.NodeManager) error
#     Begins listening on http://$host, for enabling remote web access Does NOT use HTTPS

func NewGUIMux(appLoc string, nm *nodemanager.NodeManager) *http.ServeMux
#     Creates an http.ServeMux with handlers registered

func RegisterNodeManagerHandlers(mux *http.ServeMux, nm *nodemanager.NodeManager)
#     RegisterNodeManagerHandlers - create routes for NodeManager






node
-----------------



type Node struct {
	// Has unexported fields.
}

func NewNode(config domain.NodeConfig) (*Node, error)

func (self *Node) AddRoute(routeID domain.RouteID, toPeer cipher.PubKey) error
#   toPeer must be the public key of a connected peer

func (self *Node) AddTransport(transportNode transport.ITransport)

func (self *Node) Close() error

func (self *Node) ConnectedToPeer(peer cipher.PubKey) bool

func (self *Node) DeleteRoute(routeID domain.RouteID) (err error)

func (self *Node) ExtendRoute(routeID domain.RouteID, toPeer cipher.PubKey, timeout time.Duration) (err error)
#    toPeer must be the public key of a peer connected to the current last node
#    in this route Blocks until the set route reply is received or the timeout is
#    reached



func (self *Node) GetConfig() domain.NodeConfig

func (self *Node) GetConnectedPeers() []cipher.PubKey

func (self *Node) GetRouteLastConfirmed(routeID domain.RouteID) (time.Time, error)

func (self *Node) GetTransports() []transport.ITransport

func (self *Node) RemoveTransport(transport transport.ITransport)

func (self *Node) SendMessageBackThruRoute(replyTo domain.ReplyTo, contents []byte) error
#    Blocks until message is confirmed received

func (self *Node) SendMessageThruRoute(routeId domain.RouteID, contents []byte) error
#    Blocks until message is confirmed received

func (self *Node) SendMessageToPeer(toPeer cipher.PubKey, contents []byte) (error, domain.RouteID)

func (self *Node) SetReceiveChannel(received chan domain.MeshMessage)



package connection
------------------------

package connection // import "github.com/skycoin/src/mesh/node/connection"

var ConnectionManager Connection

func (self *Connection) DeserializeMessage(msg []byte) (interface{}, error)

func (self *Connection) FragmentMessage(fullContents []byte, toPeer cipher.PubKey, transport transport.ITransport, base domain.MessageBase) []domain.UserMessage

func (self *Connection) GetMaximumContentLength(toPeer cipher.PubKey, transport transport.ITransport) uint64

func (self *Connection) SerializeMessage(msg interface{}) []byte




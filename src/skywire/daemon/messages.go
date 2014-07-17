package daemon

import (
	//"encoding/binary"
	//"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/skywire/daemon/pex"
	//"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/skywire/lib/gnet"
	"log"
	"math/rand"
	//"net"
)

//used to detect self connection; replace with public key
var MirrorConstant uint32 = rand.Uint32()

//Daemon on channel 0
//The channel 0 service manages exposing service metainformation and
//server setup and teardown
type DaemonService struct {
	Daemon         *Daemon
	Service        *gnet.Service //service for daemon
	ServiceManager *gnet.ServiceManager
}

// TODO:
// - add request packet for service list
// - add connection packet for service
// - move into daemon

func NewDaemonService(sm *gnet.ServiceManager, daemon *Daemon) *DaemonService {
	var swd DaemonService
	swd.Daemon = daemon
	swd.ServiceManager = sm
	//associate service with channel 0
	swd.Service = sm.AddService(
		[]byte("Skywire Daemon"),
		[]byte("{service:'Skywire Daemon',version=0"),
		0, &swd)

	return &swd
}

//move to daemon

func (sd *DaemonService) OnConnect(c *gnet.Connection) {
	fmt.Printf("SkywireDaemon: OnConnect, addr= %s \n", c.Addr())
}

func (sd *DaemonService) OnDisconnect(c *gnet.Connection) {
	fmt.Printf("SkywireDaemon: OnDisconnect, addr= %s \n", c.Addr())
}

func (sd *DaemonService) RegisterMessages(d *gnet.Dispatcher) {
	fmt.Printf("SkywireDaemon: RegisterMessages \n")

	var messageMap map[string](interface{}) = map[string](interface{}){
		//put messages here
		"INTR": IntroductionMessage{},
		"GETP": GetPeersMessage{},
		"GIVP": GivePeersMessage{},
		"PING": PingMessage{},
		"PONG": PongMessage{},
		"SCON": ServiceConnectMessage{},
	}
	d.RegisterMessages(messageMap)
}

// Compact representation of IP:Port
// Addresses in future can be darknet addresses or IPv6, should be strings
type IPAddr struct {
	Addr []byte // as string
}

// Returns an IPAddr from an ip:port string.  If ipv6 or invalid, error is
// returned
func NewIPAddr(addr string) (ipaddr IPAddr, err error) {
	return IPAddr{Addr: []byte(addr)}, nil
}

// Returns IPAddr as "ip:port"
func (self IPAddr) String() string {
	return string(self.Addr)
}

// Messages that perform an action when received must implement this interface.
// Process() is called after the message is pulled off of messageEvent channel.
// Messages should place themselves on the messageEvent channel in their
// Handle() method required by gnet.
//type AsyncMessage interface {
//	Process(d *Daemon)
//}

// Sent to request peers
type GetPeersMessage struct {
	c *gnet.MessageContext `enc:"-"`
}

func NewGetPeersMessage() *GetPeersMessage {
	return &GetPeersMessage{}
}

func (self *GetPeersMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	s := state.(*DaemonService)
	d := s.Daemon

	if d.Peers.Config.Disabled {
		return nil
	}
	peers := d.Peers.Peers.Peerlist.RandomPublic(d.Peers.Config.ReplyCount)
	if len(peers) == 0 {
		logger.Debug("We have no peers to send in reply")
		return nil
	}
	m := NewGivePeersMessage(peers)

	s.Service.Send(self.c.Conn, m)

	return nil
}

// Sent in response to GetPeersMessage
type GivePeersMessage struct {
	Peers []IPAddr
}

// []*pex.Peer is converted to []IPAddr for binary transmission
func NewGivePeersMessage(peers []*pex.Peer) *GivePeersMessage {
	ipaddrs := make([]IPAddr, 0, len(peers))
	for _, ps := range peers {
		ipaddr, err := NewIPAddr(ps.Addr)
		if err != nil {
			logger.Warning("GivePeersMessage skipping address %s", ps.Addr)
			logger.Warning(err.Error())
			continue
		}
		ipaddrs = append(ipaddrs, ipaddr)
	}
	return &GivePeersMessage{Peers: ipaddrs}
}

// GetPeers is required by the pex.GivePeersMessage interface.
// It returns the peers contained in the message as an array of "ip:port"
// strings.
func (self *GivePeersMessage) GetPeers() []string {
	peers := make([]string, len(self.Peers))
	for i, ipaddr := range self.Peers {
		peers[i] = ipaddr.String()
	}
	return peers
}

func (self *GivePeersMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	s := state.(*DaemonService)
	d := s.Daemon

	if d.Peers.Config.Disabled {
		return nil
	}
	peers := self.GetPeers()
	if len(peers) != 0 {
		logger.Debug("Got these peers via PEX:")
		for _, p := range peers {
			logger.Debug("\t%s", p)
		}
	}
	d.Peers.Peers.AddPeers(peers)
	return nil
}

// An IntroductionMessage is sent on first connect by both parties
type IntroductionMessage struct {
	// Mirror is a random value generated on client startup that is used
	// to identify self-connections
	Mirror uint32
	// Port is the port that this client is listening on
	Port uint16
	// Our client version
	Version int32

	// We validate the message in Handle() and cache the result for Process()
	valid bool `enc:"-"` // skip it during encoding
}

func NewIntroductionMessage(mirror uint32, version int32,
	port uint16) *IntroductionMessage {
	return &IntroductionMessage{
		Mirror:  mirror,
		Version: version,
		Port:    port,
	}
}

// Note :in future, address will be pubkey or ip:port

// Responds to an gnet.Pool event. We implement Handle() here because we
// need to control the DisconnectReason sent back to gnet.  We still implement
// Process(), where we do modifications that are not threadsafe
func (self *IntroductionMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	s := state.(*DaemonService)
	d := s.Daemon

	var err error

	addr := mc.Conn.Addr()
	// Disconnect if this is a self connection (we have the same mirror value)
	if self.Mirror == MirrorConstant {
		logger.Info("Remote mirror value %v matches ours", self.Mirror)
		d.Pool.Disconnect(mc.Conn, DisconnectSelf)
		err = DisconnectSelf
	}
	// Disconnect if not running the same version
	if self.Version != d.Config.Version {
		logger.Info("%s has different version %d. Disconnecting.",
			addr, self.Version)

		//diconnect whole peer, not just service
		d.Pool.Disconnect(mc.Conn, DisconnectInvalidVersion)
		err = DisconnectInvalidVersion
	} else {
		logger.Info("%s verified for version %d", addr, self.Version)
	}

	if err != nil {
		return nil
	}
	//weird condition if same client connects/reconnects
	delete(d.ExpectingIntroductions, mc.Conn.Addr())

	// Add the remote peer with their chosen listening port
	a := mc.Conn.Addr()
	ip, _, err := SplitAddr(a)
	if err != nil {
		// This should never happen, but the program should still work if it
		// does.
		logger.Error("Invalid Addr() for connection: %s", a)
		d.Pool.Disconnect(mc.Conn, DisconnectOtherError)
		return nil
	}

	_, err = d.Peers.Peers.AddPeer(fmt.Sprintf("%s:%d", ip, self.Port))
	if err != nil {
		logger.Error("Failed to add peer: %v", err)
	}
	return nil
}

// Sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
}

func (self *PingMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	s := state.(*DaemonService)
	//d := s.Daemon

	logger.Debug("Reply to ping from %s", mc.Conn.Addr())
	s.Service.Send(mc.Conn, &PongMessage{})
	return nil
}

// Sent in reply to a PingMessage.  No action is taken when this is received.
type PongMessage struct {
}

func (self *PongMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	//s := state.(*DaemonService)
	//d := s.Daemon

	logger.Debug("Received pong from %s", mc.Conn.Addr())
	return nil
}

func (self *Daemon) ConnectToService(Conn *gnet.Connection, Service *gnet.Service) {
	var Id [20]byte
	copy(Id[0:20], Service.Id[0:20])

	scm := ServiceConnectMessage{}
	scm.Originating = 1
	scm.ServiceIdentifer = Id
	scm.OriginChannel = Service.Channel
	scm.RemoteChannel = 0 //unknown

	self.Service.Send(Conn, &scm) //channel 0
}

type ServiceConnectMessage struct {
	//peer originating requests sets to 1
	//peer responding sets to 0
	Originating uint32

	ServiceIdentifer [20]byte //20 byte hash, identifying service
	OriginChannel    uint16   //channel of initiator
	RemoteChannel    uint16   //channel of responder

	ErrorMessage []byte //fail if error len != 0
}

func (self *ServiceConnectMessage) Handle(context *gnet.MessageContext,
	state interface{}) error {
	server := state.(*DaemonService) //service server state

	if len(self.ServiceIdentifer) > 140 {
		log.Printf("ServiceConnectMessage: Error service identifer exceeds 140 bytes, ignored")
		return nil
	}

	//message from remote for connection
	if self.Originating == 1 {

		service := server.ServiceManager.ServiceById(self.ServiceIdentifer)

		if service != nil {
			//service exists, send success message
			var scm ServiceConnectMessage
			scm.OriginChannel = self.OriginChannel
			scm.RemoteChannel = service.Channel
			scm.Originating = 0
			scm.ErrorMessage = []byte("")
			server.Service.Send(context.Conn, &scm) //channel 0
			//trigger connection event
			service.ConnectionEvent(context.Conn, self.OriginChannel)
			return nil
		}

		if server == nil {
			//server does not exist
			log.Printf("ServiceConnectMessage: no service with id exists \n")

			//failure message
			var scm ServiceConnectMessage
			scm.OriginChannel = self.OriginChannel
			scm.RemoteChannel = 0
			scm.Originating = 0
			scm.ErrorMessage = []byte("no service with id exists")
			server.Service.Send(context.Conn, &scm) //channel 0
			return nil
		}

	}
	//message reponse from remote for connection
	if self.Originating == 0 {
		if len(self.ErrorMessage) != 0 {
			log.Printf("Service Connection Failed:addr= %s, LocalChannel= %d, Remotechannel= %d \n",
				context.Conn.Addr(), self.OriginChannel, self.RemoteChannel)
			return nil
		}

		service, ok := server.ServiceManager.Services[self.RemoteChannel]

		if ok == false {
			log.Printf("service does not exist on local, LocalChannel= %d from addr= %s \n",
				self.OriginChannel, context.Conn.Addr())
		}

		service.ConnectionEvent(context.Conn, self.RemoteChannel)
		return nil
	}
	return nil
}

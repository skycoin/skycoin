package daemon

import (
	//"encoding/binary"
	//"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/aether/daemon/pex"
	//"github.com/skycoin/skycoin/src/util"
	"log"
	"math/rand"

	gnet "github.com/skycoin/skycoin/src/aether/gnet"
	//"net"
)

// MirrorConstant used to detect self connection; replace with public key
var MirrorConstant = rand.Uint32()

// DaemonService Daemon on channel 0
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

// NewDaemonService creates daemon service
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

// OnConnect callback function will be invoked when connection established
func (sd *DaemonService) OnConnect(c *gnet.Connection) {
	fmt.Printf("SkywireDaemon: OnConnect, addr= %s \n", c.Addr())
}

// OnDisconnect callback function will be invocked when connection breakdown
func (sd *DaemonService) OnDisconnect(c *gnet.Connection) {
	fmt.Printf("SkywireDaemon: OnDisconnect, addr= %s \n", c.Addr())
}

// RegisterMessages registers messages
func (sd *DaemonService) RegisterMessages(d *gnet.Dispatcher) {
	fmt.Printf("SkywireDaemon: RegisterMessages \n")

	var messageMap = map[string](interface{}){
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

// IPAddr compact representation of IP:Port
// Addresses in future can be darknet addresses or IPv6, should be strings
type IPAddr struct {
	Addr []byte // as string
}

// NewIPAddr returns an IPAddr from an ip:port string.  If ipv6 or invalid, error is
// returned
func NewIPAddr(addr string) (ipaddr IPAddr, err error) {
	return IPAddr{Addr: []byte(addr)}, nil
}

// String returns IPAddr as "ip:port"
func (ipa IPAddr) String() string {
	return string(ipa.Addr)
}

// Messages that perform an action when received must implement this interface.
// Process() is called after the message is pulled off of messageEvent channel.
// Messages should place themselves on the messageEvent channel in their
// Handle() method required by gnet.
//type AsyncMessage interface {
//	Process(d *Daemon)
//}

// GetPeersMessage sent to request peers
type GetPeersMessage struct {
	c *gnet.MessageContext `enc:"-"`
}

// NewGetPeersMessage create GetPeersMessage
func NewGetPeersMessage() *GetPeersMessage {
	return &GetPeersMessage{}
}

// Handle process message
func (gpm *GetPeersMessage) Handle(mc *gnet.MessageContext,
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

	s.Service.Send(gpm.c.Conn, m)

	return nil
}

// GivePeersMessage Sent in response to GetPeersMessage
type GivePeersMessage struct {
	Peers []IPAddr
}

// NewGivePeersMessage []*pex.Peer is converted to []IPAddr for binary transmission
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
func (gpm *GivePeersMessage) GetPeers() []string {
	peers := make([]string, len(gpm.Peers))
	for i, ipaddr := range gpm.Peers {
		peers[i] = ipaddr.String()
	}
	return peers
}

// Handle process message
func (gpm *GivePeersMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	s := state.(*DaemonService)
	d := s.Daemon

	if d.Peers.Config.Disabled {
		return nil
	}
	peers := gpm.GetPeers()
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

// NewIntroductionMessage creates introduction message
func NewIntroductionMessage(mirror uint32, version int32,
	port uint16) *IntroductionMessage {
	return &IntroductionMessage{
		Mirror:  mirror,
		Version: version,
		Port:    port,
	}
}

// Note :in future, address will be pubkey or ip:port

// Handle responds to an gnet.Pool event. We implement Handle() here because we
// need to control the DisconnectReason sent back to gnet.  We still implement
// Process(), where we do modifications that are not threadsafe
func (im *IntroductionMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	s := state.(*DaemonService)
	d := s.Daemon

	var err error

	addr := mc.Conn.Addr()
	// Disconnect if this is a self connection (we have the same mirror value)
	if im.Mirror == MirrorConstant {
		logger.Info("Remote mirror value %v matches ours", im.Mirror)
		d.Pool.Disconnect(mc.Conn, ErrDisconnectSelf)
		err = ErrDisconnectSelf
	}
	// Disconnect if not running the same version
	if im.Version != d.Config.Version {
		logger.Info("%s has different version %d. Disconnecting.",
			addr, im.Version)

		//diconnect whole peer, not just service
		d.Pool.Disconnect(mc.Conn, ErrDisconnectInvalidVersion)
		err = ErrDisconnectInvalidVersion
	} else {
		logger.Info("%s verified for version %d", addr, im.Version)
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
		d.Pool.Disconnect(mc.Conn, ErrDisconnectOtherError)
		return nil
	}

	_, err = d.Peers.Peers.AddPeer(fmt.Sprintf("%s:%d", ip, im.Port))
	if err != nil {
		logger.Error("Failed to add peer: %v", err)
	}
	return nil
}

// PingMessage sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
}

// Handle process message
func (pm *PingMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	s := state.(*DaemonService)
	//d := s.Daemon

	logger.Debug("Reply to ping from %s", mc.Conn.Addr())
	s.Service.Send(mc.Conn, &PongMessage{})
	return nil
}

// PongMessage sent in reply to a PingMessage.  No action is taken when this is received.
type PongMessage struct {
}

// Handle process message
func (pm *PongMessage) Handle(mc *gnet.MessageContext,
	state interface{}) error {
	//s := state.(*DaemonService)
	//d := s.Daemon

	logger.Debug("Received pong from %s", mc.Conn.Addr())
	return nil
}

// ConnectToService connection to service
func (dm *Daemon) ConnectToService(Conn *gnet.Connection, Service *gnet.Service) {
	var ID [20]byte
	copy(ID[0:20], Service.Id[0:20])

	scm := ServiceConnectMessage{}
	scm.Originating = 1
	scm.ServiceIdentifer = ID
	scm.OriginChannel = Service.Channel
	scm.RemoteChannel = 0 //unknown

	dm.Service.Send(Conn, &scm) //channel 0
}

// ServiceConnectMessage service connect message
type ServiceConnectMessage struct {
	//peer originating requests sets to 1
	//peer responding sets to 0
	Originating uint32

	ServiceIdentifer [20]byte //20 byte hash, identifying service
	OriginChannel    uint16   //channel of initiator
	RemoteChannel    uint16   //channel of responder

	ErrorMessage []byte //fail if error len != 0
}

// Handle process message
func (scm *ServiceConnectMessage) Handle(context *gnet.MessageContext,
	state interface{}) error {
	server := state.(*DaemonService) //service server state

	if len(scm.ServiceIdentifer) > 140 {
		log.Printf("ServiceConnectMessage: Error service identifer exceeds 140 bytes, ignored")
		return nil
	}

	//message from remote for connection
	if scm.Originating == 1 {

		service := server.ServiceManager.ServiceById(scm.ServiceIdentifer)

		if service != nil {
			//service exists, send success message
			var msg ServiceConnectMessage
			msg.OriginChannel = scm.OriginChannel
			msg.RemoteChannel = service.Channel
			msg.Originating = 0
			msg.ErrorMessage = []byte("")
			server.Service.Send(context.Conn, &msg) //channel 0
			//trigger connection event
			service.ConnectionEvent(context.Conn, scm.OriginChannel)
			return nil
		}

		if server == nil {
			//server does not exist
			log.Printf("ServiceConnectMessage: no service with id exists \n")

			//failure message
			var msg ServiceConnectMessage
			msg.OriginChannel = scm.OriginChannel
			msg.RemoteChannel = 0
			msg.Originating = 0
			msg.ErrorMessage = []byte("no service with id exists")
			server.Service.Send(context.Conn, &msg) //channel 0
			return nil
		}

	}
	//message response from remote for connection
	if scm.Originating == 0 {
		if len(scm.ErrorMessage) != 0 {
			log.Printf("Service Connection Failed:addr= %s, LocalChannel= %d, Remotechannel= %d \n",
				context.Conn.Addr(), scm.OriginChannel, scm.RemoteChannel)
			return nil
		}

		service, ok := server.ServiceManager.Services[scm.RemoteChannel]

		if ok == false {
			log.Printf("service does not exist on local, LocalChannel= %d from addr= %s \n",
				scm.OriginChannel, context.Conn.Addr())
		}

		service.ConnectionEvent(context.Conn, scm.RemoteChannel)
		return nil
	}
	return nil
}

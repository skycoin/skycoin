package daemon

import (
	"errors"
	"fmt"

	"log"
	"net"
	"strconv"
	"strings"
	"time"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/util"

	"github.com/skycoin/skycoin/src/aether/dht"       //dht library
	gnet "github.com/skycoin/skycoin/src/aether/gnet" //use local gnet
)

/*
	Todo:
	- give each daemon a pubkey and address
*/

/*
	Problems:
	- DHT Does not appear to be used anywhere!?
	- why isnt DHT being used for peer lookups?

*/

/*
	Why does daemon exist?
	- manage blacklists
	- DHT
	- peer exchange

	Can Daemon be eliminated?
	- would just be individual services
	- is there any reason that daemons should connect to each other? No
	- only advantage is single end-point for TCP/IP

	Future:
	- let daemon handle finding peers
	- services should not go through daemon?

	Right now
	- existing system works...
	- just get it working to minimum degree
	Todo:
	- peer list store applications

	Just gut everything
	- rip out peer list
	- have DHT query and callback for peer exchange
	- make connection to service in introduction message
*/
var (
	// ErrDisconneInvalidVersion invalid version error
	ErrDisconnectInvalidVersion gnet.DisconnectReason = errors.New("Invalid version")
	// ErrDisconnectIntroductionTimeout version time out
	ErrDisconnectIntroductionTimeout gnet.DisconnectReason = errors.New("Version timeout")
	// ErrDisconnectIsBlacklisted connection is blacklisted
	ErrDisconnectIsBlacklisted gnet.DisconnectReason = errors.New("Blacklisted")
	// ErrDisconnectSelf self connection
	ErrDisconnectSelf gnet.DisconnectReason = errors.New("Self connect")
	// ErrDisconnectConnectedTwice connect twice
	ErrDisconnectConnectedTwice gnet.DisconnectReason = errors.New("Already connected")

	ErrDisconnectOtherError gnet.DisconnectReason = errors.New("Incomprehensible error")

	// BlacklistOffenses a peer when they get disconnected for these
	// DisconnectReasons
	BlacklistOffenses = map[gnet.DisconnectReason]time.Duration{
		ErrDisconnectSelf:                   time.Hour * 24,
		ErrDisconnectIntroductionTimeout:    time.Hour,
		gnet.DisconnectInvalidMessageLength: time.Hour * 8,
		gnet.DisconnectMalformedMessage:     time.Hour * 8,
		gnet.DisconnectUnknownMessage:       time.Hour * 8,
	}

	logger = logging.MustGetLogger("skycoin.daemon")
)

// Config Subsystem configurations
type Config struct {
	Daemon DaemonConfig
	Peers  PeersConfig
	DHT    dht.DHTConfig //useless after config!
}

// NewConfig Returns a Config with defaults set
func NewConfig() Config {
	return Config{
		Daemon: NewDaemonConfig(),
		Peers:  NewPeersConfig(),
		DHT:    dht.NewDHTConfig(),
	}
}

func (cfg *Config) preprocess() Config {
	config := *cfg
	if config.Daemon.LocalhostOnly {
		config.Daemon.Address = LocalhostIP()
		config.DHT.Disabled = true
	}

	config.DHT.Port = config.Daemon.Port

	if config.Daemon.DisableNetworking {
		config.Peers.Disabled = true
		config.DHT.Disabled = true
		config.Daemon.DisableIncomingConnections = true
		config.Daemon.DisableOutgoingConnections = true
	} else {
		if config.Daemon.DisableIncomingConnections {
			logger.Info("Incoming connections are disabled.")
		}
		if config.Daemon.DisableOutgoingConnections {
			logger.Info("Outgoing connections are disabled.")
		}
	}
	return config
}

// DaemonConfig Configuration for the Daemon
type DaemonConfig struct {
	// Application version. TODO -- manage version better
	Version int32
	// IP Address to serve on. Leave empty for automatic assignment
	Address string
	// TCP/UDP port for connections and DHT
	Port int
	// Directory where application data is stored
	DataDirectory string
	// How often to check and initiate an outgoing connection if needed
	OutgoingRate time.Duration
	// How often to re-attempt to fill any missing private (aka required)
	// connections
	PrivateRate time.Duration
	// Number of outgoing connections to maintain
	OutgoingMax int
	// Maximum number of connections to try at once
	PendingMax int
	// How long to wait for a version packet
	IntroductionWait time.Duration
	// How often to check for peers that have decided to stop communicating
	CullInvalidRate time.Duration
	// Disable all networking activity
	DisableNetworking bool
	// Don't make outgoing connections
	DisableOutgoingConnections bool
	// Don't allow incoming connections
	DisableIncomingConnections bool
	// Run on localhost and only connect to localhost peers
	LocalhostOnly bool
}

// NewDaemonConfig create DaemonConfig
func NewDaemonConfig() DaemonConfig {
	return DaemonConfig{
		Version:                    3,
		Address:                    "",
		Port:                       6677,
		OutgoingRate:               time.Second * 5,
		PrivateRate:                time.Second * 5,
		OutgoingMax:                8,
		PendingMax:                 16, //for pex
		IntroductionWait:           time.Second * 30,
		CullInvalidRate:            time.Second * 3,
		DisableNetworking:          false,
		DisableOutgoingConnections: false, //makes random connections to new peers
		DisableIncomingConnections: false,
		LocalhostOnly:              false,
	}
}

// Daemon Stateful properties of the daemon
type Daemon struct {
	// Daemon configuration
	Config DaemonConfig

	// Components
	Pool  *gnet.ConnectionPool //what does this do
	Peers *Peers

	DHT            *dht.DHT
	ServiceManager *gnet.ServiceManager //service manager for pool
	Service        *gnet.Service        //base service for daemon

	// Separate index of outgoing connections
	OutgoingConnections map[string]*gnet.Connection //deprecate?
	// Number of connections waiting to be formed or timeout
	pendingConnections map[string]([]*gnet.Service)
	// Keep track of unsolicited clients who should notify us of their version
	ExpectingIntroductions map[string]time.Time
	// Connection failure events
	connectionErrors chan ConnectionError
}

// NewDaemon returns a Daemon with primitives allocated
func NewDaemon(config Config) *Daemon {
	config = config.preprocess()
	// c.DHT.address = c.Daemon.Address
	d := &Daemon{
		Config: config.Daemon,
		Peers:  NewPeers(config.Peers),
		ExpectingIntroductions: make(map[string]time.Time),

		// TODO -- if there are performance problems from blocking chans,
		// Its because we are connecting to more things than OutgoingMax
		// if we have private peers

		connectionErrors: make(chan ConnectionError,
			config.Daemon.OutgoingMax),
		OutgoingConnections: make(map[string]*gnet.Connection,
			config.Daemon.OutgoingMax),
		pendingConnections: make(map[string]([]*gnet.Service),
			config.Daemon.PendingMax),
	}
	d.Peers.Init()

	if config.DHT.Disabled == false {
		d.DHT = dht.NewDHT(config.DHT)
		d.DHT.Init()
	}

	//gnet set connection pool
	gnetConfig := gnet.NewConfig()
	gnetConfig.Port = uint16(d.Config.Port) //set listening port
	gnetConfig.Address = d.Config.Address
	d.Pool = gnet.NewConnectionPool(gnetConfig)

	//service manager
	d.ServiceManager = gnet.NewServiceManager(d.Pool)
	ds := NewDaemonService(d.ServiceManager, d)
	d.Service = ds.Service

	return d
}

// ConnectEvent generated when a client connects
type ConnectEvent struct {
	Addr      string
	Solicited bool
}

// ConnectionError represent a failure to connect/dial a connection, with context
type ConnectionError struct {
	Addr  string
	Error error
}

// Shutdown terminates all subsystems safely.  To stop the Daemon run loop, send a value
// over the quit channel provided to Init.  The Daemon run lopp must be stopped
// before calling this function.
func (dm *Daemon) Shutdown() {
	if dm.DHT != nil {
		dm.DHT.Shutdown()
	}

	dm.Peers.Shutdown()
	dm.Pool.Shutdown() //send disconnect message first

	dm.Pool = nil
	dm.DHT = nil
	dm.Peers = nil
}

// Runs initialization that must complete before the Start goroutine

//func (self *Daemon) Init() {
//	if !self.Config.DisableIncomingConnections {
//self.Pool.Listen()
//if err := self.Pool.StartListen(); err != nil {
//	log.Panic(err)
//}
//go self.Pool.AcceptConnections() //listen for connections
//	}
//}

// Start Main loop for peer/connection management. Send anything to quit to shut it
// down
func (dm *Daemon) Start(quit chan int) {
	if !dm.Config.DisableIncomingConnections {
		//listen for incoming
		if err := dm.Pool.StartListen(); err != nil {
			log.Panic(err)
		}
		//goroutine for accepting incoming
		go dm.Pool.AcceptConnections() //listen for connections
	}

	//fix this, should poll without delay
	messageHandlingTicker := time.Tick(time.Millisecond * 10)

	//peer exchange tickers
	clearOldPeersTicker := time.Tick(dm.Peers.Config.CullRate)
	//requestPeersTicker := time.Tick(self.Peers.Config.RequestRate)
	updateBlacklistTicker := time.Tick(dm.Peers.Config.UpdateBlacklistRate)

	//daemon tickers
	//privateConnectionsTicker := time.Tick(self.Config.PrivateRate)
	cullInvalidTicker := time.Tick(dm.Config.CullInvalidRate)
	//outgoingConnectionsTicker := time.Tick(self.Config.OutgoingRate)

main:

	for {

		select {

		//Module: Peers

		// Flush expired blacklisted peers
		case <-updateBlacklistTicker:
			if !dm.Peers.Config.Disabled {
				dm.Peers.Peers.Blacklist.Refresh()
			}
		// Request peers via PEX
		//case <-requestPeersTicker:
		//	self.Peers.requestPeers(self.Service)

		// Remove peers we haven't seen in a while
		case <-clearOldPeersTicker:
			if !dm.Peers.Config.Disabled {
				dm.Peers.Peers.Peerlist.ClearOld(dm.Peers.Config.Expiration)
			}

		// Module: Pool

		//process incoming messages
		case <-messageHandlingTicker:
			if !dm.Config.DisableNetworking {
				dm.Pool.HandleMessages()
			}
		// Process disconnections
		case r := <-dm.Pool.DisconnectQueue:
			if dm.Config.DisableNetworking {
				log.Panic("There should be nothing in the DisconnectQueue")
			}
			dm.Pool.HandleDisconnectEvent(r)

		//Module: Daemon

		// Remove connections that failed to complete the handshake
		case <-cullInvalidTicker:
			if !dm.Config.DisableNetworking {
				dm.cullInvalidConnections()
			}

		case r := <-dm.connectionErrors:
			if dm.Config.DisableNetworking {
				log.Panic("There should be no connection errors")
			}
			dm.handleConnectionError(r)

		case <-quit:
			logger.Info("Breaking From Daemon Main")
			break main
		}
	}
}

// GetListenPort returns the ListenPort for a given address.  If no port is found, 0 is
// returned
//this might be broken now
func (dm *Daemon) GetListenPort(addr string) uint16 {
	_, p, err := SplitAddr(addr)
	if err != nil {
		logger.Error("GetListenPort received invalid addr: %v", err)
		return 0
	}
	return p
}

// ConnectToAddr connects to a particular service on daemon
func (dm *Daemon) ConnectToAddr(addr string, service *gnet.Service) error {

	//addr should be ip:port and ip/port must be valid
	//if its not, the connect attempt will just fail

	//connected to daemon, connect to service
	if dm.Pool.Addresses[addr] != nil {
		c := dm.Pool.Addresses[addr]
		dm.ConnectToService(c, service)
		return nil
	}
	//not connected
	if dm.Pool.Addresses[addr] == nil {
		//only the first service connection triggers connection attempt
		if dm.pendingConnections[addr] == nil {
			dm.pendingConnections[addr] = make([]*gnet.Service, 0)
			go func() {
				_, err := dm.Pool.Connect(addr)
				if err != nil {
					dm.connectionErrors <- ConnectionError{addr, err}
				}
			}()
		}
		if service != nil {
			dm.pendingConnections[addr] = append(dm.pendingConnections[addr], service)
		}
		return nil
	}
	return nil
}

// We remove a peer from the Pex if we failed to connect
func (dm *Daemon) handleConnectionError(c ConnectionError) {
	logger.Debug("Removing %s because failed to connect: %v", c.Addr,
		c.Error)
	delete(dm.pendingConnections, c.Addr)
	dm.Peers.RemovePeer(c.Addr)
}

// Removes unsolicited connections who haven't sent a version

func (dm *Daemon) cullInvalidConnections() {
	// This method only handles the erroneous people from the DHT, but not
	// malicious nodes
	now := util.Now()
	for a, t := range dm.ExpectingIntroductions {
		// Forget about anyone that already disconnected
		if dm.Pool.Addresses[a] == nil {
			delete(dm.ExpectingIntroductions, a)
			continue
		}
		// Remove anyone that fails to send a version within introductionWait time
		if t.Add(dm.Config.IntroductionWait).Before(now) {
			logger.Info("Removing %s for not sending introduction", a)
			delete(dm.ExpectingIntroductions, a)
			dm.Pool.Disconnect(dm.Pool.Addresses[a],
				ErrDisconnectIntroductionTimeout)
			dm.Peers.RemovePeer(a)
		}
	}
}

// Called when a ConnectEvent is processed off the onConnectEvent channel
func (dm *Daemon) onConnect(c *gnet.Connection, solicited bool) {
	a := c.Addr()

	if solicited {
		logger.Info("Connected to %s as we requested", a)
	} else {
		logger.Info("Received unsolicited connection to %s", a)
	}

	serviceConList := dm.pendingConnections[a] //list of services to connect to
	delete(dm.pendingConnections, a)

	blacklisted := dm.Peers.Peers.IsBlacklisted(a)
	if blacklisted {
		logger.Info("%s is blacklisted, disconnecting", a)
		dm.Pool.Disconnect(c, ErrDisconnectIsBlacklisted)
		return
	}

	if dm.Pool.Addresses[a] != nil {
		logger.Info("Already connected to %s, disconnecting", a)
		dm.Pool.Disconnect(c, ErrDisconnectConnectedTwice)
	}

	if solicited {
		dm.OutgoingConnections[a] = c
	}
	dm.ExpectingIntroductions[a] = util.Now()
	logger.Debug("Sending introduction message to %s", a)

	m := NewIntroductionMessage(MirrorConstant, dm.Config.Version,
		dm.Pool.Config.Port)
	dm.Service.Send(c, m)

	//send connection message to each service in list
	for _, service := range serviceConList {
		dm.ConnectToService(c, service)
	}
}

// Triggered when an gnet.Connection terminates. Disconnect events are not
// pushed to a separate channel, because disconnects are already processed
// by a queue in the daemon.Run() select{}.
func (dm *Daemon) onGnetDisconnect(c *gnet.Connection,
	reason gnet.DisconnectReason) {

	a := c.Addr()
	logger.Info("%s disconnected because: %v", a, reason)
	duration, exists := BlacklistOffenses[reason]
	if exists {
		dm.Peers.Peers.AddBlacklistEntry(a, duration)
	}
	delete(dm.OutgoingConnections, a)
	delete(dm.ExpectingIntroductions, a)
}

// Triggered when an gnet.Connection is connected
//func (self *Daemon) onGnetConnect(c *gnet.Connection, solicited bool) {
//	self.onConnectEvent <- ConnectEvent{Addr: c.Addr(), Solicited: solicited}
//}

// LocalhostIP returns the address for localhost on the machine
func LocalhostIP() string {
	tt, err := net.Interfaces()
	if err != nil {
		log.Panicf("Failed to obtain localhost IP: %v", err)
		return ""
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			log.Panicf("Failed to obtain localhost IP: %v", err)
			return ""
		}
		for _, a := range aa {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
				return ipnet.IP.String()
			}
		}
	}
	log.Panicf("Failed to obtain localhost IP: No Local IP found")
	return ""
}

// Returns true if addr is a localhost address
//func IsLocalhost(addr string) bool {
//	return net.ParseIP(addr).IsLoopback()
//}

// SplitAddr splits an ip:port string to ip, port
func SplitAddr(addr string) (string, uint16, error) {
	pts := strings.Split(addr, ":")
	if len(pts) != 2 {
		return pts[0], 0, fmt.Errorf("Invalid addr %s", addr)
	}
	port64, err := strconv.ParseUint(pts[1], 10, 16)
	if err != nil {
		return pts[0], 0, fmt.Errorf("Invalid port in %s", addr)
	}
	return pts[0], uint16(port64), nil
}

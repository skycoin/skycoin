package udp

import(
    "net"
    "fmt"
    "os"
    "time"
    "sync"
    "errors"
    "strconv"
    "reflect"
    "crypto/rand"
    "encoding/json"
    "encoding/binary")

import(
	"github.com/skycoin/encoder"
	"github.com/skycoin/skycoin/src/mesh2/transport"
    "github.com/skycoin/skycoin/src/cipher")

import(
    "github.com/ccding/go-stun/stun")

type UDPConfig struct {
	transport.TransportConfig
	DatagramLength	uint16
	LocalAddress string 	// "" for default

	NumListenPorts uint16
	ListenPortMin uint16		// If 0, STUN is used
	ExternalAddress string  	// External address to use if STUN is not
	StunEndpoints []string		// STUN servers to try for NAT traversal
}

type ListenPort struct {
	externalHost net.UDPAddr
	conn *net.UDPConn
}

type UDPCommConfig struct {
	DatagramLength	uint16
	ExternalHosts []net.UDPAddr
}

type UDPTransport struct {
	config UDPConfig
	listenPorts []ListenPort
	messagesReceived chan []byte
	closing chan bool
	closeWait *sync.WaitGroup
	crypto transport.TransportCrypto

	// Thread protected variables
	lock *sync.Mutex
	connectedPeers map[cipher.PubKey]UDPCommConfig
}

func OpenUDPPort(port_index uint16, config UDPConfig, wg *sync.WaitGroup, 
				 errorChan chan error, portChan chan ListenPort) () {
	defer wg.Done()

	port := (uint16)(0)
	if config.ListenPortMin > 0 {
		port = config.ListenPortMin + port_index
	}

	udpAddr := net.JoinHostPort(config.LocalAddress, strconv.Itoa((int)(port)))
    listenAddr,resolvErr := net.ResolveUDPAddr("udp", udpAddr)
    if resolvErr != nil {
    	errorChan <- resolvErr
    	return
    }
 
    udpConn,listenErr := net.ListenUDP("udp", listenAddr)
    if listenErr != nil {
    	errorChan <- listenErr
    	return
    }

    externalHostStr := net.JoinHostPort(config.ExternalAddress, strconv.Itoa((int)(port)))
    externalHost := &net.UDPAddr{}
    externalHost, resolvErr = net.ResolveUDPAddr("udp", externalHostStr)
    if resolvErr != nil {
    	errorChan <- resolvErr
    	return
    }

	if config.ListenPortMin == 0 {
		if (config.StunEndpoints == nil) || len(config.StunEndpoints) == 0 {
			errorChan <- errors.New("No local port or STUN endpoints specified in config: no way to receive datagrams")
	    	return
		}
		var stun_success bool = false
		for _, addr := range config.StunEndpoints {
			stunClient := stun.NewClientWithConnection(udpConn)
			stunClient.SetServerAddr(addr)

			_, host, error := stunClient.Discover()
			if error != nil {
				fmt.Fprintf(os.Stderr, "STUN Error for Endpoint '%v': %v\n", addr, error)
				continue
			} else {
				externalHostStr = host.TransportAddr()
			    externalHost, resolvErr = net.ResolveUDPAddr("udp", externalHostStr)
			    if resolvErr != nil {
			    	errorChan <- resolvErr
			    	return
			    }
				stun_success = true
				break
			}
		}
		if !stun_success {
			errorChan <- errors.New("All STUN requests failed")
    		return
		}
	}

	// STUN library sets the deadlines
    udpConn.SetDeadline(time.Time{})
    udpConn.SetReadDeadline(time.Time{})
    udpConn.SetWriteDeadline(time.Time{})
	portChan <- ListenPort{*externalHost, udpConn}
}

func (self*UDPTransport) receiveMessage(buffer []byte) {
	if self.crypto != nil {
		buffer = self.crypto.Decrypt(buffer)
	}
	var v reflect.Value = reflect.New(reflect.TypeOf([]byte{}))
	_, err := encoder.DeserializeRawToValue(buffer, v)
    if err != nil {
    	fmt.Fprintf(os.Stderr, "Error on DeserializeRawToValue: %v\n", err)
        return
    }
    m, succ := (v.Elem().Interface()).(interface{})
    if !succ {
    	fmt.Fprintf(os.Stderr, "Error on Interface()\n")
    	return
    }
    recv_chan := self.messagesReceived
    if recv_chan != nil {
    	recv_chan <- m.([]byte)
    }
}

func strongUint() uint32 {
	socket_i_b := make([]byte, 4)
	n, err := rand.Read(socket_i_b[:4])
	if n != 4 || err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint32(socket_i_b)
}

func (self*UDPTransport) safeGetPeerComm(peer cipher.PubKey) (*UDPCommConfig, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	peerComm, foundPeer := self.connectedPeers[peer]
	if !foundPeer {
		return nil, false
	}
	return &peerComm, true
}

func (self*UDPTransport) listenTo(port ListenPort) {
	self.closeWait.Add(1)
	defer self.closeWait.Done()

	buffer := make([]byte, self.config.DatagramLength)

	for len(self.closing) == 0 {
		n, _, err := port.conn.ReadFromUDP(buffer)
		if err != nil {
			if len(self.closing) == 0 {
				fmt.Fprintf(os.Stderr, "Error on ReadFromUDP for %v: %v\n", port.externalHost, err)
			}
			break
		}
		self.receiveMessage(buffer[:n])
	}
}

// Blocks waiting for STUN requests, port opening
func NewUDPTransport(config UDPConfig) (*UDPTransport, error) {
	if config.DatagramLength < 32 {
		return nil, errors.New("Datagram length too short")
	}

	// Open all ports at once
	errors := make(chan error, config.NumListenPorts)
	ports := make(chan ListenPort, config.NumListenPorts)
	var portGroup sync.WaitGroup
	portGroup.Add((int)(config.NumListenPorts))
	for port_i := (uint16)(0); port_i < config.NumListenPorts; port_i++ {
		go OpenUDPPort(port_i, config, &portGroup, errors, ports)
	}
	portGroup.Wait()

	if len(errors) > 0 {
		for len(ports) > 0 {
			port := <- ports
			port.conn.Close()
		}
		return nil, <- errors
	}

	portsArray := make([]ListenPort, 0)
	for len(ports) > 0 {
		port := <- ports
		portsArray = append(portsArray, port)
	}	

	waitGroup := &sync.WaitGroup{}
	ret := &UDPTransport{
		config,
		portsArray,
		nil,	// Receive channel
		make(chan bool, 10 * len(portsArray)), // closing
		waitGroup,
		nil,	// No crypto by default
		&sync.Mutex{},
		make(map[cipher.PubKey]UDPCommConfig),
	}

	for _, port := range ret.listenPorts {
		go ret.listenTo(port)
	}

	return ret, nil
}

func (self*UDPTransport) Close() error {
	self.closeWait.Add(len(self.listenPorts))
	for i := 0;i < 10*len(self.listenPorts);i++ {
		self.closing <- true
	}

	for _, port := range self.listenPorts {
		go func (conn *net.UDPConn) {
			conn.Close()
			self.closeWait.Done()
		}(port.conn)
	}

	self.closeWait.Wait()
	return nil
}

func (self*UDPTransport) SetCrypto(crypto transport.TransportCrypto) {
	self.crypto = crypto
}

func (self*UDPTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	_, found := self.safeGetPeerComm(peer)
	return found
}

func (self*UDPTransport) GetConnectedPeers() []cipher.PubKey {
	self.lock.Lock()
	defer self.lock.Unlock()
	ret := []cipher.PubKey{}
	for key, _ := range(self.connectedPeers) {
		ret = append(ret, key)
	}
	return ret
}

func (self*UDPTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	commConfig, found := self.safeGetPeerComm(peer)
	if !found {
		fmt.Fprintf(os.Stderr, "Unknown peer passed to GetMaximumMessageSizeToPeer: %v\n", peer)
		return 0
	}
	return (uint)(commConfig.DatagramLength)
}

// May block
func (self*UDPTransport) SendMessage(toPeer cipher.PubKey, contents []byte) error {
	// Find pubkey
	peerComm, found := self.safeGetPeerComm(toPeer)
	if !found {
		return errors.New(fmt.Sprintf("Dropping message that is to an unknown peer: %v\n", toPeer))
	}

	// Check length
	if len(contents) > int(peerComm.DatagramLength) {
		return errors.New(fmt.Sprintf("Dropping message that is too large: %v > %v\n", len(contents), self.config.DatagramLength))
	}

	// Pad to length
	encoderBuffer := encoder.Serialize(contents)
	datagramBuffer := make([]byte, peerComm.DatagramLength)
	copy(datagramBuffer, encoderBuffer)

	// Apply crypto
	if self.crypto != nil {
		datagramBuffer = self.crypto.Encrypt(datagramBuffer)
	}

	// Choose a socket randomly
	fromSocketIndex := strongUint() % (uint32)(len(self.listenPorts))
	conn := self.listenPorts[fromSocketIndex].conn

	// Send datagram
	toAddrIndex := strongUint() % (uint32)(len(peerComm.ExternalHosts))
	toAddr := peerComm.ExternalHosts[toAddrIndex]

	n, err := conn.WriteToUDP(datagramBuffer, &toAddr)
	if err != nil {
		return errors.New(fmt.Sprintf("Error on WriteToUDP: %v\n", err))
	}
	if n != int(peerComm.DatagramLength) {
		return errors.New(fmt.Sprintf("WriteToUDP returned %v != %v\n", n, peerComm.DatagramLength))
	}
	return nil
}

func  (self*UDPTransport) SetReceiveChannel(received chan []byte) {
	self.messagesReceived = received
}

// UDP Transport only functions
func (self*UDPTransport) GetTransportConnectInfo() string {
	hostsArray := make([]net.UDPAddr, 0)

	for _, port := range self.listenPorts {
		hostsArray = append(hostsArray, port.externalHost)
	}

	info := UDPCommConfig{
		self.config.DatagramLength,
		hostsArray,
	}

	ret, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	return string(ret)
}

func (self*UDPTransport) ConnectToPeer(peer cipher.PubKey, connectInfo string) error {
	config := UDPCommConfig{}
	parseError := json.Unmarshal([]byte(connectInfo), &config)
	if parseError != nil {
		return parseError
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	_, connected := self.connectedPeers[peer]
	if connected {
		return errors.New(fmt.Sprintf("Already connected to peer %v", peer))
	}
	self.connectedPeers[peer] = config
	return nil
}

func (self*UDPTransport) DisconnectFromPeer(peer cipher.PubKey) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.connectedPeers, peer)
}


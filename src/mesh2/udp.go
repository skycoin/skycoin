package mesh

import(
    "net"
    "fmt"
    "os"
    "sync"
    "errors"
    "strconv"
    "encoding/json")

import(
    "github.com/skycoin/skycoin/src/cipher")

import(
    "github.com/ccding/go-stun/stun")

type UDPConfig struct {
	TransportConfig
	DatagramLength	uint64
	LocalAddress string 	// "" for default

	NumListenPorts uint16
	ListenPortMin uint16		// If 0, STUN is used
	StunEndpoints []string		// STUN servers to try for NAT traversal
}

type ListenPort struct {
	externalHost string
	conn *net.UDPConn
}

type UDPTransport struct {
	config UDPConfig
	listenPorts []ListenPort
	messagesToSend chan TransportMessage
	messagesReceived chan TransportMessage

	closing chan bool
	closeWait *sync.WaitGroup
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

	externalHost := udpConn.LocalAddr().String()

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
				externalHost = host.TransportAddr()
				stun_success = true
				break
			}
		}
		if !stun_success {
			errorChan <- errors.New("All STUN requests failed")
    		return
		}
	}

	portChan <- ListenPort{externalHost, udpConn}
}

func (self*UDPTransport) receiveMessage(buffer []byte) {
	// ...
}

func (self*UDPTransport) sendMessage(message TransportMessage) {
	// ...
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

func (self*UDPTransport) sendLoop() {
	self.closeWait.Add(1)
	defer self.closeWait.Done()

	for {
		select {
			case message := <- self.messagesToSend: {
				self.sendMessage(message)
				break
			}
			case <- self.closing:
				return
		}
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
		make(chan TransportMessage, config.SendChannelLength),
		make(chan TransportMessage, config.ReceiveChannelLength),
		make(chan bool, 10 * len(portsArray)), // closing
		waitGroup,
	}

	for _, port := range ret.listenPorts {
		go ret.listenTo(port)
	}

	go ret.sendLoop()

	return ret, nil
}

func (self*UDPTransport) Close() {
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
}

type UDPCommConfig struct {
	DatagramLength	uint64
	ExternalHosts []string
}

func (self*UDPTransport) GetTransportConnectInfo() string {
	hostsArray := make([]string, 0)

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

func (self*UDPTransport) IsReliable() bool {
	return false
}

func (self*UDPTransport) ConnectedToPeer(peer cipher.PubKey, connectInfo string) bool {
	return false
}

func (self*UDPTransport) RetransmitIntervalHint(toPeer cipher.PubKey) uint32 {
	// TODO: Implement latency tracking
	return 500
}

func (self*UDPTransport) DisconnectFromPeer(peer cipher.PubKey) {

}

func (self*UDPTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	return 0
}

func (self*UDPTransport) SendMessage(msg TransportMessage) error {
	return nil
}

func (self*UDPTransport) GetReceiveChannel() chan TransportMessage {
	return self.messagesReceived
}




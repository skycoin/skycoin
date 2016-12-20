package physical

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/ccding/go-stun/stun"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

type ListenPort struct {
	externalHost net.UDPAddr
	conn         *net.UDPConn
}

type UDPConfig struct {
	TransportConfig
	DatagramLength  uint16
	LocalAddress    string   // "" for default
	ListenPort      uint16   // If 0, STUN is used
	ExternalAddress string   // External address to use if STUN is not
	StunEndpoints   []string // STUN servers to try for NAT traversal
}

type TransportConfig struct {
	SendChannelLength uint32
}

type UDPCommConfig struct {
	DatagramLength uint16
	ExternalHost   net.UDPAddr
	CryptoKey      []byte
}

type UDPTransport struct {
	config           UDPConfig
	listenPort       ListenPort
	messagesReceived chan []byte
	closing          chan bool
	closeWait        *sync.WaitGroup
	crypto           transport.ITransportCrypto

	// Thread protected variables
	lock              *sync.Mutex
	connectedPeerKey  *cipher.PubKey
	connectedPeerConf *UDPCommConfig
}

var emptyPK cipher.PubKey = cipher.PubKey{}

func OpenUDPPort(config UDPConfig, wg *sync.WaitGroup,
	errorChan chan error, portChan chan ListenPort) {
	defer wg.Done()

	port := (uint16)(0)
	if config.ListenPort > 0 {
		port = config.ListenPort
	}

	udpAddr := net.JoinHostPort(config.LocalAddress, strconv.Itoa((int)(port)))
	listenAddr, resolvErr := net.ResolveUDPAddr("udp", udpAddr)
	if resolvErr != nil {
		errorChan <- resolvErr
		return
	}

	udpConn, listenErr := net.ListenUDP("udp", listenAddr)
	if listenErr != nil {
		errorChan <- listenErr
		return
	}

	externalHostStr := net.JoinHostPort(config.ExternalAddress, strconv.Itoa((int)(port)))
	externalHost, resolvErr := net.ResolveUDPAddr("udp", externalHostStr)
	if resolvErr != nil {
		errorChan <- resolvErr
		return
	}

	//check if one UDPConfig can have more than 1 StunEndpoint
	if config.ListenPort == 0 {
		if (config.StunEndpoints == nil) || len(config.StunEndpoints) == 0 {
			errorChan <- errors.New("No local port or STUN endpoints specified in config: no way to receive datagrams")
			return
		}
		var stun_success bool = false
		for _, addr := range config.StunEndpoints {
			stunClient := stun.NewClientWithConnection(udpConn)
			stunClient.SetServerAddr(addr)

			_, host, error := stunClient.Discover()
			if error != nil || host == nil {
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

func (s *UDPTransport) receiveMessage(buffer []byte) {
	if s.crypto != nil {
		buffer = s.crypto.Decrypt(buffer)
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
	recv_chan := s.messagesReceived
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

<<<<<<< HEAD
func (self *UDPTransport) safeGetPeerComm(peer cipher.PubKey) (*UDPCommConfig, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if peer != *self.connectedPeerKey {
		return nil, false
	}
	return self.connectedPeerConf, true
=======
func (s *UDPTransport) safeGetPeerComm(peer cipher.PubKey) (*UDPCommConfig, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if peer != *s.connectedPeerKey {
		return nil, false
	}
	return s.connectedPeerConf, true
>>>>>>> 662d87062c1592cf12ec5fd885179ac2289a3af9
}

func (s *UDPTransport) listenTo(port ListenPort) {
	s.closeWait.Add(1)
	defer s.closeWait.Done()

	buffer := make([]byte, s.config.DatagramLength)

	for len(s.closing) == 0 {
		n, _, err := port.conn.ReadFromUDP(buffer)
		if err != nil {
			if len(s.closing) == 0 {
				fmt.Fprintf(os.Stderr, "Error on ReadFromUDP for %v: %v\n", port.externalHost, err)
			}
			break
		}
		s.receiveMessage(buffer[:n])
	}
}

// Blocks waiting for STUN requests, port opening
func NewUDPTransport(config UDPConfig) (*UDPTransport, error) {
	if config.DatagramLength < 32 {
		return nil, errors.New("Datagram length too short")
	}

	// Open all ports at once
	errors := make(chan error, 1)
	ports := make(chan ListenPort, 1)
	var portGroup sync.WaitGroup
	portGroup.Add(1)
	go OpenUDPPort(config, &portGroup, errors, ports)
	portGroup.Wait()

	if len(errors) > 0 {
		for len(ports) > 0 {
			port := <-ports
			port.conn.Close()
		}
		return nil, <-errors
	}

	port := <-ports

	waitGroup := &sync.WaitGroup{}
	ret := &UDPTransport{
		config,
		port,
		nil,                 // Receive channel
		make(chan bool, 10), // closing
		waitGroup,
		nil, // No crypto by default
		&sync.Mutex{},
		&cipher.PubKey{},
		&UDPCommConfig{},
	}

	go ret.listenTo(ret.listenPort)

	return ret, nil
}

func (s *UDPTransport) Close() error {
	s.closeWait.Add(1)
	for i := 0; i < 10; i++ {
		s.closing <- true
	}

	go func(conn *net.UDPConn) {
		conn.Close()
		s.closeWait.Done()
	}(s.listenPort.conn)

	s.closeWait.Wait()
	return nil
}

func (s *UDPTransport) SetCrypto(crypto transport.ITransportCrypto) {
	s.crypto = crypto
}

func (s *UDPTransport) ConnectedToPeer(peer cipher.PubKey) bool {
	_, found := s.safeGetPeerComm(peer)
	return found
}

func (s *UDPTransport) GetConnectedPeer() cipher.PubKey {
	s.lock.Lock()
	defer s.lock.Unlock()
	return *s.connectedPeerKey
}

func (s *UDPTransport) GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint {
	commConfig, found := s.safeGetPeerComm(peer)
	if !found {
		fmt.Fprintf(os.Stderr, "Unknown peer passed to GetMaximumMessageSizeToPeer: %v\n", peer)
		return 0
	}
	return (uint)(commConfig.DatagramLength)
}

// May block
func (s *UDPTransport) SendMessage(toPeer cipher.PubKey, contents []byte, retChan chan error) error {
	var retErr error = nil
	// Find pubkey
	peerComm, found := s.safeGetPeerComm(toPeer)
	if !found {
		retErr = errors.New(fmt.Sprintf("Dropping message that is to an unknown peer: %v\n", toPeer))
		if retChan != nil {
			retChan <- retErr
		}
		return retErr
	}

	// Check length
	if len(contents) > int(peerComm.DatagramLength) {
<<<<<<< HEAD
		retErr = errors.New(fmt.Sprintf("Dropping message that is too large: %v > %v\n", len(contents), self.config.DatagramLength))
=======
		retErr = errors.New(fmt.Sprintf("Dropping message that is too large: %v > %v\n", len(contents), s.config.DatagramLength))
>>>>>>> 662d87062c1592cf12ec5fd885179ac2289a3af9
		if retChan != nil {
			retChan <- retErr
		}
		return retErr
	}

	// Pad to length
	encoderBuffer := encoder.Serialize(contents)
	datagramBuffer := make([]byte, peerComm.DatagramLength)
	copy(datagramBuffer, encoderBuffer)

	// Apply crypto
	if s.crypto != nil {
		datagramBuffer = s.crypto.Encrypt(datagramBuffer, peerComm.CryptoKey)
	}

	// Choose a socket randomly
	conn := s.listenPort.conn

	// Send datagram
	toAddr := peerComm.ExternalHost

	n, err := conn.WriteToUDP(datagramBuffer, &toAddr)
	if err != nil {
		retErr = errors.New(fmt.Sprintf("Error on WriteToUDP: %v\n", err))
		if retChan != nil {
			retChan <- retErr
		}
		return retErr
	}
	if n != int(peerComm.DatagramLength) {
		retErr = errors.New(fmt.Sprintf("WriteToUDP returned %v != %v\n", n, peerComm.DatagramLength))
		if retChan != nil {
			retChan <- retErr
		}
		return retErr
	}
	if retChan != nil {
		retChan <- nil
	}
	return nil
}

func (s *UDPTransport) SetReceiveChannel(received chan []byte) {
	s.messagesReceived = received
}

func (s *UDPTransport) safeGetCrypto() transport.ITransportCrypto {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.crypto
}

// UDP Transport only functions
func (s *UDPTransport) GetTransportConnectInfo() string {

	key := []byte{}
	crypto := s.safeGetCrypto()
	if crypto != nil {
		key = crypto.GetKey()
	}
	info := UDPCommConfig{
		s.config.DatagramLength,
		s.listenPort.externalHost,
		key,
	}

	ret, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(ret)
}

func (s *UDPTransport) ConnectToPeer(peer cipher.PubKey, connectInfo string) error {
	if peer == emptyPK {
		return errors.New(fmt.Sprintf("Cannot connect to empty peer %v", peer))
	}
	config := UDPCommConfig{}
	connectInfoRaw, decodeHexError := hex.DecodeString(connectInfo)
	if decodeHexError != nil {
		return decodeHexError
	}
	parseError := json.Unmarshal(connectInfoRaw, &config)
	if parseError != nil {
		return parseError
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if *s.connectedPeerKey != emptyPK {
		return errors.New(fmt.Sprintf("Already connected to peer %v", peer))
	}
	s.connectedPeerKey = &peer
	s.connectedPeerConf = &config
	return nil
}

func (s *UDPTransport) DisconnectFromPeer(peer cipher.PubKey) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.connectedPeerKey = &cipher.PubKey{}
	s.connectedPeerConf = &UDPCommConfig{}
}

// Create UDPTransport
func CreateNewUDPTransport(configUdp UDPConfig) *UDPTransport {
	udpTransport, createUDPError := NewUDPTransport(configUdp)
	if createUDPError != nil {
		panic(createUDPError)
	}
	return udpTransport
}

// Create Udp config
func CreateUdp(port int, externalA string) UDPConfig {
	udp := UDPConfig{}
	udp.SendChannelLength = uint32(100)
	udp.DatagramLength = uint16(512)
	udp.LocalAddress = ""
	udp.ListenPort = uint16(port)
	udp.ExternalAddress = externalA

	return udp
}

// Create info for the peer's connection.
func CreateUDPCommConfig(addr string, cryptoKey []byte) string {
	config := UDPCommConfig{}
	config.DatagramLength = uint16(512)
	address1, _ := net.ResolveUDPAddr("", addr)
	config.ExternalHost = *address1

	if cryptoKey == nil {
		var zero = make([]byte, 32, 32)
		cryptoKey = zero
	}
	if len(cryptoKey) != 32 {
		log.Panic("Error: mesh.transport.protocol, CreateUDPCommConfig, crypto key length != 32")
	}
	config.CryptoKey = cryptoKey

	src, _ := json.Marshal(&config)
	infoPeer := hex.EncodeToString(src)

	return infoPeer
}

package nodemanager

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"encoding/hex"
	"encoding/json"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	mesh "github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/transport"
	"github.com/skycoin/skycoin/src/mesh/transport/physical"
)

const LOCALHOST string = "127.0.0.1"

func init() {
	ServerConfig = CreateTestConfig(15101)

	b := []byte{234, 15, 123, 220, 185, 171, 218, 20, 130, 48, 24, 255, 214, 133, 191, 164, 211, 190, 224, 127, 105, 125, 141, 178, 226, 250, 123, 149, 229, 33, 187, 165, 27}
	p := cipher.PubKey{}
	copy(p[:], b[:])
	ServerConfig.NodeConfig.PubKey = p

	//b2 := [32]byte{48, 126, 177, 168, 139, 146, 205, 8, 191, 110, 195, 254, 184, 22, 168, 118, 237, 126, 87, 224, 171, 243, 239, 87, 106, 152, 251, 217, 120, 239, 88, 138}
	//ServerConfig.Node.ChaCha20Key = b2
}

var ServerConfig *TestConfig

type NodeManager struct {
	Port       int // nodemanager port is for testing purposes only, ports are assigned from config files
	ConfigList map[cipher.PubKey]*TestConfig
	NodesList  map[cipher.PubKey]*mesh.Node
	PubKeyList []cipher.PubKey
	//	Routes     map[RouteKey]Route
	Routes     map[RouteKey]*RouteConfig
	RouteGraph *Graph
}

//configuration in here
type NodeManagerConfig struct {
}

//add config eventually
func NewNodeManager(config *NodeManagerConfig) *NodeManager {
	nm := NodeManager{}
	return &nm
}

func NewEmptyNodeManager() *NodeManager {
	nm := NodeManager{}
	nm.ConfigList = map[cipher.PubKey]*TestConfig{}
	nm.NodesList = map[cipher.PubKey]*mesh.Node{}
	nm.PubKeyList = []cipher.PubKey{}
	nm.Routes = map[RouteKey]*RouteConfig{}
	nm.RouteGraph = NewGraph()
	return &nm
}

// run node manager, in goroutine;
// call Shutdown to stop
func (s *NodeManager) Start() {
	//	s.RebuildRouteGraph() // maybe it will be better to leave it here; then no need to call it before calling this
	for _, pubKey := range s.PubKeyList {
		config, found := s.ConfigList[pubKey]
		if !found {
			fmt.Println("No config found with PubKey", pubKey)
			panic("")
		}
		node, found := s.NodesList[pubKey]
		if !found {
			node = CreateNode(*config)
			s.NodesList[pubKey] = node
		}
		AddPeersToNode(node, *config)
	}
}

//called to trigger shutdown
func (s *NodeManager) Shutdown() {
	s.CloseAll()
}

type RouteKey struct {
	SourceNode cipher.PubKey
	TargetNode cipher.PubKey
}

type Route struct {
	SourceNode        cipher.PubKey
	TargetNode        cipher.PubKey
	RoutesToEstablish []cipher.PubKey
	Weight            int
}

func CreateTestConfig(port int) *TestConfig {
	return CreateConfig(LOCALHOST, port)
}

// Create TestConfig to the test using the functions created in the meshnet library.
func CreateConfig(address string, port int) *TestConfig {
	testConfig := &TestConfig{}
	testConfig.ExternalAddress = address
	testConfig.StartPort = port
	testConfig.Port = port
	testConfig.NodeConfig = NewNodeConfig()
	testConfig.TransportConfig = transport.CreateTransportConfig(testConfig.NodeConfig.PubKey)
	testConfig.PeerToPeers = map[string]*Peer{}

	return testConfig
}

func CreateNode(config TestConfig) *mesh.Node {
	node, createNodeError := mesh.NewNode(config.NodeConfig)
	if createNodeError != nil {
		panic(createNodeError)
	}

	return node
}

// Create public key
func CreatePubKey() cipher.PubKey {
	pub, _ := cipher.GenerateKeyPair()
	return pub
}

// Create new node config
func NewNodeConfig() mesh.NodeConfig {
	nodeConfig := mesh.NodeConfig{}
	nodeConfig.PubKey = CreatePubKey()
	//nodeConfig.ChaCha20Key = CreateChaCha20Key()
	nodeConfig.MaximumForwardingDuration = 1 * time.Minute
	nodeConfig.RefreshRouteDuration = 5 * time.Minute
	nodeConfig.ExpireRoutesInterval = 5 * time.Minute
	nodeConfig.TransportMessageChannelLength = 100

	return nodeConfig
}

// Create node list
func (s *NodeManager) CreateNodeConfigList(n int) {
	s.ConfigList = make(map[cipher.PubKey]*TestConfig)
	s.NodesList = make(map[cipher.PubKey]*mesh.Node)
	for a := 1; a <= n; a++ {
		s.AddNode()
	}
}

// Add Node to Nodes List
func (s *NodeManager) AddNode() int {
	if len(s.ConfigList) == 0 {
		s.ConfigList = make(map[cipher.PubKey]*TestConfig)
		s.NodesList = make(map[cipher.PubKey]*mesh.Node)
	}
	config := CreateTestConfig(s.Port)
	s.Port += 100 // to avoid overlaps
	s.ConfigList[config.NodeConfig.PubKey] = config
	node := CreateNode(*config)
	s.NodesList[config.NodeConfig.PubKey] = node
	s.PubKeyList = append(s.PubKeyList, config.NodeConfig.PubKey)
	index := len(s.NodesList) - 1
	return index
}

func connected(config1, config2 *TestConfig) bool {
	to := config1.NodeConfig.PubKey
	for _, peerToPeer := range config2.PeerToPeers {
		if peerToPeer.Peer == to {
			return true
		}
	}
	return false
}

func ConnectNodeToNode(config1, config2 *TestConfig) {
	if connected(config1, config2) && connected(config2, config1) {
		return
	}
	config1.AddPeerToConnect(config2)
	config1.AddRouteToEstablish(config2)
	config2.AddPeerToConnect(config1)
	config2.AddRouteToEstablish(config1)
	config1.Port++
	config2.Port++
}

// Connect the node list
func (s *NodeManager) ConnectNodes() {

	lenght := len(s.ConfigList)

	if lenght > 1 {
		for index := 0; index < lenght-1; index++ {
			pubKey1 := s.PubKeyList[index]
			pubKey2 := s.PubKeyList[index+1]
			config1 := s.ConfigList[pubKey1]
			config2 := s.ConfigList[pubKey2]
			ConnectNodeToNode(config1, config2)
			AddPeersToNode(s.NodesList[pubKey1], *config1)
		}
	}
}

// Add Routes to Node
func AddRoutesToEstablish(node *mesh.Node, routesConfigs []RouteConfig) {
	for _, routeConfig := range routesConfigs {
		addRouteToEstablish(node, routeConfig)
	}
}

func addRouteToEstablish(node *mesh.Node, routeConfig RouteConfig) {
	if len(routeConfig.Peers) == 0 {
		return
	}

	addRouteErr := node.AddRoute((domain.RouteID)(routeConfig.RouteID), routeConfig.Peers[0])
	if addRouteErr != nil {
		return
		panic(addRouteErr)
	}
	for peer := 1; peer < len(routeConfig.Peers); peer++ {
		extendErr := node.ExtendRoute((domain.RouteID)(routeConfig.RouteID), routeConfig.Peers[peer], 5*time.Second)
		if extendErr != nil {
			panic(extendErr)
		}
	}
}

// Add transport to Node
func AddPeersToNode(node *mesh.Node, config TestConfig) {

	emptyPK := cipher.PubKey{}

	// Connect
	for info, peerToPeer := range config.PeerToPeers {
		if peerToPeer.Peer == emptyPK {
			continue
		}
		addr, port := infoToAddr(info)
		udpConfig := physical.CreateUdp(port, addr)
		udpTransport := physical.CreateNewUDPTransport(udpConfig)
		connectError := udpTransport.ConnectToPeer(peerToPeer.Peer, peerToPeer.Info)
		if connectError != nil {
			panic(connectError)
		}
		transportToPeer := transport.NewTransport(udpTransport, config.TransportConfig)
		node.AddTransport(transportToPeer)
	}
}

func infoToAddr(info string) (string, int) {
	infoBytes, err := hex.DecodeString(info)
	if err != nil {
		panic(err)
		return "", 0
	}

	udp := physical.UDPCommConfig{}

	err = json.Unmarshal(infoBytes, &udp)
	if err != nil {
		panic(err)
		return "", 0
	}

	host := udp.ExternalHost
	return host.IP.String(), host.Port
}

// Returns Node by index
func (s *NodeManager) GetNodeByIndex(indexNode int) *mesh.Node {
	nodePubKey := s.PubKeyList[indexNode]
	return s.NodesList[nodePubKey]
}

// Get all transports from one node
func (s *NodeManager) GetTransportsFromNode(indexNode int) []transport.ITransport {
	nodePubKey := s.PubKeyList[indexNode]
	node := s.NodesList[nodePubKey]
	return node.GetTransports()
}

func (s *NodeManager) RemoveTransportsFromNode(indexNode int, transport transport.ITransport) {
	nodePubKey := s.PubKeyList[indexNode]
	node := s.NodesList[nodePubKey]
	node.RemoveTransport(transport)
}

// Connect node to netwotk
func (s *NodeManager) ConnectNodeToNetwork() (int, int) {
	// Create new node
	index1 := s.AddNode()
	index2 := s.ConnectNodeRandomly(index1)
	return index1, index2
}

// Connect Node Randomly
func (s *NodeManager) ConnectNodeRandomly(index1 int) int {
	var index2, rang int
	rang = len(s.ConfigList)
	for i := 0; i < 3; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		index2 = rand.Intn(rang)
		if index2 == index1 && i == 2 {
			fmt.Fprintf(os.Stderr, "Error Node %v not connected\n", index1)
			index2 = -1
			break
		} else if index2 != index1 {
			fmt.Fprintf(os.Stdout, "Connect node %v to node %v and vice versa.\n", index1, index2)
			pubKey1 := s.PubKeyList[index1]
			config1 := s.ConfigList[pubKey1]
			pubKey2 := s.PubKeyList[index2]
			config2 := s.ConfigList[pubKey2]
			ConnectNodeToNode(config1, config2)
			break
		}
	}
	return index2
}

// Create routes from a node
/*
func (s *NodeManager) BuildRoutes() {
	s.Routes = make(map[RouteKey]Route)
	for _, pubKey := range s.PubKeyList {
		s.FindRoutes(pubKey)
	}
}
*/

func (s *NodeManager) CloseAll() {
	for _, node := range s.NodesList {
		node.Close()
	}
}

func makePeer(pubKey cipher.PubKey, addr string, port int) *Peer {
	portStr := strconv.Itoa(port)
	address := addr + ":" + portStr

	peer := Peer{}
	peer.Peer = pubKey
	peer.Info = physical.CreateUDPCommConfig(address, nil)
	return &peer
}

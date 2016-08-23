package nodemanager

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	mesh "github.com/skycoin/skycoin/src/mesh2"
	"github.com/skycoin/skycoin/src/mesh2/transport/reliable"
	"github.com/skycoin/skycoin/src/mesh2/transport/udp"
)

type NodeManager struct {
	ConfigList []*mesh.TestConfig
	Port       int
	NodesList  []*mesh.Node
}

// Create TestConfig to the test using the functions created in the meshnet library.
func CreateTestConfig(port int) *mesh.TestConfig {
	testConfig := &mesh.TestConfig{}
	testConfig.Node = NewNodeConfig()
	testConfig.Reliable = reliable.CreateReliable(testConfig.Node.PubKey)
	testConfig.Udp = udp.CreateUdp(port, "127.0.0.1")

	return testConfig
}

func CreateNode(config mesh.TestConfig) *mesh.Node {
	node, createNodeError := mesh.NewNode(config.Node)
	if createNodeError != nil {
		panic(createNodeError)
	}

	return node
}

// Create public key
func createPubKey() cipher.PubKey {
	b := cipher.RandByte(33)
	return cipher.NewPubKey(b)
}

// Create ChaCha20Key
func CreateChaCha20Key() cipher.SecKey {
	b := cipher.RandByte(32)
	return cipher.NewSecKey(b)
}

// Create new node config
func NewNodeConfig() mesh.NodeConfig {
	nodeConfig := mesh.NodeConfig{}
	nodeConfig.PubKey = createPubKey()
	nodeConfig.ChaCha20Key = CreateChaCha20Key()
	nodeConfig.MaximumForwardingDuration = 1 * time.Minute
	nodeConfig.RefreshRouteDuration = 5 * time.Minute
	nodeConfig.ExpireMessagesInterval = 5 * time.Minute
	nodeConfig.ExpireRoutesInterval = 5 * time.Minute
	nodeConfig.TimeToAssembleMessage = 5 * time.Minute
	nodeConfig.TransportMessageChannelLength = 100

	return nodeConfig
}

// Create node list
func (self *NodeManager) CreateNodeConfigList(n int) {
	self.ConfigList = []*mesh.TestConfig{}
	self.NodesList = []*mesh.Node{}
	self.Port = 10000
	for a := 1; a <= n; a++ {
		self.AddNode()
	}
}

// Add Node to Nodes List
func (self *NodeManager) AddNode() int {
	if len(self.ConfigList) == 0 {
		self.ConfigList = []*mesh.TestConfig{}
		self.NodesList = []*mesh.Node{}
	}
	config := CreateTestConfig(self.Port)
	self.ConfigList = append(self.ConfigList, config)
	self.Port++
	node := CreateNode(*config)
	self.NodesList = append(self.NodesList, node)
	index := len(self.NodesList) - 1
	return index
}

// Connect the node list
func (self *NodeManager) ConnectNodes() {

	var index2, index3 int
	var lenght int = len(self.ConfigList)

	if lenght > 1 {
		for index1, config1 := range self.ConfigList {

			if index1 == 0 {
				index2 = 1
			} else {
				if index1 == lenght-1 {
					index2 = index1 - 1
					index3 = 0
				} else {
					index2 = index1 - 1
					index3 = index1 + 1
				}
			}
			config2 := self.ConfigList[index2]
			ConnectNodeToNode(config1, config2)

			if index3 > 0 {
				config3 := self.ConfigList[index3]
				ConnectNodeToNode(config1, config3)
			}
			self.NodesList[index1].AddTransportToNode(*config1)
		}
	}
}

// Connect Node1 (config1) to Node2 (config2)
func ConnectNodeToNode(config1, config2 *mesh.TestConfig) {
	var addr bytes.Buffer
	addr.WriteString(config2.Udp.ExternalAddress)
	addr.WriteString(":")
	addr.WriteString(strconv.Itoa(int(config2.Udp.ListenPortMin)))
	config1.AddPeerToConnect(addr.String(), config2)
	addr.Reset()
}

// Obtain port for to use in the creating from node
func (self *NodeManager) GetPort() int {
	port := self.Port
	self.Port++
	return port
}

// Connect node to netwotk
func (self *NodeManager) ConnectNodeToNetwork() (int, int) {
	// Create new node
	index1 := self.AddNode()
	index2 := self.ConnectNodeRandomly(index1)
	return index1, index2
}

// Connect Node Randomly
func (self *NodeManager) ConnectNodeRandomly(index1 int) int {
	var index2, rang int
	rang = len(self.ConfigList)
	for i := 0; i < 3; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		index2 = rand.Intn(rang)
		if index2 == index1 && i == 2 {
			fmt.Fprintf(os.Stderr, "Error Node %v not connected\n", index1)
			index2 = -1
			break
		} else if index2 != index1 {
			fmt.Fprintf(os.Stdout, "Connect node %v to node %v and vice versa.\n", index1, index2)
			config1 := self.ConfigList[index1]
			config2 := self.ConfigList[index2]
			ConnectNodeToNode(config1, config2)
			ConnectNodeToNode(config2, config1)
			break
		}
	}
	return index2
}

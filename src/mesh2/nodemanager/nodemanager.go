package nodemanager

import (
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	mesh "github.com/skycoin/skycoin/src/mesh2"
	"github.com/skycoin/skycoin/src/mesh2/transport/reliable"
	"github.com/skycoin/skycoin/src/mesh2/transport/udp"
)

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
